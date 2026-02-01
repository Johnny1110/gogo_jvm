package runtime

import (
	"bytes"
	"fmt"
	"github.com/Johnny1110/gogo_jvm/global"
	"github.com/Johnny1110/gogo_jvm/runtime/heap"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

// ============================================================
// Thread Manager - v0.4.0
// ============================================================
// Responsibility
// 1. manage all JVMThread instances
// 2. provide Thread Local Storage (TLS) feature
// 3. trace daemon / non-daemon thread count
// 4. provide currentThread() feature

// ThreadManager Global ThreadManager
type ThreadManager struct {
	mutex sync.RWMutex

	// all JVMThread mapper (threadID(JVM) -> JVMThread)
	threads map[int64]*JVMThread

	// goroutine ID -> JVMThread mapper（for currentThread()）
	goroutineMap map[int64]*JVMThread

	// id_seq
	nextThreadID int64

	// when it comes to 0, jvm should exit
	nonDaemonCount int32

	// main thread
	mainThread *JVMThread

	// JVM exist signal
	exitCh chan struct{}
}

// globalThreadManager: global singleton ThreadManager
var globalThreadManager = &ThreadManager{
	threads:      make(map[int64]*JVMThread),
	goroutineMap: make(map[int64]*JVMThread),
	nextThreadID: 1, // start from 1 (1 is main thread)
	exitCh:       make(chan struct{}),
}

// ============================================================
// Public API
// ============================================================

// GetThreadManager getter
func GetThreadManager() *ThreadManager {
	return globalThreadManager
}

// CurrentJVMThread get current goroutine mapping JVMThread
// this is for supporting Thread.currentThread() feature
func CurrentJVMThread() *JVMThread {
	return globalThreadManager.CurrentThread()
}

// ============================================================
// Thread Registration
// ============================================================

func (tm *ThreadManager) Register(t *JVMThread) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.threads[t.id] = t

	// update count
	if !t.daemon {
		atomic.AddInt32(&tm.nonDaemonCount, 1)
	}

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - ThreadManager: Registered thread [%d] %s (daemon=%v, total=%d)\n",
			t.id, t.name, t.daemon, len(tm.threads))
	}
}

// Unregister remove thread
func (tm *ThreadManager) Unregister(t *JVMThread) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	delete(tm.threads, t.id)

	// update count
	if !t.daemon {
		count := atomic.AddInt32(&tm.nonDaemonCount, -1)
		if count == 0 {
			// all non-daemon == 0, means main thread is gone also.
			if global.DebugMode() {
				fmt.Println("@@ DEBUG - ThreadManager: All non-daemon threads terminated, signaling exit")
			}
			close(tm.exitCh)
		}
	}

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - ThreadManager: Unregistered thread [%d] %s (remaining=%d)\n",
			t.id, t.name, len(tm.threads))
	}
}

// ============================================================
// Goroutine Mapping (Thread Local Storage)
// ============================================================

// RegisterGoroutineMapping register goroutine ID to JVMThread mapping
// this method should call before start of each thread's goroutine
func (tm *ThreadManager) RegisterGoroutineMapping(t *JVMThread) {
	gid := getGoroutineID()

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.goroutineMap[gid] = t

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - ThreadManager: Mapped goroutine %d -> thread [%d] %s\n",
			gid, t.id, t.name)
	}
}

// UnregisterGoroutineMapping remove current goroutine mapping
func (tm *ThreadManager) UnregisterGoroutineMapping() {
	gid := getGoroutineID()

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	delete(tm.goroutineMap, gid)
}

// CurrentThread get current goroutine mapped JVMThread
func (tm *ThreadManager) CurrentThread() *JVMThread {
	gid := getGoroutineID()

	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	if t, ok := tm.goroutineMap[gid]; ok {
		return t
	} else {
		return tm.mainThread
	}
}

// ============================================================
// Thread ID Allocation
// ============================================================

// AllocateThreadID allocate thread ID
func (tm *ThreadManager) AllocateThreadID() int64 {
	return atomic.AddInt64(&tm.nextThreadID, 1)
}

// GetThread get thread by jvm thread ID
func (tm *ThreadManager) GetThread(id int64) *JVMThread {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	return tm.threads[id]
}

// AllThreads get all threads
func (tm *ThreadManager) AllThreads() []*JVMThread {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	result := make([]*JVMThread, 0, len(tm.threads))

	for _, t := range tm.threads {
		result = append(result, t)
	}
	return result
}

// ThreadCount get thread count
func (tm *ThreadManager) ThreadCount() int {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	return len(tm.threads)
}

// NonDaemonCount get non-daemon count
func (tm *ThreadManager) NonDaemonCount() int32 {
	return atomic.LoadInt32(&tm.nonDaemonCount)
}

// ============================================================
// Main Thread Management
// ============================================================

// SetMainThread setup main thread
func (tm *ThreadManager) SetMainThread(t *JVMThread) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.mainThread = t
	tm.threads[1] = t // main thread ID should always be 1
	// add non daemon thread count
	atomic.AddInt32(&tm.nonDaemonCount, 1)
}

func (tm *ThreadManager) GetMainThread() *JVMThread {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	return tm.mainThread
}

// ============================================================
// JVM Exit
// ============================================================

// ExitChannel get JVM exit signal channel
func (tm *ThreadManager) ExitChannel() <-chan struct{} {
	return tm.exitCh
}

// WaitForExit waiting for all non-daemon thread done
func (tm *ThreadManager) WaitForExit() {
	<-tm.exitCh // block until get signal
}

// ============================================================
// Thread Adapter (Compatible with legacy Thread structure)
// ============================================================

var threadAdapters sync.Map // JVMThread -> *Thread

// getThreadAdapter
func getThreadAdapter(jvmThread *JVMThread) *Thread {
	if adapter, ok := threadAdapters.Load(jvmThread); ok {
		return adapter.(*Thread)
	}

	// 創建新的適配器
	adapter := &Thread{
		pc:    0,
		stack: jvmThread.jvmStack,
	}

	threadAdapters.Store(jvmThread, adapter)
	return adapter
}

// GetThreadAdapter
func (tm *ThreadManager) GetThreadAdapter(jvmThread *JVMThread) *Thread {
	return getThreadAdapter(jvmThread)
}

// syncAdapterPC
func syncAdapterPC(jvmThread *JVMThread) {
	if adapter, ok := threadAdapters.Load(jvmThread); ok {
		adapter.(*Thread).pc = jvmThread.pc
	}
}

// ============================================================
// Goroutine ID Helper
// ============================================================

// getGoroutineID get current goroutine ID
func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)

	// Stack output format: "goroutine 123 [running]:\n..."
	// we need parse "123"
	field := bytes.Fields(buf[:n])
	if len(field) < 2 {
		return 0
	}

	idStr := string(field[1])
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0
	}

	return id
}

// ============================================================
// Helper Functions for Creating Threads
// ============================================================

// CreateThread create new JVMThread
// this is for native method Thread.start0()
func CreateThread(name string, daemon bool, priority int) *JVMThread {
	id := GetThreadManager().AllocateThreadID()

	t := NewJVMThread(id, name)
	t.daemon = daemon
	t.priority = priority

	if global.DebugMode() {
		fmt.Printf("@@ DEBUG - CreateThread: %v", t)
	}

	return t
}

// CreateThreadFromJavaObject create JVMThread from java.lang.Thread
func CreateThreadFromJavaObject(threadObj *heap.Object) *JVMThread {
	id := GetThreadManager().AllocateThreadID()

	t := NewJVMThread(id, fmt.Sprintf("Thread-%d", id))
	t.javaThreadObj = threadObj

	return t
}

// ============================================================
// Initialize Main Thread
// ============================================================

// InitMainThread init main
// invoke this when jvm start
func InitMainThread() *JVMThread {
	mainThread := NewMainThread()

	tm := GetThreadManager()
	tm.SetMainThread(mainThread)

	tm.RegisterGoroutineMapping(mainThread)

	return mainThread
}

// ============================================================
// Compatibility Layer
// ============================================================

func NewThread() *Thread {
	tm := GetThreadManager()
	if tm.mainThread == nil {
		// 初始化主執行緒
		mainThread := InitMainThread()
		return getThreadAdapter(mainThread)
	}

	current := tm.CurrentThread()
	if current != nil {
		return getThreadAdapter(current)
	}

	return &Thread{
		pc:    0,
		stack: NewJVMStack(DEFAULT_STACK_SIZE),
	}
}
