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
	config, err := server.NewConfig()
	if err != nil {
		panic(err)
	}

	l, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer l.Close()

	h := handlers.NewHandlers(storage.NewStorage(), *config)

	r := chi.NewRouter()

	r.Use(server.WithCompressed)
	r.Use(l.WithLogging)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.JSONUpdateHandler)
	})
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.JSONValueHandler)
	})
	r.Get("/value/{metricType}/{metricName}", h.ValueHandler)
	r.Get("/", h.MetricsHandler)

	err = http.ListenAndServe(config.Host, func(next http.Handler) http.Handler {
		l.Sugar.Infow("The metrics collection server is running. ", "Server address: ", config.Host)
		return next
	}(r))
	if err != nil {
		panic(err)
	}
}
