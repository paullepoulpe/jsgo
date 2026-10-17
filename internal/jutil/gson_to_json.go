// Port of com/google/gson/Gson.java (Gson 2.9.1): toJson, with the adapters
// of com/google/gson/internal/bind/TypeAdapters.java and
// com/google/gson/internal/bind/ReflectiveTypeAdapterFactory.java that it uses

/*
 * Copyright (C) 2008 Google Inc. (Gson.java)
 * Copyright (C) 2011 Google Inc. (TypeAdapters.java, ReflectiveTypeAdapterFactory.java)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jutil

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// Go-only: language-adapter — Gson 2.9.1's object serialization:
// GsonBuilder (serializeNulls, disableHtmlEscaping, setPrettyPrinting,
// registerTypeHierarchyAdapter), Gson.toJson(Object) and the TypeAdapter
// factories that toJson consults, applied to the Go forms of Java values
// (PORTING.md), as shaded into the pinned closure-compiler.jar.
//
// Adapters are chosen like Gson does for the runtime class of each value:
// JsonElement trees (TypeAdapters.JSON_ELEMENT), then the registered type
// hierarchy adapters (last registered first), then the built-in adapters:
//
//   - nil, typed nil pointers/slices and nil interfaces: null (every
//     built-in adapter is null-safe);
//   - string (String), bool (Boolean), int/int8/int16/int32/int64
//     (Integer, Long, Short, Byte), uint16 (Character), float64 (Double,
//     rejecting NaN and infinities like Gson.doubleAdapter);
//   - enums (an int32 kind with a String() method, PORTING.md §4): the
//     constant's name (TypeAdapters.ENUM_FACTORY);
//   - slices and arrays (arrays, ImmutableList) and the jutil collections
//     (Collection): JSON arrays; the jutil maps (Map): JSON objects named by
//     String.valueOf(key);
//   - structs and pointers to structs: ReflectiveTypeAdapterFactory — the
//     struct's own fields in declaration order, then the fields of its
//     embedded structs (the Java superclass), each named by the Go field
//     name or a `gson:"name"` tag (@SerializedName). Fields tagged
//     `gson:"-"` (transient) and the `self` back-pointer of the
//     inheritance pattern (PORTING.md §4) are not serialized; a field whose
//     value is the object itself is skipped, as in Gson.
//
// Other Go kinds (func, chan, native map, unsigned integers other than
// uint16, float32, complex) have no Java counterpart here and throw
// UnsupportedOperationException.

// ToJson serializes src to its JSON representation; nil is "null".
//
// Java: Gson.toJson(Object)
func (g *Gson) ToJson(src any) string {
	if IsNil(src) {
		return g.ToJsonElement(JsonNull_INSTANCE)
	}
	// Java: toJson(src, src.getClass()) with a StringWriter.
	var writer strings.Builder
	g.ToJsonWithWriter(src, g.NewJsonWriter(&writer))
	return writer.String()
}

// ToJsonWithWriter writes the JSON representation of src to writer, with
// the writer's lenient, htmlSafe and serializeNulls settings replaced by
// this Gson's for the duration of the call.
//
// Java: Gson.toJson(Object, Type, JsonWriter)
func (g *Gson) ToJsonWithWriter(src any, writer *JsonWriter) {
	oldLenient := writer.IsLenient()
	writer.SetLenient(true)
	oldHtmlSafe := writer.IsHtmlSafe()
	writer.SetHtmlSafe(g.htmlSafe)
	oldSerializeNulls := writer.GetSerializeNulls()
	writer.SetSerializeNulls(g.serializeNulls)
	defer func() {
		writer.SetLenient(oldLenient)
		writer.SetHtmlSafe(oldHtmlSafe)
		writer.SetSerializeNulls(oldSerializeNulls)
	}()
	gson_wrapWriteExceptions(func() { g.writeValue(writer, src) })
}

// ToJsonElement serializes a JSON tree.
//
// Java: Gson.toJson(JsonElement)
func (g *Gson) ToJsonElement(jsonElement JsonElement) string {
	var writer strings.Builder
	jsonWriter := g.NewJsonWriter(&writer)
	// Java: toJson(JsonElement, JsonWriter)
	oldLenient := jsonWriter.IsLenient()
	jsonWriter.SetLenient(true)
	oldHtmlSafe := jsonWriter.IsHtmlSafe()
	jsonWriter.SetHtmlSafe(g.htmlSafe)
	oldSerializeNulls := jsonWriter.GetSerializeNulls()
	jsonWriter.SetSerializeNulls(g.serializeNulls)
	defer func() {
		jsonWriter.SetLenient(oldLenient)
		jsonWriter.SetHtmlSafe(oldHtmlSafe)
		jsonWriter.SetSerializeNulls(oldSerializeNulls)
	}()
	gson_wrapWriteExceptions(func() { Streams_Write(jsonElement, jsonWriter) })
	return writer.String()
}

// gson_wrapWriteExceptions is the catch clauses of Gson.toJson: an
// IOException becomes a JsonIOException and an AssertionError is rethrown
// with the Gson version in its message.
func gson_wrapWriteExceptions(write func()) {
	Try(func() {
		Try(write, func(e IIOException) {
			panic(NewJsonIOExceptionWithCause(e))
		})
	}, func(e *AssertionError) {
		err := NewAssertionError("AssertionError (GSON " + GsonBuildConfig_VERSION + "): " + javaMessageOrNull(e))
		err.InitCause(e)
		panic(err)
	})
}

// NewJsonWriter returns a JSON writer configured for this Gson's settings.
//
// Java: Gson.newJsonWriter(Writer)
func (g *Gson) NewJsonWriter(writer io.StringWriter) *JsonWriter {
	jsonWriter := NewJsonWriter(writer)
	if g.prettyPrinting {
		jsonWriter.SetIndent("  ")
	}
	jsonWriter.SetHtmlSafe(g.htmlSafe)
	jsonWriter.SetLenient(g.lenient)
	jsonWriter.SetSerializeNulls(g.serializeNulls)
	return jsonWriter
}

// writeValue writes value with the adapter Gson.getAdapter picks for its
// runtime class (TypeAdapterRuntimeTypeWrapper does the same for fields and
// elements).
func (g *Gson) writeValue(out *JsonWriter, value any) {
	if IsNil(value) {
		out.NullValue()
		return
	}
	if e, ok := value.(JsonElement); ok { // TypeAdapters.JSON_ELEMENT_FACTORY
		TypeAdapters_JSON_ELEMENT_Write(out, e)
		return
	}
	for _, f := range g.factoriesToBeAdded {
		if f.baseType(value) {
			f.typeAdapter.Write(out, value)
			return
		}
	}
	switch v := value.(type) {
	case gson_collection: // CollectionTypeAdapterFactory
		out.BeginArray()
		for _, element := range v.gsonElements() {
			g.writeValue(out, element)
		}
		out.EndArray()
		return
	case gson_map: // MapTypeAdapterFactory, complexMapKeySerialization = false
		out.BeginObject()
		keys, values := v.gsonEntries()
		for i, key := range keys {
			out.Name(ToJavaString(key))
			g.writeValue(out, values[i])
		}
		out.EndObject()
		return
	}
	g.writeReflectValue(out, reflect.ValueOf(value))
}

func (g *Gson) writeReflectValue(out *JsonWriter, rv reflect.Value) {
	switch rv.Kind() {
	case reflect.String: // TypeAdapters.STRING
		s := rv.String()
		out.ValueString(&s)
	case reflect.Bool: // TypeAdapters.BOOLEAN
		out.ValueBoolean(rv.Bool())
	case reflect.Int32:
		if s, ok := rv.Interface().(fmt.Stringer); ok && rv.Type().PkgPath() != "" { // TypeAdapters.ENUM_FACTORY
			name := s.String()
			out.ValueString(&name)
			return
		}
		out.ValueNumber(gson_javaInteger(rv.Int())) // TypeAdapters.INTEGER
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64: // INTEGER, BYTE, SHORT, LONG
		out.ValueNumber(gson_javaInteger(rv.Int()))
	case reflect.Uint16: // TypeAdapters.CHARACTER
		s := StringFromChar(uint16(rv.Uint()))
		out.ValueString(&s)
	case reflect.Float64: // Gson.doubleAdapter(false)
		d := rv.Float()
		gson_checkValidFloatingPoint(d)
		out.ValueNumber(gson_javaDouble(d))
	case reflect.Slice, reflect.Array: // ArrayTypeAdapter / CollectionTypeAdapterFactory
		if rv.Kind() == reflect.Slice && rv.IsNil() {
			out.NullValue()
			return
		}
		out.BeginArray()
		for i := 0; i < rv.Len(); i++ {
			g.writeValue(out, gson_interface(rv.Index(i)))
		}
		out.EndArray()
	case reflect.Pointer:
		if rv.IsNil() {
			out.NullValue()
			return
		}
		if rv.Elem().Kind() == reflect.Struct {
			g.writeStruct(out, rv.Elem(), rv)
			return
		}
		g.writeValue(out, gson_interface(rv.Elem()))
	case reflect.Interface:
		if rv.IsNil() {
			out.NullValue()
			return
		}
		g.writeValue(out, gson_interface(rv.Elem()))
	case reflect.Struct:
		// A struct value (not reached through a pointer): copy it so that its
		// unexported fields can be read.
		c := reflect.New(rv.Type()).Elem()
		c.Set(rv)
		g.writeStruct(out, c, reflect.Value{})
	default:
		panic(NewUnsupportedOperationExceptionMsg("Gson cannot serialize Go type " + rv.Type().String()))
	}
}

// writeStruct is ReflectiveTypeAdapterFactory.Adapter.write. self is the
// pointer the struct was reached through (for the direct-recursion check),
// or the zero Value.
func (g *Gson) writeStruct(out *JsonWriter, sv reflect.Value, self reflect.Value) {
	fields := gson_boundFields(sv.Type())
	out.BeginObject()
	for _, f := range fields {
		fieldValue := gson_interface(sv.FieldByIndex(f.index))
		if self.IsValid() && !IsNil(fieldValue) {
			if fp := reflect.ValueOf(fieldValue); fp.Kind() == reflect.Pointer && fp.Pointer() == self.Pointer() {
				// avoid direct recursion
				continue
			}
		}
		out.Name(f.name)
		g.writeValue(out, fieldValue)
	}
	out.EndObject()
}

type gson_boundField struct {
	name  string
	index []int
}

// gson_boundFields is ReflectiveTypeAdapterFactory.getBoundFields: the
// struct's own fields, then those of its embedded structs (the Java
// superclass chain), with duplicate names rejected.
func gson_boundFields(t reflect.Type) []gson_boundField {
	var fields []gson_boundField
	seen := map[string]bool{}
	var collect func(t reflect.Type, prefix []int)
	collect = func(t reflect.Type, prefix []int) {
		var embedded []int
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if sf.Anonymous && sf.Type.Kind() == reflect.Struct {
				embedded = append(embedded, i)
				continue
			}
			tag := sf.Tag.Get("gson")
			if tag == "-" || sf.Name == "self" {
				continue
			}
			name := sf.Name
			if tag != "" {
				name = tag
			}
			if seen[name] {
				panic(NewIllegalArgumentExceptionMsg(t.String() + " declares multiple JSON fields named " + name))
			}
			seen[name] = true
			fields = append(fields, gson_boundField{name: name, index: append(append([]int{}, prefix...), i)})
		}
		for _, i := range embedded {
			collect(t.Field(i).Type, append(append([]int{}, prefix...), i))
		}
	}
	collect(t, nil)
	return fields
}

// gson_interface returns the value of v as an interface, reading
// unexported struct fields the way Gson's ReflectionHelper.makeAccessible
// lets it read private Java fields.
func gson_interface(v reflect.Value) any {
	if v.CanInterface() {
		return v.Interface()
	}
	if v.CanAddr() {
		return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface()
	}
	panic(NewUnsupportedOperationExceptionMsg("Gson cannot read unaddressable Go value of type " + v.Type().String()))
}

// gson_checkValidFloatingPoint is Gson.checkValidFloatingPoint(double).
func gson_checkValidFloatingPoint(value float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		panic(NewIllegalArgumentExceptionMsg(DoubleToString(value) +
			" is not a valid double value as per JSON specification. To override this" +
			" behavior, use GsonBuilder.serializeSpecialFloatingPointValues() method."))
	}
}

// gson_javaInteger is a boxed Integer, Long, Short or Byte as a Number.
type gson_javaInteger int64

func (n gson_javaInteger) IntValue() int32      { return int32(n) }
func (n gson_javaInteger) LongValue() int64     { return int64(n) }
func (n gson_javaInteger) DoubleValue() float64 { return float64(n) }
func (n gson_javaInteger) String() string       { return strconv.FormatInt(int64(n), 10) }

// gson_javaDouble is a boxed Double as a Number.
type gson_javaDouble float64

func (n gson_javaDouble) IntValue() int32      { return D2I(float64(n)) }
func (n gson_javaDouble) LongValue() int64     { return D2L(float64(n)) }
func (n gson_javaDouble) DoubleValue() float64 { return float64(n) }
func (n gson_javaDouble) String() string       { return DoubleToString(float64(n)) }

// gson_collection is implemented by the jutil types that model a
// java.util.Collection: their elements in iteration order.
type gson_collection interface {
	gsonElements() []any
}

// gson_map is implemented by the jutil types that model a java.util.Map:
// their entries in iteration order.
type gson_map interface {
	gsonEntries() (keys []any, values []any)
}

func (l *ArrayList[T]) gsonElements() []any {
	out := make([]any, 0, l.Size())
	for e := range l.All() {
		out = append(out, e)
	}
	return out
}

func (s *LinkedHashSet[T]) gsonElements() []any {
	out := make([]any, 0, s.Size())
	for e := range s.All() {
		out = append(out, e)
	}
	return out
}

func (s *ObjectLinkedHashSet[T]) gsonElements() []any {
	var out []any
	for e := range s.All() {
		out = append(out, e)
	}
	return out
}

func (s *EnumSet[E]) gsonElements() []any {
	var out []any
	for e := range s.All() {
		out = append(out, e)
	}
	return out
}

func (h *LinkedHashMap[K, V]) gsonEntries() (keys []any, values []any) {
	for k, v := range h.All() {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func (h *ObjectLinkedHashMap[K, V]) gsonEntries() (keys []any, values []any) {
	for k, v := range h.All() {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func (m *TreeMap[K, V]) gsonEntries() (keys []any, values []any) {
	for k, v := range m.All() {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
