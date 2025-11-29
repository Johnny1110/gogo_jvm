package heap

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

type FieldRef struct {
	MemberRef
	field *Field // actual in memory address
}

func NewFieldRef(cp *RuntimeConstantPool, info *classfile.ConstantFieldRefInfo) *FieldRef {
	ref := &FieldRef{}
	ref.cp = cp
	return ref
}
