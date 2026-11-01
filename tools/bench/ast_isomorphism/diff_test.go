package main

import (
	"strings"
	"testing"
)

func diffOrFail(t *testing.T, a, b string, opts DiffOptions) *DiffResult {
	t.Helper()
	res, err := DiffBuilds(a, b, opts)
	if err != nil {
		t.Fatalf("DiffBuilds: %v", err)
	}
	return res
}

func changedText(res *DiffResult, b string) []string {
	var out []string
	for _, c := range res.Changes {
		if c.B != nil {
			out = append(out, c.Kind+":"+b[c.B[0]:c.B[1]])
		} else {
			out = append(out, c.Kind+":-")
		}
	}
	return out
}

func TestDiffBuilds(t *testing.T) {
	t.Run("renaming_only", func(t *testing.T) {
		a := "var g=1;function f(a,b){return a+g*b}window.x=f(2,3);"
		b := "var h=1;function k(x,y){return x+h*y}window.x=k(2,3);"
		res := diffOrFail(t, a, b, DiffOptions{})
		if !res.Isomorphic {
			t.Fatalf("expected isomorphic, got %v", changedText(res, b))
		}
	})

	t.Run("swapped_operands_are_a_change", func(t *testing.T) {
		// Both operands are free in f: a swap must not pass as a renaming.
		a := "var p=1,q=2;function f(){return p-q}function g(){return p*q}window.x=f()+g();"
		b := "var p=1,q=2;function f(){return q-p}function g(){return p*q}window.x=f()+g();"
		res := diffOrFail(t, a, b, DiffOptions{})
		if res.Isomorphic || res.ChangedB != 1 || res.UnitsB[0].Name != "f" {
			t.Fatalf("expected one changed function f, got %+v %v", res.UnitsB, changedText(res, b))
		}
	})

	t.Run("inserted_statement_is_localised", func(t *testing.T) {
		a := "function f(a){return a+1}function g(a){var b=a*2;return b}window.x=f(1)+g(2);"
		b := "function f(a){return a+1}function g(a){var b=a*2;console.log(b);return b}window.x=f(1)+g(2);"
		res := diffOrFail(t, a, b, DiffOptions{})
		got := changedText(res, b)
		if len(got) != 1 || got[0] != "insert:console.log(b);" || res.UnitsB[res.Changes[0].UnitB].Name != "g" {
			t.Fatalf("expected one insert in g, got %v", got)
		}
	})

	t.Run("deleted_statement", func(t *testing.T) {
		a := "function g(a){var b=a*2;sideEffect();return b}window.x=g(2);"
		b := "function g(a){var b=a*2;return b}window.x=g(2);"
		res := diffOrFail(t, a, b, DiffOptions{})
		if len(res.Changes) != 1 || res.Changes[0].Kind != "delete" || a[res.Changes[0].A[0]:res.Changes[0].A[1]] != "sideEffect();" {
			t.Fatalf("expected the deletion of sideEffect(), got %+v", res.Changes)
		}
	})

	t.Run("renamed_properties", func(t *testing.T) {
		a := "function C(){this.a1=1;this.b1=2}C.prototype.m=function(){return this.a1+this.b1};window.x=(new C).m();"
		b := "function C(){this.b1=1;this.c1=2}C.prototype.m=function(){return this.b1+this.c1};window.x=(new C).m();"
		if res := diffOrFail(t, a, b, DiffOptions{}); res.Isomorphic {
			t.Fatalf("exact properties: expected a difference")
		}
		res := diffOrFail(t, a, b, DiffOptions{RenamedProps: true, Fixed: map[string]bool{"prototype": true}})
		if !res.Isomorphic {
			t.Fatalf("renamed properties: expected isomorphic, got %v", changedText(res, b))
		}
		// A fixed (extern) name never pairs with another name.
		c := "function C(){this.a1=1;this.b1=2}C.prototype.m=function(){return this.a1+this.length};window.x=(new C).m();"
		res = diffOrFail(t, a, c, DiffOptions{RenamedProps: true, Fixed: map[string]bool{"prototype": true, "length": true}})
		if res.Isomorphic {
			t.Fatalf("expected b1 vs length to be a change")
		}
	})

	t.Run("literal_change", func(t *testing.T) {
		a := "window.x=function(a){return a?1:-0}(2);"
		b := "window.x=function(a){return a?1:0}(2);"
		res := diffOrFail(t, a, b, DiffOptions{})
		if res.Isomorphic {
			t.Fatalf("expected -0 vs 0 to be a change")
		}
	})

	t.Run("long_lines_have_exact_positions", func(t *testing.T) {
		var sb strings.Builder
		sb.WriteString("var s='a,b;c',r=/[,;]/g,t=`x${1,2}y,z`;window.x=[")
		for i := 0; i < 2000; i++ {
			sb.WriteString("1,")
		}
		a := sb.String() + "2];window.y=function(){return 3}();"
		b := sb.String() + "2];window.y=function(){return 4}();"
		res := diffOrFail(t, a, b, DiffOptions{})
		if !res.Exact {
			t.Fatalf("expected exact positions")
		}
		got := changedText(res, b)
		if len(got) != 1 || got[0] != "node:4" {
			t.Fatalf("expected the literal 4, got %v", got)
		}
	})
}
