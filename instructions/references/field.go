package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// ============================================================
// GETSTATIC - get static field, opcode = 0xB2
// ============================================================
// format: getstatic indexbyte1 indexbyte2
// op counts: 2 bytes (uint16), rt constant pool index, pointing to a FieldRef

type GETSTATIC struct {
	base.Index16Instruction // uint16, pointing to FieldRef
}

func (g *GETSTATIC) Execute(frame *runtime.Frame) {
	// 1. get const pool
	rtcp := frame.Method().Class().ConstantPool()
	// 2. get field from rtcp
	fieldRef := rtcp.GetConstant(g.Index).(*method_area.FieldRef)
	// 3. load field
	field := fieldRef.ResolvedField()
	// 4. get class (if field resolved, class should already be resolved also)
	class := field.Class()
	// 5. check class <clinit>
	if !class.InitStarted() {
		frame.RevertNextPC() // do it again after InitClass.
		InitClass(frame.Thread(), class)
		return
	}
	// 6. check static
	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 7. get slot and desc
	slotId := field.SlotId()
	staticVarSlots := class.StaticVars()
	stack := frame.OperandStack()
	descriptor := field.Descriptor()

	// 8. push val into stack
	pushFieldValue(stack, staticVarSlots, slotId, descriptor)
}

func (g *GETSTATIC) Opcode() uint8 {
	return 0xB2
}

// ============================================================
// PUTSTATIC - put static field, opcode = 0xB3
// ============================================================
// format: putstatic indexbyte1 indexbyte2
// op counts: 2 bytes, rt constant pool index, pointing to a FieldRef

type PUTSTATIC struct {
	base.Index16Instruction
}

func (p *PUTSTATIC) Execute(frame *runtime.Frame) {
	// 1. get constant pool
	cp := frame.Method().Class().ConstantPool()

	// 2. get FieldRef
	fieldRef := cp.GetConstant(p.Index).(*method_area.FieldRef)

	// 3. resolve field
	field := fieldRef.ResolvedField()

	// 4. get class
	class := field.Class()

	// 5. do <clinit> if required
	if !class.InitStarted() {
		frame.RevertNextPC()
		InitClass(frame.Thread(), class)
		return
	}

	// 6. check status
	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 7. get field's slotId and descriptor
	slotId := field.SlotId()
	slots := class.StaticVars()
	stack := frame.OperandStack()
	descriptor := field.Descriptor()

	// pop val from stack and put into staticVars
	popAndSetFieldValue(stack, slots, slotId, descriptor)
}

func (p *PUTSTATIC) Opcode() uint8 {
	return 0xB3
}

// ============================================================
// GETFIELD - get instant field, opcode = 0xB4
// ============================================================
// format: getfield indexbyte1 indexbyte2
// op counts: 2 bytes, rt instant pool index, pointing to FieldRef

type GETFIELD struct {
	base.Index16Instruction
}

func (g *GETFIELD) Execute(frame *runtime.Frame) {
	cp := frame.Method().Class().ConstantPool()
	fieldRef := cp.GetConstant(g.Index).(*method_area.FieldRef)
	field := fieldRef.ResolvedField()
	// 1. not static
	if field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	stack := frame.OperandStack()
	ref := stack.PopRef() // target object's ref

	// 2. check null
	checkNotNull(ref)

	// 3. convert ref to object pointer
	object := ref.(*heap.Object)

	// 4. get slotId and descriptor
	slotId := field.SlotId()
	slots := object.Fields()
	descriptor := field.Descriptor()

	// 5. push into stack
	pushFieldValue(stack, slots, slotId, descriptor)
}

func (g *GETFIELD) Opcode() uint8 {
	return 0xB4
}

// ============================================================
// PUTFIELD - set instant field, opcode = 0xB5
// ============================================================
// format: putfield indexbyte1 indexbyte2
// op counts: 2 bytes, rt constant pool, pointing to FieldRef

type PUTFIELD struct {
	base.Index16Instruction
}

func (p *PUTFIELD) Execute(frame *runtime.Frame) {
	cp := frame.Method().Class().ConstantPool()
	fieldRef := cp.GetConstant(p.Index).(*method_area.FieldRef)
	field := fieldRef.ResolvedField()
	if field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	slotId := field.SlotId()
	descriptor := field.Descriptor()
	stack := frame.OperandStack()

	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		// boolean, byte, char, short, int
		val := stack.PopInt()
		ref := stack.PopRef()
		checkNotNull(ref)
		ref.(*heap.Object).Fields().SetInt(slotId, val)
	case 'F':
		// float
		val := stack.PopFloat()
		ref := stack.PopRef()
		checkNotNull(ref)
		ref.(*heap.Object).Fields().SetFloat(slotId, val)
	case 'J':
		// long
		val := stack.PopLong()
		ref := stack.PopRef()
		checkNotNull(ref)
		ref.(*heap.Object).Fields().SetLong(slotId, val)
	case 'D':
		val := stack.PopDouble()
		ref := stack.PopRef()
		checkNotNull(ref)
		ref.(*heap.Object).Fields().SetDouble(slotId, val)
	case 'L', '[':
		// ref or array
		val := stack.PopRef()
		ref := stack.PopRef()
		checkNotNull(ref)
		ref.(*heap.Object).Fields().SetRef(slotId, val)
	default:
		panic("Unknown field descriptor: " + descriptor)
	}
}

func (p *PUTFIELD) Opcode() uint8 {
	return 0xB5
}
