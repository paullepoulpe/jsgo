// Port of org/kohsuke/args4j/spi/Setter.java @ args4j 2.33
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Setter writes parsed option values into a Bean property.
type Setter interface {
	AddValue(value any)
	AddStringValue(value string)
	GetType() *Class
	IsMultiValued() bool
	AsFieldSetter() *FieldSetter
}

func newSetterForField(f *Field) Setter {
	if f.final {
		panic(jutil.NewIllegalStateExceptionMsg("Cannot set value to a final field '" + f.declaringClass + "." + f.name + "'."))
	}
	return &FieldSetter{f: f, multiValued: f.type_ == classList}
}

// FieldSetter writes single or list values to a Field.
type FieldSetter struct {
	f           *Field
	multiValued bool
}

func (s *FieldSetter) IsMultiValued() bool { return s.multiValued }
func (s *FieldSetter) getValue() any       { return s.f.get() }

func (s *FieldSetter) AsFieldSetter() *FieldSetter {
	if s.multiValued {
		return &FieldSetter{f: s.f}
	}
	return s
}

func (s *FieldSetter) GetType() *Class {
	if !s.multiValued {
		return s.f.type_
	}
	if s.f.typeArgument != nil {
		return s.f.typeArgument
	}
	return classObject
}

func (s *FieldSetter) AddValue(value any) {
	if !s.multiValued {
		s.f.set(value)
		return
	}
	s.f.add(value)
}

func (s *FieldSetter) AddStringValue(value string) {
	if !s.multiValued {
		if s.f.setString != nil {
			s.f.setString(value)
			return
		}
		s.f.set(value)
		return
	}
	if s.f.addString != nil {
		s.f.addString(value)
		return
	}
	s.f.add(value)
}

func (s *FieldSetter) getValueList() []any {
	if s.multiValued {
		return s.f.getList()
	}
	o := s.f.get()
	if o == nil {
		return nil
	}
	return []any{o}
}

// Class describes the Java type of a Bean Field for handler dispatch and
// enum constant lookup. Classes are compared by pointer identity.
type Class struct {
	name        string
	isInterface bool
	isPrimitive bool
	isEnum      bool
	constants   []any
}

var (
	Class_String = &Class{name: "java.lang.String"}

	classObject      = &Class{name: "java.lang.Object"}
	classBoolean     = &Class{name: "java.lang.Boolean"}
	classInteger     = &Class{name: "java.lang.Integer"}
	classList        = &Class{name: "java.util.List", isInterface: true}
	classMap         = &Class{name: "java.util.Map", isInterface: true}
	classBooleanPrim = &Class{name: "boolean", isPrimitive: true}
	classIntPrim     = &Class{name: "int", isPrimitive: true}
)

// EnumClassOf builds a Class descriptor for a typed slice of enum values.
func EnumClassOf[E fmt.Stringer](name string, values []E) *Class {
	constants := make([]any, len(values))
	for i, v := range values {
		constants[i] = v
	}
	return &Class{name: name, isEnum: true, constants: constants}
}

func (c *Class) String() string {
	if c.isInterface {
		return "interface " + c.name
	}
	if c.isPrimitive {
		return c.name
	}
	return "class " + c.name
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
	add            func(value any)
	addString      func(value string)
	getList        func() []any
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
	return valueField(declaringClass, name, classBooleanPrim, p)
}

func NullableBooleanField(declaringClass, name string, p **bool) *Field {
	return nullableField(declaringClass, name, classBoolean, p)
}

func IntField(declaringClass, name string, p *int) *Field {
	return NewField(declaringClass, name, classIntPrim, nil, false,
		func() any { return int32(*p) },
		func(v any) { *p = int(v.(int32)) })
}

func IntegerField(declaringClass, name string, p *int32) *Field {
	return valueField(declaringClass, name, classInteger, p)
}

func NullableIntegerField(declaringClass, name string, p **int32) *Field {
	return nullableField(declaringClass, name, classInteger, p)
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
	f := &Field{
		declaringClass: declaringClass,
		name:           name,
		type_:          classList,
		typeArgument:   elementType,
		add:            func(v any) { *p = append(*p, v.(E)) },
		getList: func() []any {
			if len(*p) == 0 {
				return nil
			}
			out := make([]any, len(*p))
			for i, e := range *p {
				out[i] = e
			}
			return out
		},
	}
	if sp, ok := any(p).(*[]string); ok {
		f.addString = func(v string) {
			*sp = append(*sp, v)
		}
	}
	return f
}

func LinkedStringMapField(declaringClass, name string, p **jutil.LinkedHashMap[string, *string]) *Field {
	return NewField(declaringClass, name, classMap, Class_String, false,
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
