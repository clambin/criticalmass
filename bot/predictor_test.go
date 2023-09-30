package bot

import (
	"github.com/clambin/criticalmass/engine"
	"github.com/clambin/go-common/set"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectBestMove(t *testing.T) {
	testCases := []struct {
		name       string
		candidates []move
		wantMoves  set.Set[move]
	}{
		{
			name: "score",
			candidates: []move{
				{position: engine.Coordinate{Row: 0, Column: 0}, score: 0, remaining: 10},
				{position: engine.Coordinate{Row: 1, Column: 1}, score: 1, remaining: 9},
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 2, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 3, remaining: 7},
			},
			wantMoves: set.Create(move{position: engine.Coordinate{Row: 3, Column: 3}, score: 3, remaining: 7}),
		},
		{
			name: "remaining",
			candidates: []move{
				{position: engine.Coordinate{Row: 0, Column: 0}, score: 1, remaining: 10},
				{position: engine.Coordinate{Row: 1, Column: 1}, score: 1, remaining: 9},
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 7},
			},
			wantMoves: set.Create(move{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 7}),
		},
		{
			name: "multiple candidates",
			candidates: []move{
				{position: engine.Coordinate{Row: 0, Column: 0}, score: 1, remaining: 10},
				{position: engine.Coordinate{Row: 1, Column: 1}, score: 1, remaining: 9},
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 8},
			},
			wantMoves: set.Create(
				move{position: engine.Coordinate{Row: 2, Column: 2}, score: 1, remaining: 8},
				move{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 8},
			),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			selected := selectBestMove(tt.candidates)
			assert.True(t, tt.wantMoves.Contains(selected))

			selected = selectBestMoveAlt(tt.candidates)
			assert.True(t, tt.wantMoves.Contains(selected))
		})
	}
}

func BenchmarkSelectBestMove(b *testing.B) {
	var candidates []move
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			candidates = append(candidates, move{position: engine.Coordinate{Row: r, Column: c}, score: r % 3, remaining: c % 3})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = selectBestMove(candidates)
	}
}

func BenchmarkSelectBestMoveAlt(b *testing.B) {
	var candidates []move
	for r := 0; r < 10; r++ {
		for c := 0; c < 10; c++ {
			candidates = append(candidates, move{position: engine.Coordinate{Row: r, Column: c}, score: r % 3, remaining: c % 3})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = selectBestMoveAlt(candidates)
	}
}
