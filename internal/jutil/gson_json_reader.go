// Port of com/google/gson/stream/JsonReader.java (Gson 2.9.1), with
// com/google/gson/stream/JsonToken.java and JsonScope.java (all three files
// carry the same header)

/*
 * Copyright (C) 2010 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jutil

import "strconv"

// JsonToken is com.google.gson.stream.JsonToken.
type JsonToken int32

const (
	JsonToken_BEGIN_ARRAY JsonToken = iota
	JsonToken_END_ARRAY
	JsonToken_BEGIN_OBJECT
	JsonToken_END_OBJECT
	JsonToken_NAME
	JsonToken_STRING
	JsonToken_NUMBER
	JsonToken_BOOLEAN
	JsonToken_NULL
	JsonToken_END_DOCUMENT
)

var jsonTokenNames = [...]string{
	"BEGIN_ARRAY", "END_ARRAY", "BEGIN_OBJECT", "END_OBJECT", "NAME",
	"STRING", "NUMBER", "BOOLEAN", "NULL", "END_DOCUMENT",
}

// String is the enum constant's name().
func (t JsonToken) String() string { return jsonTokenNames[t] }

// PeekToken returns the public JsonToken without consuming it.
func (s *JsonScanner) PeekToken() (JsonToken, Throwable) {
	p, err := s.Peek()
	if err != nil {
		return 0, err
	}
	switch p {
	case JsonPeekBeginObject:
		return JsonToken_BEGIN_OBJECT, nil
	case JsonPeekEndObject:
		return JsonToken_END_OBJECT, nil
	case JsonPeekBeginArray:
		return JsonToken_BEGIN_ARRAY, nil
	case JsonPeekEndArray:
		return JsonToken_END_ARRAY, nil
	case JsonPeekSingleQuotedName, JsonPeekDoubleQuotedName, JsonPeekUnquotedName:
		return JsonToken_NAME, nil
	case JsonPeekTrue, JsonPeekFalse:
		return JsonToken_BOOLEAN, nil
	case JsonPeekNull:
		return JsonToken_NULL, nil
	case JsonPeekSingleQuoted, JsonPeekDoubleQuoted, JsonPeekUnquoted:
		return JsonToken_STRING, nil
	case JsonPeekLong, JsonPeekNumber:
		return JsonToken_NUMBER, nil
	case JsonPeekEOF:
		return JsonToken_END_DOCUMENT, nil
	default:
		return 0, NewAssertionErrorNoMessage()
	}
}

func (s *JsonScanner) unexpectedTokenErr(expected string) Throwable {
	tok, err := s.PeekToken()
	if err != nil {
		return err
	}
	return NewIllegalStateExceptionMsg("Expected " + expected + " but was " + tok.String() + s.LocationString())
}

// HasNext reports whether the current array or object has another element.
func (s *JsonScanner) HasNext() (bool, Throwable) {
	p, err := s.Peek()
	if err != nil {
		return false, err
	}
	return p != JsonPeekEndObject && p != JsonPeekEndArray && p != JsonPeekEOF, nil
}

// BeginObject consumes '{'.
func (s *JsonScanner) BeginObject() Throwable {
	p, err := s.Peek()
	if err != nil {
		return err
	}
	if p != JsonPeekBeginObject {
		return s.unexpectedTokenErr("BEGIN_OBJECT")
	}
	s.pushScope(jsonScope_EMPTY_OBJECT)
	s.peeked = JsonPeekNone
	return nil
}

// EndObject consumes '}'.
func (s *JsonScanner) EndObject() Throwable {
	p, err := s.Peek()
	if err != nil {
		return err
	}
	if p != JsonPeekEndObject {
		return s.unexpectedTokenErr("END_OBJECT")
	}
	s.stack[len(s.stack)-1] = jsonFrame{}
	s.stack = s.stack[:len(s.stack)-1]
	s.incTopIndex()
	s.peeked = JsonPeekNone
	return nil
}

// BeginArray consumes '['.
func (s *JsonScanner) BeginArray() Throwable {
	p, err := s.Peek()
	if err != nil {
		return err
	}
	if p != JsonPeekBeginArray {
		return s.unexpectedTokenErr("BEGIN_ARRAY")
	}
	s.pushScope(jsonScope_EMPTY_ARRAY)
	s.peeked = JsonPeekNone
	return nil
}

// EndArray consumes ']'.
func (s *JsonScanner) EndArray() Throwable {
	p, err := s.Peek()
	if err != nil {
		return err
	}
	if p != JsonPeekEndArray {
		return s.unexpectedTokenErr("END_ARRAY")
	}
	s.stack = s.stack[:len(s.stack)-1]
	s.incTopIndex()
	s.peeked = JsonPeekNone
	return nil
}

// NextName consumes and returns a property name.
func (s *JsonScanner) NextName() (string, Throwable) {
	p, err := s.Peek()
	if err != nil {
		return "", err
	}
	var name string
	switch p {
	case JsonPeekDoubleQuotedName:
		name, err = s.readQuotedValue('"')
	case JsonPeekSingleQuotedName:
		name, err = s.readQuotedValue('\'')
	case JsonPeekUnquotedName:
		name, err = s.readUnquotedValue()
	default:
		return "", s.unexpectedTokenErr("a name")
	}
	if err != nil {
		return "", err
	}
	s.peeked = JsonPeekNone
	s.setTopName(name)
	return name, nil
}

// NextString consumes and returns a string or number token's string form.
func (s *JsonScanner) NextString() (string, Throwable) {
	p, err := s.Peek()
	if err != nil {
		return "", err
	}
	var res string
	switch p {
	case JsonPeekDoubleQuoted:
		res, err = s.readQuotedValue('"')
	case JsonPeekSingleQuoted:
		res, err = s.readQuotedValue('\'')
	case JsonPeekUnquoted:
		res, err = s.readUnquotedValue()
	case JsonPeekLong:
		res = strconv.FormatInt(s.peekedLong, 10)
	case JsonPeekNumber:
		res = s.json[s.pos : s.pos+s.peekedNumberLength]
		s.pos += s.peekedNumberLength
	default:
		return "", s.unexpectedTokenErr("a string")
	}
	if err != nil {
		return "", err
	}
	s.peeked = JsonPeekNone
	s.incTopIndex()
	return res, nil
}

// NextBoolean consumes and returns a boolean token.
func (s *JsonScanner) NextBoolean() (bool, Throwable) {
	p, err := s.Peek()
	if err != nil {
		return false, err
	}
	switch p {
	case JsonPeekTrue:
		s.peeked = JsonPeekNone
		s.incTopIndex()
		return true, nil
	case JsonPeekFalse:
		s.peeked = JsonPeekNone
		s.incTopIndex()
		return false, nil
	default:
		return false, s.unexpectedTokenErr("a boolean")
	}
}

// NextNull consumes a null token.
func (s *JsonScanner) NextNull() Throwable {
	p, err := s.Peek()
	if err != nil {
		return err
	}
	if p != JsonPeekNull {
		return s.unexpectedTokenErr("null")
	}
	s.peeked = JsonPeekNone
	s.incTopIndex()
	return nil
}

// ConsumeScalar consumes an already-peeked Long, Number, Boolean, or Null token.
func (s *JsonScanner) ConsumeScalar() {
	if s.peeked == JsonPeekNumber {
		s.pos += s.peekedNumberLength
	}
	s.peeked = JsonPeekNone
	s.incTopIndex()
}

// SkipValue skips the next value recursively with JsonReader.skipValue() semantics.
func (s *JsonScanner) SkipValue() Throwable {
	count := 0
	for {
		p, err := s.Peek()
		if err != nil {
			return err
		}
		switch p {
		case JsonPeekBeginArray:
			s.pushScope(jsonScope_EMPTY_ARRAY)
			count++
		case JsonPeekBeginObject:
			s.pushScope(jsonScope_EMPTY_OBJECT)
			count++
		case JsonPeekEndArray, JsonPeekEndObject:
			s.stack = s.stack[:len(s.stack)-1]
			count--
		case JsonPeekUnquotedName, JsonPeekUnquoted:
			if _, err := s.readUnquotedValue(); err != nil {
				return err
			}
		case JsonPeekSingleQuoted, JsonPeekSingleQuotedName:
			if err := s.skipQuotedValue('\''); err != nil {
				return err
			}
		case JsonPeekDoubleQuoted, JsonPeekDoubleQuotedName:
			if err := s.skipQuotedValue('"'); err != nil {
				return err
			}
		case JsonPeekNumber:
			s.pos += s.peekedNumberLength
		}
		s.peeked = JsonPeekNone
		if count == 0 {
			break
		}
	}
	s.incTopIndex()
	s.setTopName("null")
	return nil
}

func (s *JsonScanner) skipTerminalAsElement(p uint8) Throwable {
	switch p {
	case JsonPeekDoubleQuoted:
		if err := s.skipQuotedValue('"'); err != nil {
			return err
		}
		s.peeked = JsonPeekNone
		s.incTopIndex()
		return nil
	case JsonPeekSingleQuoted:
		if err := s.skipQuotedValue('\''); err != nil {
			return err
		}
		s.peeked = JsonPeekNone
		s.incTopIndex()
		return nil
	case JsonPeekUnquoted:
		if _, err := s.readUnquotedValue(); err != nil {
			return err
		}
		s.peeked = JsonPeekNone
		s.incTopIndex()
		return nil
	case JsonPeekLong, JsonPeekNumber, JsonPeekTrue, JsonPeekFalse, JsonPeekNull:
		s.ConsumeScalar()
		return nil
	default:
		tok, err := s.PeekToken()
		if err != nil {
			return err
		}
		return NewIllegalStateExceptionMsg("Unexpected token: " + tok.String())
	}
}

// SkipValueAsElement consumes and syntax-validates one JSON value while updating
// path state identically to TypeAdapters.JSON_ELEMENT.read without allocating DOM nodes.
func (s *JsonScanner) SkipValueAsElement() Throwable {
	p, err := s.Peek()
	if err != nil {
		return err
	}
	switch p {
	case JsonPeekBeginObject:
		_ = s.BeginObject()
	case JsonPeekBeginArray:
		_ = s.BeginArray()
	default:
		return s.skipTerminalAsElement(p)
	}

	depth := 1
	for depth > 0 {
		for {
			more, err := s.HasNext()
			if err != nil {
				return err
			}
			if !more {
				break
			}
			if sc := s.topScope(); sc != jsonScope_EMPTY_ARRAY && sc != jsonScope_NONEMPTY_ARRAY {
				if _, err := s.NextName(); err != nil {
					return err
				}
			}
			vp, err := s.Peek()
			if err != nil {
				return err
			}
			switch vp {
			case JsonPeekBeginObject:
				_ = s.BeginObject()
				depth++
			case JsonPeekBeginArray:
				_ = s.BeginArray()
				depth++
			default:
				if err := s.skipTerminalAsElement(vp); err != nil {
					return err
				}
			}
		}
		if sc := s.topScope(); sc == jsonScope_EMPTY_ARRAY || sc == jsonScope_NONEMPTY_ARRAY {
			if err := s.EndArray(); err != nil {
				return err
			}
		} else {
			if err := s.EndObject(); err != nil {
				return err
			}
		}
		depth--
	}
	return nil
}

// JsonReader is com.google.gson.stream.JsonReader: reads a JSON encoded
// value from a JavaReader as a stream of tokens.
type JsonReader struct {
	in     JavaReader
	s      JsonScanner
	loaded bool
}

// NewJsonReader creates a new instance that reads a JSON-encoded stream from in.
//
// Java: JsonReader(Reader)
func NewJsonReader(in JavaReader) *JsonReader {
	if in == nil {
		panic(NewNullPointerExceptionMsg("in == null"))
	}
	r := &JsonReader{in: in}
	r.s.Init("", false)
	return r
}

func (r *JsonReader) ensureLoaded() {
	if r.loaded {
		return
	}
	lenient := r.s.IsLenient()
	text := CharStreams_ToString(r.in)
	r.s.Init(text, lenient)
	r.loaded = true
}

func mustNoErr(err Throwable) {
	if err != nil {
		panic(err)
	}
}

func (r *JsonReader) SetLenient(lenient bool) { r.s.SetLenient(lenient) }
func (r *JsonReader) IsLenient() bool         { return r.s.IsLenient() }
func (r *JsonReader) BeginArray()             { r.ensureLoaded(); mustNoErr(r.s.BeginArray()) }
func (r *JsonReader) EndArray()               { r.ensureLoaded(); mustNoErr(r.s.EndArray()) }
func (r *JsonReader) BeginObject()            { r.ensureLoaded(); mustNoErr(r.s.BeginObject()) }
func (r *JsonReader) EndObject()              { r.ensureLoaded(); mustNoErr(r.s.EndObject()) }
func (r *JsonReader) NextNull()               { r.ensureLoaded(); mustNoErr(r.s.NextNull()) }
func (r *JsonReader) SkipValue()              { r.ensureLoaded(); mustNoErr(r.s.SkipValue()) }
func (r *JsonReader) GetPath() string         { return r.s.GetPath() }
func (r *JsonReader) getPreviousPath() string { return r.s.GetPreviousPath() }
func (r *JsonReader) String() string          { return "JsonReader" + r.s.LocationString() }

func (r *JsonReader) HasNext() bool {
	r.ensureLoaded()
	ok, err := r.s.HasNext()
	mustNoErr(err)
	return ok
}

func (r *JsonReader) Peek() JsonToken {
	r.ensureLoaded()
	tok, err := r.s.PeekToken()
	mustNoErr(err)
	return tok
}

func (r *JsonReader) NextName() string {
	r.ensureLoaded()
	name, err := r.s.NextName()
	mustNoErr(err)
	return name
}

func (r *JsonReader) NextString() string {
	r.ensureLoaded()
	str, err := r.s.NextString()
	mustNoErr(err)
	return str
}

func (r *JsonReader) NextBoolean() bool {
	r.ensureLoaded()
	v, err := r.s.NextBoolean()
	mustNoErr(err)
	return v
}

func (r *JsonReader) Close() {
	r.loaded = true
	r.s.Close()
	r.in.Close()
}
