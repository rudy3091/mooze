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
	term.SetPrompt("\x1B[0m\x1B[31m> \x1B[0m")
	return term
}

func (m *mooze) getHorizontalLayout(w, h int) (*MoozeWindow, *MoozeWindow, *MoozeWindow) {
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

	w2 := NewMoozeWindow(urlHeight, rw, rHeight, w-rw, false)
	w2.Title("res")

	w3 := NewMoozeWindow(h-statusHeight, 0, statusHeight, w, false)
	w3.Title("status")
	w3.Content([]string{
		"url: " + m.req.url,
		"method: " + methodTypeToString(m.req.method),
	})

	return w1, w2, w3
}

func (m *mooze) getVerticalLayout(w, h int) (*MoozeWindow, *MoozeWindow, *MoozeWindow) {
	urlHeight := 1
	statusHeight := 7

	rh := (h - (urlHeight + statusHeight)) / 2
	if h%2 == 1 {
		rh = rh - 1
	}

	w1 := NewMoozeWindow(urlHeight, 0, rh, w, false)
	w1.Title("req")
	w1.Content([]string{"request body"})

	w2 := NewMoozeWindow(rh+1, 0, rh, w, false)
	w2.Title("res")

	w3 := NewMoozeWindow(h-statusHeight, 0, statusHeight, w, false)
	w3.Title("status")
	w3.Content([]string{
		"url: " + m.req.url,
		"method: " + methodTypeToString(m.req.method),
	})

	return w1, w2, w3
}

func (m *mooze) initLayout(w, h int) (*MoozeWindow, *MoozeWindow, *MoozeWindow) {
	if w > 100 {
		return m.getHorizontalLayout(w, h)
	} else {
		return m.getVerticalLayout(w, h)
	}
}

func (m *mooze) renderLayout(w ...*MoozeWindow) {
	for _, window := range w {
		m.ms.RenderWindow(window, ToStyle("white"))
	}
}

func (m *mooze) statusCode(w *MoozeWindow) {
	style := ToStyle("white")
	switch m.req.resCode / 100 {
	case 2:
		style = ToStyle("black", "green")
	case 3:
		style = ToStyle("black", "yellow")
	case 4:
		style = ToStyle("black", "red")
	case 5:
		style = ToStyle("black", "crimson")
	}

	_y := w.y + w.sizeY - (len(m.req.resStatus) + 2)
	m.ms.Print(w.x+1, _y, m.req.resStatus, style)
}

func (m *mooze) readLine() string {
	m.ms.r.ShowCursor()
	m.ms.r.MoveCursorTo(1, 1)
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
	wReq, wRes, wStatus := mooze.initLayout(w, h)

CORE:
	for {
		mooze.renderLayout(wReq, wRes, wStatus)
		mooze.statusCode(wRes)
		mooze.ms.Show()

		ev := mooze.ms.EmitEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			mooze.ms.r.ClearConsoleUnix()
			w, h = mooze.ms.Size()
			wReq, wRes, wStatus = mooze.initLayout(w, h)
			wRes.Content(mooze.req.data)
			mooze.ms.Reload()
			mooze.statusCode(wRes)
			mooze.ms.Show()

		case *tcell.EventKey:
			switch ev.Rune() {
			// exit application
			case rune(Q):
				break CORE

			// url input
			case rune(U):
				mooze.req.url = mooze.readLine()
				wStatus.content[0] =
					"url: " + mooze.req.url
				mooze.renderLayout(wReq, wRes, wStatus)
				mooze.ms.Show()

			// send Request
			case rune(CTRLS):
				// erase former response
				wRes.Content([]string{})
				mooze.ms.RenderWindow(wRes, ToStyle("red"))
				mooze.ms.Show()

				end := make(chan bool)
				defer close(end)
				var res Response

				go func() {
					res = mooze.req.Send()
					end <- true
				}()
				<-end

				defer res.Body.Close()
				rData := mooze.req.Body(res)
				mooze.req.data = mooze.req.Prettify(rData)
				wRes.Content(mooze.req.data)

				mooze.renderLayout(wReq, wRes, wStatus)
				mooze.req.resStatus = res.Status
				mooze.req.resCode = res.StatusCode
				mooze.ms.Show()
			}
		}
	}
	mooze.ms.Exit(0)
}
