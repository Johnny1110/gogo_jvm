/**
 * v0.3.1 Test - Simple Reflection Test
 *
 * 最簡單的反射測試，用於驗證基本功能
 */
public class TestGetClass {
    public static void main(String[] args) {
        // 建立物件
        TestGetClass obj = new TestGetClass();

        // 取得 Class 物件
        Class clazz = obj.getClass();

        // 印出類別名稱
        String name = clazz.getName();
        System.out.println(name);

        // 取得父類
        Class parent = clazz.getSuperclass();
        System.out.println(parent.getName());

        // 測試完成
        System.out.println(99);
    }
}