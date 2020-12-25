package mooze

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	// "strconv"
	// "strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type input int

const (
	ESCAPE    input = 27
	Q         input = 113
	BACKSPACE input = 127
)

func openTty() *os.File {
	in, err := os.OpenFile("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	return in
}

func UiInit() {
	in := openTty()

	state, err := terminal.MakeRaw(int(in.Fd()))
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := terminal.NewTerminal(screen, "")
	term.SetPrompt(string(term.Escape.Red) + "> " + string(term.Escape.Reset))

	err = syscall.SetNonblock(int(in.Fd()), true)
	if err != nil {
		panic(err)
	}

	line, err := term.ReadLine()
	if err == io.EOF {
		return
	}
	if line == "c" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	fmt.Fprintln(term, line)

	buf := make([]byte, syscall.SizeofInotifyEvent)
	str := ""
	for {
		syscall.Read(int(in.Fd()), buf)
		str = str + string(buf)
		fmt.Fprintf(term, "buf: %3d ", buf[0])
		fmt.Fprintln(term, "string:", str)

		if buf[0] == byte(Q) || buf[0] == byte(ESCAPE) {
			break
		}
		if buf[0] == byte(BACKSPACE) {
			// backspace
		}
	}

	s := "BLUE"
	NewConsoleFg("ff0000").Println("Red")
	NewConsoleFg("0000ff").Printf("color: %s\n", s)

	err = terminal.Restore(int(in.Fd()), state)
	if err != nil {
		panic(err)
	}
}
