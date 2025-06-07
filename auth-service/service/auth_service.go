package service

import (
	"auth-service/model"
	"auth-service/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *repository.UserRepository
	jwtService *JWTService
}

func NewAuthService(repo *repository.UserRepository, jwt *JWTService) *AuthService {
	return &AuthService{repo, jwt}
}

func (s *AuthService) Register(email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashed),
	}

	return s.repo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
