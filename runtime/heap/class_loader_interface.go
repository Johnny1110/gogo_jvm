package heap

// ============================================================
// v0.3.1: Interface to break circular dependency
// ============================================================
type ClassLoaderProvider interface {
	LoadClassIface(name string) interface{}
}
