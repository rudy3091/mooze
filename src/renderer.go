package mooze

import (
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Renderer struct {
	TtyCol  int
	TtyRow  int
	CursorX int
	CursorY int
}

func NewRenderer() *Renderer {
	r := &Renderer{0, 0, 0, 0}
	w, h, err := terminal.GetSize(int(openTty().Fd()))
	if err != nil {
		panic(err)
	}

	r.TtyCol = w
	r.TtyRow = h
	return r
}

func (r *Renderer) WriteChar(c []byte, fd *os.File) {
	syscall.Write(int(fd.Fd()), c)
}

func (r *Renderer) ToRawMode(fd *os.File) (*terminal.State, error) {
	state, err := terminal.MakeRaw(int(fd.Fd()))
	if err != nil {
		panic(err)
	}
	return state, err
}

func (r *Renderer) RestoreState(fd *os.File, s *terminal.State) error {
	err := terminal.Restore(int(fd.Fd()), s)
	if err != nil {
		panic(err)
	}
	return err
}

func (r *Renderer) ClearConsoleUnix(fd *os.File) {
	// for UNIX machine
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
