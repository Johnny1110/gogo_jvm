# References Instructions（引用類指令）

<br>

---

<br>

這個套件實作 JVM 的引用類指令，用於物件建立、欄位存取、方法呼叫等操作。

<br>

## 指令一覽

### 物件建立

| 指令 | Opcode | 功能 |
|------|--------|------|
| new | 0xBB | 建立物件實例 |

### 欄位存取

| 指令 | Opcode | 功能 |
|------|--------|------|
| getstatic | 0xB2 | 取得靜態欄位 |
| putstatic | 0xB3 | 設定靜態欄位 |
| getfield | 0xB4 | 取得實例欄位 |
| putfield | 0xB5 | 設定實例欄位 |

### 方法呼叫

| 指令 | Opcode | 功能 |
|------|--------|------|
| invokestatic | 0xB8 | 呼叫靜態方法 |
| invokespecial | 0xB7 | 呼叫建構子/私有方法/父類方法 (Phase v0.2.5) |
| invokevirtual | 0xB6 | 呼叫實例方法 (Phase v0.2.5) |

<br>

## new 指令

### 執行流程

```
new #2  (其中 #2 是常量池中的 ClassRef)

1. 從常量池取得 ClassRef
2. 解析 ClassRef，載入類別（如果還沒載入）
3. 檢查：不能是介面或抽象類
4. 確保類別已初始化（執行 <clinit>）
5. 呼叫 class.NewObject() 建立物件
6. 將 Object 引用 push 到棧
```

### 典型使用模式

```
new #2              // 建立物件，push 引用
dup                 // 複製引用（給建構子用）
invokespecial #3    // 呼叫 <init> 建構子
astore_1            // 存到局部變數
```

### 類別初始化機制

當發現類別尚未初始化時：

```
┌─────────────────────────────────────────────────────────────┐
│  1. RevertNextPC() - 回退 PC，讓 new 指令稍後重新執行       │
│  2. 建立 <clinit> 的 Frame 並 push 到執行緒棧              │
│  3. 解釋器執行 <clinit>                                    │
│  4. <clinit> 執行完，Frame 被 pop                          │
│  5. 重新執行 new 指令，這次 initStarted = true             │
│  6. 正常建立物件                                           │
└─────────────────────────────────────────────────────────────┘
```

<br>

## 欄位存取指令

### 靜態欄位 vs 實例欄位

```
┌──────────────────────────┐  ┌──────────────────────────┐
│  Class (方法區)           │  │  Object (堆)             │
├──────────────────────────┤  ├──────────────────────────┤
│  staticVars (Slots)      │  │  fields (Slots)          │
│  ┌─────┬─────┬─────┐     │  │  ┌─────┬─────┬─────┐    │
│  │count│ MAX │ ... │     │  │  │value│name │ ... │    │
│  └─────┴─────┴─────┘     │  │  └─────┴─────┴─────┘    │
└──────────────────────────┘  └──────────────────────────┘
         ↑                              ↑
  getstatic/putstatic            getfield/putfield
  不需要物件引用                  需要物件引用
```

### 棧變化

```
getstatic:
  執行前：[...]
  執行後：[..., value]

putstatic:
  執行前：[..., value]
  執行後：[...]

getfield:
  執行前：[..., objectref]
  執行後：[..., value]

putfield:
  執行前：[..., objectref, value]
  執行後：[...]
```

### 欄位類型處理

根據 descriptor 決定如何讀寫：

| Descriptor | Java 類型 | JVM 內部類型 | Slot 數量 |
|------------|----------|-------------|----------|
| B | byte | int32 | 1 |
| C | char | int32 | 1 |
| S | short | int32 | 1 |
| I | int | int32 | 1 |
| Z | boolean | int32 | 1 |
| F | float | float32 | 1 |
| J | long | int64 | 2 |
| D | double | float64 | 2 |
| L...;  | 引用類型 | reference | 1 |
| [ | 陣列 | reference | 1 |

<br>

## 範例：欄位存取的字節碼

### Java 程式碼

```java
class Counter {
    static int count = 0;
    int value;
    
    void increment() {
        value++;
        count++;
    }
}
```

### 編譯後的 increment() 方法

```
aload_0           // 載入 this
dup               // 複製 this
getfield #2       // 取得 this.value
iconst_1          // 載入常數 1
iadd              // value + 1
putfield #2       // 設定 this.value

getstatic #3      // 取得 Counter.count
iconst_1          // 載入常數 1
iadd              // count + 1
putstatic #3      // 設定 Counter.count

return
```

<br>

## 初始化觸發時機

以下指令會觸發類別初始化：

1. **new** - 建立物件時
2. **getstatic** - 存取靜態欄位時
3. **putstatic** - 設定靜態欄位時
4. **invokestatic** - 呼叫靜態方法時

所有這些指令都會檢查 `class.InitStarted()`，如果為 false 則先執行 `<clinit>`。

<br>

## 錯誤處理

| 錯誤 | 觸發條件 |
|------|---------|
| InstantiationError | new 介面或抽象類 |
| IncompatibleClassChangeError | getstatic/putstatic 存取非靜態欄位，或 getfield/putfield 存取靜態欄位 |
| NullPointerException | getfield/putfield 的 objectref 為 null |