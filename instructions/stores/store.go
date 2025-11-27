package stores

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// ============================================================
// ISTORE Series (int)
// ============================================================

// opcode = 0x36
type ISTORE struct{ base.Index8Instruction }

func (i *ISTORE) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	frame.LocalVars().SetInt(i.Index, val)
}

func (i *ISTORE) Opcode() uint8 {
	return 0x36
}

// ISTORE_0 opcode = 0x3B
type ISTORE_0 struct{ base.NoOperandsInstruction }

func (i *ISTORE_0) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	frame.LocalVars().SetInt(0, val)
}

func (i *ISTORE_0) Opcode() uint8 {
	return 0x3B
}

// ISTORE_1 opcode = 0x3C
type ISTORE_1 struct{ base.NoOperandsInstruction }

func (i *ISTORE_1) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	frame.LocalVars().SetInt(1, val)
}

func (i *ISTORE_1) Opcode() uint8 {
	return 0x3C
}

// ISTORE_2 opcode = 0x3D
type ISTORE_2 struct{ base.NoOperandsInstruction }

func (i *ISTORE_2) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	frame.LocalVars().SetInt(2, val)
}

func (i *ISTORE_2) Opcode() uint8 {
	return 0x3D
}

// ISTORE_3 opcode = 0x3E
type ISTORE_3 struct{ base.NoOperandsInstruction }

func (i *ISTORE_3) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	frame.LocalVars().SetInt(3, val)
}

func (i *ISTORE_3) Opcode() uint8 {
	return 0x3E
}

// ============================================================
// LSTORE Series (long)
// ============================================================

// LSTORE opcode = 0x37
type LSTORE struct{ base.Index8Instruction }

func (l *LSTORE) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopLong()
	frame.LocalVars().SetLong(l.Index, val)
}

func (i *LSTORE) Opcode() uint8 {
	return 0x37
}

// LSTORE_0 opcode = 0x3F
type LSTORE_0 struct{ base.NoOperandsInstruction }

func (l *LSTORE_0) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopLong()
	frame.LocalVars().SetLong(0, val)
}

func (i *LSTORE_0) Opcode() uint8 {
	return 0x3F
}

// LSTORE_1 opcode = 0x40
type LSTORE_1 struct{ base.NoOperandsInstruction }

func (l *LSTORE_1) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopLong()
	frame.LocalVars().SetLong(1, val)
}

func (i *LSTORE_1) Opcode() uint8 {
	return 0x40
}

// LSTORE_2 opcode = 0x41
type LSTORE_2 struct{ base.NoOperandsInstruction }

func (l *LSTORE_2) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopLong()
	frame.LocalVars().SetLong(2, val)
}

func (i *LSTORE_2) Opcode() uint8 {
	return 0x41
}

// LSTORE_3 opcode = 0x42
type LSTORE_3 struct{ base.NoOperandsInstruction }

func (l *LSTORE_3) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopLong()
	frame.LocalVars().SetLong(3, val)
}

func (i *LSTORE_3) Opcode() uint8 {
	return 0x42
}

// ============================================================
// FSTORE Series (float)
// ============================================================

// FSTORE opcode = 0x38
type FSTORE struct{ base.Index8Instruction }

func (f *FSTORE) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopFloat()
	frame.LocalVars().SetFloat(f.Index, val)
}

func (i *FSTORE) Opcode() uint8 {
	return 0x38
}

// FSTORE_0 opcode = 0x43
type FSTORE_0 struct{ base.NoOperandsInstruction }

func (f *FSTORE_0) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopFloat()
	frame.LocalVars().SetFloat(0, val)
}

func (i *FSTORE_0) Opcode() uint8 {
	return 0x43
}

// FSTORE_1 opcode = 0x44
type FSTORE_1 struct{ base.NoOperandsInstruction }

func (f *FSTORE_1) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopFloat()
	frame.LocalVars().SetFloat(1, val)
}

func (i *FSTORE_1) Opcode() uint8 {
	return 0x44
}

// FSTORE_2 opcode = 0x45
type FSTORE_2 struct{ base.NoOperandsInstruction }

func (f *FSTORE_2) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopFloat()
	frame.LocalVars().SetFloat(2, val)
}

func (i *FSTORE_2) Opcode() uint8 {
	return 0x45
}

// FSTORE_3 opcode = 0x46
type FSTORE_3 struct{ base.NoOperandsInstruction }

func (f *FSTORE_3) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopFloat()
	frame.LocalVars().SetFloat(3, val)
}

func (i *FSTORE_3) Opcode() uint8 {
	return 0x46
}

// ============================================================
// DSTORE Series (double)
// ============================================================

// DSTORE opcode = 0x39
type DSTORE struct{ base.Index8Instruction }

func (d *DSTORE) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopDouble()
	frame.LocalVars().SetDouble(d.Index, val)
}

func (i *DSTORE) Opcode() uint8 {
	return 0x39
}

// DSTORE_0 opcode = 0x47
type DSTORE_0 struct{ base.NoOperandsInstruction }

func (d *DSTORE_0) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopDouble()
	frame.LocalVars().SetDouble(0, val)
}

func (i *DSTORE_0) Opcode() uint8 {
	return 0x47
}

// DSTORE_1 opcode = 0x48
type DSTORE_1 struct{ base.NoOperandsInstruction }

func (d *DSTORE_1) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopDouble()
	frame.LocalVars().SetDouble(1, val)
}

func (i *DSTORE_1) Opcode() uint8 {
	return 0x48
}

// DSTORE_2 opcode = 0x49
type DSTORE_2 struct{ base.NoOperandsInstruction }

func (d *DSTORE_2) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopDouble()
	frame.LocalVars().SetDouble(2, val)
}

func (i *DSTORE_2) Opcode() uint8 {
	return 0x49
}

// DSTORE_3 opcode = 0x4A
type DSTORE_3 struct{ base.NoOperandsInstruction }

func (d *DSTORE_3) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopDouble()
	frame.LocalVars().SetDouble(3, val)
}

func (i *DSTORE_3) Opcode() uint8 {
	return 0x4A
}

// ============================================================
// ASTORE Series (reference)
// ============================================================

// ASTORE opcode = 0x3A
type ASTORE struct{ base.Index8Instruction }

func (a *ASTORE) Execute(frame *runtime.Frame) {
	ref := frame.OperandStack().PopRef()
	frame.LocalVars().SetRef(a.Index, ref)
}

func (i *ASTORE) Opcode() uint8 {
	return 0x3A
}

// ASTORE_0 opcode = 0x4B
type ASTORE_0 struct{ base.NoOperandsInstruction }

func (a *ASTORE_0) Execute(frame *runtime.Frame) {
	ref := frame.OperandStack().PopRef()
	frame.LocalVars().SetRef(0, ref)
}

func (i *ASTORE_0) Opcode() uint8 {
	return 0x4B
}

// ASTORE_1 opcode = 0x4C
type ASTORE_1 struct{ base.NoOperandsInstruction }

func (a *ASTORE_1) Execute(frame *runtime.Frame) {
	ref := frame.OperandStack().PopRef()
	frame.LocalVars().SetRef(1, ref)
}

func (i *ASTORE_1) Opcode() uint8 {
	return 0x4C
}

// ASTORE_2 opcode = 0x4D
type ASTORE_2 struct{ base.NoOperandsInstruction }

func (a *ASTORE_2) Execute(frame *runtime.Frame) {
	ref := frame.OperandStack().PopRef()
	frame.LocalVars().SetRef(2, ref)
}

func (i *ASTORE_2) Opcode() uint8 {
	return 0x4D
}

// ASTORE_3 opcode = 0x4E
type ASTORE_3 struct{ base.NoOperandsInstruction }

func (a *ASTORE_3) Execute(frame *runtime.Frame) {
	ref := frame.OperandStack().PopRef()
	frame.LocalVars().SetRef(3, ref)
}

func (i *ASTORE_3) Opcode() uint8 {
	return 0x4E
}
