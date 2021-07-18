package ui

import (
	"fmt"
	"strconv"
)

type ansiColor = int

var keyEscape string = "\x1b["

var reset string = keyEscape + "0m"

const (
	BLACK ansiColor = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

func Fg(s string, c ansiColor) string {
	return keyEscape + "3" + strconv.Itoa(c) + "m" + s + reset
}

func Bg(s string, c ansiColor) string {
	return keyEscape + "4" + strconv.Itoa(c) + s + reset
}

func SetFg(c ansiColor) {
	fmt.Print(keyEscape + "3" + strconv.Itoa(c) + "m")
}

func SetBg(c ansiColor) {
	fmt.Print(keyEscape + "4" + strconv.Itoa(c) + "m")
}
