/**
 * TestThrow.java
 * v0.2.10 測試案例 2: 顯式 throw
 *
 * 預期輸出:
 * 1
 * 2
 */
public class TestThrow {
    public static void main(String[] args) {
        try {
            throwIt();
            System.out.println(0);  // 不應該執行
        } catch (RuntimeException e) {
            System.out.println(1);  // 應輸出 1
        }
        System.out.println(2);  // 應輸出 2
    }

    public static void throwIt() {
        throw new RuntimeException();
    }
}