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

func (m *mooze) PrintKeymaps() {
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

	mooze.PrintKeymaps()

CORE:
	for {
		buf := make([]byte, 10)
		t.Read(buf)

		s.ClearLine()

		if buf[0] == 113 {
			s.Println("\n\033[31mTerminating...\r\033[0m")
			break CORE
		}

		// s.Print("Input:", string(buf), "\r")
		s.Print(string(buf), "\r")

		switch string(buf[0]) {
		// refresh
		case "r":
			s.ClearScreen()
			mooze.PrintKeymaps()

		// enter url input mode
		case "u":
			str := ""
			s.Println("\033[31mInput target url\r\033[0m")
			s.Print("\033[31m> \033[0m")
			t.RestoreRaw()
			// t.Read(buf)
			t.ReadString(str)
			r.Url = str
			t.MakeRaw()

		// enter method input mode
		case "m":
			s.Println("\033[31mInput request method\r\033[0m")
			s.Print("\033[31m> \033[0m")
			t.RestoreRaw()
			t.RestoreRaw()
			t.Read(buf)
			r.Method = string(buf)
			t.MakeRaw()

		// send request
		case "s":
			s.Println("\033[31mRequest Sent\r\033[0m")
			res := r.Send()
			s.Println(string(res))
			s.Print("\r")
		}
	}

	s.UnloadAlternateScreen()
}
