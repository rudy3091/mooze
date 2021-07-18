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
	for i := (idx + 1) % l; ; i = (i + 1) % l {
		if WindowStore[i].Meta.focusable {
			WindowStore[i].Meta.focused = true
			break
		}
	}

	ReloadAll()
}
