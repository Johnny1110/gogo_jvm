package references

import (
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// hacked_invoke_native temp solution for invoke native method
func hacked_invoke_native(frame *runtime.Frame, methodRef *method_area.MethodRef) bool {
	className := methodRef.ClassName()
	methodName := methodRef.Name()
	descriptor := methodRef.Descriptor()
	if nativeMethod := runtime.FindNativeMethod(className, methodName, descriptor); nativeMethod != nil {
		// call native() directly
		invokeNativeMethod(frame, nativeMethod, descriptor)
		return true
	}
	// not in native method registry
	return false
}
