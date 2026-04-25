package jscomp

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

func TestDebugFixture1(t *testing.T) {
	src := "/** @constructor */var T = function() { this.x = ''; };\n/** @type {number} */ T.prototype.x = 0;"

	script, err := parser.Parse(src)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	reg := jstype.NewRegistry()
	creator := NewTypedScopeCreator(reg)
	scope := creator.CreateScope(script)

	typ := scope.GetDeclaredType("T")
	fmt.Printf("T type: %T\n", typ)
	if fn, ok := typ.(*jstype.FunctionType); ok {
		if fn.Instance != nil {
			xType := fn.Instance.GetPropertyType("x")
			xInferred := fn.Instance.IsPropertyInferred("x")
			fmt.Printf("T.Instance.x type: %T / %v, inferred=%v\n", xType, xType, xInferred)
		} else {
			fmt.Printf("T.Instance is nil\n")
		}
	} else {
		fmt.Printf("T is not a FunctionType\n")
	}
}
