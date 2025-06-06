package handlers

//import (
//	"github.com/labstack/echo/v4"
//	"io"
//	"log"
//	"net/http"
//	"os"
//)

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

//
//func ProxyHandler(serviceEnv string) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		serviceURL := os.Getenv(serviceEnv)
//		target, _ := url.JoinPath(serviceURL, c.Request().URL.Path)
//
//		req, err := http.NewRequest(c.Request().Method, target, c.Request().Body)
//		if err != nil {
//			return err
//		}
//
//		req.Header = c.Request().Header.Clone()
//
//		client := &http.Client{}
//		resp, err := client.Do(req)
//		if err != nil {
//			return err
//		}
//		defer resp.Body.Close()
//
//		body, _ := io.ReadAll(resp.Body)
//
//		return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), body)
//	}
//}

func ProxyHandler(serviceEnv string) echo.HandlerFunc {
	target := os.Getenv(serviceEnv)
	if target == "" {
		log.Fatalf("Missing env: %s", serviceEnv)
	}
	return func(c echo.Context) error {
		url := target + c.Request().URL.Path
		log.Printf("Proxying %s request to: %s", c.Request().Method, url)

		req, err := http.NewRequest(c.Request().Method, url, c.Request().Body)
		if err != nil {
			log.Println("Error creating request:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "proxy error")
		}
		req.Header = c.Request().Header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error forwarding request:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "proxy error")
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			c.Response().Header()[k] = v
		}
		c.Response().WriteHeader(resp.StatusCode)
		_, err = io.Copy(c.Response().Writer, resp.Body)
		return err
	}
}
