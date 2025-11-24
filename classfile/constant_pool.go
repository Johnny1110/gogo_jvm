package classfile

import "fmt"

// constant pool (JVM Spec 4.4)
// Java 的不同類型信息（className, methodName, String）在常量池中用不同標籤區分
const (
	CONSTANT_Utf8               = 1
	CONSTANT_Integer            = 3
	CONSTANT_Float              = 4
	CONSTANT_Long               = 5
	CONSTANT_Double             = 6
	CONSTANT_Class              = 7
	CONSTANT_String             = 8
	CONSTANT_Fieldref           = 9
	CONSTANT_Methodref          = 10
	CONSTANT_InterfaceMethodref = 11
	CONSTANT_NameAndType        = 12
	CONSTANT_MethodHandle       = 15
	CONSTANT_MethodType         = 16
	CONSTANT_InvokeDynamic      = 18
)

// ConstantInfo interface for all constant
type ConstantInfo interface {
	// every constant know how to read themselves by ClassReader
	readInfo(reader *ClassReader)
	String() string
}

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

func readConstantInfo(reader *ClassReader, cp ConstantPool) ConstantInfo {
	tag := reader.readU1()
	info := newConstantInfo(tag, cp)
	if info == nil {
		panic(fmt.Sprintf("java.lang.ClassFormatError: constant pool tag %d", tag))
	}
	info.readInfo(reader)
	return info
}
