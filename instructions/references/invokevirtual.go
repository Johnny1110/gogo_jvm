package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// INVOKEVIRTUAL
// ex:
//
//	Animal a = new Dog();
//	a.speak();  // call Dog.speak() not Animal.speak()
//
// opcode = 0xB6
// operands: 2 bytes (constant pool index)
// stack: [objectref, arg1, arg2, ...] → [...]
type INVOKEVIRTUAL struct {
	base.Index16Instruction
}

func (i *INVOKEVIRTUAL) Execute(frame *runtime.Frame) {
	// 1. get currentClass and rtcp
	currentClass := frame.Method().Class()
	rtcp := currentClass.ConstantPool()

	// 2. load methodRef
	methodRef := rtcp.GetConstant(i.Index).(*method_area.MethodRef)

	// 3. load method
	resolvedMethod := methodRef.ResolvedMethod()

	// 4. check is static (no static)
	if resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 5. get objectref
	objectref := frame.OperandStack().PeekRefFromTop(resolvedMethod.ArgSlotCount() - 1)
	if objectref == nil {
		panic("java.lang.NullPointerException")
	}

	// 6. get object and object's class
	object := objectref.(*heap.Object)
	actualClass := object.Class().(*method_area.Class)

	// 7. dynamic binding method (if can not find in this lang, lookup to super)
	methodToCall := actualClass.GetMethod(resolvedMethod.Name(), resolvedMethod.Descriptor())

	if methodToCall == nil || methodToCall.IsAbstract() {
		panic("java.lang.AbstractMethodError")
	}

	// 8. invoke method
	InvokeMethod(frame, methodToCall)
}

func (i *INVOKEVIRTUAL) Opcode() uint8 {
	return 0xB6
}
