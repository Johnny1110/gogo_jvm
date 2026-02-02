package references

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// hacked_invoke_native temp solution for invoke native method
func hacked_invoke_native(frame *runtime.Frame, methodRef *method_area.MethodRef, isStaticCall bool) bool {
	className := methodRef.ClassName()
	methodName := methodRef.Name()
	descriptor := methodRef.Descriptor()

	if nativeMethod := runtime.FindNativeMethod(className, methodName, descriptor, isStaticCall); nativeMethod != nil {
		// call native() directly
		invokeNativeMethod(frame, nativeMethod, descriptor, isStaticCall)
		return true
	}

	// not in native method registry
	if global.DebugMode() {
		fmt.Printf("@@ DEGUB - hacked_invoke_native failed, not found native method in regisrty. \n")
	}
	return false
}
