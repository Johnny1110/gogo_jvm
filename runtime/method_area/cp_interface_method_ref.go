package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

// InterfaceMethodRef interface method ref
type InterfaceMethodRef struct {
	MemberRef
	method *Method
}

// NewInterfaceMethodRef create ref from ClassFile
func NewInterfaceMethodRef(cp *RuntimeConstantPool, refInfo *classfile.ConstantInterfaceMethodRefInfo) *InterfaceMethodRef {
	ref := &InterfaceMethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberRefInfo)
	return ref
}

// ResolvedInterfaceMethod parse
func (r *InterfaceMethodRef) ResolvedInterfaceMethod() *Method {
	if r.method == nil {
		r.resolveInterfaceMethodRef()
	}
	return r.method
}

func (r *InterfaceMethodRef) resolveInterfaceMethodRef() {
	// 1. parse class
	c := r.ResolvedClass()

	// 2. check is interface
	if !c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 3. loop up method
	method := lookupInterfaceMethod(c, r.name, r.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError")
	}

	// TODO: 訪問權限檢查
	r.method = method
}

func (r *InterfaceMethodRef) ResolvedClass() *Class {
	return r.method.Class()
}

// lookupInterfaceMethod find method
func lookupInterfaceMethod(iface *Class, name, descriptor string) *Method {
	// find in current class
	for _, method := range iface.Methods() {
		if method.Name() == name && method.Descriptor() == descriptor {
			return method
		}
	}

	if iface.superClass != nil {
		return lookupInterfaceMethod(iface.superClass, name, descriptor)
	}

	return nil
}
