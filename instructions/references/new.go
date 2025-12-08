package references

import (
	"github.com/Johnny1110/gogo_jvm/instructions/base"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// NEW Create Object

// Progress:
// 1. get ClassRef from CP
// 2. resolve ClassRef get class
// 3. check class (not interface or abs)
// 4. make sure already inited
// 5. class.NewObject()
// 6. push object's ref into stack
type NEW struct {
	base.Index16Instruction
}

func (n *NEW) Execute(frame *runtime.Frame) {
	cp := frame.Method().Class().ConstantPool()

	// 1. get ClassRef from CP
	classRef := cp.GetConstant(n.Index).(*method_area.ClassRef)

	// 2. resolve ClassRef get class
	class := classRef.ResolvedClass()

	// 3. check class (not interface or abs)
	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}

	// 4. make sure already inited
	// 如果類別還沒初始化，需要先執行 <clinit>
	if !class.InitStarted() {
		// revert PC, rerun this new inst
		frame.RevertNextPC()
		// start init class
		InitClass(frame.Thread(), class)
		return
	}

	// 5. class.NewObject()
	object := class.NewObject()

	// 6. push object's ref into stack
	frame.OperandStack().PushRef(object)
}

func (n *NEW) Opcode() uint8 {
	return 0xBB
}

// ============================================================
// Class Init
// ============================================================

// scheduleClinit settle <clinit> (if exist) method to a new Frame
func scheduleClinit(thread *runtime.Thread, class *method_area.Class) {
	clinit := class.GetClinitMethod()
	if clinit != nil {
		newFrame := thread.NewFrameWithMethod(clinit)
		thread.PushFrame(newFrame)
	}
}

// initSuperClass
// JVM standards: before init class, all parent should be init
func initSuperClass(thread *runtime.Thread, class *method_area.Class) {
	if class.IsInterface() { // skip interface
		return
	}

	superClass := class.SuperClass()
	if superClass != nil && !superClass.InitStarted() {
		InitClass(thread, superClass)
	}
}
