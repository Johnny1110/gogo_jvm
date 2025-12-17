package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
	"strings"
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

			catchClass := handler.CatchType.ResolvedClass() // (TODO: v0.2.10: we don't have ex classes, currently using FindExceptionHandlerByClassName instead)

			if catchClass.IsAssignableFrom(exClass) {
				// catchCass must be parented of exClass
				// - every catch {...} section will auto `goto` finally {...} eventually, so don't worry about only return catch entry PC.
				return handler.HandlerPC
			}
		}
	}
	return -1 // not found handler
}

// FindExceptionHandlerByClassName using className find ex handler (TODO: MVP simplify)
// handle no class exception object
// args:
// - exClassName: exception class name ex: "java/lang/ArithmeticException"
// - pc: current pc
//
// return:
// - handlerPC: handler entry, of not found return 01
func (table ExceptionTable) FindExceptionHandlerByClassName(exClassName string, pc int) int {
	for _, handler := range table {
		if pc >= handler.StartPC && pc < handler.EndPC {
			if handler.CatchType == nil {
				return handler.HandlerPC // finally
			}

			// TODO: using classname for matching
			if isExceptionAssignable(exClassName, handler.CatchTypeName) {
				return handler.HandlerPC
			}
		}
	}

	return -1
}

// isExceptionAssignable
func isExceptionAssignable(exClassName, catchTypeName string) bool {
	// fully match
	if exClassName == catchTypeName {
		return true
	}

	// remove prefix "java/lang/"
	exSimple := simplifyClassName(exClassName)
	catchSimple := simplifyClassName(catchTypeName)

	// key: son, value: all parents classes
	hierarchy := map[string][]string{
		"ArithmeticException":            {"RuntimeException", "Exception", "Throwable"},
		"NullPointerException":           {"RuntimeException", "Exception", "Throwable"},
		"ArrayIndexOutOfBoundsException": {"IndexOutOfBoundsException", "RuntimeException", "Exception", "Throwable"},
		"IndexOutOfBoundsException":      {"RuntimeException", "Exception", "Throwable"},
		"ClassCastException":             {"RuntimeException", "Exception", "Throwable"},
		"NegativeArraySizeException":     {"RuntimeException", "Exception", "Throwable"},
		"IllegalArgumentException":       {"RuntimeException", "Exception", "Throwable"},
		"RuntimeException":               {"Exception", "Throwable"},
		"Exception":                      {"Throwable"},
		"Error":                          {"Throwable"},
		"Throwable":                      {},
	}

	if parents, ok := hierarchy[exSimple]; ok {
		for _, parent := range parents {
			if parent == catchSimple {
				return true
			}
		}
	}

	return false
}

// simplifyClassName "java/lang/ArithmeticException" → "ArithmeticException"
func simplifyClassName(className string) string {
	if idx := strings.LastIndex(className, "/"); idx >= 0 {
		return className[idx+1:]
	}
	return className
}
