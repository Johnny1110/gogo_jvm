package runtime

// ============================================================
// Thread State Definition (v0.4.0)
// ============================================================

type ThreadState int

const (
	// THREAD_NEW Thread created but not run
	// Thread t = new Thread(); // State is NEW
	THREAD_NEW ThreadState = iota

	// THREAD_RUNNABLE thread is running or waiting for CPU
	THREAD_RUNNABLE

	// THREAD_BLOCKED thread is waiting for get Monitor（synchronized）
	// try to enter synchronized but lock been occupied by other thread.
	THREAD_BLOCKED

	// THREAD_WAITING waiting for ever
	// Object.wait() / Thread.join() / LockSupport.park()
	THREAD_WAITING

	// THREAD_TIMED_WAITING waiting in a limited time
	// Thread.sleep(n) / Object.wait(n) / Thread.join(n)
	THREAD_TIMED_WAITING

	// THREAD_TERMINATED thread already done
	// run() end or threw exception
	THREAD_TERMINATED
)

func (s ThreadState) String() string {
	switch s {
	case THREAD_NEW:
		return "NEW"
	case THREAD_RUNNABLE:
		return "RUNNABLE"
	case THREAD_BLOCKED:
		return "BLOCKED"
	case THREAD_WAITING:
		return "WAITING"
	case THREAD_TIMED_WAITING:
		return "TIMED_WAITING"
	case THREAD_TERMINATED:
		return "TERMINATED"
	default:
		return "UNKNOWN"
	}
}

// ToJavaOrdinal convert to Java Thread.State ordinal value
// for setup java.lang.Thread's threadStatus field
func (s ThreadState) ToJavaOrdinal() int {
	// NEW=0, RUNNABLE=1, BLOCKED=2, WAITING=3, TIMED_WAITING=4, TERMINATED=5
	return int(s)
}

// FromJavaOrdinal from Java ordinal convert to ThreadState
func FromJavaOrdinal(ordinal int) ThreadState {
	if ordinal < 0 || ordinal > 5 {
		panic("invalid enum int for java thread state")
	}

	return ThreadState(ordinal)
}

// ============================================================
// State Transition Validation
// ============================================================

// CanTransitionTo check state transform is legal or not
// this is for debug mode
func (s ThreadState) CanTransitionTo(target ThreadState) bool {
	switch s {
	case THREAD_NEW:
		// new only can trans runnable
		return target == THREAD_RUNNABLE

	case THREAD_RUNNABLE:
		// runnable can trans to (blocked, waiting, timed_waiting, terminated)
		return target == THREAD_BLOCKED ||
			target == THREAD_WAITING ||
			target == THREAD_TIMED_WAITING ||
			target == THREAD_TERMINATED

	case THREAD_BLOCKED:
		// blocked can trans to runnable (after get monitor)
		return target == THREAD_RUNNABLE

	case THREAD_WAITING, THREAD_TIMED_WAITING:
		// waiting can trans to runnable or blocked
		return target == THREAD_RUNNABLE || target == THREAD_BLOCKED

	case THREAD_TERMINATED:
		// terminated can't do anything.
		return false

	default:
		return false
	}
}

// IsAlive check thread is alive or not
// alive definition: thread already start() and not terminated yet.
func (s ThreadState) IsAlive() bool {
	return s != THREAD_NEW && s != THREAD_TERMINATED
}
