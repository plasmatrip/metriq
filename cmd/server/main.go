package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/plasmatrip/metriq/internal/backup"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/server/routing"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/storage/db"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	l, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer l.Close()

	var s storage.Repository
	if c.DSN == "" {
		s = storage.NewStorage()
	} else {
		s, err = db.NewPostgresStorage(c.DSN)
		if err != nil {
			l.Sugar.Infow("database connection error: ", err)
			os.Exit(1)
		}
		defer s.Close()
	}

	backup, err := backup.NewBackup(*c, s, l)
	if err != nil {
		l.Sugar.Panic("error initializing backup: ", err, " ", c.FileStoragePath)
	}
	backup.Start()

	server := http.Server{
		Addr: c.Host,
		Handler: func(next http.Handler) http.Handler {
			l.Sugar.Infow("The metrics collection server is running. ", "Server address: ", c.Host)
			l.Sugar.Infow("Server config: ", "store interval: ", c.StoreInterval, "backup file: ", c.FileStoragePath, "DSN: ", c.DSN)
			return next
		}(routing.NewRouter(s, *c, l)),
	}

	go server.ListenAndServe()

	<-ctx.Done()

	err = backup.Save()
	if err != nil {
		l.Sugar.Infow("error saving to backup: ", err, " ", c.FileStoragePath)
	}

	server.Shutdown(context.Background())

	os.Exit(0)
}
