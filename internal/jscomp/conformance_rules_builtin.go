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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: ConformanceRules.RestrictedNameCall.Restriction
type ConformanceRules_RestrictedNameCall_Restriction struct {
	name               rhino.Node
	restrictedCallType rhino.FunctionType
}

// Java: ConformanceRules.RestrictedNameCall(AbstractCompiler, Requirement)
func newConformanceRules_RestrictedNameCall(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_RestrictedNameCall, error) {
	r := &ConformanceRules_RestrictedNameCall{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	restrictions := make([]ConformanceRules_RestrictedNameCall_Restriction, 0, len(values))
	for _, value := range values {
		nameFromValue, hasName := conformanceRules_getNameFromValue(value)
		if !hasName {
			// NodeUtil.newQName(compiler, null) throws a NullPointerException.
			panic(jutil.NewNullPointerException())
		}
		name := NodeUtil_NewQName(compiler, nameFromValue)
		restrictedDecl, hasDecl := conformanceRules_getTypeFromValue(value)
		if name.IsNil() || !hasDecl {
			return nil, errInvalidRequirementSpec("bad prop value")
		}

		evaluated, err := conformanceRules_evaluateTypeString(compiler, restrictedDecl)
		if err != nil {
			return nil, err
		}
		restrictedCallType := evaluated.ToMaybeFunctionType()
		if restrictedCallType == nil {
			return nil, errInvalidRequirementSpec("invalid conformance type")
		}
		restrictions = append(restrictions, ConformanceRules_RestrictedNameCall_Restriction{
			name:               name,
			restrictedCallType: restrictedCallType,
		})
	}
	r.restrictions = restrictions
	return r, nil
}

func (r *ConformanceRules_RestrictedNameCall) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.RestrictedNameCall.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_RestrictedNameCall) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if ConformanceRules_ConformanceUtil_IsCallTarget(n) && n.IsQualifiedName() {
		for i := 0; i < len(r.restrictions); i++ {
			res := &r.restrictions[i]

			if n.MatchesQualifiedNameNode(res.name) {
				if !ConformanceRules_ConformanceUtil_ValidateCall(
					r.compiler, n.GetParent(), res.restrictedCallType, false) {
					return ConformanceRules_ConformanceResult_VIOLATION
				}
			} else if gp, ok := ast.AsGetProp(n); ok &&
				gp.Property() == "call" &&
				gp.Object().Node().MatchesQualifiedNameNode(res.name) {
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
func conformanceRules_getNameFromValue(specName string) (string, bool) {
	index := strings.IndexByte(specName, ':')
	if index < 1 {
		return "", false
	}
	return specName[:index], true
}

// ConformanceRules_RestrictedMethodCall is the banned property call rule.
type ConformanceRules_RestrictedMethodCall struct {
	ConformanceRules_AbstractRule

	restrictions []ConformanceRules_RestrictedMethodCall_Restriction
}

// Java: ConformanceRules.RestrictedMethodCall.Restriction
type ConformanceRules_RestrictedMethodCall_Restriction struct {
	type_              rhino.JSType
	property           string
	restrictedCallType rhino.FunctionType
}

// Java: ConformanceRules.RestrictedMethodCall(AbstractCompiler, Requirement)
func newConformanceRules_RestrictedMethodCall(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_RestrictedMethodCall, error) {
	r := &ConformanceRules_RestrictedMethodCall{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	registry := compiler.GetTypeRegistry()
	restrictions := make([]ConformanceRules_RestrictedMethodCall_Restriction, 0, len(values))
	for _, value := range values {
		declName, err := conformanceRules_removeTypeDecl(value)
		if err != nil {
			return nil, err
		}
		typeName, property, hasClassAndProp := conformanceRules_splitClassAndProperty(declName)
		restrictedDecl, hasDecl := conformanceRules_getTypeFromValue(value)
		if !hasClassAndProp || !hasDecl {
			return nil, errInvalidRequirementSpec("bad prop value")
		}

		evaluated, err := conformanceRules_evaluateTypeString(compiler, restrictedDecl)
		if err != nil {
			return nil, err
		}
		restrictedCallType := evaluated.ToMaybeFunctionType()
		if restrictedCallType == nil {
			return nil, errInvalidRequirementSpec("invalid conformance type")
		}
		restrictions = append(restrictions, ConformanceRules_RestrictedMethodCall_Restriction{
			type_:              registry.GetGlobalType(typeName),
			property:           property,
			restrictedCallType: restrictedCallType,
		})
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *ConformanceRules_RestrictedMethodCall) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.RestrictedMethodCall.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_RestrictedMethodCall) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	gp, ok := ast.AsGetProp(n)
	if !ok || !ConformanceRules_ConformanceUtil_IsCallTarget(n) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	for i := 0; i < len(r.restrictions); i++ {
		res := &r.restrictions[i]
		result := ConformanceRules_ConformanceResult_CONFORMANCE

		if gp.Property() == res.property {
			result = r.checkConformanceWithRestriction(t, gp, res, false)
		} else if gp.Property() == "call" {
			if inner, ok := ast.AsGetProp(gp.Object().Node()); ok && inner.Property() == res.property {
				// handle .call invocation
				result = r.checkConformanceWithRestriction(t, gp, res, true)
			}
		}

		if result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
			return result
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.RestrictedMethodCall.checkConformance(NodeTraversal, Node, Restriction, boolean)
func (r *ConformanceRules_RestrictedMethodCall) checkConformanceWithRestriction(
	t *NodeTraversal, gp ast.GetProp, res *ConformanceRules_RestrictedMethodCall_Restriction, isCallInvocation bool) *ConformanceRules_ConformanceResult {
	registry := t.GetCompiler().GetTypeRegistry()
	methodClassType := res.type_
	var lhs rhino.Node
	if isCallInvocation {
		lhs = ast.MustGetProp(gp.Object().Node()).Object().Node()
	} else {
		lhs = gp.Object().Node()
	}
	if methodClassType != nil && lhs.GetJSType() != nil {
		targetType := lhs.GetJSType().RestrictByNotNullOrUndefined()
		n := gp.Node()
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

// ConformanceRules_RestrictedPropertyWrite is the restricted property write
// rule.
type ConformanceRules_RestrictedPropertyWrite struct {
	ConformanceRules_AbstractRule

	restrictions []ConformanceRules_RestrictedPropertyWrite_Restriction
}

// Java: ConformanceRules.RestrictedPropertyWrite.Restriction
type ConformanceRules_RestrictedPropertyWrite_Restriction struct {
	type_          rhino.JSType
	property       string
	restrictedType rhino.JSType
}

// Java: ConformanceRules.RestrictedPropertyWrite(AbstractCompiler, Requirement)
func newConformanceRules_RestrictedPropertyWrite(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_RestrictedPropertyWrite, error) {
	r := &ConformanceRules_RestrictedPropertyWrite{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	registry := compiler.GetTypeRegistry()
	restrictions := make([]ConformanceRules_RestrictedPropertyWrite_Restriction, 0, len(values))
	for _, value := range values {
		declName, err := conformanceRules_removeTypeDecl(value)
		if err != nil {
			return nil, err
		}
		typeName, property, hasClassAndProp := conformanceRules_splitClassAndProperty(declName)
		restrictedDecl, hasDecl := conformanceRules_getTypeFromValue(value)
		if !hasClassAndProp || !hasDecl {
			return nil, errInvalidRequirementSpec("bad prop value")
		}
		restrictedType, err := conformanceRules_evaluateTypeString(compiler, restrictedDecl)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, ConformanceRules_RestrictedPropertyWrite_Restriction{
			type_:          registry.GetGlobalType(typeName),
			property:       property,
			restrictedType: restrictedType,
		})
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *ConformanceRules_RestrictedPropertyWrite) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.RestrictedPropertyWrite.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_RestrictedPropertyWrite) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if gp, ok := ast.AsGetProp(n); ok {
		if asgn, isAssign := ast.AsAssign(n.GetParent()); isAssign && asgn.Target().Node() == n {
			rhsType := asgn.Value().Node().GetJSType()
			targetType := gp.Object().Node().GetJSType()
			if rhsType != nil && targetType != nil {
				prop := gp.Property()
				var targetNotNullType rhino.JSType
				for i := range r.restrictions {
					res := &r.restrictions[i]
					if prop == res.property && !rhsType.IsSubtypeOf(res.restrictedType) {
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
	ConformanceRules_AbstractRule

	restrictions []*TemplateAstMatcher
}

// Java: ConformanceRules.BannedCodePattern(AbstractCompiler, Requirement)
func newConformanceRules_BannedCodePattern(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedCodePattern, error) {
	r := &ConformanceRules_BannedCodePattern{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	restrictions := make([]*TemplateAstMatcher, 0, len(values))
	for _, value := range values {
		parseRoot :=
			NewCompilerInput(SourceFile_FromCode("<template>", value)).GetAstRoot(compiler)
		if !parseRoot.HasOneChild() || !parseRoot.GetFirstChild().IsFunction() {
			return nil, errInvalidRequirementSpec("invalid conformance template: " + value)
		}
		templateRoot := parseRoot.GetFirstChild()
		astMatcher := NewTemplateAstMatcher(compiler, templateRoot, r.typeMatchingStrategy)
		restrictions = append(restrictions, astMatcher)
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *ConformanceRules_BannedCodePattern) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
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

// Java: ConformanceRules.CustomRuleProxy(AbstractCompiler, Requirement)
func newConformanceRules_CustomRuleProxy(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_CustomRuleProxy, error) {
	if !requirement.HasJavaClass() {
		return nil, errInvalidRequirementSpec("missing java_class")
	}
	factory, ok := conformanceRules_customRuleRegistry[requirement.GetJavaClass()]
	if !ok {
		return nil, errInvalidRequirementSpec("Unknown conformance rule: " + requirement.GetJavaClass())
	}
	rule, err := factory(compiler, requirement)
	if err != nil {
		return nil, err
	}
	return &ConformanceRules_CustomRuleProxy{customRule: rule}, nil
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

// Java: ConformanceRules.CustomRuleProxy.CUSTOM_RULE_REGISTRY
var conformanceRules_customRuleRegistry = map[string]func(AbstractCompiler, *Requirement) (CheckConformance_Rule, error){
	"com.google.javascript.jscomp.ConformanceRules$NoOp": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_NoOp(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$InferredConstCheck": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_InferredConstCheck(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedDependency": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedDependency(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedDependencyRegex": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedDependencyRegex(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedName": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedName(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedProperty": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedProperty(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedStringRegex": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedStringRegex(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$RestrictedNameCall": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_RestrictedNameCall(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$RestrictedMethodCall": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_RestrictedMethodCall(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$RestrictedPropertyWrite": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_RestrictedPropertyWrite(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedCodePattern": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedCodePattern(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanForOf": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanForOf(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$RequireUseStrict": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_RequireUseStrict(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanThrowOfNonErrorTypes(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanNullDeref": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanNullDeref(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanUnknownThis": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanUnknownThis(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanUnknownDirectThisPropsReferences": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanUnknownDirectThisPropsReferences(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanUnknownTypedClassPropsReferences": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanUnknownTypedClassPropsReferences(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanUnresolvedType": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanUnresolvedType(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$StrictBanUnresolvedType": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_StrictBanUnresolvedType(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanGlobalVars": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanGlobalVars(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedEnhance": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedEnhance(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedModsRegex": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedModsRegex(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanCreateElement": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanCreateElement(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanCreateDom": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanCreateDom(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanElementSetAttribute(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanSettingAttributes(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanExecCommand": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanExecCommand(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanStaticThis": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BanStaticThis(c, r)
	},
}

// ConformanceRules_BanForOf bans for/of loops.
type ConformanceRules_BanForOf struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.BanForOf(AbstractCompiler, Requirement)
func newConformanceRules_BanForOf(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanForOf, error) {
	r := &ConformanceRules_BanForOf{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_BanForOf) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
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
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.RequireUseStrict(AbstractCompiler, Requirement)
func newConformanceRules_RequireUseStrict(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_RequireUseStrict, error) {
	r := &ConformanceRules_RequireUseStrict{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	if len(requirement.GetValueList()) != 0 {
		return nil, errInvalidRequirementSpec("invalid value")
	}
	return r, nil
}

func (r *ConformanceRules_RequireUseStrict) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
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
	ConformanceRules_AbstractRule

	errorObjType rhino.JSType
}

// Java: ConformanceRules.BanThrowOfNonErrorTypes(AbstractCompiler, Requirement)
func newConformanceRules_BanThrowOfNonErrorTypes(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanThrowOfNonErrorTypes, error) {
	r := &ConformanceRules_BanThrowOfNonErrorTypes{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	r.errorObjType = compiler.GetTypeRegistry().GetGlobalType("Error")
	return r, nil
}

func (r *ConformanceRules_BanThrowOfNonErrorTypes) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanThrowOfNonErrorTypes.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanThrowOfNonErrorTypes) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if r.errorObjType != nil {
		if th, ok := ast.AsThrow(n); ok {
			thrown := th.Value().Node().GetJSType()
			if thrown != nil &&
				!thrown.IsUnknownType() &&
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
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.BanNullDeref(AbstractCompiler, Requirement)
func newConformanceRules_BanNullDeref(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanNullDeref, error) {
	r := &ConformanceRules_BanNullDeref{}
	if err := r.initTypeRestrictionRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_BanNullDeref) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanNullDeref.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanNullDeref) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	var violation bool
	if get, ok := ast.AsNormalGet(n); ok {
		violation = r.ReportNode(get.Object().Node())
	} else if cn, ok := ast.AsCallOrNew(n); ok && !n.IsOptChainCall() {
		violation = r.ReportNode(cn.Callee().Node())
	} else if bin, ok := ast.AsBinary(n); ok && bin.Token() == rhino.Token_IN {
		violation = r.ReportNode(bin.Right().Node())
	}

	if violation {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanNullDeref.report(Node)
func (r *ConformanceRules_BanNullDeref) ReportNode(n rhino.Node) bool {
	return n.GetJSType() != nil && ConformanceRules_AbstractTypeRestrictionRule_IsKnown(n) && r.InvalidDeref(n) && !r.IsAllowlistedType(n)
}

// Java: ConformanceRules.BanNullDeref.invalidDeref(Node)
func (r *ConformanceRules_BanNullDeref) InvalidDeref(n rhino.Node) bool {
	type_ := n.GetJSType()
	return !type_.IsAllType() && (type_.IsNullable() || type_.IsVoidable())
}

// ConformanceRules_BanUnknownThis bans unknown "this" types.
type ConformanceRules_BanUnknownThis struct {
	ConformanceRules_AbstractRule

	reports map[rhino.Node]struct{}
}

// Java: ConformanceRules.BanUnknownThis(AbstractCompiler, Requirement)
func newConformanceRules_BanUnknownThis(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanUnknownThis, error) {
	r := &ConformanceRules_BanUnknownThis{
		reports: make(map[rhino.Node]struct{}),
	}
	if err := r.initTypeRestrictionRule(compiler, requirement, true); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_BanUnknownThis) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanUnknownThis.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanUnknownThis) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsThis() {
		type_ := n.GetJSType()
		if type_ != nil && type_.IsUnknownType() && !r.IsTypeImmediatelyTightened(n) {
			root := t.GetScopeRoot()
			if _, reported := r.reports[root]; !reported {
				r.reports[root] = struct{}{}
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanUnknownDirectThisPropsReferences bans unknown type
// references of the form "this.prop".
type ConformanceRules_BanUnknownDirectThisPropsReferences struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences(AbstractCompiler, Requirement)
func newConformanceRules_BanUnknownDirectThisPropsReferences(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanUnknownDirectThisPropsReferences, error) {
	r := &ConformanceRules_BanUnknownDirectThisPropsReferences{}
	if err := r.initTypeRestrictionRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_BanUnknownDirectThisPropsReferences) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanUnknownDirectThisPropsReferences) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if gp, ok := ast.AsGetProp(n); ok &&
		ConformanceRules_BanUnknownDirectThisPropsReferences_IsKnownThis(gp.Object().Node()) &&
		ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n) &&
		!ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(n) &&
		r.IsUsed(n) &&
		!r.IsTypeImmediatelyTightened(n) &&
		!conformanceRules_isExplicitlyUnknownProp(gp) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences.isKnownThis(Node)
func ConformanceRules_BanUnknownDirectThisPropsReferences_IsKnownThis(n rhino.Node) bool {
	return n.IsThis() && !ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n)
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences.isExplicitlyUnknown(Node)
func conformanceRules_isExplicitlyUnknownProp(gp ast.GetProp) bool {
	owner := rhino.ObjectType_Cast(gp.Object().Node().GetJSType())
	if owner == nil {
		return false
	}
	prop := owner.GetSlot(gp.Property())
	return prop != nil && !prop.IsTypeInferred()
}

// ConformanceRules_BanUnknownTypedClassPropsReferences bans unknown type
// references of the form "instance.prop".
type ConformanceRules_BanUnknownTypedClassPropsReferences struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.BanUnknownTypedClassPropsReferences(AbstractCompiler, Requirement)
func newConformanceRules_BanUnknownTypedClassPropsReferences(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanUnknownTypedClassPropsReferences, error) {
	r := &ConformanceRules_BanUnknownTypedClassPropsReferences{}
	if err := r.initTypeRestrictionRule(compiler, requirement, true); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_BanUnknownTypedClassPropsReferences) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanUnknownTypedClassPropsReferences.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanUnknownTypedClassPropsReferences) CheckConformance(t *NodeTraversal, getprop rhino.Node) *ConformanceRules_ConformanceResult {
	if gp, ok := ast.AsGetProp(getprop); ok &&
		ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(getprop) &&
		r.IsUsed(getprop) &&
		!r.IsTypeImmediatelyTightened(getprop) &&
		r.IsCheckablePropertySource(gp.Object().Node()) &&
		!ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(getprop) &&
		!r.isDeclaredUnknown(gp) {
		propName := gp.Property()
		typeName := gp.Object().Node().GetJSType().String()
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
func (r *ConformanceRules_BanUnknownTypedClassPropsReferences) isDeclaredUnknown(gp ast.GetProp) bool {
	targetType := gp.Object().Node().GetJSType().RestrictByNotNullOrUndefined().ToMaybeObjectType()
	if targetType == nil {
		return false
	}

	info := targetType.GetPropertyJSDocInfo(gp.Property())
	if info == nil || !info.HasType() {
		return false
	}

	typeExprNode := info.GetType().GetRoot()
	if typeExprNode.GetToken() == rhino.Token_QMARK && !typeExprNode.HasChildren() {
		return true
	} else if typeExprNode.GetToken() == rhino.Token_PIPE {
		for child := typeExprNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.GetToken() == rhino.Token_QMARK {
				return true
			}
		}
	}

	return false
}

// ConformanceRules_BanUnresolvedType bans accessing properties from objects
// that are unresolved forward-declared type names.
type ConformanceRules_BanUnresolvedType struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.BanUnresolvedType(AbstractCompiler, Requirement)
func newConformanceRules_BanUnresolvedType(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanUnresolvedType, error) {
	r := &ConformanceRules_BanUnresolvedType{}
	if err := r.initTypeRestrictionRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_BanUnresolvedType) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanUnresolvedType.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanUnresolvedType) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if gp, ok := ast.AsGetProp(n); ok {
		nonConformingPart, hasPart := conformanceRules_getNonConformingTypePart(gp.Object().Node().GetJSType())
		if hasPart && !r.IsTypeImmediatelyTightened(n) {
			return NewConformanceRules_ConformanceResultWithNote(
				ConformanceRules_ConformanceLevel_VIOLATION,
				"Reference to type '"+nonConformingPart+"' never resolved.")
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanUnresolvedType.getNonConformingPart(JSType)
func conformanceRules_getNonConformingTypePart(type_ rhino.JSType) (string, bool) {
	if type_ == nil {
		return "", false
	}
	if type_.IsUnionType() {
		var nonConformingParts []string
		for _, part := range type_.GetUnionMembers() {
			if nonConformingPart, ok := conformanceRules_getNonConformingTypePart(part); ok {
				nonConformingParts = append(nonConformingParts, nonConformingPart)
			}
		}
		if len(nonConformingParts) != 0 {
			return strings.Join(nonConformingParts, "|"), true
		}
	} else if type_.IsNoResolvedType() {
		if refName := type_.ToObjectType().GetReferenceName(); refName != nil {
			return *refName, true
		}
	}
	return "", false
}

// ConformanceRules_StrictBanUnresolvedType bans any use of unresolved
// forward-declared types.
type ConformanceRules_StrictBanUnresolvedType struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.StrictBanUnresolvedType(AbstractCompiler, Requirement)
func newConformanceRules_StrictBanUnresolvedType(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_StrictBanUnresolvedType, error) {
	r := &ConformanceRules_StrictBanUnresolvedType{}
	if err := r.initTypeRestrictionRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_StrictBanUnresolvedType) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.StrictBanUnresolvedType.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_StrictBanUnresolvedType) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if nonConformingPart, ok := conformanceRules_getNonConformingTypePart(n.GetJSType()); ok && !r.IsTypeImmediatelyTightened(n) {
		return NewConformanceRules_ConformanceResultWithNote(
			ConformanceRules_ConformanceLevel_VIOLATION,
			"Reference to type '"+nonConformingPart+"' never resolved.")
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanGlobalVars bans global var declarations.
type ConformanceRules_BanGlobalVars struct {
	ConformanceRules_AbstractRule

	allowlistedNames map[string]struct{}
}

// Java: ConformanceRules.BanGlobalVars(AbstractCompiler, Requirement)
func newConformanceRules_BanGlobalVars(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanGlobalVars, error) {
	r := &ConformanceRules_BanGlobalVars{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	allowlisted := make(map[string]struct{}, len(values))
	for _, v := range values {
		allowlisted[v] = struct{}{}
	}
	r.allowlistedNames = allowlisted
	return r, nil
}

func (r *ConformanceRules_BanGlobalVars) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
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

	if fn, ok := ast.AsFunction(n); ok {
		return r.IsAllowlistedName(fn.NameNode().Ident())
	}

	if NodeUtil_IsNameDeclaration(n) {
		allowlisted := true
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
	if _, ok := r.allowlistedNames[name]; ok {
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
	ConformanceRules_AbstractRule

	bannedEnhancedNamespaces map[string]struct{}
}

// Java: ConformanceRules.BannedEnhance(AbstractCompiler, Requirement)
func newConformanceRules_BannedEnhance(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedEnhance, error) {
	r := &ConformanceRules_BannedEnhance{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}
	banned := make(map[string]struct{}, len(values))
	for _, v := range values {
		banned[v] = struct{}{}
	}
	r.bannedEnhancedNamespaces = banned
	return r, nil
}

func (r *ConformanceRules_BannedEnhance) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BannedEnhance.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedEnhance) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	docInfo := n.GetJSDocInfo()

	if docInfo == nil || !docInfo.HasEnhance() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	enhance := *docInfo.GetEnhance()
	if _, ok := r.bannedEnhancedNamespaces[enhance]; ok {
		return NewConformanceRules_ConformanceResultWithNote(
			ConformanceRules_ConformanceLevel_VIOLATION, "The enhanced namespace \""+enhance+"\"")
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BannedEnhance.getPrecondition()
func (r *ConformanceRules_BannedEnhance) GetPrecondition() CheckConformance_Precondition {
	return checkConformance_PreconditionBannedEnhanceScript
}

// ConformanceRules_BannedModsRegex checks that file does not include an @mods
// annotation for a banned namespace regex.
type ConformanceRules_BannedModsRegex struct {
	ConformanceRules_AbstractRule

	bannedModsRegex *jutil.Pattern
}

// Java: ConformanceRules.BannedModsRegex(AbstractCompiler, Requirement)
func newConformanceRules_BannedModsRegex(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedModsRegex, error) {
	r := &ConformanceRules_BannedModsRegex{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}

	bannedModsRegexList := requirement.GetValueList()
	if len(bannedModsRegexList) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}
	pattern, err := conformanceRules_buildPattern(bannedModsRegexList)
	if err != nil {
		return nil, err
	}
	r.bannedModsRegex = pattern
	return r, nil
}

func (r *ConformanceRules_BannedModsRegex) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
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

// Java: ConformanceRules.BannedModsRegex.getPrecondition()
func (r *ConformanceRules_BannedModsRegex) GetPrecondition() CheckConformance_Precondition {
	return checkConformance_PreconditionBannedModsRegexScript
}

// Java: ConformanceRules.GOOG_DOM
var conformanceRules_GOOG_DOM = rhino.QualifiedName_Of("goog.dom")

// Java: ConformanceRules.GOOG_DOM_TAGNAME
var conformanceRules_GOOG_DOM_TAGNAME = rhino.QualifiedName_Of("goog.dom.TagName")

// ConformanceRules_BanCreateElement bans document.createElement and similar
// methods with string literal parameter specified in value, e.g. value:
// 'script'.
type ConformanceRules_BanCreateElement struct {
	ConformanceRules_AbstractRule

	bannedTags    map[string]struct{}
	domHelperType rhino.JSType
	documentType  rhino.JSType
}

// Java: ConformanceRules.BanCreateElement(AbstractCompiler, Requirement)
func newConformanceRules_BanCreateElement(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanCreateElement, error) {
	r := &ConformanceRules_BanCreateElement{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("Specify one or more values.")
	}
	bannedTags := make(map[string]struct{}, len(values))
	for _, value := range values {
		bannedTags[jutil.StringToLowerCaseASCII(value)] = struct{}{}
	}
	r.bannedTags = bannedTags
	r.domHelperType = compiler.GetTypeRegistry().GetGlobalType("goog.dom.DomHelper")
	r.documentType = compiler.GetTypeRegistry().GetGlobalType("Document")
	return r, nil
}

func (r *ConformanceRules_BanCreateElement) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanCreateElement.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanCreateElement) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsCall(n); ok {
		if firstArg, ok := call.FirstArg(); ok {
			if tag, ok := ast.AsStringLit(firstArg.Node()); ok {
				if _, banned := r.bannedTags[jutil.StringToLowerCaseASCII(tag.Value())]; banned {
					return r.checkCreateElement(call)
				}
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanCreateElement.checkCreateElement(Node)
func (r *ConformanceRules_BanCreateElement) checkCreateElement(call ast.Call) *ConformanceRules_ConformanceResult {
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	functionName := target.Property()
	if functionName != "createElement" && functionName != "createDom" {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	srcObj := target.Object().Node()
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
// format tagname.attribute, e.g. value: 'iframe.src'.
type ConformanceRules_BanCreateDom struct {
	ConformanceRules_AbstractRule

	bannedTagAttrs [][2]string
	domHelperType  rhino.JSType
	classNameTypes rhino.JSType
}

// Java: ConformanceRules.BanCreateDom(AbstractCompiler, Requirement)
func newConformanceRules_BanCreateDom(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanCreateDom, error) {
	r := &ConformanceRules_BanCreateDom{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	bannedTagAttrs := make([][2]string, 0, len(values))
	for _, value := range values {
		tagAttr := conformanceRules_splitToList(conformanceRules_ON_DOT, value)
		if len(tagAttr) != 2 || tagAttr[0] == "" || tagAttr[1] == "" {
			return nil, errInvalidRequirementSpec("Values must be in the format tagname.attribute.")
		}
		lowercasedTag := jutil.StringToLowerCaseRoot(tagAttr[0])
		bannedTagAttrs = append(bannedTagAttrs, [2]string{lowercasedTag, tagAttr[1]})
	}
	if len(bannedTagAttrs) == 0 {
		return nil, errInvalidRequirementSpec("Specify one or more values.")
	}
	r.bannedTagAttrs = bannedTagAttrs
	registry := compiler.GetTypeRegistry()
	r.domHelperType = registry.GetGlobalType("goog.dom.DomHelper")
	r.classNameTypes = registry.CreateUnionType(
		registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_NULL_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
	return r, nil
}

func (r *ConformanceRules_BanCreateDom) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanCreateDom.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanCreateDom) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	call, ok := r.asCreateDomCall(n)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	firstArg, hasFirst := call.FirstArg()
	if !hasFirst {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	attrsNode := firstArg.Node().GetNext()
	if attrsNode.IsNil() {
		// goog.dom.createDom('iframe') is fine.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	tagNames, tagNamesIsNull := r.GetTagNames(firstArg.Node())
	attrsType := attrsNode.GetJSType()
	if attrsType == nil {
		// Type information is not available; don't run this check.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// String or array attribute sets the class.
	isClassName := !attrsType.IsUnknownType() && attrsType.IsSubtypeOf(r.classNameTypes)

	if attrsNode.IsNull() || attrsType.IsVoidType() {
		// goog.dom.createDom('iframe', null) is fine.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	hasExtraChildren := !attrsNode.GetNext().IsNil()

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
			hasExtraChildren &&
			violation != ConformanceRules_ConformanceResult_CONFORMANCE {
			return violation
		}
		objLit, isObjLit := ast.AsObjectLit(attrsNode)
		if !isObjLit {
			// Attrs is not an object literal and tagName matches or is unknown.
			if r.reportLooseTypeViolations {
				return ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION
			}
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		prop := NodeUtil_GetFirstPropMatchingKey(attrsNode, tagAttr[1])
		if !prop.IsNil() {
			if NodeUtil_IsSomeCompileTimeConstStringValue(prop) {
				// Ignore string literal values.
				continue
			}
			return violation
		}

		for attr := objLit.Node().GetFirstChild(); !attr.IsNil(); attr = attr.GetNext() {
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

// Java: ConformanceRules.BanCreateDom.getTagNames(Node)
func (r *ConformanceRules_BanCreateDom) GetTagNames(tag rhino.Node) (tagNames []string, isNull bool) {
	if s, ok := ast.AsStringLit(tag); ok {
		return []string{jutil.StringToLowerCaseRoot(s.Value())}, false
	} else if gp, ok := ast.AsGetProp(tag); ok && conformanceRules_GOOG_DOM_TAGNAME.Matches(gp.Object().Node()) {
		return []string{jutil.StringToLowerCaseRoot(gp.Property())}, false
	}
	type_ := tag.GetJSType()
	if type_ == nil || !type_.IsTemplatizedType() {
		return nil, true
	}
	typeAsObj := type_.ToMaybeObjectType()
	if *typeAsObj.GetRawType().GetDisplayName() == "goog.dom.TagName" {
		tagType := conformanceRules_getOnlyElement(typeAsObj.GetTemplateTypes())
		var names []string
		if displayName := tagType.GetDisplayName(); displayName != nil {
			names = conformanceRules_elementTagNames[*displayName]
		}
		return names, false
	}
	return nil, true
}

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
var conformanceRules_elementTagNames = map[string][]string{
	"HTMLAnchorElement":       {"a"},
	"HTMLAppletElement":       {"applet"},
	"HTMLAreaElement":         {"area"},
	"HTMLAudioElement":        {"audio"},
	"HTMLBRElement":           {"br"},
	"HTMLBaseElement":         {"base"},
	"HTMLBaseFontElement":     {"basefont"},
	"HTMLBodyElement":         {"body"},
	"HTMLButtonElement":       {"button"},
	"HTMLCanvasElement":       {"canvas"},
	"HTMLDListElement":        {"dl"},
	"HTMLDataListElement":     {"datalist"},
	"HTMLDetailsElement":      {"details"},
	"HTMLDialogElement":       {"dialog"},
	"HTMLDirectoryElement":    {"dir"},
	"HTMLDivElement":          {"div"},
	"HTMLEmbedElement":        {"embed"},
	"HTMLFieldSetElement":     {"fieldset"},
	"HTMLFontElement":         {"font"},
	"HTMLFormElement":         {"form"},
	"HTMLFrameElement":        {"frame"},
	"HTMLFrameSetElement":     {"frameset"},
	"HTMLHRElement":           {"hr"},
	"HTMLHeadElement":         {"head"},
	"HTMLHeadingElement":      {"h1", "h2", "h3", "h4", "h5", "h6"},
	"HTMLHtmlElement":         {"html"},
	"HTMLIFrameElement":       {"iframe"},
	"HTMLImageElement":        {"img"},
	"HTMLInputElement":        {"input"},
	"HTMLIsIndexElement":      {"isindex"},
	"HTMLLIElement":           {"li"},
	"HTMLLabelElement":        {"label"},
	"HTMLLegendElement":       {"legend"},
	"HTMLLinkElement":         {"link"},
	"HTMLMapElement":          {"map"},
	"HTMLMenuElement":         {"menu"},
	"HTMLMetaElement":         {"meta"},
	"HTMLMeterElement":        {"meter"},
	"HTMLModElement":          {"del", "ins"},
	"HTMLOListElement":        {"ol"},
	"HTMLObjectElement":       {"object"},
	"HTMLOptGroupElement":     {"optgroup"},
	"HTMLOptionElement":       {"option"},
	"HTMLOutputElement":       {"output"},
	"HTMLParagraphElement":    {"p"},
	"HTMLParamElement":        {"param"},
	"HTMLPreElement":          {"pre"},
	"HTMLProgressElement":     {"progress"},
	"HTMLQuoteElement":        {"blockquote", "q"},
	"HTMLScriptElement":       {"script"},
	"HTMLSelectElement":       {"select"},
	"HTMLSourceElement":       {"source"},
	"HTMLSpanElement":         {"span"},
	"HTMLStyleElement":        {"style"},
	"HTMLTableCaptionElement": {"caption"},
	"HTMLTableCellElement":    {"td", "th"},
	"HTMLTableColElement":     {"col", "colgroup"},
	"HTMLTableElement":        {"table"},
	"HTMLTableRowElement":     {"tr"},
	"HTMLTableSectionElement": {"tbody", "tfoot", "thead"},
	"HTMLTemplateElement":     {"template"},
	"HTMLTextAreaElement":     {"textarea"},
	"HTMLTitleElement":        {"title"},
	"HTMLTrackElement":        {"track"},
	"HTMLUListElement":        {"ul"},
	"HTMLVideoElement":        {"video"},
}

// Java: ConformanceRules.BanCreateDom.isCreateDomCall(Node)
func (r *ConformanceRules_BanCreateDom) asCreateDomCall(n rhino.Node) (ast.Call, bool) {
	call, ok := ast.AsCall(n)
	if !ok {
		return ast.Call{}, false
	}
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok || target.Property() != "createDom" {
		return ast.Call{}, false
	}

	srcObj := target.Object().Node()
	if conformanceRules_GOOG_DOM.Matches(srcObj) {
		return call, true
	}
	type_ := srcObj.GetJSType()
	if type_ != nil && type_.Equals(r.domHelperType) {
		return call, true
	}
	return ast.Call{}, false
}

// ConformanceRules_SecuritySensitiveAttributes checks nodes for conformance
// with banning the setting of attributes that are on the blocklist.
type ConformanceRules_SecuritySensitiveAttributes struct {
	bannedAttrs             map[string]struct{}
	globalNamespaceSupplier func() *GlobalNamespace
}

// Java: ConformanceRules.SecuritySensitiveAttributes.ALL_BANNED_ATTRS
var conformanceRules_allBannedAttrs = [...]string{
	"href",
	"rel",
	"src",
	"srcdoc",
	"action",
	"formaction",
	"sandbox",
	"icon",
	"codebase",
	"data",
}

// Java: ConformanceRules.SecuritySensitiveAttributes()
func NewConformanceRules_SecuritySensitiveAttributes() *ConformanceRules_SecuritySensitiveAttributes {
	return NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(
		conformanceRules_allBannedAttrs[:], func() *GlobalNamespace { return nil })
}

// Java: ConformanceRules.SecuritySensitiveAttributes(Collection<String>, Supplier<GlobalNamespace>)
func NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(
	bannedAttrs []string, globalNamespaceSupplier func() *GlobalNamespace) *ConformanceRules_SecuritySensitiveAttributes {
	m := make(map[string]struct{}, len(bannedAttrs))
	for _, a := range bannedAttrs {
		m[jutil.StringToLowerCaseRoot(a)] = struct{}{}
	}
	return &ConformanceRules_SecuritySensitiveAttributes{
		bannedAttrs:             m,
		globalNamespaceSupplier: globalNamespaceSupplier,
	}
}

// Java: ConformanceRules.SecuritySensitiveAttributes.contains(String)
func (s *ConformanceRules_SecuritySensitiveAttributes) Contains(attributeName string) bool {
	_, ok := s.bannedAttrs[attributeName]
	return ok
}

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
	if s.Contains(lowerName) || ConformanceRules_ConformanceUtil_IsEventHandlerAttrName(lowerName) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

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
type ConformanceRules_BanElementSetAttribute struct {
	ConformanceRules_AbstractRule

	performGlobalNamespaceAnalysis bool
	globalNamespace                *GlobalNamespace

	elementType                      rhino.JSType
	securitySensitiveAttributes      *ConformanceRules_SecuritySensitiveAttributes
	defaultDecisionForUncertainCases *ConformanceRules_ConformanceResult
}

const (
	conformanceRules_BanElementSetAttribute_ELEMENT_TYPE_NAME               = "Element"
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE                   = "setAttribute"
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS                = "setAttributeNS"
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE              = "setAttributeNode"
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE_NS           = "setAttributeNodeNS"
	conformanceRules_BanElementSetAttribute_GLOBAL_NAMESPACE_ANALYSIS_LIMIT = 10000
)

func conformanceRules_isBannedSetAttributeProp(name string) bool {
	switch name {
	case conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE,
		conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS,
		conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE,
		conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE_NS:
		return true
	default:
		return false
	}
}

// Java: ConformanceRules.BanElementSetAttribute(AbstractCompiler, Requirement)
func newConformanceRules_BanElementSetAttribute(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanElementSetAttribute, error) {
	r := &ConformanceRules_BanElementSetAttribute{
		performGlobalNamespaceAnalysis: true,
	}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	r.securitySensitiveAttributes =
		NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(requirement.GetValueList(), r.GetGlobalNamespace)
	if requirement.GetReportLooseTypeViolations() {
		r.defaultDecisionForUncertainCases = ConformanceRules_ConformanceResult_VIOLATION
	} else {
		r.defaultDecisionForUncertainCases = ConformanceRules_ConformanceResult_CONFORMANCE
	}
	r.elementType = compiler.GetTypeRegistry().GetGlobalType(conformanceRules_BanElementSetAttribute_ELEMENT_TYPE_NAME)
	return r, nil
}

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

func (r *ConformanceRules_BanElementSetAttribute) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanElementSetAttribute.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsCall(node); ok {
		return r.checkConformanceOnPropertyCall(traversal, call)
	} else if ge, ok := ast.AsGetElem(node); ok && NodeUtil_IsLValue(node) {
		return r.checkConformanceOnGetElement(traversal, ge)
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanElementSetAttribute.checkConformanceOnPropertyCall(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) checkConformanceOnPropertyCall(
	traversal *NodeTraversal, call ast.Call) *ConformanceRules_ConformanceResult {
	calledProperty, ok := r.getBannedPropertyName(call)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if calledProperty == conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE {
		firstArg, hasFirst := call.FirstArg()
		if !hasFirst || firstArg.Node().GetNext().IsNil() {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}

		if r.requirement.GetReportLooseTypeViolations() {
			return r.securitySensitiveAttributes.CheckConformanceForAttributeName(
				traversal, firstArg.Node())
		}
		return r.securitySensitiveAttributes.CheckConformanceForAttributeNameWithHighConfidence(
			traversal, firstArg.Node())
	}

	if calledProperty == conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS {
		firstArg, hasFirst := call.FirstArg()
		if !hasFirst {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		secondArg := firstArg.Node().GetNext()
		if secondArg.IsNil() || secondArg.GetNext().IsNil() {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}

		if firstArg.Node().IsNull() {
			// A null namespace makes this equivalent to a regular setAttribute call.
			if r.requirement.GetReportLooseTypeViolations() {
				return r.securitySensitiveAttributes.CheckConformanceForAttributeName(
					traversal, secondArg)
			}
			return r.securitySensitiveAttributes.CheckConformanceForAttributeNameWithHighConfidence(
				traversal, secondArg)
		}
	}

	// We haven't defined a security contract on setAttributeNode and setAttributeNodeNS yet, so
	// flag them as risky if report_loose_type_violations is set.
	return r.defaultDecisionForUncertainCases
}

// Java: ConformanceRules.BanElementSetAttribute.getBannedPropertyName(Node)
func (r *ConformanceRules_BanElementSetAttribute) getBannedPropertyName(call ast.Call) (string, bool) {
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok {
		return "", false
	}
	propertyName := target.Property()
	if !conformanceRules_isBannedSetAttributeProp(propertyName) {
		return "", false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return "", false
	}
	if r.elementType != nil && type_.RestrictByNotNullOrUndefined().IsSubtypeOf(r.elementType) {
		return propertyName, true
	}
	return "", false
}

// Java: ConformanceRules.BanElementSetAttribute.checkConformanceOnGetElement(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) checkConformanceOnGetElement(
	traversal *NodeTraversal, ge ast.GetElem) *ConformanceRules_ConformanceResult {
	key := ge.Index().Node()
	keyName :=
		ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), key, r.GetGlobalNamespace)
	if keyName != nil {
		if r.securitySensitiveAttributes.Contains(jutil.StringToLowerCaseRoot(*keyName)) ||
			(r.requirement.GetReportLooseTypeViolations() &&
				(*keyName == "innerHTML" || *keyName == "outerHTML")) {
			if r.hasElementType(ge) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if !r.requirement.GetReportLooseTypeViolations() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	keyType := key.GetJSType()
	if keyType == nil || ConformanceRules_ConformanceUtil_IsXid(keyType) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if r.hasElementType(ge) {
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
func (r *ConformanceRules_BanElementSetAttribute) hasElementType(ge ast.GetElem) bool {
	objType := ge.Object().Node().GetJSType()
	if objType == nil || r.elementType == nil {
		return false
	}
	objType = objType.RestrictByNotNullOrUndefined()
	if objType.IsUnknownType() || objType.IsEmptyType() || !objType.IsSubtypeOf(r.elementType) {
		return false
	}
	return true
}

// ConformanceRules_BanSettingAttributes can be used to ban any function like
// Element#setAttribute which takes two parameters, an attribute name and a
// value.
type ConformanceRules_BanSettingAttributes struct {
	ConformanceRules_AbstractRule

	restrictions []ConformanceRules_BanSettingAttributes_AttributeSettingRestriction
}

// Java: ConformanceRules.BanSettingAttributes.AttributeSettingRestriction
type ConformanceRules_BanSettingAttributes_AttributeSettingRestriction struct {
	type_    rhino.JSType
	property string
}

// Java: ConformanceRules.BanSettingAttributes.BANNED_ATTRS
var conformanceRules_BanSettingAttributes_BANNED_ATTRS = NewConformanceRules_SecuritySensitiveAttributes()

// Java: ConformanceRules.BanSettingAttributes(AbstractCompiler, Requirement)
func newConformanceRules_BanSettingAttributes(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanSettingAttributes, error) {
	r := &ConformanceRules_BanSettingAttributes{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	registry := compiler.GetTypeRegistry()
	restrictions := make([]ConformanceRules_BanSettingAttributes_AttributeSettingRestriction, 0, len(values))
	for _, value := range values {
		typeName, property, ok := conformanceRules_splitClassAndProperty(value)
		if !ok {
			return nil, errInvalidRequirementSpec("bad prop value")
		}
		restrictions = append(restrictions, ConformanceRules_BanSettingAttributes_AttributeSettingRestriction{
			type_:    registry.GetGlobalType(typeName),
			property: property,
		})
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *ConformanceRules_BanSettingAttributes) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanSettingAttributes.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanSettingAttributes) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsCall(node); ok {
		return r.checkConformanceOnPropertyCall(traversal, call)
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanSettingAttributes.checkConformanceOnPropertyCall(NodeTraversal, Node)
func (r *ConformanceRules_BanSettingAttributes) checkConformanceOnPropertyCall(traversal *NodeTraversal, call ast.Call) *ConformanceRules_ConformanceResult {
	if !r.matchesBannedProperty(call) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	firstArg, hasFirst := call.FirstArg()
	if !hasFirst || firstArg.Node().GetNext().IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	return conformanceRules_BanSettingAttributes_BANNED_ATTRS.CheckConformanceForAttributeName(traversal, firstArg.Node())
}

// Java: ConformanceRules.BanSettingAttributes.getBannedPropertyName(Node)
func (r *ConformanceRules_BanSettingAttributes) matchesBannedProperty(call ast.Call) bool {
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok {
		return false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return false
	}
	propertyName := target.Property()
	for i := range r.restrictions {
		restricted := &r.restrictions[i]
		if propertyName == restricted.property &&
			restricted.type_ != nil &&
			type_.IsSubtypeOf(restricted.type_) {
			return true
		}
	}
	return false
}

// ConformanceRules_BanExecCommand bans Document#execCommand with command
// names specified in value or any dynamic string.
type ConformanceRules_BanExecCommand struct {
	ConformanceRules_AbstractRule

	bannedAttrs []string
}

// Java: ConformanceRules.BanExecCommand(AbstractCompiler, Requirement)
func newConformanceRules_BanExecCommand(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanExecCommand, error) {
	r := &ConformanceRules_BanExecCommand{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	bannedAttrs := make([]string, 0, len(values))
	for _, v := range values {
		bannedAttrs = append(bannedAttrs, jutil.StringToLowerCaseRoot(v))
	}
	r.bannedAttrs = bannedAttrs
	return r, nil
}

func (r *ConformanceRules_BanExecCommand) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanExecCommand.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanExecCommand) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	call, ok := r.asBannedExecCommandCall(node)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	firstArg, hasFirst := call.FirstArg()
	if !hasFirst {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	attr, ok := ast.AsStringLit(firstArg.Node())
	if !ok {
		return ConformanceRules_ConformanceResult_VIOLATION
	}

	if slices.Contains(r.bannedAttrs, jutil.StringToLowerCaseRoot(attr.Value())) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanExecCommand.isBannedProperty(Node)
func (r *ConformanceRules_BanExecCommand) asBannedExecCommandCall(node rhino.Node) (ast.Call, bool) {
	call, ok := ast.AsCall(node)
	if !ok {
		return ast.Call{}, false
	}
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok || target.Property() != "execCommand" {
		return ast.Call{}, false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return ast.Call{}, false
	}
	documentType := r.compiler.GetTypeRegistry().GetGlobalType("Document")
	if documentType == nil || !type_.RestrictByNotNullOrUndefined().IsSubtypeOf(documentType) {
		return ast.Call{}, false
	}
	return call, true
}

// ConformanceRules_BanStaticThis checks that this is not being referenced
// directly within a static member function.
type ConformanceRules_BanStaticThis struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.BanStaticThis(AbstractCompiler, Requirement)
func newConformanceRules_BanStaticThis(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanStaticThis, error) {
	r := &ConformanceRules_BanStaticThis{}
	if err := r.initTypeRestrictionRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_BanStaticThis) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
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
	mfd, ok := ast.AsMemberFunctionDef(n.GetParent())
	return ok && mfd.IsStatic()
}
