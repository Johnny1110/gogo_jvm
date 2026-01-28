package runtime

import (
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

const DEFAULT_STACK_SIZE = 1024

type Thread struct {
	pc    int       // Program Counter
	stack *JVMStack // JVM Frame Stack
}

// NewThread create new Thread
func NewThread() *Thread {
	return &Thread{
		pc:    0,
		stack: NewJVMStack(DEFAULT_STACK_SIZE),
	}
}

func (t *Thread) PC() int {
	return t.pc
}

func (t *Thread) SetPC(pc int) {
	t.pc = pc
}

func (t *Thread) PushFrame(frame *Frame) {
	t.stack.Push(frame)
}

func (t *Thread) PopFrame() *Frame {
	return t.stack.Pop()
}

// CurrentFrame get current frame without pop
func (t *Thread) CurrentFrame() *Frame {
	return t.stack.Top()
}

// TopFrame is CurrentFrame alias
func (t *Thread) TopFrame() *Frame {
	return t.CurrentFrame()
}

func (t *Thread) IsStackEmpty() bool {
	return t.stack.IsEmpty()
}

func (t *Thread) StackDepth() uint {
	return t.stack.Size()
}

func (t *Thread) GetFrames() []*Frame {
	return t.stack.GetFrames()
}

func (t *Thread) ClearStack() {
	t.stack.Clear()
}

// NewFrame create a new Frame and put this thread as constructor's params
func (t *Thread) NewFrame(maxLocals, maxStack uint16) *Frame {
	return NewFrame(t, maxLocals, maxStack)
}

func (t *Thread) NewFrameWithMethodAndExHandler(method *method_area.Method, exHandler func(frame *Frame, ex *heap.Object)) *Frame {
	return NewFrameWithMethodAndExHandler(t, method, exHandler)
}
