# Method Area (類和方法的運行時表示)

<br>

---

<br>

## MethodDescriptor 方法描述符解析結果

將方法描述符轉化成資料結構

例如：`(IDLjava/lang/String;)V`

```go
type MethodDescriptor struct {
	parameterTypes []string
	returnType     string
}
```

* parameterTypes: ["I", "D", "Ljava/lang/String;"]
* returnType: "V"

### 描述符格式：(參數類型)返回類型

類型編碼：
```
B - byte      C - char      D - double    F - float
I - int       J - long      S - short     Z - boolean
V - void      L類名; - 引用類型    [ - 數組
```

例子：
```
()V                      → void method()
(II)I                    → int method(int, int)
(Ljava/lang/String;)V    → void method(String)
([I)V                    → void method(int[])
([[Ljava/lang/Object;)V  → void method(Object[][])
```

<br>

## ClassLoader 類加載器

### 職責：
1. 根據類名找到 .class 文件
2. 解析 ClassFile
3. 轉換為運行時 Class 結構
4. 存入 Method Area（classMap）
5. 保證類的唯一性（同一個類只加載一次）

<br>

### 類加載流程：
```
 ┌─────────────────────────────────────────────────────────┐
 │  Loading → Linking → Initialization                     │
 │                                                         │
 │  Loading:                                               │
 │    - 讀取 .class 文件                                    │
 │    - 解析成 ClassFile                                   │
 │    - 轉換成 Class                                       │
 │                                                         │
 │  Linking:                                               │
 │    - Verification: 驗證字節碼（簡化跳過）                   │
 │    - Preparation: 為靜態變量分配空間                       │
 │    - Resolution: 符號引用解析（懶加載）                     │
 │                                                         │
 │  Initialization:                                        │
 │    - 執行 <clinit> 方法（靜態初始化）                       │
 └─────────────────────────────────────────────────────────┘
```

## 執行流程圖解

當執行 `invokestatic StaticCall.add` 時：

```
invokestatic #1
      │
      ▼
┌─────────────────────────────────────────────────────────┐
│ INVOKE_STATIC.Execute(frame)                            │
│                                                         │
│   cp := frame.Method().Class().ConstantPool()          │
│   methodRef := cp.GetConstant(1).(*MethodRef)          │
│   method := methodRef.ResolvedMethod()  ◄───────────┐   │
│                                                     │   │
└─────────────────────────────────────────────────────│───┘
                                                      │
      ┌───────────────────────────────────────────────┘
      ▼
┌─────────────────────────────────────────────────────────┐
│ MethodRef.ResolvedMethod()                              │
│                                                         │
│   if r.method == nil {                                 │
│       r.resolveMethodRef()  ◄────────────────────┐      │
│   }                                              │      │
│   return r.method                                │      │
│                                                  │      │
└──────────────────────────────────────────────────│──────┘
                                                   │
      ┌────────────────────────────────────────────┘
      ▼
┌─────────────────────────────────────────────────────────┐
│ MethodRef.resolveMethodRef()                            │
│                                                         │
│   c := r.ResolvedClass()  ◄──────────────────────┐      │
│   method := lookupMethod(c, r.name, r.descriptor)│      │
│   r.method = method                              │      │
│                                                  │      │
└──────────────────────────────────────────────────│──────┘
                                                   │
      ┌────────────────────────────────────────────┘
      ▼
┌─────────────────────────────────────────────────────────┐
│ SymRef.ResolvedClass()                                  │
│                                                         │
│   if r.class == nil {                                  │
│       r.resolveClassRef()  ◄─────────────────────┐      │
│   }                                              │      │
│   return r.class                                 │      │
│                                                  │      │
└──────────────────────────────────────────────────│──────┘
                                                   │
      ┌────────────────────────────────────────────┘
      ▼
┌─────────────────────────────────────────────────────────┐
│ SymRef.resolveClassRef()                                │
│                                                         │
│   d := r.cp.class          // 當前類 (StaticCall)       │
│   loader := d.loader       // 當前類的 ClassLoader      │
│   c := loader.LoadClass(r.className)  // 加載目標類     │
│   r.class = c              // 緩存                      │
│                                                         │
└─────────────────────────────────────────────────────────┘

```

<br>

## 為什麼用「當前類的 ClassLoader」？

這是 類加載器委託模型（Class Loader Delegation） 的關鍵：

```
god := r.cp.class        // d = 當前正在執行的類
c := d.loader.LoadClass(r.className)  // 用 d 的 loader 加載新類
```

**原因：**

1. **雙親委派模型**：確保相同的類名由相同的 ClassLoader 加載
2. **類型安全**：不同 ClassLoader 加載的同名類是不同的類型
3. **命名空間隔離**：同一個類名在不同 ClassLoader 下可以是不同的類

例如：

```
┌─────────────────────────────────────────────────────────┐
│  StaticCall 類由 AppClassLoader 加載                     │
│  StaticCall 調用 Calculator.add()                       │
│  Calculator 也應該由 AppClassLoader 加載                 │
│  （使用同一個 ClassLoader）                               │
└─────────────────────────────────────────────────────────┘
```