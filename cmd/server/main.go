package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plasmatrip/metriq/internal/server/handlers"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {
	// config := server.NewConfig()

	var host string
	flag.StringVar(&host, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	fmt.Println(host)

	r := chi.NewRouter()

	handlers := handlers.NewHandlers(storage.NewStorage())

	r.Post("/update/*", handlers.UpdateHandler)
	r.Get("/value/*", handlers.ValueHandler)
	r.Get("/", handlers.MetricsHandler)

	err := http.ListenAndServe(host, r)
	// err := http.ListenAndServe(config.Host+":"+config.Port, r)
	if err != nil {
		panic(err)
	}
}
