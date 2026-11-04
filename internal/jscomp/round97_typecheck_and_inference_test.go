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
	"math"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type round97CheckResult struct {
	errors         []string
	warnings       []string
	mismatches     []string
	typedPercent   float64
	parallelJS     bool
	parallelWorker int
	hasRecEnumFlag bool
}

func runRound97TypeCheck(t *testing.T, externsCode string, scripts []*SourceFile, numThreads int) round97CheckResult {
	t.Helper()
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetCheckTypes(true)
	opts.SetClosurePass(true)
	// Intentionally leave FastMode false (default mode!) to exercise default-mode parallel TypeCheck.
	opts.SetNumParallelThreads(numThreads)
	opts.SetWarningLevel(DiagnosticGroups_CHECK_TYPES, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_WARNING)

	externs := []*SourceFile{
		SourceFile_FromCode("externs.js", externsCode),
	}
	compiler.Init(externs, scripts, opts)
	root := compiler.ParseInputs()
	if compiler.GetErrorCount() != 0 {
		t.Fatalf("parse errors (threads=%d): %v", numThreads, compiler.GetErrors())
	}
	externsRoot := root.GetFirstChild()
	jsRoot := externsRoot.GetNext()
	reg := compiler.GetTypeRegistry()
	rai := NewSemanticReverseAbstractInterpreter(reg)
	tc := NewTypeCheck(compiler, rai, reg)
	tc.ProcessForTesting(externsRoot, jsRoot)
	_, workers, canPar := tc.canCheckInParallel(jsRoot)

	formatDiags := func(diags []*JSError) []string {
		out := make([]string, len(diags))
		for i, d := range diags {
			out[i] = fmt.Sprintf("%s:%d:%d:%s:%s",
				d.sourceName, d.GetLineNumber(), d.GetCharno(), d.GetType().Key(), d.GetDescription())
		}
		return out
	}

	rawMismatches := compiler.GetTypeValidator().GetMismatches()
	mismatches := make([]string, len(rawMismatches))
	for i, m := range rawMismatches {
		mismatches[i] = fmt.Sprintf("%s<-%s", m.Found, m.Required)
	}

	return round97CheckResult{
		errors:         formatDiags(compiler.GetErrors()),
		warnings:       formatDiags(compiler.GetWarnings()),
		mismatches:     mismatches,
		typedPercent:   tc.GetTypedPercent(),
		parallelJS:     canPar,
		parallelWorker: workers,
		hasRecEnumFlag: tc.scopeCreator.hasUnevaluatedRecordEnumAlias,
	}
}

func assertRound97CheckParity(t *testing.T, seq, par round97CheckResult) {
	t.Helper()
	if len(seq.errors) != len(par.errors) {
		t.Fatalf("errors count mismatch:\nseq=%v\npar=%v", seq.errors, par.errors)
	}
	for i := range seq.errors {
		if seq.errors[i] != par.errors[i] {
			t.Fatalf("error[%d] mismatch:\nseq=%s\npar=%s", i, seq.errors[i], par.errors[i])
		}
	}
	if len(seq.warnings) != len(par.warnings) {
		t.Fatalf("warnings count mismatch:\nseq (%d)=%v\npar (%d)=%v",
			len(seq.warnings), seq.warnings, len(par.warnings), par.warnings)
	}
	for i := range seq.warnings {
		if seq.warnings[i] != par.warnings[i] {
			t.Fatalf("warning[%d] mismatch:\nseq=%s\npar=%s", i, seq.warnings[i], par.warnings[i])
		}
	}
	if len(seq.mismatches) != len(par.mismatches) {
		t.Fatalf("mismatches count mismatch: seq=%d par=%d", len(seq.mismatches), len(par.mismatches))
	}
	for i := range seq.mismatches {
		if seq.mismatches[i] != par.mismatches[i] {
			t.Fatalf("mismatch[%d] differs: seq=%s par=%s", i, seq.mismatches[i], par.mismatches[i])
		}
	}
	if math.Abs(seq.typedPercent-par.typedPercent) > 1e-9 {
		t.Fatalf("typedPercent mismatch: seq=%v par=%v", seq.typedPercent, par.typedPercent)
	}
}

func TestRound97DefaultParallelTypeCheckTemplatedRecordRace(t *testing.T) {
	externs := `
var Array;
/** @type {number} */ Array.prototype.length;
`
	buildScripts := func() []*SourceFile {
		scripts := make([]*SourceFile, 0, 12)
		scripts = append(scripts, SourceFile_FromCode("s00.js", `
/**
 * @template T
 * @param {{a: T, b: !Array<T>, nested: {val: T}}} rec
 * @return {T}
 */
function extractRec(rec) {
  return rec.nested.val;
}

/**
 * @template K, V
 * @param {{key: K, values: !Array<V>}} entry
 * @return {V}
 */
function firstEntryValue(entry) {
  return entry.values[0];
}
`))
		for i := 1; i < 12; i++ {
			code := fmt.Sprintf(`
(function() {
  var numVal_%d = extractRec({a: %d, b: [%d], nested: {val: %d}});
  /** @type {number} */ var okNum_%d = numVal_%d;
  var strVal_%d = extractRec({a: 's%d', b: ['s%d'], nested: {val: 's%d'}});
  /** @type {number} */ var badAssign_%d = strVal_%d;
  var pair_%d = firstEntryValue({key: 'k%d', values: [{p%d: %d}]});
  var missing_%d = pair_%d.nonExistentProp_%d;
})();
`, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i)
			scripts = append(scripts, SourceFile_FromCode(fmt.Sprintf("s%02d.js", i), code))
		}
		return scripts
	}

	seq := runRound97TypeCheck(t, externs, buildScripts(), 1)
	if len(seq.warnings) < 22 {
		t.Fatalf("expected >=22 warnings from type mismatches and missing properties, got %d: %v",
			len(seq.warnings), seq.warnings)
	}

	for iter := range 4 {
		par := runRound97TypeCheck(t, externs, buildScripts(), 4)
		if !par.parallelJS || par.parallelWorker < 2 {
			t.Fatalf("iter %d: expected default-mode parallel TypeCheck with >=2 workers, got ok=%v workers=%d",
				iter, par.parallelJS, par.parallelWorker)
		}
		assertRound97CheckParity(t, seq, par)
	}
}

func TestRound97ParallelTypeCheckEnumAliasReporterAndRecordFallback(t *testing.T) {
	externs := `var Object;`

	t.Run("NonRecordEnumAliasesStayParallelAndDeterministic", func(t *testing.T) {
		buildScripts := func() []*SourceFile {
			scripts := make([]*SourceFile, 0, 12)
			scripts = append(scripts, SourceFile_FromCode("e00.js", `
/** @enum {number} */
const BaseNumEnum = {A: 1, B: 2};
/** @enum {string} */
const BaseStrEnum = {X: 'x', Y: 'y'};
`))
			for i := 1; i < 12; i++ {
				var code string
				if i%2 == 0 {
					code = fmt.Sprintf(`
/** @enum {MissingEnumType_%d} */
const BadUnknownAlias_%d = BaseNumEnum;
`, i, i)
				} else {
					code = fmt.Sprintf(`
/** @enum {string} */
const MismatchedAlias_%d = BaseNumEnum;
/** @enum {number} */
const ValidAlias_%d = BaseNumEnum;
/** @type {BaseNumEnum} */
const checkElem_%d = ValidAlias_%d.A;
`, i, i, i, i)
				}
				scripts = append(scripts, SourceFile_FromCode(fmt.Sprintf("e%02d.js", i), code))
			}
			return scripts
		}

		seq := runRound97TypeCheck(t, externs, buildScripts(), 1)
		if len(seq.warnings) != 11 {
			t.Fatalf("expected 11 deterministic enum alias warnings (5 unrecognized + 6 incompatible), got %d: %v",
				len(seq.warnings), seq.warnings)
		}

		for iter := range 4 {
			par := runRound97TypeCheck(t, externs, buildScripts(), 4)
			if par.hasRecEnumFlag {
				t.Fatalf("iter %d: unexpected hasUnevaluatedRecordEnumAlias=true for non-record enum aliases", iter)
			}
			if !par.parallelJS || par.parallelWorker < 2 {
				t.Fatalf("iter %d: expected parallel TypeCheck enabled, got ok=%v workers=%d",
					iter, par.parallelJS, par.parallelWorker)
			}
			assertRound97CheckParity(t, seq, par)
		}
	})

	t.Run("RecordEnumAliasDisablesParallelCleanly", func(t *testing.T) {
		buildScripts := func() []*SourceFile {
			return []*SourceFile{
				SourceFile_FromCode("r0.js", `
/** @enum {{rec: number}} */
const BaseRecEnum = {FIRST: {rec: 1}};
`),
				SourceFile_FromCode("r1.js", `
/** @enum {{rec: number, extra: MissingRecPropType}} */
const AliasRecEnum = BaseRecEnum;
`),
				SourceFile_FromCode("r2.js", `
/** @type {number} */
const goodRecField = AliasRecEnum.FIRST.rec;
`),
				SourceFile_FromCode("r3.js", `
/** @type {string} */
const badRecField = AliasRecEnum.FIRST.rec;
`),
			}
		}

		seq := runRound97TypeCheck(t, externs, buildScripts(), 1)
		par := runRound97TypeCheck(t, externs, buildScripts(), 4)
		if !par.hasRecEnumFlag {
			t.Fatalf("expected hasUnevaluatedRecordEnumAlias=true when @enum {{...}} aliases an EnumType")
		}
		if par.parallelJS {
			t.Fatalf("expected canCheckInParallel=false when hasUnevaluatedRecordEnumAlias is set")
		}
		assertRound97CheckParity(t, seq, par)
	})
}

func TestRound97GetPropQualifiedSlotFastRejectAndNarrowing(t *testing.T) {
	externs := `var Object;`
	code := `
class Base {
  constructor() {
    /** @type {number|string} */
    this.prop = 1;
  }
  /** @return {number} */
  baseMethod() { return 1; }
}

class Sub extends Base {
  /** @return {number} */
  testThisAndSuper() {
    if (typeof this.prop === 'number') {
      /** @type {number} */ var narrowedThis = this.prop;
      return narrowedThis + super.baseMethod();
    }
    return super.baseMethod();
  }
}

/**
 * @param {{b: {c: (number|null)}}} a
 * @param {?{x: number}} optObj
 * @param {function(): {k: number}} factory
 * @return {number}
 */
function testChains(a, optObj, factory) {
  var optVal = optObj?.x;
  var callVal = factory().k;
  if (a.b.c !== null) {
    /** @type {number} */ var narrowedChain = a.b.c;
    return narrowedChain + callVal + (optVal || 0);
  }
  return callVal;
}

/**
 * @param {{b: {c: (number|null)}}} a
 * @return {string}
 */
function testBadNarrowedChain(a) {
  if (a.b.c !== null) {
    return a.b.c;
  }
  return 'ok';
}
`
	res := runRound97TypeCheck(t, externs, []*SourceFile{SourceFile_FromCode("qname.js", code)}, 1)
	if len(res.warnings) != 1 || !strings.Contains(res.warnings[0], "inconsistent return type") {
		t.Fatalf("expected exactly 1 inconsistent return type warning from testBadNarrowedChain, got: %v", res.warnings)
	}

	// Verify LookupGetPropSlotType fast-rejects in 0 heap allocations on a local
	// GETPROP chain even when globalScope has qualified vars.
	compiler := NewCompiler()
	reg := compiler.GetTypeRegistry()
	numType := reg.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)

	fn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(rhino.IR_Name("a")), rhino.IR_Block())
	getPropNode := rhino.IR_Getprop(rhino.IR_Getprop(rhino.IR_Name("a"), "b"), "c")
	script := rhino.IR_Script()
	script.AddChildToBack(fn)
	root := rhino.IR_Root(script)

	globalScope := TypedScope_CreateGlobalScope(root)
	globalScope.Declare("ns.globalQual", rhino.IR_Getprop(rhino.IR_Name("ns"), "globalQual"), numType, nil, false)
	if !globalScope.hasQualifiedVars() {
		t.Fatalf("expected globalScope.hasQualifiedVars()=true")
	}

	fnScope := NewTypedScope(globalScope, fn)
	fnScope.Declare("a", rhino.IR_Name("a"), numType, nil, false)
	if fnScope.hasQualifiedVars() {
		t.Fatalf("expected fnScope.hasQualifiedVars()=false")
	}

	var arena linkedFlowScopeArena
	entry := newLinkedFlowScopeWithArena(&arena, compiler, LinkedFlowScope_EMPTY_SCOPES, fnScope, fnScope)
	withSimpleNarrow := entry.InferSlotType("a", numType)

	allocs := testing.AllocsPerRun(100, func() {
		qName, isQName, _, _, _, _, ok := withSimpleNarrow.LookupGetPropSlotType(getPropNode)
		if !isQName || ok || qName != "" {
			t.Fatalf("expected fast-rejected miss (qName=\"\", isQName=true, ok=false), got (%q, %v, %v)", qName, isQName, ok)
		}
	})
	if allocs != 0 {
		t.Fatalf("LookupGetPropSlotType fast-reject allocs = %v, want 0", allocs)
	}

	// Now infer a qualified slot "a.b.c" and verify LookupGetPropSlotType resolves it accurately.
	withQualNarrow := withSimpleNarrow.InferSlotType("a.b.c", numType)
	qName, isQName, _, _, gotType, isInferred, ok := withQualNarrow.LookupGetPropSlotType(getPropNode)
	if !isQName || !ok || !isInferred || qName != "a.b.c" || gotType != numType {
		t.Fatalf("LookupGetPropSlotType hit = (%q, %v, %v, %v, %v), want (\"a.b.c\", true, number, true, true)",
			qName, isQName, gotType, isInferred, ok)
	}
}

func TestRound97TraverseNameOuterVsLocalAndShadowing(t *testing.T) {
	externs := `var Object; /** @type {number} */ String.prototype.length;`
	code := `
var singleOuter = 42;
function readSingleOuter() {
  /** @type {number} */ var ok = singleOuter;
  return ok;
}

var multiOuter = null;
function readMultiOuterBad() {
  // multiOuter is assigned more than once across scopes; cannot assume it is null here.
  /** @type {string} */ var bad = multiOuter;
  return bad;
}
multiOuter = {x: 1};

function readShadowedInBlock() {
  {
    let multiOuter = 'localShadow';
    /** @type {string} */ var goodShadow = multiOuter;
    if (goodShadow.length > 0) {
      return goodShadow;
    }
  }
  return '';
}
`
	res := runRound97TypeCheck(t, externs, []*SourceFile{SourceFile_FromCode("outer.js", code)}, 1)
	if len(res.warnings) != 1 || !strings.Contains(res.warnings[0], "initializing variable") {
		t.Fatalf("expected 1 warning on multiOuter read in readMultiOuterBad, got: %v", res.warnings)
	}
}

func TestRound97FlowThroughInnerScopesVsFlatFastPath(t *testing.T) {
	externs := `/** @type {undefined} */ var undefined;`
	code := `
/**
 * @param {number} a
 * @param {boolean} cond
 * @return {number}
 */
function flatFn(a, cond) {
  var hoistedUnbound;
  function unusedInnerFn() { return 1; }
  hoistedUnbound = cond ? a : hoistedUnbound;
  return hoistedUnbound !== undefined ? hoistedUnbound : a + 5;
}

/**
 * @param {number} a
 * @param {boolean} cond
 * @return {number}
 */
function flatFnUnnarrowed(a, cond) {
  var hoistedUnbound;
  hoistedUnbound = cond ? a : hoistedUnbound;
  // Intentional mismatch: hoistedUnbound is (number|undefined) when !cond.
  return hoistedUnbound;
}

/**
 * @param {number=} a
 * @return {number}
 */
function defaultParamShadow(a = 10) {
  var a;
  return a || 0;
}

/** @type {function(number): number} */
const conciseArrow = (x) => x + 1;

/** @return {function(): number} */
const nestedConciseArrow = () => () => 42;

/**
 * @param {number} n
 * @return {string}
 */
function innerScopedFn(n) {
  let acc = 0;
  for (let i = 0; i < n; i++) {
    acc += i;
  }
  try {
    {
      let blockScoped = 'str';
      if (acc > 10) {
        return blockScoped;
      }
    }
  } catch (e) {
    return 'err';
  }
  // Intentional type mismatch: returning number acc where string is required.
  return acc;
}
`
	res := runRound97TypeCheck(t, externs, []*SourceFile{SourceFile_FromCode("flow.js", code)}, 1)
	if len(res.warnings) != 2 ||
		!strings.Contains(res.warnings[0], "(number|undefined)") ||
		!strings.Contains(res.warnings[1], "required: string") {
		t.Fatalf("expected 2 inconsistent return type warnings (flatFnUnnarrowed + innerScopedFn), got: %v", res.warnings)
	}

	// Verify ControlFlowAnalysis sets hasInnerScopeRoots=false on flatFn & conciseArrow,
	// and hasInnerScopeRoots=true on nestedConciseArrow & innerScopedFn.
	compiler := NewCompiler()
	opts := NewCompilerOptions()
	compiler.Init([]*SourceFile{SourceFile_FromCode("externs.js", externs)}, []*SourceFile{SourceFile_FromCode("flow.js", code)}, opts)
	root := compiler.ParseInputs()
	script := root.GetSecondChild().GetFirstChild()

	checkHasInner := func(fnNode rhino.Node) bool {
		cfg := ComputeControlFlowGraph(compiler, fnNode, false, true)
		return cfg.hasInnerScopeRoots
	}

	var flatNode, conciseNode, nestedConciseNode, innerScopedNode rhino.Node
	for child := script.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if child.IsFunction() {
			switch child.GetFirstChild().GetString() {
			case "flatFn":
				flatNode = child
			case "innerScopedFn":
				innerScopedNode = child
			}
		} else if child.IsConst() {
			decl := child.GetFirstChild()
			switch decl.GetString() {
			case "conciseArrow":
				conciseNode = decl.GetFirstChild()
			case "nestedConciseArrow":
				nestedConciseNode = decl.GetFirstChild()
			}
		}
	}

	if checkHasInner(flatNode) {
		t.Errorf("expected flatFn hasInnerScopeRoots=false, got true")
	}
	if checkHasInner(conciseNode) {
		t.Errorf("expected conciseArrow hasInnerScopeRoots=false, got true")
	}
	if !checkHasInner(nestedConciseNode) {
		t.Errorf("expected nestedConciseArrow hasInnerScopeRoots=true, got false")
	}
	if !checkHasInner(innerScopedNode) {
		t.Errorf("expected innerScopedFn hasInnerScopeRoots=true, got false")
	}
}

func TestRound97ConcreteFlowHamtParityAndZeroAllocs(t *testing.T) {
	compiler := NewCompiler()
	reg := compiler.GetTypeRegistry()
	numType := reg.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	strType := reg.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	boolType := reg.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	numOrStr := reg.CreateUnionType(numType, strType)

	fn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), rhino.IR_Block())
	script := rhino.IR_Script()
	script.AddChildToBack(fn)
	root := rhino.IR_Root(script)
	globalScope := TypedScope_CreateGlobalScope(root)
	fnScope := NewTypedScope(globalScope, fn)

	const numVars = 64
	ids := make([]uint32, numVars)
	names := make([]string, numVars)
	for i := range numVars {
		name := fmt.Sprintf("slot_%02d", i)
		names[i] = name
		ids[i] = uint32(strpool.Intern(name))
		fnScope.Declare(name, rhino.IR_Name(name), numOrStr, nil, true)
	}

	var arena linkedFlowScopeArena
	joinOp := NewLinkedFlowScope_FlowScopeJoinOp(compiler)

	entry := newLinkedFlowScopeWithArena(&arena, compiler, LinkedFlowScope_EMPTY_SCOPES, fnScope, fnScope)
	left := entry
	right := entry
	for i := range numVars {
		left = left.InferSlotTypeByID(ids[i], names[i], numType)
		if i%2 == 0 {
			right = right.InferSlotTypeByID(ids[i], names[i], strType)
		} else {
			right = right.InferSlotTypeByID(ids[i], names[i], numType)
		}
	}
	right = right.InferSlotType("extra_right_only", boolType)

	joinOp.result = nil
	joinOp.JoinFlow(left)
	joinOp.JoinFlow(right)
	joined := joinOp.Finish()

	for i := range numVars {
		_, _, ft, inferred, ok := joined.LookupSlotTypeByID(ids[i])
		if !ok || !inferred {
			t.Fatalf("slot %d missing after join: ok=%v inferred=%v", i, ok, inferred)
		}
		want := numType
		if i%2 == 0 {
			want = numOrStr
		}
		if !ft.Equals(want) {
			t.Fatalf("slot %d joined type = %v, want %v", i, ft, want)
		}
	}

	// Verify symmetry of Equivalent & JoinFlow.
	joinOpRev := NewLinkedFlowScope_FlowScopeJoinOp(compiler)
	joinOpRev.JoinFlow(right)
	joinOpRev.JoinFlow(left)
	joinedRev := joinOpRev.Finish()
	if !joined.EqualScope(joinedRev) {
		t.Fatalf("expected joined and reverse-joined LinkedFlowScopes to be EqualScope")
	}

	// Verify 0 allocations on LookupSlotTypeByID + InferSlotTypeByID with warmed arena.
	allocs := testing.AllocsPerRun(100, func() {
		arena.reset()
		s := newLinkedFlowScopeWithArena(&arena, compiler, LinkedFlowScope_EMPTY_SCOPES, fnScope, fnScope)
		s1 := s.InferSlotTypeByID(ids[0], names[0], numType).InferSlotTypeByID(ids[1], names[1], boolType)
		s2 := s.InferSlotTypeByID(ids[0], names[0], numOrStr).InferSlotTypeByID(ids[1], names[1], boolType)
		joinOp.result = nil
		joinOp.JoinFlow(s1)
		joinOp.JoinFlow(s2)
		j := joinOp.Finish()
		_, _, ft, _, ok := j.LookupSlotTypeByID(ids[0])
		if !ok || ft != numOrStr {
			t.Fatal("unexpected joined slot type")
		}
	})
	if allocs != 0 {
		t.Fatalf("warmed concrete HAMT infer+join+lookupByID allocs = %v, want 0", allocs)
	}
}
