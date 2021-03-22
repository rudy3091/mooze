package v2

import (
	"time"
)

func Run() {
	s := NewScreen()
	s.LoadAlternateScreen()
	s.Println("mooze v2")
	time.Sleep(time.Second * 2)
	s.UnloadAlternateScreen()
}
