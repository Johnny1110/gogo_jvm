package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// INVOKESPECIAL
// 1. constructor <init>
// 2. private method
// 3. parent method (super.xxx())
// opcode = 0xB7
// operands: 2 bytes (constant pool index)
// stack: [objectref, arg1, arg2, ...] → [...]
type INVOKESPECIAL struct {
	base.Index16Instruction
}

func (i *INVOKESPECIAL) Execute(frame *runtime.Frame) {
	// 1. get lang & ctcp
	currentClass := frame.Method().Class()
	rtcp := currentClass.ConstantPool()

	// 2. get methodRef
	methodRef := rtcp.GetConstant(i.Index).(*method_area.MethodRef)

	// 3. resolve target method's lang and method
	resolvedMethod := methodRef.ResolvedMethod()
	resolvedClass := methodRef.ResolvedClass()

	// 4. check <init> must invoke by invokespecial (# README: 實作 invokespecial 時的發現)
	if resolvedMethod.Name() == "<init>" && resolvedMethod.Class() != resolvedClass {
		// why?
		// currentClass -> rtcp -> methodRef -> resolvedMethod (parent's method)
		// methodRef -> resolvedClass (this.lang)
		// this.lang != parent's method.Class()
		panic("java.lang.NoSuchMethodError")
	}

	// 5. check can not call static method
	if resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 6. get objectref(this) from stack, objectref is under args, like stack: [objectref, arg1, arg2, ...]
	objectref := frame.OperandStack().PeekRefFromTop(resolvedMethod.ArgSlotCount() - 1)
	if objectref == nil {
		panic("java.lang.NullPointerException")
	}

	// 7. determine which method to call. most of all, using resolvedMethod
	methodToCall := resolvedMethod

	// 8. call super.xxx()
	// - currentClass is son fo resolvedClass
	// - currentClass has ACC_SUPER access_flag
	// - not constructor call
	if isSubClassOf(currentClass, resolvedClass) &&
		currentClass.IsSuper() && // ACC_SUPER
		resolvedMethod.Name() != "<init>" {
		// look up start from super lang.
		methodToCall = currentClass.SuperClass().GetMethod(resolvedMethod.Name(), resolvedMethod.Descriptor())
	}

	if methodToCall == nil {
		panic("java.lang.AbstractMethodError")
	}

	// 9. invoke method
	invokeMethod(frame, methodToCall)
}

func (i *INVOKESPECIAL) Opcode() uint8 {
	return 0xB7
}
