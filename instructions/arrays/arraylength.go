package arrays

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// ARRAYLENGTH opcode: 0xBE
// access array length
// Stack: [..., arrayref] → [..., length]
type ARRAYLENGTH struct {
	base.NoOperandsInstruction
}

func (a *ARRAYLENGTH) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arrRef := stack.PopRef()

	// check null
	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	arrLen := arr.ArrayLength()

	// push len into stack
	stack.PushInt(arrLen)
}

func (a *ARRAYLENGTH) Opcode() uint8 {
	return opcodes.ARRAYLENGTH
}
