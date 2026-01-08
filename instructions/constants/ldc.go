package constants

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// ============================================================
// LDC - Load Constant
// ============================================================
// LDC load constnat from RuntimeConstantPool into OpStack

// LDC Load Constant，opcode = 0x12
// format: ldc index
// oprands: 1 byte (ConstantPool Index: 1~255)
// stack:  [...] -> [..., constValue]
//
// Supported Const:
// - int32
// - float32
// - string (v0.2.9)
// - *ClassRef: Class Const (v0.3.1)
type LDC struct {
	base.Index8Instruction // take 1 byte ConstantPool index
}

func (l *LDC) Execute(frame *runtime.Frame) {
	_ldc(frame, l.Index)
}

func (l *LDC) Opcode() uint8 {
	return opcodes.LDC
}

// LDC_W Load Constant (wide index)，opcode = 0x13
// format: ldc_w index1 index2
// oprands: 2 byte (ConstantPool Index: 1~65535)
// stack:  [...] -> [..., constValue]
type LDC_W struct {
	base.Index16Instruction // take 2 byte ConstantPool index
}

func (l *LDC_W) Execute(frame *runtime.Frame) {
	_ldc(frame, l.Index)
}

func (l *LDC_W) Opcode() uint8 {
	return opcodes.LDC_W
}

// LDC2_W Load Constant 2 (wide index)，opcode = 0x14 -> for double & long
// format: ldc2_w index1 index2
// oprands: 2 byte (ConstantPool Index: 1~65535)
// stack:  [...] -> [..., 2slotConstValue]
type LDC2_W struct {
	base.Index16Instruction // take 2 byte ConstantPool index
}

func (l *LDC2_W) Execute(frame *runtime.Frame) {
	rtcp := frame.Method().Class().ConstantPool()
	constVal := rtcp.GetConstant(l.Index)
	stack := frame.OperandStack()

	switch val := constVal.(type) {
	case int64:
		stack.PushLong(val)
	case float64:
		stack.PushDouble(val)
	default:
		panic(common.NewJavaException(frame.Method().Class().Name(), "unknown constant value: "+fmt.Sprint(val)))
	}
}

func (l *LDC2_W) Opcode() uint8 {
	return opcodes.LDC2_W
}

// ============================================================
// ldc and ldc_w common method
// ============================================================
// ConstantPool const type:
// - int32 -> push
// - float32 -> push
// - string -> v0.2.9 supported, using internString() create String Object
// - *ClassRef -> for Foo.class TODO: Reflection
func _ldc(frame *runtime.Frame, index uint) {
	rtcp := frame.Method().Class().ConstantPool()
	constVal := rtcp.GetConstant(index)
	stack := frame.OperandStack()

	switch val := constVal.(type) {
	case int32:
		stack.PushInt(val)
	case float32:
		stack.PushFloat(val)
	case string:
		// ldc #N (N pointing to a CONSTANT_String in rtcp - RuntimeConstantPool)
		// rtcp will return a Go string (UTF-8)
		// using heap.InternString() get java string reference and return.
		javaStrObj := heap.InternString(val)
		stack.PushRef(javaStrObj)
	case *method_area.ClassRef:
		// usage: Class<?> c = String.class;
		// In real JVM, Class constant need: (v0.3.1)
		// return java.lang.Class Object
		// Compile Result:
		//   ldc #N  // N pointing to CONSTANT_Class, value is "java/lang/String"
		// Process:
		// 1. parse ClassRef and resolveClass
		// 2. get java.lang.Class Object
		// 3. push to OpStack
		class := val.ResolvedClass()
		jClass := class.JClass()
		if jClass == nil {
			panic(fmt.Sprintf("Class object not initialized for: %s", class.Name()))
		}

		stack.PushRef(jClass)
	default:
		panic("java.lang.ClassFormatError: ldc with unknown constant type")
	}
}
