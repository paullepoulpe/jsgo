package jutil

// IllegalMonitorStateException is java.lang.IllegalMonitorStateException:
// thrown when a thread releases a monitor (or a ReentrantLock) it does not
// hold.
type IllegalMonitorStateException struct{ RuntimeException }

func (*IllegalMonitorStateException) illegalMonitorStateException() {}

// IIllegalMonitorStateException matches IllegalMonitorStateException.
type IIllegalMonitorStateException interface {
	IRuntimeException
	illegalMonitorStateException()
}

// NewIllegalMonitorStateException is new IllegalMonitorStateException().
func NewIllegalMonitorStateException() *IllegalMonitorStateException {
	e := &IllegalMonitorStateException{}
	InitThrowable(e, "java.lang.IllegalMonitorStateException", nil, nil)
	return e
}
