package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/common"
)

// Method in Class
//
// including:
// - method signature (name & descriptor)
// - access flags
// - bytecode
// - max Stack size
// - max LocalVars table size
// - method input params count (actually is slots count)
// - exception table (v0.2.10)
type Method struct {
	accessFlags    uint16
	name           string
	descriptor     string
	class          *Class
	maxStack       uint16
	maxLocals      uint16
	code           []byte
	argSlotCount   uint
	exceptionTable ExceptionTable // v0.2.10
}

// newMethods create from classfile
func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {
	methods := make([]*Method, len(cfMethods))
	for i, cfMethod := range cfMethods {
		methods[i] = newMethod(class, cfMethod)
	}
	return methods
}

func newMethod(class *Class, cfMethod *classfile.MemberInfo) *Method {
	method := &Method{}
	method.class = class
	method.accessFlags = cfMethod.AccessFlags()
	method.name = cfMethod.Name()
	method.descriptor = cfMethod.Descriptor()
	method.copyAttributes(cfMethod)
	method.calcArgSlotCount()
	return method
}

// copyAttributes copy from classfile.CodeAttribute
func (m *Method) copyAttributes(cfMethod *classfile.MemberInfo) {
	if codeAttr := cfMethod.CodeAttribute(); codeAttr != nil {
		m.maxStack = codeAttr.MaxStack()
		m.maxLocals = codeAttr.MaxLocals()
		m.code = codeAttr.Code()
		// v0.2.10: parse exception table:
		m.exceptionTable = newExceptionTable(codeAttr.ExceptionTable(), m.Class().ConstantPool())
	}
}

// calcArgSlotCount calculate params take slots count
// based on Descriptor
// ex: (II)V → 2 int → 2 slots
// ex: (DD)V → 2 double → 4 slots (double take 2 slots)
// ex: (JD)V → 1 long + 1 double → 4 slots
// ex: (Ljava/lang/String;I)V → 1 ref + 1 int → 2 slots
func (m *Method) calcArgSlotCount() {
	parsedDescriptor := parseMethodDescriptor(m.descriptor)
	for _, paramType := range parsedDescriptor.parameterTypes {
		m.argSlotCount++
		// long & double take 2 slots
		if paramType == "J" || paramType == "D" {
			m.argSlotCount++
		}
	}

	// non-static method first slot store this(object) reference, so we + 1
	if !m.IsStatic() {
		m.argSlotCount++
	}
}

// =============== Getters ===============

func (m *Method) Name() string        { return m.name }
func (m *Method) Descriptor() string  { return m.descriptor }
func (m *Method) Class() *Class       { return m.class }
func (m *Method) MaxStack() uint16    { return m.maxStack }
func (m *Method) MaxLocals() uint16   { return m.maxLocals }
func (m *Method) Code() []byte        { return m.code }
func (m *Method) ArgSlotCount() uint  { return m.argSlotCount }
func (m *Method) AccessFlags() uint16 { return m.accessFlags }
func (m *Method) ExceptionTable() ExceptionTable {
	return m.exceptionTable
}

// =============== Access Flags ===============

func (m *Method) IsPublic() bool       { return m.accessFlags&common.ACC_PUBLIC != 0 }
func (m *Method) IsPrivate() bool      { return m.accessFlags&common.ACC_PRIVATE != 0 }
func (m *Method) IsProtected() bool    { return m.accessFlags&common.ACC_PROTECTED != 0 }
func (m *Method) IsStatic() bool       { return m.accessFlags&common.ACC_STATIC != 0 }
func (m *Method) IsFinal() bool        { return m.accessFlags&common.ACC_FINAL != 0 }
func (m *Method) IsSynchronized() bool { return m.accessFlags&common.ACC_SYNCHRONIZED != 0 }
func (m *Method) IsBridge() bool       { return m.accessFlags&common.ACC_BRIDGE != 0 }
func (m *Method) IsVarargs() bool      { return m.accessFlags&common.ACC_VARARGS != 0 }
func (m *Method) IsNative() bool       { return m.accessFlags&common.ACC_NATIVE != 0 }
func (m *Method) IsAbstract() bool     { return m.accessFlags&common.ACC_ABSTRACT != 0 }
func (m *Method) IsStrict() bool       { return m.accessFlags&common.ACC_STRICT != 0 }

// =============== Helper ===============

// IsPSVM check is public static void main(String[] args)
func (m *Method) IsPSVM() bool {
	return m.IsPublic() && m.IsStatic() &&
		m.name == "main" && m.descriptor == "([Ljava/lang/String;)V"
}
