// Port of org/kohsuke/args4j/ClassParser.java @ args4j 2.33

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

// ClassParser is the parser for analyzing Args4J annotations in the class
// hierarchy.
//
// This can be used to feed option bindings that span across multiple
// instances.
//
// Java: class ClassParser
type ClassParser struct {
	// Go-only: language-adapter — distinct identity.
	_ byte
}

// Java: ClassParser()
func NewClassParser() *ClassParser { return &ClassParser{} }

// Bean is an option bean: an object whose fields carry the @Option and
// @Argument annotations.
//
// Go-only: language-adapter — annotations (see the package documentation).
// Java finds the annotated fields by reflection; a Go bean lists them.
type Bean interface {
	// Args4jFields returns the annotated fields of the bean in the order in
	// which ClassParser.parse visits them: for the bean's class, then for
	// each of its superclasses, the fields that class declares, in
	// declaration order (Class.getDeclaredFields()).
	Args4jFields() []AnnotatedField
}

// AnnotatedField is a field of an option bean with its args4j annotations:
// Option or Argument, or both.
//
// Go-only: language-adapter — annotations.
type AnnotatedField struct {
	Field    *Field
	Option   *Option
	Argument *Argument
}

// Java: ClassParser.parse(Object, CmdLineParser)
func (p *ClassParser) Parse(bean any, parser *CmdLineParser) {
	// recursively process all the methods/fields.
	//
	// Go-only: language-adapter — annotations: annotated methods are not
	// ported (MethodSetter), and the fields of the class hierarchy come from
	// Bean.Args4jFields in the order of Java's loops. A value that is not a
	// Bean has no annotated fields.
	b, ok := bean.(Bean)
	if !ok {
		return
	}
	for _, af := range b.Args4jFields() {
		f := af.Field
		o := af.Option
		if o != nil {
			parser.AddOption(Setters_Create(f, bean), o)
		}
		a := af.Argument
		if a != nil {
			parser.AddArgument(Setters_Create(f, bean), a)
		}
	}
}
