package runtime

import (
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
)

// Frame java:
// ┌────────────────────────────────────────┐
// │              Frame                     │
// ├────────────────────────────────────────┤
// │  lower (pointing to invoker frame)     │
// ├────────────────────────────────────────┤
// │  localVars	table	                    │
// │  ┌─────┬─────┬─────┬─────┬─────┐       │
// │  │  0  │  1  │  2  │  3  │ ... │       │
// │  └─────┴─────┴─────┴─────┴─────┘       │
// ├────────────────────────────────────────┤
// │  operandStack                          │
// │  ┌─────┐                               │
// │  │ top │ ← currentFrame (method)       │
// │  ├─────┤                               │
// │  │ ..  │                               │
// │  └─────┘                               │
// ├────────────────────────────────────────┤
// │  thread (thread.currentThread())       │
// ├────────────────────────────────────────┤
// │  nextPC (next method bytecode addr)    │
// └────────────────────────────────────────┘
type Frame struct {
	lower        *Frame // previous frame (caller frame)
	localVars    rtcore.Slots
	operandStack *OperandStack
	thread       *Thread
	nextPC       int
	method       *method_area.Method
}

// NewFrame create new Frame
// thread: target thread which is frame belongs to
func NewFrame(thread *Thread, maxLocals, maxStack uint16) *Frame {
	return &Frame{
		thread:       thread,
		localVars:    NewLocalVars(maxLocals),
		operandStack: NewOperandStack(maxStack),
	}
}

func NewFrameWithMethod(thread *Thread, method *method_area.Method) *Frame {
	return &Frame{
		thread:       thread,
		method:       method,
		localVars:    NewLocalVars(method.MaxLocals()),
		operandStack: NewOperandStack(method.MaxStack()),
	}
}

func (f *Frame) LocalVars() rtcore.Slots {
	return f.localVars
}

func (f *Frame) Thread() *Thread {
	return f.thread
}

func (f *Frame) OperandStack() *OperandStack {
	return f.operandStack
}

func (f *Frame) NextPC() int {
	return f.nextPC
}

// SetNextPC setup next instruction address (index)
// use for `for` `while` `if` `break`...
func (f *Frame) SetNextPC(nextPC int) {
	f.nextPC = nextPC
}

// Lower control caller's frame
func (f *Frame) Lower() *Frame {
	return f.lower
}

// RevertNextPC revert PC to current instruction
// Usages: new/getstatic/putstatic/invokestatic
func (f *Frame) RevertNextPC() {
	f.nextPC = f.thread.PC()
}

func (f *Frame) Method() *method_area.Method { return f.method }
