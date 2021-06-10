package args

import (
	"os"
)

type ArgInfo struct {
	Name             string
	Description      string
	InShort          byte
	IsValueNeeded    bool
	IsSingleDashOnly bool
}

type Arg struct {
	Info  ArgInfo
	Value string
}

type Argv []Arg

func GetArgs() []string {
	a := os.Args
	return a[1:]
}

func Parse() Argv {
	args := GetArgs()
	gen := Generate(FromStringSlice(GetArgs()))
	result := []Arg{}
	cnt := 0

PARSE:
	for {
		cnt++
		if cnt >= len(args) {
			break PARSE
		}

		next := gen.Next()
		single := next[1] != '-'

		if single {
			result = append(result,
				handleSingleDash(next[1], gen, &cnt))
		} else {
			result = append(result,
				handleDoubleDash(next[2:], gen, &cnt))
		}
	}

	return result
}

func hasValue(flag string) bool {
	return argsAvailable[flag].IsValueNeeded
}

func handleSingleDash(b byte, gen ch, cnt *int) Arg {
	s := string(b)
	if hasValue(s) {
		*cnt++
		return Arg{argsAvailable[s], gen.Next()}
	} else {
		return Arg{argsAvailable[s], ""}
	}
}

func handleDoubleDash(s string, gen ch, cnt *int) Arg {
	if hasValue(s) {
		*cnt++
		return Arg{argsAvailable[s], gen.Next()}
	} else {
		return Arg{argsAvailable[s], ""}
	}
}
