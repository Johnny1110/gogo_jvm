# GOGO JVM (當前本版 v0.2.3)

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
	magic        uint32 // magic number: 0xCAFEBABE, for classify .lang file (4 bytes)
	minorVersion uint16
	majorVersion uint16
	constantPool ClassFileConstantPool // constants pool
	accessFlags  uint16                // lang access flags
	thisClass    uint16                // this lang index (pointing to constantPool)
	superClass   uint16                // super lang index
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

### Phase v0.2.4: 物件建立與欄位存取 ⌛ (進行中)

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

### Phase v0.2.5：實例方法調用

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

### Phase v0.2.6：陣列支援

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
    data   interface{}  // []int32, []int64, []float32, []*Object.java 等
}
```

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


