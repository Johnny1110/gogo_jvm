package references

import (
	"github.com/Johnny1110/gogo_jvm/common"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
)

// initClass
// execute <clinit> method (if class have it)
//
// init order：
// 1. parent class's <clinit>
// 2. this class's <clinit>
//
// this func will create a new Frame to execute <clinit>
// we need call RevertNextPC() let interpreor do this `new` again after init
func initClass(thread *runtime.Thread, class *method_area.Class) {
	// mark class is doing init
	class.StartInit()

	// call <clinit>
	scheduleClinit(thread, class)

	// init parents (recursive until all parents are init)
	initSuperClass(thread, class)
}

// invokeMethod call method common func
// usage: invokestatic, invokevirtual
func invokeMethod(invokerFrame *runtime.Frame, method *method_area.Method) {
	// 1, get current thread
	thread := invokerFrame.Thread()

	// 2. create a new frame (represent new method)
	newFrame := thread.NewFrameWithMethod(method)
	thread.PushFrame(newFrame)

	// 3. pass vars
	argSlotCount := int(method.ArgSlotCount())
	if argSlotCount > 0 {
		for i := argSlotCount - 1; i >= 0; i-- {
			// the passing var will be standby in invokerFrame's stack.
			// pop it from invokerFrame's stack
			slot := invokerFrame.OperandStack().PopSlot()
			// put into newFrame's LocalVars
			newFrame.LocalVars().SetSlot(uint(i), slot)
		}
	}

	// no need to reset PC, new frame nextPC will be default 0
}

// invokeNativeMethod
// 1. get args from callerFrame
// 2. put args into a temp LocalVars
// 3. pass args to native Go func
// no need a read frame to do native method.
func invokeNativeMethod(callerFrame *runtime.Frame, nativeMethod runtime.NativeMethod, descriptor string) {
	// calculate args slot count including this.
	argSlotCount := calcArgSlotCount(descriptor) + 1 // LocalVars[0] = this, so we need + 1

	// create a temp Frame for store params (not require to push into JVMStack)
	tempFrame := runtime.NewNativeFrame(callerFrame.Thread(), uint16(argSlotCount))

	// pop args from caller op-stack put into tempFrame's LocalVars
	stack := callerFrame.OperandStack() // stack: [argN, argN-1, ..., arg1, this]
	localVars := tempFrame.LocalVars()  // localVars: [this, arg1, ... argN-1, argN]

	// push args into temp LocalVars
	for i := argSlotCount - 1; i >= 0; i-- { // i := argSlotCount - 1 -> skip this
		slot := stack.PopSlot()
		localVars.SetSlot(uint(i), slot)
	}

	// call native method
	nativeMethod(tempFrame)

	// TODO: currently we only support void native method
	// TODO: we should return val if called native method has return val.
}

// calcArgSlotCount calculate arg count by method descriptor
// Ex:
//
//	()V           → 0
//	(I)V          → 1 int
//	(J)V          → 1 long (2 slots)
//	(II)I         → 2 int (2 slots)
//	(IJD)V        → int(1) + long(2) + double(2) = 5 slots
func calcArgSlotCount(descriptor string) int {
	slotCount := 0
	i := 1 // skip '('

	// iterate args descriptors
	for descriptor[i] != ')' {
		switch descriptor[i] {
		case 'B', 'C', 'F', 'I', 'S', 'Z':
			// byte, char, float, int, short, boolean → 1 slot
			slotCount++
			i++
		case 'D', 'J':
			// double, long → 2 slots
			slotCount += 2
			i++
		case 'L':
			// Object Ref: Ljava/lang/String;
			slotCount++
			for descriptor[i] != ';' {
				i++
			}
			i++ // skip ';'
		case '[':
			// array
			slotCount++
			// skip array element type
			for descriptor[i] == '[' {
				i++ // skip '['
			}

			if descriptor[i] == 'L' { // Is Object Type
				for descriptor[i] != ';' {
					i++ // skip obj type
				}
				i++ // skip ';'
			} else { // Is Basic Type
				i++
			}
		default:
			panic(common.NewJavaException("", "Calculate method arg slot count failed, unknown descriptor type: "+string(descriptor[i])))
		}
	}
	return slotCount
}

// pushFieldValue according to descriptor, copy val from slots and put into stack
func pushFieldValue(stack *runtime.OperandStack, slots rtcore.Slots, slotId uint, descriptor string) {
	// 1. pop val from stack
	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		// boolean, byte, char, short, int -> int32
		stack.PushInt(slots.GetInt(slotId))
	case 'F':
		// float
		stack.PushFloat(slots.GetFloat(slotId))
	case 'J':
		// long
		stack.PushLong(slots.GetLong(slotId))
	case 'D':
		// double
		stack.PushDouble(slots.GetDouble(slotId))
	case 'L', '[':
		// ObjectRef OR ArrayRef
		stack.PushRef(slots.GetRef(slotId))
	default:
		panic("Unknown field descriptor: " + descriptor)
	}
}

// popAndSetFieldValue pop val from stack and put into slots
func popAndSetFieldValue(stack *runtime.OperandStack, slots rtcore.Slots, slotId uint, descriptor string) {
	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		// boolean, byte, char, short, int -> int32
		slots.SetInt(slotId, stack.PopInt())
	case 'F':
		// float
		slots.SetFloat(slotId, stack.PopFloat())
	case 'J':
		// long
		slots.SetLong(slotId, stack.PopLong())
	case 'D':
		slots.SetDouble(slotId, stack.PopDouble())
	case 'L', '[':
		// ref or array
		slots.SetRef(slotId, stack.PopRef())
	default:
		panic("Unknown field descriptor: " + descriptor)
	}
}

// checkNotNull check ref is not null
func checkNotNull(ref interface{}) {
	if ref == nil {
		panic("java.class.NullPointerException")
	}
}

// isSubClassOf check child is a subclass from parent
func isSubClassOf(child, parent *method_area.Class) bool {
	for c := child.SuperClass(); c != nil; c = c.SuperClass() {
		if c == parent {
			return true
		}
	}
	return false
}

// ============================================================
// isInstanceOf - core check isInstance func
// ============================================================
// declare:  S = object's type, T = target class
// situation-1: S is normal class
//   - T is normal class: S == T or S is T's subclass
//   - T is interface: S implements T
//
// situation-2: S is Array type
//   - T is Object: true（all array is a kind of Object）
//   - T is Cloneable: true（JLS Standard）
//   - T is Serializable: true （JLS Standard）
//   - T 是 Array type TC[]: S element can pass to TC
//   - otherwise: false
func isInstanceOf(targetObject *heap.Object, targetClass *method_area.Class) bool {

	if targetObject.IsArray() { // handle situation-2:  S is Array type
		return isArrayInstanceOf(targetObject, targetClass)
	}

	targetObjectClass := targetObject.Class()

	if targetObjectClass != nil { // handle situation-1: S is normal class
		// using Class's IsAssignableFrom()
		S := targetObjectClass.(*method_area.Class)
		return targetClass.IsAssignableFrom(S) // targetClass = S
	}

	return false
}

// isArrayInstanceOf
// - int[] instanceof Object → true
// - int[] instanceof Cloneable → true
// - int[] instanceof Serializable → true
// - String[] instanceof Object[] → true (class String extends Object)
// - int[] instanceof Object[] → false (java basic type is not extended from Object)
func isArrayInstanceOf(arrObject *heap.Object, targetClass *method_area.Class) bool {
	targetClassName := targetClass.Name()

	// int[] instanceof Object → true
	if targetClassName == "java/lang/Object" {
		return true
	}
	// int[] instanceof Cloneable → true
	if targetClassName == "java/lang/Cloneable" {
		return true
	}
	// int[] instanceof Serializable → true
	if targetClassName == "java/io/Serializable" {
		return true
	}

	// TODO: 完整實現需要陣列類型系統支援 等後續逐步優化
	if arrObject.Class() != nil { // TODO: 短解，我在建立 Object Array 時候會把 element type (class) 放菜 object.class 內
		arrClass := arrObject.Class().(*method_area.Class)
		return targetClass.IsAssignableFrom(arrClass)
	}

	return false

}
