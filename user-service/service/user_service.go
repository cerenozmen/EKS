package service

import (
	"errors"
	"user-service/model"
	"user-service/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(username, password, name string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(username, password, name string) (*model.User, error) {
	
	existingUser, _ := s.repo.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
