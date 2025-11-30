# Constant Pool

<br>

---

<br>

## Java 中的常量池總共分為 3 類（依 JVM 層級）

### 1. Class File Constant Pool（class 檔常量池）

位在 .class 檔裡，是 __編譯後靜態存在__ 的常量池。

用來存放：

* 字串常量（CONSTANT_String）
* 數值常量（Integer, Float, Long, Double）
* 類、欄位、方法的符號引用
* 名稱與描述符等

這是 Java 中最核心、最重要的 constant pool。

<br>


這個例子展示了哪些內容會進入常量池:

```java
public class ConstantPoolDemo {
    // 1. 字符串字面量 -> 進入常量池
    private static final String GREETING = "Hello, World!";
    
    // 2. 數字字面量 -> 進入常量池
    private static final int NUMBER = 42;
    private static final double PI = 3.14159;
    
    // 3. 類名、方法名、字段名 -> 都會進入常量池
    private String message;
    
    public static void main(String[] args) {
        // 4. 方法調用信息 -> 進入常量池
        System.out.println(GREETING);
        
        // 5. 類引用 -> 進入常量池
        ConstantPoolDemo demo = new ConstantPoolDemo();
        
        // 6. 字段訪問 -> 進入常量池
        demo.message = "Test";
        
        // 7. 方法調用 -> 進入常量池
        demo.printMessage();
    }
    
    public void printMessage() {
        System.out.println(message);
    }
}
```

<br>
<br>

### 2. Runtime Constant Pool（運行時常量池）

## Constant Pool 的生成時機

是 class file constant pool 載入後（由 ClassLoader + JVM）放入 method area 的結構

動態版本的常量池

可以增加新常量，例如：

* `String.intern()` 加入的新字串
* `invokedynamic` 執行時生成的常量

JVM 執行時使用這個 pool 來解析符號引用 → 變成直接引用（指標）

__每個 Class 在 JVM 裡都有自己的 runtime constant pool。__

比較直觀一點看：

```
 ┌────────────────────────────────────────────────────────┐
 │  ClassFileConstantPool（Compile）             │
 │  #1 Methodref → class=#2, nameAndType=#3               │
 │  #2 Class → name=#4                                    │
 │  #3 NameAndType → name=#5, desc=#6                     │
 │  #4 Utf8 → "Calculator"                                │
 │  #5 Utf8 → "add"                                       │
 │  #6 Utf8 → "(II)I"                                     │
 └────────────────────────────────────────────────────────┘
                        ↓ parse
 ┌────────────────────────────────────────────────────────┐
 │  RumtimeConstantPool                                 │
 │  #1 MethodRef → pointing to Calculator.add()           │
 └────────────────────────────────────────────────────────┘
```

<br>
<br>

### 3. String Constant Pool（字串常量池 / String Intern Pool）

JVM 專門為字串建立的池

與 runtime constant pool 不同，是 JVM 內部獨立維護的一張 hash table

透過以下方式加入 pool：

* 字面量（"abc"）
* `String.intern()`

用於字串重複利用，提高效率與節省記憶體。