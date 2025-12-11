package arrays

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// xaload/xastore series inst
//  ┌────────────────┬────────┬──────────────────────────────┐
//  │ iaload/iastore │ 0x2E/4F│ int[]                        │
//  │ laload/lastore │ 0x2F/50│ long[]                       │
//  │ faload/fastore │ 0x30/51│ float[]                      │
//  │ daload/dastore │ 0x31/52│ double[]                     │
//  │ aaload/aastore │ 0x32/53│ Object[]                     │
//  │ baload/bastore │ 0x33/54│ byte[]/boolean[]             │
//  │ caload/castore │ 0x34/55│ char[]                       │
//  │ saload/sastore │ 0x35/56│ short[]                      │
//  └────────────────┴────────┴──────────────────────────────┘

// ============================================================
// xALOAD series - load element into stack from array
// Stack: [..., arrayref, index] → [..., value]
// ============================================================

// IALOAD load int from array
// opcode = 0x2E
type IALOAD struct{ base.NoOperandsInstruction }

func (i *IALOAD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arr, index := checkArrayAndIndex(stack)
	val := arr.GetArrayInt(index)
	stack.PushInt(val)
}

func (i *IALOAD) Opcode() uint8 { return opcodes.IALOAD }

// LALOAD load long from array
// opcode = 0x2F
type LALOAD struct{ base.NoOperandsInstruction }

func (l *LALOAD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arr, index := checkArrayAndIndex(stack)
	val := arr.GetArrayLong(index)
	stack.PushLong(val)
}

func (l *LALOAD) Opcode() uint8 { return opcodes.LALOAD }

// FALOAD load float from array
// opcode = 0x30
type FALOAD struct{ base.NoOperandsInstruction }

func (f *FALOAD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arr, index := checkArrayAndIndex(stack)
	val := arr.GetArrayFloat(index)
	stack.PushFloat(val)
}

func (f *FALOAD) Opcode() uint8 { return opcodes.FALOAD }

// DALOAD load double from array
// opcode = 0x31
type DALOAD struct{ base.NoOperandsInstruction }

func (d *DALOAD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arr, index := checkArrayAndIndex(stack)
	val := arr.GetArrayDouble(index)
	stack.PushDouble(val)
}

func (d *DALOAD) Opcode() uint8 { return opcodes.DALOAD }

// AALOAD load Ref from array
// opcode = 0x32
type AALOAD struct{ base.NoOperandsInstruction }

func (a *AALOAD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arr, index := checkArrayAndIndex(stack)
	ref := arr.GetArrayRef(index)
	stack.PushRef(ref)
}

func (a *AALOAD) Opcode() uint8 { return opcodes.AALOAD }

// BALOAD load byte/boolean from array (push int32 into stack)
// opcode = 0x33
type BALOAD struct{ base.NoOperandsInstruction }

func (b *BALOAD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arr, index := checkArrayAndIndex(stack)
	val := arr.GetArrayByte(index)
	stack.PushInt(int32(val)) // byte -> int32
}

func (b *BALOAD) Opcode() uint8 { return opcodes.BALOAD }

// CALOAD load char from array
// opcode = 0x34
type CALOAD struct{ base.NoOperandsInstruction }

func (c *CALOAD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arr, index := checkArrayAndIndex(stack)
	val := arr.GetArrayChar(index)
	stack.PushInt(int32(val)) // char -> int32
}

func (c *CALOAD) Opcode() uint8 { return opcodes.CALOAD }

// SALOAD load short from array
// opcode = 0x35
type SALOAD struct{ base.NoOperandsInstruction }

func (s *SALOAD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	arr, index := checkArrayAndIndex(stack)
	val := arr.GetArrayShort(index)
	stack.PushInt(int32(val)) // short -> int32
}

func (s *SALOAD) Opcode() uint8 { return opcodes.SALOAD }

// ============================================================
// xASTORE series - store element into array
// Stack: [..., arrayref, index, value] → [...]
// ============================================================

// IASTORE store int into array
// opcode = 0x4F
type IASTORE struct{ base.NoOperandsInstruction }

func (i *IASTORE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	arr.SetArrayInt(index, val)
}

func (i *IASTORE) Opcode() uint8 { return opcodes.IASTORE }

// LASTORE store long into array
// opcode = 0x50
type LASTORE struct{ base.NoOperandsInstruction }

func (l *LASTORE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopLong()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	arr.SetArrayLong(index, val)
}

func (l *LASTORE) Opcode() uint8 { return opcodes.LASTORE }

// FASTORE store float into array
// opcode = 0x51
type FASTORE struct{ base.NoOperandsInstruction }

func (f *FASTORE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopFloat()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	arr.SetArrayFloat(index, val)
}

func (f *FASTORE) Opcode() uint8 { return opcodes.FASTORE }

// DASTORE store double into array
// opcode = 0x52
type DASTORE struct{ base.NoOperandsInstruction }

func (d *DASTORE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopDouble()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	arr.SetArrayDouble(index, val)
}

func (d *DASTORE) Opcode() uint8 { return opcodes.DASTORE }

// AASTORE store ref into array
// opcode = 0x53
type AASTORE struct{ base.NoOperandsInstruction }

func (a *AASTORE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	ref := stack.PopRef()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	// TODO: 完整實現需要檢查類型相容性
	if ref == nil {
		arr.SetArrayRef(index, nil)
	} else {
		arr.SetArrayRef(index, ref.(*heap.Object))
	}
}

func (a *AASTORE) Opcode() uint8 { return opcodes.AASTORE }

// BASTORE store byte/boolean into array
// opcode = 0x54
type BASTORE struct{ base.NoOperandsInstruction }

func (b *BASTORE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	arr.SetArrayByte(index, int8(val)) // int32 -> int8
}

func (b *BASTORE) Opcode() uint8 { return opcodes.BASTORE }

// CASTORE store char into array
// opcode = 0x55
type CASTORE struct{ base.NoOperandsInstruction }

func (c *CASTORE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	arr.SetArrayChar(index, uint16(val)) // int32 -> char(uint16)
}

func (c *CASTORE) Opcode() uint8 { return opcodes.CASTORE }

// SASTORE store short into array
// opcode = 0x56
type SASTORE struct{ base.NoOperandsInstruction }

func (s *SASTORE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	arr.SetArrayShort(index, int16(val)) // int32 -> short(int16)
}

func (s *SASTORE) Opcode() uint8 { return opcodes.SASTORE }
