package main

import (
	"net/http"

	"github.com/plasmatrip/metriq/internal/server"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {
	handlers := server.NewHandlers(storage.NewStorage())

	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, handlers.UpdateHandler)

	err := http.ListenAndServe(server.Address+":"+server.Port, mux)
	if err != nil {
		panic(err)
	}
}
