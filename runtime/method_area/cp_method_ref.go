package method_area

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
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
// return (method, java-error)
func (r *MethodRef) ResolvedMethod() (resolvedMethod *Method, javaErr *heap.Object) {
	if r.method == nil { // using cache, lazy load.
		err := r.resolveMethodRef()
		if err != nil {
			return nil, err
		}
	}
	return r.method, nil
}

func (r *MethodRef) resolveMethodRef() *heap.Object {
	// 1. parse class (load class if not loaded)
	class := r.ResolvedClass()

	// 2. check is interface or not.
	if class.IsInterface() {
		errClass := class.Loader().LoadClass("java/lang/IncompatibleClassChangeError", false)
		if errClass == nil {
			panic("java.lang.IncompatibleClassChangeError")
		} else {
			return heap.NewExceptionObject(errClass, fmt.Sprintf("class %s is a interface", class.name))
		}

	}

	// 3. find method
	method := lookupMethod(class, r.name, r.descriptor)
	if method == nil {
		fmt.Printf("@@ DEBUG - resolveMethodRef failed, jClass = %s \n", r.class.jClass)
		errClass := class.Loader().LoadClass("java/lang/NoSuchMethodError", false)
		if errClass == nil {
			panic("java.lang.NoSuchMethodError: " + r.className + "." + r.name + r.descriptor)
		} else {
			return heap.NewExceptionObject(errClass, fmt.Sprintf("java.lang.NoSuchMethodError: "+r.className+"."+r.name+r.descriptor))
		}

	}

	r.method = method

	return nil
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

// lookupMethodInClass find method in class (helper for Object methods)
func lookupMethodInClass(class *Class, methodName, methodDescriptor string) *Method {
	for c := class; c != nil; c = c.superClass {
		for _, method := range c.Methods() {
			if method.Name() == methodName && method.Descriptor() == methodDescriptor {
				return method
			}
		}
	}

	return nil
}
