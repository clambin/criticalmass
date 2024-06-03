package gui

import (
	"github.com/clambin/criticalmass/engine"
	"github.com/faiface/pixel/pixelgl"
)

var _ Player = &HumanPlayer{}

type HumanPlayer struct {
	win  *pixelgl.Window
	grid *Grid
}

func (h HumanPlayer) Choose() (pos engine.Coordinate, found bool) {
	if h.win.JustPressed(pixelgl.MouseButtonLeft) {
		pos, found = h.grid.FindCell(h.win.MousePosition())
	}
	return pos, found
}
