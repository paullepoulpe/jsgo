// Port of org/kohsuke/args4j/spi/IntOptionHandler.java @ args4j 2.33

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

// IntOptionHandler is the Integer OptionHandler (a OneArgumentOptionHandler).
//
// Java: class spi.IntOptionHandler extends OneArgumentOptionHandler<Integer>
type IntOptionHandler struct {
	OneArgumentOptionHandlerBase
}

// NewIntOptionHandler returns the handler as an OptionHandler, so that it is
// an OptionHandlerClass (Go-only: language-adapter).
//
// Java: IntOptionHandler(CmdLineParser, OptionDef, Setter)
func NewIntOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &IntOptionHandler{}
	h.InitOneArgumentOptionHandlerBase(h, parser, option, setter)
	return h
}

// Parse returns an int32 (a boxed Integer).
//
// Java: IntOptionHandler.parse(String)
func (h *IntOptionHandler) Parse(argument string) any {
	return jutil.IntegerParseInt(argument, 10)
}
