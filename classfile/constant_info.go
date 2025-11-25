package classfile

import (
	"fmt"
	"math"
)

// ConstantInfo interface for all constant
type ConstantInfo interface {
	// every constant know how to read themselves by ClassReader
	readInfo(reader *ClassReader)
	String() string
}

func readConstantInfo(reader *ClassReader, cp ConstantPool) ConstantInfo {
	tagUint8 := reader.readU1()
	if tag, err := uint8ToTag(tagUint8); err == nil {
		info := newConstantInfo(tag, cp)
		if info == nil {
			panic(fmt.Sprintf("java.lang.ClassFormatError: constant pool tag %d", tag))
		}
		info.readInfo(reader)
		return info
	} else {
		panic("java.lang.ClassFormatError: " + err.Error())
	}
}

func newConstantInfo(tag ConstantTag, cp ConstantPool) ConstantInfo {
	switch tag {
	case CONSTANT_Utf8:
		return &ConstantUtf8Info{}
	case CONSTANT_Integer:
		return &ConstantIntegerInfo{}
	case CONSTANT_Float:
		return &ConstantFloatInfo{}
	case CONSTANT_Long:
		return &ConstantLongInfo{}
	case CONSTANT_Double:
		return &ConstantDoubleInfo{}
	case CONSTANT_Class:
		return &ConstantClassInfo{cp: cp}
	case CONSTANT_String:
		return &ConstantStringInfo{cp: cp}
	case CONSTANT_Fieldref:
		return &ConstantFieldrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_Methodref:
		return &ConstantMethodrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_InterfaceMethodref:
		return &ConstantInterfaceMethodrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_NameAndType:
		return &ConstantNameAndTypeInfo{}
	default:
		// MVP 階段暫不支持 MethodHandle, MethodType, InvokeDynamic
		return nil
	}
}

// ConstantUtf8Info UTF8 string constant
type ConstantUtf8Info struct {
	str string
}

func (c *ConstantUtf8Info) readInfo(reader *ClassReader) {
	// read 2 unit (2 bytes) as length
	length := uint32(reader.readU2())
	bytes := reader.readBytes(length)
	c.str = string(bytes)
}

func (c *ConstantUtf8Info) String() string {
	return c.str
}

// ConstantIntegerInfo Integer (32bit) constant
type ConstantIntegerInfo struct {
	val int32
}

func (c *ConstantIntegerInfo) readInfo(reader *ClassReader) {
	// read 4 unit (4 bytes)
	c.val = int32(reader.readU4())
}

func (c *ConstantIntegerInfo) String() string {
	return fmt.Sprintf("%d", c.val)
}

// ConstantFloatInfo Float (32bit) constant
type ConstantFloatInfo struct {
	val float32
}

func (c *ConstantFloatInfo) readInfo(reader *ClassReader) {
	bytes := reader.readU4()
	c.val = math.Float32frombits(bytes)
}

func (c *ConstantFloatInfo) String() string {
	return fmt.Sprintf("%f", c.val)
}

// ConstantLongInfo Long (64bit) constant
type ConstantLongInfo struct {
	val int64
}

func (c *ConstantLongInfo) readInfo(reader *ClassReader) {
	c.val = int64(reader.readU8())
}

func (c *ConstantLongInfo) String() string {
	return fmt.Sprintf("%d", c.val)
}

// ConstantDoubleInfo (64bit) constant
type ConstantDoubleInfo struct {
	val float64
}

func (c *ConstantDoubleInfo) readInfo(reader *ClassReader) {
	bytes := reader.readU8()
	c.val = math.Float64frombits(bytes)
}

func (c *ConstantDoubleInfo) String() string {
	return fmt.Sprintf("%f", c.val)
}

// ConstantClassInfo Class Constant
type ConstantClassInfo struct {
	cp        ConstantPool
	nameIndex uint16 // point to UTF8 constant
}

func (c *ConstantClassInfo) readInfo(reader *ClassReader) {
	c.nameIndex = reader.readU2()
}

func (c *ConstantClassInfo) String() string {
	return c.Name()
}

func (c *ConstantClassInfo) Name() string {
	return getUtf8(c.cp, c.nameIndex)
}

// ConstantStringInfo String constant
type ConstantStringInfo struct {
	cp          ConstantPool
	stringIndex uint16
}
