package jutil

import "strconv"

// Go-only: language-adapter — java.net.URISyntaxException (JDK 21).

// URISyntaxException is java.net.URISyntaxException, a checked exception.
type URISyntaxException struct {
	JavaException
	input  string
	reason string
	index  int
}

func (*URISyntaxException) uriSyntaxException() {}

// IURISyntaxException matches URISyntaxException.
type IURISyntaxException interface {
	IException
	uriSyntaxException()
}

// NewURISyntaxException is new URISyntaxException(input, reason, index);
// index -1 means unknown.
func NewURISyntaxException(input, reason string, index int) *URISyntaxException {
	e := &URISyntaxException{input: input, reason: reason, index: index}
	msg := reason
	if index > -1 {
		msg += " at index " + strconv.Itoa(index)
	}
	msg += ": " + input
	InitThrowable(e, "java.net.URISyntaxException", &msg, nil)
	return e
}

// GetInput is URISyntaxException.getInput().
func (e *URISyntaxException) GetInput() string { return e.input }

// GetReason is URISyntaxException.getReason().
func (e *URISyntaxException) GetReason() string { return e.reason }

// GetIndex is URISyntaxException.getIndex().
func (e *URISyntaxException) GetIndex() int { return e.index }
