package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"todo-service/internal/model"
	"todo-service/internal/nats"
	"todo-service/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	service   service.TodoService
	publisher *nats.Publisher
}

func NewTodoHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/todo/create", h.CreateTodo)
	e.GET("/todos", h.GetTodos)
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	userID, err := h.extractUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	todo := new(model.Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	todo.UserID = userID

	if err := h.service.CreateTodo(todo); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// ✅ Publish event to NATS
	if h.publisher != nil {
		message := fmt.Sprintf(`{"event":"todo.created","user_id":"%s","title":"%s","id":"%d"}`, todo.UserID, todo.Title, todo.ID)
		h.publisher.Publish("todo.created", message)
	}

	return c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodos(c echo.Context) error {
	userID, err := h.extractUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	todos, err := h.service.GetUserTodos(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) extractUserID(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return "", echo.ErrUnauthorized
	}

	return claims["user_id"].(string), nil
}
