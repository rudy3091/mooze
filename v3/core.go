package v3

import (
	"log"
	"os"
	"syscall"
	"time"

	"github.com/rudy3091/mooze/v3/ui"
	"golang.org/x/term"
)

func openTty() *os.File {
	in, err := os.OpenFile("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		log.Fatal("cannot open /dev/tty")
		panic(err)
	}
	return in
}

func makeRaw() (int, *term.State) {
	fd := int(openTty().Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal("cannot make terminal raw")
		panic(err)
	}
	return fd, state
}

func restoreRaw(fd int, s *term.State) {
	term.Restore(fd, s)
}

func Run() {
	ui.LoadAlternateScreen()
	fd, state := makeRaw()

	defer restoreRaw(fd, state)
	defer ui.UnloadAlternateScreen()

	ui.KeyBindings()
	time.Sleep(time.Second * 2)
}
