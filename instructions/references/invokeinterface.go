package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// ============================================================
// INVOKEINTERFACE - invoke interface method
// ============================================================
// opcode: 0xB9
// operands: 2 bytes (constant pool index) + 1 byte (count) + 1 byte (0)
// stack: [objectref, arg1, arg2, ...] → [...]
//
// count: number of argument slots + 1 (for 'this')
// last 0: historical, must be 0

type INVOKEINTERFACE struct {
	base.Index16Instruction // constant pool index, pointing to InterfaceMethodRef
}

func (i *INVOKEINTERFACE) FetchOperands(reader *base.BytecodeReader) {
	i.Index = uint(reader.ReadUint16())
	_ = reader.ReadUint8() // count - read but not used
	_ = reader.ReadUint8() // must be 0 - read and ignored
}

func (i *INVOKEINTERFACE) Execute(frame *runtime.Frame) {
	// 1. get current class and runtime constant pool
	currentClass := frame.Method().Class()
	rtcp := currentClass.ConstantPool()

	// 2. get InterfaceMethodRef from constant pool
	methodRef := rtcp.GetConstant(i.Index).(*method_area.InterfaceMethodRef)

	// 3. resolve interface method
	resolvedMethod := methodRef.ResolvedInterfaceMethod()

	// 4. check: method must not be static or <init>
	if resolvedMethod.IsStatic() || resolvedMethod.Name() == "<init>" {
		exceptionObj := NewIncompatibleClassChangeError(frame, "interface method can not be static or <init>")
		ThrowException(frame, exceptionObj)
		return
	}

	// 5. get objectref from operand stack
	// objectref is under args: stack = [..., objectref, arg1, arg2, ...]
	objectref := frame.OperandStack().PeekRefFromTop(resolvedMethod.ArgSlotCount() - 1)

	// 6. null check
	if objectref == nil {
		exceptionObj := NewNullPointerException(frame)
		ThrowException(frame, exceptionObj)
		return
	}

	// 7. get object's actual class
	object := objectref.(*heap.Object)
	actualClass := object.Class().(*method_area.Class)

	// 8. lookup method in actual class (dynamic binding)
	// This is the key difference from invokevirtual:
	// We need to search from the actual class, not use vtable
	methodToCall := lookupMethodForInvokeInterface(actualClass, resolvedMethod.Name(), resolvedMethod.Descriptor())

	// 9. check if method found
	if methodToCall == nil {
		exceptionObj := NewAbstractMethodError(frame, actualClass.Name(), resolvedMethod.Name(), resolvedMethod.Descriptor())
		ThrowException(frame, exceptionObj)
		return
	}

	// 10. check if method is abstract
	if methodToCall.IsAbstract() {
		exceptionObj := NewAbstractMethodError(frame, actualClass.Name(), resolvedMethod.Name(), resolvedMethod.Descriptor())
		ThrowException(frame, exceptionObj)
		return
	}

	// 11. check access: interface methods must be public
	if !methodToCall.IsPublic() {
		exceptionObj := NewIllegalAccessError(frame, actualClass.Name(), resolvedMethod.Name(), resolvedMethod.Descriptor())
		ThrowException(frame, exceptionObj)
		return
	}

	// 12. invoke method
	invokeMethod(frame, methodToCall)
}

func (i *INVOKEINTERFACE) Opcode() uint8 {
	return 0xB9
}

// ============================================================
// lookupMethodForInvokeInterface
// ============================================================
// Search method implementation in class and its superclasses,
// then in interfaces (for Java 8 default methods)
//
// Search order:
//  1. Current class methods
//  2. Superclass methods (recursively)
//  3. Interface default methods (for Java 8+)
func lookupMethodForInvokeInterface(class *method_area.Class, name, descriptor string) *method_area.Method {
	// 1. search in class and super class.
	method := class.GetMethod(name, descriptor)
	if method != nil {
		return method
	}

	// 2: Search in interfaces (for default methods)
	// This handles Java 8 default methods
	return lookupMethodInInterfaces(class, name, descriptor)
}

// lookupMethodInInterfaces searches for method in interfaces
// This is needed for Java 8 default method support
func lookupMethodInInterfaces(class *method_area.Class, name, descriptor string) *method_area.Method {
	// 1. Search in directly implemented interfaces
	for _, iface := range class.Interfaces() {
		// 1-1. Iterate all methods of this interface
		for _, method := range iface.Methods() {
			// name, desc, and not abs, then return it.
			if method.Name() == name && method.Descriptor() == descriptor && !method.IsAbstract() {
				return method
			}
		}

		// 1-2. Recursively search in super-interfaces
		method := lookupMethodInInterfaces(iface, name, descriptor)
		if method != nil {
			return method
		}
	}

	// 2. Search in superclass's interfaces
	if class.SuperClass() != nil {
		return lookupMethodInInterfaces(class.SuperClass(), name, descriptor)
	}

	return nil
}
