package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server"
	"github.com/plasmatrip/metriq/internal/server/handlers"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {
	config := server.NewConfig()

	r := chi.NewRouter()

	l, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer l.Close()

	h := handlers.NewHandlers(storage.NewStorage())

	r.Post("/update/*", l.WithLogging(h.UpdateHandler))
	r.Get("/value/*", l.WithLogging(h.ValueHandler))
	r.Get("/", l.WithLogging(h.MetricsHandler))

	err = http.ListenAndServe(config.Host, func(next http.Handler) http.Handler {
		l.Sugar.Infow("The metrics collection server is running. ", "Server address: ", config.Host)
		return next
	}(r))
	if err != nil {
		panic(err)
	}
}
