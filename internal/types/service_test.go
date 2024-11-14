package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CheckType(t *testing.T) {
	t.Run("Gauge type test", func(t *testing.T) {
		assert.True(t, checkType(Gauge))
	})

	t.Run("Counter type test", func(t *testing.T) {
		assert.True(t, checkType(Counter))
	})

	t.Run("Wrong type test", func(t *testing.T) {
		assert.False(t, checkType("wrong type"))
	})
}

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
}
