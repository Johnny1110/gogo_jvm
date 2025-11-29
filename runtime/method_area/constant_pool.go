package method_area

import (
	"github.com/Johnny1110/gogo_jvm/classfile"
)

// ConstantPool runtime constant pool
//
// ClassFile's ConstantPool is static symbol ref (just string stuff)
// Runtime's ConstantPool convert static symbol ref to direct ref (to a actual method / field)
//
// symbol ref vs direct ref：
// ┌────────────────────────────────────────────────────────┐
// │  ClassFile ConstantPool（Compile）                      │
// │  #1 Methodref → class=#2, nameAndType=#3               │
// │  #2 Class → name=#4                                    │
// │  #3 NameAndType → name=#5, desc=#6                    │
// │  #4 Utf8 → "Calculator"                               │
// │  #5 Utf8 → "add"                                      │
// │  #6 Utf8 → "(II)I"                                    │
// └────────────────────────────────────────────────────────┘
//                        ↓ parse
// ┌────────────────────────────────────────────────────────┐
// │  Rumtime Constant Pool                                 │
// │  #1 MethodRef → pointing to Calculator.add()           │
// └────────────────────────────────────────────────────────┘

type Constant interface {
}

type RuntimeConstantPool struct {
	class  *Class
	consts []Constant
}

// newConstantPool create RuntimeConstantPool from classfile.ConstantPool
func newConstantPool(class *Class, cfCp classfile.ConstantPool) *RuntimeConstantPool {
	cpCount := len(cfCp)
	consts := make([]Constant, cpCount)
	rtCp := &RuntimeConstantPool{class: class, consts: consts}

	// like classfile CP, start from index 1
	for i := 1; i < cpCount; i++ {
		cpInfo := cfCp[i]
		switch cpInfo.(type) {
		case *classfile.ConstantIntegerInfo:
			intInfo := cpInfo.(*classfile.ConstantIntegerInfo)
			consts[i] = intInfo.Value()
		case *classfile.ConstantFloatInfo:
			floatInfo := cpInfo.(*classfile.ConstantFloatInfo)
			consts[i] = floatInfo.Value()
		case *classfile.ConstantLongInfo:
			longInfo := cpInfo.(*classfile.ConstantLongInfo)
			consts[i] = longInfo.Value()
			i++ // long take 2 positions
		case *classfile.ConstantDoubleInfo:
			doubleInfo := cpInfo.(*classfile.ConstantDoubleInfo)
			consts[i] = doubleInfo.Value()
			i++ // double take 2 positions
		case *classfile.ConstantStringInfo:
			stringInfo := cpInfo.(*classfile.ConstantStringInfo)
			consts[i] = stringInfo.String()
		case *classfile.ConstantClassInfo:
			classInfo := cpInfo.(*classfile.ConstantClassInfo)
			consts[i] = NewClassRef(rtCp, classInfo.Name())
		case *classfile.ConstantFieldRefInfo:
			fieldRefInfo := cpInfo.(*classfile.ConstantFieldRefInfo)
			consts[i] = NewFieldRef(rtCp, fieldRefInfo)
		case *classfile.ConstantMethodRefInfo:
			methodRefInfo := cpInfo.(*classfile.ConstantMethodRefInfo)
			consts[i] = NewMethodRef(rtCp, methodRefInfo)
		case *classfile.ConstantInterfaceMethodRefInfo:
			methodRefInfo := cpInfo.(*classfile.ConstantInterfaceMethodRefInfo)
			consts[i] = NewInterfaceMethodRef(rtCp, methodRefInfo)
			// Utf8 and NameAndType are not required to put in runtime constant pool, they are used by others.
		}
	}

	return rtCp
}

// GetConstant get constant
func (cp *RuntimeConstantPool) GetConstant(index uint) Constant {
	if c := cp.consts[index]; c != nil {
		return c
	}
	panic("No constant at index")
}

// Class get class
func (cp *RuntimeConstantPool) Class() *Class {
	return cp.class
}
