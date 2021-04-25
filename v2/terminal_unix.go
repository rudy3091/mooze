package v2

import (
	// "fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type StdReadWriter struct {
	io.Reader
	io.Writer
}

type TerminalUnix struct {
	W            int
	H            int
	In           *os.File
	State        *terminal.State
	Prompt       *terminal.Terminal
	UrlPrompt    *terminal.Terminal
	MethodPrompt *terminal.Terminal
	BodyPrompt   *terminal.Terminal
	HeaderPrompt *terminal.Terminal
}

func NewTerminalUnix() *TerminalUnix {
	fd := openTty()
	w, h, err := terminal.GetSize(int(fd.Fd()))
	if err != nil {
		w = 0
		h = 0
	}

	return &TerminalUnix{
		W:  w,
		H:  h,
		In: fd,
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

func (t *TerminalUnix) GetWindowResizeChan() (chan os.Signal, chan bool) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGWINCH)
	return sigs, done
}

func (t *TerminalUnix) HandleResize() {
	w, h, err := terminal.GetSize(int(t.In.Fd()))
	if err != nil {
		w = 0
		h = 0
	}

	t.W = w
	t.H = h
}
