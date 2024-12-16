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
	"github.com/plasmatrip/metriq/internal/storage/mem"
)

func main() {
	// создаем контекст и прослишиваем сигналы прекращения работы программы
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	controller := controller.NewController(mem.NewStorage(), *cfg)

	var wg sync.WaitGroup

	// запускаем горутину сбора метрик
	// в цикле читаем тикер и отмену контекста, при отмене контекста заверашем горутину
	// по тикеру обновляем метрики
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

	// запускаем горутину сбора метрик c помощью пакета gopsutil
	// в цикле читаем тикер и отмену контекста, при отмене контекста заверашем горутину
	// по тикеру обновляем метрики
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

	// запускаем горутину отправки метрик
	// в цикле читаем тикер и отмену контекста, при отмене контекста заверашем горутину
	// по тикеру отправляем метрики на сервер
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

	// ждем отмены контекста
	<-ctx.Done()

	// ожидаем завершения горутин и выходим из программы
	wg.Wait()

	os.Exit(0)
}
