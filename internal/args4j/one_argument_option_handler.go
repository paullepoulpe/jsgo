// Port of org/kohsuke/args4j/spi/OneArgumentOptionHandler.java @ args4j 2.33

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

// OneArgumentOptionHandler is an OptionHandler which handles an option with
// exactly one argument, like "-foo bar".
//
// T is the type of the Setter-class.
//
// Java: abstract class spi.OneArgumentOptionHandler<T> extends OptionHandler<T>
type OneArgumentOptionHandler interface {
	OptionHandler

	// Parse parses a string to a real value of Type <T>. argument is the
	// String value to parse.
	//
	// It throws NumberFormatException if parsing is not possible, and
	// CmdLineException if the parsing encounters a failure that should be
	// reported to the user.
	//
	// Java: OneArgumentOptionHandler.parse(String)
	Parse(argument string) any
}

// OneArgumentOptionHandlerBase holds the implementation of
// OneArgumentOptionHandler.
type OneArgumentOptionHandlerBase struct {
	OptionHandlerBase
}

// Java: OneArgumentOptionHandler(CmdLineParser, OptionDef, Setter)
func (h *OneArgumentOptionHandlerBase) InitOneArgumentOptionHandlerBase(self OneArgumentOptionHandler, parser *CmdLineParser, option OptionDef, setter Setter) {
	h.InitOptionHandlerBase(self, parser, option, setter)
}

// Java: OneArgumentOptionHandler.getDefaultMetaVariable()
func (h *OneArgumentOptionHandlerBase) GetDefaultMetaVariable() *string {
	s := "N"
	return &s
}

// Java: OneArgumentOptionHandler.parseArguments(Parameters)
func (h *OneArgumentOptionHandlerBase) ParseArguments(params Parameters) int {
	token := params.GetParameter(0)
	jutil.Try(func() {
		value := h.self.(OneArgumentOptionHandler).Parse(token)
		h.Setter.AddValue(value)
	}, func(ex jutil.INumberFormatException) {
		panic(NewCmdLineExceptionWithLocalizable(h.Owner, SpiMessages_ILLEGAL_OPERAND, params.GetParameter(-1), token))
	})
	return 1
}
