package main

import (
	"auth-service/handlers"
	"auth-service/models"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Connect to PostgreSQL
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal("Failed to connect to database")
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.User{}); err != nil {
		e.Logger.Fatal("Failed to migrate database")
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)

	// Routes
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
	e.GET("/validate", authHandler.ValidateToken)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}
