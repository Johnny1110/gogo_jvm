package constants

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// ============================================================
// BIPUSH & SIPUSH
// ============================================================

// BIPUSH PUSH 1 bytes signed byte (int) into stack
// opcodes = 0x10
type BIPUSH struct {
	val int8 // operands: 1 byte signed int
}

func (b *BIPUSH) FetchOperands(reader *base.BytecodeReader) {
	b.val = reader.ReadInt8()
}

func (b *BIPUSH) Execute(frame *runtime.Frame) {
	// int8 -> int32
	frame.OperandStack().PushInt(int32(b.val))
}

func (b *BIPUSH) Opcode() uint8 {
	return 0x10
}

// SIPUSH PUSH 2 bytes signed int into stack
// opcodes = 0x11
type SIPUSH struct {
	val int16 // operands: 2 bytes
}

func (s *SIPUSH) FetchOperands(reader *base.BytecodeReader) {
	s.val = reader.ReadInt16()
}

func (s *SIPUSH) Execute(frame *runtime.Frame) {
	// int16 -> int32
	frame.OperandStack().PushInt(int32(s.val))
}

func (s *SIPUSH) Opcode() uint8 {
	return 0x11
}
