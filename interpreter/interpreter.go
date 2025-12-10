package interpreter

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/instructions"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
	"os"
)

// Interpret Bytecode interpret
func Interpret(method *method_area.Method, debug bool) {
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

	mainMethodFrame := thread.TopFrame()

	// check is end
	// when func returned, stack will be empty (for main method)
	// or current frame is not origin frame (for not main method)
	for !thread.IsStackEmpty() {
		// get current frame
		frame := thread.CurrentFrame()

		// get bytecode from frame's method
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
		instruction.FetchOperands(reader) // fetch (index, offset) if required
		frame.SetNextPC(reader.PC())      // update PC (to next instruction)

		if debug {
			fmt.Println("<--------------------------------------------------------------------------------->")
			printDebug(pc, instruction, frame)
		}

		// Execute: perform instruction
		instruction.Execute(frame)
	}

	if debug {
		fmt.Println("================================================================")
		fmt.Println("GOGO JVM: Thread's JVMFrameStack is empty before exist LocalVarsTable:")
		for i, slot := range mainMethodFrame.LocalVars() {
			fmt.Printf("* Slot - %d:\n", i)

			if slot.Ref != nil {
				fmt.Printf("\t <REF>: %v \n", slot.Ref)
				fmt.Printf("\t\t\t Field Details: %v \n", slot.Ref.(*heap.Object).Fields())
			} else {
				fmt.Printf("\t <NUM>: %v \n", slot.Num)
			}
			fmt.Printf("\n")
		}
	}
}

// printDebug print debug info
func printDebug(pc int, inst base.Instruction, frame *runtime.Frame) {
	opName := opcodes.OpcodeNames[inst.Opcode()]
	if opName == "" {
		opName = fmt.Sprintf("unknown(0x%02X)", inst.Opcode())
	}

	fmt.Printf("method: %s (%s) | PC:%3d | %-12s | Stack: ", frame.Method().Name(), frame.Method().Class().Name(), pc, opName)
	printStack(frame.OperandStack())
	fmt.Println()
	printLocalVars(frame.LocalVars())
}

func printLocalVars(vars rtcore.Slots) {
	fmt.Printf("* LocarVars=%v \n", vars)
}

// printStack print stack details
func printStack(stack *runtime.OperandStack) {
	if stack == nil {
		panic("Interpret error, OperandStack is nil!")
	}
	currentSize, maxSize := stack.Size()
	fmt.Printf("[currentSize=%d, maxSize=%d]\n", currentSize, maxSize)
	fmt.Printf("stack: %v \n", stack)
}
