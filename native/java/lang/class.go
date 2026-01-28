package lang

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/exception"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"strings"
)

// ============================================================
// java.lang.Class Native Methods - v0.3.1
// ============================================================
func init() {
	fmt.Println("@@ Debug - init Native java/lang/Class")

	// basic reflection
	runtime.Register("java/lang/Class", "getPrimitiveClass", "(Ljava/lang/String;)Ljava/lang/Class;", getPrimitiveClass)
	runtime.Register("java/lang/Class", "getName0", "()Ljava/lang/String;", getName0)           // old JDK
	runtime.Register("java/lang/Class", "initClassName", "()Ljava/lang/String;", initClassName) // new JDK
	runtime.Register("java/lang/Class", "getSuperclass", "()Ljava/lang/Class;", getSuperclass)
	runtime.Register("java/lang/Class", "getInterfaces0", "()[Ljava/lang/Class;", getInterfaces0)
	runtime.Register("java/lang/Class", "getComponentType", "()Ljava/lang/Class;", getComponentType)

	// type etc.
	runtime.Register("java/lang/Class", "isInterface", "()Z", isInterface)
	runtime.Register("java/lang/Class", "isArray", "()Z", isArray)
	runtime.Register("java/lang/Class", "isPrimitive", "()Z", isPrimitive)

	// dynamic loading
	runtime.Register("java/lang/Class", "forName0", "(Ljava/lang/String;ZLjava/lang/ClassLoader;Ljava/lang/Class;)Ljava/lang/Class;", forName0)

	// new
	runtime.Register("java/lang/Class", "newInstance", "()Ljava/lang/Object;", newInstance)

	// others
	runtime.Register("java/lang/Class", "desiredAssertionStatus0", "(Ljava/lang/Class;)Z", desiredAssertionStatus0)
	runtime.Register("java/lang/Class", "isAssignableFrom", "(Ljava/lang/Class;)Z", isAssignableFrom)
}

// ============================================================
// getPrimitiveClass - Class.getPrimitiveClass(String)
// ============================================================

// getPrimitiveClass
// Java signature: static native Class<?> getPrimitiveClass(String name);
// get primitive class like int.class ...
func getPrimitiveClass(frame *runtime.Frame) {
	nameObj := frame.LocalVars().GetRef(0)
	if nameObj == nil {
		frame.OperandStack().PushRef(nil)
		return
	}

	// to go string
	goName := heap.GoString(nameObj.(*heap.Object))

	// get classLoader
	classLoader := frame.Method().Class().Loader()
	primitiveClass := classLoader.GetPrimitiveClass(goName)
	if primitiveClass == nil {
		panic("java.lang.ClassNotFoundException: " + goName)
	}

	// return java.lang.Class Object
	frame.OperandStack().PushRef(primitiveClass.JClass())
}

// ============================================================
// getName0 - Class.getName0()
// ============================================================
// Java signature: private native String getName0();
// get class name (full name)
func getName0(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)
	// get this class
	class := this.Extra().(*method_area.Class)
	javaName := class.JavaName()

	// create Java String and return it.
	jString := heap.InternString(javaName, class.Loader())
	frame.OperandStack().PushRef(jString)
}

// ============================================================
// initClassName - Class.initClassName()
// ============================================================
// Java signature: private native String initClassName();
// get class full name（new version JDK）
func initClassName(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)

	fmt.Printf("@@ DEBUG - initClassName, this: %s\n", this.String())

	class := this.Extra().(*method_area.Class)

	fmt.Printf("@@ DEBUG - initClassName, class: %s\n", class.String())

	// JVM inner: java/lang/String
	// Java API: java.lang.String
	javaName := class.JavaName()

	// create Java String and return
	jString := heap.InternString(javaName, class.Loader())

	// cache the name

	frame.OperandStack().PushRef(jString)
}

// ============================================================
// getSuperclass - Class.getSuperclass()
// ============================================================
// Java signature: public native Class<? super T> getSuperclass();
// get parent Class Object
func getSuperclass(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)
	class := this.Extra().(*method_area.Class)

	// return nil if interface (interface don't have parent)
	if class.IsInterface() {
		frame.OperandStack().PushRef(nil)
		return
	}
	// primitive type don't have parent
	if class.IsPrimitive() {
		frame.OperandStack().PushRef(nil)
		return
	}
	// java.lang.Object don't have super class
	superClass := class.SuperClass()
	if superClass == nil {
		frame.OperandStack().PushRef(nil)
		return
	}

	frame.OperandStack().PushRef(superClass.JClass())
}

// ============================================================
// getInterfaces0 - Class.getInterfaces0()
// ============================================================
// Java signature: private native Class<?>[] getInterfaces0();
// get all implements interfaces
func getInterfaces0(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)
	class := this.Extra().(*method_area.Class)

	interfaces := class.Interfaces()
	interfaceCount := int32(len(interfaces))

	// get Class[] array type
	loader := frame.Method().Class().Loader()
	classArrayClass := loader.LoadClass("[Ljava/lang/Class;", false)
	classArray := heap.NewRefArray(classArrayClass, interfaceCount)

	for i, iface := range interfaces {
		classArray.SetArrayRef(int32(i), iface.JClass())
	}

	frame.OperandStack().PushRef(classArray)
}

// ============================================================
// getComponentType - Class.getComponentType()
// ============================================================
// Java signature: public native Class<?> getComponentType();
// get array element type, return null if not array
func getComponentType(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)
	class := this.Extra().(*method_area.Class)

	if !class.IsArray() {
		frame.OperandStack().PushRef(nil)
		return
	}

	componentClass := class.ComponentClass()
	if componentClass == nil {
		frame.OperandStack().PushRef(nil)
		return
	}

	frame.OperandStack().PushRef(componentClass.JClass())
}

// ============================================================
// isInterface - Class.isInterface()
// ============================================================
// Java signature: public native boolean isInterface();
func isInterface(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)
	class := this.Extra().(*method_area.Class)
	if class.IsInterface() {
		frame.OperandStack().PushTrue()
	} else {
		frame.OperandStack().PushFalse()
	}
}

// ============================================================
// isArray - Class.isArray()
// ============================================================
// Java signature: public native boolean isArray();
func isArray(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)
	class := this.Extra().(*method_area.Class)

	if class.IsArray() {
		frame.OperandStack().PushTrue()
	} else {
		frame.OperandStack().PushFalse()
	}
}

// ============================================================
// isPrimitive - Class.isPrimitive()
// ============================================================
// Java signature: public native boolean isPrimitive();
func isPrimitive(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)
	class := this.Extra().(*method_area.Class)

	if class.IsPrimitive() {
		frame.OperandStack().PushTrue()
	} else {
		frame.OperandStack().PushFalse()
	}
}

// ============================================================
// forName0 - Class.forName0()
// ============================================================
// Java signature: private static native Class<?> forName0(String name, boolean initialize, ClassLoader loader, Class<?> caller);
// Dynamic Load Class
func forName0(frame *runtime.Frame) {
	// Params:
	// [0] name - class name (Java: java.lang.String)
	// [1] initialize - is init?
	// [2] loader - ClassLoader
	// [3] caller - caller Class

	nameObj := frame.LocalVars().GetRef(0)
	if nameObj == nil {
		frame.NativeThrow(exception.NewNullPointerException(frame))
		return
	}

	// get class name
	javaName := heap.GoString(nameObj.(*heap.Object))
	// java.lang.String → java/lang/String
	jvmName := strings.ReplaceAll(javaName, ".", "/")

	loader := frame.Method().Class().Loader()
	class := loader.LoadClass(jvmName, false)

	// TODO: 根據 initialize 參數決定是否執行 <clinit>, MVP 簡化：總是初始化

	// return jCLass Object
	frame.OperandStack().PushRef(class.JClass())
}

// ============================================================
// newInstance - Class.newInstance() (Deprecated in Java 9+)
// ============================================================
// Java signature: public T newInstance();
// using non-constructor create new class
func newInstance(frame *runtime.Frame) {
	this := frame.LocalVars().GetThis().(*heap.Object)
	class := this.Extra().(*method_area.Class)

	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationException")
	}

	obj := class.NewObject()

	// find non-constructor init method
	constructor := class.GetMethod("<init>", "()V")
	if constructor == nil {
		panic("java.lang.InstantiationException: no default constructor")
	}

	// TODO: 呼叫建構子 這需要建立新的 Frame 並執行 <init>
	// TODO: 簡化：直接返回未初始化的物件 這在實際應用中是不正確的！

	frame.OperandStack().PushRef(obj)
}

// ============================================================
// desiredAssertionStatus0 - Class.desiredAssertionStatus0()
// ============================================================
// Java signature: private static native boolean desiredAssertionStatus0(Class<?> clazz);
// TODO:（MVP 簡化：總是返回 false）
func desiredAssertionStatus0(frame *runtime.Frame) {
	frame.OperandStack().PushFalse() // assertions disabled
}

// ============================================================
// isAssignableFrom - Class.isAssignableFrom(Class)
// ============================================================
// Java signature: public native boolean isAssignableFrom(Class<?> cls);
func isAssignableFrom(frame *runtime.Frame) {
	// this -- -- -- -- -- -- -- -- -- -- -- -- --
	this := frame.LocalVars().GetThis()
	if this == nil {
		frame.NativeThrow(exception.NewNullPointerException(frame))
		return
	}
	thisObj := this.(*heap.Object)
	thisClass := thisObj.Extra().(*method_area.Class)

	// other -- -- -- -- -- -- -- -- -- -- -- -- --
	clsRef := frame.LocalVars().GetRef(1)
	if clsRef == nil {
		frame.NativeThrow(exception.NewNullPointerException(frame))
		return
	}
	clsObj := clsRef.(*heap.Object)
	otherClass := clsObj.Extra().(*method_area.Class)

	// other is sub of this
	if thisClass.IsAssignableFrom(otherClass) {
		frame.OperandStack().PushTrue()
	} else {
		frame.OperandStack().PushFalse()
	}
}
