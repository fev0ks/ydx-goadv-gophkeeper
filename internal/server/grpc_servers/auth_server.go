package grpc_servers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"ydx-goadv-gophkeeper/internal/server/services"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/model"
	"ydx-goadv-gophkeeper/internal/model/errs"
	pb "ydx-goadv-gophkeeper/internal/proto"
)

type authServer struct {
	log *zap.SugaredLogger
	pb.UnimplementedAuthServer
	userService  services.UserService
	tokenService services.TokenService
}

func NewAuthServer(userService services.UserService, tokenService services.TokenService) pb.AuthServer {
	return &authServer{
		log:          logger.NewLogger("auth-server"),
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *authServer) Register(ctx context.Context, authData *pb.AuthData) (*pb.TokenData, error) {
	if err := validateAuthData(authData); err != nil {
		return nil, err
	}

	user := &model.User{Username: authData.Username, Password: []byte(authData.Password)}
	id, err := s.userService.CreateUser(ctx, user)
	if errors.Is(err, errs.ErrUserAlreadyExist) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		s.log.Errorf("failed to create user: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create user: %v", err))
	}
	return s.genToken(id)
}

func (s *authServer) Login(ctx context.Context, authData *pb.AuthData) (*pb.TokenData, error) {
	if err := validateAuthData(authData); err != nil {
		return nil, err
	}
	user, err := s.userService.GetUser(ctx, authData.Username)
	if errors.Is(err, errs.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		s.log.Errorf("failed to get user: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user: %v", err))
	}
	ok, err := s.userService.ValidatePassword(ctx, user, authData.Password)
	if err != nil {
		s.log.Errorf("failed to check user password: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to check user password: %v", err))
	}
	if !ok {
		s.log.Warn("password is incorrect")
		return nil, status.Error(codes.InvalidArgument, "password is incorrect")
	}
	return s.genToken(user.Id)
}

func validateAuthData(authData *pb.AuthData) error {
	if len(authData.Username) == 0 || len(authData.Password) == 0 {
		return status.Error(codes.InvalidArgument, "invalid username/password format: must be nonempty")
	}
	return nil
}

func (s *authServer) genToken(id int32) (*pb.TokenData, error) {
	expireAt := time.Now().UTC().Add(time.Hour)
	token, err := s.tokenService.Generate(id, expireAt)
	if err != nil {
		s.log.Errorf("error on register: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	s.log.Infof("generate token successfully: %v", zap.Time("expireAt", expireAt))
	return &pb.TokenData{Token: token, ExpireAt: timestamppb.New(expireAt)}, nil
}
