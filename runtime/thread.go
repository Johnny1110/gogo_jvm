package runtime

type Thread struct {
	pc int
}

func (t *Thread) PC() int {
	return t.pc
}
