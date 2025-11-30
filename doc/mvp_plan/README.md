# MVP Plan

<br>

---

<br>

## 理論基礎

**需要理解的核心概念：**

1. Class 文件結構 - Java 編譯後的 .class 文件格式

* 魔數（Magic Number）和版本號
* 常量池（Constant Pool）
* 訪問標誌、類索引、接口索引
* 字段表、方法表、屬性表

<br>

2. JVM 內存模型

* 堆（Heap）：存放對象實例
* 方法區（Method Area）：存放類信息
* 虛擬機棧（VM Stack）：每個線程的運行時數據
* 程序計數器（PC Register）：當前執行的字節碼位置
* 本地方法棧（Native Method Stack）

<br>

3. 字節碼指令集

* 基本指令類型：加載/存儲、運算、類型轉換、控制流程等
* 操作數棧的概念


<br>
<br>

## 階段 1：Class 文件解析器

Class 文件是 JVM 的輸入，就像編譯器需要先做詞法分析一樣，JVM 需要先能讀懂 Class 文件。


**目標結構：**
```
gogo_jvm/
├── classfile/
│   ├── class_file.go      // Class 文件主結構
│   ├── class_reader.go    // 讀取二進制數據
│   ├── constant_pool.go   // 常量池解析
│   ├── member_info.go     // 字段和方法信息
│   └── attribute_info.go  // 屬性表解析
└── main.go
```

**MVP 功能：**

* 能夠讀取並解析簡單的 .class 文件
* 印出類的基本信息（類名、方法、字段）

<br>
<br>

## 階段 2：運行時數據區

這是 JVM 執行字節碼的「工作空間」。就像 CPU 需要寄存器和內存，JVM 需要棧和堆來存儲運行時數據。


**目標結構：**
```
├── runtime/
│   ├── thread.go          // 線程
│   ├── stack.go           // 虛擬機棧
│   ├── frame.go           // 棧幀
│   ├── local_vars.go      // 局部變量表
│   ├── operand_stack.go   // 操作數棧
│   ├── trcore/
│   │   └── slot.go        // 定義 Slot 與 Slots 資料結構
│   │
│   ├── heap/
│   │   └── object.go      // 對象表示
│   │ 
│   └──method_area/
│        ├── class_loader.go      
│        ├── constant_pool.go  // runtime constant pool 
│        ├── rtma_class.go     // method-area class
│        ├── rtma_field.go     // method-area field
│        ├── rtma_method.go    // method-area method
│        ├── cp_class_ref.go   // constant pool class 直接引用
│        ├── cp_field_ref.go   // constant pool field 直接引用
│        └── cp_method_ref.go  // constant pool method 直接引用
```

**關鍵設計決策：**

* 棧幀如何表示方法調用
* 操作數棧如何實現（用 slice）
* 局部變量表的索引方式

<br>
<br>

## 階段 3：字節碼解釋器

解釋器負責執行字節碼指令。這就像 CPU 執行機器碼一樣。


**目標結構：**
```
├── instructions/
│   ├── base/              // 指令基礎結構
│   │   ├── instruction.go
│   │   └── bytecode_reader.go
│   ├── constants/         // 常量加載指令
│   ├── loads/            // 變量加載指令
│   ├── stores/          // 變量存儲指令
│   ├── stack/           // 棧操作指令
│   ├── math/            // 算術指令
│   ├── conversions/     // 類型轉換
│   ├── comparisons/     // 比較指令
│   ├── control/         // 控制流指令
│   └── references/      // 引用類指令（簡化版）
└── interpreter/
    └── interpreter.go    // 解釋器主循環
```

**MVP 目標：**

能執行簡單的 Java 程序，如：

```java
public class HelloWorld {
    public static void main(String[] args) {
        int a = 1;
        int b = 2;
        int c = a + b;
        System.out.println(c);  // 初期可以用特殊處理
    }
}
```

<br>
<br>

## 階段 4：類加載器

類加載器負責將 Class 文件加載到內存並創建類對象。這實現了 Java 的動態加載特性。

```
├── classpath/
│   ├── classpath.go      // 類路徑抽象
│   ├── entry.go          // 路徑項
│   └── entry_dir.go      // 目錄類路徑
└── classloader/
└── classloader.go     // 類加載邏輯
```

<br>
<br>

## 階段 5：方法調用

方法調用涉及參數傳遞、棧幀創建、返回值處理等，是 JVM 的核心功能之一。(複雜)

**需要實現：**

* invokestatic（靜態方法調用）
* return 指令族
* 方法參數傳遞
* 返回值處理

<br>

**MVP 最終目標**

能夠運行：

```java
public class Fibonacci {
    public static void main(String[] args) {
        int n = 10;
        int result = fib(n);
        System.out.println(result);
    }
    
    public static int fib(int n) {
        if (n <= 1) return n;
        return fib(n-1) + fib(n-2);
    }
}
```


<br>
<br>