// Port of com/google/javascript/jscomp/base/format/IllegalFormatPrecisionException.java @ closure-compiler v20260415

package base

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IllegalFormatPrecisionException is com.google.javascript.jscomp.base.format.IllegalFormatPrecisionException.
type IllegalFormatPrecisionException struct{ jutil.RuntimeException }

func (*IllegalFormatPrecisionException) isIllegalFormatPrecisionException() {}

// IIllegalFormatPrecisionException matches IllegalFormatPrecisionException and its subclasses in jutil.Try.
type IIllegalFormatPrecisionException interface {
	jutil.IRuntimeException
	isIllegalFormatPrecisionException()
}

// NewIllegalFormatPrecisionException constructs the exception.
//
// Java: IllegalFormatPrecisionException(int)
func NewIllegalFormatPrecisionException(precision int) *IllegalFormatPrecisionException {
	e := &IllegalFormatPrecisionException{}
	m := "Precision: " + strconv.Itoa(precision)
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.IllegalFormatPrecisionException", &m, nil)
	return e
}
