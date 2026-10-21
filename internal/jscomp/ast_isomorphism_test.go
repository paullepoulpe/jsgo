package jscomp

import (
	"strings"
	"testing"
)

func TestVerifyAstAlphaIsomorphism(t *testing.T) {
	t.Run("identical", func(t *testing.T) {
		code := "var a=1;function f(b,c){return a+b*c}window.x=f(2,3);\n"
		stats, err := VerifyAstAlphaIsomorphism(code, code)
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
		if stats.Functions != 1 || stats.Returns != 1 {
			t.Fatalf("unexpected stats: %+v", stats)
		}
	})

	t.Run("local_var_and_closure_alpha_bijection", func(t *testing.T) {
		codeA := "var g=1;window.f=function(a,b){var c=a+g,d=b*2;return function(e){return c+d+e}};\n"
		codeB := "var h=1;window.f=function(x,y){var m=x+h,n=y*2;return function(z){return m+n+z}};\n"
		stats, err := VerifyAstAlphaIsomorphism(codeA, codeB)
		if err != nil {
			t.Fatalf("expected alpha-isomorphic ASTs to pass, got: %v", err)
		}
		if stats.RenamedVarsDiff != 6 {
			t.Fatalf("expected 6 renamed vars, got %d", stats.RenamedVarsDiff)
		}
	})

	t.Run("label_alpha_bijection", func(t *testing.T) {
		codeA := "function f(a){a:{for(var b=0;b<10;b++)if(b===a)break a;return 1}return 0}\n"
		codeB := "function f(x){b:{for(var y=0;y<10;y++)if(y===x)break b;return 1}return 0}\n"
		if _, err := VerifyAstAlphaIsomorphism(codeA, codeB); err != nil {
			t.Fatalf("expected label alpha-bijection to pass, got: %v", err)
		}
	})

	t.Run("reject_bijection_violation", func(t *testing.T) {
		codeA := "function f(a,b){return a+a}\n"
		codeB := "function f(a,b){return a+b}\n"
		_, err := VerifyAstAlphaIsomorphism(codeA, codeB)
		if err == nil || !strings.Contains(err.Error(), "bijection violation") {
			t.Fatalf("expected bijection violation error, got: %v", err)
		}
	})

	t.Run("reject_scope_depth_mismatch", func(t *testing.T) {
		codeA := "var g=1;function f(a){return g}\n"
		codeB := "var g=1;function f(a){return a}\n"
		_, err := VerifyAstAlphaIsomorphism(codeA, codeB)
		if err == nil || !strings.Contains(err.Error(), "lexical scope depth mismatch") {
			t.Fatalf("expected scope depth mismatch error, got: %v", err)
		}
	})

	t.Run("reject_extern_mismatch", func(t *testing.T) {
		codeA := "window.foo=1;\n"
		codeB := "document.foo=1;\n"
		_, err := VerifyAstAlphaIsomorphism(codeA, codeB)
		if err == nil || !strings.Contains(err.Error(), "undeclared global/extern NAME mismatch") {
			t.Fatalf("expected extern mismatch error, got: %v", err)
		}
	})

	t.Run("reject_dead_assignment_diff", func(t *testing.T) {
		codeA := "function f(){var a=0;return 1}\n"
		codeB := "function f(){var a=0;a=2;return 1}\n"
		_, err := VerifyAstAlphaIsomorphism(codeA, codeB)
		if err == nil || !strings.Contains(err.Error(), "child count mismatch") {
			t.Fatalf("expected child count mismatch error, got: %v", err)
		}
	})
}
