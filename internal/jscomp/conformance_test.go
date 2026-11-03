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

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestConformance_CyclicConstStringInference(t *testing.T) {
	compiler := NewCompilerWithErrorManager(NewBlackHoleErrorManager())
	opts := NewCompilerOptions()
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "")},
		[]*SourceFile{SourceFile_FromCode("in.js", `
const a = b;
const b = a;
/** @const */ var selfCycle = selfCycle;
const ns = {};
/** @const */ ns.x = ns.y;
/** @const */ ns.y = ns.x;
const okBase = 'hello';
const okAlias = okBase;
use(a, selfCycle, ns.x, okAlias);
`)},
		opts)
	compiler.ParseInputs()
	registry := compiler.GetTypeRegistry()
	stringType := registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)

	script := compiler.GetJsRoot().GetFirstChild()
	scope := NewSyntacticScopeCreator(compiler).CreateScope(script, nil).(*Scope)
	gns := NewGlobalNamespace(compiler, script)
	gnsSupplier := func() *GlobalNamespace { return gns }

	// Attach string JSType to ns.x and ns.y GETPROPs so the global-namespace branch runs.
	var attachStringTypes func(rhino.Node)
	attachStringTypes = func(n rhino.Node) {
		if n.IsGetProp() && (n.MatchesQualifiedName("ns.x") || n.MatchesQualifiedName("ns.y")) {
			n.SetJSType(stringType)
		}
		for c := range n.Children() {
			attachStringTypes(c)
		}
	}
	attachStringTypes(script)

	useCall := script.GetLastChild().GetFirstChild()
	argA := useCall.GetSecondChild()
	argSelf := argA.GetNext()
	argNsX := argSelf.GetNext()
	argOk := argNsX.GetNext()

	if val, ok := inferConformanceStringValue(scope, argA, gnsSupplier); ok || val != "" {
		t.Fatalf("cyclic const a/b = (%q, %v), want (\"\", false)", val, ok)
	}
	if val, ok := inferConformanceStringValue(scope, argSelf, gnsSupplier); ok || val != "" {
		t.Fatalf("self-cyclic var = (%q, %v), want (\"\", false)", val, ok)
	}
	if val, ok := inferConformanceStringValue(scope, argNsX, gnsSupplier); ok || val != "" {
		t.Fatalf("cyclic global ns.x/ns.y = (%q, %v), want (\"\", false)", val, ok)
	}
	if got := ConformanceRules_ConformanceUtil_InferStringValue(scope, argA, gnsSupplier); got != nil {
		t.Fatalf("InferStringValue(cyclic a) = %v, want nil", *got)
	}
	nsXName := "ns.x"
	if got := ConformanceRules_ConformanceUtil_GetValueFromGlobalName(scope, gns, &nsXName); got != nil {
		t.Fatalf("GetValueFromGlobalName(ns.x) = %v, want nil", *got)
	}

	val, ok := inferConformanceStringValue(scope, argOk, gnsSupplier)
	if !ok || val != "hello" {
		t.Fatalf("non-cyclic const alias = (%q, %v), want (\"hello\", true)", val, ok)
	}

	// Non-cyclic inference must perform zero heap allocations.
	allocs := testing.AllocsPerRun(100, func() {
		v, ok := inferConformanceStringValue(scope, argOk, gnsSupplier)
		if !ok || v != "hello" {
			t.Fatal("unexpected result")
		}
	})
	if allocs != 0 {
		t.Fatalf("inferConformanceStringValue non-cyclic allocated %.2f objects/run, want 0", allocs)
	}
}
