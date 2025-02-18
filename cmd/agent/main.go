// This package contains the main function for the metrics collection agent. The
// agent is responsible for collecting various metrics about the application and
// sending them to the server. The main function initializes the agent by
// creating a context to listen for termination signals and setting up goroutines
// for collecting application metrics, collecting system metrics using gopsutil,
// and sending metrics to the server. The agent uses a configuration module to
// determine the polling intervals for collecting metrics and the server address.
// The agent gracefully shuts down all goroutines and exits when it receives a
// termination signal.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/plasmatrip/metriq/internal/agent/config"
	"github.com/plasmatrip/metriq/internal/agent/controller"
	"github.com/plasmatrip/metriq/internal/storage/mem"
)

// main initializes the metrics collection agent. It sets up a context to listen
// for termination signals and configures the agent using settings from the
// configuration module. The function starts three goroutines: one for collecting
// application metrics, one for collecting system metrics using gopsutil, and
// one for sending metrics to the server. Each goroutine runs periodically based
// on configuration-defined intervals. The function waits for a stop signal to
// gracefully shut down all goroutines and exit.
func main() {
	// Create a context to listen for termination signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	controller := controller.NewController(mem.NewStorage(), *cfg)

	var wg sync.WaitGroup

	// start goroutine to collect application metrics
	// in an infinite loop, it reads from a ticker and the context's Done channel
	// when the context is canceled, the goroutine exits
	// on each ticker event, it updates the metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				controller.UpdateMetrics(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start a goroutine to collect metrics using the gopsutil package.
	// In a loop, it reads from a ticker and the context's Done channel.
	// When the context is canceled, the goroutine exits.
	// On each ticker event, it updates the metrics.
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := controller.UpdatePSMetrics(ctx)
				if err != nil {
					fmt.Println("error while collecting system utilization metrics using gopsutil: ", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// start goroutine to send metrics to the server
	// in a loop, it reads from a ticker and the context's Done channel
	// when the context is canceled, the goroutine exits
	// on each ticker event, it sends the collected metrics to the server
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
		defer ticker.Stop()

		for i := 0; i < cfg.RateLimit; i++ {
			go controller.SendMetricsWorker(ctx, i)
		}

		for {
			select {
			case <-ticker.C:
				controller.Works <- controller.SendMetricsBatch
			case result := <-controller.Results:
				if result.Err != nil {
					fmt.Println("error sending metrics to server: ", result.Err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	fmt.Printf(`The metrics collection agent is running.
The interval for collecting metrics is %d seconds, the interval for transmitting metrics to the server is %d seconds.
Server address: %s
`, cfg.PollInterval, cfg.ReportInterval, cfg.Host)

	// wait for the context to be canceled
	<-ctx.Done()

	// wait for all goroutines to finish and exit the program
	wg.Wait()

	os.Exit(0)
}
