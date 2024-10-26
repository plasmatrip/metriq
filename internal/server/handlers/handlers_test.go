package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plasmatrip/metriq/internal/server"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandlers(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
		url  string
	}{
		{
			name: "Status ok test",
			url:  server.Address + ":" + server.Port + "/update/gauge/metric/100",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "No name metrics test",
			url:  server.Address + ":" + server.Port + "/update/gauge//100",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name: "Wrong metrics type",
			url:  server.Address + ":" + server.Port + "/update/gaaauge/metric/100",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name: "Wrong value",
			url:  server.Address + ":" + server.Port + "/update/counter/metric/100.5",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name: "Wrong value",
			url:  server.Address + ":" + server.Port + "/update/counter/metric/aa",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}
	handlers := NewHandlers(storage.NewStorage())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)

			w := httptest.NewRecorder()
			w.Header().Set(`Content-Type`, `text/plain`)
			handlers.UpdateHandler(w, request)

			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
		})
	}
}
