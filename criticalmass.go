package main

import (
	"criticalmass/gui"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	ui := gui.NewUI(3, 9)
	pixelgl.Run(ui.Run)
}
