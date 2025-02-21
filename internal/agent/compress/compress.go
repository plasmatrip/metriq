// Compress takes a byte slice of data and returns a new byte slice with the data
// compressed using the gzip algorithm. This method is useful for reducing the size
// of data before storage or transmission over a network. The function initializes
// a gzip writer with the best compression level, writes the input data to the writer,
// and then closes the writer to flush and finalize the compression. If any errors
// occur during this process, such as during the initialization of the gzip writer,
// writing to it, or closing it, the function returns an error. Otherwise, it returns
// the compressed data as a byte slice. The compressed data is in gzip format, which
// is a commonly used format for data compression.

package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

// Compress takes a byte slice of data and returns a new byte slice with the same data,
// but compressed using the gzip algorithm. The returned byte slice is a gzip
// formatted stream, which can be safely written to a file or sent over a network
// connection. If the function encounters any errors during the compression process,
// it will return an error.
//
// Note that the returned byte slice is not the same as the original data, but
// rather a gzip formatted stream that can be decompressed to obtain the original
// data.
//
// The function takes a byte slice as input and returns a new byte slice with the
// compressed data, or an error if the compression fails.
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
