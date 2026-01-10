package arrays

import (
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// checkArrayAndIndex check ref not null and index
func checkArrayAndIndex(stack *runtime.OperandStack) (*heap.Object, int32) {
	index := stack.PopInt()
	arrRef := stack.PopRef()

	if arrRef == nil {
		panic("java.lang.NullPointerException")
	}

	arr := arrRef.(*heap.Object)
	return arr, index
}

// getArrayClassName generate array class name according to element's name
// int -> [I
// java/lang/String → [Ljava/lang/String;
// [I → [[I
func getArrayClassName(className string) string {
	if len(className) > 0 && className[0] == '[' {
		return "[" + className
	}

	// basic type
	switch className {
	case "void":
		panic("cannot create array of void")
	case "boolean":
		return "[Z"
	case "byte":
		return "[B"
	case "char":
		return "[C"
	case "short":
		return "[S"
	case "int":
		return "[I"
	case "long":
		return "[J"
	case "float":
		return "[F"
	case "double":
		return "[D"
	default:
		// Ref Type
		return "[L" + className + ";"
	}
}
