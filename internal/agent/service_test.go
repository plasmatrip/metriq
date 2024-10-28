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
			Mu:      sync.Mutex{},
			Storage: map[string]storage.Metric{"metric": {MetricType: types.Gauge, Value: 100}, "counter": {MetricType: types.Counter, Value: 100}},
		},
	}
}

func TestService_SendMetrics(t *testing.T) {
	mock := NewMockStorage()
	mock.Update("metric", storage.Metric{MetricType: types.Gauge, Value: 100})
	mock.Update("counter", storage.Metric{MetricType: types.Counter, Value: 100})

	controller := NewController(&MockStorage{})

	t.Run("Send metrics test", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method, "Only POST requests are allowed!")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		err := controller.SendMetrics(server.URL)
		assert.NoError(t, err, "No error expected when sending metrics")
	})

}
