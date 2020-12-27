package util

import (
	"unicode/utf8"
)

// returns true if byte slice contains Ascii character
func IsAscii(b []byte) bool {
	if b[1] != 0 {
		return false
	}
	return true
}

// returns first non-zero index in byte slice
func GetLength(b []byte) int {
	for i, v := range b {
		if v == 0 {
			return i
		}
	}
	return len(b)
}

// returns rune data from byte slice
func BytesToRune(b []byte) rune {
	r, n := utf8.DecodeRune(b[0:4])
	if n != 0 {
		return r
	} else {
		panic("decoding failed")
	}
}
