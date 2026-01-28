package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// INVOKE_STATIC
// opcode = 0xB8
type INVOKE_STATIC struct {
	base.Index16Instruction
}

func (i *INVOKE_STATIC) Execute(frame *runtime.Frame) {
	// 1. get RuntimeConstantPool from current frame
	cp := frame.Method().Class().ConstantPool()

	// 2. get method reference, index is target methodRef index which is already loaded in RuntimeConstantPool.
	methodRef := cp.GetConstant(i.Index).(*method_area.MethodRef)

	// 3. parse method ref, get target method
	resolvedMethod, err := methodRef.ResolvedMethod()
	if err != nil {
		frame.JavaThrow(err)
	}

	// 4. make sure it's a static method
	if !resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	// 5. TODO: 類初始化（<clinit>，如果類還沒初始化，需要先執行 <clinit> MVP 階段暫時跳過

	// 6. call method
	invokeMethod(frame, resolvedMethod)
}

func (is *INVOKE_STATIC) Opcode() uint8 {
	return 0xB8
}
