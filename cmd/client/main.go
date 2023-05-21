package main

import (
	"context"
	"io"
	"os"

	clients "ydx-goadv-gophkeeper/internal/client"
	"ydx-goadv-gophkeeper/internal/client/configs"
	"ydx-goadv-gophkeeper/internal/client/model"
	"ydx-goadv-gophkeeper/internal/client/services"
	"ydx-goadv-gophkeeper/internal/client/terminal"
	"ydx-goadv-gophkeeper/pkg/logger"
	"ydx-goadv-gophkeeper/pkg/pb"
	intsrv "ydx-goadv-gophkeeper/pkg/services"
	"ydx-goadv-gophkeeper/pkg/shutdown"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"

	configPathEnvVar  = "CONFIG"
	defaultConfigPath = "cmd/client/config.json"
)

func main() {
	log := logger.NewLogger("main")
	log.Infof("Server args: %s", os.Args[1:])
	ctx, ctxClose := context.WithCancel(context.Background())
	configPath := os.Getenv(configPathEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
	}
	appConfig, err := configs.InitAppConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	exitHandler := shutdown.NewExitHandlerWithCtx(ctxClose)

	tokenHolder := &model.TokenHolder{}

	grpcConn, err := clients.CreateGrpcConnection(appConfig.ServerPort, tokenHolder)
	if err != nil {
		log.Fatalf("failed to create grpc connection: %v", err)
	}
	exitHandler.ToClose = []io.Closer{grpcConn}

	authService := services.NewAuthService(pb.NewAuthClient(grpcConn), tokenHolder)
	fileService := intsrv.NewFileService()
	cryptoService := services.NewCryptService(appConfig.PrivateKey)
	resourceService := services.NewResourceService(pb.NewResourcesClient(grpcConn), fileService, cryptoService)
	exit := shutdown.ProperExitDefer(exitHandler)

	commandProcessor := terminal.NewCommandParser(buildVersion, buildDate, authService, resourceService, exitHandler)
	commandProcessor.Start(exit)
	<-ctx.Done()
}
