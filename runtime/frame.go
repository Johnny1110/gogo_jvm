package runtime

import "github.com/Johnny1110/gogo_jvm/runtime/java"

// Frame java:
// ┌────────────────────────────────────────┐
// │              Frame                     │
// ├────────────────────────────────────────┤
// │  lower (指向下一幀，形成鏈表)             │
// ├────────────────────────────────────────┤
// │  localVars (局部變量表)                  │
// │  ┌─────┬─────┬─────┬─────┬─────┐       │
// │  │  0  │  1  │  2  │  3  │ ... │       │
// │  └─────┴─────┴─────┴─────┴─────┘       │
// ├────────────────────────────────────────┤
// │  operandStack (操作數棧)                │
// │  ┌─────┐                               │
// │  │ top │ ← 棧頂                         │
// │  ├─────┤                               │
// │  │     │                               │
// │  └─────┘                               │
// ├────────────────────────────────────────┤
// │  thread (所屬線程的引用)                 │
// ├────────────────────────────────────────┤
// │  nextPC (下一條要執行的指令地址)          │
// └────────────────────────────────────────┘
type Frame struct {
	lower        *Frame // previous frame (caller frame)
	localVars    java.Slots
	operandStack *OperandStack
	thread       *Thread
	nextPC       int
	method       *java.Method
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

func NewFrameWithMethod(thread *Thread, method *java.Method) *Frame {
	return &Frame{
		thread:       thread,
		method:       method,
		localVars:    NewLocalVars(method.MaxLocals()),
		operandStack: NewOperandStack(method.MaxStack()),
	}
}

func (f *Frame) LocalVars() java.Slots {
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
func (f *Frame) RevertNextPC() {
	f.nextPC = f.thread.PC()
}

func (f *Frame) Method() *java.Method { return f.method }
