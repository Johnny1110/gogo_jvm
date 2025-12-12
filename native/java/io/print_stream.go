package io

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/runtime"
)

// PrintStream.println implementation

// ============================================================
// java/io/PrintStream Native Methods
// ============================================================

// init register all different println overloading
func init() {
	fmt.Println("@@ Debug - init Native java/io/PrintStream")
	runtime.Register("java/io/PrintStream", "println", "()V", println)
	runtime.Register("java/io/PrintStream", "println", "(Z)V", printlnBoolean)
	runtime.Register("java/io/PrintStream", "println", "(C)V", printlnChar)
	runtime.Register("java/io/PrintStream", "println", "(I)V", printlnInt)
	runtime.Register("java/io/PrintStream", "println", "(J)V", printlnLong)
	runtime.Register("java/io/PrintStream", "println", "(F)V", printlnFloat)
	runtime.Register("java/io/PrintStream", "println", "(D)V", printlnDouble)
	// TODO: String require String support first
	// native.Register("java/io/PrintStream", "println", "(Ljava/lang/String;)V", printlnString)
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
func println(frame *runtime.Frame) {
	fmt.Println()
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
func printlnBoolean(frame *runtime.Frame) {
	val := frame.LocalVars().GetInt(1)
	if val != 0 {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
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
func printlnChar(frame *runtime.Frame) {
	val := frame.LocalVars().GetInt(1)
	fmt.Println(string(rune(val)))
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
func printlnInt(frame *runtime.Frame) {
	val := frame.LocalVars().GetInt(1)
	fmt.Println(val)
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
func printlnLong(frame *runtime.Frame) {
	val := frame.LocalVars().GetLong(1)
	fmt.Println(val)
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
func printlnFloat(frame *runtime.Frame) {
	val := frame.LocalVars().GetFloat(1)
	fmt.Println(val)
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
func printlnDouble(frame *runtime.Frame) {
	val := frame.LocalVars().GetDouble(1)
	fmt.Println(val)
}

// ============================================================
// TODO: println(String) - 印字串
// ============================================================
// 這個需要等我們實現 String 物件的處理後才能完成
// 目前的挑戰：
//  1. Java String 是一個物件，內部有 char[] 陣列
//  2. 需要從 String 物件中提取 char[]
//  3. 將 char[] 轉換成 Go string
func printlnString(frame *runtime.Frame) {
	//ref := frame.LocalVars().GetRef(1)
	//if ref == nil {
	//	fmt.Println("null")
	//	return
	//}
	//jStr := ref.(*heap.Object)
	//goStr := StringToGoString(jStr) // TODO
	//fmt.Println(goStr)
	fmt.Println("(X) printlnString not support right now.")
}
