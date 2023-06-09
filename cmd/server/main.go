package main

import (
	"context"
	"os"

	"ydx-goadv-gophkeeper/internal/server/configs"
	servers "ydx-goadv-gophkeeper/internal/server/grpc_servers"
	"ydx-goadv-gophkeeper/internal/server/repositories"
	"ydx-goadv-gophkeeper/internal/server/services"
	"ydx-goadv-gophkeeper/pkg/logger"
	intsrv "ydx-goadv-gophkeeper/pkg/services"
	"ydx-goadv-gophkeeper/pkg/shutdown"
)

var (
	BuildVersion      = "N/A"
	BuildDate         = "N/A"
	BuildCommit       = "N/A"
	configPathEnvVar  = "CONFIG"
	defaultConfigPath = "cmd/server/config.json"
)

func main() {
	log := logger.NewLogger("main")
	log.Infof("Build version: %s", BuildVersion)
	log.Debugf("Build date: %s", BuildDate)
	log.Debugf("Build commit: %s", BuildCommit)

	ctx, ctxCancel := context.WithCancel(context.Background())

	log.Infof("Server args: %s", os.Args[1:])
	configPath := os.Getenv(configPathEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
	}
	appConfig, err := configs.InitAppConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	dbProvider, err := repositories.NewPgProvider(ctx, appConfig)
	if err != nil {
		log.Fatalln(err)
	}
	exitHandler := shutdown.NewExitHandlerWithCtx(ctxCancel)

	userRepo := repositories.NewUserRepository(dbProvider)
	resRepo := repositories.NewResourceRepository(dbProvider)

	userSrv := services.NewUserService(userRepo)
	resSrv := services.NewResourceService(resRepo)
	tokenSrv := services.NewTokenService(appConfig.TokenKey)
	fileProcessor := intsrv.NewFileService()

	authServer := servers.NewAuthServer(userSrv, tokenSrv)
	resourcesServer := servers.NewResourcesServer(resSrv, fileProcessor, exitHandler)

	serverManager, err := servers.NewServerManager(tokenSrv)
	if err != nil {
		log.Fatalf("failed to init grpc server: %v", err)
	}
	serverManager.RegisterResourcesServer(resourcesServer)
	serverManager.RegisterAuthServer(authServer)
	server, err := serverManager.Start(appConfig.ServerPort)
	if err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
	exitHandler.ShutdownGrpcServerBeforeExit(server)
	exit := exitHandler.ProperExitDefer()
	<-exit
	log.Info("Program is going to be closed")
	<-ctx.Done()
}
