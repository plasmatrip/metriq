package main

import (
	"time"

	"github.com/plasmatrip/metriq/internal/agent"
	"github.com/plasmatrip/metriq/internal/storage"
)

func main() {
	// httpClient := http.Client()
	controller := agent.NewSender(storage.NewStorage())

	readChan := make(chan struct{})
	sendChan := make(chan struct{})

	go func() {
		for {
			readChan <- struct{}{}
			time.Sleep(agent.ReadTimeout * time.Second)
		}
	}()
	go func() {
		for {
			sendChan <- struct{}{}
			time.Sleep(agent.SendTimeout * time.Second)
		}
	}()

	for {
		select {
		case <-readChan:
			controller.UpdateMetrics()
		case <-sendChan:
			controller.SendMetrics()
		}
	}
}
