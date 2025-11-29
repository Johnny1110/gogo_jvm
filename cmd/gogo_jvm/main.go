package main

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/interpreter"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"io/ioutil"
	"os"
)

// ============================================================
// Gogo JVM - Phase 2.3
// ============================================================
//
// 執行流程：
// 1. 讀取 .class 文件
// 2. 使用 ClassLoader 加載類
// 3. 找到 main 方法
// 4. 執行解釋器

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	classFilePath := os.Args[1]
	debug := len(os.Args) > 2 && os.Args[2] == "-debug"

	// 獲取類路徑（.class 文件所在目錄）
	classPath := getClassPath(classFilePath)
	className := getClassName(classFilePath)

	fmt.Printf("ClassPath: %s\n", classPath)
	fmt.Printf("ClassName: %s\n", className)
	fmt.Println("============================================")

	// 創建 ClassLoader
	loader := method_area.NewClassLoader(classPath)

	// 加載類
	class := loader.LoadClass(className)

	// 找到 main 方法
	mainMethod := class.GetMainMethod()
	if mainMethod == nil {
		fmt.Println("Error: No main method found!")
		fmt.Println("Main method signature must be: public static void main(String[] args)")
		os.Exit(1)
	}

	fmt.Printf("\n=== Executing %s.main() ===\n\n", className)

	// 執行
	interpreter.Interpret(mainMethod, debug)

	fmt.Println("\n=== Execution completed ===")
}

// getClassPath 從文件路徑獲取類路徑
func getClassPath(filePath string) string {
	// 找到最後一個 / 的位置
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			return filePath[:i]
		}
	}
	return "."
}

// getClassName 從文件路徑獲取類名
func getClassName(filePath string) string {
	// 找到文件名
	start := 0
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			start = i + 1
			break
		}
	}
	name := filePath[start:]

	// 移除 .class 後綴
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

// readClassFile 讀取類文件（備用，如果不用 ClassLoader）
func readClassFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
