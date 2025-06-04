package main

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"todo-service/handlers"
	"todo-service/models"
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
	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		e.Logger.Fatal("Failed to migrate database")
	}

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(db)

	// JWT Middleware - Updated configuration
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &jwt.MapClaims{}
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	// Routes with JWT auth - Updated middleware
	r := e.Group("/todos")
	r.Use(echojwt.WithConfig(config)) // Changed to use echojwt
	r.POST("", todoHandler.CreateTodo)
	r.GET("", todoHandler.GetTodos)
	r.GET("/:id", todoHandler.GetTodo)
	r.PUT("/:id", todoHandler.UpdateTodo)
	r.DELETE("/:id", todoHandler.DeleteTodo)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
