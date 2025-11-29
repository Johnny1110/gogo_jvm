package heap

import "github.com/Johnny1110/gogo_jvm/classfile"

// MethodRef 方法引用
// 例如：invokestatic Calculator.add → 需要解析 add 方法
type MethodRef struct {
	MemberRef
	method *Method // 解析後的方法（緩存）
}

// newMethodRef 從 ClassFile 創建方法引用
func newMethodRef(cp *RuntimeConstantPool, refInfo *classfile.ConstantMethodRefInfo) *MethodRef {
	ref := &MethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberRefInfo)
	return ref
}

// ResolvedMethod 解析方法引用（這是 invokestatic 的核心！）
func (r *MethodRef) ResolvedMethod() *Method {
	if r.method == nil {
		r.resolveMethodRef()
	}
	return r.method
}

func (r *MethodRef) resolveMethodRef() {
	// 1. 解析類
	c := r.ResolvedClass()

	// 2. 檢查是否是接口
	if c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 3. 查找方法
	method := lookupMethod(c, r.name, r.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError: " + r.className + "." + r.name + r.descriptor)
	}

	r.method = method
}

// lookupMethod 查找方法（包含繼承）
func lookupMethod(c *Class, name, descriptor string) *Method {
	method := lookupMethodInClass(c, name, descriptor)
	if method != nil {
		return method
	}
	// TODO: 在父類中遞歸查找
	return nil
}

func lookupMethodInClass(c *Class, name, descriptor string) *Method {
	for _, method := range c.methods {
		if method.name == name && method.descriptor == descriptor {
			return method
		}
	}
	return nil
}

//```
//
//**這是 `invokestatic` 指令的核心！** 解析流程：
//```
//invokestatic #1  (調用 Calculator.add)
//│
//▼
//MethodRef.ResolvedMethod()
//│
//├─ 1. ResolvedClass() → 加載 Calculator 類
//│
//├─ 2. 檢查不是接口
//│
//└─ 3. lookupMethod("add", "(II)I")
//│
//└─ 在 Calculator.methods 中查找
//│
//▼
//返回 Method 對象
