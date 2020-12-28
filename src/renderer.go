package mooze

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/RudyPark3091/mooze/src/util"
	"golang.org/x/crypto/ssh/terminal"
)

/*
 * col, y
 *   -------------------------->
 * r | Screen
 * o |
 * w |
 * , |
 * x |
 *   V
 */
type Renderer struct {
	tty *os.File

	// row of cursor position
	CursorX int
	// column of cursor position
	CursorY int
}

func NewRenderer() *Renderer {
	r := &Renderer{openTty(), 1, 1}
	return r
}

func (r *Renderer) TtyCol() int {
	w, _, err := terminal.GetSize(int(openTty().Fd()))
	if err != nil {
		panic(err)
	}
	return w
}

func (r *Renderer) TtyRow() int {
	_, h, err := terminal.GetSize(int(openTty().Fd()))
	if err != nil {
		panic(err)
	}
	return h
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

	if r.TtyCol() <= r.CursorY {
		r.CursorX += 1
		r.CursorY += offset
	} else {
		r.CursorY += offset
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
	fmt.Print("\x1B[?25l")
}

func (r *Renderer) ShowCursor() {
	fmt.Print("\x1B[?25h")
}

// move cursor to (x, y): x row & y col
func (r *Renderer) MoveCursorTo(x, y int) {
	fmt.Printf("\x1B[%d;%dH", x, y)
}

func (r *Renderer) MoveCursorLeft() {
	if r.CursorY > 1 {
		r.CursorY -= 1
		r.MoveCursorTo(r.CursorX, r.CursorY)
	}
}

func (r *Renderer) MoveCursorRight() {
	if r.CursorY < r.TtyCol() {
		r.CursorY += 1
		r.MoveCursorTo(r.CursorX, r.CursorY)
	}
}

func (r *Renderer) MoveCursorUp() {
	if r.CursorX > 1 {
		r.CursorX -= 1
		r.MoveCursorTo(r.CursorX, r.CursorY)
	}
}

func (r *Renderer) MoveCursorDown() {
	if r.CursorX < r.TtyRow() {
		r.CursorX += 1
		r.MoveCursorTo(r.CursorX, r.CursorY)
	}
}

func (r *Renderer) ClearLine() {
	fmt.Print("\x1B[2K")
}

func (r *Renderer) Backspace() {
	r.MoveCursorLeft()
	// renders " " (space)
	r.WriteChar([]byte{32, 0, 0, 0})
	r.MoveCursorLeft()
}

func (r *Renderer) RenderTextTo(x, y int, s string, a ...interface{}) {
	r.MoveCursorTo(x, y)
	r.ClearLine()
	fmt.Printf(s, a...)
	r.MoveCursorTo(r.CursorX, r.CursorY)
}

func (r *Renderer) RenderTextNoClear(x, y int, s string, a ...interface{}) {
	r.MoveCursorTo(x, y)
	fmt.Printf(s, a...)
	r.MoveCursorTo(r.CursorX, r.CursorY)
}

func (r *Renderer) RenderWindow(w Window) {
	r.RenderTextNoClear(w.x, w.y, w.charTopLeft)
	r.RenderTextNoClear(w.x+w.sizeX, w.y, w.charBottomLeft)
	r.RenderTextNoClear(w.x, w.y+w.sizeY, w.charTopRight)
	r.RenderTextNoClear(w.x+w.sizeX, w.y+w.sizeY, w.charBottomRight)

	for i := 0; i < w.sizeY-1; i++ {
		r.RenderTextNoClear(w.x, w.y+i+1, w.charHorizontal)
		r.RenderTextNoClear(w.x+w.sizeX, w.y+i+1, w.charHorizontal)
	}

	for i := 0; i < w.sizeX-1; i++ {
		r.RenderTextNoClear(w.x+i+1, w.y, w.charVertical)
		r.RenderTextNoClear(w.x+i+1, w.y+w.sizeY, w.charVertical)
	}
}
