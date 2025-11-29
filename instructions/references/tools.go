package references

import (
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// InvokeMethod call method common func
// usage: invokestatic, invokevirtual
func InvokeMethod(invokerFrame *runtime.Frame, method *heap.Method) {
	// 1, get thread
	thread := invokerFrame.Thread()

	// 2. create a new frame
	maxLocal := method.MaxLocals()
	maxStack := method.MaxStack()
	newFrame := thread.NewFrame(maxLocal, maxStack)
	thread.PushFrame(newFrame)

	// 3. pass vars
	argSlotCount := int(method.ArgSlotCount())
	if argSlotCount > 0 {
		for i := argSlotCount - 1; i >= 0; i-- {
			// the passing var will be standby in invokerFrame's stack.
			// pop it from invokerFrame's stack
			slot := invokerFrame.OperandStack().PopSlot()
			// put into newFrame's LocalVars
			newFrame.LocalVars().SetSlot(uint(i), slot)
		}
	}

	// no need to reset PC, new frame nextPC will be default 0
}
