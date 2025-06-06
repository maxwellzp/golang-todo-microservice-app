package repository

import (
	"database/sql"
	"todo-service/internal/model"
)

type TodoRepository interface {
	Create(todo *model.Todo) error
	FindByUserID(userID string) ([]model.Todo, error)
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *model.Todo) error {
	query := `
		INSERT INTO todos (user_id, title, completed)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, todo.UserID, todo.Title, todo.Completed).
		Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)
}

func (r *todoRepository) FindByUserID(userID string) ([]model.Todo, error) {
	rows, err := r.db.Query("SELECT id, user_id, title, completed, created_at, updated_at FROM todos WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}
