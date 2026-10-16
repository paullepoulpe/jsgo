package jutil

// IllegalCharsetNameException is
// java.nio.charset.IllegalCharsetNameException: thrown when a string that
// is not a legal charset name is used as such. Its message is the name.
type IllegalCharsetNameException struct {
	IllegalArgumentException
	charsetName string
}

func (*IllegalCharsetNameException) illegalCharsetNameException() {}

// IIllegalCharsetNameException is the catch marker for
// IllegalCharsetNameException.
type IIllegalCharsetNameException interface {
	IIllegalArgumentException
	illegalCharsetNameException()
}

// NewIllegalCharsetNameException is new IllegalCharsetNameException(String).
func NewIllegalCharsetNameException(charsetName string) *IllegalCharsetNameException {
	e := &IllegalCharsetNameException{charsetName: charsetName}
	InitThrowable(e, "java.nio.charset.IllegalCharsetNameException", &charsetName, nil)
	return e
}

// GetCharsetName is IllegalCharsetNameException.getCharsetName().
func (e *IllegalCharsetNameException) GetCharsetName() string { return e.charsetName }

// requireLegalCharsetName applies the Javadoc's rules for legal charset names:
// non-empty, made of the letters A-Z and a-z, the digits 0-9, '-', '+',
// '.', ':' and '_', and beginning with a letter or a digit.
func requireLegalCharsetName(name string) {
	if name == "" {
		panic(NewIllegalCharsetNameException(name))
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		letterOrDigit := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		if letterOrDigit {
			continue
		}
		if i > 0 && (c == '-' || c == '+' || c == '.' || c == ':' || c == '_') {
			continue
		}
		panic(NewIllegalCharsetNameException(name))
	}
}

// Charset_IsSupported is Charset.isSupported(String): whether the named
// charset is supported, throwing IllegalCharsetNameException for an illegal
// name.
//
// Go-only: language-adapter — only the standard charsets (charset.go) are
// supported; a JDK also supports many others (windows-1252, Shift_JIS, ...),
// for which this returns false.
func Charset_IsSupported(charsetName string) bool {
	requireLegalCharsetName(charsetName)
	lower := StringToLowerCaseASCII(charsetName)
	for _, c := range standardCharsets {
		if lower == StringToLowerCaseASCII(c.name) {
			return true
		}
		for _, a := range c.aliases {
			if lower == StringToLowerCaseASCII(a) {
				return true
			}
		}
	}
	return false
}
