# Gogo JVM

<br>

---

<br>


Using GoLang to implement a simple JVM

<br>

## Structure

<br>

1. [JVM Entry](cmd/gogo_jvm)
2. [Classfile](classfile)
3. [JVM Runtime Data Area](runtime)
   * [Method-Area](runtime/method_area)
   * [Heap](runtime/heap)
   * [Slot](runtime/rtcore)
4. [instructions](instructions)
5. [interpreter](interpreter)
6. [native](native)

<br>
<br>


## Version

<br>

### MVP Plan (v0.1~v0.2)

這是一切的開始，起初只是好奇心，想做一個能跑 `System.out.println()` 的 JVM。到後面停不下來探索知識的慾望。
MVP 時期的筆記有些雜亂，移到 MVP Plan Documents 做紀念。

* [MVP Plan Documents](doc/mvp_plan)

<br>
<br>


### 總覽：版本規劃

| 大版本 | 主題 | 子版本數 | 目標                  |
|--------|------|----------|---------------------|
| v0.1.x | ClassFile 解析 | 1 | 了解 .class 結構基礎      |
| v0.2.x | 核心執行引擎 | 11 | 單執行緒 Java 程式完整執行    |
| v0.3.x | 物件模型完善 | 4 | GC/多執行緒 共同前置        |
| v0.4.x | 多執行緒 | 6 | 併發程式支援 + Safe Point |
| v0.5.x | GC 專題 | 7 | 深入研究各種 GC 演算法       |

<br>
<br>


### v0.2 - MVP JVM Project (Core Engine) ⌛

v0.2 目標是讓 GOGO_JVM 能執行大部分 Single-Thread 的 Java 程式。

* ✅️ v0.2.1: [Basic Instructions](doc/version/v0_2/v0_2_1.md)
* ✅️ v0.2.2: [Simple Runtime Thread](doc/version/v0_2/v0_2_2.md)
* ✅️ v0.2.3: [Basic Method Area Implementation](doc/version/v0_2/v0_2_3.md)
* ✅️ v0.2.4: [Define Object Structure And `new` Instruction](doc/version/v0_2/v0_2_4.md)
* ✅️ v0.2.5: [Invoke Object Method `invokespecial` And `invokevirtual`](doc/version/v0_2/v0_2_5.md)
* ✅️ v0.2.6: [Support Basic Type Array](doc/version/v0_2/v0_2_6.md) 
* ✅️ v0.2.7: [Native Method And `System.out.println`](doc/version/v0_2/v0_2_7.md) 
* ✅️ v0.2.8: [Class Enhancement - `instanceof` / `checkcast` / `anewarray`](doc/version/v0_2/v0_2_8.md)
* ✅️ v0.2.9: [Support `java.lang.String`](doc/version/v0_2/v0_2_9.md)
* ⌛ v0.2.10: [Support ` try-catch-finally`](doc/version/v0_2/v0_2_10.md)
* ⛔ v0.2.11: [Support interface](doc/version/v0_2/v0_2_11.md)

<br>
<br>

### v0.3 - Object Model (Pre-Work for GC/Multi-Thread)

v0.3 目標是能支援 Class 物件與基本反射，標準物件頭 (Mark Word)，
Reference 類型支援，完整的物件生命週期。後續可以實現多執行緒和 GC。

* ⛔ v0.3.0: [Basic Reflection](doc/version/v0_3/v0_3_0.md)
* ⛔ v0.3.1: [Revamp - Object Header](doc/version/v0_3/v0_3_1.md)
* ⛔ v0.3.2: [Support Reference Type `java.lang.ref`](doc/version/v0_3/v0_3_2.md)
* ⛔ v0.3.3: [Object Lifecycle (Object from creation to Destroy)](doc/version/v0_3/v0_3_3.md)

<br>
<br>

### v0.4 - Multi-Thread  (Pre-Work for GC)

v0.4 目標是實現 `java.lang.Thread` 類與基本執行緒管理。

* ⛔ v0.4.0: [Thread Basic](doc/version/v0_4/v0_4_0.md)
* ⛔ v0.4.1: [Support `synchronized`](doc/version/v0_4/v0_4_1.md)
* ⛔ v0.4.2: [Support `wait` / `notify`](doc/version/v0_4/v0_4_2.md)
* ⛔ v0.4.3: [Support `interrupt`](doc/version/v0_4/v0_4_3.md)
* ⛔ v0.4.4: [Support `volatile` And Memory Model](doc/version/v0_4/v0_4_4.md)
* ⛔ v0.4.5: [Safe Point (Pre-Work for GC)](doc/version/v0_4/v0_4_5.md)

<br>
<br>

### v0.5 - GC

v0.5 深入研究實現各種 GC 演算法

* ⛔ [v0.5.0: Heap Structure](doc/version/v0_5/v0_5_0.md)
* ⛔ [v0.5.1: Mark-Sweep GC](doc/version/v0_5/v0_5_1.md)
* ⛔ [v0.5.2: Mark-Compact GC](doc/version/v0_5/v0_5_2.md)
* ⛔ [v0.5.3: Copying GC](doc/version/v0_5/v0_5_3.md)
* ⛔ [v0.5.4: GC (Generational)](doc/version/v0_5/v0_5_4.md)
* ⛔ [v0.5.5: Concurrent GC Basic](doc/version/v0_5/v0_5_5.md)
* ⛔ [v0.5.6: GC Advance (G1 / ZGC)](doc/version/v0_5/v0_5_6.md)

<br>
<br>

### 後續可能涉足的領域

**類加載器完善**


目前：單一 ClassLoader
目標：完整的類加載器層次結構

```
目標：
├── 雙親委派模型
├── 類的隔離與共享
├── 動態類加載
├── 模組系統（Java 9+ Module）
└── OSGi 概念
```

實現內容：
* Bootstrap ClassLoader
* Extension ClassLoader  
* Application ClassLoader
* 自定義 ClassLoader
* Class Unloading（配合 GC）

<br>
<br>

**JIT 編譯器**

目前：解釋執行（慢）
目標：熱點代碼編譯成機器碼（快） - 理解現代 VM

```
目標：
├── 編譯器後端設計
├── 寄存器分配
├── 指令選擇
├── 內聯優化（Inlining）
├── 逃逸分析（Escape Analysis）
└── On-Stack Replacement (OSR)
```

實現路徑：
* 熱點偵測（計數器）
* 基本 IR（中間表示）設計
* 簡單的機器碼生成（x86-64）
* 方法內聯
* 逃逸分析 + 標量替換

<br>
<br>

**調試與監控工具**

目標：實現 JVMTI / JMX 類似功能

```
目標：
├── 斷點機制
├── 單步執行
├── 變數查看
├── 性能分析（Profiling）
└── 記憶體分析（Heap Dump）
```

實現內容：
* 調試協議（類似 JDWP）
* 斷點表管理
* Stack Trace 生成
* Heap Histogram
* GC 日誌與分析

<br>
<br>

### 開發邊學習隨心記:

* Java Constant Pools: [link](doc/tips/constant_pool.md)
* Java Method Descriptor: [link](doc/tips/descriptor.md)
* ClassLoader: 為什麼類別載入時的 static field slot ID 計算時不需要考慮父類別？: [link](doc/tips/class_loader_cal_field_slot_id.md)