package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

// FieldRef
type FieldRef struct {
	MemberRef
	field *Field // resolved field (cached)
}

func NewFieldRef(cp *RuntimeConstantPool, refInfo *classfile.ConstantFieldRefInfo) *FieldRef {
	ref := &FieldRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberRefInfo)
	return ref
}

func (r *FieldRef) ResolvedField() *Field {
	if r.field == nil {
		r.resolveFieldRef()
	}
	return r.field
}

func (r *FieldRef) resolveFieldRef() {
	class := r.ResolvedClass()
	field := lookupField(class, r.name, r.descriptor)
	if field == nil {
		panic("java.lang.NoSuchFieldError: " + r.name)
	}
	r.field = field
}

func lookupField(class *Class, fieldName, fieldDescriptor string) *Field {
	for _, field := range class.Fields() {
		if field.Name() == fieldName && field.Descriptor() == fieldDescriptor {
			return field
		}
	}

	if class.superClass != nil {
		return lookupField(class.superClass, fieldName, fieldDescriptor)
	}

	return nil
}
