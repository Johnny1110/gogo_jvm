# Bytecode interpreter 解釋器

<br>

這是 JVM 的核心，解釋器負責執行字節碼指令。

**工作流程**

Fetch-Decode-Execute 循環：

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│   ┌─────────┐    ┌─────────┐    ┌─────────┐               │
│   │  Fetch  │ →  │ Decode  │ →  │ Execute │               │
│   │ 取指令   │    │  解碼   │    │   執行   │               │
│   └─────────┘    └─────────┘    └─────────┘               │
│        ↑                              │                    │
│        └──────────────────────────────┘                    │
│                   循環直到方法返回                           │
│                                                             │
│   Fetch:  從 PC 指向的位置讀取 opcode                        │
│   Decode: 根據 opcode 找到對應的指令實現                      │
│   Execute: 執行指令（操作棧、局部變量、PC）                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 半成品版本

裡我先簡化處理，因為完整實現需要 ClassLoader。等我們做到方法調用時，會回來修正 Frame 的設計。

目標：
* 單方法執行，code 作為參數                               │
* 驗證 Fetch-Decode-Execute 正確

```go
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

		// Fetch: 1 byte opcode
		reader.Reset(code, pc)
		opcode := reader.ReadUint8()

		// Decode:
		instruction := instructions.NewInstruction(opcode)
		instruction.FetchOperands(reader) // fetch index, offset if required
		frame.SetNextPC(reader.PC())      // update PC (to next instruction)

		if debug {
			printDebug(pc, opcode, instruction, frame)
		}

		// Execute: perform instruction
		instruction.Execute(frame)
	}
}
```

<br>

注意：現在無論怎麼 `loop()` 你會發現使用的 bytecode 都是一開始傳入的版本，不無會隨 call function 載入新的 code (__因為我們還沒有實現 call function 的邏輯__)。