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
	"io"
	"math"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func runRound98ColorsPipeline(
	t *testing.T,
	externsFiles, srcFiles []*SourceFile,
	numThreads int,
	serOpts *SerializationOptions,
) (*Compiler, []astColorSnapshot, *colors.ColorRegistry) {
	t.Helper()
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	errMgr := newTypeCheckDump_errorManager()
	compiler.SetErrorManager(errMgr)
	options := typeCheckDump_options()
	options.SetNumParallelThreads(numThreads)
	compiler.Init(externsFiles, srcFiles, options)
	compiler.Parse()
	serializationDump_typeCheck(compiler, options)
	NewRemoveCastNodes(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	NewConvertTypesToColors(compiler, serOpts).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	snaps := collectAstColorSnapshots(compiler.GetJsRoot().GetParent())
	return compiler, snaps, compiler.GetColorRegistry()
}

func TestRound98_ConvertTypesToColors_ParallelParityAndReconCache(t *testing.T) {
	externs := []*SourceFile{
		SourceFile_FromCode("externs1.js", `
/** @constructor */ function Object() {}
/** @constructor */ function Function() {}
/** @constructor @param {*=} x @return {string} */ function String(x) {}
/** @constructor @param {*=} x @return {number} */ function Number(x) {}
/** @constructor @param {*=} x @return {boolean} */ function Boolean(x) {}
/** @type {number} */ var extGlobalNum;
`),
		SourceFile_FromCode("externs_summary.i.js", `
/** @fileoverview @typeSummary */
/** @interface */
function SummaryIface() {}
/** @return {string} */
SummaryIface.prototype.summaryMethod = function() {};
`),
		SourceFile_FromCode("externs2.js", `
/** @constructor */
function DomWidget() {}
/** @type {string} */
DomWidget.prototype.label;
/** @param {number} n @return {boolean} */
DomWidget.prototype.update = function(n) {};
`),
	}

	srcs := []*SourceFile{
		SourceFile_FromCode("s0.js", `
/** @license Preserved Header */
/** @interface */
class Renderable {
  /** @return {string} */
  render() {}
}
/** @implements {Renderable} */
class Box {
  constructor() {
    /** @type {number} */ this.width = 10;
    /** @type {string} */ this.label = 'box';
  }
  /** @override @noinline @return {string} */
  render() { return this.label; }
}
`),
		SourceFile_FromCode("s1.js", `
class FancyBox extends Box {
  constructor() {
    super();
    /** @type {boolean} */ this.active = true;
  }
  /** @nosideeffects @return {number} */
  area() { return this.width * 2; }
}
/** @desc Greeting @meaning ui_greeting */
var MSG_HELLO = 'hello';
`),
		SourceFile_FromCode("s2.js", `
/**
 * @param {!Box|!DomWidget} item
 * @param {*} raw
 * @return {string}
 */
function inspectItem(item, raw) {
  var casted = /** @type {number} */ (raw);
  return item.label + casted;
}
`),
		SourceFile_FromCode("s3.js", `
/** @type {number} */
var badAssignMismatch = 'not_a_number';
/** @param {!FancyBox} fb @return {string} */
function useFancy(fb) {
  return fb.render() + fb.area() + extGlobalNum;
}
`),
		SourceFile_FromCode("s4.js", `
/** @enum {number} */
const Mode = { FAST: 1, SAFE: 2 };
/** @param {Mode} m @return {number} */
function modeVal(m) {
  var w = new DomWidget();
  w.update(m);
  return m;
}
`),
		SourceFile_FromCode("s5.js", `
function propBurst() {
  var a = {p0: 1, p1: 'a', p2: true, p3: 4, p4: 'b', p5: 6};
  var b = {p0: 2, p1: 'c', p2: false, p3: 8, p4: 'd', p5: 9};
  return a.p0 + b.p0 + a.p3 + b.p5 + (a.p2 ? a.p1 : b.p4).length;
}
`),
	}

	makeFiles := func() ([]*SourceFile, []*SourceFile) {
		ex := make([]*SourceFile, len(externs))
		for i, f := range externs {
			ex[i] = SourceFile_FromCode(f.GetName(), mustGetSourceCode(t, f))
		}
		in := make([]*SourceFile, len(srcs))
		for i, f := range srcs {
			in[i] = SourceFile_FromCode(f.GetName(), mustGetSourceCode(t, f))
		}
		return ex, in
	}

	debugSerOpts := NewSerializationOptions().SetIncludeDebugInfo(true).SetRunValidation(true).Build()
	ex1, in1 := makeFiles()
	ex4, in4 := makeFiles()
	_, snapsSeq, regSeq := runRound98ColorsPipeline(t, ex1, in1, 1, debugSerOpts)
	_, snapsPar, regPar := runRound98ColorsPipeline(t, ex4, in4, 4, debugSerOpts)

	assertIdenticalColorGraphs(t, snapsSeq, snapsPar, regSeq, regPar, true)
}

func mustGetSourceCode(t *testing.T, sf *SourceFile) string {
	t.Helper()
	return sf.GetCode()
}

func TestRound98_ConvertTypesToColors_UntypedRemoveTypesParallelParity(t *testing.T) {
	run := func(numThreads int) (string, []astColorSnapshot) {
		compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
		opts := NewCompilerOptions()
		opts.SetNumParallelThreads(numThreads)
		externs := []*SourceFile{
			SourceFile_FromCode("ext1.js", "/** @type {number} */ var extA;"),
			SourceFile_FromCode("ext2.js", "/** @type {string} */ var extB;"),
		}
		srcs := []*SourceFile{
			SourceFile_FromCode("a.js", "/** @noinline @param {number} x */ function f(x) { return x + 1; }"),
			SourceFile_FromCode("b.js", "/** @license MIT */ /** @nosideeffects */ function g() { return f(2); }"),
			SourceFile_FromCode("c.js", "/** @desc d @meaning m */ var MSG_X = 'hi';"),
			SourceFile_FromCode("d.js", "var z = g() + MSG_X;"),
		}
		compiler.Init(externs, srcs, opts)
		root := compiler.ParseInputs()
		if compiler.GetErrorCount() != 0 {
			t.Fatalf("unexpected parse errors: %v", compiler.GetErrors())
		}
		NewConvertTypesToColors(compiler, NewSerializationOptions().SetIncludeDebugInfo(false).Build()).Process(
			root.GetFirstChild(), root.GetSecondChild())
		return compiler.ToSource(), collectAstColorSnapshots(root)
	}

	srcSeq, snapsSeq := run(1)
	srcPar, snapsPar := run(4)
	if srcSeq != srcPar {
		t.Fatalf("ToSource mismatch:\nseq:\n%s\npar:\n%s", srcSeq, srcPar)
	}
	if len(snapsSeq) != len(snapsPar) {
		t.Fatalf("snapshot count mismatch: seq=%d par=%d", len(snapsSeq), len(snapsPar))
	}
	for i := range snapsSeq {
		if snapsSeq[i] != snapsPar[i] {
			t.Fatalf("snapshot[%d] mismatch: seq=%+v par=%+v", i, snapsSeq[i], snapsPar[i])
		}
	}
}

type round98NormalizeResult struct {
	source    string
	constTags []string
	errors    []string
	warnings  []string
}

func runRound98Normalize(t *testing.T, externsCode []string, scriptsCode []string, numThreads int, googleConv bool) round98NormalizeResult {
	t.Helper()
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	opts := NewCompilerOptions()
	opts.SetNumParallelThreads(numThreads)
	if googleConv {
		opts.SetCodingConvention(NewGoogleCodingConvention())
	}
	externs := make([]*SourceFile, len(externsCode))
	for i, c := range externsCode {
		externs[i] = SourceFile_FromCode(fmt.Sprintf("ext%d.js", i), c)
	}
	srcs := make([]*SourceFile, len(scriptsCode))
	for i, c := range scriptsCode {
		srcs[i] = SourceFile_FromCode(fmt.Sprintf("in%d.js", i), c)
	}
	compiler.Init(externs, srcs, opts)
	root := compiler.ParseInputs()
	if compiler.GetErrorCount() != 0 {
		t.Fatalf("parse errors: %v", compiler.GetErrors())
	}
	NewConvertTypesToColors(compiler, NewSerializationOptions().SetIncludeDebugInfo(false).Build()).Process(
		root.GetFirstChild(), root.GetSecondChild())
	NewNormalize(compiler, false).Process(root.GetFirstChild(), root.GetSecondChild())

	var constTags []string
	var walk func(n rhino.Node)
	walk = func(n rhino.Node) {
		if n.IsName() {
			constTags = append(constTags, fmt.Sprintf("%s:%v", n.GetString(), n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)))
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c)
		}
	}
	walk(root)

	fmtDiags := func(ds []*JSError) []string {
		out := make([]string, len(ds))
		for i, d := range ds {
			out[i] = fmt.Sprintf("%s:%d:%s:%s", d.sourceName, d.GetLineNumber(), d.GetType().Key(), d.GetDescription())
		}
		return out
	}

	return round98NormalizeResult{
		source:    compiler.ToSource(),
		constTags: constTags,
		errors:    fmtDiags(compiler.GetErrors()),
		warnings:  fmtDiags(compiler.GetWarnings()),
	}
}

func TestRound98_Normalize_MultiScriptParallelParity(t *testing.T) {
	externs := []string{
		"var extShared; var EXT_CONST_VAL;",
		"function extFn(a) {}",
	}
	scripts := []string{
		`
/** @const */ var GLOBAL_ANNOT_CONST = 10;
var CONV_CONST_ONE = 20;
var dupGlobal = 1;
var extShared = 99;
`,
		`
var dupGlobal = 2;
function forAndWhileTest(n) {
  var total = 0;
  outer: inner: for (var i = 0, j = 10; i < n; i++, j--) {
    total += i + j + GLOBAL_ANNOT_CONST + CONV_CONST_ONE;
  }
  while (total > 100) {
    total -= 1;
  }
  return total;
}
`,
		`
var dupGlobal;
function localDedupAndShadow(x) {
  var localDup = 1;
  var localDup = 2;
  var localDup;
  const LOCAL_C = 5;
  {
    let x = 99;
    localDup += x + LOCAL_C;
  }
  try {
    throw new Error('e');
  } catch (errVar) {
    var errVar = null;
  }
  var fnExprHolder = function namedExpr(k) {
    return k > 0 ? namedExpr(k - 1) : 0;
  };
  return localDup + fnExprHolder(2);
}
`,
		`
function logicalAssignAndCommas(obj, a, b,) {
  obj.x ||= [1, 2,];
  obj.y &&= {k: a,};
  obj.z ??= extFn(b,);
  var arrow = (v) => v + EXT_CONST_VAL + GLOBAL_ANNOT_CONST;
  return arrow(1);
}
`,
		`
function shorthandConstCheck() {
  var CONV_CONST_TWO = 42;
  return {CONV_CONST_ONE, CONV_CONST_TWO, normalKey: dupGlobal};
}
`,
	}

	for _, googleConv := range []bool{false, true} {
		seq := runRound98Normalize(t, externs, scripts, 1, googleConv)
		par := runRound98Normalize(t, externs, scripts, 4, googleConv)
		if seq.source != par.source {
			t.Fatalf("googleConv=%v source mismatch:\nseq:\n%s\npar:\n%s", googleConv, seq.source, par.source)
		}
		if !slices.Equal(seq.constTags, par.constTags) {
			t.Fatalf("googleConv=%v constTags mismatch:\nseq=%v\npar=%v", googleConv, seq.constTags, par.constTags)
		}
		if !slices.Equal(seq.errors, par.errors) {
			t.Fatalf("googleConv=%v errors mismatch:\nseq=%v\npar=%v", googleConv, seq.errors, par.errors)
		}
		if !slices.Equal(seq.warnings, par.warnings) {
			t.Fatalf("googleConv=%v warnings mismatch:\nseq=%v\npar=%v", googleConv, seq.warnings, par.warnings)
		}
	}
}

func TestRound98_Normalize_AssertOnChangeParity(t *testing.T) {
	runAssert := func(numThreads int) (panicked bool, msg string) {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				msg = fmt.Sprint(r)
			}
		}()
		compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
		opts := NewCompilerOptions()
		opts.SetNumParallelThreads(numThreads)
		opts.SetCodingConvention(NewGoogleCodingConvention())
		externs := []*SourceFile{SourceFile_FromCode("ext.js", "")}
		srcs := []*SourceFile{
			SourceFile_FromCode("s0.js", "var ok = 1;"),
			SourceFile_FromCode("s1.js", "var UNMARKED_CONST = 2; use(UNMARKED_CONST);"),
			SourceFile_FromCode("s2.js", "var ok2 = 3;"),
		}
		compiler.Init(externs, srcs, opts)
		root := compiler.ParseInputs()
		// Run only b2 via Normalize with assertOnChange=true on unnormalized AST containing constant name.
		norm := NewNormalize(compiler, true)
		if !norm.runPropagateAndDedupParallel(root.GetFirstChild(), root.GetSecondChild()) {
			sc2 := NewSyntacticScopeCreatorWithRedeclarationHandler(compiler, NewNormalize_DuplicateDeclarationHandler(norm))
			(NodeTraversal_Builder{
				compiler:     compiler,
				callback:     NewNormalize_PropagateConstantPropertyOverVars(compiler, true),
				scopeCreator: sc2,
			}).TraverseRoots(root.GetFirstChild(), root.GetSecondChild())
		}
		return false, ""
	}

	p1, m1 := runAssert(1)
	p4, m4 := runAssert(4)
	if !p1 || !p4 {
		t.Fatalf("expected panic in both seq (%v) and par (%v)", p1, p4)
	}
	if !strings.Contains(m1, "Should be const but not marked as const") ||
		!strings.Contains(m4, "Should be const but not marked as const") {
		t.Fatalf("unexpected panic messages:\nseq=%s\npar=%s", m1, m4)
	}
}

type round98InferenceSnapshot struct {
	errors       []string
	warnings     []string
	mismatches   []string
	typedPercent float64
	nodeTypes    []string
}

func runRound98TypeCheckFull(t *testing.T, externs, scripts []*SourceFile, numThreads int) round98InferenceSnapshot {
	t.Helper()
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	opts := NewCompilerOptions()
	opts.SetCheckTypes(true)
	opts.SetClosurePass(true)
	opts.SetNumParallelThreads(numThreads)
	opts.SetWarningLevel(DiagnosticGroups_CHECK_TYPES, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_WARNING)

	exCopy := make([]*SourceFile, len(externs))
	for i, f := range externs {
		exCopy[i] = SourceFile_FromCode(f.GetName(), mustGetSourceCode(t, f))
	}
	inCopy := make([]*SourceFile, len(scripts))
	for i, f := range scripts {
		inCopy[i] = SourceFile_FromCode(f.GetName(), mustGetSourceCode(t, f))
	}

	compiler.Init(exCopy, inCopy, opts)
	root := compiler.ParseInputs()
	if compiler.GetErrorCount() != 0 {
		t.Fatalf("parse errors (threads=%d): %v", numThreads, compiler.GetErrors())
	}

	serializationDump_typeCheck(compiler, opts)

	fmtDiags := func(ds []*JSError) []string {
		out := make([]string, len(ds))
		for i, d := range ds {
			out[i] = fmt.Sprintf("%s:%d:%d:%s:%s",
				d.sourceName, d.GetLineNumber(), d.GetCharno(), d.GetType().Key(), d.GetDescription())
		}
		return out
	}

	rawMis := compiler.GetTypeValidator().GetMismatches()
	mismatches := make([]string, len(rawMis))
	for i, m := range rawMis {
		mismatches[i] = fmt.Sprintf("%s<-%s", m.Found, m.Required)
	}

	var nodeTypes []string
	var walk func(n rhino.Node)
	walk = func(n rhino.Node) {
		if jt := n.GetJSType(); jt != nil {
			nodeTypes = append(nodeTypes, fmt.Sprintf("%s:%s", n.GetToken(), jt.String()))
		} else {
			nodeTypes = append(nodeTypes, fmt.Sprintf("%s:<nil>", n.GetToken()))
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c)
		}
	}
	walk(root)

	return round98InferenceSnapshot{
		errors:       fmtDiags(compiler.GetErrors()),
		warnings:     fmtDiags(compiler.GetWarnings()),
		mismatches:   mismatches,
		typedPercent: compiler.GetErrorManager().GetTypedPercent(),
		nodeTypes:    nodeTypes,
	}
}

func TestRound98_TypeInferenceAndScopeCreator_ParallelFirstOrderAndScopeOrderParity(t *testing.T) {
	externs := []*SourceFile{
		SourceFile_FromCode("ext0.js", `
/** @const */ var goog = {};
/** @param {string} ns */ goog.provide = function(ns) {};
/** @param {string} ns */ goog.require = function(ns) {};
/** @param {string} ns */ goog.requireType = function(ns) {};
/** @param {string} ns */ goog.module = function(ns) {};
/** @param {function(?):?} fn */ goog.loadModule = function(fn) {};
`),
		SourceFile_FromCode("ext1.js", `
/** @enum {string} */
var ExtKind = { ALPHA: 'a', BETA: 'b' };
`),
	}

	scripts := []*SourceFile{
		SourceFile_FromCode("t0_enums.js", `
/** @enum {number} */
var GlobalStatus = { OK: 1, ERR: 2 };
/** @typedef {{id: number, status: GlobalStatus}} */
var StatusRecord;
`),
		SourceFile_FromCode("t1_enum_aliases.js", `
/** @const */
var StatusAlias = GlobalStatus;
/** @const */
var ExtKindAlias = ExtKind;

/** @param {StatusAlias} s @param {ExtKindAlias} k @return {number} */
function checkNonNullableAliases(s, k) {
  /** @type {!GlobalStatus} */ var nonNullS = s;
  /** @type {!ExtKind} */ var nonNullK = k;
  return nonNullS + nonNullK.length;
}

function localShadowEnumTest() {
  /** @type {? {OK: number}} */
  var GlobalStatus = null;
  /** @const */
  var LocalNotEnum = GlobalStatus;
  /** @enum {number} */
  var InnerEnum = { ONE: 1 };
  /** @const */
  var InnerEnumAlias = InnerEnum;
  /** @param {InnerEnumAlias} v @return {!InnerEnum} */
  function takeInner(v) { return v; }
  return LocalNotEnum ? takeInner(InnerEnum.ONE) : InnerEnum.ONE;
}
`),
		SourceFile_FromCode("t2_escaped_vars.js", `
function escapedAndMultiAssigned(flag) {
  var single = 42;
  var multi = 1;
  if (flag) {
    multi = 'str';
  }
  var escapedSimple = 10;
  var escapedHolder = { count: 0 };
  var arrowFn = (delta) => delta + single;
  function mutator() {
    escapedSimple = 'changed_in_closure';
    escapedHolder.count = 'bad_count';
    return arrowFn(1);
  }
  mutator();
  return { single: single, multi: multi, esc: escapedSimple, cnt: escapedHolder.count };
}
`),
		SourceFile_FromCode("t3_classes_and_fields.js", `
class BaseService {
  /** @type {number} */
  baseField = 1;
  static staticCount = 0;
  static {
    BaseService.staticCount += 1;
  }
  /** @param {string} name */
  constructor(name) {
    /** @type {string} */
    this.name = name;
  }
}

const computedKey = 'dynamicField';
class SubService extends BaseService {
  [computedKey] = this.baseField + 2;
  /** @param {string} name */
  constructor(name) {
    super(name);
  }
  /** @return {number} */
  total() {
    for (let i = 0; i < 3; i++) {
      switch (i) {
        case 0: {
          let step = this.baseField;
          if (step > 0) return step;
        }
      }
    }
    return 0;
  }
}
`),
		SourceFile_FromCode("t4_iife_and_provides.js", `
goog.provide('app.ns');
(function() {
  /** @param {number} n @return {string} */
  app.ns.formatNum = function(n) { return '' + n; };
})();
/** @type {string} */
var formattedOk = app.ns.formatNum(123);
/** @type {number} */
var intentionalWarning = app.ns.formatNum('wrong_arg_type');
`),
		SourceFile_FromCode("t5_bundled_module.js", `
goog.loadModule(function(exports) {
  'use strict';
  goog.module('app.modA');
  /** @param {number} x @return {number} */
  function doubleIt(x) { return x * 2; }
  exports = { doubleIt: doubleIt };
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('app.modB');
  const modA = goog.require('app.modA');
  const unusedType = goog.requireType('app.modA');
  /** @type {number} */
  const res = modA.doubleIt(21);
  exports.res = res;
  return exports;
});
`),
	}

	seq := runRound98TypeCheckFull(t, externs, scripts, 1)
	par := runRound98TypeCheckFull(t, externs, scripts, 4)

	if !slices.Equal(seq.errors, par.errors) {
		t.Fatalf("errors mismatch:\nseq=%v\npar=%v", seq.errors, par.errors)
	}
	if !slices.Equal(seq.warnings, par.warnings) {
		t.Fatalf("warnings mismatch:\nseq=%v\npar=%v", seq.warnings, par.warnings)
	}
	if len(seq.warnings) == 0 {
		t.Fatalf("expected intentional type warning in test fixture")
	}
	if !slices.Equal(seq.mismatches, par.mismatches) {
		t.Fatalf("mismatches mismatch:\nseq=%v\npar=%v", seq.mismatches, par.mismatches)
	}
	if math.Abs(seq.typedPercent-par.typedPercent) > 1e-9 {
		t.Fatalf("typedPercent mismatch: seq=%v par=%v", seq.typedPercent, par.typedPercent)
	}
	if !slices.Equal(seq.nodeTypes, par.nodeTypes) {
		t.Fatalf("AST node types mismatch (len seq=%d par=%d)", len(seq.nodeTypes), len(par.nodeTypes))
	}
}

func TestRound98_ControlFlowAnalysis_AstPositionPeakReuse(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("function hugeFn(x) {\n  var acc = 0;\n")
	for i := 0; i < 220; i++ {
		fmt.Fprintf(&sb, "  if (x > %d) { acc += %d; }\n", i, i+1)
	}
	sb.WriteString("  return acc;\n}\n")
	for i := 0; i < 12; i++ {
		fmt.Fprintf(&sb, "function smallFn%d(a) { try { if (a) return a + %d; } finally { a = 0; } return 0; }\n", i, i)
	}

	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	opts := NewCompilerOptions()
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("ext.js", "")},
		[]*SourceFile{SourceFile_FromCode("cfg.js", sb.String())},
		opts,
	)
	root := compiler.ParseInputs()
	script := root.GetSecondChild().GetFirstChild()

	var reusable ReusableCfg
	fnCount := 0
	for fn := script.GetFirstChild(); !fn.IsNil(); fn = fn.GetNext() {
		if !fn.IsFunction() {
			continue
		}
		fnCount++
		gotCfg := reusable.Compute(compiler, fn, false, true)

		wantCfg := ComputeControlFlowGraph(compiler, fn, false, true)

		if len(gotCfg.NodesSlice()) != len(wantCfg.NodesSlice()) {
			t.Fatalf("fn #%d node count mismatch: got=%d want=%d",
				fnCount, len(gotCfg.NodesSlice()), len(wantCfg.NodesSlice()))
		}
		if fnCount == 1 && (reusable.cfa.astPosition != nil || reusable.cfa.astPositionPeak != 0) {
			t.Fatalf("expected hugeFn to drop oversized astPosition map immediately, got nil=%v peak=%d",
				reusable.cfa.astPosition == nil, reusable.cfa.astPositionPeak)
		}
		if fnCount >= 2 && (reusable.cfa.astPosition == nil || reusable.cfa.astPositionPeak == 0 || reusable.cfa.astPositionPeak > cfaAstPosReuseMaxSize) {
			t.Fatalf("expected compact astPosition map reuse on smallFn #%d, got nil=%v peak=%d",
				fnCount, reusable.cfa.astPosition == nil, reusable.cfa.astPositionPeak)
		}
		reusable.release()
	}
	if fnCount != 13 {
		t.Fatalf("expected 13 functions, visited %d", fnCount)
	}
}
