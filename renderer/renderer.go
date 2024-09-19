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

func RenderStats(s *stats.Stats, p *player.Player, lvl *level.Level, buf *bytes.Buffer) {
	timeLeft := s.GetTimeLeft()
	buf.WriteString(fmt.Sprintf("\033[%d;%dH-- Stats --", lvl.Height+1, 1))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHFPS: %.2f", lvl.Height+2, 1, s.FPS))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHHealth: %03d", lvl.Height+3, 1, p.Health))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHAmmunition: %03d", lvl.Height+4, 1, p.Ammunition))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHTime Left: %d", lvl.Height+5, 1, timeLeft))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHCurrent Wave: %d", lvl.Height+6, 1, s.CurrentWave))
	buf.WriteString(fmt.Sprintf("\033[%d;%dHTotal Waves Cleared: %d", lvl.Height+7, 1, s.WavesCleared))
}
