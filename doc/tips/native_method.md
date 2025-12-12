# Native 方法機制設計

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

<br>
<br>

## 設計架構

Native 方法執行流程

```
Java Code                                                      
─────────                                                      
  System.out.println("Hello");                                   
                                                                 
       ↓ 編譯                                                    
                                                                 
Bytecode                                                       
─────────                                                      
  getstatic System.out                                           
  ldc "Hello"                                                    
  invokevirtual PrintStream.println(String)                      
                                                                 
       ↓ 執行                                                    
                                                                 
JVM                                                            
─────────                                                      
  1. 發現是 native 方法                                          
  2. 查找 NativeMethodRegistry                                   
  3. 找到對應的 Go func                                          
  4. 執行 Go func（真正印出文字）   
```