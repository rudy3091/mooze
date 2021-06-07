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

func Parse(a []string) Argv {
	return []Arg{
		{ArgInfo{"name", "desc", 'n', true, true}, "v"},
		{argsAvailable["shell"], "sh"},
	}
}
