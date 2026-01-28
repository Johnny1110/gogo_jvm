/**
 * v0.3.3 Test: Array clone()
 * 
 * Arrays are always Cloneable (JLS requirement).
 * Array.clone() creates a shallow copy:
 * - For primitive arrays: values are copied
 * - For reference arrays: references are copied (shallow)
 * 
 * Expected output:
 * 1
 * 2
 * 3
 * 4
 * 5
 * 99
 */
public class TestArrayClone {
    public static void main(String[] args) {
        // Test 1-3: Primitive array clone
        testPrimitiveArrayClone();
        
        // Test 4-5: Reference array clone (shallow copy behavior)
        testReferenceArrayClone();
        
        System.out.println(99);
    }
    
    static void testPrimitiveArrayClone() {
        int[] original = {1, 2, 3, 4, 5};
        int[] cloned = original.clone();
        
        // Test 1: Clone is a new array (different reference)
        if (original != cloned) {
            System.out.println(1);
        }
        
        // Test 2: Values are copied correctly
        boolean allSame = true;
        for (int i = 0; i < original.length; i++) {
            if (original[i] != cloned[i]) {
                allSame = false;
                break;
            }
        }
        if (allSame) {
            System.out.println(2);
        }
        
        // Test 3: Modifying clone doesn't affect original
        cloned[0] = 100;
        if (original[0] == 1) {
            System.out.println(3);
        }
    }
    
    static void testReferenceArrayClone() {
        SimpleBean[] original = new SimpleBean[2];
        original[0] = new SimpleBean("TestA");
        original[1] = new SimpleBean("TestB");
        
        SimpleBean[] cloned = original.clone();
        
        // Test 4: Clone is a new array
        if (original != cloned) {
            System.out.println(4);
        }
        
        // Test 5: Shallow copy - references point to same objects
        // Modifying the object through cloned array affects original
        // (This is the expected behavior of shallow copy)
        if (original[0] == cloned[0]) {
            System.out.println(5);
        }
    }
}
