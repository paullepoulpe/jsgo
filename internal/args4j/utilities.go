// Port of org/kohsuke/args4j/Utilities.java @ args4j 2.33

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

// Utilities are misc utility methods. Don't make this class visible to the
// outside world. When we switch to JDK 1.7, re-check the sense of this
// class.
//
// Java: class Utilities
type Utilities struct{}

// Utilities_CheckNonNull is similar to Objects.requireNonNull(). But this one
// is available for JDK 1.6 which is the current target of args4j. I didn't
// want to break compatibility with JDK 1.6. obj is the object to check for
// null value; name is the object name. If obj is null, then an exception is
// constructed from this name.
//
// Java: Utilities.checkNonNull(Object, String)
func Utilities_CheckNonNull(obj any, name string) {
	if jutil.IsNil(obj) {
		panic(jutil.NewNullPointerExceptionMsg(name + " is null"))
	}
}
