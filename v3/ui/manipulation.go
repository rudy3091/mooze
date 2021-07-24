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

func NextItem() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print("something gone wrong")
	}

	w := WindowStore[idx]
	w.Meta.cursor += 1
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
	w.Meta.cursor -= 1
	w.Render()
}
