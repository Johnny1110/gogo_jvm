package lang

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// Object class's native methods (hashCode etc.)

func init() {
	fmt.Println("@@ Debug - init Native java/lang/Object")
	runtime.Register("java/lang/Object", "hashCode", "()I", objectHashCode)
	runtime.Register("java/lang/Object", "getClass", "()Ljava/lang/Class;", objectGetClass)
}

// objectHashCode native implementation of Object.hashCode()
// Java signature: public native int hashCode();
// return Object's identity hash code (store in Object's markWord)
func objectHashCode(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis()
	if this == nil {
		// should not be happened
		frame.OperandStack().PushInt(0)
		return
	}

	obj := this.(*heap.Object)
	hash := obj.HashCode()
	frame.OperandStack().PushInt(hash)
}

func objectGetClass(frame *runtime.Frame) {
	// TODO
	panic(common.NewJavaException("Object", "objectGetClass not implemented - wait for v0.3.1"))
}
