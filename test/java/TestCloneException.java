/**
 * v0.3.3 Test: Object.clone() throws CloneNotSupportedException
 * 
 * When a class does NOT implement Cloneable, calling clone() should
 * throw CloneNotSupportedException.
 * 
 * Expected output:
 * 1
 * 99
 */
public class TestCloneException {
    // This class does NOT implement Cloneable
    
    public static void main(String[] args) {
        TestCloneException obj = new TestCloneException();
        
        try {
            // This should throw CloneNotSupportedException
            // because TestCloneException doesn't implement Cloneable
            Object cloned = obj.clone();
            
            // Should NOT reach here
            System.out.println(0);
            
        } catch (CloneNotSupportedException e) {
            // Expected: clone() should throw exception
            System.out.println(1);
        }
        
        System.out.println(99);
    }
    
    // Override clone to make it public (so we can call it)
    // but don't implement Cloneable
    @Override
    public Object clone() throws CloneNotSupportedException {
        return super.clone();  // This should throw because we're not Cloneable
    }
}
