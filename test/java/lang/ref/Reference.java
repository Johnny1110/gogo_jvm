package java.lang.ref;

/**
 * Abstract base class for reference objects.  This class defines the
 * operations common to all reference objects.  Because reference objects are
 * implemented in close cooperation with the garbage collector, this class may
 * not be subclassed directly.
 * @param <T> the type of the referent
 *
 * @author   Mark Reinhold
 * @since    1.2
 * @sealedGraph
 */
public abstract class Reference<T> {

    private T referent;

    volatile ReferenceQueue<? super T> queue;

    /* The link in a ReferenceQueue's list of Reference objects.
         *
         * When registered: null
         *        enqueued: next element in queue (or this if last)
         *        dequeued: this (marking FinalReferences as inactive)
         *    unregistered: null
    */
    volatile Reference next;


    Reference(T referent, ReferenceQueue<? super T> queue) {
        this.referent = referent;
        this.queue = (queue == null) ? ReferenceQueue.NULL : queue;
    }


    Reference(T referent) {
        this(referent, null);
    }

    /**
     * Tests if the referent of this reference object is {@code obj}.
     * Using a {@code null} {@code obj} returns {@code true} if the
     * reference object has been cleared.
     *
     * @param  obj the object to compare with this reference object's referent
     * @return {@code true} if {@code obj} is the referent of this reference object
     * @since 16
     */
    public final boolean refersTo(T obj) {
        return refersToImpl(obj);
    }

    /* Implementation of refersTo(), overridden for phantom references.
     * This method exists only to avoid making refersTo0() virtual. Making
     * refersTo0() virtual has the undesirable effect of C2 often preferring
     * to call the native implementation over the intrinsic.
     */
    boolean refersToImpl(T obj) {
        return refersTo0(obj);
    }

    private native boolean refersTo0(Object o);

    public T get() {
            return this.referent;
    }

    public void clear() {
            clear0();
        }

    /* Implementation of clear(), also used by enqueue().  A simple
    * assignment of the referent field won't do for some garbage
    * collectors.
    */
    private native void clear0();

    public boolean enqueue() {
       clear0();               // Intentionally clear0() rather than clear()
       return this.queue.enqueue(this);
    }
}