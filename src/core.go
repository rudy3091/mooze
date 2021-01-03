// TODO: fix http response rune-width

package mooze

import (
	"io"
	"os"
	"runtime"

	"github.com/gdamore/tcell"
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
	req     *MoozeRequest
	history string // will update to file
	os      string
}

func NewMooze() *mooze {
	return &mooze{
		term:    initTerminal(),
		ms:      NewMoozeScreen(),
		req:     NewMoozeRequest(),
		history: "",
		os:      runtime.GOOS,
	}
}

func initTerminal() *terminal.Terminal {
	term := terminal.NewTerminal(
		struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, "",
	)
	term.SetPrompt(NewColorContext("ff5555").Colorize("> "))
	return term
}

func (m *mooze) initLayout(w, h int) []*MoozeWindow {
	window := []*MoozeWindow{}
	urlHeight := 1
	statusHeight := 7
	rHeight := h - (urlHeight + statusHeight)

	rw := 0
	if w%2 == 1 {
		rw = w/2 + 1
	} else {
		rw = w / 2
	}

	w1 := NewMoozeWindow(urlHeight, 0, rHeight, rw, false)
	w1.Title("req")
	w1.Content([]string{"request body"})
	window = append(window, w1)

	w2 := NewMoozeWindow(urlHeight, rw, rHeight, w-rw, false)
	w2.Title("res")
	window = append(window, w2)

	w3 := NewMoozeWindow(h-statusHeight, 0, statusHeight, w, false)
	w3.Title("status")
	w3.Content([]string{
		"url: " + m.req.url,
		"method: " + methodTypeToString(m.req.method),
	})
	window = append(window, w3)

	return window
}

func (m *mooze) renderLayout(w []*MoozeWindow) {
	for _, window := range w {
		m.ms.RenderWindow(window, ToStyle("white"))
	}
}

func (m *mooze) statusCode(w *MoozeWindow) {
	style := ToStyle("white")
	switch m.req.resCode / 100 {
	case 2:
		style = ToStyle("white", "green")
	case 3:
		style = ToStyle("white", "yellow")
	case 4:
		style = ToStyle("white", "red")
	case 5:
		style = ToStyle("white", "crimson")
	}

	_y := w.y + w.sizeY - (len(m.req.resStatus) + 2)
	m.ms.Print(w.x+1, _y, m.req.resStatus, style)
}

func (m *mooze) readLine() string {
	m.ms.r.ShowCursor()
	m.ms.r.MoveCursorTo(1, 1)
	m.ms.r.ClearLine()
	l, err := m.term.ReadLine()
	if err != nil {
		panic(err)
	}
	m.ms.r.MoveCursorTo(1, 1)
	m.ms.r.ClearLine()
	m.ms.r.HideCursor()
	return l
}

func Run() {
	// applications state
	mooze := NewMooze()
	mooze.ms.InitScreen(false)
	defer mooze.ms.Exit(1)

	w, h := mooze.ms.Size()
	layout := mooze.initLayout(w, h)

CORE:
	for {
		mooze.renderLayout(layout)
		mooze.statusCode(layout[1])
		mooze.ms.Show()

		ev := mooze.ms.EmitEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			w, h = mooze.ms.Size()
			layout = mooze.initLayout(w, h)
			mooze.renderLayout(layout)
			mooze.ms.Reload()
		case *tcell.EventKey:
			switch ev.Rune() {
			case rune(Q):
				break CORE

			case rune(U):
				mooze.req.url = mooze.readLine()
				layout[2].content[0] =
					"url: " + mooze.req.url
				mooze.renderLayout(layout)
				mooze.ms.Show()

			case rune(CTRLS):
				res := mooze.req.Send()
				defer res.Body.Close()
				rData := mooze.req.Body(res)
				layout[1].Content(mooze.req.Prettify(rData))
				mooze.renderLayout(layout)
				mooze.req.resStatus = res.Status
				mooze.req.resCode = res.StatusCode
				mooze.ms.Show()
			}
		}
	}
	mooze.ms.Exit(0)
}
