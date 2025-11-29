package java

import "math"

// Slot is JVM runtime basic data unit
// JVM standard: local vars table and op stack all made by Slot.
// each Slot can store a 32-bit data
// 64-bit data (long. double), should using 2 Slot
type Slot struct {
	// num store number type data (int32, float32)
	// for long or double, we split to 2 slots.
	num int32

	// ref store reference type (pointer)
	// GC need to know the object references, to iterate
	ref *Object
}

// JVM defined long/double, should take slot[n] and slot[n+1]
// and also should read slot[n] and slot[n+1]

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
	s[index].num = val
}

func (s Slots) GetInt(index uint) int32 {
	return s[index].num
}

// =============== Float ===============

func (s Slots) SetFloat(index uint, val float32) {
	bits := math.Float32bits(val)
	s[index].num = int32(bits)
}

func (s Slots) GetFloat(index uint) float32 {
	bits := uint32(s[index].num)
	return math.Float32frombits(bits)
}

// =============== Long ===============

func (s Slots) SetLong(index uint, val int64) {
	// 低 32 位
	s[index].num = int32(val)
	// 高 32 位
	s[index+1].num = int32(val >> 32)
}

func (s Slots) GetLong(index uint) int64 {
	low := uint32(s[index].num)
	high := uint32(s[index+1].num)
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

func (s Slots) SetRef(index uint, ref *Object) {
	s[index].ref = ref
}

func (s Slots) GetRef(index uint) *Object {
	return s[index].ref
}

// =========== Slot Operation ============

func (s Slots) SetSlot(index uint, slot Slot) {
	s[index] = slot
}

func (s Slots) GetSlot(index uint) Slot {
	return s[index]
}

// GetThis get this ref
// for Object methods, slot[0] always be this
func (s Slots) GetThis() interface{} {
	return s.GetRef(0)
}
