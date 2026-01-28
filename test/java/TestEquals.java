/**
 * v0.3.3 Test: Object.equals() default behavior
 * 
 * The default Object.equals() compares references (identity equality).
 * Two objects are equal only if they are the same object.
 * 
 * Expected output:
 * 1
 * 2
 * 3
 * 4
 * 99
 */
public class TestEquals {
    public static void main(String[] args) {
        // Test 1: Same reference should be equal
        Object obj1 = new Object();
        Object obj3 = obj1;  // Same reference
        
        if (obj1.equals(obj3)) {
            System.out.println(1);  // Should print
        }
        
        // Test 2: Different objects should NOT be equal (default behavior)
        Object obj2 = new Object();
        
        if (!obj1.equals(obj2)) {
            System.out.println(2);  // Should print
        }
        
        // Test 3: equals(null) should return false
        if (!obj1.equals(null)) {
            System.out.println(3);  // Should print
        }
        
        // Test 4: Reflexive property - object equals itself
        if (obj1.equals(obj1)) {
            System.out.println(4);  // Should print
        }
        
        System.out.println(99);  // Test completed
    }
}
