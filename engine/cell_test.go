package engine_test

import (
	"github.com/clambin/criticalmass/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCell(t *testing.T) {
	c := engine.Cell{
		Critical: 9,
	}

	count := 0
	for !c.IsCritical() {
		count++
		ok := c.Add(engine.PlayerA, 1, false)
		assert.True(t, ok)
	}

	assert.Equal(t, 9, count)
	assert.Zero(t, c.Remaining())

	ok := c.Add(engine.PlayerB, 1, false)
	assert.False(t, ok)
}
