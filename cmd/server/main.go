package main

import (
	"net/http"

	"github.com/plasmatrip/metriq/internal/server"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, server.UpdateHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
