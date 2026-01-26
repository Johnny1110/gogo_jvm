package classfile

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/common"
)

const HOLY_MAGIC = 0xCAFEBABE

// ClassFile represent Java class file format (Java Standard)
// every .class file should follow this
type ClassFile struct {
	magic        uint32 // magic number: 0xCAFEBABE, for classify .class file (4 bytes)
	minorVersion uint16
	majorVersion uint16
	constantPool ClassFileConstantPool // constants pool
	accessFlags  uint16                // class access flags
	thisClass    uint16                // this class index (pointing to constantPool)
	superClass   uint16                // super class index
	interfaces   []uint16              // implemented interfaces index
	fields       []*MemberInfo         // fields table
	methods      []*MemberInfo         // methods table
	attributes   []AttributeInfo       // attributes table
}

// Parse parse class file
func Parse(classData []byte) (cf *ClassFile, err error) {
	defer func() { // prevent panic error happened, convert to error.
		if r := recover(); r != nil {
			err = fmt.Errorf("parse class file error: %v", r)
		}
	}()

	cr := &ClassReader{classData}
	cf = &ClassFile{}
	cf.read(cr)
	return
}

// fill the ClassFile by input ClassReader (must be ordered like this)
func (cf *ClassFile) read(reader *ClassReader) {
	cf.readAndCheckMagic(reader)
	cf.readAndCheckVersion(reader)
	cf.constantPool = readConstantPool(reader)
	cf.accessFlags = reader.readU2()
	cf.thisClass = reader.readU2()
	cf.superClass = reader.readU2()
	cf.interfaces = reader.readU2Table()
	cf.fields = readMembers(reader, cf.constantPool)
	cf.methods = readMembers(reader, cf.constantPool)
	cf.attributes = readAttributes(reader, cf.constantPool)
}

// readAndCheckMagic read and check magic number
// must be 0xCAFEBABE -> James Gosling favourite coffee shop
func (cf *ClassFile) readAndCheckMagic(reader *ClassReader) {
	fileMagic := reader.readU4() // magic number take 4 bytes
	if fileMagic != HOLY_MAGIC {
		panic("java.lang.ClassFormatError: Invalid magic number")
	}
	cf.magic = fileMagic
}

// readAndCheckVersion read and check version number
// different JVM version support different features, high version class file can not run on low version JVM.
func (cf *ClassFile) readAndCheckVersion(reader *ClassReader) {
	cf.minorVersion = reader.readU2()
	cf.majorVersion = reader.readU2()

	// gogo-jvm support from Java 1.0 ~ Java 8 (version no 45-52)
	// java 8 is 52.0, java 7 is 51.0 ...
	switch cf.majorVersion {
	case 45, 46, 47, 48, 49, 50, 51, 52, 53:
		return
	}

	panic(fmt.Sprintf("java.lang.UnsupportedClassVersionError: %d.%d",
		cf.majorVersion, cf.minorVersion))
}

// ====================================================================

func (cf *ClassFile) ClassName() string {
	return cf.constantPool.getClassName(cf.thisClass)
}

func (cf *ClassFile) SuperClassName() string {
	if cf.superClass == 0 {
		// java.lang.Object has no parent (super)
		return ""
	}

	return cf.constantPool.getClassName(cf.superClass)
}

func (cf *ClassFile) InterfaceNames() []string {
	interfaceNames := make([]string, len(cf.interfaces))
	for i, interfaceIdx := range cf.interfaces {
		interfaceNames[i] = cf.constantPool.getClassName(interfaceIdx)
	}
	return interfaceNames
}

func (cf *ClassFile) Fields() []*MemberInfo {
	return cf.fields
}

func (cf *ClassFile) Methods() []*MemberInfo {
	return cf.methods
}

func (cf *ClassFile) ConstantPool() ClassFileConstantPool {
	return cf.constantPool
}

func (cf *ClassFile) AccessFlags() uint16 {
	return cf.accessFlags
}

func (cf *ClassFile) IsPublic() bool {
	return cf.accessFlags&common.ACC_PUBLIC != 0
}

func (cf *ClassFile) IsFinal() bool {
	return cf.accessFlags&common.ACC_FINAL != 0
}

func (cf *ClassFile) IsInterface() bool {
	return cf.accessFlags&common.INTERFACE != 0
}

func (cf *ClassFile) IsAbstract() bool {
	return cf.accessFlags&common.ACC_ABSTRACT != 0
}

// SourceFileAttribute
func (cf *ClassFile) SourceFileAttribute() *SourceFileAttribute {
	for _, attr := range cf.attributes {
		if sfAttr, ok := attr.(*SourceFileAttribute); ok {
			return sfAttr
		}
	}
	return nil
}

// MajorVersion 主版本號
func (cf *ClassFile) MajorVersion() uint16 {
	return cf.majorVersion
}

// MinorVersion 次版本號
func (cf *ClassFile) MinorVersion() uint16 {
	return cf.minorVersion
}

// GetMainMethod get main function
// main func sign must be `public static void main(String[] args)`
func (cf *ClassFile) GetMainMethod() *MemberInfo {
	for _, method := range cf.methods {
		if method.IsPSVM() {
			return method
		}
	}
	return nil
}

// for debug usage
func (cf *ClassFile) String() string {
	return fmt.Sprintf(`ClassFile {
	Magic: %X
	Version: %d.%d
	ClassFileConstantPool: %d items
	AccessFlags: 0x%04X
	ThisClass: %s
	SuperClass: %s
	Interfaces: %v
	Fields: %d
	Methods: %d
	Attributes: %d
}`,
		cf.magic,
		cf.majorVersion, cf.minorVersion,
		len(cf.constantPool),
		cf.accessFlags,
		cf.ClassName(),
		cf.SuperClassName(),
		cf.InterfaceNames(),
		len(cf.fields),
		len(cf.methods),
		len(cf.attributes))
}
