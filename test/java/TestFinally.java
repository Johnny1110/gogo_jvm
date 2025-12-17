/**
 * TestFinally.java
 * v0.2.10 測試案例 5: finally 區塊
 *
 * 測試 finally 無論是否發生異常都會執行
 *
 * 預期輸出:
 * 1
 * 2
 * 3
 * 4
 */
public class TestFinally {
    public static void main(String[] args) {
        try {
            System.out.println(1);  // 應輸出 1
            int x = 1 / 0;          // 觸發異常
            System.out.println(99); // 不應該執行
        } catch (ArithmeticException e) {
            System.out.println(2);  // 應輸出 2
        } finally {
            System.out.println(3);  // 應輸出 3 (finally 一定執行)
        }
        System.out.println(4);      // 應輸出 4
    }
}