package classfile

// AttributeInfo attribute interface
// Attribute is a extendable structure in a class, using for store extra info
// like: method's bytecode, SourceFileName...
type AttributeInfo interface {
	// every attribute should know how to read them self.
}

func readAttributes(reader *ClassReader, cp ConstantPool) []AttributeInfo {
	attributesCount := reader.readU2()
	attributes := make([]AttributeInfo, attributesCount)
	for i := range attributes {
		attributes[i] = readAttribute(reader, cp)
	}
	return attributes
}

func readAttribute(reader *ClassReader, cp ConstantPool) AttributeInfo {
	attrNameIndex := reader.readU2()
	attrName := getUtf8(cp, attrNameIndex)
	attrLength := reader.readU4()

	// create attribute instance according to attribute name.
	attrInfo := newAttributeInfo(attrName, attrLength, cp)
	if attrInfo == nil {
		// unknown attr, skip
		attrInfo = &UnparsedAttribute{
			name:   attrName,
			length: attrLength,
			info:   reader.readBytes(attrLength),
		}
	} else {
		// read attr content
		attrInfo.(interface {
			readInfo(*ClassReader)
		}).readInfo(reader)
	}

	return attrInfo
}

func newAttributeInfo(attrName string, attrLength uint32, cp ConstantPool) AttributeInfo {
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
		// TODO: MVP Phase - only implement part of it.
		return nil
	}
}

// UnparsedAttribute Unparsed Attr
type UnparsedAttribute struct {
	name   string
	length uint32
	info   []byte
}

// CodeAttribute including method's bytecode
type CodeAttribute struct {
	cp             ConstantPool
	maxStack       uint16              // max depth of stack
	maxLocals      uint16              // size of 局部變量表
	code           []byte              // bytecode
	exceptionTable []*ExceptionHandler // exception table
	attributes     []AttributeInfo     // attribute table (nested)
}

func (c *CodeAttribute) readInfo(reader *ClassReader) {
	// stack
	c.maxStack = reader.readU2()
	// local vars table
	c.maxLocals = reader.readU2()
	// bytecode
	codeLength := reader.readU4()
	c.code = reader.readBytes(codeLength)
	// exception
	c.exceptionTable = readExceptionTable(reader)
	// attr
	c.attributes = readAttributes(reader, c.cp)
}

// MaxStack return max depth of operation stack
func (c *CodeAttribute) MaxStack() uint16 {
	return c.maxStack
}

// MaxLocals return 局部變量表大小
func (c *CodeAttribute) MaxLocals() uint16 {
	return c.maxLocals
}

// Code return bytecode
func (c *CodeAttribute) Code() []byte {
	return c.code
}

type ExceptionHandler struct {
	startPc   uint16 // try { position
	endPc     uint16 // try } position
	handlerPc uint16 // catch position
	catchType uint16 // capture exception type (index of ConstantPool)
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

// ConstantValueAttribute represent constant value (static final)
type ConstantValueAttribute struct {
	constantValueIndex uint16
}

func (c *ConstantValueAttribute) readInfo(reader *ClassReader) {
	c.constantValueIndex = reader.readU2()
}

// ExceptionsAttribute represent exceptions that will be able to throw
type ExceptionsAttribute struct {
	exceptionsIndexTable []uint16
}

func (e *ExceptionsAttribute) readInfo(reader *ClassReader) {
	e.exceptionsIndexTable = reader.readU2Table()
}

// SourceFileAttribute source file name
type SourceFileAttribute struct {
	cp              ConstantPool
	sourceFileIndex uint16
}

func (s *SourceFileAttribute) readInfo(reader *ClassReader) {
	s.sourceFileIndex = reader.readU2()
}

func (s *SourceFileAttribute) FileName() string {
	return getUtf8(s.cp, s.sourceFileIndex)
}

// LineNumberTableAttribute line number table
// bytecode and line number mapping, for debug usage
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

// LocalVariableTableAttribute 局部變量表
// for debug, mark down local vars info
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
