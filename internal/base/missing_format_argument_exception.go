// Port of com/google/javascript/jscomp/base/format/MissingFormatArgumentException.java @ closure-compiler v20260415

package base

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// MissingFormatArgumentException is com.google.javascript.jscomp.base.format.MissingFormatArgumentException.
type MissingFormatArgumentException struct{ jutil.RuntimeException }

func (*MissingFormatArgumentException) isMissingFormatArgumentException() {}

// IMissingFormatArgumentException matches MissingFormatArgumentException and its subclasses in jutil.Try.
type IMissingFormatArgumentException interface {
	jutil.IRuntimeException
	isMissingFormatArgumentException()
}

// NewMissingFormatArgumentException constructs the exception.
//
// Java: MissingFormatArgumentException(String)
func NewMissingFormatArgumentException(message string) *MissingFormatArgumentException {
	e := &MissingFormatArgumentException{}
	m := message
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.MissingFormatArgumentException", &m, nil)
	return e
}
