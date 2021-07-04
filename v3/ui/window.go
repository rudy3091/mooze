package ui

type Window struct {
	x             int
	y             int
	w             int
	h             int
	title         string
	content       []string
	isTransparent bool

	frameVertical    string
	frameHorizontal  string
	frameTopLeft     string
	frameTopRight    string
	frameBottomLeft  string
	frameBottomRight string
}
