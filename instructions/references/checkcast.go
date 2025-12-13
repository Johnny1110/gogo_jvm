package references

import (
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// ============================================================
// CHECKCAST - object cast
// ============================================================

// opcode = 0xC0
// format: checkcast indexbyte1 indexbyte2
// operands: 2 bytes (constant pool index pointing to ClassRef)
// stack: [..., objectref] → [..., objectref] (no change)

// 1. null won't throw any exception（null can cast to any type）
// 2. if not match IsAssignableFrom() ->  ClassCastException
type CHECKCAST struct {
	base.Index16Instruction
}

func (c *CHECKCAST) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()

	objectref := stack.PopRef() // will be pushed back later!

	if objectref == nil {
		stack.PushRef(nil)
		return
	}

	// get cast target
	rtcp := frame.Method().Class().ConstantPool()
	targetClassRef := rtcp.GetConstant(c.Index).(*method_area.ClassRef)
	targetClass := targetClassRef.ResolvedClass()

	object := objectref.(*heap.Object)

	if !isInstanceOf(object, targetClass) {
		panic(common.NewJavaException("", "java.lang.ClassCastException"))
	}

	// push back
	stack.PushRef(objectref)
}

func (c *CHECKCAST) Opcode() uint8 {
	// opcode = 0xC0
	return opcodes.CHECKCAST
}
