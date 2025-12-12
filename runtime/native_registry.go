package runtime

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/common"
)

// ============================================================
// Native Method Registry
// ============================================================

// NativeMethod is all native method's signature
// why need input frame?
//  1. read args from LocalVars (including this)
//  2. put return val into caller's op-stack
type NativeMethod func(frame *Frame)

// registry Native Method Registry
// key: className + methodName + descriptor
var registry = map[string]NativeMethod{}

// makeKey java support method overloading like: System.out.println(String s), System.out.println(Integer i)
// so we need add descriptor to make a method as unique
func makeKey(className, methodName, descriptor string) string {
	return fmt.Sprintf("%s~%s~%s", className, methodName, descriptor)
}

// Register a method into Native Method Registry
func Register(className, methodName, descriptor string, method NativeMethod) {
	key := makeKey(className, methodName, descriptor)
	registry[key] = method
}

// FindNativeMethod find method in Native Method Registry
func FindNativeMethod(className, methodName, descriptor string) NativeMethod {
	key := makeKey(className, methodName, descriptor)
	if method, ok := registry[key]; ok {
		return method
	}

	// TODO: 特殊處理: 空的 native 方法
	// 某些 native 方法在 MVP 階段可以忽略，例如：
	//   - registerNatives() - JVM 內部初始化用
	//   - initIDs() - JNI 初始化用
	// 這些方法對 MVP 階段沒用，先回一個 nil 未來再說
	if methodName == "registerNatives" || methodName == "initIDs" {
		return emptyNativeMethod
	}

	panic(common.NewJavaException(className, fmt.Sprintf("method %s not found", methodName)))
}

// emptyNativeMethod for MVP Phase: ignore some native methods
func emptyNativeMethod(frame *Frame) {
	// do nothing
}
