package method_area

import (
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"strings"
)

// ============================================================
// Reference Type Integration - v0.3.2
// ============================================================
// This file contains helper functions to integrate Reference types
// (SoftReference, WeakReference, PhantomReference, ReferenceQueue)
// into the class loading and object creation process.
//
// Key integration points:
// 1. IsReferenceClass() - Check if a class is a Reference type
// 2. IsReferenceQueueClass() - Check if a class is ReferenceQueue
// 3. InitializeReferenceObject() - Initialize ReferenceData in Object.extra
// 4. InitializeReferenceQueueObject() - Initialize ReferenceQueueData in Object.extra

// ============================================================
// Reference Class Detection
// ============================================================

// Reference class names (JVM internal format)
const (
	ClassNameReference        = "java/lang/ref/Reference"
	ClassNameSoftReference    = "java/lang/ref/SoftReference"
	ClassNameWeakReference    = "java/lang/ref/WeakReference"
	ClassNamePhantomReference = "java/lang/ref/PhantomReference"
	ClassNameReferenceQueue   = "java/lang/ref/ReferenceQueue"
)

// IsReferenceClass checks if the class is a Reference type or subclass
// Returns true for: Reference, SoftReference, WeakReference, PhantomReference
func IsReferenceClass(class *Class) bool {
	if class == nil {
		return false
	}

	// Check class name directly
	switch class.name {
	case ClassNameReference,
		// actually the refs below are all extended from ClassNameReference (java/lang/ref/Reference)
		ClassNameSoftReference,
		ClassNameWeakReference,
		ClassNamePhantomReference:
		return true
	}

	// Check if class extends Reference (for custom Reference subclasses)
	return isSubclassOfReference(class)
}

// isSubclassOfReference checks if class is a subclass of java.lang.ref.Reference
func isSubclassOfReference(class *Class) bool {
	return isSubclassOfName(class, ClassNameReference)
}

// IsReferenceQueueClass checks if the class is ReferenceQueue
func IsReferenceQueueClass(class *Class) bool {
	return isSubclassOfName(class, ClassNameReferenceQueue)
}

// GetReferenceType returns the ReferenceType for a given class
// Returns RefTypeNone if not a Reference class
func GetReferenceType(class *Class) heap.ReferenceType {
	if class == nil {
		return heap.RefTypeNone
	}

	// Check by class name
	switch class.name {
	case ClassNameSoftReference:
		return heap.RefTypeSoft
	case ClassNameWeakReference:
		return heap.RefTypeWeak
	case ClassNamePhantomReference:
		return heap.RefTypePhantom
	}

	// If not found check parent class for subclasses.
	for c := class.superClass; c != nil; c = c.superClass {
		switch c.name {
		case ClassNameSoftReference:
			return heap.RefTypeSoft
		case ClassNameWeakReference:
			return heap.RefTypeWeak
		case ClassNamePhantomReference:
			return heap.RefTypePhantom
		}
	}

	// default:
	return heap.RefTypeNone
}

// ============================================================
// Reference Object Initialization
// ============================================================

// InitializeReferenceObject initializes the ReferenceData for a Reference object
// This should be called when the Reference constructor is invoked
//
// Parameters:
//   - refObj: the Reference object being constructed
//   - referent: the object being referenced (T in Reference<T>)
//   - queue: the ReferenceQueue (can be nil)
//
// The function determines the Reference type from refObj's class and
// creates the appropriate ReferenceData.
func InitializeReferenceObject(refObj *heap.Object, referent *heap.Object, queue *heap.Object) {
	if refObj == nil {
		return
	}

	// Get the class
	class := refObj.Class().(*Class)
	refType := GetReferenceType(class)

	if refType == heap.RefTypeNone {
		// Not a Reference type, shouldn't happen
		return
	}

	// Create and set ReferenceData
	refData := heap.NewReferenceData(refType, referent, queue)
	refObj.SetExtra(refData)
}

// InitializeReferenceQueueObject initializes the ReferenceQueueData for a ReferenceQueue object
// This should be called when the ReferenceQueue constructor is invoked
func InitializeReferenceQueueObject(queueObj *heap.Object) {
	if queueObj == nil {
		return
	}

	// Create and set ReferenceQueueData
	queueData := heap.NewReferenceQueueData()
	queueObj.SetReferenceQueueData(queueData)
}

// ============================================================
// Constructor Argument Extraction Helpers
// ============================================================
// These help extract constructor arguments for Reference initialization.
// Reference constructors have two forms:
//   - Reference(T referent)
//   - Reference(T referent, ReferenceQueue<? super T> queue)

// ReferenceConstructorInfo holds information about a Reference constructor
type ReferenceConstructorInfo struct {
	HasQueue bool               // true if constructor has queue parameter
	RefType  heap.ReferenceType // type of reference
}

// ParseReferenceConstructor parses a constructor descriptor to determine
// if it's a Reference constructor and what parameters it has
//
// Reference constructors:
//   - SoftReference(T): (Ljava/lang/Object;)V
//   - SoftReference(T, ReferenceQueue): (Ljava/lang/Object;Ljava/lang/ref/ReferenceQueue;)V
//   - WeakReference(T): (Ljava/lang/Object;)V
//   - WeakReference(T, ReferenceQueue): (Ljava/lang/Object;Ljava/lang/ref/ReferenceQueue;)V
//   - PhantomReference(T, ReferenceQueue): (Ljava/lang/Object;Ljava/lang/ref/ReferenceQueue;)V
//     Note: PhantomReference ONLY has the two-argument constructor
func ParseReferenceConstructor(className string, descriptor string) *ReferenceConstructorInfo {
	refType := heap.RefTypeNone

	switch className {
	case ClassNameSoftReference:
		refType = heap.RefTypeSoft
	case ClassNameWeakReference:
		refType = heap.RefTypeWeak
	case ClassNamePhantomReference:
		refType = heap.RefTypePhantom
	default:
		return nil
	}

	// Check descriptor for queue parameter
	hasQueue := strings.Contains(descriptor, "Ljava/lang/ref/ReferenceQueue;")

	return &ReferenceConstructorInfo{
		HasQueue: hasQueue,
		RefType:  refType,
	}
}

// IsReferenceQueueConstructor checks if a method is a ReferenceQueue constructor
func IsReferenceQueueConstructor(className string, methodName string) bool {
	return className == ClassNameReferenceQueue && methodName == "<init>"
}

// ------------------------------------------------------------------------
// Support Methods
// ------------------------------------------------------------------------

// isSubclassOfName checks if class is a subclass of `ofClassName`
func isSubclassOfName(class *Class, ofClassName string) bool {
	for c := class.superClass; c != nil; c = c.superClass {
		if c.name == ofClassName {
			return true
		}
	}

	return false
}
