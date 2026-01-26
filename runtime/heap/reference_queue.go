package heap

import "sync"

// ============================================================
// ReferenceQueue - v0.3.2
// ============================================================
// Java: java.lang.ref.ReferenceQueue<T>
//
// ReferenceQueue is a queue for Reference objects that have been
// cleared by the garbage collector. When a Reference's referent
// becomes eligible for GC, the Reference is enqueued here.
//
// Usage Pattern:
//
//	ReferenceQueue<Object> queue = new ReferenceQueue<>();
//	WeakReference<Object> ref = new WeakReference<>(obj, queue);
//
//	Later, after GC clears the reference:
//	Reference<? extends Object> clearedRef = queue.poll();
//	if (clearedRef != null) {
//	    // Handle cleanup
//	}
//
// Implementation:
// ReferenceQueue uses a singly-linked list of Reference objects.
// The linking is done through ReferenceData.Next field.
//
//	┌─────────────────────────────────────────────────────────────┐
//	│  ReferenceQueue Object (heap.Object)                        │
//	├─────────────────────────────────────────────────────────────┤
//	│  extra: *ReferenceQueueData                                 │
//	│         ┌─────────────────────────────────────────────────┐ │
//	│         │  Head ──► Ref1 ──► Ref2 ──► Ref3 ──► nil        │ │
//	│         │  Tail ───────────────────────┘                  │ │
//	│         │  Length: 3                                      │ │
//	│         └─────────────────────────────────────────────────┘ │
//	└─────────────────────────────────────────────────────────────┘

// ============================================================
// ReferenceQueueData - Stored in Object.extra
// ============================================================

// ReferenceQueueData stores the queue data structure
type ReferenceQueueData struct {
	// Head points to the first Reference in the queue
	Head *Object
	// Tail points to the last Reference in the queue
	// Used for O(1) enqueue operation
	Tail *Object
	// Length is the number of References in the queue
	Length int
	// Lock provides thread-safety for queue operations
	// In real JVM, this uses more sophisticated synchronization
	lock sync.Mutex
}

// ============================================================
// ReferenceQueueData Factory
// ============================================================

// NewReferenceQueueData creates a new empty ReferenceQueueData
func NewReferenceQueueData() *ReferenceQueueData {
	return &ReferenceQueueData{
		Head:   nil,
		Tail:   nil,
		Length: 0,
	}
}

// ============================================================
// ReferenceQueueData Methods
// ============================================================

// Enqueue adds a Reference to the queue
// Returns true if successfully enqueued, false if already enqueued or inactive
//
// Queue Structure After Enqueue:
//
//	Before: Head → [A] → [B] → nil, Tail → [B]
//	After:  Head → [A] → [B] → [C] → nil, Tail → [C]
func (q *ReferenceQueueData) Enqueue(ref *Object) bool {
	if ref == nil {
		return false
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	refData := ref.GetReferenceData()
	if refData == nil {
		return false
	}

	// Can only enqueue from Active or Pending state
	if refData.State != RefStateActive && refData.State != RefStatePending {
		return false
	}

	// Add to tail of queue
	refData.Next = nil
	if q.Tail == nil {
		// queue is empty.
		q.Head = ref
		q.Tail = ref
	} else {
		tail := q.Tail.GetReferenceData()
		if tail != nil {
			tail.Next = ref
		}
		q.Tail = ref // ref will be new tail
	}

	// Update state
	refData.State = RefStateEnqueued
	q.Length++

	return true
}

// Poll removes and returns the head of the queue, or nil if empty
// This is a non-blocking operation
//
// Queue Structure After Poll:
//
//	Before: Head → [A] → [B] → [C] → nil, Tail → [C]
//	After:  Head → [B] → [C] → nil, Tail → [C]
//	Returns: [A]
func (q *ReferenceQueueData) Poll() *Object {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.Head == nil {
		return nil
	}

	ref := q.Head
	refData := ref.GetReferenceData() // ref.data will be return

	if refData != nil {
		q.Head = refData.Next
		if q.Head == nil { // already empty
			q.Tail = nil
		}

		// Update state and clear Next pointer
		refData.State = RefStateInactive
		refData.Next = nil
	} else {
		// Corrupted reference, just remove it
		q.Head = nil
		q.Tail = nil
	}

	q.Length--

	if q.Length < 0 {
		q.Length = 0
	}

	return ref
}

// Remove removes and returns the head of the queue
// In real JVM, this can block with optional timeout
// For MVP, we implement it as non-blocking (same as Poll)
//
// Parameters:
//   - timeout: timeout in milliseconds (0 = no wait, -1 = wait forever)
//
// Returns: Reference object or nil if queue is empty/timeout
func (q *ReferenceQueueData) Remove(timeout int64) *Object {
	// MVP simplification: non-blocking implementation
	// TODO: v0.4.x - implement proper blocking with wait/notify
	return q.Poll()
}

// IsEmpty returns true if the queue has no elements
func (q *ReferenceQueueData) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.Length == 0
}

// Size returns the number of elements in the queue
func (q *ReferenceQueueData) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.Length
}

// ============================================================
// Object Extension Methods for ReferenceQueue
// ============================================================

// IsReferenceQueue checks if this object is a ReferenceQueue
func (o *Object) IsReferenceQueue() bool {
	if o.extra == nil {
		return false
	}

	_, ok := o.extra.(*ReferenceQueueData)
	return ok
}

// GetReferenceQueueData returns the ReferenceQueueData if this is a ReferenceQueue
func (o *Object) GetReferenceQueueData() *ReferenceQueueData {
	if o.extra == nil {
		return nil
	}

	if t, ok := o.extra.(*ReferenceQueueData); ok {
		return t
	}

	return nil
}

// SetReferenceQueueData sets the ReferenceQueueData for this object
func (o *Object) SetReferenceQueueData(data *ReferenceQueueData) {
	o.extra = data
}

// ============================================================
// Special Queue Constants (similar to Java implementation)
// ============================================================

// In Java, there are special queue instances:
// - ReferenceQueue.NULL: indicates no queue was registered
// - ReferenceQueue.ENQUEUED: indicates reference was enqueued (marker)
//
// For MVP, we use nil to represent "no queue" and rely on
// ReferenceState to track enqueued status.
