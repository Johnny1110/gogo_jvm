package interpreter

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/instructions"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/rtda/heap"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"os"
)

// Interpret Bytecode interpret
// code: bytecode - from Code Attributes
// maxLocals: size of LocalVars
// maxStack: max size of opStack
// args: method args (if exists)
// debug: display debug message info
func Interpret(method *heap.Method, debug bool) {
	// 1. create thread
	thread := runtime.NewThread()

	// 2. create frame
	frame := thread.NewFrameWithMethod(method)
	thread.PushFrame(frame)

	// 3. start execute
	loop(thread, debug)
}

// loop interpreter main logic
// Fetch -> Decode -> Execute -> Fetch ...
func loop(thread *runtime.Thread, debug bool) {
	reader := &base.BytecodeReader{}

	// check is end
	// when func returned, stack will be empty (for main method)
	// or current frame is not origin frame (for not main method)
	for !thread.IsStackEmpty() {
		// get current frame
		frame := thread.CurrentFrame()

		// get bytecode from frame
		bytecode := frame.Method().Code()

		// calculate PC
		pc := frame.NextPC()
		thread.SetPC(pc)

		// Fetch: 1 byte opcodes
		reader.Reset(bytecode, pc)
		opcode := reader.ReadUint8()

		// Decode:
		instruction, err := instructions.NewInstruction(opcode)
		if err != nil {
			fmt.Printf("Error parsing instruction: %s\n", err)
			os.Exit(1)
		}
		instruction.FetchOperands(reader) // fetch index, offset if required
		frame.SetNextPC(reader.PC())      // update PC (to next instruction)

		if debug {
			printDebug(pc, opcode, instruction, frame)
		}

		// Execute: perform instruction
		instruction.Execute(frame)
	}
}

// printDebug print debug info
func printDebug(pc int, opcode uint8, inst base.Instruction, frame *runtime.Frame) {
	opName := opcodes.OpcodeNames[opcode]
	if opName == "" {
		opName = fmt.Sprintf("unknown(0x%02X)", opcode)
	}

	fmt.Printf("PC:%3d | %-12s | Stack: ", pc, opName)
	printStack(frame.OperandStack())
	fmt.Println()
}

// printStack print stack details
func printStack(stack *runtime.OperandStack) {
	currentSize, maxSize := stack.Size()
	fmt.Printf("[currentSize=%d, maxSize=%d]\n]", currentSize, maxSize)
}
