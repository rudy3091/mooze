package ui

import (
	"errors"
	"fmt"
)

var LensWindow *Window

// TODO: remove this (for debug)
func GetCurrentFocus() (int, error) {
	for i, w := range WindowStore {
		if w.Meta.focused {
			return i, nil
		}
	}
	return -1, errors.New("cannot find focused window")
}

func getCurrentFocus() (int, error) {
	for i, w := range WindowStore {
		if w.Meta.focused {
			return i, nil
		}
	}
	return -1, errors.New("cannot find focused window")
}

func RotateFocus() {
	l := len(WindowStore)

	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print("something gone wrong")
		return
	}

	w := WindowStore[idx]
	w.Meta.focused = false
	w.Render()
	for i := (idx + 1) % l; ; i = (i + 1) % l {
		if WindowStore[i].Meta.focusable {
			WindowStore[i].Meta.focused = true
			WindowStore[i].Render()
			LensWindow.Content(
				[]string{WindowStore[i].content[WindowStore[i].Meta.cursor]},
			).Render()
			break
		}
	}
}

func NextItem() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print("something gone wrong")
	}

	w := WindowStore[idx]
	if w.Meta.page*(w.h-2)+w.Meta.cursor < len(w.content)-1 {
		w.Meta.cursor += 1
	}

	if w.Meta.cursor == w.h-2 {
		w.Meta.cursor = 1
		w.Meta.page += 1
		w.Clear()
	}

	lens := LensWindow
	lens.Content([]string{w.content[w.Meta.cursor]})
	lens.Render()
	w.Render()
}

func PrevItem() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print("something gone wrong")
	}

	w := WindowStore[idx]
	if w.Meta.cursor <= 1 && w.Meta.page != 0 {
		w.Meta.cursor = w.h - 2
		w.Meta.page -= 1
		w.Clear()
	}

	if w.Meta.cursor > 0 {
		w.Meta.cursor -= 1
	}

	lens := LensWindow
	lens.Content([]string{w.content[w.Meta.cursor]})
	lens.Render()
	w.Render()
}

func ScrollHalfDown() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print("something gone wrong")
	}

	w := WindowStore[idx]
	if w.Meta.cursor+w.h/2 < len(w.content) {
		w.Meta.cursor += w.h/2 - 1
	} else {
		w.Meta.cursor = len(w.content) - 1
	}

	lens := LensWindow
	lens.Content([]string{w.content[w.Meta.cursor]})
	lens.Render()
	w.Render()
}

func ScrollHalfUp() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print("something gone wrong")
	}

	w := WindowStore[idx]
	if w.Meta.cursor-w.h/2-1 > 0 {
		w.Meta.cursor -= w.h/2 - 1
	} else {
		w.Meta.cursor = 0
	}

	lens := LensWindow
	lens.Content([]string{w.content[w.Meta.cursor]})
	lens.Render()
	w.Render()
}
