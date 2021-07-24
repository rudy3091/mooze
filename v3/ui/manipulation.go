package ui

import (
	"errors"
	"fmt"
)

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

	WindowStore[idx].Meta.focused = false
	WindowStore[idx].Render()
	for i := (idx + 1) % l; ; i = (i + 1) % l {
		if WindowStore[i].Meta.focusable {
			WindowStore[i].Meta.focused = true
			WindowStore[i].Render()
			break
		}
	}
}

func NextItem(lens *Window) {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print("something gone wrong")
	}

	w := WindowStore[idx]
	if w.Meta.cursor < len(w.content)-1 {
		w.Meta.cursor += 1
	}
	lens.Content([]string{w.content[w.Meta.cursor]})
	lens.Render()
	w.Render()
}

func PrevItem(lens *Window) {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print("something gone wrong")
	}

	w := WindowStore[idx]
	if w.Meta.cursor > 0 {
		w.Meta.cursor -= 1
	}
	lens.Content([]string{w.content[w.Meta.cursor]})
	lens.Render()
	w.Render()
}
