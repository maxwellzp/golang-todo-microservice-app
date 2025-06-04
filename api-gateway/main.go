package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServiceConfig struct {
	Prefix string
	Target string
	Auth   bool
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configuration for services
	services := []ServiceConfig{
		{
			Prefix: "/auth",
			Target: "http://auth-service:8081",
			Auth:   false, // Auth service doesn't need JWT validation
		},
		{
			Prefix: "/todos",
			Target: "http://todo-service:8082",
			Auth:   true,
		},
		{
			Prefix: "/notify",
			Target: "http://notification-service:8083",
			Auth:   true,
		},
	}

	// JWT secret (must match auth service)
	jwtSecret := "your-strong-secret-key" // In production, use environment variable

	// Create reverse proxies for each service
	for _, service := range services {
		target, _ := url.Parse(service.Target)
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Configure the proxy director
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host
			req.URL.Host = target.Host
			req.URL.Scheme = target.Scheme
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		}

		// Create middleware chain
		middlewares := []echo.MiddlewareFunc{
			middleware.RemoveTrailingSlash(),
			middleware.BodyLimit("2M"),
		}

		// Add JWT middleware if needed
		if service.Auth {
			middlewares = append(middlewares, jwtMiddleware(jwtSecret))
		}

		// Register the route with middlewares
		e.Group(service.Prefix, middlewares...).Any("/*", echo.WrapHandler(proxy))
	}

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// jwtMiddleware validates JWT tokens for protected routes
func jwtMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token method")
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// Store token in context for downstream services
			c.Set("user", token)

			// Read and rewind the request body so it can be read again by the proxy
			if c.Request().Body != nil {
				body, err := io.ReadAll(c.Request().Body)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read request body")
				}
				c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
			}

			return next(c)
		}
	}
}
