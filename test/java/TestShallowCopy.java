/**
 * v0.3.3 Test: Shallow copy behavior demonstration
 * 
 * Object.clone() performs shallow copy by default:
 * - Primitive fields: values are copied
 * - Reference fields: references are copied (not the objects)
 * 
 * This test demonstrates the shallow copy behavior.
 * 
 * Expected output:
 * 1
 * 2
 * 3
 * 4
 * 99
 */
public class TestShallowCopy implements Cloneable {
    public int primitiveValue;
    public int[] arrayValue;  // Reference to array
    
    public TestShallowCopy(int pv, int[] av) {
        this.primitiveValue = pv;
        this.arrayValue = av;
    }
    
    @Override
    public Object clone() throws CloneNotSupportedException {
        return super.clone();
    }
    
    public static void main(String[] args) {
        try {
            int[] sharedArray = {10, 20, 30};
            TestShallowCopy original = new TestShallowCopy(42, sharedArray);
            TestShallowCopy cloned = (TestShallowCopy) original.clone();
            
            // Test 1: Primitive value is copied (independent)
            cloned.primitiveValue = 100;
            if (original.primitiveValue == 42) {
                System.out.println(1);
            }
            
            // Test 2: Reference field points to SAME array (shallow copy)
            if (original.arrayValue == cloned.arrayValue) {
                System.out.println(2);
            }
            
            // Test 3: Modifying array through clone affects original
            // (because they share the same array)
            cloned.arrayValue[0] = 999;
            if (original.arrayValue[0] == 999) {
                System.out.println(3);
            }
            
            // Test 4: Assigning new array to clone doesn't affect original
            cloned.arrayValue = new int[]{1, 2, 3};
            if (original.arrayValue[0] == 999) {  // Original still has the old array
                System.out.println(4);
            }
            
            System.out.println(99);
            
        } catch (CloneNotSupportedException e) {
            System.out.println(0);
        }
    }
}
