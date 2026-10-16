// Port of com/google/javascript/jscomp/base/format/UnknownFormatConversionException.java @ closure-compiler v20260415

package base

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// UnknownFormatConversionException is com.google.javascript.jscomp.base.format.UnknownFormatConversionException.
type UnknownFormatConversionException struct{ jutil.RuntimeException }

func (*UnknownFormatConversionException) isUnknownFormatConversionException() {}

// IUnknownFormatConversionException matches UnknownFormatConversionException and its subclasses in jutil.Try.
type IUnknownFormatConversionException interface {
	jutil.IRuntimeException
	isUnknownFormatConversionException()
}

// NewUnknownFormatConversionException constructs the exception.
//
// Java: UnknownFormatConversionException(String)
func NewUnknownFormatConversionException(message string) *UnknownFormatConversionException {
	e := &UnknownFormatConversionException{}
	m := message
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.UnknownFormatConversionException", &m, nil)
	return e
}
