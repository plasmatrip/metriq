package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/plasmatrip/metriq/internal/agent"
	"github.com/plasmatrip/metriq/internal/server"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {

	controller := agent.NewController(storage.NewStorage())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			controller.UpdateMetrics()
			time.Sleep(agent.ReadTimeout * time.Second)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := controller.SendMetrics(server.URL); err != nil {
				fmt.Print(err)
			}
			time.Sleep(agent.SendTimeout * time.Second)
		}
	}()
	wg.Wait()
}
