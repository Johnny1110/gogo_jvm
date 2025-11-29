package heap

import "github.com/Johnny1110/gogo_jvm/classfile"

type Method struct {
	maxLocals    uint16
	maxStack     uint16
	argSlotCount uint16
	code         []byte
	name         string
	descriptor   string
	class        *Class
	accessFlags  uint16
}

func (m *Method) MaxLocals() uint16 {
	return m.maxLocals
}

func (m *Method) MaxStack() uint16 {
	return m.maxStack
}

func (m *Method) ArgSlotCount() uint16 {
	return m.argSlotCount
}

func (m *Method) Code() []byte {
	return m.code
}

func (m *Method) Name() string {
	return m.name
}

func (m *Method) Descriptor() string {
	return m.descriptor
}

func (m *Method) Class() *Class {
	return m.class
}

func (m *Method) IsStatic() bool {
	return m.accessFlags&classfile.ACC_STATIC != 0
}
