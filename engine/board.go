package engine

import "fmt"

// A Board holds all cells and implements the game logic
type Board struct {
	Rows    int
	Columns int
	Cells   map[Coordinate]*Cell
}

func New(rows, columns int) *Board {
	board := Board{
		Rows:    rows,
		Columns: columns,
	}
	coordinates := board.Coordinates()
	cells := make(map[Coordinate]*Cell, len(coordinates))
	for _, pos := range coordinates {
		cells[pos] = &Cell{Critical: 1 + len(board.neighbours(pos))}
	}
	board.Cells = cells
	return &board
}

var (
	candidates = []Coordinate{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
)

func (b *Board) neighbours(pos Coordinate) []Coordinate {
	neighbours := make([]Coordinate, 0, len(candidates))
	for _, candidate := range candidates {
		p := Coordinate{Row: pos.Row, Column: pos.Column}.Add(candidate)
		if b.isValid(p) {
			neighbours = append(neighbours, p)
		}
	}
	return neighbours
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

// GetCriticalCells returns the position of all cells that have reached critical mass and the cells that will get
// a remainder of their mass after explosion
func (b *Board) GetCriticalCells() (criticals []Coordinate) {
	for pos, cell := range b.Cells {
		if cell.IsCritical() {
			criticals = append(criticals, pos)
			criticals = append(criticals, b.neighbours(pos)...)
		}
	}
	return
}

// ProcessCriticalCells explodes any cell that has achieved critical mass
func (b *Board) ProcessCriticalCells() {
	for _, pos := range b.Coordinates() {
		if cell := b.Cells[pos]; cell.IsCritical() {
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
	scores := b.Score()
	if scores[PlayerA] > 0 && scores[PlayerB] == 0 {
		return PlayerA, true
	}
	if scores[PlayerA] == 0 && scores[PlayerB] > 0 {
		return PlayerB, true
	}
	return
}

// GameOver determines if the game is over
func (b *Board) GameOver() bool {
	for _, c := range b.Cells {
		if c.Count == 0 {
			return false
		}
	}
	if len(b.GetCriticalCells()) > 0 {
		return false
	}
	scores := b.Score()
	return scores[PlayerA] == 0 || scores[PlayerB] == 0
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

func (b *Board) Score() (scores map[Player]int) {
	scores = map[Player]int{
		PlayerA: 0,
		PlayerB: 0,
	}
	for _, cell := range b.Cells {
		scores[cell.Owner] += cell.Count
	}
	return scores
}
