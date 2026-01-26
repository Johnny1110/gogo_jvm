package heap

// ============================================================
// Reference Processor - v0.3.2
// ============================================================
// This module handles the interaction between GC and Reference objects.
//
// In a real JVM, the Reference Processor is invoked during GC to:
// 1. Discover References whose referents are only softly/weakly reachable
// 2. Clear the referents based on reference type
// 3. Enqueue the References to their registered queues
//
// For MVP (v0.3.2), we provide the basic structure.
// TODO: Full GC integration will be implemented in v0.5.x.
//
// GC Reference Processing Phases:
//
//	┌─────────────────────────────────────────────────────────────┐
//	│  Phase 1: Mark                                              │
//	│  - Mark all strongly reachable objects                      │
//	│  - Discover Reference objects, add to pending lists         │
//	├─────────────────────────────────────────────────────────────┤
//	│  Phase 2: Process SoftReferences                            │
//	│  - If memory is tight: clear referents, add to pending      │
//	│  - If memory is OK: keep referents alive                    │
//	├─────────────────────────────────────────────────────────────┤
//	│  Phase 3: Process WeakReferences                            │
//	│  - Clear all referents (unconditionally)                    │
//	│  - Add to pending list                                      │
//	├─────────────────────────────────────────────────────────────┤
//	│  Phase 4: Finalization                                      │
//	│  - Process objects with finalize() methods                  │
//	├─────────────────────────────────────────────────────────────┤
//	│  Phase 5: Process PhantomReferences                         │
//	│  - Enqueue phantom references (referent already finalized)  │
//	├─────────────────────────────────────────────────────────────┤
//	│  Phase 6: Sweep / Enqueue                                   │
//	│  - Reclaim unreachable objects                              │
//	│  - Reference Handler enqueues pending references            │
//	└─────────────────────────────────────────────────────────────┘

// ============================================================
// Pending Lists
// ============================================================

// PendingList holds References discovered during GC marking phase
// These are separated by type for ordered processing
type PendingList struct {
	// SoftRefs holds discovered SoftReferences
	// Processed first, may or may not be cleared based on memory
	SoftRefs []*Object
	// WeakRefs holds discovered WeakReferences
	// Always cleared during GC
	WeakRefs []*Object
	// PhantomRefs holds discovered PhantomReferences
	// Processed after finalization
	PhantomRefs []*Object
}

// NewPendingList creates an empty PendingList
func NewPendingList() *PendingList {
	return &PendingList{
		SoftRefs:    make([]*Object, 0),
		WeakRefs:    make([]*Object, 0),
		PhantomRefs: make([]*Object, 0),
	}
}

// Clear resets the pending list for the next GC cycle
func (pl *PendingList) Clear() {
	pl.SoftRefs = pl.SoftRefs[:0]
	pl.WeakRefs = pl.WeakRefs[:0]
	pl.PhantomRefs = pl.PhantomRefs[:0]
}

// ============================================================
// Reference Processor
// ============================================================

// ReferenceProcessor manages Reference processing during GC
type ReferenceProcessor struct {
	// pending holds References discovered during current GC cycle
	pending *PendingList
	// memoryPressure indicates current memory pressure (0.0 - 1.0)
	// Used to decide whether to clear SoftReferences
	memoryPressure float64
}

// NewReferenceProcessor creates a new ReferenceProcessor
func NewReferenceProcessor() *ReferenceProcessor {
	return &ReferenceProcessor{
		pending:        NewPendingList(),
		memoryPressure: 0.0,
	}
}

// ============================================================
// Discovery Phase (called during GC Mark)
// ============================================================

// DiscoverReference is called when GC encounters a Reference object
// It categorizes the Reference for later processing
//
// This should be called during the Mark phase when:
// 1. The Reference object itself is reachable
// 2. The referent may or may not be reachable
func (rp *ReferenceProcessor) DiscoverReference(ref *Object) {
	if ref == nil {
		return
	}

	refData := ref.GetReferenceData()
	if refData == nil {
		return
	}

	// Only discover Active references
	if refData.State != RefStateActive {
		return
	}

	// Categorize by type
	switch refData.RefType {
	case RefTypeSoft:
		rp.pending.SoftRefs = append(rp.pending.SoftRefs, ref)
	case RefTypeWeak:
		rp.pending.WeakRefs = append(rp.pending.WeakRefs, ref)
	case RefTypePhantom:
		rp.pending.PhantomRefs = append(rp.pending.PhantomRefs, ref)
	}
}

// ============================================================
// Processing Phases
// ============================================================

// SetMemoryPressure sets the current memory pressure
// This affects SoftReference clearing decisions
func (rp *ReferenceProcessor) SetMemoryPressure(pressure float64) {
	if pressure < 0 {
		pressure = 0
	}
	if pressure > 1 {
		pressure = 1
	}
	rp.memoryPressure = pressure
}

// ProcessSoftReferences processes discovered SoftReferences
// Clears referents based on memory pressure and LRU
//
// HotSpot's policy (simplified):
// - If free_heap < soft_ref_threshold: clear all soft refs
// - Otherwise: clear based on LRU timestamp
//
// Our MVP policy:
// - If memoryPressure > 0.8: clear all soft refs
// - Otherwise: keep all soft refs
func (rp *ReferenceProcessor) ProcessSoftReferences() {
	for _, ref := range rp.pending.SoftRefs {
		refData := ref.GetReferenceData()
		if refData == nil {
			continue
		}

		// Simple policy: clear if memory pressure is high
		if rp.memoryPressure > 0.8 {
			rp.clearReference(ref, refData)
		}
		// TODO: Implement LRU-based clearing for intermediate pressure
	}
}

// ProcessWeakReferences processes discovered WeakReferences
// Unconditionally clears all referents
func (rp *ReferenceProcessor) ProcessWeakReferences() {
	for _, ref := range rp.pending.WeakRefs {
		refData := ref.GetReferenceData()
		if refData == nil {
			continue
		}
		rp.clearReference(ref, refData)
	}
}

// ProcessPhantomReferences processes discovered PhantomReferences
// Marks them as pending for enqueue (referent is NOT cleared in Java 8)
//
// Note: In Java 9+, phantom referents ARE cleared. We follow Java 8 behavior.
func (rp *ReferenceProcessor) ProcessPhantomReferences() {
	for _, ref := range rp.pending.PhantomRefs {
		refData := ref.GetReferenceData()
		if refData == nil {
			continue
		}

		// PhantomReference must have a queue
		if refData.Queue != nil {
			refData.State = RefStatePending
			// Note: We do NOT clear refData.Referent for PhantomReference (Java 8)
			// In Java 9+, it would be cleared here: rp.clearReference(ref, refData)
		}
	}
}

// clearReference clears the referent and marks for enqueue
func (rp *ReferenceProcessor) clearReference(ref *Object, refData *ReferenceData) {
	// Clear the referent <T> -> nil
	refData.Referent = nil
	// Mark for enqueue if queue is registered
	if refData.Queue != nil {
		refData.State = RefStatePending // pending enqueue.
	} else {
		refData.State = RefStateInactive // not require enqueue.
	}
}

// ============================================================
// Enqueue Phase (called after GC)
// ============================================================

// EnqueuePendingReferences enqueues all pending References to their queues
// This should be called after GC completes
//
// In real JVM, this is done by the Reference Handler thread.
// For MVP, we do it synchronously.
func (rp *ReferenceProcessor) EnqueuePendingReferences() {
	// Process all reference types
	totalRefsSize := len(rp.pending.SoftRefs) + len(rp.pending.WeakRefs) + len(rp.pending.PhantomRefs)
	allRefs := make([]*Object, 0, totalRefsSize)
	allRefs = append(allRefs, rp.pending.SoftRefs...)
	allRefs = append(allRefs, rp.pending.WeakRefs...)
	allRefs = append(allRefs, rp.pending.PhantomRefs...)

	for _, ref := range allRefs {
		refData := ref.GetReferenceData()
		if refData == nil {
			continue
		}

		// Only enqueue Pending references with a queue
		if refData.State == RefStatePending && refData.Queue != nil {
			queueData := refData.Queue.GetReferenceQueueData()
			if queueData != nil {
				queueData.Enqueue(ref) // Enqueue() will change ref state to RefStateEnqueued(2)
			}
		}
	}

	// Clear pending list (soft, week, phantom) for next GC cycle
	rp.pending.Clear()
}

// ============================================================
// Full GC Cycle (convenience method)
// ============================================================

// ProcessAllReferences runs the complete reference processing cycle
// Call this after the Mark phase and before the Sweep phase
func (rp *ReferenceProcessor) ProcessAllReferences() {
	// Order matters: Soft -> Weak -> Phantom -> enqueue
	rp.ProcessSoftReferences()
	rp.ProcessWeakReferences()
	rp.ProcessPhantomReferences()
	rp.EnqueuePendingReferences()
}

// ============================================================
// Global Reference Processor Instance
// ============================================================
// globalRefProcessor is the singleton ReferenceProcessor
// Used by GC when processing references
var globalRefProcessor *ReferenceProcessor

// GetReferenceProcessor returns the global ReferenceProcessor instance
func GetReferenceProcessor() *ReferenceProcessor {
	if globalRefProcessor == nil {
		globalRefProcessor = NewReferenceProcessor()
	}
	return globalRefProcessor
}
