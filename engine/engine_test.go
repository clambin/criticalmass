package engine_test

import (
	"criticalmass/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEngine(t *testing.T) {
	board := engine.New(3, 3)

	for i := 0; i < 3; i++ {
		board.Add(engine.PlayerA, engine.Position{Row: 0, Column: 0})
		board.Add(engine.PlayerB, engine.Position{Row: 0, Column: 2})
	}
	_, won := board.Winner()
	assert.False(t, won)

	assert.Equal(t, `A3  0 B3 
 0  0  0 
 0  0  0 
`, board.Dump())

	board.Add(engine.PlayerA, engine.Position{Row: 0, Column: 0})

	c := board.GetCriticals()
	require.Len(t, c, 4)
	assert.Equal(t, engine.Position{Row: 0, Column: 0}, c[0])

	board.ProcessCriticals()

	assert.Equal(t, `A1 A1 B3 
A1 A1  0 
 0  0  0 
`, board.Dump())

	board.Add(engine.PlayerB, engine.Position{Row: 0, Column: 2})
	board.ProcessCriticals()

	assert.Equal(t, `A1 B2 B1 
A1 B2 B1 
 0  0  0 
`, board.Dump())

}

func TestEngine_SuddenDeath(t *testing.T) {
	board := engine.New(4, 4)

	moves := map[engine.Player]engine.Position{
		engine.PlayerA: {0, 0},
		engine.PlayerB: {3, 3},
	}

	player := engine.Player(engine.PlayerA)

	//index := 0
	for !board.GameOver() {
		board.Add(player, moves[player])
		board.ProcessCriticals()

		//fmt.Println(board.Dump())
		//fmt.Printf("==================== - %d\n", index)
		//index++

		player = player.NextPlayer()
	}

	winner, won := board.Winner()
	assert.True(t, won)
	assert.Equal(t, engine.Player(engine.PlayerA), winner)
	assert.Equal(t, 63, board.Sum())
}
