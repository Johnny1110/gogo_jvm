package method_area

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
	"io/ioutil"
)

type ClassLoader struct {
	classPath string            // .class files dir
	classMap  map[string]*Class // loaded classes（Method Area）- key: className
}

// NewClassLoader create class loader
func NewClassLoader(classPath string) *ClassLoader {
	return &ClassLoader{
		classPath: classPath,
		classMap:  make(map[string]*Class),
	}
}

// LoadClass load class
// name: class's full name, like "java/lang/Object" or "Calculator"
func (loader *ClassLoader) LoadClass(name string, debug bool) *Class {
	// 1. check is loaded or not, return cache.
	if class, ok := loader.classMap[name]; ok {
		return class
	}

	// 2. load class
	return loader.loadNonArrayClass(name, debug)
}

// loadNonArrayClass load non array class
func (loader *ClassLoader) loadNonArrayClass(name string, debug bool) *Class {
	// 1. read .class
	classBytecode, err := loader.readClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}

	// 2. parse ClassFile bytecode to class object
	class := loader.defineClass(classBytecode, debug)

	// 3. do link（Verification and Preparation）
	link(class)

	fmt.Printf("[ClassLoader] Loaded: %s\n", name)
	return class
}

// readClass read .class file
func (loader *ClassLoader) readClass(name string) ([]byte, error) {
	// try read from classpath
	// classname format: java/lang/Object → java/lang/Object.class
	filePath := loader.classPath + "/" + name + ".class"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		// try find at current path
		filePath = loader.classPath + "/" + name + ".class"
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			// try with class name
			data, err = ioutil.ReadFile(name + ".class")
		}
	}
	return data, err
}

// defineClass define class (create Class from bytecode)
func (loader *ClassLoader) defineClass(data []byte, debug bool) *Class {
	// 1. parse ClassFile
	cf, err := classfile.Parse(data)
	if err != nil {
		panic("java.lang.ClassFormatError: " + err.Error())
	}

	if debug {
		classfile.Debug(cf, true)
	}

	// 2. convert classfile to Class object
	class := newClass(cf)
	// make class linked with this ClassLoader
	class.loader = loader

	// 3. load super class
	loader.resolveSuperClass(class)
	// 4. load interface
	loader.resolveInterfaces(class)

	// after step 3 and 4, all interfaces and parent, grandparent will be loaded into this ClassLoader

	// 5. store into Method Area (this class will never be load again)
	loader.classMap[class.name] = class

	return class
}

// resolveSuperClass load super class
func (loader *ClassLoader) resolveSuperClass(class *Class) {
	// TODO: MVP phase -> skip this (currently we don't need extends)
	//if class.name != "java/lang/Object" && class.superClassName != "" {
	//	// recursive load parent class
	//	class.superClass = class.loader.LoadClass(class.superClassName)
	//}
}

// resolveInterfaces load interface
func (loader *ClassLoader) resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, ifaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(ifaceName, false)
		}
	}
}

// link
func link(class *Class) {
	// 1. Verification - TODO: skip
	verify(class)

	// 2.Preparation - allocate space for static const
	prepare(class)
}

// verify Verification check bytecode, prevent jvm from crash
func verify(class *Class) {
	// TODO: 字節碼驗證 MVP 階段跳過
}

// prepare Preparation
// allocate space for static const
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)
	calcStaticFieldSlotIds(class)
	allocAndInitVars(class)
}

// calcInstanceFieldSlotIds calculate instance fields slot ID
func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	// from parent
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}

	class.instanceSlotCount = slotId
}

// calcStaticFieldSlotIds calculate static fields slot ID
func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}

// allocAndInitVars allocate and init static & non-static vars
func allocAndInitVars(class *Class) {
	class.instanceVars = rtcore.NewSlots(class.instanceSlotCount)
	class.staticVars = rtcore.NewSlots(class.staticSlotCount)
	// TODO: 初始化 static final 常量
}
