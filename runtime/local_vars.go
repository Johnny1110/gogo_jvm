package runtime

import "math"

// LocalVars 局部變量表
type LocalVars []Slot

// NewLocalVars LocalVars must be pre-allocate, size is defined in class file.
func NewLocalVars(maxLocals uint16) LocalVars {
	if maxLocals > 0 {
		return make(LocalVars, maxLocals)
	}
	return nil
}

// =========== int Operation ============

func (lv LocalVars) SetInt(index uint, val int32) {
	lv[index].num = val
}

func (lv LocalVars) GetInt(index uint) int32 {
	return lv[index].num
}

// =========== float Operation ============
// float 32 and int32 both are 32-bit

func (lv LocalVars) SetFloat(index uint, val float32) {
	bits := math.Float32bits(val)
	// make float32 to int32
	lv[index].num = int32(bits)
}

func (lv LocalVars) GetFloat(index uint) float32 {
	bits := uint32(lv[index].num)
	return math.Float32frombits(bits)
}

// =========== long Operation ============
// long take 2 position in []slot

func (lv LocalVars) SetLong(index uint, val int64) {
	// low 32 bit put into first slot
	lv[index].num = int32(val)
	// high 32 bit put into  sec slot index+1
	lv[index+1].num = int32(val >> 32)
}

func (lv LocalVars) GetLong(index uint) int64 {
	low := uint32(lv[index].num)
	high := uint32(lv[index+1].num)
	return int64(high)<<32 | int64(low)
}

// =========== double Operation ============
// double take 2 position in []slot

func (lv LocalVars) SetDouble(index uint, val float64) {
	bits := math.Float64bits(val)
	// make double to long and store as long
	lv.SetLong(index, int64(bits))
}

func (lv LocalVars) GetDouble(index uint) float64 {
	bits := uint64(lv.GetLong(index))
	return math.Float64frombits(bits)
}

// =========== Reference Operation ============

func (lv LocalVars) SetRef(index uint, ref interface{}) {
	lv[index].ref = ref
}

func (lv LocalVars) GetRef(index uint) interface{} {
	return lv[index].ref
}

// =========== Slot Operation ============

func (lv LocalVars) SetSlot(index uint, slot Slot) {
	lv[index] = slot
}

func (lv LocalVars) GetSlot(index uint) Slot {
	return lv[index]
}

// GetThis get this ref
// for Object methods, slot[0] always be this
func (lv LocalVars) GetThis() interface{} {
	return lv.GetRef(0)
}
