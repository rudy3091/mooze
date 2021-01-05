package util

import (
	"time"
)

type StopWatch struct {
	Laps        []time.Time
	StartTime   time.Time
	StopTime    time.Time
	Duration    int64
	DurationSec float64
}

func NewStopWatch() *StopWatch {
	return &StopWatch{nil, time.Now(), time.Now(), 0, 0}
}

func (s *StopWatch) Start() {
	s.StartTime = time.Now()
}

func (s *StopWatch) Stop() int64 {
	s.StopTime = time.Now()
	s.Duration = s.StopTime.UnixNano() - s.StartTime.UnixNano()
	s.DurationSec = float64(s.Duration) / 1e9
	return s.StopTime.UnixNano() - s.StartTime.UnixNano()
}

func (s *StopWatch) Lap() {
	s.Laps = append(s.Laps, time.Now())
}

func (s *StopWatch) Reset() {
	s.Laps = nil
	s.StartTime = time.Now()
	s.StopTime = time.Now()
	s.Duration = 0
	s.DurationSec = 0
}
