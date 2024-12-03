package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CheckValue(t *testing.T) {
	t.Run("Float64 value", func(t *testing.T) {
		_, err := CheckValue(Gauge, "100")
		assert.NoError(t, err)
	})
	t.Run("Int64 value", func(t *testing.T) {
		_, err := CheckValue(Counter, "100")
		assert.NoError(t, err)
	})
	t.Run("Wrong value", func(t *testing.T) {
		_, err := CheckValue(Counter, "100.5")
		assert.Error(t, err)
	})
	t.Run("Wrong value", func(t *testing.T) {
		_, err := CheckValue(Gauge, "aa")
		assert.Error(t, err)
	})
	t.Run("Wrong type", func(t *testing.T) {
		_, err := CheckValue("SomeType", "aa")
		assert.Error(t, err)
	})
}

func TestService_CheckMetricType(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
		wanrErr bool
	}{
		{
			name:    "Gauge type test",
			value:   Gauge,
			wantErr: false,
		},
		{
			name:    "Counter type test",
			value:   Counter,
			wantErr: false,
		},
		{
			name:    "Wrong type test",
			value:   "SomeType",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckMetricType(test.value)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
