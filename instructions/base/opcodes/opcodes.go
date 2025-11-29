package opcodes

// JVM Opcode define
const (
	NOP          = 0x00
	ACONST_NULL  = 0x01
	ICONST_M1    = 0x02
	ICONST_0     = 0x03
	ICONST_1     = 0x04
	ICONST_2     = 0x05
	ICONST_3     = 0x06
	ICONST_4     = 0x07
	ICONST_5     = 0x08
	LCONST_0     = 0x09
	LCONST_1     = 0x0A
	FCONST_0     = 0x0B
	FCONST_1     = 0x0C
	FCONST_2     = 0x0D
	DCONST_0     = 0x0E
	DCONST_1     = 0x0F
	BIPUSH       = 0x10
	SIPUSH       = 0x11
	ILOAD        = 0x15
	LLOAD        = 0x16
	FLOAD        = 0x17
	DLOAD        = 0x18
	ALOAD        = 0x19
	ILOAD_0      = 0x1A
	ILOAD_1      = 0x1B
	ILOAD_2      = 0x1C
	ILOAD_3      = 0x1D
	LLOAD_0      = 0x1E
	LLOAD_1      = 0x1F
	LLOAD_2      = 0x20
	LLOAD_3      = 0x21
	FLOAD_0      = 0x22
	FLOAD_1      = 0x23
	FLOAD_2      = 0x24
	FLOAD_3      = 0x25
	DLOAD_0      = 0x26
	DLOAD_1      = 0x27
	DLOAD_2      = 0x28
	DLOAD_3      = 0x29
	ALOAD_0      = 0x2A
	ALOAD_1      = 0x2B
	ALOAD_2      = 0x2C
	ALOAD_3      = 0x2D
	ISTORE       = 0x36
	LSTORE       = 0x37
	FSTORE       = 0x38
	DSTORE       = 0x39
	ASTORE       = 0x3A
	ISTORE_0     = 0x3B
	ISTORE_1     = 0x3C
	ISTORE_2     = 0x3D
	ISTORE_3     = 0x3E
	LSTORE_0     = 0x3F
	LSTORE_1     = 0x40
	LSTORE_2     = 0x41
	LSTORE_3     = 0x42
	FSTORE_0     = 0x43
	FSTORE_1     = 0x44
	FSTORE_2     = 0x45
	FSTORE_3     = 0x46
	DSTORE_0     = 0x47
	DSTORE_1     = 0x48
	DSTORE_2     = 0x49
	DSTORE_3     = 0x4A
	ASTORE_0     = 0x4B
	ASTORE_1     = 0x4C
	ASTORE_2     = 0x4D
	ASTORE_3     = 0x4E
	POP          = 0x57
	POP2         = 0x58
	DUP          = 0x59
	SWAP         = 0x5F
	IADD         = 0x60
	LADD         = 0x61
	FADD         = 0x62
	DADD         = 0x63
	ISUB         = 0x64
	LSUB         = 0x65
	FSUB         = 0x66
	DSUB         = 0x67
	IMUL         = 0x68
	LMUL         = 0x69
	FMUL         = 0x6A
	DMUL         = 0x6B
	IDIV         = 0x6C
	LDIV         = 0x6D
	FDIV         = 0x6E
	DDIV         = 0x6F
	IREM         = 0x70
	LREM         = 0x71
	FREM         = 0x72
	DREM         = 0x73
	INEG         = 0x74
	LNEG         = 0x75
	FNEG         = 0x76
	DNEG         = 0x77
	IINC         = 0x84
	IFEQ         = 0x99
	IFNE         = 0x9A
	IFLT         = 0x9B
	IFGE         = 0x9C
	IFGT         = 0x9D
	IFLE         = 0x9E
	IF_ICMPEQ    = 0x9F
	IF_ICMPNE    = 0xA0
	IF_ICMPLT    = 0xA1
	IF_ICMPGE    = 0xA2
	IF_ICMPGT    = 0xA3
	IF_ICMPLE    = 0xA4
	IF_ACMPEQ    = 0xA5
	IF_ACMPNE    = 0xA6
	GOTO         = 0xA7
	IRETURN      = 0xAC
	LRETURN      = 0xAD
	FRETURN      = 0xAE
	DRETURN      = 0xAF
	ARETURN      = 0xB0
	RETURN       = 0xB1
	IFNULL       = 0xC6
	IFNONNULL    = 0xC7
	INVOKESTATIC = 0xB8
	LDC          = 0x12
	LDC_W        = 0x13
	LDC2_W       = 0x14
)

// OpcodeNames for debug display
var OpcodeNames = map[uint8]string{
	NOP: "nop", ACONST_NULL: "aconst_null",
	ICONST_M1: "iconst_m1", ICONST_0: "iconst_0", ICONST_1: "iconst_1",
	ICONST_2: "iconst_2", ICONST_3: "iconst_3", ICONST_4: "iconst_4",
	ICONST_5: "iconst_5", LCONST_0: "lconst_0", LCONST_1: "lconst_1",
	FCONST_0: "fconst_0", FCONST_1: "fconst_1", FCONST_2: "fconst_2",
	DCONST_0: "dconst_0", DCONST_1: "dconst_1",
	BIPUSH: "bipush", SIPUSH: "sipush",
	ILOAD: "iload", LLOAD: "lload", FLOAD: "fload", DLOAD: "dload", ALOAD: "aload",
	ILOAD_0: "iload_0", ILOAD_1: "iload_1", ILOAD_2: "iload_2", ILOAD_3: "iload_3",
	ISTORE: "istore", LSTORE: "lstore", FSTORE: "fstore", DSTORE: "dstore", ASTORE: "astore",
	ISTORE_0: "istore_0", ISTORE_1: "istore_1", ISTORE_2: "istore_2", ISTORE_3: "istore_3",
	IADD: "iadd", LADD: "ladd", FADD: "fadd", DADD: "dadd",
	ISUB: "isub", LSUB: "lsub", FSUB: "fsub", DSUB: "dsub",
	IMUL: "imul", LMUL: "lmul", FMUL: "fmul", DMUL: "dmul",
	IDIV: "idiv", LDIV: "ldiv", FDIV: "fdiv", DDIV: "ddiv",
	IREM: "irem", LREM: "lrem", FREM: "frem", DREM: "drem",
	INEG: "ineg", LNEG: "lneg", FNEG: "fneg", DNEG: "dneg",
	IINC: "iinc",
	IFEQ: "ifeq", IFNE: "ifne", IFLT: "iflt", IFGE: "ifge", IFGT: "ifgt", IFLE: "ifle",
	IF_ICMPEQ: "if_icmpeq", IF_ICMPNE: "if_icmpne", IF_ICMPLT: "if_icmplt",
	IF_ICMPGE: "if_icmpge", IF_ICMPGT: "if_icmpgt", IF_ICMPLE: "if_icmple",
	IF_ACMPEQ: "if_acmpeq", IF_ACMPNE: "if_acmpne",
	GOTO:    "goto",
	IRETURN: "ireturn", LRETURN: "lreturn", FRETURN: "freturn",
	DRETURN: "dreturn", ARETURN: "areturn", RETURN: "return",
	IFNULL: "ifnull", IFNONNULL: "ifnonnull",
	INVOKESTATIC: "invokestatic", LDC: "ldc", LDC_W: "ldc_w", LDC2_W: "ldc2_w",
}
