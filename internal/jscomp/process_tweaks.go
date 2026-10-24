// Port of com/google/javascript/jscomp/ProcessTweaks.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ProcessTweaks processes goog.tweak primitives. Checks that:
//
//   - parameters to goog.tweak.register* are literals of the correct type.
//   - the parameter to goog.tweak.get* is a string literal.
//   - parameters to goog.tweak.overrideDefaultValue are literals of the
//     correct type.
//   - tweak IDs passed to goog.tweak.get* and goog.tweak.overrideDefaultValue
//     correspond to registered tweaks.
//   - all calls to goog.tweak.register* and goog.tweak.overrideDefaultValue
//     are within the top-level context.
//   - each tweak is registered only once.
//   - calls to goog.tweak.overrideDefaultValue occur before the call to the
//     corresponding goog.tweak.register* function.
type ProcessTweaks struct {
	compiler    AbstractCompiler
	stripTweaks bool
}

// processTweaks_ID_MATCHER_MatchesAllOf is
// ID_MATCHER.matchesAllOf(sequence), where ID_MATCHER is
// CharMatcher.inRange('a', 'z').or(CharMatcher.inRange('A', 'Z')).or(CharMatcher.anyOf("0123456789_.")).
//
// Go-only: language-adapter — Guava CharMatcher, over the UTF-16 chars of
// the sequence.
func processTweaks_ID_MATCHER_MatchesAllOf(sequence string) bool {
	for _, c := range jutil.StringToUTF16(sequence) {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_' || c == '.') {
			return false
		}
	}
	return true
}

// Warnings and Errors.
var ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR = DiagnosticType_Error(
	"JSC_TWEAK_MULTIPLY_REGISTERED_ERROR", "Tweak {0} has already been registered.")

var ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR = DiagnosticType_Error("JSC_NON_LITERAL_TWEAK_ID_ERROR", "tweak ID must be a string literal")

var ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING = DiagnosticType_Warning(
	"JSC_INVALID_TWEAK_DEFAULT_VALUE_WARNING",
	"tweak {0} registered with {1} must have a default value that is a literal of type {2}")

var ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR = DiagnosticType_Error(
	"JSC_NON_GLOBAL_TWEAK_INIT_ERROR",
	"tweak declaration {0} must occur in the global scope")

var ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING = DiagnosticType_Warning(
	"JSC_TWEAK_WRONG_GETTER_TYPE_WARNING",
	"tweak getter function {0} used for tweak registered using {1}")

var ProcessTweaks_INVALID_TWEAK_ID_ERROR = DiagnosticType_Error(
	"JSC_INVALID_TWEAK_ID_ERROR",
	"tweak ID contains illegal characters. Only letters, numbers, _ and . are allowed")

// ProcessTweaks_TweakFunction is an enum of goog.tweak functions.
type ProcessTweaks_TweakFunction int32

const (
	ProcessTweaks_TweakFunction_REGISTER_BOOLEAN ProcessTweaks_TweakFunction = iota
	ProcessTweaks_TweakFunction_REGISTER_NUMBER
	ProcessTweaks_TweakFunction_REGISTER_STRING
	ProcessTweaks_TweakFunction_GET_BOOLEAN
	ProcessTweaks_TweakFunction_GET_NUMBER
	ProcessTweaks_TweakFunction_GET_STRING
)

// processTweaks_TweakFunction_null stands for a null TweakFunction.
//
// Go-only: language-adapter — Java null for an enum-typed field.
const processTweaks_TweakFunction_null ProcessTweaks_TweakFunction = -1

var processTweaks_TweakFunction_names = [...]string{
	"REGISTER_BOOLEAN", "REGISTER_NUMBER", "REGISTER_STRING", "GET_BOOLEAN", "GET_NUMBER", "GET_STRING",
}

// processTweaks_TweakFunction_fields holds the enum constants' fields,
// indexed by ordinal.
type processTweaks_TweakFunction_fields struct {
	name             string
	expectedTypeName *string
	validNodeTypeA   rhino.Token
	validNodeTypeB   rhino.Token
	registerFunction ProcessTweaks_TweakFunction
}

var processTweaks_TweakFunction_table = [...]processTweaks_TweakFunction_fields{
	newProcessTweaks_TweakFunctionWithTypeB("goog.tweak.registerBoolean", "boolean", rhino.Token_TRUE, rhino.Token_FALSE),
	newProcessTweaks_TweakFunction("goog.tweak.registerNumber", "number", rhino.Token_NUMBER),
	newProcessTweaks_TweakFunction("goog.tweak.registerString", "string", rhino.Token_STRINGLIT),
	newProcessTweaks_TweakFunctionWithRegisterFunction("goog.tweak.getBoolean", ProcessTweaks_TweakFunction_REGISTER_BOOLEAN),
	newProcessTweaks_TweakFunctionWithRegisterFunction("goog.tweak.getNumber", ProcessTweaks_TweakFunction_REGISTER_NUMBER),
	newProcessTweaks_TweakFunctionWithRegisterFunction("goog.tweak.getString", ProcessTweaks_TweakFunction_REGISTER_STRING),
}

// ProcessTweaks_TweakFunction_values is TweakFunction.values().
func ProcessTweaks_TweakFunction_values() []ProcessTweaks_TweakFunction {
	return []ProcessTweaks_TweakFunction{
		ProcessTweaks_TweakFunction_REGISTER_BOOLEAN,
		ProcessTweaks_TweakFunction_REGISTER_NUMBER,
		ProcessTweaks_TweakFunction_REGISTER_STRING,
		ProcessTweaks_TweakFunction_GET_BOOLEAN,
		ProcessTweaks_TweakFunction_GET_NUMBER,
		ProcessTweaks_TweakFunction_GET_STRING,
	}
}

// Name is Enum.name().
func (f ProcessTweaks_TweakFunction) Name() string { return processTweaks_TweakFunction_names[f] }

// String is Enum.toString().
func (f ProcessTweaks_TweakFunction) String() string { return f.Name() }

// Ordinal is Enum.ordinal().
func (f ProcessTweaks_TweakFunction) Ordinal() int { return int(f) }

// Java: ProcessTweaks.TweakFunction(String, String, Token)
func newProcessTweaks_TweakFunction(name string, expectedTypeName string, validNodeTypeA rhino.Token) processTweaks_TweakFunction_fields {
	return newProcessTweaks_TweakFunctionWithAll(name, &expectedTypeName, validNodeTypeA, rhino.Token_EMPTY, processTweaks_TweakFunction_null)
}

// Java: ProcessTweaks.TweakFunction(String, String, Token, Token)
func newProcessTweaks_TweakFunctionWithTypeB(
	name string, expectedTypeName string, validNodeTypeA rhino.Token, validNodeTypeB rhino.Token) processTweaks_TweakFunction_fields {
	return newProcessTweaks_TweakFunctionWithAll(name, &expectedTypeName, validNodeTypeA, validNodeTypeB, processTweaks_TweakFunction_null)
}

// Java: ProcessTweaks.TweakFunction(String, TweakFunction)
func newProcessTweaks_TweakFunctionWithRegisterFunction(name string, registerFunction ProcessTweaks_TweakFunction) processTweaks_TweakFunction_fields {
	return newProcessTweaks_TweakFunctionWithAll(name, nil, rhino.Token_EMPTY, rhino.Token_EMPTY, registerFunction)
}

// Java: ProcessTweaks.TweakFunction(String, String, Token, Token, TweakFunction)
func newProcessTweaks_TweakFunctionWithAll(
	name string,
	expectedTypeName *string,
	validNodeTypeA rhino.Token,
	validNodeTypeB rhino.Token,
	registerFunction ProcessTweaks_TweakFunction) processTweaks_TweakFunction_fields {
	f := processTweaks_TweakFunction_fields{}
	f.name = name
	f.expectedTypeName = expectedTypeName
	f.validNodeTypeA = validNodeTypeA
	f.validNodeTypeB = validNodeTypeB
	f.registerFunction = registerFunction
	return f
}

// Java: ProcessTweaks.TweakFunction.isValidNodeType(Token)
func (f ProcessTweaks_TweakFunction) IsValidNodeType(type_ rhino.Token) bool {
	return type_ == processTweaks_TweakFunction_table[f].validNodeTypeA || type_ == processTweaks_TweakFunction_table[f].validNodeTypeB
}

// Java: ProcessTweaks.TweakFunction.isCorrectRegisterFunction(TweakFunction)
func (f ProcessTweaks_TweakFunction) IsCorrectRegisterFunction(registerFunction ProcessTweaks_TweakFunction) bool {
	if registerFunction == processTweaks_TweakFunction_null {
		panic(jutil.NewNullPointerException())
	}
	return processTweaks_TweakFunction_table[f].registerFunction == registerFunction
}

// Java: ProcessTweaks.TweakFunction.isGetterFunction()
func (f ProcessTweaks_TweakFunction) IsGetterFunction() bool {
	return processTweaks_TweakFunction_table[f].registerFunction != processTweaks_TweakFunction_null
}

// Java: ProcessTweaks.TweakFunction.getName()
func (f ProcessTweaks_TweakFunction) GetName() string {
	return processTweaks_TweakFunction_table[f].name
}

// Java: ProcessTweaks.TweakFunction.getExpectedTypeName()
func (f ProcessTweaks_TweakFunction) GetExpectedTypeName() *string {
	return processTweaks_TweakFunction_table[f].expectedTypeName
}

// Java: ProcessTweaks.TweakFunction.createDefaultValueNode()
func (f ProcessTweaks_TweakFunction) CreateDefaultValueNode() *rhino.Node {
	switch f {
	case ProcessTweaks_TweakFunction_REGISTER_BOOLEAN:
		return rhino.IR_FalseNode()
	case ProcessTweaks_TweakFunction_REGISTER_NUMBER:
		return rhino.IR_Number(0)
	case ProcessTweaks_TweakFunction_REGISTER_STRING:
		return rhino.IR_String("")
	default:
		panic(jutil.NewIllegalStateException())
	}
}

// A map of function name -> TweakFunction.
var processTweaks_TWEAK_FUNCTIONS_MAP = func() *jutil.LinkedHashMap[string, ProcessTweaks_TweakFunction] {
	m := jutil.NewLinkedHashMap[string, ProcessTweaks_TweakFunction]()
	for _, func_ := range ProcessTweaks_TweakFunction_values() {
		m.Put(func_.GetName(), func_)
	}
	return m
}()

// Java: ProcessTweaks(AbstractCompiler, boolean)
func NewProcessTweaks(compiler AbstractCompiler, stripTweaks bool) *ProcessTweaks {
	p := &ProcessTweaks{}
	p.compiler = compiler
	p.stripTweaks = stripTweaks
	return p
}

// Java: ProcessTweaks.process(Node, Node)
func (p *ProcessTweaks) Process(externs *rhino.Node, root *rhino.Node) {
	result := p.CollectTweaks(root)

	if p.stripTweaks {
		p.StripAllCalls(result.tweakInfos)
	}
}

// StripAllCalls removes all CALL nodes in the given TweakInfos, replacing
// calls to getter functions with the tweak's default value.
//
// Java: ProcessTweaks.stripAllCalls(Map<String, TweakInfo>)
func (p *ProcessTweaks) StripAllCalls(tweakInfos *jutil.LinkedHashMap[string, *ProcessTweaks_TweakInfo]) {
	for tweakInfo := range tweakInfos.Values() {
		isRegistered := tweakInfo.IsRegistered()
		for functionCall := range tweakInfo.functionCalls.All() {
			callNode := functionCall.callNode
			parent := callNode.GetParent()
			if functionCall.tweakFunc.IsGetterFunction() {
				var newValue *rhino.Node
				if isRegistered {
					newValue = tweakInfo.GetDefaultValueNode().CloneNode()
				} else {
					// When we find a getter of an unregistered tweak, there has
					// already been a warning about it, so now just use a default
					// value when stripping.
					registerFunction := processTweaks_TweakFunction_table[functionCall.tweakFunc].registerFunction
					newValue = registerFunction.CreateDefaultValueNode()
				}
				callNode.ReplaceWith(newValue)
				p.compiler.ReportChangeToEnclosingScope(parent)
			} else {
				voidZeroNode := rhino.IR_VoidNode(rhino.IR_Number(0).Srcref(callNode)).Srcref(callNode)
				callNode.ReplaceWith(voidZeroNode)
				p.compiler.ReportChangeToEnclosingScope(parent)
			}
		}
	}
}

// CollectTweaks finds all calls to goog.tweak functions and emits
// warnings/errors if any of the calls have issues.
//
// Returns a map of TweakInfo structures, keyed by tweak ID.
//
// Java: ProcessTweaks.collectTweaks(Node)
func (p *ProcessTweaks) CollectTweaks(root *rhino.Node) *ProcessTweaks_CollectTweaksResult {
	pass := newProcessTweaks_CollectTweaks(p)
	NodeTraversal_Traverse(p.compiler, root, pass)

	tweakInfos := pass.allTweaks
	for tweakInfo := range tweakInfos.Values() {
		tweakInfo.EmitAllWarnings()
	}
	return newProcessTweaks_CollectTweaksResult(tweakInfos)
}

// ProcessTweaks_CollectTweaksResult is CollectTweaksResult.
type ProcessTweaks_CollectTweaksResult struct {
	tweakInfos *jutil.LinkedHashMap[string, *ProcessTweaks_TweakInfo]
}

// Java: ProcessTweaks.CollectTweaksResult(Map<String, TweakInfo>)
func newProcessTweaks_CollectTweaksResult(tweakInfos *jutil.LinkedHashMap[string, *ProcessTweaks_TweakInfo]) *ProcessTweaks_CollectTweaksResult {
	r := &ProcessTweaks_CollectTweaksResult{}
	r.tweakInfos = tweakInfos
	return r
}

// ProcessTweaks_CollectTweaks processes all calls to goog.tweak functions.
type ProcessTweaks_CollectTweaks struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *ProcessTweaks

	allTweaks *jutil.LinkedHashMap[string, *ProcessTweaks_TweakInfo]
}

// Java: ProcessTweaks.CollectTweaks()
func newProcessTweaks_CollectTweaks(outer *ProcessTweaks) *ProcessTweaks_CollectTweaks {
	c := &ProcessTweaks_CollectTweaks{outer: outer}
	c.allTweaks = jutil.NewLinkedHashMap[string, *ProcessTweaks_TweakInfo]()
	return c
}

// Java: ProcessTweaks.CollectTweaks.visit(NodeTraversal, Node, Node)
func (c *ProcessTweaks_CollectTweaks) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsCall() {
		return
	}

	callName := n.GetFirstChild().GetQualifiedName()
	tweakFunc := processTweaks_TweakFunction_null
	if callName != nil {
		if f, ok := processTweaks_TWEAK_FUNCTIONS_MAP.Lookup(*callName); ok {
			tweakFunc = f
		}
	}
	if tweakFunc == processTweaks_TweakFunction_null {
		return
	}

	// Ensure the first parameter (the tweak ID) is a string literal.
	tweakIdNode := n.GetSecondChild()
	if !tweakIdNode.IsStringLit() {
		c.outer.compiler.Report(JSError_MakeWithNode(tweakIdNode, ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR))
		return
	}
	tweakId := tweakIdNode.GetString()

	// Make sure there is a TweakInfo structure for it.
	tweakInfo := c.allTweaks.ComputeIfAbsent(tweakId, func(tweakId string) *ProcessTweaks_TweakInfo {
		return newProcessTweaks_TweakInfo(c.outer, tweakId)
	})

	switch tweakFunc {
	case ProcessTweaks_TweakFunction_REGISTER_BOOLEAN, ProcessTweaks_TweakFunction_REGISTER_NUMBER, ProcessTweaks_TweakFunction_REGISTER_STRING:
		// Ensure the ID contains only valid characters.
		if !processTweaks_ID_MATCHER_MatchesAllOf(tweakId) {
			c.outer.compiler.Report(JSError_MakeWithNode(tweakIdNode, ProcessTweaks_INVALID_TWEAK_ID_ERROR))
		}

		// Ensure tweaks are registered in the global scope.
		if !t.InGlobalHoistScope() {
			c.outer.compiler.Report(JSError_MakeWithNode(n, ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR, tweakId))
			break
		}

		// Ensure tweaks are registered only once.
		if tweakInfo.IsRegistered() {
			c.outer.compiler.Report(JSError_MakeWithNode(n, ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR, tweakId))
			break
		}

		tweakDefaultValueNode := tweakIdNode.GetNext().GetNext()
		tweakInfo.AddRegisterCall(tweakFunc, n, tweakDefaultValueNode)
	case ProcessTweaks_TweakFunction_GET_BOOLEAN, ProcessTweaks_TweakFunction_GET_NUMBER, ProcessTweaks_TweakFunction_GET_STRING:
		tweakInfo.AddGetterCall(tweakFunc, n)
	}
}

// ProcessTweaks_TweakFunctionCall holds information about a call to a
// goog.tweak function.
type ProcessTweaks_TweakFunctionCall struct {
	tweakFunc ProcessTweaks_TweakFunction
	callNode  *rhino.Node
	valueNode *rhino.Node
}

// Java: ProcessTweaks.TweakFunctionCall(TweakFunction, Node)
func newProcessTweaks_TweakFunctionCall(tweakFunc ProcessTweaks_TweakFunction, callNode *rhino.Node) *ProcessTweaks_TweakFunctionCall {
	return newProcessTweaks_TweakFunctionCallWithValueNode(tweakFunc, callNode, nil)
}

// Java: ProcessTweaks.TweakFunctionCall(TweakFunction, Node, Node)
func newProcessTweaks_TweakFunctionCallWithValueNode(tweakFunc ProcessTweaks_TweakFunction, callNode *rhino.Node, valueNode *rhino.Node) *ProcessTweaks_TweakFunctionCall {
	c := &ProcessTweaks_TweakFunctionCall{}
	c.callNode = callNode
	c.tweakFunc = tweakFunc
	c.valueNode = valueNode
	return c
}

// ProcessTweaks_TweakInfo stores information about a single tweak.
type ProcessTweaks_TweakInfo struct {
	outer *ProcessTweaks

	tweakId          string
	functionCalls    *jutil.ArrayList[*ProcessTweaks_TweakFunctionCall]
	registerCall     *ProcessTweaks_TweakFunctionCall
	defaultValueNode *rhino.Node
}

// Java: ProcessTweaks.TweakInfo(String)
func newProcessTweaks_TweakInfo(outer *ProcessTweaks, tweakId string) *ProcessTweaks_TweakInfo {
	i := &ProcessTweaks_TweakInfo{outer: outer}
	i.tweakId = tweakId
	i.functionCalls = jutil.NewArrayList[*ProcessTweaks_TweakFunctionCall]()
	return i
}

// EmitAllWarnings: if this tweak is registered, then looks for type
// warnings in default value parameters and getter functions. If it is not
// registered, emits an error for each function call.
//
// Java: ProcessTweaks.TweakInfo.emitAllWarnings()
func (i *ProcessTweaks_TweakInfo) EmitAllWarnings() {
	if i.IsRegistered() {
		i.EmitAllTypeWarnings()
	}
}

// EmitAllTypeWarnings emits a warning for each default value parameter
// that has the wrong type and for each getter function that was used for
// the wrong type of tweak.
//
// Java: ProcessTweaks.TweakInfo.emitAllTypeWarnings()
func (i *ProcessTweaks_TweakInfo) EmitAllTypeWarnings() {
	for call := range i.functionCalls.All() {
		valueNode := call.valueNode
		tweakFunc := call.tweakFunc
		registerFunc := i.registerCall.tweakFunc
		if valueNode != nil {
			// For register* and overrideDefaultValue calls, ensure the default
			// value is a literal of the correct type.
			if !registerFunc.IsValidNodeType(valueNode.GetToken()) {
				i.outer.compiler.Report(
					JSError_MakeWithNode(
						valueNode,
						ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING,
						i.tweakId,
						registerFunc.GetName(),
						jutil.ToJavaString(registerFunc.GetExpectedTypeName())))
			}
		} else if tweakFunc.IsGetterFunction() {
			// For getter calls, ensure the correct getter was used.
			if !tweakFunc.IsCorrectRegisterFunction(registerFunc) {
				i.outer.compiler.Report(
					JSError_MakeWithNode(
						call.callNode,
						ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING,
						tweakFunc.GetName(),
						registerFunc.GetName()))
			}
		}
	}
}

// Java: ProcessTweaks.TweakInfo.addRegisterCall(TweakFunction, Node, Node)
func (i *ProcessTweaks_TweakInfo) AddRegisterCall(tweakFunc ProcessTweaks_TweakFunction, callNode *rhino.Node, defaultValueNode *rhino.Node) {
	i.registerCall = newProcessTweaks_TweakFunctionCallWithValueNode(tweakFunc, callNode, defaultValueNode)
	i.functionCalls.Add(i.registerCall)
}

// Java: ProcessTweaks.TweakInfo.addGetterCall(TweakFunction, Node)
func (i *ProcessTweaks_TweakInfo) AddGetterCall(tweakFunc ProcessTweaks_TweakFunction, callNode *rhino.Node) {
	i.functionCalls.Add(newProcessTweaks_TweakFunctionCall(tweakFunc, callNode))
}

// Java: ProcessTweaks.TweakInfo.isRegistered()
func (i *ProcessTweaks_TweakInfo) IsRegistered() bool {
	return i.registerCall != nil
}

// Java: ProcessTweaks.TweakInfo.getDefaultValueNode()
func (i *ProcessTweaks_TweakInfo) GetDefaultValueNode() *rhino.Node {
	jutil.CheckState(i.IsRegistered())
	// Use calls to goog.tweak.overrideDefaultValue() first.
	if i.defaultValueNode != nil {
		return i.defaultValueNode
	}
	// Use the value passed to the register function next.
	if i.registerCall.valueNode != nil {
		return i.registerCall.valueNode
	}
	// Otherwise, use the default value for the tweak's type.
	return i.registerCall.tweakFunc.CreateDefaultValueNode()
}
