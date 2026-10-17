// Port of org/kohsuke/args4j/IllegalAnnotationError.java @ args4j 2.33

/*
 * Copyright (c) 2013 Kohsuke Kawaguchi and other contributors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 */

package args4j

import "github.com/paullepoulpe/jsgo/internal/jutil"

const illegalAnnotationError_className = "com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.IllegalAnnotationError"

// IllegalAnnotationError signals an incorrect use of args4j annotations.
//
// This only happens when there's something wrong with the way you use args4j
// in your code--not when the user's arguments are wrong. Therefore, this
// class is an Error.
//
// Java: class IllegalAnnotationError extends Error
type IllegalAnnotationError struct {
	jutil.JavaLangError
}

func (*IllegalAnnotationError) illegalAnnotationError() {}

// IIllegalAnnotationError matches IllegalAnnotationError.
type IIllegalAnnotationError interface {
	jutil.IError
	illegalAnnotationError()
}

// Java: IllegalAnnotationError(String)
func NewIllegalAnnotationError(message string) *IllegalAnnotationError {
	e := &IllegalAnnotationError{}
	jutil.InitThrowable(e, illegalAnnotationError_className, &message, nil)
	return e
}

// Java: IllegalAnnotationError(String, Throwable)
func NewIllegalAnnotationErrorWithMessageAndCause(message string, cause jutil.Throwable) *IllegalAnnotationError {
	e := &IllegalAnnotationError{}
	jutil.InitThrowable(e, illegalAnnotationError_className, &message, cause)
	return e
}

// Java: IllegalAnnotationError(Throwable)
func NewIllegalAnnotationErrorWithCause(cause jutil.Throwable) *IllegalAnnotationError {
	e := &IllegalAnnotationError{}
	jutil.InitThrowable(e, illegalAnnotationError_className, nil, cause)
	return e
}
