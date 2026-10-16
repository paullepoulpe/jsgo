// Port of org/kohsuke/args4j/spi/BooleanOptionHandler.java @ args4j 2.33

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

// BooleanOptionHandler is the Boolean OptionHandler.
//
// Java: class spi.BooleanOptionHandler extends OptionHandler<Boolean>
type BooleanOptionHandler struct {
	OptionHandlerBase
}

var booleanOptionHandler_ACCEPTABLE_VALUES = []string{"true", "on", "yes", "1",
	"false", "off", "no", "0"}

// NewBooleanOptionHandler returns the handler as an OptionHandler, so that
// it is an OptionHandlerClass (Go-only: language-adapter).
//
// Java: BooleanOptionHandler(CmdLineParser, OptionDef, Setter)
func NewBooleanOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &BooleanOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	return h
}

// Java: BooleanOptionHandler.parseArguments(Parameters)
func (h *BooleanOptionHandler) ParseArguments(params Parameters) int {
	if h.Option.IsArgument() {
		// Go-only: String.toLowerCase() uses the default locale, which is
		// English in the port (see locale.go): it lowercases as
		// toLowerCase(Locale.ROOT) does.
		valueStr := jutil.StringToLowerCaseRoot(params.GetParameter(0))
		index := -1
		for i, v := range booleanOptionHandler_ACCEPTABLE_VALUES {
			if v == valueStr {
				index = i
				break
			}
		}
		if index == -1 {
			panic(NewCmdLineExceptionWithLocalizable(h.Owner, SpiMessages_ILLEGAL_BOOLEAN, valueStr))
		}
		h.Setter.AddValue(index < len(booleanOptionHandler_ACCEPTABLE_VALUES)/2)
		return 1
	} else {
		h.Setter.AddValue(true)
		return 0
	}
}

// Java: BooleanOptionHandler.getDefaultMetaVariable()
func (h *BooleanOptionHandler) GetDefaultMetaVariable() *string {
	return nil
}
