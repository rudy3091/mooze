package mooze

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Renderer struct {
	TtyCol int
	TtyRow int

	// column of cursor position
	CursorX int
	// row of cursor position
	CursorY int
}

func NewRenderer() *Renderer {
	r := &Renderer{0, 0, 1, 1}
	w, h, err := terminal.GetSize(int(openTty().Fd()))
	if err != nil {
		panic(err)
	}

	r.TtyCol = w
	r.TtyRow = h
	return r
}

func (r *Renderer) ReadChar(fd *os.File, buf []byte) (int, error) {
	return syscall.Read(int(fd.Fd()), buf)
}

func (r *Renderer) WriteChar(buf []byte) {
	fmt.Fprint(os.Stdout, string(buf[0]))
	if r.TtyCol <= r.CursorX {
		r.CursorY += 1
		r.CursorX = 1
	} else {
		r.CursorX += 1
	}
}

func (r *Renderer) ToRawMode(fd *os.File) *terminal.State {
	state, err := terminal.MakeRaw(int(fd.Fd()))
	if err != nil {
		panic(err)
	}
	return state
}

func (r *Renderer) RestoreState(fd *os.File, s *terminal.State) {
	err := terminal.Restore(int(fd.Fd()), s)
	if err != nil {
		panic(err)
	}
}

func (r *Renderer) ClearConsoleUnix() {
	// for UNIX machine
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (r *Renderer) UseNonblockIo(fd *os.File, b bool) {
	err := syscall.SetNonblock(int(fd.Fd()), b)
	if err != nil {
		panic(err)
	}
}

func (r *Renderer) HideCursor() {
	fmt.Print("\\e[?25l")
}

func (r *Renderer) ShowCursor() {
	fmt.Print("\\e[?25h")
}

func (r *Renderer) TargetStdout(x, y int, s string, a ...interface{}) {
	fmt.Printf("\x1B[%d;%dH", x, y)
	fmt.Print("\x1B[2K")
	fmt.Printf(s, a...)
	fmt.Printf("\x1B[%d;%dH", r.CursorY, r.CursorX)
}
