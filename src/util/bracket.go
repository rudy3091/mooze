package util

import (
	"sync"
)

type bracketChecker struct {
	// S stands for "("
	S int
	// M stands for "{"
	M int
	// L stands for "["
	L int
}

// singleton pattern
var instance *bracketChecker
var mtx sync.Mutex

func BracketChecker() *bracketChecker {
	mtx.Lock()
	defer mtx.Unlock()

	if instance == nil {
		instance = &bracketChecker{0, 0, 0}
	}
	return instance
}

func (b *bracketChecker) Parse(s string) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			b.S += 1
		case ')':
			b.S -= 1
		case '{':
			b.M += 1
		case '}':
			b.M -= 1
		case '[':
			b.L += 1
		case ']':
			b.L -= 1
		}
	}
}
