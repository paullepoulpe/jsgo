// Port of com/google/javascript/jscomp/AstAnalyzerTest.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Unit tests for AstAnalyzer.
//
// The JUnit @Parameterized classes become one Go test each, with one
// subtest per case named "#<index> <case>" as the Java parameter name is.

// astAnalyzerTest_AnalysisCase is AstAnalyzerTest.AnalysisCase.
type astAnalyzerTest_AnalysisCase struct {
	expect bool
	js     string
	// Go-only: language-adapter — a nullable Token (hasToken false is null).
	token                rhino.Token
	hasToken             bool
	globalRegExp         bool
	assumeGettersArePure bool
	assumeBuiltinsPure   bool
}

func (c *astAnalyzerTest_AnalysisCase) Expect(b bool) *astAnalyzerTest_AnalysisCase {
	c.expect = b
	return c
}

func (c *astAnalyzerTest_AnalysisCase) Js(s string) *astAnalyzerTest_AnalysisCase {
	c.js = s
	return c
}

func (c *astAnalyzerTest_AnalysisCase) Token(t rhino.Token) *astAnalyzerTest_AnalysisCase {
	c.token = t
	c.hasToken = true
	return c
}

func (c *astAnalyzerTest_AnalysisCase) GlobalRegExp(b bool) *astAnalyzerTest_AnalysisCase {
	c.globalRegExp = b
	return c
}

func (c *astAnalyzerTest_AnalysisCase) AssumeGettersArePure(b bool) *astAnalyzerTest_AnalysisCase {
	c.assumeGettersArePure = b
	return c
}

func (c *astAnalyzerTest_AnalysisCase) AssumeBuiltinsPure(b bool) *astAnalyzerTest_AnalysisCase {
	c.assumeBuiltinsPure = b
	return c
}

func (c *astAnalyzerTest_AnalysisCase) String() string {
	token := "null"
	if c.hasToken {
		token = c.token.String()
	}
	return fmt.Sprintf("%s node in `%s` -> %t", token, c.js, c.expect)
}

func kase() *astAnalyzerTest_AnalysisCase {
	return &astAnalyzerTest_AnalysisCase{assumeBuiltinsPure: true}
}

// astAnalyzerTest_ParseHelper provides methods for parsing and accessing
// the compiler used for the parsing.
type astAnalyzerTest_ParseHelper struct {
	t                             *testing.T
	useTypesForLocalOptimizations bool
	hasGlobalRegexpReferences     bool
	assumeGettersArePure          bool
	assumeBuiltinsPure            bool
	typeRegistry                  *rhino.JSTypeRegistry
	accessorSummary               *AccessorSummary
}

func newAstAnalyzerTest_ParseHelper(t *testing.T) *astAnalyzerTest_ParseHelper {
	accessors := jutil.NewLinkedHashMap[string, AccessorSummary_PropertyAccessKind]()
	accessors.Put("getter", AccessorSummary_PropertyAccessKind_GETTER_ONLY)
	accessors.Put("setter", AccessorSummary_PropertyAccessKind_SETTER_ONLY)
	return &astAnalyzerTest_ParseHelper{
		t:                             t,
		useTypesForLocalOptimizations: false,
		hasGlobalRegexpReferences:     true,
		assumeGettersArePure:          true,
		assumeBuiltinsPure:            true,
		accessorSummary:               AccessorSummary_Create(accessors),
	}
}

// Java: AstAnalyzerTest.ParseHelper.newCompiler()
func (h *astAnalyzerTest_ParseHelper) newCompiler() *Compiler {
	options := NewCompilerOptions()

	// To allow octal literals such as 0123 to be parsed.
	options.SetStrictModeInput(false)
	options.SetWarningLevel(DiagnosticGroups_ES5_STRICT, CheckLevel_OFF)

	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)

	compiler := NewCompiler()
	compiler.InitOptions(options)

	return compiler
}

// Java: AstAnalyzerTest.ParseHelper.parseInternal(String)
func (h *astAnalyzerTest_ParseHelper) parseInternal(js string) rhino.Node {
	h.t.Helper()
	compiler := h.newCompiler()
	n := compiler.ParseTestCode(js)
	if errors := compiler.GetErrors(); len(errors) != 0 {
		h.t.Fatalf("expected to be empty, but was %v", errors)
	}
	h.typeRegistry = compiler.GetTypeRegistry()
	return n
}

// parseFirst parses a string of JavaScript and returns the first node found
// with the given token in a preorder DFS.
func (h *astAnalyzerTest_ParseHelper) parseFirst(token rhino.Token, js string) rhino.Node {
	n := astAnalyzerTest_findFirst(token, h.parseInternal(js))
	if n.IsNil() {
		h.t.Fatalf("no %s node in %q", token, js)
	}
	return n
}

func (h *astAnalyzerTest_ParseHelper) parseCase(kase *astAnalyzerTest_AnalysisCase) rhino.Node {
	h.hasGlobalRegexpReferences = kase.globalRegExp
	h.assumeGettersArePure = kase.assumeGettersArePure
	h.assumeBuiltinsPure = kase.assumeBuiltinsPure
	root := h.parseInternal(kase.js)
	if !kase.hasToken {
		return root.GetFirstChild()
	}
	n := astAnalyzerTest_findFirst(kase.token, root)
	if n.IsNil() {
		h.t.Fatalf("no %s node in %q", kase.token, kase.js)
	}
	return n
}

func (h *astAnalyzerTest_ParseHelper) getAstAnalyzer() *AstAnalyzer {
	return NewAstAnalyzer(
		AstAnalyzer_Options_Builder_().
			SetUseTypesForLocalOptimization(h.useTypesForLocalOptimizations).
			SetHasRegexpGlobalReferences(h.hasGlobalRegexpReferences).
			SetAssumeGettersArePure(h.assumeGettersArePure).
			SetAssumeKnownBuiltinsArePure(h.assumeBuiltinsPure).
			Build(),
		h.typeRegistry,
		h.accessorSummary)
}

// astAnalyzerTest_findFirst does a preorder DFS, returning the first node
// found that has the given token (nil for Optional.empty()).
func astAnalyzerTest_findFirst(token rhino.Token, root rhino.Node) rhino.Node {
	stack := []rhino.Node{root}

	for len(stack) != 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.GetToken() == token {
			return cur
		}

		for child := cur.GetLastChild(); !child.IsNil(); child = child.GetPrevious() {
			stack = append(stack, child)
		}
	}

	return rhino.NilNode
}

func astAnalyzerTest_runCases(
	t *testing.T, cases []*astAnalyzerTest_AnalysisCase, query func(*AstAnalyzer, rhino.Node) bool) {
	for i, kase := range cases {
		t.Run(fmt.Sprintf("#%d %s", i, kase), func(t *testing.T) {
			helper := newAstAnalyzerTest_ParseHelper(t)
			node := helper.parseCase(kase)
			analyzer := helper.getAstAnalyzer()
			if got := query(analyzer, node); got != kase.expect {
				t.Errorf("got %t, want %t", got, kase.expect)
			}
		})
	}
}

// Java: AstAnalyzerTest.MayEffectMutableStateTest.mayEffectMutableState()
func TestAstAnalyzer_MayEffectMutableStateTest(t *testing.T) {
	astAnalyzerTest_runCases(t, astAnalyzerTest_MayEffectMutableStateTest_cases,
		(*AstAnalyzer).MayEffectMutableState)
}

// Java: AstAnalyzerTest.MayHaveSideEffects.test()
func TestAstAnalyzer_MayHaveSideEffects(t *testing.T) {
	astAnalyzerTest_runCases(t, astAnalyzerTest_MayHaveSideEffects_cases,
		(*AstAnalyzer).MayHaveSideEffects)
}

// Java: AstAnalyzerTest.NodeTypeMayHaveSideEffects.test()
func TestAstAnalyzer_NodeTypeMayHaveSideEffects(t *testing.T) {
	astAnalyzerTest_runCases(t, astAnalyzerTest_NodeTypeMayHaveSideEffects_cases,
		(*AstAnalyzer).NodeTypeMayHaveSideEffects)
}

func astAnalyzerTest_assertBool(t *testing.T, what string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %t, want %t", what, got, want)
	}
}

// Java: AstAnalyzerTest.FunctionCallHasSideEffects.testCallSideEffects()
func TestAstAnalyzer_FunctionCallHasSideEffects_testCallSideEffects(t *testing.T) {
	helper := newAstAnalyzerTest_ParseHelper(t)

	// Parens force interpretation as an expression.
	newXDotMethodCall := helper.parseFirst(rhino.Token_CALL, "(new x().method());")
	astAnalyzer := helper.getAstAnalyzer()
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(newXDotMethodCall), true)

	newExpr := newXDotMethodCall.GetFirstFirstChild()
	jutil.CheckState(newExpr.IsNew())
	flags := rhino.NewNode_SideEffectFlags()

	// No side effects, local result
	flags.ClearAllFlags()
	newExpr.SetSideEffectFlagsObject(flags)
	flags.ClearAllFlags()
	newXDotMethodCall.SetSideEffectFlagsObject(flags)

	// Cannot determine this evaluates to a local value (even though it does in practice).
	astAnalyzerTest_assertBool(t, "evaluatesToLocalValue", NodeUtil_EvaluatesToLocalValue(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "mayHaveSideEffects", astAnalyzer.MayHaveSideEffects(newXDotMethodCall), false)

	// Modifies this, local result
	flags.ClearAllFlags()
	newExpr.SetSideEffectFlagsObject(flags)
	flags.ClearAllFlags()
	flags.SetMutatesThis()
	newXDotMethodCall.SetSideEffectFlagsObject(flags)

	astAnalyzerTest_assertBool(t, "evaluatesToLocalValue", NodeUtil_EvaluatesToLocalValue(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "mayHaveSideEffects", astAnalyzer.MayHaveSideEffects(newXDotMethodCall), false)

	// Modifies this, non-local result
	flags.ClearAllFlags()
	newExpr.SetSideEffectFlagsObject(flags)
	flags.ClearAllFlags()
	flags.SetMutatesThis()
	newXDotMethodCall.SetSideEffectFlagsObject(flags)

	astAnalyzerTest_assertBool(t, "evaluatesToLocalValue", NodeUtil_EvaluatesToLocalValue(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "mayHaveSideEffects", astAnalyzer.MayHaveSideEffects(newXDotMethodCall), false)

	// No modifications, non-local result
	flags.ClearAllFlags()
	newExpr.SetSideEffectFlagsObject(flags)
	flags.ClearAllFlags()
	newXDotMethodCall.SetSideEffectFlagsObject(flags)

	astAnalyzerTest_assertBool(t, "evaluatesToLocalValue", NodeUtil_EvaluatesToLocalValue(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "mayHaveSideEffects", astAnalyzer.MayHaveSideEffects(newXDotMethodCall), false)

	// The new modifies global state, no side-effect call
	// This call could be removed, but not the new.
	flags.ClearAllFlags()
	flags.SetMutatesGlobalState()
	newExpr.SetSideEffectFlagsObject(flags)
	flags.ClearAllFlags()
	newXDotMethodCall.SetSideEffectFlagsObject(flags)

	// This does evaluate to a local value but NodeUtil does not know that
	astAnalyzerTest_assertBool(t, "evaluatesToLocalValue", NodeUtil_EvaluatesToLocalValue(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(newXDotMethodCall), false)
	astAnalyzerTest_assertBool(t, "mayHaveSideEffects", astAnalyzer.MayHaveSideEffects(newXDotMethodCall), true)
}

// Java: AstAnalyzerTest.FunctionCallHasSideEffects.testStringMethodCallSideEffects_noTypesForLocalOptimizations()
func TestAstAnalyzer_FunctionCallHasSideEffects_testStringMethodCallSideEffects_noTypesForLocalOptimizations(t *testing.T) {
	helper := newAstAnalyzerTest_ParseHelper(t)

	xDotReplaceCall := helper.parseFirst(rhino.Token_CALL, "x.replace(/xyz/g, '');")
	astAnalyzer := helper.getAstAnalyzer()
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(xDotReplaceCall), true)

	helper.hasGlobalRegexpReferences = false
	astAnalyzer = helper.getAstAnalyzer()
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(xDotReplaceCall), true)

	xNode := xDotReplaceCall.GetFirstFirstChild()
	xNode.SetJSType(helper.typeRegistry.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", astAnalyzer.FunctionCallHasSideEffects(xDotReplaceCall), true)
}

// Java: AstAnalyzerTest.FunctionCallHasSideEffects.testTypeBasedStringMethodCallSideEffects_useTypesForLocalOptimziations()
func TestAstAnalyzer_FunctionCallHasSideEffects_testTypeBasedStringMethodCallSideEffects_useTypesForLocalOptimziations(t *testing.T) {
	helper := newAstAnalyzerTest_ParseHelper(t)
	helper.useTypesForLocalOptimizations = true
	helper.hasGlobalRegexpReferences = false

	xDotReplaceCall := helper.parseFirst(rhino.Token_CALL, "x.replace(/xyz/g, '');")
	xDotReplaceCallStringType := xDotReplaceCall.CloneTree()
	xDotReplaceCallStringColor := xDotReplaceCall.CloneTree()

	stringType := helper.typeRegistry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	xDotReplaceCallStringType.GetFirstFirstChild().SetJSType(stringType)
	xDotReplaceCallStringColor.GetFirstFirstChild().SetColor(colors.StandardColors_STRING)
	astAnalyzer := helper.getAstAnalyzer()

	astAnalyzerTest_assertBool(t, "untyped", astAnalyzer.FunctionCallHasSideEffects(xDotReplaceCall), true)
	astAnalyzerTest_assertBool(t, "string type", astAnalyzer.FunctionCallHasSideEffects(xDotReplaceCallStringType), false)
	astAnalyzerTest_assertBool(t, "string color", astAnalyzer.FunctionCallHasSideEffects(xDotReplaceCallStringColor), false)
}

// Java: AstAnalyzerTest.BuiltInFunctionWithoutSideEffects.cases()
var astAnalyzerTest_BuiltInFunctionWithoutSideEffects_cases = []string{
	"Object", "Array", "String", "Number", "BigInt", "Boolean", "RegExp", "Error", "Symbol",
}

// Java: AstAnalyzerTest.BuiltInFunctionWithoutSideEffects.test_assumingPureBuiltins()
func TestAstAnalyzer_BuiltInFunctionWithoutSideEffects_test_assumingPureBuiltins(t *testing.T) {
	for i, functionName := range astAnalyzerTest_BuiltInFunctionWithoutSideEffects_cases {
		t.Run(fmt.Sprintf("#%d %s", i, functionName), func(t *testing.T) {
			parseHelper := newAstAnalyzerTest_ParseHelper(t)
			parseHelper.assumeBuiltinsPure = true
			func_ := parseHelper.parseFirst(rhino.Token_CALL, fmt.Sprintf("%s(1);", functionName))
			astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", parseHelper.getAstAnalyzer().FunctionCallHasSideEffects(func_), false)
		})
	}
}

// Java: AstAnalyzerTest.BuiltInFunctionWithoutSideEffects.test_noAssumePureBuiltins()
func TestAstAnalyzer_BuiltInFunctionWithoutSideEffects_test_noAssumePureBuiltins(t *testing.T) {
	for i, functionName := range astAnalyzerTest_BuiltInFunctionWithoutSideEffects_cases {
		t.Run(fmt.Sprintf("#%d %s", i, functionName), func(t *testing.T) {
			parseHelper := newAstAnalyzerTest_ParseHelper(t)
			parseHelper.assumeBuiltinsPure = false
			func_ := parseHelper.parseFirst(rhino.Token_CALL, fmt.Sprintf("%s(1);", functionName))
			astAnalyzerTest_assertBool(t, "functionCallHasSideEffects", parseHelper.getAstAnalyzer().FunctionCallHasSideEffects(func_), true)
		})
	}
}

// Java: AstAnalyzerTest.ConstructorCallHasSideEffects.byDefaultAConstructorCallHasSideEffects()
func TestAstAnalyzer_ConstructorCallHasSideEffects_byDefaultAConstructorCallHasSideEffects(t *testing.T) {
	parseHelper := newAstAnalyzerTest_ParseHelper(t)
	newNode := parseHelper.parseFirst(rhino.Token_NEW, "new SomeClass();")
	astAnalyzer := parseHelper.getAstAnalyzer()
	// we know nothing about the class being instantiated, so assume side effects occur.
	astAnalyzerTest_assertBool(t, "constructorCallHasSideEffects", astAnalyzer.ConstructorCallHasSideEffects(newNode), true)
}

// Java: AstAnalyzerTest.ConstructorCallHasSideEffects.constructorCallMarkedAsNoSideEffectsHasNone()
func TestAstAnalyzer_ConstructorCallHasSideEffects_constructorCallMarkedAsNoSideEffectsHasNone(t *testing.T) {
	parseHelper := newAstAnalyzerTest_ParseHelper(t)
	newNode := parseHelper.parseFirst(rhino.Token_NEW, "new SomeClass();")
	astAnalyzer := parseHelper.getAstAnalyzer()
	// simulate PureFunctionIdentifier marking the call as having no side effects.
	flags := rhino.NewNode_SideEffectFlags()
	flags.ClearAllFlags()
	newNode.SetSideEffectFlagsObject(flags)

	astAnalyzerTest_assertBool(t, "constructorCallHasSideEffects", astAnalyzer.ConstructorCallHasSideEffects(newNode), false)
}

// Java: AstAnalyzerTest.ConstructorCallHasSideEffects.modifyingALocalArgumentIsNotASideEffect()
func TestAstAnalyzer_ConstructorCallHasSideEffects_modifyingALocalArgumentIsNotASideEffect(t *testing.T) {
	parseHelper := newAstAnalyzerTest_ParseHelper(t)
	// object literal is considered a local value. Modifying it isn't a side effect, since
	// nothing else looks at it.
	newNode := parseHelper.parseFirst(rhino.Token_NEW, "new SomeClass({});")
	astAnalyzer := parseHelper.getAstAnalyzer()
	// simulate PureFunctionIdentifier marking the call as only modifying its arguments
	flags := rhino.NewNode_SideEffectFlags()
	flags.ClearAllFlags()
	flags.SetMutatesArguments()
	newNode.SetSideEffectFlagsObject(flags)

	astAnalyzerTest_assertBool(t, "constructorCallHasSideEffects", astAnalyzer.ConstructorCallHasSideEffects(newNode), false)
}

// Java: AstAnalyzerTest.ConstructorCallHasSideEffects.modifyingANonLocalArgumentIsASideEffect()
func TestAstAnalyzer_ConstructorCallHasSideEffects_modifyingANonLocalArgumentIsASideEffect(t *testing.T) {
	parseHelper := newAstAnalyzerTest_ParseHelper(t)
	// variable name is a non-local value. Modifying it is a side effect.
	newNode := parseHelper.parseFirst(rhino.Token_NEW, "new SomeClass(x);")
	astAnalyzer := parseHelper.getAstAnalyzer()
	// simulate PureFunctionIdentifier marking the call as only modifying its arguments
	flags := rhino.NewNode_SideEffectFlags()
	flags.ClearAllFlags()
	flags.SetMutatesArguments()
	newNode.SetSideEffectFlagsObject(flags)

	astAnalyzerTest_assertBool(t, "constructorCallHasSideEffects", astAnalyzer.ConstructorCallHasSideEffects(newNode), true)
}

// Java: AstAnalyzerTest.ConstructorsKnownToHaveNoSideEffects.cases()
var astAnalyzerTest_ConstructorsKnownToHaveNoSideEffects_cases = []string{
	"Array", "Date", "Error", "Object", "RegExp", "XMLHttpRequest",
}

// Java: AstAnalyzerTest.ConstructorsKnownToHaveNoSideEffects.noSideEffectsForKnownConstructor()
func TestAstAnalyzer_ConstructorsKnownToHaveNoSideEffects_noSideEffectsForKnownConstructor(t *testing.T) {
	for i, constructorName := range astAnalyzerTest_ConstructorsKnownToHaveNoSideEffects_cases {
		t.Run(fmt.Sprintf("#%d %s", i, constructorName), func(t *testing.T) {
			parseHelper := newAstAnalyzerTest_ParseHelper(t)
			newNode := parseHelper.parseFirst(rhino.Token_NEW, fmt.Sprintf("new %s();", constructorName))
			astAnalyzer := parseHelper.getAstAnalyzer()
			// we know nothing about the class being instantiated, so assume side effects occur.
			astAnalyzerTest_assertBool(t, "constructorCallHasSideEffects", astAnalyzer.ConstructorCallHasSideEffects(newNode), false)
		})
	}
}

// Go-only: test-helper — the AstAnalyzer.Options builder generated by
// AutoBuilder (AutoBuilder_AstAnalyzer_Options_Builder) has no Java test;
// these check the generated code's behaviour.
func TestAstAnalyzer_OptionsBuilder(t *testing.T) {
	o := AstAnalyzer_Options_Builder_().
		SetUseTypesForLocalOptimization(true).
		SetAssumeGettersArePure(false).
		SetHasRegexpGlobalReferences(true).
		SetAssumeKnownBuiltinsArePure(false).
		Build()
	if !o.UseTypesForLocalOptimization() || o.AssumeGettersArePure() || !o.HasRegexpGlobalReferences() || o.AssumeKnownBuiltinsArePure() {
		t.Errorf("unexpected options: %+v", o)
	}
	jsaThrows[*jutil.IllegalArgumentException](t, func() {
		NewAstAnalyzer(o, nil, nil)
	})
}

// Java: AstAnalyzerTest.MayEffectMutableStateTest.cases()
var astAnalyzerTest_MayEffectMutableStateTest_cases = []*astAnalyzerTest_AnalysisCase{
	kase().Js("i++").Token(rhino.Token_INC).Expect(true),
	kase().Js("[b, [a, i++]]").Token(rhino.Token_ARRAYLIT).Expect(true),
	kase().Js("i=3").Token(rhino.Token_ASSIGN).Expect(true),
	kase().Js("[0, i=3]").Token(rhino.Token_ARRAYLIT).Expect(true),
	kase().Js("b()").Token(rhino.Token_CALL).Expect(true),
	kase().Js("b?.()").Token(rhino.Token_OPTCHAIN_CALL).Expect(true),
	kase().Js("void b()").Token(rhino.Token_VOID).Expect(true),
	kase().Js("[1, b()]").Token(rhino.Token_ARRAYLIT).Expect(true),
	kase().Js("b.b=4").Token(rhino.Token_ASSIGN).Expect(true),
	kase().Js("b.b--").Token(rhino.Token_DEC).Expect(true),
	kase().Js("i--").Token(rhino.Token_DEC).Expect(true),
	kase().Js("a[0][i=4]").Token(rhino.Token_GETELEM).Expect(true),
	kase().Js("a?.[0][i=4]").Token(rhino.Token_OPTCHAIN_GETELEM).Expect(true),
	kase().Js("a += 3").Token(rhino.Token_ASSIGN_ADD).Expect(true),
	kase().Js("a ||= b").Token(rhino.Token_ASSIGN_OR).Expect(true),
	kase().Js("a &&= b").Token(rhino.Token_ASSIGN_AND).Expect(true),
	kase().Js("a ??= b").Token(rhino.Token_ASSIGN_COALESCE).Expect(true),
	kase().Js("a, b, z += 4").Token(rhino.Token_COMMA).Expect(true),
	kase().Js("a ? c : d++").Token(rhino.Token_HOOK).Expect(true),
	kase().Js("a ?? b++").Token(rhino.Token_COALESCE).Expect(true),
	kase().Js("a + c++").Token(rhino.Token_ADD).Expect(true),
	kase().Js("a + c - d()").Token(rhino.Token_SUB).Expect(true),
	kase().Js("function foo() {}").Token(rhino.Token_FUNCTION).Expect(true),
	kase().Js("while(true);").Token(rhino.Token_WHILE).Expect(true),
	kase().Js("if(true){a()}").Token(rhino.Token_IF).Expect(true),
	kase().Js("if(true){a}").Token(rhino.Token_IF).Expect(false),
	kase().Js("(function() { })").Token(rhino.Token_FUNCTION).Expect(true),
	kase().Js("(function() { i++ })").Token(rhino.Token_FUNCTION).Expect(true),
	kase().Js("[function a(){}]").Token(rhino.Token_ARRAYLIT).Expect(true),
	kase().Js("a").Token(rhino.Token_NAME).Expect(false),
	kase().Js("[b, c [d, [e]]]").Token(rhino.Token_ARRAYLIT).Expect(true),
	kase().Js("({a: x, b: y, c: z})").Token(rhino.Token_OBJECTLIT).Expect(true),
	// Note: RegExp objects are not immutable, for instance, the exec
	// method maintains state for "global" searches.
	kase().Js("/abc/gi").Token(rhino.Token_REGEXP).Expect(true),
	kase().Js("'a'").Token(rhino.Token_STRINGLIT).Expect(false),
	kase().Js("0").Token(rhino.Token_NUMBER).Expect(false),
	kase().Js("1n").Token(rhino.Token_BIGINT).Expect(false),
	kase().Js("a + c").Token(rhino.Token_ADD).Expect(false),
	kase().Js("'c' + a[0]").Token(rhino.Token_ADD).Expect(false),
	kase().Js("a[0][1]").Token(rhino.Token_GETELEM).Expect(false),
	kase().Js("a?.[0][1]").Token(rhino.Token_OPTCHAIN_GETELEM).Expect(false),
	kase().Js("'a' + c").Token(rhino.Token_ADD).Expect(false),
	kase().Js("'a' + a.name").Token(rhino.Token_ADD).Expect(false),
	kase().Js("1, 2, 3").Token(rhino.Token_COMMA).Expect(false),
	kase().Js("a, b, 3").Token(rhino.Token_COMMA).Expect(false),
	kase().Js("(function(a, b) {  })").Token(rhino.Token_FUNCTION).Expect(true),
	kase().Js("a ? c : d").Token(rhino.Token_HOOK).Expect(false),
	kase().Js("a ?? b").Token(rhino.Token_COALESCE).Expect(false),
	kase().Js("'1' + navigator.userAgent").Token(rhino.Token_ADD).Expect(false),
	kase().Js("new RegExp('foobar', 'i')").Token(rhino.Token_NEW).Expect(true),
	kase().Js("new RegExp(SomethingWacky(), 'i')").Token(rhino.Token_NEW).Expect(true),
	kase().Js("new Array()").Token(rhino.Token_NEW).Expect(true),
	kase().Js("new Array").Token(rhino.Token_NEW).Expect(true),
	kase().Js("new Array(4)").Token(rhino.Token_NEW).Expect(true),
	kase().Js("new Array('a', 'b', 'c')").Token(rhino.Token_NEW).Expect(true),
	kase().Js("new SomeClassINeverHeardOf()").Token(rhino.Token_NEW).Expect(true),

	// Getters and setters - object rest and spread
	kase().Js("({...x});").Token(rhino.Token_OBJECT_SPREAD).AssumeGettersArePure(false).Expect(true),
	kase().Js("const {...x} = y;").Token(rhino.Token_OBJECT_REST).AssumeGettersArePure(false).Expect(true),
	kase().Js("({...x});").Token(rhino.Token_OBJECT_SPREAD).AssumeGettersArePure(true).Expect(false),
	kase().Js("const {...x} = y;").Token(rhino.Token_OBJECT_REST).AssumeGettersArePure(true).Expect(false),
	kase().Js("({...f().x});").Token(rhino.Token_OBJECT_SPREAD).AssumeGettersArePure(true).Expect(true),
	kase().Js("({...f().x} = y);").Token(rhino.Token_OBJECT_REST).AssumeGettersArePure(true).Expect(true),

	// the presence of `a` affects what gets put into `x`
	kase().Js("({a, ...x} = y);").Token(rhino.Token_STRING_KEY).AssumeGettersArePure(true).Expect(true),

	// Getters and setters
	kase().Js("x.getter;").Token(rhino.Token_GETPROP).Expect(true),
	kase().Js("x?.getter;").Token(rhino.Token_OPTCHAIN_GETPROP).Expect(true),
	// Overapproximation to avoid inspecting the parent.
	kase().Js("x.setter;").Token(rhino.Token_GETPROP).Expect(true),
	kase().Js("x?.setter;").Token(rhino.Token_OPTCHAIN_GETPROP).Expect(true),
	kase().Js("x.normal;").Token(rhino.Token_GETPROP).Expect(false),
	kase().Js("x?.normal;").Token(rhino.Token_OPTCHAIN_GETPROP).Expect(false),
	kase().Js("const {getter} = x;").Token(rhino.Token_STRING_KEY).Expect(true),
	// Overapproximation to avoid inspecting the parent.
	kase().Js("const {setter} = x;").Token(rhino.Token_STRING_KEY).Expect(false),
	kase().Js("const {normal} = x;").Token(rhino.Token_STRING_KEY).Expect(false),
	kase().Js("x.getter = 0;").Token(rhino.Token_GETPROP).Expect(true),
	kase().Js("x.setter = 0;").Token(rhino.Token_GETPROP).Expect(true),
	kase().Js("x.normal = 0;").Token(rhino.Token_GETPROP).Expect(false),

	// Default values delegates to children.
	kase().Js("({x = 0} = y);").Token(rhino.Token_DEFAULT_VALUE).Expect(false),
	kase().Js("([x = 0] = y);").Token(rhino.Token_DEFAULT_VALUE).Expect(false),
	kase().Js("function f(x = 0) { };").Token(rhino.Token_DEFAULT_VALUE).Expect(false),

	// Dynamic import can mutate global state
	kase().Js("import('./module.js')").Token(rhino.Token_DYNAMIC_IMPORT).Expect(true),
}

// Java: AstAnalyzerTest.MayHaveSideEffects.cases()
var astAnalyzerTest_MayHaveSideEffects_cases = []*astAnalyzerTest_AnalysisCase{
	// Cases in need of differentiation.
	kase().Expect(false).Js("[1]"),
	kase().Expect(false).Js("[1, 2]"),
	kase().Expect(false).Js("[1, 2]").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("i++"),
	kase().Expect(true).Js("[b, [a, i++]]"),
	kase().Expect(true).Js("i=3"),
	kase().Expect(true).Js("[0, i=3]"),
	kase().Expect(true).Js("b()"),
	kase().Expect(true).Js("[1, b()]"),
	kase().Expect(true).Js("b.b=4"),
	kase().Expect(true).Js("b.b--"),
	kase().Expect(true).Js("i--"),
	kase().Expect(true).Js("a[0][i=4]"),
	kase().Expect(true).Js("a += 3"),
	kase().Expect(true).Js("a, b, z += 4"),
	kase().Expect(true).Js("a ? c : d++"),
	kase().Expect(true).Js("a ?? b++"),
	kase().Expect(true).Js("a + c++"),
	kase().Expect(true).Js("a + c - d()"),
	kase().Expect(true).Js("a + c - d()"),
	kase().Expect(true).Js("function foo() {}"),
	kase().Expect(true).Js("class Foo {}"),
	kase().Expect(true).Js("while(true);"),
	kase().Expect(true).Js("if(true){a()}"),
	kase().Expect(false).Js("if(true){a}"),
	kase().Expect(false).Js("(function() { })"),
	kase().Expect(false).Js("(function() { i++ })"),
	kase().Expect(false).Js("[function a(){}]"),
	kase().Expect(false).Js("(class { })"),
	kase().Expect(false).Js("(class { method() { i++ } })"),
	kase().Expect(true).Js("(class { [computedName()]() {} })"),
	kase().Expect(false).Js("(class { [computedName]() {} })"),
	kase().Expect(false).Js("(class Foo extends Bar { })"),
	kase().Expect(true).Js("(class extends foo() { })"),
	kase().Expect(false).Js("a"),
	kase().Expect(false).Js("a").AssumeBuiltinsPure(false),
	kase().Expect(false).Js("a.b"),
	kase().Expect(false).Js("a.b.c"),
	kase().Expect(false).Js("[b, c [d, [e]]]"),
	kase().Expect(false).Js("({a: x, b: y, c: z})"),
	kase().Expect(false).Js("({a, b, c})"),
	kase().Expect(false).Js("/abc/gi"),
	kase().Expect(false).Js("'a'"),
	kase().Expect(false).Js("0"),
	kase().Expect(false).Js("a + c"),
	kase().Expect(false).Js("'c' + a[0]"),
	kase().Expect(false).Js("a[0][1]"),
	kase().Expect(false).Js("'a' + c"),
	kase().Expect(false).Js("'a' + a.name"),
	kase().Expect(false).Js("1, 2, 3"),
	kase().Expect(false).Js("a, b, 3"),
	kase().Expect(false).Js("(function(a, b) {  })"),
	kase().Expect(false).Js("a ? c : d"),
	kase().Expect(false).Js("a ?? b"),
	kase().Expect(false).Js("'1' + navigator.userAgent"),
	kase().Expect(false).Js("`template`"),
	kase().Expect(false).Js("`template${name}`"),
	kase().Expect(false).Js("`${name}template`"),
	kase().Expect(true).Js("`${naming()}template`"),
	kase().Expect(true).Js("templateFunction`template`"),
	kase().Expect(true).Js("st = `${name}template`"),
	kase().Expect(true).Js("tempFunc = templateFunction`template`"),
	kase().Expect(false).Js("new RegExp('foobar', 'i')"),
	kase().Expect(true).Js("new RegExp(SomethingWacky(), 'i')"),
	kase().Expect(false).Js("new Array()"),
	kase().Expect(false).Js("new Array"),
	kase().Expect(false).Js("new Array(4)"),
	kase().Expect(false).Js("new Array('a', 'b', 'c')"),
	kase().Expect(true).Js("new Array()").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("new SomeClassINeverHeardOf()"),
	kase().Expect(true).Js("new SomeClassINeverHeardOf()"),
	kase().Expect(false).Js("({}).foo = 4"),
	kase().Expect(false).Js("([]).foo = 4"),
	kase().Expect(false).Js("(function() {}).foo = 4"),
	kase().Expect(true).Js("this.foo = 4"),
	kase().Expect(true).Js("a.foo = 4"),
	kase().Expect(true).Js("(function() { return n; })().foo = 4"),
	kase().Expect(true).Js("([]).foo = bar()"),
	kase().Expect(false).Js("undefined"),
	kase().Expect(false).Js("void 0"),
	kase().Expect(true).Js("void foo()"),
	kase().Expect(false).Js("-Infinity"),
	kase().Expect(false).Js("Infinity"),
	kase().Expect(false).Js("NaN"),
	kase().Expect(false).Js("({}||[]).foo = 2;"),
	kase().Expect(false).Js("(true ? {} : []).foo = 2;"),
	kase().Expect(false).Js("({},[]).foo = 2;"),
	kase().Expect(true).Js("delete a.b"),
	kase().Expect(false).Js("Math.random();"),
	kase().Expect(true).Js("Math.random();").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("Math.random(seed);"),
	kase().Expect(true).Js("Math.random(seed);").AssumeBuiltinsPure(false),
	kase().Expect(false).Js("Math.sin(1, 10);"),
	kase().Expect(true).Js("Math.sin(1, 10);").AssumeBuiltinsPure(false),
	kase().Expect(false).Js("[1, 1].foo;"),
	kase().Expect(true).Js("export var x = 0;"),
	kase().Expect(true).Js("export let x = 0;"),
	kase().Expect(true).Js("export const x = 0;"),
	kase().Expect(true).Js("export class X {};"),
	kase().Expect(true).Js("export function x() {};"),
	kase().Expect(true).Js("export {x};"),

	// ARRAYLIT-ITER_SPREAD
	kase().Expect(false).Js("[...[]]"),
	kase().Expect(false).Js("[...[1]]"),
	kase().Expect(true).Js("[...[i++]]"),
	kase().Expect(false).Js("[...'string']"),
	kase().Expect(false).Js("[...`templatelit`]"),
	kase().Expect(false).Js("[...`templatelit ${safe}`]"),
	kase().Expect(true).Js("[...`templatelit ${unsafe()}`]"),
	kase().Expect(true).Js("[...f()]"),
	kase().Expect(true).Js("[...5]"),
	kase().Expect(true).Js("[...null]"),
	kase().Expect(true).Js("[...true]"),

	// CALL-ITER_SPREAD
	kase().Expect(false).Js("Math.sin(...[])"),
	kase().Expect(false).Js("Math.sin(...[1])"),
	kase().Expect(true).Js("Math.sin(...[i++])"),
	kase().Expect(false).Js("Math.sin(...'string')"),
	kase().Expect(false).Js("Math.sin(...`templatelit`)"),
	kase().Expect(false).Js("Math.sin(...`templatelit ${safe}`)"),
	kase().Expect(true).Js("Math.sin(...`templatelit ${unsafe()}`)"),
	kase().Expect(true).Js("Math.sin(...f())"),
	kase().Expect(true).Js("Math.sin(...5)"),
	kase().Expect(true).Js("Math.sin(...null)"),
	kase().Expect(true).Js("Math.sin(...true)"),

	// NEW-ITER_SPREAD
	kase().Expect(false).Js("new Object(...[])"),
	kase().Expect(false).Js("new Object(...[1])"),
	kase().Expect(true).Js("new Object(...[i++])"),
	kase().Expect(false).Js("new Object(...'string')"),
	kase().Expect(false).Js("new Object(...`templatelit`)"),
	kase().Expect(false).Js("new Object(...`templatelit ${safe}`)"),
	kase().Expect(true).Js("new Object(...`templatelit ${unsafe()}`)"),
	kase().Expect(true).Js("new Object(...f())"),
	kase().Expect(true).Js("new Object(...5)"),
	kase().Expect(true).Js("new Object(...null)"),
	kase().Expect(true).Js("new Object(...true)"),

	// OBJECT_SPREAD
	// These could all invoke getters.
	kase().Expect(true).Js("({...x})"),
	kase().Expect(true).Js("({...{}})"),
	kase().Expect(true).Js("({...{a:1}})"),
	kase().Expect(true).Js("({...{a:i++}})"),
	kase().Expect(true).Js("({...{a:f()}})"),
	kase().Expect(true).Js("({...f()})"),

	// OBJECT_REST
	// This could invoke getters.
	kase().Expect(true).Token(rhino.Token_OBJECT_REST).Js("({...x} = something)"),
	// the presence of `a` affects what goes into `x`
	kase().Expect(true).Token(rhino.Token_STRING_KEY).Js("({a, ...x} = something)"),

	// ITER_REST
	// We currently assume all iterable-rests are side-effectful.
	kase().Expect(true).Token(rhino.Token_ITER_REST).Js("([...x] = 'safe')"),
	kase().Expect(false).Token(rhino.Token_ITER_REST).Js("(function(...x) { })"),

	// Context switch
	kase().Expect(true).Token(rhino.Token_AWAIT).Js("async function f() { await 0; }"),
	kase().Expect(true).Token(rhino.Token_FOR_AWAIT_OF).Js("(async()=>{ for await (let x of []) {} })"),
	kase().Expect(true).Token(rhino.Token_THROW).Js("function f() { throw 'something'; }"),
	kase().Expect(true).Token(rhino.Token_YIELD).Js("function* f() { yield 'something'; }"),
	kase().Expect(true).Token(rhino.Token_YIELD).Js("function* f() { yield* 'something'; }"),

	// Enhanced for loop
	// These edge cases are actually side-effect free. We include them to confirm we just give
	// up
	// on enhanced for loops.
	kase().Expect(true).Js("for (const x in []) { }"),
	kase().Expect(true).Js("for (const x of []) { }"),

	// COMPUTED_PROP - OBJECTLIT
	kase().Expect(false).Token(rhino.Token_COMPUTED_PROP).Js("({[a]: x})"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_PROP).Js("({[a()]: x})"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_PROP).Js("({[a]: x()})"),

	// computed property getters and setters are modeled as COMPUTED_PROP with an
	// annotation to indicate getter or setter.
	kase().Expect(false).Token(rhino.Token_COMPUTED_PROP).Js("({ get [a]() {} })"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_PROP).Js("({ get [a()]() {} })"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_PROP).Js("({ set [a](x) {} })"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_PROP).Js("({ set [a()](x) {} })"),

	// COMPUTED_PROP - CLASS
	kase().Expect(false).Token(rhino.Token_COMPUTED_PROP).Js("class C { [a]() {} }"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_PROP).Js("class C { [a()]() {} }"),

	// computed property getters and setters are modeled as COMPUTED_PROP with an
	// annotation to indicate getter or setter.
	kase().Expect(false).Token(rhino.Token_COMPUTED_PROP).Js("class C { get [a]() {} }"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_PROP).Js("class C { get [a()]() {} }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_PROP).Js("class C { set [a](x) {} }"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_PROP).Js("class C { set [a()](x) {} }"),

	// GETTER_DEF
	kase().Expect(false).Token(rhino.Token_GETTER_DEF).Js("({ get a() {} })"),
	kase().Expect(false).Token(rhino.Token_GETTER_DEF).Js("class C { get a() {} }"),

	// Getter use
	kase().Expect(true).Token(rhino.Token_GETPROP).Js("x.getter;"),
	kase().Expect(true).Token(rhino.Token_OPTCHAIN_GETPROP).Js("x?.getter;"),
	// Overapproximation because to avoid inspecting the parent.
	kase().Expect(true).Token(rhino.Token_GETPROP).Js("x.setter;"),
	kase().Expect(true).Token(rhino.Token_OPTCHAIN_GETPROP).Js("x?.setter;"),
	kase().Expect(false).Token(rhino.Token_GETPROP).Js("x.normal;"),
	kase().Expect(false).Token(rhino.Token_OPTCHAIN_GETPROP).Js("x?.normal;"),
	kase().Expect(true).Token(rhino.Token_STRING_KEY).Js("({getter} = foo());"),
	kase().Expect(false).Token(rhino.Token_STRING_KEY).Js("({setter} = foo());"),
	kase().Expect(false).Token(rhino.Token_STRING_KEY).Js("({normal} = foo());"),

	// SETTER_DEF
	kase().Expect(false).Token(rhino.Token_SETTER_DEF).Js("({ set a(x) {} })"),
	kase().Expect(false).Token(rhino.Token_SETTER_DEF).Js("class C { set a(x) {} }"),

	// SETTER_USE
	// Overapproximation because to avoid inspecting the parent.
	kase().Expect(true).Token(rhino.Token_GETPROP).Js("x.getter = 0;"),
	kase().Expect(true).Token(rhino.Token_GETPROP).Js("x.setter = 0;"),
	kase().Expect(false).Token(rhino.Token_GETPROP).Js("x.normal = 0;"),

	// MEMBER_FUNCTION_DEF
	kase().Expect(false).Token(rhino.Token_MEMBER_FUNCTION_DEF).Js("({ a(x) {} })"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FUNCTION_DEF).Js("class C { a(x) {} }"),

	// MEMBER_FIELD_DEF
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { x=2; }"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { x; }"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { x }"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { x \n y }"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { static x=2; }"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { static x; }"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { static x }"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { static x \n static y }"),
	kase().Expect(false).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { x = alert(1); }"),
	kase().Expect(true).Token(rhino.Token_MEMBER_FIELD_DEF).Js("class C { static x = alert(1); }"),

	// COMPUTED_FIELD_DEF
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { [x]; }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { ['x']=2; }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { 'x'=2; }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { 1=2; }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { static [x]; }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { static ['x']=2; }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { static 'x'=2; }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { static 1=2; }"),
	kase().Expect(false).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { ['x'] = alert(1); }"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { static ['x'] = alert(1); }"),
	kase().Expect(true).Token(rhino.Token_COMPUTED_FIELD_DEF).Js("class C { static [alert(1)] = 2; }"),

	// CLASS_STATIC_BLOCK
	kase().Expect(false).Token(rhino.Token_BLOCK).Js("class C { static {} }"),
	kase().Expect(false).Token(rhino.Token_BLOCK).Js("class C { static { [1]; } }"),
	kase().Expect(true).Token(rhino.Token_BLOCK).Js("class C { static { let x; } }"),
	kase().Expect(true).Token(rhino.Token_BLOCK).Js("class C { static { const x =1 ; } }"),
	kase().Expect(true).Token(rhino.Token_BLOCK).Js("class C { static { var x; } }"),
	kase().Expect(true).Token(rhino.Token_BLOCK).Js("class C { static { this.x = 1; } }"),
	kase().Expect(true).Token(rhino.Token_BLOCK).Js("class C { static { function f() { } } }"),
	kase().Expect(false).Token(rhino.Token_BLOCK).Js("class C { static { (function () {} )} }"),
	kase().Expect(false).Token(rhino.Token_BLOCK).Js("class C { static { ()=>{} } }"),

	// SUPER calls
	kase().Expect(false).Token(rhino.Token_SUPER).Js("super()"),
	kase().Expect(false).Token(rhino.Token_SUPER).Js("super.foo()"),

	// Some common Object methods are assumed to be side-effect free when called
	// with their standard builtin arity.
	kase().Expect(false).Js("o.toString()"),
	kase().Expect(false).Js("o.toString(16)"),
	kase().Expect(false).Js("(0).toString(16)"),
	kase().Expect(false).Js("(0n).toString(16)"),
	kase().Expect(false).Js("o?.toString(16)"),
	kase().Expect(false).Js("(0)?.toString(16)"),
	kase().Expect(false).Js("(0n)?.toString(16)"),
	kase().Expect(false).Js("o.valueOf()"),
	kase().Expect(true).Js("o.toString(16, 2)"),
	kase().Expect(true).Js("o.valueOf(1)"),
	// These builtin-method assumptions are disabled when assumeKnownBuiltinsArePure is false.
	kase().Expect(true).Js("o.toString()").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("o.toString(16)").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("(0).toString(16)").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("(0n).toString(16)").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("o?.toString(16)").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("(0)?.toString(16)").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("(0n)?.toString(16)").AssumeBuiltinsPure(false),
	kase().Expect(true).Js("o.valueOf()").AssumeBuiltinsPure(false),
	// other methods depend on the extern definitions
	kase().Expect(true).Js("o.watch()"),

	// A RegExp Object by itself doesn't have any side-effects
	kase().Expect(false).Js("/abc/gi").GlobalRegExp(true),
	kase().Expect(false).Js("/abc/gi").GlobalRegExp(false),

	// RegExp instance methods have global side-effects, so whether they are
	// considered side-effect free depends on whether the global properties
	// are referenced and whether we assume that the built-in methods are not overwritten
	// with impure variants.
	kase().Expect(true).Js("(/abc/gi).test('')").GlobalRegExp(true),
	kase().Expect(false).Js("(/abc/gi).test('')").GlobalRegExp(false),
	kase().Expect(true).Js("(/abc/gi).test('')").GlobalRegExp(false).AssumeBuiltinsPure(false),
	kase().Expect(true).Js("(/abc/gi).test(a)").GlobalRegExp(true),
	kase().Expect(false).Js("(/abc/gi).test(b)").GlobalRegExp(false),
	kase().Expect(true).Js("(/abc/gi).exec('')").GlobalRegExp(true),
	kase().Expect(false).Js("(/abc/gi).exec('')").GlobalRegExp(false),

	// Some RegExp object method that may have side-effects.
	kase().Expect(true).Js("(/abc/gi).foo('')").GlobalRegExp(true),
	kase().Expect(true).Js("(/abc/gi).foo('')").GlobalRegExp(false),

	// Try the string RegExp ops.
	kase().Expect(true).Js("''.match('a')").GlobalRegExp(true),
	kase().Expect(false).Js("''.match('a')").GlobalRegExp(false),
	kase().Expect(true).Js("''.match(/(a)/)").GlobalRegExp(true),
	kase().Expect(false).Js("''.match(/(a)/)").GlobalRegExp(false),
	kase().Expect(true).Js("''.replace('a')").GlobalRegExp(true),
	kase().Expect(false).Js("''.replace('a')").GlobalRegExp(false),
	kase().Expect(true).Js("''.search('a')").GlobalRegExp(true),
	kase().Expect(false).Js("''.search('a')").GlobalRegExp(false),
	kase().Expect(true).Js("''.split('a')").GlobalRegExp(true),
	kase().Expect(false).Js("''.split('a')").GlobalRegExp(false),

	// Some non-RegExp string op that may have side-effects.
	kase().Expect(true).Js("''.foo('a')").GlobalRegExp(true),
	kase().Expect(true).Js("''.foo('a')").GlobalRegExp(false),

	// 'a' might be a RegExp object with the 'g' flag, in which case
	// the state might change by running any of the string ops.
	// Specifically, using these methods resets the "lastIndex" if used
	// in combination with a RegExp instance "exec" method.
	kase().Expect(true).Js("''.match(a)").GlobalRegExp(true),
	kase().Expect(true).Js("''.match(a)").GlobalRegExp(false),
	kase().Expect(true).Js("'a'.replace(/a/, function (s) {alert(s)})").GlobalRegExp(false),
	kase().Expect(false).Js("'a'.replace(/a/, 'x')").GlobalRegExp(false),

	// Dynamic import changes global state
	kase().Expect(true).Token(rhino.Token_DYNAMIC_IMPORT).Js("import('./module.js')"),

	// Symbol() call expressions are side-effect free
	kase().Expect(false).Js("Symbol()"),
	kase().Expect(false).Js("Symbol('desc')"),
	kase().Expect(true).Js("Symbol(foo())"),
	kase().Expect(true).Js("Symbol()").AssumeBuiltinsPure(false),
}

// Java: AstAnalyzerTest.NodeTypeMayHaveSideEffects.cases()
var astAnalyzerTest_NodeTypeMayHaveSideEffects_cases = []*astAnalyzerTest_AnalysisCase{
	kase().Js("x = y").Token(rhino.Token_ASSIGN).Expect(true),
	kase().Js("x += y").Token(rhino.Token_ASSIGN_ADD).Expect(true),
	kase().Js("delete x.y").Token(rhino.Token_DELPROP).Expect(true),
	kase().Js("x++").Token(rhino.Token_INC).Expect(true),
	kase().Js("x--").Token(rhino.Token_DEC).Expect(true),
	kase().Js("for (prop in obj) {}").Token(rhino.Token_FOR_IN).Expect(true),
	kase().Js("for (item of iterable) {}").Token(rhino.Token_FOR_OF).Expect(true),
	// name declaration has a side effect when a value is actually assigned
	kase().Js("var x = 1;").Token(rhino.Token_NAME).Expect(true),
	kase().Js("let x = 1;").Token(rhino.Token_NAME).Expect(true),
	kase().Js("const x = 1;").Token(rhino.Token_NAME).Expect(true),
	// don't consider name declaration to have a side effect if there's no assignment
	kase().Js("var x;").Token(rhino.Token_NAME).Expect(false),
	kase().Js("let x;").Token(rhino.Token_NAME).Expect(false),

	// destructuring declarations and assignments are always considered side effectful even
	// when empty
	kase().Js("var {x} = {};").Token(rhino.Token_DESTRUCTURING_LHS).Expect(true),
	kase().Js("var {} = {};").Token(rhino.Token_DESTRUCTURING_LHS).Expect(true),
	kase().Js("var [x] = [];").Token(rhino.Token_DESTRUCTURING_LHS).Expect(true),
	kase().Js("var {y: [x]} = {};").Token(rhino.Token_OBJECT_PATTERN).Expect(false),
	kase().Js("var {y: [x]} = {};").Token(rhino.Token_ARRAY_PATTERN).Expect(false),
	kase().Js("[x] = arr;").Token(rhino.Token_ASSIGN).Expect(true),

	// NOTE: CALL and NEW nodes are delegated to functionCallHasSideEffects() and
	// constructorCallHasSideEffects(), respectively. The cases below are just a few
	// representative examples that are convenient to test here.          // in general function and constructor calls are assumed to have side effects
	kase().Js("foo();").Token(rhino.Token_CALL).Expect(true),
	kase().Js("foo?.();").Token(rhino.Token_OPTCHAIN_CALL).Expect(true),
	kase().Js("new Foo();").Token(rhino.Token_NEW).Expect(true),
	// Object() is known not to have side-effects, though
	kase().Js("Object();").Token(rhino.Token_CALL).Expect(false),
	kase().Js("Object?.();").Token(rhino.Token_OPTCHAIN_CALL).Expect(false),
	kase().Js("new Object();").Token(rhino.Token_NEW).Expect(false),
	// Symbol() is known not to have side-effects
	kase().Js("Symbol();").Token(rhino.Token_CALL).Expect(false),
	kase().Js("Symbol('desc');").Token(rhino.Token_CALL).Expect(false),
	// TAGGED_TEMPLATELIT is just a special syntax for a CALL.
	kase().Js("foo`template`;").Token(rhino.Token_TAGGED_TEMPLATELIT).Expect(true),

	// NOTE: ITER_REST and ITER_SPREAD are delegated to NodeUtil.iteratesImpureIterable()
	// Test cases here are just easy to test representative examples.
	kase().Js("[...[1, 2, 3]]").Token(rhino.Token_ITER_SPREAD).Expect(false),
	// unknown iterable, so assume side-effects
	kase().Js("[...someIterable]").Token(rhino.Token_ITER_SPREAD).Expect(true),
	// we just assume the rhs of an array pattern may have iteration side-effects
	// without looking too closely.
	kase().Js("let [...rest] = [1, 2, 3];").Token(rhino.Token_ITER_REST).Expect(true),
	// ITER_REST in parameter list does not trigger iteration at the function definition, so
	// it has no side effects.
	kase().Js("function foo(...rest) {}").Token(rhino.Token_ITER_REST).Expect(false),

	// defining a class or a function is not considered to be a side-effect
	kase().Js("function foo() {}").Token(rhino.Token_FUNCTION).Expect(false),
	kase().Js("class Foo {}").Token(rhino.Token_CLASS).Expect(false),

	// arithmetic, logic, and bitwise operations do not have side-effects
	kase().Js("x + y").Token(rhino.Token_ADD).Expect(false),
	kase().Js("x || y").Token(rhino.Token_OR).Expect(false),
	kase().Js("x | y").Token(rhino.Token_BITOR).Expect(false),
	kase().Js("x ?? y").Token(rhino.Token_COALESCE).Expect(false),

	// Getters and setters
	kase().Js("({...x});").Token(rhino.Token_OBJECT_SPREAD).Expect(true),
	kase().Js("const {...x} = y;").Token(rhino.Token_OBJECT_REST).Expect(true),
	kase().Js("y.getter;").Token(rhino.Token_GETPROP).Expect(true),
	kase().Js("y?.getter;").Token(rhino.Token_OPTCHAIN_GETPROP).Expect(true),
	kase().Js("y.setter;").Token(rhino.Token_GETPROP).Expect(true),
	kase().Js("y?.setter;").Token(rhino.Token_OPTCHAIN_GETPROP).Expect(true),
	kase().Js("y.normal;").Token(rhino.Token_GETPROP).Expect(false),
	kase().Js("y?.normal;").Token(rhino.Token_OPTCHAIN_GETPROP).Expect(false),
	kase().Js("const {getter} = y;").Token(rhino.Token_STRING_KEY).Expect(true),
	kase().Js("const {setter} = y;").Token(rhino.Token_STRING_KEY).Expect(false),
	kase().Js("const {normal} = y;").Token(rhino.Token_STRING_KEY).Expect(false),
	kase().Js("y.getter = 0;").Token(rhino.Token_GETPROP).Expect(true),
	kase().Js("y.setter = 0;").Token(rhino.Token_GETPROP).Expect(true),
	kase().Js("y.normal = 0;").Token(rhino.Token_GETPROP).Expect(false),

	// Dynamic import causes side effects
	kase().Js("import('./module.js')").Token(rhino.Token_DYNAMIC_IMPORT).Expect(true),
}
