package v3

import (
	"fmt"
	"time"

	"github.com/rudy3091/mooze/v3/ui"
)

func Run() {
	ui.LoadAlternateScreen()
	fmt.Println("Hello v3")
	time.Sleep(time.Second * 2)
	ui.UnloadAlternateScreen()
}
