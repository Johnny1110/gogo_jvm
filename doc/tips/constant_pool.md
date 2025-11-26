# Constant Pool

<br>

---

<br>

## Constant Pool 的生成時機

**Constant Pool 是 Java 編譯器（javac）在編譯階段創建的。**


這個例子展示了哪些內容會進入常量池:

```java
public class ConstantPoolDemo {
    // 1. 字符串字面量 -> 進入常量池
    private static final String GREETING = "Hello, World!";
    
    // 2. 數字字面量 -> 進入常量池
    private static final int NUMBER = 42;
    private static final double PI = 3.14159;
    
    // 3. 類名、方法名、字段名 -> 都會進入常量池
    private String message;
    
    public static void main(String[] args) {
        // 4. 方法調用信息 -> 進入常量池
        System.out.println(GREETING);
        
        // 5. 類引用 -> 進入常量池
        ConstantPoolDemo demo = new ConstantPoolDemo();
        
        // 6. 字段訪問 -> 進入常量池
        demo.message = "Test";
        
        // 7. 方法調用 -> 進入常量池
        demo.printMessage();
    }
    
    public void printMessage() {
        System.out.println(message);
    }
}
```