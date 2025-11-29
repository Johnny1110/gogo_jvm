package runtime

import "github.com/Johnny1110/gogo_jvm/runtime/rtcore"

// NewLocalVars LocalVars must be pre-allocate, size is defined in class file.
func NewLocalVars(maxLocals uint16) rtcore.Slots {
	if maxLocals > 0 {
		return make(rtcore.Slots, maxLocals)
	}
	return nil
}
