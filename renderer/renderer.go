package renderer

import (
	"bytes"
	"fmt"
	"tergom/bullet"
	"tergom/enemy"
	"tergom/level"
	"tergom/player"
	"tergom/stats"
)

func RenderLevel(lvl *level.Level, buf *bytes.Buffer) {
	for h := 0; h < lvl.Height; h++ {
		for w := 0; w < lvl.Width; w++ {
			switch lvl.Data[h][w] {
			case level.NOTHING:
				buf.WriteString("  ")
			case level.WALL:
				buf.WriteString("██")
			}
		}
		buf.WriteString("\n")
	}
}

func RenderBullet(b *bullet.Bullet, buf *bytes.Buffer) {
	if b.Position.Y < 0 || b.Position.Y >= b.Level.Height || b.Position.X < 0 || b.Position.X >= b.Level.Width {
		return
	}
	if b.IsPlayerBullet {
		buf.WriteString(fmt.Sprintf("\033[%d;%dH| ", b.Position.Y+1, b.Position.X*2+1))
	} else {
		buf.WriteString(fmt.Sprintf("\033[%d;%dHK ", b.Position.Y+1, b.Position.X*2+1))
	}
}

func RenderPlayer(p *player.Player, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("\033[%d;%dH◯◯", p.Position.Y+1, p.Position.X*2+1))
}

func RenderEnemy(e *enemy.Enemy, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("\033[%d;%dHEE", e.Position.Y+1, e.Position.X*2+1))
}

func RenderStats(s *stats.Stats, p *player.Player, lvl *level.Level, buf *bytes.Buffer, width int) {
	timeLeft := s.GetTimeLeft()

	leftX := 1
	rightX := width/2 + 1
	baseY := lvl.Height + 1

	buf.WriteString(fmt.Sprintf("\033[%d;%dH-- Player Stats --", baseY, leftX))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHHealth: %03d", baseY+1, leftX, p.Health))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHAmmunition: %03d", baseY+2, leftX, p.Ammunition))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHScore: %d", baseY+3, leftX, p.Score))

	buf.WriteString(fmt.Sprintf("\033[%d;%dH-- Game Stats --", baseY, rightX))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHFPS: %.2f", baseY+1, rightX, s.FPS))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHTime Left: %d", baseY+2, rightX, timeLeft))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHCurrent Wave: %d", baseY+3, rightX, s.CurrentWave))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHTotal Waves Cleared: %d", baseY+4, rightX, s.WavesCleared))

	buf.WriteString(fmt.Sprintf("\033[%d;%dH-- LOGS --", baseY, rightX*2))
	startIndex := len(s.Logs) - 4
	if startIndex < 0 {
		startIndex = 0
	}
	for i := 0; i < 4; i++ {
		logIndex := startIndex + i
		var log string
		if logIndex < len(s.Logs) {
			log = s.Logs[logIndex]
		} else {
			log = ""
		}
		buf.WriteString(fmt.Sprintf("\033[%d;%dH%-50s", baseY+i+1, rightX*2, log))
	}
}
