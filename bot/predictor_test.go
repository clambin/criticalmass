package bot

import (
	"github.com/clambin/criticalmass/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectBestMove(t *testing.T) {
	testCases := []struct {
		name       string
		candidates []move
		wantMoves  []move
	}{
		{
			name: "based on score",
			candidates: []move{
				{position: engine.Coordinate{Row: 0, Column: 0}, score: 0, remaining: 10},
				{position: engine.Coordinate{Row: 1, Column: 1}, score: 1, remaining: 9},
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 2, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 3, remaining: 7},
			},
			wantMoves: []move{{position: engine.Coordinate{Row: 3, Column: 3}, score: 3, remaining: 7}},
		},
		{
			name: "based on remaining",
			candidates: []move{
				{position: engine.Coordinate{Row: 0, Column: 0}, score: 1, remaining: 10},
				{position: engine.Coordinate{Row: 1, Column: 1}, score: 1, remaining: 9},
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 7},
			},
			wantMoves: []move{{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 7}},
		},
		{
			name: "multiple best",
			candidates: []move{
				{position: engine.Coordinate{Row: 0, Column: 0}, score: 1, remaining: 10},
				{position: engine.Coordinate{Row: 1, Column: 1}, score: 1, remaining: 9},
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 8},
			},
			wantMoves: []move{
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 8},
			},
		},
		{
			name: "all best",
			candidates: []move{
				{position: engine.Coordinate{Row: 0, Column: 0}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 1, Column: 1}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 8},
			},
			wantMoves: []move{
				{position: engine.Coordinate{Row: 0, Column: 0}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 1, Column: 1}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 2, Column: 2}, score: 1, remaining: 8},
				{position: engine.Coordinate{Row: 3, Column: 3}, score: 1, remaining: 8},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, tt.wantMoves, selectBestMove(tt.candidates))
			assert.Contains(t, tt.wantMoves, selectBestMoveAlt(tt.candidates))
		})
	}
}

func BenchmarkSelectBestMove(b *testing.B) {
	var candidates []move
	for r := range 10 {
		for c := range 10 {
			candidates = append(candidates, move{position: engine.Coordinate{Row: r, Column: c}, score: r % 3, remaining: c % 3})
		}
	}

	b.ResetTimer()
	b.Run("selectBestMove", func(b *testing.B) {
		for range b.N {
			_ = selectBestMove(candidates)
		}
	})
	b.Run("selectBestMoveAlt", func(b *testing.B) {
		for range b.N {
			_ = selectBestMoveAlt(candidates)
		}
	})
}
