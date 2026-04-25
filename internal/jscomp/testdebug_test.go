package jscomp

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

func TestDebugLocallyInferred(t *testing.T) {
	src := "/** @constructor */ function F() {}\n/** @type {number} */ F.prototype.z;\n/** @const */ var x = {}; /** @type {F} */ x.y;\nfunction f() { x.y.z = 'abc'; }"

	script, err := parser.Parse(src)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	reg := jstype.NewRegistry()
	creator := NewTypedScopeCreator(reg)
	scope := creator.CreateScope(script)

	// Check F.Instance.z
	fType := scope.GetDeclaredType("F")
	fmt.Printf("F type: %T\n", fType)
	if fn, ok := fType.(*jstype.FunctionType); ok && fn.Instance != nil {
		zType := fn.Instance.GetPropertyType("z")
		zInferred := fn.Instance.IsPropertyInferred("z")
		fmt.Printf("F.Instance.z type: %T / %v, inferred=%v\n", zType, zType, zInferred)
	}

	// Check x.y namespace prop type
	xyType := scope.LookupNamespaceProp("x", "y")
	fmt.Printf("x.y type: %T / %v\n", xyType, xyType)
}
