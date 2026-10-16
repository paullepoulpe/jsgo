package jutil

// Go-only: language-adapter — java.text.ParseException.

// ParseException is java.text.ParseException.
type ParseException struct {
	JavaException
	errorOffset int
}

func (*ParseException) parseException() {}

// IParseException matches ParseException.
type IParseException interface {
	IException
	parseException()
}

// NewParseException is new ParseException(String s, int errorOffset).
func NewParseException(s string, errorOffset int) *ParseException {
	e := &ParseException{errorOffset: errorOffset}
	InitThrowable(e, "java.text.ParseException", msgPtr(s), nil)
	return e
}

// GetErrorOffset is ParseException.getErrorOffset().
func (e *ParseException) GetErrorOffset() int { return e.errorOffset }
