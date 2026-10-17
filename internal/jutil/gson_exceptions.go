// Port of com/google/gson/JsonParseException.java (Gson 2.9.1), with
// com/google/gson/JsonSyntaxException.java, com/google/gson/JsonIOException.java
// and com/google/gson/stream/MalformedJsonException.java

/*
 * Copyright (C) 2008 Google Inc. (JsonParseException.java, JsonIOException.java)
 * Copyright (C) 2010 Google Inc. (JsonSyntaxException.java, MalformedJsonException.java)
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

// Go-only: language-adapter — Gson's exception classes
// (com.google.gson.JsonParseException, JsonSyntaxException, JsonIOException,
// com.google.gson.stream.MalformedJsonException) and java.io.EOFException,
// which Gson's JsonReader throws.
//
// The pinned closure-compiler.jar shades Gson 2.9.1 with jarjar under
// com.google.javascript.jscomp.jarjar, so Throwable.toString() (and every
// message that embeds a cause, such as "JSON parse exception: " + ex in
// SourceMapObjectParser) carries the shaded class names. GsonPackage is that
// prefix.

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

// NewJsonParseException is new JsonParseException(String).
func NewJsonParseException(msg string) *JsonParseException {
	e := &JsonParseException{}
	InitThrowable(e, GsonPackage+".JsonParseException", msgPtr(msg), nil)
	return e
}

// NewJsonParseExceptionWithCause is new JsonParseException(Throwable).
func NewJsonParseExceptionWithCause(cause Throwable) *JsonParseException {
	e := &JsonParseException{}
	InitThrowable(e, GsonPackage+".JsonParseException", nil, cause)
	return e
}

// JsonSyntaxException is com.google.gson.JsonSyntaxException.
type JsonSyntaxException struct{ JsonParseException }

func (*JsonSyntaxException) jsonSyntaxException() {}

// IJsonSyntaxException matches JsonSyntaxException.
type IJsonSyntaxException interface {
	IJsonParseException
	jsonSyntaxException()
}

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

func (*JsonIOException) jsonIOException() {}

// IJsonIOException matches JsonIOException.
type IJsonIOException interface {
	IJsonParseException
	jsonIOException()
}

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
