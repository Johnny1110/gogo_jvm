package instructions

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/instructions/constants"
	"github.com/Johnny1110/gogo_jvm/instructions/control"
	"github.com/Johnny1110/gogo_jvm/instructions/ipush"
	"github.com/Johnny1110/gogo_jvm/instructions/loads"
	"github.com/Johnny1110/gogo_jvm/instructions/math"
	"github.com/Johnny1110/gogo_jvm/instructions/stores"
)

// ============================================================
// instructions factory（decoder）
// ============================================================

// !!: 預先建立好的 instructions 都是無狀態的，不可以把帶狀態 (index, offset) 的指令建立單例
var (
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
	iload_0     = &loads.ILOAD_0{}
	iload_1     = &loads.ILOAD_1{}
	iload_2     = &loads.ILOAD_2{}
	iload_3     = &loads.ILOAD_3{}
	lload_0     = &loads.LLOAD_0{}
	lload_1     = &loads.LLOAD_1{}
	lload_2     = &loads.LLOAD_2{}
	lload_3     = &loads.LLOAD_3{}
	fload_0     = &loads.FLOAD_0{}
	fload_1     = &loads.FLOAD_1{}
	fload_2     = &loads.FLOAD_2{}
	fload_3     = &loads.FLOAD_3{}
	dload_0     = &loads.DLOAD_0{}
	dload_1     = &loads.DLOAD_1{}
	dload_2     = &loads.DLOAD_2{}
	dload_3     = &loads.DLOAD_3{}
	aload_0     = &loads.ALOAD_0{}
	aload_1     = &loads.ALOAD_1{}
	aload_2     = &loads.ALOAD_2{}
	aload_3     = &loads.ALOAD_3{}
	istore_0    = &stores.ISTORE_0{}
	istore_1    = &stores.ISTORE_1{}
	istore_2    = &stores.ISTORE_2{}
	istore_3    = &stores.ISTORE_3{}
	lstore_0    = &stores.LSTORE_0{}
	lstore_1    = &stores.LSTORE_1{}
	lstore_2    = &stores.LSTORE_2{}
	lstore_3    = &stores.LSTORE_3{}
	fstore_0    = &stores.FSTORE_0{}
	fstore_1    = &stores.FSTORE_1{}
	fstore_2    = &stores.FSTORE_2{}
	fstore_3    = &stores.FSTORE_3{}
	dstore_0    = &stores.DSTORE_0{}
	dstore_1    = &stores.DSTORE_1{}
	dstore_2    = &stores.DSTORE_2{}
	dstore_3    = &stores.DSTORE_3{}
	astore_0    = &stores.ASTORE_0{}
	astore_1    = &stores.ASTORE_1{}
	astore_2    = &stores.ASTORE_2{}
	astore_3    = &stores.ASTORE_3{}
	iadd        = &math.IADD{}
	ladd        = &math.LADD{}
	fadd        = &math.FADD{}
	dadd        = &math.DADD{}
	isub        = &math.ISUB{}
	lsub        = &math.LSUB{}
	fsub        = &math.FSUB{}
	dsub        = &math.DSUB{}
	imul        = &math.IMUL{}
	lmul        = &math.LMUL{}
	fmul        = &math.FMUL{}
	dmul        = &math.DMUL{}
	idiv        = &math.IDIV{}
	ldiv        = &math.LDIV{}
	fdiv        = &math.FDIV{}
	ddiv        = &math.DDIV{}
	irem        = &math.IREM{}
	lrem        = &math.LREM{}
	frem        = &math.FREM{}
	drem        = &math.DREM{}
	ineg        = &math.INEG{}
	lneg        = &math.LNEG{}
	fneg        = &math.FNEG{}
	dneg        = &math.DNEG{}
	ireturn     = &control.IRETURN{}
	lreturn     = &control.LRETURN{}
	freturn     = &control.FRETURN{}
	dreturn     = &control.DRETURN{}
	areturn     = &control.ARETURN{}
	_return     = &control.RETURN{}
)

// NewInstruction return instruction based on input opcode
func NewInstruction(opcode byte) base.Instruction {
	switch opcode {
	// const instructions
	case base.NOP:
		return nop
	case base.ACONST_NULL:
		return aconst_null
	case base.ICONST_M1:
		return iconst_m1
	case base.ICONST_0:
		return iconst_0
	case base.ICONST_1:
		return iconst_1
	case base.ICONST_2:
		return iconst_2
	case base.ICONST_3:
		return iconst_3
	case base.ICONST_4:
		return iconst_4
	case base.ICONST_5:
		return iconst_5
	case base.LCONST_0:
		return lconst_0
	case base.LCONST_1:
		return lconst_1
	case base.FCONST_0:
		return fconst_0
	case base.FCONST_1:
		return fconst_1
	case base.FCONST_2:
		return fconst_2
	case base.DCONST_0:
		return dconst_0
	case base.DCONST_1:
		return dconst_1
	case base.BIPUSH:
		return &ipush.BIPUSH{}
	case base.SIPUSH:
		return &ipush.SIPUSH{}

	// load instructions
	case base.ILOAD:
		return &loads.ILOAD{}
	case base.LLOAD:
		return &loads.LLOAD{}
	case base.FLOAD:
		return &loads.FLOAD{}
	case base.DLOAD:
		return &loads.DLOAD{}
	case base.ALOAD:
		return &loads.ALOAD{}
	case base.ILOAD_0:
		return iload_0
	case base.ILOAD_1:
		return iload_1
	case base.ILOAD_2:
		return iload_2
	case base.ILOAD_3:
		return iload_3
	case base.LLOAD_0:
		return lload_0
	case base.LLOAD_1:
		return lload_1
	case base.LLOAD_2:
		return lload_2
	case base.LLOAD_3:
		return lload_3
	case base.FLOAD_0:
		return fload_0
	case base.FLOAD_1:
		return fload_1
	case base.FLOAD_2:
		return fload_2
	case base.FLOAD_3:
		return fload_3
	case base.DLOAD_0:
		return dload_0
	case base.DLOAD_1:
		return dload_1
	case base.DLOAD_2:
		return dload_2
	case base.DLOAD_3:
		return dload_3
	case base.ALOAD_0:
		return aload_0
	case base.ALOAD_1:
		return aload_1
	case base.ALOAD_2:
		return aload_2
	case base.ALOAD_3:
		return aload_3

	// store instructions
	case base.ISTORE:
		return &stores.ISTORE{}
	case base.LSTORE:
		return &stores.LSTORE{}
	case base.FSTORE:
		return &stores.FSTORE{}
	case base.DSTORE:
		return &stores.DSTORE{}
	case base.ASTORE:
		return &stores.ASTORE{}
	case base.ISTORE_0:
		return istore_0
	case base.ISTORE_1:
		return istore_1
	case base.ISTORE_2:
		return istore_2
	case base.ISTORE_3:
		return istore_3
	case base.LSTORE_0:
		return lstore_0
	case base.LSTORE_1:
		return lstore_1
	case base.LSTORE_2:
		return lstore_2
	case base.LSTORE_3:
		return lstore_3
	case base.FSTORE_0:
		return fstore_0
	case base.FSTORE_1:
		return fstore_1
	case base.FSTORE_2:
		return fstore_2
	case base.FSTORE_3:
		return fstore_3
	case base.DSTORE_0:
		return dstore_0
	case base.DSTORE_1:
		return dstore_1
	case base.DSTORE_2:
		return dstore_2
	case base.DSTORE_3:
		return dstore_3
	case base.ASTORE_0:
		return astore_0
	case base.ASTORE_1:
		return astore_1
	case base.ASTORE_2:
		return astore_2
	case base.ASTORE_3:
		return astore_3

	// // math instructions
	case base.IADD:
		return iadd
	case base.LADD:
		return ladd
	case base.FADD:
		return fadd
	case base.DADD:
		return dadd
	case base.ISUB:
		return isub
	case base.LSUB:
		return lsub
	case base.FSUB:
		return fsub
	case base.DSUB:
		return dsub
	case base.IMUL:
		return imul
	case base.LMUL:
		return lmul
	case base.FMUL:
		return fmul
	case base.DMUL:
		return dmul
	case base.IDIV:
		return idiv
	case base.LDIV:
		return ldiv
	case base.FDIV:
		return fdiv
	case base.DDIV:
		return ddiv
	case base.IREM:
		return irem
	case base.LREM:
		return lrem
	case base.FREM:
		return frem
	case base.DREM:
		return drem
	case base.INEG:
		return ineg
	case base.LNEG:
		return lneg
	case base.FNEG:
		return fneg
	case base.DNEG:
		return dneg
	case base.IINC:
		return &math.IINC{}

	// compare instructions
	case base.IFEQ:
		return &control.IFEQ{}
	case base.IFNE:
		return &control.IFNE{}
	case base.IFLT:
		return &control.IFLT{}
	case base.IFGE:
		return &control.IFGE{}
	case base.IFGT:
		return &control.IFGT{}
	case base.IFLE:
		return &control.IFLE{}
	case base.IF_ICMPEQ:
		return &control.IF_ICMPEQ{}
	case base.IF_ICMPNE:
		return &control.IF_ICMPNE{}
	case base.IF_ICMPLT:
		return &control.IF_ICMPLT{}
	case base.IF_ICMPGE:
		return &control.IF_ICMPGE{}
	case base.IF_ICMPGT:
		return &control.IF_ICMPGT{}
	case base.IF_ICMPLE:
		return &control.IF_ICMPLE{}
	case base.IF_ACMPEQ:
		return &control.IF_ACMPEQ{}
	case base.IF_ACMPNE:
		return &control.IF_ACMPNE{}
	case base.GOTO:
		return &control.GOTO{}
	case base.IFNULL:
		return &control.IFNULL{}
	case base.IFNONNULL:
		return &control.IFNONNULL{}

	// return instructions
	case base.IRETURN:
		return ireturn
	case base.LRETURN:
		return lreturn
	case base.FRETURN:
		return freturn
	case base.DRETURN:
		return dreturn
	case base.ARETURN:
		return areturn
	case base.RETURN:
		return _return

	default:
		panic(fmt.Sprintf("Unsupported opcode: 0x%02X", opcode))
	}
}
