// Port of com/google/gson/Gson.java (Gson 2.9.1): fromJson, with
// TypeAdapters.JSON_ELEMENT of com/google/gson/internal/bind/TypeAdapters.java
// and Streams.write of com/google/gson/internal/Streams.java
// Port of com/google/gson/TypeAdapter.java (Gson 2.9.1)

/*
 * Copyright (C) 2008 Google Inc. (Gson.java)
 * Copyright (C) 2010 Google Inc. (Streams.java)
 * Copyright (C) 2011 Google Inc. (TypeAdapters.java, TypeAdapter.java)
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
