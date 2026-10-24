// Port of com/google/javascript/jscomp/diagnostic/LogsGson.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package diagnostic

import "github.com/paullepoulpe/jsgo/internal/jutil"

// LogsGson is a Gson instance tailored to generating logs output.
//
// Java: public final class LogsGson (only static members).
type LogsGson struct{}

// LogsGson_Able is a conversion to a preferred logging representation.
//
// When an implementor is logged, the result of this method will be JSON
// stringified in place of the implementor.
//
// Use this interface when the default Gson serialization for a class is bad
// and the class should always be logged another way. Primitive wrappers and
// entity type are ideal candidates.
//
// Java: public interface LogsGson.Able
type LogsGson_Able interface {
	// Java: Able.toLogsGson()
	ToLogsGson() any
}

// LogsGson_ToJson serializes o with LogsGson_GSON.
//
// Java: LogsGson.toJson(Object)
func LogsGson_ToJson(o any) string {
	return LogsGson_GSON.ToJson(o)
}

// LogsGson_GSON is Gson with two type hierarchy adapters: LogsGson_Able
// values are written as the JSON of their ToLogsGson(), and Guava
// Multimaps as the JSON of their asMap().
var LogsGson_GSON = jutil.NewGsonBuilder().
	RegisterTypeHierarchyAdapter(
		func(value any) bool { _, ok := value.(LogsGson_Able); return ok }, // Java: Able.class
		&logsGson_ableTypeAdapter{}).
	RegisterTypeHierarchyAdapter(
		func(value any) bool { _, ok := value.(jutil.JavaMultimap); return ok }, // Java: Multimap.class
		&logsGson_multimapTypeAdapter{}).
	Create()

// logsGson_ableTypeAdapter is the anonymous TypeAdapter<Able>.
type logsGson_ableTypeAdapter struct{}

// Java: new TypeAdapter<Able>() {...}.read(JsonReader)
func (*logsGson_ableTypeAdapter) Read(in *jutil.JsonReader) any {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: new TypeAdapter<Able>() {...}.write(JsonWriter, Able)
func (*logsGson_ableTypeAdapter) Write(out *jutil.JsonWriter, value any) {
	json := LogsGson_GSON.ToJson(value.(LogsGson_Able).ToLogsGson())
	out.JsonValue(&json)
}

// logsGson_multimapTypeAdapter is the anonymous TypeAdapter<Multimap<?, ?>>.
type logsGson_multimapTypeAdapter struct{}

// Java: new TypeAdapter<Multimap<?, ?>>() {...}.read(JsonReader)
func (*logsGson_multimapTypeAdapter) Read(in *jutil.JsonReader) any {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: new TypeAdapter<Multimap<?, ?>>() {...}.write(JsonWriter, Multimap<?, ?>)
func (*logsGson_multimapTypeAdapter) Write(out *jutil.JsonWriter, value any) {
	json := LogsGson_GSON.ToJson(value.(jutil.JavaMultimap).AsMapObject())
	out.JsonValue(&json)
}

// NewLogsGson is the private constructor, which always throws.
//
// Java: LogsGson()
func NewLogsGson() *LogsGson {
	panic(jutil.NewAssertionErrorNoMessage())
}
