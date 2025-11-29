package interpreter

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/instructions"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// Interpret Bytecode interpret
// code: bytecode - from Code Attributes
// maxLocals: size of LocalVars
// maxStack: max size of opStack
// args: method args (if exists)
// debug: display debug message info
func Interpret(code []byte, maxLocals, maxStack uint16, debug bool) {
	// 1. create thread
	thread := runtime.NewThread()

	// 2. create frame
	frame := thread.NewFrame(maxLocals, maxStack)
	thread.PushFrame(frame)

	// 3. start execute
	loop(thread, code, debug)
}

// InterpretWithArgs run with args（for testing）
func InterpretWithArgs(code []byte, maxLocals, maxStack uint16, args []int32, debug bool) int32 {
	// 1. create thread
	thread := runtime.NewThread()

	// 2. create frame
	frame := thread.NewFrame(maxLocals, maxStack)
	thread.PushFrame(frame)

	// 3. put args into LocalVars
	for i, arg := range args {
		frame.LocalVars().SetInt(uint(i), arg)
	}

	// 4. start execute
	loop(thread, code, debug)

	// 5. return result（if have）
	size, _ := frame.OperandStack().Size()
	if size > 0 {
		// actually main method should be void, this is only for testing
		return frame.OperandStack().PopInt()
	}

	return 0
}

// loop interpreter main logic
// Fetch -> Decode -> Execute -> Fetch ...
func loop(thread *runtime.Thread, code []byte, debug bool) {
	reader := &base.BytecodeReader{}

	// check is end
	// when func returned, stack will be empty (for main method)
	// or current frame is not origin frame (for not main method)
	for !thread.IsStackEmpty() {
		// get current frame
		frame := thread.CurrentFrame()

		// calculate PC
		pc := frame.NextPC()
		thread.SetPC(pc)

		// Fetch: 1 byte opcodes
		reader.Reset(code, pc)
		opcode := reader.ReadUint8()

		// Decode:
		instruction, err := instructions.NewInstruction(opcode)
		if err != nil {
			fmt.Errorf("error while parsing instruction: %s", err)
			panic("instruction parse error")
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
