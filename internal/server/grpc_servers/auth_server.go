package grpc_servers

import (
	"context"
	"errors"
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

	User := model.User{Username: authData.Username, Password: authData.Password}
	id, err := s.userService.Register(ctx, User)

	if err == nil {
		return s.genToken(id)
	}
	if errors.Is(err, errs.ErrUserAlreadyExist) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	return nil, status.Error(codes.Internal, "internal error")
}

func (s *authServer) Login(ctx context.Context, authData *pb.AuthData) (*pb.TokenData, error) {
	if err := validateAuthData(authData); err != nil {
		return nil, err
	}
	User := model.User{Username: authData.Username, Password: authData.Password}
	id, err := s.userService.Login(ctx, User)
	if err == nil {
		return s.genToken(id)
	}
	if errors.Is(err, errs.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return nil, status.Error(codes.Internal, "internal error")
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
