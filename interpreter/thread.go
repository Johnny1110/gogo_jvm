package interpreter

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/instructions"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/references"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// ============================================================
// Thread Executor - v0.4.0
// ============================================================
// solve runtime -> interpreter dependency

func init() {
	// setup thread executor
	runtime.SetThreadExecutor(executeThread)
}

// executeThread export method
// this will be called by JVMThread.runInternal()
func executeThread(jvmThread *runtime.JVMThread, runMethod *method_area.Method) {
	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - executeThread: Executing Thread [%d] %s, method: %s.%s\n",
			jvmThread.ID(), jvmThread.Name(),
			runMethod.Class().Name(), runMethod.Name())
	}

	// create method call frame
	frame := runtime.NewFrameWithMethodAndExHandler(
		getThreadAdapter(jvmThread),
		runMethod,
		references.ThrowException,
	)

	// setup this (java.lang.Thread)
	frame.LocalVars().SetRef(0, jvmThread.JavaThreadObj())

	jvmThread.PushFrame(frame)

	// loop read execute opcode
	loopForThread(jvmThread)
}

// loopForThread thread interpreter
func loopForThread(jvmThread *runtime.JVMThread) {
	reader := &base.BytecodeReader{}

	for !jvmThread.IsStackEmpty() {
		// get current Frame
		frame := jvmThread.CurrentFrame()

		// read method byte code
		bytecode := frame.Method().Code()

		// calculate PC
		pc := frame.NextPC()
		jvmThread.SetPC(pc)

		// TODO: remove this
		syncAdapterPC(jvmThread)

		// Fetch
		reader.Reset(bytecode, pc)
		opcode := reader.ReadUint8()

		// Decode
		instruction, err := instructions.NewInstruction(opcode)
		if err != nil {
			fmt.Printf("Error parsing instruction: %s\n", err)
			return
		}

		instruction.FetchOperands(reader)
		frame.SetNextPC(reader.PC())

		if global.DebugMode() {
			fmt.Printf("Thread [%s] PC:%3d | %s\n", jvmThread.Name(), pc, instruction)
		}

		// Execute
		instruction.Execute(frame)
	}

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - Thread [%s] execution completed\n", jvmThread.Name())
	}
}

func getThreadAdapter(jvmThread *runtime.JVMThread) *runtime.Thread {
	return runtime.GetThreadManager().GetThreadAdapter(jvmThread)
}

func syncAdapterPC(jvmThread *runtime.JVMThread) {
	adapter := getThreadAdapter(jvmThread)
	if adapter != nil {
		adapter.SetPC(jvmThread.PC())
	}
}
