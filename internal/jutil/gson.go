// Port of com/google/gson/Gson.java (Gson 2.9.1): fromJson and toJson, with
// com/google/gson/internal/bind/TypeAdapters.java,
// com/google/gson/internal/bind/ReflectiveTypeAdapterFactory.java,
// and Streams.write of com/google/gson/internal/Streams.java
// Port of com/google/gson/GsonBuilder.java (Gson 2.9.1)
// Port of com/google/gson/TypeAdapter.java (Gson 2.9.1)
// Port of com/google/gson/JsonParseException.java (Gson 2.9.1), with
// com/google/gson/JsonSyntaxException.java, com/google/gson/JsonIOException.java
// and com/google/gson/stream/MalformedJsonException.java

/*
 * Copyright (C) 2008 Google Inc. (Gson.java, GsonBuilder.java, JsonParseException.java, JsonIOException.java)
 * Copyright (C) 2010 Google Inc. (Streams.java, JsonSyntaxException.java, MalformedJsonException.java)
 * Copyright (C) 2011 Google Inc. (TypeAdapters.java, TypeAdapter.java, ReflectiveTypeAdapterFactory.java)
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
	"sync"
)

// Go-only: language-adapter — the parts of com.google.gson.Gson (2.9.1, as
// shaded into the pinned closure-compiler.jar) that Closure uses on JSON
// trees: Gson.fromJson(String, JsonElement/JsonObject/JsonArray.class) and
// the JsonElement type adapter (TypeAdapters.JSON_ELEMENT, wrapped by
// TypeAdapters.newTypeHierarchyFactory) with Streams.write.
//
// Gson.fromJson reads with a lenient JsonReader, so comments, unquoted and
// single-quoted strings and names, ';' separators, "=" / "=>" name
// separators, the ")]}'\n" prefix and empty array elements (null) are all
// accepted, exactly as in Java. Only the trailing-data check
// (assertFullConsumption) runs in strict mode.

// gsonBuildConfig_VERSION is GsonBuildConfig.VERSION of the shaded Gson.
const gsonBuildConfig_VERSION = "2.9.1"

// Gson is com.google.gson.Gson with its default configuration, or the
// configuration of a GsonBuilder (gson_to_json.go).
type Gson struct {
	lenient bool

	// The GsonBuilder settings that Gson.toJson reads (gson_to_json.go).
	serializeNulls     bool
	htmlSafe           bool
	prettyPrinting     bool
	factoriesToBeAdded []*gson_typeHierarchyFactory
}

// NewGson is new Gson().
//
// Java: Gson()
func NewGson() *Gson {
	return &Gson{lenient: false, htmlSafe: true}
}

// gsonRequestedType identifies the Class passed to fromJson.
type gsonRequestedType int

const (
	gsonRequestedJsonElement gsonRequestedType = iota
	gsonRequestedJsonObject
	gsonRequestedJsonArray
)

func (t gsonRequestedType) getName() string {
	switch t {
	case gsonRequestedJsonObject:
		return GsonPackage + ".JsonObject"
	case gsonRequestedJsonArray:
		return GsonPackage + ".JsonArray"
	}
	return GsonPackage + ".JsonElement"
}

func (t gsonRequestedType) isInstance(e JsonElement) bool {
	switch t {
	case gsonRequestedJsonObject:
		_, ok := e.(*JsonObject)
		return ok
	case gsonRequestedJsonArray:
		_, ok := e.(*JsonArray)
		return ok
	}
	return true
}

// FromJsonElement is gson.fromJson(json, JsonElement.class); nil is Java
// null (an empty document).
//
// Java: Gson.fromJson(String, Class)
func (g *Gson) FromJsonElement(json string) JsonElement {
	e, err := g.TryFromJsonElement(json)
	if err != nil {
		panic(err)
	}
	return e
}

// TryFromJsonElement parses json as JsonElement, returning a JsonParseException
// instead of panicking on malformed JSON.
func (g *Gson) TryFromJsonElement(json string) (JsonElement, IJsonParseException) {
	return g.fromJsonStringErr(json, gsonRequestedJsonElement)
}

// FromJsonObject is gson.fromJson(json, JsonObject.class); nil is Java null
// (an empty document).
//
// Java: Gson.fromJson(String, Class)
func (g *Gson) FromJsonObject(json string) *JsonObject {
	o, err := g.TryFromJsonObject(json)
	if err != nil {
		panic(err)
	}
	return o
}

// TryFromJsonObject parses json as JsonObject, returning a JsonParseException
// instead of panicking on malformed JSON.
func (g *Gson) TryFromJsonObject(json string) (*JsonObject, IJsonParseException) {
	e, err := g.fromJsonStringErr(json, gsonRequestedJsonObject)
	if err != nil || e == nil {
		return nil, err
	}
	return e.(*JsonObject), nil
}

// FromJsonArray is gson.fromJson(json, JsonArray.class); nil is Java null
// (an empty document).
//
// Java: Gson.fromJson(String, Class)
func (g *Gson) FromJsonArray(json string) *JsonArray {
	a, err := g.TryFromJsonArray(json)
	if err != nil {
		panic(err)
	}
	return a
}

// TryFromJsonArray parses json as JsonArray, returning a JsonParseException
// instead of panicking on malformed JSON.
func (g *Gson) TryFromJsonArray(json string) (*JsonArray, IJsonParseException) {
	e, err := g.fromJsonStringErr(json, gsonRequestedJsonArray)
	if err != nil || e == nil {
		return nil, err
	}
	return e.(*JsonArray), nil
}

// Java: Gson.fromJson(String, Type) and Gson.fromJson(Reader, Type)
func (g *Gson) fromJsonStringErr(json string, typeOfT gsonRequestedType) (JsonElement, IJsonParseException) {
	var s JsonScanner
	s.Init(json, true)
	_, emptyEOF, err := s.DoPeek(true)
	if err != nil {
		return nil, gson_wrapScannerErr(err)
	}
	if emptyEOF {
		return nil, nil
	}
	object, err := s.ReadJsonElement()
	if err != nil {
		return nil, gson_wrapScannerErr(err)
	}
	if object != nil && !typeOfT.isInstance(object) {
		return nil, newJsonSyntaxException("Expected a " + typeOfT.getName() +
			" but was " + jsonElementClassName(object) + "; at path " + s.GetPreviousPath())
	}
	if object == nil {
		return nil, nil
	}
	s.SetLenient(g.lenient)
	endTok, err := s.Peek()
	if err != nil {
		switch e := err.(type) {
		case *malformedJsonException:
			return nil, NewJsonSyntaxExceptionWithCause(e)
		case IIOException:
			return nil, newJsonIOExceptionWithCause(e)
		default:
			panic(err)
		}
	}
	if endTok != JsonPeekEOF {
		return nil, newJsonSyntaxException("JSON document was not fully consumed.")
	}
	return object, nil
}

func gson_wrapScannerErr(err Throwable) IJsonParseException {
	switch e := err.(type) {
	case IEOFException, IIllegalStateException, IIOException:
		return NewJsonSyntaxExceptionWithCause(e)
	case *AssertionError:
		wrapped := NewAssertionError("AssertionError (GSON " + gsonBuildConfig_VERSION + "): " + javaMessageOrNull(e))
		wrapped.InitCause(e)
		panic(wrapped)
	default:
		panic(err)
	}
}

// NewJsonReader returns a new JSON reader configured for the settings on
// this Gson instance.
//
// Java: Gson.newJsonReader(Reader)
func (g *Gson) NewJsonReader(reader JavaReader) *JsonReader {
	jsonReader := NewJsonReader(reader)
	jsonReader.SetLenient(g.lenient)
	return jsonReader
}

// ReadJsonElement reads one Gson JsonElement tree iteratively (TypeAdapters.JSON_ELEMENT.read).
func (s *JsonScanner) ReadJsonElement() (JsonElement, Throwable) {
	p, err := s.Peek()
	if err != nil {
		return nil, err
	}
	readScalar := func(tok uint8) (JsonElement, Throwable) {
		switch tok {
		case JsonPeekDoubleQuoted, JsonPeekSingleQuoted, JsonPeekUnquoted:
			str, err := s.NextString()
			if err != nil {
				return nil, err
			}
			return NewJsonPrimitiveString(str), nil
		case JsonPeekLong, JsonPeekNumber:
			str, err := s.NextString()
			if err != nil {
				return nil, err
			}
			return NewJsonPrimitiveNumber(NewLazilyParsedNumber(str)), nil
		case JsonPeekTrue:
			s.ConsumeScalar()
			return NewJsonPrimitiveBoolean(true), nil
		case JsonPeekFalse:
			s.ConsumeScalar()
			return NewJsonPrimitiveBoolean(false), nil
		case JsonPeekNull:
			s.ConsumeScalar()
			return JsonNull_INSTANCE, nil
		default:
			jt, err := s.PeekToken()
			if err != nil {
				return nil, err
			}
			return nil, NewIllegalStateExceptionMsg("Unexpected token: " + jt.String())
		}
	}

	var current JsonElement
	switch p {
	case JsonPeekBeginArray:
		_ = s.BeginArray()
		current = NewJsonArray()
	case JsonPeekBeginObject:
		_ = s.BeginObject()
		current = NewJsonObject()
	default:
		return readScalar(p)
	}

	var stack []JsonElement
	for {
		for {
			more, err := s.HasNext()
			if err != nil {
				return nil, err
			}
			if !more {
				break
			}
			var name string
			_, isObj := current.(*JsonObject)
			if isObj {
				name, err = s.NextName()
				if err != nil {
					return nil, err
				}
			}
			vp, err := s.Peek()
			if err != nil {
				return nil, err
			}
			var value JsonElement
			isNesting := false
			switch vp {
			case JsonPeekBeginArray:
				_ = s.BeginArray()
				value = NewJsonArray()
				isNesting = true
			case JsonPeekBeginObject:
				_ = s.BeginObject()
				value = NewJsonObject()
				isNesting = true
			default:
				value, err = readScalar(vp)
				if err != nil {
					return nil, err
				}
			}
			if isObj {
				current.(*JsonObject).Add(name, value)
			} else {
				current.(*JsonArray).Add(value)
			}
			if isNesting {
				stack = append(stack, current)
				current = value
			}
		}
		if _, isArr := current.(*JsonArray); isArr {
			if err := s.EndArray(); err != nil {
				return nil, err
			}
		} else {
			if err := s.EndObject(); err != nil {
				return nil, err
			}
		}
		if len(stack) == 0 {
			return current, nil
		}
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}

// typeAdapters_JSON_ELEMENT_Write is TypeAdapters.JSON_ELEMENT.write.
//
// Java: TypeAdapters.JSON_ELEMENT.write(JsonWriter, JsonElement)
func typeAdapters_JSON_ELEMENT_Write(out *JsonWriter, value JsonElement) {
	if value == nil {
		out.NullValue()
		return
	}
	switch v := value.(type) {
	case *JsonNull:
		out.NullValue()
	case *JsonPrimitive:
		if v.IsNumber() {
			out.ValueNumber(v.getAsNumber())
		} else if v.IsBoolean() {
			out.ValueBoolean(v.GetAsBoolean())
		} else {
			s := v.GetAsString()
			out.ValueString(&s)
		}
	case *JsonArray:
		out.BeginArray()
		for e := range v.All() {
			typeAdapters_JSON_ELEMENT_Write(out, e)
		}
		out.EndArray()
	case *JsonObject:
		out.BeginObject()
		for k, elem := range v.EntrySet().All() {
			out.Name(k)
			typeAdapters_JSON_ELEMENT_Write(out, elem)
		}
		out.EndObject()
	default:
		panic(NewIllegalArgumentExceptionMsg("Couldn't write class " + jsonElementClassName(value)))
	}
}

// Streams_Write writes the JSON element to the writer, recursively.
//
// Java: Streams.write(JsonElement, JsonWriter)
func Streams_Write(element JsonElement, writer *JsonWriter) {
	typeAdapters_JSON_ELEMENT_Write(writer, element)
}

// GsonWritable is implemented by Go types that stream their own JSON
// representation directly to JsonWriter without reflection.
type GsonWritable interface {
	WriteGson(out *JsonWriter)
}

// TypeAdapter is com.google.gson.TypeAdapter: converts Java objects to and
// from JSON.
type TypeAdapter interface {
	// Write writes one JSON value (an array, object, string, number,
	// boolean or null) for value.
	//
	// Java: TypeAdapter.write(JsonWriter, T)
	Write(out *JsonWriter, value any)

	// Read reads one JSON value and converts it to a Java object.
	//
	// Java: TypeAdapter.read(JsonReader)
	Read(in *JsonReader) any
}

// ToJson serializes src to its JSON representation; nil is "null".
//
// Java: Gson.toJson(Object)
func (g *Gson) ToJson(src any) string {
	if src == nil {
		return g.ToJsonElement(JsonNull_INSTANCE)
	}
	var writer strings.Builder
	g.ToJsonWriter(src, g.NewJsonWriter(&writer))
	return writer.String()
}

// ToJsonWriter writes the JSON representation of src to writer, with
// the writer's lenient, htmlSafe and serializeNulls settings replaced by
// this Gson's for the duration of the call.
//
// Java: Gson.toJson(Object, Type, JsonWriter)
func (g *Gson) ToJsonWriter(src any, writer *JsonWriter) {
	oldLenient := writer.IsLenient()
	writer.SetLenient(true)
	oldHtmlSafe := writer.isHtmlSafe()
	writer.SetHtmlSafe(g.htmlSafe)
	oldSerializeNulls := writer.getSerializeNulls()
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
	oldLenient := jsonWriter.IsLenient()
	jsonWriter.SetLenient(true)
	oldHtmlSafe := jsonWriter.isHtmlSafe()
	jsonWriter.SetHtmlSafe(g.htmlSafe)
	oldSerializeNulls := jsonWriter.getSerializeNulls()
	jsonWriter.SetSerializeNulls(g.serializeNulls)
	defer func() {
		jsonWriter.SetLenient(oldLenient)
		jsonWriter.SetHtmlSafe(oldHtmlSafe)
		jsonWriter.SetSerializeNulls(oldSerializeNulls)
	}()
	gson_wrapWriteExceptions(func() { Streams_Write(jsonElement, jsonWriter) })
	return writer.String()
}

func gson_wrapWriteExceptions(write func()) {
	Try(func() {
		Try(write, func(e IIOException) {
			panic(newJsonIOExceptionWithCause(e))
		})
	}, func(e *AssertionError) {
		err := NewAssertionError("AssertionError (GSON " + gsonBuildConfig_VERSION + "): " + javaMessageOrNull(e))
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
		jsonWriter.setIndent("  ")
	}
	jsonWriter.SetHtmlSafe(g.htmlSafe)
	jsonWriter.SetLenient(g.lenient)
	jsonWriter.SetSerializeNulls(g.serializeNulls)
	return jsonWriter
}

func (g *Gson) writeValue(out *JsonWriter, value any) {
	if value == nil {
		out.NullValue()
		return
	}
	if len(g.factoriesToBeAdded) == 0 {
		if w, ok := value.(GsonWritable); ok {
			w.WriteGson(out)
			return
		}
	}
	if e, ok := value.(JsonElement); ok {
		if IsNil(e) {
			out.NullValue()
			return
		}
		typeAdapters_JSON_ELEMENT_Write(out, e)
		return
	}
	if len(g.factoriesToBeAdded) > 0 {
		if IsNil(value) {
			out.NullValue()
			return
		}
		for _, f := range g.factoriesToBeAdded {
			if f.baseType(value) {
				f.typeAdapter.Write(out, value)
				return
			}
		}
		if w, ok := value.(GsonWritable); ok {
			w.WriteGson(out)
			return
		}
	}
	switch v := value.(type) {
	case string:
		out.ValueString(&v)
		return
	case *string:
		out.ValueString(v)
		return
	case bool:
		out.ValueBoolean(v)
		return
	case int:
		out.ValueNumber(gson_javaInteger(v))
		return
	case int32:
		out.ValueNumber(gson_javaInteger(v))
		return
	case int64:
		out.ValueNumber(gson_javaInteger(v))
		return
	case int16:
		out.ValueNumber(gson_javaInteger(v))
		return
	case int8:
		out.ValueNumber(gson_javaInteger(v))
		return
	case uint16:
		s := StringFromChar(v)
		out.ValueString(&s)
		return
	case float64:
		gson_checkValidFloatingPoint(v)
		out.ValueNumber(gson_javaDouble(v))
		return
	case []any:
		if v == nil {
			out.NullValue()
			return
		}
		out.BeginArray()
		for _, element := range v {
			g.writeValue(out, element)
		}
		out.EndArray()
		return
	case []string:
		if v == nil {
			out.NullValue()
			return
		}
		out.BeginArray()
		for i := range v {
			if len(g.factoriesToBeAdded) == 0 {
				out.ValueString(&v[i])
			} else {
				g.writeValue(out, v[i])
			}
		}
		out.EndArray()
		return
	case []*string:
		if v == nil {
			out.NullValue()
			return
		}
		out.BeginArray()
		for _, element := range v {
			if element == nil {
				out.NullValue()
			} else if len(g.factoriesToBeAdded) == 0 {
				out.ValueString(element)
			} else {
				g.writeValue(out, *element)
			}
		}
		out.EndArray()
		return
	case []int32:
		if v == nil {
			out.NullValue()
			return
		}
		out.BeginArray()
		for _, element := range v {
			if len(g.factoriesToBeAdded) == 0 {
				out.ValueNumber(gson_javaInteger(element))
			} else {
				g.writeValue(out, element)
			}
		}
		out.EndArray()
		return
	case gson_collection:
		v.writeGsonCollection(g, out)
		return
	case gson_map:
		v.writeGsonMap(g, out)
		return
	}
	if IsNil(value) {
		out.NullValue()
		return
	}
	g.writeReflectValue(out, reflect.ValueOf(value))
}

func (g *Gson) writeReflectValue(out *JsonWriter, rv reflect.Value) {
	switch rv.Kind() {
	case reflect.String:
		s := rv.String()
		out.ValueString(&s)
	case reflect.Bool:
		out.ValueBoolean(rv.Bool())
	case reflect.Int32:
		if s, ok := rv.Interface().(fmt.Stringer); ok && rv.Type().PkgPath() != "" {
			name := s.String()
			out.ValueString(&name)
			return
		}
		out.ValueNumber(gson_javaInteger(rv.Int()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
		out.ValueNumber(gson_javaInteger(rv.Int()))
	case reflect.Uint16:
		s := StringFromChar(uint16(rv.Uint()))
		out.ValueString(&s)
	case reflect.Float64:
		d := rv.Float()
		gson_checkValidFloatingPoint(d)
		out.ValueNumber(gson_javaDouble(d))
	case reflect.Slice, reflect.Array:
		if rv.Kind() == reflect.Slice && rv.IsNil() {
			out.NullValue()
			return
		}
		out.BeginArray()
		for i := 0; i < rv.Len(); i++ {
			g.writeValue(out, rv.Index(i).Interface())
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
		g.writeValue(out, rv.Elem().Interface())
	case reflect.Interface:
		if rv.IsNil() {
			out.NullValue()
			return
		}
		g.writeValue(out, rv.Elem().Interface())
	case reflect.Struct:
		g.writeStruct(out, rv, reflect.Value{})
	default:
		panic(NewUnsupportedOperationExceptionMsg("Gson cannot serialize Go type " + rv.Type().String()))
	}
}

func (g *Gson) writeStruct(out *JsonWriter, sv reflect.Value, self reflect.Value) {
	fields := gson_boundFields(sv.Type())
	out.BeginObject()
	for _, f := range fields {
		fieldValue := sv.FieldByIndex(f.index).Interface()
		if self.IsValid() && !IsNil(fieldValue) {
			if fp := reflect.ValueOf(fieldValue); fp.Kind() == reflect.Pointer && fp.Pointer() == self.Pointer() {
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

var gsonBoundFieldsCache sync.Map // map[reflect.Type][]gson_boundField

func gson_boundFields(t reflect.Type) []gson_boundField {
	if cached, ok := gsonBoundFieldsCache.Load(t); ok {
		return cached.([]gson_boundField)
	}
	var fields []gson_boundField
	seen := map[string]bool{}
	var collect func(t reflect.Type, prefix []int)
	collect = func(t reflect.Type, prefix []int) {
		var embedded []int
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if sf.Anonymous && sf.Type.Kind() == reflect.Struct {
				if sf.IsExported() {
					embedded = append(embedded, i)
				}
				continue
			}
			tag := sf.Tag.Get("gson")
			if tag == "-" || !sf.IsExported() {
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
	gsonBoundFieldsCache.Store(t, fields)
	return fields
}

func gson_checkValidFloatingPoint(value float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		panic(NewIllegalArgumentExceptionMsg(DoubleToString(value) +
			" is not a valid double value as per JSON specification. To override this" +
			" behavior, use GsonBuilder.serializeSpecialFloatingPointValues() method."))
	}
}

type gson_javaInteger int64

func (n gson_javaInteger) IntValue() int32      { return int32(n) }
func (n gson_javaInteger) LongValue() int64     { return int64(n) }
func (n gson_javaInteger) DoubleValue() float64 { return float64(n) }
func (n gson_javaInteger) String() string       { return strconv.FormatInt(int64(n), 10) }

type gson_javaDouble float64

func (n gson_javaDouble) IntValue() int32      { return D2I(float64(n)) }
func (n gson_javaDouble) LongValue() int64     { return D2L(float64(n)) }
func (n gson_javaDouble) DoubleValue() float64 { return float64(n) }
func (n gson_javaDouble) String() string       { return DoubleToString(float64(n)) }

type gson_collection interface {
	writeGsonCollection(g *Gson, out *JsonWriter)
}

type gson_map interface {
	writeGsonMap(g *Gson, out *JsonWriter)
}

func (s *LinkedHashSet[T]) writeGsonCollection(g *Gson, out *JsonWriter) {
	if s == nil {
		out.NullValue()
		return
	}
	out.BeginArray()
	for e := range s.All() {
		g.writeValue(out, e)
	}
	out.EndArray()
}

func (s *ObjectLinkedHashSet[T]) writeGsonCollection(g *Gson, out *JsonWriter) {
	if s == nil {
		out.NullValue()
		return
	}
	out.BeginArray()
	for e := range s.All() {
		g.writeValue(out, e)
	}
	out.EndArray()
}

func (t *orderedTable[K, V, M]) writeGsonMap(g *Gson, out *JsonWriter) {
	if t == nil {
		out.NullValue()
		return
	}
	out.BeginObject()
	for k, v := range t.all() {
		out.Name(ToJavaString(k))
		g.writeValue(out, v)
	}
	out.EndObject()
}

func (h *LinkedHashMap[K, V]) writeGsonMap(g *Gson, out *JsonWriter) {
	h.table().writeGsonMap(g, out)
}

func (h *ObjectLinkedHashMap[K, V]) writeGsonMap(g *Gson, out *JsonWriter) {
	h.table().writeGsonMap(g, out)
}

// Go-only: language-adapter — com.google.gson.GsonBuilder (Gson 2.9.1, as
// shaded into the pinned closure-compiler.jar): the settings that
// Gson.toJson uses (gson_to_json.go).

// gson_typeHierarchyFactory is the factory that
// TypeAdapters.newTypeHierarchyFactory(Class, TypeAdapter) returns: it
// applies typeAdapter to every value that is an instance of baseType.
type gson_typeHierarchyFactory struct {
	baseType    func(value any) bool
	typeAdapter TypeAdapter
}

// GsonBuilder is com.google.gson.GsonBuilder (the settings toJson uses).
type GsonBuilder struct {
	hierarchyFactories []*gson_typeHierarchyFactory
	serializeNullsFlag bool
	prettyPrinting     bool
	escapeHtmlChars    bool
}

// NewGsonBuilder is new GsonBuilder(), whose configuration is new Gson()'s.
//
// Java: GsonBuilder()
func NewGsonBuilder() *GsonBuilder {
	return &GsonBuilder{escapeHtmlChars: true}
}

// serializeNulls makes Gson serialize null fields and map values.
//
// Java: GsonBuilder.serializeNulls()
func (b *GsonBuilder) serializeNulls() *GsonBuilder {
	b.serializeNullsFlag = true
	return b
}

// DisableHtmlEscaping makes Gson pass <, >, &, = and ' through as is.
//
// Java: GsonBuilder.disableHtmlEscaping()
func (b *GsonBuilder) DisableHtmlEscaping() *GsonBuilder {
	b.escapeHtmlChars = false
	return b
}

// SetPrettyPrinting makes Gson indent its output by two spaces per level.
//
// Java: GsonBuilder.setPrettyPrinting()
func (b *GsonBuilder) SetPrettyPrinting() *GsonBuilder {
	b.prettyPrinting = true
	return b
}

// RegisterTypeHierarchyAdapter registers typeAdapter for baseType and all
// its subtypes. baseType stands for the Class<?> argument: it reports
// whether a value is an instance of the class (`value instanceof Base`).
//
// Java: GsonBuilder.registerTypeHierarchyAdapter(Class<?>, Object)
func (b *GsonBuilder) RegisterTypeHierarchyAdapter(baseType func(value any) bool, typeAdapter TypeAdapter) *GsonBuilder {
	CheckArgument(typeAdapter != nil)
	b.hierarchyFactories = append(b.hierarchyFactories,
		&gson_typeHierarchyFactory{baseType: baseType, typeAdapter: typeAdapter})
	return b
}

// Create creates a Gson instance with the builder's configuration. Type
// hierarchy adapters registered later take precedence.
//
// Java: GsonBuilder.create()
func (b *GsonBuilder) Create() *Gson {
	factories := make([]*gson_typeHierarchyFactory, 0, len(b.hierarchyFactories))
	for i := len(b.hierarchyFactories) - 1; i >= 0; i-- { // Collections.reverse
		factories = append(factories, b.hierarchyFactories[i])
	}
	return &Gson{
		lenient:            false,
		serializeNulls:     b.serializeNullsFlag,
		htmlSafe:           b.escapeHtmlChars,
		prettyPrinting:     b.prettyPrinting,
		factoriesToBeAdded: factories,
	}
}

// ---------------------------------------------------------------------------
// Gson exceptions

// GsonPackage is the (shaded) Java package of Gson's classes in the pinned
// closure-compiler.jar.
const GsonPackage = "com.google.javascript.jscomp.jarjar.com.google.gson"

// eofException is java.io.EOFException.
type eofException struct{ IOException }

func (*eofException) isEOFException() {}

// IEOFException matches EOFException.
type IEOFException interface {
	IIOException
	isEOFException()
}

// newEOFException is new EOFException(String).
func newEOFException(message string) *eofException {
	e := &eofException{}
	InitThrowable(e, "java.io.EOFException", msgPtr(message), nil)
	return e
}

// malformedJsonException is com.google.gson.stream.MalformedJsonException
// (an IOException).
type malformedJsonException struct{ IOException }

// newMalformedJsonException is new MalformedJsonException(String).
func newMalformedJsonException(msg string) *malformedJsonException {
	e := &malformedJsonException{}
	InitThrowable(e, GsonPackage+".stream.MalformedJsonException", msgPtr(msg), nil)
	return e
}

// jsonParseException is com.google.gson.JsonParseException (a
// RuntimeException).
type jsonParseException struct{ RuntimeException }

func (*jsonParseException) isJsonParseException() {}

// IJsonParseException matches JsonParseException and its subclasses.
type IJsonParseException interface {
	IRuntimeException
	isJsonParseException()
}

// jsonSyntaxException is com.google.gson.JsonSyntaxException.
type jsonSyntaxException struct{ jsonParseException }

// newJsonSyntaxException is new JsonSyntaxException(String).
func newJsonSyntaxException(msg string) *jsonSyntaxException {
	e := &jsonSyntaxException{}
	InitThrowable(e, GsonPackage+".JsonSyntaxException", msgPtr(msg), nil)
	return e
}

// NewJsonSyntaxExceptionWithCause is new JsonSyntaxException(Throwable).
func NewJsonSyntaxExceptionWithCause(cause Throwable) IJsonParseException {
	e := &jsonSyntaxException{}
	InitThrowable(e, GsonPackage+".JsonSyntaxException", nil, cause)
	return e
}

// jsonIOException is com.google.gson.JsonIOException.
type jsonIOException struct{ jsonParseException }

// newJsonIOExceptionWithCause is new JsonIOException(Throwable).
func newJsonIOExceptionWithCause(cause Throwable) *jsonIOException {
	e := &jsonIOException{}
	InitThrowable(e, GsonPackage+".JsonIOException", nil, cause)
	return e
}

// javaMessageOrNull is "" + t.getMessage(): the message, or "null" when
// the Java message is null.
func javaMessageOrNull(t Throwable) string {
	b := t.throwable()
	if !b.hasMsg {
		return "null"
	}
	return b.message
}
