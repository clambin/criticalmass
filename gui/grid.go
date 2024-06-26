package gui

import (
	"fmt"
	"github.com/clambin/criticalmass/engine"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/colornames"
	"image/color"
)

type Grid struct {
	Board          *engine.Board
	Rects          [][]pixel.Rect
	animating      bool
	animationsLeft int
	cellsToAnimate []engine.Coordinate
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
			g.DrawCell(win, imd, g.Board.Cells[engine.Coordinate{Row: row, Column: column}], g.Rects[row][column], g.isAnimatedCell(engine.Coordinate{Row: row, Column: column}))
		}
	}
	imd.Draw(win)

	g.maybeStopAnimation()
}

func (g *Grid) isAnimatedCell(position engine.Coordinate) bool {
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
	g.cellsToAnimate = g.Board.GetCriticalCells()
	g.animating = len(g.cellsToAnimate) > 0
	if g.animating {
		g.animationsLeft = 7
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
	g.Board.ProcessCriticalCells()
	//fmt.Printf("board now has %d items\n", g.Board.Sum())
	//fmt.Printf("scores: %v\n", g.Board.Score())
}

func (g *Grid) DrawCell(win pixel.Target, imd *imdraw.IMDraw, cell *engine.Cell, rect pixel.Rect, animate bool) {
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

func (g *Grid) drawText(win pixel.Target, pos pixel.Vec, cell *engine.Cell, txt string) {
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

func (g *Grid) FindCell(pos pixel.Vec) (engine.Coordinate, bool) {
	for r, row := range g.Rects {
		for c, cell := range row {
			if cell.Contains(pos) {
				return engine.Coordinate{Row: r, Column: c}, true
			}
		}
	}
	return engine.Coordinate{}, false
}
