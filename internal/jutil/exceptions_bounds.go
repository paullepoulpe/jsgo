package jutil

import "fmt"

// Go-only: language-adapter — the exact index-out-of-bounds exceptions the
// JDK 21 runtime throws from String, StringBuilder and the collection
// classes (their messages reach Closure's diagnostics, e.g. MALFORMED_REGEXP).

// NewStringIndexOutOfBoundsExceptionForIndex is the exception thrown by
// String.charAt(int) and StringBuilder.charAt(int) / setCharAt(int, char):
// "Index 3 out of bounds for length 3".
func NewStringIndexOutOfBoundsExceptionForIndex(index, length int) *StringIndexOutOfBoundsException {
	e := &StringIndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.StringIndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("Index %d out of bounds for length %d", index, length)), nil)
	return e
}

// NewStringIndexOutOfBoundsExceptionForRange is the exception thrown by
// String.substring(int, int) and String.substring(int):
// "Range [2, 1) out of bounds for length 3".
func NewStringIndexOutOfBoundsExceptionForRange(begin, end, length int) *StringIndexOutOfBoundsException {
	e := &StringIndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.StringIndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("Range [%d, %d) out of bounds for length %d", begin, end, length)), nil)
	return e
}

// NewIndexOutOfBoundsExceptionNoMessage is new IndexOutOfBoundsException().
func NewIndexOutOfBoundsExceptionNoMessage() *IndexOutOfBoundsException {
	e := &IndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.IndexOutOfBoundsException", nil, nil)
	return e
}

// NewIndexOutOfBoundsExceptionForRange is the exception thrown by
// Preconditions.checkFromToIndex with the IOOBE formatter, e.g. by
// StringBuilder.append(CharSequence, int, int):
// "Range [2, 5) out of bounds for length 3".
func NewIndexOutOfBoundsExceptionForRange(from, to, length int) *IndexOutOfBoundsException {
	e := &IndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.IndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("Range [%d, %d) out of bounds for length %d", from, to, length)), nil)
	return e
}
