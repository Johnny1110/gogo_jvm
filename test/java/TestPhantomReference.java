import java.lang.ref.PhantomReference;
import java.lang.ref.ReferenceQueue;

/**
 * v0.3.2 Test - PhantomReference (Java 8 Compatible)
 *
 * Tests PhantomReference functionality:
 * 1. get() always returns null (key characteristic!)
 * 2. Must be created with a ReferenceQueue
 * 3. enqueue() and poll() work correctly
 *
 * PhantomReference is the weakest reference type:
 * - You can NEVER get the referent back
 * - It's only used to track WHEN an object is reclaimed
 * - Commonly used for resource cleanup (replacing finalize())
 */
public class TestPhantomReference {
    public static void main(String[] args) {
        System.out.println(1);  // Test start marker

        // ============================================================
        // Test 1: get() always returns null (CRITICAL TEST!)
        // ============================================================
        ReferenceQueue<Object> queue = new ReferenceQueue<>();
        Object obj = new Object();
        PhantomReference<Object> phantomRef = new PhantomReference<>(obj, queue);

        // THE KEY PROPERTY: get() ALWAYS returns null for PhantomReference
        if (phantomRef.get() == null) {
            System.out.println(11);  // Pass - this is the expected behavior!
        } else {
            System.out.println(0);   // Fail - get() should NEVER return non-null
        }

        // ============================================================
        // Test 2: get() returns null even when object is still alive
        // ============================================================
        // obj is still strongly referenced here
        if (obj != null && phantomRef.get() == null) {
            System.out.println(12);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 3: Queue empty initially
        // ============================================================
        if (queue.poll() == null) {
            System.out.println(13);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 4: Manual enqueue
        // ============================================================
        // Note: In real usage, GC enqueues phantom refs after finalization
        // For testing, we enqueue manually
        boolean enqueued = phantomRef.enqueue();

        if (enqueued) {
            System.out.println(14);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 5: poll() from queue after enqueue
        // ============================================================
        Object polled = queue.poll();
        if (polled != null) {
            System.out.println(15);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 6: Queue empty after poll
        // ============================================================
        if (queue.poll() == null) {
            System.out.println(16);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 7: get() still returns null after all operations
        // ============================================================
        if (phantomRef.get() == null) {
            System.out.println(17);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 8: clear() works (even though get() is always null)
        // ============================================================
        ReferenceQueue<Object> queue2 = new ReferenceQueue<>();
        Object obj2 = new Object();
        PhantomReference<Object> phantomRef2 = new PhantomReference<>(obj2, queue2);

        phantomRef2.clear();  // Should not throw

        // get() is still null (as always for PhantomReference)
        if (phantomRef2.get() == null) {
            System.out.println(18);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        System.out.println(99);  // Test end marker
    }
}