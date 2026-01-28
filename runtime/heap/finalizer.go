package heap

import "sync"

// ============================================================
// v0.3.3: Finalization Support
// ============================================================
// This file implements the finalization mechanism for objects
// that override the finalize() method.
//
// Design based on HotSpot JVM:
// 1. Objects with non-trivial finalize() are registered at creation time
// 2. When GC determines object is unreachable, it's moved to finalization queue
// 3. Finalizer thread processes the queue and calls finalize()
// 4. After finalization, object can be collected in next GC cycle
//
// Note: This is a simplified implementation for educational purposes.
// Real JVMs have much more sophisticated finalization handling.

// ============================================================
// Finalizable Object Wrapper
// ============================================================

// FinalizableState represents the finalization state of an object
type FinalizableState int

const (
	// FinalizableActive - Object is alive, finalize() not yet called
	FinalizableActive FinalizableState = iota
	// FinalizablePending - Object is unreachable, waiting for finalization
	FinalizablePending
	// FinalizableFinalized - finalize() has been called
	FinalizableFinalized
)

// FinalizableObject wraps an object that needs finalization
type FinalizableObject struct {
	// The object that needs finalization
	Object *Object
	// Current finalization state
	State FinalizableState
	// Has finalize() been called? (can only be called once)
	FinalizeCalled bool
}

// ============================================================
// Finalization Queue
// ============================================================
// FinalizationQueue manages objects waiting for finalization
// Thread-safe for concurrent access

type FinalizationQueue struct {
	// Registered objects that may need finalization
	// Key: Object pointer, Value: FinalizableObject wrapper
	registered map[*Object]*FinalizableObject
	// Queue of objects pending finalization
	pending []*FinalizableObject
	// Mutex for thread safety
	mu sync.RWMutex
	// Statistics
	stats FinalizationStats
}

// FinalizationStats tracks finalization statistics
type FinalizationStats struct {
	// Total objects registered for finalization
	TotalRegistered int64
	// Objects currently pending finalization
	CurrentPending int64
	// Total finalize() calls completed
	TotalFinalized int64
}

// Global finalization queue instance
var GlobalFinalizationQueue = NewFinalizationQueue()

// NewFinalizationQueue creates a new finalization queue
func NewFinalizationQueue() *FinalizationQueue {
	return &FinalizationQueue{
		registered: make(map[*Object]*FinalizableObject),
		pending:    make([]*FinalizableObject, 0),
		stats: FinalizationStats{
			TotalRegistered: 0,
			CurrentPending:  0,
			TotalFinalized:  0,
		},
	}
}

// ============================================================
// Registration Methods
// ============================================================

// Register adds an object to the finalization system
// This should be called when creating an object whose class has
// a non-trivial finalize() method
func (fq *FinalizationQueue) Register(obj *Object) {
	if obj == nil {
		return
	}

	fq.mu.Lock()
	defer fq.mu.Unlock()

	// Don't register twice
	if _, exists := fq.registered[obj]; exists {
		return
	}

	wrapper := &FinalizableObject{
		Object:         obj,
		State:          FinalizableActive,
		FinalizeCalled: false,
	}

	fq.registered[obj] = wrapper
	fq.stats.TotalRegistered++
}

// Unregister removes an object from the finalization system
// This is called when an object is collected after finalization
func (fq *FinalizationQueue) Unregister(obj *Object) {
	if obj == nil {
		return
	}

	fq.mu.Lock()
	defer fq.mu.Unlock()

	delete(fq.registered, obj)
}

// ============================================================
// GC Integration Methods
// ============================================================

// MarkPending moves an unreachable object to the pending queue
// This is called by GC when it finds an unreachable finalizable object
func (fq *FinalizationQueue) MarkPending(obj *Object) bool {
	if obj == nil {
		return false
	}

	fq.mu.Lock()
	defer fq.mu.Unlock()

	wrapper, exists := fq.registered[obj]
	if !exists {
		return false
	}

	// Can only mark as pending if currently active
	if wrapper.State != FinalizableActive {
		return false
	}

	wrapper.State = FinalizablePending
	fq.pending = append(fq.pending, wrapper)
	fq.stats.CurrentPending++

	return true
}

// GetPending returns the next object pending finalization
// Returns nil if queue is empty
func (fq *FinalizationQueue) GetPending() *Object {
	fq.mu.Lock()
	defer fq.mu.Unlock()

	if len(fq.pending) == 0 {
		return nil
	}

	// Dequeue first pending object
	wrapper := fq.pending[0]
	fq.pending = fq.pending[1:]
	fq.stats.CurrentPending--

	return wrapper.Object
}

// MarkFinalized marks an object as having had finalize() called
func (fq *FinalizationQueue) MarkFinalized(obj *Object) {
	if obj == nil {
		return
	}

	fq.mu.Lock()
	defer fq.mu.Unlock()

	wrapper, exists := fq.registered[obj]
	if !exists {
		return
	}

	wrapper.State = FinalizableFinalized
	wrapper.FinalizeCalled = true
	fq.stats.TotalFinalized++
}

// ============================================================
// Query Methods
// ============================================================

// IsRegistered checks if an object is registered for finalization
func (fq *FinalizationQueue) IsRegistered(obj *Object) bool {
	if obj == nil {
		return false
	}

	fq.mu.RLock()
	defer fq.mu.RUnlock()

	_, exists := fq.registered[obj]
	return exists
}

// IsFinalized checks if an object has already been finalized
func (fq *FinalizationQueue) IsFinalized(obj *Object) bool {
	if obj == nil {
		return false
	}

	fq.mu.RLock()
	defer fq.mu.RUnlock()

	wrapper, exists := fq.registered[obj]
	return exists && wrapper.FinalizeCalled
}

// PendingCount returns the number of objects pending finalization
func (fq *FinalizationQueue) PendingCount() int {
	fq.mu.RLock()
	defer fq.mu.RUnlock()

	return len(fq.pending)
}

// GetStats returns a copy of the finalization statistics
func (fq *FinalizationQueue) GetStats() FinalizationStats {
	fq.mu.RLock()
	defer fq.mu.RUnlock()

	return fq.stats
}

// ============================================================
// Helper Functions
// ============================================================

// ShouldRegisterForFinalization determines if an object needs
// to be registered for finalization based on its class
//
// This is a helper function that should be called during object creation
// The actual check requires access to method_area.Class, so the caller
// must perform the HasNonTrivialFinalizer() check
func ShouldRegisterForFinalization(hasNonTrivialFinalizer bool) bool {
	return hasNonTrivialFinalizer
}
