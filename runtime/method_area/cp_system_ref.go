package method_area

// SymRef cp_ref base type
// for ClassRef, FieldRef, MethodRef
type SymRef struct {
	cp        *RuntimeConstantPool
	className string
	class     *Class
}

// ResolvedClass load class (lazy loading)
// lazy loading: only load 1 time, after that return cached class
func (r *SymRef) ResolvedClass() *Class {
	if r.class == nil {
		r.resolveClassRef() // actual load
	}
	return r.class // return cached class
}

func (r *SymRef) resolveClassRef() {
	// 1. Get class by class's RuntimeConstantPool
	class := r.cp.Class() // main class will be loaded at least
	// using ClassLoader load class
	c := class.Loader().LoadClass(r.className, false)
	// TODO: 訪問權限檢查
	r.class = c
}

// ClassName getter
func (r *SymRef) ClassName() string {
	return r.className
}
