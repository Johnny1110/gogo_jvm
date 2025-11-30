package classfile

// ClassFileConstantPool this is classfile's constant pool not runtime's constant pool
// 注意：常量池索引從 1 開始，而不是 0, 這是 JVM 規範的歷史遺留設計
type ClassFileConstantPool []ConstantInfo

func readConstantPool(reader *ClassReader) ClassFileConstantPool {
	constantPoolCount := int(reader.readU2())

	// why constantPoolCount not constantPoolCount-1?
	// because constantPool index start from 1, 0 is invalid.
	constantPool := make([]ConstantInfo, constantPoolCount)

	for i := 1; i < constantPoolCount; i++ {
		constantPool[i] = readConstantInfo(reader, constantPool)

		// caution: long and double take 2 positions (這是一個歷史遺留問題，為了在 32 位系統上對齊)
		// long 與 double 是 64 bit 需要拆成兩個 32 bit.
		switch constantPool[i].(type) {
		case *ConstantLongInfo, *ConstantDoubleInfo:
			i++
		}
	}

	return constantPool
}

func (cp ClassFileConstantPool) getClassName(index uint16) string {
	if classInfo, ok := cp[index].(*ConstantClassInfo); ok {
		return classInfo.Name()
	}
	panic("Not a class info")
}

func (cp ClassFileConstantPool) getNameAndType(index uint16) (string, string) {
	if ntInfo, ok := cp[index].(*ConstantNameAndTypeInfo); ok {
		name := getUtf8(cp, ntInfo.nameIndex)
		descriptor := getUtf8(cp, ntInfo.descriptorIndex)
		return name, descriptor
	}
	panic("Not a name and type info")
}

func (cp ClassFileConstantPool) getMemberRef(index int) (className, name, descriptor string) {
	// field ref
	if memberRef, ok := cp[index].(*ConstantFieldRefInfo); ok {
		return cp.resolveMemberRef(&memberRef.ConstantMemberRefInfo)
	}
	// method ref
	if memberRef, ok := cp[index].(*ConstantMethodRefInfo); ok {
		return cp.resolveMemberRef(&memberRef.ConstantMemberRefInfo)
	}
	// interface method ref
	if memberRef, ok := cp[index].(*ConstantInterfaceMethodRefInfo); ok {
		return cp.resolveMemberRef(&memberRef.ConstantMemberRefInfo)
	}
	panic("Not a member ref")
}

func (cp ClassFileConstantPool) resolveMemberRef(memberRef *ConstantMemberRefInfo) (className, name, descriptor string) {
	className = cp.getClassName(memberRef.classIndex)
	name, descriptor = cp.getNameAndType(memberRef.nameAndTypeIndex)
	return
}
