package heap

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

// MethodRef Method Reference
// ex: invokestatic Calculator.add -> nedd parse add()
// 1. find Class by name
// 2. find method in class
// 3. catch result
type MethodRef struct {
	MemberRef
	method *Method // actual in memory address
}

// newMethodRef make method ref by ClassFile
func NewMethodRef(cp *RuntimeConstantPool, refInfo *classfile.ConstantMethodRefInfo) *MethodRef {
	ref := &MethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberRefInfo)
	return ref
}

// ResolvedMethod parse method ref
// core for: invokestatic/invokevirtual
func (r *MethodRef) ResolvedMethod() *Method {
	if r.method == nil {
		r.resolveMethodRef()
	}
	return r.method
}

func (r *MethodRef) resolveMethodRef() {
	// 1. parse class
	c := r.ResolvedClass()

	// 2. check is interface（interface should using InterfaceMethodRef）
	if c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 3. find method
	method := lookupMethod(c, r.name, r.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError: " + r.className + "." + r.name + r.descriptor)
	}

	// TODO: 訪問權限檢查
	r.method = method
}

func (r *MethodRef) ResolvedClass() *Class {
	return nil
}

// lookupMethod find method（include extend methods）
func lookupMethod(c *Class, name, descriptor string) *Method {
	// find in current class
	method := lookupMethodInClass(c, name, descriptor)
	if method != nil {
		return method
	}

	// TODO: 實現繼承查找 (父類別)
	// c.super

	return nil
}

// lookupMethodInClass find method in target method
func lookupMethodInClass(c *Class, name, descriptor string) *Method {
	for _, method := range c.Methods() {
		if method.Name() == name && method.Descriptor() == descriptor {
			return method
		}
	}
	return nil
}
