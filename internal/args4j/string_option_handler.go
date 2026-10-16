// Port of org/kohsuke/args4j/spi/StringOptionHandler.java @ args4j 2.33

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

// StringOptionHandler is the String OptionHandler.
//
// Subclasses embed it and call InitStringOptionHandler.
//
// Java: class spi.StringOptionHandler extends OptionHandler<String>
type StringOptionHandler struct {
	OptionHandlerBase
}

// NewStringOptionHandler returns the handler as an OptionHandler, so that it
// is an OptionHandlerClass (Go-only: language-adapter).
//
// Java: StringOptionHandler(CmdLineParser, OptionDef, Setter)
func NewStringOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &StringOptionHandler{}
	h.InitStringOptionHandler(h, parser, option, setter)
	return h
}

// InitStringOptionHandler is the constructor as subclasses call it; self is
// the subclass instance.
//
// Java: StringOptionHandler(CmdLineParser, OptionDef, Setter)
func (h *StringOptionHandler) InitStringOptionHandler(self OptionHandler, parser *CmdLineParser, option OptionDef, setter Setter) {
	h.InitOptionHandlerBase(self, parser, option, setter)
}

// Java: StringOptionHandler.parseArguments(Parameters)
func (h *StringOptionHandler) ParseArguments(params Parameters) int {
	h.Setter.AddValue(params.GetParameter(0))
	return 1
}

// Java: StringOptionHandler.getDefaultMetaVariable()
func (h *StringOptionHandler) GetDefaultMetaVariable() *string {
	s := SpiMessages_DEFAULT_META_STRING_OPTION_HANDLER.Format()
	return &s
}
