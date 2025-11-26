package main

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/classfile"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gogo_jvm <path to java classfile>")
		fmt.Println("Example: gogo_jvm HelloWorld.class")
		os.Exit(1)
	}

	className := os.Args[1]

	// read class file
	classData, err := ioutil.ReadFile(className)
	if err != nil {
		fmt.Printf("Error reading classfile: %v\n", err)
		os.Exit(1)
	}

	// parsing class file
	cf, err := classfile.Parse(classData)
	if err != nil {
		fmt.Printf("Failed to parsing classfile: %v\n", err)
		os.Exit(1)
	}

	classfile.Debug(cf, true) // print classfile details
}
