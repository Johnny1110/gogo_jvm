package classfile

import "fmt"

// getUtf8 get UTF8 String from ConstantPool
func getUtf8(cp ConstantPool, index uint16) string {
	if utf8Info, ok := cp[index].(*ConstantUtf8Info); ok {
		return utf8Info.str
	}
	panic(fmt.Sprintf("Wrong constant pool index: %d", index))
}

func Debug(cf *ClassFile) {
	fmt.Println("=== CLASS FILE INFORMATION ===")
	fmt.Printf("Version: %d.%d\n", cf.MajorVersion(), cf.MinorVersion())
	fmt.Printf("Class Name: %s\n", cf.ClassName())
	fmt.Printf("Super Class: %s\n", cf.SuperClassName())
	fmt.Printf("Interfaces: %v\n", cf.InterfaceNames())

	fmt.Printf("Access Flags: 0x%04X", cf.AccessFlags())
	if cf.IsPublic() {
		fmt.Print(" [public]")
	}
	if cf.IsFinal() {
		fmt.Print(" [final]")
	}
	if cf.IsInterface() {
		fmt.Print(" [interface]")
	}
	if cf.IsAbstract() {
		fmt.Print(" [abstract]")
	}

	fmt.Println()

	// print source file name
	if sourceFile := cf.SourceFileAttribute(); sourceFile != nil {
		fmt.Printf("Source File: %s\n", sourceFile.FileName())
	}

	// print fields
	fmt.Printf("\n=== FIELDS (%d) ===\n", len(cf.Fields()))
	for _, field := range cf.Fields() {
		printMemberInfo("Field", field)
	}

	// print methods
	fmt.Printf("\n=== METHODS (%d) ===\n", len(cf.Methods()))
	for _, method := range cf.Methods() {
		printMemberInfo("Method", method)

		// if is main(), print extra info
		if method.IsPSVM() {
			if codeAttr := method.CodeAttribute(); codeAttr != nil {
				fmt.Printf("   - Max Stack: %d\n", codeAttr.MaxStack())
				fmt.Printf("   - Max Locals: %d\n", codeAttr.MaxLocals())
				fmt.Printf("   - Code Length: %d bytes\n", len(codeAttr.Code()))
			}
		}
	}

	// check PSVM exists
	if mainMethod := cf.GetMainMethod(); mainMethod != nil {
		fmt.Println("\n✓ Found main method!")
	} else {
		fmt.Println("\n✗ No main method found")
	}
}

func printMemberInfo(memberType string, member *MemberInfo) {
	fmt.Printf("  %s: %s %s", memberType, member.Name(), member.Descriptor())

	// print access flags
	flags := []string{}
	if member.IsPublic() {
		flags = append(flags, "public")
	}
	if member.IsPrivate() {
		flags = append(flags, "private")
	}
	if member.IsProtected() {
		flags = append(flags, "protected")
	}
	if member.IsStatic() {
		flags = append(flags, "static")
	}
	if member.IsFinal() {
		flags = append(flags, "final")
	}
	if member.IsAbstract() {
		flags = append(flags, "abstract")
	}
	if member.IsNative() {
		flags = append(flags, "native")
	}

	if len(flags) > 0 {
		fmt.Printf(" [%s]", join(flags, ", "))
	}
	fmt.Println()
}

func join(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
