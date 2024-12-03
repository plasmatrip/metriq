package routing

import (
	"github.com/go-chi/chi/v5"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/compress"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/server/handlers"
	"github.com/plasmatrip/metriq/internal/storage"
)

func NewRouter(s storage.Repository, c config.Config, l *logger.Logger) *chi.Mux {
	h := handlers.NewHandlers(s, c, *l)

	r := chi.NewRouter()

	r.Use(compress.WithCompressed)
	r.Use(l.WithLogging)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.JSONUpdate)
	})
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", h.JSONUpdates)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.JSONValue)
	})
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.Update)
	r.Get("/value/{metricType}/{metricName}", h.Value)
	r.Get("/", h.Metrics)
	r.Route("/ping", func(r chi.Router) {
		r.Get("/", h.Ping)
	})

	return r
}
