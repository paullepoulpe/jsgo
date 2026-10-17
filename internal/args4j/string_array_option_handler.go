// Port of org/kohsuke/args4j/spi/StringArrayOptionHandler.java @ args4j 2.33

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

// StringArrayOptionHandler is an OptionHandler for greedily mapping a list
// of tokens into a collection of Strings (such as String[], List<String>,
// etc.).
//
// How it works: example for parameter -s, which is type String[]:
//
//	java -jar aaa.jar -s banan hruska jablko
//	java -jar aaa.jar -s banan "hruska jablko"
//	java -jar aaa.jar -s "banan hruska jablko"
//	java -jar aaa.jar -s banan hruska jablko -l 4 -r
//	java -jar aaa.jar -t 222 -s banan hruska jablko -r
//
// All of them result in a single string array that contains three tokens:
// banan, hruska, and jablko.
//
// This OptionHandler scans for parameter which begins with -. If found, it
// will stop.
//
// Java: class spi.StringArrayOptionHandler extends OptionHandler<String>
type StringArrayOptionHandler struct {
	OptionHandlerBase
}

// NewStringArrayOptionHandler returns the handler as an OptionHandler, so
// that it is an OptionHandlerClass (Go-only: language-adapter).
//
// Java: StringArrayOptionHandler(CmdLineParser, OptionDef, Setter)
func NewStringArrayOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &StringArrayOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	return h
}

// GetDefaultMetaVariable returns "STRING[]".
//
// Java: StringArrayOptionHandler.getDefaultMetaVariable()
func (h *StringArrayOptionHandler) GetDefaultMetaVariable() *string {
	s := SpiMessages_DEFAULT_META_STRING_ARRAY_OPTION_HANDLER.Format()
	return &s
}

// ParseArguments tries to parse String[] argument from Parameters.
//
// Java: StringArrayOptionHandler.parseArguments(Parameters)
func (h *StringArrayOptionHandler) ParseArguments(params Parameters) int {
	counter := 0
	for ; counter < params.Size(); counter++ {
		param := params.GetParameter(counter)

		if jutil.StringStartsWith(param, "-") {
			break
		}

		// String.split(" "): the regex matches only the literal space.
		for _, p := range jutil.StringSplitLiteral(param, " ") {
			h.Setter.AddValue(p)
		}
	}

	return counter
}
