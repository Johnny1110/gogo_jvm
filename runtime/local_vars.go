package runtime

import (
	"github.com/Johnny1110/gogo_jvm/runtime/java"
)

// NewLocalVars LocalVars must be pre-allocate, size is defined in class file.
func NewLocalVars(maxLocals uint16) java.Slots {
	if maxLocals > 0 {
		return make(java.Slots, maxLocals)
	}
	return nil
}
