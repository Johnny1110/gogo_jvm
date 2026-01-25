package heap

import (
	"fmt"
	"time"
)

// ============================================================
// Reference Type - v0.3.2
// ============================================================
// Java Reference Type: java.lang.ref.Reference
//
// Reference hierarchy:
//   Reference<T> (abstract base)
//   ├── SoftReference<T>   - GC when memory insufficient
//   ├── WeakReference<T>   - GC at next collection
//   └── PhantomReference<T> - get() always returns null
//
// Reference Object uses Object.extra to store ReferenceData
// This is similar to how we handle arrays (extra = []int32) and
// exceptions (extra = *ExceptionData)

// ============================================================
// Reference Type Constants
// ============================================================

// ReferenceType represents the type of reference
type ReferenceType uint8

const (
	RefTypeNone    ReferenceType = 0 // Not a reference type
	RefTypeSoft    ReferenceType = 1 // SoftReference
	RefTypeWeak    ReferenceType = 2 // WeakReference
	RefTypePhantom ReferenceType = 3 // PhantomReference
)

// String returns the string representation of ReferenceType
func (rt ReferenceType) String() string {
	switch rt {
	case RefTypeSoft:
		return "SoftReference"
	case RefTypeWeak:
		return "WeakReference"
	case RefTypePhantom:
		return "PhantomReference"
	default:
		return "Unknown"
	}
}

// ============================================================
// Reference State Constants
// ============================================================

// ReferenceState represents the lifecycle state of a reference
type ReferenceState uint8

const (
	// RefStateActive - Initial state when reference is created
	// referent is accessible via get() (except PhantomReference)
	RefStateActive ReferenceState = 0

	// RefStatePending - GC has cleared the referent, waiting to be enqueued
	// referent = nil, waiting for Reference Handler to enqueue
	RefStatePending ReferenceState = 1

	// RefStateEnqueued - Reference is in the queue
	// Can be retrieved via ReferenceQueue.poll() or remove()
	RefStateEnqueued ReferenceState = 2

	// RefStateInactive - Final state, reference is no longer useful
	// Either: no queue registered, or already removed from queue
	RefStateInactive ReferenceState = 3
)

// String returns the string representation of ReferenceState
func (rs ReferenceState) String() string {
	switch rs {
	case RefStateActive:
		return "Active"
	case RefStatePending:
		return "Pending"
	case RefStateEnqueued:
		return "Enqueued"
	case RefStateInactive:
		return "Inactive"
	default:
		return "Unknown"
	}
}

// ============================================================
// ReferenceData - Stored in Object.extra
// ============================================================
// ReferenceData stores reference-specific data in Object.extra
// This structure is created when a Reference object is instantiated
//
// Memory Layout:
//
//	┌─────────────────────────────────────────────────────────────┐
//	│  Reference Object (heap.Object)                             │
//	├─────────────────────────────────────────────────────────────┤
//	│  markWord: uint64         (Object header)                   │
//	│  class:    interface{}    (→ java/lang/ref/WeakReference)   │
//	│  fields:   Slots          (Java fields, if any)             │
//	│  extra:    interface{}    (→ *ReferenceData) <-- HERE       │
//	└─────────────────────────────────────────────────────────────┘

type ReferenceData struct {
	// RefType indicates the type of reference (Soft/Weak/Phantom)
	RefType ReferenceType
	// State indicates the current lifecycle state
	State ReferenceState
	// Referent is the object being referenced
	// This is the "T" in Reference<T>
	// Will be set to nil when GC clears the reference
	Referent *Object
	// Queue is the ReferenceQueue this reference is registered with
	// Can be nil if no queue was provided in constructor
	Queue *Object
	// Next is used for linking references in ReferenceQueue
	// Forms a singly-linked list within the queue
	Next *Object
	// ========== SoftReference specific fields ==========
	// Timestamp records the last access time (for LRU in SoftReference)
	// Used by GC to decide which SoftReferences to clear first
	// Only meaningful for RefTypeSoft
	Timestamp int64
}

// ============================================================
// ReferenceData Factory Functions
// ============================================================

// NewReferenceData creates a new ReferenceData for a Reference object
func NewReferenceData(refType ReferenceType, referent *Object, queue *Object) *ReferenceData {
	data := &ReferenceData{
		RefType:  refType,
		State:    RefStateActive,
		Referent: referent,
		Queue:    queue,
		Next:     nil,
	}
	// Initialize timestamp for SoftReference (LRU tracking)
	if refType == RefTypeSoft {
		data.Timestamp = time.Now().UnixNano()
	}

	return data
}

// NewSoftReferenceData creates ReferenceData for SoftReference
func NewSoftReferenceData(referent *Object, queue *Object) *ReferenceData {
	return NewReferenceData(RefTypeSoft, referent, queue)
}

// NewWeakReferenceData creates ReferenceData for WeakReference
func NewWeakReferenceData(referent *Object, queue *Object) *ReferenceData {
	return NewReferenceData(RefTypeWeak, referent, queue)
}

// NewPhantomReferenceData creates ReferenceData for PhantomReference
// Note: PhantomReference MUST have a queue (enforced at Java level)
func NewPhantomReferenceData(referent *Object, queue *Object) *ReferenceData {
	return NewReferenceData(RefTypePhantom, referent, queue)
}

// ============================================================
// ReferenceData Methods
// ============================================================

// Get returns the referent object
// For PhantomReference, always returns nil (per Java spec)
func (rd *ReferenceData) Get() *Object {
	if rd.RefType == RefTypePhantom {
		return nil // PhantomReference always return nil
	}

	// Update timestamp for SoftReference (LRU tracking)
	if rd.RefType == RefTypeSoft && rd.Referent != nil {
		rd.Timestamp = time.Now().UnixNano()
	}

	return rd.Referent
}

// Clear sets the referent to nil
// This can be called explicitly by user code or by GC
func (rd *ReferenceData) Clear() {
	fmt.Printf("@@ DEBUG - heap/reference Clear in..\n")
	rd.Referent = nil
}

// IsEnqueued returns true if this reference is in the queue
func (rd *ReferenceData) IsEnqueued() bool {
	return rd.State == RefStateEnqueued
}

// ============================================================
// Object Extension Methods for Reference
// ============================================================

// IsReference checks if this object is a Reference type
// Returns true if Object.extra contains *ReferenceData
func (o *Object) IsReference() bool {
	if o.extra == nil {
		return false
	}

	_, ok := o.extra.(*ReferenceData)
	return ok
}

// GetReferenceData returns the ReferenceData if this is a Reference object
// Returns nil if not a Reference
func (o *Object) GetReferenceData() *ReferenceData {
	if o.extra == nil {
		return nil
	}

	if data, ok := o.extra.(*ReferenceData); ok {
		return data
	} else {
		return nil
	}
}

// SetReferenceData sets the ReferenceData for this object
// Used during Reference object construction
func (o *Object) SetReferenceData(data *ReferenceData) {
	o.extra = data
}
