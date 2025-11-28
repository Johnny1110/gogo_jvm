package loads

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// ============================================================
// ILOAD Series (int)
// ============================================================

// ILOAD load int from LocalVars into stack (copy)
// opcodes = 0x15
type ILOAD struct{ base.Index8Instruction }

func (i *ILOAD) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetInt(i.Index)
	frame.OperandStack().PushInt(val)
}

func (i *ILOAD) Opcode() uint8 {
	return 0x15
}

// ILOAD_0 load data from LocalVars by index 0
// opcodes = 0x1A
type ILOAD_0 struct{ base.NoOperandsInstruction }

func (i *ILOAD_0) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetInt(0)
	frame.OperandStack().PushInt(val)
}

func (i *ILOAD_0) Opcode() uint8 {
	return 0x1A
}

// ILOAD_1 load data from LocalVars by index 1
// opcodes = 0x1B
type ILOAD_1 struct{ base.NoOperandsInstruction }

func (i *ILOAD_1) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetInt(1)
	frame.OperandStack().PushInt(val)
}

func (i *ILOAD_1) Opcode() uint8 {
	return 0x1B
}

// ILOAD_2 load data from LocalVars by index 2
// opcodes = 0x1C
type ILOAD_2 struct{ base.NoOperandsInstruction }

func (i *ILOAD_2) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetInt(2)
	frame.OperandStack().PushInt(val)
}

func (i *ILOAD_2) Opcode() uint8 {
	return 0x1C
}

// ILOAD_3 load data from LocalVars by index 3
// opcodes = 0x1D
type ILOAD_3 struct{ base.NoOperandsInstruction }

func (i *ILOAD_3) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetInt(3)
	frame.OperandStack().PushInt(val)
}

func (i *ILOAD_3) Opcode() uint8 {
	return 0x1D
}

// ============================================================
// LLOAD Series (long)
// ============================================================

// LLOAD load long from LocalVars
// opcodes = 0x16
type LLOAD struct{ base.Index8Instruction }

func (l *LLOAD) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetLong(l.Index)
	frame.OperandStack().PushLong(val)
}

func (i *LLOAD) Opcode() uint8 {
	return 0x16
}

// LLOAD_0 opcodes = 0x1E
type LLOAD_0 struct{ base.NoOperandsInstruction }

func (l *LLOAD_0) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetLong(0)
	frame.OperandStack().PushLong(val)
}

func (i *LLOAD_0) Opcode() uint8 {
	return 0x1E
}

// LLOAD_1 opcodes = 0x1F
type LLOAD_1 struct{ base.NoOperandsInstruction }

func (l *LLOAD_1) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetLong(1)
	frame.OperandStack().PushLong(val)
}

func (i *LLOAD_1) Opcode() uint8 {
	return 0x1F
}

// LLOAD_2 opcodes = 0x20
type LLOAD_2 struct{ base.NoOperandsInstruction }

func (l *LLOAD_2) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetLong(2)
	frame.OperandStack().PushLong(val)
}

func (i *LLOAD_2) Opcode() uint8 {
	return 0x20
}

// LLOAD_3 opcodes = 0x21
type LLOAD_3 struct{ base.NoOperandsInstruction }

func (l *LLOAD_3) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetLong(3)
	frame.OperandStack().PushLong(val)
}

func (i *LLOAD_3) Opcode() uint8 {
	return 0x21
}

// ============================================================
// FLOAD Series (float)
// ============================================================

// FLOAD opcodes = 0x17
type FLOAD struct{ base.Index8Instruction }

func (f *FLOAD) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetFloat(f.Index)
	frame.OperandStack().PushFloat(val)
}

func (i *FLOAD) Opcode() uint8 {
	return 0x17
}

// FLOAD_0 opcodes = 0x22
type FLOAD_0 struct{ base.NoOperandsInstruction }

func (f *FLOAD_0) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetFloat(0)
	frame.OperandStack().PushFloat(val)
}

func (i *FLOAD_0) Opcode() uint8 {
	return 0x22
}

// FLOAD_1 opcodes = 0x23
type FLOAD_1 struct{ base.NoOperandsInstruction }

func (f *FLOAD_1) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetFloat(1)
	frame.OperandStack().PushFloat(val)
}

func (i *FLOAD_1) Opcode() uint8 {
	return 0x23
}

// FLOAD_2 opcodes = 0x24
type FLOAD_2 struct{ base.NoOperandsInstruction }

func (f *FLOAD_2) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetFloat(2)
	frame.OperandStack().PushFloat(val)
}

func (i *FLOAD_2) Opcode() uint8 {
	return 0x24
}

// FLOAD_3 opcodes = 0x25
type FLOAD_3 struct{ base.NoOperandsInstruction }

func (f *FLOAD_3) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetFloat(3)
	frame.OperandStack().PushFloat(val)
}

func (i *FLOAD_3) Opcode() uint8 {
	return 0x25
}

// ============================================================
// DLOAD Series (double)
// ============================================================

// DLOAD opcodes = 0x18
type DLOAD struct{ base.Index8Instruction }

func (d *DLOAD) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetDouble(d.Index)
	frame.OperandStack().PushDouble(val)
}

func (i *DLOAD) Opcode() uint8 {
	return 0x18
}

// DLOAD_0 opcodes = 0x26
type DLOAD_0 struct{ base.NoOperandsInstruction }

func (d *DLOAD_0) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetDouble(0)
	frame.OperandStack().PushDouble(val)
}

func (i *DLOAD_0) Opcode() uint8 {
	return 0x26
}

// DLOAD_1 opcodes = 0x27
type DLOAD_1 struct{ base.NoOperandsInstruction }

func (d *DLOAD_1) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetDouble(1)
	frame.OperandStack().PushDouble(val)
}

func (i *DLOAD_1) Opcode() uint8 {
	return 0x27
}

// DLOAD_2 opcodes = 0x28
type DLOAD_2 struct{ base.NoOperandsInstruction }

func (d *DLOAD_2) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetDouble(2)
	frame.OperandStack().PushDouble(val)
}

func (i *DLOAD_2) Opcode() uint8 {
	return 0x28
}

// DLOAD_3 opcodes = 0x29
type DLOAD_3 struct{ base.NoOperandsInstruction }

func (d *DLOAD_3) Execute(frame *runtime.Frame) {
	val := frame.LocalVars().GetDouble(3)
	frame.OperandStack().PushDouble(val)
}

func (i *DLOAD_3) Opcode() uint8 {
	return 0x29
}

// ============================================================
// ALOAD 系列 (reference)
// ============================================================

// ALOAD opcodes = 0x19
type ALOAD struct{ base.Index8Instruction }

func (a *ALOAD) Execute(frame *runtime.Frame) {
	ref := frame.LocalVars().GetRef(a.Index)
	frame.OperandStack().PushRef(ref)
}

func (i *ALOAD) Opcode() uint8 {
	return 0x19
}

// ALOAD_0 opcodes = 0x2A
// in practical: aload_0 load `this` ref
type ALOAD_0 struct{ base.NoOperandsInstruction }

func (a *ALOAD_0) Execute(frame *runtime.Frame) {
	ref := frame.LocalVars().GetRef(0)
	frame.OperandStack().PushRef(ref)
}

func (i *ALOAD_0) Opcode() uint8 {
	return 0x2A
}

// ALOAD_1 opcodes = 0x2B
type ALOAD_1 struct{ base.NoOperandsInstruction }

func (a *ALOAD_1) Execute(frame *runtime.Frame) {
	ref := frame.LocalVars().GetRef(1)
	frame.OperandStack().PushRef(ref)
}

func (i *ALOAD_1) Opcode() uint8 {
	return 0x2B
}

// ALOAD_2 opcodes = 0x2C
type ALOAD_2 struct{ base.NoOperandsInstruction }

func (a *ALOAD_2) Execute(frame *runtime.Frame) {
	ref := frame.LocalVars().GetRef(2)
	frame.OperandStack().PushRef(ref)
}

func (i *ALOAD_2) Opcode() uint8 {
	return 0x2C
}

// ALOAD_3 opcodes = 0x2D
type ALOAD_3 struct{ base.NoOperandsInstruction }

func (a *ALOAD_3) Execute(frame *runtime.Frame) {
	ref := frame.LocalVars().GetRef(3)
	frame.OperandStack().PushRef(ref)
}

func (i *ALOAD_3) Opcode() uint8 {
	return 0x2D
}
