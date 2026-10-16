// Port of com/google/javascript/jscomp/base/format/IllegalFormatConversionException.java @ closure-compiler v20260415

package base

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IllegalFormatConversionException is com.google.javascript.jscomp.base.format.IllegalFormatConversionException.
type IllegalFormatConversionException struct{ jutil.RuntimeException }

func (*IllegalFormatConversionException) isIllegalFormatConversionException() {}

// IIllegalFormatConversionException matches IllegalFormatConversionException and its subclasses in jutil.Try.
type IIllegalFormatConversionException interface {
	jutil.IRuntimeException
	isIllegalFormatConversionException()
}

// NewIllegalFormatConversionException constructs the exception.
//
// Java: IllegalFormatConversionException(char, Class<? extends Object>)
//
// Go-only: language-adapter — the Class argument is passed as its
// Class.toString() form (e.g. "class java.lang.String").
func NewIllegalFormatConversionException(conversionType uint16, class1 string) *IllegalFormatConversionException {
	e := &IllegalFormatConversionException{}
	m := "ConversionType = " + jutil.StringFromChar(conversionType) + " class = " + class1
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.IllegalFormatConversionException", &m, nil)
	return e
}
