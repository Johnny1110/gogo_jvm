package math

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// ============================================================
// ADD Series
// ============================================================

// IADD int
// opcode = 0x60
type IADD struct{ base.NoOperandsInstruction }

func (i *IADD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	result := v1 + v2
	stack.PushInt(result)
}

func (i *IADD) Opcode() uint8 {
	return 0x60
}

// LADD long
// opcode = 0x61
type LADD struct{ base.NoOperandsInstruction }

func (l *LADD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	result := v1 + v2
	stack.PushLong(result)
}

func (i *LADD) Opcode() uint8 {
	return 0x61
}

// FADD float
// opcode = 0x62
type FADD struct{ base.NoOperandsInstruction }

func (f *FADD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	result := v1 + v2
	stack.PushFloat(result)
}

func (i *FADD) Opcode() uint8 {
	return 0x62
}

// DADD double
// opcode = 0x63
type DADD struct{ base.NoOperandsInstruction }

func (d *DADD) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopDouble()
	v1 := stack.PopDouble()
	result := v1 + v2
	stack.PushDouble(result)
}

func (i *DADD) Opcode() uint8 {
	return 0x63
}

// ============================================================
// SUB Series
// ============================================================

// ISUB int
// opcode = 0x64
type ISUB struct{ base.NoOperandsInstruction }

func (i *ISUB) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt() // subtrahend
	v1 := stack.PopInt() // minuend
	result := v1 - v2
	stack.PushInt(result)
}

func (i *ISUB) Opcode() uint8 {
	return 0x64
}

// LSUB long
// opcode = 0x65
type LSUB struct{ base.NoOperandsInstruction }

func (l *LSUB) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	result := v1 - v2
	stack.PushLong(result)
}

func (i *LSUB) Opcode() uint8 {
	return 0x65
}

// FSUB float
// opcode = 0x66
type FSUB struct{ base.NoOperandsInstruction }

func (f *FSUB) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	result := v1 - v2
	stack.PushFloat(result)
}

func (i *FSUB) Opcode() uint8 {
	return 0x66
}

// DSUB double
// opcode = 0x67
type DSUB struct{ base.NoOperandsInstruction }

func (d *DSUB) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopDouble()
	v1 := stack.PopDouble()
	result := v1 - v2
	stack.PushDouble(result)
}

func (i *DSUB) Opcode() uint8 {
	return 0x67
}

// ============================================================
// MUL Series
// ============================================================

// IMUL int
// opcode = 0x68
type IMUL struct{ base.NoOperandsInstruction }

func (i *IMUL) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	result := v1 * v2
	stack.PushInt(result)
}

func (i *IMUL) Opcode() uint8 {
	return 0x68
}

// LMUL long
// opcode = 0x69
type LMUL struct{ base.NoOperandsInstruction }

func (l *LMUL) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	result := v1 * v2
	stack.PushLong(result)
}

func (i *LMUL) Opcode() uint8 {
	return 0x69
}

// FMUL float
// opcode = 0x6A
type FMUL struct{ base.NoOperandsInstruction }

func (f *FMUL) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	result := v1 * v2
	stack.PushFloat(result)
}

func (i *FMUL) Opcode() uint8 {
	return 0x6A
}

// DMUL double
// opcode = 0x6B
type DMUL struct{ base.NoOperandsInstruction }

func (d *DMUL) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopDouble()
	v1 := stack.PopDouble()
	result := v1 * v2
	stack.PushDouble(result)
}

func (i *DMUL) Opcode() uint8 {
	return 0x6B
}

// ============================================================
// DIV Series
// ============================================================
// Warning: When performing integer division, if the divisor is 0, an ArithmeticException should be thrown.
// TODO: we using panic instead temporary (MVP Phase)

// IDIV int
// opcode = 0x6C
type IDIV struct{ base.NoOperandsInstruction }

func (i *IDIV) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v2 == 0 {
		panic("java.lang.ArithmeticException: / by zero")
	}
	result := v1 / v2
	stack.PushInt(result)
}

func (i *IDIV) Opcode() uint8 {
	return 0x6C
}

// LDIV long
// opcode = 0x6D
type LDIV struct{ base.NoOperandsInstruction }

func (l *LDIV) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	if v2 == 0 {
		panic("java.lang.ArithmeticException: / by zero")
	}
	result := v1 / v2
	stack.PushLong(result)
}

func (i *LDIV) Opcode() uint8 {
	return 0x6D
}

// FDIV float
// opcode = 0x6E
// Floating-point division does not throw an exception; dividing by 0 yields Infinity or NaN.
type FDIV struct{ base.NoOperandsInstruction }

func (f *FDIV) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	result := v1 / v2
	stack.PushFloat(result)
}

func (i *FDIV) Opcode() uint8 {
	return 0x6E
}

// DDIV double
// opcode = 0x6F
type DDIV struct{ base.NoOperandsInstruction }

func (d *DDIV) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopDouble()
	v1 := stack.PopDouble()
	result := v1 / v2
	stack.PushDouble(result)
}

func (i *DDIV) Opcode() uint8 {
	return 0x6F
}

// ============================================================
// REM Mod Series
// ============================================================

// IREM int
// opcode = 0x70
type IREM struct{ base.NoOperandsInstruction }

func (i *IREM) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	if v2 == 0 {
		panic("java.lang.ArithmeticException: / by zero")
	}
	result := v1 % v2
	stack.PushInt(result)
}

func (i *IREM) Opcode() uint8 {
	return 0x70
}

// LREM long
// opcode = 0x71
type LREM struct{ base.NoOperandsInstruction }

func (l *LREM) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	if v2 == 0 {
		panic("java.lang.ArithmeticException: / by zero")
	}
	result := v1 % v2
	stack.PushLong(result)
}

func (i *LREM) Opcode() uint8 {
	return 0x71
}

// FREM float
// opcode = 0x72
type FREM struct{ base.NoOperandsInstruction }

func (f *FREM) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	// Go 的 float32 沒有直接的 % 運算符
	// 使用 math.Remainder 或自己計算
	result := float32(float64(v1) - float64(int32(v1/v2))*float64(v2))
	stack.PushFloat(result)
}

func (i *FREM) Opcode() uint8 {
	return 0x72
}

// DREM double
// opcode = 0x73
type DREM struct{ base.NoOperandsInstruction }

func (d *DREM) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopDouble()
	v1 := stack.PopDouble()
	result := v1 - float64(int64(v1/v2))*v2
	stack.PushDouble(result)
}

func (i *DREM) Opcode() uint8 {
	return 0x73
}

// ============================================================
// NEG Negative
// ============================================================

// INEG int
// opcode = 0x74
type INEG struct{ base.NoOperandsInstruction }

func (i *INEG) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopInt()
	stack.PushInt(-val)
}

func (i *INEG) Opcode() uint8 {
	return 0x74
}

// LNEG long
// opcode = 0x75
type LNEG struct{ base.NoOperandsInstruction }

func (l *LNEG) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopLong()
	stack.PushLong(-val)
}

func (i *LNEG) Opcode() uint8 {
	return 0x75
}

// FNEG float
// opcode = 0x76
type FNEG struct{ base.NoOperandsInstruction }

func (f *FNEG) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopFloat()
	stack.PushFloat(-val)
}

func (i *FNEG) Opcode() uint8 {
	return 0x76
}

// DNEG double
// opcode = 0x77
type DNEG struct{ base.NoOperandsInstruction }

func (d *DNEG) Execute(frame *runtime.Frame) {
	stack := frame.OperandStack()
	val := stack.PopDouble()
	stack.PushDouble(-val)
}

func (i *DNEG) Opcode() uint8 {
	return 0x77
}

// ============================================================
// IINC increasing val in LocalVars
// ============================================================
// special inst: access LocalVars directly, not using stack.
// use case: for loop i++
// faster than using stack

// IINC add const to a var in LocalVars
// opcode = 0x84
// format: iinc index const
// ex: iinc 1 1  （把局部變量1加1，即 i++）
type IINC struct {
	Index uint  // index of LocalVars
	Const int32 // increasing amt
}

func (i *IINC) FetchOperands(reader *base.BytecodeReader) {
	i.Index = uint(reader.ReadUint8())
	i.Const = int32(reader.ReadInt8())
}

func (i *IINC) Execute(frame *runtime.Frame) {
	localVars := frame.LocalVars()
	val := localVars.GetInt(i.Index) // get val from LocalVars
	val += i.Const                   // val + const
	localVars.SetInt(i.Index, val)
}

func (i *IINC) Opcode() uint8 {
	return 0x84
}
