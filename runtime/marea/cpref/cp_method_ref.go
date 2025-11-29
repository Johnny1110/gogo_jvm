package cpref

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/runtime/java"
	"github.com/Johnny1110/gogo_jvm/runtime/marea"
)

// MethodRef Method Reference
// ex: invokestatic Calculator.add -> nedd parse add()
// 1. find Class by name
// 2. find method in class
// 3. catch result
type MethodRef struct {
	MemberRef
	method *java.Method // actual in memory address
}

// newMethodRef make method ref by ClassFile
func newMethodRef(cp *marea.RuntimeConstantPool, refInfo *classfile.ConstantMethodRefInfo) *MethodRef {
	ref := &MethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberRefInfo)
	return ref
}

// ResolvedMethod parse method ref
// core for: invokestatic/invokevirtual
func (r *MethodRef) ResolvedMethod() *java.Method {
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

	// 3. 查找方法
	method := lookupMethod(c, r.name, r.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError: " + r.className + "." + r.name + r.descriptor)
	}

	// TODO: 訪問權限檢查
	r.method = method
}

func (r *MethodRef) ResolvedClass() interface{} {

}
