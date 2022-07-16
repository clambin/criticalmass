package bot

import (
	"criticalmass/engine"
	"sort"
)

type Bot interface {
	Choose() (engine.Position, bool)
}

type ExploderBot struct {
	Board  *engine.Board
	Player engine.Player
}

var _ Bot = &ExploderBot{}

func (e ExploderBot) Choose() (pos engine.Position, found bool) {
	return e.getMostCritical()
}

func (e ExploderBot) getMostCritical() (pos engine.Position, found bool) {
	type critical struct {
		position  engine.Position
		remaining int
	}

	var criticals []critical
	for _, p := range e.Board.PossibleMoves(e.Player) {
		criticals = append(criticals, critical{
			position:  p,
			remaining: e.Board.Cells[p.Row][p.Column].Remaining(),
		})

	}

	if len(criticals) == 0 {
		return
	}
	sort.Slice(criticals, func(i, j int) bool { return criticals[i].remaining < criticals[j].remaining })
	return criticals[0].position, true
}
