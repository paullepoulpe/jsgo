package bundle

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestCompilerBasedTransformerUntransformedSkipsSourceMapAndTransformedEmitsValidV3Map(t *testing.T) {
	supplier := NewCompilerSupplier()

	// 1. Untransformed ES5 input: Transformed must be false and SourceMap must be empty.
	es5Code := "var answer = 42;\nfunction greet(name) {\n  return 'hello ' + name;\n}\n"
	untransformedRes := supplier.Compile("dir/es5_passthrough.js", es5Code)
	if untransformedRes.Transformed {
		t.Fatalf("expected Transformed == false for ES5 pass-through input, got true")
	}
	if untransformedRes.SourceMap != "" {
		t.Fatalf("expected empty SourceMap for untransformed input, got %q", untransformedRes.SourceMap)
	}
	if len(untransformedRes.Errors) != 0 {
		t.Fatalf("unexpected errors on ES5 input: %v", untransformedRes.Errors)
	}

	// Verify Transpiler.Transform returns the original Source pointer (and exact raw code) untouched when untransformed.
	tr := NewTranspiler(supplier, "es6_runtime")
	origSrc := NewSource(SourceConfig{Path: jutil.Path_Of("dir/es5_passthrough.js"), Code: es5Code})
	gotSrc := tr.Transform(origSrc)
	if gotSrc != origSrc || gotSrc.Code() != es5Code || gotSrc.SourceMap() != "" {
		t.Fatalf("expected Transform to return original *Source untouched on untransformed input, got %+v", gotSrc)
	}

	// 2. Transformed ES6 input: Transformed must be true and SourceMap must be a valid V3 source map.
	es6Code := "const add = (a, b) => a + b;\n"
	transformedRes := supplier.Compile("dir/es6_input.js", es6Code)
	if !transformedRes.Transformed {
		t.Fatalf("expected Transformed == true for ES6 arrow function input")
	}
	if transformedRes.SourceMap == "" {
		t.Fatalf("expected non-empty V3 SourceMap for transformed ES6 input")
	}
	if !strings.Contains(transformedRes.Source, "var add = function(a, b)") {
		t.Fatalf("unexpected transpiled Source: %q", transformedRes.Source)
	}

	var v3 struct {
		Version        int      `json:"version"`
		File           string   `json:"file"`
		Sources        []string `json:"sources"`
		SourcesContent []string `json:"sourcesContent"`
		Mappings       string   `json:"mappings"`
	}
	if err := json.Unmarshal([]byte(transformedRes.SourceMap), &v3); err != nil {
		t.Fatalf("invalid JSON in V3 SourceMap: %v\nraw: %s", err, transformedRes.SourceMap)
	}
	if v3.Version != 3 {
		t.Fatalf("SourceMap version = %d, want 3", v3.Version)
	}
	if v3.File != "dir/es6_input.js" {
		t.Fatalf("SourceMap file = %q, want %q", v3.File, "dir/es6_input.js")
	}
	if len(v3.Sources) != 1 || v3.Sources[0] != "dir/es6_input.js" {
		t.Fatalf("SourceMap sources = %v, want [\"dir/es6_input.js\"]", v3.Sources)
	}
	if len(v3.SourcesContent) != 1 || v3.SourcesContent[0] != es6Code {
		t.Fatalf("SourceMap sourcesContent = %q, want [%q]", v3.SourcesContent, es6Code)
	}
	if v3.Mappings == "" {
		t.Fatalf("expected non-empty VLQ mappings in V3 SourceMap")
	}

	// 3. AlwaysTransformed supplier (Coverage): even ES5 input must emit SourceMap.
	covSupplier := NewCoverageCompilerSupplier()
	covRes := covSupplier.Compile("dir/cov.js", "var x = 1;\n")
	if !covRes.Transformed || covRes.SourceMap == "" {
		t.Fatalf("AlwaysTransformed supplier: Transformed=%v, SourceMap len=%d; want true and non-empty", covRes.Transformed, len(covRes.SourceMap))
	}
}
