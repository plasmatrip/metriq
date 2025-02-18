package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/storage/mem"
	"github.com/plasmatrip/metriq/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingHandler(t *testing.T) {
	h := NewHandlers(mem.NewStorage(), config.Config{}, logger.Logger{})
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", h.Ping)
	serv := httptest.NewServer(mux)
	defer serv.Close()

	t.Run("Ping test", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodGet, serv.URL+"/ping", nil)
		assert.NoError(t, err)

		res, err := serv.Client().Do(request)
		assert.NoError(t, err)
		defer res.Body.Close()
	})

}

func TestValueHandler(t *testing.T) {
	tests := []struct {
		name string
		want int
		url  string
	}{
		{
			name: "Status ok test",
			url:  "/value/gauge/metric",
			want: http.StatusOK,
		},
		{
			name: "Wrong metrics name",
			url:  "/value/gauge/wrong",
			want: http.StatusNotFound,
		},
		{
			name: "Wrong metrics type",
			url:  "/value/gaaauge/metric",
			want: http.StatusBadRequest,
		},
	}

	storage := mem.NewStorage()
	storage.SetMetric(context.TODO(), "metric", types.Metric{MetricType: types.Gauge, Value: float64(100)})
	storage.SetMetric(context.TODO(), "counter", types.Metric{MetricType: types.Counter, Value: int64(100)})

	log, err := logger.NewLogger()
	require.NoError(t, err)

	h := NewHandlers(storage, config.Config{}, log)
	mux := http.NewServeMux()
	mux.HandleFunc("/value/{metricType}/{metricName}", h.Value)
	serv := httptest.NewServer(mux)
	defer serv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, serv.URL+test.url, nil)
			assert.NoError(t, err)

			res, err := serv.Client().Do(request)
			require.NotNil(t, res)
			assert.NoError(t, err)
			assert.Equal(t, test.want, res.StatusCode, res)
			defer res.Body.Close()
		})
	}
}

func TestJSONValueHandler(t *testing.T) {
	tests := []struct {
		name string
		want int
		data map[string]interface{}
	}{
		{
			name: "Status ok test",
			want: http.StatusOK,
			data: map[string]interface{}{"id": "metric", "type": "gauge"},
		},
		{
			name: "Wrong metrics name",
			want: http.StatusNotFound,
			data: map[string]interface{}{"id": "wrong", "type": "gauge"},
		},
		{
			name: "Wrong metrics type",
			want: http.StatusBadRequest,
			data: map[string]interface{}{"id": "metric", "type": "wrong"},
		},
		{
			name: "Wrong data",
			want: http.StatusBadRequest,
			data: map[string]interface{}{"metric": "wrong"},
		},
		{
			name: "Wrong JSON",
			want: http.StatusBadRequest,
			data: map[string]interface{}{"metric": 42},
		},
	}

	storage := mem.NewStorage()
	storage.SetMetric(context.TODO(), "metric", types.Metric{MetricType: types.Gauge, Value: float64(100)})
	storage.SetMetric(context.TODO(), "counter", types.Metric{MetricType: types.Counter, Value: int64(100)})

	log, err := logger.NewLogger()
	require.NoError(t, err)

	h := NewHandlers(storage, config.Config{}, log)
	mux := http.NewServeMux()
	mux.HandleFunc("/value", h.JSONValue)
	serv := httptest.NewServer(mux)
	defer serv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, err := json.Marshal(test.data)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, serv.URL+"/value", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			res, err := serv.Client().Do(request)
			require.NotNil(t, res)
			assert.NoError(t, err)
			assert.Equal(t, test.want, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestJSONUpdateHandler(t *testing.T) {
	tests := []struct {
		name string
		want int
		data map[string]interface{}
	}{
		{
			name: "Status ok test",
			want: http.StatusOK,
			data: map[string]interface{}{"id": "metric", "type": "gauge", "value": 10},
		},
		{
			name: "Wrong metrics name",
			want: http.StatusNotFound,
			data: map[string]interface{}{"id": "wrong", "type": "gauge", "value": 10},
		},
		{
			name: "Wrong metrics type",
			want: http.StatusBadRequest,
			data: map[string]interface{}{"id": "metric", "type": "wrong", "delta": 10},
		},
		{
			name: "Wrong data",
			want: http.StatusBadRequest,
			data: map[string]interface{}{"metric": "wrong", "delta": 10},
		},
		{
			name: "Wrong JSON",
			want: http.StatusBadRequest,
			data: map[string]interface{}{"metric": 42, "delta": "aa"},
		},
	}

	storage := mem.NewStorage()
	storage.SetMetric(context.TODO(), "metric", types.Metric{MetricType: types.Gauge, Value: float64(100)})
	storage.SetMetric(context.TODO(), "counter", types.Metric{MetricType: types.Counter, Value: int64(100)})

	log, err := logger.NewLogger()
	require.NoError(t, err)

	h := NewHandlers(storage, config.Config{}, log)
	mux := http.NewServeMux()
	mux.HandleFunc("/update", h.JSONUpdate)
	serv := httptest.NewServer(mux)
	defer serv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, err := json.Marshal(test.data)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, serv.URL+"/update", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			res, err := serv.Client().Do(request)
			// require.NotNil(t, res)
			assert.NoError(t, err)
			assert.Equal(t, test.want, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name string
		want int
		url  string
	}{
		{
			name: "Status ok test",
			url:  "/update/gauge/metric/100",
			want: http.StatusOK,
		},
		{
			name: "No name metrics test",
			url:  "/update/gauge//100",
			want: http.StatusNotFound,
		},
		{
			name: "Wrong metrics type",
			url:  "/update/gaaauge/metric/100",
			want: http.StatusBadRequest,
		},
		{
			name: "Wrong value",
			url:  "/update/counter/metric/100.5",
			want: http.StatusBadRequest,
		},
		{
			name: "Wrong value",
			url:  "/update/counter/metric/aa",
			want: http.StatusBadRequest,
		},
	}
	h := NewHandlers(mem.NewStorage(), config.Config{}, logger.Logger{})
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", h.Update)
	serv := httptest.NewServer(mux)
	defer serv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request, err := http.NewRequest(http.MethodPost, serv.URL+test.url, nil)
			assert.NoError(t, err)

			res, err := serv.Client().Do(request)
			require.NotNil(t, res)
			assert.NoError(t, err)
			defer res.Body.Close()
		})
	}
}

func BenchmarkUpdateHandler(b *testing.B) {
	h := NewHandlers(mem.NewStorage(), config.Config{}, logger.Logger{})
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", h.Update)
	serv := httptest.NewServer(mux)
	defer serv.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		request, err := http.NewRequest(http.MethodPost, serv.URL+"/update/gauge/metric/100", nil)
		assert.NoError(b, err)

		res, err := serv.Client().Do(request)
		require.NotNil(b, res)
		assert.NoError(b, err)
		defer res.Body.Close()
	}
}
