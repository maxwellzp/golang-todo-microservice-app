package service

import (
	"todo-service/internal/model"
	"todo-service/internal/repository"
)

type TodoService interface {
	CreateTodo(todo *model.Todo) error
	GetUserTodos(userID string) ([]model.Todo, error)
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(todo *model.Todo) error {
	return s.repo.Create(todo)
}

func (s *todoService) GetUserTodos(userID string) ([]model.Todo, error) {
	return s.repo.FindByUserID(userID)
}
