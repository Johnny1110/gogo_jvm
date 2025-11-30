package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
)

// Class instant in runtime method area
//
// ClassFile（Compile）   →    Class（Runtime）
// ┌─────────────────┐        ┌─────────────────┐
// │ constantPool    │   →    │ constantPool    │  RuntimeConstantPool
// │ accessFlags     │   →    │ accessFlags     │
// │ thisClass       │   →    │ name            │  class name
// │ superClass      │   →    │ superClass      │  ref to super Class
// │ interfaces      │   →    │ interfaces      │  ref to interface Class[]
// │ fields          │   →    │ fields          │  runtime fields
// │ methods         │   →    │ methods         │  runtime methods
// └─────────────────┘        └─────────────────┘
type Class struct {
	accessFlags       uint16
	name              string // className (ex: java/lang/Object)
	superClassName    string
	interfaceNames    []string
	constantPool      *RuntimeConstantPool // ConstantPool Runtime
	fields            []*Field
	methods           []*Method
	loader            *ClassLoader // 加載此類的 ClassLoader
	superClass        *Class       // parent class ref
	interfaces        []*Class     // interface refs
	instanceSlotCount uint         // 實例變量佔用的 slot 數量
	staticSlotCount   uint         // 類變量佔用的 slot 數量
	instanceVars      rtcore.Slots // class's non-static vars
	staticVars        rtcore.Slots // class's static vars
}

// newClass create Class from classfile.ClassFile
func newClass(cf *classfile.ClassFile) *Class {
	c := &Class{}
	c.accessFlags = cf.AccessFlags()
	c.name = cf.ClassName()
	c.superClassName = cf.SuperClassName()
	c.interfaceNames = cf.InterfaceNames()
	c.constantPool = newRuntimeConstantPool(c, cf.ConstantPool())
	c.fields = newFields(c, cf.Fields())
	c.methods = newMethods(c, cf.Methods())
	return c
}

// =============== Getters ===============

func (c *Class) Name() string                       { return c.name }
func (c *Class) SuperClassName() string             { return c.superClassName }
func (c *Class) InterfaceNames() []string           { return c.interfaceNames }
func (c *Class) ConstantPool() *RuntimeConstantPool { return c.constantPool }
func (c *Class) Fields() []*Field                   { return c.fields }
func (c *Class) Methods() []*Method                 { return c.methods }
func (c *Class) Loader() *ClassLoader               { return c.loader }
func (c *Class) SuperClass() *Class                 { return c.superClass }
func (c *Class) StaticVars() rtcore.Slots           { return c.staticVars }
func (c *Class) AccessFlags() uint16                { return c.accessFlags }
func (c *Class) InstanceSlotCount() uint            { return c.instanceSlotCount }

// =============== Access Flags ===============

func (c *Class) IsPublic() bool     { return c.accessFlags&common.ACC_PUBLIC != 0 }
func (c *Class) IsFinal() bool      { return c.accessFlags&common.ACC_FINAL != 0 }
func (c *Class) IsSuper() bool      { return c.accessFlags&common.ACC_SUPER != 0 }
func (c *Class) IsInterface() bool  { return c.accessFlags&common.ACC_INTERFACE != 0 }
func (c *Class) IsAbstract() bool   { return c.accessFlags&common.ACC_ABSTRACT != 0 }
func (c *Class) IsSynthetic() bool  { return c.accessFlags&common.ACC_SYNTHETIC != 0 }
func (c *Class) IsAnnotation() bool { return c.accessFlags&common.ACC_ANNOTATION != 0 }
func (c *Class) IsEnum() bool       { return c.accessFlags&common.ACC_ENUM != 0 }

// =============== Method Lookup ===============

// GetMainMethod get `public static void main([]String args) {}`
func (c *Class) GetMainMethod() *Method {
	return c.getStaticMethod("main", "([Ljava/lang/String;)V")
}

// getStaticMethod get static method
func (c *Class) getStaticMethod(name, descriptor string) *Method {
	for _, method := range c.methods {
		if method.IsStatic() && method.name == name && method.descriptor == descriptor {
			return method
		}
	}
	return nil
}

// GetMethod get method (including parent and grandparent...)
func (c *Class) GetMethod(name, descriptor string) *Method {
	// try to find in current class
	for _, method := range c.methods {
		if method.name == name && method.descriptor == descriptor {
			return method
		}
	}

	// try to find in parent and grandparent...
	if c.superClass != nil {
		return c.superClass.GetMethod(name, descriptor)
	}

	return nil
}

func (c *Class) GetStaticMethod(name, descriptor string) *Method {
	return c.getStaticMethod(name, descriptor)
}
