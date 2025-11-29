package heap

import "github.com/Johnny1110/gogo_jvm/classfile"

// MemberRef 成員引用基類（字段和方法的共同基礎）
type MemberRef struct {
	SymRef
	name       string // 成員名
	descriptor string // 描述符
}

// copyMemberRefInfo 從 ClassFile 複製成員引用信息
func (r *MemberRef) copyMemberRefInfo(refInfo *classfile.ConstantMemberRefInfo) {
	r.className = refInfo.ClassName()
	r.name, r.descriptor = refInfo.NameAndDescriptor()
}

// Name 獲取成員名
func (r *MemberRef) Name() string {
	return r.name
}

// Descriptor 獲取描述符
func (r *MemberRef) Descriptor() string {
	return r.descriptor
}
