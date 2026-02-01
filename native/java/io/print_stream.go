package io

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
)

// PrintStream.println implementation

// ============================================================
// java/io/PrintStream Native Methods
// ============================================================

// init register all different println overloading
func init() {
	if global.DebugMode() {
		fmt.Println("@@ Debug - init Native java/io/PrintStream")
	}

	runtime.Register("java/io/PrintStream", "println", "()V", println)
	runtime.Register("java/io/PrintStream", "println", "(Z)V", printlnBoolean)
	runtime.Register("java/io/PrintStream", "println", "(C)V", printlnChar)
	runtime.Register("java/io/PrintStream", "println", "(I)V", printlnInt)
	runtime.Register("java/io/PrintStream", "println", "(J)V", printlnLong)
	runtime.Register("java/io/PrintStream", "println", "(F)V", printlnFloat)
	runtime.Register("java/io/PrintStream", "println", "(D)V", printlnDouble)
	// v0.2.9 supported - string print
	runtime.Register("java/io/PrintStream", "println", "(Ljava/lang/String;)V", printlnString)
}

// ============================================================
// println() - only print "\n"
// ============================================================
// Java: System.out.println();
// Descriptor: ()V
//
// LocalVars:
//
//	[0] = this (PrintStream Ref)
func println(frame *runtime.Frame) (ex *heap.Object) {
	fmt.Println()

	return
}

// ============================================================
// println(boolean) - print bool
// ============================================================
// Java: System.out.println(true);
// Descriptor: (Z)V
//
// LocalVars:
//
//	[0] = this (PrintStream Ref)
//	[1] = boolean value (JVM present by int: 0=false, non 0=true)
func printlnBoolean(frame *runtime.Frame) (ex *heap.Object) {
	val := frame.LocalVars().GetInt(1)
	if val != 0 {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}

	return
}

// ============================================================
// println(char) - print char
// ============================================================
// Java: System.out.println('A');
// Descriptor: (C)V
//
// LocalVars:
//
//	[0] = this
//	[1] = char 值 (JVM - int: scope = 0-65535)
//
// Java char is 16-bit unsigned (UTF-16)
func printlnChar(frame *runtime.Frame) (ex *heap.Object) {
	val := frame.LocalVars().GetInt(1)
	fmt.Println(string(rune(val)))

	return
}

// ============================================================
// println(int) - print int
// ============================================================
// Java: System.out.println(42);
// Descriptor: (I)V
//
// LocalVars:
//
//	[0] = this (PrintStream Ref)
//	[1] = int
func printlnInt(frame *runtime.Frame) (ex *heap.Object) {
	val := frame.LocalVars().GetInt(1)
	fmt.Println(val)

	return
}

// ============================================================
// println(long) - print long
// ============================================================
// Java: System.out.println(123456789L);
// Descriptor: (J)V
//
// LocalVars:
//
//	[0] = this
//	[1-2] = long 值（take 2 slots）
func printlnLong(frame *runtime.Frame) (ex *heap.Object) {
	val := frame.LocalVars().GetLong(1)
	fmt.Println(val)

	return
}

// ============================================================
// println(float) - print float
// ============================================================
// Java: System.out.println(3.14f);
// Descriptor: (F)V
//
// LocalVars:
//
//	[0] = this
//	[1] = float value
func printlnFloat(frame *runtime.Frame) (ex *heap.Object) {
	val := frame.LocalVars().GetFloat(1)
	fmt.Println(val)

	return
}

// ============================================================
// println(double) - print double
// ============================================================
// Java: System.out.println(3.14159265358979);
// Descriptor: (D)V
//
// LocalVars:
//
//	[0] = this
//	[1-2] = double 值（take 2 slots）
func printlnDouble(frame *runtime.Frame) (ex *heap.Object) {
	val := frame.LocalVars().GetDouble(1)
	fmt.Println(val)

	return
}

// ============================================================
// println(String) - print String Object
// ============================================================
// Java: System.out.println("Hello");
// Descriptor: (Ljava/lang/String;)V
//
// LocalVars:
//
//	[0] = this (PrintStream Ref)
//	[1] = String Object Ref
//
// String Object Structure (MVP):
// ┌─────────────────────────────┐
// │ String Object               │
// ├─────────────────────────────┤
// │ extra → char[] Object       │
// │           ↓                 │
// │    []uint16 (UTF-16 data)   │
// └─────────────────────────────┘
func printlnString(frame *runtime.Frame) (ex *heap.Object) {
	strRef := frame.LocalVars().GetRef(1)

	if strRef == nil {
		fmt.Println("null")
		return
	}

	strObject := strRef.(*heap.Object)
	goStr := heap.GoString(strObject)

	fmt.Println(goStr)

	return
}
