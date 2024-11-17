package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/plasmatrip/metriq/internal/backup"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/server/routing"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, os.Interrupt)
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

	s := storage.NewStorage()

	backup, err := backup.NewBackup(*c, s, l)
	if err != nil {
		l.Sugar.Panic("error initializing backup: ", err, " ", c.FileStoragePath)
	}
	backup.Start()

	server := http.Server{
		Addr: c.Host,
		Handler: func(next http.Handler) http.Handler {
			l.Sugar.Infow("The metrics collection server is running. ", "Server address: ", c.Host)
			return next
		}(routing.NewRouter(s, *c, l)),
	}

	go server.ListenAndServe()

	<-ctx.Done()

	// err = http.ListenAndServe(c.Host, func(next http.Handler) http.Handler {
	// 	l.Sugar.Infow("The metrics collection server is running. ", "Server address: ", c.Host)
	// 	return next
	// }(routing.NewRouter(s, l, *c)))
	// if err != nil {
	// 	panic(err)
	// }

	l.Sugar.Infow("save metrics on exit")
	err = backup.Save()
	if err != nil {
		l.Sugar.Infow("error saving to backup: ", err, " ", c.FileStoragePath)
	}

	os.Exit(0)
}
