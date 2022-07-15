package engine

// A Player of the game. Currently, we support only two players.
type Player int

const (
	PlayerA = iota
	PlayerB
)

func Players() []Player {
	return []Player{PlayerA, PlayerB}
}

func (p Player) NextPlayer() Player {
	return Player((int(p) + 1) % len(Players()))
}

func (p Player) String() string {
	return []string{"A", "B"}[p]
}
