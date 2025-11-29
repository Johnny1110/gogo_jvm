package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/marea/cpref"
)

// INVOKE_STATIC
// opcode = 0xB8
type INVOKE_STATIC struct {
	base.Index16Instruction
}

func (i *INVOKE_STATIC) Execute(frame *runtime.Frame) {
	// 1. get runtime CP from current frame
	cp := frame.Method().Class().ConstantPool()

	// 2. get method reference
	methodRef := cp.GetConstant(i.Index).(*cpref.MethodRef)

	// 3. parse method ref, get target method
	resolvedMethod := methodRef.ResolvedMethod()

	// 4. make sure it's a static method
	if !resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 5. TODO: 類初始化（<clinit>）
	// 如果類還沒初始化，需要先執行 <clinit>
	// MVP 階段暫時跳過

	// 6. call method
	InvokeMethod(frame, resolvedMethod)
}

func (is *INVOKE_STATIC) Opcode() uint8 {
	return 0xB8
}
