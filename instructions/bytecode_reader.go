package instructions

type BytecodeReader struct {
	code []byte // bytecode (from Code Attribute)
	pc   int    // program counter
}

// Reset reader
// when exec new func or redirect
func (br *BytecodeReader) Reset(code []byte, pc int) {
	br.code = code
	br.pc = pc
}

// PC return current pc
func (br *BytecodeReader) PC() int {
	return br.pc
}

// SetPC set PC (using for JUMP)
func (br *BytecodeReader) SetPC(pc int) {
	br.pc = pc
}

// ReadUint8 read 1 byte (unsign)
// usage: read opcode, bipush...
func (br *BytecodeReader) ReadUint8() uint8 {
	b := br.code[br.pc]
	br.pc++
	return b
}

// ReadInt8 read 1 byte (sign)
func (br *BytecodeReader) ReadInt8() int8 {
	return int8(br.ReadUint8())
}

// ReadUint16 read 2 bytes (unsign Big-Endian)
// highBits on the first, lowBits is sec.
func (br *BytecodeReader) ReadUint16() uint16 {
	highBits := uint16(br.ReadUint8())
	lowBits := uint16(br.ReadUint8())
	return (highBits << 8) | lowBits
}

// ReadInt16 read 2 bytes
func (br *BytecodeReader) ReadInt16() int16 {
	return int16(br.ReadUint16())
}

// ReadInt32 read 4 bytes (sign)
func (br *BytecodeReader) ReadInt32() int32 {
	highBits := int32(br.ReadUint16())
	lowBits := int32(br.ReadUint16())
	return (highBits << 16) | lowBits
}

// ReadInt32s read many int32
// count: read count
func (br *BytecodeReader) ReadInt32s(count int32) []int32 {
	result := make([]int32, count)
	for i := int32(0); i < count; i++ {
		result[i] = br.ReadInt32()
	}
	return result
}

// SkipPadding skip padding bytes
// tableswitch and lookupswitch need align 4 bytes, this is for efficiency (access aligned data is fast)
// ex: if opcode at index 5, we need skip 3 bytes, start from index 8.
func (br *BytecodeReader) SkipPadding() {
	br.pc += (4 - br.pc%4) % 4
}
