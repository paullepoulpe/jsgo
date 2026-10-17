// Port of test/com/google/debugging/sourcemap/TestJsonBuilder.java @ closure-compiler v20261015

/*
 * Copyright 2011 The Closure Compiler Authors.
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

package sourcemap_test

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// testJsonBuilder is a minimal builder for a Map representing a source map
// JSON file.
//
// This builder is designed and intended for tests. It does not check any
// invariants. It has only the bare minimum type checks. It does not support
// round-tripping of Java types.
//
// Go-only: language-adapter — the Java builds an ImmutableMap that the tests
// serialize with GSON.toJson; here the map is a Gson JsonObject (same
// insertion order) and gsonToJson serializes it the way Gson's default
// toJson does (compact, HTML-safe).
type testJsonBuilder struct {
	internal *jutil.JsonObject
	sections *jutil.JsonArray
}

func testJsonBuilder_Create() *testJsonBuilder {
	return &testJsonBuilder{internal: jutil.NewJsonObject(), sections: jutil.NewJsonArray()}
}

func jsonInt(i int) jutil.JsonElement {
	return jutil.NewJsonPrimitiveNumber(jutil.NewLazilyParsedNumber(jutil.ToJavaString(i)))
}

func jsonStrings(ss ...string) jutil.JsonElement {
	a := jutil.NewJsonArray()
	for _, s := range ss {
		a.Add(jutil.NewJsonPrimitiveString(s))
	}
	return a
}

func (b *testJsonBuilder) setVersion(version int) *testJsonBuilder {
	b.internal.Add("version", jsonInt(version))
	return b
}

func (b *testJsonBuilder) setFile(file string) *testJsonBuilder {
	b.internal.Add("file", jutil.NewJsonPrimitiveString(file))
	return b
}

func (b *testJsonBuilder) setLineCount(count int) *testJsonBuilder {
	b.internal.Add("lineCount", jsonInt(count))
	return b
}

func (b *testJsonBuilder) setMappings(mappings string) *testJsonBuilder {
	b.internal.Add("mappings", jutil.NewJsonPrimitiveString(mappings))
	return b
}

func (b *testJsonBuilder) setSourceRoot(root string) *testJsonBuilder {
	b.internal.Add("sourceRoot", jutil.NewJsonPrimitiveString(root))
	return b
}

func (b *testJsonBuilder) setSources(sources ...string) *testJsonBuilder {
	b.internal.Add("sources", jsonStrings(sources...))
	return b
}

func (b *testJsonBuilder) setSourcesContent(contents ...string) *testJsonBuilder {
	b.internal.Add("sourcesContent", jsonStrings(contents...))
	return b
}

func (b *testJsonBuilder) setNames(names ...string) *testJsonBuilder {
	b.internal.Add("names", jsonStrings(names...))
	return b
}

func (b *testJsonBuilder) addSection(line, column int, m *testJsonBuilder) *testJsonBuilder {
	section := jutil.NewJsonObject()
	offset := jutil.NewJsonObject()
	offset.Add("line", jsonInt(1))
	offset.Add("column", jsonInt(2))
	section.Add("offset", offset)
	section.Add("map", m.build())
	b.sections.Add(section)
	b.internal.Add("sections", b.sections.DeepCopy())
	return b
}

func (b *testJsonBuilder) setCustomProperty(name string, value jutil.JsonElement) *testJsonBuilder {
	b.internal.Add(name, value)
	return b
}

func (b *testJsonBuilder) build() *jutil.JsonObject {
	return b.internal
}

// gsonToJson is new Gson().toJson(map).
func gsonToJson(e jutil.JsonElement) string {
	var sb strings.Builder
	w := jutil.NewJsonWriter(&sb)
	w.SetHtmlSafe(true)
	w.SetLenient(true)
	jutil.Streams_Write(e, w)
	return sb.String()
}

// assertThrows is org.junit.Assert.assertThrows: body must throw an
// exception assignable to E.
func assertThrows[E any](t *testing.T, body func()) E {
	t.Helper()
	var caught E
	ok := false
	jutil.Try(body, func(e E) {
		caught = e
		ok = true
	})
	if !ok {
		t.Fatalf("expected %T to be thrown", (*E)(nil))
	}
	return caught
}

// deref is the Java value of a nullable String ("<null>" for null, for
// test messages).
func deref(s *string) string {
	if s == nil {
		return "<null>"
	}
	return *s
}
