/*
 * stdout text & background colorization
 */

package console

import (
	"fmt"
	"strconv"
)

type Color struct {
	R int
	G int
	B int
}

// for changing text color
// uses RGB color
type ConsoleFg struct {
	TermFgR int
	TermFgG int
	TermFgB int
}

// for changing text background color
type ConsoleBg struct {
	TermBgR int
	TermBgG int
	TermBgB int
}

// for changing text color & background color
type Console struct {
	TermFgR int
	TermFgG int
	TermFgB int

	TermBgR int
	TermBgG int
	TermBgB int
}

// #ffffff to (255, 255, 255)
func HexToColor(h string) Color {
	r, _ := strconv.ParseInt(h[0:2], 16, 0)
	g, _ := strconv.ParseInt(h[2:4], 16, 0)
	b, _ := strconv.ParseInt(h[4:6], 16, 0)
	return Color{int(r), int(g), int(b)}
}

func NewConsoleFg(s string) *ConsoleFg {
	c := HexToColor(s)
	return &ConsoleFg{c.R, c.G, c.B}
}

func NewConsoleBg(s string) *ConsoleBg {
	c := HexToColor(s)
	return &ConsoleBg{c.R, c.G, c.B}
}

func NewConsole(s ...string) *Console {
	fc := HexToColor(s[0])
	bc := HexToColor(s[1])

	return &Console{
		fc.R, fc.G, fc.B,
		bc.R, bc.G, bc.B,
	}
}

func (c ConsoleFg) Println(a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm",
		c.TermFgR, c.TermFgG, c.TermFgB,
	)
	fmt.Print(a...)
	fmt.Println("\x1B[0m")
}

func (c ConsoleBg) Println(a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm",
		c.TermBgR, c.TermBgG, c.TermBgB,
	)
	fmt.Print(a...)
	fmt.Println("\x1B[0m")
}

func (c Console) Println(a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm",
		c.TermFgR, c.TermFgG, c.TermFgB,
		c.TermBgR, c.TermBgG, c.TermBgB,
	)
	fmt.Print(a...)
	fmt.Println("\x1B[0m")
}

func (c ConsoleFg) Print(a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm",
		c.TermFgR, c.TermFgG, c.TermFgB,
	)
	fmt.Print(a...)
	fmt.Print("\x1B[0m")
}

func (c ConsoleBg) Print(a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm",
		c.TermBgR, c.TermBgG, c.TermBgB,
	)
	fmt.Print(a...)
	fmt.Print("\x1B[0m")
}

func (c Console) Print(a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm",
		c.TermFgR, c.TermFgG, c.TermFgB,
		c.TermBgR, c.TermBgG, c.TermBgB,
	)
	fmt.Print(a...)
	fmt.Print("\x1B[0m")
}

func (c ConsoleFg) Printf(s string, a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm",
		c.TermFgR, c.TermFgG, c.TermFgB,
	)
	fmt.Printf(s, a...)
	fmt.Print("\x1B[0m")
}

func (c ConsoleBg) Printf(s string, a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm",
		c.TermBgR, c.TermBgG, c.TermBgB,
	)
	fmt.Printf(s, a...)
	fmt.Print("\x1B[0m")
}

func (c Console) Printf(s string, a ...interface{}) {
	fmt.Printf(
		"\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm",
		c.TermFgR, c.TermFgG, c.TermFgB,
		c.TermBgR, c.TermBgG, c.TermBgB,
	)
	fmt.Printf(s, a...)
	fmt.Print("\x1B[0m")
}

func (c ConsoleFg) Colorize(s string) string {
	return "\x1B[38;2;" +
		string(c.TermFgR) + ";" +
		string(c.TermFgG) + ";" +
		string(c.TermFgB) + "m" +
		s + "\x1B[0m"
}

func (c ConsoleBg) Colorize(s string) string {
	return "\x1B[48;2;" +
		string(c.TermBgR) + ";" +
		string(c.TermBgG) + ";" +
		string(c.TermBgB) + "m" +
		s + "\x1B[0m"
}

func (c Console) Colorize(s string) string {
	return "\x1B[38;2;" +
		string(c.TermFgR) + ";" +
		string(c.TermFgG) + ";" +
		string(c.TermFgB) + "m" +
		"\x1B[48;2;" +
		string(c.TermBgR) + ";" +
		string(c.TermBgG) + ";" +
		string(c.TermBgB) + "m" +
		s + "\x1B[0m"
}
