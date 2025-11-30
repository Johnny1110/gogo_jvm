package classfile

import (
	"fmt"
	"strings"
)

// getUtf8 get UTF8 String from ClassFileConstantPool
func getUtf8(cp ClassFileConstantPool, index uint16) string {
	if utf8Info, ok := cp[index].(*ConstantUtf8Info); ok {
		return utf8Info.str
	}
	panic(fmt.Sprintf("Wrong constants pool index for UTF-8 tag: %d", index))
}

func Debug(cf *ClassFile, printCP bool) {
	fmt.Println("=== CLASS FILE INFORMATION ===")
	fmt.Printf("Version: %d.%d\n", cf.MajorVersion(), cf.MinorVersion())
	fmt.Printf("Class Name: %s\n", cf.ClassName())
	fmt.Printf("Super Class: %s\n", cf.SuperClassName())
	fmt.Printf("Interfaces: %v\n", cf.InterfaceNames())
	fmt.Printf("Constant Pool Size: %v\n", len(cf.ConstantPool())-1)

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

	if printCP {
		printConstantPool(cf.constantPool)
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

func printConstantPool(pool ClassFileConstantPool) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 100))
	fmt.Printf("CONSTANT POOL (Total: %d entries, Index: 1-%d)\n", len(pool)-1, len(pool)-1)
	fmt.Printf("%s\n", strings.Repeat("=", 100))

	// 注意：常量池索引從 1 開始，0 位置是空的
	for idx := 1; idx < len(pool); idx++ {
		constant := pool[idx]
		if constant == nil {
			// Long 和 Double 佔兩個位置，第二個位置是 nil
			fmt.Printf("#%-4d [EMPTY - part of long/double]\n", idx)
			continue
		}

		// 打印基本信息
		fmt.Printf("#%-4d ", idx)

		switch c := constant.(type) {
		case *ConstantUtf8Info:
			printUtf8Info(c)

		case *ConstantIntegerInfo:
			printIntegerInfo(c)

		case *ConstantFloatInfo:
			printFloatInfo(c)

		case *ConstantLongInfo:
			printLongInfo(c)
			idx++ // Long 佔兩個位置

		case *ConstantDoubleInfo:
			printDoubleInfo(c)
			idx++ // Double 佔兩個位置

		case *ConstantClassInfo:
			printClassInfo(c, pool)

		case *ConstantStringInfo:
			printStringInfo(c, pool)

		case *ConstantFieldRefInfo:
			printFieldrefInfo(c, pool, idx)

		case *ConstantMethodRefInfo:
			printMethodrefInfo(c, pool, idx)

		case *ConstantInterfaceMethodRefInfo:
			printInterfaceMethodrefInfo(c, pool, idx)

		case *ConstantNameAndTypeInfo:
			printNameAndTypeInfo(c, pool)

		default:
			fmt.Printf("%-20s [Unknown constants type]\n", "Unknown")
		}
	}

	fmt.Printf("%s\n\n", strings.Repeat("=", 100))
}

// UTF-8 常量
func printUtf8Info(c *ConstantUtf8Info) {
	fmt.Printf("%-20s value=\"%s\"\n", "Utf8", c.str)
}

// 整數常量
func printIntegerInfo(c *ConstantIntegerInfo) {
	fmt.Printf("%-20s value=%d (0x%08X)\n", "Integer", c.val, c.val)
}

// 浮點數常量
func printFloatInfo(c *ConstantFloatInfo) {
	fmt.Printf("%-20s value=%f\n", "Float", c.val)
}

// Long 常量
func printLongInfo(c *ConstantLongInfo) {
	fmt.Printf("%-20s value=%d (0x%016X)\n", "Long", c.val, c.val)
}

// Double 常量
func printDoubleInfo(c *ConstantDoubleInfo) {
	fmt.Printf("%-20s value=%f\n", "Double", c.val)
}

// Class 引用
func printClassInfo(c *ConstantClassInfo, pool ClassFileConstantPool) {
	className := c.Name()
	fmt.Printf("%-20s name_index=#%-3d -> \"%s\"\n",
		"Class", c.nameIndex, className)
}

// String 常量
func printStringInfo(c *ConstantStringInfo, pool ClassFileConstantPool) {
	str := getUtf8(pool, c.stringIndex)
	fmt.Printf("%-20s string_index=#%-3d -> \"%s\"\n",
		"String", c.stringIndex, str)
}

// 字段引用
func printFieldrefInfo(c *ConstantFieldRefInfo, pool ClassFileConstantPool, idx int) {
	className, name, descriptor := pool.getMemberRef(idx)
	fmt.Printf("%-20s class=#%-3d, name_type=#%-3d\n",
		"Fieldref", c.classIndex, c.nameAndTypeIndex)
	fmt.Printf("%-20s └─> %s.%s:%s\n", "", className, name, descriptor)
}

// 方法引用
func printMethodrefInfo(c *ConstantMethodRefInfo, pool ClassFileConstantPool, idx int) {
	className, name, descriptor := pool.getMemberRef(idx)
	fmt.Printf("%-20s class=#%-3d, name_type=#%-3d\n",
		"Methodref", c.classIndex, c.nameAndTypeIndex)
	fmt.Printf("%-20s └─> %s.%s%s\n", "", className, name, descriptor)
}

// 接口方法引用
func printInterfaceMethodrefInfo(c *ConstantInterfaceMethodRefInfo, pool ClassFileConstantPool, idx int) {
	className, name, descriptor := pool.getMemberRef(idx)
	fmt.Printf("%-20s class=#%-3d, name_type=#%-3d\n",
		"InterfaceMethodref", c.classIndex, c.nameAndTypeIndex)
	fmt.Printf("%-20s └─> %s.%s%s\n", "", className, name, descriptor)
}

// 名稱和類型
func printNameAndTypeInfo(c *ConstantNameAndTypeInfo, pool ClassFileConstantPool) {
	name := getUtf8(pool, c.nameIndex)
	descriptor := getUtf8(pool, c.descriptorIndex)
	fmt.Printf("%-20s name=#%-3d, desc=#%-3d\n",
		"NameAndType", c.nameIndex, c.descriptorIndex)
	fmt.Printf("%-20s └─> %s:%s\n", "", name, descriptor)
}

// PrintConstantPoolWithReferences 打印常量池及其引用關係圖
func PrintConstantPoolWithReferences(pool ClassFileConstantPool) {
	printConstantPool(pool)

	// 打印引用關係圖
	fmt.Println("REFERENCE GRAPH:")
	fmt.Println(strings.Repeat("-", 100))

	for idx := 1; idx < len(pool); idx++ {
		constant := pool[idx]
		if constant == nil {
			continue
		}

		var refs []string

		switch c := constant.(type) {
		case *ConstantClassInfo:
			refs = append(refs, fmt.Sprintf("→ #%d (Utf8: %s)", c.nameIndex, c.Name()))

		case *ConstantStringInfo:
			refs = append(refs, fmt.Sprintf("→ #%d (Utf8: %s)", c.stringIndex, getUtf8(pool, c.stringIndex)))

		case *ConstantFieldRefInfo:
			refs = append(refs, fmt.Sprintf("→ #%d (Class)", c.classIndex))
			refs = append(refs, fmt.Sprintf("→ #%d (NameAndType)", c.nameAndTypeIndex))

		case *ConstantMethodRefInfo:
			refs = append(refs, fmt.Sprintf("→ #%d (Class)", c.classIndex))
			refs = append(refs, fmt.Sprintf("→ #%d (NameAndType)", c.nameAndTypeIndex))

		case *ConstantInterfaceMethodRefInfo:
			refs = append(refs, fmt.Sprintf("→ #%d (Class)", c.classIndex))
			refs = append(refs, fmt.Sprintf("→ #%d (NameAndType)", c.nameAndTypeIndex))

		case *ConstantNameAndTypeInfo:
			refs = append(refs, fmt.Sprintf("→ #%d (Utf8: %s)", c.nameIndex, getUtf8(pool, c.nameIndex)))
			refs = append(refs, fmt.Sprintf("→ #%d (Utf8: %s)", c.descriptorIndex, getUtf8(pool, c.descriptorIndex)))
		}

		if len(refs) > 0 {
			fmt.Printf("#%-3d ", idx)
			for i, ref := range refs {
				if i == 0 {
					fmt.Printf("%s\n", ref)
				} else {
					fmt.Printf("     %s\n", ref)
				}
			}
		}
	}

	fmt.Println(strings.Repeat("-", 100))
}

// PrintConstantPoolSummary 打印常量池統計摘要
func PrintConstantPoolSummary(pool ClassFileConstantPool) {
	counts := make(map[string]int)

	for idx := 1; idx < len(pool); idx++ {
		constant := pool[idx]
		if constant == nil {
			continue
		}

		switch constant.(type) {
		case *ConstantUtf8Info:
			counts["Utf8"]++
		case *ConstantIntegerInfo:
			counts["Integer"]++
		case *ConstantFloatInfo:
			counts["Float"]++
		case *ConstantLongInfo:
			counts["Long"]++
		case *ConstantDoubleInfo:
			counts["Double"]++
		case *ConstantClassInfo:
			counts["Class"]++
		case *ConstantStringInfo:
			counts["String"]++
		case *ConstantFieldRefInfo:
			counts["Fieldref"]++
		case *ConstantMethodRefInfo:
			counts["Methodref"]++
		case *ConstantInterfaceMethodRefInfo:
			counts["InterfaceMethodref"]++
		case *ConstantNameAndTypeInfo:
			counts["NameAndType"]++
		}
	}

	fmt.Println("\nCONSTANT POOL SUMMARY:")
	fmt.Println(strings.Repeat("-", 50))

	total := 0
	for typeName, count := range counts {
		fmt.Printf("%-20s: %d\n", typeName, count)
		total += count
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("%-20s: %d\n", "Total Entries", total)
	fmt.Printf("%-20s: %d bytes (estimated)\n", "Memory Usage", estimatePoolSize(pool))
}

// estimatePoolSize 估算常量池佔用的內存大小
func estimatePoolSize(pool ClassFileConstantPool) int {
	size := 0

	for idx := 1; idx < len(pool); idx++ {
		constant := pool[idx]
		if constant == nil {
			continue
		}

		// 基本開銷：標籤（1字節）
		size += 1

		switch c := constant.(type) {
		case *ConstantUtf8Info:
			size += 2 + len(c.str) // length(2) + string bytes
		case *ConstantIntegerInfo:
			size += 4
		case *ConstantFloatInfo:
			size += 4
		case *ConstantLongInfo:
			size += 8
		case *ConstantDoubleInfo:
			size += 8
		case *ConstantClassInfo:
			size += 2 // name_index
		case *ConstantStringInfo:
			size += 2 // string_index
		case *ConstantFieldRefInfo:
			size += 4 // class_index(2) + name_and_type_index(2)
		case *ConstantMethodRefInfo:
			size += 4
		case *ConstantInterfaceMethodRefInfo:
			size += 4
		case *ConstantNameAndTypeInfo:
			size += 4 // name_index(2) + descriptor_index(2)
		}
	}

	return size
}

// GetConstantPoolEntryInfo 獲取指定索引的常量池項詳細信息
func GetConstantPoolEntryInfo(pool ClassFileConstantPool, index int) string {
	if index == 0 || int(index) >= len(pool) {
		return fmt.Sprintf("Invalid index: %d", index)
	}

	constant := pool[index]
	if constant == nil {
		return fmt.Sprintf("#%d: [EMPTY - part of long/double]", index)
	}

	var info strings.Builder
	info.WriteString(fmt.Sprintf("#%d: ", index))

	switch c := constant.(type) {
	case *ConstantUtf8Info:
		info.WriteString(fmt.Sprintf("Utf8[\"%s\"]", c.str))

	case *ConstantIntegerInfo:
		info.WriteString(fmt.Sprintf("Integer[%d]", c.val))

	case *ConstantFloatInfo:
		info.WriteString(fmt.Sprintf("Float[%f]", c.val))

	case *ConstantLongInfo:
		info.WriteString(fmt.Sprintf("Long[%d]", c.val))

	case *ConstantDoubleInfo:
		info.WriteString(fmt.Sprintf("Double[%f]", c.val))

	case *ConstantClassInfo:
		info.WriteString(fmt.Sprintf("Class[name=#%d:\"%s\"]", c.nameIndex, c.Name()))

	case *ConstantStringInfo:
		str := getUtf8(pool, c.stringIndex)
		info.WriteString(fmt.Sprintf("String[#%d:\"%s\"]", c.stringIndex, str))

	case *ConstantFieldRefInfo:
		className, name, descriptor := pool.getMemberRef(index)
		info.WriteString(fmt.Sprintf("Fieldref[%s.%s:%s]", className, name, descriptor))

	case *ConstantMethodRefInfo:
		className, name, descriptor := pool.getMemberRef(index)
		info.WriteString(fmt.Sprintf("Methodref[%s.%s%s]", className, name, descriptor))

	case *ConstantInterfaceMethodRefInfo:
		className, name, descriptor := pool.getMemberRef(index)
		info.WriteString(fmt.Sprintf("InterfaceMethodref[%s.%s%s]", className, name, descriptor))

	case *ConstantNameAndTypeInfo:
		name := getUtf8(pool, c.nameIndex)
		descriptor := getUtf8(pool, c.descriptorIndex)
		info.WriteString(fmt.Sprintf("NameAndType[%s:%s]", name, descriptor))

	default:
		info.WriteString("Unknown")
	}

	return info.String()
}
