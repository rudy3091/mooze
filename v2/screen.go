package v2

import (
	"fmt"
)

func LoadAlternateScreen() {
	fmt.Println("\033[?1049h\033[H")
}

func UnloadAlternateScreen() {
	fmt.Println("\033[?1049l")
}
