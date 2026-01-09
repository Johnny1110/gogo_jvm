/**
 * v0.3.1 Test - Basic Reflection
 *
 * 測試 java.lang.Class 的基本功能
 */
public class TestReflection {
    public static void main(String[] args) {
        // ============================================================
        // Test 1: Object.getClass() 和 Class.getName()
        // ============================================================
        System.out.println(1);  // 標記測試開始

        TestReflection obj = new TestReflection();
        Class clazz = obj.getClass();

        // 印出類別名稱
        // 預期：TestReflection
        String name = clazz.getName();
        System.out.println(name);

        // ============================================================
        // Test 2: Class.getSuperclass()
        // ============================================================
        System.out.println(2);  // 標記測試開始

        Class parent = clazz.getSuperclass();
        if (parent != null) {
            System.out.println(parent.getName());
            // 預期：java.lang.Object
        } else {
            System.out.println(0);  // 錯誤
        }

        // ============================================================
        // Test 3: isPrimitive(), isArray(), isInterface()
        // ============================================================
        System.out.println(3);  // 標記測試開始

        // TestReflection 不是基本類型、陣列、介面
        if (!clazz.isPrimitive()) {
            System.out.println(31);  // 正確
        }
        if (!clazz.isArray()) {
            System.out.println(32);  // 正確
        }
        if (!clazz.isInterface()) {
            System.out.println(33);  // 正確
        }

        // ============================================================
        // Test 4: 陣列的 getClass()
        // ============================================================
        System.out.println(4);  // 標記測試開始

        int[] arr = new int[5];
        Class arrClass = arr.getClass();
        System.out.println(arrClass.getName());
        // 預期：[I

        if (arrClass.isArray()) {
            System.out.println(41);  // 正確
        }

        // ============================================================
        // Test 5: .class 字面量
        // ============================================================
        System.out.println(5);  // 標記測試開始

        Class strClass = String.class;
        System.out.println(strClass.getName());
        // 預期：java.lang.String

        // ============================================================
        // 測試完成
        // ============================================================
        System.out.println(99);  // 測試結束標記
    }
}