package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/plasmatrip/metriq/internal/agent"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {
	config, err := agent.NewConfig()
	if err != nil {
		panic(err)
	}
	controller := agent.NewController(storage.NewStorage())

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
			if err := controller.SendMetrics("http://" + config.Host); err != nil {
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
}
