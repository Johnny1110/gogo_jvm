package classfile

// ConstantPool
// 注意：常量池索引從 1 開始，而不是 0, 這是 JVM 規範的歷史遺留設計
type ConstantPool []ConstantInfo

func readConstantPool(reader *ClassReader) ConstantPool {
	constantPoolCount := int(reader.readU2())

	// why constantPoolCount not constantPoolCount-1?
	// because constantPool index start from 1, 0 is invalid.
	constantPool := make([]ConstantInfo, constantPoolCount)

	for i := 1; i < constantPoolCount; i++ {
		constantPool[i] = readConstantInfo(reader, constantPool)

		// caution: long and double take 2 positions (這是另一個歷史遺留問題，為了在 32 位系統上對齊)
		switch constantPool[i].(type) {
		case *ConstantLongInfo, *ConstantDoubleInfo:
			i++
		}
	}

	return constantPool
}
