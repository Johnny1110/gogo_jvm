# Bytecode instructions 解釋器

<br>

---

<br>

## 解釋器的核心：Fetch-Decode-Execute 循環

```
┌─────────────────────────────────────────────────────────────────┐
│                    解釋器工作原理                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│    ┌──────────────────────────────────────────────────────┐     │
│    │                                                      │     │
│    │    ┌─────────┐    ┌─────────┐    ┌─────────┐         │     │
│    │    │  Fetch  │ →  │ Decode  │ →  │ Execute │         │     │
│    │    └─────────┘    └─────────┘    └─────────┘         │     │
│    │         ↑                              │             │     │
│    │         └──────────────────────────────┘             │     │
│    │                    LOOP                              │     │
│    └──────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────────┘
```

1. Fetch：從 PC 指向的位置讀取 opcode
2. Decode：根據 opcode 找到對應的指令實現
3. Execute：執行指令，可能修改棧、局部變量、PC
4. 重複直到方法返回

<br>
<br>

## JVM 指令集設計哲學

**為什麼用單字節 Opcode？**

```
┌─────────────────────────────────────────────────────────────────┐
│  JVM 指令格式                                                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌──────────┬──────────────────────────────────┐              │
│   │  Opcode  │         Operands (可選)           │              │
│   │  1 byte  │        0 ~ N bytes               │              │
│   └──────────┴──────────────────────────────────┘              │
│                                                                 │
│   例子：                                                        │
│   ┌──────┐                                                     │
│   │ 0x03 │                    iconst_0 (無操作數)               │
│   └──────┘                                                     │
│                                                                 │
│   ┌──────┬──────┐                                              │
│   │ 0x10 │  42  │              bipush 42 (1字節操作數)          │
│   └──────┴──────┘                                              │
│                                                                 │
│   ┌──────┬──────┬──────┐                                       │
│   │ 0x11 │ 0x01 │ 0x00 │       sipush 256 (2字節操作數)         │
│   └──────┴──────┴──────┘                                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

為什麼？

* 1 byte = 256 種可能的 opcode
* JVM 實際用了約 200 個
* 指令緊湊，class 文件更小
* 網絡傳輸更快（過去在 Java Applet 流行的時代很重要）


<br>

## 指令命名規則

**JVM 指令名有規律可循：**

```
┌─────────────────────────────────────────────────────────────────┐
│  指令命名模式：<類型前綴><操作>                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  類型前綴：                                                      │
│  ┌─────┬───────────┐                                           │
│  │  i  │  int      │                                           │
│  │  l  │  long     │                                           │
│  │  f  │  float    │                                           │
│  │  d  │  double   │                                           │
│  │  a  │  reference│  (address)                                │
│  │  b  │  byte     │                                           │
│  │  c  │  char     │                                           │
│  │  s  │  short    │                                           │
│  └─────┴───────────┘                                           │
│                                                                 │
│  例子：                                                         │
│  iadd = int add      (整數加法)                                 │
│  ladd = long add     (長整數加法)                               │
│  fadd = float add    (浮點加法)                                 │
│  iload = int load    (載入整數到棧)                              │
│  aload = address load (載入引用到棧)                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

<br>
<br>

---

<br>
<br>

## ByteCodeReader - 字節碼讀取器

**設計背景：**

每個方法的 Code 屬性中包含字節碼（bytecode），解釋器需要按順序讀取這些字節碼並執行。

<br>

BytecodeReader 封裝了字節碼的讀取邏輯，提供：

* 讀取不同長度的數據（1, 2, 4 字節）
* 跟蹤當前讀取位置（PC）
* 支持跳轉（修改 PC）

<br>

與 ClassReader 的區別：

* ClassReader：讀取整個 .class 文件，只讀一次
* BytecodeReader：讀取方法的字節碼，可能來回跳轉

```
 ┌────────────────────────────────────────────────────────────┐
 │                    BytecodeReader                          │
 ├────────────────────────────────────────────────────────────┤
 │  code: [0x03, 0x3C, 0x04, 0x3D, 0x1B, 0x1C, 0x60, 0xAC]   │
 │                     ↑                                      │
 │                    pc=0                                    │
 │                                                            │
 │  ReadUint8() → 返回 0x03，pc 變成 1                         │
 │  ReadUint8() → 返回 0x3C，pc 變成 2                         │
 │  ...                                                       │
 └────────────────────────────────────────────────────────────┘
```

<br>
<br>

## Instruction 指令接口

**設計背景：**

這是 JVM 指令的抽象表示。每一條 JVM 指令都需要實現這個接口。

* FetchOperands：負責「讀取」（從字節碼讀操作數）
* Execute：負責「執行」（操作棧和局部變量）

**解釋器的工作流程：**

```
for {                                                      
    pc := frame.NextPC()                                   
    reader.SetPC(pc)                                       
                                                              
    opcode := reader.ReadUint8()       1. 讀取 opcode       
    inst := decodeInstruction(opcode)  2. 解碼得到指令       
    inst.FetchOperands(reader)         3. 讀取操作數         
    inst.Execute(frame)                4. 執行指令          
}  
```

<br>
<br>

## 建立一些基礎指令 - Const (把 Const PUSH 到 Stack)

<br>

### CONST 指令

把一些常用的 int (-1~5) long, float double 壓入 Stack 時使用。

<br>

### IPUSH 指令

使用 CONST 無法滿足時，可以使用 `BIPUSH` 和 `SIPUSH`

* `BIPUSH`: Byte Immediate PUSH
  - 操作數是 1 byte 有符號數
  - 範圍：-128 ~ 127

* `SIPUSH`: Short Immediate PUSH
  - 操作數是 2 bytes 有符號數
  - 範圍：-32768 ~ 32767

舉例：
```java
int a = 100;     編譯成 bipush 100
int b = 1000;    編譯成 sipush 1000
int c = 100000;  編譯成 ldc（從常量池載入）
```

<br>

### LOAD 指令

JVM 中最常用的指令之一

**工作原理：**

把值從局部變量表複製一份放入 Stack
```
 ┌─────────────────────────────────────────────────────────┐
 │  局部變量表                    操作數棧                   │
 │  ┌─────┬─────┬─────┐         ┌─────┐                   │
 │  │  5  │ 10  │ 15  │         │     │                   │
 │  └─────┴─────┴─────┘         └─────┘                   │
 │    [0]   [1]   [2]                                     │
 │                                                        │
 │  執行 iload_1 後：                                      │
 │  ┌─────┬─────┬─────┐         ┌─────┐                   │
 │  │  5  │ 10  │ 15  │         │ 10  │ ← 從 [1] 複製     │
 │  └─────┴─────┴─────┘         └─────┘                   │
 │    [0]   [1]   [2]                                     │
 │                                                        │
 │  注意：局部變量表的值不變（是複製，不是移動）                 │
 └─────────────────────────────────────────────────────────┘
```

<br>


### STORE 指令

**工作原理：**

把值從 stack 移動到 LocalVars 中

```
 ┌─────────────────────────────────────────────────────────┐
 │  局部變量表                    操作數棧                  │
 │  ┌─────┬─────┬─────┐         ┌─────┐                   │
 │  │  5  │  ?  │  ?  │         │ 42  │ ← 棧頂             │
 │  └─────┴─────┴─────┘         └─────┘                   │
 │    [0]   [1]   [2]                                     │
 │                                                        │
 │  執行 istore_1 後：                                      │
 │  ┌─────┬─────┬─────┐         ┌─────┐                   │
 │  │  5  │ 42  │  ?  │         │     │ ← 棧變空          │
 │  └─────┴─────┴─────┘         └─────┘                   │
 │    [0]   [1]   [2]                                     │
 │              ↑                                         │
 │             新值                                        │
 └─────────────────────────────────────────────────────────┘
```

<br>

### 算術指令 (math)

JVM 的算術運算完全基於操作數棧

例如計算 `3 + 5`：

```
 ┌─────────────────────────────────────────────────────────┐
 │  初始狀態              壓入 3            壓入 5          │
 │  ┌─────┐            ┌─────┐          ┌─────┐           │
 │  │     │            │  3  │          │  5  │ ← top     │
 │  └─────┘            └─────┘          ├─────┤           │
 │                                      │  3  │           │
 │                                      └─────┘           │
 │                                                        │
 │  執行 iadd 後：                                         │
 │  ┌─────┐                                               │
 │  │  8  │ ← 3 + 5 的結果                                │
 │  └─────┘                                               │
 └─────────────────────────────────────────────────────────┘
```

注意：
* 運算結果會自動處理溢出（Java 定義的行為）
* 除法時除數為 0 會拋出 ArithmeticException

<br>

### 跳轉與比較指令

比較大小與跳行相關

舉一個間單的範例，有一個簡單的 if 條件式：

ex: `if (a < b) { ... }`

編譯過後為：
```
compiled:
   iload_0       // load a
   iload_1       // load b
   if_icmpge L1  // if a >= b jump to L1（skip if block）
   ...           // if block logics
   L1:           // if ending
```

看一下 `if_icmpge` 指令 source code 實現：

```go
// IF_ICMPGE jump if v1 >= v2
// opcodes = 0xA2
type IF_ICMPGE struct{ base.BranchInstruction }

func (i *IF_ICMPGE) Execute(frame *runtime.Frame) {
  stack := frame.OperandStack()
  v2 := stack.PopInt() // pop 出 a
  v1 := stack.PopInt() // pop 出 b
  if v1 >= v2 {        // 進行比較 
    branch(frame, i.Offset) // 達成跳行條件，執行跳行
  }
}
```

跳行實作：

```go
// branch helper func: perform jump
func branch(frame *runtime.Frame, offset int) {
	pc := frame.Thread().PC() // Thread PC 記錄起始位置
	nextPC := pc + offset     // 起始位置 + 偏移量
	frame.SetNextPC(nextPC)   // 設定新的 PC 給 frame，下一次 Fetch 指令時，從新指定的 nextPC 開始
}
```

### RETURN 系列指令

當方法執行完畢時，需要返回到調用者

返回指令的工作：

1. 從當前棧幀的操作數棧彈出返回值（如果有）
2. 彈出當前棧幀
3. 把返回值壓入調用者棧幀的操作數棧（如果有）

```
 ┌─────────────────────────────────────────────────────────────┐
 │  方法 foo() 調用 bar()，bar() 返回 42                        │
 │                                                             │
 │  bar() 執行 ireturn 前：        bar() 返回後：               │
 │  ┌─────────────────┐           ┌─────────────────┐         │
 │  │ bar() frame     │           │ foo() frame     │         │
 │  │ opStack: [42]   │    →      │ opStack: [42]   │ ← 返回值│
 │  ├─────────────────┤           └─────────────────┘         │
 │  │ foo() frame     │                                       │
 │  │ opStack: []     │           bar 的棧幀已彈出              │
 │  └─────────────────┘                                       │
 └─────────────────────────────────────────────────────────────┘
```

<br>
<br>

## Instruction Factory 指令集工廠

使用工廠模式維護 Instruction 實現的取得

需要注意的是，指令在這裡分成兩類，一類是無狀態的 `NoOperandsInstruction`，一類是帶有狀態的 `BranchInstruction` 與 `IndexInstruction`。

* 無狀態的 Instruction 可以做成單例，節省運行成本。
* 有狀態的 Instruction (帶有 index, offset 等資訊) 則必須每次初始化一個新的

<br>
<br>

### INVOKE_STATIC 調用靜態方法

opcode = 0xB8

這是實現方法調用的核心指令

```
┌─────────────────────────────────────────────────────────┐
│  1. 從常量池獲取方法引用（MethodRef）                       │
│  2. 解析方法引用，找到目標方法（Method）                     │
│  3. 創建新的棧幀（Frame）                                 │
│  4. 從調用者的操作數棧彈出參數，放入新棧幀的局部變量表          │
│  5. 將新棧幀壓入 JVM Stack                                │
│  6. 下次循環時，解釋器就會執行新方法的字節碼                   │
└─────────────────────────────────────────────────────────┘
```

例子：invokestatic Calculator.add(II)I

```
  invokestatic #1  (調用 Calculator.add)
  │
  ▼
  MethodRef.ResolvedMethod() // 加載 add() 法
  │
  ├─ 1. ResolvedClass() → 加載 Calculator 類
  │
  ├─ 2. 檢查不是接口
  │
  └─ 3. lookupMethod("add", "(II)I")
  │
  └─ 在 Calculator.methods 中查找
  │
  ▼
  返回 Method 對象
```

<br>

這裡我想著重補充一個概念，上面提到 __從調用者的操作數棧彈出參數，放入新棧幀的局部變量表__。

* 在 `Method` 中我們可以由 `argSlotCount` 參數得知 invoke 該方法時需要傳遞多少參數。
* 在先前的計算中，理應把需要傳入目標方法得參數們 (slots) 都照順序押入 stack 中。
* 在建立 newFrame 時，就需要把參數 pop 出來，倒著放進 newFrame 的 LocalVars 中。

```
調用前：                      調用後：
 ┌─────────────┐              ┌─────────────┐
 │ main Frame  │              │  add Frame  │ ← 新的當前幀
 │ stack: [3,5]│              │ locals[0]=3 │
 └─────────────┘              │ locals[1]=5 │
                              ├─────────────┤
                              │ main Frame  │
                              │ stack: []   │
                              └─────────────┘
```

<br>

source code: 
```go
// InvokeMethod call method common func
// usage: invokestatic, invokevirtual
func InvokeMethod(invokerFrame *runtime.Frame, method *method_area.Method) {
	// 1, get current thread
	thread := invokerFrame.Thread()

	// 2. create a new frame (represent new method)
	newFrame := thread.NewFrameWithMethod(method)
	thread.PushFrame(newFrame)

	// 3. pass vars, we can know how many params should pass in by method.ArgSlotCount()
	argSlotCount := int(method.ArgSlotCount())
	if argSlotCount > 0 {
		for i := argSlotCount - 1; i >= 0; i-- {
			// the passing var will be standby in invokerFrame's stack.
			// pop it from invokerFrame's stack
			slot := invokerFrame.OperandStack().PopSlot()
			// put into newFrame's LocalVars
			newFrame.LocalVars().SetSlot(uint(i), slot)
		}
	}

	// no need to reset PC, new frame nextPC will be default 0
}
```

<br>
<br>

