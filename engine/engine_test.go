package engine_test

import (
	"github.com/clambin/criticalmass/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine_SuddenDeath(t *testing.T) {
	board := engine.New(4, 4)

	moves := map[engine.Player]engine.Coordinate{
		engine.PlayerA: {0, 0},
		engine.PlayerB: {3, 3},
	}

	player := engine.Player(engine.PlayerA)

	for !board.GameOver() {
		board.Add(player, moves[player])
		for len(board.GetCriticalCells()) > 0 {
			board.ProcessCriticalCells()
		}
		player = player.NextPlayer()
	}

	winner, won := board.Winner()
	assert.True(t, won)
	assert.Equal(t, engine.Player(engine.PlayerA), winner)
	assert.Equal(t, 60, board.Sum())
}

func Benchmark_SuddenDeath(b *testing.B) {
	moves := map[engine.Player]engine.Coordinate{
		engine.PlayerA: {0, 0},
		engine.PlayerB: {3, 3},
	}
	player := engine.Player(engine.PlayerA)

	for range b.N {
		board := engine.New(4, 4)
		for !board.GameOver() {
			board.Add(player, moves[player])
			for len(board.GetCriticalCells()) > 0 {
				board.ProcessCriticalCells()
			}
			player = player.NextPlayer()
		}
	}
}
