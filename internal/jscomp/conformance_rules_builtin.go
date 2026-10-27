// Port of com/google/javascript/jscomp/ConformanceRules.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: ConformanceRules.RestrictedNameCall.Restriction
type ConformanceRules_RestrictedNameCall_Restriction struct {
	name               rhino.Node
	restrictedCallType rhino.FunctionType
}

// Java: ConformanceRules.RestrictedNameCall.Restriction(Node, FunctionType)
func NewConformanceRules_RestrictedNameCall_Restriction(name rhino.Node, restrictedCallType rhino.FunctionType) *ConformanceRules_RestrictedNameCall_Restriction {
	r := &ConformanceRules_RestrictedNameCall_Restriction{}
	r.name = name
	r.restrictedCallType = restrictedCallType
	return r
}

// Java: ConformanceRules.RestrictedNameCall(AbstractCompiler, Requirement)
func NewConformanceRules_RestrictedNameCall(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_RestrictedNameCall {
	r := &ConformanceRules_RestrictedNameCall{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}

	builder := []*ConformanceRules_RestrictedNameCall_Restriction{}
	for _, value := range requirement.GetValueList() {
		var name rhino.Node
		if nameFromValue := ConformanceRules_RestrictedNameCall_GetNameFromValue(value); nameFromValue != nil {
			name = NodeUtil_NewQName(compiler, *nameFromValue)
		} else {
			// NodeUtil.newQName(compiler, null) throws a NullPointerException.
			panic(jutil.NewNullPointerException())
		}
		restrictedDecl := ConformanceRules_ConformanceUtil_GetTypeFromValue(value)
		if name.IsNil() || restrictedDecl == nil {
			panic(NewCheckConformance_InvalidRequirementSpec("bad prop value"))
		}

		restrictedCallType :=
			ConformanceRules_ConformanceUtil_EvaluateTypeString(compiler, *restrictedDecl).ToMaybeFunctionType()
		if restrictedCallType == nil {
			panic(NewCheckConformance_InvalidRequirementSpec("invalid conformance type"))
		}
		builder = append(builder, NewConformanceRules_RestrictedNameCall_Restriction(name, restrictedCallType))
	}
	r.restrictions = builder
	return r
}

// Java: ConformanceRules.RestrictedNameCall.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_RestrictedNameCall) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if ConformanceRules_ConformanceUtil_IsCallTarget(n) && n.IsQualifiedName() {
		// TODO(johnlenz): restrict to global names
		for i := 0; i < len(r.restrictions); i++ {
			res := r.restrictions[i]

			if n.MatchesQualifiedNameNode(res.name) {
				if !ConformanceRules_ConformanceUtil_ValidateCall(
					r.compiler, n.GetParent(), res.restrictedCallType, false) {
					return ConformanceRules_ConformanceResult_VIOLATION
				}
			} else if n.IsGetProp() &&
				n.GetString() == "call" &&
				n.GetFirstChild().MatchesQualifiedNameNode(res.name) {
				if !ConformanceRules_ConformanceUtil_ValidateCall(
					r.compiler, n.GetParent(), res.restrictedCallType, true) {
					return ConformanceRules_ConformanceResult_VIOLATION
				}
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.RestrictedNameCall.getNameFromValue(String)
func ConformanceRules_RestrictedNameCall_GetNameFromValue(specName string) *string {
	index := jutil.StringIndexOfChar(specName, ':')
	if index < 1 {
		return nil
	}
	s := jutil.StringSubstring(specName, 0, index)
	return &s
}

// ConformanceRules_RestrictedMethodCall is the banned property call rule.
type ConformanceRules_RestrictedMethodCall struct {
	ConformanceRules_AbstractRuleBase

	restrictions []*ConformanceRules_RestrictedMethodCall_Restriction
}

// Java: ConformanceRules.RestrictedMethodCall.Restriction
type ConformanceRules_RestrictedMethodCall_Restriction struct {
	type_              rhino.JSType
	property           string
	restrictedCallType rhino.FunctionType
}

// Java: ConformanceRules.RestrictedMethodCall.Restriction(JSType, String, FunctionType)
func NewConformanceRules_RestrictedMethodCall_Restriction(
	type_ rhino.JSType, property string, restrictedCallType rhino.FunctionType) *ConformanceRules_RestrictedMethodCall_Restriction {
	r := &ConformanceRules_RestrictedMethodCall_Restriction{}
	r.type_ = type_
	r.property = property
	r.restrictedCallType = restrictedCallType
	return r
}

// Java: ConformanceRules.RestrictedMethodCall(AbstractCompiler, Requirement)
func NewConformanceRules_RestrictedMethodCall(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_RestrictedMethodCall {
	r := &ConformanceRules_RestrictedMethodCall{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)

	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}

	registry := compiler.GetTypeRegistry()
	builder := []*ConformanceRules_RestrictedMethodCall_Restriction{}
	for _, value := range requirement.GetValueList() {
		type_ :=
			ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(ConformanceRules_ConformanceUtil_RemoveTypeDecl(value))
		property :=
			ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(ConformanceRules_ConformanceUtil_RemoveTypeDecl(value))
		restrictedDecl := ConformanceRules_ConformanceUtil_GetTypeFromValue(value)
		if type_ == nil || property == nil || restrictedDecl == nil {
			panic(NewCheckConformance_InvalidRequirementSpec("bad prop value"))
		}

		restrictedCallType :=
			ConformanceRules_ConformanceUtil_EvaluateTypeString(compiler, *restrictedDecl).ToMaybeFunctionType()
		if restrictedCallType == nil {
			panic(NewCheckConformance_InvalidRequirementSpec("invalid conformance type"))
		}
		builder = append(builder, NewConformanceRules_RestrictedMethodCall_Restriction(registry.GetGlobalType(*type_), *property, restrictedCallType))
	}

	r.restrictions = builder
	return r
}

// Java: ConformanceRules.RestrictedMethodCall.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_RestrictedMethodCall) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if !n.IsGetProp() || !ConformanceRules_ConformanceUtil_IsCallTarget(n) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	for i := 0; i < len(r.restrictions); i++ {
		res := r.restrictions[i]
		result := ConformanceRules_ConformanceResult_CONFORMANCE

		if r.MatchesProp(n, res) {
			result = r.CheckConformanceWithRestriction(t, n, res, false)
		} else if n.GetString() == "call" && r.MatchesProp(n.GetFirstChild(), res) {
			// handle .call invocation
			result = r.CheckConformanceWithRestriction(t, n, res, true)
		}

		// TODO(johnlenz): should "apply" always be a possible violation?
		if result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
			return result
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.RestrictedMethodCall.checkConformance(NodeTraversal, Node, Restriction, boolean)
func (r *ConformanceRules_RestrictedMethodCall) CheckConformanceWithRestriction(
	t *NodeTraversal, n rhino.Node, res *ConformanceRules_RestrictedMethodCall_Restriction, isCallInvocation bool) *ConformanceRules_ConformanceResult {
	registry := t.GetCompiler().GetTypeRegistry()
	methodClassType := res.type_
	var lhs rhino.Node
	if isCallInvocation {
		lhs = n.GetFirstFirstChild()
	} else {
		lhs = n.GetFirstChild()
	}
	if methodClassType != nil && lhs.GetJSType() != nil {
		targetType := lhs.GetJSType().RestrictByNotNullOrUndefined()
		if ConformanceRules_ConformanceUtil_IsLooseType(targetType) || ConformanceRules_IsLooseObject(targetType, registry) {
			if r.reportLooseTypeViolations &&
				!ConformanceRules_ConformanceUtil_ValidateCall(
					r.compiler, n.GetParent(), res.restrictedCallType, isCallInvocation) {
				return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
			}
		} else if targetType.IsSubtypeOf(methodClassType) {
			if !ConformanceRules_ConformanceUtil_ValidateCall(
				r.compiler, n.GetParent(), res.restrictedCallType, isCallInvocation) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.RestrictedMethodCall.matchesProp(Node, Restriction)
func (r *ConformanceRules_RestrictedMethodCall) MatchesProp(n rhino.Node, res *ConformanceRules_RestrictedMethodCall_Restriction) bool {
	return n.IsGetProp() && n.GetString() == res.property
}

// ConformanceRules_RestrictedPropertyWrite is the restricted property write
// rule.
type ConformanceRules_RestrictedPropertyWrite struct {
	ConformanceRules_AbstractRuleBase

	restrictions []*ConformanceRules_RestrictedPropertyWrite_Restriction
}

// Java: ConformanceRules.RestrictedPropertyWrite.Restriction
type ConformanceRules_RestrictedPropertyWrite_Restriction struct {
	type_          rhino.JSType
	property       string
	restrictedType rhino.JSType
}

// Java: ConformanceRules.RestrictedPropertyWrite.Restriction(JSType, String, JSType)
func NewConformanceRules_RestrictedPropertyWrite_Restriction(
	type_ rhino.JSType, property string, restrictedType rhino.JSType) *ConformanceRules_RestrictedPropertyWrite_Restriction {
	r := &ConformanceRules_RestrictedPropertyWrite_Restriction{}
	r.type_ = type_
	r.property = property
	r.restrictedType = restrictedType
	return r
}

// Java: ConformanceRules.RestrictedPropertyWrite(AbstractCompiler, Requirement)
func NewConformanceRules_RestrictedPropertyWrite(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_RestrictedPropertyWrite {
	r := &ConformanceRules_RestrictedPropertyWrite{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)

	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}

	registry := compiler.GetTypeRegistry()
	builder := []*ConformanceRules_RestrictedPropertyWrite_Restriction{}
	for _, value := range requirement.GetValueList() {
		type_ :=
			ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(ConformanceRules_ConformanceUtil_RemoveTypeDecl(value))
		property :=
			ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(ConformanceRules_ConformanceUtil_RemoveTypeDecl(value))
		restrictedDecl := ConformanceRules_ConformanceUtil_GetTypeFromValue(value)
		if type_ == nil || property == nil || restrictedDecl == nil {
			panic(NewCheckConformance_InvalidRequirementSpec("bad prop value"))
		}
		restrictedType := ConformanceRules_ConformanceUtil_EvaluateTypeString(compiler, *restrictedDecl)
		builder = append(builder, NewConformanceRules_RestrictedPropertyWrite_Restriction(registry.GetGlobalType(*type_), *property, restrictedType))
	}

	r.restrictions = builder
	return r
}

// Java: ConformanceRules.RestrictedPropertyWrite.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_RestrictedPropertyWrite) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsGetProp() && NodeUtil_IsLhsOfAssign(n) {
		rhsType := n.GetNext().GetJSType()
		targetType := n.GetFirstChild().GetJSType()
		if rhsType != nil && targetType != nil {
			var targetNotNullType rhino.JSType
			for _, res := range r.restrictions {
				if n.GetString() == res.property {
					if !rhsType.IsSubtypeOf(res.restrictedType) {
						if ConformanceRules_ConformanceUtil_IsLooseType(targetType) {
							if r.reportLooseTypeViolations {
								return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
							}
						} else {
							if targetNotNullType == nil {
								targetNotNullType = targetType.RestrictByNotNullOrUndefined()
							}
							if targetNotNullType.IsSubtypeOf(res.type_) {
								return ConformanceRules_ConformanceResult_VIOLATION
							}
						}
					}
				}
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BannedCodePattern is the banned code pattern rule.
type ConformanceRules_BannedCodePattern struct {
	ConformanceRules_AbstractRuleBase

	restrictions []*TemplateAstMatcher
}

// Java: ConformanceRules.BannedCodePattern(AbstractCompiler, Requirement)
func NewConformanceRules_BannedCodePattern(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BannedCodePattern {
	r := &ConformanceRules_BannedCodePattern{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)

	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}

	builder := []*TemplateAstMatcher{}
	for _, value := range requirement.GetValueList() {
		parseRoot :=
			NewCompilerInput(SourceFile_FromCode("<template>", value)).GetAstRoot(compiler)
		if !parseRoot.HasOneChild() || !parseRoot.GetFirstChild().IsFunction() {
			panic(NewCheckConformance_InvalidRequirementSpec("invalid conformance template: " + value))
		}
		templateRoot := parseRoot.GetFirstChild()
		astMatcher := NewTemplateAstMatcher(compiler, templateRoot, r.typeMatchingStrategy)
		builder = append(builder, astMatcher)
	}

	r.restrictions = builder
	return r
}

// Java: ConformanceRules.BannedCodePattern.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedCodePattern) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	possibleViolation := false
	for i := 0; i < len(r.restrictions); i++ {
		matcher := r.restrictions[i]
		if matcher.Matches(n) {
			if matcher.IsLooseMatch() {
				possibleViolation = true
			} else {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	if possibleViolation && r.reportLooseTypeViolations {
		return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_CustomRuleProxy is a custom rule proxy, for rules that we
// load dynamically.
type ConformanceRules_CustomRuleProxy struct {
	customRule CheckConformance_Rule
}

// NewConformanceRules_CustomRuleProxy panics with InvalidRequirementSpec (a
// checked exception).
//
// Java: ConformanceRules.CustomRuleProxy(AbstractCompiler, Requirement)
func NewConformanceRules_CustomRuleProxy(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_CustomRuleProxy {
	r := &ConformanceRules_CustomRuleProxy{}
	if !requirement.HasJavaClass() {
		panic(NewCheckConformance_InvalidRequirementSpec("missing java_class"))
	}
	r.customRule = r.CreateRule(compiler, requirement)
	return r
}

// Java: ConformanceRules.CustomRuleProxy.getPrecondition()
func (r *ConformanceRules_CustomRuleProxy) GetPrecondition() CheckConformance_Precondition {
	return r.customRule.GetPrecondition()
}

// Java: ConformanceRules.CustomRuleProxy.check(NodeTraversal, Node, LibraryLevelNonAllowlistedConformanceViolationsBehavior)
func (r *ConformanceRules_CustomRuleProxy) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	r.customRule.Check(t, n, behavior)
}

// Java: ConformanceRules.CustomRuleProxy.createRule(AbstractCompiler, Requirement)
func (r *ConformanceRules_CustomRuleProxy) CreateRule(compiler AbstractCompiler, requirement *Requirement) CheckConformance_Rule {
	factory, ok := conformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY.Lookup(requirement.GetJavaClass())
	if !ok {
		panic(NewCheckConformance_InvalidRequirementSpec("Unknown conformance rule: " + requirement.GetJavaClass()))
	}
	return factory.Create(compiler, requirement)
}

// ConformanceRules_CustomRuleProxy_CustomRuleFactory creates a custom rule
// for a requirement. Create panics with InvalidRequirementSpec (a checked
// exception).
//
// Java: ConformanceRules.CustomRuleProxy.CustomRuleFactory
type ConformanceRules_CustomRuleProxy_CustomRuleFactory interface {
	// Java: ConformanceRules.CustomRuleProxy.CustomRuleFactory.create(AbstractCompiler, Requirement)
	Create(compiler AbstractCompiler, requirement *Requirement) CheckConformance_Rule
}

// ConformanceRules_CustomRuleProxy_CustomRuleFactoryFunc adapts a function to
// ConformanceRules_CustomRuleProxy_CustomRuleFactory.
//
// Go-only: language-adapter — the Java method references (NoOp::new, ...).
type ConformanceRules_CustomRuleProxy_CustomRuleFactoryFunc func(compiler AbstractCompiler, requirement *Requirement) CheckConformance_Rule

// Java: ConformanceRules.CustomRuleProxy.CustomRuleFactory.create(AbstractCompiler, Requirement)
func (f ConformanceRules_CustomRuleProxy_CustomRuleFactoryFunc) Create(compiler AbstractCompiler, requirement *Requirement) CheckConformance_Rule {
	return f(compiler, requirement)
}

// Java: ConformanceRules.CustomRuleProxy.CUSTOM_RULE_REGISTRY
var conformanceRules_CustomRuleProxy_CUSTOM_RULE_REGISTRY = conformanceRules_CustomRuleProxy_CreateCustomRuleRegistry()

// Java: ConformanceRules.CustomRuleProxy.createCustomRuleRegistry()
func conformanceRules_CustomRuleProxy_CreateCustomRuleRegistry() *jutil.LinkedHashMap[string, ConformanceRules_CustomRuleProxy_CustomRuleFactory] {
	builder := jutil.NewLinkedHashMap[string, ConformanceRules_CustomRuleProxy_CustomRuleFactory]()
	put := func(className string, factory ConformanceRules_CustomRuleProxy_CustomRuleFactoryFunc) {
		builder.Put(className, factory)
	}
	put("com.google.javascript.jscomp.ConformanceRules$NoOp",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_NoOp(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$InferredConstCheck",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_InferredConstCheck(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BannedDependency",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BannedDependency(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BannedDependencyRegex",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BannedDependencyRegex(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BannedName",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BannedName(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BannedProperty",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BannedProperty(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BannedStringRegex",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BannedStringRegex(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$RestrictedNameCall",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_RestrictedNameCall(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$RestrictedMethodCall",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_RestrictedMethodCall(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$RestrictedPropertyWrite",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_RestrictedPropertyWrite(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BannedCodePattern",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BannedCodePattern(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanForOf",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanForOf(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$RequireUseStrict",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_RequireUseStrict(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanThrowOfNonErrorTypes(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanNullDeref",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanNullDeref(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanUnknownThis",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanUnknownThis(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanUnknownDirectThisPropsReferences",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanUnknownDirectThisPropsReferences(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanUnknownTypedClassPropsReferences",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanUnknownTypedClassPropsReferences(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanUnresolvedType",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanUnresolvedType(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$StrictBanUnresolvedType",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_StrictBanUnresolvedType(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanGlobalVars",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanGlobalVars(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BannedEnhance",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BannedEnhance(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BannedModsRegex",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BannedModsRegex(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanCreateElement",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanCreateElement(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanCreateDom",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanCreateDom(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanElementSetAttribute(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanSettingAttributes(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanExecCommand",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanExecCommand(c, r)
		})
	put("com.google.javascript.jscomp.ConformanceRules$BanStaticThis",
		func(c AbstractCompiler, r *Requirement) CheckConformance_Rule {
			return NewConformanceRules_BanStaticThis(c, r)
		})
	return builder
}

// ConformanceRules_BanForOf bans for/of loops.
type ConformanceRules_BanForOf struct {
	ConformanceRules_AbstractRuleBase
}

// Java: ConformanceRules.BanForOf(AbstractCompiler, Requirement)
func NewConformanceRules_BanForOf(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanForOf {
	r := &ConformanceRules_BanForOf{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.BanForOf.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanForOf) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsForOf() || n.IsForAwaitOf() {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_RequireUseStrict is the require "use strict" rule.
type ConformanceRules_RequireUseStrict struct {
	ConformanceRules_AbstractRuleBase
}

// Java: ConformanceRules.RequireUseStrict(AbstractCompiler, Requirement)
func NewConformanceRules_RequireUseStrict(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_RequireUseStrict {
	r := &ConformanceRules_RequireUseStrict{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	if len(requirement.GetValueList()) != 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("invalid value"))
	}
	return r
}

// Java: ConformanceRules.RequireUseStrict.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_RequireUseStrict) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsScript() && !n.IsUseStrict() {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanThrowOfNonErrorTypes bans throw of non-error object
// types.
type ConformanceRules_BanThrowOfNonErrorTypes struct {
	ConformanceRules_AbstractRuleBase

	errorObjType rhino.JSType
}

// Java: ConformanceRules.BanThrowOfNonErrorTypes(AbstractCompiler, Requirement)
func NewConformanceRules_BanThrowOfNonErrorTypes(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanThrowOfNonErrorTypes {
	r := &ConformanceRules_BanThrowOfNonErrorTypes{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	r.errorObjType = compiler.GetTypeRegistry().GetGlobalType("Error")
	return r
}

// Java: ConformanceRules.BanThrowOfNonErrorTypes.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanThrowOfNonErrorTypes) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if r.errorObjType != nil && n.IsThrow() {
		thrown := n.GetFirstChild().GetJSType()
		if thrown != nil {
			// Allow vague types, as is typical of re-throws of exceptions
			if !thrown.IsUnknownType() &&
				!thrown.IsAllType() &&
				!thrown.IsEmptyType() &&
				!thrown.IsSubtypeOf(r.errorObjType) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanNullDeref bans dereferencing null or undefined types.
type ConformanceRules_BanNullDeref struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

// Java: ConformanceRules.BanNullDeref(AbstractCompiler, Requirement)
func NewConformanceRules_BanNullDeref(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanNullDeref {
	r := &ConformanceRules_BanNullDeref{}
	r.InitConformanceRules_AbstractTypeRestrictionRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.BanNullDeref.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanNullDeref) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {

	var violation bool
	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_GETELEM, rhino.Token_NEW, rhino.Token_CALL:
		violation = r.ReportNode(n.GetFirstChild())
	case rhino.Token_IN:
		violation = r.ReportNode(n.GetLastChild())
	default:
		violation = false
	}

	if violation {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ReportNode is BanNullDeref's report(Node).
//
// Go-only: language-adapter — named ReportNode, not Report: the inherited
// overload report(Node, ConformanceResult, ...) keeps the name Report, which
// the AbstractRule interface requires.
//
// Java: ConformanceRules.BanNullDeref.report(Node)
func (r *ConformanceRules_BanNullDeref) ReportNode(n rhino.Node) bool {
	return n.GetJSType() != nil && ConformanceRules_AbstractTypeRestrictionRule_IsKnown(n) && r.InvalidDeref(n) && !r.IsAllowlistedType(n)
}

// InvalidDeref reports whether the type is known to be invalid to
// dereference.
//
// Java: ConformanceRules.BanNullDeref.invalidDeref(Node)
func (r *ConformanceRules_BanNullDeref) InvalidDeref(n rhino.Node) bool {
	type_ := n.GetJSType()
	// TODO(johnlenz): top type should not be allowed here
	return !type_.IsAllType() && (type_.IsNullable() || type_.IsVoidable())
}

// ConformanceRules_BanUnknownThis bans unknown "this" types.
type ConformanceRules_BanUnknownThis struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase

	// LinkedIdentityHashSet<Node>
	reports *jutil.LinkedHashSet[rhino.Node]
}

// Java: ConformanceRules.BanUnknownThis(AbstractCompiler, Requirement)
func NewConformanceRules_BanUnknownThis(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanUnknownThis {
	r := &ConformanceRules_BanUnknownThis{}
	r.reports = jutil.NewLinkedHashSet[rhino.Node]()
	r.InitConformanceRules_AbstractTypeRestrictionRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.BanUnknownThis.tsIsAllowlisted()
func (r *ConformanceRules_BanUnknownThis) TsIsAllowlisted() bool {
	// We expect TypeScript to already check it.
	// LINT.IfChange
	return true
	// LINT.ThenChange(//depot/google3/java/com/google/javascript/modules/librarydepsconformancechecker/conformance_checker.go)
}

// Java: ConformanceRules.BanUnknownThis.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanUnknownThis) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsThis() {
		type_ := n.GetJSType()
		if type_ != nil && type_.IsUnknownType() && !r.IsTypeImmediatelyTightened(n) {
			root := t.GetScopeRoot()
			if !r.reports.Contains(root) {
				r.reports.Add(root)
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanUnknownDirectThisPropsReferences bans unknown type
// references of the form "this.prop" unless
//   - it is immediately cast,
//   - it is a @template type (until template type restricts are enabled) or
//   - the value is unused.
//   - the "this" type is unknown (as this is expected to be used with
//     BanUnknownThis which would have already reported the root cause).
type ConformanceRules_BanUnknownDirectThisPropsReferences struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences(AbstractCompiler, Requirement)
func NewConformanceRules_BanUnknownDirectThisPropsReferences(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanUnknownDirectThisPropsReferences {
	r := &ConformanceRules_BanUnknownDirectThisPropsReferences{}
	r.InitConformanceRules_AbstractTypeRestrictionRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanUnknownDirectThisPropsReferences) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsGetProp() &&
		ConformanceRules_BanUnknownDirectThisPropsReferences_IsKnownThis(n.GetFirstChild()) && // not a cascading unknown
		ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n) &&
		!ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(n) &&
		r.IsUsed(n) && // skip most assignments, etc
		!r.IsTypeImmediatelyTightened(n) &&
		!ConformanceRules_BanUnknownDirectThisPropsReferences_IsExplicitlyUnknown(n) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences.isKnownThis(Node)
func ConformanceRules_BanUnknownDirectThisPropsReferences_IsKnownThis(n rhino.Node) bool {
	return n.IsThis() && !ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n)
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences.isExplicitlyUnknown(Node)
func ConformanceRules_BanUnknownDirectThisPropsReferences_IsExplicitlyUnknown(n rhino.Node) bool {
	owner := rhino.ObjectType_Cast(n.GetFirstChild().GetJSType())
	var prop *rhino.Property
	if owner != nil {
		prop = owner.GetSlot(n.GetString())
	}
	return prop != nil && !prop.IsTypeInferred()
}

// ConformanceRules_BanUnknownTypedClassPropsReferences bans unknown type
// references of the form "instance.prop" unless
//   - (a) it is immediately cast/asserted, or
//   - (b) it is a @template type (until template type restrictions are
//     enabled), or
//   - (c) the value is unused, or
//   - (d) the source object type is unknown (to avoid error cascades)
type ConformanceRules_BanUnknownTypedClassPropsReferences struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

// Java: ConformanceRules.BanUnknownTypedClassPropsReferences(AbstractCompiler, Requirement)
func NewConformanceRules_BanUnknownTypedClassPropsReferences(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanUnknownTypedClassPropsReferences {
	r := &ConformanceRules_BanUnknownTypedClassPropsReferences{}
	r.InitConformanceRules_AbstractTypeRestrictionRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.BanUnknownTypedClassPropsReferences.tsIsAllowlisted()
func (r *ConformanceRules_BanUnknownTypedClassPropsReferences) TsIsAllowlisted() bool {
	// LINT.IfChange
	return true
	// LINT.ThenChange(//depot/google3/java/com/google/javascript/modules/librarydepsconformancechecker/conformance_checker.go)
}

// Java: ConformanceRules.BanUnknownTypedClassPropsReferences.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanUnknownTypedClassPropsReferences) CheckConformance(t *NodeTraversal, getprop rhino.Node) *ConformanceRules_ConformanceResult {
	if getprop.IsGetProp() &&
		ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(getprop) &&
		r.IsUsed(getprop) && // skip most assignments, etc
		!r.IsTypeImmediatelyTightened(getprop) &&
		r.IsCheckablePropertySource(getprop.GetFirstChild()) && // not a cascading unknown
		!ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(getprop) &&
		!r.IsDeclaredUnknown(getprop) {
		propName := getprop.GetString()
		typeName := getprop.GetFirstChild().GetJSType().String()
		return NewConformanceRules_ConformanceResultWithNote(
			ConformanceRules_ConformanceLevel_VIOLATION,
			"The property \""+propName+"\" on type \""+typeName+"\"")
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanUnknownTypedClassPropsReferences.isCheckablePropertySource(Node)
func (r *ConformanceRules_BanUnknownTypedClassPropsReferences) IsCheckablePropertySource(n rhino.Node) bool {
	return ConformanceRules_AbstractTypeRestrictionRule_IsKnown(n) &&
		!ConformanceRules_AbstractTypeRestrictionRule_IsTop(n) &&
		r.IsClassType(n) &&
		!r.IsNativeObjectType(n) &&
		!r.IsAllowlistedType(n)
}

// Java: ConformanceRules.BanUnknownTypedClassPropsReferences.isClassType(Node)
func (r *ConformanceRules_BanUnknownTypedClassPropsReferences) IsClassType(n rhino.Node) bool {
	type_ := n.GetJSType().RestrictByNotNullOrUndefined().ToMaybeObjectType()
	if type_ == nil || !type_.IsInstanceType() {
		return false
	}

	ctor := type_.GetConstructor()
	if ctor == nil {
		return false
	}

	info := ctor.GetJSDocInfo()
	source := ctor.GetSource()

	return (info != nil && info.IsConstructorOrInterface()) ||
		(!source.IsNil() && source.IsClass())
}

// Java: ConformanceRules.BanUnknownTypedClassPropsReferences.isDeclaredUnknown(Node)
func (r *ConformanceRules_BanUnknownTypedClassPropsReferences) IsDeclaredUnknown(n rhino.Node) bool {
	target := n.GetFirstChild()
	targetType := target.GetJSType().RestrictByNotNullOrUndefined().ToMaybeObjectType()
	if targetType == nil {
		return false
	}

	info := targetType.GetPropertyJSDocInfo(n.GetString())
	if info == nil || !info.HasType() {
		return false
	}

	expr := info.GetType()
	typeExprNode := expr.GetRoot()
	if typeExprNode.GetToken() == rhino.Token_QMARK && !typeExprNode.HasChildren() {
		return true
	} else if typeExprNode.GetToken() == rhino.Token_PIPE {
		// Might be a union type including ? that's collapsed during checking.
		for child := typeExprNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.GetToken() == rhino.Token_QMARK {
				return true
			}
		}
	}

	return false
}

// ConformanceRules_BanUnresolvedType bans accessing properties from objects
// that are unresolved forward-declared type names. For legacy reasons this
// is allowed but causes unexpected weaknesses in the type inference.
type ConformanceRules_BanUnresolvedType struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

// Java: ConformanceRules.BanUnresolvedType(AbstractCompiler, Requirement)
func NewConformanceRules_BanUnresolvedType(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanUnresolvedType {
	r := &ConformanceRules_BanUnresolvedType{}
	r.InitConformanceRules_AbstractTypeRestrictionRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.BanUnresolvedType.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanUnresolvedType) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsGetProp() {
		target := n.GetFirstChild()
		type_ := target.GetJSType()
		nonConformingPart := ConformanceRules_BanUnresolvedType_GetNonConformingPart(type_)
		if nonConformingPart != nil && !r.IsTypeImmediatelyTightened(n) {
			return NewConformanceRules_ConformanceResultWithNote(
				ConformanceRules_ConformanceLevel_VIOLATION,
				"Reference to type '"+*nonConformingPart+"' never resolved.")
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanUnresolvedType.getNonConformingPart(JSType)
func ConformanceRules_BanUnresolvedType_GetNonConformingPart(type_ rhino.JSType) *string {
	if type_ == nil {
		return nil
	}
	if type_.IsUnionType() {
		// unwrap union types which might contain unresolved type name
		// references for example {Foo|undefined}
		var nonConformingParts []string
		for _, part := range type_.GetUnionMembers() {
			nonConformingPart := ConformanceRules_BanUnresolvedType_GetNonConformingPart(part)
			if nonConformingPart != nil {
				nonConformingParts = append(nonConformingParts, *nonConformingPart)
			}
		}
		if nonConformingParts != nil {
			joined := strings.Join(nonConformingParts, "|")
			return &joined
		}
	} else if type_.IsNoResolvedType() {
		noResolvedType := type_.ToObjectType()
		return noResolvedType.GetReferenceName()
	}
	return nil
}

// ConformanceRules_StrictBanUnresolvedType bans any use of unresolved
// forward-declared types.
type ConformanceRules_StrictBanUnresolvedType struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

// Java: ConformanceRules.StrictBanUnresolvedType(AbstractCompiler, Requirement)
func NewConformanceRules_StrictBanUnresolvedType(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_StrictBanUnresolvedType {
	r := &ConformanceRules_StrictBanUnresolvedType{}
	r.InitConformanceRules_AbstractTypeRestrictionRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.StrictBanUnresolvedType.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_StrictBanUnresolvedType) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	nonConformingPart := ConformanceRules_BanUnresolvedType_GetNonConformingPart(n.GetJSType())
	if nonConformingPart != nil && !r.IsTypeImmediatelyTightened(n) {
		return NewConformanceRules_ConformanceResultWithNote(
			ConformanceRules_ConformanceLevel_VIOLATION,
			"Reference to type '"+*nonConformingPart+"' never resolved.")
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanGlobalVars bans global var declarations.
type ConformanceRules_BanGlobalVars struct {
	ConformanceRules_AbstractRuleBase

	allowlistedNames *jutil.LinkedHashSet[string]
}

// Java: ConformanceRules.BanGlobalVars(AbstractCompiler, Requirement)
func NewConformanceRules_BanGlobalVars(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanGlobalVars {
	r := &ConformanceRules_BanGlobalVars{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)

	r.allowlistedNames = jutil.NewLinkedHashSet(requirement.GetValueList()...)
	return r
}

// Java: ConformanceRules.BanGlobalVars.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanGlobalVars) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if t.InGlobalScope() &&
		NodeUtil_IsDeclaration(n) &&
		!n.GetBooleanProp(rhino.Node_IS_NAMESPACE) &&
		!r.IsAllowlisted(n) {
		enclosingScript := NodeUtil_GetEnclosingScript(n)
		if !enclosingScript.IsNil() &&
			(enclosingScript.GetBooleanProp(rhino.Node_GOOG_MODULE) ||
				enclosingScript.GetBooleanProp(rhino.Node_ES6_MODULE) ||
				enclosingScript.GetInputId().Equals(r.compiler.GetSyntheticCodeInputId())) {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanGlobalVars.isAllowlisted(Node)
func (r *ConformanceRules_BanGlobalVars) IsAllowlisted(n rhino.Node) bool {
	if n.IsFromExterns() {
		return true
	}

	if n.IsFunction() {
		return r.IsAllowlistedName(n.GetFirstChild().GetString())
	}

	if NodeUtil_IsNameDeclaration(n) {
		allowlisted := true // for lambda access
		NodeUtil_VisitLhsNodesInNode(
			n,
			func(name rhino.Node) {
				if !r.IsAllowlistedName(name.GetString()) {
					allowlisted = false
				}
			})
		return allowlisted
	}

	return false
}

// Java: ConformanceRules.BanGlobalVars.isAllowlistedName(String)
func (r *ConformanceRules_BanGlobalVars) IsAllowlistedName(name string) bool {
	if r.allowlistedNames.Contains(name) {
		return true
	}
	// names created by JSCompiler internals
	return name == "$jscomp" ||
		strings.HasPrefix(name, "$jscomp$compprop") ||
		ClosureRewriteModule_IsModuleContent(name) ||
		ClosureRewriteModule_IsModuleExport(name) ||
		ScopedAliases_IsScopedAliases(name)
}

// ConformanceRules_BannedEnhance checks that file does not include an
// @enhance annotation for a banned namespace.
type ConformanceRules_BannedEnhance struct {
	ConformanceRules_AbstractRuleBase

	bannedEnhancedNamespaces *jutil.LinkedHashSet[string]
}

// Java: ConformanceRules.BannedEnhance(AbstractCompiler, Requirement)
func NewConformanceRules_BannedEnhance(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BannedEnhance {
	r := &ConformanceRules_BannedEnhance{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}
	r.bannedEnhancedNamespaces = jutil.NewLinkedHashSet(requirement.GetValueList()...)
	return r
}

// Java: ConformanceRules.BannedEnhance.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedEnhance) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	docInfo := n.GetJSDocInfo()

	if docInfo == nil || !docInfo.HasEnhance() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	for banned := range r.bannedEnhancedNamespaces.All() {
		if *docInfo.GetEnhance() == banned {
			return NewConformanceRules_ConformanceResultWithNote(
				ConformanceRules_ConformanceLevel_VIOLATION, "The enhanced namespace \""+banned+"\"")
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// conformanceRules_BannedEnhance_isScript is the Node::isScript method
// reference of BannedEnhance.getPrecondition() (a non-capturing lambda: one
// instance).
var conformanceRules_BannedEnhance_isScript CheckConformance_Precondition = NewCheckConformance_PreconditionFunc(
	(rhino.Node).IsScript)

// Java: ConformanceRules.BannedEnhance.getPrecondition()
func (r *ConformanceRules_BannedEnhance) GetPrecondition() CheckConformance_Precondition {
	return conformanceRules_BannedEnhance_isScript
}

// ConformanceRules_BannedModsRegex checks that file does not include an @mods
// annotation for a banned namespace regex.
type ConformanceRules_BannedModsRegex struct {
	ConformanceRules_AbstractRuleBase

	bannedModsRegex *jutil.Pattern
}

// Java: ConformanceRules.BannedModsRegex(AbstractCompiler, Requirement)
func NewConformanceRules_BannedModsRegex(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BannedModsRegex {
	r := &ConformanceRules_BannedModsRegex{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)

	bannedModsRegexList := requirement.GetValueList()
	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}
	r.bannedModsRegex = ConformanceRules_BuildPattern(bannedModsRegexList)
	return r
}

// Java: ConformanceRules.BannedModsRegex.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedModsRegex) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	docInfo := n.GetJSDocInfo()

	if docInfo == nil || !docInfo.HasMods() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if r.bannedModsRegex.Matcher(*docInfo.GetMods()).Find() {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// conformanceRules_BannedModsRegex_isScript is the Node::isScript method
// reference of BannedModsRegex.getPrecondition() (a non-capturing lambda:
// one instance).
var conformanceRules_BannedModsRegex_isScript CheckConformance_Precondition = NewCheckConformance_PreconditionFunc(
	(rhino.Node).IsScript)

// Java: ConformanceRules.BannedModsRegex.getPrecondition()
func (r *ConformanceRules_BannedModsRegex) GetPrecondition() CheckConformance_Precondition {
	return conformanceRules_BannedModsRegex_isScript
}

// Java: ConformanceRules.GOOG_DOM
var conformanceRules_GOOG_DOM = rhino.QualifiedName_Of("goog.dom")

// Java: ConformanceRules.GOOG_DOM_TAGNAME
var conformanceRules_GOOG_DOM_TAGNAME = rhino.QualifiedName_Of("goog.dom.TagName")

// ConformanceRules_BanCreateElement bans document.createElement and similar
// methods with string literal parameter specified in value, e.g. value:
// 'script'. The purpose of banning these is that they don't provide the type
// information which hinders other rules. Authors should use e.g.
// goog.dom.createElement(goog.dom.TagName.SCRIPT) which returns
// HTMLScriptElement.
type ConformanceRules_BanCreateElement struct {
	ConformanceRules_AbstractRuleBase

	bannedTags    *jutil.LinkedHashSet[string]
	domHelperType rhino.JSType
	documentType  rhino.JSType
}

// Java: ConformanceRules.BanCreateElement(AbstractCompiler, Requirement)
func NewConformanceRules_BanCreateElement(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanCreateElement {
	r := &ConformanceRules_BanCreateElement{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	r.bannedTags = jutil.NewLinkedHashSet[string]()
	for _, value := range requirement.GetValueList() {
		r.bannedTags.Add(jutil.StringToLowerCaseASCII(value))
	}
	if r.bannedTags.IsEmpty() {
		panic(NewCheckConformance_InvalidRequirementSpec("Specify one or more values."))
	}
	r.domHelperType = compiler.GetTypeRegistry().GetGlobalType("goog.dom.DomHelper")
	r.documentType = compiler.GetTypeRegistry().GetGlobalType("Document")
	return r
}

// Java: ConformanceRules.BanCreateElement.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanCreateElement) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsCall() {
		tag := n.GetSecondChild()
		if !tag.IsNil() &&
			tag.IsStringLit() &&
			r.bannedTags.Contains(jutil.StringToLowerCaseASCII(tag.GetString())) {
			return r.CheckCreateElement(n)
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanCreateElement.checkCreateElement(Node)
func (r *ConformanceRules_BanCreateElement) CheckCreateElement(n rhino.Node) *ConformanceRules_ConformanceResult {
	target := n.GetFirstChild()
	if !target.IsGetProp() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	functionName := target.GetString()
	if "createElement" != functionName && "createDom" != functionName {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	srcObj := target.GetFirstChild()
	if conformanceRules_GOOG_DOM.Matches(srcObj) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	type_ := srcObj.GetJSType()
	if type_ == nil || ConformanceRules_ConformanceUtil_IsLooseType(type_) {
		if r.reportLooseTypeViolations {
			return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
		}
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if (r.domHelperType != nil && r.domHelperType.IsSubtypeOf(type_)) ||
		(r.documentType != nil && r.documentType.IsSubtypeOf(type_)) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanCreateDom bans goog.dom.createDom and
// goog.dom.DomHelper#createDom with parameters specified in value in the
// format tagname.attribute, e.g. value: 'iframe.src'. Tag name might be also
// * to ban the attribute in any tag. Note that string literal values
// assigned to banned attributes are allowed as they couldn't be attacker
// controlled.
type ConformanceRules_BanCreateDom struct {
	ConformanceRules_AbstractRuleBase

	bannedTagAttrs [][2]string
	domHelperType  rhino.JSType
	classNameTypes rhino.JSType
}

// Java: ConformanceRules.BanCreateDom(AbstractCompiler, Requirement)
func NewConformanceRules_BanCreateDom(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanCreateDom {
	r := &ConformanceRules_BanCreateDom{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	bannedTagAttrs := [][2]string{}
	for _, value := range requirement.GetValueList() {
		tagAttr := conformanceRules_splitToList(conformanceRules_ON_DOT, value)
		if len(tagAttr) != 2 || tagAttr[0] == "" || tagAttr[1] == "" {
			panic(NewCheckConformance_InvalidRequirementSpec("Values must be in the format tagname.attribute."))
		}
		lowercasedTag := jutil.StringToLowerCaseRoot(tagAttr[0])
		bannedTagAttrs = append(bannedTagAttrs, [2]string{lowercasedTag, tagAttr[1]})
	}
	r.bannedTagAttrs = bannedTagAttrs
	if len(r.bannedTagAttrs) == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("Specify one or more values."))
	}
	r.domHelperType = compiler.GetTypeRegistry().GetGlobalType("goog.dom.DomHelper")
	r.classNameTypes =
		compiler.
			GetTypeRegistry().
			CreateUnionType(
				compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_STRING_TYPE),
				compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_ARRAY_TYPE),
				compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_NULL_TYPE),
				compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_VOID_TYPE))
	return r
}

// Java: ConformanceRules.BanCreateDom.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanCreateDom) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if !r.IsCreateDomCall(n) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if n.GetChildCount() < 3 {
		// goog.dom.createDom('iframe') is fine.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	tagNames, tagNamesIsNull := r.GetTagNames(n.GetSecondChild())
	attrs := n.GetChildAtIndex(2)
	attrsType := attrs.GetJSType()

	// TODO(tbreisacher): Remove this after the typechecker understands ES6.
	if attrsType == nil {
		// Type information is not available; don't run this check.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// String or array attribute sets the class.
	isClassName := !attrsType.IsUnknownType() && attrsType.IsSubtypeOf(r.classNameTypes)

	if attrs.IsNull() || (attrsType != nil && attrsType.IsVoidType()) {
		// goog.dom.createDom('iframe', null) is fine.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	for _, tagAttr := range r.bannedTagAttrs {
		if !tagNamesIsNull && !slices.Contains(tagNames, tagAttr[0]) && tagAttr[0] != "*" {
			continue
		}
		var violation *ConformanceRules_ConformanceResult
		if !tagNamesIsNull || tagAttr[0] == "*" {
			violation = ConformanceRules_ConformanceResult_VIOLATION
		} else if r.reportLooseTypeViolations {
			violation = ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION
		} else {
			violation = ConformanceRules_ConformanceResult_CONFORMANCE
		}
		if isClassName {
			if tagAttr[1] != "class" {
				continue
			}
			return violation
		}
		if tagAttr[1] == "textContent" &&
			n.GetChildCount() > 3 &&
			violation != ConformanceRules_ConformanceResult_CONFORMANCE {
			return violation
		}
		if !attrs.IsObjectLit() {
			// Attrs is not an object literal and tagName matches or is unknown.
			if r.reportLooseTypeViolations {
				return ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION
			}
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		prop := NodeUtil_GetFirstPropMatchingKey(attrs, tagAttr[1])
		if !prop.IsNil() {
			if NodeUtil_IsSomeCompileTimeConstStringValue(prop) {
				// Ignore string literal values.
				continue
			}
			return violation
		}

		for attr := attrs.GetFirstChild(); !attr.IsNil(); attr = attr.GetNext() {
			if attr.IsComputedProp() {
				// We don't know if the computed property matches 'src' or not
				if r.reportLooseTypeViolations {
					return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
				}
				return ConformanceRules_ConformanceResult_CONFORMANCE
			}
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// GetTagNames returns the (nullable) collection of tag names; isNull is
// Java's null.
//
// Java: ConformanceRules.BanCreateDom.getTagNames(Node)
func (r *ConformanceRules_BanCreateDom) GetTagNames(tag rhino.Node) (tagNames []string, isNull bool) {
	if tag.IsStringLit() {
		return []string{jutil.StringToLowerCaseRoot(tag.GetString())}, false
	} else if tag.IsGetProp() && conformanceRules_GOOG_DOM_TAGNAME.Matches(tag.GetFirstChild()) {
		return []string{jutil.StringToLowerCaseRoot(tag.GetString())}, false
	}
	// TODO(jakubvrana): Support union, e.g. {!TagName<!HTMLDivElement>|!TagName<!HTMLBRElement>}.
	type_ := tag.GetJSType()
	if type_ == nil || !type_.IsTemplatizedType() {
		return nil, true
	}
	typeAsObj := type_.ToMaybeObjectType()
	if *typeAsObj.GetRawType().GetDisplayName() == "goog.dom.TagName" {
		tagType := conformanceRules_getOnlyElement(typeAsObj.GetTemplateTypes())
		// ImmutableMultimap.get: an empty collection for a missing (or null)
		// key.
		names := []string{}
		if displayName := tagType.GetDisplayName(); displayName != nil {
			if v, ok := conformanceRules_BanCreateDom_ELEMENT_TAG_NAMES.Lookup(*displayName); ok {
				names = v
			}
		}
		return names, false
	}
	return nil, true
}

// conformanceRules_getOnlyElement is Guava's Iterables.getOnlyElement.
//
// Go-only: language-adapter — Guava Iterables.
func conformanceRules_getOnlyElement(types []rhino.JSType) rhino.JSType {
	if len(types) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if len(types) == 1 {
		return types[0]
	}
	var sb strings.Builder
	sb.WriteString("expected one element but was: <")
	sb.WriteString(jutil.ToJavaString(types[0]))
	i := 1
	for ; i < 5 && i < len(types); i++ {
		sb.WriteString(", ")
		sb.WriteString(jutil.ToJavaString(types[i]))
	}
	if i < len(types) {
		sb.WriteString(", ...")
	}
	sb.WriteByte('>')
	panic(jutil.NewIllegalArgumentExceptionMsg(sb.String()))
}

// Java: ConformanceRules.BanCreateDom.ELEMENT_TAG_NAMES (an ImmutableMultimap)
var conformanceRules_BanCreateDom_ELEMENT_TAG_NAMES = func() *jutil.LinkedHashMap[string, []string] {
	m := jutil.NewLinkedHashMap[string, []string]()
	put := func(k, v string) {
		m.Put(k, append(m.Get(k), v))
	}
	put("HTMLAnchorElement", "a")
	put("HTMLAppletElement", "applet")
	put("HTMLAreaElement", "area")
	put("HTMLAudioElement", "audio")
	put("HTMLBRElement", "br")
	put("HTMLBaseElement", "base")
	put("HTMLBaseFontElement", "basefont")
	put("HTMLBodyElement", "body")
	put("HTMLButtonElement", "button")
	put("HTMLCanvasElement", "canvas")
	put("HTMLDListElement", "dl")
	put("HTMLDataListElement", "datalist")
	put("HTMLDetailsElement", "details")
	put("HTMLDialogElement", "dialog")
	put("HTMLDirectoryElement", "dir")
	put("HTMLDivElement", "div")
	put("HTMLEmbedElement", "embed")
	put("HTMLFieldSetElement", "fieldset")
	put("HTMLFontElement", "font")
	put("HTMLFormElement", "form")
	put("HTMLFrameElement", "frame")
	put("HTMLFrameSetElement", "frameset")
	put("HTMLHRElement", "hr")
	put("HTMLHeadElement", "head")
	put("HTMLHeadingElement", "h1")
	put("HTMLHeadingElement", "h2")
	put("HTMLHeadingElement", "h3")
	put("HTMLHeadingElement", "h4")
	put("HTMLHeadingElement", "h5")
	put("HTMLHeadingElement", "h6")
	put("HTMLHtmlElement", "html")
	put("HTMLIFrameElement", "iframe")
	put("HTMLImageElement", "img")
	put("HTMLInputElement", "input")
	put("HTMLIsIndexElement", "isindex")
	put("HTMLLIElement", "li")
	put("HTMLLabelElement", "label")
	put("HTMLLegendElement", "legend")
	put("HTMLLinkElement", "link")
	put("HTMLMapElement", "map")
	put("HTMLMenuElement", "menu")
	put("HTMLMetaElement", "meta")
	put("HTMLMeterElement", "meter")
	put("HTMLModElement", "del")
	put("HTMLModElement", "ins")
	put("HTMLOListElement", "ol")
	put("HTMLObjectElement", "object")
	put("HTMLOptGroupElement", "optgroup")
	put("HTMLOptionElement", "option")
	put("HTMLOutputElement", "output")
	put("HTMLParagraphElement", "p")
	put("HTMLParamElement", "param")
	put("HTMLPreElement", "pre")
	put("HTMLProgressElement", "progress")
	put("HTMLQuoteElement", "blockquote")
	put("HTMLQuoteElement", "q")
	put("HTMLScriptElement", "script")
	put("HTMLSelectElement", "select")
	put("HTMLSourceElement", "source")
	put("HTMLSpanElement", "span")
	put("HTMLStyleElement", "style")
	put("HTMLTableCaptionElement", "caption")
	put("HTMLTableCellElement", "td")
	put("HTMLTableCellElement", "th")
	put("HTMLTableColElement", "col")
	put("HTMLTableColElement", "colgroup")
	put("HTMLTableElement", "table")
	put("HTMLTableRowElement", "tr")
	put("HTMLTableSectionElement", "tbody")
	put("HTMLTableSectionElement", "tfoot")
	put("HTMLTableSectionElement", "thead")
	put("HTMLTemplateElement", "template")
	put("HTMLTextAreaElement", "textarea")
	put("HTMLTitleElement", "title")
	put("HTMLTrackElement", "track")
	put("HTMLUListElement", "ul")
	put("HTMLVideoElement", "video")
	return m
}()

// Java: ConformanceRules.BanCreateDom.isCreateDomCall(Node)
func (r *ConformanceRules_BanCreateDom) IsCreateDomCall(n rhino.Node) bool {
	if !n.IsCall() {
		return false
	}
	target := n.GetFirstChild()
	if !target.IsGetProp() {
		return false
	}
	if "createDom" != target.GetString() {
		return false
	}

	srcObj := target.GetFirstChild()
	if conformanceRules_GOOG_DOM.Matches(srcObj) {
		return true
	}
	type_ := srcObj.GetJSType()
	if type_ == nil {
		return false
	}
	if type_.Equals(r.domHelperType) {
		return true
	}
	return false
}

// ConformanceRules_SecuritySensitiveAttributes checks nodes for conformance
// with banning the setting of attributes that are on the blocklist.
type ConformanceRules_SecuritySensitiveAttributes struct {
	bannedAtrrs             *jutil.LinkedHashSet[string]
	globalNamespaceSupplier func() *GlobalNamespace
}

// ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS is the
// security-sensitive attributes that are banned from being set.
//
// Making updates to these attributes requires a new JSCompiler release. You
// must test the change using a global presubmit "at head" and update any
// affected allowlists. See go/jscompiler-global-presubmit and
// go/tsjs-conformance-team-docs.
//
// Java: ConformanceRules.SecuritySensitiveAttributes.ALL_BANNED_ATTRS
var ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS = jutil.NewLinkedHashSet(
	"href",
	"rel",
	"src",
	"srcdoc",
	"action",
	"formaction",
	"sandbox",
	"icon",
	"codebase",
	"data")

// Java: ConformanceRules.SecuritySensitiveAttributes()
func NewConformanceRules_SecuritySensitiveAttributes() *ConformanceRules_SecuritySensitiveAttributes {
	return NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(
		ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS.ToSlice(), func() *GlobalNamespace { return nil })
}

// Java: ConformanceRules.SecuritySensitiveAttributes(Supplier<GlobalNamespace>)
func NewConformanceRules_SecuritySensitiveAttributesWithSupplier(globalNamespaceSupplier func() *GlobalNamespace) *ConformanceRules_SecuritySensitiveAttributes {
	return NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(
		ConformanceRules_SecuritySensitiveAttributes_ALL_BANNED_ATTRS.ToSlice(), globalNamespaceSupplier)
}

// Java: ConformanceRules.SecuritySensitiveAttributes(Collection<String>)
func NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrs(bannedAtrrs []string) *ConformanceRules_SecuritySensitiveAttributes {
	return NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(
		bannedAtrrs, func() *GlobalNamespace { return nil })
}

// Java: ConformanceRules.SecuritySensitiveAttributes(Collection<String>, Supplier<GlobalNamespace>)
func NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(
	bannedAtrrs []string, globalNamespaceSupplier func() *GlobalNamespace) *ConformanceRules_SecuritySensitiveAttributes {
	s := &ConformanceRules_SecuritySensitiveAttributes{}
	s.bannedAtrrs = jutil.NewLinkedHashSet[string]()
	for _, a := range bannedAtrrs {
		s.bannedAtrrs.Add(jutil.StringToLowerCaseRoot(a))
	}
	s.globalNamespaceSupplier = globalNamespaceSupplier
	return s
}

// Contains checks if a attribute name is on the security banlist. Callers
// should make sure the attribute name is lower-cased, as attribute names are
// case-insensitve in HTML.
//
// Java: ConformanceRules.SecuritySensitiveAttributes.contains(String)
func (s *ConformanceRules_SecuritySensitiveAttributes) Contains(attributeName string) bool {
	return s.bannedAtrrs.Contains(attributeName)
}

// CheckConformanceForAttributeName checks, given a NodeTraversal and Node, if
// the attribute violates conformance.
//
// A violation is returned if the attribute name cannot be determined (and it
// is not an xid), if the attribute is on a list of banned attributes, or if
// it begins with the letters "on". Otherwise, it is a conforming attribute.
//
// Java: ConformanceRules.SecuritySensitiveAttributes.checkConformanceForAttributeName(NodeTraversal, Node)
func (s *ConformanceRules_SecuritySensitiveAttributes) CheckConformanceForAttributeName(
	traversal *NodeTraversal, attrName rhino.Node) *ConformanceRules_ConformanceResult {
	literalName :=
		ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), attrName, s.globalNamespaceSupplier)
	if literalName == nil {
		// xid() obfuscates attribute names, thus never clashing with security-sensitive attributes.
		if ConformanceRules_ConformanceUtil_IsXid(attrName.GetJSType()) {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		return ConformanceRules_ConformanceResult_VIOLATION
	}

	lowerName := jutil.StringToLowerCaseRoot(*literalName)

	if s.bannedAtrrs.Contains(lowerName) ||
		ConformanceRules_ConformanceUtil_IsEventHandlerAttrName(lowerName) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// CheckConformanceForAttributeNameWithHighConfidence checks, given a
// NodeTraversal and Node, if the attribute violates conformance.
//
// A violation is returned only if the attribute name can be statically
// determined and is on the list of banned attributes.
//
// Java: ConformanceRules.SecuritySensitiveAttributes.checkConformanceForAttributeNameWithHighConfidence(NodeTraversal, Node)
func (s *ConformanceRules_SecuritySensitiveAttributes) CheckConformanceForAttributeNameWithHighConfidence(
	traversal *NodeTraversal, attrName rhino.Node) *ConformanceRules_ConformanceResult {
	literalName :=
		ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), attrName, s.globalNamespaceSupplier)
	if literalName == nil {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if s.Contains(jutil.StringToLowerCaseRoot(*literalName)) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanElementSetAttribute bans Element#setAttribute with
// attribute names specified in value or any dynamic string.
//
// Using the element access syntax to set properties of Element is considered
// to be equivalent to calling Element#setAttribute. E.g., element[prop] =
// value is treated as element.setAttribute(prop, value) when element has the
// type Element.
type ConformanceRules_BanElementSetAttribute struct {
	ConformanceRules_AbstractRuleBase

	performGlobalNamespaceAnalysis bool
	globalNamespace                *GlobalNamespace

	elementType                      rhino.JSType
	securitySensitiveAttributes      *ConformanceRules_SecuritySensitiveAttributes
	defaultDecisionForUncertainCases *ConformanceRules_ConformanceResult
}

// Java: ConformanceRules.BanElementSetAttribute.ELEMENT_TYPE_NAME
const conformanceRules_BanElementSetAttribute_ELEMENT_TYPE_NAME = "Element"

// Java: ConformanceRules.BanElementSetAttribute.SET_ATTRIBUTE
const conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE = "setAttribute"

// Java: ConformanceRules.BanElementSetAttribute.SET_ATTRIBUTE_NS
const conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS = "setAttributeNS"

// Java: ConformanceRules.BanElementSetAttribute.SET_ATTRIBUTE_NODE
const conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE = "setAttributeNode"

// Java: ConformanceRules.BanElementSetAttribute.SET_ATTRIBUTE_NODE_NS
const conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE_NS = "setAttributeNodeNS"

// Java: ConformanceRules.BanElementSetAttribute.BANNED_PROPERTIES
var conformanceRules_BanElementSetAttribute_BANNED_PROPERTIES = jutil.NewLinkedHashSet(
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE,
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS,
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE,
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE_NS)

// The cap for script size, beyond which we give up performing global
// namespace analysis due to exccessive performance cost. The value is purely
// herustic and subject to future regulation.
//
// Java: ConformanceRules.BanElementSetAttribute.GLOBAL_NAMESPACE_ANALYSIS_LIMIT
const conformanceRules_BanElementSetAttribute_GLOBAL_NAMESPACE_ANALYSIS_LIMIT = 10000

// NewConformanceRules_BanElementSetAttribute creates a custom checker to ban
// Element#setAttribute based on conformance a requirement spec. It panics
// with InvalidRequirementSpec if the requirement is malformed.
//
// Java: ConformanceRules.BanElementSetAttribute(AbstractCompiler, Requirement)
func NewConformanceRules_BanElementSetAttribute(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanElementSetAttribute {
	r := &ConformanceRules_BanElementSetAttribute{}
	r.performGlobalNamespaceAnalysis = true
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	r.securitySensitiveAttributes =
		NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(requirement.GetValueList(), r.GetGlobalNamespace)
	if requirement.GetReportLooseTypeViolations() {
		r.defaultDecisionForUncertainCases = ConformanceRules_ConformanceResult_VIOLATION
	} else {
		r.defaultDecisionForUncertainCases = ConformanceRules_ConformanceResult_CONFORMANCE
	}
	r.elementType = compiler.GetTypeRegistry().GetGlobalType(conformanceRules_BanElementSetAttribute_ELEMENT_TYPE_NAME)
	return r
}

// GetGlobalNamespace builds GlobalNamespace starting from the Root node of
// the input, excluding externs. Skip the build if the input AST is too large.
//
// Java: ConformanceRules.BanElementSetAttribute.getGlobalNamespace()
func (r *ConformanceRules_BanElementSetAttribute) GetGlobalNamespace() *GlobalNamespace {
	if r.performGlobalNamespaceAnalysis && r.globalNamespace == nil {
		if NodeUtil_CountAstSizeUpToLimit(r.compiler.GetJsRoot(), conformanceRules_BanElementSetAttribute_GLOBAL_NAMESPACE_ANALYSIS_LIMIT) >=
			conformanceRules_BanElementSetAttribute_GLOBAL_NAMESPACE_ANALYSIS_LIMIT {
			r.performGlobalNamespaceAnalysis = false
		} else {
			r.globalNamespace = NewGlobalNamespace(r.compiler, r.compiler.GetJsRoot())
		}
	}
	return r.globalNamespace
}

// Java: ConformanceRules.BanElementSetAttribute.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if node.IsCall() {
		return r.CheckConformanceOnPropertyCall(traversal, node)
	} else if node.IsGetElem() && NodeUtil_IsLValue(node) {
		return r.CheckConformanceOnGetElement(traversal, node)
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanElementSetAttribute.checkConformanceOnPropertyCall(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) CheckConformanceOnPropertyCall(
	traversal *NodeTraversal, callNode rhino.Node) *ConformanceRules_ConformanceResult {
	calledProperty := r.GetBannedPropertyName(callNode)
	if calledProperty == nil {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if *calledProperty == conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE {
		// If the setAttribute function is called with less than two arguments it's either a false
		// positive or uncompilable code.
		if callNode.GetChildCount() < 3 {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}

		if r.requirement.GetReportLooseTypeViolations() {
			return r.securitySensitiveAttributes.CheckConformanceForAttributeName(
				traversal, callNode.GetSecondChild())
		} else {
			return r.securitySensitiveAttributes.CheckConformanceForAttributeNameWithHighConfidence(
				traversal, callNode.GetSecondChild())
		}
	}

	if *calledProperty == conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS {
		// If the setAttributeNS function is called with less than three arguments it's either a
		// false positive or uncompilable code.
		if callNode.GetChildCount() < 4 {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}

		if callNode.GetSecondChild().IsNull() {
			// A null namespace makes this equivalent to a regular setAttribute call.
			if r.requirement.GetReportLooseTypeViolations() {
				return r.securitySensitiveAttributes.CheckConformanceForAttributeName(
					traversal, callNode.GetChildAtIndex(2))
			} else {
				return r.securitySensitiveAttributes.CheckConformanceForAttributeNameWithHighConfidence(
					traversal, callNode.GetChildAtIndex(2))
			}
		}
	}

	// We haven't defined a security contract on setAttributeNode and setAttributeNodeNS yet, so
	// flag them as risky if report_loose_type_violations is set.
	return r.defaultDecisionForUncertainCases
}

// GetBannedPropertyName returns the Optional<String> as a nullable string.
//
// Java: ConformanceRules.BanElementSetAttribute.getBannedPropertyName(Node)
func (r *ConformanceRules_BanElementSetAttribute) GetBannedPropertyName(callNode rhino.Node) *string {
	target := callNode.GetFirstChild()
	if !target.IsGetProp() {
		return nil
	}
	propertyName := target.GetString()
	if !conformanceRules_BanElementSetAttribute_BANNED_PROPERTIES.Contains(propertyName) {
		return nil
	}
	type_ := target.GetFirstChild().GetJSType()
	if type_ == nil {
		return nil
	}
	if r.elementType != nil && type_.RestrictByNotNullOrUndefined().IsSubtypeOf(r.elementType) {
		return &propertyName
	}
	return nil
}

// CheckConformanceOnGetElement checks if there is any security implication
// in a GetElem node.
//
// Returns true if the checked Node is setting a security sensitive property
// on an Element object through the bracket syntax.
//
// Java: ConformanceRules.BanElementSetAttribute.checkConformanceOnGetElement(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) CheckConformanceOnGetElement(
	traversal *NodeTraversal, getElementNode rhino.Node) *ConformanceRules_ConformanceResult {
	key := getElementNode.GetSecondChild()
	keyName :=
		ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), key, r.GetGlobalNamespace)
	if keyName != nil {
		if r.securitySensitiveAttributes.Contains(jutil.StringToLowerCaseRoot(*keyName)) ||
			(r.requirement.GetReportLooseTypeViolations() &&
				// innerHTML and outerHTML are not attributes but properties. For simplicity we
				// check them here instead of having another rule.
				(*keyName == "innerHTML" || *keyName == "outerHTML")) {
			if r.HasElementType(getElementNode) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if !r.requirement.GetReportLooseTypeViolations() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// We have seen code using other types of keys (e.g., numbers) for element access. Only
	// report violations if the key is explicitly typed as string or the union of string and
	// other types.
	keyType := key.GetJSType()
	if keyType == nil || ConformanceRules_ConformanceUtil_IsXid(keyType) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if r.HasElementType(getElementNode) {
		if keyType.IsString() {
			return ConformanceRules_ConformanceResult_VIOLATION
		} else if keyType.IsUnionType() {
			for _, alternate := range keyType.ToMaybeUnionType().GetAlternates() {
				if alternate.IsString() {
					return ConformanceRules_ConformanceResult_VIOLATION
				}
			}
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanElementSetAttribute.hasElementType(Node)
func (r *ConformanceRules_BanElementSetAttribute) HasElementType(getElementNode rhino.Node) bool {
	objType := getElementNode.GetFirstChild().GetJSType()

	// Do not further check the node if there's no type information available.
	if objType == nil || r.elementType == nil {
		return false
	}
	objType = objType.RestrictByNotNullOrUndefined()
	// If the type of the object is not known to be a subtype of Element, the code is not
	// equivalent to setAttribute. Without this check, we will get overwhelmed by false
	// positives.
	if objType.IsUnknownType() || objType.IsEmptyType() || !objType.IsSubtypeOf(r.elementType) {
		return false
	}

	return true
}

// ConformanceRules_BanSettingAttributes can be used to ban any function like
// Element#setAttribute which takes two parameters, an attribute name and a
// value.
//
// The specific attribute names that are banned are specified in
// SecuritySensitiveAttributes.BANNED_ATTRS.
type ConformanceRules_BanSettingAttributes struct {
	ConformanceRules_AbstractRuleBase

	restrictions []*ConformanceRules_BanSettingAttributes_AttributeSettingRestriction
}

// Java: ConformanceRules.BanSettingAttributes.AttributeSettingRestriction
type ConformanceRules_BanSettingAttributes_AttributeSettingRestriction struct {
	type_    rhino.JSType
	property string
}

// Java: ConformanceRules.BanSettingAttributes.AttributeSettingRestriction(JSType, String)
func NewConformanceRules_BanSettingAttributes_AttributeSettingRestriction(
	type_ rhino.JSType, property string) *ConformanceRules_BanSettingAttributes_AttributeSettingRestriction {
	a := &ConformanceRules_BanSettingAttributes_AttributeSettingRestriction{}
	a.type_ = type_
	a.property = property
	return a
}

// Java: ConformanceRules.BanSettingAttributes.BANNED_ATTRS
var conformanceRules_BanSettingAttributes_BANNED_ATTRS = NewConformanceRules_SecuritySensitiveAttributes()

// NewConformanceRules_BanSettingAttributes creates a custom checker to ban a
// function like Element#setAttribute based on a conformance requirement spec.
//
// Names in value fields indicate the functions that should be blocked.
// Panics with InvalidRequirementSpec if the requirement is malformed.
//
// Java: ConformanceRules.BanSettingAttributes(AbstractCompiler, Requirement)
func NewConformanceRules_BanSettingAttributes(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanSettingAttributes {
	r := &ConformanceRules_BanSettingAttributes{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)

	if len(requirement.GetValueList()) == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}

	registry := compiler.GetTypeRegistry()
	builder := []*ConformanceRules_BanSettingAttributes_AttributeSettingRestriction{}
	for _, value := range requirement.GetValueList() {
		type_ := ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(value)
		property := ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(value)
		if type_ == nil || property == nil {
			panic(NewCheckConformance_InvalidRequirementSpec("bad prop value"))
		}
		builder = append(builder, NewConformanceRules_BanSettingAttributes_AttributeSettingRestriction(registry.GetGlobalType(*type_), *property))
	}

	r.restrictions = builder
	return r
}

// Java: ConformanceRules.BanSettingAttributes.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanSettingAttributes) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if node.IsCall() {
		return r.CheckConformanceOnPropertyCall(traversal, node)
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanSettingAttributes.checkConformanceOnPropertyCall(NodeTraversal, Node)
func (r *ConformanceRules_BanSettingAttributes) CheckConformanceOnPropertyCall(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	calledProperty := r.GetBannedPropertyName(node)
	if calledProperty == nil {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// This node is conformant if it has less than three children (the function and two
	// arguments). It's either reading the attribute, a false positive, or uncompilable code.
	if node.GetChildCount() < 3 {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	return conformanceRules_BanSettingAttributes_BANNED_ATTRS.CheckConformanceForAttributeName(traversal, node.GetSecondChild())
}

// GetBannedPropertyName returns the Optional<String> as a nullable string.
//
// Java: ConformanceRules.BanSettingAttributes.getBannedPropertyName(Node)
func (r *ConformanceRules_BanSettingAttributes) GetBannedPropertyName(node rhino.Node) *string {
	target := node.GetFirstChild()
	if !target.IsGetProp() {
		return nil
	}
	type_ := target.GetFirstChild().GetJSType()
	if type_ == nil {
		return nil
	}
	propertyName := target.GetString()
	for _, restricted := range r.restrictions {
		if propertyName == restricted.property &&
			restricted.type_ != nil &&
			type_.IsSubtypeOf(restricted.type_) {
			return &propertyName
		}
	}
	return nil
}

// ConformanceRules_BanExecCommand bans Document#execCommand with command
// names specified in value or any dynamic string.
type ConformanceRules_BanExecCommand struct {
	ConformanceRules_AbstractRuleBase

	bannedAttrs []string
}

// NewConformanceRules_BanExecCommand creates a custom checker to ban
// Document#execCommand based on a conformance requirement spec. Names in
// value fields indicate the attribute names that should be blocked. Panics
// with InvalidRequirementSpec if the requirement is malformed.
//
// Java: ConformanceRules.BanExecCommand(AbstractCompiler, Requirement)
func NewConformanceRules_BanExecCommand(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanExecCommand {
	r := &ConformanceRules_BanExecCommand{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	r.bannedAttrs = []string{}
	for _, v := range requirement.GetValueList() {
		r.bannedAttrs = append(r.bannedAttrs, jutil.StringToLowerCaseRoot(v))
	}
	return r
}

// Java: ConformanceRules.BanExecCommand.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanExecCommand) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if !r.IsBannedProperty(node) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// If the function is called without arguments then it's either a false positive or
	// uncompilable code.
	if node.GetChildCount() < 2 {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	attr := node.GetSecondChild()
	if !attr.IsStringLit() {
		return ConformanceRules_ConformanceResult_VIOLATION
	}

	attrName := attr.GetString()
	if slices.Contains(r.bannedAttrs, jutil.StringToLowerCaseRoot(attrName)) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanExecCommand.isBannedProperty(Node)
func (r *ConformanceRules_BanExecCommand) IsBannedProperty(node rhino.Node) bool {
	if !node.IsCall() {
		return false
	}
	target := node.GetFirstChild()
	if !target.IsGetProp() {
		return false
	}
	propertyName := target.GetString()
	if propertyName != "execCommand" {
		return false
	}
	type_ := target.GetFirstChild().GetJSType()
	if type_ == nil {
		return false
	}
	documentType := r.compiler.GetTypeRegistry().GetGlobalType("Document")
	if documentType == nil || !type_.RestrictByNotNullOrUndefined().IsSubtypeOf(documentType) {
		return false
	}
	return true
}

// ConformanceRules_BanStaticThis checks that this is not being referenced
// directly within a static member function.
type ConformanceRules_BanStaticThis struct {
	ConformanceRules_AbstractTypeRestrictionRuleBase
}

// Java: ConformanceRules.BanStaticThis(AbstractCompiler, Requirement)
func NewConformanceRules_BanStaticThis(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BanStaticThis {
	r := &ConformanceRules_BanStaticThis{}
	r.InitConformanceRules_AbstractTypeRestrictionRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.BanStaticThis.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanStaticThis) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if !n.IsThis() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	enclosingFunction := NodeUtil_GetEnclosingNonArrowFunction(n)
	if !enclosingFunction.IsNil() && ConformanceRules_BanStaticThis_IsStaticMethod(enclosingFunction) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanStaticThis.isStaticMethod(Node)
func ConformanceRules_BanStaticThis_IsStaticMethod(n rhino.Node) bool {
	jutil.CheckArgument(n.IsFunction())
	parent := n.GetParent()
	return !parent.IsNil() && parent.IsMemberFunctionDef() && parent.IsStaticMember()
}
