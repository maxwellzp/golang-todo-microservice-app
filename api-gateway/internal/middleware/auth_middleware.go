package middleware

import (
	"net/http"
	"strings"

	"api-gateway/internal/utils"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid token"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		return next(c)
	}
}
