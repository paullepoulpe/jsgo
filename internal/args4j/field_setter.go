// Port of org/kohsuke/args4j/spi/FieldSetter.java @ args4j 2.33

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

// FieldSetter is a Setter that sets to a Field.
//
// Java: final class spi.FieldSetter implements Getter, Setter
type FieldSetter struct {
	f    *Field
	bean any
}

// Java: FieldSetter(Object, Field)
func NewFieldSetter(bean any, f *Field) *FieldSetter {
	s := &FieldSetter{}
	s.bean = bean
	s.f = f
	return s
}

// Java: FieldSetter.getType()
func (s *FieldSetter) GetType() *Class {
	return s.f.GetType()
}

// Java: FieldSetter.isMultiValued()
func (s *FieldSetter) IsMultiValued() bool {
	// a field can only store one value. a collection field is handled via
	// MultiValueFieldSetter
	return false
}

// Java: FieldSetter.asFieldSetter()
func (s *FieldSetter) AsFieldSetter() *FieldSetter {
	return NewFieldSetter(s.bean, s.f)
}

// Java: FieldSetter.asAnnotatedElement()
func (s *FieldSetter) AsAnnotatedElement() AnnotatedElement {
	return s.f
}

// Java: FieldSetter.addValue(Object)
func (s *FieldSetter) AddValue(value any) {
	// Go-only: language-adapter — reflection: a Go Field has no access
	// control, so Java's retry after setAccessible(true) (and its
	// IllegalAccessError) cannot happen.
	s.f.Set(s.bean, value)
}

// AddStringValue implements StringSetter without boxing value into any.
func (s *FieldSetter) AddStringValue(value string) {
	if s.f.setString != nil {
		s.f.setString(value)
		return
	}
	s.AddValue(value)
}

// Java: FieldSetter.getValue()
func (s *FieldSetter) GetValue() any {
	// Go-only: language-adapter — reflection (see AddValue).
	return s.f.Get(s.bean)
}

// Java: FieldSetter.getValueList()
func (s *FieldSetter) GetValueList() []any {
	return s.asList(s.GetValue())
}

// Java: FieldSetter.asList(Object)
func (s *FieldSetter) asList(o any) []any {
	if o != nil {
		return []any{o}
	}
	return []any{}
}
