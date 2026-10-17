// Port of org/kohsuke/args4j/OptionHandlerFilter.java @ args4j 2.33

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

// OptionHandlerFilter selects OptionHandler.
//
// For example, we use this to let the caller specify which options are
// printed, and which ones aren't.
//
// See CmdLineParser.printExample(OptionHandlerFilter) and
// CmdLineParser.printUsage(Writer, ResourceBundle, OptionHandlerFilter).
//
// Java: interface OptionHandlerFilter
type OptionHandlerFilter interface {
	// Select reports whether to choose o (true) or to
	// ignore/discard/disregard it (false). o is never nil. Internally options
	// (like -r) and arguments (others) are treated uniformly as
	// OptionHandler. See OptionDef.isArgument() to distinguish them.
	//
	// Java: OptionHandlerFilter.select(OptionHandler)
	Select(o OptionHandler) bool
}

// OptionHandlerFilterFunc adapts a function (a Java lambda) to
// OptionHandlerFilter.
//
// Go-only: language-adapter
type OptionHandlerFilterFunc func(o OptionHandler) bool

// Select calls f.
func (f OptionHandlerFilterFunc) Select(o OptionHandler) bool { return f(o) }

// optionHandlerFilter_anon is an anonymous OptionHandlerFilter class.
type optionHandlerFilter_anon struct {
	sel func(o OptionHandler) bool
}

func (f *optionHandlerFilter_anon) Select(o OptionHandler) bool { return f.sel(o) }

// OptionHandlerFilter_ALL prints all defined options in the example.
//
// Java: OptionHandlerFilter.ALL
var OptionHandlerFilter_ALL OptionHandlerFilter = &optionHandlerFilter_anon{func(o OptionHandler) bool {
	return true
}}

// OptionHandlerFilter_PUBLIC prints all non-hidden options.
//
// This would only be useful with a small number of options.
//
// Java: OptionHandlerFilter.PUBLIC
var OptionHandlerFilter_PUBLIC OptionHandlerFilter = &optionHandlerFilter_anon{func(o OptionHandler) bool {
	return !o.GetOptionHandlerBase().Option.Hidden()
}}

// OptionHandlerFilter_REQUIRED prints all required options.
//
// Java: OptionHandlerFilter.REQUIRED
var OptionHandlerFilter_REQUIRED OptionHandlerFilter = &optionHandlerFilter_anon{func(o OptionHandler) bool {
	return o.GetOptionHandlerBase().Option.Required()
}}
