package enemy

import (
	"math/rand"
	"tergom/bullet"
	"tergom/level"
	"time"
)

type Position struct {
	X int
	Y int
}

type Enemy struct {
	Position       Position
	Health         int
	LastShoot      time.Time
	ShootDelay     time.Duration
	Level          *level.Level
	MoveDirection  int // 1 for right, -1 for left
	MoveInterval   time.Duration
	LastMove       time.Time
	MoveCooldown   time.Duration
	SpeedIncrement float64
}

func NewEnemy(x, y int, lvl *level.Level, speedIncrement float64) *Enemy {
	direction := rand.Intn(2)
	if direction == 0 {
		direction = -1
	} else {
		direction = 1
	}
	return &Enemy{
		Position:       Position{X: x, Y: y},
		Health:         10, // for the basic enemies
		LastShoot:      time.Now(),
		ShootDelay:     2 * time.Second,
		Level:          lvl,
		MoveDirection:  direction,
		MoveInterval:   500 * time.Millisecond, // Initial move interval: 0.5s
		LastMove:       time.Now(),
		MoveCooldown:   500 * time.Millisecond,
		SpeedIncrement: speedIncrement,
	}
}

func (e *Enemy) Move(dx, dy int) {
	e.Position.X += dx
	e.Position.Y += dy
}

func (e *Enemy) CanShoot() bool {
	return time.Since(e.LastShoot) >= e.ShootDelay
}

func (e *Enemy) Shoot() *bullet.Bullet {
	e.LastShoot = time.Now()
	return bullet.NewBullet(e.Position.X, e.Position.Y+1, 1, e.Level, false)
}

func (e *Enemy) TakeDamage(damage int) {
	e.Health -= damage
}

func (e *Enemy) UpdateMovement() {
	if time.Since(e.LastMove) >= e.MoveCooldown {
		e.Move(e.MoveDirection, 0)

		if e.Position.X <= 1 || e.Position.X >= e.Level.Width-2 {
			e.MoveDirection *= -1
			e.Move(0, 1)
		}
		e.LastMove = time.Now()
	}
}

func (e *Enemy) IncreaseSpeed() {
	if e.MoveCooldown > 100*time.Millisecond {
		e.MoveCooldown = time.Duration(float64(e.MoveCooldown) * e.SpeedIncrement)
	}
}
