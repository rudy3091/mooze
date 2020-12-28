// TODO: stop displaying rune & buffer code - for debug
// TODO: terminal width processing - what if input exceeds TtyCol?
//       -> Maybe I should change input methology terminal.Newterminal()?
// TODO: reactive terminal - change the layout if terminal size varies

package mooze

import (
	"io"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/RudyPark3091/mooze/src/util"
	"golang.org/x/crypto/ssh/terminal"
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
	term    *terminal.Terminal
	history string
	header  string
	body    string
	url     string
	method  methodtype
	os      string
	message string
}

func NewMooze() *mooze {
	return &mooze{
		tty:     openTty(),
		term:    InitTerminal(),
		history: "",
		header:  "",
		body:    "",
		url:     "",
		method:  GET,
		os:      runtime.GOOS,
		message: "",
	}
}

func InitTerminal() *terminal.Terminal {
	term := terminal.NewTerminal(
		struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, "",
	)
	term.SetPrompt(NewColorContext("ff5555").Colorize("> "))
	return term
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

	defer r.ClearConsoleUnix()
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

	// status window
	sx := 7
	sy := r.TtyCol() - 1
	statusWindow := NewWindow(r.TtyRow()-sx-1, 1, sx, sy, "ffaaaa")

	// drawing status window
	r.RenderWindow(statusWindow)
	r.RenderTextNoClear(r.TtyRow()-2, 3, "\x1B[4m\x1B[31mu\x1B[0mrl: %s", FColors.yellow.Colorize(mooze.url))
	r.RenderTextNoClear(r.TtyRow()-3, 3, "\x1B[4m\x1B[31mm\x1B[0method: %s",
		FColors.yellow.Colorize(methodTypeToString(GET)))
	r.RenderTextNoClear(r.TtyRow()-4, 3, "\x1B[4m\x1B[31mh\x1B[0meader: %s", mooze.header)
	r.RenderTextNoClear(r.TtyRow()-5, 3, "\x1B[4m\x1B[31mb\x1B[0mody: %s", mooze.body)
	r.RenderTextNoClear(r.TtyRow()-6, 3, "history: %s", mooze.history)
	r.RenderTextNoClear(r.TtyRow()-7, 3, FColors.yellow.Colorize("send: ctrl + s"))

	sBarBg := NewColorContext("bfbfbf", "555555")

CORE:
	for {
		// CORE LOGIC
		// ----------------------------------------------------------------------
		// Render StatusBar at (TtyRow(), 1)
		nowBg := NewColorContext(MoozeStatusBar.Now.Hex, "555555")
		r.RenderTextTo(r.TtyRow(), 1, sBarBg.Colorize(strings.Repeat(" ", r.TtyCol())))
		r.RenderTextNoClear(
			r.TtyRow(), 1,
			NewColorContext("555555", MoozeStatusBar.Now.Hex).Colorize("%s")+nowBg.Colorize("\ue0b0"),
			MoozeStatusBar.Now.Name,
		)
		r.RenderTextNoClear(
			r.TtyRow(), r.TtyCol()-6,
			nowBg.Colorize("\ue0b2")+
				NewColorContext("555555", MoozeStatusBar.Now.Hex).Colorize("%-2d, %-2d"),
			r.CursorX, r.CursorY,
		)

		if wflag {
			r.ShowCursor()
		} else {
			r.HideCursor()
			defer r.ShowCursor()
		}
		buf := make([]byte, syscall.SizeofInotifyEvent)

		if !wflag {
			r.ReadChar(tty, buf)
		}
		if wflag {
			line, _ := mooze.term.ReadLine()
			r.MoveCursorTo(1, 1)
			r.ClearLine()

			MoozeStatusBar.Now = MoozeStatusBar.Normal
			mooze.message = line
			switch {
			case f.url:
				mooze.url = mooze.message
				r.RenderTextNoClear(
					r.TtyRow()-2, 8,
					NewColorContext("ffff55").Colorize(mooze.url),
				)
				req.url = mooze.url
				f.url = false
			}
			wflag = false
		}
		rn := util.BytesToRune(buf)

		r.RenderTextTo(2, 1, "Rune num: %-10d Buffer: %d", rn, buf)

		/* process user inputs
		 * -------------------------------------------------------------------
		 * FIXME 1: Escape && Arrow key input has same rune value
		 * FIXED 2: typing Ctrl-j makes unintentional new line
		 * FIXED 3: if buffer's length == 1 Cursor coordinate not works
		 * FIXME 4: no need to re-render all screen in every input
		 * FIXME 5: can't type 'u' in url input mode
		 *
		 * input mode:
		 * - u for url
		 * - m for method
		 * - h for header
		 * - b for body
		 *
		 * send request
		 * - Ctrl + s
		 *
		 * quit application
		 * - Ctrl + q
		 *
		 */

		switch rn {
		// exit application
		case rune(CTRLQ):
			r.ClearConsoleUnix()
			break CORE

		case rune(ESCAPE):
			msg = ""
			switch {
			case f.url:
				f.url = false
			case f.method:
				f.method = false
			case f.header:
				f.header = false
			case f.history:
				f.history = false
			}
			MoozeStatusBar.Now = MoozeStatusBar.Normal
			r.ClearLine()
			r.CursorX = 1
			r.CursorY = 1
			wflag = false

		// send request
		case rune(CTRLS):
			res := req.Send()
			defer res.Body.Close()
			rbytes := req.Body(res)

			// render response body to console
			r.RestoreState(tty, state)
			r.RenderTextTo(6, 1, string(rbytes))
			h.Write(string(rbytes))
			state = r.ToRawMode(tty)

			wflag = false

		// process user input
		case rune(ENTER):
			MoozeStatusBar.Now = MoozeStatusBar.Normal
			r.ClearLine()
			r.CursorX = 1
			r.CursorY = 1
			// digest message
			mooze.message = msg
			msg = ""

			// set request infos
			switch {
			case f.url:
				mooze.url = mooze.message
				r.RenderTextNoClear(
					r.TtyRow()-2, 8,
					NewColorContext("ffff55").Colorize(mooze.url),
				)
				req.url = mooze.url
				f.url = false

			case f.method:
				mooze.method = 0

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
			MoozeStatusBar.Now = MoozeStatusBar.Url
			wflag = true
		}
	}

	r.RestoreState(tty, state)
}
