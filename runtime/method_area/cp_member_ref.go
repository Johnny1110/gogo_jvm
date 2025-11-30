package method_area

import "github.com/Johnny1110/gogo_jvm/classfile"

// MemberRef base ref for all FieldRef and MethodRef
type MemberRef struct {
	SymRef
	name       string
	descriptor string
}

// copyMemberRefInfo copy into from ClassFile
func (r *MemberRef) copyMemberRefInfo(refInfo *classfile.ConstantMemberRefInfo) {
	r.className = refInfo.ClassName()
	r.name, r.descriptor = refInfo.NameAndDescriptor()
}

// Name getter
func (r *MemberRef) Name() string {
	return r.name
}

// Descriptor getter
func (r *MemberRef) Descriptor() string {
	return r.descriptor
}
