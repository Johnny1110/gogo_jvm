import java.lang.ref.WeakReference;
import java.lang.ref.ReferenceQueue;
import java.lang.ref.Reference;

/**
 * v0.3.2 Test - ReferenceQueue (Java 8 Compatible)
 *
 * Tests ReferenceQueue functionality:
 * 1. Empty queue poll() returns null
 * 2. enqueue() adds reference to queue
 * 3. poll() removes and returns reference
 * 4. Multiple enqueue/poll operations
 */
public class TestReferenceQueue {
    public static void main(String[] args) {
        System.out.println(1);  // Test start marker

        // ============================================================
        // Test 1: Empty queue poll() returns null
        // ============================================================
        ReferenceQueue<Object> queue = new ReferenceQueue<>();

        Reference<?> ref = queue.poll();
        if (ref == null) {
            System.out.println(11);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 2: Create WeakReference with queue
        // ============================================================
        Object obj1 = new Object();
        WeakReference<Object> weakRef1 = new WeakReference<>(obj1, queue);

        // Queue should still be empty (not enqueued yet)
        if (queue.poll() == null) {
            System.out.println(12);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 3: Manual enqueue and poll
        // ============================================================
        weakRef1.clear();
        weakRef1.enqueue();

        Reference<?> polled1 = queue.poll();
        if (polled1 != null) {
            System.out.println(13);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 4: Queue is empty after poll
        // ============================================================
        if (queue.poll() == null) {
            System.out.println(14);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        // ============================================================
        // Test 5: Multiple references in queue
        // ============================================================
        ReferenceQueue<Object> queue2 = new ReferenceQueue<>();

        Object obj2 = new Object();
        Object obj3 = new Object();
        Object obj4 = new Object();

        WeakReference<Object> ref2 = new WeakReference<>(obj2, queue2);
        WeakReference<Object> ref3 = new WeakReference<>(obj3, queue2);
        WeakReference<Object> ref4 = new WeakReference<>(obj4, queue2);

        // Enqueue all three
        ref2.clear();
        ref2.enqueue();
        ref3.clear();
        ref3.enqueue();
        ref4.clear();
        ref4.enqueue();

        // Poll all three
        int count = 0;
        while (queue2.poll() != null) {
            count++;
        }

        if (count == 3) {
            System.out.println(15);  // Pass
        } else {
            System.out.println(0);   // Fail: expected 3, got count
        }

        // ============================================================
        // Test 6: Cannot enqueue twice
        // ============================================================
        ReferenceQueue<Object> queue3 = new ReferenceQueue<>();
        Object obj5 = new Object();
        WeakReference<Object> ref5 = new WeakReference<>(obj5, queue3);

        ref5.clear();
        boolean first = ref5.enqueue();   // Should succeed
        boolean second = ref5.enqueue();  // Should fail (already enqueued)

        // After first poll, state becomes inactive, can't enqueue again
        queue3.poll();
        boolean third = ref5.enqueue();   // Should fail (inactive)

        if (first && !second && !third) {
            System.out.println(16);  // Pass
        } else {
            System.out.println(0);   // Fail
        }

        System.out.println(99);  // Test end marker
    }
}