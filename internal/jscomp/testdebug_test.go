package jscomp

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

func TestDebugLocallyInferred(t *testing.T) {
	src := "/** @constructor */ function F() {}\n/** @type {number} */ F.prototype.z;\n/** @const */ var x = {}; /** @type {F} */ x.y;\nfunction f() { x.y.z = 'abc'; }\n/** @param {number} x */ function g(x) {}\nfunction h() { g(x.y.z); }"

	opts := NewCompilerOptions()
	diags, _ := CheckSourceWithOptions("test.js", src, opts)
	fmt.Printf("Got %d diagnostics\n", len(diags))
	for _, d := range diags {
		fmt.Printf("  %d:%d %s %s\n", d.Line, d.Col, d.Code, d.Message)
	}
	_ = jstype.Unknown
}
