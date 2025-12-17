package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// ============================================================
// ATHROW - throw exception
// ============================================================
// opcode: 0xBF (191)
// operands: 無
// stack: ..., objectref(exception)

type ATHROW struct {
	base.NoOperandsInstruction
}

func (a *ATHROW) Execute(frame *runtime.Frame) {
	// 1. pop exception
	exceptionRef := frame.OperandStack().PopRef()

	if exceptionRef == nil {
		exceptionObj := NewNullPointerException(frame)
		ThrowException(frame, exceptionObj)
		return
	}

	exceptionObj := exceptionRef.(*heap.Object)
	ThrowException(frame, exceptionObj)
}

func (a *ATHROW) Opcode() uint8 {
	return opcodes.ATHROW
}
