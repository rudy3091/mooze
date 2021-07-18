package ui

import (
	"fmt"
	"strconv"
)

func LoadAlternateScreen() {
	fmt.Print("\x1b[?1049h\033[H")
}

func UnloadAlternateScreen() {
	fmt.Print("\x1b[?1049l")
}

// prints everything and append carriage return
func Print(v ...interface{}) {
	fmt.Print(v...)
	fmt.Print("\r")
}

// prints everything and append carriage return and newline
func Println(v ...interface{}) {
	fmt.Print(v...)
	fmt.Print("\r\n")
}

func ShowKeyBindings() {
	main := Fg(
		"Mooze: Yet another REST api test tool for command-line users",
		MAGENTA)
	Println(main)

	Println("Request")
	Println("- " + Fg("u", RED) + "rl: ")
	Println("- " + Fg("m", RED) + "ethod: ")
	Println("- " + Fg("b", RED) + "ody: ")
	Println("- " + Fg("h", RED) + "eader: ")

	Println("Operations")
	Println("- " + Fg("r", RED) + "efresh screen")
	Println("- " + Fg("s", RED) + "end request")
	Println("- " + Fg("q", RED) + "uit")
}

// x: vertical, y: horizontal
func MoveCursorTo(x, y int) {
	_x := strconv.Itoa(x)
	_y := strconv.Itoa(y)
	fmt.Print("\x1b[" + _x + ";" + _y + "H")
}

func HideCursor() {
	fmt.Print("\x1b[?25l")
}

func ShowCursor() {
	fmt.Print("\x1b[?25h")
}
