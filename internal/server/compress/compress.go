package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/plasmatrip/metriq/internal/logger"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	w.Header().Set("Content-Encoding", "gzip")
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(data []byte) (int, error) {
	return c.zw.Write(data)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c *compressReader) Read(data []byte) (int, error) {
	return c.zr.Read(data)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// WithCompression устанавливает сжатие
func WithCompression(log logger.Logger) func(next http.Handler) http.Handler {
	log.Sugar.Debug("compression started")

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Устанавливаем сжатие
			ow := w
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportGzip := strings.Contains(acceptEncoding, "gzip")
			if supportGzip {
				cw := newCompressWriter(w)
				defer func() {
					if err := cw.Close(); err != nil {
						log.Sugar.Infow("failed to close compress writer", "error", err)
					}
				}()
				ow = cw
			}

			// Проверяем, что клиент отправляет сжатый контент
			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					log.Sugar.Infow("failed to create compress reader", "error", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer func() {
					if err := cr.Close(); err != nil {
						log.Sugar.Infow("failed to close compress reader", "error", err)
					}
				}()
				r.Body = cr
			}

			next.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
