package gui

import (
	"criticalmass/engine"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
)

type Grid struct {
	Board          *engine.Board
	Rects          [][]pixel.Rect
	CurrentPlayer  engine.Player
	animating      bool
	animationsLeft int
	cellsToAnimate []engine.Position
}

func NewGrid(board *engine.Board) *Grid {
	return &Grid{
		Board: board,
		Rects: makeRects(board.Rows, board.Columns),
	}
}

func makeRects(rows, columns int) (rects [][]pixel.Rect) {
	cells := make([][]pixel.Rect, rows)
	for r := 0; r < rows; r++ {
		cells[r] = make([]pixel.Rect, columns)
		for c := 0; c < columns; c++ {
			x := float64(c) * pixelsPerCell
			y := float64(r) * pixelsPerCell
			cells[r][c] = pixel.R(x, y, x+pixelsPerCell, y+pixelsPerCell)
		}
	}
	return cells
}

func (g *Grid) Draw(win pixel.Target) {
	g.maybeSetupAnimation()

	imd := imdraw.New(nil)
	for row := 0; row < g.Board.Rows; row++ {
		for column := 0; column < g.Board.Columns; column++ {
			g.DrawCell(win, imd, g.Board.Cells[row][column], g.Rects[row][column], g.isAnimatedCell(engine.Position{Row: row, Column: column}))
		}
	}
	imd.Draw(win)

	g.maybeStopAnimation()
}

func (g *Grid) isAnimatedCell(position engine.Position) bool {
	for _, pos := range g.cellsToAnimate {
		if pos == position {
			return true
		}
	}
	return false
}

func (g *Grid) maybeSetupAnimation() {
	if g.animating {
		return
	}
	g.cellsToAnimate = g.Board.GetCriticals()
	g.animating = len(g.cellsToAnimate) > 0
	if g.animating {
		g.animationsLeft = 5
	} else {
		g.cellsToAnimate = nil
	}
}

func (g *Grid) maybeStopAnimation() {
	if !g.animating {
		return
	}
	if g.animationsLeft > 0 {
		g.animationsLeft--
		return
	}
	g.animating = false
	g.Board.ProcessCriticals()
	fmt.Printf("board now has %d items\n", g.Board.Sum())
}

func (g Grid) DrawCell(win pixel.Target, imd *imdraw.IMDraw, cell *engine.Cell, rect pixel.Rect, animate bool) {
	if cell.Count > 0 {
		p := rect.Max.Sub(rect.Min).Scaled(.5).Add(rect.Min)
		g.drawText(win, p, cell, fmt.Sprintf("%d", cell.Count))
	}

	thickness := 5.0
	if animate {
		thickness *= float64(g.animationsLeft)
	}

	imd.Color = colornames.Darkgrey
	imd.Push(rect.Min, rect.Max)
	imd.Rectangle(thickness)
}

func (g Grid) drawText(win pixel.Target, pos pixel.Vec, cell *engine.Cell, txt string) {
	t := text.New(pos, text.Atlas7x13)
	var textColor color.RGBA
	if cell.Owner == engine.PlayerA {
		textColor = colornames.Aqua
	} else {
		textColor = colornames.Yellow
	}
	t.Color = textColor
	_, _ = fmt.Fprint(t, txt)
	t.Draw(win, pixel.IM)
}

func (g Grid) FindCell(pos pixel.Vec) (engine.Position, bool) {
	for r, row := range g.Rects {
		for c, cell := range row {
			if cell.Contains(pos) {
				return engine.Position{Row: r, Column: c}, true
			}
		}
	}
	return engine.Position{}, false
}

func (g *Grid) DoMove(pos engine.Position) {
	if g.Board.Add(g.CurrentPlayer, pos) {
		g.CurrentPlayer = g.CurrentPlayer.NextPlayer()
	}
}
