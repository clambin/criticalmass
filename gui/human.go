package gui

import (
	"github.com/clambin/criticalmass/engine"
	"github.com/gopxl/pixel/v2"
	pixelgl "github.com/gopxl/pixel/v2/backends/opengl"
)

var _ Player = &HumanPlayer{}

type HumanPlayer struct {
	win  *pixelgl.Window
	grid *Grid
}

func (h HumanPlayer) Choose() (pos engine.Coordinate, found bool) {
	if h.win.JustPressed(pixel.MouseButtonLeft) {
		pos, found = h.grid.FindCell(h.win.MousePosition())
	}
	return pos, found
}
