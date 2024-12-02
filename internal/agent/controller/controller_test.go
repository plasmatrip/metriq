package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/plasmatrip/metriq/internal/agent/config"
	"github.com/plasmatrip/metriq/internal/storage/mem"
	"github.com/plasmatrip/metriq/internal/types"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	mem.MemStorage
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		mem.MemStorage{
			Mu:      sync.RWMutex{},
			Storage: map[string]types.Metric{},
		},
	}
}

func TestService_SendMetrics(t *testing.T) {
	ctx := context.Background()
	mock := NewMockStorage()
	mock.SetMetric(ctx, "metric", types.Metric{MetricType: types.Gauge, Value: 100})
	mock.SetMetric(ctx, "counter", types.Metric{MetricType: types.Counter, Value: 100})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Only POST requests are allowed!")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Content-Encoding"), "application/gzip")
		w.WriteHeader(http.StatusOK)
	}))

	controller := NewController(mock, config.Config{Host: strings.Split(server.URL, "//")[1]})
	controller.Client = *server.Client()

	t.Run("Send metrics test", func(t *testing.T) {
		err := controller.SendMetrics()
		assert.NoError(t, err, "No error expected when sending metrics")
	})

	server.Close()

	t.Run("Send metrics error test", func(t *testing.T) {
		err := controller.SendMetrics()
		assert.Error(t, err, "Error expected if metrics are not sending")
	})
}

func TestService_UpdateMetrics(t *testing.T) {
	mock := NewMockStorage()

	controller := NewController(mock, config.Config{})

	ctx := context.Background()

	t.Run("Send metrics test", func(t *testing.T) {
		metrics, err := mock.Metrics(ctx)
		assert.NoError(t, err)
		assert.Empty(t, metrics)
		controller.UpdateMetrics(ctx)
		metrics, err = mock.Metrics(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, metrics)
		controller.UpdateMetrics(ctx)
		newMetrics, err := mock.Metrics(ctx)
		assert.NoError(t, err)
		assert.NotEqual(t, metrics, newMetrics)
	})
}
