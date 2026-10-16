// Port of com/google/javascript/jscomp/FlagUsageException.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// FlagUsageException is an exception thrown when command-line flags are used
// incorrectly.
type FlagUsageException struct{ jutil.RuntimeException }

func (*FlagUsageException) flagUsageException() {}

// IFlagUsageException matches FlagUsageException.
type IFlagUsageException interface {
	jutil.IRuntimeException
	flagUsageException()
}

// Java: FlagUsageException.serialVersionUID
const FlagUsageException_serialVersionUID int64 = 1

// Java: FlagUsageException(String)
func NewFlagUsageException(message string) *FlagUsageException {
	e := &FlagUsageException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.FlagUsageException", &message, nil)
	return e
}
