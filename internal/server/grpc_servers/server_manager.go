package grpc_servers

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "ydx-goadv-gophkeeper/internal/api/proto"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/server/interceptors"
	"ydx-goadv-gophkeeper/internal/server/services"
)

const (
	registerMethod = "/gophkeeper.Auth/Register"
	loginMethod    = "/gophkeeper.Auth/Login"
)

type ServerManager interface {
	RegisterAuthServer(authServer pb.AuthServer)
	RegisterResourcesServer(resServer pb.ResourcesServer)
	Start(port string) (*grpc.Server, error)
}

type serverManager struct {
	log    *zap.SugaredLogger
	server *grpc.Server
}

func NewServerManager(
	tokenService services.TokenService,
) ServerManager {
	tokenValidator := interceptors.NewRequestTokenProcessor(tokenService, registerMethod, loginMethod)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(tokenValidator.TokenInterceptor()),
		grpc.StreamInterceptor(tokenValidator.TokenStreamInterceptor()),
	)
	return &serverManager{
		log:    logger.NewLogger("server-mnr"),
		server: server,
	}
}

func (s *serverManager) Start(port string) (*grpc.Server, error) {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		s.log.Errorf("failed to listen '%s': %v", port, err)
		return nil, err
	}

	go func() {
		s.log.Infof("proto server started on %s", port)
		if err := s.server.Serve(listen); err != nil {
			s.log.Fatalf("failed to start server: %v", err)
		}
	}()
	return s.server, nil
}

func (s *serverManager) RegisterAuthServer(authServer pb.AuthServer) {
	pb.RegisterAuthServer(s.server, authServer)
}

func (s *serverManager) RegisterResourcesServer(resServer pb.ResourcesServer) {
	pb.RegisterResourcesServer(s.server, resServer)
}
