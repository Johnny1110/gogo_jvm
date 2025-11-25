package classfile

import "fmt"

// getUtf8 get UTF8 String from ConstantPool
func getUtf8(cp ConstantPool, index uint16) string {
	if utf8Info, ok := cp[index].(*ConstantUtf8Info); ok {
		return utf8Info.str
	}
	panic(fmt.Sprintf("Wrong constant pool index: %d", index))
}
