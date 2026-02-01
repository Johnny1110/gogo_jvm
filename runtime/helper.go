package runtime

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

func runtimeHandleException(frame *Frame, className string, msg string) {
	exClass := frame.LoaderClass(className)
	if exClass == nil {
		panic(fmt.Sprintf("%s: %s", className, msg))
	} else {
		exObj := heap.NewExceptionObject(exClass, msg)
		frame.JavaThrow(exObj)
	}
}
