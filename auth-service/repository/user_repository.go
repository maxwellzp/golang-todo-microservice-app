package repository

import (
	"auth-service/model"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(user *model.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	row := r.db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
