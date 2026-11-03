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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Setter writes parsed option values into a Bean property.
type Setter interface {
	IsMultiValued() bool
	DefaultHandler() OptionHandlerClass
}

// TypedSetter writes parsed option values of type T without interface boxing.
type TypedSetter[T any] interface {
	Setter
	AddTypedValue(value T)
}

type defaultValueSource interface {
	defaultValueString() (string, bool)
}

// TypeSpec describes how args4j parses a single operand of type T.
type TypeSpec[T any] struct {
	defaultHandler OptionHandlerClass
}

var (
	BoolSpec   = &TypeSpec[bool]{defaultHandler: newBooleanOptionHandler}
	Int32Spec  = &TypeSpec[int32]{defaultHandler: NewIntOptionHandler}
	StringSpec = &TypeSpec[string]{defaultHandler: NewStringOptionHandler}
)

// EnumClassOf builds a typed TypeSpec for an enum's values slice.
func EnumClassOf[E fmt.Stringer](values []E) *TypeSpec[E] {
	return &TypeSpec[E]{defaultHandler: newEnumOptionHandler(values)}
}

func enumName[E fmt.Stringer](constant E) string {
	if n, ok := any(constant).(interface{ Name() string }); ok {
		return n.Name()
	}
	return constant.String()
}

// Field binds an option bean's Go field to a typed Setter.
type Field struct {
	setter   Setter
	finalErr string
}

func newSetterForField(f *Field) Setter {
	if f.finalErr != "" {
		panic(jutil.NewIllegalStateExceptionMsg(f.finalErr))
	}
	return f.setter
}

type scalarSetter[T any] struct {
	p       *T
	handler OptionHandlerClass
}

func (s *scalarSetter[T]) IsMultiValued() bool                { return false }
func (s *scalarSetter[T]) DefaultHandler() OptionHandlerClass { return s.handler }
func (s *scalarSetter[T]) AddTypedValue(v T)                  { *s.p = v }
func (s *scalarSetter[T]) defaultValueString() (string, bool) { return jutil.ToJavaString(*s.p), true }

type trackedScalarSetter[T any] struct {
	p         *T
	isDefault *bool
	handler   OptionHandlerClass
}

func (s *trackedScalarSetter[T]) IsMultiValued() bool                { return false }
func (s *trackedScalarSetter[T]) DefaultHandler() OptionHandlerClass { return s.handler }
func (s *trackedScalarSetter[T]) AddTypedValue(v T)                  { *s.p, *s.isDefault = v, false }
func (s *trackedScalarSetter[T]) defaultValueString() (string, bool) {
	return jutil.ToJavaString(*s.p), true
}

type nullableSetter[T any] struct {
	p       **T
	handler OptionHandlerClass
}

func (s *nullableSetter[T]) IsMultiValued() bool                { return false }
func (s *nullableSetter[T]) DefaultHandler() OptionHandlerClass { return s.handler }
func (s *nullableSetter[T]) AddTypedValue(v T)                  { *s.p = &v }
func (s *nullableSetter[T]) defaultValueString() (string, bool) {
	if *s.p == nil {
		return "", false
	}
	return jutil.ToJavaString(**s.p), true
}

type sliceSetter[T any] struct {
	p       *[]T
	handler OptionHandlerClass
}

func (s *sliceSetter[T]) IsMultiValued() bool                { return true }
func (s *sliceSetter[T]) DefaultHandler() OptionHandlerClass { return s.handler }
func (s *sliceSetter[T]) AddTypedValue(v T)                  { *s.p = append(*s.p, v) }
func (s *sliceSetter[T]) defaultValueString() (string, bool) {
	vals := *s.p
	switch len(vals) {
	case 0:
		return "", false
	case 1:
		return jutil.ToJavaString(vals[0]), true
	default:
		var buf strings.Builder
		buf.WriteByte('[')
		for i, v := range vals {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(jutil.ToJavaString(v))
		}
		buf.WriteByte(']')
		return buf.String(), true
	}
}

type intSetter struct {
	p *int
}

func (s *intSetter) IsMultiValued() bool                { return false }
func (s *intSetter) DefaultHandler() OptionHandlerClass { return NewIntOptionHandler }
func (s *intSetter) AddTypedValue(v int32)              { *s.p = int(v) }
func (s *intSetter) defaultValueString() (string, bool) { return strconv.Itoa(*s.p), true }

type stringMapSetter struct {
	p **jutil.OrderedMap[string, string]
}

func (s *stringMapSetter) IsMultiValued() bool                { return false }
func (s *stringMapSetter) DefaultHandler() OptionHandlerClass { return StringMapHandler }
func (s *stringMapSetter) defaultValueString() (string, bool) {
	if *s.p == nil {
		return "", false
	}
	return formatStringMapDefault(*s.p), true
}

func formatStringMapDefault(m *jutil.OrderedMap[string, string]) string {
	if m.IsEmpty() {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range m.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(k)
		sb.WriteByte('=')
		if v == "" {
			sb.WriteString("null")
		} else {
			sb.WriteString(v)
		}
	}
	sb.WriteByte('}')
	return sb.String()
}

// ScalarField binds a non-nullable scalar field *p using spec.
func ScalarField[T any](spec *TypeSpec[T], p *T) *Field {
	return &Field{setter: &scalarSetter[T]{p: p, handler: spec.defaultHandler}}
}

// NullableField binds a nullable pointer field **p using spec.
func NullableField[T any](spec *TypeSpec[T], p **T) *Field {
	return &Field{setter: &nullableSetter[T]{p: p, handler: spec.defaultHandler}}
}

// SliceField binds a multi-valued slice field *[]p using spec.
func SliceField[T any](spec *TypeSpec[T], p *[]T) *Field {
	return &Field{setter: &sliceSetter[T]{p: p, handler: spec.defaultHandler}}
}

type optStringSetter struct {
	p *base.OptString
}

func (s *optStringSetter) IsMultiValued() bool                { return false }
func (s *optStringSetter) DefaultHandler() OptionHandlerClass { return NewStringOptionHandler }
func (s *optStringSetter) AddTypedValue(v string)             { *s.p = base.SomeString(v) }
func (s *optStringSetter) defaultValueString() (string, bool) { return s.p.Get() }

func Bool(p *bool) *Field                { return ScalarField(BoolSpec, p) }
func NullableBool(p **bool) *Field       { return NullableField(BoolSpec, p) }
func Int(p *int) *Field                  { return &Field{setter: &intSetter{p: p}} }
func Int32(p *int32) *Field              { return ScalarField(Int32Spec, p) }
func NullableInt32(p **int32) *Field     { return NullableField(Int32Spec, p) }
func String(p *string) *Field            { return ScalarField(StringSpec, p) }
func OptString(p *base.OptString) *Field { return &Field{setter: &optStringSetter{p: p}} }
func StringSlice(p *[]string) *Field     { return SliceField(StringSpec, p) }
func StringMap(p **jutil.OrderedMap[string, string]) *Field {
	return &Field{setter: &stringMapSetter{p: p}}
}

// TrackedString binds *p and clears *isDefault whenever the option is set.
func TrackedString(p *string, isDefault *bool) *Field {
	return &Field{setter: &trackedScalarSetter[string]{p: p, isDefault: isDefault, handler: NewStringOptionHandler}}
}
