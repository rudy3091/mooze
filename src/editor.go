package mooze

import (
	"github.com/gdamore/tcell"
)

type MoozeEditor struct {
	cursorX int
	cursorY int
}

func NewMoozeEditor() *MoozeEditor {
	return &MoozeEditor{0, 0}
}

func (e *MoozeEditor) readLine(m *mooze) {
	width, height := m.ms.Size()
	w := NewMoozeWindow(5, 5, height-10, width-10, false)
	m.ms.RenderWindow(w, ToStyle("black", "green"))
	m.ms.Show()

LOOP:
	for {
		ev := m.ms.EmitEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			m.ms.Reload()

		case *tcell.EventKey:
			r := ev.Rune()
			// Ctrl + Enter exits editor
			if r == rune(10) {
				break LOOP
			}
			w.Content(append(w.content, string(r)))
			m.ms.RenderWindow(w, ToStyle("black", "green"))
			m.ms.Show()
		}
	}
}
