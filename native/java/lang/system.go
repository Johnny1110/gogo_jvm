package lang

import (
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// System class's native methods

func init() {
	runtime.Register("java/lang/System", "arraycopy", "(Ljava/lang/Object;ILjava/lang/Object;II)V", systemArraycopy)
	runtime.Register("java/lang/System", "currentTimeMillis", "()J", systemCurrentTimeMillis)
}

func systemArraycopy(frame *runtime.Frame) {
	// TODO
	panic(common.NewJavaException("System", "systemArraycopy not implemented"))
}

func systemCurrentTimeMillis(frame *runtime.Frame) {
	// TODO
	panic(common.NewJavaException("System", "systemCurrentTimeMillis not implemented"))
}
