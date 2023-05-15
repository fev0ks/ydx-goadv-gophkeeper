package services

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"ydx-goadv-gophkeeper/internal/model"
	"ydx-goadv-gophkeeper/internal/server/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (int32, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
	ValidatePassword(_ context.Context, user *model.User, password string) (bool, error)
}

type userService struct {
	repo repositories.UserRepository
	ctx  context.Context
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) (int32, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, 8)
	if err != nil {
		return 0, err
	}
	newUser := &model.User{
		Username: user.Username,
		Password: hashedPassword,
	}
	return s.repo.CreateUser(ctx, newUser)
}

func (s *userService) GetUser(ctx context.Context, username string) (*model.User, error) {
	return s.repo.GetUser(ctx, username)
}

func (us *userService) ValidatePassword(_ context.Context, user *model.User, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
