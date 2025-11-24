package classfile

import (
	"encoding/binary"
)

// ClassReader encapsulate []byte，provide reading different type of data func
// Class file is binary data, including different length of data type.
// ClassReader provider interface to read those data
type ClassReader struct {
	data []byte
}

// readU1 read a un-sign byte (u1 in JVM spec)
func (r *ClassReader) readU1() uint8 {
	val := r.data[0]
	r.data = r.data[1:]
	return val
}

// readU2 read 2 bytes (u2 in JVM spec)
// Class file usingBig-Endian
func (r *ClassReader) readU2() uint16 {
	val := binary.BigEndian.Uint16(r.data)
	r.data = r.data[2:]
	return val
}

// readU4 read 4 bytes (u4 in JVM spec)
func (r *ClassReader) readU4() uint32 {
	val := binary.BigEndian.Uint32(r.data)
	r.data = r.data[4:]
	return val
}

// readU8 read 8 bytes (u8 in JVM spec)
func (r *ClassReader) readU8() uint64 {
	val := binary.BigEndian.Uint64(r.data)
	r.data = r.data[8:]
	return val
}

// readU2Table read uint16 table
// Class 文件中很多地方都是先有一個計數值，然後是對應數量的項
func (r *ClassReader) readU2Table() []uint16 {
	count := r.readU2()
	table := make([]uint16, count)
	for i := range table {
		table[i] = r.readU2()
	}
	return table
}

// readBytes read certain length of bytes
func (r *ClassReader) readBytes(length uint32) []byte {
	bytes := r.data[:length]
	r.data = r.data[length:]
	return bytes
}

func (r *ClassReader) remaining() int {
	return len(r.data)
}
