package interpreter

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/instructions"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"testing"
)

// executeAndGetLocal0 執行字節碼並返回 locals[0] 的值
func executeAndGetLocal0(code []byte, maxLocals, maxStack uint16, debug bool) int32 {
	thread := runtime.NewThread()
	frame := thread.NewFrame(maxLocals, maxStack)
	thread.PushFrame(frame)

	reader := &base.BytecodeReader{}

	loopCount := 0

	for !thread.IsStackEmpty() {
		loopCount++
		fmt.Printf("<loop> -----------------------count: %v --------------------------- <loop> \n", loopCount)

		currentFrame := thread.CurrentFrame()

		pc := currentFrame.NextPC()
		thread.SetPC(pc)
		reader.Reset(code, pc)
		opcode := reader.ReadUint8()
		inst, _ := instructions.NewInstruction(opcode)
		inst.FetchOperands(reader)
		currentFrame.SetNextPC(reader.PC())

		if debug {
			name := opcodes.OpcodeNames[opcode]
			if name == "" {
				name = "???"
			}
			println("threadPC:", thread.PC(), "nextPC:", pc, "OP:", name)
			fmt.Printf("LocalVars: %v \n", frame.LocalVars())
		}

		inst.Execute(currentFrame)
	}

	return frame.LocalVars().GetInt(0)
}

// TestSimpleAdd 測試簡單加法: 1 + 2 = 3
func TestSimpleAdd(t *testing.T) {
	code := []byte{
		0x04, // iconst_1
		0x05, // iconst_2
		0x60, // iadd
		0x3B, // istore_0
		0xB1, // return
	}

	result := executeAndGetLocal0(code, 1, 2, false)
	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}
	t.Log("✓ 1 + 2 = 3")
}

// TestSubtract 測試減法: 10 - 3 = 7
func TestSubtract(t *testing.T) {
	code := []byte{
		0x10, 10, // bipush 10
		0x06, // iconst_3
		0x64, // isub
		0x3B, // istore_0
		0xB1, // return
	}

	result := executeAndGetLocal0(code, 1, 2, false)
	if result != 7 {
		t.Errorf("Expected 7, got %d", result)
	}
	t.Log("✓ 10 - 3 = 7")
}

// TestMultiply 測試乘法: 6 * 7 = 42
func TestMultiply(t *testing.T) {
	code := []byte{
		0x10, 6, // bipush 6
		0x10, 7, // bipush 7
		0x68, // imul
		0x3B, // istore_0
		0xB1, // return
	}

	result := executeAndGetLocal0(code, 1, 2, false)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
	t.Log("✓ 6 * 7 = 42")
}

// TestDivide 測試除法: 100 / 10 = 10
func TestDivide(t *testing.T) {
	code := []byte{
		0x10, 100, // bipush 100
		0x10, 10, // bipush 10
		0x6C, // idiv
		0x3B, // istore_0
		0xB1, // return
	}

	result := executeAndGetLocal0(code, 1, 2, false)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}
	t.Log("✓ 100 / 10 = 10")
}

// TestNegation 測試取負: -(42) = -42
func TestNegation(t *testing.T) {
	code := []byte{
		opcodes.BIPUSH, 42, // bipush 42
		opcodes.INEG,     // ineg
		opcodes.ISTORE_0, // istore_0
		opcodes.RETURN,   // return
	}

	result := executeAndGetLocal0(code, 1, 2, false)
	if result != -42 {
		t.Errorf("Expected -42, got %d", result)
	}
	t.Log("✓ -(42) = -42")
}

// TestLocalVariables 測試局部變量操作
// int a = 5; int b = 10; int c = a + b;
func TestLocalVariables(t *testing.T) {
	code := []byte{
		0x08,     // iconst_5   -> stack: [5]
		0x3B,     // istore_0   -> locals[0] = 5
		0x10, 10, // bipush 10  -> stack: [10]
		0x3C, // istore_1   -> locals[1] = 10
		0x1A, // iload_0    -> stack: [5]
		0x1B, // iload_1    -> stack: [5, 10]
		0x60, // iadd       -> stack: [15]
		0x3B, // istore_0   -> locals[0] = 15 (結果)
		0xB1, // return
	}

	result := executeAndGetLocal0(code, 3, 2, false)
	if result != 15 {
		t.Errorf("Expected 15, got %d", result)
	}
	t.Log("✓ a=5, b=10, c=a+b => c=15")
}

// TestIINC 測試 iinc 指令
// int i = 0; i++; i++; i++;
func TestIINC(t *testing.T) {
	code := []byte{
		0x03,       // iconst_0
		0x3B,       // istore_0  -> i = 0
		0x84, 0, 1, // iinc 0, 1 -> i = 1
		0x84, 0, 1, // iinc 0, 1 -> i = 2
		0x84, 0, 1, // iinc 0, 1 -> i = 3
		0xB1, // return
	}

	result := executeAndGetLocal0(code, 1, 1, false)
	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}
	t.Log("✓ i=0; i++; i++; i++ => i=3")
}

// TestComplexExpression 測試複雜表達式
// result = (10 + 20) * 3 - 5 = 85
func TestComplexExpression(t *testing.T) {
	code := []byte{
		0x10, 10, // bipush 10
		0x10, 20, // bipush 20
		0x60, // iadd -> 30
		0x06, // iconst_3
		0x68, // imul -> 90
		0x08, // iconst_5
		0x64, // isub -> 85
		0x3B, // istore_0
		0xB1, // return
	}

	result := executeAndGetLocal0(code, 1, 3, false)
	if result != 85 {
		t.Errorf("Expected 85, got %d", result)
	}
	t.Log("✓ (10 + 20) * 3 - 5 = 85")
}

// TestSimpleIf 測試簡單 if 語句
// 使用 IFEQ: if (value == 0) result = 1; else result = 2;
func TestSimpleIf(t *testing.T) {
	// 測試 value = 0 的情況
	code := []byte{
		opcodes.ICONST_0,         // 0: iconst_0      value = 0
		opcodes.IFNE, 0x00, 0x06, // 1: ifne +6       如果 != 0 跳到 7 --> fetch 2 bytes 當作 offset 這裡會是 6，pop 出 val 0，當 val != 0 時會跳到 初始 pc+6 = 0+6 = 6, 但是這裡為 0 所以不跳行
		opcodes.ICONST_1, // 4: iconst_1      result = 1
		opcodes.ISTORE_0, // 5: istore_0
		opcodes.RETURN,   // 6: return
		opcodes.ICONST_2, // 7: iconst_2      result = 2
		opcodes.ISTORE_0, // 8: istore_0
		opcodes.RETURN,   // 9: return
	}

	result := executeAndGetLocal0(code, 1, 1, true)
	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
	t.Log("✓ if (0 == 0) => 1")
}

// TestSimpleIf 測試簡單 if 語句 - 2
// 使用 IFEQ: if (value == 0) result = 1; else result = 2;
func TestSimpleIf_2(t *testing.T) {
	// 測試 value = 1 的情況
	code := []byte{
		opcodes.ICONST_1,         // 0: iconst_1      value = 1
		opcodes.IFNE, 0x00, 0x06, // 1: ifne +6       如果 != 0 跳到 7 --> fetch 2 bytes 當作 offset 這裡會是 6，pop 出 val 0，當 val != 0 時會跳到 thread.pc+6 = 1+6 = 7
		opcodes.ICONST_1, // 4: iconst_1      result = 1
		opcodes.ISTORE_0, // 5: istore_0
		opcodes.RETURN,   // 6: return
		opcodes.ICONST_2, // 7: iconst_2      result = 2
		opcodes.ISTORE_0, // 8: istore_0
		opcodes.RETURN,   // 9: return
	}

	result := executeAndGetLocal0(code, 1, 1, true)
	if result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
	t.Log("✓ if (1 != 0) => 2")
}

// TestSimpleLoop 測試簡單循環
// sum = 0; for (i = 1; i <= 3; i++) sum += i;
// 計算 1+2+3 = 6
func TestSimpleLoop(t *testing.T) {
	// 仔細計算每個指令的 PC 位置
	// PC  Instruction
	// 0   iconst_0
	// 1   istore_0
	// 2   iconst_1
	// 3   istore_1
	// 4   iload_1
	// 5   iconst_3
	// 6   if_icmpgt (3 bytes: opcodes + 2 byte offset)
	// 9   iload_0
	// 10  iload_1
	// 11  iadd
	// 12  istore_0
	// 13  iinc (3 bytes)
	// 16  goto (3 bytes) - 要跳到 PC=4, offset = 4 - 16 = -12
	// 19  return

	code := []byte{
		opcodes.ICONST_0, // 0: iconst_0 -> 推 0 到 stack
		opcodes.ISTORE_0, // 1: istore_0 -> sum = 0 -> pop 0 出來，放進 LocalVars[0]
		opcodes.ICONST_1, // 2: iconst_1 -> 推 1 到 stack
		opcodes.ISTORE_1, // 3: istore_1   i = 1  -> pop 1 出來，放進 LocalVars[1]

		// []LocalVars = {0, 1}
		// LocalVars 第一個元素是 sum
		// LocalVars 第二個元素是 i

		opcodes.ILOAD_1,               // 4: iload_1    load i -> 把 i load 到 stack 中
		opcodes.ICONST_3,              // 5: iconst_3   load 3 -> 把 3 推入 stack
		opcodes.IF_ICMPGT, 0x00, 0x0D, // 6: if_icmpgt +13(0x0D)  如果 i>3 跳到 19 (return)

		opcodes.ILOAD_0,  // 9: iload_0    load sum
		opcodes.ILOAD_1,  // 10: iload_1   load i
		opcodes.IADD,     // 11: iadd      sum + i
		opcodes.ISTORE_0, // 12: istore_0  sum = sum + i

		opcodes.IINC, 0x01, 0x01, // 13: iinc 1, 1  i++

		opcodes.GOTO, 0xFF, 0xF4, // 16: goto -12   跳到 4 (16 + (-12) = 4)

		opcodes.RETURN, // 19: return
	}

	result := executeAndGetLocal0(code, 2, 2, true)
	if result != 6 {
		t.Errorf("Expected 6, got %d", result)
	}
	t.Log("✓ sum(1..3) = 6")
}
