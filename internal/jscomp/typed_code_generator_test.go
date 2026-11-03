package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestTypedCodeGenerator_FindMethodOwnerGuards(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitCompilerOptionsIfTesting()
	registry := compiler.GetTypeRegistry()
	gen := NewTypedCodeGenerator(nil, compiler.GetOptions(), registry)

	// 1. Non-ObjectType global type on prototype assignment (e.g., primitive/union typedef)
	// must return nil cleanly instead of panicking on an unchecked ObjectType assertion.
	registry.DeclareType(nil, "NonObjectGlobal", registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	nonObjRoot := compiler.ParseTestCode("NonObjectGlobal.prototype.m = function() {};")
	nonObjFn := ast.MustScript(nonObjRoot).MustFirstItem().MustExprStmt().Expr().MustAssign().Value().Node()
	if got := gen.FindMethodOwner(nonObjFn); got != nil {
		t.Fatalf("FindMethodOwner(non-object global) = %v, want nil", got)
	}

	// 2. Valid constructor instance ObjectType on prototype assignment returns instance type.
	fooName := "FooCtor"
	var ctor rhino.FunctionType
	func() {
		closer := registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		ctor = registry.CreateConstructorType(fooName, rhino.NewNode(rhino.Token_FUNCTION), registry.CreateParameters(), nil, nil, false)
	}()
	registry.DeclareType(nil, fooName, ctor.GetInstanceType())
	ctorRoot := compiler.ParseTestCode("FooCtor.prototype.m = function() {};")
	ctorFn := ast.MustScript(ctorRoot).MustFirstItem().MustExprStmt().Expr().MustAssign().Value().Node()
	if got := gen.FindMethodOwner(ctorFn); got != ctor.GetInstanceType() {
		t.Fatalf("FindMethodOwner(FooCtor.prototype.m) = %v, want %v", got, ctor.GetInstanceType())
	}

	// 3. Parent CLASS node with nil JSType must return nil instead of dereferencing nil.
	clsRoot := compiler.ParseTestCode("class UnannotatedClass {}")
	cls := ast.MustScript(clsRoot).MustFirstItem().MustClass()
	clsNode := cls.Node()
	clsBody := cls.Body().Node()
	if got := gen.FindMethodOwner(clsBody); got != nil {
		t.Fatalf("FindMethodOwner(untyped class child) = %v, want nil", got)
	}

	// 4. Parent CLASS node with constructor FunctionType returns its instance type.
	clsCtorName := "AnnotatedClass"
	var clsCtor rhino.FunctionType
	func() {
		closer := registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		clsCtor = registry.CreateConstructorType(clsCtorName, clsNode, registry.CreateParameters(), nil, nil, false)
	}()
	clsNode.SetJSType(clsCtor)
	if got := gen.FindMethodOwner(clsBody); got != clsCtor.GetInstanceType() {
		t.Fatalf("FindMethodOwner(typed class child) = %v, want %v", got, clsCtor.GetInstanceType())
	}
}
