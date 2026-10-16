package jutil

// ConcurrentModificationException is
// java.util.ConcurrentModificationException: thrown by fail-fast iterators
// (JavaHashMap, JavaHashSet) when the collection is structurally modified
// other than through the iterator itself.
type ConcurrentModificationException struct{ RuntimeException }

func (*ConcurrentModificationException) concurrentModificationException() {}

type IConcurrentModificationException interface {
	IRuntimeException
	concurrentModificationException()
}

// NewConcurrentModificationException is new
// ConcurrentModificationException().
func NewConcurrentModificationException() *ConcurrentModificationException {
	e := &ConcurrentModificationException{}
	InitThrowable(e, "java.util.ConcurrentModificationException", nil, nil)
	return e
}
