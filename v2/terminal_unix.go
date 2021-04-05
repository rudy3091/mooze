package v2

import (
	// "fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type StdReadWriter struct {
	io.Reader
	io.Writer
}

type TerminalUnix struct {
	In           *os.File
	State        *terminal.State
	Prompt       *terminal.Terminal
	UrlPrompt    *terminal.Terminal
	MethodPrompt *terminal.Terminal
	BodyPrompt   *terminal.Terminal
	HeaderPrompt *terminal.Terminal
}

func NewTerminalUnix() *TerminalUnix {
	return &TerminalUnix{
		In: openTty(),
		Prompt: terminal.NewTerminal(
			StdReadWriter{os.Stdin, os.Stdout},
			FgGreen("> "),
		),
		UrlPrompt: terminal.NewTerminal(
			StdReadWriter{os.Stdin, os.Stdout},
			FgGreen("url: > "),
		),
		MethodPrompt: terminal.NewTerminal(
			StdReadWriter{os.Stdin, os.Stdout},
			FgGreen("method: > "),
		),
		BodyPrompt: terminal.NewTerminal(
			StdReadWriter{os.Stdin, os.Stdout},
			FgGreen("body: > "),
		),
		HeaderPrompt: terminal.NewTerminal(
			StdReadWriter{os.Stdin, os.Stdout},
			FgGreen("header: > "),
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
	line, err := t.Prompt.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}

func (t *TerminalUnix) ReadStringTyped(ts string) (string, error) {
	switch ts {
	case "url":
		return t.ReadUrlString()
	case "method":
		return t.ReadMethodString()
	case "body":
		return t.ReadBodyString()
	case "header":
		return t.ReadHeaderString()
	default:
		return t.ReadString()
	}
}

func (t *TerminalUnix) ReadUrlString() (string, error) {
	line, err := t.UrlPrompt.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}

func (t *TerminalUnix) ReadMethodString() (string, error) {
	line, err := t.MethodPrompt.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}

func (t *TerminalUnix) ReadBodyString() (string, error) {
	line, err := t.BodyPrompt.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}

func (t *TerminalUnix) ReadHeaderString() (string, error) {
	line, err := t.HeaderPrompt.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}
