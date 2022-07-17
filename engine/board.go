package engine

import "fmt"

// A Board holds all cells and implements the game logic
type Board struct {
	Rows    int
	Columns int
	Cells   map[Coordinate]*Cell
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
	cells := make(map[Coordinate]*Cell)
	for _, pos := range b.Coordinates() {
		cells[pos] = &Cell{Critical: 1 + len(b.neighbours(pos))}
	}
	b.Cells = cells
}

func (b *Board) neighbours(pos Coordinate) (neighbours []Coordinate) {
	candidates := []Coordinate{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	neighbours = make([]Coordinate, 0, len(candidates))
	for _, candidate := range candidates {
		p := Coordinate{Row: pos.Row, Column: pos.Column}.Add(candidate)
		if b.isValid(p) {
			neighbours = append(neighbours, p)
		}
	}
	return
}

func (b *Board) isValid(p Coordinate) bool {
	return p.Row >= 0 && p.Row < b.Rows &&
		p.Column >= 0 && p.Column < b.Columns
}

func (b *Board) Add(player Player, position Coordinate) bool {
	c, found := b.Cells[position]
	if !found {
		return false
	}
	return c.Add(player, 1, false)
}

// GetCriticals returns the position of all cells that have reached critical mass and the cells that will get
// a remainder of their mass after explosion
func (b *Board) GetCriticals() (criticals []Coordinate) {
	for pos, cell := range b.Cells {
		if cell.IsCritical() {
			criticals = append(criticals, pos)
			criticals = append(criticals, b.neighbours(pos)...)
		}
	}
	return
}

// ProcessCriticals explodes any cell that has achieved critical mass
func (b *Board) ProcessCriticals() {
	for _, pos := range b.Coordinates() {
		cell := b.Cells[pos]
		if cell.IsCritical() {
			ns := b.neighbours(pos)
			cell.Add(cell.Owner, -len(ns), true)
			for _, n := range ns {
				b.Cells[n].Add(cell.Owner, 1, true)
			}
		}
	}
}

// Winner determines if someone's won the game. Winning means a player has occupied all cells on the board.
func (b *Board) Winner() (winner Player, won bool) {
	counts := make([]int, PlayerB+1)
	for _, cell := range b.Cells {
		if cell.Count == 0 {
			return
		}
		counts[cell.Owner]++
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
func (b *Board) GameOver() bool {
	_, won := b.Winner()
	return won
}

func (b *Board) Dump() (dump string) {
	currentRow := 0
	for _, pos := range b.Coordinates() {
		if pos.Row != currentRow {
			dump += "\n"
			currentRow = pos.Row
		}

		cell := b.Cells[pos]
		var name string
		if cell.Count == 0 {
			name = " "
		} else {
			name = cell.Owner.String()
		}
		dump += fmt.Sprintf("%s%d ", name, cell.Count)

	}
	dump += "\n"
	return
}

func (b *Board) Sum() (sum int) {
	for _, cell := range b.Cells {
		sum += cell.Count
	}
	return
}

func (b *Board) Score() (scores map[Player]int) {
	scores = make(map[Player]int)
	for _, player := range Players() {
		scores[player] = 0
	}

	for _, cell := range b.Cells {
		if cell.Count == 0 {
			continue
		}
		current := scores[cell.Owner]
		current += cell.Count
		scores[cell.Owner] = current
	}
	return
}

func (b *Board) PossibleMoves(player Player) (moves []Coordinate) {
	for pos, cell := range b.Cells {
		if cell.Count == 0 || cell.Owner == player {
			moves = append(moves, pos)
		}
	}
	return
}

func (b *Board) Copy() *Board {
	cells := make(map[Coordinate]*Cell, b.Rows)
	for pos, cell := range b.Cells {
		cells[pos] = &Cell{
			Owner:    cell.Owner,
			Count:    cell.Count,
			Critical: cell.Critical,
		}
	}

	return &Board{
		Rows:    b.Rows,
		Columns: b.Columns,
		Cells:   cells,
	}
}

func (b *Board) Coordinates() (result []Coordinate) {
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Columns; c++ {
			result = append(result, Coordinate{Row: r, Column: c})
		}
	}
	return
}
