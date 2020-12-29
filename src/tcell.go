package mooze

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type MoozeWindow struct {
	x        int
	y        int
	sizeX    int
	sizeY    int
	useMouse bool
}

var _screen tcell.Screen

func (w *MoozeWindow) initScreen() {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = s.Init(); err != nil {
		panic(err)
	}
	if w.useMouse {
		s.EnableMouse()
	} else {
		s.DisableMouse()
	}
	_screen = s
}

func runeWidth(r rune) int {
	return runewidth.RuneWidth(r)
}
