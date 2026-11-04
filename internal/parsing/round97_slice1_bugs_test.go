/*
 * Copyright 2026 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parsing

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestRound97TTLNonNameCalleeNoPanic(t *testing.T) {
	jt := newJsDocInfoParserTest(t)

	// IIFE and member-call expressions have a non-NAME callee and must report
	// invalid TTL warnings instead of panicking in GetCallName.
	jt.parse(
		`@template T := (() => "number")() =: */`,
		"Bad type annotation. Invalid type transformation expression")
	jt.parse(
		`@template T := obj.union('a', 'b') =: */`,
		"Bad type annotation. Invalid type transformation expression")
	jt.parse(
		`@template T := cond(obj.eq('a', 'b'), 'number', 'string') =: */`,
		"Bad type annotation. Invalid boolean expression",
		"Bad type annotation. Invalid expression inside conditional")
	jt.parse(
		`@template T := union((() => 'a')(), 'b') =: */`,
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside union type")
}

func TestRound97TTLNativeTypeExprUnparseableString(t *testing.T) {
	jt := newJsDocInfoParserTest(t)

	// Empty or unparseable JSDoc type strings return NilNode from
	// JsDocInfoParser_ParseTypeString and must emit TTL warnings instead of
	// dereferencing NilNode.
	jt.parse(
		`@template T := typeExpr("") =: */`,
		"Bad type annotation. Invalid native type expression",
		"Bad type annotation. Invalid expression inside typeExpr")
	jt.parse(
		`@template T := typeExpr("<<<") =: */`,
		"Bad type annotation. Invalid native type expression",
		"Bad type annotation. Invalid expression inside typeExpr")
}

func TestRound97TTLPredicateAtTypePositionJavaParity(t *testing.T) {
	// Preserves Java v20261028 oracle parity (pinned by irf/ttl1.js): boolean
	// predicate keywords at type position throw IllegalStateException("Invalid
	// type transformation expression").
	for _, ttl := range []string{
		`eq("number", "string")`,
		`sub("number", "string")`,
		`union(eq("a", "b"), "number")`,
	} {
		func() {
			defer func() {
				r := recover()
				if _, ok := r.(jutil.IIllegalStateException); !ok {
					t.Fatalf("ParseTypeTransformation(%q) panic = %T (%v), want IIllegalStateException", ttl, r, r)
				}
			}()
			src := rhino.NewSimpleSourceFile("ttl_test.js", rhino.StaticSourceFile_SourceKind_STRONG)
			p := NewTypeTransformationParser(ttl, src, rhino.ErrorReporter_NULL_INSTANCE, 1, 0)
			p.ParseTypeTransformation()
		}()
	}
}

func TestRound97IRFactoryLiteralsPolish(t *testing.T) {
	pt := newParserTest(t)

	// String escapes: hex, 4-digit unicode, braced code point, line continuation, and UTF-8 / surrogate pairs.
	script := pt.parseWarning(
		"var s = '\\x41\\u0042\\u{1F600}\\\r\nZ'; var q = 'hé😀';",
		"String continuations are not recommended. See https://google.github.io/styleguide/jsguide.html#features-strings-no-line-continuations")
	firstInit, _ := ast.MustNameDeclaration(script.GetFirstChild()).MustFirstDeclarator().MustName().Init()
	sVal := ast.MustStringLit(firstInit.Node()).Value()
	if want := "AB😀Z"; sVal != want {
		t.Fatalf("escaped string value = %q, want %q", sVal, want)
	}
	secondInit, _ := ast.MustNameDeclaration(script.GetSecondChild()).MustFirstDeclarator().MustName().Init()
	qVal := ast.MustStringLit(secondInit.Node()).Value()
	if want := "hé😀"; qVal != want {
		t.Fatalf("utf-8/surrogate string value = %q, want %q", qVal, want)
	}

	// Numeric & BigInt literal normalization with bases and numeric separators.
	numScript := pt.parse("const a = [0b1010_0001, 0o7_7, 0x1a_2b, 1_000.5e1, 0b10_10n, 0o7_7n, 0x1_fn, 0n, 1_234n];")
	arrInit, _ := ast.MustNameDeclaration(numScript.GetFirstChild()).MustFirstDeclarator().MustName().Init()
	arr := ast.MustArrayLit(arrInit.Node())
	var elems []rhino.Node
	for elem := range arr.Elements() {
		elems = append(elems, elem.Node())
	}
	if len(elems) != 9 {
		t.Fatalf("got %d array elements, want 9", len(elems))
	}
	wantNums := []float64{0b10100001, 0o77, 0x1a2b, 10005}
	for i, want := range wantNums {
		if got := elems[i].GetDouble(); got != want {
			t.Errorf("elems[%d].GetDouble() = %v, want %v", i, got, want)
		}
	}
	wantBigInts := []string{"10", "63", "31", "0", "1234"}
	for i, want := range wantBigInts {
		if got := elems[4+i].GetBigInt().String(); got != want {
			t.Errorf("elems[%d].GetBigInt() = %q, want %q", 4+i, got, want)
		}
	}
}
