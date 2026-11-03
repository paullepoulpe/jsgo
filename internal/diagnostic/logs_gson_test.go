// Port of test/com/google/javascript/jscomp/diagnostic/LogsGsonTest.java @ closure-compiler v20261028

/*
 * Copyright 2021 The Closure Compiler Authors.
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

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestAdapter_write_logsGsonable(t *testing.T) {
	assertLogsGson(t, []any{&TestLogsGsonable{}}, "[[0]]")
}

type TestLogsGsonable struct{}

func (*TestLogsGsonable) ToLogsGson() any {
	return []int32{0} // Java: ImmutableList.of(0)
}

type testMultimap struct {
	m *jutil.OrderedMap[string, *jutil.OrderedSet[any]]
}

func newTestMultimap() *testMultimap {
	return &testMultimap{m: jutil.NewOrderedMap[string, *jutil.OrderedSet[any]]()}
}

func (mm *testMultimap) Put(k string, v any) {
	s, ok := mm.m.Lookup(k)
	if !ok {
		s = jutil.NewOrderedSet[any]()
		mm.m.Put(k, s)
	}
	s.Add(v)
}

func (mm *testMultimap) AsMapObject() any { return mm.m }

func TestAdapter_write_multimap(t *testing.T) {
	// Given
	// Java: ImmutableMultimap.builder().put("foo", 0).put("foo", 1).put("bar", 2).build();
	value := newTestMultimap()
	value.Put("foo", int32(0))
	value.Put("foo", int32(1))
	value.Put("bar", int32(2))

	// Then
	assertLogsGson(t, value, "{\"foo\":[0,1],\"bar\":[2]}")
}

func assertLogsGson(t *testing.T, value any, want string) {
	t.Helper()
	if got := LogsGson_ToJson(value); got != want {
		t.Errorf("LogsGson.toJson = %s, want %s", got, want)
	}
}
