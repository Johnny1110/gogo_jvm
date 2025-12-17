package references

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"os"
)

// ============================================================
// Handle Exception Common Func
// ============================================================
// ThrowException
// usage: idiv, aaload, checkcast ...
func ThrowException(frame *runtime.Frame, exceptionObj *heap.Object) {
	thread := frame.Thread()
	// 2. find ex handler and process exception
	if !handleException(thread, exceptionObj) {
		// can not find handler to process exception -> UncaughtException
		handleUncaughtException(thread, exceptionObj)
	}
}

// ============================================================
// Exception Handler Core Logic
// ============================================================

// handleException process exception return success catch or not
// args:
//   - currentThread
//   - exceptionObj: target exception
//
// return:
//   - true: found handler and goto target PC
//   - false: can not find handler until searched all JVMStack
func handleException(currentThread *runtime.Thread, exceptionObj *heap.Object) bool {
	for {
		frame := currentThread.CurrentFrame()
		pc := frame.CurrentPC() // the pc where the error thrown

		fmt.Printf("@@ DEBUG - handleException frame.CurrentPC() -> %v, opcode: 0x%02X (%s) \n", pc, frame.Method().Code()[pc], opcodes.OpcodeNames[frame.Method().Code()[pc]])

		// try to find handler in current frame (method) by ExceptionTable
		handlerPC := findExceptionHandler(frame, exceptionObj, pc)

		if handlerPC >= 0 { // found handler in current frame (method)
			handleCatch(frame, exceptionObj, handlerPC)
			return true
		} else { // no handler found in current frame (method)
			currentThread.PopFrame()
			if currentThread.IsStackEmpty() {
				// handler not found until popped all frames (method)
				fmt.Println("@@ DEBUG - Warning! handleException not found matching catch until popped all JVMFrameStack.")
				return false
			}
		}
	}
}

// findExceptionHandler find exception handler in method's exception table
func findExceptionHandler(frame *runtime.Frame, exceptionObj *heap.Object, pc int) int {
	method := frame.Method()
	exTable := method.ExceptionTable()

	if exTable == nil {
		return -1
	}

	// try to get ex class
	exClass := getExceptionClass(exceptionObj)
	return exTable.FindExceptionHandler(exClass, pc)
}

// handleCatch handle cache
// args:
// - frame: current method
// - exceptionObj: exception
// - handlerPC: handler entry point PC
func handleCatch(frame *runtime.Frame, exceptionObj *heap.Object, handlerPC int) {
	// we will handle this exception in this frame (method)
	stack := frame.OperandStack()
	// 1. clear op-stack (clean calculated trash result)
	stack.Clear()
	// 2. push exception into stack
	stack.PushRef(exceptionObj)
	// 3. jump to handler
	frame.SetNextPC(handlerPC)

}

// handleUncaughtException handle uncaught exception
// print ex and os.exit()
func handleUncaughtException(thread *runtime.Thread, exceptionObj *heap.Object) {
	className := getExceptionClassName(exceptionObj)
	message := getExceptionMessage(exceptionObj)

	if message != "" {
		fmt.Fprintf(os.Stderr, "Exception in thread \"main\" %s: %s\n", className, message)
	} else {
		fmt.Fprintf(os.Stderr, "Exception in thread \"main\" %s\n", className)
	}

	// TODO:
	// printStackTrace(thread, exceptionObj)

	// 終止程式
	os.Exit(1)
}

// ============================================================
// Tools
// ============================================================

// getExceptionClass get ex's class TODO: this is MVP simplify
func getExceptionClass(exceptionObj *heap.Object) *method_area.Class {
	if exceptionObj.Class() != nil {
		return exceptionObj.Class().(*method_area.Class)
	}

	return nil
}

// getExceptionClassName get ex class name TODO: this is MVP simplify
func getExceptionClassName(exceptionObj *heap.Object) string {
	exClass := getExceptionClass(exceptionObj)
	if exClass != nil {
		return exClass.Name()
	}

	fmt.Println("@@ DEBUG - getExceptionClassName can not find ex class.")

	return "java/lang/Throwable"
}

// getExceptionMessage get ex message
func getExceptionMessage(exceptionObj *heap.Object) string {
	if extra := exceptionObj.Extra(); extra != nil {
		if info, ok := extra.(*heap.ExceptionData); ok {
			return info.Message
		}
	}
	return ""
}

// ============================================================
// Java Common Exception Factory
// ============================================================

// NewArithmeticException
func NewArithmeticException(frame *runtime.Frame, message string) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/ArithmeticException", false)
	return heap.NewExceptionObject(exClass, message)
}

// NewNullPointerException
func NewNullPointerException(frame *runtime.Frame) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/NullPointerException", false)
	return heap.NewExceptionObject(exClass, "")
}

// NewArrayIndexOutOfBoundsException
func NewArrayIndexOutOfBoundsException(frame *runtime.Frame, index int32) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/ArrayIndexOutOfBoundsException", false)
	return heap.NewExceptionObject(exClass, fmt.Sprintf("%d", index))
}

// NewClassCastException
func NewClassCastException(frame *runtime.Frame, from, to string) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/ClassCastException", false)
	return heap.NewExceptionObject(exClass, fmt.Sprintf("%s cannot be cast to %s", from, to))
}

// NewNegativeArraySizeException
func NewNegativeArraySizeException(frame *runtime.Frame, size int32) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/NegativeArraySizeException", false)
	return heap.NewExceptionObject(exClass, fmt.Sprintf("%d", size))
}
