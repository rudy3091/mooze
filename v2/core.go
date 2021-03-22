package v2

import ()

type mooze struct {
	Url    string
	Method string
}

func Run() {
	mooze := &mooze{}
	s := NewScreen()
	t := NewTerminalUnix()

	t.MakeRaw()

	defer t.RestoreRaw()

	s.LoadAlternateScreen()
	s.Println("mooze v2\r")
	s.Println("\033[H")
	s.Println("u: enter url input mode\r")
	s.Println("m: enter request method input mode\r")
	s.Println("b: enter request body input mode\r")
	s.Println("q: quit\r")

CORE:
	for {
		// s.ReadNumber(&n)
		buf := make([]byte, 10)
		t.Read(buf)

		s.ClearLine()
		s.Print("Input:", string(buf), "\r")

		if buf[0] == 113 {
			s.Println("\n\033[31mTerminating...\r\033[0m")
			break CORE
		}

		s.Print("Input:", string(buf), "\r")

		switch string(buf[0]) {
		case "u":
			s.Println("\033[31mInput target url\r\033[0m")
			t.RestoreRaw()
			t.Read(buf)
			mooze.Url = string(buf)
			t.MakeRaw()

		case "m":
			s.Println("\033[31mInput request method\r\033[0m")
			t.RestoreRaw()
			t.RestoreRaw()
			t.Read(buf)
			mooze.Method = string(buf)
			t.MakeRaw()
		}
	}

	s.UnloadAlternateScreen()
}
