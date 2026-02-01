package rtcore

import (
	"fmt"
	"math"
)

// Slot is JVM runtime basic data unit
// JVM standard: local vars table and op stack all made by Slot.
// each Slot can store a 32-bit data
// 64-bit data (long. double), should using 2 Slot
type Slot struct {
	// Num store Number type data (int32, float32)
	// for long or double, we split to 2 slots.
	Num int32

	// ref store reference type (pointer)
	// GC need to know the object references, to iterate
	Ref interface{}
}

// JVM defined long/double, should take rtcore[n] and rtcore[n+1]
// and also should read rtcore[n] and rtcore[n+1]

type Slots []Slot

// NewSlots create Slots with maxLen
func NewSlots(slotCount uint) Slots {
	if slotCount > 0 {
		return make(Slots, slotCount)
	}
	return nil
}

// =============== Int ===============

func (s Slots) SetInt(index uint, val int32) {
	s[index].Num = val
}

func (s Slots) GetInt(index uint) int32 {
	return s[index].Num
}

// =============== Boolean ===============

func (s Slots) SetBoolean(index uint, boolean bool) {
	if boolean {
		s[index].Num = 1
	} else {
		s[index].Num = 0
	}
}

func (s Slots) GetBoolean(index uint) bool {
	return s[index].Num == 1
}

// =============== Float ===============

func (s Slots) SetFloat(index uint, val float32) {
	bits := math.Float32bits(val)
	s[index].Num = int32(bits)
}

func (s Slots) GetFloat(index uint) float32 {
	bits := uint32(s[index].Num)
	return math.Float32frombits(bits)
}

// =============== Long ===============

func (s Slots) SetLong(index uint, val int64) {
	// 低 32 位
	s[index].Num = int32(val)
	// 高 32 位
	s[index+1].Num = int32(val >> 32)
}

func (s Slots) GetLong(index uint) int64 {
	low := uint32(s[index].Num)
	high := uint32(s[index+1].Num)
	return int64(high)<<32 | int64(low)
}

// =============== Double ===============

func (s Slots) SetDouble(index uint, val float64) {
	bits := math.Float64bits(val)
	s.SetLong(index, int64(bits))
}

func (s Slots) GetDouble(index uint) float64 {
	bits := uint64(s.GetLong(index))
	return math.Float64frombits(bits)
}

// =============== Ref ===============

func (s Slots) SetRef(index uint, ref interface{}) {
	s[index].Ref = ref
}

func (s Slots) GetRef(index uint) interface{} {
	return s[index].Ref
}

// =========== Slot Operation ============

func (s Slots) SetSlot(index uint, slot Slot) {
	s[index] = slot
}

func (s Slots) GetSlot(index uint) Slot {
	return s[index]
}

// GetThis get this ref
// for Object methods, rtcore[0] always be this
func (s Slots) GetThis() interface{} {
	return s.GetRef(0)
}

func (s Slot) String() string {
	if s.Ref != nil {
		return fmt.Sprintf("Ref<%T>", s.Ref)
	}
	return fmt.Sprintf("Num<%d>", s.Num)
}
