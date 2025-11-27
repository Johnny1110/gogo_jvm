package interpreter

import "testing"

func TestSimpleAdd(t *testing.T) {
	// 字節碼：
	// iconst_1    (0x04) - 將 1 壓入棧
	// iconst_2    (0x05) - 將 2 壓入棧
	// iadd        (0x60) - 彈出兩個數相加，結果壓入棧
	// ireturn     (0xAC) - 返回 int 結果

	// 但是因為我們只有一個棧幀，ireturn 會導致棧空
	// 所以我們用 istore_0 存儲結果，然後手動檢查

	code := []byte{
		0x04, // iconst_1
		0x05, // iconst_2
		0x60, // iadd
		0x3B, // istore_0 (存入局部變量 index: 0)
		0xB1, // return (void 返回)
	}

	// maxLocals=1 (存結果), maxStack=2 (最多兩個數)
	result := InterpretWithArgs(code, 2, 2, []int32{}, true)

	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}
	t.Log("✓ 1 + 2 = 3")
}
