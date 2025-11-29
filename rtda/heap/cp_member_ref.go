package heap

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

type MemberRef struct {
	cp               *RuntimeConstantPool
	classIndex       uint
	nameAndTypeIndex uint
	name             string
	descriptor       string
	className        string
}

func (r *MemberRef) copyMemberRefInfo(c *classfile.ConstantMemberRefInfo) {
	r.classIndex = uint(c.ClassIndex())
	r.nameAndTypeIndex = uint(c.NameAndTypeIndex())
	r.name = r.cp.GetConstant(r.classIndex).(string)
	r.descriptor = r.cp.GetConstant(r.nameAndTypeIndex).(string)
	// TODO
}
