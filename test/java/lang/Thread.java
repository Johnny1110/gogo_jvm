package java.lang;

public class Thread implements Runnable {

    public static final int MIN_PRIORITY = 1;

    /**
     * The default priority that is assigned to a thread.
     */
    public static final int NORM_PRIORITY = 5;

    /**
     * The maximum priority that a thread can have.
     */
    public static final int MAX_PRIORITY = 10;

    // Additional fields for platform threads.
    // All fields, except task, are accessed directly by the VM.
    private static class FieldHolder {
        final ThreadGroup group;
        final Runnable task;
        final long stackSize;
        volatile int priority;
        volatile boolean daemon;
        volatile int threadStatus;

        FieldHolder(ThreadGroup group,
                    Runnable task,
                    long stackSize,
                    int priority,
                    boolean daemon) {
            this.group = group;
            this.task = task;
            this.stackSize = stackSize;
            this.priority = priority;
            if (daemon)
                this.daemon = true;
        }
    }

    private final FieldHolder holder;

    public Thread(Runnable task) {
        this(null, null, 0, task, 0);
    }

    Thread(ThreadGroup g, String name, int characteristics, Runnable task, long stackSize) {

        Thread parent = currentThread();
        boolean attached = (parent == this);   // primordial or JNI attached
        if (attached) {
            this.holder = new FieldHolder(g, task, stackSize, NORM_PRIORITY, false);
        } else {
            this.holder = new FieldHolder(g, task, stackSize, NORM_PRIORITY, false);
            // TODO
        }
    }

    public void start() {
        start0();
    }

    private native void start0();

    public final void join() throws InterruptedException {
        join(0);
    }

    public final void join(long millis) throws InterruptedException {
        if (millis < 0) {
            throw new IllegalArgumentException("timeout value is negative");
        }

    }

    public static void sleep(long millis) throws InterruptedException {
        if (millis < 0) {
            throw new IllegalArgumentException("timeout value is negative");
        }

        sleep0(millis);
    }

    private static native void sleep0(long millis) throws InterruptedException;

    public static native Thread currentThread();

    @Override
    public void run() {
        Runnable task = holder.task;
        if (task != null) {
            task.run();
        }
    }

}