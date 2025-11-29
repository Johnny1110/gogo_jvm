package main

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/instructions"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"io/ioutil"
	"os"
)

// ============================================================
// Gogo JVM - Phase 2.2
// ============================================================

// Read .class file and execute
// ┌─────────────────────────────────────────────────────────┐
// │  1. read .class file                                    │
// │  2. Parse to ClassFile struct                           │
// │  3. find main()                                         │
// │  4. Get bytecode, maxLocals, maxStack                   │
// │  5. Cread Thread and Frame                              │
// │  6. Do Fetch-Decode-Execute loop                        │
// │  7. print result                                        │
// └─────────────────────────────────────────────────────────┘

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	className := os.Args[1]

	// read class file
	classData, err := ioutil.ReadFile(className)
	if err != nil {
		fmt.Printf("Error reading classfile: %v\n", err)
		os.Exit(1)
	}

	classFilePath := os.Args[1]
	debug := len(os.Args) > 2 && os.Args[2] == "-debug"

	// 1. read classfile
	classData, err = ioutil.ReadFile(classFilePath)
	if err != nil {
		fmt.Printf("Error reading class file: %v\n", err)
		os.Exit(1)
	}

	// 2. parse classfile
	cf, err := classfile.Parse(classData)
	if err != nil {
		fmt.Printf("Error parsing class file: %v\n", err)
		os.Exit(1)
	}

	classfile.Debug(cf, true)

	// 3. find main()
	mainMethod := cf.GetMainMethod()
	if mainMethod == nil {
		fmt.Println("Error: No main method found!")
		fmt.Println("Main method signature must be: public static void main(String[] args)")
		os.Exit(1)
	}

	// 4. get code attribute
	codeAttr := mainMethod.CodeAttribute()
	if codeAttr == nil {
		fmt.Println("Error: main method has no Code attribute!")
		os.Exit(1)
	}

	bytecode := codeAttr.Code()
	maxLocals := codeAttr.MaxLocals()
	maxStack := codeAttr.MaxStack()

	if debug {
		fmt.Println("Bytecode:")
		printBytecode(bytecode)
		fmt.Println("============================================")
	}

	// 5. execute Fetch-Decode-Execute
	interpret(bytecode, maxLocals, maxStack, debug)
}

func interpret(bytecode []byte, maxLocals uint16, maxStack uint16, debug bool) {
	thread := runtime.NewThread()

	frame := thread.NewFrame(maxLocals, maxStack)
	thread.PushFrame(frame)

	reader := &base.BytecodeReader{}

	loopCount := 0

	for !thread.IsStackEmpty() {
		loopCount++
		fmt.Printf("====== <loop count: %v> ======= \n", loopCount)

		currentFrame := thread.CurrentFrame()
		pc := currentFrame.NextPC()
		thread.SetPC(pc)

		// Fetch
		reader.Reset(bytecode, pc)
		opcode := reader.ReadUint8()

		// Decode
		inst := instructions.NewInstruction(opcode)
		inst.FetchOperands(reader)
		currentFrame.SetNextPC(reader.PC())

		if debug {
			opName := opcodes.OpcodeNames[opcode]
			if opName == "" {
				opName = fmt.Sprintf("0x%02X", opcode)
			}
			fmt.Printf("  PC:%3d | %-12s | Locals: %v\n",
				pc, opName, getLocalsSnapshot(currentFrame, maxLocals))
		}

		// Execute
		inst.Execute(frame)

		// print final
		fmt.Println("Final local variables:")
		for i := uint16(0); i < maxLocals; i++ {
			val := frame.LocalVars().GetInt(uint(i))
			if val != 0 {
				fmt.Printf("  locals[%d] = %d\n", i, val)
			}
		}
	}
}

func getLocalsSnapshot(frame *runtime.Frame, maxLocals uint16) any {
	snapshot := make([]int32, maxLocals)
	for i := uint16(0); i < maxLocals; i++ {
		snapshot[i] = frame.LocalVars().GetInt(uint(i))
	}
	return snapshot
}

func printBytecode(bytecode []byte) {
	for i, b := range bytecode {
		if i > 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X ", b)
	}
	fmt.Println()
}

func printUsage() {
	fmt.Println("Gogo JVM - A simple JVM implementation in Go")
	fmt.Println()
	fmt.Println("Usage: gogo_jvm <classfile> [-debug]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gogo_jvm Test.class")
	fmt.Println("  gogo_jvm Test.class -debug")
}
