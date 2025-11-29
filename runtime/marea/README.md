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