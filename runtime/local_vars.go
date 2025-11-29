package runtime

import "github.com/Johnny1110/gogo_jvm/rtda/heap"

// NewLocalVars LocalVars must be pre-allocate, size is defined in class file.
func NewLocalVars(maxLocals uint16) heap.Slots {
	if maxLocals > 0 {
		return make(heap.Slots, maxLocals)
	}
	return nil
}
