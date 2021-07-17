package ui

import (
	"fmt"
	"strings"
)

var WindowStore = []*Window{}

type Window struct {
	x             int
	y             int
	w             int
	h             int
	title         string
	content       []string
	isTransparent bool

	frameVertical    string
	frameHorizontal  string
	frameTopLeft     string
	frameTopRight    string
	frameBottomLeft  string
	frameBottomRight string

	Meta *WindowMeta /* window metadata */
}

// WindowMeta includes information about window manipulation
type WindowMeta struct {
	focusable  bool /* if window is focusable */
	focused    bool /* if window is focused */
	selectable bool /* if window has selectable items */
	cursor     int  /* index of cursor */
	page       int  /* current page number if content is paged */
}

// x: vertical coordinate of window top-left point
// y: horizontal coordinate of window top-left point
// h: window's vertical length (including frame)
// w: window's horizontal length (including frame)
func NewWindow(x, y, h, w int) *Window {
	win := &Window{
		x:       x,
		y:       y,
		w:       w - 1,
		h:       h - 1,
		title:   "",
		content: []string{"test"},

		frameVertical:    "\u2500",
		frameHorizontal:  "\u2502",
		frameTopLeft:     "\u250c",
		frameTopRight:    "\u2510",
		frameBottomLeft:  "\u2514",
		frameBottomRight: "\u2518",

		Meta: &WindowMeta{
			focusable:  false,
			focused:    false,
			selectable: false,
			cursor:     0, /* 0-based index */
			page:       0,
		},
	}

	WindowStore = append(WindowStore, win)
	return win
}

func (w *Window) Title(title string) *Window {
	w.title = title
	return w
}

func (w *Window) Content(content []string) *Window {
	w.content = content
	return w
}

func (w *Window) Fill() {
	MoveCursorTo(w.x+1, w.y+1)
	for i, con := range w.content {
		if i+1 >= w.h {
			MoveCursorTo(w.x+i, w.y+1)
			fmt.Print(strings.Repeat(" ", w.w-1))

			// replace this to display page number
			// and others when pagination ready
			MoveCursorTo(w.x+i, w.y+1)
			fmt.Print("...")
			break
		}

		if len(con) >= w.w {
			fmt.Print(con[0 : w.w-3])
			fmt.Print("..")
		} else {
			fmt.Print(con)
		}

		MoveCursorTo(w.x+i+2, w.y+1)
	}
}

func (w *Window) Render() {
	// upper horizontal
	MoveCursorTo(w.x, w.y)
	for i := 0; i < w.w; i++ {
		fmt.Print(w.frameVertical)
	}

	// lower horizontal
	MoveCursorTo(w.x+w.h, w.y)
	for i := 0; i < w.w; i++ {
		fmt.Print(w.frameVertical)
	}

	// left vertical
	MoveCursorTo(w.x, w.y)
	for i := 0; i <= w.h; i++ {
		fmt.Print(w.frameHorizontal)
		MoveCursorTo(w.x+i, w.y)
	}

	// right vertical
	MoveCursorTo(w.x, w.y+w.w)
	for i := 0; i <= w.h; i++ {
		fmt.Print(w.frameHorizontal)
		MoveCursorTo(w.x+i, w.y+w.w)
	}

	// top left
	MoveCursorTo(w.x, w.y)
	fmt.Print(w.frameTopLeft)
	// top right
	MoveCursorTo(w.x, w.y+w.w)
	fmt.Print(w.frameTopRight)
	// bottom left
	MoveCursorTo(w.x+w.h, w.y)
	fmt.Print(w.frameBottomLeft)
	// bottom right
	MoveCursorTo(w.x+w.h, w.y+w.w)
	fmt.Print(w.frameBottomRight)

	// Render Title
	if w.title != "" && len(w.title) < w.w {
		MoveCursorTo(w.x, w.y+1)
		fmt.Print(w.title)
	}

	// Render Content
	w.Fill()
}
