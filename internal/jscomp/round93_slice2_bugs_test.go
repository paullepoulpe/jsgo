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
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestRound93FastModeExpectSuperTypeNoRace(t *testing.T) {
	// 1. Verify IsParallelReadOnly() and ExpectSuperType safety when parallelReadOnly is toggled
	// (including when subObject has nil GetConstructor() or an uncached constructor).
	for _, parallelReadOnly := range []bool{false, true} {
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetCheckTypes(true)
		opts.SetWarningLevel(DiagnosticGroups_CHECK_TYPES, CheckLevel_WARNING)
		compiler.InitOptions(opts)
		reg := compiler.GetTypeRegistry()
		closer := reg.GetResolver().OpenForDefinition()

		baseACtor := reg.CreateConstructorType("BaseA", rhino.NilNode, nil, nil, nil, false)
		baseBCtor := reg.CreateConstructorType("BaseB", rhino.NilNode, nil, nil, nil, false)
		subCtor := reg.CreateConstructorType("Sub", rhino.NilNode, nil, nil, nil, false)
		baseA := baseACtor.GetInstanceType()
		baseB := baseBCtor.GetInstanceType()
		intermProto := reg.CreateObjectType(base.SomeString("Sub.prototype"), baseA)
		subObj := reg.CreateObjectType(base.SomeString("SubInstance"), intermProto)
		closer.Close()

		errNode := rhino.IR_Name("Sub")
		reg.SetParallelReadOnly(parallelReadOnly)
		if got := reg.IsParallelReadOnly(); got != parallelReadOnly {
			t.Fatalf("IsParallelReadOnly() = %v, want %v", got, parallelReadOnly)
		}
		compiler.GetTypeValidator().ExpectSuperType(errNode, baseB, subObj)
		compiler.GetTypeValidator().ExpectSuperType(errNode, baseB, subCtor.GetInstanceType())
		reg.SetParallelReadOnly(false)
		if got := reg.IsParallelReadOnly(); got {
			t.Fatal("IsParallelReadOnly() remained true after SetParallelReadOnly(false)")
		}
		if len(compiler.GetWarnings()) == 0 {
			t.Fatalf("parallelReadOnly=%v: expected superclass mismatch warning", parallelReadOnly)
		}
	}

	// 2. Concurrent --fast_mode TypeCheck across scripts with mismatched goog.inherits / @extends
	// and concurrent readers of Sub and its prototype chain (must be race-free under -race).
	for iter := range 4 {
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetCheckTypes(true)
		opts.SetClosurePass(true)
		opts.SetFastMode(true)
		opts.SetNumParallelThreads(4)
		opts.SetWarningLevel(DiagnosticGroups_CHECK_TYPES, CheckLevel_WARNING)

		inputs := []*SourceFile{
			SourceFile_FromCode("s0.js", `
/** @const */ var goog = {};
goog.inherits = function(childCtor, parentCtor) {};
/** @constructor */ function BaseA() {}
/** @return {number} */ BaseA.prototype.getA = function() { return 1; };
/** @constructor */ function BaseB() {}
/** @return {string} */ BaseB.prototype.getB = function() { return 'b'; };
/** @constructor @extends {BaseA} */ function Sub() {}
/** @constructor */ function SubNoExtends() {}
`),
			SourceFile_FromCode("s1.js", "goog.inherits(Sub, BaseB); goog.inherits(SubNoExtends, BaseB);\n"),
			SourceFile_FromCode("s2.js", "goog.inherits(Sub, BaseB); goog.inherits(SubNoExtends, BaseB);\n"),
			SourceFile_FromCode("s3.js", "var inst1 = new Sub(); var v1 = inst1.getA(); var p1 = Sub.prototype; var u1 = new SubNoExtends();\n"),
			SourceFile_FromCode("s4.js", "goog.inherits(Sub, BaseB); goog.inherits(SubNoExtends, BaseB);\n"),
			SourceFile_FromCode("s5.js", "/** @type {!BaseA} */ var asBase = new Sub(); var v2 = asBase.getA(); var p2 = SubNoExtends.prototype;\n"),
			SourceFile_FromCode("s6.js", "goog.inherits(Sub, BaseB); goog.inherits(SubNoExtends, BaseB);\n"),
			SourceFile_FromCode("s7.js", "var inst2 = new Sub(); var ctorRef = Sub; var ctorRef2 = SubNoExtends;\n"),
		}
		compiler.Init([]*SourceFile{SourceFile_FromCode("externs.js", "var Object; Function.prototype.call; Function.prototype.apply; Function.prototype.bind;")}, inputs, opts)
		root := compiler.ParseInputs()
		if compiler.GetErrorCount() != 0 {
			t.Fatalf("iter %d parse errors: %v", iter, compiler.GetErrors())
		}
		externsRoot := root.GetFirstChild()
		jsRoot := externsRoot.GetNext()
		reg := compiler.GetTypeRegistry()
		rai := NewSemanticReverseAbstractInterpreter(reg)
		tc := NewTypeCheck(compiler, rai, reg)
		tc.ProcessForTesting(externsRoot, jsRoot)
		if _, workers, ok := tc.canCheckInParallel(jsRoot); !ok || workers < 2 {
			t.Fatalf("iter %d: expected parallel TypeCheck with >=2 workers, got ok=%v workers=%d", iter, ok, workers)
		}

		foundMismatch := false
		for _, w := range compiler.GetWarnings() {
			if strings.Contains(w.GetDescription(), "mismatch in declaration of superclass type") {
				foundMismatch = true
				break
			}
		}
		if !foundMismatch {
			t.Fatalf("iter %d: expected superclass mismatch warning, got warnings=%v", iter, compiler.GetWarnings())
		}
	}
}

func TestRound93TypeTransformationTemplateTypeOfEdgeCases(t *testing.T) {
	setupTTL := func() (*Compiler, *TypeTransformation, *OrderedStringMap[rhino.JSType], *OrderedStringMap[string]) {
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		compiler.InitOptions(opts)
		reg := compiler.GetTypeRegistry()
		extra := compiler.ParseTestCode("var dummy;")
		var scopeAny any = NewTypedScopeCreator(compiler).
			CreateScope(rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(extra)), nil)
		scope := scopeAny.(*TypedScope)
		var typeVars OrderedStringMap[rhino.JSType]
		arrType := reg.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
		numType := reg.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
		typeVars.Put("ARRNUM", reg.CreateTemplatizedType(arrType, numType))
		var nameVars OrderedStringMap[string]
		return compiler, NewTypeTransformation(compiler, scope), &typeVars, &nameVars
	}

	// 1. Negative index `templateTypeOf(ARRNUM, -1)` must report INDEX_OUTOFBOUNDS
	// and return UNKNOWN_TYPE without panicking.
	{
		compiler, tt, typeVars, nameVars := setupTTL()
		reg := compiler.GetTypeRegistry()
		negCall := compiler.ParseTestCode("templateTypeOf(ARRNUM, -1);").GetFirstFirstChild()
		closer := reg.GetResolver().OpenForDefinition()
		got := tt.Eval(negCall, typeVars, nameVars)
		closer.Close()

		want := reg.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		if !got.Equals(want) {
			t.Fatalf("Eval(templateTypeOf(ARRNUM, -1)) = %v, want %v", got, want)
		}
		warnings := compiler.GetWarnings()
		if len(warnings) != 1 || warnings[0].GetType() != TypeTransformation_INDEX_OUTOFBOUNDS {
			t.Fatalf("Eval(templateTypeOf(ARRNUM, -1)) warnings = %v, want 1 TypeTransformation_INDEX_OUTOFBOUNDS", warnings)
		}
		if !strings.Contains(warnings[0].GetDescription(), "found -1") {
			t.Fatalf("Eval(templateTypeOf(ARRNUM, -1)) description = %q, want 'found -1'", warnings[0].GetDescription())
		}
	}

	// 2. Malformed direct AST calls (<2 args, non-number 2nd arg, unary NEG(NAME), huge float) must not panic.
	malformedCases := []struct {
		name      string
		call      rhino.Node
		wantBound bool
	}{
		{
			name: "zero_args",
			call: rhino.IR_Call(rhino.IR_Name("templateTypeOf")),
		},
		{
			name: "one_arg",
			call: rhino.IR_Call(rhino.IR_Name("templateTypeOf"), rhino.IR_Name("ARRNUM")),
		},
		{
			name: "non_number_index",
			call: rhino.IR_Call(rhino.IR_Name("templateTypeOf"), rhino.IR_Name("ARRNUM"), rhino.IR_String("bad")),
		},
		{
			name: "unary_neg_non_number",
			call: rhino.IR_Call(rhino.IR_Name("templateTypeOf"), rhino.IR_Name("ARRNUM"), rhino.IR_Neg(rhino.IR_Name("bad"))),
		},
		{
			name:      "huge_float_index",
			call:      rhino.IR_Call(rhino.IR_Name("templateTypeOf"), rhino.IR_Name("ARRNUM"), rhino.IR_Number(1e20)),
			wantBound: true,
		},
	}

	for _, tc := range malformedCases {
		t.Run(tc.name, func(t *testing.T) {
			compiler, tt, typeVars, nameVars := setupTTL()
			reg := compiler.GetTypeRegistry()
			closer := reg.GetResolver().OpenForDefinition()
			got := tt.Eval(tc.call, typeVars, nameVars)
			closer.Close()

			want := reg.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
			if !got.Equals(want) {
				t.Fatalf("Eval(%s) = %v, want %v", tc.name, got, want)
			}
			warnings := compiler.GetWarnings()
			if tc.wantBound {
				if len(warnings) != 1 || warnings[0].GetType() != TypeTransformation_INDEX_OUTOFBOUNDS {
					t.Fatalf("Eval(%s) warnings = %v, want TypeTransformation_INDEX_OUTOFBOUNDS", tc.name, warnings)
				}
			}
		})
	}
}

func TestRound93AstValidatorAccessorAndComputedPatternFixes(t *testing.T) {
	validateScript := func(src string, mutate func(script rhino.Node)) []string {
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		compiler.InitOptions(opts)
		script := NewCompilerInput(SourceFile_FromCode("test.js", src)).GetAstRoot(compiler)
		if compiler.GetErrorCount() != 0 {
			t.Fatalf("parse errors for %q: %v", src, compiler.GetErrors())
		}
		if mutate != nil {
			mutate(script)
		}
		var violations []string
		validator := NewAstValidatorWithViolationFunc(
			compiler,
			func(message string, _ rhino.Node) {
				violations = append(violations, message)
			},
			/* validateScriptFeatures= */ false,
			/* shouldValidateRequiredInlinings= */ false)
		validator.SetTypeValidationMode(AstValidator_TypeInfoValidation_NONE)
		validator.ValidateScript(script)
		return violations
	}

	// 1. Malformed getters/setters with childless (NilNode) and non-FUNCTION children must not panic.
	accessorMutations := []struct {
		name   string
		mutate func(objLit rhino.Node)
	}{
		{
			name: "childless_getter_def",
			mutate: func(objLit rhino.Node) {
				objLit.GetFirstChild().RemoveChildren()
			},
		},
		{
			name: "getter_with_number_child",
			mutate: func(objLit rhino.Node) {
				getter := objLit.GetFirstChild()
				getter.GetFirstChild().ReplaceWith(rhino.IR_Number(1).Srcref(getter))
			},
		},
		{
			name: "setter_with_empty_child",
			mutate: func(objLit rhino.Node) {
				setter := objLit.GetSecondChild()
				setter.GetFirstChild().ReplaceWith(rhino.IR_Empty().Srcref(setter))
			},
		},
		{
			name: "getter_with_3child_non_function",
			mutate: func(objLit rhino.Node) {
				getter := objLit.GetFirstChild()
				forNode := rhino.NewNodeWith3Children(
					rhino.Token_FOR_OF,
					rhino.IR_Empty(),
					rhino.IR_Empty(),
					rhino.IR_Block()).SrcrefTree(getter)
				getter.GetFirstChild().ReplaceWith(forNode)
			},
		},
		{
			name: "childless_computed_getter",
			mutate: func(objLit rhino.Node) {
				compGet := objLit.GetChildAtIndex(2)
				compGet.RemoveChildren()
			},
		},
		{
			name: "setter_with_number_child",
			mutate: func(objLit rhino.Node) {
				setter := objLit.GetSecondChild()
				setter.GetFirstChild().ReplaceWith(rhino.IR_Number(1).Srcref(setter))
			},
		},
	}

	for _, tc := range accessorMutations {
		t.Run(tc.name, func(t *testing.T) {
			violations := validateScript(
				"({ get g() { return 1; }, set s(x) {}, get [k]() { return 1; }, set [k](x) {} });",
				func(script rhino.Node) {
					objLit := script.GetFirstFirstChild()
					tc.mutate(objLit)
				})
			if len(violations) == 0 {
				t.Fatalf("%s: expected at least one AST validation violation, got none", tc.name)
			}
		})
	}

	// Direct computed getter/setter validator helpers with childless / non-FUNCTION nodes must not panic.
	{
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		compiler.InitOptions(opts)
		var violations []string
		validator := NewAstValidatorWithViolationFunc(
			compiler,
			func(message string, _ rhino.Node) {
				violations = append(violations, message)
			},
			false,
			false)
		validator.SetTypeValidationMode(AstValidator_TypeInfoValidation_NONE)

		compSet := rhino.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, rhino.IR_Name("k"), rhino.IR_Number(1))
		compSet.PutBooleanProp(rhino.Node_COMPUTED_PROP_SETTER, true)
		validator.validateObjectLitComputedPropSetKey(ast.MustComputedProp(compSet))
		if len(violations) == 0 {
			t.Fatal("expected violation for computed setter with NUMBER child")
		}

		violations = nil
		compGetEmpty := rhino.NewNode(rhino.Token_COMPUTED_PROP)
		compGetEmpty.PutBooleanProp(rhino.Node_COMPUTED_PROP_GETTER, true)
		validator.validateObjectLitComputedPropGetKey(ast.MustComputedProp(compGetEmpty))
		if len(violations) == 0 {
			t.Fatal("expected violation for childless computed getter")
		}
	}

	// 2. `const {[k]: o.p} = obj;` must report "Invalid child for CONST node", whereas
	// `({[k]: o.p} = obj);` is valid in an ASSIGN destructuring target.
	{
		// Build `const {[k]: o.p} = obj;` by replacing `x` with `o.p` in `const {[k]: x} = obj;`
		// because the parser rejects member expressions in lexical destructuring bindings.
		violations := validateScript("const {[k]: x} = obj;", func(script rhino.Node) {
			constNode := script.GetFirstChild()
			lhs := constNode.GetFirstChild()
			objPat := lhs.GetFirstChild()
			compProp := objPat.GetFirstChild()
			target := compProp.GetLastChild()
			getProp := rhino.IR_Getprop(rhino.IR_Name("o"), "p").SrcrefTree(target)
			target.ReplaceWith(getProp)
		})
		if !slices.Contains(violations, "Invalid child for CONST node") {
			t.Fatalf("expected 'Invalid child for CONST node', got %v", violations)
		}
	}
	{
		violations := validateScript("({[k]: o.p} = obj);", nil)
		if len(violations) != 0 {
			t.Fatalf("expected valid assignment destructuring ({[k]: o.p} = obj), got violations: %v", violations)
		}
	}

	// 3. Malformed COMPUTED_PROP inside OBJECT_PATTERN with 0 children and 1 child must not panic.
	for _, keepChildren := range []int{0, 1} {
		t.Run(fmt.Sprintf("computed_pattern_children_%d", keepChildren), func(t *testing.T) {
			violations := validateScript("let {[k]: x} = obj;", func(script rhino.Node) {
				compProp := script.GetFirstChild().GetFirstChild().GetFirstChild().GetFirstChild()
				if keepChildren == 0 {
					compProp.RemoveChildren()
				} else {
					compProp.GetLastChild().Detach()
				}
			})
			if len(violations) == 0 {
				t.Fatalf("keepChildren=%d: expected arity violation, got none", keepChildren)
			}
		})
	}
}

func TestRound93BundledGoogModuleMalformedTraversal(t *testing.T) {
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetClosurePass(true)
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetWarningLevel(DiagnosticGroups_MISSING_REQUIRE, CheckLevel_WARNING)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "")},
		[]*SourceFile{SourceFile_FromCode("bundle.js", `
goog.loadModule(function(exports) {
  goog.module();
  var x = new UndeclaredThing();
  return exports;
});
goog.loadModule(function(exports) {
  goog.module(123);
  return exports;
});
goog.loadModule(function(exports) {
  goog.module('valid.bundled.Mod');
  exports.val = 42;
  return exports;
});
`)},
		opts)
	root := compiler.ParseInputs()
	externsRoot := root.GetFirstChild()
	jsRoot := externsRoot.GetNext()

	NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(externsRoot, jsRoot)
	moduleMap := compiler.GetModuleMetadataMap()
	if moduleMap == nil {
		t.Fatal("expected non-nil ModuleMetadataMap")
	}
	if moduleMap.ModuleByGoogNamespace("valid.bundled.Mod") == nil {
		t.Fatal("expected valid.bundled.Mod to be recorded in ModuleMetadataMap")
	}

	// None of the AbstractModuleCallback passes may panic on malformed bundled goog.module() calls.
	NewCheckMissingRequires(compiler, moduleMap).Process(externsRoot, jsRoot)
	NewCheckClosureImports(compiler, moduleMap).Process(externsRoot, jsRoot)
	NewClosureCheckModule(compiler, moduleMap).Process(externsRoot, jsRoot)
}
