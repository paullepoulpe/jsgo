package jutil

import "strconv"

// Go-only: language-adapter — java.nio.charset.CharacterCodingException and
// its subclass MalformedInputException (both IOExceptions).

// CharacterCodingException is java.nio.charset.CharacterCodingException.
type CharacterCodingException struct{ IOException }

func (*CharacterCodingException) characterCodingException() {}

// ICharacterCodingException matches CharacterCodingException and its
// subclasses.
type ICharacterCodingException interface {
	IIOException
	characterCodingException()
}

// MalformedInputException is java.nio.charset.MalformedInputException.
type MalformedInputException struct {
	CharacterCodingException
	inputLength int
}

// NewMalformedInputException is new MalformedInputException(inputLength);
// its message is "Input length = <inputLength>".
func NewMalformedInputException(inputLength int) *MalformedInputException {
	e := &MalformedInputException{inputLength: inputLength}
	msg := "Input length = " + strconv.Itoa(inputLength)
	InitThrowable(e, "java.nio.charset.MalformedInputException", &msg, nil)
	return e
}

// GetInputLength is MalformedInputException.getInputLength().
func (e *MalformedInputException) GetInputLength() int { return e.inputLength }
