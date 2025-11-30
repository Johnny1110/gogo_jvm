# Stack 系列指令

<br>

---

<br>

## Stack 操作指令實作:

```
┌────────────┬────────┬────────────────────────────────────┐ 
│ 指令        │ Opcode │ 功能                               │ 
├────────────┼────────┼────────────────────────────────────┤  
│ DUP        │ 0x59   │ 複製棧頂 1 slot                      │  
│ DUP_X1     │ 0x5A   │ 複製並插入到第 2 個下面                │  
│ DUP_X2     │ 0x5B   │ 複製並插入到第 3 個下面                │  
│ DUP2       │ 0x5C   │ 複製棧頂 2 slots                     │  
│ DUP2_X1    │ 0x5D   │ 複製 2 個並插入到第 3 個下面           │
│ DUP2_X2    │ 0x5E   │ 複製 2 個並插入到第 4 個下面           │
│ POP        │ 0x57   │ 彈出 1 slot                         │  
│ POP2       │ 0x58   │ 彈出 2 slots                        │  
│ SWAP       │ 0x5F   │ 交換棧頂 2 個 slots                  │  
└────────────┴────────┴────────────────────────────────────┘  
```

* 使用 PopSlot/PushSlot 操作，不關心具體類型                   
* dup 複製的是 Slot（包含 Num 和 Ref）      
* 對於 reference，複製的是指標，不是物件


<br>

**為什麼需要 dup 指令?**

Java 物件建立的字節碼模式：

java:

```java
Counter c = new Counter();
```

bytecode:

```
new              // 建立物件，push Ref 到 stack 中
dup              // 複製 Top Stack Ref
invokespecial    // 呼叫 <init> 建構子
astore_1         // 存到 LocalVars 表
```

<br>

slot stack 的變化：                                                    
```
  ┌──────────────┬──────────────┬──────────────┬────────────┐  
  │  new         │  dup         │ invokespecial│  astore    │
  │              │              │  (popped 1)  │ (popped 1) │
  │  ┌─────┐     │  ┌─────┐     │  ┌─────┐     │  ┌─────┐   │
  │  │ ref │     │  │ ref │     │  │ ref │     │  │     │   │
  │  └─────┘     │  ├─────┤     │  └─────┘     │  └─────┘   │
  │              │  │ ref │     │              │            │
  │              │  └─────┘     │              │            │
  └──────────────┴──────────────┴──────────────┴────────────┘
```

如果沒有 dup： 
* `invokespecial` 消耗掉唯一的 ref
* `astore` 就沒有引用可以存了

<br>
<br>

## dup

操作：複製 top stack 的 1 個 slot，再 push 一次。

```
執行前：          執行後：                                      
  ┌─────────┐      ┌─────────┐                                  
  │  value  │ top  │  value  │ top  ← 複製的                    
  ├─────────┤      ├─────────┤                                  
  │   ...   │      │  value  │      ← 原本的                    
  └─────────┘      ├─────────┤                                  
                   │   ...   │                                  
                   └─────────┘                                  
```

實作方式：   

```
slot := stack.PopSlot()   // 彈出                            
stack.PushSlot(slot)      // 放回去                          
stack.PushSlot(slot)      // 再放一次
```