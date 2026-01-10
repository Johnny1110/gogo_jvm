package method_area

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/classfile"
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
	"io/ioutil"
)

type ClassLoader struct {
	classPath string            // .class files dir
	classMap  map[string]*Class // loaded classes（Method Area）- key: className

	// v0.3.1: Reflection Support - 解決雞生蛋問題
	// lClassClass is "java/lang/Class" 的 Class (metadata)
	// for create java.lang.Class tp all Objects
	jlClassClass *Class
	// for primitive type, key: "void", "int", "long", etc.
	primitiveClasses map[string]*Class
}

// NewClassLoader create class loader
func NewClassLoader(classPath string) *ClassLoader {
	loader := &ClassLoader{
		classPath:        classPath,
		classMap:         make(map[string]*Class),
		primitiveClasses: make(map[string]*Class),
	}
	// v0.3.1: init reflection system
	loader.initReflection()

	return loader
}

// ============================================================
// v0.3.1: Reflection Initialization (Bootstrap)
// ============================================================
// initReflection init reflection system
// problems:
// - creating String.class Object, need java.lang.Class Object
// - BUT java.lang.Class is a Class Object either
// solutions:
// - load all basic type without jClass
// - inject required jClass
func (loader *ClassLoader) initReflection() {
	fmt.Println("@@ Debug - [ClassLoader] Initializing reflection system...")

	// ---------------------------------------------------------------------------
	// 1: load "java/lang/Class" but not create jClass (Object), jlClassClass.jClass = nil
	jlClassClass := loader.loadBasicClass("java/lang/Class")
	loader.jlClassClass = jlClassClass
	// ---------------------------------------------------------------------------
	// 2: create jClass for "java/lang/Class"
	jlClassClass.jClass = loader.createJClassObject(jlClassClass)
	// ---------------------------------------------------------------------------
	// 3: create jClass for other loaded classes
	for _, class := range loader.classMap {
		if class.jClass == nil {
			class.jClass = loader.createJClassObject(class)
		}
	}
	// ---------------------------------------------------------------------------
	// 4: init primitive classes
	loader.initPrimitiveClasses()
	// ---------------------------------------------------------------------------
	fmt.Println("@@ Debug - [ClassLoader] Reflection system initialized.")
}

// loadBasicClass load basic class not create  jClass
func (loader *ClassLoader) loadBasicClass(name string) *Class {
	// if loaded just return
	if class, ok := loader.classMap[name]; ok {
		return class
	}

	// read class bytecode
	classBytecode, err := loader.readClass(name)
	if err != nil {
		fmt.Printf("read class %s error: %v \n", name, err)
		panic("java.lang.ClassNotFoundException: " + name)
	}

	// parse classfile
	cf, err := classfile.Parse(classBytecode)
	if err != nil {
		fmt.Printf("parse class %s error: %v \n", name, err)
		panic("java.lang.ClassFormatError: " + err.Error())
	}

	// create class without jClass
	class := newClass(cf)
	class.loader = loader

	// store into classMap
	loader.classMap[name] = class

	// load super
	loader.resolveSuperClass(class)
	// load interfaces
	loader.resolveInterfaces(class)
	// link（Verification and Preparation）
	link(class)

	fmt.Printf("@@ Debug - [ClassLoader] Loaded (basic): %s\n", name)
	return class
}

// initPrimitiveClasses init primitive classes
// basic type didn't have .class, we hardcode this
func (loader *ClassLoader) initPrimitiveClasses() {
	primitiveTypes := []string{
		"void", "boolean", "byte", "char", "short",
		"int", "long", "float", "double",
	}

	for _, typeName := range primitiveTypes {
		// create basic type Class struct
		primitiveClass := &Class{
			name:        typeName,
			accessFlags: common.ACC_PUBLIC, // basic type is all public
			loader:      loader,
			// superClass, methods, fields is nil
		}

		// create jClass -> java.lang.Class Object
		primitiveClass.jClass = loader.createJClassObject(primitiveClass)
		// store into ClassLoader's primitiveClasses
		loader.primitiveClasses[typeName] = primitiveClass

		fmt.Printf("@@ Debug - [ClassLoader] Created primitive class: %s\n", typeName)
	}
}

// createJClassObject create java.lang.Class (Object) for target class
func (loader *ClassLoader) createJClassObject(class *Class) *heap.Object {
	jlClassClass := loader.jlClassClass // this is java.lang.Class's *class

	// create "Class" Obj
	classObj := &heap.Object{}
	classObj.SetMarkWord(heap.InitialMarkWord)

	// ==============================================================
	// !!! 下面這兩個不要混淆了，SetClass 是這個 class 的類型，extra 才是實際 mirror class
	// set up type（java.lang.Class）
	classObj.SetClass(jlClassClass)
	// Class Object's extra is method_area.Class (Mirror)
	classObj.SetExtra(class)
	// ==============================================================

	// TODO: 如果 java.lang.Class 有實例欄位，需要分配空間。目前不處理 java.lang.Class 的欄位
	if jlClassClass != nil && jlClassClass.instanceSlotCount > 0 {
		classObj.SetFields(rtcore.NewSlots(jlClassClass.instanceSlotCount))
	}

	return classObj
}

// GetPrimitiveClass get basic type's Class
func (loader *ClassLoader) GetPrimitiveClass(name string) *Class {
	return loader.primitiveClasses[name]
}

// ============================================================

// LoadClass load class
// name: class's full name, like "java/lang/Object" or "Calculator"
func (loader *ClassLoader) LoadClass(name string, debug bool) *Class {
	// 1. check primitiveClass first (v0.3.1)
	if primitiveClass, ok := loader.primitiveClasses[name]; ok {
		return primitiveClass
	}

	// 2. check is loaded or not, return cache.
	if class, ok := loader.classMap[name]; ok {
		return class
	}

	// check array type (v0.3.1)
	if len(name) > 0 && name[0] == '[' {
		// 3. load array class
		return loader.loadArrayClass(name)
	} else {
		// 4. load normal class
		return loader.loadNonArrayClass(name, debug)
	}
}

func (loader *ClassLoader) LoadClassIface(name string) interface{} {
	return loader.LoadClass(name, false)
}

// loadArrayClass load array class
// array class is dynamic generate, no need .class file
func (loader *ClassLoader) loadArrayClass(name string) *Class {
	fmt.Printf("@@ Debug - [ClassLoader] Loading array class: %s\n", name)

	arrayClass := &Class{
		name:        name,
		accessFlags: common.ACC_PUBLIC, // array class is public
		loader:      loader,
		superClass:  loader.LoadClass("java/lang/Object", false),
		interfaces: []*Class{
			loader.LoadClass("java/lang/Cloneable", false),
			loader.LoadClass("java/io/Serializable", false),
		},
		initStarted: true, // array no need init
	}

	// parse elements
	componentClassName := GetComponentClassName(name)
	if componentClassName != "" {
		if isPrimitiveTypeName(componentClassName) {
			// primitive type
			arrayClass.componentClass = loader.GetPrimitiveClass(componentClassName)
		} else {
			// normal class type
			arrayClass.componentClass = loader.LoadClass(componentClassName, false)
		}
	}

	// create jClass for arrayClass
	arrayClass.jClass = loader.createJClassObject(arrayClass)
	// cache
	loader.classMap[name] = arrayClass

	fmt.Printf("@@ Debug - [ClassLoader] Loaded array class: %s (component: %s)\n",
		name, componentClassName)

	return arrayClass
}

// isPrimitiveTypeName check is primitive
func isPrimitiveTypeName(name string) bool {
	switch name {
	case "void", "boolean", "byte", "char", "short",
		"int", "long", "float", "double":
		return true
	}
	return false
}

// loadNonArrayClass load non array class
func (loader *ClassLoader) loadNonArrayClass(name string, debug bool) *Class {
	// 1. read .class
	classBytecode, err := loader.readClass(name)
	if err != nil {
		fmt.Printf("read class %s error: %v \n", name, err)
		panic("java.lang.ClassNotFoundException: " + name)
	}

	// 2. parse ClassFile bytecode to class object
	class := loader.defineClass(classBytecode, debug)

	// 3. do link（Verification and Preparation）
	link(class)

	// 4. create java.lang.Class (v0.3.1)
	if class.jClass == nil && loader.jlClassClass != nil {
		class.jClass = loader.createJClassObject(class)
	}

	fmt.Printf("@@ Debug - [ClassLoader] Loaded: %s\n", name)
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
	if class.name != "java/lang/Object" && class.superClassName != "" {
		// recursive load parent class
		class.superClass = class.loader.LoadClass(class.superClassName, false)
	}
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
	calcInstanceFieldSlotIds(class) // for object
	calcStaticFieldSlotIds(class)   // for class
	allocAndInitStaticVars(class)   // only init static, instant will be alloc when create object
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
	// why not include parent class ? -> find tips in heap/README.md (為什麼類別的 static field slot ID 計算時不需要考慮父類別？)
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

// allocAndInitStaticVars allocate and init static-vars
func allocAndInitStaticVars(class *Class) {
	class.staticVars = rtcore.NewSlots(class.staticSlotCount)
	// TODO: 初始化 static final 常量
}
