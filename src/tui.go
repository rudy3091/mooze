/*
 * user input & ANSI terminal colors
 */
package mooze

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

type input rune

const (
	CTRLA     input = 1
	CTRLJ     input = 10
	CTRLQ     input = 17
	CTRLS     input = 19
	CTRLENTER input = 10

	TAB input = 9
	// CTRL M == ENTER on UNIX
	ENTER     input = 13
	ESCAPE    input = 27
	BACKSPACE input = 127

	SHIFTA input = 65
	SHIFTU input = 85

	A input = 97
	I input = 105
	Q input = 113
	U input = 117
)

type Color struct {
	R int
	G int
	B int
}

type textEffect int

const (
	NORMAL     textEffect = 0
	BOLD       textEffect = 1
	DIM        textEffect = 2
	UNDERLINED textEffect = 4
	BLINK      textEffect = 5
)

// for changing text color
// uses RGB color
type Fg struct {
	color  Color
	effect textEffect
}

// for changing text background color
type Bg struct {
	color  Color
	effect textEffect
}

// for changing text color & background color
type ColorContext struct {
	fg Fg
	bg Bg
}

type FrontColorEnum struct {
	red     ColorContext
	blue    ColorContext
	green   ColorContext
	yellow  ColorContext
	cyan    ColorContext
	orange  ColorContext
	magenta ColorContext
	gray    ColorContext
	black   ColorContext
	emerald ColorContext
	purple  ColorContext
	tomato  ColorContext
	coral   ColorContext
}

// ffffff to Color{255, 255, 255}
func HexToColor(h string) Color {
	r, _ := strconv.ParseInt(h[0:2], 16, 0)
	g, _ := strconv.ParseInt(h[2:4], 16, 0)
	b, _ := strconv.ParseInt(h[4:6], 16, 0)
	return Color{int(r), int(g), int(b)}
}

func UndefinedColor() Color {
	return Color{-1, -1, -1}
}

func NewColorContext(s ...string) *ColorContext {
	var fg Fg
	var bg Bg

	fg = Fg{HexToColor(s[0]), 0}
	if len(s) == 1 {
		bg = Bg{UndefinedColor(), 0}
	} else {
		bg = Bg{HexToColor(s[1]), 0}
	}

	return &ColorContext{fg, bg}
}

func (c ColorContext) HasBg() bool {
	if c.bg.color != UndefinedColor() {
		return true
	}
	return false
}

func (c ColorContext) Println(a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm",
		c.fg.color.R, c.fg.color.G, c.fg.color.B,
		c.bg.color.R, c.bg.color.G, c.bg.color.B,
	)
	fmt.Print(a...)
	fmt.Println("\x1B[0m")
}

func (c ColorContext) Print(a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm",
		c.fg.color.R, c.fg.color.G, c.fg.color.B,
		c.bg.color.R, c.bg.color.G, c.bg.color.B,
	)
	fmt.Print(a...)
	fmt.Print("\x1B[0m")
}

func (c ColorContext) Printf(s string, a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm",
		c.fg.color.R, c.fg.color.G, c.fg.color.B,
		c.bg.color.R, c.bg.color.G, c.bg.color.B,
	)
	fmt.Printf(s, a...)
	fmt.Print("\x1B[0m")
}

func (c ColorContext) Colorize(s string) string {
	ret := "\x1B[38;2;" +
		strconv.Itoa(c.fg.color.R) + ";" +
		strconv.Itoa(c.fg.color.G) + ";" +
		strconv.Itoa(c.fg.color.B) + "m"

	if c.HasBg() {
		ret += "\x1B[48;2;" +
			strconv.Itoa(c.bg.color.R) + ";" +
			strconv.Itoa(c.bg.color.G) + ";" +
			strconv.Itoa(c.bg.color.B) + "m"
	}
	ret += s + "\x1B[0m"
	return ret
}

type Window struct {
	x               int
	y               int
	sizeX           int
	sizeY           int
	charHorizontal  string
	charVertical    string
	charTopLeft     string
	charTopRight    string
	charBottomLeft  string
	charBottomRight string
}

func NewWindow(x, y, sizeX, sizeY int) Window {
	return Window{x, y, sizeX, sizeY, "\u2500", "\u2502", "\u250c", "\u2510", "\u2514", "\u2518"}
}

// returns file descriptor of /dev/tty
func openTty() *os.File {
	tty, err := os.OpenFile("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	return tty
}
