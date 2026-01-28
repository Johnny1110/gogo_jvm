/**
 * v0.3.3 Test: Object.clone() with Cloneable interface
 * 
 * Tests:
 * 1. Cloning creates a new object (different reference)
 * 2. Cloned object has same field values (shallow copy)
 * 3. Modifying clone doesn't affect original
 * 
 * Expected output:
 * 1
 * 2
 * 3
 * 99
 */
public class TestClone implements Cloneable {
    public int value;
    public String name;
    
    public TestClone(int v, String n) {
        this.value = v;
        this.name = n;
    }
    
    @Override
    public Object clone() throws CloneNotSupportedException {
        return super.clone();
    }
    
    public static void main(String[] args) {
        try {
            TestClone original = new TestClone(42, "original");
            TestClone cloned = (TestClone) original.clone();
            
            // Test 1: clone() creates a new object (different reference)
            if (original != cloned) {
                System.out.println(1);
            }
            
            // Test 2: Field values are the same
            if (original.value == cloned.value) {
                System.out.println(2);
            }
            
            // Test 3: Modifying clone doesn't affect original
            cloned.value = 100;
            if (original.value == 42) {
                System.out.println(3);
            }
            
            System.out.println(99);
            
        } catch (CloneNotSupportedException e) {
            System.out.println(0);  // Should NOT reach here
        }
    }
}
