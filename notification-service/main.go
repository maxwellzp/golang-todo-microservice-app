package main

import (
	"log"
	"notification-service/internal/handler"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with system environment variables.")
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL is not set")
	}

	e := echo.New()

	// Initialize NATS handler
	natsHandler, err := handler.NewNATSHandler(natsURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}

	// Run NATS subscriptions
	natsHandler.Subscribe("todo.created", func(msg string) {
		log.Println("🔔 Received todo.created event:", msg)
	})

	log.Println("✅ Notification service is running on :8083")
	e.Logger.Fatal(e.Start(":8083"))
}
