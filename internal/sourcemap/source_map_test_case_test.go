// Port of test/com/google/debugging/sourcemap/SourceMapTestCase.java @ closure-compiler v20260415

package sourcemap

import (
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// sourceMapTestCase_DetailLevel stands for jscomp's SourceMap.DetailLevel,
// which is not ported yet.
type sourceMapTestCase_DetailLevel string

const (
	detailLevel_ALL     sourceMapTestCase_DetailLevel = "ALL"
	detailLevel_SYMBOLS sourceMapTestCase_DetailLevel = "SYMBOLS"
)

// sourceMapTestCase is the abstract SourceMapTestCase: compiles JS with
// source maps enabled and checks the map. Compiling needs jscomp.Compiler,
// which is not ported yet, so compile skips the calling test.
type sourceMapTestCase struct {
	t                              *testing.T
	validateColumns                bool
	detailLevel                    sourceMapTestCase_DetailLevel
	sourceMapIncludeSourcesContent bool
}

func newSourceMapTestCase(t *testing.T) *sourceMapTestCase {
	tc := &sourceMapTestCase{t: t, validateColumns: true}
	tc.setUp()
	return tc
}

func (tc *sourceMapTestCase) disableColumnValidation() {
	tc.validateColumns = false
}

type sourceMapTestCase_RunResult struct {
	generatedSource string
	sourceMap       interface {
		AppendTo(out io.StringWriter, name *string)
	}
	sourceMapFileContent string
}

type sourceMapTestCase_Token struct {
	tokenName string
	inputName string
	position  *FilePosition
}

func (tc *sourceMapTestCase) setUp() {
	tc.detailLevel = detailLevel_ALL
}

// checkSourceMap creates a source map for the given JS code and asserts it
// is equal to the expected golden map.
func (tc *sourceMapTestCase) checkSourceMap(js string, expectedMap *jutil.JsonObject) {
	tc.checkSourceMapWithFileName("testcode", js, expectedMap)
}

func (tc *sourceMapTestCase) checkSourceMapWithFileName(fileName, js string, expectedMap *jutil.JsonObject) {
	t := tc.t
	result := tc.compile(js, fileName)
	// We round-trip the expectation to coerce any of its types in the same
	// way as the actual map.
	roundTrippedExpectedMap := jutil.NewGson().FromJsonObject(gsonToJson(expectedMap))

	resultMap := jutil.NewGson().FromJsonObject(result.sourceMapFileContent)
	if !jsonMapEquals(resultMap, roundTrippedExpectedMap) {
		t.Fatalf("%s\nmap %s\nwant %s", result.generatedSource, resultMap, roundTrippedExpectedMap)
	}

	resultMapMap := jutil.NewGson().FromJsonObject(tc.getSourceMap(result))
	if !jsonMapEquals(resultMapMap, roundTrippedExpectedMap) {
		t.Fatalf("%s\nmap %s\nwant %s", result.generatedSource, resultMapMap, roundTrippedExpectedMap)
	}
}

// jsonMapEquals is the equality of the Map<String, ?> that Gson
// deserializes: unordered objects, ordered lists, numbers as doubles.
func jsonMapEquals(a, b jutil.JsonElement) bool {
	switch {
	case a.IsJsonNull() || b.IsJsonNull():
		return a.IsJsonNull() && b.IsJsonNull()
	case a.IsJsonObject():
		if !b.IsJsonObject() || a.GetAsJsonObject().Size() != b.GetAsJsonObject().Size() {
			return false
		}
		for k, v := range a.GetAsJsonObject().EntrySet().All() {
			w := b.GetAsJsonObject().Get(k)
			if w == nil || !jsonMapEquals(v, w) {
				return false
			}
		}
		return true
	case a.IsJsonArray():
		if !b.IsJsonArray() || a.GetAsJsonArray().Size() != b.GetAsJsonArray().Size() {
			return false
		}
		for i := 0; i < a.GetAsJsonArray().Size(); i++ {
			if !jsonMapEquals(a.GetAsJsonArray().Get(i), b.GetAsJsonArray().Get(i)) {
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

func (tc *sourceMapTestCase) getSourceMap(result *sourceMapTestCase_RunResult) string {
	var sb strings.Builder
	result.sourceMap.AppendTo(&sb, jutil.StringPtr("testcode"))
	return sb.String()
}

// findTokensInInputs finds the all the __XX__ tokens in the given
// JavaScript string.
func findTokensInInputs(inputs *jutil.LinkedHashMap[string, string]) *jutil.LinkedHashMap[string, *sourceMapTestCase_Token] {
	tokens := jutil.NewLinkedHashMap[string, *sourceMapTestCase_Token]()

	for k, v := range inputs.All() {
		findTokensInto(tokens, k, v)
	}

	return tokens
}

// findTokens finds the all the __XX__ tokens in the given JavaScript string.
func findTokens(src string) *jutil.LinkedHashMap[string, *sourceMapTestCase_Token] {
	tokens := jutil.NewLinkedHashMap[string, *sourceMapTestCase_Token]()

	findTokensInto(tokens, "", src)

	return tokens
}

// findTokensInto finds the all the __XX__ tokens in the given JavaScript
// string.
func findTokensInto(
	tokens *jutil.LinkedHashMap[string, *sourceMapTestCase_Token], inputName string, jsString string) *jutil.LinkedHashMap[string, *sourceMapTestCase_Token] {

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
				token := &sourceMapTestCase_Token{
					jutil.StringFromUTF16(tokenName), inputName, NewFilePosition(currentLine, currentPosition)}
				tokens.Put(token.tokenName, token)
			}

			i = j
		}
	}

	return tokens
}

func (tc *sourceMapTestCase) compileAndCheck(js string) {
	inputName := "testcode"
	result := tc.compile(js, inputName)
	tc.check(inputName, js, result.generatedSource, result.sourceMapFileContent)
}

func (tc *sourceMapTestCase) check(inputName, input, output, sourceMapFileContent string) {
	inputMap := jutil.NewLinkedHashMap[string, string]()
	inputMap.Put(inputName, input)
	tc.checkInputs(inputMap, output, sourceMapFileContent, nil)
}

func (tc *sourceMapTestCase) checkInputs(
	originalInputs *jutil.LinkedHashMap[string, string],
	generatedSource string,
	sourceMapFileContent string,
	supplier SourceMapSupplier) {
	t := tc.t
	t.Helper()
	// Find all instances of the __XXX__ pattern in the original source code.
	originalTokens := findTokensInInputs(originalInputs)

	// Find all instances of the __XXX__ pattern in the generated source
	// code.
	resultTokens := findTokens(generatedSource)

	// Ensure that the generated instances match via the source map to the
	// original source code.

	// Ensure the token counts match.
	if resultTokens.Size() != originalTokens.Size() {
		t.Fatalf("result has %d tokens, original %d", resultTokens.Size(), originalTokens.Size())
	}

	var reader SourceMapping
	jutil.Try(func() {
		reader = SourceMapConsumerFactory_ParseWithSupplier(sourceMapFileContent, supplier)
	}, func(e ISourceMapParseException) {
		t.Fatalf("unexpected exception: %v", e)
	})

	// Map the tokens from the generated source back to the input source and
	// ensure that the map is correct.
	for token := range resultTokens.Values() {
		mapping := reader.GetMappingForLine(token.position.GetLine()+1, token.position.GetColumn()+1)

		if mapping == nil {
			t.Fatalf("no mapping for token %s", token.tokenName)
		}

		// Find the associated token in the input source.
		inputToken := originalTokens.Get(token.tokenName)
		if inputToken == nil {
			t.Fatalf("no input token %s", token.tokenName)
		}
		if inputToken.inputName != mapping.GetOriginalFile() {
			t.Fatalf("token %s: file %s, want %s", token.tokenName, mapping.GetOriginalFile(), inputToken.inputName)
		}

		// Ensure that the map correctly points to the token (we add 1 to
		// normalize versus the Rhino line number indexing scheme).
		if inputToken.position.GetLine()+1 != mapping.GetLineNumber() {
			t.Fatalf("Checking line for token %s: %d, want %d", token.tokenName, mapping.GetLineNumber(), inputToken.position.GetLine()+1)
		}

		start := inputToken.position.GetColumn() + 1
		if strings.HasPrefix(inputToken.tokenName, "STR") {
			// include the preceding quote.
			start--
		}

		if tc.validateColumns {
			if mapping.GetColumnPosition() != start {
				t.Fatalf("Checking column for token %s: %d, want %d", token.tokenName, mapping.GetColumnPosition(), start)
			}
		}

		// Ensure that if the token name does not being with an 'STR'
		// (meaning a string) it has an original name.
		if !strings.HasPrefix(inputToken.tokenName, "STR") {
			if mapping.GetIdentifier() == "" {
				t.Fatalf("missing name for %s", inputToken.tokenName)
			}
		}

		// Ensure that if the mapping has a name, it matches the token.
		if mapping.GetIdentifier() != "" {
			if "__"+inputToken.tokenName+"__" != mapping.GetIdentifier() {
				t.Fatalf("name %s, want __%s__", mapping.GetIdentifier(), inputToken.tokenName)
			}
		}
	}
}

func (tc *sourceMapTestCase) compile(js, fileName string) *sourceMapTestCase_RunResult {
	return tc.compile2(js, fileName, nil, nil)
}

// compile2 compiles one or two inputs with source maps. It needs
// jscomp.Compiler, CompilerOptions (source map output path, format, detail
// level, sources content) and a warnings guard that turns off
// CHECK_USELESS_CODE, none of which are ported yet.
func (tc *sourceMapTestCase) compile2(js1, fileName1 string, js2, fileName2 *string) *sourceMapTestCase_RunResult {
	tc.t.Skip("needs jscomp.Compiler (M3)")
	return nil
}
