package v2

import (
	"fmt"
)

type Screen struct {
}

func NewScreen() *Screen {
	return &Screen{}
}

func (s Screen) LoadAlternateScreen() {
	fmt.Println("\033[?1049h\033[H")
}

func (s Screen) UnloadAlternateScreen() {
	fmt.Println("\033[?1049l")
}

func (s Screen) Println(v ...interface{}) {
	fmt.Println(v...)
}
