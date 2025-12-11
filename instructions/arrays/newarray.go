package arrays

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// create array
//  new array：
//  ┌────────────────┬────────┬──────────────────────────────┐
//  │ newarray       │ 0xBC   │ create basic type array      │
//  │ anewarray      │ 0xBD   │ create ref type array (TODO) │
//  │ multianewarray │ 0xC5   │ create 2D array（TODO）  	 │
//  └────────────────┴────────┴──────────────────────────────┘

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

	var arr *heap.Object
	switch n.atype {
	case AT_BOOLEAN, AT_BYTE:
		arr = heap.NewByteArray(nil, count)
	case AT_SHORT:
		arr = heap.NewShortArray(nil, count)
	case AT_CHAR:
		arr = heap.NewCharArray(nil, count)
	case AT_INT:
		arr = heap.NewIntArray(nil, count)
	case AT_LONG:
		arr = heap.NewLongArray(nil, count)
	case AT_FLOAT:
		arr = heap.NewFloatArray(nil, count)
	case AT_DOUBLE:
		arr = heap.NewDoubleArray(nil, count)
	default:
		panic("Invalid atype for NEWARRAY operandCode")
	}

	// push ref into stack
	stack.PushRef(arr)
}

func (n *NEWARRAY) Opcode() uint8 {
	return opcodes.NEWARRAY
}
