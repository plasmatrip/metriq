package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"syscall"
	"time"

	"github.com/plasmatrip/metriq/internal/agent/compress"
	"github.com/plasmatrip/metriq/internal/agent/config"
	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/types"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type Result struct {
	Err error
}

type Controller struct {
	Repo    storage.Repository
	Client  http.Client
	cfg     config.Config
	Works   chan func() error
	Results chan Result
}

func NewController(repo storage.Repository, cfg config.Config) *Controller {
	return &Controller{Repo: repo, Client: http.Client{Timeout: cfg.ClientTimeout}, cfg: cfg, Works: make(chan func() error), Results: make(chan Result)}
}

func (c Controller) SendMetricsWorker(ctx context.Context, idx int) {
	for {
		select {
		case work := <-c.Works:
			fmt.Printf("worker %d start\n", idx)
			err := work()
			result := Result{}
			if err != nil {
				result.Err = err
			}
			c.Results <- result
			fmt.Printf("worker %d end\n", idx)
		case <-ctx.Done():
			fmt.Printf("worker %d stop\n", idx)
			return
		}
	}
}

func (c Controller) SendMetricsBatch() error {
	fmt.Println("send: ", time.Now().Format("15:04:05"))

	metrics, err := c.Repo.Metrics(context.Background())
	if len(metrics) == 0 {
		return nil
	}

	if err != nil {
		return err
	}

	sMetrics := make([]models.Metrics, 0, len(metrics))
	for mName, metric := range metrics {
		sMetrics = append(sMetrics, metric.Convert(mName))
	}

	data, err := json.Marshal(sMetrics)
	if err != nil {
		return err
	}

	data, err = compress.Compress(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+c.cfg.Host+"/updates", bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "application/gzip")

	// если есть ключ, хэшируем request body
	if len(c.cfg.Key) > 0 {
		copyBody, err := req.GetBody()
		if err != nil {
			return err
		}

		hash, err := c.Sum(copyBody)
		if err != nil {
			return err
		}

		req.Header.Set("HashSHA256", hash)
	}

	// в цикле пытаемся отправить на серевер метрики
	// количество попыток, интервал в сек между попытками настраивается в конфиге
	retryCount := 0
	wait := c.cfg.StartRetryInterval
	for {
		resp, err := c.Client.Do(req)
		if errors.Is(err, syscall.ECONNREFUSED) && retryCount < c.cfg.MaxRetries {
			time.Sleep(wait)
			retryCount++
			wait += c.cfg.RetryInterval
			continue
		}
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	}
}

func (c Controller) SendMetrics() error {
	metrics, err := c.Repo.Metrics(context.Background())
	if err != nil {
		return err
	}

	for mName, metric := range metrics {
		jMetric := metric.Convert(mName)
		data, err := json.Marshal(jMetric)
		if err != nil {
			return err
		}

		data, err = compress.Compress(data)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, "http://"+c.cfg.Host+"/update", bytes.NewReader(data))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "application/gzip")

		resp, err := c.Client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}
	return nil
}

func (c Controller) UpdateMetrics(ctx context.Context) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	c.Repo.SetMetric(ctx, "Alloc", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Alloc)})
	c.Repo.SetMetric(ctx, "TotalAlloc", types.Metric{MetricType: types.Gauge, Value: float64(rtm.TotalAlloc)})
	c.Repo.SetMetric(ctx, "Sys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Sys)})
	c.Repo.SetMetric(ctx, "Lookups", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Lookups)})
	c.Repo.SetMetric(ctx, "Mallocs", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Mallocs)})
	c.Repo.SetMetric(ctx, "Frees", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Frees)})
	c.Repo.SetMetric(ctx, "HeapAlloc", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapAlloc)})
	c.Repo.SetMetric(ctx, "HeapSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapSys)})
	c.Repo.SetMetric(ctx, "HeapIdle", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapIdle)})
	c.Repo.SetMetric(ctx, "HeapInuse", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapInuse)})
	c.Repo.SetMetric(ctx, "HeapReleased", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapReleased)})
	c.Repo.SetMetric(ctx, "HeapObjects", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapObjects)})
	c.Repo.SetMetric(ctx, "StackInuse", types.Metric{MetricType: types.Gauge, Value: float64(rtm.StackInuse)})
	c.Repo.SetMetric(ctx, "StackSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.StackSys)})
	c.Repo.SetMetric(ctx, "MSpanInuse", types.Metric{MetricType: types.Gauge, Value: float64(rtm.MSpanInuse)})
	c.Repo.SetMetric(ctx, "MSpanSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.MSpanSys)})
	c.Repo.SetMetric(ctx, "MCacheInuse", types.Metric{MetricType: types.Gauge, Value: float64(rtm.MCacheInuse)})
	c.Repo.SetMetric(ctx, "MCacheSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.MCacheSys)})
	c.Repo.SetMetric(ctx, "BuckHashSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.BuckHashSys)})
	c.Repo.SetMetric(ctx, "GCSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.GCSys)})
	c.Repo.SetMetric(ctx, "OtherSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.OtherSys)})
	c.Repo.SetMetric(ctx, "NextGC", types.Metric{MetricType: types.Gauge, Value: float64(rtm.NextGC)})
	c.Repo.SetMetric(ctx, "LastGC", types.Metric{MetricType: types.Gauge, Value: float64(rtm.LastGC)})
	c.Repo.SetMetric(ctx, "PauseTotalNs", types.Metric{MetricType: types.Gauge, Value: float64(rtm.PauseTotalNs)})
	c.Repo.SetMetric(ctx, "NumGC", types.Metric{MetricType: types.Gauge, Value: float64(rtm.NumGC)})
	c.Repo.SetMetric(ctx, "NumForcedGC", types.Metric{MetricType: types.Gauge, Value: float64(rtm.NumForcedGC)})
	c.Repo.SetMetric(ctx, "GCCPUFraction", types.Metric{MetricType: types.Gauge, Value: rtm.GCCPUFraction})
	c.Repo.SetMetric(ctx, "RandomValue", types.Metric{MetricType: types.Gauge, Value: rand.Float64()})
}

func (c Controller) UpdatePSMetrics(ctx context.Context) error {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	cpu, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}

	c.Repo.SetMetric(ctx, "TotalMemory", types.Metric{MetricType: types.Gauge, Value: float64(mem.Total)})
	c.Repo.SetMetric(ctx, "FreeMemory", types.Metric{MetricType: types.Gauge, Value: float64(mem.Free)})
	c.Repo.SetMetric(ctx, "CPUutilization1", types.Metric{MetricType: types.Gauge, Value: cpu[0]})
	return nil
}
