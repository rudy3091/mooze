package ui

import "fmt"

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
}

func NewWindow() *Window {
	return &Window{
		x:       0,
		y:       0,
		w:       20,
		h:       20,
		title:   "test",
		content: []string{"test"},

		frameVertical:    "\u2500",
		frameHorizontal:  "\u2502",
		frameTopLeft:     "\u250c",
		frameTopRight:    "\u2510",
		frameBottomLeft:  "\u2514",
		frameBottomRight: "\u2518",
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
}
