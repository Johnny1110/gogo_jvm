# Reference 系列指令

<br>

---

<br>

## new (0xBB)

* 字節碼格式：`new indexbyte1 indexbyte2`
* 操作數: 2 bytes（常量池索引，指向 ClassRef)

<br>

`new` 指令流程:

1. 讀取操作數，取得常量池索引
2. 從 RuntimeConstantPool 取得 `ClassRef`
3. 解析 ClassRef，載入目標類別（如果還沒載入）得到 `*Class`
4. 檢查類別:
   - 不能是介面（interface）
   - 不能是抽象類（abstract）
5. 檢查類別是否已初始化
   - 如果沒有，觸發類別初始化（執行 <clinit>）
6. 呼叫 class.NewObject() 建立物件
7. 將 Object 引用 push 到操作數棧

<br>

__注意: new 只分配記憶體，不執行建構子!__，建構子由後續的 invokespecial <init> 呼叫。


<br>

### 執行 new 時發現類別沒初始化，怎麼辦?

1. `RevertNextPC(`)` - 讓 `new` 指令下次重新執行
2. 建立 `<clinit>` 的 Frame 並 push 到棧
3. 解釋器下一輪會執行 `<clinit>`
4. `<clinit>` 執行完 return，Frame 被 pop
5. 解釋器回到原本的 Frame，重新執行 new
6. 這次 `class.InitStarted()` 為 true，正常建立物件

<br>

**執行順序圖：**

```

  Stack:          Stack:          Stack:             
  ┌───────┐      ┌────────┐       ┌───────┐           
  │ main  │  →   │<clinit>│   →   │ main  │          
  │ new.. │      ├────────┤       │ new.. │           
  └───────┘      │ main   │       └───────┘            
                 └────────┘                              
    發現未初始化    執行 <clinit>    <clinit> 完成    
    RevertNextPC                   重新執行 new            

```

<br>



