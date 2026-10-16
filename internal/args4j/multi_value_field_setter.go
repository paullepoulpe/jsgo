// Port of org/kohsuke/args4j/spi/MultiValueFieldSetter.java @ args4j 2.33

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

// MultiValueFieldSetter is a Setter that sets multiple values to a
// collection Field.
//
// Java: final class spi.MultiValueFieldSetter implements Getter, Setter
type MultiValueFieldSetter struct {
	bean any
	f    *Field
}

// Java: MultiValueFieldSetter(Object, Field)
func NewMultiValueFieldSetter(bean any, f *Field) *MultiValueFieldSetter {
	s := &MultiValueFieldSetter{}
	s.bean = bean
	s.f = f

	if !Class_List.IsAssignableFrom(f.GetType()) {
		panic(NewIllegalAnnotationError(SpiMessages_ILLEGAL_FIELD_SIGNATURE.Format(f.GetType())))
	}
	return s
}

// Java: MultiValueFieldSetter.isMultiValued()
func (s *MultiValueFieldSetter) IsMultiValued() bool {
	return true
}

// Java: MultiValueFieldSetter.asFieldSetter()
func (s *MultiValueFieldSetter) AsFieldSetter() *FieldSetter {
	return NewFieldSetter(s.bean, s.f)
}

// Java: MultiValueFieldSetter.asAnnotatedElement()
func (s *MultiValueFieldSetter) AsAnnotatedElement() AnnotatedElement {
	return s.f
}

// Java: MultiValueFieldSetter.getType()
func (s *MultiValueFieldSetter) GetType() *Class {
	// TODO: compute this correctly
	if t := s.f.GetGenericTypeArgument(); t != nil {
		return t
	}
	return Class_Object
}

// Java: MultiValueFieldSetter.addValue(Object)
func (s *MultiValueFieldSetter) AddValue(value any) {
	// Go-only: language-adapter — reflection: a Go Field has no access
	// control, so Java's retry after setAccessible(true) (and its
	// IllegalAccessError) cannot happen.
	s.doAddValue(s.bean, value)
}

// Java: MultiValueFieldSetter.doAddValue(Object, Object)
func (s *MultiValueFieldSetter) doAddValue(bean any, value any) {
	o := s.f.Get(bean)
	if o == nil {
		o = jutil.NewArrayList[any]()
		s.f.Set(bean, o)
	}
	l, ok := o.(List)
	if !ok {
		panic(NewIllegalAnnotationError(SpiMessages_ILLEGAL_LIST.Format(s.f)))
	}

	l.Add(value)
}

// Java: MultiValueFieldSetter.getValueList()
func (s *MultiValueFieldSetter) GetValueList() []any {
	o := s.f.Get(s.bean)
	if o == nil {
		return nil
	}
	l := o.(List)
	values := make([]any, l.Size())
	for i := range values {
		values[i] = l.Get(i)
	}
	return values
}
