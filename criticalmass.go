package main

import (
	"github.com/clambin/criticalmass/gui"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	ui := gui.NewUI(4, 8)
	pixelgl.Run(ui.Run)
}
