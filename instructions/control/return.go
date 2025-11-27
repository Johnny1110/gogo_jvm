package control

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// RETURN control void
// opcode = 0xB1
// usage: void methods
type RETURN struct{ base.NoOperandsInstruction }

func (r *RETURN) Execute(frame *runtime.Frame) {
	currentThread := frame.Thread()
	// void method don't have control value, just pop frame
	currentThread.PopFrame()
}

func (r *RETURN) Opcode() uint8 {
	return 0xB1
}

// IRETURN control int
// opcode = 0xAC
// also use for: boolean, byte, char, short（they are all int in JVM）
type IRETURN struct{ base.NoOperandsInstruction }

// Execute IRETURN
// move int from current frame's stack to caller frame's stack
func (i *IRETURN) Execute(frame *runtime.Frame) {
	currentThread := frame.Thread()
	currentFrame := currentThread.PopFrame()       // pop current frame
	callerFrame := currentThread.TopFrame()        // peek caller/invoker frame
	retVal := currentFrame.OperandStack().PopInt() // pop int from currentFrame
	callerFrame.OperandStack().PushInt(retVal)     // push int to caller frame's stack
}

func (r *IRETURN) Opcode() uint8 {
	return 0xAC
}

// LRETURN control long
// opcode = 0xAD
type LRETURN struct{ base.NoOperandsInstruction }

func (l *LRETURN) Execute(frame *runtime.Frame) {
	thread := frame.Thread()
	currentFrame := thread.PopFrame()
	invokerFrame := thread.TopFrame()
	retVal := currentFrame.OperandStack().PopLong()
	invokerFrame.OperandStack().PushLong(retVal)
}

func (r *LRETURN) Opcode() uint8 {
	return 0xAD
}

// FRETURN control float
// opcode = 0xAE
type FRETURN struct{ base.NoOperandsInstruction }

func (f *FRETURN) Execute(frame *runtime.Frame) {
	thread := frame.Thread()
	currentFrame := thread.PopFrame()
	invokerFrame := thread.TopFrame()
	retVal := currentFrame.OperandStack().PopFloat()
	invokerFrame.OperandStack().PushFloat(retVal)
}

func (r *FRETURN) Opcode() uint8 {
	return 0xAE
}

// DRETURN control double
// opcode = 0xAF
type DRETURN struct{ base.NoOperandsInstruction }

func (d *DRETURN) Execute(frame *runtime.Frame) {
	thread := frame.Thread()
	currentFrame := thread.PopFrame()
	invokerFrame := thread.TopFrame()
	retVal := currentFrame.OperandStack().PopDouble()
	invokerFrame.OperandStack().PushDouble(retVal)
}

func (r *DRETURN) Opcode() uint8 {
	return 0xAF
}

// ARETURN control Ref
// opcode = 0xB0
type ARETURN struct{ base.NoOperandsInstruction }

func (a *ARETURN) Execute(frame *runtime.Frame) {
	thread := frame.Thread()
	currentFrame := thread.PopFrame()
	invokerFrame := thread.TopFrame()
	retVal := currentFrame.OperandStack().PopRef()
	invokerFrame.OperandStack().PushRef(retVal)
}

func (r *ARETURN) Opcode() uint8 {
	return 0xB0
}
