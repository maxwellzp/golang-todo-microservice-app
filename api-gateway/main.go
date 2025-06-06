package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"api-gateway/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using default env.")
	}

	e := echo.New()

	routes.InitRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
