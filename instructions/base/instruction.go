package base

import (
	"github.com/Johnny1110/gogo_jvm/runtime"
)

type Instruction interface {
	// FetchOperands read operands from bytecode
	// some inst don't have operands (like: iconst_0), this method do nothing
	// sone isnt have operands (like: bipush), need read some data from next few bytes
	FetchOperands(reader *BytecodeReader)

	// Execute execute inst
	// access LocalVars and Stack from Frame
	Execute(frame *runtime.Frame)

	Opcode() uint8
}

// ================== define prototype =======================

// NoOperandsInstruction no operands inst
// usage:
// - iconst_0, iconst_1, ..., iconst_5
// - iadd, isub, imul, idiv
// - ireturn, control
// those inst FetchOperands do nothing
type NoOperandsInstruction struct{}

func (*NoOperandsInstruction) FetchOperands(reader *BytecodeReader) {
	// do nothing...
}

// BranchInstruction JUMP
// including a 2 bytes offset
// usage:
// goto, if_icmpeq, if_icmpne
// jump target = current pc + offset
// warning: offset could be negative (jump forward or backward)
type BranchInstruction struct {
	Offset int //  could be negative !!
}

func (b *BranchInstruction) FetchOperands(reader *BytecodeReader) {
	b.Offset = int(reader.ReadUint16()) // 2 bytes offset
}

// Index8Instruction single byte index inst
// operands is 1 byte index
// usage:
// iload, istore（current index < 256）
type Index8Instruction struct {
	Index uint
}

func (i *Index8Instruction) FetchOperands(reader *BytecodeReader) {
	i.Index = uint(reader.ReadUint8())
}

// Index16Instruction 2 bytes index inst
// operands is 2 byte index
// usage:
// getfield, invokevirtual（constant_pool index）
type Index16Instruction struct {
	Index uint
}

func (i *Index16Instruction) FetchOperands(reader *BytecodeReader) {
	i.Index = uint(reader.ReadUint16())
}
