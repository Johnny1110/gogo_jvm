package stack

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

type DUP struct{ base.NoOperandsInstruction }

func (d *DUP) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	slot := stack.PopSlot()
	stack.PushSlot(slot)
	stack.PushSlot(slot)
}

func (d *DUP) Opcode() uint8 {
	return 0x59
}

type DUP_X1 struct{ base.NoOperandsInstruction }

func (d *DUP_X1) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	slot1 := stack.PopSlot()
	slot2 := stack.PopSlot()
	stack.PushSlot(slot1)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

func (d *DUP_X1) Opcode() uint8 {
	return 0x5A
}

type DUP_X2 struct{ base.NoOperandsInstruction }

func (d *DUP_X2) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	slot1 := stack.PopSlot()
	slot2 := stack.PopSlot()
	slot3 := stack.PopSlot()
	stack.PushSlot(slot1)
	stack.PushSlot(slot3)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

func (d *DUP_X2) Opcode() uint8 {
	return 0x5B
}

type DUP2 struct{ base.NoOperandsInstruction }

func (d *DUP2) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	slot1 := stack.PopSlot()
	slot2 := stack.PopSlot()
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

func (d *DUP2) Opcode() uint8 {
	return 0x5C
}

type DUP2_X1 struct{ base.NoOperandsInstruction }

func (d *DUP2_X1) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	slot1 := stack.PopSlot()
	slot2 := stack.PopSlot()
	slot3 := stack.PopSlot()
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
	stack.PushSlot(slot3)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

func (d *DUP2_X1) Opcode() uint8 {
	return 0x5D
}

type DUP2_X2 struct{ base.NoOperandsInstruction }

func (d *DUP2_X2) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	slot1 := stack.PopSlot()
	slot2 := stack.PopSlot()
	slot3 := stack.PopSlot()
	slot4 := stack.PopSlot()
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
	stack.PushSlot(slot4)
	stack.PushSlot(slot3)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

func (d *DUP2_X2) Opcode() uint8 {
	return 0x5E
}

type POP struct{ base.NoOperandsInstruction }

func (p *POP) Execute(frame *runtime.Frame) {
	frame.OperandStack().PopSlot()
}

func (p *POP) Opcode() uint8 {
	return 0x57
}

type POP2 struct{ base.NoOperandsInstruction }

func (p *POP2) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	stack.PopSlot()
	stack.PopSlot()
}

func (p *POP2) Opcode() uint8 {
	return 0x58
}

type SWAP struct{ base.NoOperandsInstruction }

func (s *SWAP) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	slot1 := stack.PopSlot()
	slot2 := stack.PopSlot()
	stack.PushSlot(slot1)
	stack.PushSlot(slot2)
}

func (s *SWAP) Opcode() uint8 {
	return 0x5F
}
