package v2

import (
	// "fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type TerminalUnix struct {
	In     *os.File
	State  *terminal.State
	Prompt *terminal.Terminal
}

func NewTerminalUnix() *TerminalUnix {
	return &TerminalUnix{
		In: openTty(),
		Prompt: terminal.NewTerminal(
			struct {
				io.Reader
				io.Writer
			}{os.Stdin, os.Stdout}, "",
		),
	}
}

func openTty() *os.File {
	in, err := os.OpenFile("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	return in
}

func (t *TerminalUnix) MakeRaw() {
	state, err := terminal.MakeRaw(int(t.In.Fd()))
	if err != nil {
		panic(err)
	}
	t.State = state
}

func (t *TerminalUnix) RestoreRaw() {
	terminal.Restore(int(t.In.Fd()), t.State)
}

func (t *TerminalUnix) MakeNonblock() {
	err := syscall.SetNonblock(int(t.In.Fd()), true)
	if err != nil {
		panic(err)
	}
}

func (t *TerminalUnix) RestoreNonblock() {
}

func (t *TerminalUnix) Read(buf []byte) []byte {
	syscall.Read(int(t.In.Fd()), buf)
	return buf
}

func (t *TerminalUnix) ReadString() (string, error) {
	// var s string
	// _, err := fmt.Scanln(&s)
	// if err != nil {
	// 	panic(err)
	// }
	// return s

	line, err := t.Prompt.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}
