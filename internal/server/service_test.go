package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CheckType(t *testing.T) {
	t.Run("Gauge type test", func(t *testing.T) {
		assert.True(t, CheckType(Gauge))
	})

	t.Run("Counter type test", func(t *testing.T) {
		assert.True(t, CheckType(Counter))
	})

	t.Run("Wrong type test", func(t *testing.T) {
		assert.False(t, CheckType("wrong type"))
	})
}

func TestService_CheckValue(t *testing.T) {
	t.Run("Float64 value", func(t *testing.T) {
		assert.NoError(t, CheckValue(Gauge, "100"))
	})
	t.Run("Int64 value", func(t *testing.T) {
		assert.NoError(t, CheckValue(Counter, "100"))
	})
	t.Run("Wrong value", func(t *testing.T) {
		assert.Error(t, CheckValue(Counter, "100.5"))
	})
	t.Run("Wrong value", func(t *testing.T) {
		assert.Error(t, CheckValue(Gauge, "aa"))
	})
}
