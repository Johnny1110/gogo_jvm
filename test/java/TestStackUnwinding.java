/**
 * TestStackUnwinding.java
 * v0.2.10 測試案例 3: 棧展開 (Stack Unwinding)
 *
 * 測試異常跨越多層方法傳播
 *
 * 預期輸出:
 * 1
 * 2
 * 3
 * 4
 * 5
 */
public class TestStackUnwinding {
    public static void main(String[] args) {
        try {
            methodA();
            System.out.println(99);  // 不應該執行
        } catch (Exception e) {
            System.out.println(4);  // 應輸出 4
        }
        System.out.println(5);  // 應輸出 5
    }

    public static void methodA() {
        System.out.println(1);  // 應輸出 1
        methodB();
        System.out.println(99);  // 不應該執行
    }

    public static void methodB() {
        System.out.println(2);  // 應輸出 2
        methodC();
        System.out.println(99);  // 不應該執行
    }

    public static void methodC() {
        System.out.println(3);  // 應輸出 3
        throw new RuntimeException();
    }
}