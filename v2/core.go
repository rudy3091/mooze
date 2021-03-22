package v2

import (
	"time"
)

func Run() {
	n := 0
	s := NewScreen()

	s.LoadAlternateScreen()
	s.Println("mooze v2")
	s.ReadNumber(&n)
	s.Println("input number: ", n)
	time.Sleep(time.Second * 2)
	s.UnloadAlternateScreen()
}
