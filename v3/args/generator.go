package args

type T string

type ch chan T

func FromStringSlice(s []string) []T {
	t := []T{}
	for _, str := range s {
		t = append(t, T(str))
	}
	return t
}

func (c ch) Next() T {
	return <-c
}

func Generate(source []T) ch {
	yield := make(chan T)
	idx := 0

	go func() {
		for {
			if idx < len(source) {
				yield <- source[idx]
				idx++
			}
		}
	}()

	return yield
}
