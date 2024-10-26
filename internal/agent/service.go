package agent

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime"

	"github.com/plasmatrip/metriq/internal/storage"
)

type Controller struct {
	Repo   storage.Repository
	Client http.Client
}

func NewController(repo storage.Repository) *Controller {
	return &Controller{Repo: repo, Client: http.Client{}}
}

func (c *Controller) SendMetrics(server string) error {
	for v, k := range c.Repo.GetGauges() {
		url := fmt.Sprint(server, "/update/gauge", "/", v, "/", k)
		c.send(url)
	}

	url := fmt.Sprint(server, "/update/counter/", storage.PollCount, "/", c.Repo.GetCounter())
	if err := c.send(url); err != nil {
		return err
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

	c.Repo.UpdateGauge("Alloc", storage.Gauge(rtm.Alloc))
	c.Repo.UpdateGauge("TotalAlloc", storage.Gauge(rtm.TotalAlloc))
	c.Repo.UpdateGauge("Sys", storage.Gauge(rtm.Sys))
	c.Repo.UpdateGauge("Lookups", storage.Gauge(rtm.Lookups))
	c.Repo.UpdateGauge("Mallocs", storage.Gauge(rtm.Mallocs))
	c.Repo.UpdateGauge("Frees", storage.Gauge(rtm.Frees))
	c.Repo.UpdateGauge("HeapAlloc", storage.Gauge(rtm.HeapAlloc))
	c.Repo.UpdateGauge("HeapSys", storage.Gauge(rtm.HeapSys))
	c.Repo.UpdateGauge("HeapIdle", storage.Gauge(rtm.HeapIdle))
	c.Repo.UpdateGauge("HeapInuse", storage.Gauge(rtm.HeapInuse))
	c.Repo.UpdateGauge("HeapReleased", storage.Gauge(rtm.HeapReleased))
	c.Repo.UpdateGauge("HeapObjects", storage.Gauge(rtm.HeapObjects))
	c.Repo.UpdateGauge("StackInuse", storage.Gauge(rtm.StackInuse))
	c.Repo.UpdateGauge("StackSys", storage.Gauge(rtm.StackSys))
	c.Repo.UpdateGauge("MSpanInuse", storage.Gauge(rtm.MSpanInuse))
	c.Repo.UpdateGauge("MSpanSys", storage.Gauge(rtm.MSpanSys))
	c.Repo.UpdateGauge("MCacheInuse", storage.Gauge(rtm.MCacheInuse))
	c.Repo.UpdateGauge("MCacheSys", storage.Gauge(rtm.MCacheSys))
	c.Repo.UpdateGauge("BuckHashSys", storage.Gauge(rtm.BuckHashSys))
	c.Repo.UpdateGauge("GCSys", storage.Gauge(rtm.GCSys))
	c.Repo.UpdateGauge("OtherSys", storage.Gauge(rtm.OtherSys))
	c.Repo.UpdateGauge("NextGC", storage.Gauge(rtm.NextGC))
	c.Repo.UpdateGauge("LastGC", storage.Gauge(rtm.LastGC))
	c.Repo.UpdateGauge("PauseTotalNs", storage.Gauge(rtm.PauseTotalNs))
	c.Repo.UpdateGauge("NumGC", storage.Gauge(rtm.NumGC))
	c.Repo.UpdateGauge("NumForcedGC", storage.Gauge(rtm.NumForcedGC))
	c.Repo.UpdateGauge("GCCPUFraction", storage.Gauge(rtm.GCCPUFraction))
	c.Repo.UpdateGauge("RandomValue", storage.Gauge(rand.Float64()))
}
