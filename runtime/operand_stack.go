package runtime

import (
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
	"math"
)

type OperandStack struct {
	writePtr uint // equals to current stack element count
	slots    []rtcore.Slot
}

// NewOperandStack crate stack with max size
func NewOperandStack(maxStack uint16) *OperandStack {
	if maxStack > 0 {
		return &OperandStack{
			slots:    make([]rtcore.Slot, maxStack),
			writePtr: 0, // init writer pointer
		}
	}
	return nil
}

func (os *OperandStack) Size() (current, max uint) {
	return os.writePtr, uint(len(os.slots))
}

// ============ Basic Operations ============

func (os *OperandStack) PushInt(val int32) {
	os.slots[os.writePtr].Num = val
	os.writePtr++
}

func (os *OperandStack) PopInt() int32 {
	os.writePtr--
	return os.slots[os.writePtr].Num
}

func (os *OperandStack) PushFloat(val float32) {
	bits := math.Float32bits(val)
	os.slots[os.writePtr].Num = int32(bits)
	os.writePtr++
}

func (os *OperandStack) PopFloat() float32 {
	os.writePtr--
	bits := uint32(os.slots[os.writePtr].Num)
	return math.Float32frombits(bits)
}

func (os *OperandStack) PushLong(val int64) {
	// low 32 bit
	os.slots[os.writePtr].Num = int32(val)
	os.writePtr++
	// high 32 bit
	os.slots[os.writePtr].Num = int32(val >> 32)
	os.writePtr++
}

func (os *OperandStack) PopLong() int64 {
	os.writePtr--
	highBits := uint32(os.slots[os.writePtr].Num)
	os.writePtr--
	lowBits := uint32(os.slots[os.writePtr].Num)
	return int64(highBits)<<32 | int64(lowBits)
}

func (os *OperandStack) PushDouble(val float64) {
	bits := math.Float64bits(val)
	os.PushLong(int64(bits))
}

func (os *OperandStack) PopDouble() float64 {
	bits := uint64(os.PopLong())
	return math.Float64frombits(bits)
}

func (os *OperandStack) PushRef(Ref interface{}) {
	os.slots[os.writePtr].Ref = Ref
	os.writePtr++
}

func (os *OperandStack) PopRef() interface{} {
	os.writePtr--
	Ref := os.slots[os.writePtr].Ref
	os.slots[os.writePtr].Ref = nil // GC
	return Ref
}

func (os *OperandStack) PushSlot(slot rtcore.Slot) {
	os.slots[os.writePtr] = slot
	os.writePtr++
}

func (os *OperandStack) PopSlot() rtcore.Slot {
	os.writePtr--
	return os.slots[os.writePtr]
}

// ================= support methods ===================

// PeekRefFromTop control Ref from top of stack (no pop, just peek)
// top -> n = 0
func (os *OperandStack) PeekRefFromTop(n uint) interface{} {
	return os.slots[os.writePtr-1-n].Ref
}

func (os *OperandStack) PushBoolean(val bool) {
	if val {
		os.PushInt(1)
	} else {
		os.PushInt(0)
	}
}

func (os *OperandStack) PopBoolean() bool {
	return os.PopInt() != 0
}

func (os *OperandStack) Clear() {
	os.writePtr = 0 // reset write pointer
	// clean rtcore.Ref for GC
	for i := range os.slots {
		os.slots[i].Ref = nil
	}
}
