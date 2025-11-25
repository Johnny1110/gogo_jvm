package classfile

import "fmt"

// constant pool (JVM Spec 4.4)
// Java 的不同類型信息（className, methodName, String）在常量池中用不同標籤區分
type ConstantTag uint8

const (
	CONSTANT_Utf8               ConstantTag = 1
	CONSTANT_Integer            ConstantTag = 3
	CONSTANT_Float              ConstantTag = 4
	CONSTANT_Long               ConstantTag = 5
	CONSTANT_Double             ConstantTag = 6
	CONSTANT_Class              ConstantTag = 7
	CONSTANT_String             ConstantTag = 8
	CONSTANT_Fieldref           ConstantTag = 9
	CONSTANT_Methodref          ConstantTag = 10
	CONSTANT_InterfaceMethodref ConstantTag = 11
	CONSTANT_NameAndType        ConstantTag = 12
	CONSTANT_MethodHandle       ConstantTag = 15
	CONSTANT_MethodType         ConstantTag = 16
	CONSTANT_InvokeDynamic      ConstantTag = 18
)

func uint8ToTag(input uint8) (ConstantTag, error) {
	switch input {
	case 1:
		return CONSTANT_Utf8, nil
	case 3:
		return CONSTANT_Integer, nil
	case 4:
		return CONSTANT_Float, nil
	case 5:
		return CONSTANT_Long, nil
	case 6:
		return CONSTANT_Double, nil
	case 7:
		return CONSTANT_Class, nil
	case 8:
		return CONSTANT_String, nil
	case 9:
		return CONSTANT_Fieldref, nil
	case 10:
		return CONSTANT_Methodref, nil
	case 11:
		return CONSTANT_InterfaceMethodref, nil
	case 12:
		return CONSTANT_NameAndType, nil
	case 15:
		return CONSTANT_MethodHandle, nil
	case 16:
		return CONSTANT_MethodType, nil
	case 18:
		return CONSTANT_InvokeDynamic, nil
	default:
		return 0, fmt.Errorf("unknown tag: %d", input)
	}
}
