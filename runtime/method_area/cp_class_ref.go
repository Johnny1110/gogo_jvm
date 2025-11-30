package method_area

// ClassRef Class Reference
type ClassRef struct {
	SymRef
}

func NewClassRef(cp *RuntimeConstantPool, className string) *ClassRef {
	ref := &ClassRef{}
	ref.cp = cp
	ref.className = className
	return ref
}
