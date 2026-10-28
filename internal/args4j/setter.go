// Port of org/kohsuke/args4j/spi/Setter.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/Getter.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/Setters.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/FieldSetter.java @ args4j 2.33
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

import (
	"fmt"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Setter writes parsed option values into a Bean property.
type Setter interface {
	AddValue(value any)
	GetType() *Class
	IsMultiValued() bool
	AsFieldSetter() *FieldSetter
	AsAnnotatedElement() AnnotatedElement
}

// StringSetter is an optional fast path implemented by Setters that accept
// unboxed Go strings without interface boxing.
type StringSetter interface {
	AddStringValue(value string)
}

// AnnotatedElement represents the Field carrying an Option or Argument.
type AnnotatedElement interface {
	annotatedElement()
}

// Getter reads the current default value(s) of a Bean property for usage output.
type Getter interface {
	GetValueList() []any
}

func newSetterForField(f *Field, bean any) Setter {
	if f.IsFinal() {
		panic(jutil.NewIllegalStateExceptionMsg("Cannot set value to a final field '" + f.GetDeclaringClassName() + "." + f.GetName() + "'."))
	}
	if Class_List.IsAssignableFrom(f.GetType()) {
		return NewMultiValueFieldSetter(bean, f)
	}
	return NewFieldSetter(bean, f)
}

// FieldSetter writes single values to a Field.
type FieldSetter struct {
	f    *Field
	bean any
}

func NewFieldSetter(bean any, f *Field) *FieldSetter {
	return &FieldSetter{f: f, bean: bean}
}

func (s *FieldSetter) GetType() *Class                      { return s.f.GetType() }
func (s *FieldSetter) IsMultiValued() bool                  { return false }
func (s *FieldSetter) AsFieldSetter() *FieldSetter          { return NewFieldSetter(s.bean, s.f) }
func (s *FieldSetter) AsAnnotatedElement() AnnotatedElement { return s.f }
func (s *FieldSetter) AddValue(value any)                   { s.f.Set(s.bean, value) }
func (s *FieldSetter) GetValue() any                        { return s.f.Get(s.bean) }

func (s *FieldSetter) AddStringValue(value string) {
	if s.f.setString != nil {
		s.f.setString(value)
		return
	}
	s.AddValue(value)
}

func (s *FieldSetter) GetValueList() []any {
	if v := s.GetValue(); v != nil {
		return []any{v}
	}
	return []any{}
}

// MultiValueFieldSetter appends multiple values to a List Field.
type MultiValueFieldSetter struct {
	bean any
	f    *Field
}

func NewMultiValueFieldSetter(bean any, f *Field) *MultiValueFieldSetter {
	if !Class_List.IsAssignableFrom(f.GetType()) {
		panic(NewIllegalAnnotationError(spiMsgIllegalFieldSignature.format(f.GetType())))
	}
	return &MultiValueFieldSetter{bean: bean, f: f}
}

func (s *MultiValueFieldSetter) IsMultiValued() bool                  { return true }
func (s *MultiValueFieldSetter) AsFieldSetter() *FieldSetter          { return NewFieldSetter(s.bean, s.f) }
func (s *MultiValueFieldSetter) AsAnnotatedElement() AnnotatedElement { return s.f }

func (s *MultiValueFieldSetter) GetType() *Class {
	if t := s.f.GetGenericTypeArgument(); t != nil {
		return t
	}
	return Class_Object
}

func (s *MultiValueFieldSetter) AddValue(value any) {
	s.doAddValue(s.bean, value)
}

func (s *MultiValueFieldSetter) AddStringValue(value string) {
	if s.f.addString != nil {
		s.f.addString(value)
		return
	}
	s.doAddValue(s.bean, value)
}

func (s *MultiValueFieldSetter) doAddValue(bean any, value any) {
	o := s.f.Get(bean)
	if o == nil {
		o = jutil.NewArrayList[any]()
		s.f.Set(bean, o)
	}
	l, ok := o.(List)
	if !ok {
		panic(NewIllegalAnnotationError(spiMsgIllegalList.format(s.f)))
	}
	l.Add(value)
}

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

// Class describes the Java type of a Bean Field for handler dispatch and
// enum constant lookup. Classes are compared by pointer identity.
type Class struct {
	name        string
	isInterface bool
	isPrimitive bool
	isEnum      bool
	constants   []any
	supertypes  []*Class
}

var (
	Class_Object     = &Class{name: "java.lang.Object"}
	Class_Enum       = &Class{name: "java.lang.Enum"}
	Class_String     = &Class{name: "java.lang.String"}
	Class_Boolean    = &Class{name: "java.lang.Boolean"}
	Class_Integer    = &Class{name: "java.lang.Integer"}
	Class_Collection = &Class{name: "java.util.Collection", isInterface: true}
	Class_List       = &Class{name: "java.util.List", isInterface: true, supertypes: []*Class{Class_Collection}}
	Class_Map        = &Class{name: "java.util.Map", isInterface: true}
	Class_boolean    = &Class{name: "boolean", isPrimitive: true}
	Class_int        = &Class{name: "int", isPrimitive: true}
)

// NewEnumClass creates the Class descriptor for a Java enum with its constants
// in ordinal order.
func NewEnumClass(name string, constants []any) *Class {
	return &Class{name: name, isEnum: true, constants: constants, supertypes: []*Class{Class_Enum}}
}

// EnumClassOf builds a Class descriptor for a typed slice of enum values.
func EnumClassOf[E fmt.Stringer](name string, values []E) *Class {
	constants := make([]any, len(values))
	for i, v := range values {
		constants[i] = v
	}
	return NewEnumClass(name, constants)
}

func (c *Class) GetName() string { return c.name }

func (c *Class) String() string {
	if c.isInterface {
		return "interface " + c.name
	}
	if c.isPrimitive {
		return c.name
	}
	return "class " + c.name
}

func (c *Class) GetEnumConstants() []any {
	if !c.isEnum {
		return nil
	}
	return slices.Clone(c.constants)
}

func (c *Class) IsAssignableFrom(cls *Class) bool {
	if c == cls {
		return true
	}
	if cls.isPrimitive {
		return false
	}
	if c == Class_Object {
		return true
	}
	for _, s := range cls.supertypes {
		if c.IsAssignableFrom(s) {
			return true
		}
	}
	return false
}

func enumName(constant any) string {
	if n, ok := constant.(interface{ Name() string }); ok {
		return n.Name()
	}
	return constant.(fmt.Stringer).String()
}

// Field binds an option bean's Go field to its declared Class and typed accessors.
type Field struct {
	declaringClass string
	name           string
	type_          *Class
	typeArgument   *Class
	final          bool
	get            func() any
	set            func(value any)
	setString      func(value string)
	addString      func(value string)
}

func NewField(declaringClass, name string, type_, typeArgument *Class, final bool, get func() any, set func(value any)) *Field {
	return &Field{
		declaringClass: declaringClass,
		name:           name,
		type_:          type_,
		typeArgument:   typeArgument,
		final:          final,
		get:            get,
		set:            set,
	}
}

func (f *Field) GetDeclaringClassName() string  { return f.declaringClass }
func (f *Field) GetName() string                { return f.name }
func (f *Field) GetType() *Class                { return f.type_ }
func (f *Field) GetGenericTypeArgument() *Class { return f.typeArgument }
func (f *Field) IsFinal() bool                  { return f.final }
func (f *Field) Get(_ any) any                  { return f.get() }
func (f *Field) Set(_ any, value any)           { f.set(value) }
func (f *Field) String() string                 { return f.type_.GetName() + " " + f.declaringClass + "." + f.name }
func (*Field) annotatedElement()                {}

// List is the collection interface used by MultiValueFieldSetter.
type List interface {
	Add(e any) bool
	Size() int
	Get(index int) any
}

// Map is the key-value interface populated by MapOptionHandler.
type Map interface {
	Put(key string, value *string) *string
}

func valueField[T any](declaringClass, name string, type_ *Class, p *T) *Field {
	return NewField(declaringClass, name, type_, nil, false,
		func() any { return *p },
		func(v any) { *p = v.(T) })
}

func nullableField[T any](declaringClass, name string, type_ *Class, p **T) *Field {
	return NewField(declaringClass, name, type_, nil, false,
		func() any {
			if *p == nil {
				return nil
			}
			return **p
		},
		func(v any) {
			if v == nil {
				*p = nil
				return
			}
			x := v.(T)
			*p = &x
		})
}

func BooleanField(declaringClass, name string, p *bool) *Field {
	return valueField(declaringClass, name, Class_boolean, p)
}

func NullableBooleanField(declaringClass, name string, p **bool) *Field {
	return nullableField(declaringClass, name, Class_Boolean, p)
}

func IntField(declaringClass, name string, p *int) *Field {
	return NewField(declaringClass, name, Class_int, nil, false,
		func() any { return int32(*p) },
		func(v any) { *p = int(v.(int32)) })
}

func IntegerField(declaringClass, name string, p *int32) *Field {
	return valueField(declaringClass, name, Class_Integer, p)
}

func NullableIntegerField(declaringClass, name string, p **int32) *Field {
	return nullableField(declaringClass, name, Class_Integer, p)
}

func StringField(declaringClass, name string, p *string) *Field {
	f := valueField(declaringClass, name, Class_String, p)
	f.setString = func(v string) { *p = v }
	return f
}

func NullableStringField(declaringClass, name string, p **string) *Field {
	f := nullableField(declaringClass, name, Class_String, p)
	f.setString = func(v string) { *p = &v }
	return f
}

func EnumField[E any](declaringClass, name string, enumClass *Class, p *E) *Field {
	return valueField(declaringClass, name, enumClass, p)
}

func NullableEnumField[E any](declaringClass, name string, enumClass *Class, p **E) *Field {
	return nullableField(declaringClass, name, enumClass, p)
}

func ListField[E any](declaringClass, name string, elementType *Class, p *[]E) *Field {
	sl := &sliceList[E]{p}
	f := NewField(declaringClass, name, Class_List, elementType, false,
		func() any { return sl },
		func(v any) {
			l := v.(List)
			s := make([]E, l.Size())
			for i := range s {
				s[i] = l.Get(i).(E)
			}
			*p = s
		})
	if sp, ok := any(p).(*[]string); ok {
		f.addString = func(v string) {
			*sp = append(*sp, v)
		}
	}
	return f
}

type sliceList[E any] struct{ p *[]E }

func (l *sliceList[E]) Add(e any) bool {
	*l.p = append(*l.p, e.(E))
	return true
}
func (l *sliceList[E]) Size() int         { return len(*l.p) }
func (l *sliceList[E]) Get(index int) any { return (*l.p)[index] }

func StringMapField(declaringClass, name string, p **jutil.JavaHashMap[string, *string]) *Field {
	return NewField(declaringClass, name, Class_Map, Class_String, false,
		func() any {
			if *p == nil {
				return nil
			}
			return *p
		},
		func(v any) {
			if v == nil {
				*p = nil
				return
			}
			*p = v.(*jutil.JavaHashMap[string, *string])
		})
}

func LinkedStringMapField(declaringClass, name string, p **jutil.LinkedHashMap[string, *string]) *Field {
	return NewField(declaringClass, name, Class_Map, Class_String, false,
		func() any {
			if *p == nil {
				return nil
			}
			return *p
		},
		func(v any) {
			if v == nil {
				*p = nil
				return
			}
			*p = v.(*jutil.LinkedHashMap[string, *string])
		})
}
