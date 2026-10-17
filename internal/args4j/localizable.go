// Port of org/kohsuke/args4j/Localizable.java @ args4j 2.33

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

// Localizable is a message that can be formatted with arguments and locale.
//
// The message is implicitly given by this object and can contain formatting
// similar to the MessageFormat class.
//
// Java: interface Localizable
type Localizable interface {
	// FormatWithLocale formats the implicitly given message by this object
	// with the given locale. args are the arguments to use for formatting;
	// see MessageFormat.format(java.lang.String, java.lang.Object...).
	//
	// Java: Localizable.formatWithLocale(Locale, Object...)
	FormatWithLocale(locale *Locale, args ...any) string

	// Format formats the implicitly given message by this object with the
	// default locale.
	//
	// Java: Localizable.format(Object...)
	Format(args ...any) string
}
