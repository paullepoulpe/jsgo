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
	return n.GetJSType() != nil && conformanceRules_AbstractTypeRestrictionRule_isKnown(n) && r.InvalidDeref(n) && !r.IsAllowlistedType(n)
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
		conformanceRules_AbstractTypeRestrictionRule_isUnknown(n) &&
		!conformanceRules_AbstractTypeRestrictionRule_isTemplateType(n) &&
		r.IsUsed(n) &&
		!r.IsTypeImmediatelyTightened(n) &&
		!conformanceRules_isExplicitlyUnknownProp(gp) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanUnknownDirectThisPropsReferences.isKnownThis(Node)
func ConformanceRules_BanUnknownDirectThisPropsReferences_IsKnownThis(n rhino.Node) bool {
	return n.IsThis() && !conformanceRules_AbstractTypeRestrictionRule_isUnknown(n)
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
		conformanceRules_AbstractTypeRestrictionRule_isUnknown(getprop) &&
		r.IsUsed(getprop) &&
		!r.IsTypeImmediatelyTightened(getprop) &&
		r.IsCheckablePropertySource(gp.Object().Node()) &&
		!conformanceRules_AbstractTypeRestrictionRule_isTemplateType(getprop) &&
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
	return conformanceRules_AbstractTypeRestrictionRule_isKnown(n) &&
		!conformanceRules_AbstractTypeRestrictionRule_isTop(n) &&
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
