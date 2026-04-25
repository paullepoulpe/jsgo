package harness_test

import (
	"fmt"
	"testing"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

func TestDebugFixture(t *testing.T) {
	// Test case 1: simple inferred write on constructor
	src1 := `/** @constructor */ function Foo() {}
function f() {
  var x = new Foo();
  x.onload = function() {};
}`
	opts := jscomp.NewCompilerOptions()
	opts.SetWarningLevel(jscomp.GroupStrictCheckTypes, jscomp.LevelWarning)
	fmt.Println("=== Test 1: x.onload = function(){} ===")
	diags, _ := jscomp.CheckSourceWithOptions("test.js", src1, opts)
	fmt.Printf("diags=%d\n", len(diags))
	for _, d := range diags {
		fmt.Printf("  %d:%d [%s] %s\n", d.Line, d.Col, d.Code, d.Message)
	}

	// Test case 2: Object read 
	src2 := `/** @constructor */
function Foo() {
  this.prop = 123;
}
/** @param {!Object} x */
function f(x) {
  return x.prop;
}`
	fmt.Println("=== Test 2: !Object read ===")
	diags, _ = jscomp.CheckSourceWithOptions("test.js", src2, opts)
	fmt.Printf("diags=%d\n", len(diags))
	for _, d := range diags {
		fmt.Printf("  %d:%d [%s] %s\n", d.Line, d.Col, d.Code, d.Message)
	}
}
