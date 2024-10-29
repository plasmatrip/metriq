package agent

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime"

	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/types"
)

type Controller struct {
	Repo   storage.Repository
	Client http.Client
}

func NewController(repo storage.Repository) *Controller {
	return &Controller{Repo: repo, Client: http.Client{}}
}

func (c *Controller) SendMetrics(server string) error {
	for name, metric := range c.Repo.GetAll() {
		var path string
		switch metric.MetricType {
		case types.Gauge:
			path = "/update/gauge/"
		case types.Counter:
			path = "/update/counter/"
		}
		if err := c.send(fmt.Sprint(server, path, name, "/", metric.Value)); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) send(url string) error {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

func (c *Controller) UpdateMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	c.Repo.Update("Alloc", storage.Metric{MetricType: types.Gauge, Value: rtm.Alloc})
	c.Repo.Update("TotalAlloc", storage.Metric{MetricType: types.Gauge, Value: rtm.TotalAlloc})
	c.Repo.Update("Sys", storage.Metric{MetricType: types.Gauge, Value: rtm.Sys})
	c.Repo.Update("Lookups", storage.Metric{MetricType: types.Gauge, Value: rtm.Lookups})
	c.Repo.Update("Mallocs", storage.Metric{MetricType: types.Gauge, Value: rtm.Mallocs})
	c.Repo.Update("Frees", storage.Metric{MetricType: types.Gauge, Value: rtm.Frees})
	c.Repo.Update("HeapAlloc", storage.Metric{MetricType: types.Gauge, Value: rtm.HeapAlloc})
	c.Repo.Update("HeapSys", storage.Metric{MetricType: types.Gauge, Value: rtm.HeapSys})
	c.Repo.Update("HeapIdle", storage.Metric{MetricType: types.Gauge, Value: rtm.HeapIdle})
	c.Repo.Update("HeapInuse", storage.Metric{MetricType: types.Gauge, Value: rtm.HeapInuse})
	c.Repo.Update("HeapReleased", storage.Metric{MetricType: types.Gauge, Value: rtm.HeapReleased})
	c.Repo.Update("HeapObjects", storage.Metric{MetricType: types.Gauge, Value: rtm.HeapObjects})
	c.Repo.Update("StackInuse", storage.Metric{MetricType: types.Gauge, Value: rtm.StackInuse})
	c.Repo.Update("StackSys", storage.Metric{MetricType: types.Gauge, Value: rtm.StackSys})
	c.Repo.Update("MSpanInuse", storage.Metric{MetricType: types.Gauge, Value: rtm.MSpanInuse})
	c.Repo.Update("MSpanSys", storage.Metric{MetricType: types.Gauge, Value: rtm.MSpanSys})
	c.Repo.Update("MCacheInuse", storage.Metric{MetricType: types.Gauge, Value: rtm.MCacheInuse})
	c.Repo.Update("MCacheSys", storage.Metric{MetricType: types.Gauge, Value: rtm.MCacheSys})
	c.Repo.Update("BuckHashSys", storage.Metric{MetricType: types.Gauge, Value: rtm.BuckHashSys})
	c.Repo.Update("GCSys", storage.Metric{MetricType: types.Gauge, Value: rtm.GCSys})
	c.Repo.Update("OtherSys", storage.Metric{MetricType: types.Gauge, Value: rtm.OtherSys})
	c.Repo.Update("NextGC", storage.Metric{MetricType: types.Gauge, Value: rtm.NextGC})
	c.Repo.Update("LastGC", storage.Metric{MetricType: types.Gauge, Value: rtm.LastGC})
	c.Repo.Update("PauseTotalNs", storage.Metric{MetricType: types.Gauge, Value: rtm.PauseTotalNs})
	c.Repo.Update("NumGC", storage.Metric{MetricType: types.Gauge, Value: rtm.NumGC})
	c.Repo.Update("NumForcedGC", storage.Metric{MetricType: types.Gauge, Value: rtm.NumForcedGC})
	c.Repo.Update("GCCPUFraction", storage.Metric{MetricType: types.Gauge, Value: rtm.GCCPUFraction})
	c.Repo.Update("RandomValue", storage.Metric{MetricType: types.Gauge, Value: rand.Float64()})
}
