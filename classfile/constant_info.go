package classfile

import (
	"fmt"
	"math"
)

// ConstantInfo interface for all constants
type ConstantInfo interface {
	// every constants know how to read themselves by ClassReader
	readInfo(reader *ClassReader)
	String() string
	Tag() ConstantTag
}

func readConstantInfo(reader *ClassReader, cp ClassFileConstantPool) ConstantInfo {
	tagUint8 := reader.readU1()
	if tag, err := uint8ToTag(tagUint8); err == nil {
		info := newConstantInfo(tag, cp)
		if info == nil {
			panic(fmt.Sprintf("java.lang.ClassFormatError: constants pool tag %d", tag))
		}
		info.readInfo(reader)
		return info
	} else {
		panic("java.lang.ClassFormatError: " + err.Error())
	}
}

func newConstantInfo(tag ConstantTag, cp ClassFileConstantPool) ConstantInfo {
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
		return &ConstantFieldRefInfo{ConstantMemberRefInfo{cp: cp}}
	case CONSTANT_Methodref:
		return &ConstantMethodRefInfo{ConstantMemberRefInfo{cp: cp}}
	case CONSTANT_InterfaceMethodref:
		return &ConstantInterfaceMethodRefInfo{ConstantMemberRefInfo{cp: cp}}
	case CONSTANT_NameAndType:
		return &ConstantNameAndTypeInfo{}
	default:
		// TODO: MVP Phase ignore: MethodHandle, MethodType, InvokeDynamic
		return nil
	}
}

// ConstantUtf8Info UTF8 string constants
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

func (c *ConstantUtf8Info) Tag() ConstantTag {
	return CONSTANT_Utf8
}

// ConstantIntegerInfo Integer (32bit) constants
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

func (c *ConstantIntegerInfo) Tag() ConstantTag {
	return CONSTANT_Integer
}

func (c *ConstantIntegerInfo) Value() int32 {
	return c.val
}

// ConstantFloatInfo Float (32bit) constants
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

func (c *ConstantFloatInfo) Tag() ConstantTag {
	return CONSTANT_Float
}

func (c *ConstantFloatInfo) Value() float32 {
	return c.val
}

// ConstantLongInfo Long (64bit) constants
type ConstantLongInfo struct {
	val int64
}

func (c *ConstantLongInfo) readInfo(reader *ClassReader) {
	c.val = int64(reader.readU8())
}

func (c *ConstantLongInfo) String() string {
	return fmt.Sprintf("%d", c.val)
}

func (c *ConstantLongInfo) Tag() ConstantTag {
	return CONSTANT_Long
}

func (c *ConstantLongInfo) Value() int64 {
	return c.val
}

// ConstantDoubleInfo (64bit) constants
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

func (c *ConstantDoubleInfo) Tag() ConstantTag {
	return CONSTANT_Double
}

func (c *ConstantDoubleInfo) Value() float64 {
	return c.val
}

// ConstantClassInfo Class Constant
type ConstantClassInfo struct {
	cp        ClassFileConstantPool
	nameIndex uint16 // point to UTF8 constants
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

func (c *ConstantClassInfo) Tag() ConstantTag {
	return CONSTANT_Class
}

// ConstantStringInfo String constants
type ConstantStringInfo struct {
	cp          ClassFileConstantPool
	stringIndex uint16
}

func (c *ConstantStringInfo) readInfo(reader *ClassReader) {
	c.stringIndex = reader.readU2()
}

func (c *ConstantStringInfo) String() string {
	return getUtf8(c.cp, c.stringIndex)
}

func (c *ConstantStringInfo) Tag() ConstantTag {
	return CONSTANT_String
}

// ConstantMemberRefInfo Member Ref constants (fields and methods)
type ConstantMemberRefInfo struct {
	cp               ClassFileConstantPool
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (c *ConstantMemberRefInfo) ClassName() string {
	classInfo := c.cp[c.classIndex].(*ConstantClassInfo)
	return getUtf8(c.cp, classInfo.nameIndex)
}

func (c *ConstantMemberRefInfo) readInfo(reader *ClassReader) {
	c.classIndex = reader.readU2()
	c.nameAndTypeIndex = reader.readU2()
}

func (c *ConstantMemberRefInfo) ClassIndex() uint16 {
	return c.classIndex
}

func (c *ConstantMemberRefInfo) NameAndTypeIndex() uint16 {
	return c.nameAndTypeIndex
}

func (c *ConstantMemberRefInfo) NameAndDescriptor() (string, string) {
	con := c.cp[c.nameAndTypeIndex].(*ConstantNameAndTypeInfo)
	name := getUtf8(c.cp, con.nameIndex)
	desc := getUtf8(c.cp, con.descriptorIndex)
	return name, desc
}

// ConstantFieldRefInfo field Ref constants
type ConstantFieldRefInfo struct {
	ConstantMemberRefInfo
}

func (c *ConstantFieldRefInfo) String() string {
	return fmt.Sprintf("FieldRef: class=%d, nameAndType=%d", c.classIndex, c.nameAndTypeIndex)
}

func (c *ConstantFieldRefInfo) Tag() ConstantTag {
	return CONSTANT_Fieldref
}

// ConstantMethodRefInfo method ref constants
type ConstantMethodRefInfo struct {
	ConstantMemberRefInfo
}

func (c *ConstantMethodRefInfo) String() string {
	return fmt.Sprintf("MethodRef: class=%d, nameAndType=%d", c.classIndex, c.nameAndTypeIndex)
}

func (c *ConstantMethodRefInfo) Tag() ConstantTag {
	return CONSTANT_Methodref
}

// ConstantInterfaceMethodRefInfo interface method ref constants
type ConstantInterfaceMethodRefInfo struct {
	ConstantMemberRefInfo
}

func (c *ConstantInterfaceMethodRefInfo) String() string {
	return fmt.Sprintf("InterfaceMethodRef: class=%d, nameAndType=%d", c.classIndex, c.nameAndTypeIndex)
}

func (c *ConstantInterfaceMethodRefInfo) Tag() ConstantTag {
	return CONSTANT_InterfaceMethodref
}

// ConstantNameAndTypeInfo name and type
type ConstantNameAndTypeInfo struct {
	nameIndex       uint16
	descriptorIndex uint16
}

func (c *ConstantNameAndTypeInfo) readInfo(reader *ClassReader) {
	c.nameIndex = reader.readU2()
	c.descriptorIndex = reader.readU2()
}

func (c *ConstantNameAndTypeInfo) String() string {
	return fmt.Sprintf("NameAndType: name=%d, descriptor=%d", c.nameIndex, c.descriptorIndex)
}

func (c *ConstantNameAndTypeInfo) Tag() ConstantTag {
	return CONSTANT_NameAndType
}

func (c *ConstantNameAndTypeInfo) NameIndex() uint16 {
	return c.nameIndex
}

func (c *ConstantNameAndTypeInfo) DescriptorIndex() uint16 {
	return c.descriptorIndex
}
