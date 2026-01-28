/**
 * v0.3.3 Test: hashCode() and equals() contract
 * 
 * The contract between hashCode() and equals():
 * 1. If a.equals(b) is true, then a.hashCode() == b.hashCode() must be true
 * 2. If a.equals(b) is false, hashCode() may or may not be different
 * 3. Same object must return same hashCode (consistency)
 * 
 * For Object's default implementation:
 * - equals() uses reference equality (this == obj)
 * - hashCode() returns identity hash code
 * 
 * Expected output:
 * 1
 * 2
 * 3
 * 4
 * 99
 */
public class TestHashCodeEquals {
    public static void main(String[] args) {
        Object obj1 = new Object();
        Object obj2 = new Object();
        Object obj3 = obj1;  // Same reference as obj1
        
        // Test 1: Same object has consistent hashCode
        int hash1 = obj1.hashCode();
        int hash2 = obj1.hashCode();
        if (hash1 == hash2) {
            System.out.println(1);
        }
        
        // Test 2: If equals() is true (same reference), hashCode must be same
        // obj1.equals(obj3) is true because they're the same reference
        if (obj1.equals(obj3) && obj1.hashCode() == obj3.hashCode()) {
            System.out.println(2);
        }
        
        // Test 3: Different objects typically have different hashCodes
        // (Not required by contract, but highly probable for identity hash)
        if (obj1.hashCode() != obj2.hashCode()) {
            System.out.println(3);
        }
        
        // Test 4: hashCode is positive (our implementation guarantees this)
        if (obj1.hashCode() > 0 && obj2.hashCode() > 0) {
            System.out.println(4);
        }
        
        System.out.println(99);
    }
}
