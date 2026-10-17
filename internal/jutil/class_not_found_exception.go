package jutil

// ClassNotFoundException is java.lang.ClassNotFoundException (a checked
// exception, extending ReflectiveOperationException).
type ClassNotFoundException struct{ JavaException }

func (*ClassNotFoundException) classNotFoundException() {}

// IClassNotFoundException is the catch marker for ClassNotFoundException.
type IClassNotFoundException interface {
	IException
	classNotFoundException()
}

// NewClassNotFoundException is new ClassNotFoundException(String).
func NewClassNotFoundException(message string) *ClassNotFoundException {
	e := &ClassNotFoundException{}
	InitThrowable(e, "java.lang.ClassNotFoundException", &message, nil)
	return e
}
