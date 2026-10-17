// Port of com/google/gson/TypeAdapter.java (Gson 2.9.1)

/*
 * Copyright (C) 2011 Google Inc.
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

// Go-only: language-adapter — com.google.gson.TypeAdapter (Gson 2.9.1, as
// shaded into the pinned closure-compiler.jar).

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
