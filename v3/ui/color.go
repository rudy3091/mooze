package ui

type ansiColor = byte

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

func addFgColor(s string, c ansiColor) string {
	return keyEscape + "3" + string(c) + s + reset
}

func addBgColor(s string, c ansiColor) string {
	return keyEscape + "4" + string(c) + s + reset
}
