// Port of test/com/google/debugging/sourcemap/SourceMapConsumerV3Test.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

func basicMapJSON() string {
	return gsonToJson(
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings("AAAAA,QAASA,UAAS,EAAG;").
			setSources("testcode").
			setNames("__BASIC__").
			build())
}

func assertSources(t *testing.T, got []*string, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("sources = %d entries, want %v", len(got), want)
	}
	for i := range want {
		if deref(got[i]) != want[i] {
			t.Fatalf("sources[%d] = %s, want %s", i, deref(got[i]), want[i])
		}
	}
}

func TestSources(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(basicMapJSON())

	assertSources(t, consumer.GetOriginalSources(), "testcode")
	if consumer.GetSourceRoot() != nil {
		t.Errorf("sourceRoot = %q, want null", *consumer.GetSourceRoot())
	}
}

func TestSectionsFormat_parsesSuccessfully(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(
		gsonToJson(
			testJsonBuilder_Create().
				setVersion(3).
				setFile("testcode").
				addSection(
					1,
					2,
					testJsonBuilder_Create().
						setVersion(3).
						setFile("part_a").
						setMappings("AAAAA,QAASA,UAAS,EAAG;").
						setSources("testcode.js").
						setNames("foo")).
				build()))
}

func TestSectionsFormat_fileNameIsOptional_inSectionAndTopLevel(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(
		gsonToJson(
			testJsonBuilder_Create().
				setVersion(3).
				// .setFile("testcode")
				addSection(
					1,
					2,
					testJsonBuilder_Create().
						setVersion(3).
						// .setFile("part_a")
						setMappings("AAAAA,QAASA,UAAS,EAAG;").
						setSources("testcode.js").
						setNames("foo")).
				build()))
}

func TestSectionsFormat_ignoresEmptySources(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(
		gsonToJson(
			testJsonBuilder_Create().
				setVersion(3).
				setFile("testcode").
				setSources().
				addSection(
					1,
					2,
					testJsonBuilder_Create().
						setVersion(3).
						setFile("part_a").
						setMappings("AAAAA,QAASA,UAAS,EAAG;").
						setSources("testcode.js").
						setNames("foo")).
				build()))
}

func TestSourcesWithRoot(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(
		gsonToJson(
			testJsonBuilder_Create().
				setVersion(3).
				setFile("testcode").
				setLineCount(1).
				setMappings("AAAAA,QAASA,UAAS,EAAG;").
				setSourceRoot("http://server/path/").
				setSources("testcode").
				setNames("__BASIC__").
				build()))

	//By default sourceRoot is not prepended
	assertSources(t, consumer.GetOriginalSources(), "testcode")
	if deref(consumer.GetSourceRoot()) != "http://server/path/" {
		t.Errorf("sourceRoot = %s", deref(consumer.GetSourceRoot()))
	}
}

func TestExtensions(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(
		gsonToJson(
			testJsonBuilder_Create().
				setVersion(3).
				setFile("testcode").
				setLineCount(1).
				setMappings("AAAAA,QAASA,UAAS,EAAG;").
				setSources("testcode").
				setNames("__BASIC__").
				setCustomProperty("x_org_int", jsonInt(2)).
				setCustomProperty("x_org_array", jutil.NewJsonArray()).
				build()))

	exts := consumer.GetExtensions()

	if len(exts) != 2 {
		t.Fatalf("extensions = %v", exts)
	}
	extByKey := make(map[string]any, len(exts))
	for _, ext := range exts {
		extByKey[ext.Key] = ext.Value
	}
	if _, ok := extByKey["org_int"]; ok {
		t.Errorf("unexpected key org_int")
	}
	if got := extByKey["x_org_int"].(jutil.JsonElement).GetAsInt(); got != 2 {
		t.Errorf("x_org_int = %d", got)
	}
	if !extByKey["x_org_array"].(*jutil.JsonArray).IsEmpty() {
		t.Errorf("x_org_array not empty")
	}
}

func TestSourceMappingExactMatch(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(basicMapJSON())

	mapping, _, ok := consumer.GetMappingForLineValue(1, 1)

	if !ok {
		t.Fatal("mapping is null")
	}
	if mapping.GetLineNumber() != 1 {
		t.Errorf("line = %d", mapping.GetLineNumber())
	}
	if mapping.GetPrecision() != sourcemap.Mapping_OriginalMapping_Precision_EXACT {
		t.Errorf("precision = %s", mapping.GetPrecision())
	}
}

func TestSourceMappingApproximatedLine(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(
		gsonToJson(
			testJsonBuilder_Create().
				setVersion(3).
				setMappings(
					";;;;;;;;;;;;;;;;;;IAAMA,K,GACL,eAAaC,EAAb,EAAiB;AAAA;;AAAA;;AAChB,OAAKA,EAAL,GAAUA,EAAV;AACA,C;;IAEIC,S;;;;;;;AACL,qBAAYD,EAAZ,EAAgB;AAAA;;AAAA,6BACTA,EADS;AAEf;;;EAHsBD,K;;AAKxB,IAAIG,CAAC,GAAG,IAAID,SAAJ,CAAc,UAAd,CAAR").
				setSourcesContent(
					"class Shape {\n"+
						"\tconstructor (id) {\n"+
						"\t\tthis.id = id;\n"+
						"\t}\n"+
						"}\n"+
						"class Rectangle extends Shape {\n"+
						"\tconstructor(id) {\n"+
						"\t\tsuper(id);\n"+
						"\t}\n"+
						"}\n"+
						"var s = new Rectangle(\"Shape ID\");\n").
				setSources("testcode").
				setNames("Shape", "id", "Rectangle", "s").
				build()))

	mapping, _, ok := consumer.GetMappingForLineValue(40, 10)

	if !ok {
		t.Fatal("mapping is null")
	}
	// The Previous line mapping was retrieved, and thus it is "approximated"
	if mapping.GetLineNumber() != 9 {
		t.Errorf("line = %d", mapping.GetLineNumber())
	}
	if mapping.GetPrecision() != sourcemap.Mapping_OriginalMapping_Precision_APPROXIMATE_LINE {
		t.Errorf("precision = %s", mapping.GetPrecision())
	}
}

func TestLargeMappings(t *testing.T) {
	// Generate a mapping with many entries to trigger array resizing.
	var mappings strings.Builder
	for i := 0; i < 2000; i++ {
		// AAAA = [0, 0, 0, 0]
		mappings.WriteString("AAAA")
		if i%2 == 0 {
			mappings.WriteString("A") // 5th value
		}
		mappings.WriteString(",")
	}
	mappings.WriteString(";")

	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(
		gsonToJson(
			testJsonBuilder_Create().
				setVersion(3).
				setFile("testcode").
				setLineCount(1).
				setMappings(mappings.String()).
				setSources("testcode").
				setNames("foo").
				build()))

	if consumer.GetLineCount() != 1 {
		t.Errorf("lineCount = %d", consumer.GetLineCount())
	}
	if _, _, ok := consumer.GetMappingForLineValue(1, 1); !ok {
		t.Error("mapping is null")
	}
}

func TestMixedEntryTypes(t *testing.T) {
	// Test a mix of 1, 4, and 5 value entries.
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(
		gsonToJson(
			testJsonBuilder_Create().
				setVersion(3).
				setFile("testcode").
				setLineCount(1).
				// col 0: A (unmapped)
				// col 1: CAAA (mapped to src index 0, line 0, col 0 - 4 values)
				// col 3: EAAAA (mapped to src index 0, line 0, col 0, name index 0 - 5 values)
				setMappings("A,CAAA,EAAAA;").
				setSources("testcode").
				setNames("name1").
				build()))

	m1, _, ok1 := consumer.GetMappingForLineValue(1, 1) // col 0
	m2, _, ok2 := consumer.GetMappingForLineValue(1, 2) // col 1
	m4, _, ok4 := consumer.GetMappingForLineValue(1, 4) // col 3

	if ok1 {
		t.Errorf("m1 = %v", m1)
	}
	if !ok2 || !ok4 {
		t.Fatalf("ok2 = %v, ok4 = %v", ok2, ok4)
	}

	if m2.GetOriginalFile() != "testcode" || m2.GetIdentifier() != "" {
		t.Errorf("m2 = %q %q", m2.GetOriginalFile(), m2.GetIdentifier())
	}
	if m4.GetOriginalFile() != "testcode" || m4.GetIdentifier() != "name1" {
		t.Errorf("m4 = %q %q", m4.GetOriginalFile(), m4.GetIdentifier())
	}
}

func TestReverseMapping(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(basicMapJSON())

	reverse := consumer.GetReverseMapping("testcode", 0, 0)
	if len(reverse) == 0 {
		t.Error("reverse mapping is empty")
	}
}

func TestBinarySearchEdgeCases(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	// Test with different number of entries to verify binary search 'mid'
	// calculation.
	for numEntries := 1; numEntries < 20; numEntries++ {
		var mappings strings.Builder
		for i := 0; i < numEntries; i++ {
			// VLQ for delta 1 is 'C'
			// AAAA = [0,0,0,0]
			// CAAA = [1,0,0,0]
			// Previous col starts at 0.
			// First entry: AAAA -> col 0
			// Subsequent entries: CAAA -> col 1, 2, 3...
			if i == 0 {
				mappings.WriteString("AAAA")
			} else {
				mappings.WriteString("CAAA")
			}
			mappings.WriteString(",")
		}
		mappings.WriteString(";")

		consumer.Parse(
			gsonToJson(
				testJsonBuilder_Create().
					setVersion(3).
					setFile("testcode").
					setLineCount(1).
					setMappings(mappings.String()).
					setSources("testcode").
					setNames("foo").
					build()))

		for col := 0; col < numEntries; col++ {
			mapping, _, ok := consumer.GetMappingForLineValue(1, col+1)
			if !ok {
				t.Fatalf("%d entries, col %d: null", numEntries, col)
			}
			// All map to source line 0, col 0 (OriginalMapping uses 1-indexed)
			if mapping.GetLineNumber() != 1 || mapping.GetColumnPosition() != 1 {
				t.Fatalf("%d entries, col %d: %d:%d", numEntries, col, mapping.GetLineNumber(), mapping.GetColumnPosition())
			}
		}
	}
}

type entryVisitorFunc func(sourceName, symbolName *string, sourcePos, startPos, endPos *sourcemap.FilePosition)

func (f entryVisitorFunc) Visit(sourceName, symbolName *string, sourcePos, startPos, endPos *sourcemap.FilePosition) {
	f(sourceName, symbolName, sourcePos, startPos, endPos)
}

func TestVisitMappings(t *testing.T) {
	consumer := sourcemap.NewSourceMapConsumerV3()
	consumer.Parse(basicMapJSON())

	var symbols []string
	consumer.VisitMappings(entryVisitorFunc(
		func(sourceName, symbolName *string, sourcePos, startPos, endPos *sourcemap.FilePosition) {
			if symbolName != nil {
				symbols = append(symbols, *symbolName)
			}
		}))

	if strings.Join(symbols, ",") != "__BASIC__,__BASIC__" {
		t.Errorf("symbols = %v", symbols)
	}
}

func TestInvalidEntryValuesThrows(t *testing.T) {
	// A VLQ of length 2, e.g. "AA" -> [0, 0]
	invalidMappings := "AA;"

	consumer := sourcemap.NewSourceMapConsumerV3()
	expected := assertThrows[*sourcemap.SourceMapParseException](t, func() {
		consumer.Parse(
			gsonToJson(
				testJsonBuilder_Create().
					setVersion(3).
					setFile("testcode").
					setLineCount(1).
					setMappings(invalidMappings).
					setSources("testcode").
					setNames("foo").
					build()))
	})

	if !strings.Contains(expected.GetMessage(), "Unexpected number of values for entry:2") {
		t.Errorf("message = %q", expected.GetMessage())
	}
}
