package util

import "testing"

func TestBracketChecker(t *testing.T) {
	b1 := BracketChecker()
	b2 := BracketChecker()
	if b1 != b2 {
		t.Error("Singleton Instance Error")
	}
}

func TestParse(t *testing.T) {
	b := BracketChecker()

	b.Parse("{{{")
	if b.M != 3 {
		t.Error()
	}

	b.Parse("}}}")
	if b.M != 0 {
		t.Error()
	}
	if b.L != 0 || b.S != 0 {
		t.Error()
	}

	b.Parse("[]")
	if b.L != 0 {
		t.Error()
	}

	b.Parse("(()())(())[{}][{{]}}")
	if b.S != 0 || b.M != 0 || b.L != 0 {
		t.Error()
	}
}
