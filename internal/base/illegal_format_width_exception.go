// Port of com/google/javascript/jscomp/base/format/IllegalFormatWidthException.java @ closure-compiler v20260415

package base

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IllegalFormatWidthException is com.google.javascript.jscomp.base.format.IllegalFormatWidthException.
type IllegalFormatWidthException struct{ jutil.RuntimeException }

func (*IllegalFormatWidthException) isIllegalFormatWidthException() {}

// IIllegalFormatWidthException matches IllegalFormatWidthException and its subclasses in jutil.Try.
type IIllegalFormatWidthException interface {
	jutil.IRuntimeException
	isIllegalFormatWidthException()
}

// NewIllegalFormatWidthException constructs the exception.
//
// Java: IllegalFormatWidthException(int)
func NewIllegalFormatWidthException(width int) *IllegalFormatWidthException {
	e := &IllegalFormatWidthException{}
	m := "Width:" + strconv.Itoa(width)
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.IllegalFormatWidthException", &m, nil)
	return e
}
