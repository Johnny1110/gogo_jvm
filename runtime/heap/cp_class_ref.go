package heap

// ClassRef class ref
// new Calculator → need resolve Calculator class
type ClassRef struct {
	SymRef
}

// newClassRef create class ref
func newClassRef(cp *RuntimeConstantPool, className string) *ClassRef {
	ref := &ClassRef{}
	ref.cp = cp
	ref.className = className
	return ref
}
