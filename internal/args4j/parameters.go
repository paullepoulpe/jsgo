// Port of org/kohsuke/args4j/spi/Parameters.java @ args4j 2.33

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

// Parameters is a list of arguments.
//
// Object of this interface is passed to OptionHandlers to make it easy/safe
// to parse additional parameters for options.
//
// Java: interface spi.Parameters
type Parameters interface {
	// GetParameter gets the additional parameter to this option.
	//
	// idx: specifying 0 will retrieve the token next to the option. For
	// example, if the command line looks like -o abc -d x, then
	// getParameter(0) for -o returns abc and getParameter(1) will return -d.
	//
	// It always returns a valid String. If an attempt is made to access a
	// non-existent index, this method throws appropriate CmdLineException.
	//
	// Java: Parameters.getParameter(int)
	GetParameter(idx int) string

	// Size is the number of remaining tokens.
	//
	// Java: Parameters.size()
	Size() int
}
