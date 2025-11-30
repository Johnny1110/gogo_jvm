# JVM Runtime Data Area - HEAP

<br>

---

<br>

這是 Java 運行時區域的 HEAP 區域，這個區域的資料是 Thread 共享的。

所有 Java 物件實例都在 Heap 上分配記憶體。

<br>

## Object 物件結構

### 設計原理

每個 Java 物件在 JVM 中都表示為一個 `Object` 結構：

```
┌─────────────────────────────────────────┐
│ Object                                  │
├─────────────────────────────────────────┤
│ class ────────────────────────► Class   │  物件的類別元資訊
├─────────────────────────────────────────┤
│ fields (Slots)                          │  物件的實例欄位
│ ┌──────┬──────┬──────┬──────┬──────┐   │
│ │ [0]  │ [1]  │ [2]  │ [3]  │ ...  │   │
│ └──────┴──────┴──────┴──────┴──────┘   │
├─────────────────────────────────────────┤
│ extra ─────────────────► interface{}    │  特殊用途（陣列元素等）
└─────────────────────────────────────────┘
```

<br>

### 為什麼 class 使用 interface{} 類型？

這是為了避免 Go 套件的循環依賴問題：

```
如果直接使用 *method_area.Class：
  heap/object.go  ──import──►  method_area/
  method_area/    ──import──►  heap/object.go  （如果 Class 要建立 Object）
  
  這會造成循環依賴，Go 編譯器不允許！

解決方案：
  heap/object.go 只依賴 rtcore 套件
  class 欄位使用 interface{} 類型
  使用時再做類型斷言：obj.Class().(*method_area.Class)
```

這也符合真實 JVM 的設計思想：物件頭（Object Header）只存放類型指標，不包含類型的完整資訊。

<br>

### 繼承情況下的記憶體佈局

```java
class Animal { 
    int age; 
}

class Dog extends Animal { 
    String name; 
}
```

當建立 `Dog` 物件時，fields 的佈局如下：

```
Dog 物件的 fields:
┌─────────────────┬─────────────────┐
│ [0] age (int)   │ [1] name (ref)  │
│   ← 父類欄位     │   ← 子類欄位     │
└─────────────────┴─────────────────┘

slotId 在 ClassLoader.prepare() 階段計算：
  Animal.age.slotId = 0
  Dog.name.slotId = 1
```

這個設計讓子類物件可以直接用父類的 slotId 存取父類欄位，支援多型。

<br>

### extra 欄位的用途

`extra` 是一個預留欄位，用於特殊類型的物件：

| 物件類型 | extra 內容 |
|---------|-----------|
| 一般物件 | nil |
| int[] 陣列 | []int32 |
| long[] 陣列 | []int64 |
| Object[] 陣列 | []*Object |
| String 物件 | Go string（可選優化） |
| Class 物件 | *method_area.Class（反射用） |

這個設計讓我們可以在不修改基本結構的情況下，支援陣列和特殊物件。

<br>

### 與其他元件的關係

```
┌────────────────────────────────────────────────────────────┐
│                         JVM 架構                            │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Method Area (方法區)              Heap (堆)                │
│  ┌─────────────────────┐          ┌─────────────────────┐  │
│  │ Class               │◄─────────│ Object              │  │
│  │ ├─ fields metadata  │  class   │ ├─ class            │  │
│  │ ├─ methods          │  引用     │ ├─ fields (Slots)   │  │
│  │ └─ staticVars       │          │ └─ extra            │  │
│  └─────────────────────┘          └─────────────────────┘  │
│           │                                ▲               │
│           │ Field.slotId                   │               │
│           ▼                                │               │
│  ┌─────────────────────┐                   │               │
│  │ Field               │───────────────────┘               │
│  │ ├─ slotId           │  用於定位欄位在 Object.fields       │
│  │ └─ descriptor       │  中的位置                          │
│  └─────────────────────┘                                   │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

<br>

## 物件建立流程

當執行 `new Counter()` 時：

```
字節碼:
  new #2              // 1. 分配記憶體
  dup                 // 2. 複製引用（給建構子用）
  invokespecial #3    // 3. 呼叫 <init> 建構子
  astore_1            // 4. 存到局部變數

步驟詳解:
1. new 指令:
   - 從常量池取得 ClassRef
   - 解析並載入類別
   - 計算需要的 slot 數量
   - 呼叫 NewObject(class, slotCount)
   - push 引用到操作數棧

2. dup 指令:
   - 複製棧頂的引用
   - 因為 invokespecial 會消耗一個引用

3. invokespecial 指令:
   - 呼叫建構子 <init>
   - 初始化物件欄位

4. astore_1 指令:
   - 將引用存到局部變數表
```

<br>

## 未來擴展

### Phase v0.2.6：陣列支援

陣列也是物件，但有特殊結構：

```go
// 陣列物件使用 extra 存放元素
arr := NewObject(intArrayClass, 0)  // fields 不需要
arr.SetExtra(make([]int32, length)) // 元素存在 extra
```

### GC 支援預留

Object 結構已經為未來的垃圾回收預留了擴展空間：

```go
// 未來可能新增的欄位
type Object struct {
    class    interface{}
    fields   rtcore.Slots
    extra    interface{}
    // markBit  uint8       // GC 標記位
    // hashCode int32       // 物件 hash code 快取
    // monitor  *Monitor    // synchronized 鎖
}
```

<br>