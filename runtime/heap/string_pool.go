package heap

// ============================================================
// String Pool - String Interning
// ============================================================

// internedStrings Global String Pool
// key: Go string (UTF-8)
// value: Java String Object (contains UTF-16 char[])
var internedStrings = map[string]*Object{}

// ============================================================
// encode transform func
// ============================================================

// utf8ToUtf16 GoString(UTF-8) to JavaString(UTF-16)
// UTF-8  (1-4 bytes per char)
// UTF-16 (fixed 2 bytes per char)
// TODO: 目前不處理超過 U+FFFF 的字元 (emoji)，完整實現需要處理 Surrogate Pair
func utf8ToUtf16(s string) []uint16 {
	// []rune auto handle UTF-8 -> Unicode code point
	runes := []rune(s)
	chars := make([]uint16, len(runes)) // fixed 2 bytes
	for idx, r := range runes {
		chars[idx] = uint16(r)
	}
	return chars
}

// utf16ToUtf8 JavaString(UTF-16) to GoString(UTF-8)
func utf16ToUtf8(chars []uint16) string {
	runes := make([]rune, len(chars))
	for idx, char := range chars {
		runes[idx] = rune(char)
	}
	return string(runes)
}

// ============================================================
// String Object (Create & Access)
// ============================================================
// Java String Object (Real JVM)：
// ┌─────────────────────────────────────┐
// │ String Object                       │
// ├─────────────────────────────────────┤
// │ class → java/lang/String            │
// │ fields:                             │
// │   [0] value → char[] Object         │
// │   [1] hash  → int (cached)          │
// └─────────────────────────────────────┘
//
// MVP Simplify:
// ┌─────────────────────────────────────┐
// │ String Object (Hack)                │
// ├─────────────────────────────────────┤
// │ class → nil (TODO)                  │
// │ fields → nil                        │
// │ extra → char[] Object               │  <- store char[] into extra
// └─────────────────────────────────────┘

// NewJString create java string
// args: goStr: Go string (UTF-8)
// return: *Object -> Java String Object
func NewJString(goStr string) *Object {
	// 1. utf-8 to utf-16
	chars := utf8ToUtf16(goStr)
	charLen := int32(len(chars))
	// 2. create []char Array Object
	charJArr := NewCharArray(nil, charLen)
	// 3. copy chars into charJArr
	copy(charJArr.Chars(), chars)
	// 4. create string object
	strObject := &Object{
		class: nil, // // TODO: in real JVM, class should be java/lang/String
		extra: charJArr,
	}

	return strObject
}

// GoString extract UTF-8 string from a java String Object
// args: Java String Object
// return: Go String (UTF-8)
func GoString(strObject *Object) string {
	if strObject == nil {
		return "null"
	}
	// take extra from char[]
	charArr := strObject.Extra()
	if charArr == nil {
		return "null"
	}
	// get utf16 char[] data
	charArrObj := charArr.(*Object)
	utf16Chars := charArrObj.Chars()
	// UTF-16 -> UTF-8
	return utf16ToUtf8(utf16Chars)
}

// ============================================================
// String Interning
// ============================================================

// InternString string -> internedStrings
// if string already in pool, return existing java string object
// if string not in pool, create new java string and put it into internedStrings
// args: goStr (UTF-8)
// return: Java String Object (from internedStrings pool)
func InternString(goStr string) *Object {
	if internedObj, ok := internedStrings[goStr]; ok {
		// already in pool
		return internedObj
	}

	// create new String Object
	strObj := NewJString(goStr)
	// in pool
	internedStrings[goStr] = strObj

	return strObj
}

// IsJString check object is java string
// check object.Extra() must be char[]
// TODO: in real JVM, should check object.class == &Class -> java.lang.String
func IsJString(obj *Object) bool {
	if obj == nil {
		return false
	}

	extra := obj.Extra() // this should be a java array object
	if extra == nil {
		return false
	}

	if charArr, ok := extra.(*Object); ok {
		_, isCharArray := charArr.Extra().([]uint16)
		// TODO: should also check charArr.class is java.lang.String
		return isCharArray
	}

	return false
}
