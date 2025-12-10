package references

import (
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
)

// InitClass
// execute <clinit> method (if class have it)
//
// init order：
// 1. parent class's <clinit>
// 2. this class's <clinit>
//
// this func will create a new Frame to execute <clinit>
// we need call RevertNextPC() let interpreor do this `new` again after init
func InitClass(thread *runtime.Thread, class *method_area.Class) {
	// mark class is doing init
	class.StartInit()

	// call <clinit>
	scheduleClinit(thread, class)

	// init parents (recursive until all parents are init)
	initSuperClass(thread, class)
}

// InvokeMethod call method common func
// usage: invokestatic, invokevirtual
func InvokeMethod(invokerFrame *runtime.Frame, method *method_area.Method) {
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
