/**
 * v0.3.1 Test - Array Reflection Test
 *
 * 測試陣列類型的反射功能
 */
public class TestArrayClass {
    public static void main(String[] args) {
        // 測試 1: int[] 的 Class
        int[] intArr = new int[3];
        Class intArrClass = intArr.getClass();
        System.out.println(intArrClass.getName());  // 預期: [I

        // 測試 2: isArray()
        if (intArrClass.isArray()) {
            System.out.println(1);  // 正確
        } else {
            System.out.println(0);  // 錯誤
        }

        // 測試 3: getSuperclass()
        // 陣列的父類是 java.lang.Object
        Class parent = intArrClass.getSuperclass();
        System.out.println(parent.getName());  // 預期: java.lang.Object

        // 測試 4: String[] 的 Class
        String[] strArr = new String[2];
        Class strArrClass = strArr.getClass();
        System.out.println(strArrClass.getName());  // 預期: [Ljava.lang.String;

        // 測試 5: int[][][] 三維陣列
        int[][][] int2dArr = new int[2][3][3];
        Class int2dArrClass = int2dArr.getClass();
        System.out.println(int2dArrClass.getName());  // 預期: [[[I

        // 測試完成
        System.out.println(99);
    }
}