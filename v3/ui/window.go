package ui

type Window struct {
	x             int
	y             int
	w             int
	h             int
	title         string
	content       []string
	isTransparent bool

	// frameVertical    rune
	// frameHorizontal  rune
	// frameTopLeft     rune
	// frameTopRight    rune
	// frameBottomLeft  rune
	// frameBottomRight rune
}
