package runtime

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================
// JVMThread - v0.4.0 Multi-Thread Support
// ============================================================
// JVMThread is represent GOGO JVM inner thread
// JVMThread & java.lang.Thread they are ref to each other (Mirror)
//
//   Java Layer: java.lang.Thread (heap.Object)
//                        ↑↓
//    JVM Layer: JVMThread (Go struct)
//                        ↓
//   Go Runtime: goroutine

const DEFAULT_STACK_SIZE = 1024
const DEFAULT_PRIORITY = 5

// JVMThread JVM inner thread structure
type JVMThread struct {
	// basic-info
	id   int64
	name string

	// runtime state
	state ThreadState
	pc    int // program counter

	// jvm frame stack (method call stack)
	jvmStack *JVMStack

	// mirror object
	javaThreadObj *heap.Object

	// thread attrs
	daemon   bool
	priority int // 1~10

	// sync (from go perspective)
	// using by join()
	done chan struct{}

	// protect thread state transform
	stateMutex sync.RWMutex

	// interrupted TODO v0.4.0
	interrupted int32 // using atomic，0=false, 1=true
}

// ============================================================
// Constructor
// ============================================================

// NewJVMThread create new JVMThread
// init state is NEW
func NewJVMThread(id int64, name string) *JVMThread {
	return &JVMThread{
		id:       id,
		name:     name,
		state:    THREAD_NEW,
		pc:       0,
		jvmStack: NewJVMStack(DEFAULT_STACK_SIZE),
		daemon:   false,
		priority: DEFAULT_PRIORITY, // NORM_PRIORITY
		done:     make(chan struct{}),
	}
}

// NewMainThread create main thread
// Main thread is special, will be created when JVM starting
func NewMainThread() *JVMThread {
	mainT := NewJVMThread(1, "main")
	mainT.state = THREAD_RUNNABLE // main 執行緒直接是 RUNNABLE
	return mainT
}

// ============================================================
// Getters
// ============================================================

func (t *JVMThread) ID() int64 {
	return t.id
}

func (t *JVMThread) Name() string {
	return t.name
}

func (t *JVMThread) State() ThreadState {
	t.stateMutex.RLock()
	defer t.stateMutex.RUnlock()
	return t.state
}

func (t *JVMThread) PC() int {
	return t.pc
}

func (t *JVMThread) IsDaemon() bool {
	return t.daemon
}

func (t *JVMThread) Priority() int {
	return t.priority
}

func (t *JVMThread) JavaThreadObj() *heap.Object {
	return t.javaThreadObj
}

func (t *JVMThread) Done() <-chan struct{} {
	return t.done
}

// ============================================================
// Setters
// ============================================================

func (t *JVMThread) SetName(name string) {
	t.name = name
}

func (t *JVMThread) SetPC(pc int) {
	t.pc = pc
}

// SetDaemon can only call before start()
func (t *JVMThread) SetDaemon(daemon bool) {
	if t.state != THREAD_NEW {
		panic("java.lang.IllegalThreadStateException: cannot set daemon after start")
	}

	t.daemon = daemon
}

// SetPriority 1~10
func (t *JVMThread) SetPriority(priority int) {
	if priority < 1 || priority > 10 {
		runtimeHandleException(t.CurrentFrame(), "java/lang/IllegalArgumentException", "priority out of range 1~10")
		return
	}

	t.priority = priority
}

func (t *JVMThread) SetJavaThreadObj(obj *heap.Object) {
	t.javaThreadObj = obj
}

func (t *JVMThread) SetState(newState ThreadState) {
	t.stateMutex.Lock()
	defer t.stateMutex.Unlock()

	if global.DebugMode() {
		oldState := t.state
		if !oldState.CanTransitionTo(newState) {
			fmt.Printf("@@ WARNING - Invalid thread state transition: %s -> %s\n",
				oldState.String(), newState.String())
		}

		fmt.Printf("@@ DEBUG - Thread [%s] state: %s -> %s\n",
			t.name, oldState.String(), newState.String())
	}

	t.state = newState
}

// ============================================================
// JVM Stack Operations
// ============================================================

func (t *JVMThread) PushFrame(frame *Frame) {
	t.jvmStack.Push(frame)
}

func (t *JVMThread) PopFrame() *Frame {
	return t.jvmStack.Pop()
}

// CurrentFrame get current frame without pop
func (t *JVMThread) CurrentFrame() *Frame {
	return t.jvmStack.Top()
}

// TopFrame is CurrentFrame alias
func (t *JVMThread) TopFrame() *Frame {
	return t.CurrentFrame()
}

func (t *JVMThread) IsStackEmpty() bool {
	return t.jvmStack.IsEmpty()
}

func (t *JVMThread) StackDepth() uint {
	return t.jvmStack.Size()
}

func (t *JVMThread) GetFrames() []*Frame {
	return t.jvmStack.GetFrames()
}

func (t *JVMThread) ClearStack() {
	t.jvmStack.Clear()
}

// ============================================================
// Frame Creation
// ============================================================

func (t *JVMThread) NewFrameWithMethodAndExHandler(method *method_area.Method, exHandler func(frame *Frame, ex *heap.Object)) *Frame {
	return NewFrameWithMethodAndExHandler(t.asThread(), method, exHandler)
}

func (t *JVMThread) asThread() *Thread {
	return getThreadAdapter(t)
}

// ============================================================
// Thread Lifecycle
// ============================================================

func (t *JVMThread) IsAlive() bool {
	return t.State().IsAlive()
}

// Start java thread
// this method will be call by native Thread.start0()
func (t *JVMThread) Start(runMethod *method_area.Method) {
	if t.state != THREAD_NEW {
		runtimeHandleException(t.CurrentFrame(), "java/lang/IllegalThreadStateException", "thread already started")
		return
	}

	t.SetState(THREAD_RUNNABLE)

	// register to globalThreadManager
	globalThreadManager.Register(t)

	go t.runInternal(runMethod)
}

// runInternal thread inner logic
func (t *JVMThread) runInternal(runMethod *method_area.Method) {
	// register goroutine ID
	globalThreadManager.RegisterGoroutineMapping(t)
	defer globalThreadManager.UnregisterGoroutineMapping()

	defer func() {

		// this is for handle uncached error
		if r := recover(); r != nil {
			if global.DebugMode() {
				fmt.Printf("@@ Thread [%s] uncaught exception: %v\n", t.name, r)
			}

			t.SetState(THREAD_TERMINATED) // end of thread lifecycle

			// close done channel, wake up all join() waiter
			close(t.done)

			globalThreadManager.Unregister(t)
		}
	}()

	// execute Thread.run() method
	t.executeRun(runMethod)
}

// executeRun execute Thread run() method
func (t *JVMThread) executeRun(runMethod *method_area.Method) {
	if threadExecutor != nil {
		threadExecutor(t, runMethod)
	} else {
		panic("Thread executor not set! Call SetThreadExecutor() first.")
	}
}

// Join wait thread running result
// millis: waiting time(ms), 0 -> infinity
func (t *JVMThread) Join(millis int64) {
	if millis == 0 {
		// waiting
		<-t.done // this will block the thread until somewhere close(t.done)
	} else {
		// timed waiting
		select {
		case <-t.done: // this will block the thread until somewhere close(t.done)
			// thread already done
		case <-timeAfter(millis): // this will block the thread until time after millis
			// time out
		}
	}
}

// ============================================================
// Time Helper
// ============================================================

func timeAfter(millis int64) <-chan time.Time {
	return time.After(time.Duration(millis) * time.Millisecond)
}

// ============================================================
// Interrupt (TODO: v0.4.3)
// ============================================================

// Interrupt
func (t *JVMThread) Interrupt() {
	atomic.StoreInt32(&t.interrupted, 1)
	// TODO: v0.4.3
}

// IsInterrupted check is interrupted
func (t *JVMThread) IsInterrupted() bool {
	return atomic.LoadInt32(&t.interrupted) == 1
}

// ClearInterrupt check and clean interrupt
func (t *JVMThread) ClearInterrupt() bool {
	return atomic.SwapInt32(&t.interrupted, 0) == 1
}

// ============================================================
// Thread Executor Callback
// ============================================================
// avoid cycle dependency（runtime -> interpreter -> runtime）

type ThreadExecutor func(thread *JVMThread, method *method_area.Method)

var threadExecutor ThreadExecutor

// SetThreadExecutor setup thread executor
// should be call before interpreter
func SetThreadExecutor(executor ThreadExecutor) {
	threadExecutor = executor
}
