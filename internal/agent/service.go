package agent

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime"

	"github.com/plasmatrip/metriq/internal/server"
	"github.com/plasmatrip/metriq/internal/storage"
)

type Controller struct {
	Repo   storage.Repository
	Client http.Client
}

func NewController(repo storage.Repository) *Controller {
	return &Controller{Repo: repo, Client: http.Client{}}
}

func (c *Controller) SendMetrics() {

	for v, k := range c.Repo.GetAll() {
		mType := server.Gauge
		if v == storage.PollCount {
			mType = server.Counter
		}

		url := fmt.Sprint("http://", server.Address, ":", server.Port, "/update/", mType, "/", v, "/", k)

		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		req.Header.Set("Content-Type", "text/plain")

		resp, err := c.Client.Do(req)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer resp.Body.Close()

		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}

func (c *Controller) UpdateMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	c.Repo.Update("Alloc", storage.Gauge(rtm.Alloc))
	c.Repo.Update("TotalAlloc", storage.Gauge(rtm.TotalAlloc))
	c.Repo.Update("Sys", storage.Gauge(rtm.Sys))
	c.Repo.Update("Lookups", storage.Gauge(rtm.Lookups))
	c.Repo.Update("Mallocs", storage.Gauge(rtm.Mallocs))
	c.Repo.Update("Frees", storage.Gauge(rtm.Frees))
	c.Repo.Update("HeapAlloc", storage.Gauge(rtm.HeapAlloc))
	c.Repo.Update("HeapSys", storage.Gauge(rtm.HeapSys))
	c.Repo.Update("HeapIdle", storage.Gauge(rtm.HeapIdle))
	c.Repo.Update("HeapInuse", storage.Gauge(rtm.HeapInuse))
	c.Repo.Update("HeapReleased", storage.Gauge(rtm.HeapReleased))
	c.Repo.Update("HeapObjects", storage.Gauge(rtm.HeapObjects))
	c.Repo.Update("StackInuse", storage.Gauge(rtm.StackInuse))
	c.Repo.Update("StackSys", storage.Gauge(rtm.StackSys))
	c.Repo.Update("MSpanInuse", storage.Gauge(rtm.MSpanInuse))
	c.Repo.Update("MSpanSys", storage.Gauge(rtm.MSpanSys))
	c.Repo.Update("MCacheInuse", storage.Gauge(rtm.MCacheInuse))
	c.Repo.Update("MCacheSys", storage.Gauge(rtm.MCacheSys))
	c.Repo.Update("BuckHashSys", storage.Gauge(rtm.BuckHashSys))
	c.Repo.Update("GCSys", storage.Gauge(rtm.GCSys))
	c.Repo.Update("OtherSys", storage.Gauge(rtm.OtherSys))
	c.Repo.Update("NextGC", storage.Gauge(rtm.NextGC))
	c.Repo.Update("LastGC", storage.Gauge(rtm.LastGC))
	c.Repo.Update("PauseTotalNs", storage.Gauge(rtm.PauseTotalNs))
	c.Repo.Update("NumGC", storage.Gauge(rtm.NumGC))
	c.Repo.Update("NumForcedGC", storage.Gauge(rtm.NumForcedGC))
	c.Repo.Update("GCCPUFraction", storage.Gauge(rtm.GCCPUFraction))
	c.Repo.Update("RandomValue", storage.Gauge(rand.Float64()))
}
