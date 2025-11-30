# JVM Runtime 基礎結構

<br>

---

<br>

## 1. Slot (槽位）

**JVM 需要一個統一的方式來存儲不同類型的數據。**

問題：如何用統一的結構存儲 int, float, long, double, reference?

解決方案：Slot（32-bit 槽位）

```
  ┌──────────┐                                       
  │   Slot   │  = 32 bits                           
  └──────────┘                                       
       │                                             
       ├── int       → 1 個 slot                       
       ├── float     → 1 個 slot                       
       ├── reference → 1 個 slot（指標）              
       ├── long      → 2 個 slots（64-bit）            
       └── double    → 2 個 slots（64-bit）            
```

<br>

**設計思想：**

* JVM 規範定義局部變量表和操作數棧都是由「槽位」(Slot) 組成的。
* 每個 Slot 可以存放一個 32-bit 的數據。
* 對於 64-bit 的數據（long, double），需要連續的兩個 Slot。


<br>
<br>