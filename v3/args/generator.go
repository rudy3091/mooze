package args

type ch chan string

func FromStringSlice(s []string) []string {
	t := []string{}
	for _, str := range s {
		t = append(t, string(str))
	}
	return t
}

func (c ch) Next() string {
	return <-c
}

func Generate(source []string) ch {
	yield := make(chan string)
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
