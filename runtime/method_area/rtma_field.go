package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/common"
)

type Field struct {
	accessFlags     uint16
	name            string
	descriptor      string
	class           *Class // belongs to
	slotId          uint   // index in slot
	constValueIndex uint   // ConstantValue attributes index (for static final) could be found in class's RuntimeConstantPool
}

func (f *Field) copyAttributes(cfFiledInfo *classfile.MemberInfo) {
	if cvAttr := cfFiledInfo.ConstantValueAttribute(); cvAttr != nil {
		f.constValueIndex = uint(cvAttr.ConstantValueIndex())
	}
}

// newFields create []*Field from ClassFile's []*MemberInfo
func newFields(class *Class, cfFields []*classfile.MemberInfo) []*Field {
	fields := make([]*Field, len(cfFields))
	for i, cfField := range cfFields {
		fields[i] = newField(class, cfField)
	}
	return fields
}

func newField(class *Class, cfFieldInfo *classfile.MemberInfo) *Field {
	field := &Field{}
	field.class = class
	field.accessFlags = cfFieldInfo.AccessFlags()
	field.name = cfFieldInfo.Name()
	field.descriptor = cfFieldInfo.Descriptor()
	field.copyAttributes(cfFieldInfo)
	return field
}

// =============== Getters ===============

func (f *Field) Name() string        { return f.name }
func (f *Field) Descriptor() string  { return f.descriptor }
func (f *Field) Class() *Class       { return f.class }
func (f *Field) SlotId() uint        { return f.slotId }
func (f *Field) AccessFlags() uint16 { return f.accessFlags }

// =============== Access Flags ===============

func (f *Field) IsPublic() bool    { return f.accessFlags&common.ACC_PUBLIC != 0 }
func (f *Field) IsPrivate() bool   { return f.accessFlags&common.ACC_PRIVATE != 0 }
func (f *Field) IsProtected() bool { return f.accessFlags&common.ACC_PROTECTED != 0 }
func (f *Field) IsStatic() bool    { return f.accessFlags&common.ACC_STATIC != 0 }
func (f *Field) IsFinal() bool     { return f.accessFlags&common.ACC_FINAL != 0 }
func (f *Field) IsVolatile() bool  { return f.accessFlags&common.ACC_VOLATILE != 0 }
func (f *Field) IsTransient() bool { return f.accessFlags&common.ACC_TRANSIENT != 0 }
func (f *Field) IsSynthetic() bool { return f.accessFlags&common.ACC_SYNTHETIC != 0 }
func (f *Field) IsEnum() bool      { return f.accessFlags&common.ACC_ENUM != 0 }

// isLongOrDouble take 2 slots
func (f *Field) isLongOrDouble() bool {
	return f.descriptor == "J" || f.descriptor == "D"
}
