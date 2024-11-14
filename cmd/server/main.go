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

	h := handlers.NewHandlers(storage.NewStorage())

	r := chi.NewRouter()

	r.Post("/update/{metricType}/{metricName}/{metricValue}", l.WithLogging(h.UpdateHandler))
	r.Post("/update/", l.WithLogging(h.JSONUpdateHandler))
	r.Post("/update", l.WithLogging(h.JSONUpdateHandler))
	r.Post("/value/", l.WithLogging(h.JSONValueHandler))
	r.Post("/value", l.WithLogging(h.JSONValueHandler))
	r.Get("/value/{metricType}/{metricName}", l.WithLogging(h.ValueHandler))
	r.Get("/", l.WithLogging(h.MetricsHandler))

	err = http.ListenAndServe(config.Host, func(next http.Handler) http.Handler {
		l.Sugar.Infow("The metrics collection server is running. ", "Server address: ", config.Host)
		return next
	}(r))
	if err != nil {
		panic(err)
	}
}
