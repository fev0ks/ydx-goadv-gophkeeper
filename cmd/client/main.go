package main

import (
	"io"

	clients "ydx-goadv-gophkeeper/internal/client"
	"ydx-goadv-gophkeeper/internal/client/model"
	"ydx-goadv-gophkeeper/internal/client/services"
	"ydx-goadv-gophkeeper/internal/client/terminal"
	"ydx-goadv-gophkeeper/internal/logger"
	pb "ydx-goadv-gophkeeper/internal/proto"
	intsrv "ydx-goadv-gophkeeper/internal/services"
	"ydx-goadv-gophkeeper/internal/shutdown"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	log := logger.NewLogger("main")
	exitHandler := shutdown.NewExitHandler()

	tokenHolder := &model.TokenHolder{}

	grpcConn, err := clients.CreateGrpcConnection(":3200", tokenHolder)
	if err != nil {
		log.Fatalf("failed to create grpc connection: %v", err)
	}
	exitHandler.ToClose = []io.Closer{grpcConn}

	authService := services.NewAuthService(pb.NewAuthClient(grpcConn), tokenHolder)
	fileService := intsrv.NewFileService()
	resourceService := services.NewResourceService(pb.NewResourcesClient(grpcConn), fileService)
	shutdown.ProperExitDefer(exitHandler)

	commandProcessor := terminal.NewCommandParser(buildVersion, buildDate, authService, resourceService)
	commandProcessor.Start()
}
