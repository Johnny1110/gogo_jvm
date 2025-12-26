// ============================================================
// TestInterface.java - 基本介面調用測試
// ============================================================
// 測試目標：
//   1. 基本的 invokeinterface 指令
//   2. 多型調用（同一介面，不同實現）
//
// 預期輸出：
//   1
//   2
//   99

interface Greeting {
    void sayHello();
}

class EnglishGreeting implements Greeting {
    public void sayHello() {
        System.out.println("Hello");
    }
}

class ChineseGreeting implements Greeting {
    public void sayHello() {
        System.out.println("你好");
    }
}

public class TestInterface {
    public static void main(String[] args) {
        Greeting g1 = new EnglishGreeting();
        Greeting g2 = new ChineseGreeting();

        g1.sayHello();
        g2.sayHello();

        System.out.println("Test done");
    }
}