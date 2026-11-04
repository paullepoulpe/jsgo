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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type round96ConstNodeSnapshot struct {
	inExterns bool
	scriptIdx int
	token     rhino.Token
	name      string
	line      int
	charno    int
	declared  bool
	inferred  bool
}

func round96CollectConstSnapshots(externsRoot, jsRoot rhino.Node) []round96ConstNodeSnapshot {
	var out []round96ConstNodeSnapshot
	var walk func(n rhino.Node, inExterns bool, scriptIdx int)
	walk = func(n rhino.Node, inExterns bool, scriptIdx int) {
		if n.IsName() || n.IsImportStar() {
			out = append(out, round96ConstNodeSnapshot{
				inExterns: inExterns,
				scriptIdx: scriptIdx,
				token:     n.GetToken(),
				name:      n.GetString(),
				line:      n.GetLineno(),
				charno:    n.GetCharno(),
				declared:  n.IsDeclaredConstantVar(),
				inferred:  n.IsInferredConstantVar(),
			})
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c, inExterns, scriptIdx)
		}
	}
	sIdx := 0
	for s := externsRoot.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
		walk(s, true, sIdx)
		sIdx++
	}
	sIdx = 0
	for s := jsRoot.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
		walk(s, false, sIdx)
		sIdx++
	}
	return out
}

func round96RunInferConsts(t *testing.T, externsSrc string, scripts []string, numThreads int) ([]round96ConstNodeSnapshot, *Compiler) {
	t.Helper()
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	opts.SetNumParallelThreads(numThreads)
	externs := []*SourceFile{
		SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS+externsSrc),
	}
	inputs := make([]*SourceFile, len(scripts))
	for i, code := range scripts {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("script_%02d.js", i), code)
	}
	compiler.Init(externs, inputs, opts)
	compiler.ParseInputs()
	if compiler.HasErrors() {
		t.Fatalf("unexpected parse errors: %v", formatJSErrors(compiler.GetErrors()))
	}
	NewInferConsts(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	return round96CollectConstSnapshots(compiler.GetExternsRoot(), compiler.GetJsRoot()), compiler
}

type round96TrackingScopeCreator struct {
	delegate *SyntacticScopeCreator
	scopes   []*Scope
}

func (s *round96TrackingScopeCreator) CreateScope(n rhino.Node, parent AbstractScope) AbstractScope {
	var parentScope *Scope
	if parent != nil {
		parentScope = parent.(*Scope)
	}
	sc := s.delegate.CreateScopeOfScope(n, parentScope)
	s.scopes = append(s.scopes, sc)
	return sc
}

func TestRound96_InferConsts_O1ScopeFrameLookupAndReset(t *testing.T) {
	src := `
function outer(paramConst, paramMutated) {
  paramMutated = 1;
  hoistedReader();
  var earlyLocal = 10;
  var lateLocal = 20;
  function hoistedReader() {
    return earlyLocal + paramConst;
  }
  {
    const blockConst = 1;
    let blockLet = 2;
    {
      let innerMutated = 3;
      for (const item of [1, 2, 3]) {
        try {
          if (item > 0) {
            const deepConst = blockConst + blockLet + item;
            innerMutated += deepConst;
          }
        } catch (err) {
          const catchLocal = err;
          if (catchLocal) {}
        }
      }
    }
  }
  return lateLocal + paramMutated;
}
`
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	opts.SetNumParallelThreads(1)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "")},
		[]*SourceFile{SourceFile_FromCode("in.js", src)},
		opts,
	)
	compiler.ParseInputs()
	if compiler.HasErrors() {
		t.Fatalf("unexpected parse errors: %v", formatJSErrors(compiler.GetErrors()))
	}

	infer := NewInferConsts(compiler)
	sc := NewSyntacticScopeCreator(compiler)
	tracker := &round96TrackingScopeCreator{delegate: sc}
	collector := newInferConstsCollector(infer, compiler, sc)
	NodeTraversal_Builder{
		compiler:     compiler,
		callback:     collector,
		scopeCreator: tracker,
	}.Traverse(compiler.GetJsRoot())
	for i := range collector.globalStates {
		if collector.globalStates[i].v != nil {
			infer.considerVarState(&collector.globalStates[i])
		}
	}

	if len(tracker.scopes) < 6 {
		t.Fatalf("expected nested scopes to be created, got %d", len(tracker.scopes))
	}
	for i, sc := range tracker.scopes {
		if sc.inferFrameIdx != -1 {
			t.Fatalf("scope[%d] (%v) leaked inferFrameIdx=%d after ExitScope", i, sc.GetRootNode(), sc.inferFrameIdx)
		}
	}

	snaps := round96CollectConstSnapshots(compiler.GetExternsRoot(), compiler.GetJsRoot())
	byName := make(map[string]round96ConstNodeSnapshot)
	for _, sn := range snaps {
		prev := byName[sn.name]
		prev.name = sn.name
		prev.declared = prev.declared || sn.declared
		prev.inferred = prev.inferred || sn.inferred
		byName[sn.name] = prev
	}
	expectInferred := map[string]bool{
		"outer":        true,
		"paramConst":   true,
		"paramMutated": false,
		"earlyLocal":   false, // read via hoistedReader() before `var earlyLocal = 10;`
		"lateLocal":    true,
		"blockConst":   true,
		"blockLet":     true,
		"innerMutated": false,
		"item":         false, // `for (const item of ...)` is declared const, not inferred const
		"deepConst":    true,
		"catchLocal":   true,
	}
	if !byName["item"].declared {
		t.Errorf("var \"item\" declared=false, want true")
	}
	for name, want := range expectInferred {
		sn, ok := byName[name]
		if !ok {
			t.Fatalf("missing snapshot for %q", name)
		}
		if sn.inferred != want {
			t.Errorf("var %q inferred=%v, want %v", name, sn.inferred, want)
		}
	}
}

func TestRound96_InferConsts_ParallelVsSequentialParity(t *testing.T) {
	externsSrc := `
/** @const */ var EXT_DECLARED_CONST;
const EXT_ES6_CONST = 1;
var EXT_MUTABLE;
function use(...args) {}
`
	scripts := []string{
		// Script 0: initializes alreadyInitGlobalVar BEFORE cross-script hoisted jump chain!
		`var alreadyInitGlobalVar = 100;
hoistedStep1();
function localWork0(a, [b = () => c], c) {
  const localConst0 = a + c;
  let localMut0 = 1;
  localMut0++;
  return localConst0 + localMut0 + b();
}`,
		// Script 1: mutual recursion trigger + global vars mutated later or in unreferenced hoisted fn
		`var mutatedAcrossScripts = 1;
var mutatedInUnreferencedFn = 2;
/** @const */ var globalDocConstNoInit;
/** @const */ var globalDocConstMutated = 1;
cycleA();`,
		// Script 2: second hop of hoisted chain reading globals declared in later scripts
		`function hoistedStep2() {
  use(alreadyInitGlobalVar, earlyGlobalVar, globalLetInLaterScript, globalConstInLaterScript, GlobalClassInLaterScript);
}
function localWork2(items) {
  var earlyLocal2Ref = helper2();
  var localVarAfterCall = 5;
  function helper2() { return localVarAfterCall; }
  const stableLocal2 = 9;
  for (const x of items) { use(x, stableLocal2, earlyLocal2Ref); }
}`,
		// Script 3: mutual recursion A -> B across scripts
		`function cycleA() {
  return cycleB();
}`,
		// Script 4: hoisted function NEVER called at top-level that mutates a global var
		`function neverCalledHoistedMutator() {
  mutatedInUnreferencedFn++;
}
{
  const blockScoped4 = 42;
  use(blockScoped4);
}`,
		// Script 5: first hop of cross-script hoisted chain (0 -> 5 -> 2)
		`function hoistedStep1() {
  hoistedStep2();
}`,
		// Script 6: global class referenced by hoistedStep2 in Script 2
		`class GlobalClassInLaterScript {
  method(p) {
    const mConst = p;
    let mLet = 0;
    for (var loopIdx = 0; loopIdx < 2; loopIdx++) {
      mLet += loopIdx;
    }
    return mConst + mLet;
  }
}`,
		// Script 7: earlyGlobalVar initialized AFTER Script 0's hoisted chain ran + cycleB reading afterCycleVar
		`var earlyGlobalVar = 200;
var safeGlobalVar = 300;
function cycleB() {
  if (EXT_MUTABLE) {
    return cycleA();
  }
  return afterCycleVar;
}`,
		// Script 8: mutates globals from Script 1, declares global let referenced by hoistedStep2
		`mutatedAcrossScripts = 99;
globalDocConstMutated = 88;
let globalLetInLaterScript = 400;
use(safeGlobalVar, EXT_DECLARED_CONST, EXT_ES6_CONST);`,
		// Script 9: global const + afterCycleVar initialized after Script 1 called cycleA()
		`const globalConstInLaterScript = 500;
var afterCycleVar = 600;
function localDestructuring9(obj) {
  const {x, y: [z = 1]} = obj;
  let reassigned9 = x;
  reassigned9 ||= z;
  return x + z + reassigned9;
}`,
	}

	wantSnaps, _ := round96RunInferConsts(t, externsSrc, scripts, 1)
	for _, threads := range []int{2, 4, 8} {
		gotSnaps, _ := round96RunInferConsts(t, externsSrc, scripts, threads)
		if !slices.Equal(gotSnaps, wantSnaps) {
			t.Fatalf("InferConsts mismatch at threads=%d (len got=%d want=%d)", threads, len(gotSnaps), len(wantSnaps))
		}
	}

	lookup := func(inExterns bool, name string) (round96ConstNodeSnapshot, bool) {
		var merged round96ConstNodeSnapshot
		found := false
		for _, sn := range wantSnaps {
			if sn.inExterns == inExterns && sn.name == name {
				found = true
				merged.name = name
				merged.declared = merged.declared || sn.declared
				merged.inferred = merged.inferred || sn.inferred
			}
		}
		return merged, found
	}

	checkExtern := map[string]bool{
		"EXT_DECLARED_CONST": true,
		"EXT_ES6_CONST":      true,
		"EXT_MUTABLE":        false,
	}
	for name, wantDecl := range checkExtern {
		sn, ok := lookup(true, name)
		if !ok {
			t.Fatalf("missing extern snapshot %q", name)
		}
		if sn.declared != wantDecl {
			t.Errorf("extern %q declared=%v, want %v", name, sn.declared, wantDecl)
		}
	}

	checkGlobalInferred := map[string]bool{
		"alreadyInitGlobalVar":     true,
		"earlyGlobalVar":           false,
		"safeGlobalVar":            true,
		"afterCycleVar":            false,
		"mutatedAcrossScripts":     false,
		"mutatedInUnreferencedFn":  false,
		"globalDocConstNoInit":     false,
		"globalDocConstMutated":    false,
		"globalLetInLaterScript":   true,
		"globalConstInLaterScript": true,
		"GlobalClassInLaterScript": true,
		"hoistedStep1":             true,
		"hoistedStep2":             true,
	}
	for name, wantInf := range checkGlobalInferred {
		sn, ok := lookup(false, name)
		if !ok {
			t.Fatalf("missing JS snapshot %q", name)
		}
		if sn.inferred != wantInf {
			t.Errorf("global %q inferred=%v, want %v", name, sn.inferred, wantInf)
		}
	}
}

func TestRound96_LintChecks_18PassParallelCloneAndParity(t *testing.T) {
	c0 := NewCompiler()
	opts0 := NewCompilerOptions()
	opts0.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	c0.InitOptions(opts0)

	combined, ok := defaultPass_lintChecks.Create(c0).(*CombinedCompilerPass)
	if !ok {
		t.Fatalf("expected *CombinedCompilerPass from defaultPass_lintChecks")
	}
	if len(combined.callbacks) != 18 {
		t.Fatalf("expected 18 lint check callbacks, got %d", len(combined.callbacks))
	}
	if !combined.canCloneForParallel() {
		for i, w := range combined.callbacks {
			if _, cloneable := w.callback.(parallelCloneableCallback); !cloneable {
				t.Errorf("lint callback[%d] %T does not implement parallelCloneableCallback", i, w.callback)
			}
		}
		t.Fatalf("expected defaultPass_lintChecks.canCloneForParallel() == true")
	}

	scripts := []string{
		// Script 0: static @private property on OrphanCtor WITHOUT OrphanCtor defined in Script 0.
		// Must NOT leak OrphanCtor from Script 4 when workers reuse Check*PrivateProperties!
		`var OrphanCtor = {};
/** @private */ OrphanCtor.staticProp_ = 1;`,
		// Script 1: class DupClass with unmodified @private prop and unused @private prop.
		`class DupClass {
  constructor() {
    /** @private */ this.neverModified_ = 1;
    /** @private */ this.neverRead_ = 2;
  }
  getVal() {
    return this.neverModified_;
  }
}`,
		// Script 2: static @private prop on class that is never read nor modified.
		`class StaticHolder {
  constructor() {}
}
/** @private */ StaticHolder.unusedStatic_ = 10;`,
		// Script 3: empty statement + useless block lint checks alongside clean class.
		`class CleanClass {
  constructor() {
    /** @private @const */ this.ok_ = 1;
  }
  read() { return this.ok_; }
}
function extraSemi() {};`,
		// Script 4: defines constructor OrphanCtor (must stay isolated to Script 4!).
		`/** @constructor */
function OrphanCtor() {
  /** @private */ this.mutatedAndRead_ = 1;
}
OrphanCtor.prototype.bump = function() {
  this.mutatedAndRead_ = 2;
  return this.mutatedAndRead_;
};`,
		// Script 5: @interface with @private prop (should not report missing @const on @interface).
		`/** @interface */
function MyIface() {}
/** @private {number} */
MyIface.prototype.ifaceProp_;`,
		// Script 6: reuses name DupClass in another script, mutating and reading its private prop.
		`/** @constructor */
function DupClass() {
  /** @private */ this.neverModified_ = 1;
}
DupClass.prototype.mutate = function() {
  this.neverModified_ = 2;
  return this.neverModified_;
};`,
		// Script 7: another class with missing @const private property.
		`class Widget7 {
  constructor() {
    /** @private */ this.fixed_ = 'abc';
  }
  label() { return this.fixed_; }
}`,
	}

	runLint := func(numThreads int) ([]string, []string) {
		comp := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		opts.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
		opts.SetNumParallelThreads(numThreads)
		inputs := make([]*SourceFile, len(scripts))
		for i, s := range scripts {
			inputs[i] = SourceFile_FromCode(fmt.Sprintf("lint_%02d.js", i), s)
		}
		comp.Init([]*SourceFile{SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS)}, inputs, opts)
		comp.ParseInputs()
		if comp.HasErrors() {
			t.Fatalf("unexpected parse errors: %v", formatJSErrors(comp.GetErrors()))
		}
		defaultPass_lintChecks.Create(comp).Process(comp.GetExternsRoot(), comp.GetJsRoot())
		return formatJSErrors(comp.GetWarnings()), formatJSErrors(comp.GetErrors())
	}

	wantWarns, wantErrs := runLint(1)
	if len(wantErrs) != 0 {
		t.Fatalf("unexpected errors: %v", wantErrs)
	}
	joinedWarns := strings.Join(wantWarns, "\n")
	if !strings.Contains(joinedWarns, CheckConstPrivateProperties_MISSING_CONST_PROPERTY.Key()) {
		t.Fatalf("expected %s warning in %v", CheckConstPrivateProperties_MISSING_CONST_PROPERTY.Key(), wantWarns)
	}
	if !strings.Contains(joinedWarns, CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY.Key()) {
		t.Fatalf("expected %s warning in %v", CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY.Key(), wantWarns)
	}
	if strings.Contains(joinedWarns, "staticProp_") {
		t.Fatalf("unexpected warning on non-constructor OrphanCtor.staticProp_ in Script 0: %v", wantWarns)
	}

	for _, threads := range []int{2, 4, 8} {
		gotWarns, gotErrs := runLint(threads)
		if !slices.Equal(gotErrs, wantErrs) {
			t.Fatalf("threads=%d lint errors mismatch:\ngot:  %v\nwant: %v", threads, gotErrs, wantErrs)
		}
		if !slices.Equal(gotWarns, wantWarns) {
			t.Fatalf("threads=%d lint warnings mismatch:\ngot:  %v\nwant: %v", threads, gotWarns, wantWarns)
		}
	}
}

func TestRound96_ImplicitNullabilityCheck_ZeroAllocAndParallel(t *testing.T) {
	comp := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	opts.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_WARNING)
	comp.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS)},
		[]*SourceFile{SourceFile_FromCode("in.js", "var plain = 1 + 2;")},
		opts,
	)
	comp.ParseInputs()

	analyzerCombined := defaultPass_analyzerChecks.Create(comp).(*CombinedCompilerPass)
	if !analyzerCombined.canCloneForParallel() {
		t.Fatalf("expected defaultPass_analyzerChecks.canCloneForParallel() == true")
	}

	chk := NewImplicitNullabilityCheck(comp)
	trav := (NodeTraversal_Builder{compiler: comp, callback: chk}).Build()
	script := comp.GetJsRoot().GetFirstChild()
	varNode := script.GetFirstChild()

	allocsNilFind := testing.AllocsPerRun(100, func() {
		res := ImplicitNullabilityCheck_FindImplicitNullabilityResults(nil, trav)
		if len(res) != 0 {
			t.Fatal("expected empty slice")
		}
	})
	if allocsNilFind != 0 {
		t.Fatalf("expected 0 allocs for FindImplicitNullabilityResults(nil, t), got %v", allocsNilFind)
	}

	allocsPlainVisit := testing.AllocsPerRun(100, func() {
		chk.Visit(trav, varNode, script)
	})
	if allocsPlainVisit != 0 {
		t.Fatalf("expected 0 allocs for ImplicitNullabilityCheck.Visit without JSDoc, got %v", allocsPlainVisit)
	}
}

func TestRound96_PrivateMemberNestedObjectLitSuper(t *testing.T) {
	parseFirstExpr := func(code string) rhino.Node {
		comp := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
		comp.Init(
			[]*SourceFile{SourceFile_FromCode("ext.js", "")},
			[]*SourceFile{SourceFile_FromCode("in.js", code)},
			opts,
		)
		comp.ParseInputs()
		if comp.HasErrors() {
			t.Fatalf("parse error on %q: %v", code, formatJSErrors(comp.GetErrors()))
		}
		return comp.GetJsRoot().GetFirstChild().GetFirstChild()
	}

	noOuterSuperCases := []string{
		`({ m() { return super.x; } });`,
		`({ get a() { return super.x; } });`,
		`({ set a(v) { super.x = v; } });`,
		`({ ['k']() { return super.x; } });`,
		`({ get ['k']() { return super.x; } });`,
		`({ set ['k'](v) { super.x = v; } });`,
	}
	for _, code := range noOuterSuperCases {
		stmt := parseFirstExpr(code)
		if NodeUtil_ReferencesSuper(stmt) {
			t.Errorf("NodeUtil_ReferencesSuper(%s) = true, want false (object-literal method binds own super)", code)
		}
	}

	hasOuterSuperCases := []string{
		`({ [super.k]() { return super.own; } });`,
		`({ fn: () => super.x });`,
		`({ val: super.x });`,
	}
	for _, code := range hasOuterSuperCases {
		stmt := parseFirstExpr(code)
		if !NodeUtil_ReferencesSuper(stmt) {
			t.Errorf("NodeUtil_ReferencesSuper(%s) = false, want true", code)
		}
	}

	// Verify end-to-end ES2022 private method rewriting preserves `super` inside nested
	// object-literal methods/getters/setters while rewriting `super` in computed keys and arrows!
	src := `
class Base {
  static get baseKey() { return 'k'; }
  get baseProp() { return 42; }
}
class Sub extends Base {
  #priv() {
    return {
      m() { return super.ownMethodProp; },
      get g() { return super.ownGetterProp; },
      set s(v) { super.ownSetterProp = v; },
      [super.baseProp]() { return super.ownComputedBodyProp; },
      arrow: () => super.baseProp,
    };
  }
  run() {
    return this.#priv();
  }
}
`
	comp := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetAssumeStaticInheritanceIsNotUsed(true)
	externs := NewTestExternsBuilder().AddObject().AddReflect().AddJSCompLibraries().Build()
	comp.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		[]*SourceFile{SourceFile_FromCode("in.js", src)},
		opts,
	)
	if comp.HasErrors() {
		t.Fatalf("unexpected compile errors: %v", formatJSErrors(comp.GetErrors()))
	}
	out := comp.ToSource()
	for _, preserved := range []string{
		"super.ownMethodProp",
		"super.ownGetterProp",
		"super.ownSetterProp",
		"super.ownComputedBodyProp",
	} {
		if !strings.Contains(out, preserved) {
			t.Errorf("expected %q to remain untouched inside object-literal method/accessor, got:\n%s", preserved, out)
		}
	}
	if strings.Contains(out, "super.baseProp") {
		t.Errorf("expected outer private-method super.baseProp to be rewritten, got:\n%s", out)
	}
}

func round96CompileToES5AndValidate(t *testing.T, name string, code string) string {
	t.Helper()
	comp := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	opts.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	externs := CompilerTestCase_DEFAULT_EXTERNS +
		NewTestExternsBuilder().
			AddMath().
			AddObject().
			AddReflect().
			AddJSCompLibraries().
			AddExtra("/** @constructor */ function WeakMap() {}\nWeakMap.prototype.get = function(k) {};\nWeakMap.prototype.set = function(k, v) {};\nWeakMap.prototype.has = function(k) {};\nvar a, b, i, k, x, p, ns = {}, items = [], obj = {}, args = [], res; function getBase() { return function() {}; } function use(...v) {}\n").
			Build()
	comp.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		[]*SourceFile{SourceFile_FromCode(name+".js", code)},
		opts,
	)
	if comp.HasErrors() {
		t.Fatalf("%s: unexpected compile errors: %v", name, formatJSErrors(comp.GetErrors()))
	}
	NewAstValidator(comp).ValidateRoot(comp.GetRoot())
	return comp.ToSource()
}

func TestRound96_LabelClimbingAndConciseArrowAcrossTranspilers(t *testing.T) {
	cases := []struct {
		name         string
		src          string
		wantContains string
	}{
		{
			name: "destructuring_labeled_for_of_outer_vars",
			src: `lbl1: lbl2: for ([a, b] of items) {
  if (a) break lbl1;
}`,
			wantContains: "lbl1:",
		},
		{
			name: "destructuring_labeled_for_in_outer_vars",
			src: `lbl1: for ({x: ns.prop} in obj) {
  continue lbl1;
}`,
			wantContains: "lbl1:",
		},
		{
			name:         "exp_assign_labeled_getprop",
			src:          `lbl1: lbl2: ns.val **= 2;`,
			wantContains: "Math.pow",
		},
		{
			name:         "exp_assign_labeled_getelem",
			src:          `lbl1: ns[k] **= 3;`,
			wantContains: "Math.pow",
		},
		{
			name:         "exp_assign_concise_arrow_getprop",
			src:          `const f = (o) => o.val **= 2; use(f);`,
			wantContains: "return",
		},
		{
			name:         "exp_assign_nested_concise_arrow_getelem",
			src:          `const f = (o) => (idx) => o[idx] **= 2; use(f);`,
			wantContains: "Math.pow",
		},
		{
			name: "exp_assign_private_field_es2022_to_es5",
			src: `class Counter {
  #count = 3;
  pow2() {
    return this.#count **= 2;
  }
}
use(new Counter().pow2());`,
			wantContains: "Math.pow",
		},
		{
			name:         "spread_call_labeled_side_effect_receiver",
			src:          `lbl1: lbl2: obj.getMaker().invoke(1, ...args);`,
			wantContains: ".apply(",
		},
		{
			name:         "spread_call_concise_arrow_side_effect_receiver",
			src:          `const callIt = (o, xs) => o.getMaker().invoke(1, ...xs); use(callIt);`,
			wantContains: ".apply(",
		},
		{
			name:         "computed_prop_labeled_assignment",
			src:          `lbl1: lbl2: res = {[k]: 1, fixed: 2};`,
			wantContains: "lbl1:",
		},
		{
			name:         "computed_prop_labeled_for_in_header",
			src:          `lbl1: for (var key in {[k]: 1}) { use(key); }`,
			wantContains: "lbl1:",
		},
		{
			name:         "computed_prop_concise_arrow",
			src:          `const makeComputed = (key) => ({[key]: 42}); use(makeComputed);`,
			wantContains: "return",
		},
		{
			name: "class_normalize_and_rewrite_labeled_qname_assign",
			src: `lbl1: lbl2: ns.MyCls = class {
  static s = 10;
  method() { return 20; }
};`,
			wantContains: "ns.MyCls",
		},
		{
			name: "class_normalize_and_rewrite_labeled_var_complex_extends",
			src: `lbl1: var Derived = class extends getBase() {
  method() { return 30; }
}; use(Derived);`,
			wantContains: "getBase()",
		},
		{
			name: "class_normalize_and_rewrite_concise_arrow_with_fields_and_private",
			src: `const makeFactoryClass = (Base) => class extends Base {
  static tag = 'ok';
  #secret = 7;
  readSecret() { return this.#secret; }
};
use(makeFactoryClass);`,
			wantContains: "return",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := round96CompileToES5AndValidate(t, tc.name, tc.src)
			if tc.wantContains != "" && !strings.Contains(out, tc.wantContains) {
				t.Fatalf("expected output to contain %q, got:\n%s", tc.wantContains, out)
			}
		})
	}
}

func TestRound96_AngularPass_AnonymousAndExportDefault(t *testing.T) {
	runAngular := func(code string) ([]string, string) {
		comp := NewCompiler()
		opts := NewCompilerOptions()
		opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		opts.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
		opts.SetAngularPass(true)
		comp.Init(
			[]*SourceFile{SourceFile_FromCode("externs.js", CompilerTestCase_DEFAULT_EXTERNS)},
			[]*SourceFile{SourceFile_FromCode("mod.js", code)},
			opts,
		)
		comp.ParseInputs()
		if comp.HasErrors() {
			t.Fatalf("unexpected parse errors on %q: %v", code, formatJSErrors(comp.GetErrors()))
		}
		NewAngularPass(comp).Process(comp.GetExternsRoot(), comp.GetJsRoot())
		errs := make([]string, len(comp.GetErrors()))
		for i, e := range comp.GetErrors() {
			errs[i] = e.Type().Key()
		}
		return errs, comp.ToSource()
	}

	nonGlobalOrBlockCases := []string{
		`export default /** @ngInject */ function(a, b) {}`,
		`export default class { /** @ngInject */ constructor(a, b) {} }`,
		`export default class { /** @ngInject */ foo(a, b) {} }`,
		`(class { /** @ngInject */ constructor(a, b) {} });`,
		`(class { /** @ngInject */ foo(a, b) {} });`,
		`(class NamedExpr { /** @ngInject */ foo(a, b) {} });`,
	}
	for _, code := range nonGlobalOrBlockCases {
		errs, out := runAngular(code)
		if !slices.Contains(errs, AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR.Key()) {
			t.Errorf("code=%q: expected error %s, got %v (out=%s)", code, AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR.Key(), errs, out)
		}
		if strings.Contains(out, "null.prototype") || strings.Contains(out, "null[") {
			t.Errorf("code=%q: emitted bogus null target in output: %s", code, out)
		}
	}

	nonQNameCases := []string{
		`var foo = {}; foo['bar'] = class { /** @ngInject */ constructor(a) {} };`,
		`var foo = {}; foo['bar'] = class { /** @ngInject */ method(a) {} };`,
	}
	for _, code := range nonQNameCases {
		errs, out := runAngular(code)
		if !slices.Contains(errs, AngularPass_INJECTED_FUNCTION_ON_NON_QNAME.Key()) {
			t.Errorf("code=%q: expected error %s, got %v (out=%s)", code, AngularPass_INJECTED_FUNCTION_ON_NON_QNAME.Key(), errs, out)
		}
		if strings.Contains(out, "null.prototype") {
			t.Errorf("code=%q: emitted bogus null.prototype in output: %s", code, out)
		}
	}
}
