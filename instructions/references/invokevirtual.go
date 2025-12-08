package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
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
	// TODO
}

func (i *INVOKEVIRTUAL) Opcode() uint8 {
	return 0xB7
}
