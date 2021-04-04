package v2

import (
	"strconv"
)

type mooze struct {
	screen       *Screen
	terminalUnix *Terminal
	request      *Request
}

func NewMooze(s *Screen, t Terminal, r *Request) *mooze {
	return &mooze{s, &t, r}
}

func (m *mooze) Init() {
	m.screen.Println("\033[35mMooze: Yet another REST api test tool for command-line users\033[0m\r")
	m.screen.Println("request \033[31mu\033[0mrl: " + m.request.Url + "\r")
	m.screen.Println("request \033[31mm\033[0method: " + m.request.Method + "\r")
	m.screen.Println("request \033[31mb\033[0mody: " + m.request.Body + "\r")
	m.screen.Println("\033[31ms\033[0mend request\r")
	m.screen.Println("\033[31mq\033[0muit\r")
	m.screen.Println("\r")
}

func (m *mooze) Refresh() {
	m.screen.ClearScreen()
	m.Init()
}

// OpenSelection takes string slice and returns
// selected item's index
func (m *mooze) OpenSelection(opt []string) int {
	t := NewTerminalUnix()
	idx := 0
	l := len(opt)
	m.screen.HideCursor()

SELECT:
	for {
		m.Refresh()
		for i, s := range opt {
			if i == idx {
				m.screen.Println(FgBlack(BgGreen("\r" + s)))
			} else {
				m.screen.Println("\r" + s)
			}
		}
		m.screen.Print("\r")

		buf := make([]byte, 10)
		t.Read(buf)

		switch string(buf[0]) {
		case "h":
			fallthrough
		case "j":
			if idx+1 < l {
				idx += 1
			}

		case "l":
			fallthrough
		case "k":
			if idx > 0 {
				idx -= 1
			}

		// enter key
		case string([]byte{13}):
			fallthrough
		case "q":
			m.Refresh()
			break SELECT
		}
	}
	m.screen.ShowCursor()
	return idx
}

func Run() {
	// initializing components
	s := NewScreen()
	t := NewTerminalUnix()
	r := NewRequest()
	mooze := NewMooze(s, t, r)

	t.MakeRaw()
	defer t.RestoreRaw()

	s.LoadAlternateScreen()
	s.ClearScreen()

	mooze.Init()

CORE:
	for {
		buf := make([]byte, 10)
		t.Read(buf)

		// s.ClearLine()

		if buf[0] == 113 {
			s.Println("\n\033[31mTerminating...\r\033[0m")
			break CORE
		}

		// s.Print("Input:", string(buf), "\r")
		s.Print(string(buf))

		switch string(buf[0]) {
		// refresh
		case "r":
			mooze.Refresh()

		// enter url input mode
		case "u":
			s.ClearLine()
			s.Println("\r\033[31mInput target url\r\033[0m")
			// s.Print("\r\033[31m> \033[0m")
			strBuf, err := t.ReadStringTyped("url")
			if err != nil {
				s.Print(err, "\r")
			}
			r.Url = strBuf
			mooze.Refresh()

		// enter request body input mode
		case "b":
			s.ClearLine()
			s.Println("\r\033[31mInput request body\r\033[0m")
			// s.Print("\r\033[31m> \033[0m")
			strBuf, err := t.ReadStringTyped("body")
			if err != nil {
				s.Print(err, "\r")
			}
			r.Body = strBuf
			mooze.Refresh()

		// enter method input mode
		case "m":
			s.ClearLine()
			s.Println("\r\033[31mInput request method\r\033[0m")
			// s.Print("\r\033[31m> \033[0m")
			strBuf, err := t.ReadStringTyped("method")
			if err != nil {
				s.Print(err, "\r")
			}
			r.Method = strBuf
			mooze.Refresh()

		// send request
		case "s":
			s.ClearLine()
			s.Println("\r\033[31mRequest Sent\n\r\033[0m")
			res, code, err := r.Send()
			// must call RestoreRaw to see not messy response
			t.RestoreRaw()
			s.Println("\r\033[31mGot Response:\r\033[0m")
			if err != nil {
				s.Print(err, "\r")
			}
			s.Println(r.Json(res))
			codeNumber := code[0:3]
			switch i, _ := strconv.Atoi(codeNumber); i / 100 {
			case 2:
				s.Println(FgBlack(BgGreen(code)))
				break
			case 4, 5:
				s.Println(FgBlack(BgRed(code)))
				break
			default:
				s.Println(FgBlack(BgBlue(code)))
			}
			s.Print("\n\r")
			t.MakeRaw()

		case "t":
			opts := []string{
				"test1",
				"test2",
				"test3",
				"test4",
			}
			s.Println(opts[mooze.OpenSelection(opts)] + "\r")
		}
	}

	s.UnloadAlternateScreen()
}
