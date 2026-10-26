// Port of org/kohsuke/args4j/ExampleMode.java @ args4j 2.33

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

// ExampleMode is used with CmdLineParser.printExample(ExampleMode).
//
// Deprecated: use OptionHandlerFilter.
//
// Java: enum ExampleMode implements OptionHandlerFilter
type ExampleMode int32

const (
	// ExampleMode_ALL prints all defined options in the example.
	//
	// This would only be useful with small numbers of options.
	ExampleMode_ALL ExampleMode = iota

	// ExampleMode_REQUIRED prints all required option.
	ExampleMode_REQUIRED
)

var exampleMode_names = [...]string{"ALL", "REQUIRED"}

// String is ExampleMode.name().
func (m ExampleMode) String() string { return exampleMode_names[m] }

// Java: ExampleMode.select(OptionHandler)
func (m ExampleMode) Select(o OptionHandler) bool {
	switch m {
	case ExampleMode_ALL:
		return true
	default: // REQUIRED
		return o.GetOptionHandlerBase().Option.Required()
	}
}
