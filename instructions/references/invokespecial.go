package references

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// INVOKESPECIAL
// 1. constructor <init>
// 2. private method
// 3. parent method (super.xxx())
// opcode = 0xB7
// operands: 2 bytes (constant pool index)
// stack: [objectref, arg1, arg2, ...] → [...]
type INVOKESPECIAL struct {
	base.Index16Instruction
}

func (i *INVOKESPECIAL) Execute(frame *runtime.Frame) {
	// 1. get lang & ctcp
	currentClass := frame.Method().Class()
	rtcp := currentClass.ConstantPool()

	// 2. get methodRef
	methodRef := rtcp.GetConstant(i.Index).(*method_area.MethodRef)

	// 3. resolve target method's lang and method
	resolvedMethod := methodRef.ResolvedMethod()

	// ============================================================
	// v0.3.1: Hack - handle native method invoke for private methods
	// private native methods (like Class.initClassName) use invokespecial
	// ============================================================
	if resolvedMethod.IsNative() {
		if hacked_invoke_native(frame, methodRef) {
			return
		} else {
			fmt.Printf("@@ DEBUG - INVOKESPECIAL hacked_invoke_native failed, method: %s\n", resolvedMethod.Name())
			panic("INVOKESPECIAL Hacked invoke method failed")
		}
	}

	// resolve class
	resolvedClass := methodRef.ResolvedClass()

	// 4. check <init> must invoke by invokespecial (# README: 實作 invokespecial 時的發現)
	if resolvedMethod.Name() == "<init>" && resolvedMethod.Class() != resolvedClass {
		// why?
		// currentClass -> rtcp -> methodRef -> resolvedMethod (parent's method)
		// methodRef -> resolvedClass (this.lang)
		// this.lang != parent's method.Class()
		panic("java.lang.NoSuchMethodError")
	}

	// 5. check can not call static method
	if resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 6. get objectref(this) from stack, objectref is under args, like stack: [objectref, arg1, arg2, ...]
	objectref := frame.OperandStack().PeekRefFromTop(resolvedMethod.ArgSlotCount() - 1)
	if objectref == nil {
		panic("java.lang.NullPointerException")
	}

	// ============================================================
	// v0.3.2: Reference/ReferenceQueue Constructor Hook
	// Initialize ReferenceData/ReferenceQueueData when constructing
	// Reference or ReferenceQueue objects
	// ============================================================
	if resolvedMethod.Name() == "<init>" {
		initializeReferenceIfNeeded(frame, resolvedClass, objectref, resolvedMethod)
	}

	// 7. determine which method to call. most of all, using resolvedMethod
	methodToCall := resolvedMethod

	// 8. call super.xxx()
	// - currentClass is son fo resolvedClass
	// - currentClass has ACC_SUPER access_flag
	// - not constructor call
	if isSubClassOf(currentClass, resolvedClass) &&
		currentClass.IsSuper() && // ACC_SUPER
		resolvedMethod.Name() != "<init>" {
		// look up start from super lang.
		methodToCall = currentClass.SuperClass().GetMethod(resolvedMethod.Name(), resolvedMethod.Descriptor())
	}

	if methodToCall == nil {
		panic("java.lang.AbstractMethodError")
	}

	// 9. invoke method
	invokeMethod(frame, methodToCall)
}

func (i *INVOKESPECIAL) Opcode() uint8 {
	return 0xB7
}

// ============================================================
// v0.3.2: Reference Initialization Hook
// ============================================================

// initializeReferenceIfNeeded checks if the object being constructed is a
// Reference or ReferenceQueue, and initializes the corresponding data structure.
//
// This is called BEFORE the constructor runs, so we need to extract
// constructor arguments from the operand stack.
//
// Reference constructors:
//   - SoftReference(T referent)
//   - SoftReference(T referent, ReferenceQueue queue)
//   - WeakReference(T referent)
//   - WeakReference(T referent, ReferenceQueue queue)
//   - PhantomReference(T referent, ReferenceQueue queue)  // queue param is mandatory here!!
//
// ReferenceQueue constructors:
//   - ReferenceQueue()  // no arguments
func initializeReferenceIfNeeded(frame *runtime.Frame, class *method_area.Class, objectref interface{}, method *method_area.Method) {
	// Check if it's a Reference or ReferenceQueue class
	if method_area.IsReferenceQueueClass(class) {
		// Initialize ReferenceQueue<T>
		obj := objectref.(*heap.Object)
		initializeReferenceQueue(obj)
		return
	}

	if method_area.IsReferenceClass(class) {
		// Initialize Reference<T>
		obj := objectref.(*heap.Object)
		initializeReference(frame, class, obj, method)
		return
	}
}

// initializeReferenceQueue initializes ReferenceQueueData for a ReferenceQueue object
func initializeReferenceQueue(obj *heap.Object) {
	// Only initialize if not already initialized
	if obj.GetReferenceQueueData() != nil {
		return
	}

	fmt.Printf("@@ Debug - [v0.3.2] Initializing ReferenceQueue\n")
	method_area.InitializeReferenceQueueObject(obj)
}

// initializeReference initializes ReferenceData for a Reference object
func initializeReference(frame *runtime.Frame, class *method_area.Class, obj *heap.Object, method *method_area.Method) {
	// Only initialize if not already initialized
	if obj.GetReferenceData() != nil {
		return
	}

	// Parse constructor to determine parameters
	constructorInfo := method_area.ParseReferenceConstructor(class.Name(), method.Descriptor())
	if constructorInfo == nil {
		return
	}

	// Extract arguments from operand stack
	// Stack layout: [objectref, referent, (optional)queue]
	// ArgSlotCount includes 'this', so:
	//   - 2 args (this + referent): ArgSlotCount = 2
	//   - 3 args (this + referent + queue): ArgSlotCount = 3
	stack := frame.OperandStack()
	argCount := method.ArgSlotCount()

	var referent *heap.Object
	var queue *heap.Object

	if constructorInfo.HasQueue {
		// Constructor: Reference(T referent, ReferenceQueue queue)
		// Stack: [..., objectref, referent, queue]
		// Peek from top: queue is at top-1 (index 0), referent is at top-2 (index 1)
		queueRef := stack.PeekRefFromTop(0)
		if queueRef != nil {
			queue = queueRef.(*heap.Object)
		}
		referentRef := stack.PeekRefFromTop(1)
		if referentRef != nil {
			referent = referentRef.(*heap.Object)
		}
	} else {
		// Constructor: Reference(T referent)
		// Stack: [..., objectref, referent]
		// Peek from top: referent is at top-1 (index 0)
		referentRef := stack.PeekRefFromTop(0)
		if referentRef != nil {
			referent = referentRef.(*heap.Object)
		}
	}

	fmt.Printf("@@ Debug - [v0.3.2] Initializing %s (referent=%v, queue=%v, argCount=%d)\n",
		constructorInfo.RefType, referent != nil, queue != nil, argCount)

	// Initialize the Reference
	method_area.InitializeReferenceObject(obj, referent, queue)
}
