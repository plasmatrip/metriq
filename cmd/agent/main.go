package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/plasmatrip/metriq/internal/agent/config"
	"github.com/plasmatrip/metriq/internal/agent/controller"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	controller := controller.NewController(storage.NewStorage(), *config)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			controller.UpdateMetrics()
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := controller.SendMetrics(); err != nil {
				fmt.Println("Error: ", err)
			}
			time.Sleep(time.Duration(config.ReportInterval) * time.Second)
		}
	}()

	fmt.Printf(`The metrics collection agent is running.
The interval for collecting metrics is %d seconds, the interval for transmitting metrics to the server is %d seconds.
Server address: %s
`, config.PollInterval, config.ReportInterval, config.Host)

	<-ctx.Done()

	wg.Wait()

	os.Exit(0)
}
