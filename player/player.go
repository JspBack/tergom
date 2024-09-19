package player

import (
	"tergom/bullet"
	"tergom/level"
	"time"
)

type Position struct {
	X, Y int
}

type Player struct {
	Position      Position
	Level         *level.Level
	Health        int
	Bullets       []*bullet.Bullet
	Ammunition    int
	lastShootTime time.Time
	shootCooldown time.Duration
}

func NewPlayer(lvl *level.Level, x, y int) *Player {
	return &Player{
		Position:      Position{X: x, Y: y},
		Level:         lvl,
		Health:        100,
		Bullets:       make([]*bullet.Bullet, 5),
		Ammunition:    100,
		lastShootTime: time.Now(),
		shootCooldown: 100 * time.Millisecond,
	}
}

func (p *Player) Move(dx, dy int) {
	newX, newY := p.Position.X+dx, p.Position.Y+dy
	if !p.Level.IsWall(newX, newY) {
		p.Position.X = newX
		p.Position.Y = newY
	}
}

func (p *Player) Shoot() {
	now := time.Now()
	if p.Ammunition > 0 && now.Sub(p.lastShootTime) >= p.shootCooldown {
		p.lastShootTime = now

		b := bullet.NewBullet(p.Position.X, p.Position.Y-1, -1, p.Level, true)
		for i, blt := range p.Bullets {
			if blt == nil {
				p.Bullets[i] = b
				p.Ammunition--
				break
			}
		}
	}
}

func (p *Player) TakeDamage(damage int) {
	p.Health -= damage
	if p.Health <= 0 {
		p.Health = 0
	}
}
