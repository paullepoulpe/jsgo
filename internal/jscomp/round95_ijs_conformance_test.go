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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type round95IjsCase struct {
	CompilerTestCaseBase
}

func newRound95IjsCase() *round95IjsCase {
	c := &round95IjsCase{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func (c *round95IjsCase) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.AllowExternsChanges()
	c.IgnoreWarnings(
		ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE,
		ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE,
	)
}

func (c *round95IjsCase) GetProcessor(compiler *Compiler) CompilerPass {
	return NewConvertToTypedInterface(compiler)
}

func (c *round95IjsCase) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

func TestRound95IjsConformance_IjsShapesAndEdgeCases(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound95IjsCase() })
}

func (c *round95IjsCase) TestMultiDeclaratorAndExportedDestructuring() {
	c.Test(
		"const a = 1, /** string */ b = init(); const {c, d = init()} = fn();",
		"/** @const {number} */ var a; /** @const {string} */ var b; /** @const {UnusableType} */ var c; /** @const {UnusableType} */ var d;",
	)

	c.Test(
		"let [x = compute(), ...rest] = makeArr();",
		"/** @const {UnusableType} */ var x; /** @const {UnusableType} */ var rest;",
	)

	// Exported destructuring (`export const {b, c} = fn();`) must wrap each split
	// VAR in its own single-child EXPORT node with valid source location info
	// (no JSC_MISSING_LINE_INFO or multi-child EXPORT crash in CodePrinter).
	c.DisableCompareJsDoc()
	c.Test(
		"export const {b, c = init()} = fn();",
		"export var b; export var c;",
	)
	c.compareJsDoc = true
	script := ast.MustRoot(c.GetLastCompiler().GetJsRoot()).MustFirstItem().MustScript()
	modBody := script.MustFirstItem().MustModuleBody()
	exportNode := modBody.MustFirstStatement().Node()
	AssertThat(exportNode.IsExport()).IsTrue()
	AssertThat(exportNode.GetJSDocInfo() != nil && exportNode.GetJSDocInfo().IsConstant()).IsTrue()
}

func (c *round95IjsCase) TestEmptyDestructuringAndGoogRequirePreservation() {
	c.Test("const {} = sideEffect(); let [] = otherSideEffect();", "")

	c.TestSame(
		`goog.module('m');
const {} = goog.require('dep');
`)
}

func (c *round95IjsCase) TestDuplicateAliasDestructuringRemoval() {
	// Non-trailing duplicate in ARRAY_PATTERN must leave an elision hole rather
	// than detaching and shifting `Bar` from index 1 to index 0.
	c.Test(
		`/** @const {number} */ var Foo = 1;
const [Foo, Bar] = ns.pair;`,
		`/** @const {number} */ var Foo;
const [, Bar] = ns.pair;`,
	)

	// When all elements in an ARRAY_PATTERN are duplicate aliases, delete the
	// entire destructuring statement.
	c.Test(
		`/** @const {number} */ var Foo = 1;
/** @const {string} */ var Bar = 'x';
const [Foo, Bar] = ns.pair;`,
		`/** @const {number} */ var Foo;
/** @const {string} */ var Bar;`,
	)

	// Duplicate alias with default value and object rest in OBJECT_PATTERN.
	c.Test(
		`/** @const {number} */ var Foo = 1;
const {Foo = 0, Bar} = ns.obj;`,
		`/** @const {number} */ var Foo;
const {Bar} = ns.obj;`,
	)

	c.Test(
		`/** @const {number} */ var Foo = 1;
const {Bar, ...Foo} = ns.obj;`,
		`/** @const {number} */ var Foo;
const {Bar} = ns.obj;`,
	)
}

func (c *round95IjsCase) TestUntypedRestAndOptionalConstructorParams() {
	c.Test(
		`class C {
  /**
   * @param {?=} opt
   * @param {...} args
   */
  constructor(opt, ...args) {
    /** @const */ this.opt = opt;
    /** @const */ this.args = args;
  }
}`,
		`class C {
  /**
   * @param {?=} opt
   * @param {...} args
   */
  constructor(opt, ...args) {}
}
/** @const {!Array<?>} */ C.prototype.args;
/** @const {?|undefined} */ C.prototype.opt;`,
	)
}

func (c *round95IjsCase) TestEnumAndPolymerBehaviorNonStringKeyMembers() {
	c.Test(
		`/** @enum {string} */
const E = {
  A: 'ok',
  LONG: '0123456789abc',
  ...Other,
  get B() { return 'b'; }
};`,
		`/** @enum {string} */
const E = {
  A: 'ok',
  LONG: '01234567..',
  get B() {}
};`,
	)

	c.Test(
		`/** @polymerBehavior */
const MyBehavior = {
  properties: {
    foo: {type: String, value: 'x', ...extra},
    get bar() { return 1; }
  }
};`,
		`/** @polymerBehavior */
const MyBehavior = {
  properties: {
    foo: {type: String},
    get bar() {}
  }
};`,
	)
}

func (c *round95IjsCase) TestES2022StaticBlocksComputedFieldsExportedClassesAndMalformedCalls() {
	c.Test(
		`class Foo {
  static {
    var leaked = init();
    sideEffect();
  }
  #priv;
  [Symbol.iterator] = makeIter();
  [Symbol.asyncIterator](x = init()) { return x; }
}`,
		`class Foo {
  /** @const @type {UnusableType} */
  #priv;
  [Symbol.iterator];
  [Symbol.asyncIterator](x = void 0) {}
}`,
	)

	c.Test(
		`export class ExportedCls {
  method() { return compute(); }
}
export default class {
  x = init();
  method() { return compute(); }
}`,
		`export class ExportedCls {
  method() {}
}
export default class {
  method() {}
}`,
	)

	c.Test(
		`export default function() { return compute(); }`,
		`export default function() {}`,
	)

	// Optional chaining expressions are runtime-only; optional-chain RHS is not a static alias.
	c.Test(
		`a?.b;
a?.b();
a?.[0];
const val = a?.b;`,
		`/** @const {UnusableType} */ var val;`,
	)

	// Malformed 0-arg and 1-arg Closure primitives must not panic or corrupt literal arguments.
	c.Test(
		`goog.provide();
goog.require();
goog.requireType();
/** @define {boolean} */ goog.define();
/** @define {boolean} */ goog.define('flag.only');
/** @define {boolean} */ goog.define('flag.withDefault', true);`,
		`goog.provide();
goog.require();
goog.requireType();
/** @define {boolean} */ goog.define();
/** @define {boolean} */ goog.define('flag.withDefault', false);`,
	)
}

func TestRound95IjsConformance_DirectNilNodeAndDetachedGuards(t *testing.T) {
	if got := jsdocForRhs(rhino.NilNode, nil); got != nil {
		t.Fatalf("jsdocForRhs(NilNode) = %v, want nil", got)
	}
	bareCast := rhino.NewNodeWithChild(rhino.Token_CAST, rhino.IR_Name("unknownExpr"))
	if got := jsdocForRhs(bareCast, nil); got != nil {
		t.Fatalf("jsdocForRhs(bare CAST) = %v, want nil", got)
	}
	emptyEquals := rhino.NewNodeWithChild(rhino.Token_EQUALS, rhino.IR_Empty())
	if got := cloneJSDocOperandOrUnknown(emptyEquals); got.GetToken() != rhino.Token_QMARK {
		t.Fatalf("cloneJSDocOperandOrUnknown(EQUALS(EMPTY)) = %v, want QMARK", got.GetToken())
	}
	if got := jsdocForName(nil, nil); got != nil {
		t.Fatalf("jsdocForName(nil) = %v, want nil", got)
	}
	if val, ok := inferConformanceStringValue(nil, rhino.NilNode, nil); ok || val != "" {
		t.Fatalf("inferConformanceStringValue(NilNode) = (%q, %v), want (\"\", false)", val, ok)
	}

	var rule conformanceAbstractRule
	if rule.isTypeImmediatelyTightened(rhino.NilNode) {
		t.Fatal("isTypeImmediatelyTightened(NilNode) = true, want false")
	}
	detachedName := rhino.IR_Name("x")
	if rule.isTypeImmediatelyTightened(detachedName) {
		t.Fatal("isTypeImmediatelyTightened(detached) = true, want false")
	}
	if rule.isNodeUsed(rhino.NilNode) {
		t.Fatal("isNodeUsed(NilNode) = true, want false")
	}
	if rule.isNodeUsed(detachedName) {
		t.Fatal("isNodeUsed(detached) = true, want false")
	}
	detachedAssignLhs := rhino.IR_Name("lhs")
	_ = rhino.IR_Assign(detachedAssignLhs, rhino.IR_Number(1))
	if rule.isNodeUsed(detachedAssignLhs) {
		t.Fatal("isNodeUsed(lhs of detached ASSIGN) = true, want false")
	}

	preconditions := []ConformanceRules_BannedProperty_RequirementPrecondition{
		ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY,
		ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE,
		ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE,
		ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ,
		ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL,
	}
	detachedGetProp := rhino.IR_Getprop(rhino.IR_Name("obj"), "prop")
	for _, pre := range preconditions {
		if pre.ShouldCheck(rhino.NilNode) {
			t.Fatalf("%s.ShouldCheck(NilNode) = true, want false", pre)
		}
		// Detached GETPROP has no parent; READ/WRITE/CALL preconditions must not panic.
		_ = pre.ShouldCheck(detachedGetProp)
	}

	var bp ConformanceRules_BannedProperty
	if ty := bp.extractType(rhino.NilNode); ty != nil {
		t.Fatalf("extractType(NilNode) = %v, want nil", ty)
	}
	detachedKey := rhino.IR_StringKey("k")
	if ty := bp.extractType(detachedKey); ty != nil {
		t.Fatalf("extractType(detached STRING_KEY) = %v, want nil", ty)
	}

	sourcelessScript := rhino.NewNode(rhino.Token_SCRIPT)
	depRule := &conformanceBannedDependencyRule{paths: []string{"bad/"}}
	if res := depRule.checkNode(nil, rhino.NilNode); res != ConformanceRules_ConformanceResult_CONFORMANCE {
		t.Fatalf("BannedDependency(NilNode) = %v, want CONFORMANCE", res)
	}
	if res := depRule.checkNode(nil, sourcelessScript); res != ConformanceRules_ConformanceResult_CONFORMANCE {
		t.Fatalf("BannedDependency(sourceless SCRIPT) = %v, want CONFORMANCE", res)
	}

	depRegexPattern, err := buildConformancePattern([]string{"bad/.*"})
	if err != nil {
		t.Fatalf("buildConformancePattern: %v", err)
	}
	depRegexRule := &conformanceBannedDependencyRegexRule{pathRegexp: depRegexPattern}
	if res := depRegexRule.checkNode(nil, sourcelessScript); res != ConformanceRules_ConformanceResult_CONFORMANCE {
		t.Fatalf("BannedDependencyRegex(sourceless SCRIPT) = %v, want CONFORMANCE", res)
	}

	if isRootOfQualifiedNameGlobal(nil, rhino.NilNode) {
		t.Fatal("isRootOfQualifiedNameGlobal(NilNode) = true, want false")
	}

	if isStaticClassMethod(rhino.NilNode) {
		t.Fatal("isStaticClassMethod(NilNode) = true, want false")
	}
	detachedFn := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block())
	if isStaticClassMethod(detachedFn) {
		t.Fatal("isStaticClassMethod(detached FUNCTION) = true, want false")
	}

	// Non-getprop callee `(fn1 || fn2)()` must return false when `thisType` is required
	// without panicking on positional child access.
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	fnTypeExpr, err := evaluateConformanceTypeString(compiler, "function(this:Object)")
	if err != nil {
		t.Fatalf("evaluateConformanceTypeString: %v", err)
	}
	orCallee := rhino.NewNodeWith2Children(rhino.Token_OR, rhino.IR_Name("fn1"), rhino.IR_Name("fn2"))
	callNode := rhino.IR_Call(orCallee)
	if validateConformanceCallThis(ast.MustCallOrNew(callNode), fnTypeExpr.ToMaybeFunctionType(), false) {
		t.Fatal("validateConformanceCallThis((fn1 || fn2)()) = true, want false")
	}
}

type round95ConformanceCase struct {
	CompilerTestCaseBase

	configuration string
}

func newRound95ConformanceCase() *round95ConformanceCase {
	c := &round95ConformanceCase{}
	externs := checkConformanceTest_EXTERNS + `
/** @constructor */ function Document() {}
/** @param {string=} tag @return {!Element} */ Document.prototype.createElement = function(tag) {};
/** @param {string=} cmd @param {boolean=} ui @param {string=} val */
Document.prototype.execCommand = function(cmd, ui, val) {};
/** @constructor */ function Element() {}
/** @param {string=} name @param {*=} value */ Element.prototype.setAttribute = function(name, value) {};
/** @param {?string} ns @param {string} name @param {*} value */
Element.prototype.setAttributeNS = function(ns, name, value) {};
/** @const {!Document} */ var document;
/** @const */ goog.dom = {};
/** @param {...*} var_args @return {!Element} */ goog.dom.createElement = function(var_args) {};
/** @param {...*} var_args @return {!Element} */ goog.dom.createDom = function(var_args) {};
/** @constructor */ goog.dom.DomHelper = function() {};
/** @param {...*} var_args @return {!Element} */ goog.dom.DomHelper.prototype.createElement = function(var_args) {};
/** @param {...*} var_args @return {!Element} */ goog.dom.DomHelper.prototype.createDom = function(var_args) {};
`
	c.InitCompilerTestCaseBaseWithExterns(c, externs)
	return c
}

func (c *round95ConformanceCase) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableClosurePass()
	c.EnableClosurePassForExpected()
	c.EnableRewriteClosureCode()
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
}

func (c *round95ConformanceCase) GetProcessor(compiler *Compiler) CompilerPass {
	builder := ConformanceConfig_NewBuilder()
	jutil.Try(func() {
		TextFormat_Merge(c.configuration, builder)
	}, func(e *TextFormat_ParseException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	reportingMode := CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	return NewCheckConformanceWithConfigs(compiler, []*ConformanceConfig{builder.Build()}, &reportingMode)
}

func TestRound95IjsConformance_ConformanceOptChainSpreadAndSecurity(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound95ConformanceCase() })
}

func (c *round95ConformanceCase) TestBannedNameAndPropertyOptionalChaining() {
	c.configuration = `
requirement: {
  type: BANNED_NAME_CALL
  value: 'eval'
  error_message: 'eval call is banned'
}
requirement: {
  type: BANNED_PROPERTY_CALL
  value: 'Box.prototype.danger'
  error_message: 'danger call is banned'
}
`
	c.TestWarning("eval();", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning("eval?.();")
	c.TestNoWarning("const ref = eval;")
	c.TestWarning(
		`class Box { danger() {} }
/** @param {!Box} b */
function run(b) { b.danger(); }`,
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestNoWarning(
		`class Box { danger() {} }
/** @param {!Box} b */
function run(b) { b.danger?.(); }`,
	)

	// Java ConformanceRules.BannedProperty / BannedPropertyCall only match
	// STRING_KEY, GETPROP, GETELEM, and COMPUTED_PROP (not OPTCHAIN_GETPROP / OPTCHAIN_GETELEM).
	c.configuration = `
requirement: {
  type: BANNED_PROPERTY
  value: 'Arguments.prototype.callee'
  value: 'DictBox.prototype.secret'
  error_message: 'banned prop'
}
requirement: {
  type: BANNED_PROPERTY_CALL
  value: 'Box.prototype.danger'
  error_message: 'danger call is banned'
}
`
	c.TestWarning("function f() { return arguments.callee; }", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning("function f() { return arguments?.callee; }")
	c.TestWarning(
		`/** @dict */ class DictBox {}
/** @param {!DictBox} d */
function f(d) { return d['secret']; }`,
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestNoWarning(
		`/** @dict */ class DictBox {}
/** @param {?DictBox} d */
function f(d) { return d?.['secret']; }`,
	)
	c.TestNoWarning(
		`class Box { danger() {} }
/** @param {?Box} b */
function run(b) { b?.danger(); }`,
	)
}

func (c *round95ConformanceCase) TestRestrictedCallOptionalChainingAndSpreadArgs() {
	c.configuration = `
requirement: {
  type: RESTRICTED_NAME_CALL
  value: 'C.m:function(number)'
  error_message: 'C.m requires a single number'
}
requirement: {
  type: RESTRICTED_METHOD_CALL
  value: 'C.prototype.inst:function(number)'
  error_message: 'inst requires a single number'
}
requirement: {
  type: RESTRICTED_METHOD_CALL
  value: 'TimerHost.prototype.setTimeout:function(Function, ...?)'
  error_message: 'setTimeout requires leading Function'
}
`
	decl := `/** @constructor */ function C() {}
/** @param {...*} a */ C.m = function(a) {};
/** @param {...*} a */ C.prototype.inst = function(a) {};
/** @constructor */ function TimerHost() {}
/** @param {...*} a */ TimerHost.prototype.setTimeout = function(a) {};
`
	c.TestWarning(decl+"C.m('bad');", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(decl + "C.m?.(1);")
	c.TestNoWarning(decl + "C.m?.('bad');")
	c.TestNoWarning(decl + "C.m.call?.(C, 'bad');")
	c.TestNoWarning(decl + "C.m?.call(C, 'bad');")

	// Spread arguments match Java ConformanceRules.validateParameterList:
	// SPREAD nodes carry null JSType (mapped to UNKNOWN_TYPE) while explicit
	// positional arguments and max-arity bounds are still validated normally.
	c.TestNoWarning(
		decl + "function okTimeout(/** !TimerHost */ host, /** !Function */ fn, /** !Array<*> */ params) { host.setTimeout(fn, 10, ...params); }",
	)
	c.TestWarning(
		decl+"function badTimeout(/** !TimerHost */ host, /** !Array<*> */ params) { host.setTimeout('notFn', 10, ...params); }",
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestWarning(
		decl+"/** @type {!Array<number>} */ const nums = [2]; C.m(1, ...nums);",
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestWarning(
		decl+"/** @type {!Array<*>} */ const mixed = [C, 'bad']; C.m.call(...mixed);",
		CheckConformance_CONFORMANCE_VIOLATION,
	)

	c.TestNoWarning(decl + "function ok(/** ?C */ x) { x?.inst(1); }")
	c.TestNoWarning(decl + "function ok1(/** ?C */ x) { x?.inst('bad'); }")
	c.TestNoWarning(decl + "function ok2(/** !C */ x) { x.inst?.('bad'); }")
	c.TestNoWarning(decl + "function ok3(/** ?C */ x) { x?.inst.call(x, 'bad'); }")
	c.TestWarning(
		decl+"function bad4(/** !C */ x, /** !Array<number> */ rest) { x.inst(1, ...rest); }",
		CheckConformance_CONFORMANCE_VIOLATION,
	)
}

func (c *round95ConformanceCase) TestDomSecurityOptionalChainingSpreadAndTemplates() {
	c.configuration = `
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateElement'
  value: 'script'
  error_message: 'BanCreateElement'
}
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanCreateDom'
  value: 'iframe.src'
  report_loose_type_violations: true
  error_message: 'BanCreateDom'
}
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute'
  value: 'href'
  value: 'src'
  report_loose_type_violations: true
  error_message: 'BanElementSetAttribute'
}
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes'
  value: 'Wrapper.prototype.attr'
  error_message: 'BanSettingAttributes'
}
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanExecCommand'
  value: 'insertHTML'
  error_message: 'BanExecCommand'
}
`
	// BanCreateElement: optional-chain parity and empty/spread non-crashes.
	c.TestWarning("function f(/** !Document */ doc) { doc.createElement('script'); }", CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning("function f(/** ?Document */ doc) { doc?.createElement('script'); }")
	c.TestNoWarning("goog.dom?.createDom('script');")
	c.TestNoWarning("function f(/** !Document */ doc, /** !Array<string> */ tags) { doc.createElement(...tags); }")
	c.TestNoWarning("function f(/** !Document */ doc) { doc.createElement(); }")

	// BanCreateDom: optional-chain parity, spread attrs under report_loose_type_violations, and constant template literal.
	c.TestNoWarning("function f(/** string */ u) { goog.dom?.createDom('iframe', {'src': u}); }")
	c.TestNoWarning("function f(/** ?goog.dom.DomHelper */ dh, /** string */ u) { dh?.createDom('iframe', {'src': u}); }")
	c.TestWarning(
		"function f(/** !Array<!Object> */ rest) { goog.dom.createDom('iframe', ...rest); }",
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
	)
	c.TestNoWarning("goog.dom.createDom('iframe', {'src': `https://safe.example`});")

	// BanElementSetAttribute: optional-chain parity, template literal attr name, spread args, and 0-arg safety.
	c.TestNoWarning("function f(/** ?Element */ el, /** string */ u) { el?.setAttribute(`href`, u); }")
	c.TestNoWarning("function f(/** ?Element */ el, /** string */ u) { el?.setAttributeNS(null, 'src', u); }")
	c.TestNoWarning("function f(/** ?Element */ el, /** !Array<*> */ args) { el?.setAttribute(...args); }")
	c.TestWarning(
		"function f(/** !Element */ el, /** string */ u) { el.setAttribute(`href`, u); }",
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestWarning(
		"function f(/** !Element */ el, /** !Array<*> */ args) { el.setAttribute(...args); }",
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestNoWarning("function f(/** !Element */ el) { el.setAttribute(); }")

	// BanSettingAttributes: optional-chain parity, event-handler attr, spread args, and 1-arg non-spread safety.
	wrapperDecl := `class Wrapper {
  /** @param {...*} args */ attr(...args) {}
}
`
	c.TestNoWarning(wrapperDecl + "function f(/** ?Wrapper */ w, /** string */ u) { w?.attr('href', u); }")
	c.TestNoWarning(wrapperDecl + "function f(/** !Wrapper */ w, /** function() */ fn) { w.attr?.('onclick', fn); }")
	c.TestNoWarning(wrapperDecl + "function f(/** ?Wrapper */ w, /** !Array<*> */ args) { w?.attr(...args); }")
	c.TestWarning(
		wrapperDecl+"function f(/** !Wrapper */ w, /** function() */ fn) { w.attr('onclick', fn); }",
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestWarning(
		wrapperDecl+"function f(/** !Wrapper */ w, /** !Array<*> */ args) { w.attr(...args); }",
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestNoWarning(wrapperDecl + "function f(/** !Wrapper */ w) { w.attr('href'); }")

	// BanExecCommand: optional-chain parity, spread args, and 0-arg safety.
	c.TestNoWarning("function f(/** ?Document */ doc, /** string */ html) { doc?.execCommand('insertHTML', false, html); }")
	c.TestNoWarning("function f(/** ?Document */ doc, /** !Array<*> */ args) { doc?.execCommand(...args); }")
	c.TestWarning(
		"function f(/** !Document */ doc, /** !Array<*> */ args) { doc.execCommand(...args); }",
		CheckConformance_CONFORMANCE_VIOLATION,
	)
	c.TestNoWarning("function f(/** !Document */ doc) { doc.execCommand(); }")
}

func (c *round95ConformanceCase) TestTypeRulesAndStaticComputedMethodThis() {
	// Optional chaining safely avoids BanNullDeref on nullable receivers/callees.
	c.configuration = `
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanNullDeref'
  error_message: 'BanNullDeref'
}
`
	c.TestNoWarning(`class Holder { constructor() { /** @type {number} */ this.x = 1; } }
function safe(/** ?Holder */ h, /** ?function(): number */ fn) {
  return (h?.x ?? 0) + (fn?.() ?? 0);
}`)
	c.TestWarning(`class Holder { constructor() { /** @type {number} */ this.x = 1; } }
function unsafe(/** ?Holder */ h) {
  return h.x;
}`, CheckConformance_CONFORMANCE_VIOLATION)

	// BanUnknownDirectThisPropsReferences catches `this.prop` when used, but not `this?.prop`.
	c.configuration = `
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownDirectThisPropsReferences'
  error_message: 'BanUnknownDirectThis'
}
`
	c.TestWarning(`class WithUnknownThisProp {
  constructor() { this.prop; }
  read() { alert(this.prop); }
}`, CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(`class WithOptChainThisProp {
  constructor() { this.prop; }
  read() { alert(this?.prop); }
}`)

	// BanUnknownTypedClassPropsReferences catches `inst.unknownProp`, but not `inst?.unknownProp` or `event?.data?.focusNext`.
	c.configuration = `
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownTypedClassPropsReferences'
  error_message: 'BanUnknownClassProp'
}
`
	c.TestWarning(`class TypedCls {}
function read(/** !TypedCls */ inst) { alert(inst.unknownProp); }
`, CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(`class TypedCls {}
class CustomEvt { constructor() { /** @type {?} */ this.data = null; } }
function read(/** ?TypedCls */ inst, /** ?CustomEvt */ event) {
  alert(inst?.unknownProp);
  alert(event?.data?.focusNext);
}
`)

	// BanUnresolvedType catches `u.prop`, but not `u?.prop`.
	c.configuration = `
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnresolvedType'
  error_message: 'BanUnresolvedType'
}
`
	c.TestWarning(`goog.forwardDeclare('UnresolvedThing');
/** @param {!UnresolvedThing} u */
function readUnresolved(u) { return u.prop; }
`, CheckConformance_CONFORMANCE_VIOLATION)
	c.TestNoWarning(`goog.forwardDeclare('UnresolvedThing');
/** @param {?UnresolvedThing} u */
function readUnresolved(u) { return u?.prop; }
`)

	// BanStaticThis catches `this` inside static computed class methods `static [k]() { ... }`.
	c.configuration = `
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanStaticThis'
  error_message: 'BanStaticThis'
}
`
	c.TestWarning(`const k = 'run';
/** @dict */
class StaticComp {
  static [k]() { return this; }
}`, CheckConformance_CONFORMANCE_VIOLATION)
}

func TestConvertToTypedInterfaceParallelDeterministicParity(t *testing.T) {
	files := []*SourceFile{
		SourceFile_FromCode("mod_a.js", `goog.module('pkg.a');
const b = goog.require('pkg.b');
var unused = 1, /** @type {number} */ EXPORTED_NUM = 42;
class Alpha {
  constructor() {
    /** @const */ this.num = EXPORTED_NUM;
    /** @const */ this.dep = b;
  }
  /** @param {string} s @return {number} */
  compute(s) { return s.length + this.num; }
}
exports = {Alpha, EXPORTED_NUM};
`),
		SourceFile_FromCode("mod_b.js", `goog.module('pkg.b');
/** @enum {number} */
const Mode = {FAST: 1, SAFE: 2};
/** @typedef {{mode: !Mode, label: string}} */
let Config;
exports = {Mode, Config};
`),
		SourceFile_FromCode("legacy_ns.js", `goog.provide('legacy.ns');
goog.provide('legacy.ns.Sub');
/** @const */
legacy.ns.Sub = class {
  constructor() {
    /** @type {boolean} */
    this.ready = true;
  }
};
/** @param {!legacy.ns.Sub} x @return {boolean} */
legacy.ns.check = function(x) { return x.ready; };
`),
		SourceFile_FromCode("warn_case.js", `goog.module('pkg.warn');
class Holder {
  constructor() {
    this.untypedProp = unknownFactory();
  }
}
exports = Holder;
`),
	}

	compileWithThreads := func(threads int) (string, []string) {
		opts := NewCompilerOptions()
		opts.SetIncrementalChecks(CompilerOptions_IncrementalCheckMode_GENERATE_IJS)
		opts.SetCodingConvention(NewClosureCodingConvention())
		opts.SetPrettyPrint(true)
		opts.SetPreserveTypeAnnotations(true)
		opts.SetNumParallelThreads(threads)

		comp := NewCompiler()
		var freshInputs []*SourceFile
		for _, f := range files {
			freshInputs = append(freshInputs, SourceFile_FromCode(f.GetName(), f.GetCode()))
		}
		res := comp.Compile([]*SourceFile{}, freshInputs, opts)
		var diags []string
		for _, e := range res.Errors {
			diags = append(diags, e.String())
		}
		for _, w := range res.Warnings {
			diags = append(diags, w.String())
		}
		return comp.ToSource(), diags
	}

	seqOut, seqDiags := compileWithThreads(1)
	for range 4 {
		parOut, parDiags := compileWithThreads(4)
		if parOut != seqOut {
			t.Fatalf("parallel GENERATE_IJS output mismatch:\ngot:\n%s\nwant:\n%s", parOut, seqOut)
		}
		if len(parDiags) != len(seqDiags) {
			t.Fatalf("parallel GENERATE_IJS diag count = %d (%v), want %d (%v)", len(parDiags), parDiags, len(seqDiags), seqDiags)
		}
		for i := range seqDiags {
			if parDiags[i] != seqDiags[i] {
				t.Fatalf("parallel GENERATE_IJS diag[%d] = %q, want %q", i, parDiags[i], seqDiags[i])
			}
		}
	}
}

func TestRound95IjsConformance_ParallelCheckConformanceDeterminism(t *testing.T) {
	cfgText := `
requirement: {
  type: BANNED_NAME
  value: 'execScript'
  error_message: 'no execScript'
}
requirement: {
  type: BANNED_CODE_PATTERN
  value: '/** @param {string} s */ function template(s) { window.alert(s); }'
  error_message: 'no alert pattern'
}
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanUnknownThis'
  error_message: 'no unknown this'
  allowlist: 'allowed_first.js'
}
requirement: {
  type: CUSTOM
  java_class: 'com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute'
  value: 'src'
  report_loose_type_violations: true
  error_message: 'no setAttribute src'
}
requirement: {
  type: RESTRICTED_METHOD_CALL
  value: 'Element.prototype.setAttribute:function(string, string)'
  error_message: 'setAttribute args must be strings'
}
`
	builder := ConformanceConfig_NewBuilder()
	if err := TextFormat_MergeErr(cfgText, builder); err != nil {
		t.Fatalf("TextFormat_MergeErr: %v", err)
	}
	cfg := builder.Build()

	externs := checkConformanceTest_EXTERNS + `
/** @constructor */ function Arguments() {}
/** @constructor */ function Element() {}
/** @param {string=} name @param {*=} value */ Element.prototype.setAttribute = function(name, value) {};
/** @param {string} s */ function execScript(s) {}
/** @param {string} s */ window.alert = function(s) {};
`
	srcs := []*SourceFile{
		SourceFile_FromCode("allowed_first.js", `
/** @const {string} */ var GLOBAL_SRC_ATTR = 'src';
var allowedGlobalThis = this;
function allowedFn() { return this.x; }
`),
		SourceFile_FromCode("s1.js", `
var globalThisHit1 = this;
function badThis1() { return this.x + this.y; }
execScript('1');
`),
		SourceFile_FromCode("s2.js", `
/** @param {function(string):*} execScript */
function shadowedExec(execScript) { return execScript('ok'); }
window.alert('bad');
`),
		SourceFile_FromCode("s3.js", `
/** @param {!Element} el */
function domOps(el) {
  el.setAttribute(GLOBAL_SRC_ATTR, 'https://example.com');
  el.setAttribute('data-id', 123);
}
`),
		SourceFile_FromCode("s4.js", `
var globalThisHit2 = this;
function badThis2() { return this.z; }
`),
	}

	runWithThreads := func(threads int) []string {
		opts := NewCompilerOptions()
		opts.SetCheckTypes(true)
		opts.SetChecksOnly(true)
		opts.SetConformanceConfigs([]*ConformanceConfig{cfg})
		opts.SetNumParallelThreads(threads)

		comp := NewCompiler()
		extInputs := []*SourceFile{SourceFile_FromCode("externs.js", externs)}
		var freshSrcs []*SourceFile
		for _, f := range srcs {
			freshSrcs = append(freshSrcs, SourceFile_FromCode(f.GetName(), f.GetCode()))
		}
		res := comp.Compile(extInputs, freshSrcs, opts)
		var diags []string
		for _, e := range res.Errors {
			diags = append(diags, e.String())
		}
		for _, w := range res.Warnings {
			diags = append(diags, w.String())
		}
		return diags
	}

	seqDiags := runWithThreads(1)
	if len(seqDiags) < 4 {
		t.Fatalf("expected >= 4 conformance diagnostics, got %d: %v", len(seqDiags), seqDiags)
	}
	for range 4 {
		parDiags := runWithThreads(4)
		if len(parDiags) != len(seqDiags) {
			t.Fatalf("parallel conformance diag count = %d (%v), want %d (%v)", len(parDiags), parDiags, len(seqDiags), seqDiags)
		}
		for i := range seqDiags {
			if parDiags[i] != seqDiags[i] {
				t.Fatalf("parallel conformance diag[%d] = %q, want %q", i, parDiags[i], seqDiags[i])
			}
		}
	}
}

func TestRound95IjsConformance_ParallelConstCheckAndIjsDeterminism(t *testing.T) {
	externs := checkConformanceTest_EXTERNS + `
/** @const {number} */ var EXT_CONST;
`
	srcs := []*SourceFile{
		SourceFile_FromCode("g1.js", `
/** @const {number} */ var GLOBAL_K = 1;
function f1() {
  /** @const {number} */ var localA = 10;
  localA = 20;
  EXT_CONST = 2;
  /** @const {number} */ var localB = 30;
  localB++;
}
`),
		SourceFile_FromCode("g2.js", `
goog.provide('ns.alpha');
/** @const {number} */ ns.alpha.X = 1;
goog.scope(function() {
  var a = ns.alpha;
  /** @type {number} */
  ns.alpha.Y = a.X;
});
GLOBAL_K = 99;
`),
		SourceFile_FromCode("m1.js", `
goog.module('mod.one');
const Alpha = goog.require('ns.alpha');
/** @const {number} */ exports.val = Alpha.Y;
`),
		SourceFile_FromCode("m2.js", `
goog.module('mod.two');
const Alpha = goog.require('ns.alpha');
const UnusedOne = goog.require('mod.one');
/** @type {number} */ exports.out = Alpha.X;
`),
	}

	runChecks := func(threads int) []string {
		opts := NewCompilerOptions()
		opts.SetCheckTypes(true)
		opts.SetChecksOnly(true)
		opts.SetClosurePass(true)
		opts.SetWarningLevel(DiagnosticGroups_CONST, CheckLevel_WARNING)
		opts.SetNumParallelThreads(threads)

		comp := NewCompiler()
		var fresh []*SourceFile
		for _, f := range srcs {
			fresh = append(fresh, SourceFile_FromCode(f.GetName(), f.GetCode()))
		}
		res := comp.Compile([]*SourceFile{SourceFile_FromCode("externs.js", externs)}, fresh, opts)
		var diags []string
		for _, e := range res.Errors {
			diags = append(diags, e.String())
		}
		for _, w := range res.Warnings {
			diags = append(diags, w.String())
		}
		return diags
	}

	seqDiags := runChecks(1)
	if len(seqDiags) < 4 {
		t.Fatalf("expected >= 4 const diagnostics, got %d: %v", len(seqDiags), seqDiags)
	}
	for range 4 {
		parDiags := runChecks(4)
		if len(parDiags) != len(seqDiags) {
			t.Fatalf("parallel const diag count = %d (%v), want %d (%v)", len(parDiags), parDiags, len(seqDiags), seqDiags)
		}
		for i := range seqDiags {
			if parDiags[i] != seqDiags[i] {
				t.Fatalf("parallel const diag[%d] = %q, want %q", i, parDiags[i], seqDiags[i])
			}
		}
	}

	runIjs := func(threads int) string {
		opts := NewCompilerOptions()
		opts.SetIncrementalChecks(CompilerOptions_IncrementalCheckMode_GENERATE_IJS)
		opts.SetClosurePass(true)
		opts.SetPrettyPrint(true)
		opts.SetNumParallelThreads(threads)

		comp := NewCompiler()
		var fresh []*SourceFile
		for _, f := range srcs {
			fresh = append(fresh, SourceFile_FromCode(f.GetName(), f.GetCode()))
		}
		comp.Compile([]*SourceFile{SourceFile_FromCode("externs.js", externs)}, fresh, opts)
		return comp.ToSource()
	}

	seqIjs := runIjs(1)
	for range 4 {
		parIjs := runIjs(4)
		if parIjs != seqIjs {
			t.Fatalf("parallel GENERATE_IJS mismatch:\ngot:\n%s\nwant:\n%s", parIjs, seqIjs)
		}
	}
}

func TestRound95IjsConformance_ParallelGatherModuleMetadataDeterminism(t *testing.T) {
	externs := `/** @const */ var goog = {};
/** @param {string} ns */ goog.provide = function(ns) {};
/** @param {string} ns */ goog.module = function(ns) {};
goog.module.declareLegacyNamespace = function() {};
/** @param {function(?): ?} fn */ goog.loadModule = function(fn) {};
/** @param {string} ns @return {?} */ goog.require = function(ns) {};
/** @param {string} ns */ goog.requireType = function(ns) {};
`
	srcs := []*SourceFile{
		SourceFile_FromCode("f0.js", `goog.provide('dup.ns'); /** @const */ dup.ns.A = 1;`),
		SourceFile_FromCode("f1.js", `goog.module('mod.leg'); goog.module.declareLegacyNamespace(); exports = {x: 1};`),
		SourceFile_FromCode("f2.js", `goog.provide('dup.ns'); goog.requireType(123);`),
		SourceFile_FromCode("f3.js", `goog.loadModule(function(exports) { 'use strict'; goog.module('mod.leg'); return exports; });`),
		SourceFile_FromCode("f4.js", `goog.module('foo$2etoggles'); exports.TOGGLE_alpha = true;`),
		SourceFile_FromCode("f5.js", `goog.module('use.toggles'); const t = goog.require('foo$2etoggles'); const v = t.TOGGLE_alpha; const bad = t.bad_name; exports = v;`),
	}

	run := func(threads int) []string {
		opts := NewCompilerOptions()
		opts.SetChecksOnly(true)
		opts.SetClosurePass(true)
		opts.SetNumParallelThreads(threads)

		comp := NewCompiler()
		var fresh []*SourceFile
		for _, f := range srcs {
			fresh = append(fresh, SourceFile_FromCode(f.GetName(), f.GetCode()))
		}
		res := comp.Compile([]*SourceFile{SourceFile_FromCode("externs.js", externs)}, fresh, opts)
		var out []string
		for _, e := range res.Errors {
			out = append(out, "E:"+e.String())
		}
		for _, w := range res.Warnings {
			out = append(out, "W:"+w.String())
		}
		if mm := comp.GetModuleMetadataMap(); mm != nil {
			for k, m := range mm.GetModulesByPath().All() {
				out = append(out, "P:"+k+":"+m.ModuleType().String())
			}
			for k, m := range mm.GetModulesByGoogNamespace().All() {
				src, _ := m.RootNode().SourceFileName()
				out = append(out, "N:"+k+":"+m.ModuleType().String()+":"+src)
			}
		}
		return out
	}

	seq := run(1)
	if len(seq) < 8 {
		t.Fatalf("expected >= 8 diagnostics+metadata entries, got %d: %v", len(seq), seq)
	}
	for range 4 {
		par := run(4)
		if len(par) != len(seq) {
			t.Fatalf("parallel count = %d (%v), want %d (%v)", len(par), par, len(seq), seq)
		}
		for i := range seq {
			if par[i] != seq[i] {
				t.Fatalf("entry[%d] = %q, want %q", i, par[i], seq[i])
			}
		}
	}
}

func TestRound95IjsConformance_ParallelVarCheckDeterminism(t *testing.T) {
	const externs = `/** @externs */
var knownExtern = 1;
/** @const */ var dupNs = dupNs || {};
/** @const */ var dupNs = dupNs || {};
orphanedExternProp.foo;
`
	var srcs []*SourceFile
	for i := range 12 {
		idx := string(rune('a' + i))
		var code string
		switch i % 4 {
		case 0:
			code = "function f_" + idx + "() { return sharedMissing + only_" + idx + " + typeof maybeGlobal; }\n"
		case 1:
			code = "var v_" + idx + " = sharedMissing + knownExtern;\n"
		case 2:
			code = "goog.provide('vc.ns." + idx + "');\nvc.ns." + idx + ".val = only_" + idx + " + sharedMissing;\n"
		case 3:
			code = "goog.module('vc.mod." + idx + "');\nexports = function() { return sharedMissing; };\n"
		}
		srcs = append(srcs, SourceFile_FromCode("vc_"+idx+".js", code))
	}

	run := func(threads int) []string {
		opts := NewCompilerOptions()
		opts.SetChecksOnly(true)
		opts.SetClosurePass(true)
		opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
		opts.SetWarningLevel(DiagnosticGroups_UNDEFINED_VARIABLES, CheckLevel_WARNING)
		opts.SetNumParallelThreads(threads)

		comp := NewCompiler()
		var fresh []*SourceFile
		for _, f := range srcs {
			fresh = append(fresh, SourceFile_FromCode(f.GetName(), f.GetCode()))
		}
		res := comp.Compile([]*SourceFile{SourceFile_FromCode("externs.js", externs)}, fresh, opts)
		var out []string
		for _, e := range res.Errors {
			out = append(out, "E:"+e.String())
		}
		for _, w := range res.Warnings {
			out = append(out, "W:"+w.String())
		}
		return out
	}

	seq := run(1)
	if len(seq) < 8 {
		t.Fatalf("expected >= 8 diagnostics, got %d: %v", len(seq), seq)
	}
	for range 4 {
		par := run(4)
		if len(par) != len(seq) {
			t.Fatalf("parallel count = %d (%v), want %d (%v)", len(par), par, len(seq), seq)
		}
		for i := range seq {
			if par[i] != seq[i] {
				t.Fatalf("entry[%d] = %q, want %q", i, par[i], seq[i])
			}
		}
	}
}
