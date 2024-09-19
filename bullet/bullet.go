package bullet

import "tergom/level"

type Position struct {
	X, Y int
}

type Bullet struct {
	Position       Position
	Direction      int
	Level          *level.Level
	IsPlayerBullet bool
}

func NewBullet(x, y, direction int, lvl *level.Level, isPlayer bool) *Bullet {
	return &Bullet{
		Position:       Position{X: x, Y: y},
		Direction:      direction, // -1 for up, 1 for down
		Level:          lvl,
		IsPlayerBullet: isPlayer,
	}
}
