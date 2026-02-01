package lang

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// System class's native methods

func init() {
	if global.DebugMode() {
		fmt.Println("@@ Debug - init Native java/lang/System")
	}
	runtime.Register("java/lang/System", "arraycopy", "(Ljava/lang/Object;ILjava/lang/Object;II)V", systemArraycopy)
	runtime.Register("java/lang/System", "currentTimeMillis", "()J", systemCurrentTimeMillis)
}

func systemArraycopy(frame *runtime.Frame) (ex *heap.Object) {
	// TODO
	panic(common.NewJavaException("System", "systemArraycopy not implemented"))
}

func systemCurrentTimeMillis(frame *runtime.Frame) (ex *heap.Object) {
	// TODO
	panic(common.NewJavaException("System", "systemCurrentTimeMillis not implemented"))
}
