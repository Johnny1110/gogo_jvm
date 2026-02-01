package runtime

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
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
	currentPC    int
	method       *method_area.Method

	exHandler func(frame *Frame, ex *heap.Object)
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

// please using NewFrameWithMethodAndExHandler if possible.
func NewFrameWithMethod(thread *Thread, method *method_area.Method) *Frame {
	return &Frame{
		thread:       thread,
		method:       method,
		localVars:    NewLocalVars(method.MaxLocals()),
		operandStack: NewOperandStack(method.MaxStack()),
	}
}

func NewFrameWithMethodAndExHandler(thread *Thread,
	method *method_area.Method,
	exHandler func(frame *Frame, ex *heap.Object)) *Frame {
	return &Frame{
		thread:       thread,
		method:       method,
		localVars:    NewLocalVars(method.MaxLocals()),
		operandStack: NewOperandStack(method.MaxStack()),
		exHandler:    exHandler,
	}
}

// NewNativeFrame create frame for native method, not require to push into JVMStack
func NewNativeFrame(thread *Thread, maxLocals uint16) *Frame {
	return &Frame{
		thread:       thread,
		localVars:    NewLocalVars(maxLocals),
		operandStack: nil, // native method don't need op-stack, just run native Go func()
	}
}

// NewNativeFrameWithStack create frame for native method with return value support
// v0.3.0: support return val native call
func NewNativeFrameWithStack(thread *Thread, maxLocals uint16, returnType string) *Frame {
	var opStack *OperandStack

	if returnType == "V" || returnType == "" {
		// void: not require opStack
		opStack = nil
	} else {
		stackSize := uint16(1)
		if len(returnType) > 0 && (returnType[0] == 'J' || returnType[0] == 'D') {
			stackSize = 2 // long/double need 2 slots
		}

		opStack = NewOperandStack(stackSize)
	}

	return &Frame{
		thread:       thread,
		localVars:    NewLocalVars(maxLocals),
		operandStack: opStack,
	}
}

func NewNativeFrameWithStackAndExHandler(callerFrame *Frame,
	maxLocals uint16,
	returnType string,
	exHandler func(frame *Frame, ex *heap.Object)) *Frame {
	var opStack *OperandStack

	if returnType == "V" || returnType == "" {
		// void: not require opStack
		opStack = nil
	} else {
		stackSize := uint16(1)
		if len(returnType) > 0 && (returnType[0] == 'J' || returnType[0] == 'D') {
			stackSize = 2 // long/double need 2 slots
		}

		opStack = NewOperandStack(stackSize)
	}

	return &Frame{
		thread:       callerFrame.Thread(),
		localVars:    NewLocalVars(maxLocals),
		operandStack: opStack,
		method:       callerFrame.method,
		exHandler:    exHandler,
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

func (f *Frame) CurrentPC() int {
	return f.currentPC
}

// SetNextPC setup next instruction address (index)
// use for `for` `while` `if` `break`...
func (f *Frame) SetNextPC(nextPC int) {
	f.currentPC = f.nextPC
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

func (f *Frame) JavaThrow(ex *heap.Object) {
	if ex == nil {
		panic("JavaThrow called with nil exception.")
	}

	exData, ok := ex.Extra().(*heap.ExceptionData)
	if !ok {
		panic("JavaThrow called with a non exception object.")
	}

	if f.exHandler == nil {
		if global.DebugMode() {
			fmt.Printf("@@ DEBUG - current frame didn't have exHandler, panic directly.\n")
		}
		panic(exData.Message)
	}

	// handle ex.
	f.exHandler(f, ex)
}

func (f *Frame) LoaderClass(className string) *method_area.Class {
	return f.method.Class().Loader().LoadClass(className, false)
}
