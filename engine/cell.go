package engine

// A Cell represents one field on the game board
type Cell struct {
	Owner    Player
	Count    int
	Critical int
}

// Add increases the count of the provided field. If the cell is not taken, and the player doesn't own the cell,
// then the move isn't made. Returns false if the move is illegal (i.e. the cell is owned by another player.
func (c *Cell) Add(player Player, count int, critical bool) bool {
	if !critical && c.Count > 0 && c.Owner != player {
		return false
	}
	c.Count += count
	c.Owner = player
	return true
}

// IsCritical returns true if the cell has achieved critical mass
func (c *Cell) IsCritical() bool {
	return c.Count >= c.Critical
}
