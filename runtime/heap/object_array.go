package heap

import "fmt"

// extend content for object.go

// =============== Array Constructors ===============
// Java Array is a Object also, different data type using different array type.

// NewByteArray create []byte or []bool
func NewByteArray(class interface{}, length int32) *Object {
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		extra:    make([]int8, length),
	}
}

// NewShortArray create short[] array
func NewShortArray(class interface{}, length int32) *Object {
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		extra:    make([]int16, length),
	}
}

// NewIntArray create int[] array
func NewIntArray(class interface{}, length int32) *Object {
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		extra:    make([]int32, length),
	}
}

// NewLongArray create long[] array
func NewLongArray(class interface{}, length int32) *Object {
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		extra:    make([]int64, length),
	}
}

// NewCharArray create char[] array
// Java char is 16-bit unsigned
func NewCharArray(class interface{}, length int32) *Object {
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		extra:    make([]uint16, length),
	}
}

// NewFloatArray create float[] array
func NewFloatArray(class interface{}, length int32) *Object {
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		extra:    make([]float32, length),
	}
}

// NewDoubleArray create double[] array
func NewDoubleArray(class interface{}, length int32) *Object {
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		extra:    make([]float64, length),
	}
}

// NewRefArray create ref array (Object[], String[], other class[])
func NewRefArray(class interface{}, length int32) *Object {
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		extra:    make([]*Object, length),
	}
}

// =============== Array Check ===============

// IsArray check is array
func (o *Object) IsArray() bool {
	switch o.extra.(type) {
	case []int8, []int16, []uint16, []int32, []int64,
		[]float32, []float64, []*Object:
		return true
	default:
		return false
	}
}

func CheckIndex(arrLen, index int32) {
	if index < 0 || index >= arrLen {
		panic("java.lang.ArrayIndexOutOfBoundsException")
	}
}

func (o *Object) checkArrayIndex(index int32) {
	if !o.IsArray() {
		panic("Object is not an array")
	}
	if index < 0 || index >= o.ArrayLength() {
		panic("java.lang.ArrayIndexOutOfBoundsException")
	}
}

// =============== Array Len ===============

func (o *Object) ArrayLength() int32 {
	switch arr := o.extra.(type) {
	case []int8:
		return int32(len(arr))
	case []int16:
		return int32(len(arr))
	case []int32:
		return int32(len(arr))
	case []int64:
		return int32(len(arr))
	case []uint16: // chars
		return int32(len(arr))
	case []float32:
		return int32(len(arr))
	case []float64:
		return int32(len(arr))
	case []*Object:
		return int32(len(arr))
	default:
		panic(fmt.Sprintf("Can't get length of array of type %T", arr))
	}
}

func (o *Object) ArrayType() string {
	switch arr := o.extra.(type) {
	case []int8:
		return "byte/bool"
	case []int16:
		return "short"
	case []int32:
		return "int"
	case []int64:
		return "long"
	case []uint16: // chars
		return "char"
	case []float32:
		return "float"
	case []float64:
		return "double"
	case []*Object:
		return "object"
	default:
		panic(fmt.Sprintf("Can't get length of array of type %T", arr))
	}
}

// =============== Array Getter/Setter ===============

// Bytes ---------------------------------------------

// Bytes get byte[] ot bool
func (o *Object) Bytes() []int8 {
	return o.extra.([]int8)
}

// GetArrayByte byte[] or bool[] by index
func (o *Object) GetArrayByte(index int32) int8 {
	o.checkArrayIndex(index)
	return o.extra.([]int8)[index]
}

// SetArrayByte set byte/bool by index
func (o *Object) SetArrayByte(index int32, val int8) {
	o.checkArrayIndex(index)
	o.extra.([]int8)[index] = val
}

// Short ---------------------------------------------

// Shorts get short[]
func (o *Object) Shorts() []int16 {
	return o.extra.([]int16)
}

// GetArrayShort get short by index
func (o *Object) GetArrayShort(index int32) int16 {
	o.checkArrayIndex(index)
	return o.extra.([]int16)[index]
}

// SetArrayShort set short by index
func (o *Object) SetArrayShort(index int32, val int16) {
	o.checkArrayIndex(index)
	o.extra.([]int16)[index] = val
}

// Int ---------------------------------------------

// Ints get int[]
func (o *Object) Ints() []int32 {
	return o.extra.([]int32)
}

// GetArrayInt get int by index
func (o *Object) GetArrayInt(index int32) int32 {
	o.checkArrayIndex(index)
	return o.extra.([]int32)[index]
}

// SetArrayInt set int by index
func (o *Object) SetArrayInt(index int32, val int32) {
	o.checkArrayIndex(index)
	o.extra.([]int32)[index] = val
}

// Long ---------------------------------------------

// Longs get long[]
func (o *Object) Longs() []int64 {
	return o.extra.([]int64)
}

// GetArrayLong get long by index
func (o *Object) GetArrayLong(index int32) int64 {
	o.checkArrayIndex(index)
	return o.extra.([]int64)[index]
}

// SetArrayLong set long by index
func (o *Object) SetArrayLong(index int32, val int64) {
	o.checkArrayIndex(index)
	o.extra.([]int64)[index] = val
}

// Char ---------------------------------------------

// Chars get char[]
func (o *Object) Chars() []uint16 {
	return o.extra.([]uint16)
}

// GetArrayChar get char by index
func (o *Object) GetArrayChar(index int32) uint16 {
	o.checkArrayIndex(index)
	return o.extra.([]uint16)[index]
}

// SetArrayChar set char by index
func (o *Object) SetArrayChar(index int32, val uint16) {
	o.checkArrayIndex(index)
	o.extra.([]uint16)[index] = val
}

// Float ---------------------------------------------

// Floats get float[]
func (o *Object) Floats() []float32 {
	return o.extra.([]float32)
}

// GetArrayFloat get float by index
func (o *Object) GetArrayFloat(index int32) float32 {
	o.checkArrayIndex(index)
	return o.extra.([]float32)[index]
}

// SetArrayFloat set float by index
func (o *Object) SetArrayFloat(index int32, val float32) {
	o.checkArrayIndex(index)
	o.extra.([]float32)[index] = val
}

// Double ---------------------------------------------

// Doubles get double[]
func (o *Object) Doubles() []float64 {
	return o.extra.([]float64)
}

// GetArrayDouble get double by index
func (o *Object) GetArrayDouble(index int32) float64 {
	o.checkArrayIndex(index)
	return o.extra.([]float64)[index]
}

// SetArrayDouble set double by index
func (o *Object) SetArrayDouble(index int32, val float64) {
	o.checkArrayIndex(index)
	o.extra.([]float64)[index] = val
}

// Reference ------------------------------------------

// Refs get ref[]
func (o *Object) Refs() []*Object {
	return o.extra.([]*Object)
}

// GetArrayRef get Object by index
func (o *Object) GetArrayRef(index int32) *Object {
	o.checkArrayIndex(index)
	return o.extra.([]*Object)[index]
}

// SetArrayRef set Object by index
func (o *Object) SetArrayRef(index int32, ref *Object) {
	o.checkArrayIndex(index)
	o.extra.([]*Object)[index] = ref
}
