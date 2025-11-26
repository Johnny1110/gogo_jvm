# Java 方法描述符（Method Descriptor）

<br>

JVM 需要一種**緊湊且無歧義**的方式來描述類型。這種格式：

1. **節省空間**：比完整文字短很多
2. **容易解析**：從左到右讀一次就能解析完成
3. **唯一性**：同樣的方法簽名一定產生同樣的描述符

<br>

---

<br>

## 類型描述符對照表

### 基本類型（Primitive Types）

| 描述符 | Java 類型 | 記憶方式 |
|--------|-----------|----------|
| `B` | byte | **B**yte |
| `C` | char | **C**har |
| `D` | double | **D**ouble |
| `F` | float | **F**loat |
| `I` | int | **I**nt |
| `J` | long | 因為 L 被物件用了，取 lon**J** |
| `S` | short | **S**hort |
| `Z` | boolean | 因為 B 被 byte 用了，取 booleanZ（想成 true/false 的 Zero/nonZero） |

### 引用類型（Reference Types）

| 描述符 | Java 類型 |
|--------|-----------|
| `Ljava/lang/String;` | String |
| `Ljava/lang/Object;` | Object |
| `Lcom/example/MyClass;` | com.example.MyClass |

**格式**：`L` + 全限定名（用 `/` 代替 `.`）+ `;`

### 陣列類型

| 描述符 | Java 類型 |
|--------|-----------|
| `[I` | int[] |
| `[D` | double[] |
| `[[I` | int[][] |
| `[Ljava/lang/String;` | String[] |
| `[[Ljava/lang/Object;` | Object[][] |

**規則**：幾個 `[` 就是幾維陣列

---

## 更多方法描述符範例

| 方法簽名 | 描述符 |
|----------|--------|
| `void run()` | `()V` |
| `int getAge()` | `()I` |
| `void setName(String name)` | `(Ljava/lang/String;)V` |
| `int add(int a, int b)` | `(II)I` |
| `String concat(String a, String b)` | `(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;` |
| `void process(int[] arr)` | `([I)V` |
| `double[][] createMatrix(int rows, int cols)` | `(II)[[D` |
| `Object get(String key, boolean flag)` | `(Ljava/lang/String;Z)Ljava/lang/Object;` |

---