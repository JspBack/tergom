package enemy

import (
	"math/rand"
	"tergom/level"
	"time"
)

type Wave struct {
	Enemies        []*Enemy
	Number         int
	SpeedIncrement float64
}

func NewWave(number int, lvl *level.Level, speedIncrement float64) *Wave {
	wave := &Wave{
		Enemies:        make([]*Enemy, 0),
		Number:         number,
		SpeedIncrement: speedIncrement,
	}
	wave.generateEnemies(lvl)
	return wave
}

func (w *Wave) generateEnemies(lvl *level.Level) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	columns := lvl.Width - 2
	rows := 3

	for i := 0; i < w.Number; i++ {
		x := rng.Intn(columns) + 1
		y := rng.Intn(rows) + 1
		enemy := NewEnemy(x, y, lvl, w.SpeedIncrement)
		w.Enemies = append(w.Enemies, enemy)
	}
}

func (w *Wave) IsCleared() bool {
	for _, enemy := range w.Enemies {
		if enemy.Health > 0 {
			return false
		}
	}
	return true
}

func (w *Wave) UpdateEnemies() {
	for _, enemy := range w.Enemies {
		if enemy.Health > 0 {
			enemy.UpdateMovement()
		}
	}
}

func (w *Wave) IncreaseEnemiesSpeed() {
	for _, enemy := range w.Enemies {
		enemy.IncreaseSpeed()
	}
}
