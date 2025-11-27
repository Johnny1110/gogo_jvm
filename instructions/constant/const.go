package constant

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// ============================================================
// NOP
// ============================================================

// NOP extend NoOperandsInstruction
// opcode = 0x00
// do nothing
type NOP struct {
	base.NoOperandsInstruction
}

func (n *NOP) Execute(frame *runtime.Frame) {
	// still do nothing...
}

func (n *NOP) Opcode() uint8 {
	return 0x00
}

// ============================================================
// CONST_NULL
// ============================================================

// ACONST_NULL
// opcode = 0x01
// create null Ref
type ACONST_NULL struct {
	base.NoOperandsInstruction
}

func (a *ACONST_NULL) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushRef(nil)
}

func (n *ACONST_NULL) Opcode() uint8 {
	return 0x01
}

// ============================================================
// ICONST Series (Integer)
// ============================================================
// PUSH int into stack

// ICONST_M1 PUSH -1 into stack (M1 = Minus 1)
// opcode = 0x02
type ICONST_M1 struct {
	base.NoOperandsInstruction
}

func (i *ICONST_M1) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushInt(-1)
}

func (i *ICONST_M1) Opcode() uint8 {
	return 0x02
}

// ICONST_0 PUSH 0
// opcode = 0x03
type ICONST_0 struct {
	base.NoOperandsInstruction
}

func (i *ICONST_0) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushInt(0)
}

func (i *ICONST_0) Opcode() uint8 {
	return 0x03
}

// ICONST_1 PUSH 1
// opcode = 0x04
type ICONST_1 struct{ base.NoOperandsInstruction }

func (i *ICONST_1) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushInt(1)
}

func (i *ICONST_1) Opcode() uint8 {
	return 0x04
}

// ICONST_2 PUSH 2
// opcode = 0x05
type ICONST_2 struct{ base.NoOperandsInstruction }

func (i *ICONST_2) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushInt(2)
}

func (i *ICONST_2) Opcode() uint8 {
	return 0x05
}

// ICONST_3 PUSH 3
// opcode = 0x06
type ICONST_3 struct{ base.NoOperandsInstruction }

func (i *ICONST_3) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushInt(3)
}

func (i *ICONST_3) Opcode() uint8 {
	return 0x06
}

// ICONST_4 PUSH 4
// opcode = 0x07
type ICONST_4 struct{ base.NoOperandsInstruction }

func (i *ICONST_4) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushInt(4)
}

func (i *ICONST_4) Opcode() uint8 {
	return 0x07
}

// ICONST_5 PUSH 5
// opcode = 0x08
type ICONST_5 struct{ base.NoOperandsInstruction }

func (i *ICONST_5) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushInt(5)
}

func (i *ICONST_5) Opcode() uint8 {
	return 0x08
}

// ============================================================
// LCONST Series (Long)
// ============================================================
// PUSH long into stack

// LCONST_0 PUSH 0
// opcode = 0x09
type LCONST_0 struct{ base.NoOperandsInstruction }

func (l *LCONST_0) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushLong(0)
}

func (i *LCONST_0) Opcode() uint8 {
	return 0x09
}

// LCONST_1 PUSH 1
// opcode = 0x0A
type LCONST_1 struct{ base.NoOperandsInstruction }

func (l *LCONST_1) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushLong(1)
}

func (i *LCONST_1) Opcode() uint8 {
	return 0x0A
}

// ============================================================
// FCONST Series (Float)
// ============================================================

// FCONST_0 PUSH 0.0
// opcode = 0x0B
type FCONST_0 struct{ base.NoOperandsInstruction }

func (f *FCONST_0) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushFloat(0.0)
}

func (i *FCONST_0) Opcode() uint8 {
	return 0x0B
}

// FCONST_1 PUSH 1.0
// opcode = 0x0C
type FCONST_1 struct{ base.NoOperandsInstruction }

func (f *FCONST_1) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushFloat(1.0)
}

func (i *FCONST_1) Opcode() uint8 {
	return 0x0C
}

// FCONST_2 PUSH 2.0
// opcode = 0x0D
type FCONST_2 struct{ base.NoOperandsInstruction }

func (f *FCONST_2) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushFloat(2.0)
}

func (i *FCONST_2) Opcode() uint8 {
	return 0x0D
}

// ============================================================
// DCONST Series (Double)
// ============================================================

// DCONST_0 PUSH 0.0
// opcode = 0x0E
type DCONST_0 struct{ base.NoOperandsInstruction }

func (d *DCONST_0) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushDouble(0.0)
}

func (i *DCONST_0) Opcode() uint8 {
	return 0x0E
}

// DCONST_1 PUSH 1
// opcode = 0x0F
type DCONST_1 struct{ base.NoOperandsInstruction }

func (d *DCONST_1) Execute(frame *runtime.Frame) {
	frame.OperandStack().PushDouble(1.0)
}

func (i *DCONST_1) Opcode() uint8 {
	return 0x0F
}
