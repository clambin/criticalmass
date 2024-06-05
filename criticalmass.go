package main

import (
	"github.com/clambin/criticalmass/gui"
	pixelgl "github.com/gopxl/pixel/v2/backends/opengl"
)

func main() {
	ui := gui.NewUI(4, 8)
	pixelgl.Run(ui.Run)
}
