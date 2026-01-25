package java.lang.ref;

public class ReferenceQueue<T> {

    private static class Null extends ReferenceQueue<Object> {
        public Null() { super(0); }

        @Override
        boolean enqueue(Reference<?> r) {
            return false;
        }
    }

    static final ReferenceQueue<Object> NULL = new Null();
    static final ReferenceQueue<Object> ENQUEUED = new Null();

    private volatile Reference<? extends T> head;
    private long queueLength = 0;

        /**
         * Constructs a new reference-object queue.
         */
        public ReferenceQueue() {

        }

        ReferenceQueue(int dummy) {

        }

    final boolean enqueue0(Reference<? extends T> r) { // must hold lock
        // Check that since getting the lock this reference hasn't already been
        // enqueued (and even then removed)
        ReferenceQueue<?> queue = r.queue;
        if ((queue == NULL) || (queue == ENQUEUED)) {
            return false;
        }

        // Self-loop end, so if a FinalReference it remains inactive.
        r.next = (head == null) ? r : head;
        head = r;
        queueLength++;
        // Update r.queue *after* adding to list, to avoid race
        // with concurrent enqueued checks and fast-path poll().
        // Volatiles ensure ordering.
        r.queue = ENQUEUED;

        signal();
        return true;
    }

    void signal() {
            // TODO
        }

    boolean enqueue(Reference<? extends T> r) { /* Called only by Reference class */
        try {
            return enqueue0(r);
        } finally {
        }
    }

    public Reference<? extends T> poll() {
            if (headIsNull())
                return null;
            //lock.lock();
            try {
                return poll0();
            } finally {
                //lock.unlock();
            }
    }

    final boolean headIsNull() {
            return head == null;
    }

    final Reference<? extends T> poll0() { // must hold lock
        Reference<? extends T> r = head;
        if (r != null) {
            r.queue = NULL;
            // Update r.queue *before* removing from list, to avoid
            // race with concurrent enqueued checks and fast-path
            // poll().  Volatiles ensure ordering.
            Reference<? extends T> rn = r.next;
            // Handle self-looped next as end of list designator.
            head = (rn == r) ? null : rn;
            // Self-loop next rather than setting to null, so if a
            // FinalReference it remains inactive.
            r.next = r;
            queueLength--;

            return r;
        }
        return null;
    }
}