# GOGO JVM (當前本版 v0.2.5)

<br>

---

<br>


## Process

1. 讀取 .class 文件
2. 使用 ClassLoader 加載類
3. 找到 main 方法
4. 執行解釋器

<br>

目前實現到方法區，使用 ClassLoader 來載入類別，並使用解釋器直接執行 main 方法。

<br>

## Phase 

### Phase v0.1.1: Classfile ✅️(完成)

目標: 完成 java classfile 解讀功能，能夠把編譯好的 .class 檔案建構成 Classfile:

```go
type ClassFile struct {
	magic        uint32 // magic number: 0xCAFEBABE, for classify .class file (4 bytes)
	minorVersion uint16
	majorVersion uint16
	constantPool ClassFileConstantPool // constants pool
	accessFlags  uint16                // class access flags
	thisClass    uint16                // this class index (pointing to constantPool)
	superClass   uint16                // super class index
	interfaces   []uint16              // implemented interfaces index
	fields       []*MemberInfo         // fields table
	methods      []*MemberInfo         // methods table
	attributes   []AttributeInfo       // attributes table
}
```

<br>

---

<br>


### Phase v0.2.1: Basic Instructions ✅️(完成)

目標：

實現一些基礎的指令:

* CONST 
* IPUSH 
* LOAD 
* STORE 
* 算數 
* 邏輯跳轉 
* 返回
* Instruction Factory

<br>

---

<br>

### Phase v0.2.2: Runtime Thread ✅️(完成)

目標：

1. 實現 JVM Runtime Thread 相關 (Thread 私有)
    * Thread
    * JVMFrameStack
    * Frame
    * OprandsStack
    * LocalVars
    * Slot

<br>

可以運行簡單的 add() 方法:

```java
public class TestAdd {
    public static void main(String[] args) {
        int a = 1;
        int b = 2;
        int c = a + b;  // 結果存在 locals[3]
    }
}
```

<br>

---

<br>

### Phase v0.2.3: 靜態方法呼叫 ✅️(完成)

目標：

1. 實現 JVM 方法區
   * ClassLoader
   * RuntimeConstantPool
   * ClassRef, FieldRef, MethodRef
   * Class, Field, Method
2. 實現 `INVOKE_STATIC` 指令 (靜態方法呼叫)

<br>

* 測試目標：

```java
public class Fibonacci {
    public static void main(String[] args) {
        int result = fib(10);
    }

    public static int fib(int n) {
        if (n <= 1) return n;
            return fib(n-1) + fib(n-2);
    }
 }
```

<br>

---

<br>

### Phase v0.2.4: 物件建立與欄位存取 ✅️(完成)

目標：能夠執行 `new` 建立物件，並存取物件欄位

* **需要實現的指令**：

| 指令       | Opcode | 功能             |
|------------|--------|----------------|
| new        | `0xBB`   | 建立物件實例     |
| getstatic  | `0xB2`   | 取得靜態欄位     |
| putstatic  | `0xB3`   | 設定靜態欄位     |
| getfield   | `0xB4`   | 取得實例欄位     |
| putfield   | `0xB5`   | 設定實例欄位     |

* 需要完善的模組：
  * /runtime/heap/object.go
  * `slotId` 在 ClassLoader.prepare() 階段計算還未完成，需要補齊

* 測試目標：

```go
public class TestObject {
    public static void main(String[] args) {
        Counter c = new Counter();
        c.value = 10;
        int x = c.value;  // x = 10
    }
}

class Counter {
    int value;
}
```

<br>

---

<br>

### Phase v0.2.5：實例方法調用 ✅️(完成)

**目標**：支援 `invokevirtual` 和 `invokespecial`

**需要實現的指令**：

| 指令 | Opcode | 功能 |
|------|--------|------|
| `invokespecial` | `0xB7` | 呼叫建構子、私有方法、父類方法 |
| `invokevirtual` | `0xB6` | 呼叫實例方法（支援多型） |
| `dup` | `0x59` | 複製棧頂元素（new 後常用） |

<br>

**設計要點**：

invokevirtual 的多型查找：

1. 從物件的實際類型開始查找方法
2. 若找不到，往父類查找
3. 這是 Java 多型的核心機制

<br>

測試目標：

```java
public class TestMethod {
    public static void main(String[] args) {
        Counter c = new Counter();
        c.increment();
        c.increment();
        int x = c.getValue();  // x = 2
    }
}

class Counter {
    int value;
    public Counter() { this.value = 0; }
    public void increment() { this.value++; }
    public int getValue() { return this.value; }
}
```

<br>

---

<br>

### Phase v0.2.6：陣列支援 ⌛ (進行中)

目標：支援基本型別陣列和物件陣列

**需要實現的指令：**

| 指令              | Opcode        | 功能               |
|-------------------|---------------|--------------------|
| newarray          | `0xBC`          | 建立基本型別陣列   |
| anewarray         | `0xBD`          | 建立物件陣列       |
| arraylength       | `0xBE`          | 取得陣列長度       |
| iaload / iastore  | `0x2E` / `0x4F`   | int 陣列讀寫       |
| aaload / aastore  | `0x32` / `0x53`   | 物件陣列讀寫       |
| xaload / xastore  | （多組 Opcode）| 其他型別陣列讀寫   |

<br>

設計要點：陣列也是物件，但有特殊結構

```go
// 
type ArrayObject struct {
    class  *Class
    length int
    data   interface{}  // []int32, []int64, []float32, []*Object 等
}
```

<br>

**JVM 陣列的特殊性**

1. 陣列不是由 ClassLoader 從 .class 載入的 → JVM 動態生成陣列類別
2. 陣列類別的命名規則:
   ```
   I        → int[]                                         
   [D        → double[]                                       
   [[I       → int[][]                                        
   [Ljava/lang/String;  → String[]
   ```
3. 陣列資料存在 `Object.extra` 欄位
   * 普通物件：`fields` 存欄位，`extra = nil`
   * 陣列物件：`fields = nil`，`extra` 存陣列資料

<br>

**陣列物件的記憶體佈局**

```
Normal Object: 
┌─────────────────────────────────────────┐
│ Object                                  │
├─────────────────────────────────────────┤
│ class ────────────────────────► Counter │
│ fields: [slot0, slot1, ...]             │
│ extra: nil                              │
└─────────────────────────────────────────┘

Array Object: 
┌─────────────────────────────────────────┐
│ Object (Array)                          │
├─────────────────────────────────────────┤
│ class ────────────────────────► [I      │
│ fields: nil                             │
│ extra: []int32{100, 200, 0, 0, 0}       │
└─────────────────────────────────────────┘
```

<br>

**為什麼不同類型用不同的底層陣列？**

```go
// byte[]/boolean[] → []int8   (1 byte per element)
// short[]          → []int16  (2 bytes per element)
// char[]           → []uint16 (2 bytes per element, unsigned)
// int[]            → []int32  (4 bytes per element)
// long[]           → []int64  (8 bytes per element)
// float[]          → []float32
// double[]         → []float64
// Object[]         → []*Object
```

原因: 記憶體效率考量，如果所有陣列都用 []int64，一個 `byte[1000000]` 會浪費 7MB 記憶體。

<br>

**`newarray` 的 `atype` 參數**

```
atype = 4  → boolean[]
atype = 5  → char[]
atype = 6  → float[]
atype = 7  → double[]
atype = 8  → byte[]
atype = 9  → short[]
atype = 10 → int[]
atype = 11 → long[]
```

<br>

**陣列元素存取**

| 載入指令 | Opcode | 存入指令 | Opcode | 類型 |
|----------|--------|----------|--------|------|
| iaload | 0x2E | iastore | 0x4F | int[] |
| laload | 0x2F | lastore | 0x50 | long[] |
| faload | 0x30 | fastore | 0x51 | float[] |
| daload | 0x31 | dastore | 0x52 | double[] |
| aaload | 0x32 | aastore | 0x53 | Object[] |
| baload | 0x33 | bastore | 0x54 | byte[]/boolean[] |
| caload | 0x34 | castore | 0x55 | char[] |
| saload | 0x35 | sastore | 0x56 | short[] |

<br>

**陣列長度**

| 指令 | Opcode | 說明 |
|------|--------|------|
| arraylength | 0xBE | 取得陣列長度 |

<br>

**指令執行流程**

* `newarray` 執行流程:
   ```
   Stack: [..., count] → [..., arrayref]
   ```
  1. 從 stack pop 出 `count` (陣列長度)
  2. 檢查 `count >= 0`，否則 `NegativeArraySizeException`
  3. 根據 `atype` 建立對應類型的陣列
  4. 將陣列引用 push 到 stack

<br>

* `iaload` 執行流程:
   ```
   Stack: [..., arrayref, index] → [..., value]
   ```
  1. 從 stack pop 出 `index`
  2. 從 stack pop 出 `arrayref`
  3. 檢查 `arrayref != null`，否則 `NullPointerException`
  4. 檢查 `0 <= index < length`，否則 `ArrayIndexOutOfBoundsException`
  5. 取得 `arr[index]` 的值
  6. 將值 push 到 stack

<br>

* `iastore` 執行流程:
   ```
   Stack: [..., arrayref, index, value] → [...]
   ```
  1. 從 stack pop 出 `value`
  2. 從 stack pop 出 `index`
  3. 從 stack pop 出 `arrayref`
  4. 檢查 `arrayref != null`
  5. 檢查 `index` 合法
  6. 設定 `arr[index] = value`

<br>

**編譯範例**

```java

編譯前：------------------------------------------

int[] arr = new int[5];
arr[0] = 100;
int x = arr[0];

編譯後：------------------------------------------

bipush 5          // push 5 (陣列長度)
newarray int      // 建立 int[5]，結果是陣列引用
astore_1          // 存入 locals[1]

aload_1           // 載入陣列引用
iconst_0          // push 0 (索引)
bipush 100        // push 100 (值)
iastore           // arr[0] = 100

aload_1           // 載入陣列引用
iconst_0          // push 0 (索引)
iaload            // 載入 arr[0]
istore_2          // x = arr[0]
-------------------------------------------------
```

<br>

**MVP 限制**

1. **陣列類別為 nil**：完整實現需要 ClassLoader 動態生成陣列類別（如 `[I`）
2. **未實現 anewarray**：引用類型陣列（如 `String[]`）需要額外處理
3. **未實現 multianewarray**：多維陣列需要遞迴建立

<br>

測試目標：

```java
public class TestArray {
    public static void main(String[] args) {
        int[] arr = new int[5];
        arr[0] = 10;
        arr[1] = 20;
        int sum = arr[0] + arr[1];  // sum = 30
    }
}
```

<br>

---

<br>

### Phase v0.2.7：Native Method 與 `System.out.println`

目標：實現最基本的 Native Method 機制，讓程式能輸出結果

設計要點：

```go
// native/registry.go
type NativeMethod func(frame *runtime.Frame)

var registry = map[string]NativeMethod{}

func Register(className, methodName, descriptor string, method NativeMethod) {
    key := className + "~" + methodName + "~" + descriptor
    registry[key] = method
}

// 註冊 System.out.println
func init() {
    Register("java/io/PrintStream", "println", "(I)V", printlnInt)
}

func printlnInt(frame *runtime.Frame) {
    val := frame.LocalVars().GetInt(1)  // [0]=this, [1]=參數
    fmt.Println(val)
}
```

<br>

測試目標：


```java
public class HelloWorld {
    public static void main(String[] args) {
        System.out.println(42);      // 輸出: 42
        System.out.println(fib(10)); // 輸出: 55
    }
    
    public static int fib(int n) {
        if (n <= 1) return n;
        return fib(n-1) + fib(n-2);
    }
}
```

<br>

---

<br>


