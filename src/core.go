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
	wflag := false
	w, h := mooze.ms.Size()

CORE:
	for {
		status := struct {
			width  int
			height int
		}{w - 1, 6}
		window := NewMoozeWindow(h-status.height-1, 1, status.height, status.width-1, false)
		window.Title("hi")
		if !wflag {
			mooze.ms.RenderWindow(
				window,
				ToStyle("red"),
			)
		}
		mooze.ms.Show()

		ev := mooze.ms.EmitEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			w, h = mooze.ms.Size()
			mooze.ms.Clear()
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
