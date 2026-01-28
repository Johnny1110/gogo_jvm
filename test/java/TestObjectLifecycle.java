/**
 * v0.3.3 Comprehensive Test: Object Lifecycle
 * 
 * This test covers all Object lifecycle methods implemented in v0.3.3:
 * - hashCode() (from v0.3.0)
 * - equals() (default reference equality)
 * - clone() (with Cloneable check)
 * 
 * Expected output:
 * === hashCode tests ===
 * 1
 * 2
 * === equals tests ===
 * 3
 * 4
 * 5
 * === clone tests ===
 * 6
 * 7
 * 8
 * === All tests passed ===
 * 99
 */
public class TestObjectLifecycle {
    public static void main(String[] args) {
        System.out.println("=== hashCode tests ===");
        testHashCode();
        
        System.out.println("=== equals tests ===");
        testEquals();
        
        System.out.println("=== clone tests ===");
        testClone();
        
        System.out.println("=== All tests passed ===");
        System.out.println(99);
    }
    
    // ============================================
    // hashCode() tests
    // ============================================
    static void testHashCode() {
        Object obj = new Object();
        
        // Test 1: hashCode is consistent
        int h1 = obj.hashCode();
        int h2 = obj.hashCode();
        if (h1 == h2) {
            System.out.println(1);
        }
        
        // Test 2: hashCode is positive
        if (h1 > 0) {
            System.out.println(2);
        }
    }
    
    // ============================================
    // equals() tests
    // ============================================
    static void testEquals() {
        Object a = new Object();
        Object b = new Object();
        Object c = a;  // same reference
        
        // Test 3: Reflexive - object equals itself
        if (a.equals(a)) {
            System.out.println(3);
        }
        
        // Test 4: Same reference -> equals true
        if (a.equals(c)) {
            System.out.println(4);
        }
        
        // Test 5: Different objects -> equals false (default behavior)
        if (!a.equals(b)) {
            System.out.println(5);
        }
    }
    
    // ============================================
    // clone() tests
    // ============================================
    static void testClone() {
        try {
            // Use an array (arrays are always Cloneable)
            int[] original = {1, 2, 3};
            int[] cloned = original.clone();
            
            // Test 6: Clone is a new object
            if (original != cloned) {
                System.out.println(6);
            }
            
            // Test 7: Clone has same values
            boolean same = true;
            for (int i = 0; i < original.length; i++) {
                if (original[i] != cloned[i]) {
                    same = false;
                }
            }
            if (same) {
                System.out.println(7);
            }
            
            // Test 8: Modifying clone doesn't affect original
            cloned[0] = 999;
            if (original[0] == 1) {
                System.out.println(8);
            }
            
        } catch (Exception e) {
            System.out.println(0);  // Should not reach here
        }
    }
}
