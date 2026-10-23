// Port of org/kohsuke/args4j/spi/Setter.java @ args4j 2.33

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

// Setter is the abstraction of the value setter.
//
// This abstracts away the difference between a field and a setter method,
// which object we are setting the value to, and/or how we handle collection
// fields differently.
//
// See Getter.
//
// Go-only: language-adapter — generics: Setter<T> takes its values as any,
// boxed as the package documentation describes.
//
// Java: interface spi.Setter<T>
type Setter interface {
	// AddValue adds/sets a value to the property of the option bean.
	//
	// A Setter object has an implicit knowledge about the property it's
	// setting, and the instance of the option bean. It may throw a
	// CmdLineException.
	//
	// Java: Setter.addValue(T)
	AddValue(value any)

	// GetType gets the type of the underlying method/field.
	//
	// Java: Setter.getType()
	GetType() *Class

	// IsMultiValued is whether this setter is intrinsically multi-valued.
	//
	// When parsing arguments (instead of options), intrinsically multi-valued
	// setters consume all the remaining arguments. So, if the setter can
	// store multiple values, this method should return true.
	//
	// This characteristics of a setter doesn't affect option parsing at all;
	// any options can be specified multiple times. In many cases, this is a
	// no-op--but when your shell script expands multiple environment
	// variables (each of which may contain options), tolerating such
	// redundant options can be useful.
	//
	// Java: Setter.isMultiValued()
	IsMultiValued() bool

	// AsFieldSetter returns, if this setter encapsulates a field, the direct
	// access to that field as FieldSetter. This method serves two purposes:
	//
	//  1. This lets OptionHandlers bypass the collection/array handling of
	//     fields. This is useful if you're defining an option handler that
	//     produces array or collection from a single argument.
	//  2. The other is to retrieve the current value of the field (via
	//     FieldSetter.getValueList()).
	//
	// It returns nil if this setter wraps a method.
	//
	// Java: Setter.asFieldSetter()
	AsFieldSetter() *FieldSetter

	// AsAnnotatedElement returns the AnnotatedElement by which you can access
	// annotations written on this setter.
	//
	// This is the same AnnotatedElement that had Option/Argument.
	//
	// This enables OptionHandler to further tweak its behavior based on
	// additional annotations.
	//
	// Java: Setter.asAnnotatedElement()
	AsAnnotatedElement() AnnotatedElement
}

// StringSetter is an optional fast path implemented by Setters that accept
// unboxed Go strings without interface boxing.
type StringSetter interface {
	AddStringValue(value string)
}

// AnnotatedElement is java.lang.reflect.AnnotatedElement: the field (a
// *Field) that carried the Option or Argument annotation.
//
// Go-only: language-adapter — reflection.
type AnnotatedElement interface {
	annotatedElement()
}
