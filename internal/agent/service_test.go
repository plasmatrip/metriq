package agent

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/types"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	storage.MemStorage
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		storage.MemStorage{
			Mu:      sync.RWMutex{},
			Storage: map[string]types.Metric{"metric": {MetricType: types.Gauge, Value: 100}, "counter": {MetricType: types.Counter, Value: 100}},
		},
	}
}

func TestService_SendMetrics(t *testing.T) {
	mock := NewMockStorage()
	mock.SetMetric("metric", types.Metric{MetricType: types.Gauge, Value: 100})
	mock.SetMetric("counter", types.Metric{MetricType: types.Counter, Value: 100})

	controller := NewController(&MockStorage{}, Config{})

	t.Run("Send metrics test", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method, "Only POST requests are allowed!")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		err := controller.SendMetrics()
		assert.NoError(t, err, "No error expected when sending metrics")
	})
}
