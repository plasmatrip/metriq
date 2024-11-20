package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return b.Bytes(), nil
}
