/**
 * v0.3.1 Test - Class Literal Test
 *
 * 測試 .class 字面量（LDC 指令支援 Class 常量）
 */
public class TestClassLiteral {
    public static void main(String[] args) {
        // 測試 1: 基本的 .class 字面量
        Class strClass = String.class;
        System.out.println(strClass.getName());  // 預期: java.lang.String

        // 測試 2: Object.class
        Class objClass = Object.class;
        System.out.println(objClass.getName());  // 預期: java.lang.Object

        // 測試 3: 陣列的 .class 字面量
        Class intArrClass = int[].class;
        System.out.println(intArrClass.getName());  // 預期: [I

        // 測試 4: 自己的 .class
        Class myClass = TestClassLiteral.class;
        System.out.println(myClass.getName());  // 預期: TestClassLiteral

        // 測試 5: 比較兩種方式取得的 Class 是否相同
        TestClassLiteral obj = new TestClassLiteral();
        Class clazz1 = obj.getClass();
        Class clazz2 = TestClassLiteral.class;

        if (clazz1 == clazz2) {
            System.out.println(1);  // 正確：同一個 Class 物件
        } else {
            System.out.println(0);  // 錯誤
        }

        // 測試完成
        System.out.println(99);
    }
}