package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// ============================================================
// INSTANCEOF - java instantceof
// ============================================================
// opcode = 0xC1
// format: instanceof indexbyte1 indexbyte2
// operands: 2 bytes (runtime constant pool index pointing to ClassRef)
// stack: [..., objectref] → [..., result (int)]

// Conditions:
// 1. if objectref is null -> push false (0)
// 2. if objectref is T's instance → push true (1)
// 3. otherwise → push false (0)

// is "T's instance" define:
// - T is class: objectref.class == T or T's subclass
// - T is interface: objectref.class implements T
// - T is array type: special handle
type INSTANCEOF struct {
	base.Index16Instruction
}

func (i *INSTANCEOF) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	objectref := stack.PopRef() // objectref

	if objectref == nil {
		stack.PushInt(0) // hit Conditions-1 -> false
		return
	}

	rtcp := frame.Method().Class().ConstantPool()
	TRef := rtcp.GetConstant(i.Index).(*method_area.ClassRef)
	T := TRef.ResolvedClass()

	object := objectref.(*heap.Object)
	if isInstanceOf(object, T) {
		stack.PushInt(1) // hit Conditions-1 -> true
		return
	}

	// otherwise -> false
	stack.PushInt(0)
}

func (i *INSTANCEOF) Opcode() uint8 {
	// 0xC1
	return opcodes.INSTANCEOF
}
