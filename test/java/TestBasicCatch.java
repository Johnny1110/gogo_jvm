/**
 * TestBasicCatch.java
 * v0.2.10 測試案例 1: 基本 try-catch
 */
public class TestBasicCatch {
    public static void main(String[] args) {
        try {
            int x = 1 / 0;  // 觸發 ArithmeticException
            System.out.println(0);  // 不應該執行
        } catch (ArithmeticException e) {
            System.out.println(1);  // 應輸出 1
        }
        System.out.println(2);  // 應輸出 2
    }
}