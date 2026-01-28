package lang

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// Object class's native methods (hashCode etc.)

func init() {
	fmt.Println("@@ Debug - init Native java/lang/Thread")
	runtime.Register("java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", threadCurrentThread)
	runtime.Register("java/lang/Thread", "sleep", "(J)V", threadSleep)
}

func threadCurrentThread(frame *runtime.Frame) (ex *heap.Object) {
	// TODO
	panic(common.NewJavaException("Thread", "threadCurrentThread not implemented"))
}

func threadSleep(frame *runtime.Frame) (ex *heap.Object) {
	// TODO
	panic(common.NewJavaException("Thread", "threadSleep not implemented"))
}
