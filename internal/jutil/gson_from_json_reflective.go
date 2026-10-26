// Port of com/google/gson/Gson.java (Gson 2.9.1): fromJson(JsonReader, Type),
// with ReflectiveTypeAdapterFactory.Adapter.read of
// com/google/gson/internal/bind/ReflectiveTypeAdapterFactory.java

/*
 * Copyright (C) 2008 Google Inc. (Gson.java)
 * Copyright (C) 2011 Google Inc. (ReflectiveTypeAdapterFactory.java)
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
	"reflect"
	"strings"
)

// Go-only: language-adapter — Gson 2.9.1's reflective deserialization
// (ReflectiveTypeAdapterFactory.Adapter.read) of an object whose fields are
// Strings, as Gson.fromJson(JsonReader, Class) does for such a class. The Go
// struct's fields are named like toJson names them (gson_to_json.go): the
// `gson` tag (@SerializedName's value) or the field name; a `gsonalt` tag
// lists @SerializedName's alternates, comma-separated. Fields must be
// string or *string (a Java String; nil is null); other kinds throw
// UnsupportedOperationException.

// FromJsonReaderReflective is Gson.fromJson(JsonReader, Class) for a class
// read by the reflective adapter: construct is the class's no-argument
// constructor, returning a pointer to the struct. It returns nil for a JSON
// null and for an empty document (JSON 1.5 compatibility), the constructed
// object otherwise.
//
// Java: Gson.fromJson(JsonReader, Type)
func (g *Gson) FromJsonReaderReflective(reader *JsonReader, construct func() any) (result any) {
	isEmpty := true
	oldLenient := reader.IsLenient()
	reader.SetLenient(true)
	defer reader.SetLenient(oldLenient)
	Try(func() {
		reader.Peek()
		isEmpty = false
		result = gson_reflectiveRead(reader, construct)
	}, func(t Throwable) {
		switch e := t.(type) {
		case IEOFException:
			if isEmpty {
				result = nil
				return
			}
			panic(NewJsonSyntaxExceptionWithCause(e))
		case IIllegalStateException:
			panic(NewJsonSyntaxExceptionWithCause(e))
		case IIOException:
			panic(NewJsonSyntaxExceptionWithCause(e))
		case *AssertionError:
			error_ := NewAssertionError("AssertionError (GSON " + GsonBuildConfig_VERSION + "): " + javaMessageOrNull(e))
			error_.InitCause(e)
			panic(error_)
		default:
			panic(t)
		}
	})
	return result
}

// gson_reflectiveRead is ReflectiveTypeAdapterFactory.Adapter.read.
func gson_reflectiveRead(in *JsonReader, construct func() any) any {
	if in.Peek() == JsonToken_NULL {
		in.NextNull()
		return nil
	}
	instance := construct()
	sv := reflect.ValueOf(instance).Elem()
	fields := gson_readBoundFields(sv.Type())
	Try(func() {
		in.BeginObject()
		for in.HasNext() {
			name := in.NextName()
			index, ok := fields[name]
			if !ok {
				in.SkipValue()
			} else {
				gson_readStringField(in, sv.FieldByIndex(index))
			}
		}
	}, func(e IIllegalStateException) {
		panic(NewJsonSyntaxExceptionWithCause(e))
	})
	in.EndObject()
	return instance
}

// gson_readBoundFields maps each JSON name (and alternate) of the struct's
// fields to the field (ReflectiveTypeAdapterFactory.getBoundFields).
func gson_readBoundFields(t reflect.Type) map[string][]int {
	fields := map[string][]int{}
	for _, f := range gson_boundFields(t) {
		fields[f.name] = f.index
		sf := t.FieldByIndex(f.index)
		if alt := sf.Tag.Get("gsonalt"); alt != "" {
			for _, a := range strings.Split(alt, ",") {
				if _, dup := fields[a]; dup {
					panic(NewIllegalArgumentExceptionMsg(t.String() + " declares multiple JSON fields named " + a))
				}
				fields[a] = f.index
			}
		}
	}
	return fields
}

// gson_readStringField reads a value with TypeAdapters.STRING (booleans are
// coerced to strings) and stores it in the field; a null is stored unless
// the field is a primitive (a Go string).
func gson_readStringField(in *JsonReader, fv reflect.Value) {
	var value *string
	switch in.Peek() {
	case JsonToken_NULL:
		in.NextNull()
	case JsonToken_BOOLEAN:
		s := "false"
		if in.NextBoolean() {
			s = "true"
		}
		value = &s
	default:
		s := in.NextString()
		value = &s
	}
	switch fv.Kind() {
	case reflect.String:
		if value != nil {
			fv.SetString(*value)
		}
	case reflect.Pointer:
		if fv.Type().Elem().Kind() != reflect.String {
			panic(NewUnsupportedOperationExceptionMsg("Gson reflective read of " + fv.Type().String()))
		}
		fv.Set(reflect.ValueOf(value))
	default:
		panic(NewUnsupportedOperationExceptionMsg("Gson reflective read of " + fv.Type().String()))
	}
}
