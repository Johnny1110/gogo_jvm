package heap

import "sync/atomic"

// ============================================================
// Mark Word - Object Header (64 bits)
// ============================================================

// bit const
const (
	// Lock State - Position and Mask
	LockStateBits = 2
	LockStateMask = 0x03 // bits [1:0] -> 0000 0011

	// Biased - Position and Mask
	BiasedBitPos  = 2
	BiasedBitMask = 0x04 // bit [2] -> 0000 0100

	// Age - Position and Mask
	AgeShift = 3
	AgeBits  = 4
	AgeMask  = 0x78 // bits [6:3] -> 0111 1000

	// HashCode - Position and Mask
	// identity hashCode 從 Mark Word 的第 7 個 bit 開始存放
	// bit index: [63 .............. 37  36 .............. 7  6 .... 0]
	//      							 ↑                 ↑
	//                            hashCode 高位        hashCode 低位
	HashCodeShift = 7
	HashCodeBits  = 31                                  // hash code take 31 bits
	HashCodeMask  = uint64(0x7FFFFFFF) << HashCodeShift // bits [37:7]
	// 0x7FFFFFFF = 0111 1111 1111 1111 1111 1111 1111 1111  (31 個 1)
	// 左移:  << 7
	// HashCodeMask 為:
	// 0000...000 0111 1111 ... 1111 0000000
	//   ↑         ↑                ↑      ↑
	//   其他    bit 37            bit 7   Lock, Biased, Age

	// Max GC Age
	MaxAge = 15
)

// ============================================================
// Lock State const
// ============================================================

const (
	// LockStateUnlocked Non-Lock (01)
	LockStateUnlocked = 0x01

	// LockStateLightLock (00)
	LockStateLightLock = 0x00

	// LockStateHeavyLock (10)
	LockStateHeavyLock = 0x02

	// LockStateGCMarked GC Marked (11)
	LockStateGCMarked = 0x03
)

// ============================================================
// Mark Word init value
// ============================================================

// InitialMarkWord (init when new object) Mark Word
// lock = 01 (non-lock), biased = 0, age = 0, hashCode = 0
const InitialMarkWord = uint64(LockStateUnlocked)

// ============================================================
// Hash Code Operations
// ============================================================

// HashCode get object's identity hash code
// create if not exist, with CAS access
func (o *Object) HashCode() int32 {
	for {
		mark := atomic.LoadUint64(&o.markWord) // safe read.
		hash := int32((mark & HashCodeMask) >> HashCodeShift)

		if hash != 0 {
			return hash
		}

		newHash := generateHashCode()
		newMark := (mark &^ HashCodeMask) | (uint64(newHash) << HashCodeShift)
		// (mark &^ HashCodeMask)
		// &^ = AND NOT（清除）
		// 把 Mark Word 中 舊的 hashCode bits 清成 0
		// 其他 bit（lock, age, biased）完全保留

		// (uint64(newHash) << HashCodeShift)
		// 把新的 hash 左移到 bits [37:7]

		// (清掉舊 hash) | (新 hash) - > OR 合併

		// CAS update, make sure thread safe.
		if atomic.CompareAndSwapUint64(&o.markWord, mark, newMark) {
			return newHash
		} else {
			//  CAS failed, means another thread already set, reload again.
			continue
		}
	}
}

// SetHashCode setup hash code (internal usage)
func (o *Object) SetHashCode(hash int32) {
	for {
		mark := atomic.LoadUint64(&o.markWord)
		newMark := (mark &^ HashCodeMask) | (uint64(hash&0x7FFFFFFF) << HashCodeShift)
		// (mark erase old hashcode) OR (new hash)
		if atomic.CompareAndSwapUint64(&o.markWord, mark, newMark) {
			return
		}
	}
}

// HasHashCode check hash code
func (o *Object) HasHashCode() bool {
	mark := atomic.LoadUint64(&o.markWord)
	return (mark & HashCodeMask) != 0
}

// ============================================================
// GC Age Operations
// ============================================================

// GCAge get GC age (0-15)
func (o *Object) GCAge() uint8 {
	mark := atomic.LoadUint64(&o.markWord)
	return uint8((mark & AgeMask) >> AgeShift)
}

// SetGCAge set GC age
func (o *Object) SetGCAge(age uint8) {
	if age > MaxAge {
		age = MaxAge
	}

	for {
		mark := atomic.LoadUint64(&o.markWord)
		newMark := (mark &^ AgeMask) | (uint64(age) << AgeShift)
		if atomic.CompareAndSwapUint64(&o.markWord, mark, newMark) {
			return
		}
	}
}

// IncrementAge increase GC Age（Minor GC 後存活時呼叫）
// Max Age = 15, can not gt 15.
func (o *Object) IncrementAge() {
	for {
		mark := atomic.LoadUint64(&o.markWord)
		age := uint8((mark & AgeMask) >> AgeShift)
		if age >= MaxAge {
			return // do nothing, object is too old.
		}

		newMark := (mark &^ AgeMask) | (uint64(age+1) << AgeShift)
		if atomic.CompareAndSwapUint64(&o.markWord, mark, newMark) {
			return
		}
	}
}

// ============================================================
// Lock State Operations (TODO v0.4.x Multi-Thread)
// ============================================================
// LockState get lock state
// - 01: Normal or Biased Lock
// - 00: Lightweight Lock
// - 10: Heavyweight Lock
// - 11: GC Marked
func (o *Object) LockState() uint8 {
	mark := atomic.LoadUint64(&o.markWord)
	return uint8(mark & LockStateMask)
}

// SetLockState set up lock state
func (o *Object) SetLockState(state uint8) {
	for {
		mark := atomic.LoadUint64(&o.markWord)
		newMark := (mark &^ LockStateMask) | uint64(state&LockStateMask) // uint64(xxxx xxxx & 0000 0011)
		if atomic.CompareAndSwapUint64(&o.markWord, mark, newMark) {
			return
		}
	}
}

// IsLocked check object is lock
func (o *Object) IsLocked() bool {
	state := o.LockState()
	return state == LockStateLightLock || state == LockStateHeavyLock
}

// IsBiased check object is biased (non-lock state &&  biased == 1)
func (o *Object) IsBiased() bool {
	mark := atomic.LoadUint64(&o.markWord)
	lockState := mark & LockStateMask
	biased := (mark & BiasedBitMask) >> BiasedBitPos
	return lockState == LockStateUnlocked && biased == 1
}

// ============================================================
// GC Mark Operations (TODO v0.5.x GC)
// ============================================================

// MarkForGC mark object into GC State
// return origin mark word (restore it after GC)
func (o *Object) MarkForGC() uint64 {
	for {
		mark := atomic.LoadUint64(&o.markWord)
		newMark := (mark &^ LockStateMask) | uint64(LockStateGCMarked)
		if atomic.CompareAndSwapUint64(&o.markWord, mark, newMark) {
			return mark // return origin (restore)
		}
	}
}

// RestoreMarkWord after GC, restore mark word
func (o *Object) RestoreMarkWord(originalMark uint64) {
	atomic.StoreUint64(&o.markWord, originalMark)
}

// IsGCMarked check obj is GC state
func (o *Object) IsGCMarked() bool {
	return o.LockState() == LockStateGCMarked
}

// ============================================================
// Debug
// ============================================================
// MarkWord get origin mark word（for testing）
func (o *Object) MarkWord() uint64 {
	return atomic.LoadUint64(&o.markWord)
}
