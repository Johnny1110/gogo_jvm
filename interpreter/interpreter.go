package interpreter

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/instructions"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/instructions/references"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
	"os"
)

// Interpret Bytecode interpret
func Interpret(method *method_area.Method, debug bool) {
	// 1. create main thread
	thread := runtime.NewMainThread()

	// 2. create frame
	frame := thread.NewFrameWithMethodAndExHandler(method, references.ThrowException)
	thread.PushFrame(frame)

	// 3. start execute
	loop(thread, debug)
}

// loop interpreter main logic
// Fetch -> Decode -> Execute -> Fetch ...
func loop(thread *runtime.JVMThread, debug bool) {
	mainMethodFrame := thread.TopFrame()

	executeLoop(thread, func(frame *runtime.Frame, pc int, instruction base.Instruction) {
		if global.DebugMode() {
			fmt.Printf("@@ DEBUG - interoreter loop, frame method: %s, class: %s \n", frame.Method().Name(), frame.Method().Class().Name())
		}

		if debug {
			fmt.Println("<--------------------------------------------------------------------------------->")
			printDebug(pc, instruction, frame)
		}
	}, func(err error) {
		fmt.Printf("Error parsing instruction: %s\n", err)
		os.Exit(1)
	})

	if debug {
		fmt.Println("================================================================")
		fmt.Println("GOGO JVM: Thread's JVMFrameStack is empty before exist LocalVarsTable:")
		for i, slot := range mainMethodFrame.LocalVars() {
			fmt.Printf("* Slot - %d:\n", i)

			if slot.Ref != nil {
				fmt.Printf("\t <REF>: %v \n", slot.Ref)
				obj := slot.Ref.(*heap.Object)
				if obj.IsArray() {
					fmt.Printf("\t\t\t Array Elements: %v \n", obj.Extra())
				} else {
					fmt.Printf("\t\t\t Object Field Details: %v \n", slot.Ref.(*heap.Object).Fields())
				}
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

// executeLoop common execution loop for both main thread and other threads
func executeLoop(thread *runtime.JVMThread, debugCallback func(*runtime.Frame, int, base.Instruction), errorCallback func(error)) {
	reader := &base.BytecodeReader{}

	for !thread.IsStackEmpty() {
		frame := thread.CurrentFrame()
		bytecode := frame.Method().Code()
		pc := frame.NextPC()
		thread.SetPC(pc)

		reader.Reset(bytecode, pc)
		opcode := reader.ReadUint8()

		instruction, err := instructions.NewInstruction(opcode)
		if err != nil {
			errorCallback(err)
			return
		}

		instruction.FetchOperands(reader)
		frame.SetNextPC(reader.PC())

		if debugCallback != nil {
			debugCallback(frame, pc, instruction)
		}

		instruction.Execute(frame)
	}
}
