package controller

import (
	"bytes"
	"encoding/json"
	"errors"
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
)

type Controller struct {
	Repo   storage.Repository
	Client http.Client
	config config.Config
}

func NewController(repo storage.Repository, config config.Config) *Controller {
	return &Controller{Repo: repo, Client: http.Client{Timeout: config.ClientTimeout}, config: config}
}

func (c Controller) SendMetricsBatch() error {
	metrics, err := c.Repo.Metrics()
	if err != nil {
		return err
	}

	sMetrics := models.SMetrics{}
	sMetrics.Metrics = make([]models.Metrics, 0, len(metrics))
	for mName, metric := range metrics {
		sMetrics.Metrics = append(sMetrics.Metrics, metric.Convert(mName))
	}

	data, err := json.Marshal(sMetrics.Metrics)
	if err != nil {
		return err
	}

	data, err = compress.Compress(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+c.config.Host+"/updates", bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "application/gzip")

	retryCount := 0
	wait := c.config.StartRetryInterval
	for {
		resp, err := c.Client.Do(req)
		if errors.Is(err, syscall.ECONNREFUSED) && retryCount < c.config.MaxRetries {
			time.Sleep(wait)
			retryCount++
			wait += c.config.RetryInterval
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
	metrics, err := c.Repo.Metrics()
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

		req, err := http.NewRequest(http.MethodPost, "http://"+c.config.Host+"/update", bytes.NewReader(data))
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

func (c Controller) UpdateMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	c.Repo.SetMetric("Alloc", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Alloc)})
	c.Repo.SetMetric("TotalAlloc", types.Metric{MetricType: types.Gauge, Value: float64(rtm.TotalAlloc)})
	c.Repo.SetMetric("Sys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Sys)})
	c.Repo.SetMetric("Lookups", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Lookups)})
	c.Repo.SetMetric("Mallocs", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Mallocs)})
	c.Repo.SetMetric("Frees", types.Metric{MetricType: types.Gauge, Value: float64(rtm.Frees)})
	c.Repo.SetMetric("HeapAlloc", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapAlloc)})
	c.Repo.SetMetric("HeapSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapSys)})
	c.Repo.SetMetric("HeapIdle", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapIdle)})
	c.Repo.SetMetric("HeapInuse", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapInuse)})
	c.Repo.SetMetric("HeapReleased", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapReleased)})
	c.Repo.SetMetric("HeapObjects", types.Metric{MetricType: types.Gauge, Value: float64(rtm.HeapObjects)})
	c.Repo.SetMetric("StackInuse", types.Metric{MetricType: types.Gauge, Value: float64(rtm.StackInuse)})
	c.Repo.SetMetric("StackSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.StackSys)})
	c.Repo.SetMetric("MSpanInuse", types.Metric{MetricType: types.Gauge, Value: float64(rtm.MSpanInuse)})
	c.Repo.SetMetric("MSpanSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.MSpanSys)})
	c.Repo.SetMetric("MCacheInuse", types.Metric{MetricType: types.Gauge, Value: float64(rtm.MCacheInuse)})
	c.Repo.SetMetric("MCacheSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.MCacheSys)})
	c.Repo.SetMetric("BuckHashSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.BuckHashSys)})
	c.Repo.SetMetric("GCSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.GCSys)})
	c.Repo.SetMetric("OtherSys", types.Metric{MetricType: types.Gauge, Value: float64(rtm.OtherSys)})
	c.Repo.SetMetric("NextGC", types.Metric{MetricType: types.Gauge, Value: float64(rtm.NextGC)})
	c.Repo.SetMetric("LastGC", types.Metric{MetricType: types.Gauge, Value: float64(rtm.LastGC)})
	c.Repo.SetMetric("PauseTotalNs", types.Metric{MetricType: types.Gauge, Value: float64(rtm.PauseTotalNs)})
	c.Repo.SetMetric("NumGC", types.Metric{MetricType: types.Gauge, Value: float64(rtm.NumGC)})
	c.Repo.SetMetric("NumForcedGC", types.Metric{MetricType: types.Gauge, Value: float64(rtm.NumForcedGC)})
	c.Repo.SetMetric("GCCPUFraction", types.Metric{MetricType: types.Gauge, Value: rtm.GCCPUFraction})
	c.Repo.SetMetric("RandomValue", types.Metric{MetricType: types.Gauge, Value: rand.Float64()})
}
