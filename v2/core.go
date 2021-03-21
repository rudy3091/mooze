package v2

import (
	"fmt"
	"time"
)

func Run() {
	LoadAlternateScreen()
	fmt.Println("mooze v2")
	time.Sleep(time.Second * 2)
	UnloadAlternateScreen()
}
