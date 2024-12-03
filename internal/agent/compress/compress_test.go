package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Uncompress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed init uncompress reder: %v", err)
	}

	_, err = b.ReadFrom(w)
	if err != nil {
		return nil, fmt.Errorf("failed read data to uncompress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed uncompress data: %v", err)
	}

	return b.Bytes(), nil
}

func TestCompress(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "Valid test 1",
			value: "test string",
			want:  "test string",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			compress, err := Compress([]byte(test.value))
			assert.NoError(t, err)
			uncompress, err := Uncompress(compress)
			assert.NoError(t, err)
			assert.Equal(t, test.want, string(uncompress))
		})
	}
}
