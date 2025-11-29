package method_area

// ClassRef 類引用
// 例如：new Calculator → 需要解析 Calculator 類
type ClassRef struct {
	SymRef
}

func NewClassRef(cp *RuntimeConstantPool, className string) *ClassRef {
	ref := &ClassRef{}
	ref.cp = cp
	ref.className = className
	return ref
}

//```
//
//用於 `new`、`checkcast`、`instanceof` 等指令。
//
//---
//
//## 符號引用 vs 直接引用
//
//這是 JVM 規範中的核心概念：
//```
//編譯時（ClassFile）                    運行時（解析後）
//┌─────────────────────────┐           ┌─────────────────────────┐
//│ ConstantPool            │           │ RuntimeConstantPool     │
//│ #1 Methodref            │           │ #1 MethodRef            │
//│    class=#2             │    →      │    method ──────────────┼──→ Method 對象
//│    nameAndType=#3       │           │                         │    (內存地址)
//│ #2 Class                │           │ #2 ClassRef             │
//│    name=#4              │    →      │    class ───────────────┼──→ Class 對象
//│ #3 NameAndType          │           │                         │
//│    name=#5              │           │                         │
//│    descriptor=#6        │           │                         │
//│ #4 Utf8 "Calculator"    │           │                         │
//│ #5 Utf8 "add"           │           │                         │
//│ #6 Utf8 "(II)I"         │           │                         │
//└─────────────────────────┘           └─────────────────────────┘
//
//符號引用：字符串形式的引用（"Calculator", "add", "(II)I"）
//直接引用：內存中的指針（*Method, *Class）
