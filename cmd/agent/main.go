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
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/template"
	"time"

	_ "net/http/pprof"

	"github.com/plasmatrip/metriq/internal/agent/config"
	"github.com/plasmatrip/metriq/internal/agent/controller"
	"github.com/plasmatrip/metriq/internal/agent/grpc"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/storage/mem"
)

var buildVersion string
var buildDate string
var buildCommit string

const buildInfo = `
	Agent build info
	Build version: {{if .BuildVersion}}{{.BuildVersion}}{{else}}"N/A"{{end}}
	Build date: {{if .BuildDate}}{{.BuildDate}}{{else}}"N/A"{{end}}
	Build commit: {{if .BuildCommit}}{{.BuildCommit}}{{else}}"N/A"{{end}}
	
`

// main initializes the metrics collection agent. It sets up a context to listen
// for termination signals and configures the agent using settings from the
// configuration module. The function starts three goroutines: one for collecting
// application metrics, one for collecting system metrics using gopsutil, and
// one for sending metrics to the server. Each goroutine runs periodically based
// on configuration-defined intervals. The function waits for a stop signal to
// gracefully shut down all goroutines and exit.
func main() {
	// Create a context to listen for termination signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	t := template.Must(template.New("buildInfo").Parse(buildInfo))

	data := struct {
		BuildVersion string
		BuildDate    string
		BuildCommit  string
	}{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}

	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Close()

	stor := mem.NewStorage()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	controller := controller.NewController(stor, *cfg)

	var grpcClient *grpc.GRPCClient
	if cfg.EnableGRPC {
		grpcClient, err = grpc.NewGRPCClient(cfg.GRPCPort, stor, *cfg, log)
		if err != nil {
			log.Sugar.Errorw("failed to create grpc client", "error", err)
			return
		}
		log.Sugar.Infow("grpc client created", "port", cfg.GRPCPort)
		defer grpcClient.Close()
	}

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
					log.Sugar.Errorw("error while collecting system utilization metrics using gopsutil", "error", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// start goroutine to send metrics to the server
	// in a loop, it reads from a ticker and the context's Done channel
	// when the context is canceled, the goroutine exits
	// on each ticker event, it sends the collected metrics to the server±~
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
		defer ticker.Stop()

		for i := 0; i < cfg.RateLimit; i++ {
			go controller.SendMetricsWorker(ctx, &wg, i)
		}

		for {
			select {
			case <-ticker.C:
				if grpcClient != nil {
					controller.Works <- grpcClient.SendMetrics
				} else {
					controller.Works <- controller.SendMetricsBatch
				}
			case result := <-controller.Results:
				if result.Err != nil {
					log.Sugar.Errorw("error sending metrics to server", "error", result.Err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Sugar.Info("The metrics collection agent has started")
	log.Sugar.Infow("Agent config", "poll interval", cfg.PollInterval, "report interval", cfg.ReportInterval, "server address", cfg.Host)

	// wait for the context to be canceled
	<-ctx.Done()

	// wait for all goroutines to finish and exit the program
	wg.Wait()

	log.Sugar.Info("The agent has been shut down gracefully")
}
