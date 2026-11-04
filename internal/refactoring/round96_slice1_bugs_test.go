package refactoring

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestRound96SuggestedFixDeleteAfterDestructuringOrComputedProp(t *testing.T) {
	t.Run("DestructuringLhsSibling", func(t *testing.T) {
		input := "const {a} = foo(), b = 1;"
		compiler := suggestedFixTest_getCompiler(input)
		root := suggestedFixTest_compileToScriptRoot(compiler)
		lastDecl := root.GetFirstChild().GetLastChild()

		fix := NewSuggestedFix_Builder().Delete(lastDecl).Build()
		wantReplacement := CodeReplacement_Create(jutil.StringIndexOf(input, ", b = 1"), jutil.StringLength(", b = 1"), "")
		SuggestedFixes_AssertReplacement(fix, wantReplacement)

		got := ApplyToCode([]*SuggestedFix{fix}, map[string]string{"test": input})["test"]
		if got != "const {a} = foo();" {
			t.Fatalf("ApplyToCode = %q, want %q", got, "const {a} = foo();")
		}
	})

	t.Run("ComputedPropObjectLitSibling", func(t *testing.T) {
		input := "const o = {[k]: 1, b: 2};"
		compiler := suggestedFixTest_getCompiler(input)
		root := suggestedFixTest_compileToScriptRoot(compiler)
		objLit := root.GetFirstFirstChild().GetFirstChild()
		lastProp := objLit.GetLastChild()

		fix := NewSuggestedFix_Builder().Delete(lastProp).Build()
		wantReplacement := CodeReplacement_Create(jutil.StringIndexOf(input, ", b: 2"), jutil.StringLength(", b: 2"), "")
		SuggestedFixes_AssertReplacement(fix, wantReplacement)

		got := ApplyToCode([]*SuggestedFix{fix}, map[string]string{"test": input})["test"]
		if got != "const o = {[k]: 1};" {
			t.Fatalf("ApplyToCode = %q, want %q", got, "const o = {[k]: 1};")
		}
	})

	t.Run("ComputedPropObjectPatternSibling", func(t *testing.T) {
		input := "const {[k]: a, b} = obj;"
		compiler := suggestedFixTest_getCompiler(input)
		root := suggestedFixTest_compileToScriptRoot(compiler)
		objPat := root.GetFirstFirstChild().GetFirstChild()
		lastProp := objPat.GetLastChild()

		fix := NewSuggestedFix_Builder().Delete(lastProp).Build()
		got := ApplyToCode([]*SuggestedFix{fix}, map[string]string{"test": input})["test"]
		if got != "const {[k]: a} = obj;" {
			t.Fatalf("ApplyToCode = %q, want %q", got, "const {[k]: a} = obj;")
		}
	})
}

func TestRound96MatchersConstructorWithNameZeroChildNode(t *testing.T) {
	input := "/** @constructor */ debugger;"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	stmt := root.GetFirstChild()
	if stmt.HasChildren() {
		t.Fatalf("expected 0-child statement, got %v", stmt)
	}
	if Matchers_ConstructorWithName("Foo").Matches(stmt, nil) {
		t.Fatalf("expected Matchers_ConstructorWithName on 0-child debugger node to be false")
	}

	b := rhino.NewJSDocInfo_Builder()
	b.RecordConstructor()
	empty := rhino.IR_Empty()
	empty.SetJSDocInfo(b.Build())
	if Matchers_ConstructorWithName("Foo").Matches(empty, nil) {
		t.Fatalf("expected Matchers_ConstructorWithName on 0-child EMPTY node to be false")
	}
}
