package util

import (
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	n := time.Now()
	s := NewStopWatch()
	time.Sleep(time.Millisecond * 50)
	s.Start()
	if s.StartTime.UnixNano() <= n.UnixNano() {
		t.Error()
	}
}

func TestStop(t *testing.T) {
	s := NewStopWatch()
	s.Start()
	time.Sleep(time.Millisecond * 50)
	s.Stop()

	if s.StopTime.UnixNano() < s.StartTime.UnixNano() {
		t.Error()
	}
	// millisecond = 1e6 unixnano
	if s.Duration < 50*1e6 {
		t.Error()
	}
	if s.DurationSec < 0.05 {
		t.Error()
	}
}

func TestLap(t *testing.T) {
	s := NewStopWatch()
	s.Start()
	s.Lap()
	time.Sleep(time.Millisecond * 50)
	s.Lap()
	time.Sleep(time.Millisecond * 50)
	s.Stop()

	if len(s.Laps) != 2 {
		t.Error()
	}
	if s.Laps[0].UnixNano() >= s.Laps[1].UnixNano() {
		t.Error()
	}
}

func TestReset(t *testing.T) {
	s := NewStopWatch()
	s.Start()
	s.Lap()
	s.Stop()
	s.Reset()

	if len(s.Laps) != 0 {
		t.Error()
	}
	if s.Duration != 0 {
		t.Error()
	}
	if s.DurationSec != 0 {
		t.Error()
	}
}
