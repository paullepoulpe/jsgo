// Port of com/google/gson/stream/JsonWriter.java (Gson 2.9.1): value(long)
// and flush()

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

// Go-only: language-adapter — JsonWriter.value(long) and JsonWriter.flush()
// (Gson), extending gson_json_writer.go.

// ValueLong encodes value.
//
// Java: JsonWriter.value(long)
func (w *JsonWriter) ValueLong(value int64) *JsonWriter {
	w.writeDeferredName()
	w.beforeValue()
	w.write(strconv.FormatInt(value, 10))
	return w
}

// Flush ensures all buffered data is written to the underlying Writer (when
// it has a Flush method) and flushes that writer.
//
// Java: JsonWriter.flush()
func (w *JsonWriter) Flush() {
	if w.stackSize == 0 {
		panic(NewIllegalStateExceptionMsg("JsonWriter is closed."))
	}
	if f, ok := w.out.(interface{ Flush() }); ok {
		f.Flush()
	}
}
