package v2

import (
	"fmt"
)

type Screen struct {
}

func NewScreen() *Screen {
	return &Screen{}
}

func (s *Screen) LoadAlternateScreen() {
	fmt.Println("\033[?1049h\033[H")
}

func (s *Screen) UnloadAlternateScreen() {
	fmt.Println("\033[?1049l")
}

func (s *Screen) Println(v ...interface{}) {
	fmt.Println(v...)
}

func (s *Screen) Print(v ...interface{}) {
	fmt.Print(v...)
}

func (s *Screen) MoveCursorTo(x, y int) {
	fmt.Println("\033[" + string(x) + ";" + string(y) + "H")
}

func (s *Screen) ClearLine() {
	fmt.Print("\033[2K")
}

func (s *Screen) ClearScreen() {
	s.MoveCursorTo(0, 0)
	fmt.Print("\033[2J")
}

func (s *Screen) HideCursor() {
	fmt.Print("\033[?25l")
}

func (s *Screen) ShowCursor() {
	fmt.Print("\033[?25h")
}
