package mooze

import (
	"io"
	"os"
	"runtime"

	"github.com/RudyPark3091/mooze/src/util"
	"github.com/gdamore/tcell"
	"golang.org/x/crypto/ssh/terminal"
)

type flags struct {
	url     bool
	method  bool
	header  bool
	body    bool
	history bool
}

func NewFlags() *flags {
	return &flags{false, false, false, false, false}
}

type mooze struct {
	term    *terminal.Terminal
	ms      *MoozeScreen
	req     *MoozeRequest
	editor  *MoozeEditor
	history string // will update
	os      string
	mode    *flags
}

func NewMooze() *mooze {
	return &mooze{
		term:    initTerminal(),
		ms:      NewMoozeScreen(),
		req:     NewMoozeRequest(),
		editor:  NewMoozeEditor(),
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
	urlHeight := 0
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
	urlHeight := 0
	statusHeight := 7

	rh := (h - (urlHeight + statusHeight)) / 2
	if h%2 == 1 {
		rh = rh - 1
	}

	w1 := NewMoozeWindow(urlHeight, 0, rh, w, false)
	w1.Title("req")

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

func (m *mooze) readLine(x, y int) string {
	m.ms.r.ShowCursor()
	m.ms.r.MoveCursorTo(x, y)
	l, err := m.term.ReadLine()
	if err != nil {
		panic(err)
	}
	m.ms.r.HideCursor()
	return l
}

func Run() {
	mooze := NewMooze()
	mooze.ms.InitScreen(false)
	defer mooze.ms.Exit(1)

	w, h := mooze.ms.Size()
	wReq, wRes, wStatus := mooze.initLayout(w, h)

	prompt := "\x1B[0m\x1B[31m> \x1B[0m"

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
			wReq.Content(mooze.req.Prettify([]byte(mooze.req.body)))
			wRes.Content(mooze.req.resData)
			mooze.statusCode(wRes)
			mooze.ms.Reload()

		case *tcell.EventKey:
			switch ev.Rune() {
			// exit application
			case rune(Q):
				break CORE

			// url input
			case rune(U):
				mooze.term.SetPrompt("\x1B[0m\x1B[42m\x1B[30m> ")
				wUrl := NewMoozeWindow(h/2-3, w/2-w/6, 3, w/3, false)
				wUrl.Title("type target url")
				mooze.ms.RenderWindow(wUrl, ToStyle("black", "green"))
				mooze.ms.Show()
				line := mooze.readLine(h/2-1, w/2-w/6+2)
				if line != "" {
					mooze.req.url = line
				}
				wStatus.content[0] =
					"url: " + mooze.req.url
				mooze.renderLayout(wReq, wRes, wStatus)
				mooze.ms.Reload()

			// "select" request method (post, get ...)
			case rune(M):
				mooze.term.SetPrompt("\x1B[0m\x1B[42m\x1B[30m> ")
				wMSelect := NewMoozeWindow(h/2-5, w/2-10, 2+3, 20, false)
				// 1: GET
				// 2: POST
				// 3: PUT
				// 4: PATCH
				// 5: DELETE
				wMSelect.Title("type method")
				wMSelect.Content([]string{"", "1: GET", "2: POST"})
				mooze.ms.RenderWindow(wMSelect, ToStyle("black", "green"))
				mooze.ms.Show()
				n := mooze.readLine(h/2-3, w/2-10+2)
				if n == "" {
					n = "1"
				}
				ni := util.ToInteger(n) - 1
				mooze.req.method = methodtype(ni)
				wStatus.content[1] =
					"method: " + methodTypeToString(mooze.req.method)
				mooze.ms.Reload()

			// body
			case rune(B):
				mooze.term.SetPrompt("\x1B[0m\x1B[42m\x1B[30m> ")
				mooze.req.body = mooze.editor.readLine(mooze)
				wReq.Content(mooze.req.Prettify([]byte(mooze.req.body)))
				// x := 7
				// y := 7
				// wBodyInput := NewMoozeWindow(5, 5, h-10, w-10, false)
				// wBodyInput.Title("type request body")
				// mooze.ms.RenderWindow(wBodyInput, ToStyle("black", "green"))
				// mooze.ms.Show()
				// x = 7
				// y = 7
				// bodyBuf := ""
				// for {
				// 	line := mooze.readLine(x, y)
				// 	if line == "" {
				// 		break
				// 	}
				// 	bodyBuf += line
				// 	x += 1
				// }
				// mooze.req.body = bodyBuf
				// wReq.Content(mooze.req.Prettify([]byte(bodyBuf)))
				mooze.ms.Reload()

			// options
			case rune(O):
				wOption := NewMoozeWindow(h/2-10, w/2-20, 20, 40, false)
				wOption.Content([]string{"options"})
				mooze.ms.RenderWindow(wOption, ToStyle("white", "blue"))
				mooze.ms.Show()
				mooze.readLine(1, 1)

			// send Request
			case rune(CTRLS):
				mooze.term.SetPrompt(prompt)
				// erase former response
				wRes.Content([]string{})
				// window frame to blink red
				// until server responds
				mooze.ms.RenderWindow(wRes, ToStyle("red"))
				jsonData := mooze.req.ParseJson(mooze.req.body)
				mooze.ms.Show()

				end := make(chan bool)
				defer close(end)
				var res Response

				go func() {
					res = mooze.req.Send(mooze.req.method, ReqArgs{
						h:   "application/json",
						buf: jsonData,
					})
					end <- true
				}()
				<-end

				defer res.Body.Close()
				rData := mooze.req.ResBody(res)
				mooze.req.resData = mooze.req.Prettify(rData)
				wRes.Content(mooze.req.resData)

				mooze.renderLayout(wReq, wRes, wStatus)
				mooze.req.resStatus = res.Status
				mooze.req.resCode = res.StatusCode
				mooze.ms.Show()
			}
		}
	}
	mooze.ms.Exit(0)
}
