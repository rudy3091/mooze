package mooze

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type input int

const (
	ENTER     input = 13
	ESCAPE    input = 27
	Q         input = 113
	BACKSPACE input = 127
)

func openTty() *os.File {
	in, err := os.OpenFile("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	return in
}

func UiInit() {
	var err error
	in := openTty()
	r := NewRenderer()

	state := r.ToRawMode(in)

	// clears console
	// every screen context should be printed after this line
	r.ClearConsoleUnix()

	r.UseNonblockIo(in, true)

	buf := make([]byte, syscall.SizeofInotifyEvent)
	str := ""
	for {
		r.ReadChar(in, buf)
		str = str + string(buf)

		if buf[0] == byte(Q) || buf[0] == byte(ESCAPE) {
			r.ClearConsoleUnix()
			break
		}
		if buf[0] == byte(ENTER) {
			fmt.Print("\x1B[2K")
			r.CursorX = 0
			r.CursorY = 1
		}
		if buf[0] == byte(BACKSPACE) {
			// backspace
		}

		r.WriteChar(buf)
		// fmt.Print("\x1B[50;1H")
		// fmt.Print("\x1B[2K")
		// fmt.Print("Cursor Position: ")
		// fmt.Printf("%d, %d", r.CursorX, r.CursorY)
		// fmt.Printf("\x1B[%d;%dH", r.CursorY, r.CursorX)
		r.TargetStdout(r.TtyRow, 1, "Cursor Coord: %3d, %3d", r.CursorX, r.CursorY)
	}

	s := "BLUE"
	NewConsoleFg("ff0000").Println("Red")
	NewConsoleFg("0000ff").Printf("color: %s\n", s)

	err = terminal.Restore(int(in.Fd()), state)
	if err != nil {
		panic(err)
	}
}
