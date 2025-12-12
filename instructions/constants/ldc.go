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
// String Pool - for string interning
// ============================================================

// LDC Load Constant，opcode = 0x12
// format: ldc index
// oprands: 1 byte (ConstantPool Index: 1~255)
// stack:  [...] -> [..., constValue]
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
// - string -> TODO: using go string right now, change to String Object in future
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
		// In real JVM, String constant need: (TODO)
		// 1. create a java.lang.String Object
		// 2. store string into Object's char[]
		// 3. String Interning
		// IN MVP Phase: simplify
		// we push nil ignore it
		javaStr := internString(val)
		stack.PushRef(javaStr)
	case *method_area.ClassRef:
		// usage: Class<?> c = String.class;
		// In real JVM, Class constant need: (TODO)
		// return java.lang.Class Object
		// IN MVP Phase: simplify
		panic("java.lang.ClassFormatError: ldc with Class constant not supported")
	default:
		panic("java.lang.ClassFormatError: ldc with unknown constant type")
	}
}

// ============================================================
// internString - Simplify
// ============================================================
// 在真實 JVM 中，字串駐留是一個重要的優化: (TODO)
//   - 相同內容的字串只創建一個物件
//   - 所有引用都指向同一個物件
//   - 節省記憶體，加快比較速度
//
// MVP 階段的簡化實現：
//   - 創建一個特殊的物件，將 Go string 存在 extra 字段
//   - 不做真正的駐留（每次都創建新物件）
func internString(goStr string) *heap.Object {
	obj := &heap.Object{}
	obj.SetExtra(goStr)
	return obj
}
