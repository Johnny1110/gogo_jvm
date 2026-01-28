package exception

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

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

func NewIncompatibleClassChangeError(frame *runtime.Frame, message string) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/IncompatibleClassChangeError", false)
	return heap.NewExceptionObject(exClass, fmt.Sprintf("%s", message))
}

func NewAbstractMethodError(frame *runtime.Frame, className, methodName, descriptor string) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/NewAbstractMethodError", false)
	return heap.NewExceptionObject(exClass, fmt.Sprintf("%s.%s%s", className, methodName, descriptor))
}

func NewIllegalAccessError(frame *runtime.Frame, className, methodName, descriptor string) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/IllegalAccessError", false)
	return heap.NewExceptionObject(exClass, fmt.Sprintf("%s.%s%s", className, methodName, descriptor))
}

func NewCloneNotSupportedException(frame *runtime.Frame, classname string) *heap.Object {
	exClass := frame.Method().Class().Loader().LoadClass("java/lang/CloneNotSupportedException", false)
	return heap.NewExceptionObject(exClass, fmt.Sprintf("%s.%s", classname, "clone()"))
}

func NewExceptionObject(exClass *method_area.Class, msg string) *heap.Object {
	return heap.NewExceptionObject(exClass, msg)
}
