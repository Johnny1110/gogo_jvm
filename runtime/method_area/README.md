# Method Area (類和方法的運行時表示)

<br>

---

<br>

這是 Java 運行時區域的 Method Area，這個區域的資料是 Thread 共享的。

<br>

## 重要核心組成

* [ClassLoader](class_loader.go) -> 負責將 class 完整載入 (從 classfile bytecode 開始解析)
* [RuntimeConstantPool](constant_pool.go) -> 運行時常量池 (每一個 class 都有一個自己專用，用於存放 class constant)
* [Class](class.go) -> Class 的實例，被 ClassLoader 載入完畢後的 classfile 就變成它，並存在 ClassLoader 裡面 (同一個類只會被載入一次)
* [Field](field.go) -> Class 的 field 實例
* [Method](method.go) -> Class 的 method 實例
* [ClassRef](cp_class_ref.go) -> 指向 Class 實例的參考
* [FieldRef](cp_field_ref.go) -> 指向 Field 實例的參考
* [MethodRef](cp_method_ref.go) -> 指向 Method 實例的參考
* [MethodDescriptor](method_descriptor_parser.go) -> 方法的 descriptor 解析，將純字面量解構成資料結構。

<br>
<br>

接下來的重點就是逐步拆解每一個重要的組成單元，看看他們是如何組在一起就能拼湊出 JVM Method Area 的。

<br>
<br>

## ClassLoader 類別加載器

ClassLoader 是 JVM 中負責 載入 class 檔案 到記憶體並建立 Class 物件的核心元件。

它知道去哪裡找 .class，怎麼讀、怎麼驗證、怎麼放到 JVM 裡。

**ClassLoader 的主要任務**

1. 將 .class 載入 JVM（Load）

   * 從 dir，JAR 等地方把 .class bytecode 資料讀出來。
   * 這裡會使用遞迴的方式一直向上載入 parent class 與 interfaces 直到所依賴的上層類別全部被載入為止。

2. 驗證 class 是否合法（Verify）

   * 確保 bytecode 合法，不會破壞 JVM。

3. 轉換成 JVM 裡的資料結構（Prepare / Resolve）

   * 透過 classfile 建立 RuntimeConstantPool，分配靜態變數空間等。

4. 建立 java.lang.Class 物件

   * 這個物件是 Class 在 JVM 內的 metadata (其實就放在 ClassLoader 內的一個 map 裏做緩存)。

<br>
<br>

### 分配靜態變數空間

分配靜態變數空間是在 `link` 階段實現，也就是在一個 class 已經被建立出來之後，我們仍需要知道一個類別的每一個靜態變量對應 `staticVars` 的哪一個索引。

```go
func prepare(class *Class) {
	// 計算非靜態 Fields slot 數量
	calcInstanceFieldSlotIds(class)
    // 計算靜態 Fields slot 數量
	calcStaticFieldSlotIds(class)
	// 分配並初始化 Vars 表
    allocAndInitVars(class)
}
```

<br>

**計算非靜態 Fields slot 數量**

```go
func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	
	if class.superClass != nil {
		// 父類別的先佔走一定數量
		slotId = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() { // 只看非 static fields
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() { // long 跟 double ˊ佔 2 個 slot 空間，多加一格
				slotId++
			}
		}
	}

	class.instanceSlotCount = slotId
}
```

<br>

**計算靜態 Fields slot 數量**

```go
func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() { // 只看 static fields
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() { // long 跟 double ˊ佔 2 個 slot 空間，多加一格
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}
```

<br>

**分配並初始化 Vars 表**

```go
func allocAndInitVars(class *Class) {
   class.instanceVars = rtcore.NewSlots(class.instanceSlotCount)
   class.staticVars = rtcore.NewSlots(class.staticSlotCount)
   // TODO: 初始化 static final 常量
}
```

<br>
<br>

## RuntimeConstantPool 運行時常量池

RuntimeConstantPool 是 ClassFileConstantPool 載入後（由 ClassLoader + JVM）放入 method area 的結構。
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

結構：

```go
type RuntimeConstantPool struct {
   class  *Class        // 所屬的 class
   consts []Constant    // 常量表
}
````

跟 ClassFileConstantPool 一樣，索引要從 1 開始。

在實際代碼 `func newRuntimeConstantPool(class *Class, cfCp classfile.ClassFileConstantPool) *RuntimeConstantPool { ... }`

我們可以看到基本上會把 ClassFileConstantPool 中的每一個 ConstantInfo 轉換成更能直接可用的 Constant 並存入 `consts` 中

* RuntimeConstantPool 的常量表索引跟 ClassFileConstantPool 保持一致。
* Utf8 與 NameAndType 類型的 `ConstantInfo` 不需要處理，因為他們在 `MemberConstantInfo` 轉換階段已經被利用完，沒有剩餘價值了。

<br>
<br>

## Class

```
 ClassFile（編譯時）    →    Class（運行時）
 ┌─────────────────┐        ┌─────────────────┐
 │ constantPool    │   →    │ constantPool    │  運行時常量池
 │ accessFlags     │   →    │ accessFlags     │
 │ thisClass       │   →    │ name            │  直接存類名
 │ superClass      │   →    │ superClass      │  指向父類 Class
 │ interfaces      │   →    │ interfaces      │  指向接口 Class[]
 │ fields          │   →    │ fields          │  運行時字段
 │ methods         │   →    │ methods         │  運行時方法
 └─────────────────┘        └─────────────────┘
```

<br>
<br>
<br>
<br>
<br>
<br>
<br>
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