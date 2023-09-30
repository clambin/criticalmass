package bot

import (
	"cmp"
	"github.com/clambin/criticalmass/engine"
	"math/rand"
	"slices"
)

type Predictor struct {
	Board  *engine.Board
	Player engine.Player
}

var _ Bot = &Predictor{}

type move struct {
	position  engine.Coordinate
	score     int
	remaining int
}

func (m move) Compare(other move) int {
	if order := cmp.Compare(m.score, other.score); order != 0 {
		return -order
	}
	return cmp.Compare(m.remaining, other.remaining)
}

func (e Predictor) Choose() (engine.Coordinate, bool) {
	var scores []move

	for _, p := range e.Board.PossibleMoves(e.Player) {
		b := e.Board.Copy()
		b.Add(e.Player, p)
		s := b.Score()
		scores = append(scores, move{
			position:  p,
			score:     s[e.Player] - s[e.Player.NextPlayer()],
			remaining: e.Board.Cells[p].Remaining(),
		})
	}
	var pos engine.Coordinate
	var found bool
	if len(scores) != 0 {
		found = true
		pos = selectBestMove(scores).position
	}

	return pos, found
}

func selectBestMove(candidates []move) move {
	slices.SortFunc(candidates, func(a, b move) int {
		return a.Compare(b)
	})

	var equalCount int
	for _, candidate := range candidates[1:] {
		if candidate.Compare(candidates[0]) != 0 {
			break
		}
		equalCount++
	}

	if equalCount == 0 {
		return candidates[0]
	}

	return candidates[rand.Intn(equalCount+1)]
}

func selectBestMoveAlt(candidates []move) move {
	var bestMoves []move

	for _, candidate := range candidates {
		if len(bestMoves) == 0 {
			bestMoves = append(bestMoves, candidate)
			continue
		}

		order := candidate.Compare(bestMoves[0])
		if order == 1 {
			continue
		}

		if order == -1 {
			bestMoves = bestMoves[0:0]
		}

		bestMoves = append(bestMoves, candidate)
	}

	if len(bestMoves) == 1 {
		return bestMoves[0]
	}

	return bestMoves[rand.Intn(len(bestMoves))]
}
