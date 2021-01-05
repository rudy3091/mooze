package util

import (
	"testing"
)

func TestIsAscii(t *testing.T) {
	// check if _b is true
	check := func(_b bool) {
		if !_b {
			t.Error()
		}
	}

	b := IsAscii([]byte{65, 0, 0, 0})
	check(b == true)
	b = IsAscii([]byte{127, 0, 0, 0})
	check(b == true)
	b = IsAscii([]byte{1, 0, 0, 0})
	check(b == true)
	b = IsAscii([]byte{97, 0, 0, 0})
	check(b == true)
	b = IsAscii([]byte{127, 20, 0, 0})
	check(b == false)
	b = IsAscii([]byte{0, 20, 0, 0})
	check(b == false)
}

func TestGetLength(t *testing.T) {
	check := func(_b bool) {
		if !_b {
			t.Error()
		}
	}

	b := GetLength([]byte{65, 0, 0, 0})
	check(b == 1)
	b = GetLength([]byte{127, 0, 0, 0})
	check(b == 1)
	b = GetLength([]byte{12, 10, 3, 0})
	check(b == 3)
	b = GetLength([]byte{127, 20, 0, 0})
	check(b == 2)
	b = GetLength([]byte{0, 20, 0, 0})
	check(b == 0)
}

func TestBytesToRune(t *testing.T) {
	r := BytesToRune([]byte{65, 0, 0, 0})
	if r != 'A' {
		t.Error()
	}
	r = BytesToRune([]byte{97, 0, 0, 0})
	if r != 'a' {
		t.Error()
	}
	// asian letter with 2 rune width
	r = BytesToRune([]byte{237, 149, 156, 0})
	if r != '한' {
		t.Error()
	}
	r = BytesToRune([]byte{234, 184, 128, 0})
	if r != '글' {
		t.Error()
	}
}
