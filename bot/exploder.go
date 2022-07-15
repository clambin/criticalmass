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

	opponent := e.Player.NextPlayer()

	var criticals []critical
	for r, row := range e.Board.Cells {
		for c, cell := range row {
			if cell.Count == 0 || cell.Owner != opponent {
				criticals = append(criticals, critical{
					position:  engine.Position{Row: r, Column: c},
					remaining: cell.Critical - cell.Count,
				})
			}
		}
	}
	if len(criticals) == 0 {
		return
	}
	sort.Slice(criticals, func(i, j int) bool { return criticals[i].remaining < criticals[j].remaining })
	return criticals[0].position, true
}
