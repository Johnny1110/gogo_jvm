package instructions

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/instructions/arrays"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/base/opcodes"
	"github.com/Johnny1110/gogo_jvm/instructions/constants"
	"github.com/Johnny1110/gogo_jvm/instructions/control"
	"github.com/Johnny1110/gogo_jvm/instructions/loads"
	"github.com/Johnny1110/gogo_jvm/instructions/math"
	"github.com/Johnny1110/gogo_jvm/instructions/references"
	"github.com/Johnny1110/gogo_jvm/instructions/stack"
	"github.com/Johnny1110/gogo_jvm/instructions/stores"
)

// ============================================================
// instructions factory（decoder）
// ============================================================

// !!: 預先建立好的 instructions 都是無狀態的，不可以把帶狀態 (index, offset) 的指令建立單例
var (
	// ============ Constants ============
	nop         = &constants.NOP{}
	aconst_null = &constants.ACONST_NULL{}
	iconst_m1   = &constants.ICONST_M1{}
	iconst_0    = &constants.ICONST_0{}
	iconst_1    = &constants.ICONST_1{}
	iconst_2    = &constants.ICONST_2{}
	iconst_3    = &constants.ICONST_3{}
	iconst_4    = &constants.ICONST_4{}
	iconst_5    = &constants.ICONST_5{}
	lconst_0    = &constants.LCONST_0{}
	lconst_1    = &constants.LCONST_1{}
	fconst_0    = &constants.FCONST_0{}
	fconst_1    = &constants.FCONST_1{}
	fconst_2    = &constants.FCONST_2{}
	dconst_0    = &constants.DCONST_0{}
	dconst_1    = &constants.DCONST_1{}

	// ============ Loads ============
	iload_0 = &loads.ILOAD_0{}
	iload_1 = &loads.ILOAD_1{}
	iload_2 = &loads.ILOAD_2{}
	iload_3 = &loads.ILOAD_3{}
	lload_0 = &loads.LLOAD_0{}
	lload_1 = &loads.LLOAD_1{}
	lload_2 = &loads.LLOAD_2{}
	lload_3 = &loads.LLOAD_3{}
	fload_0 = &loads.FLOAD_0{}
	fload_1 = &loads.FLOAD_1{}
	fload_2 = &loads.FLOAD_2{}
	fload_3 = &loads.FLOAD_3{}
	dload_0 = &loads.DLOAD_0{}
	dload_1 = &loads.DLOAD_1{}
	dload_2 = &loads.DLOAD_2{}
	dload_3 = &loads.DLOAD_3{}
	aload_0 = &loads.ALOAD_0{}
	aload_1 = &loads.ALOAD_1{}
	aload_2 = &loads.ALOAD_2{}
	aload_3 = &loads.ALOAD_3{}

	// ============ Stores ============
	istore_0 = &stores.ISTORE_0{}
	istore_1 = &stores.ISTORE_1{}
	istore_2 = &stores.ISTORE_2{}
	istore_3 = &stores.ISTORE_3{}
	lstore_0 = &stores.LSTORE_0{}
	lstore_1 = &stores.LSTORE_1{}
	lstore_2 = &stores.LSTORE_2{}
	lstore_3 = &stores.LSTORE_3{}
	fstore_0 = &stores.FSTORE_0{}
	fstore_1 = &stores.FSTORE_1{}
	fstore_2 = &stores.FSTORE_2{}
	fstore_3 = &stores.FSTORE_3{}
	dstore_0 = &stores.DSTORE_0{}
	dstore_1 = &stores.DSTORE_1{}
	dstore_2 = &stores.DSTORE_2{}
	dstore_3 = &stores.DSTORE_3{}
	astore_0 = &stores.ASTORE_0{}
	astore_1 = &stores.ASTORE_1{}
	astore_2 = &stores.ASTORE_2{}
	astore_3 = &stores.ASTORE_3{}

	// ============ Stack ============
	pop     = &stack.POP{}
	pop2    = &stack.POP2{}
	dup     = &stack.DUP{}
	dup_x1  = &stack.DUP_X1{}
	dup_x2  = &stack.DUP_X2{}
	dup2    = &stack.DUP2{}
	dup2_x1 = &stack.DUP2_X1{}
	dup2_x2 = &stack.DUP2_X2{}
	swap    = &stack.SWAP{}

	// ============ Math ============

	iadd = &math.IADD{}
	ladd = &math.LADD{}
	fadd = &math.FADD{}
	dadd = &math.DADD{}
	isub = &math.ISUB{}
	lsub = &math.LSUB{}
	fsub = &math.FSUB{}
	dsub = &math.DSUB{}
	imul = &math.IMUL{}
	lmul = &math.LMUL{}
	fmul = &math.FMUL{}
	dmul = &math.DMUL{}
	idiv = &math.IDIV{}
	ldiv = &math.LDIV{}
	fdiv = &math.FDIV{}
	ddiv = &math.DDIV{}
	irem = &math.IREM{}
	lrem = &math.LREM{}
	frem = &math.FREM{}
	drem = &math.DREM{}
	ineg = &math.INEG{}
	lneg = &math.LNEG{}
	fneg = &math.FNEG{}
	dneg = &math.DNEG{}

	// ============ Reference ============
	athrow = &references.ATHROW{}

	// ============ Control / Return ============
	ireturn = &control.IRETURN{}
	lreturn = &control.LRETURN{}
	freturn = &control.FRETURN{}
	dreturn = &control.DRETURN{}
	areturn = &control.ARETURN{}
	_return = &control.RETURN{}

	// ============ Arrays ============
	iaload      = &arrays.IALOAD{}
	laload      = &arrays.LALOAD{}
	faload      = &arrays.FALOAD{}
	daload      = &arrays.DALOAD{}
	aaload      = &arrays.AALOAD{}
	baload      = &arrays.BALOAD{}
	caload      = &arrays.CALOAD{}
	saload      = &arrays.SALOAD{}
	iastore     = &arrays.IASTORE{}
	lastore     = &arrays.LASTORE{}
	fastore     = &arrays.FASTORE{}
	dastore     = &arrays.DASTORE{}
	aastore     = &arrays.AASTORE{}
	bastore     = &arrays.BASTORE{}
	castore     = &arrays.CASTORE{}
	sastore     = &arrays.SASTORE{}
	arraylength = &arrays.ARRAYLENGTH{}
)

// NewInstruction return instruction based on input opcodes
func NewInstruction(opcode byte) (base.Instruction, error) {
	switch opcode {
	// const instructions
	case opcodes.NOP:
		return nop, nil
	case opcodes.ACONST_NULL:
		return aconst_null, nil
	case opcodes.ICONST_M1:
		return iconst_m1, nil
	case opcodes.ICONST_0:
		return iconst_0, nil
	case opcodes.ICONST_1:
		return iconst_1, nil
	case opcodes.ICONST_2:
		return iconst_2, nil
	case opcodes.ICONST_3:
		return iconst_3, nil
	case opcodes.ICONST_4:
		return iconst_4, nil
	case opcodes.ICONST_5:
		return iconst_5, nil
	case opcodes.LCONST_0:
		return lconst_0, nil
	case opcodes.LCONST_1:
		return lconst_1, nil
	case opcodes.FCONST_0:
		return fconst_0, nil
	case opcodes.FCONST_1:
		return fconst_1, nil
	case opcodes.FCONST_2:
		return fconst_2, nil
	case opcodes.DCONST_0:
		return dconst_0, nil
	case opcodes.DCONST_1:
		return dconst_1, nil
	case opcodes.BIPUSH:
		return &constants.BIPUSH{}, nil
	case opcodes.SIPUSH:
		return &constants.SIPUSH{}, nil
	case opcodes.LDC:
		return &constants.LDC{}, nil
	case opcodes.LDC_W:
		return &constants.LDC_W{}, nil
	case opcodes.LDC2_W:
		return &constants.LDC2_W{}, nil

	// load instructions
	case opcodes.ILOAD:
		return &loads.ILOAD{}, nil
	case opcodes.LLOAD:
		return &loads.LLOAD{}, nil
	case opcodes.FLOAD:
		return &loads.FLOAD{}, nil
	case opcodes.DLOAD:
		return &loads.DLOAD{}, nil
	case opcodes.ALOAD:
		return &loads.ALOAD{}, nil
	case opcodes.ILOAD_0:
		return iload_0, nil
	case opcodes.ILOAD_1:
		return iload_1, nil
	case opcodes.ILOAD_2:
		return iload_2, nil
	case opcodes.ILOAD_3:
		return iload_3, nil
	case opcodes.LLOAD_0:
		return lload_0, nil
	case opcodes.LLOAD_1:
		return lload_1, nil
	case opcodes.LLOAD_2:
		return lload_2, nil
	case opcodes.LLOAD_3:
		return lload_3, nil
	case opcodes.FLOAD_0:
		return fload_0, nil
	case opcodes.FLOAD_1:
		return fload_1, nil
	case opcodes.FLOAD_2:
		return fload_2, nil
	case opcodes.FLOAD_3:
		return fload_3, nil
	case opcodes.DLOAD_0:
		return dload_0, nil
	case opcodes.DLOAD_1:
		return dload_1, nil
	case opcodes.DLOAD_2:
		return dload_2, nil
	case opcodes.DLOAD_3:
		return dload_3, nil
	case opcodes.ALOAD_0:
		return aload_0, nil
	case opcodes.ALOAD_1:
		return aload_1, nil
	case opcodes.ALOAD_2:
		return aload_2, nil
	case opcodes.ALOAD_3:
		return aload_3, nil

	// store instructions
	case opcodes.ISTORE:
		return &stores.ISTORE{}, nil
	case opcodes.LSTORE:
		return &stores.LSTORE{}, nil
	case opcodes.FSTORE:
		return &stores.FSTORE{}, nil
	case opcodes.DSTORE:
		return &stores.DSTORE{}, nil
	case opcodes.ASTORE:
		return &stores.ASTORE{}, nil
	case opcodes.ISTORE_0:
		return istore_0, nil
	case opcodes.ISTORE_1:
		return istore_1, nil
	case opcodes.ISTORE_2:
		return istore_2, nil
	case opcodes.ISTORE_3:
		return istore_3, nil
	case opcodes.LSTORE_0:
		return lstore_0, nil
	case opcodes.LSTORE_1:
		return lstore_1, nil
	case opcodes.LSTORE_2:
		return lstore_2, nil
	case opcodes.LSTORE_3:
		return lstore_3, nil
	case opcodes.FSTORE_0:
		return fstore_0, nil
	case opcodes.FSTORE_1:
		return fstore_1, nil
	case opcodes.FSTORE_2:
		return fstore_2, nil
	case opcodes.FSTORE_3:
		return fstore_3, nil
	case opcodes.DSTORE_0:
		return dstore_0, nil
	case opcodes.DSTORE_1:
		return dstore_1, nil
	case opcodes.DSTORE_2:
		return dstore_2, nil
	case opcodes.DSTORE_3:
		return dstore_3, nil
	case opcodes.ASTORE_0:
		return astore_0, nil
	case opcodes.ASTORE_1:
		return astore_1, nil
	case opcodes.ASTORE_2:
		return astore_2, nil
	case opcodes.ASTORE_3:
		return astore_3, nil

	// reference
	case opcodes.ATHROW:
		return athrow, nil

	// Array Load Instructions
	case opcodes.IALOAD:
		return iaload, nil
	case opcodes.LALOAD:
		return laload, nil
	case opcodes.FALOAD:
		return faload, nil
	case opcodes.DALOAD:
		return daload, nil
	case opcodes.AALOAD:
		return aaload, nil
	case opcodes.BALOAD:
		return baload, nil
	case opcodes.CALOAD:
		return caload, nil
	case opcodes.SALOAD:
		return saload, nil
	// Array Store Instructions
	case opcodes.IASTORE:
		return iastore, nil
	case opcodes.LASTORE:
		return lastore, nil
	case opcodes.FASTORE:
		return fastore, nil
	case opcodes.DASTORE:
		return dastore, nil
	case opcodes.AASTORE:
		return aastore, nil
	case opcodes.BASTORE:
		return bastore, nil
	case opcodes.CASTORE:
		return castore, nil
	case opcodes.SASTORE:
		return sastore, nil

	// Stack Instructions
	case opcodes.POP:
		return pop, nil
	case opcodes.POP2:
		return pop2, nil
	case opcodes.DUP:
		return dup, nil
	case opcodes.DUP_X1:
		return dup_x1, nil
	case opcodes.DUP_X2:
		return dup_x2, nil
	case opcodes.DUP2:
		return dup2, nil
	case opcodes.DUP2_X1:
		return dup2_x1, nil
	case opcodes.DUP2_X2:
		return dup2_x2, nil
	case opcodes.SWAP:
		return swap, nil

	// math instructions
	case opcodes.IADD:
		return iadd, nil
	case opcodes.LADD:
		return ladd, nil
	case opcodes.FADD:
		return fadd, nil
	case opcodes.DADD:
		return dadd, nil
	case opcodes.ISUB:
		return isub, nil
	case opcodes.LSUB:
		return lsub, nil
	case opcodes.FSUB:
		return fsub, nil
	case opcodes.DSUB:
		return dsub, nil
	case opcodes.IMUL:
		return imul, nil
	case opcodes.LMUL:
		return lmul, nil
	case opcodes.FMUL:
		return fmul, nil
	case opcodes.DMUL:
		return dmul, nil
	case opcodes.IDIV:
		return idiv, nil
	case opcodes.LDIV:
		return ldiv, nil
	case opcodes.FDIV:
		return fdiv, nil
	case opcodes.DDIV:
		return ddiv, nil
	case opcodes.IREM:
		return irem, nil
	case opcodes.LREM:
		return lrem, nil
	case opcodes.FREM:
		return frem, nil
	case opcodes.DREM:
		return drem, nil
	case opcodes.INEG:
		return ineg, nil
	case opcodes.LNEG:
		return lneg, nil
	case opcodes.FNEG:
		return fneg, nil
	case opcodes.DNEG:
		return dneg, nil
	case opcodes.IINC:
		return &math.IINC{}, nil

	// compare instructions
	case opcodes.IFEQ:
		return &control.IFEQ{}, nil
	case opcodes.IFNE:
		return &control.IFNE{}, nil
	case opcodes.IFLT:
		return &control.IFLT{}, nil
	case opcodes.IFGE:
		return &control.IFGE{}, nil
	case opcodes.IFGT:
		return &control.IFGT{}, nil
	case opcodes.IFLE:
		return &control.IFLE{}, nil
	case opcodes.IF_ICMPEQ:
		return &control.IF_ICMPEQ{}, nil
	case opcodes.IF_ICMPNE:
		return &control.IF_ICMPNE{}, nil
	case opcodes.IF_ICMPLT:
		return &control.IF_ICMPLT{}, nil
	case opcodes.IF_ICMPGE:
		return &control.IF_ICMPGE{}, nil
	case opcodes.IF_ICMPGT:
		return &control.IF_ICMPGT{}, nil
	case opcodes.IF_ICMPLE:
		return &control.IF_ICMPLE{}, nil
	case opcodes.IF_ACMPEQ:
		return &control.IF_ACMPEQ{}, nil
	case opcodes.IF_ACMPNE:
		return &control.IF_ACMPNE{}, nil
	case opcodes.GOTO:
		return &control.GOTO{}, nil
	case opcodes.IFNULL:
		return &control.IFNULL{}, nil
	case opcodes.IFNONNULL:
		return &control.IFNONNULL{}, nil

	// return instructions
	case opcodes.IRETURN:
		return ireturn, nil
	case opcodes.LRETURN:
		return lreturn, nil
	case opcodes.FRETURN:
		return freturn, nil
	case opcodes.DRETURN:
		return dreturn, nil
	case opcodes.ARETURN:
		return areturn, nil
	case opcodes.RETURN:
		return _return, nil

	// References Instructions

	// Object.java Creation / array
	case opcodes.NEW:
		return &references.NEW{}, nil
	case opcodes.NEWARRAY:
		return &arrays.NEWARRAY{}, nil
	case opcodes.ANEWARRAY:
		return &arrays.ANEWARRAY{}, nil
	case opcodes.MULTIANEWARRAY:
		return &arrays.MULTIANEWARRAY{}, nil
	case opcodes.ARRAYLENGTH:
		return arraylength, nil
	case opcodes.INSTANCEOF:
		return &references.INSTANCEOF{}, nil
	case opcodes.CHECKCAST:
		return &references.CHECKCAST{}, nil

	// Field Access
	case opcodes.GETSTATIC:
		return &references.GETSTATIC{}, nil
	case opcodes.PUTSTATIC:
		return &references.PUTSTATIC{}, nil
	case opcodes.GETFIELD:
		return &references.GETFIELD{}, nil
	case opcodes.PUTFIELD:
		return &references.PUTFIELD{}, nil

	// invoke
	case opcodes.INVOKESTATIC:
		return &references.INVOKE_STATIC{}, nil
	case opcodes.INVOKEVIRTUAL:
		return &references.INVOKEVIRTUAL{}, nil
	case opcodes.INVOKESPECIAL:
		return &references.INVOKESPECIAL{}, nil
	case opcodes.INVOKEINTERFACE:
		return &references.INVOKEINTERFACE{}, nil

	default:
		return nop, fmt.Errorf("unsupported opcodes: 0x%02X", opcode)
	}
}
