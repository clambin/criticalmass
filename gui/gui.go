package gui

import (
	"fmt"
	"github.com/clambin/criticalmass/bot"
	"github.com/clambin/criticalmass/engine"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)

type Bot interface {
	Choose() (engine.Coordinate, bool)
}

type UI struct {
	X    float64
	Y    float64
	Name string
	Grid *Grid
	BotA Bot
	BotB Bot
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
		BotA: &bot.Predictor{Board: b, Player: engine.PlayerA},
		BotB: &bot.Predictor{Board: b, Player: engine.PlayerB},
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
		if !ui.Grid.animating {
			ui.ProcessEvents(win)
		}
		win.SetTitle(fmt.Sprintf("%s (%.1f FPS)", ui.Name, 1.0/time.Since(timestamp).Seconds()))
		timestamp = time.Now()
		<-ticker.C
	}
	ticker.Stop()
}

func (ui *UI) ProcessEvents(win *pixelgl.Window) {
	switch ui.Grid.CurrentPlayer {
	case engine.PlayerA:
		var (
			pos   engine.Coordinate
			found bool
		)
		if ui.BotA != nil {
			pos, found = ui.BotA.Choose()
		} else if win.JustPressed(pixelgl.MouseButtonLeft) {
			pos, found = ui.Grid.FindCell(win.MousePosition())
		}
		if found {
			ui.Grid.DoMove(pos)
		}
	case engine.PlayerB:
		if pos, found := ui.BotB.Choose(); found {
			ui.Grid.DoMove(pos)
		}
	}
}
