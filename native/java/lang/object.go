package lang

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/exception"
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

	// v0.3.3: equals (non-native in standard Java, but we provide default implementation)
	// Note: In standard Java, equals() is not native, but we intercept it for efficiency
	// Actually equals() is a normal Java method, so we don't register it here
	// The default Object.equals() in Java bytecode will be executed normally

	// v0.3.3: clone (with Cloneable check)
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
	fmt.Printf("@@ DEBUG - Native objectGetClass obj = %s \n", obj)
	fmt.Printf("@@ DEBUG - Native objectGetClass obj.Extra() = %v \n", obj.Extra())

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
//
// v0.3.3 Implementation:
// 1. Check if the object's class implements Cloneable interface
// 2. If not Cloneable, throw CloneNotSupportedException
// 3. Create a shallow copy of the object
// 4. For arrays, copy the array data (still shallow for reference arrays)
//
// Contract:
// - x.clone() != x (must be a new object)
// - x.clone().getClass() == x.getClass() (same type)
// - x.clone().equals(x) is typically true (but not required)
func objectClone(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis()
	if this == nil {
		frame.ThrowException(exception.NewNullPointerException(frame))
		return
	}

	obj := this.(*heap.Object)
	class := obj.Class().(*method_area.Class)

	// v0.3.3: Check if implements Cloneable interface
	// Arrays are always Cloneable (JLS requirement)
	if !class.IsCloneable() {
		frame.ThrowException(exception.NewCloneNotSupportedException(frame, class.Name()))
		return
	}

	var cloned *heap.Object

	// Handle array cloning specially
	if obj.IsArray() {
		cloned = cloneArrayObject(obj, class)
	} else {
		// Clone regular object
		cloned = cloneRegularObject(obj, class)
	}

	frame.OperandStack().PushRef(cloned)
}

// cloneRegularObject creates a shallow copy of a regular (non-array) object
func cloneRegularObject(obj *heap.Object, class *method_area.Class) *heap.Object {
	cloned := class.NewObject()

	// Copy all fields (shallow copy)
	// This copies primitive values directly, and copies references (not the objects they point to)
	srcFields := obj.Fields()
	dstFields := cloned.Fields()

	if srcFields != nil && dstFields != nil {
		copy(dstFields, srcFields)
	} else {
		// should not be happened.
		panic("Clone object not initialized")
	}

	// Copy markWord state (except hashCode which should be regenerated)
	// Actually, the cloned object should have its own identity hashCode
	// So we don't copy the markWord - the new object starts with InitialMarkWord
	return cloned
}

// cloneArrayObject creates a shallow copy of an array object
func cloneArrayObject(arr *heap.Object, class *method_area.Class) *heap.Object {
	cloned := &heap.Object{}
	cloned.SetMarkWord(heap.InitialMarkWord)
	cloned.SetClass(class)
	// Copy array data based on type
	cloned.SetExtra(cloneArrayData(arr))

	return cloned
}

// cloneArrayData copy array
func cloneArrayData(arr *heap.Object) interface{} {
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
		panic(fmt.Sprintf("Unknown array type: %T", data))
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
