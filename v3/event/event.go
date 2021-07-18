package event

import (
	"syscall"
)

type Key int

const (
	_ = iota

	KeyU
	KeyM
	KeyB
	KeyH

	KeyR
	KeyS
	KeyQ
)

type Event struct {
	T   Key
	Buf []byte
}

type Buffer []byte

// length returns index of last non-zero
// element of the buffer
func (b Buffer) length() int {
	i := len(b) - 1
	for i >= 0 {
		if b[i] != 0 {
			return i
		}
	}
	return 0 /* empty buffer */
}

// isEsc checks the buffer indicates ESC key input
func (b Buffer) isEsc() bool {
	if b.length() == 1 && b[0] == 27 {
		return true
	} else {
		return false
	}
}

// hasEscape checks the buffer has escape code
// e.g. true for ArrowUp ([27 91 65])
func (b Buffer) hasEscape() bool {
	if len(b) != 0 && b[0] == 27 {
		return true
	} else {
		return false
	}
}

func EmitKeyEvent(fd int, ch chan Event) {
	for {
		buf := make([]byte, 10)
		syscall.Read(fd, buf)
		ch <- Event{
			T:   KeyS,
			Buf: buf,
		}
	}
}
