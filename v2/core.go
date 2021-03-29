package v2

import ()

type mooze struct {
	screen       *Screen
	terminalUnix *Terminal
	request      *Request
}

func NewMooze(s *Screen, t Terminal, r *Request) *mooze {
	return &mooze{s, &t, r}
}

func (m *mooze) Init() {
	m.screen.Println("\033[31mMooze: Yet another REST api test tool for command-line users\033[0m\r")
	m.screen.Println("u: enter url input mode\r")
	m.screen.Println("m: enter request method input mode\r")
	m.screen.Println("b: enter request body input mode\r")
	m.screen.Println("q: quit\r")
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
			s.ClearScreen()
			mooze.Init()

		// enter url input mode
		case "u":
			s.Println("\r\033[31mInput target url\r\033[0m")
			// s.Print("\r\033[31m> \033[0m")
			strBuf, err := t.ReadStringTyped("url")
			if err != nil {
				s.Print(err, "\r")
			}
			r.Url = strBuf

		// enter request body input mode
		case "b":
			s.Println("\r\033[31mInput request body\r\033[0m")
			// s.Print("\r\033[31m> \033[0m")
			strBuf, err := t.ReadStringTyped("body")
			if err != nil {
				s.Print(err, "\r")
			}
			r.Body = strBuf

		// enter method input mode
		case "m":
			s.Println("\r\033[31mInput request method\r\033[0m")
			// s.Print("\r\033[31m> \033[0m")
			strBuf, err := t.ReadStringTyped("method")
			if err != nil {
				s.Print(err, "\r")
			}
			r.Method = strBuf

		// send request
		case "s":
			s.Println("\r\033[31mRequest Sent\n\r\033[0m")
			res, err := r.Send()
			// must call RestoreRaw to see not messy response
			t.RestoreRaw()
			s.Println("\r\033[31mGot Response:\r\033[0m")
			if err != nil {
				s.Print(err, "\r")
			}
			s.Println(r.Json(res))
			s.Print("\n\r")
			t.MakeRaw()
		}
	}

	s.UnloadAlternateScreen()
}
