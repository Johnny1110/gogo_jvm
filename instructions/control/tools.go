package control

import "github.com/Johnny1110/gogo_jvm/runtime"

// branch helper func: perform jump
func branch(frame *runtime.Frame, offset int) {
	pc := frame.Thread().PC()
	nextPC := pc + offset
	frame.SetNextPC(nextPC)
}
