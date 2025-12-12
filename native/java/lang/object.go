package lang

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// Object class's native methods (hashCode etc.)

func init() {
	fmt.Println("@@ Debug - init Native java/lang/Object")
	runtime.Register("java/lang/Object", "hashCode", "()I", objectHashCode)
	runtime.Register("java/lang/Object", "getClass", "()Ljava/lang/Class;", objectGetClass)
}

func objectHashCode(frame *runtime.Frame) {
	// TODO
	panic(common.NewJavaException("Object", "objectHashCode not implemented"))
}

func objectGetClass(frame *runtime.Frame) {
	// TODO
	panic(common.NewJavaException("Object", "objectGetClass not implemented"))
}
