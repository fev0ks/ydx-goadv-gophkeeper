package services

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/internal/logger"
	"ydx-goadv-gophkeeper/internal/server/repositories"
)

type HealthChecker struct {
	log *zap.SugaredLogger
	ctx context.Context
	db  repositories.DBProvider
}

func NewHealthChecker(ctx context.Context, db repositories.DBProvider) HealthChecker {
	return HealthChecker{logger.NewLogger("health-checker"), ctx, db}
}

// CheckDBHandler - проверка состояния соединения с базой данных
func (hc *HealthChecker) CheckDBHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := hc.db.HealthCheck(hc.ctx)
		if err != nil {
			log.Printf("failed db health check: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
