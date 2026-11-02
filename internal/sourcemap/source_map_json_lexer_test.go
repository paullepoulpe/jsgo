package sourcemap

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestSourceMapJsonScanner_StreamStringsAndNumbers(t *testing.T) {
	var lex jutil.JsonScanner
	lex.Init(`[0, -42, 2147483647, -2147483648, "hello\n\u0041\u03b1\ud83d\ude00"]`, true)
	if err := lex.BeginArray(); err != nil {
		t.Fatal(err)
	}
	wants := []string{
		"0",
		"-42",
		strconv.Itoa(math.MaxInt32),
		strconv.Itoa(math.MinInt32),
		"hello\nA\u03b1\U0001f600",
	}
	for _, want := range wants {
		got, err := lex.NextString()
		if err != nil {
			t.Fatalf("NextString() unexpected error: %v", err)
		}
		if got != want {
			t.Fatalf("NextString() = %q, want %q", got, want)
		}
	}
	if err := lex.EndArray(); err != nil {
		t.Fatal(err)
	}

	for _, bad := range []string{`"\u12g4"`, `"\u12"`} {
		lex.Init(bad, true)
		if _, err := lex.NextString(); err == nil {
			t.Fatalf("NextString(%s) succeeded, want error", bad)
		}
	}
}

func TestSourceMapObjectParser_ZeroExtraAllocsLargeMappings(t *testing.T) {
	largeMappings := strings.Repeat("AAAA,CAAC;", 1024)
	srcMap := `{"version":3,"file":"out.js","sourceRoot":"","sources":["in.js"],"names":["foo"],"mappings":"` + largeMappings + `","unknown":{"nested":[1,2,"skip\\nme"]}}`
	allocs := testing.AllocsPerRun(100, func() {
		obj := SourceMapObjectParser_Parse(srcMap)
		if obj.mappings.text != largeMappings {
			t.Fatal("mappings mismatch")
		}
	})
	if allocs > 4 {
		t.Fatalf("SourceMapObjectParser_Parse allocated %.1f objects, want <= 4 (SourceMapObject + sources + names slices only)", allocs)
	}
}
