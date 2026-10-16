// Port of com/google/javascript/jscomp/base/format/IllegalFormatFlagsException.java @ closure-compiler v20260415

package base

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IllegalFormatFlagsException is com.google.javascript.jscomp.base.format.IllegalFormatFlagsException.
type IllegalFormatFlagsException struct{ jutil.RuntimeException }

func (*IllegalFormatFlagsException) isIllegalFormatFlagsException() {}

// IIllegalFormatFlagsException matches IllegalFormatFlagsException and its subclasses in jutil.Try.
type IIllegalFormatFlagsException interface {
	jutil.IRuntimeException
	isIllegalFormatFlagsException()
}

// NewIllegalFormatFlagsException constructs the exception.
//
// Java: IllegalFormatFlagsException(String)
func NewIllegalFormatFlagsException(message string) *IllegalFormatFlagsException {
	e := &IllegalFormatFlagsException{}
	m := message
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.IllegalFormatFlagsException", &m, nil)
	return e
}
