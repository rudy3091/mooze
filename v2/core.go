package v2

import (
	"strconv"
	"strings"
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
	h := m.request.ParseHeaders()

	m.screen.Println("\033[35mMooze: Yet another REST api test tool for command-line users\033[0m\r")
	m.screen.Println("Request" + "\r")
	m.screen.Println("- " + FgRed("u") + "rl: " + FgBlue(m.request.Url) + "\r")
	m.screen.Println("- " + FgRed("m") + "ethod: " + FgBlue(m.request.Method) + "\r")
	m.screen.Println("- " + FgRed("b") + "ody: " + FgBlue(m.request.Body) + "\r")
	m.screen.Println("- " + FgRed("h") + "eader:\n\r" + FgBlue(h) + "\r")
	m.screen.Println("Operations" + "\r")
	m.screen.Println("- " + FgRed("r") + "efresh screen\r")
	m.screen.Println("- " + FgRed("s") + "end request\r")
	m.screen.Println("- " + FgRed("q") + "uit\r")
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

	m.Refresh()
	for i, s := range opt {
		if i == idx {
			m.screen.Println(FgBlack(BgGreen("\r" + s)))
		} else {
			m.screen.Println("\r" + s)
		}
	}
	m.screen.Print("\r")
	for i := 0; i < l; i++ {
		m.screen.MoveCursorUp()
	}

SELECT:
	for {
		buf := make([]byte, 10)
		t.Read(buf)

		// process inputs
		switch string(buf[0]) {

		// case "h":
		// 	idx = 0

		case "j":
			if idx+1 >= l {
				continue
			}
			m.screen.ClearLine()
			m.screen.Print(opt[idx], "\r")
			idx += 1
			m.screen.MoveCursorDown()
			m.screen.Print(FgBlack(BgGreen("\r"+opt[idx])), "\r")

		// case "l":
		// 	idx = l - 1

		case "k":
			if idx <= 0 {
				continue
			}
			m.screen.ClearLine()
			m.screen.Print(opt[idx], "\r")
			idx -= 1
			m.screen.MoveCursorUp()
			m.screen.Print(FgBlack(BgGreen("\r"+opt[idx])), "\r")

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

		// enter request header input mode
		case "h":
			s.ClearLine()
			s.Println(FgRed("Input request Header"), "\r")
			s.Println(FgRed("format: (key): (value)"), "\r")
			strBuf, err := t.ReadStringTyped("header")
			if err != nil {
				s.Print(err, "\r")
			}
			if strings.Contains(strBuf, ":") {
				splitted := strings.Split(strBuf, ":")
				r.Headers[splitted[0]] = splitted[1]
				s.Println(FgRed("...\n\rHeader Added"), "\n\r")
			} else {
				s.Println(FgRed("...\n\rInvalid format!"), "\n\r")
			}
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

			codeNumber := "100"
			if len(code) > 3 {
				codeNumber = code[0:3]
			}

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
