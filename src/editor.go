package mooze

import (
	"github.com/gdamore/tcell"
)

type MoozeEditor struct {
	CursorX int
	CursorY int

	Content string
}

func NewMoozeEditor() *MoozeEditor {
	return &MoozeEditor{0, 0, ""}
}

func (e *MoozeEditor) readLine(m *mooze) string {
	style := ToStyle("black", "green")
	width, height := m.ms.Size()
	w := NewMoozeWindow(5, 5, height-10, width-10, false)

	w.Content([]string{"Ctrl + Enter to finish editor",
		"(fn + Enter on macOS)"})
	m.ms.RenderWindow(w, style)
	m.ms.Show()
	w.Content([]string{})

	e.Content = ""
	e.CursorX = 6
	e.CursorY = 5
	buf := ""
	line := 0

	for {
		if buf == "" {
			m.ms.RenderWindow(w, style)
		}
		m.ms.PrintInsideWindow(
			w, e.CursorX, e.CursorY+2, "_", style,
		)
		ev := m.ms.EmitEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			m.ms.Reload()

		case *tcell.EventKey:
			r := ev.Rune()
			switch r {
			// Ctrl + Enter exits editor
			case 10:
				w.ContentAppend([]string{buf})
				return e.Content

			// enter
			case 13:
				w.ContentAppend([]string{buf})
				buf = ""
				e.Content += "\n"

				line += 1
				e.CursorX += 1
				e.CursorY = 5
				m.ms.RenderWindow(w, style)
				m.ms.PrintInsideWindow(w, e.CursorX, e.CursorY+1, "_", style)
				m.ms.Show()

			// backspace
			case 127:
				if len(buf) > 0 {
					e.Content = e.Content[0 : len(e.Content)-1]
					buf = buf[0 : len(buf)-1]
					m.ms.PrintInsideWindow(
						w, e.CursorX, e.CursorY, "_", style,
					)
					m.ms.PrintInsideWindow(
						w, e.CursorX, e.CursorY+1, " ", style,
					)
					m.ms.PrintInsideWindow(
						w, e.CursorX, e.CursorY+2, " ", style,
					)
					e.CursorY -= 1
				}
				m.ms.Show()

			// normal key input
			default:
				buf += string(r)
				e.Content += string(r)
				e.CursorY += 1
				m.ms.PrintInsideWindow(
					w, e.CursorX, e.CursorY, string(r), style,
				)
				m.ms.Show()
			}
		}
	}
}
