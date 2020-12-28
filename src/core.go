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
	body    string
	url     string
	method  string
	context string
	os      string
	message string
}

func NewMooze() *mooze {
	return &mooze{
		openTty(),
		"",           // history
		"",           // header
		"",           // body
		"",           // url
		"GET",        // method
		"",           // context
		runtime.GOOS, // os
		"",           // message
	}
}

func Run() {
	// applications state
	mooze := NewMooze()
	// input mode state
	f := NewFlags()

	tty := openTty()
	r := NewRenderer()
	h := NewHistoryWriter()
	req := NewRequest("", GET, "", "")

	state := r.ToRawMode(tty)
	defer r.RestoreState(tty, state)
	r.UseNonblockIo(tty, true)
	defer r.UseNonblockIo(tty, false)

	// clears console
	// every screen context should be printed after this line
	r.ClearConsoleUnix()

	msg := ""
	// renders character to a console if wflag is true
	wflag := false

	// renderwindow test
	w := Window{10, 10, 10, 10, "\u2500", "\u2502", "\u250c", "\u2510", "\u2514", "\u2518"}
	r.RenderWindow(w)

CORE:
	for {
		if wflag {
			r.ShowCursor()
		} else {
			r.HideCursor()
			defer r.ShowCursor()
		}
		buf := make([]byte, syscall.SizeofInotifyEvent)
		r.RenderTextTo(r.TtyRow-1, 1, "\x1B[4m\x1B[31mu\x1B[0mrl: %s", mooze.url)
		r.RenderTextTo(r.TtyRow-2, 1, "\x1B[4m\x1B[31mm\x1B[0method: %s", mooze.method)
		r.RenderTextTo(r.TtyRow-3, 1, "\x1B[4m\x1B[31mh\x1B[0meader: %s", mooze.header)
		r.RenderTextTo(r.TtyRow-4, 1, "\x1B[4m\x1B[31mb\x1B[0mody: %s", mooze.body)
		r.RenderTextTo(r.TtyRow-5, 1, "history: %s", mooze.history)

		r.RenderTextTo(r.TtyRow-6, 1, "\x1B[4m\x1B[38;2;255;50;15msend\x1B[0m")

		r.ReadChar(tty, buf)
		rn := util.BytesToRune(buf)

		r.RenderTextTo(2, 1, "Rune num: %d", rn)
		r.RenderTextNoClear(2, 20, "Buffer: %d", buf)

		// drawing window
		// horizontal line
		r.HorizontalLine(5, 1, r.TtyCol)
		r.HorizontalLine(r.TtyRow-8, 1, r.TtyCol)
		// vertical line
		r.VerticalLine(5, 1, r.TtyRow-13)
		r.VerticalLine(5, r.TtyCol, r.TtyRow-13)

		if wflag {
			r.WriteChar(buf)
		}

		/* process user inputs
		 * -------------------------------------------------------------------
		 * FIXME 1: Escape && Arrow key input has same rune value
		 * FIXED 2: typing Ctrl-j makes unintentional new line
		 * FIXME 3: if buffer's length == 1 Cursor coordinate not works
		 *
		 * input mode:
		 * - u for url
		 * - m for method
		 * - h for header
		 * - b for body
		 *
		 * send request
		 * - Ctrls
		 *
		 * quit application
		 * - Ctrlq
		 *
		 */

		switch rn {
		// exit application
		case rune(CTRLQ):
			r.ClearConsoleUnix()
			break CORE

		// send request
		case rune(CTRLS):
			res := req.Send()
			defer res.Body.Close()
			rbytes := req.Body(res)

			// render response body to console
			r.RestoreState(tty, state)
			r.RenderTextTo(6, 3, string(rbytes))
			h.Write(string(rbytes))
			state = r.ToRawMode(tty)

			wflag = false

		// process user input
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

			// set request infos
			switch {
			case f.url:
				mooze.url = mooze.message
				r.RenderTextTo(
					r.TtyRow-1, 1,
					NewColorContext("ffff55").Colorize("url: "+mooze.url),
				)
				r.RenderTextTo(3, 1, "\x1B[2K")
				req.url = mooze.url
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
			// \x1B[?25h rendered by ShowCursor not shown in console
			r.MoveCursorLeft()
			r.Backspace()
			if len(msg) > 1 {
				msg = msg[0 : len(msg)-1]
			} else {
				msg = ""
			}

		// FIXED 2
		// Ctrl-j == Ctrl-Enter
		// do nothing
		case rune(CTRLJ):
			wflag = false

		// get Request Url from user
		case rune(U):
			f.url = true
			r.RenderTextTo(3, 1, NewColorContext("ff8888").Colorize("-- URL --"))
			wflag = true

		// appending input character into message
		default:
			if wflag {
				msg = msg + string(rn)
			}
		}

		// output
		// ----------------------------------------------------------------------
		r.RenderTextTo(
			r.TtyRow, 1,
			NewColorContext("ff0000", "ffffff").Colorize("Cursor Coord: %3d, %3d"),
			r.CursorX, r.CursorY,
		)
	}

	r.RestoreState(tty, state)
}
