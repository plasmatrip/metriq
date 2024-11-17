package main

import (
	"net/http"

	"github.com/plasmatrip/metriq/internal/backup"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/server/routing"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	storage := storage.NewStorage()

	backup, err := backup.NewBackup(config.FileStoragePath, storage, logger)
	if err != nil {
		logger.Sugar.Panic("Error initializing backup: ", err, " ", config.FileStoragePath)
	}
	backup.Start(config.StoreInterval, config.Restore)

	err = http.ListenAndServe(config.Host, func(next http.Handler) http.Handler {
		logger.Sugar.Infow("The metrics collection server is running. ", "Server address: ", config.Host)
		return next
	}(routing.NewRouter(logger, storage, *config)))
	if err != nil {
		panic(err)
	}
}
