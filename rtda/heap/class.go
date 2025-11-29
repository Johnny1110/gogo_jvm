package heap

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

type Class struct {
	instanceSlotCount uint
	methods           []*Method
	constantPool      *RuntimeConstantPool
	name              string
	superClassName    string
	superClass        *Class
	loader            *ClassLoader
	interfaceNames    []string
	interfaces        []*Class
	fields            []*Field
	staticSlotCount   uint
	staticVars        Slots
}

func newClass(cf *classfile.ClassFile) *Class {
	class := &Class{}
	return class
}

func (c *Class) IsInterface() bool {
	return false
}

func (c *Class) InstanceSlotCount() uint {
	return c.instanceSlotCount
}

func (c *Class) Methods() []*Method {
	return c.methods
}

func (c *Class) ConstantPool() *RuntimeConstantPool {
	return c.constantPool
}

func (c *Class) GetMainMethod() *Method {
	return nil
}
