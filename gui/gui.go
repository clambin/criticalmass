package gui

import (
	"fmt"
	"github.com/clambin/criticalmass/bot"
	"github.com/clambin/criticalmass/engine"
	"github.com/gopxl/pixel/v2"
	pixelgl "github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
	"slices"
	"strings"
	"time"
)

type Player interface {
	Choose() (engine.Coordinate, bool)
}

type UI struct {
	X             float64
	Y             float64
	Name          string
	Grid          *Grid
	PlayerA       Player
	PlayerB       Player
	CurrentPlayer engine.Player
}

const border = 20
const pixelsPerCell = 100

func NewUI(rows, cols int) (ui *UI) {
	b := engine.New(rows, cols)
	return &UI{
		X:       float64(cols * pixelsPerCell),
		Y:       float64(rows * pixelsPerCell),
		Name:    "critical mass",
		Grid:    NewGrid(b),
		PlayerA: &HumanPlayer{}, //&bot.Predictor{Board: b, Player: engine.PlayerA},
		PlayerB: &bot.ExploderBot{Board: b, Player: engine.PlayerB},
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

	if humanPlayer, ok := ui.PlayerA.(*HumanPlayer); ok {
		humanPlayer.win = win
		humanPlayer.grid = ui.Grid
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timestamp := time.Now()

	for !ui.Grid.Board.GameOver() {
		win.Clear(colornames.Black)
		ui.Grid.Draw(win)
		win.Update()
		if !ui.Grid.animating {
			ui.ProcessEvents(win)
		}
		win.SetTitle(ui.title(timestamp))
		timestamp = time.Now()
		<-ticker.C

	}

	winner, _ := ui.Grid.Board.Winner()
	fmt.Printf("player %s won! %s", winner.String(), strings.Join(ui.scores(), ", "))
}

func (ui *UI) title(timestamp time.Time) string {
	return fmt.Sprintf(
		"%s (%.1f FPS) - %s",
		ui.Name,
		1.0/time.Since(timestamp).Seconds(),
		strings.Join(ui.scores(), ", "),
	)
}

func (ui *UI) scores() []string {
	scores := ui.Grid.Board.Score()
	scoreStrings := make([]string, 0, len(scores))
	for player, score := range scores {
		scoreStrings = append(scoreStrings, fmt.Sprintf("%s:%d", player, score))
	}
	slices.Sort(scoreStrings)
	return scoreStrings
}

func (ui *UI) ProcessEvents(_ *pixelgl.Window) {
	var player *Player
	switch ui.CurrentPlayer {
	case engine.PlayerA:
		player = &ui.PlayerA
	case engine.PlayerB:
		player = &ui.PlayerB
	default:
		return
	}
	if pos, found := (*player).Choose(); found {
		if ui.Grid.Board.Add(ui.CurrentPlayer, pos) {
			ui.CurrentPlayer = ui.CurrentPlayer.NextPlayer()
		}
	}
}
