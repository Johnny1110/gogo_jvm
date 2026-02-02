package interpreter

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
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
		jvmThread,
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
	executeLoop(jvmThread, func(frame *runtime.Frame, pc int, instruction base.Instruction) {
		if global.DebugMode() {
			fmt.Printf("@@ DEBUG - loopFor Thread -> Thread [%s] PC:%v | %s\n", jvmThread.Name(), pc, opcodes.OpcodeNames[instruction.Opcode()])
		}
	}, func(err error) {
		fmt.Printf("Error parsing instruction: %s\n", err)
	})

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - Thread [%s] execution completed\n", jvmThread.Name())
	}
}
