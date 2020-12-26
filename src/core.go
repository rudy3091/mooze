package mooze

import (
	"os"
	"runtime"
	"syscall"

	"github.com/RudyPark3091/mooze/src/util"
)

type flags struct {
	url     bool
	method  bool
	header  bool
	history bool
}

func NewFlags() *flags {
	return &flags{false, false, false, false}
}

type mooze struct {
	tty     *os.File
	history string
	header  string
	url     string
	method  string
	context string
	os      string
	message string
}

func NewMooze() *mooze {
	return &mooze{
		openTty(),
		"",
		"",
		"",
		"",
		"",
		runtime.GOOS,
		"",
	}
}

func Run() {
	mooze := new(mooze)
	f := NewFlags()
	tty := openTty()
	r := NewRenderer()

	state := r.ToRawMode(tty)

	// clears console
	// every screen context should be printed after this line
	r.ClearConsoleUnix()

	msg := ""

CORE:
	for {
		// renders character to a console if wflag is true
		wflag := true
		buf := make([]byte, syscall.SizeofInotifyEvent)
		r.ReadChar(tty, buf)
		rn := util.BytesToRune(buf)

		r.RenderTextTo(50, 50, string(rn))
		r.RenderTextTo(2, 1, "Buffer: %d", rn)

		// process user inputs
		// -------------------------------------------------------------------
		// FIXME 1: Escape && Arrow key input has same rune value
		// FIXED 2: typing Ctrl-j twice makes unintentional new line
		// FIXME 3: if buffer's length == 1 Cursor coordinate not works

		switch rn {
		// exit application
		case rune(Q):
			r.ClearConsoleUnix()
			break CORE

		// process input message
		case rune(ENTER):
			r.ClearLine()
			r.CursorX = 1
			r.CursorY = 1
			// digest message
			mooze.message = msg
			msg = ""

			r.RenderTextTo(
				3, 1,
				NewColorContext("ffff55").Colorize(mooze.message),
			)
			switch {
			case f.url:
				mooze.url = mooze.message
				r.RenderTextTo(
					r.TtyRow-1, 1,
					NewColorContext("ffff55").Colorize("url: "+mooze.url),
				)
				r.RenderTextTo(3, 1, "\x1B[2K")
				f.url = false

			case f.method:
				mooze.method = mooze.message

			case f.header:
				mooze.header = mooze.message

			case f.history:
				mooze.history = mooze.message
			}
			wflag = false

		case rune(BACKSPACE):
			r.MoveCursorLeft()
			r.WriteChar([]byte{32, 0, 0, 0})
			r.MoveCursorLeft()
			if len(msg) > 1 {
				msg = msg[0 : len(msg)-1]
			} else {
				msg = ""
			}
			wflag = false

		case rune(TAB):
			wflag = false

		// FIXED 2
		case rune(CTRLJ):
			wflag = false

		// get Request Url from user
		case rune(U):
			f.url = true
			r.RenderTextTo(3, 1, NewColorContext("88ff88").Colorize("-- URL --"))
			wflag = false

		// appending input character into message
		default:
			msg = msg + string(rn)
		}

		// output
		// ----------------------------------------------------------------------
		r.RenderTextTo(
			r.TtyRow, 1,
			NewColorContext("ff0000", "ffffff").Colorize("Cursor Coord: %3d, %3d"),
			r.CursorX, r.CursorY,
		)
		if !wflag {
			continue
		}
		r.WriteChar(buf)
	}

	r.RestoreState(tty, state)
}
