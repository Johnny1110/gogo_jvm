package main

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/interpreter"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"os"
)

// ============================================================
// Gogo JVM v0.2.3
// ============================================================
func main() {
	if len(os.Args) < 2 {
		// first arg is java
		printUsage()
		os.Exit(1)
	}
	classFilePath := os.Args[1]
	debug := len(os.Args) > 2 && os.Args[2] == "-debug"

	// get class path (.class file's dir）
	classPath := getClassPath(classFilePath)
	className := getClassName(classFilePath)

	if debug {
		fmt.Println("============================================")
		fmt.Printf("ClassPath: %s\n", classPath)
		fmt.Printf("ClassName: %s\n", className)
		fmt.Println("============================================")
	}

	// create ClassLoader
	loader := method_area.NewClassLoader(classPath)

	// let ClassLoader load class
	class := loader.LoadClass(className)

	// find main()
	mainMethod := class.GetMainMethod()
	if mainMethod == nil {
		fmt.Println("Error: No main method found")
		fmt.Println("Main method signature must be: public static void main(String[] args)")
		os.Exit(1)
	}

	// start run
	interpreter.Interpret(mainMethod, debug)

	fmt.Println("GOGO JVM exit")
}

// getClassPath
func getClassPath(filePath string) string {
	// find last '/' position
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			return filePath[:i]
		}
	}

	// if input path is like "TestLoopSum.class", the return class path will be "."
	return "."
}

// getClassName
func getClassName(filePath string) string {
	// find dir name
	start := 0
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			start = i + 1
			break
		}
	}
	name := filePath[start:]

	// remove ".class" suffix
	if len(name) > 6 && name[len(name)-6:] == ".class" {
		name = name[:len(name)-6]
	}
	return name
}

func printUsage() {
	fmt.Println("Gogo JVM - A simple JVM implementation in Go")
	fmt.Println()
	fmt.Println("Usage: gogo_jvm <classfile> [-debug]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gogo_jvm SimpleAdd.class")
	fmt.Println("  gogo_jvm SimpleAdd.class -debug")
}
