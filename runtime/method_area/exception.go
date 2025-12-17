package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

// ============================================================
// Exception Table - Runtime Exception Handler
// ============================================================
// v0.2.10 - handle exception
//
// Document: https://reurl.cc/dqq5yy
//
// catch_type = 0 (finally)

// ExceptionHandler runtime exception handler
type ExceptionHandler struct {
	StartPC       int
	EndPC         int
	HandlerPC     int
	CatchType     *ClassRef // could be nil (finally)
	CatchTypeName string    // not required in real JVM, this is for gogo JVM (no real ex class)
}

// ExceptionTable from method
type ExceptionTable []*ExceptionHandler

// newExceptionTable create ExceptionTable from classfile (map classfile to runtime)
// args:
// - cfExceptionTable: ExceptionHandler from classfile
// - rtcp: runtime constant pool (for parse catch_type)
func newExceptionTable(cfExceptionTable []*classfile.ExceptionHandler, rtcp *RuntimeConstantPool) ExceptionTable {
	if len(cfExceptionTable) == 0 {
		return nil
	}

	table := make(ExceptionTable, len(cfExceptionTable))
	for idx, cfHandler := range cfExceptionTable {
		table[idx] = &ExceptionHandler{
			StartPC:       int(cfHandler.StartPc()),
			EndPC:         int(cfHandler.EndPc()),
			HandlerPC:     int(cfHandler.HandlerPc()),
			CatchType:     getCatchType(cfHandler.CatchType(), rtcp),
			CatchTypeName: getCatchTypeName(cfHandler.CatchType(), rtcp),
		}
	}

	return table
}

// getCatchType get catch type from rtcp
// return *ClassRef: class ref (nil is finally)
func getCatchType(catchTypeIndex uint16, rtcp *RuntimeConstantPool) *ClassRef {
	if catchTypeIndex == 0 {
		// finally (0)
		return nil
	}

	return rtcp.GetConstant(uint(catchTypeIndex)).(*ClassRef)
}

// getCatchTypeName get catch ex class name
func getCatchTypeName(catchTypeIndex uint16, cp *RuntimeConstantPool) string {
	if catchTypeIndex == 0 {
		return "" // finally
	}
	classRef := cp.GetConstant(uint(catchTypeIndex)).(*ClassRef)
	return classRef.ClassName()
}

// FindExceptionHandler find matched ex handler
// args:
// - exClass: exObject class
// - pc: current pc (where error happened)
//
// return:
// - handlerPC: return handler position (pc) if found, otherwise return -1
func (table ExceptionTable) FindExceptionHandler(exClass *Class, pc int) int {
	for _, handler := range table {
		// rule-1: pc should between startPC & endPC
		if pc >= handler.StartPC && pc < handler.EndPC {
			// rule-2: check match type:
			if handler.CatchType == nil { // finally (0)
				return handler.HandlerPC
			}

			catchClass := handler.CatchType.ResolvedClass()

			if catchClass.IsAssignableFrom(exClass) {
				// catchCass must be parented of exClass
				// - every catch {...} section will auto `goto` finally {...} eventually, so don't worry about only return catch entry PC.
				return handler.HandlerPC
			}
		}
	}
	return -1 // not found handler
}
