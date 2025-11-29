package java

type Method struct {
	maxLocals    uint16
	maxStack     uint16
	argSlotCount uint16
	code         []byte
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
