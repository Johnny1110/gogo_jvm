# Class file parsing

<br>

ClassReader、ConstantPool、ConstantTag、ConstantInfo 這四個概念是 JVM 載入 `.class` 檔案時的核心組件。

<br>

### 為什麼需要這些組件？

當我們寫完 Java 程式並編譯後，產生的 `.class` 檔案是一種二進制格式。JVM 要執行這個程式，必須：

1. 讀取這個二進制檔案
2. 解析其中的結構
3. 理解裡面的常量（字串、類別名、方法名等）

這就是這四個組件各自的職責。


### 各組件的角色與關係

```
┌─────────────────────────────────────────────────────────┐
│                    .class 檔案 (二進制)                   │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
                   ┌───────────────┐
                   │  ClassReader  │  ← 負責「讀取」二進制資料
                   └───────────────┘
                           │
                           │ 讀取並解析
                           ▼
                   ┌───────────────┐
                   │ ConstantPool  │  ← 常量池：儲存所有常量的「容器」
                   └───────────────┘
                           │
                           │ 包含多個
                           ▼
              ┌────────────────────────┐
              │  ConstantInfo (陣列)   │  ← 每一個常量的「資料結構」
              └────────────────────────┘
                           │
                           │ 每個 ConstantInfo 有一個
                           ▼
                   ┌───────────────┐
                   │  ConstantTag  │  ← 標記這個常量的「類型」
                   └───────────────┘
```

<br>
<br>

### 1. ClassReader (二進制讀取器)

作用：負責從 `.class` 檔案讀取原始的二進制資料。

`.class` 檔案是按照特定順序排列的 bytes。ClassReader 封裝了讀取邏輯，我們可以指揮他讀取多少個單位的 bytes。

```go
gotype ClassReader struct {
    data []byte  // 原始二進制資料
}
// 讀取 1 個 byte (u1)
func (r *ClassReader) ReadUint8() uint8
// 讀取 2 個 bytes (u2)  
func (r *ClassReader) ReadUint16() uint16
// 讀取 4 個 bytes (u4)
func (r *ClassReader) ReadUint32() uint32
...
```

<br>
<br>

### 2. ConstantPool — 常量池

作用：儲存 class 檔案中所有的「常量」資料。

想像一下，你的 Java 程式碼中有：

```java
public class Hello {
    public static void main(String[] args) {
        System.out.println("Hello World");
    }
}
```

這段程式碼涉及到很多「名稱」：

* 類別名：Hello、System、PrintStream
* 方法名：main、println
* 字串："Hello World"
* 描述符：([Ljava/lang/String;)V

<br>

如果每次用到都重複儲存，會很浪費空間。所以 JVM 設計了常量池：

把所有常量集中存放，用索引來引用。

```go
type ConstantPool struct {
    infos []ConstantInfo  // 索引從 1 開始（索引 0 不用）
}

// 根據索引取得常量
func (cp *ConstantPool) GetConstantInfo(index uint16) ConstantInfo
```

<br>
<br>

### 3. ConstantTag — 常量標籤

作用：標識常量的類型。

常量池中的常量有很多種類型：

# 常量池（Constant Pool）Tag 對照表

| Tag | 常量類型 (Constant Type)          | 用途說明 |
|-----|-----------------------------------|----------|
| 1   | CONSTANT_Utf8                     | UTF-8 編碼的字串 |
| 3   | CONSTANT_Integer                  | int 字面量 |
| 4   | CONSTANT_Float                    | float 字面量 |
| 5   | CONSTANT_Long                     | long 字面量 |
| 6   | CONSTANT_Double                   | double 字面量 |
| 7   | CONSTANT_Class                    | 類別或介面的符號引用 |
| 8   | CONSTANT_String                   | 字串字面量 |
| 9   | CONSTANT_Fieldref                 | 欄位的符號引用 |
| 10  | CONSTANT_Methodref                | 方法的符號引用 |
| 11  | CONSTANT_InterfaceMethodref       | 介面方法的符號引用 |
| 12  | CONSTANT_NameAndType              | 名稱和類型描述符 |

> 為什麼從 1 開始且跳過 2？ 這是歷史原因，tag=2 原本是 Unicode 但被廢棄了。

<br>
<br>

### 4. ConstantInfo — 常量資訊

作用：表示常量池中每一個常量的具體資料結構。

不同類型的常量，結構不同：

```go
type ConstantInfo interface {
    readInfo(reader *ClassReader)  // 讀取自己的資料
}

// UTF-8 字串常量
type ConstantUtf8Info struct {
    value string
}

// 類別引用常量
type ConstantClassInfo struct {
    nameIndex uint16  // 指向常量池中的 UTF-8 常量
}

// 方法引用常量
type ConstantMethodrefInfo struct {
    classIndex       uint16  // 指向 Class 常量
    nameAndTypeIndex uint16  // 指向 NameAndType 常量
}
```

<br>

**關鍵概念：常量之間可以互相引用！**

<br>
<br>

## 它們如何協作？一個完整的例子

假設要解析 `System.out.println` 這個方法呼叫：

常量池內容（簡化）：

```
┌───────┬──────────────────────┬─────────────────────────────┐
│ Index │ Tag                  │ 內容                         │
├───────┼──────────────────────┼─────────────────────────────┤
│ 1     │ Methodref            │ classIndex=2, nameType=4    │
│ 2     │ Class                │ nameIndex=3                 │
│ 3     │ Utf8                 │ "java/io/PrintStream"       │
│ 4     │ NameAndType          │ nameIndex=5, descIndex=6    │
│ 5     │ Utf8                 │ "println"                   │
│ 6     │ Utf8                 │ "(Ljava/lang/String;)V"     │
└───────┴──────────────────────┴─────────────────────────────┘
```

**解析流程**：

1. ClassReader 讀取到 Methodref（索引 1）

2. 要知道是哪個類別的方法 → 查 classIndex=2
   → 找到 Class 常量 → 查 nameIndex=3  
   → 找到 "java/io/PrintStream"

3. 要知道方法名和簽名 → 查 nameType=4
   → 找到 NameAndType → 查 nameIndex=5 和 descIndex=6
   → 找到 "println" 和 "(Ljava/lang/String;)V"

4. 最終得知：呼叫 PrintStream 類別的 println 方法

<br>
<br>
<br>
<br>

## 總結

| 組件            | 職責說明             | 比喻 |
|-----------------|----------------------|------|
| ClassReader     | 讀取二進制資料        | 翻譯官 |
| ConstantPool    | 儲存所有常量的容器    | 字典 |
| ConstantTag     | 標識常量的類型        | 字典中的詞性標記 |
| ConstantInfo    | 每個常量的具體資料    | 字典中的每一個詞條 |

> ClassReader 讀取二進制 → 根據 ConstantTag 判斷類型 → 建立對應的 ConstantInfo → 放入 ConstantPool

這種設計的好處是解耦和可擴展：新增常量類型只需加新的 Tag 和 Info，不影響其他部分。
