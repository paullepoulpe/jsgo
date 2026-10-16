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

// Gson is com.google.gson.Gson with its default configuration.
type Gson struct {
	lenient bool
}

// NewGson is new Gson().
//
// Java: Gson()
func NewGson() *Gson {
	return &Gson{lenient: false}
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
	return g.fromJsonString(json, gsonRequestedJsonElement)
}

// FromJsonObject is gson.fromJson(json, JsonObject.class); nil is Java null
// (an empty document).
//
// Java: Gson.fromJson(String, Class)
func (g *Gson) FromJsonObject(json string) *JsonObject {
	e := g.fromJsonString(json, gsonRequestedJsonObject)
	if e == nil {
		return nil
	}
	return e.(*JsonObject)
}

// FromJsonArray is gson.fromJson(json, JsonArray.class); nil is Java null
// (an empty document).
//
// Java: Gson.fromJson(String, Class)
func (g *Gson) FromJsonArray(json string) *JsonArray {
	e := g.fromJsonString(json, gsonRequestedJsonArray)
	if e == nil {
		return nil
	}
	return e.(*JsonArray)
}

// Java: Gson.fromJson(String, Type) and Gson.fromJson(Reader, Type)
func (g *Gson) fromJsonString(json string, typeOfT gsonRequestedType) JsonElement {
	reader := NewStringReader(json)
	jsonReader := g.NewJsonReader(reader)
	object := g.fromJsonReader(jsonReader, typeOfT)
	gson_assertFullConsumption(object, jsonReader)
	return object
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

func gson_assertFullConsumption(obj JsonElement, reader *JsonReader) {
	Try(func() {
		if obj != nil && reader.Peek() != JsonToken_END_DOCUMENT {
			panic(NewJsonSyntaxException("JSON document was not fully consumed."))
		}
	}, func(t Throwable) {
		switch e := t.(type) {
		case IMalformedJsonException:
			panic(NewJsonSyntaxExceptionWithCause(e))
		case IIOException:
			panic(NewJsonIOExceptionWithCause(e))
		default:
			panic(t)
		}
	})
}

// fromJsonReader reads the next JSON value from reader as typeOfT. Returns
// nil (Java null) if the reader is at EOF.
//
// Java: Gson.fromJson(JsonReader, Type)
func (g *Gson) fromJsonReader(reader *JsonReader, typeOfT gsonRequestedType) (result JsonElement) {
	isEmpty := true
	oldLenient := reader.IsLenient()
	reader.SetLenient(true)
	func() {
		defer reader.SetLenient(oldLenient)
		Try(func() {
			reader.Peek()
			isEmpty = false
			result = gson_typeHierarchyAdapterRead(reader, typeOfT)
		}, func(t Throwable) {
			switch e := t.(type) {
			case IEOFException:
				// For compatibility with JSON 1.5 and earlier, we return null
				// for empty documents instead of throwing.
				if isEmpty {
					result = nil
					return
				}
				panic(NewJsonSyntaxExceptionWithCause(e))
			case IIllegalStateException:
				panic(NewJsonSyntaxExceptionWithCause(e))
			case IIOException:
				// TODO(inder): Figure out whether it is indeed right to rethrow
				// this as JsonSyntaxException
				panic(NewJsonSyntaxExceptionWithCause(e))
			case *AssertionError:
				error_ := NewAssertionError("AssertionError (GSON " + GsonBuildConfig_VERSION + "): " + javaMessageOrNull(e))
				error_.InitCause(e)
				panic(error_)
			default:
				panic(t)
			}
		})
	}()
	return result
}

// gson_typeHierarchyAdapterRead is the read method of the adapter that
// TypeAdapters.newTypeHierarchyFactory(JsonElement.class, JSON_ELEMENT)
// creates for the requested type.
func gson_typeHierarchyAdapterRead(in *JsonReader, requestedType gsonRequestedType) JsonElement {
	result := TypeAdapters_JSON_ELEMENT_Read(in)
	if result != nil && !requestedType.isInstance(result) {
		panic(NewJsonSyntaxException("Expected a " + requestedType.getName() +
			" but was " + jsonElementClassName(result) + "; at path " + in.GetPreviousPath()))
	}
	return result
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
	// Either JsonArray or JsonObject
	var current JsonElement
	peeked := in.Peek()

	current = typeAdapters_tryBeginNesting(in, peeked)
	if current == nil {
		return typeAdapters_readTerminal(in, peeked)
	}

	var stack []JsonElement

	for {
		for in.HasNext() {
			var name string
			// Name is only used for JSON object members
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

		// End current element
		if _, ok := current.(*JsonArray); ok {
			in.EndArray()
		} else {
			in.EndObject()
		}

		if len(stack) == 0 {
			return current
		} else {
			// Continue with enclosing element
			current = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
	}
}

// TypeAdapters_JSON_ELEMENT_Write is TypeAdapters.JSON_ELEMENT.write.
//
// Java: TypeAdapters.JSON_ELEMENT.write(JsonWriter, JsonElement)
func TypeAdapters_JSON_ELEMENT_Write(out *JsonWriter, value JsonElement) {
	if value == nil || value.IsJsonNull() {
		out.NullValue()
	} else if value.IsJsonPrimitive() {
		primitive := value.GetAsJsonPrimitive()
		if primitive.IsNumber() {
			out.ValueNumber(primitive.GetAsNumber())
		} else if primitive.IsBoolean() {
			out.ValueBoolean(primitive.GetAsBoolean())
		} else {
			s := primitive.GetAsString()
			out.ValueString(&s)
		}
	} else if value.IsJsonArray() {
		out.BeginArray()
		for e := range value.GetAsJsonArray().All() {
			TypeAdapters_JSON_ELEMENT_Write(out, e)
		}
		out.EndArray()
	} else if value.IsJsonObject() {
		out.BeginObject()
		for k, v := range value.GetAsJsonObject().EntrySet().All() {
			out.Name(k)
			TypeAdapters_JSON_ELEMENT_Write(out, v)
		}
		out.EndObject()
	} else {
		panic(NewIllegalArgumentExceptionMsg("Couldn't write class " + jsonElementClassName(value)))
	}
}

// Streams_Write writes the JSON element to the writer, recursively.
//
// Java: Streams.write(JsonElement, JsonWriter)
func Streams_Write(element JsonElement, writer *JsonWriter) {
	TypeAdapters_JSON_ELEMENT_Write(writer, element)
}
