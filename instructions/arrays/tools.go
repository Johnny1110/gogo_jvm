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
