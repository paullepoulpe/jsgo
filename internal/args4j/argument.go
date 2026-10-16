// Port of org/kohsuke/args4j/Argument.java @ args4j 2.33

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

// Argument of the command line.
//
// This works mostly like Option except the following differences.
//
//  1. Arguments have an index about their relative position on the command
//     line.
//
// Go-only: language-adapter — annotations: the annotation is a struct whose
// fields are its elements; a field left at its zero value has the
// element's default (see the package documentation).
//
// Java: @interface Argument
type Argument struct {
	// See Option.Usage.
	Usage string

	// See Option.MetaVar.
	MetaVar string

	// See Option.Required.
	Required bool

	// See Option.Hidden.
	Hidden bool

	// See Option.Handler.
	Handler OptionHandlerClass

	// Position of the argument.
	//
	// If you define multiple single value properties to bind to arguments,
	// they should have index=0, index=1, index=2, ... and so on.
	//
	// Multi value properties bound to arguments must be always the last
	// entry.
	Index int

	MultiValued bool
}
