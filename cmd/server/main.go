package main

import (
	"context"
	"os"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/server/configs"
	servers "ydx-goadv-gophkeeper/internal/server/grpc_servers"
	"ydx-goadv-gophkeeper/internal/server/repositories"
	"ydx-goadv-gophkeeper/internal/server/services"
	"ydx-goadv-gophkeeper/internal/shutdown"
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

	ctx := context.Background()

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
	exitHandler := shutdown.NewExitHandler()

	userRepo := repositories.NewUserRepository(dbProvider)
	resRepo := repositories.NewResourceRepository(dbProvider)

	userSrv := services.NewUserService(userRepo)
	resSrv := services.NewResourceService(resRepo)
	tokenSrv := services.NewTokenService(appConfig.TokenKey)
	fileProcessor := services.NewFileProcessor()

	authServer := servers.NewAuthServer(userSrv, tokenSrv)
	resourcesServer := servers.NewResourcesServer(resSrv, fileProcessor)

	serverManager := servers.NewServerManager(tokenSrv)
	serverManager.RegisterResourcesServer(resourcesServer)
	serverManager.RegisterAuthServer(authServer)
	server, err := serverManager.Start(appConfig.ServerPort)
	if err != nil {
		log.Fatalf("failed to start proto server: %v", err)
	}
	exitHandler.ShutdownGrpcServerBeforeExit(server)
	shutdown.ProperExitDefer(exitHandler)
	<-ctx.Done()
}
