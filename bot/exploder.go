package bot

import (
	"cmp"
	"github.com/clambin/criticalmass/engine"
	"slices"
)

type ExploderBot struct {
	Board  *engine.Board
	Player engine.Player
}

func (e ExploderBot) Choose() (engine.Coordinate, bool) {
	return e.getMostCritical()
}

func (e ExploderBot) getMostCritical() (pos engine.Coordinate, found bool) {
	type critical struct {
		position  engine.Coordinate
		remaining int
	}

	var criticals []critical
	for _, p := range e.Board.PossibleMoves(e.Player) {
		criticals = append(criticals, critical{
			position:  p,
			remaining: e.Board.Cells[p].Remaining(),
		})

	}

	if len(criticals) == 0 {
		return engine.Coordinate{}, false
	}
	slices.SortFunc(criticals, func(a, b critical) int {
		return cmp.Compare(a.remaining, b.remaining)
	})

	return criticals[0].position, true
}
