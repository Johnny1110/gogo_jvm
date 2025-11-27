package control

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// ============================================================
// GOTO jump without condition
// ============================================================

// GOTO jump without condition
// opcode = 0xA7
// operands: 2 bytes (signed)
// jump target = current PC + Offsets
type GOTO struct{ base.BranchInstruction }

func (g *GOTO) Execute(frame *runtime.Frame) {
	branch(frame, g.Offset)
}

func (g *GOTO) Opcode() uint8 {
	return 0xA7
}

// ============================================================
// IF_ICMP: compare 2 int and jump
// ============================================================

// ex: if (a < b) { ... }
// compiled:
//   iload_0       // load a
//   iload_1       // load b
//   if_icmpge L1  // if a >= b jump to L1（skip if block）
//   ...           // if block logics
//   L1:           // if ending

// IF_ICMPEQ jump if v1 == v2
// opcode = 0x9F
type IF_ICMPEQ struct{ base.BranchInstruction }

func (i *IF_ICMPEQ) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v1 == v2 {
		branch(frame, i.Offset)
	}
}

func (g *IF_ICMPEQ) Opcode() uint8 {
	return 0x9F
}

// IF_ICMPNE jump if v1 != v2
// opcode = 0xA0
type IF_ICMPNE struct{ base.BranchInstruction }

func (i *IF_ICMPNE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v1 != v2 {
		branch(frame, i.Offset)
	}
}

func (g *IF_ICMPNE) Opcode() uint8 {
	return 0xA0
}

// IF_ICMPLT jump if v1 < v2
// opcode = 0xA1
type IF_ICMPLT struct{ base.BranchInstruction }

func (i *IF_ICMPLT) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v1 < v2 {
		branch(frame, i.Offset)
	}
}

func (g *IF_ICMPLT) Opcode() uint8 {
	return 0xA1
}

// IF_ICMPGE jump if v1 >= v2
// opcode = 0xA2
type IF_ICMPGE struct{ base.BranchInstruction }

func (i *IF_ICMPGE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v1 >= v2 {
		branch(frame, i.Offset)
	}
}

func (g *IF_ICMPGE) Opcode() uint8 {
	return 0xA2
}

// IF_ICMPGT jump if v1 > v2
// opcode = 0xA3
type IF_ICMPGT struct{ base.BranchInstruction }

func (i *IF_ICMPGT) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v1 > v2 {
		branch(frame, i.Offset)
	}
}

func (g *IF_ICMPGT) Opcode() uint8 {
	return 0xA3
}

// IF_ICMPLE jump if v1 <= v2
// opcode = 0xA4
type IF_ICMPLE struct{ base.BranchInstruction }

func (i *IF_ICMPLE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v1 <= v2 {
		branch(frame, i.Offset)
	}
}

func (g *IF_ICMPLE) Opcode() uint8 {
	return 0xA4
}

// ============================================================
// IF Series: Compare int and 0 then jump
// ============================================================
// Usage: if (a == 0), if (flag), while (count > 0) ...

// IFEQ jump if val == 0
// opcode = 0x99
type IFEQ struct{ base.BranchInstruction }

func (i *IFEQ) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	if val == 0 {
		branch(frame, i.Offset)
	}
}

func (g *IFEQ) Opcode() uint8 {
	return 0x99
}

// IFNE jump if val != 0
// opcode = 0x9A
type IFNE struct{ base.BranchInstruction }

func (i *IFNE) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	if val != 0 {
		branch(frame, i.Offset)
	}
}

func (g *IFNE) Opcode() uint8 {
	return 0x9A
}

// IFLT jump if val < 0
// opcode = 0x9B
type IFLT struct{ base.BranchInstruction }

func (i *IFLT) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	if val < 0 {
		branch(frame, i.Offset)
	}
}

func (g *IFLT) Opcode() uint8 {
	return 0x9B
}

// IFGE jump if val >= 0
// opcode = 0x9C
type IFGE struct{ base.BranchInstruction }

func (i *IFGE) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	if val >= 0 {
		branch(frame, i.Offset)
	}
}

func (g *IFGE) Opcode() uint8 {
	return 0x9C
}

// IFGT jump if val > 0
// opcode = 0x9D
type IFGT struct{ base.BranchInstruction }

func (i *IFGT) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	if val > 0 {
		branch(frame, i.Offset)
	}
}

func (g *IFGT) Opcode() uint8 {
	return 0x9D
}

// IFLE jump if val <= 0
// opcode = 0x9E
type IFLE struct{ base.BranchInstruction }

func (i *IFLE) Execute(frame *runtime.Frame) {
	val := frame.OperandStack().PopInt()
	if val <= 0 {
		branch(frame, i.Offset)
	}
}

func (g *IFLE) Opcode() uint8 {
	return 0x9E
}

// ============================================================
// IF_ACMP Series: compare 2 ref then jump
// ============================================================

// IF_ACMPEQ jump if both ref equals
// opcode = 0xA5
type IF_ACMPEQ struct{ base.BranchInstruction }

func (i *IF_ACMPEQ) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	ref2 := stack.PopRef()
	ref1 := stack.PopRef()
	if ref1 == ref2 {
		branch(frame, i.Offset)
	}
}

func (g *IF_ACMPEQ) Opcode() uint8 {
	return 0xA5
}

// IF_ACMPNE jump if both ref are not equals
// opcode = 0xA6
type IF_ACMPNE struct{ base.BranchInstruction }

func (i *IF_ACMPNE) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	ref2 := stack.PopRef()
	ref1 := stack.PopRef()
	if ref1 != ref2 {
		branch(frame, i.Offset)
	}
}

func (g *IF_ACMPNE) Opcode() uint8 {
	return 0xA6
}

// ============================================================
// IFNULL / IFNONNULL: check ref is null
// ============================================================

// IFNULL jump if ref is nil
// opcode = 0xC6
type IFNULL struct{ base.BranchInstruction }

func (i *IFNULL) Execute(frame *runtime.Frame) {
	ref := frame.OperandStack().PopRef()
	if ref == nil {
		branch(frame, i.Offset)
	}
}

func (g *IFNULL) Opcode() uint8 {
	return 0xC6
}

// IFNONNULL jump if ref is not nil
// opcode = 0xC7
type IFNONNULL struct{ base.BranchInstruction }

func (i *IFNONNULL) Execute(frame *runtime.Frame) {
	ref := frame.OperandStack().PopRef()
	if ref != nil {
		branch(frame, i.Offset)
	}
}

func (g *IFNONNULL) Opcode() uint8 {
	return 0xC7
}
