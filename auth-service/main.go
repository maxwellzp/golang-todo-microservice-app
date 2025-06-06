package main

import (
	"auth-service/handler"
	"auth-service/repository"
	"auth-service/service"
	"auth-service/util"
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", util.GetPostgresDSN())
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)
	jwtService := service.NewJWTService(os.Getenv("JWT_SECRET"))
	authService := service.NewAuthService(repo, jwtService)
	authHandler := handler.NewAuthHandler(authService)

	e := echo.New()
	// Group auth routes under /auth
	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Println("Auth service running on port", port)
	log.Fatal(e.Start(":" + port))
}
