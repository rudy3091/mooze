package ui

import (
	"io"
	"os"

	"golang.org/x/term"
)

var PopupWindow *Window
var Input = false
var inputTerm = term.NewTerminal(struct {
	io.Reader
	io.Writer
}{os.Stdin, os.Stdout}, "> ")

func OpenPopup(tw, th int) string {
	Input = true
	defer func() { Input = false }()
	defer ClosePopup()

	width := tw / 2
	height := th / 2

	x := th / 4
	y := tw / 4

	// for i := 0; i < height-3; i++ {
	// 	content = append(content, "")
	// }
	// content = append(content, "> _")

	content := []string{}
	PopupWindow = NewWindow(x, y, height, width).
		Title("popup").
		Content(content)

	PopupWindow.Clear()
	PopupWindow.Render()
	PopupWindow.Focus()

	MoveCursorTo(x+height-2, y+1)
	ShowCursor()
	defer HideCursor()

	line, err := inputTerm.ReadLine()
	if err != nil {
		// handle error
	}
	return line
}

func ClosePopup() {
	PopupWindow.Close()
	WindowStore[0].Focus()
	ReloadAll()
}

func OpenPopupSelect(tw, th int) {
	width := tw / 4

	x := th/4 + 5
	y := tw/4 + 10

	PopupWindow = NewWindow(x, y, 10, width).
		Title("select").
		Content([]string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
		})

	PopupWindow.Focus()
	PopupWindow.Clear()
	PopupWindow.Render()
}
