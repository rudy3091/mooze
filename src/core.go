package mooze

import (
	"fmt"
	"os"
	"syscall"

	"github.com/RudyPark3091/mooze/src/util"
)

type mooze struct {
	tty      *os.File
	renderer *Renderer
	history  *History
	url      string
	method   string
	context  string
	os       string
	message  string
}

func Run() {
	mooze := new(mooze)
	tty := openTty()
	r := NewRenderer()

	state := r.ToRawMode(tty)

	// clears console
	// every screen context should be printed after this line
	r.ClearConsoleUnix()

	str := ""

CORE:
	for {
		buf := make([]byte, syscall.SizeofInotifyEvent)
		r.ReadChar(tty, buf)
		rn := util.BytesToRune(buf)

		str = str + string(rn)
		r.RenderTextTo(50, 50, string(rn))

		r.RenderTextTo(2, 1, "Buffer: %+v", buf)

		// process user inputs
		// -------------------------------------------------------------------
		switch rn {
		// exit application
		case rune(Q):
			r.ClearConsoleUnix()
			break CORE
		case rune(I):
			fmt.Print("\x1B[4m")

			r.UseNonblockIo(tty, false)
			r.RestoreState(tty, state)
		case rune(ENTER):
			r.ClearLine()
			r.CursorX = 0
			r.CursorY = 1
			// digest message
			mooze.message = str
			str = ""

			r.UseNonblockIo(tty, true)
			state = r.ToRawMode(tty)

			r.RenderTextTo(
				r.TtyRow-1, 1,
				NewColorContext("ffff00").Colorize(mooze.message),
			)
		case rune(BACKSPACE):
			// backspace
			r.MoveCursorLeft()
			r.WriteChar([]byte{32, 0, 0, 0, 0})
			r.MoveCursorLeft()
			str = str[0 : len(str)-2]
			continue
		}

		// output
		// ----------------------------------------------------------------------
		r.WriteChar(buf)
		r.RenderTextTo(
			r.TtyRow, 1,
			NewColorContext("ff0000", "ffffff").Colorize("Cursor Coord: %3d, %3d"),
			r.CursorX, r.CursorY,
		)
	}

	r.RestoreState(tty, state)
}
