package gui

import (
	"criticalmass/engine"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)

type UI struct {
	X    float64
	Y    float64
	Name string
	Grid *Grid
}

const border = 20
const pixelsPerCell = 100

func NewUI(rows, cols int) (ui *UI) {
	b := engine.New(rows, cols)
	return &UI{
		X:    float64(cols * pixelsPerCell),
		Y:    float64(rows * pixelsPerCell),
		Name: "critical mass",
		Grid: NewGrid(b),
	}
}

func (ui *UI) Run() {
	cfg := pixelgl.WindowConfig{
		Title:  ui.Name,
		Bounds: pixel.R(-border, -border, ui.X+border, ui.Y+border),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	timestamp := time.Now()
	for !win.Closed() {
		win.Clear(colornames.Black)
		ui.Grid.Draw(win)
		win.Update()
		ui.ProcessEvents(win)
		win.SetTitle(fmt.Sprintf("%s (%.1f FPS)", ui.Name, 1.0/time.Since(timestamp).Seconds()))
		timestamp = time.Now()
		<-ticker.C
	}
	ticker.Stop()
}

func (ui *UI) ProcessEvents(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		pos := win.MousePosition()
		if p, found := ui.Grid.FindCell(pos); found {
			fmt.Printf("clicked %v. matches cell (%d,%d)\n", pos, p.Row, p.Column)
			ui.Grid.DoMove(p)
		}
	}
}