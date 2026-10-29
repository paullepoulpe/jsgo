// Port of test/com/google/debugging/sourcemap/SourceMapTestCase.java @ closure-compiler v20261021

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

// Package sourcemaptesting holds the source-map test support that compiles
// JavaScript: com.google.debugging.sourcemap.SourceMapTestCase, the base of
// SourceMapGeneratorV3Test (package sourcemap) and of jscomp's SourceMapTest,
// with the jscomp.testing.TestErrorManager it uses.
//
// Go-only: test-helper — Java's test classpath lets a debugging.sourcemap
// test class import jscomp, which imports debugging.sourcemap. Go test files
// cannot be imported, and sourcemap's own tests cannot import jscomp, so this
// shared base is an importable package (like rhino/testing's rhinotesting),
// used by the external test packages sourcemap_test and jscomp_test.
package sourcemaptesting

import (
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// SourceMapTestCase is the abstract SourceMapTestCase: subclasses embed
// SourceMapTestCaseBase and give the format and consumer to test.
//
// Java: public abstract class SourceMapTestCase
type SourceMapTestCase interface {
	// Java: SourceMapTestCase.getSourceMapFormat()
	GetSourceMapFormat() jscomp.SourceMap_Format
	// Java: SourceMapTestCase.getSourceMapConsumer()
	GetSourceMapConsumer() sourcemap.SourceMapConsumer
	// Java: SourceMapTestCase.getCompilerOptions()
	GetCompilerOptions() *jscomp.CompilerOptions

	GetSourceMapTestCaseBase() *SourceMapTestCaseBase
}

// SourceMapTestCaseBase holds SourceMapTestCase's fields and methods.
type SourceMapTestCaseBase struct {
	self SourceMapTestCase

	// Go-only: test-helper — the Go test running the current test method;
	// assertions fail it.
	t testing.TB

	validateColumns bool

	DetailLevel                    jscomp.SourceMap_DetailLevel
	SourceMapIncludeSourcesContent bool
}

// InitSourceMapTestCaseBase constructs the base of self, running on the Go
// test t.
func (b *SourceMapTestCaseBase) InitSourceMapTestCaseBase(self SourceMapTestCase, t testing.TB) {
	b.self = self
	b.t = t
	b.validateColumns = true
	b.DetailLevel = jscomp.SourceMap_DetailLevel_ALL
	b.SourceMapIncludeSourcesContent = false
}

// GetSourceMapTestCaseBase gives access to the base.
func (b *SourceMapTestCaseBase) GetSourceMapTestCaseBase() *SourceMapTestCaseBase { return b }

// SetT gives the base the Go test running the current test method.
func (b *SourceMapTestCaseBase) SetT(t testing.TB) { b.t = t }

// Java: SourceMapTestCase.disableColumnValidation()
func (b *SourceMapTestCaseBase) DisableColumnValidation() {
	b.validateColumns = false
}

// Java: SourceMapTestCase.EXTERNS
var SourceMapTestCase_EXTERNS = []*jscomp.SourceFile{jscomp.SourceFile_FromCode("externs", "")}

// Java: SourceMapTestCase.RunResult
type SourceMapTestCase_RunResult struct {
	GeneratedSource      string
	SourceMap            *jscomp.SourceMap
	SourceMapFileContent string
}

// Java: SourceMapTestCase.Token
type SourceMapTestCase_Token struct {
	tokenName string
	inputName string
	position  *sourcemap.FilePosition
}

// Java: SourceMapTestCase.Token(String, String, FilePosition)
func newSourceMapTestCase_Token(tokenName string, inputName string, position *sourcemap.FilePosition) *SourceMapTestCase_Token {
	return &SourceMapTestCase_Token{tokenName: tokenName, inputName: inputName, position: position}
}

// Java: SourceMapTestCase.setUp()
func (b *SourceMapTestCaseBase) SetUp() {
	b.DetailLevel = jscomp.SourceMap_DetailLevel_ALL
}

// CheckSourceMap creates a source map for the given JS code and asserts it
// is equal to the expected golden map.
//
// Java: SourceMapTestCase.checkSourceMap(String, ImmutableMap<String, ?>)
func (b *SourceMapTestCaseBase) CheckSourceMap(js string, expectedMap *jutil.JsonObject) {
	b.t.Helper()
	b.CheckSourceMapWithFileName("testcode", js, expectedMap)
}

// Java: SourceMapTestCase.checkSourceMap(String, String, ImmutableMap<String, ?>)
func (b *SourceMapTestCaseBase) CheckSourceMapWithFileName(fileName string, js string, expectedMap *jutil.JsonObject) {
	b.t.Helper()
	result := b.Compile(js, fileName)
	// We round-trip the expectation to coerce any of its types in the same way as the actual map.
	roundTrippedExpectedMap := jutil.NewGson().FromJsonObject(jutil.NewGson().ToJsonElement(expectedMap))

	resultMap := jutil.NewGson().FromJsonObject(result.SourceMapFileContent)
	if !sourceMapTestCase_jsonMapEquals(resultMap, roundTrippedExpectedMap) {
		b.t.Fatalf("%s\nexpected: %s\nbut was : %s", result.GeneratedSource, roundTrippedExpectedMap, resultMap)
	}

	resultMapMap := jutil.NewGson().FromJsonObject(b.GetSourceMap(result))
	if !sourceMapTestCase_jsonMapEquals(resultMapMap, roundTrippedExpectedMap) {
		b.t.Fatalf("%s\nexpected: %s\nbut was : %s", result.GeneratedSource, roundTrippedExpectedMap, resultMapMap)
	}
}

// sourceMapTestCase_jsonMapEquals is Map.equals on two maps that Gson
// deserialized as Map<String, ?>: LinkedTreeMaps compared by entries,
// ArrayLists in order, and every JSON number a Double.
//
// Go-only: language-adapter — Gson.fromJson(json, Map<String, ?>) is the
// JsonElement tree here, compared the way the Java objects compare.
func sourceMapTestCase_jsonMapEquals(a, b jutil.JsonElement) bool {
	switch {
	case a.IsJsonNull() || b.IsJsonNull():
		return a.IsJsonNull() && b.IsJsonNull()
	case a.IsJsonObject():
		if !b.IsJsonObject() || a.GetAsJsonObject().Size() != b.GetAsJsonObject().Size() {
			return false
		}
		for k, v := range a.GetAsJsonObject().EntrySet().All() {
			w := b.GetAsJsonObject().Get(k)
			if w == nil || !sourceMapTestCase_jsonMapEquals(v, w) {
				return false
			}
		}
		return true
	case a.IsJsonArray():
		if !b.IsJsonArray() || a.GetAsJsonArray().Size() != b.GetAsJsonArray().Size() {
			return false
		}
		for i := 0; i < a.GetAsJsonArray().Size(); i++ {
			if !sourceMapTestCase_jsonMapEquals(a.GetAsJsonArray().Get(i), b.GetAsJsonArray().Get(i)) {
				return false
			}
		}
		return true
	}
	if !b.IsJsonPrimitive() {
		return false
	}
	pa, pb := a.GetAsJsonPrimitive(), b.GetAsJsonPrimitive()
	if pa.IsNumber() && pb.IsNumber() {
		x, _ := strconv.ParseFloat(pa.GetAsString(), 64)
		y, _ := strconv.ParseFloat(pb.GetAsString(), 64)
		return x == y
	}
	return pa.IsNumber() == pb.IsNumber() &&
		pa.IsBoolean() == pb.IsBoolean() &&
		pa.GetAsString() == pb.GetAsString()
}

// Java: SourceMapTestCase.getSourceMap(RunResult)
func (b *SourceMapTestCaseBase) GetSourceMap(result *SourceMapTestCase_RunResult) string {
	var sb strings.Builder
	result.SourceMap.AppendTo(&sb, jutil.StringPtr("testcode"))
	return sb.String()
}

// SourceMapTestInput is a named JS input file for Check/CheckWithSupplier.
type SourceMapTestInput struct {
	Name string
	Code string
}

type sourceMapTokenMap struct {
	order []string
	byKey map[string]*SourceMapTestCase_Token
}

func newSourceMapTokenMap() *sourceMapTokenMap {
	return &sourceMapTokenMap{byKey: make(map[string]*SourceMapTestCase_Token)}
}

func (m *sourceMapTokenMap) put(k string, v *SourceMapTestCase_Token) {
	if _, ok := m.byKey[k]; !ok {
		m.order = append(m.order, k)
	}
	m.byKey[k] = v
}

// sourceMapTestCase_findTokensInInputs finds the all the __XX__ tokens in
// the given JavaScript string.
//
// Java: SourceMapTestCase.findTokens(Map<String, String>)
func sourceMapTestCase_findTokensInInputs(inputs []SourceMapTestInput) *sourceMapTokenMap {
	tokens := newSourceMapTokenMap()

	for _, in := range inputs {
		sourceMapTestCase_findTokensInto(tokens, in.Name, in.Code)
	}

	return tokens
}

// sourceMapTestCase_findTokens finds the all the __XX__ tokens in the given
// JavaScript string.
//
// Java: SourceMapTestCase.findTokens(String)
func sourceMapTestCase_findTokens(src string) *sourceMapTokenMap {
	tokens := newSourceMapTokenMap()

	sourceMapTestCase_findTokensInto(tokens, "", src)

	return tokens
}

// sourceMapTestCase_findTokensInto finds the all the __XX__ tokens in the
// given JavaScript string.
//
// Java: SourceMapTestCase.findTokens(Map<String, Token>, String, String)
func sourceMapTestCase_findTokensInto(
	tokens *sourceMapTokenMap, inputName string, jsString string) *sourceMapTokenMap {

	js := jutil.StringToUTF16(jsString)
	currentLine := 0
	positionOffset := 0

	for i := 0; i < len(js); i++ {
		current := js[i]

		if current == '\n' {
			positionOffset = i + 1
			currentLine++
			continue
		}

		if current == '_' && (i < len(js)-5) {
			// Check for the _ token.
			if js[i+1] != '_' {
				continue
			}

			// Loop until we have another _ token.
			var tokenName []uint16

			j := i + 2
			for ; j < len(js); j++ {
				if js[j] == '_' {
					break
				}

				tokenName = append(tokenName, js[j])
			}

			if len(tokenName) > 0 {
				currentPosition := i - positionOffset
				token :=
					newSourceMapTestCase_Token(
						jutil.StringFromUTF16(tokenName), inputName, sourcemap.NewFilePosition(currentLine, currentPosition))
				tokens.put(token.tokenName, token)
			}

			i = j
		}
	}

	return tokens
}

// Java: SourceMapTestCase.compileAndCheck(String)
func (b *SourceMapTestCaseBase) CompileAndCheck(js string) {
	b.t.Helper()
	inputName := "testcode"
	result := b.Compile(js, inputName)
	b.CheckOfInput(inputName, js, result.GeneratedSource, result.SourceMapFileContent)
}

// Java: SourceMapTestCase.check(String, String, String, String)
func (b *SourceMapTestCaseBase) CheckOfInput(inputName string, input string, output string, sourceMapFileContent string) {
	b.t.Helper()
	b.Check([]SourceMapTestInput{{Name: inputName, Code: input}}, output, sourceMapFileContent)
}

// Java: SourceMapTestCase.check(Map<String, String>, String, String)
func (b *SourceMapTestCaseBase) Check(
	originalInputs []SourceMapTestInput, generatedSource string, sourceMapFileContent string) {
	b.t.Helper()
	b.CheckWithSupplier(originalInputs, generatedSource, sourceMapFileContent, nil)
}

// Java: SourceMapTestCase.check(Map<String, String>, String, String, SourceMapSupplier)
func (b *SourceMapTestCaseBase) CheckWithSupplier(
	originalInputs []SourceMapTestInput,
	generatedSource string,
	sourceMapFileContent string,
	supplier sourcemap.SourceMapSupplier) {
	t := b.t
	t.Helper()
	// Find all instances of the __XXX__ pattern in the original
	// source code.
	originalTokens := sourceMapTestCase_findTokensInInputs(originalInputs)

	// Find all instances of the __XXX__ pattern in the generated
	// source code.
	resultTokens := sourceMapTestCase_findTokens(generatedSource)

	// Ensure that the generated instances match via the source map
	// to the original source code.

	// Ensure the token counts match.
	if len(resultTokens.byKey) != len(originalTokens.byKey) {
		t.Fatalf("expected to have size: %d\nbut was: %d", len(originalTokens.byKey), len(resultTokens.byKey))
	}

	var reader sourcemap.SourceMapping
	jutil.Try(func() {
		reader = sourcemap.SourceMapConsumerFactory_ParseWithSupplier(sourceMapFileContent, supplier)
	}, func(e sourcemap.ISourceMapParseException) {
		panic(jutil.NewRuntimeExceptionWithMessageAndCause("unexpected exception", e))
	})

	// Map the tokens from the generated source back to the
	// input source and ensure that the map is correct.
	for _, tokenKey := range resultTokens.order {
		token := resultTokens.byKey[tokenKey]
		mapping := reader.GetMappingForLine(token.position.GetLine()+1, token.position.GetColumn()+1)

		if mapping == nil {
			t.Fatalf("expected not to be null (mapping for token %s)", token.tokenName)
		}

		// Find the associated token in the input source.
		inputToken := originalTokens.byKey[token.tokenName]
		if inputToken == nil {
			t.Fatalf("expected not to be null (input token %s)", token.tokenName)
		}
		if inputToken.inputName != mapping.GetOriginalFile() {
			t.Fatalf("expected: %s\nbut was : %s", mapping.GetOriginalFile(), inputToken.inputName)
		}

		// Ensure that the map correctly points to the token (we add 1
		// to normalize versus the Rhino line number indexing scheme).
		if inputToken.position.GetLine()+1 != mapping.GetLineNumber() {
			t.Fatalf("Checking line for token %s \nexpected: %d\nbut was : %d",
				token.tokenName, mapping.GetLineNumber(), inputToken.position.GetLine()+1)
		}

		start := inputToken.position.GetColumn() + 1
		if strings.HasPrefix(inputToken.tokenName, "STR") {
			// include the preceding quote.
			start--
		}

		if b.validateColumns {
			if mapping.GetColumnPosition() != start {
				t.Fatalf("Checking column for token %s \nexpected: %d\nbut was : %d",
					token.tokenName, start, mapping.GetColumnPosition())
			}
		}

		// Ensure that if the token name does not being with an 'STR' (meaning a
		// string) it has an original name.
		if !strings.HasPrefix(inputToken.tokenName, "STR") {
			if mapping.GetIdentifier() == "" {
				t.Fatalf("missing name for %s", inputToken.tokenName)
			}
		}

		// Ensure that if the mapping has a name, it matches the token.
		if mapping.GetIdentifier() != "" {
			if "__"+inputToken.tokenName+"__" != mapping.GetIdentifier() {
				t.Fatalf("expected: %s\nbut was : __%s__", mapping.GetIdentifier(), inputToken.tokenName)
			}
		}
	}
}

// Java: SourceMapTestCase.compile(String, String)
func (b *SourceMapTestCaseBase) Compile(js string, fileName string) *SourceMapTestCase_RunResult {
	b.t.Helper()
	return b.CompileOfTwo(js, fileName, nil, nil)
}

// CompileOfTwo compiles js1 (and js2, when js2 and fileName2 are not null)
// with source maps.
//
// Java: SourceMapTestCase.compile(String, String, String, String)
func (b *SourceMapTestCaseBase) CompileOfTwo(js1 string, fileName1 string, js2 *string, fileName2 *string) *SourceMapTestCase_RunResult {
	t := b.t
	t.Helper()
	compiler := jscomp.NewCompilerWithErrorManager(NewTestErrorManager())
	options := b.self.GetCompilerOptions()

	options.SetChecksOnly(true)

	inputs := []*jscomp.SourceFile{jscomp.SourceFile_FromCode(fileName1, js1)}

	if js2 != nil && fileName2 != nil {
		inputs =
			[]*jscomp.SourceFile{
				jscomp.SourceFile_FromCode(fileName1, js1), jscomp.SourceFile_FromCode(*fileName2, *js2)}
	}

	result := compiler.Compile(SourceMapTestCase_EXTERNS, inputs, options)

	if len(result.Errors) != 0 {
		t.Fatalf("compilation failed with errors\nexpected: []\nbut was : %s", jutil.ToJavaString(result.Errors))
	}
	if len(result.Warnings) != 0 {
		t.Fatalf("compilation failed with warnings\nexpected: []\nbut was : %s", jutil.ToJavaString(result.Warnings))
	}

	if !result.Success {
		t.Fatalf("compilation failed (other reason)")
	}
	source := compiler.ToSource()

	var sb strings.Builder
	result.SourceMap.Validate(true)
	result.SourceMap.AppendTo(&sb, jutil.StringPtr("testcode"))

	rr := &SourceMapTestCase_RunResult{}
	rr.GeneratedSource = source
	rr.SourceMap = result.SourceMap
	rr.SourceMapFileContent = sb.String()
	return rr
}

// sourceMapTestCase_uselessCodeGuard is getCompilerOptions' anonymous
// WarningsGuard: it turns CHECK_USELESS_CODE off.
type sourceMapTestCase_uselessCodeGuard struct {
	jscomp.WarningsGuardBase
}

// Java: SourceMapTestCase.getCompilerOptions().new WarningsGuard() {...}.level(JSError)
func (*sourceMapTestCase_uselessCodeGuard) Level(error *jscomp.JSError) *jscomp.CheckLevel {
	if jscomp.DiagnosticGroups_CHECK_USELESS_CODE.Matches(error) {
		return jscomp.CheckLevel_OFF.Ptr()
	}
	return nil
}

// Java: SourceMapTestCase.getCompilerOptions()
func (b *SourceMapTestCaseBase) GetCompilerOptions() *jscomp.CompilerOptions {
	options := jscomp.NewCompilerOptions()
	options.SetSourceMapOutputPath(jutil.StringPtr("testcode_source_map.out"))
	options.SetSourceMapFormat(b.self.GetSourceMapFormat())
	options.SetSourceMapDetailLevel(b.DetailLevel)
	options.SetSourceMapIncludeSourcesContent(b.SourceMapIncludeSourcesContent)
	options.AddWarningsGuard(&sourceMapTestCase_uselessCodeGuard{})
	return options
}
