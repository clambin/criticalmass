package bot

import (
	"github.com/clambin/criticalmass/engine"
	"sort"
)

type ExploderBot struct {
	Board  *engine.Board
	Player engine.Player
}

func (e ExploderBot) Choose() (pos engine.Coordinate, found bool) {
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
		return
	}
	sort.Slice(criticals, func(i, j int) bool { return criticals[i].remaining < criticals[j].remaining })
	return criticals[0].position, true
}
