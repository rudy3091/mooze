package ui

import (
	"fmt"
	"strconv"
	"strings"
)

var WindowStore = []*Window{}
var windowIndex uint64 = 0

type Window struct {
	id            uint64
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
	focusable       bool /* if window is focusable */
	focused         bool /* if window is focused */
	selectable      bool /* if window has selectable items */
	cursor          int  /* row index of cursor */
	page            int  /* current page number if content is paged */
	horizontalIndex int  /* content's horizontal scroll amount */
	onSelect        []func()
}

// x: vertical coordinate of window top-left point
// y: horizontal coordinate of window top-left point
// h: window's vertical length (including frame)
// w: window's horizontal length (including frame)
func NewWindow(x, y, h, w int) *Window {
	windowIndex += 1
	win := &Window{
		id:      windowIndex,
		x:       x,
		y:       y,
		w:       w - 1,
		h:       h - 1,
		title:   "",
		content: []string{},

		frameVertical:    "\u2500",
		frameHorizontal:  "\u2502",
		frameTopLeft:     "\u250c",
		frameTopRight:    "\u2510",
		frameBottomLeft:  "\u2514",
		frameBottomRight: "\u2518",

		Meta: &WindowMeta{
			focusable:       true,
			focused:         false,
			selectable:      true,
			cursor:          0, /* 0-based index */
			page:            0,
			horizontalIndex: 0,
		},
	}

	WindowStore = append(WindowStore, win)
	return win
}

func NewRoundWindow(x, y, h, w int) *Window {
	win := NewWindow(x, y, h, w)
	win.frameTopLeft = "\u256d"
	win.frameTopRight = "\u256e"
	win.frameBottomLeft = "\u2570"
	win.frameBottomRight = "\u256f"

	WindowStore = append(WindowStore, win)
	return win
}

func (w *Window) ClearWithFrame() {
	for i := 0; i < w.h+1; i++ {
		MoveCursorTo(w.x+i, w.y)
		for j := 0; j < w.w+1; j++ {
			fmt.Print(" ")
		}
	}
}

func (w *Window) Clear() {
	for i := 1; i < w.h; i++ {
		MoveCursorTo(w.x+i, w.y+1)
		for j := 1; j < w.w; j++ {
			fmt.Print(" ")
		}
	}
}

func (w *Window) Close() {
	for i, win := range WindowStore {
		if win.id == w.id {
			WindowStore = append(WindowStore[:i], WindowStore[i+1:]...)
			w.ClearWithFrame()
			return
		}
	}
}

func (w *Window) OnSelect(fns []func()) *Window {
	w.Meta.onSelect = fns
	return w
}

func Select() {
	idx, _ := getCurrentFocus()
	w := WindowStore[idx]
	if !w.Meta.selectable || w.Meta.cursor >= len(w.Meta.onSelect) {
		return
	}

	fn := w.Meta.page*(w.h-2) + w.Meta.cursor
	w.Meta.onSelect[fn]()
}

func (w *Window) Title(title string) *Window {
	w.title = title
	return w
}

func (w *Window) Content(content []string) *Window {
	w.content = content
	return w
}

func (w *Window) Append(line string) *Window {
	w.content = append(w.content, line)
	return w
}

func (w *Window) Focus() *Window {
	if !w.Meta.focusable {
		return w
	}

	for _, win := range WindowStore {
		win.Meta.focused = false
	}
	w.Meta.focused = true
	return w
}

// returns absolute cursor position index
func (w *Window) getOverallCursorIndex() int {
	return (w.h-1)*w.Meta.page + w.Meta.cursor
}

func (w *Window) getPages() int {
	if w.h != 2 {
		return len(w.content)/(w.h-1) + 1
	} else {
		return len(w.content) / (w.h - 1)
	}
}

func (w *Window) Disable() *Window {
	w.Meta.focusable = false
	return w
}

func (w *Window) Enable() *Window {
	w.Meta.focusable = true
	return w
}

func (w *Window) Resize(height, width int) *Window {
	w.h = height - 1
	w.w = width - 1
	return w
}

func (w *Window) Relocate(x, y int) *Window {
	w.x = x
	w.y = y
	return w
}

func (w *Window) Fill() {
	MoveCursorTo(w.x+1, w.y+1)
	pageStart := w.Meta.page * (w.h - 1)

	for i, con := range w.content[pageStart:] {
		if i+1 >= w.h {
			// MoveCursorTo(w.x+i, w.y+1)
			// fmt.Print(strings.Repeat(" ", w.w-1))
			break
		}

		MoveCursorTo(w.x+i+1, w.y+1)
		if w.Meta.focused && w.Meta.cursor == i {
			SetFg(BLACK)
			SetBg(CYAN)
		} else if w.Meta.focused {
			fmt.Print(reset)
			SetFg(CYAN)
		}

		// // former version without horizontal index
		// if len(con) >= w.w {
		// 	fmt.Print(con[0 : w.w-3])
		// 	fmt.Print("..")
		// } else {
		// 	// pad whitespaces
		// 	fmt.Print(con + strings.Repeat(" ", w.w-1-len(con)))
		// }

		// window's horizontal scroll amount
		hi := w.Meta.horizontalIndex
		if len(con)-hi >= w.w {
			fmt.Print(con[hi : w.w-3+hi])
			fmt.Print("..")
		} else if hi >= len(con) {
			fmt.Print(strings.Repeat(" ", w.w-1))
		} else {
			// pad whitespaces
			fmt.Print(con[hi:] + strings.Repeat(" ", w.w-1-len(con)+hi))
		}

		MoveCursorTo(w.x+i+2, w.y+1)
	}
}

func (w *Window) Render() {
	if w.Meta.focused {
		SetFg(CYAN)
	}
	defer fmt.Print(reset)

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

	pageNumToStr := strconv.Itoa(w.Meta.page + 1)
	MoveCursorTo(w.x+w.h, w.y+w.w-len(pageNumToStr)-5)
	fmt.Println(pageNumToStr + " of " +
		strconv.Itoa(w.getPages()))

	// Render Content
	w.Fill()
}

func ReloadAll() {
	for _, w := range WindowStore {
		w.Render()
	}
}
