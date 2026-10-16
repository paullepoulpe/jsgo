// Port of com/google/javascript/jscomp/base/format/IllegalFormatCodePointException.java @ closure-compiler v20260415

package base

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IllegalFormatCodePointException is com.google.javascript.jscomp.base.format.IllegalFormatCodePointException.
type IllegalFormatCodePointException struct{ jutil.RuntimeException }

func (*IllegalFormatCodePointException) isIllegalFormatCodePointException() {}

// IIllegalFormatCodePointException matches IllegalFormatCodePointException and its subclasses in jutil.Try.
type IIllegalFormatCodePointException interface {
	jutil.IRuntimeException
	isIllegalFormatCodePointException()
}

// NewIllegalFormatCodePointException constructs the exception.
//
// Java: IllegalFormatCodePointException(int)
func NewIllegalFormatCodePointException(codePoint int) *IllegalFormatCodePointException {
	e := &IllegalFormatCodePointException{}
	m := "codePoint:" + strconv.Itoa(codePoint)
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.IllegalFormatCodePointException", &m, nil)
	return e
}
