package bot

import (
	"github.com/clambin/criticalmass/engine"
	"math/rand"
	"sort"
)

type Predictor struct {
	Board  *engine.Board
	Player engine.Player
}

var _ Bot = &Predictor{}

func (e Predictor) Choose() (pos engine.Coordinate, found bool) {
	type score struct {
		position  engine.Coordinate
		score     int
		remaining int
	}
	var scores []score

	for _, p := range e.Board.PossibleMoves(e.Player) {
		b := e.Board.Copy()
		b.Add(e.Player, p)
		s := b.Score()
		scores = append(scores, score{
			position:  p,
			score:     s[e.Player] - s[e.Player.NextPlayer()],
			remaining: e.Board.Cells[p].Remaining(),
		})
	}
	if len(scores) == 0 {
		return pos, false
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score >= scores[j].score && scores[i].remaining <= scores[j].remaining
	})

	var highestScores []score
	for _, s := range scores {
		if s.score == scores[0].score && s.remaining == scores[0].remaining {
			highestScores = append(highestScores, s)
		}
	}

	return highestScores[rand.Intn(len(highestScores))].position, true
}
