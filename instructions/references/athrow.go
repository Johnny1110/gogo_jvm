package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
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
	// TODO
}

func (a *ATHROW) Opcode() uint8 {
	return opcodes.ATHROW
}
