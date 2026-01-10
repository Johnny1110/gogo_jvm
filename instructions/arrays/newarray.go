package arrays

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// create array
//  new array：
//  ┌────────────────┬────────┬────────────────────────────────┐
//  │ newarray       │ 0xBC   │ create basic type array        │
//  │ anewarray      │ 0xBD   │ create ref type array (v0.2.8) │
//  │ multianewarray │ 0xC5   │ create 2D array（TODO）  	   │
//  └────────────────┴────────┴────────────────────────────────┘

// newarray - atype const
// for distinguish which type of array should be created
const (
	AT_BOOLEAN = 4
	AT_CHAR    = 5
	AT_FLOAT   = 6
	AT_DOUBLE  = 7
	AT_BYTE    = 8
	AT_SHORT   = 9
	AT_INT     = 10
	AT_LONG    = 11
)

// NEWARRAY create basic array (not Object Ref Array), opcode = 0xBC
// format: newarray atype
// atype is a byte, could be 4 ~ 11 (bool ~ long)
type NEWARRAY struct {
	atype uint8 // array type
}

func (n *NEWARRAY) FetchOperands(reader *base.BytecodeReader) {
	n.atype = reader.ReadUint8()
}

func (n *NEWARRAY) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	// get array count
	count := stack.PopInt()

	if count < 0 {
		panic("java.lang.NegativeArraySizeException")
	}

	classLoader := frame.Method().Class().Loader()
	if classLoader == nil {
		fmt.Printf("NEWARRAY error, classLoader not found.")
		panic("java.lang.ClassLoaderException")
	}

	var arr *heap.Object
	switch n.atype {
	case AT_BOOLEAN:
		// v0.3.1: dynamic load array class
		c := classLoader.LoadClass("[Z", false)
		arr = heap.NewByteArray(c, count)
	case AT_BYTE:
		c := classLoader.LoadClass("[B", false)
		arr = heap.NewByteArray(c, count)
	case AT_SHORT:
		c := classLoader.LoadClass("[S", false)
		arr = heap.NewShortArray(c, count)
	case AT_CHAR:
		c := classLoader.LoadClass("[C", false)
		arr = heap.NewCharArray(c, count)
	case AT_INT:
		c := classLoader.LoadClass("[I", false)
		arr = heap.NewIntArray(c, count)
	case AT_LONG:
		c := classLoader.LoadClass("[J", false)
		arr = heap.NewLongArray(c, count)
	case AT_FLOAT:
		c := classLoader.LoadClass("[F", false)
		arr = heap.NewFloatArray(c, count)
	case AT_DOUBLE:
		c := classLoader.LoadClass("[D", false)
		arr = heap.NewDoubleArray(c, count)
	default:
		panic("Invalid atype for NEWARRAY operandCode")
	}

	// push ref into stack
	stack.PushRef(arr)
}

func (n *NEWARRAY) Opcode() uint8 {
	return opcodes.NEWARRAY
}

// ----------------------------------------------------------------------------

// ANEWARRAY create Ref type Array (TODO: v3.x version will revamp this)
// opcode = 0xBD
// format: anewarray indexbyte1 indexbyte2
// operands: 2 bytes (constant pool index pointing to ClassRef)
// stack: [..., count] → [..., arrayref]
// usage: create object array, like: String[], Object[], MyClass[] ...
type ANEWARRAY struct {
	base.Index16Instruction // 2 bytes -> type index in rtcp
}

func (a *ANEWARRAY) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	count := stack.PopInt() // array size
	if count < 0 {
		panic(common.NewJavaException("", "java.lang.NegativeArraySizeException"))
	}

	// get type from rtcp
	rtcp := frame.Method().Class().ConstantPool()
	elementClassRef := rtcp.GetConstant(a.Index).(*method_area.ClassRef)
	elementClassName := elementClassRef.ResolvedClass().Name()

	// create array
	// class will be store in object's class field
	// "java/lang/String" ->"[Ljava/lang/String;"
	arrayClassName := "[L" + elementClassName + ";"
	classLoader := frame.Method().Class().Loader()
	arrayClass := classLoader.LoadClass(arrayClassName, false)
	array := heap.NewRefArray(arrayClass, count)

	// push array ref into stack
	stack.PushRef(array)
}

func (a *ANEWARRAY) Opcode() uint8 {
	// opcode = 0xBD
	return opcodes.ANEWARRAY
}

// ----------------------------------------------------------------------------
// ANEWARRAY create 2D Array
// multianewarray indexbyte1 indexbyte2 dimensions
// opcode: 0xC5
// operands: 3 bytes total
//   - 2 bytes: constant pool index (pointing to array type, like "[[I")
//   - 1 byte: dimensions
type MULTIANEWARRAY struct {
	base.Index16Instruction
	dimensions uint8
}

func (m *MULTIANEWARRAY) FetchOperands(reader *base.BytecodeReader) {
	m.Index = uint(reader.ReadUint16())
	m.dimensions = reader.ReadUint8()
}

func (m *MULTIANEWARRAY) Execute(frame *runtime.Frame) {
	// 1. get array class type
	rtcp := frame.Method().Class().ConstantPool()
	classRef := rtcp.GetConstant(m.Index).(*method_area.ClassRef)
	arrayClass := classRef.ResolvedClass()

	// 2. pop length of all elements
	stack := frame.OperandStack()
	counts := popAndCheckCounts(stack, int(m.dimensions))
	// count will be like:
	// new int[2][3] → stack: [..., 2, 3] → counts: [2, 3]

	// 3. create all dimensional array
	arr := newMultiDimensionalArray(counts, arrayClass)

	// 4. push array ref to stack
	stack.PushRef(arr)
}

// popAndCheckCounts 彈出各維度的長度並檢查 pop every dimensions len from stack and check
// stack order: 最外層維度在最下面，最內層在最上面
// ex: new int[2][3] → stack: [..., 2, 3] → counts: [2, 3]
func popAndCheckCounts(stack *runtime.OperandStack, dimensions int) []int32 {
	counts := make([]int32, dimensions)

	// pop from stack.
	for i := dimensions - 1; i >= 0; i-- {
		counts[i] = stack.PopInt()
		if counts[i] < 0 {
			// should not happen
			panic("java.lang.NegativeArraySizeException")
		}
	}

	return counts
}

// newMultiDimensionalArray create multi-dimension array (recursive)
// counts: dimensions len, ex [2, 3] -> new Type[2][3]
// arrayClass: array type, ex: "[[I" or "[[[Ljava/lang/String;"
func newMultiDimensionalArray(counts []int32, arrayClass *method_area.Class) *heap.Object {
	// current dimensional count
	count := counts[0]
	// create current dimensional array
	arr := heap.NewRefArray(arrayClass, count)

	if len(counts) > 1 { // if it has more
		refs := arr.Refs() // get current array's ref

		// get array type (sub array type)
		// if current arrayClass is "[[Ljava/lang/String;", ComponentClass() should be "[Ljava/lang/String;"
		componentClass := arrayClass.ComponentClass()

		for i := int32(0); i < count; i++ { // each current row should create new sub array.
			subArr := newMultiDimensionalArray(counts[1:], componentClass)
			refs[i] = subArr
		}
	}

	return arr
}

func (m *MULTIANEWARRAY) Opcode() uint8 {
	return opcodes.MULTIANEWARRAY
}
