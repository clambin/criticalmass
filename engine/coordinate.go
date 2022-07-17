package engine

type Coordinate struct {
	Row, Column int
}

func (p Coordinate) Add(p2 Coordinate) Coordinate {
	return Coordinate{
		Row:    p.Row + p2.Row,
		Column: p.Column + p2.Column,
	}
}
