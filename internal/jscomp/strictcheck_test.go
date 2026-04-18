// Port of com/google/javascript/jscomp/TypeCheckStrictTest.java
// @ closure-compiler v20260415 (subset; arithmetic / bitwise /
// comparison branches only).
//
// Verifies that the strict-mode operand validators stay silent
// under the default CompilerOptions and fire when the
// STRICT_CHECK_TYPES group is turned on. The Java test class covers
// the same expectation — TypeValidator must be gated behind the
// group; every fixture in the extracted corpus was generated with
// the default OFF level, so enabling here would otherwise regress
// the oracle match.

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/diag"
)

func TestStrictChecksOffByDefault(t *testing.T) {
	cases := []string{
		"var x = 'abc' - 1;",        // string - number
		"var x = 'abc' * 2;",        // string * number
		"var x = 'abc' & 1;",        // string & number
		"var x = 1 >>> 'abc';",      // number >>> string
		"var x = true < 'asdf';",    // mixed-type relative comparison
	}
	for _, src := range cases {
		diags, err := CheckSource("t.js", src)
		if err != nil {
			t.Fatalf("parse: %v (%q)", err, src)
		}
		for _, d := range diags {
			if d.Code == diag.JSC_TYPE_MISMATCH {
				t.Errorf("%q: unexpected default-mode TYPE_MISMATCH at %d:%d: %s",
					src, d.Line, d.Col, d.Message)
			}
		}
	}
}

func TestStrictChecksFireWhenEnabled(t *testing.T) {
	opts := NewCompilerOptions()
	opts.SetWarningLevel(GroupStrictCheckTypes, LevelError)
	cases := []string{
		"var x = 'abc' - 1;",
		"var x = 'abc' * 2;",
		"var x = 1 >>> 'abc';",
	}
	for _, src := range cases {
		diags, err := CheckSourceWithOptions("t.js", src, opts)
		if err != nil {
			t.Fatalf("parse: %v (%q)", err, src)
		}
		found := false
		for _, d := range diags {
			if d.Code == diag.JSC_TYPE_MISMATCH {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%q: expected strict JSC_TYPE_MISMATCH under strictCheckTypes", src)
		}
	}
}
