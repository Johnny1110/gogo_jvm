package cpref

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/runtime/marea"
)

type MemberRef struct {
	cp               *marea.RuntimeConstantPool
	classIndex       uint16
	nameAndTypeIndex uint16
	name             string
	descriptor       string
	className        string
}

func (r *MemberRef) copyMemberRefInfo(c *classfile.ConstantMemberRefInfo) {
	r.classIndex = c.ClassIndex()
	r.nameAndTypeIndex = c.NameAndTypeIndex()

	// TODO
}
