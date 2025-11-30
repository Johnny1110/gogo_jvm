package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

// MethodRef a pointer, pointing to a method in method area, which can be cached.
type MethodRef struct {
	MemberRef
	method *Method // can cache the method after first time ResolvedMethod()
}

// NewMethodRef create method ref from RuntimeConstantPool
func NewMethodRef(cp *RuntimeConstantPool, refInfo *classfile.ConstantMethodRefInfo) *MethodRef {
	ref := &MethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberRefInfo)
	return ref
}

// ResolvedMethod lazy loading
func (r *MethodRef) ResolvedMethod() *Method {
	if r.method == nil { // using cache, lazy load.
		r.resolveMethodRef()
	}
	return r.method
}

func (r *MethodRef) resolveMethodRef() {
	// 1. parse class (load class if not loaded)
	class := r.ResolvedClass()

	// 2. check is interface or not.
	if class.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 3. find method
	method := lookupMethod(class, r.name, r.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError: " + r.className + "." + r.name + r.descriptor)
	}

	r.method = method
}

// lookupMethod find method (including find in parent class)
func lookupMethod(c *Class, methodName, methodDescriptor string) *Method {
	method := lookupMethodInClass(c, methodName, methodDescriptor)
	if method != nil {
		return method
	}

	if c.superClass != nil {
		return lookupMethod(c.superClass, methodName, methodDescriptor)
	}

	return nil
}

func lookupMethodInClass(c *Class, methodName, methodDescriptor string) *Method {
	for _, method := range c.Methods() {
		if method.Name() == methodName && method.Descriptor() == methodDescriptor {
			return method
		}
	}
	return nil
}
