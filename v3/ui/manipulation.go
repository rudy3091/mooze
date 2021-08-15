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
		fmt.Print(err)
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

func (w *Window) focusNext() {
	// if next item index exceeds window boundary
	// increase page index
	if w.Meta.cursor+1 > w.h-2 && w.Meta.page+1 < w.getPages() {
		w.Meta.cursor = 0
		w.Meta.page += 1
		// should redraw
		w.Clear()
		w.Render()
	} else if w.getOverallCursorIndex()+1 >= len(w.content) {
		// do nothing if there's no more entries
		return
	} else {
		w.Meta.cursor += 1
	}
}

func (w *Window) focusPrev() {
	// if previous item index exceeds window boundary
	// increase page index
	if w.Meta.cursor == 0 && w.Meta.page != 0 {
		w.Meta.cursor = w.h - 2
		w.Meta.page -= 1
		// should redraw
		w.Clear()
		w.Render()
	} else if w.getOverallCursorIndex() == 0 {
		// do nothing if cursor is on first entry
		return
	} else {
		w.Meta.cursor -= 1
	}
}

func (w *Window) focusNextPage() {
	if w.Meta.page < w.getPages()-1 {
		w.Meta.page += 1
		w.Meta.cursor = 0
	}
}

func (w *Window) focusPrevPage() {
	if w.Meta.page != 0 {
		w.Meta.page -= 1
		w.Meta.cursor = 0
	}
}

func NextItem() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print(err)
	}

	w := WindowStore[idx]
	w.focusNext()

	lens := LensWindow
	lens.Content([]string{w.content[w.getOverallCursorIndex()]})
	lens.Render()
	w.Render()
}

func PrevItem() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print(err)
	}

	w := WindowStore[idx]
	w.focusPrev()

	lens := LensWindow
	lens.Content([]string{w.content[w.getOverallCursorIndex()]})
	lens.Render()
	w.Render()
}

func NextPage() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print(err)
	}

	w := WindowStore[idx]
	w.focusNextPage()

	lens := LensWindow
	lens.Content([]string{w.content[w.getOverallCursorIndex()]})
	lens.Render()
	w.Clear()
	w.Render()
}

func PrevPage() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print(err)
	}

	w := WindowStore[idx]
	w.focusPrevPage()

	lens := LensWindow
	lens.Content([]string{w.content[w.getOverallCursorIndex()]})
	lens.Render()
	w.Clear()
	w.Render()
}

func ScrollHalfDown() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print(err)
	}

	w := WindowStore[idx]
	if w.Meta.page == w.getPages()-1 {
		w.Meta.cursor = len(w.content)%(w.h-1) - 1
	} else if w.Meta.cursor <= (w.h-2)/2 {
		w.Meta.cursor += w.h/2 - 1
	} else {
		w.Meta.cursor = w.h - 2
	}

	lens := LensWindow
	lens.Content([]string{w.content[w.getOverallCursorIndex()]})
	lens.Render()
	w.Render()
}

func ScrollHalfUp() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print(err)
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

func ScrollRight() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print(err)
	}

	w := WindowStore[idx]
	if w.Meta.horizontalIndex > 0 {
		w.Meta.horizontalIndex -= 1
	}

	lens := LensWindow
	lens.Content([]string{w.content[w.Meta.cursor]})
	lens.Render()
	w.Render()
}

func ScrollLeft() {
	idx, err := getCurrentFocus()
	if err != nil {
		// handle error
		MoveCursorTo(1, 1)
		fmt.Print(err)
	}

	w := WindowStore[idx]
	w.Meta.horizontalIndex += 1

	lens := LensWindow
	lens.Content([]string{w.content[w.Meta.cursor]})
	lens.Render()
	w.Render()
}
