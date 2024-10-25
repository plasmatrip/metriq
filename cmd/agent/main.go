package main

import (
	"sync"
	"time"

	"github.com/plasmatrip/metriq/internal/agent"
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
			// controller.Repo.Print()
			time.Sleep(agent.ReadTimeout * time.Second)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			controller.SendMetrics()
			time.Sleep(agent.SendTimeout * time.Second)
		}
	}()
	wg.Wait()
}
