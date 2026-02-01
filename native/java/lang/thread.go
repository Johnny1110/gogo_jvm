package lang

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	jvmruntime "github.com/Johnny1110/gogo_jvm/runtime"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"runtime"
	"time"
)

// ============================================================
// java.lang.Thread Native Methods - v0.4.0
// ============================================================

func init() {
	if global.DebugMode() {
		fmt.Println("@@ Debug - init Native java/lang/Thread")
	}
	// core
	jvmruntime.Register("java/lang/Thread", "start0", "()V", threadStart0)
	jvmruntime.Register("java/lang/Thread", "sleep", "(J)V", threadSleep)
	jvmruntime.Register("java/lang/Thread", "yield", "()V", threadYield)
	jvmruntime.Register("java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", threadCurrentThread)
	jvmruntime.Register("java/lang/Thread", "isAlive", "()Z", threadIsAlive)

	// attrs
	jvmruntime.Register("java/lang/Thread", "setPriority0", "(I)V", threadSetPriority0)
	jvmruntime.Register("java/lang/Thread", "isInterrupted", "(Z)Z", threadIsInterrupted)
	jvmruntime.Register("java/lang/Thread", "interrupt0", "()V", threadInterrupt0)

	// other
	jvmruntime.Register("java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", threadHoldsLock)
	jvmruntime.Register("java/lang/Thread", "getThreads", "()[Ljava/lang/Thread;", threadGetThreads)
	jvmruntime.Register("java/lang/Thread", "dumpThreads", "([Ljava/lang/Thread;)[[Ljava/lang/StackTraceElement;", threadDumpThreads)
}

// ============================================================
// start0 - start a thread
// ============================================================
// Java signature: private native void start0();
// Thread.start() inner native call
// Stack: [this] → []
func threadStart0(frame *jvmruntime.Frame) (ex *heap.Object) {
	// 1. get this (java.lang.Thread)
	this := frame.LocalVars().GetThis()
	if this == nil {
		return genExceptionObj(frame, "java/lang/NullPointerException", "")
	}

	threadObj := this.(*heap.Object)

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - threadStart0: Starting thread object %v\n", threadObj)
	}

	// create jvm thread based on java thread obj
	jvmThread := jvmruntime.CreateThreadFromJavaObject(threadObj)

	// build mirror
	threadObj.SetExtra(jvmThread)         // Java -> JVM
	jvmThread.SetJavaThreadObj(threadObj) // JVM -> Java

	// get thread run method
	threadClass := threadObj.Class().(*method_area.Class)
	runMethod := threadClass.GetMethod("run", "()V")

	if runMethod == nil {
		// should be happened
		return genExceptionObj(frame, "java/lang/NoSuchMethodError", "run() method not found")
	}

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - threadStart0: Found run() method in class %s\n", threadClass.Name())
	}

	jvmThread.Start(runMethod)

	return nil
}

// ============================================================
// sleep - thread sleep
// ============================================================
// Java signature: public static native void sleep(long millis) throws InterruptedException;
//
// stop current thread (milli)
// not release any lock
//
// Stack: [millis (long)] → []
func threadSleep(frame *jvmruntime.Frame) (ex *heap.Object) {
	// 1. egt millis
	millis := frame.LocalVars().GetLong(0)

	if millis < 0 {
		return genExceptionObj(frame, "java/lang/IllegalArgumentException", "timeout value is negative")
	}

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - threadSleep: Sleeping for %d ms\n", millis)
	}

	// 2. get current thread
	currentThread := jvmruntime.CurrentJVMThread()
	if currentThread != nil {
		// set current thread to timed_waiting
		currentThread.SetState(jvmruntime.THREAD_TIMED_WAITING)
	}

	// 3. sleep
	time.Sleep(time.Duration(millis) * time.Millisecond)

	// 4. recover
	if currentThread != nil {
		currentThread.SetState(jvmruntime.THREAD_RUNNABLE)
	}

	// TODO: v0.4.3 檢查中斷狀態，如果被中斷則拋出 InterruptedException
	return nil
}

// ============================================================
// yield - give up CPU
// ============================================================
// Java signature: public static native void yield();
//
// Stack: [] → []
func threadYield(frame *jvmruntime.Frame) (ex *heap.Object) {
	if global.DebugMode() {
		fmt.Println("@@ DEBUG - threadYield: Yielding CPU")
	}

	runtime.Gosched()
	return nil
}

// ============================================================
// currentThread - 取得當前執行緒
// ============================================================
// Java signature: public static native Thread currentThread();
//
// return current running thread
//
// Stack: [] → [Thread]
func threadCurrentThread(frame *jvmruntime.Frame) (ex *heap.Object) {
	// get currentThread
	currentThread := jvmruntime.CurrentJVMThread()

	if currentThread == nil {
		if global.DebugMode() {
			fmt.Println("@@ WARNING - threadCurrentThread: No current thread found!")
		}
		frame.OperandStack().PushRef(nil)
		return
	}

	javaThreadObj := currentThread.JavaThreadObj()

	if javaThreadObj == nil {
		if global.DebugMode() {
			fmt.Println("@@ WARNING - threadCurrentThread: Current thread has no Java object!")
			frame.OperandStack().PushRef(nil)
			return nil
		}
	}

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - threadCurrentThread: Returning thread [%d] %s\n",
			currentThread.ID(), currentThread.Name())
	}

	frame.OperandStack().PushRef(javaThreadObj)
	return nil
}

// ============================================================
// isAlive - check current thread is alive
// ============================================================
// Java signature: public final native boolean isAlive();
//
// Stack: [this] → [boolean]
func threadIsAlive(frame *jvmruntime.Frame) (ex *heap.Object) {
	this := frame.LocalVars().GetThis()
	if this == nil {
		frame.OperandStack().PushInt(0)
		return nil
	}
	threadObj := this.(*heap.Object)

	jvmThread := getJVMThreadFromObject(threadObj)

	if jvmThread == nil {
		// not start yet, or already GC
		frame.OperandStack().PushInt(0)
		return nil
	}

	if jvmThread.IsAlive() {
		frame.OperandStack().PushInt(1)
	} else {
		frame.OperandStack().PushInt(0)
	}

	return nil
}

// ============================================================
// setPriority0
// ============================================================
// Java signature: private native void setPriority0(int newPriority);
//
// Stack: [this, newPriority] → []
func threadSetPriority0(frame *jvmruntime.Frame) (ex *heap.Object) {
	this := frame.LocalVars().GetThis()
	newPriority := frame.LocalVars().GetInt(1)

	if this == nil {
		return nil
	}

	threadObj := this.(*heap.Object)

	// get jvmThread
	jvmThread := getJVMThreadFromObject(threadObj)
	if jvmThread != nil {
		jvmThread.SetPriority(int(newPriority))
	}

	if global.DebugMode() {
		// goroutine not support priority
		fmt.Printf("@@ DEBUG - threadSetPriority0: Set priority to %d (note: no effect on goroutine)\n", newPriority)
	}

	return nil
}

// ============================================================
// isInterrupted - check interrupted
// ============================================================
// Java signature: private native boolean isInterrupted(boolean ClearInterrupted);
//
// Stack: [this, clearInterrupted] → [boolean]
func threadIsInterrupted(frame *jvmruntime.Frame) (ex *heap.Object) {
	this := frame.LocalVars().GetThis()
	clearInterrupted := frame.LocalVars().GetBoolean(1)

	if this == nil {
		frame.OperandStack().PushInt(0)
		return nil
	}
	threadObj := this.(*heap.Object)

	jvmThread := getJVMThreadFromObject(threadObj)
	if jvmThread == nil {
		frame.OperandStack().PushInt(0)
		return nil
	}

	var interrupted bool
	if clearInterrupted {
		interrupted = jvmThread.ClearInterrupt()
	} else {
		interrupted = jvmThread.IsInterrupted()
	}

	frame.OperandStack().PushBoolean(interrupted)

	return nil
}

// ============================================================
// interrupt0
// ============================================================
// Java signature: private native void interrupt0();
//
// Stack: [this] → []
func threadInterrupt0(frame *jvmruntime.Frame) (ex *heap.Object) {
	this := frame.LocalVars().GetThis()
	if this == nil {
		return nil
	}
	threadObj := this.(*heap.Object)

	jvmThread := getJVMThreadFromObject(threadObj)
	if jvmThread != nil {
		jvmThread.Interrupt()
	}

	return nil
}

// ============================================================
// holdsLock - check if holds any lock
// ============================================================
// Java signature: public static native boolean holdsLock(Object obj);
//
// Stack: [obj] → [boolean]
func threadHoldsLock(frame *jvmruntime.Frame) (ex *heap.Object) {
	// TODO: v0.4.1 實現 synchronized 後再完善
	if global.DebugMode() {
		fmt.Println("@@ WARNING - threadHoldsLock: Not implemented (v0.4.1)")
	}
	frame.OperandStack().PushInt(0)
	return nil
}

// ============================================================
// getThreads
// ============================================================
// Java signature: private static native Thread[] getThreads();
//
// Stack: [] → [Thread[]]
func threadGetThreads(frame *jvmruntime.Frame) (ex *heap.Object) {
	// TODO: 返回所有執行緒的陣列
	if global.DebugMode() {
		fmt.Println("@@ WARNING - threadGetThreads: Not fully implemented")
	}
	frame.OperandStack().PushRef(nil)
	return nil
}

// ============================================================
// dumpThreads
// ============================================================
// Java signature: private static native StackTraceElement[][] dumpThreads(Thread[] threads);
//
// Stack: [threads] → [StackTraceElement[][]]
func threadDumpThreads(frame *jvmruntime.Frame) (ex *heap.Object) {
	// TODO: 返回堆疊追蹤
	if global.DebugMode() {
		fmt.Println("@@ WARNING - threadDumpThreads: Not implemented")
	}
	frame.OperandStack().PushRef(nil)
	return nil
}

// ============================================================
// Helper Functions
// ============================================================

// getJVMThreadFromObject get JVMThread from java.lang.Thread
func getJVMThreadFromObject(threadObj *heap.Object) *jvmruntime.JVMThread {
	if threadObj == nil {
		return nil
	}

	extra := threadObj.Extra()
	if extra == nil {
		return nil
	}

	jvmThread, ok := extra.(*jvmruntime.JVMThread)
	if !ok {
		return nil
	}

	return jvmThread
}
