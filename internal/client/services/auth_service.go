package services

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "ydx-goadv-gophkeeper/internal/api/proto"

	"ydx-goadv-gophkeeper/internal/client/model"
	"ydx-goadv-gophkeeper/internal/logger"
)

type AuthService interface {
	Register(ctx context.Context, username string, password string) (*pb.TokenData, error)
	Login(ctx context.Context, username string, password string) (*pb.TokenData, error)
}

type authService struct {
	log              *zap.SugaredLogger
	authClient       pb.AuthClient
	refreshTokenOnce sync.Once
	tokenHolder      *model.TokenHolder
}

func NewAuthService(
	client pb.AuthClient,
	tokenHolder *model.TokenHolder,
) AuthService {
	return &authService{
		log:         logger.NewLogger("auth-service"),
		authClient:  client,
		tokenHolder: tokenHolder,
	}
}

func (s *authService) Register(ctx context.Context, username string, password string) (*pb.TokenData, error) {
	tokenData, err := s.authClient.Register(ctx, &pb.AuthData{
		Username: username,
		Password: password,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok && e.Code() == codes.AlreadyExists {
			return nil, errors.New(e.Message())
		}
		s.log.Errorf("failed to register: %v", err)
		return nil, err
	}
	s.tokenHolder.Set(tokenData.Token)

	return tokenData, nil
}

func (s *authService) Login(ctx context.Context, username string, password string) (*pb.TokenData, error) {
	tokenData, err := s.authClient.Login(
		ctx,
		&pb.AuthData{
			Username: username,
			Password: password,
		},
	)

	if err != nil {
		if statusErr, ok := status.FromError(err); ok && statusErr.Code() == codes.NotFound {
			return nil, errors.New(statusErr.Message())
		}
		return nil, err
	}

	s.tokenHolder.Set(tokenData.Token)
	return tokenData, nil
}
