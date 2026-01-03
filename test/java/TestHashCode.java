/**
 * TestHashCode.java - v0.3.0 測試
 *
 * 測試目標：
 * 1. Object.hashCode() 返回正整數
 * 2. 同一物件多次呼叫返回相同值（一致性）
 * 3. 不同物件返回不同值（高機率）
 */
public class TestHashCode {
    public static void main(String[] args) {
        Object obj1 = new Object();
        Object obj2 = new Object();

        int hash1 = obj1.hashCode();
        int hash2 = obj2.hashCode();
        int hash1Again = obj1.hashCode();

        // Test 1: hash1 應該等於 hash1Again（一致性）
        if (hash1 == hash1Again) {
            System.out.println("hash1 == hash1Again");
        }

        // Test 2: hash1 應該（很大機率）不等於 hash2（唯一性）
        if (hash1 != hash2) {
            System.out.println("hash1 != hash2");
        }

        // Test 3: hashCode 應該是正數
        if (hash1 > 0 && hash2 > 0) {
            System.out.println("hash1 and hash2 both are gt 0");
        }

        // 輸出實際的 hash 值供驗證
        System.out.println(hash1);
        System.out.println(hash2);

        System.out.println("DONE");
    }
}