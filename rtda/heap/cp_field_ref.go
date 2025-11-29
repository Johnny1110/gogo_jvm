package heap

import "github.com/Johnny1110/gogo_jvm/classfile"

// FieldRef 字段引用
// 例如：getstatic System.out → 需要解析 out 字段
type FieldRef struct {
	MemberRef
	field *Field // 解析後的字段（緩存）
}

func newFieldRef(cp *RuntimeConstantPool, refInfo *classfile.ConstantFieldRefInfo) *FieldRef {
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
	c := r.ResolvedClass()
	field := lookupField(c, r.name, r.descriptor)
	if field == nil {
		panic("java.lang.NoSuchFieldError: " + r.name)
	}
	r.field = field
}

func lookupField(c *Class, name, descriptor string) *Field {
	for _, field := range c.fields {
		if field.name == name && field.descriptor == descriptor {
			return field
		}
	}
	// TODO: 在接口和父類中查找
	return nil
}
