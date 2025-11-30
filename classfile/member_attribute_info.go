package classfile

// AttributeInfo attribute interface
// Attribute is a extendable java in a class, using for store extra info
// like: method's bytecode, SourceFileName...
// AttributeInfo attribute interface
type AttributeInfo interface{}

func readAttributes(reader *ClassReader, cp ClassFileConstantPool) []AttributeInfo {
	attributesCount := reader.readU2()
	attributes := make([]AttributeInfo, attributesCount)
	for i := range attributes {
		attributes[i] = readAttribute(reader, cp)
	}
	return attributes
}

func readAttribute(reader *ClassReader, cp ClassFileConstantPool) AttributeInfo {
	attrNameIndex := reader.readU2()
	attrName := getUtf8(cp, attrNameIndex)
	attrLength := reader.readU4()

	attrInfo := newAttributeInfo(attrName, attrLength, cp)
	if attrInfo == nil {
		attrInfo = &UnparsedAttribute{
			name:   attrName,
			length: attrLength,
			info:   reader.readBytes(attrLength),
		}
	} else {
		attrInfo.(interface {
			readInfo(*ClassReader)
		}).readInfo(reader)
	}

	return attrInfo
}

func newAttributeInfo(attrName string, attrLength uint32, cp ClassFileConstantPool) AttributeInfo {
	switch attrName {
	case "Code":
		return &CodeAttribute{cp: cp}
	case "ConstantValue":
		return &ConstantValueAttribute{}
	case "Exceptions":
		return &ExceptionsAttribute{}
	case "SourceFile":
		return &SourceFileAttribute{cp: cp}
	case "LineNumberTable":
		return &LineNumberTableAttribute{}
	case "LocalVariableTable":
		return &LocalVariableTableAttribute{}
	default:
		return nil
	}
}

// UnparsedAttribute
type UnparsedAttribute struct {
	name   string
	length uint32
	info   []byte
}

// CodeAttribute
type CodeAttribute struct {
	cp             ClassFileConstantPool
	maxStack       uint16
	maxLocals      uint16
	code           []byte
	exceptionTable []*ExceptionHandler
	attributes     []AttributeInfo
}

func (c *CodeAttribute) readInfo(reader *ClassReader) {
	c.maxStack = reader.readU2()
	c.maxLocals = reader.readU2()
	codeLength := reader.readU4()
	c.code = reader.readBytes(codeLength)
	c.exceptionTable = readExceptionTable(reader)
	c.attributes = readAttributes(reader, c.cp)
}

func (c *CodeAttribute) MaxStack() uint16  { return c.maxStack }
func (c *CodeAttribute) MaxLocals() uint16 { return c.maxLocals }
func (c *CodeAttribute) Code() []byte      { return c.code }

type ExceptionHandler struct {
	startPc   uint16
	endPc     uint16
	handlerPc uint16
	catchType uint16
}

func readExceptionTable(reader *ClassReader) []*ExceptionHandler {
	exceptionTableLength := reader.readU2()
	exceptionTable := make([]*ExceptionHandler, exceptionTableLength)
	for i := range exceptionTable {
		exceptionTable[i] = &ExceptionHandler{
			startPc:   reader.readU2(),
			endPc:     reader.readU2(),
			handlerPc: reader.readU2(),
			catchType: reader.readU2(),
		}
	}
	return exceptionTable
}

// ConstantValueAttribute
type ConstantValueAttribute struct {
	constantValueIndex uint16
}

func (c *ConstantValueAttribute) readInfo(reader *ClassReader) {
	c.constantValueIndex = reader.readU2()
}

func (c *ConstantValueAttribute) ConstantValueIndex() uint16 {
	return c.constantValueIndex
}

// ExceptionsAttribute
type ExceptionsAttribute struct {
	exceptionsIndexTable []uint16
}

func (e *ExceptionsAttribute) readInfo(reader *ClassReader) {
	e.exceptionsIndexTable = reader.readU2Table()
}

// SourceFileAttribute
type SourceFileAttribute struct {
	cp              ClassFileConstantPool
	sourceFileIndex uint16
}

func (s *SourceFileAttribute) readInfo(reader *ClassReader) {
	s.sourceFileIndex = reader.readU2()
}

func (s *SourceFileAttribute) FileName() string {
	return getUtf8(s.cp, s.sourceFileIndex)
}

// LineNumberTableAttribute
type LineNumberTableAttribute struct {
	lineNumberTable []*LineNumberTableEntry
}

type LineNumberTableEntry struct {
	startPc    uint16
	lineNumber uint16
}

func (l *LineNumberTableAttribute) readInfo(reader *ClassReader) {
	lineNumberTableLength := reader.readU2()
	l.lineNumberTable = make([]*LineNumberTableEntry, lineNumberTableLength)
	for i := range l.lineNumberTable {
		l.lineNumberTable[i] = &LineNumberTableEntry{
			startPc:    reader.readU2(),
			lineNumber: reader.readU2(),
		}
	}
}

// LocalVariableTableAttribute
type LocalVariableTableAttribute struct {
	localVariableTable []*LocalVariableTableEntry
}

type LocalVariableTableEntry struct {
	startPc         uint16
	length          uint16
	nameIndex       uint16
	descriptorIndex uint16
	index           uint16
}

func (l *LocalVariableTableAttribute) readInfo(reader *ClassReader) {
	localVarTableLength := reader.readU2()
	l.localVariableTable = make([]*LocalVariableTableEntry, localVarTableLength)
	for i := range l.localVariableTable {
		l.localVariableTable[i] = &LocalVariableTableEntry{
			startPc:         reader.readU2(),
			length:          reader.readU2(),
			nameIndex:       reader.readU2(),
			descriptorIndex: reader.readU2(),
			index:           reader.readU2(),
		}
	}
}
