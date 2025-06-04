package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"todo-service/models"
)

type TodoHandler struct {
	DB *gorm.DB
}

type CreateTodoRequest struct {
	Title       string    `json:"title" validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"max=1000"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTodoRequest struct {
	Title       string    `json:"title" validate:"omitempty,min=3,max=255"`
	Description string    `json:"description" validate:"omitempty,max=1000"`
	Completed   bool      `json:"completed"`
	DueDate     time.Time `json:"due_date"`
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{DB: db}
}

// Extract user ID from JWT token
func (h *TodoHandler) getUserIDFromToken(c echo.Context) (uint, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))
	return userID, nil
}

// CreateTodo creates a new todo item
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	req := new(CreateTodoRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	todo := models.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}

	if err := h.DB.Create(&todo).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create todo")
	}

	return c.JSON(http.StatusCreated, todo)
}

// GetTodos returns all todos for the authenticated user
func (h *TodoHandler) GetTodos(c echo.Context) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	var todos []models.Todo
	if err := h.DB.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch todos")
	}

	return c.JSON(http.StatusOK, todos)
}

// GetTodo returns a single todo item
func (h *TodoHandler) GetTodo(c echo.Context) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
	}

	var todo models.Todo
	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	}

	return c.JSON(http.StatusOK, todo)
}

// UpdateTodo updates a todo item
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
	}

	req := new(UpdateTodoRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var todo models.Todo
	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	}

	// Update fields if they are provided in the request
	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	todo.Completed = req.Completed
	if !req.DueDate.IsZero() {
		todo.DueDate = req.DueDate
	}

	if err := h.DB.Save(&todo).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update todo")
	}

	return c.JSON(http.StatusOK, todo)
}

// DeleteTodo deletes a todo item
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
	}

	result := h.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Todo{})
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete todo")
	}
	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	}

	return c.NoContent(http.StatusNoContent)
}
