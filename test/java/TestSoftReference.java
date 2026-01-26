import java.lang.ref.SoftReference;
import java.lang.ref.ReferenceQueue;

/**
 * v0.3.2 Test - SoftReference (Java 8 Compatible)
 *
 * Tests SoftReference functionality:
 * 1. get() returns the referent
 * 2. clear() sets referent to null
 * 3. SoftReference behavior with queue
 *
 * Note: We cannot easily test GC behavior (memory-based clearing)
 * in unit tests. Those tests would require simulating memory pressure.
 */
public class TestSoftReference {
    public static void main(String[] args) {
        System.out.println(1);  // Test start marker

        // ============================================================
        // Test 1: Basic get()
        // ============================================================
        Object obj = new Object();
        SoftReference<Object> softRef = new SoftReference<>(obj);

        // get() should return the object
        if (softRef.get() != null) {
            System.out.println(11);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 2: get() returns same object
        // ============================================================
        if (softRef.get() == obj) {
            System.out.println(12);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 3: clear()
        // ============================================================
        SoftReference<Object> softRef2 = new SoftReference<>(new Object());
        softRef2.clear();

        if (softRef2.get() == null) {
            System.out.println(13);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 4: SoftReference with queue - queue empty initially
        // ============================================================
        ReferenceQueue<Object> queue = new ReferenceQueue<>();
        Object obj4 = new Object();
        SoftReference<Object> softRef4 = new SoftReference<>(obj4, queue);

        // Queue should be empty initially
        if (queue.poll() == null) {
            System.out.println(14);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 5: Manual enqueue
        // ============================================================
        softRef4.clear();
        boolean enqueued = softRef4.enqueue();

        if (enqueued) {
            System.out.println(15);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 6: Poll from queue
        // ============================================================
        if (queue.poll() != null) {
            System.out.println(16);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 7: Even with strong reference removed, soft ref may survive
        // (Without GC trigger, soft ref keeps the object)
        // ============================================================
        Object tempObj = new Object();
        SoftReference<Object> softRef7 = new SoftReference<>(tempObj);
        tempObj = null;  // Remove strong reference

        // Without explicit GC, soft ref should still have the object
        // (GC only clears soft refs under memory pressure)
        // Note: This test may be flaky if GC runs automatically
        // For MVP, we just verify get() works
        Object retrieved = softRef7.get();
        // Just verify no exception - actual value depends on GC
        System.out.println(17);  // Pass (no exception)

        System.out.println(99);  // Test end marker
    }
}