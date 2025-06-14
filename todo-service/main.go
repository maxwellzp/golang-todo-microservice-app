package main

import (
	"database/sql"
	"log"
	"os"
	"todo-service/internal/handler"
	"todo-service/internal/nats"
	"todo-service/internal/repository"
	"todo-service/internal/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// ✅ NATS Publisher
	publisher, err := nats.NewPublisher()
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer publisher.Close()

	// ✅ Inject all dependencies
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService, publisher)

	e := echo.New()
	todoHandler.RegisterRoutes(e)

	log.Println("✅ TODO service started on :8082")
	e.Logger.Fatal(e.Start(":8082"))
}
