# Native 方法

<br>

---

<br>

## 為什麼需要 Native 方法？

有些功能 Java bytecode 無法實現：

* 印出文字到螢幕 (需要作業系統 API)
* 讀取檔案 (需要作業系統 API)
* 取得系統時間
* 執行緒操作

這些功能由「Native 方法」提供

* Java 側宣告 native 方法
* JVM 用 C/C++（或 Go）實現

```java
public class System {                                          
    public static native void print(String s);  // native method
} 
```

<br>
<br>

## Native Method 運作機制

當 Java 編譯器看到 native 關鍵字：

```java
public class System {
    public static native void arraycopy(Object src, int srcPos, Object dest, int destPos, int length);
}
```

編譯後，這個方法的 `access_flags` 會包含 `ACC_NATIVE (0x0100)`，且 **沒有 Code 屬性**。

<br>

### JVM 如何處理 native 方法

`invokevirtual PrintStream.println(I)V`

1. `methodRef.ResolveMethod()` 找到目標方法
2. 檢查 `method.IsNative()`
    * **false** -> 正常建立 frame，執行 bytecode
    * **true** -> 查找 Native Method 註冊表

<br>

Native Registry 是一個 KV 鍵值對：

```
key(string) -> func

ex: 
"println(I)V" -> printlnInt()
```

<br>

## `System.out.println` 的完整調用鏈

分析 `System.out.println(42);` 整個執行過程。

### Step 1: 編譯後的 Bytecode

```
getstatic     #2   // Field - 將參考 `java/lang/System.out:Ljava/io/PrintStream;` 推入 stack (可以在 java.lang.System 裡找到這個 PrintStream final const)
bipush        42   // 將 int32(42) 推入 stack 中，當前 stack: [PrintStreamRef, int32(42)]
invokevirtual #3   // Method - java/io/PrintStream.println:(I)V -> invoke PrintStream..println:(I)V
```

<br>

### Step 2: `getstatic` 執行

```
執行前 Stack: []

getstatic java/lang/System.out
   1. 解析 FieldRef → System 類的 `out` Field
   2. 檢查 System 類是否初始化（觸發 <clinit>）
   3. 從 System.staticVars 取得 `out` Field 的值
   4. Push `out` Ref 到 stack
   
執行後 Stack: [PrintStream ref]
```

**問題**: System 類的 `out` Field 從哪來？

在真實 JVM 中，`System.<clinit>` 會調用 `initializeSystemClass()`，這是一個 native 方法，由 JVM 啟動時設置。

**MVP 策略**：我們直接在載入 System 類時，手動設置 `out` 字段指向一個模擬的 PrintStream 物件。

<br>

### Step 3: `bipush` 42

```
執行前 Stack: [PrintStream ref]

bipush 42
    └─ Push 42 到棧

執行後 Stack: [PrintStream ref, 42]
```

<br>

### Step 4: invokevirtual 執行

```
執行前 Stack: [PrintStream ref, 42]

invokevirtual PrintStream.println(I)V
    │
    ├─ 1. 解析 MethodRef
    ├─ 2. 從棧頂往下數 argSlotCount 個位置，取得 objectref (PrintStream ref)
    │     argSlotCount = 2 (this + int參數)
    │     objectref = PrintStream ref
    ├─ 3. 動態綁定: 從 objectref 的實際類型查找 println
    ├─ 4. 發現 println 是 native 方法
    ├─ 5. 查找 Native 註冊表
    └─ 6. 調用 Go 函數 printlnInt(frame)

執行後 Stack: []  (void 方法，無返回值)
```

<br>
<br>

## 設計架構

```
gogo_jvm/
├── native/
│   ├── registry.go          # Native 方法註冊表
│   ├── init.go               # 自動註冊所有 native 方法
│   └── java/
│       ├── lang/
│       │   ├── system.go     # System 類的 native 方法
│       │   └── object.go     # Object 類的 native 方法 (hashCode 等)
│       └── io/
│           └── printstream.go # PrintStream.println 的實現
```