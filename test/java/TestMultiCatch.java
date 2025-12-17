/**
 * TestMultiCatch.java
 * v0.2.10 測試案例 4: 多層 catch
 *
 * 測試多個 catch 區塊的匹配
 *
 * 預期輸出:
 * 10
 * 20
 * 30
 */
public class TestMultiCatch {
    public static void main(String[] args) {
        test(1);  // ArithmeticException
        test(2);  // NullPointerException
        test(3);  // 其他 Exception
    }

    public static void test(int type) {
        try {
            if (type == 1) {
                int x = 1 / 0;  // ArithmeticException
            } else if (type == 2) {
                throw new NullPointerException();
            } else {
                throw new RuntimeException();
            }
        } catch (ArithmeticException e) {
            System.out.println(10);
        } catch (NullPointerException e) {
            System.out.println(20);
        } catch (Exception e) {
            System.out.println(30);
        }
    }
}