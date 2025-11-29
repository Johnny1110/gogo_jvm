package classfile

// MemberInfo represent field and method.
// in a class, field and method java are the same.
// - access_flags     : access flag
// - name_index       : index of name in constants pool
// - descriptor_index : descriptor in constants pool
// - attributes       : attributes list
type MemberInfo struct {
	cp              ConstantPool
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []AttributeInfo // could contains: Code, ConstantValue, Exceptions, SourceFile, LineNumberTable, LocalVariableTable
}

func readMembers(reader *ClassReader, cp ConstantPool) []*MemberInfo {
	memberCount := reader.readU2()
	members := make([]*MemberInfo, memberCount)
	for i := range members {
		members[i] = readMember(reader, cp)
	}
	return members
}

func readMember(reader *ClassReader, cp ConstantPool) *MemberInfo {
	return &MemberInfo{
		cp:              cp,
		accessFlags:     reader.readU2(),
		nameIndex:       reader.readU2(),
		descriptorIndex: reader.readU2(),
		attributes:      readAttributes(reader, cp),
	}
}

// Name get member name
func (m *MemberInfo) Name() string {
	return getUtf8(m.cp, m.nameIndex)
}

func (m *MemberInfo) Descriptor() string {
	return getUtf8(m.cp, m.descriptorIndex)
}

func (m *MemberInfo) AccessFlags() uint16 {
	return m.accessFlags
}

func (m *MemberInfo) IsPublic() bool {
	return m.accessFlags&ACC_PUBLIC != 0
}

func (m *MemberInfo) IsPrivate() bool {
	return m.accessFlags&ACC_PRIVATE != 0
}

func (m *MemberInfo) IsProtected() bool {
	return m.accessFlags&ACC_PROTECTED != 0
}

func (m *MemberInfo) IsStatic() bool {
	return m.accessFlags&ACC_STATIC != 0
}

func (m *MemberInfo) IsFinal() bool {
	return m.accessFlags&ACC_FINAL != 0
}

func (m *MemberInfo) IsSynchronized() bool {
	return m.accessFlags&ACC_SYNCHRONIZED != 0
}

func (m *MemberInfo) IsVolatile() bool {
	return m.accessFlags&ACC_VOLATILE != 0
}

func (m *MemberInfo) IsTransient() bool {
	return m.accessFlags&ACC_TRANSIENT != 0
}

func (m *MemberInfo) IsNative() bool {
	return m.accessFlags&ACC_NATIVE != 0
}

func (m *MemberInfo) IsAbstract() bool {
	return m.accessFlags&ACC_ABSTRACT != 0
}

func (m *MemberInfo) IsSynthetic() bool {
	return m.accessFlags&ACC_SYNTHETIC != 0
}

func (m *MemberInfo) IsEnum() bool {
	return m.accessFlags&ACC_ENUM != 0
}

// Get Code attribute (method's bytecode)
func (m *MemberInfo) CodeAttribute() *CodeAttribute {
	for _, attr := range m.attributes {
		if codeAttr, ok := attr.(*CodeAttribute); ok {
			return codeAttr
		}
	}
	return nil
}

// Get ConstantValue attribute (static final field value)
func (m *MemberInfo) ConstantValueAttribute() *ConstantValueAttribute {
	for _, attr := range m.attributes {
		if cvAttr, ok := attr.(*ConstantValueAttribute); ok {
			return cvAttr
		}
	}
	return nil
}

// IsPSVM is public static void main([]String args) {...}
func (m *MemberInfo) IsPSVM() bool {
	return m.Name() == "main" && m.Descriptor() == "([Ljava/lang/String;)V" &&
		m.IsStatic() && m.IsPublic()
}
