# Constant 相關指令

<br>

---

<br>

## CONST 指令

把一些常用的 int (-1~5) long, float double 壓入 Stack 時使用。

<br>
<br>

## IPUSH 指令

使用 CONST 無法滿足時，可以使用 `BIPUSH` 和 `SIPUSH`

* `BIPUSH`: Byte Immediate PUSH
    - 操作數是 1 byte 有符號數
    - 範圍：-128 ~ 127

* `SIPUSH`: Short Immediate PUSH
    - 操作數是 2 bytes 有符號數
    - 範圍：-32768 ~ 32767

舉例：
```java
int a = 100;     編譯成 bipush 100
int b = 1000;    編譯成 sipush 1000
int c = 100000;  編譯成 ldc（從常量池載入）
```

<br>
<br>

## LDC 指令

LDC = Load Constant (載入常量)，這是 JVM 用來將 RuntimeConstantPool 中的 const 載入到 op-stack 的指令。

<br>

### 為什麼需要 LDC？

```
iconst_0 ~ iconst_5    → 只能載入 0 到 5
bipush                 → 只能載入 -128 到 127（1 byte）
sipush                 → 只能載入 -32768 到 32767（2 bytes）
```

* 當 `javac` 編譯時，大數字和浮點數，字串會被存入 ClassfileConstantPool。
* 當 ClassLoader 載入 Class 時，會將 ClassfileConstantPool 裡面的常數 const 映射到 RuntimeConstantPool。
* 運行時如果需要載入這些大的數字，浮點數，或字串就必須從 Class.RuntimeConstantPool 裡面載入。


<br>

```
int x = 1000000;           // 太大，bipush/sipush 放不下
float f = 3.14f;           // float 沒有專門的 push 指令
double d = 3.14159265358;  // double
String s = "Hello";        // 字串
```

這些值會被存放在 RuntimeConstantPool 中，然後用 `LDC` 載入。

<br>

### LDC 指令家族

| 指令      | Opcode | 說明                                         |
|-----------|--------|----------------------------------------------|
| ldc       | 0x12   | 載入常量池項目（1 byte 索引，值佔 1 slot）     |
| ldc_w     | 0x13   | 載入常量池項目（2 byte 索引，值佔 1 slot）     |
| ldc2_w    | 0x14   | 載入常量池項目（2 byte 索引，值佔 2 slots）    |

<br>

<br>

### 可載入的常量類型

| 類型         | 說明                                   |
|--------------|----------------------------------------|
| `int`        | 整數常量，push `int`                   |
| `float`      | 浮點數常量，push `float`               |
| `String`     | 字串常量，push `String` 物件參考       |
| `Class`      | 類別常量（反射用），push `Class` 參考  |
| `MethodType` | 方法型別（Java 7+）                    |
| `MethodHandle` | 方法句柄（Java 7+）                 |


<br>
<br>

### MVP 階段的限制

**String 常量**

完整的 JVM 需要：

1. 創建 java.lang.String 物件
2. 將字串存入物件的 char[] 字段
3. 字串駐留（String Interning）

MVP 階段：

* 創建特殊物件，將 Go string 存在 extra 字段。

<br>

**Class 常量**

* 用於 Foo.class 語法（反射），MVP 階段暫不支援。

<br>


