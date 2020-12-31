package mooze

import (
	"github.com/gdamore/tcell"
	"io"
	"os"
	"runtime"

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
	term    *terminal.Terminal
	ms      *MoozeScreen
	history string // will update to file
	os      string

	header  string
	body    string
	url     string
	method  methodtype
	message string
}

func NewMooze() *mooze {
	return &mooze{
		term:    InitTerminal(),
		ms:      NewMoozeScreen(),
		method:  GET,
		os:      runtime.GOOS,
		history: "",
		header:  "",
		body:    "",
		url:     "",
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

func initLayout(w, h int) []*MoozeWindow {
	window := []*MoozeWindow{}
	urlHeight := 1
	statusHeight := 5
	rHeight := h - (urlHeight + statusHeight)

	rw := 0
	if w%2 == 1 {
		rw = w/2 + 1
	} else {
		rw = w / 2
	}

	w1 := NewMoozeWindow(urlHeight, 0, rHeight, rw, false)
	w1.Title("req")
	window = append(window, w1)

	w2 := NewMoozeWindow(urlHeight, rw, rHeight, w-rw, false)
	w2.Title("res")
	window = append(window, w2)

	w3 := NewMoozeWindow(h-statusHeight, 0, statusHeight, w, false)
	w3.Title("status")
	window = append(window, w3)

	return window
}

func (m *mooze) renderLayout(w []*MoozeWindow) {
	for _, window := range w {
		m.ms.RenderWindow(window, ToStyle("red"))
	}
}

func Run() {
	// applications state
	mooze := NewMooze()
	mooze.ms.InitScreen(false)

	// // input mode state
	// f := NewFlags()

	// r := NewRenderer()
	// h := NewHistoryWriter()
	// req := NewRequest("", GET, "", "")

	// defer r.ClearConsoleUnix()

	// r.ClearConsoleUnix()

	// msg := ""
	// wflag := false
	w, h := mooze.ms.Size()
	layout := initLayout(w, h)

CORE:
	for {
		mooze.renderLayout(layout)
		mooze.ms.Show()

		ev := mooze.ms.EmitEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			w, h = mooze.ms.Size()
			layout = initLayout(w, h)
			mooze.renderLayout(layout)
			mooze.ms.Reload()
		case *tcell.EventKey:
			switch ev.Rune() {
			case 'q':
				break CORE
			}
		}
	}
	mooze.ms.Exit(0)
}
