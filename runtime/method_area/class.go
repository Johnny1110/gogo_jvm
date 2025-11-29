package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
)

// Class 運行時類結構
// 這是 ClassFile 在 Method Area 中的運行時表示
//
// ClassFile（編譯時）    →    Class（運行時）
// ┌─────────────────┐        ┌─────────────────┐
// │ constantPool    │   →    │ constantPool    │  運行時常量池
// │ accessFlags     │   →    │ accessFlags     │
// │ thisClass       │   →    │ name            │  直接存類名
// │ superClass      │   →    │ superClass      │  指向父類 Class
// │ interfaces      │   →    │ interfaces      │  指向接口 Class[]
// │ fields          │   →    │ fields          │  運行時字段
// │ methods         │   →    │ methods         │  運行時方法
// └─────────────────┘        └─────────────────┘
type Class struct {
	accessFlags       uint16
	name              string               // 類名（全限定名，如 java/lang/Object）
	superClassName    string               // 父類名
	interfaceNames    []string             // 接口名列表
	constantPool      *RuntimeConstantPool // 運行時常量池
	fields            []*Field             // 字段列表
	methods           []*Method            // 方法列表
	loader            *ClassLoader         // 加載此類的 ClassLoader
	superClass        *Class               // 父類引用（解析後）
	interfaces        []*Class             // 接口引用列表（解析後）
	instanceSlotCount uint                 // 實例變量佔用的 rtcore 數量
	staticSlotCount   uint                 // 類變量佔用的 rtcore 數量
	staticVars        rtcore.Slots         // 靜態變量（類變量）
}

// newClass 從 ClassFile 創建 Class
func newClass(cf *classfile.ClassFile) *Class {
	c := &Class{}
	c.accessFlags = cf.AccessFlags()
	c.name = cf.ClassName()
	c.superClassName = cf.SuperClassName()
	c.interfaceNames = cf.InterfaceNames()
	c.constantPool = newConstantPool(c, cf.ConstantPool())
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

// GetMainMethod 獲取 main 方法
func (c *Class) GetMainMethod() *Method {
	return c.getStaticMethod("main", "([Ljava/lang/String;)V")
}

// getStaticMethod 獲取靜態方法
func (c *Class) getStaticMethod(name, descriptor string) *Method {
	for _, method := range c.methods {
		if method.IsStatic() && method.name == name && method.descriptor == descriptor {
			return method
		}
	}
	return nil
}

// GetMethod 獲取方法（包含繼承查找，暫時簡化）
func (c *Class) GetMethod(name, descriptor string) *Method {
	// 先在當前類查找
	for _, method := range c.methods {
		if method.name == name && method.descriptor == descriptor {
			return method
		}
	}
	// TODO: 在父類查找（繼承）
	return nil
}

// GetStaticMethod 公開版本
func (c *Class) GetStaticMethod(name, descriptor string) *Method {
	return c.getStaticMethod(name, descriptor)
}

func (c *Class) InstanceSlotCount() uint {
	return c.instanceSlotCount
}
