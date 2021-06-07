package ui

import (
	"fmt"
)

func LoadAlternateScreen() {
	fmt.Print("\033[?1049h\033H")
}

func UnloadAlternateScreen() {
	fmt.Print("\033[?1049l")
}
