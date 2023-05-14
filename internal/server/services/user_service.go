package services

import (
	"context"

	"ydx-goadv-gophkeeper/internal/model"
	"ydx-goadv-gophkeeper/internal/server/repositories"
)

type UserService interface {
	Register(ctx context.Context, info model.User) (int32, error)
	Login(ctx context.Context, info model.User) (int32, error)
}

type userService struct {
	repo repositories.UserRepository
	ctx  context.Context
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, info model.User) (int32, error) {
	return s.repo.Register(ctx, info)
}

func (s *userService) Login(ctx context.Context, info model.User) (int32, error) {
	return s.repo.Login(ctx, info)
}
