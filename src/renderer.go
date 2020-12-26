package mooze

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/RudyPark3091/mooze/src/util"
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

func (r *Renderer) HandleResize() {
	w, h, err := terminal.GetSize(int(openTty().Fd()))
	if err != nil {
		panic(err)
	}

	r.TtyCol = w
	r.TtyRow = h
}

func (r *Renderer) ReadChar(fd *os.File, buf []byte) (int, error) {
	return syscall.Read(int(fd.Fd()), buf)
}

func (r *Renderer) WriteChar(buf []byte) {
	fmt.Fprint(os.Stdout, string(util.BytesToRune(buf)))
	offset := 0
	if util.IsAscii(buf) {
		offset = 1
	} else {
		offset = 2
	}
	if r.TtyCol <= r.CursorX {
		r.CursorY += offset
		r.CursorX = offset
	} else {
		r.CursorX += offset
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

func (r *Renderer) MoveCursorTo(x, y int) {
	fmt.Printf("\x1B[%d;%dH", x, y)
}

func (r *Renderer) MoveCursorLeft() {
	if r.CursorX > 2 {
		r.CursorX -= 1
		r.MoveCursorTo(r.CursorY, r.CursorX)
	}
}

func (r *Renderer) MoveCursorRight() {
	if r.CursorX < r.TtyCol {
		r.CursorX += 1
		r.MoveCursorTo(r.CursorY, r.CursorX)
	}
}

func (r *Renderer) MoveCursorUp() {
	if r.CursorY > 2 {
		r.CursorY -= 1
		r.MoveCursorTo(r.CursorY, r.CursorX)
	}
}

func (r *Renderer) MoveCursorDown() {
	if r.CursorY < r.TtyRow {
		r.CursorY += 1
		r.MoveCursorTo(r.CursorY, r.CursorX)
	}
}

func (r *Renderer) ClearLine() {
	fmt.Print("\x1B[2K")
}

func (r *Renderer) RenderTextTo(x, y int, s string, a ...interface{}) {
	r.MoveCursorTo(x, y)
	r.ClearLine()
	fmt.Printf(s, a...)
	r.MoveCursorTo(r.CursorY, r.CursorX)
}
