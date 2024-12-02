package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/plasmatrip/metriq/internal/agent/config"
	"github.com/plasmatrip/metriq/internal/agent/controller"
	"github.com/plasmatrip/metriq/internal/storage/mem"
)

func main() {
	ctx := context.Background()

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	controller := controller.NewController(mem.NewStorage(), *config)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			controller.UpdateMetrics(ctx)
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {

			if err := controller.SendMetricsBatch(); err != nil {
				fmt.Println("Error: ", err)
			}

			time.Sleep(time.Duration(config.ReportInterval) * time.Second)
		}
	}()

	fmt.Printf(`The metrics collection agent is running.
The interval for collecting metrics is %d seconds, the interval for transmitting metrics to the server is %d seconds.
Server address: %s
`, config.PollInterval, config.ReportInterval, config.Host)

	wg.Wait()

	os.Exit(0)
}
