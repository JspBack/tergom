package game

import (
	"bytes"
	"fmt"
	"os"
	"tergom/bullet"
	"tergom/enemy"
	"tergom/level"
	"tergom/player"
	"tergom/renderer"
	"tergom/stats"
	"tergom/utils"
	"time"

	"github.com/eiannone/keyboard"
)

type Game struct {
	isRunning      bool
	level          *level.Level
	drawBuf        *bytes.Buffer
	stats          *stats.Stats
	player         *player.Player
	inputChan      chan keyboard.Key
	closeChan      chan struct{}
	currentWave    *enemy.Wave
	waveNumber     int
	enemyBullets   []*bullet.Bullet
	gameOver       bool
	gameOverMsg    string
	transitionBuf  *bytes.Buffer
	speedIncrement float64
}

func NewGame(width, height int) *Game {
	lvl := level.NewLevel(width, height)
	return &Game{
		level:          lvl,
		drawBuf:        new(bytes.Buffer),
		stats:          stats.NewStats(),
		player:         player.NewPlayer(lvl, width/2, height-2),
		inputChan:      make(chan keyboard.Key),
		closeChan:      make(chan struct{}),
		waveNumber:     1,
		enemyBullets:   make([]*bullet.Bullet, 0),
		transitionBuf:  new(bytes.Buffer),
		speedIncrement: 0.95, // Enemies move 5% faster each wave
	}
}

func (g *Game) Start() {
	g.isRunning = true
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	fmt.Print("\033[?25l")       // Hide cursor
	defer fmt.Print("\033[?25h") // Show cursor

	go g.listenInput()

	g.stats.InitializeTimer(60)
	g.spawnWave()

	g.loop()
}

func (g *Game) spawnWave() {
	g.currentWave = enemy.NewWave(g.waveNumber*5, g.level, g.speedIncrement)
}

func (g *Game) listenInput() {
	for g.isRunning {
		_, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}
		g.inputChan <- key
	}
}

func (g *Game) loop() {
	ticker := time.NewTicker(time.Millisecond * 33) // ~30fps
	defer ticker.Stop()

	for g.isRunning {
		select {
		case key := <-g.inputChan:
			g.handleKey(key)
		case <-ticker.C:
			g.update()
			g.render()
			g.stats.Update()
			if g.checkGameOver() {
				g.isRunning = false
			}
		case <-g.closeChan:
			return
		}
	}

	g.showGameOver()
}

func (g *Game) handleKey(key keyboard.Key) {
	if g.gameOver {
		return
	}

	switch key {
	case keyboard.KeyEsc:
		g.isRunning = false
	case keyboard.KeyArrowUp:
		g.player.Move(0, -1)
	case keyboard.KeyArrowDown:
		g.player.Move(0, 1)
	case keyboard.KeyArrowLeft:
		g.player.Move(-1, 0)
	case keyboard.KeyArrowRight:
		g.player.Move(1, 0)
	case keyboard.KeySpace:
		g.player.Shoot()
	}
}

func (g *Game) update() {
	if g.gameOver {
		return
	}

	// Update player bullets
	for i, b := range g.player.Bullets {
		if b == nil {
			continue
		}
		b.Position.Y += b.Direction

		if b.Position.Y < 0 || b.Position.Y >= g.level.Height || g.level.IsWall(b.Position.X, b.Position.Y) {
			g.player.Bullets[i] = nil
			continue
		}

		for _, enemy := range g.currentWave.Enemies {
			if enemy.Health > 0 && enemy.Position.X == b.Position.X && enemy.Position.Y == b.Position.Y {
				enemy.TakeDamage(10)
				g.player.Bullets[i] = nil
				break
			}
		}
	}

	// Update enemies
	g.currentWave.UpdateEnemies()

	// Enemy shooting
	for _, enemy := range g.currentWave.Enemies {
		if enemy.Health <= 0 {
			continue
		}
		if enemy.CanShoot() {
			eb := enemy.Shoot()
			g.enemyBullets = append(g.enemyBullets, eb)
		}
	}

	// Update enemy bullets
	for i := 0; i < len(g.enemyBullets); i++ {
		b := g.enemyBullets[i]
		if b == nil {
			continue
		}
		b.Position.Y += b.Direction

		if b.Position.Y < 0 || b.Position.Y >= g.level.Height || g.level.IsWall(b.Position.X, b.Position.Y) {
			g.enemyBullets[i] = nil
			continue
		}

		if b.Position.X == g.player.Position.X && b.Position.Y == g.player.Position.Y {
			g.player.TakeDamage(10)
			g.enemyBullets[i] = nil
			continue
		}
	}

	// clean enemy bullets
	cleanBullets := make([]*bullet.Bullet, 0)
	for _, b := range g.enemyBullets {
		if b != nil {
			cleanBullets = append(cleanBullets, b)
		}
	}
	g.enemyBullets = cleanBullets

	// is enemy reached the bottom
	for _, enemy := range g.currentWave.Enemies {
		if enemy.Health > 0 && enemy.Position.Y >= g.level.Height-2 {
			g.player.TakeDamage(10)
			enemy.Health = 0
		}
	}

	if g.currentWave.IsCleared() {
		g.stats.SetWaveCompleted(g.waveNumber)
		g.waveNumber++
		g.spawnWave()
		g.currentWave.IncreaseEnemiesSpeed()
		return
	}

	// player & enemy collisions
	for _, enemy := range g.currentWave.Enemies {
		if enemy.Health > 0 && enemy.Position.X == g.player.Position.X && enemy.Position.Y == g.player.Position.Y {
			g.player.TakeDamage(20)
			enemy.TakeDamage(20)
			if g.player.Health <= 0 {
				break
			}
		}
	}

}

func (g *Game) render() {
	g.drawBuf.Reset()
	fmt.Fprint(os.Stdout, "\033[H")

	renderer.RenderLevel(g.level, g.drawBuf)
	renderer.RenderPlayer(g.player, g.drawBuf)

	for _, bullet := range g.player.Bullets {
		if bullet == nil {
			continue
		}
		renderer.RenderBullet(bullet, g.drawBuf)
	}

	for _, eb := range g.enemyBullets {
		renderer.RenderBullet(eb, g.drawBuf)
	}

	for _, enemy := range g.currentWave.Enemies {
		if enemy.Health > 0 {
			renderer.RenderEnemy(enemy, g.drawBuf)
		}
	}

	renderer.RenderStats(g.stats, g.player, g.level, g.drawBuf)
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}

func (g *Game) checkGameOver() bool {
	if g.player.Health <= 0 {
		g.gameOver = true
		g.gameOverMsg = "Game Over! You lost."
		return true
	}

	if g.stats.GetTimeLeft() <= 0 {
		g.gameOver = true
		g.gameOverMsg = "Time's Up! You lost."
		return true
	}

	if g.player.Ammunition <= 0 && !g.hasActiveBullets() {
		g.gameOver = true
		g.gameOverMsg = "Out of Ammunition! You lost."
		return true
	}

	return false
}

func (g *Game) hasActiveBullets() bool {
	for _, b := range g.player.Bullets {
		if b != nil {
			return true
		}
	}
	return false
}

func (g *Game) showGameOver() {
	utils.ClearScreen()
	width, height := g.level.Width, g.level.Height
	msg := g.gameOverMsg

	x := (width*2)/2 - len(msg)/2
	y := height / 2

	fmt.Printf("\033[%d;%dH%s", y, x, msg)
	fmt.Printf("\033[%d;%dHPress Enter to return to the menu.", y+2, (width*2)/2-20)

	if g.stats.WavesCleared > 0 {
		infoMsg := fmt.Sprintf("Waves Cleared: %d", g.stats.WavesCleared)
		fmt.Printf("\033[%d;%dH%s", y+4, (width*2)/2-len(infoMsg)/2, infoMsg)
	}

	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			continue
		}
		if key == keyboard.KeyEnter || char == '\r' {
			break
		}
	}
}
