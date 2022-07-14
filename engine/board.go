package engine

import "fmt"

// A Board holds all cells and implements the game logic
type Board struct {
	Rows    int
	Columns int
	Cells   [][]*Cell
}

func New(rows, columns int) (board *Board) {
	board = &Board{
		Rows:    rows,
		Columns: columns,
	}
	board.makeCells()
	return
}

func (b *Board) makeCells() {
	cells := make([][]*Cell, b.Rows)
	for r := 0; r < b.Rows; r++ {
		cells[r] = make([]*Cell, b.Columns)
		for c := 0; c < b.Columns; c++ {
			cells[r][c] = &Cell{Critical: 1 + len(b.neighbours(r, c))}
		}
	}
	b.Cells = cells
}

type Position struct {
	Row, Column int
}

func (p Position) Add(p2 Position) Position {
	return Position{
		Row:    p.Row + p2.Row,
		Column: p.Column + p2.Column,
	}
}

func (b Board) neighbours(row, col int) (neighbours []Position) {
	candidates := []Position{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, candidate := range candidates {
		p := Position{Row: row, Column: col}.Add(candidate)
		if b.isValid(p) {
			neighbours = append(neighbours, p)
		}
	}
	return
}

func (b Board) isValid(p Position) bool {
	return p.Row >= 0 && p.Row < b.Rows &&
		p.Column >= 0 && p.Column < b.Columns
}

func (b *Board) Add(player Player, position Position) bool {
	if !b.isValid(position) {
		return false
	}

	return b.Cells[position.Row][position.Column].Add(player, 1, false)
}

// GetCriticals returns the position of all cells that have reached critical mass and the cells that will get
// a remainder of their mass after explosion
func (b Board) GetCriticals() (criticals []Position) {
	for r, row := range b.Cells {
		for c, cell := range row {
			if cell.IsCritical() {
				p := Position{Row: r, Column: c}
				criticals = append(criticals, p)
				criticals = append(criticals, b.neighbours(r, c)...)
			}
		}
	}
	return
}

// ProcessCriticals explodes any cell that has achieved critical mass
func (b *Board) ProcessCriticals() {
	for r, row := range b.Cells {
		for c, cell := range row {
			if cell.IsCritical() {
				ns := b.neighbours(r, c)
				cell.Add(cell.Owner, -len(ns) /*-1*/, true)
				for _, n := range ns {
					b.Cells[n.Row][n.Column].Add(cell.Owner, 1, true)
				}
			}
		}
	}
}

// Winner determines if someone's won the game. Winning means a player has occupied all cells on the board.
func (b Board) Winner() (winner Player, won bool) {
	counts := make([]int, PlayerB+1)
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.Count == 0 {
				return
			}
			counts[cell.Owner]++
		}
	}
	won = true
	if counts[PlayerA] == 0 {
		winner = PlayerB
	} else {
		winner = PlayerA
	}
	return
}

// GameOver determines if the game is over
func (b Board) GameOver() bool {
	_, won := b.Winner()
	return won
}

func (b Board) Dump() (dump string) {
	ownerNames := []string{"A", "B"}
	for _, row := range b.Cells {
		for _, cell := range row {
			var name string
			if cell.Count == 0 {
				name = " "
			} else {
				name = ownerNames[cell.Owner]
			}
			dump += fmt.Sprintf("%s%d ", name, cell.Count)
		}
		dump += "\n"
	}
	return
}

func (b Board) Sum() (sum int) {
	for _, row := range b.Cells {
		for _, cell := range row {
			sum += cell.Count
		}
	}
	return
}
