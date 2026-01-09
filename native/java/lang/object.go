package lang

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// Object class's native methods (hashCode etc.)

func init() {
	fmt.Println("@@ Debug - init Native java/lang/Object")
	// v0.3.0
	runtime.Register("java/lang/Object", "hashCode", "()I", objectHashCode)
	// v0.3.1
	runtime.Register("java/lang/Object", "getClass", "()Ljava/lang/Class;", objectGetClass)
	// (MVP Clone)
	runtime.Register("java/lang/Object", "clone", "()Ljava/lang/Object;", objectClone)
	// TODO: v0.4.x: notify/notifyAll/wait
	runtime.Register("java/lang/Object", "notify", "()V", objectNotify)
	runtime.Register("java/lang/Object", "notifyAll", "()V", objectNotifyAll)
	runtime.Register("java/lang/Object", "wait", "(J)V", objectWait)
}

// ============================================================
// hashCode - Object.hashCode()
// ============================================================
// objectHashCode v0.3.0 - native implementation of Object.hashCode()
// Java signature: public native int hashCode();
// return Object's identity hash code (store in Object's markWord)
func objectHashCode(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis()
	if this == nil {
		// should not be happened
		frame.OperandStack().PushInt(0)
		return
	}

	obj := this.(*heap.Object)
	hash := obj.HashCode()
	frame.OperandStack().PushInt(hash)
}

// ============================================================
// getClass - Object.getClass()
// ============================================================
// Java signature: public final native Class<?> getClass();
// return runtime class type
//
// Reflection Entry:
//
//	Object obj = new String("hello");
//	Class<?> c = obj.getClass();  // -> java.lang.String  Class (Object)
func objectGetClass(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis()
	if this == nil {
		panic("java.lang.NullPointerException")
	}

	obj := this.(*heap.Object)

	// class metadata
	class := obj.Class().(*method_area.Class)
	// get jClass (Object) from class
	fmt.Printf("@@ DEBUG - Native objectGetClass  obj.Class() = %s \n", class.Name())
	jClass := class.JClass()
	if jClass == nil {
		panic(fmt.Sprintf("Class object not initialized for: %s", class.Name()))
	}

	// return java.lang.Class (Object)
	frame.OperandStack().PushRef(jClass)
}

// ============================================================
// clone - Object.clone()
// ============================================================
// Java signature: protected native Object clone() throws CloneNotSupportedException;
// TODO: MVP 簡化版 - 只複製基本物件，不處理 Cloneable 檢查
func objectClone(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis()
	if this == nil {
		panic("java.lang.NullPointerException")
	}

	obj := this.(*heap.Object)
	class := obj.Class().(*method_area.Class)

	// TODO: 檢查是否實現 Cloneable 介面 如果沒有實現，應該拋出 CloneNotSupportedException MVP 簡化：跳過檢查
	cloned := class.NewObject()

	// copy all fields (shallow copy)
	srcFields := obj.Fields()
	dstFields := cloned.Fields()

	if srcFields != nil && dstFields != nil {
		copy(dstFields, srcFields)
	}

	// if is array: copy extra
	if obj.IsArray() {
		cloned.SetExtra(cloneArray(obj))
	}

	frame.OperandStack().PushRef(cloned)
}

// cloneArray copy array
func cloneArray(arr *heap.Object) interface{} {
	switch data := arr.Extra().(type) {
	case []int8:
		cloned := make([]int8, len(data))
		copy(cloned, data)
		return cloned
	case []int16:
		cloned := make([]int16, len(data))
		copy(cloned, data)
		return cloned
	case []int32:
		cloned := make([]int32, len(data))
		copy(cloned, data)
		return cloned
	case []int64:
		cloned := make([]int64, len(data))
		copy(cloned, data)
		return cloned
	case []uint16:
		cloned := make([]uint16, len(data))
		copy(cloned, data)
		return cloned
	case []float32:
		cloned := make([]float32, len(data))
		copy(cloned, data)
		return cloned
	case []float64:
		cloned := make([]float64, len(data))
		copy(cloned, data)
		return cloned
	case []*heap.Object:
		cloned := make([]*heap.Object, len(data))
		copy(cloned, data) // !!! shallow copy, only copy ref.
		return cloned
	default:
		return nil
	}
}

// ============================================================
// Synchronization Methods - TODO: v0.4.x
// ============================================================

// objectNotify - Object.notify()
// Java signature: public final native void notify();
func objectNotify(frame *runtime.Frame) {
	fmt.Println("Warning - Object.notify() not implemented (v0.4.x)")
}

// objectNotifyAll - Object.notifyAll()
// Java signature: public final native void notifyAll();
func objectNotifyAll(frame *runtime.Frame) {
	fmt.Println("Warning - Object.notifyAll() not implemented (v0.4.x)")
}

// objectWait - Object.wait(long timeout)
// Java signature: public final native void wait(long timeout);
func objectWait(frame *runtime.Frame) {
	fmt.Println("Warning - Object.wait() not implemented (v0.4.x)")
}
