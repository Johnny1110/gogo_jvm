package heap

import (
	"sync"
	"time"
)

// ============================================================
// Identity Hash Code Generator
// ============================================================
// Using Marsaglia XorShift algo to generate identity hash code
// Same as HotSpot JVM Default (hashCode=5)

var (
	// TODO: 目前使用全域鎖，v0.4.x 多執行緒階段可改為 Thread-Local
	hashMutex sync.Mutex
	hashState uint32
)

func init() {
	hashState = uint32(time.Now().UnixNano())
	// make sure seed is not 0（XorShift）
	if hashState == 0 {
		hashState = 1
	}
}

// generateHashCode generate identity hash code
// Using XorShift32 Algo
// return 31 bits int (0x00000001 ~ 0x7FFFFFFF)
func generateHashCode() int32 {
	hashMutex.Lock()
	defer hashMutex.Unlock()

	// XorShift32 Algo
	x := hashState
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	hashState = x

	// make sure 31 bits and not 0
	// hashCode = 0 not calculated, so actual val must > 0.
	result := int32(x & 0x7FFFFFFF)
	if result == 0 {
		result = 1
	}
	return result
}
