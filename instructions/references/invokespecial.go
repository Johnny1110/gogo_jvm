package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// INVOKESPECIAL
// 1. constructor <init>
// 2. private method
// 3. parent method (super.xxx())
// opcode = 0xB7
// operands: 2 bytes (constant pool index)
// stack: [objectref, arg1, arg2, ...] → [...]
type INVOKESPECIAL struct {
	base.Index16Instruction
}

func (i *INVOKESPECIAL) Execute(frame *runtime.Frame) {
	// TODO
}

func (i *INVOKESPECIAL) Opcode() uint8 {
	return 0xB7
}
