// Port of com/google/javascript/jscomp/base/format/MissingFormatWidthException.java @ closure-compiler v20260415

package base

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// MissingFormatWidthException is com.google.javascript.jscomp.base.format.MissingFormatWidthException.
type MissingFormatWidthException struct{ jutil.RuntimeException }

func (*MissingFormatWidthException) isMissingFormatWidthException() {}

// IMissingFormatWidthException matches MissingFormatWidthException and its subclasses in jutil.Try.
type IMissingFormatWidthException interface {
	jutil.IRuntimeException
	isMissingFormatWidthException()
}

// NewMissingFormatWidthException constructs the exception.
//
// Java: MissingFormatWidthException(String)
func NewMissingFormatWidthException(msg string) *MissingFormatWidthException {
	e := &MissingFormatWidthException{}
	m := msg
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.MissingFormatWidthException", &m, nil)
	return e
}
