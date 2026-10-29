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

// GsonBuildConfig_VERSION is GsonBuildConfig.VERSION of the shaded Gson.
const GsonBuildConfig_VERSION = "2.9.1"

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
	reader := NewStringReader(json)
	jsonReader := g.NewJsonReader(reader)
	object, err := g.fromJsonReaderErr(jsonReader, typeOfT)
	if err != nil {
		return nil, err
	}
	if err := gson_assertFullConsumptionErr(object, jsonReader); err != nil {
		return nil, err
	}
	return object, nil
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

func gson_assertFullConsumptionErr(obj JsonElement, reader *JsonReader) (parseErr IJsonParseException) {
	if obj == nil {
		return nil
	}
	Try(func() {
		if reader.Peek() != JsonToken_END_DOCUMENT {
			parseErr = NewJsonSyntaxException("JSON document was not fully consumed.")
		}
	}, func(t Throwable) {
		switch e := t.(type) {
		case IMalformedJsonException:
			parseErr = NewJsonSyntaxExceptionWithCause(e)
		case IIOException:
			parseErr = NewJsonIOExceptionWithCause(e)
		default:
			panic(t)
		}
	})
	return parseErr
}

// fromJsonReaderErr reads the next JSON value from reader as typeOfT. Returns
// (nil, nil) if the reader is at EOF.
func (g *Gson) fromJsonReaderErr(reader *JsonReader, typeOfT gsonRequestedType) (result JsonElement, parseErr IJsonParseException) {
	isEmpty := true
	oldLenient := reader.IsLenient()
	reader.SetLenient(true)
	defer reader.SetLenient(oldLenient)
	Try(func() {
		reader.Peek()
		isEmpty = false
		result, parseErr = gson_typeHierarchyAdapterReadErr(reader, typeOfT)
	}, func(t Throwable) {
		switch e := t.(type) {
		case IEOFException:
			if isEmpty {
				result = nil
				return
			}
			parseErr = NewJsonSyntaxExceptionWithCause(e)
		case IIllegalStateException:
			parseErr = NewJsonSyntaxExceptionWithCause(e)
		case IIOException:
			parseErr = NewJsonSyntaxExceptionWithCause(e)
		case *AssertionError:
			error_ := NewAssertionError("AssertionError (GSON " + GsonBuildConfig_VERSION + "): " + javaMessageOrNull(e))
			error_.InitCause(e)
			panic(error_)
		default:
			panic(t)
		}
	})
	return result, parseErr
}

func gson_typeHierarchyAdapterReadErr(in *JsonReader, requestedType gsonRequestedType) (JsonElement, IJsonParseException) {
	result := TypeAdapters_JSON_ELEMENT_Read(in)
	if result != nil && !requestedType.isInstance(result) {
		return nil, NewJsonSyntaxException("Expected a " + requestedType.getName() +
			" but was " + jsonElementClassName(result) + "; at path " + in.GetPreviousPath())
	}
	return result, nil
}

// typeAdapters_tryBeginNesting tries to begin reading a JSON array or JSON
// object, returning nil if the next element is neither of those.
func typeAdapters_tryBeginNesting(in *JsonReader, peeked JsonToken) JsonElement {
	switch peeked {
	case JsonToken_BEGIN_ARRAY:
		in.BeginArray()
		return NewJsonArray()
	case JsonToken_BEGIN_OBJECT:
		in.BeginObject()
		return NewJsonObject()
	default:
		return nil
	}
}

// typeAdapters_readTerminal reads a JsonElement which cannot have any
// nested elements.
func typeAdapters_readTerminal(in *JsonReader, peeked JsonToken) JsonElement {
	switch peeked {
	case JsonToken_STRING:
		return NewJsonPrimitiveString(in.NextString())
	case JsonToken_NUMBER:
		number := in.NextString()
		return NewJsonPrimitiveNumber(NewLazilyParsedNumber(number))
	case JsonToken_BOOLEAN:
		return NewJsonPrimitiveBoolean(in.NextBoolean())
	case JsonToken_NULL:
		in.NextNull()
		return JsonNull_INSTANCE
	default:
		// When read(JsonReader) is called with JsonReader in invalid state
		panic(NewIllegalStateExceptionMsg("Unexpected token: " + peeked.String()))
	}
}

// TypeAdapters_JSON_ELEMENT_Read is TypeAdapters.JSON_ELEMENT.read: reads a
// whole JSON value into a tree, iteratively.
//
// Java: TypeAdapters.JSON_ELEMENT.read(JsonReader)
func TypeAdapters_JSON_ELEMENT_Read(in *JsonReader) JsonElement {
	peeked := in.Peek()
	current := typeAdapters_tryBeginNesting(in, peeked)
	if current == nil {
		return typeAdapters_readTerminal(in, peeked)
	}

	var stack []JsonElement

	for {
		for in.HasNext() {
			var name string
			if _, ok := current.(*JsonObject); ok {
				name = in.NextName()
			}

			peeked = in.Peek()
			value := typeAdapters_tryBeginNesting(in, peeked)
			isNesting := value != nil
			if value == nil {
				value = typeAdapters_readTerminal(in, peeked)
			}

			if a, ok := current.(*JsonArray); ok {
				a.Add(value)
			} else {
				current.(*JsonObject).Add(name, value)
			}

			if isNesting {
				stack = append(stack, current)
				current = value
			}
		}

		if _, ok := current.(*JsonArray); ok {
			in.EndArray()
		} else {
			in.EndObject()
		}

		if len(stack) == 0 {
			return current
		}
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}

// TypeAdapters_JSON_ELEMENT_Write is TypeAdapters.JSON_ELEMENT.write.
//
// Java: TypeAdapters.JSON_ELEMENT.write(JsonWriter, JsonElement)
func TypeAdapters_JSON_ELEMENT_Write(out *JsonWriter, value JsonElement) {
	if value == nil {
		out.NullValue()
		return
	}
	switch v := value.(type) {
	case *JsonNull:
		out.NullValue()
	case *JsonPrimitive:
		if v.IsNumber() {
			out.ValueNumber(v.GetAsNumber())
		} else if v.IsBoolean() {
			out.ValueBoolean(v.GetAsBoolean())
		} else {
			s := v.GetAsString()
			out.ValueString(&s)
		}
	case *JsonArray:
		out.BeginArray()
		for e := range v.All() {
			TypeAdapters_JSON_ELEMENT_Write(out, e)
		}
		out.EndArray()
	case *JsonObject:
		out.BeginObject()
		for k, elem := range v.EntrySet().All() {
			out.Name(k)
			TypeAdapters_JSON_ELEMENT_Write(out, elem)
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
	TypeAdapters_JSON_ELEMENT_Write(writer, element)
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
	if IsNil(src) {
		return g.ToJsonElement(JsonNull_INSTANCE)
	}
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

func (g *Gson) writeValue(out *JsonWriter, value any) {
	if IsNil(value) {
		out.NullValue()
		return
	}
	if e, ok := value.(JsonElement); ok {
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
	case gson_collection:
		out.BeginArray()
		for _, element := range v.gsonElements() {
			g.writeValue(out, element)
		}
		out.EndArray()
		return
	case gson_map:
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

func gson_boundFields(t reflect.Type) []gson_boundField {
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
	gsonElements() []any
}

type gson_map interface {
	gsonEntries() (keys []any, values []any)
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
	serializeNulls     bool
	prettyPrinting     bool
	escapeHtmlChars    bool
}

// NewGsonBuilder is new GsonBuilder(), whose configuration is new Gson()'s.
//
// Java: GsonBuilder()
func NewGsonBuilder() *GsonBuilder {
	return &GsonBuilder{escapeHtmlChars: true}
}

// SerializeNulls makes Gson serialize null fields and map values.
//
// Java: GsonBuilder.serializeNulls()
func (b *GsonBuilder) SerializeNulls() *GsonBuilder {
	b.serializeNulls = true
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
		serializeNulls:     b.serializeNulls,
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

// EOFException is java.io.EOFException.
type EOFException struct{ IOException }

func (*EOFException) eofException() {}

// IEOFException matches EOFException.
type IEOFException interface {
	IIOException
	eofException()
}

// NewEOFException is new EOFException(String).
func NewEOFException(message string) *EOFException {
	e := &EOFException{}
	InitThrowable(e, "java.io.EOFException", msgPtr(message), nil)
	return e
}

// MalformedJsonException is com.google.gson.stream.MalformedJsonException
// (an IOException).
type MalformedJsonException struct{ IOException }

func (*MalformedJsonException) malformedJsonException() {}

// IMalformedJsonException matches MalformedJsonException.
type IMalformedJsonException interface {
	IIOException
	malformedJsonException()
}

// NewMalformedJsonException is new MalformedJsonException(String).
func NewMalformedJsonException(msg string) *MalformedJsonException {
	e := &MalformedJsonException{}
	InitThrowable(e, GsonPackage+".stream.MalformedJsonException", msgPtr(msg), nil)
	return e
}

// JsonParseException is com.google.gson.JsonParseException (a
// RuntimeException).
type JsonParseException struct{ RuntimeException }

func (*JsonParseException) jsonParseException() {}

// IJsonParseException matches JsonParseException and its subclasses.
type IJsonParseException interface {
	IRuntimeException
	jsonParseException()
}

// JsonSyntaxException is com.google.gson.JsonSyntaxException.
type JsonSyntaxException struct{ JsonParseException }

// NewJsonSyntaxException is new JsonSyntaxException(String).
func NewJsonSyntaxException(msg string) *JsonSyntaxException {
	e := &JsonSyntaxException{}
	InitThrowable(e, GsonPackage+".JsonSyntaxException", msgPtr(msg), nil)
	return e
}

// NewJsonSyntaxExceptionWithCause is new JsonSyntaxException(Throwable).
func NewJsonSyntaxExceptionWithCause(cause Throwable) *JsonSyntaxException {
	e := &JsonSyntaxException{}
	InitThrowable(e, GsonPackage+".JsonSyntaxException", nil, cause)
	return e
}

// JsonIOException is com.google.gson.JsonIOException.
type JsonIOException struct{ JsonParseException }

// NewJsonIOExceptionWithCause is new JsonIOException(Throwable).
func NewJsonIOExceptionWithCause(cause Throwable) *JsonIOException {
	e := &JsonIOException{}
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
