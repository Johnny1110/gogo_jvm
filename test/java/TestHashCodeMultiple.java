/**
 * TestHashCodeMultiple.java - v0.3.0 測試
 *
 * 測試多個物件的 hashCode 分佈
 */
public class TestHashCodeMultiple {
    public static void main(String[] args) {
        Object[] objects = new Object[5];

        // 建立 5 個物件
        for (int i = 0; i < 5; i++) {
            objects[i] = new Object();
        }

        // 印出每個物件的 hashCode
        for (int i = 0; i < 5; i++) {
            System.out.println(objects[i].hashCode());
        }

        // 驗證一致性：再次取得 hashCode 應該相同
        int hash0 = objects[0].hashCode();
        int hash0Again = objects[0].hashCode();

        if (hash0 == hash0Again) {
            System.out.println("passed!");  // 一致性測試通過
        } else {
            System.out.println("not pass!");
        }
    }
}