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
	moves := e.getPossibleMoves()

	if len(moves) == 0 {
		return engine.Coordinate{}, false
	}

	return selectBestMove(moves).position, true
}

func (e Predictor) getPossibleMoves() []move {
	var moves []move
	for _, p := range e.Board.PossibleMoves(e.Player) {
		b := e.Board.Copy()
		b.Add(e.Player, p)
		s := b.Score()
		moves = append(moves, move{
			position:  p,
			score:     s[e.Player] - s[e.Player.NextPlayer()],
			remaining: e.Board.Cells[p].Remaining(),
		})
	}
	return moves
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
	candidates = candidates[:equalCount+1]
	return candidates[rand.Intn(len(candidates))]
}

func selectBestMoveAlt(candidates []move) move {
	bestMoves := make([]move, 0, len(candidates))

	for _, candidate := range candidates {
		if len(bestMoves) == 0 {
			bestMoves = append(bestMoves, candidate)
			continue
		}

		// how does the candidate compare to our current best?
		order := candidate.Compare(bestMoves[0])

		// worse. discard it.
		if order == 1 {
			continue
		}

		// better. discard current best moves (keeping the slice)
		if order == -1 {
			bestMoves = bestMoves[0:0]
		}

		// add the candidate to the list of best moves
		bestMoves = append(bestMoves, candidate)
	}

	// only one best move.
	if len(bestMoves) == 1 {
		return bestMoves[0]
	}

	// multiple best moves. pick a random one.
	return bestMoves[rand.Intn(len(bestMoves))]
}
