package engine_test

import (
	"github.com/clambin/criticalmass/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoard(t *testing.T) {
	board := engine.New(3, 3)

	for i := 0; i < 3; i++ {
		board.Add(engine.PlayerA, engine.Coordinate{Row: 0, Column: 0})
		board.Add(engine.PlayerB, engine.Coordinate{Row: 0, Column: 2})
	}
	_, won := board.Winner()
	assert.False(t, won)

	assert.Equal(t, `A3  0 B3 
 0  0  0 
 0  0  0 
`, board.Dump())

	board.Add(engine.PlayerA, engine.Coordinate{Row: 0, Column: 0})

	c := board.GetCriticals()
	require.Len(t, c, 4)
	assert.Equal(t, engine.Coordinate{Row: 0, Column: 0}, c[0])

	board.ProcessCriticals()

	assert.Equal(t, `A1 A1 B3 
A1 A1  0 
 0  0  0 
`, board.Dump())

	board.Add(engine.PlayerB, engine.Coordinate{Row: 0, Column: 2})
	board.ProcessCriticals()

	assert.Equal(t, `A1 B2 B1 
A1 B2 B1 
 0  0  0 
`, board.Dump())

	scores := board.Score()
	score, ok := scores[engine.PlayerA]
	require.True(t, ok)
	assert.Equal(t, 2, score)

	score, ok = scores[engine.PlayerB]
	require.True(t, ok)
	assert.Equal(t, 6, score)
}

func TestBoard_Copy(t *testing.T) {
	c := engine.New(9, 9)

	ok := c.Add(engine.PlayerA, engine.Coordinate{Row: 0, Column: 0})
	assert.True(t, ok)

	moves := c.PossibleMoves(engine.PlayerB)

	c2 := c.Copy()
	require.NotNil(t, c2)

	assert.Equal(t, slice2map(moves), slice2map(c2.PossibleMoves(engine.PlayerB)))
}

func slice2map[T comparable](list []T) (result map[T]struct{}) {
	result = make(map[T]struct{})
	for _, t := range list {
		result[t] = struct{}{}
	}
	return
}
