package jutil

// InterruptedException is java.lang.InterruptedException (a checked
// exception).
type InterruptedException struct{ JavaException }

func (*InterruptedException) interruptedException() {}

// IInterruptedException matches InterruptedException.
type IInterruptedException interface {
	IException
	interruptedException()
}

// NewInterruptedException is new InterruptedException().
func NewInterruptedException() *InterruptedException {
	e := &InterruptedException{}
	InitThrowable(e, "java.lang.InterruptedException", nil, nil)
	return e
}
