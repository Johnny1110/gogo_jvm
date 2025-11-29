package runtime

// NewLocalVars LocalVars must be pre-allocate, size is defined in class file.
func NewLocalVars(maxLocals uint16) Slots {
	if maxLocals > 0 {
		return make(Slots, maxLocals)
	}
	return nil
}
