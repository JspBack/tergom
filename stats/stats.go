package stats

import (
	"sync"
	"time"
)

const MAX_SAMPLES = 100

type Stats struct {
	Start        time.Time
	Frames       int
	FPS          float64
	TimeLeft     int
	CurrentWave  int
	WavesCleared int
	mutex        sync.Mutex
	timerStart   time.Time
	Logs         []string
}

func NewStats() *Stats {
	return &Stats{
		Start:        time.Now(),
		Frames:       0,
		FPS:          0.0,
		TimeLeft:     60,
		CurrentWave:  1,
		WavesCleared: 0,
		Logs:         []string{"DEBUG: Game started"},
	}
}

func (s *Stats) InitializeTimer(seconds int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TimeLeft = seconds
	s.timerStart = time.Now()
}

func (s *Stats) Update() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Frames++
	elapsed := time.Since(s.Start).Seconds()
	if s.Frames >= MAX_SAMPLES {
		s.FPS = float64(s.Frames) / elapsed
		s.Frames = 0
		s.Start = time.Now()
	}

	if time.Since(s.timerStart) >= time.Second {
		if s.TimeLeft > 0 {
			s.TimeLeft--
		}
		s.timerStart = time.Now()
	}
}

func (s *Stats) GetTimeLeft() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.TimeLeft
}

func (s *Stats) SetWaveCompleted(waveNumber int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if waveNumber > s.WavesCleared {
		s.WavesCleared = waveNumber
		s.CurrentWave = waveNumber + 1
	}
}

func (s *Stats) IncreaseTime(seconds int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TimeLeft += seconds
}

func (s *Stats) AddLog(log string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Logs = append(s.Logs, log)
}
