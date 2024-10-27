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
	controller := agent.NewController(*agent.NewConfig(), storage.NewStorage())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			controller.UpdateMetrics()
			time.Sleep(time.Duration(controller.Config.PollInterval) * time.Second)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := controller.SendMetrics(server.URL); err != nil {
				fmt.Print(err)
			}
			time.Sleep(time.Duration(controller.Config.ReportInterval) * time.Second)
		}
	}()
	wg.Wait()
}
