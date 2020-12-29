/*
 * user input & ANSI terminal colors
 */
package mooze

import (
	"fmt"
	"strconv"
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
	B input = 98
	H input = 104
	I input = 105
	M input = 109
	Q input = 113
	R input = 114
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
	red     *ColorContext
	blue    *ColorContext
	green   *ColorContext
	yellow  *ColorContext
	cyan    *ColorContext
	orange  *ColorContext
	magenta *ColorContext
	gray    *ColorContext
	black   *ColorContext
	emerald *ColorContext
	purple  *ColorContext
	tomato  *ColorContext
	coral   *ColorContext
}

var FColors = FrontColorEnum{
	red:     NewColorContext("ff0000"),
	blue:    NewColorContext("0000ff"),
	green:   NewColorContext("00ff00"),
	yellow:  NewColorContext("ffff00"),
	cyan:    NewColorContext("00ffff"),
	orange:  NewColorContext("f39c12"),
	magenta: NewColorContext("ff00ff"),
}

type BackColorEnum struct {
	red     *ColorContext
	blue    *ColorContext
	green   *ColorContext
	yellow  *ColorContext
	cyan    *ColorContext
	orange  *ColorContext
	magenta *ColorContext
	gray    *ColorContext
	black   *ColorContext
	emerald *ColorContext
	purple  *ColorContext
	tomato  *ColorContext
	coral   *ColorContext
}

var BColors = BackColorEnum{
	red:     NewColorContext("ff0000"),
	blue:    NewColorContext("0000ff"),
	green:   NewColorContext("00ff00"),
	yellow:  NewColorContext("ffff00"),
	cyan:    NewColorContext("00ffff"),
	orange:  NewColorContext("f39c12"),
	magenta: NewColorContext("ff00ff"),
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
	X               int
	Y               int
	SizeX           int
	SizeY           int
	CharHorizontal  string
	CharVertical    string
	CharTopLeft     string
	CharTopRight    string
	CharBottomLeft  string
	CharBottomRight string
	Focused         bool
	FrameColor      *ColorContext
	Content         []string
}

func NewWindow(x, y, sizeX, sizeY int, hex string) *Window {
	return &Window{
		x, y, sizeX, sizeY,
		"\u2500",
		"\u2502",
		"\u250c",
		"\u2510",
		"\u2514",
		"\u2518",
		false,
		NewColorContext(hex),
		make([]string, x),
	}
}

func (w *Window) Focus() {
	w.Focused = true
	w.CharHorizontal = "\u2501"
	w.CharVertical = "\u2503"
	w.CharTopLeft = "\u250f"
	w.CharTopRight = "\u2513"
	w.CharBottomLeft = "\u2517"
	w.CharBottomRight = "\u251b"
	NewRenderer().RenderWindow(w)
}

func (w *Window) Unfocus() {
	w.Focused = false
	w.CharHorizontal = "\u2500"
	w.CharVertical = "\u2502"
	w.CharTopLeft = "\u250c"
	w.CharTopRight = "\u2510"
	w.CharBottomLeft = "\u2514"
	w.CharBottomRight = "\u2518"
	NewRenderer().RenderWindow(w)
}

type Mode struct {
	Hex  string
	Name string
}

type StatusBar struct {
	Now    *Mode
	Normal *Mode
	Url    *Mode
}

func NewStatusBar() *StatusBar {
	normal := &Mode{"88ee88", "   NORMAL   "}
	url := &Mode{"eeee88", "  URL input  "}
	return &StatusBar{
		Now:    normal,
		Normal: normal,
		Url:    url,
	}
}

var MoozeStatusBar = NewStatusBar()
