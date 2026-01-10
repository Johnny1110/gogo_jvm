# private 方法會是由 `invokespecial` 處理而非 `invokevirtual`

<br>

---

<br>

我在實現 v0.3.1 的反射基礎功能時遇到一個問題：

我寫了一個簡易版的 Class.java:

```java
package java.lang;

public class Class {

    // Cache the name to reduce the number of calls into the VM.
    // This field would be set by VM itself during initClassName call.
    private transient String name;

    public String getName() {
            String name = this.name;
            return name != null ? name : initClassName();
    }

    private native String initClassName();
   
}
```

因為反射需要使用 `getName()` 方法，所以我複製了 source code 關於 `getName()` 的這一部分。並且在本地方法註冊表中 (native/java/lang/class.go)
 中註冊了 `initClassName()` 作為本地方法。

我原本預期當測試的 java code 需要呼叫 `getName()` 時會經過 `invokevirtual` 並被 `hacked_invoke_virtual()` 攔截轉倒給 native 方法處理。

結果卻是被交給 `invokespecial` 處理了，因為先前版本我對 `invokespecial` 的印象都是只作用於 `<clinit>` 方法。

<br>

實際上 `invokespecial` 的定義為：

1. constructor <init>
2. private method
3. parent method (super.xxx())

<br>
<br>

在 Java 中，private 方法永遠使用 `invokespecial` 而不是 `invokevirtual`，因為：

1. private 方法不能被覆寫，所以不需要動態綁定 (不需要去父類別找， private 一定在自己身上)
2. 編譯器在編譯時就確定要呼叫哪個方法
3. 這比 `invokevirtual` 的動態查找更高效