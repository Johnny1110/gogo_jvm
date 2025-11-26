package runtime

// Slot is JVM runtime basic data unit
// JVM standard: local vars table and op stack all made by Slot.
// each Slot can store a 32-bit data
// 64-bit data (long. double), should using 2 Slot
type Slot struct {
	// num store number type data (int32, float32)
	// for long or double, we split to 2 slots.
	num int32

	// ref store reference type (pointer)
	// GC need to know the object references, to iterate
	// TODO: we should replace interface{} to *Object
	ref interface{}
}

// JVM defined long/double, should take slot[n] and slot[n+1]
// and also should read slot[n] and slot[n+1]
