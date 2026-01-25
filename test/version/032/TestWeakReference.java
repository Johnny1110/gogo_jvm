import java.lang.ref.WeakReference;
import java.lang.ref.ReferenceQueue;

/**
 * v0.3.2 Test - Basic WeakReference (Java 8 Compatible)
 *
 * Tests basic WeakReference functionality:
 * 1. get() returns the referent
 * 2. clear() sets referent to null
 * 3. enqueue() and poll() work correctly
 *
 * Note: isEnqueued() is package-private in Java 8, so we don't use it.
 */
public class TestWeakReference {
    public static void main(String[] args) {
        System.out.println(1);  // Test start marker

        // ============================================================
        // Test 1: Basic get()
        // ============================================================
        Object obj = new Object();
        WeakReference<Object> weakRef = new WeakReference<>(obj);

        System.out.println("check point 1");

        // get() should return the object while strong ref exists
        if (weakRef.get() != null) {
            System.out.println(11);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 2: get() returns the same object
        // ============================================================
        if (weakRef.get() == obj) {
            System.out.println(12);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 3: clear()
        // ============================================================
        WeakReference<Object> weakRef2 = new WeakReference<>(new Object());
        weakRef2.clear();

        // After clear(), get() should return null
        if (weakRef2.get() == null) {
            System.out.println(13);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 4: WeakReference with queue - enqueue and poll
        // ============================================================
        ReferenceQueue<Object> queue = new ReferenceQueue<>();
        Object obj4 = new Object();
        WeakReference<Object> weakRef4 = new WeakReference<>(obj4, queue);

        // Queue should be empty initially
        if (queue.poll() == null) {
            System.out.println(14);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 5: Manual enqueue()
        // ============================================================
        weakRef4.clear();
        boolean enqueued = weakRef4.enqueue();

        if (enqueued) {
            System.out.println(15);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 6: poll() from queue after enqueue
        // ============================================================
        Object polled = queue.poll();
        if (polled != null) {
            System.out.println(16);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 7: Queue should be empty after poll
        // ============================================================
        if (queue.poll() == null) {
            System.out.println(17);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 8: Cannot enqueue twice (second enqueue returns false)
        // ============================================================
        ReferenceQueue<Object> queue2 = new ReferenceQueue<>();
        WeakReference<Object> weakRef5 = new WeakReference<>(new Object(), queue2);
        weakRef5.clear();

        boolean first = weakRef5.enqueue();   // Should succeed
        boolean second = weakRef5.enqueue();  // Should fail (already enqueued)

        if (first && !second) {
            System.out.println(18);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        System.out.println(99);  // Test end marker
    }
}