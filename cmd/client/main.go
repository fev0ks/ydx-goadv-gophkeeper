package main

import (
	"io"
	"os"

	pb "ydx-goadv-gophkeeper/api/proto"
	clients "ydx-goadv-gophkeeper/internal/client"
	"ydx-goadv-gophkeeper/internal/client/configs"
	"ydx-goadv-gophkeeper/internal/client/model"
	"ydx-goadv-gophkeeper/internal/client/services"
	"ydx-goadv-gophkeeper/internal/client/terminal"
	"ydx-goadv-gophkeeper/internal/logger"
	intsrv "ydx-goadv-gophkeeper/internal/services"
	"ydx-goadv-gophkeeper/internal/shutdown"
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
	configPath := os.Getenv(configPathEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
	}
	appConfig, err := configs.InitAppConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	exitHandler := shutdown.NewExitHandler()

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
	shutdown.ProperExitDefer(exitHandler)

	commandProcessor := terminal.NewCommandParser(buildVersion, buildDate, authService, resourceService)
	commandProcessor.Start()
}
