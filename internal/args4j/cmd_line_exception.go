// Port of org/kohsuke/args4j/CmdLineException.java @ args4j 2.33

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

// cmdLineException_className is the class name Throwable.toString() prints:
// the pinned closure-compiler.jar shades args4j into
// com.google.javascript.jscomp.jarjar.
const cmdLineException_className = "com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.CmdLineException"

// CmdLineException signals an error in the user input.
//
// Java: class CmdLineException extends Exception
type CmdLineException struct {
	jutil.JavaException

	parser *CmdLineParser

	// The optional localized message.
	//
	// See Throwable.getLocalizedMessage().
	localizedMessage *string
}

func (*CmdLineException) cmdLineException() {}

// ICmdLineException matches CmdLineException and its subclasses.
type ICmdLineException interface {
	jutil.IException
	cmdLineException()
	GetLocalizedMessage() string
	GetParser() *CmdLineParser
}

func newCmdLineException(message *string, cause jutil.Throwable) *CmdLineException {
	e := &CmdLineException{}
	jutil.InitThrowable(e, cmdLineException_className, message, cause)
	return e
}

// Deprecated: Use NewCmdLineExceptionWithParser.
//
// Java: CmdLineException(String)
func NewCmdLineException(message string) *CmdLineException {
	return newCmdLineException(&message, nil)
}

// Deprecated: Use NewCmdLineExceptionWithParserMessageAndCause.
//
// Java: CmdLineException(String, Throwable)
func NewCmdLineExceptionWithMessageAndCause(message string, cause jutil.Throwable) *CmdLineException {
	return newCmdLineException(&message, cause)
}

// Deprecated: Use NewCmdLineExceptionWithParserAndCause.
//
// Java: CmdLineException(Throwable)
func NewCmdLineExceptionWithCause(cause jutil.Throwable) *CmdLineException {
	return newCmdLineException(nil, cause)
}

// Java: CmdLineException(CmdLineParser, Localizable, String...)
func NewCmdLineExceptionWithLocalizable(parser *CmdLineParser, message Localizable, args ...string) *CmdLineException {
	objects := make([]any, len(args))
	for i, a := range args {
		objects[i] = a
	}
	m := message.FormatWithLocale(Locale_ENGLISH, objects...)
	e := newCmdLineException(&m, nil)
	localized := message.Format(objects...)
	e.localizedMessage = &localized
	e.parser = parser
	return e
}

// Deprecated: Use NewCmdLineExceptionWithLocalizable.
//
// Java: CmdLineException(CmdLineParser, String)
func NewCmdLineExceptionWithParser(parser *CmdLineParser, message string) *CmdLineException {
	e := newCmdLineException(&message, nil)
	e.parser = parser
	return e
}

// Java: CmdLineException(CmdLineParser, String, Throwable)
func NewCmdLineExceptionWithParserMessageAndCause(parser *CmdLineParser, message string, cause jutil.Throwable) *CmdLineException {
	e := newCmdLineException(&message, cause)
	e.parser = parser
	return e
}

// Java: CmdLineException(CmdLineParser, Throwable)
func NewCmdLineExceptionWithParserAndCause(parser *CmdLineParser, cause jutil.Throwable) *CmdLineException {
	e := newCmdLineException(nil, cause)
	e.parser = parser
	return e
}

// Java: CmdLineException.getLocalizedMessage()
func (e *CmdLineException) GetLocalizedMessage() string {
	if e.localizedMessage != nil {
		return *e.localizedMessage
	} else {
		return e.GetMessage()
	}
}

// Error is Throwable.toString(), which prints getLocalizedMessage().
func (e *CmdLineException) Error() string {
	if jutil.GetMessageOrNull(e) == nil && e.localizedMessage == nil {
		return cmdLineException_className
	}
	return cmdLineException_className + ": " + e.GetLocalizedMessage()
}

// Obtains the CmdLineParser that triggered an exception.
//
// Unless you have legacy OptionHandler that doesn't pass in this information
// when it throws an exception, this method should always return a non-null
// value.
//
// Java: CmdLineException.getParser()
func (e *CmdLineException) GetParser() *CmdLineParser {
	return e.parser
}
