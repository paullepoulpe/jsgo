package jutil

// MissingResourceException is java.util.MissingResourceException (thrown by
// ResourceBundle.getString for an unknown key).
type MissingResourceException struct {
	RuntimeException
	className string
	key       string
}

func (*MissingResourceException) missingResourceException() {}

// IMissingResourceException is the catch marker for MissingResourceException.
type IMissingResourceException interface {
	IRuntimeException
	missingResourceException()
}

// NewMissingResourceException is new MissingResourceException(String s,
// String className, String key).
func NewMissingResourceException(s, className, key string) *MissingResourceException {
	e := &MissingResourceException{className: className, key: key}
	InitThrowable(e, "java.util.MissingResourceException", &s, nil)
	return e
}

// GetClassName is MissingResourceException.getClassName().
func (e *MissingResourceException) GetClassName() string { return e.className }

// GetKey is MissingResourceException.getKey().
func (e *MissingResourceException) GetKey() string { return e.key }
