package ref

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// ============================================================
// java.lang.ref Native Methods - v0.3.2
// ============================================================

// Reference<T> is the abstract base class for all reference types.
// Key methods:
//   - get(): returns the referent (or null if cleared)
//   - clear(): clears the referent
//   - enqueue(): adds this reference to its queue
//   - isEnqueued(): checks if this reference is in a queue

func init() {
	fmt.Println("@@ Debug - init Native java/lang/ref/Reference")

	// Core Reference methods
	runtime.Register("java/lang/ref/Reference", "get", "()Ljava/lang/Object;", referenceGet)
	runtime.Register("java/lang/ref/Reference", "clear", "()V", referenceClear)
	runtime.Register("java/lang/ref/Reference", "clear0", "()V", referenceClear) // JDK internal variant
	runtime.Register("java/lang/ref/Reference", "enqueue", "()Z", referenceEnqueue)
	runtime.Register("java/lang/ref/Reference", "isEnqueued", "()Z", referenceIsEnqueued)

	// JDK 9+ uses refersTo for phantom reference check
	runtime.Register("java/lang/ref/Reference", "refersTo", "(Ljava/lang/Object;)Z", referenceRefersTo)
	runtime.Register("java/lang/ref/Reference", "refersTo0", "(Ljava/lang/Object;)Z", referenceRefersTo)

	// Clone is not supported for Reference
	runtime.Register("java/lang/ref/Reference", "clone", "()Ljava/lang/Object;", referenceClone)
}

// ============================================================
// get - Reference.get()
// ============================================================
// Java signature: public T get()
//
// Returns the referent, or null if:
// - The referent has been cleared (by GC or clear())
// - This is a PhantomReference (always returns null)
//
// For SoftReference, this also updates the timestamp for LRU tracking.
//
// Stack:
//
//	[this] → [referent or null]
func referenceGet(frame *runtime.Frame) {

	fmt.Printf("@@ DEBUG - native referenceGet in ... \n")

	this := frame.LocalVars().GetThis()
	if this == nil {
		frame.OperandStack().PushRef(nil)
		return
	}

	obj := this.(*heap.Object)
	refData := obj.GetReferenceData()
	if refData == nil {
		// Not properly initialized as a Reference
		frame.OperandStack().PushRef(nil)
		return
	}

	// Use ReferenceData.Get() which handles:
	// - PhantomReference always returning null
	// - SoftReference timestamp update
	referent := refData.Get()

	fmt.Printf("@@ DEBUG - native referenceGet > referent: %v \n", referent)

	frame.OperandStack().PushRef(referent)
}

// ============================================================
// clear - Reference.clear()
// ============================================================
// Java signature: public void clear()
//
// Clears the referent. After this call, get() will return null.
// This method can be called by:
// - User code explicitly
// - GC when the referent becomes unreachable
//
// Note: Clearing does NOT enqueue the reference. The reference
// must still be explicitly enqueued or will be enqueued by GC.
//
// Stack:
//
//	[this] → []
func referenceClear(frame *runtime.Frame) {
	fmt.Printf("@@ DEBUG - native referenceClear in...\n")

	this := frame.LocalVars().GetThis()
	if this == nil {
		fmt.Printf("@@ DEBUG - native referenceClear > ref is null (high level error)...\n")
		return
	}

	obj := this.(*heap.Object)
	refData := obj.GetReferenceData()

	if refData == nil {
		fmt.Printf("@@ DEBUG - native referenceClear > refData is null (high level error)...\n")
		return
	}

	// Referent = nil

	obj.Fields()[0].Ref = nil // TODO: ? test clean first filed (java referent)

	refData.Clear()
}

// ============================================================
// enqueue - Reference.enqueue()
// ============================================================
// Java signature: public boolean enqueue()
//
// Adds this reference to its registered queue, if:
// - A queue was registered at construction time
// - The reference hasn't already been enqueued
//
// Returns true if successfully enqueued, false otherwise.
//
// Stack:
//
//	[this] → [boolean]
func referenceEnqueue(frame *runtime.Frame) {
	if refData, ok := getThisReferenceData(frame); ok {
		// Check if queue is registered
		if refData.Queue == nil {
			frame.OperandStack().PushFalse()
			return
		}

		// Check current state - can only enqueue from Active or Pending
		if refData.State != heap.RefStateActive && refData.State != heap.RefStatePending {
			frame.OperandStack().PushFalse()
			return
		}

		// Get the queue data and enqueue
		queueData := refData.Queue.GetReferenceQueueData()
		if queueData == nil {
			frame.OperandStack().PushFalse()
			return
		}

		this, _ := getThis(frame)
		frame.OperandStack().PushBoolean(queueData.Enqueue(this))

	} else {
		frame.OperandStack().PushFalse()
		return
	}
}

// ============================================================
// isEnqueued - Reference.isEnqueued()
// ============================================================
// Java signature: public boolean isEnqueued()
//
// Returns true if this reference has been enqueued.
// Note: Once removed from the queue, this returns false.
//
// Stack:
//
//	[this] → [boolean]
func referenceIsEnqueued(frame *runtime.Frame) {
	if refData, ok := getThisReferenceData(frame); ok {
		if refData.State == heap.RefStateEnqueued {
			frame.OperandStack().PushTrue()
			return
		} else {
			frame.OperandStack().PushFalse()
			return
		}
	} else {
		frame.OperandStack().PushFalse()
		return
	}
}

// ============================================================
// refersTo - Reference.refersTo(Object)
// ============================================================
// Java signature: public final boolean refersTo(T obj) (JDK 16+)
//
// Tests if this reference refers to the given object.
// This is useful for PhantomReference where get() always returns null.
//
// Stack:
//
//	[this, obj] → [boolean]
func referenceRefersTo(frame *runtime.Frame) {
	refData, ok := getThisReferenceData(frame)
	if !ok {
		frame.OperandStack().PushFalse()
		return
	}

	obj := frame.LocalVars().GetRef(1)

	if refData.Referent == obj {
		frame.OperandStack().PushTrue()
		return
	} else {
		frame.OperandStack().PushFalse()
		return
	}
}

// ============================================================
// clone - Reference.clone()
// ============================================================
// Java signature: protected Object clone() throws CloneNotSupportedException
//
// Reference objects cannot be cloned. Always throws CloneNotSupportedException.
//
// Stack:
//
//	[this] → throws CloneNotSupportedException
func referenceClone(frame *runtime.Frame) {
	// Reference objects cannot be cloned
	// In real JVM, this throws CloneNotSupportedException
	// For MVP, we panic with a descriptive message
	panic("java.lang.CloneNotSupportedException: Reference objects cannot be cloned")
}

// ------------------------------------------------------------
// private methods
// ------------------------------------------------------------
func getThis(frame *runtime.Frame) (*heap.Object, bool) {
	thisIface := frame.LocalVars().GetThis()
	if thisIface == nil {
		return nil, false
	}
	return thisIface.(*heap.Object), true
}

func getThisReferenceData(frame *runtime.Frame) (*heap.ReferenceData, bool) {
	if this, ok := getThis(frame); ok {
		refData := this.GetReferenceData()
		if refData != nil {
			return refData, true
		} else {
			return nil, false
		}
	} else {
		return nil, false
	}
}
