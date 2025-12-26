package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

// ============================================================
// InterfaceMethodRef - interface method reference
// ============================================================
// Used by invokeinterface instruction
// Similar to MethodRef, but for interface methods
//
// Key differences from MethodRef:
//   - Must reference an interface, not a class
//   - Used with invokeinterface, not invokevirtual
//   - Cannot use vtable optimization

type InterfaceMethodRef struct {
	MemberRef
	method *Method // resolved method (cached)
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

// lookupInterfaceMethod find method in interface and super-interfaces
// This searches:
//  1. Methods declared in the interface itself
//  2. Methods in super-interfaces (recursively)
//  3. Methods in java.lang.Object (interfaces implicitly inherit Object methods)
func lookupInterfaceMethod(iface *Class, name, descriptor string) *Method {
	// Step 1: Search in this interface's methods
	for _, method := range iface.Methods() {
		if method.Name() == name && method.Descriptor() == descriptor {
			return method
		}
	}

	// Step 2: Search in super-interfaces
	for _, superIface := range iface.Interfaces() {
		method := lookupInterfaceMethod(superIface, name, descriptor)
		if method != nil {
			return method
		}
	}

	// Step 3: Search in java.lang.Object
	// Interface methods like toString(), hashCode() come from Object
	if iface.SuperClass() != nil {
		return lookupMethodInClass(iface.SuperClass(), name, descriptor)
	}

	// return nothing
	return nil
}
