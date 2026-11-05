// Port of com/google/javascript/jscomp/ConformanceRules.java @ closure-compiler v20261028

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
	"sync"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type restrictedNameCallSpec struct {
	name               rhino.Node
	restrictedCallType rhino.FunctionType
}

func newConformanceRules_RestrictedNameCall(compiler AbstractCompiler, requirement *Requirement) (*conformanceRestrictedNameCallRule, error) {
	r := &conformanceRestrictedNameCallRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	restrictions := make([]restrictedNameCallSpec, 0, len(values))
	for _, value := range values {
		nameFromValue, restrictedDecl, ok := splitConformanceNameAndType(value)
		if !ok {
			return nil, errInvalidRequirementSpec("bad prop value")
		}
		name := NodeUtil_NewQName(compiler, nameFromValue)
		if name.IsNil() {
			return nil, errInvalidRequirementSpec("bad prop value")
		}

		evaluated, err := evaluateConformanceTypeString(compiler, restrictedDecl)
		if err != nil {
			return nil, err
		}
		restrictedCallType := evaluated.ToMaybeFunctionType()
		if restrictedCallType == nil {
			return nil, errInvalidRequirementSpec("invalid conformance type")
		}
		restrictions = append(restrictions, restrictedNameCallSpec{
			name:               name,
			restrictedCallType: restrictedCallType,
		})
	}
	r.restrictions = restrictions
	return r, nil
}

func (r *conformanceRestrictedNameCallRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if !n.IsNil() && isConformanceCallTarget(n) {
		for i := range r.restrictions {
			res := &r.restrictions[i]

			if n.IsQualifiedName() && n.MatchesQualifiedNameNode(res.name) {
				if !validateConformanceCall(r.compiler, n.GetParent(), res.restrictedCallType, false) {
					return ConformanceRules_ConformanceResult_VIOLATION
				}
			} else if gp, ok := ast.AsGetProp(n); ok &&
				gp.Property() == "call" &&
				gp.Object().Node().MatchesQualifiedNameNode(res.name) {
				if !validateConformanceCall(r.compiler, n.GetParent(), res.restrictedCallType, true) {
					return ConformanceRules_ConformanceResult_VIOLATION
				}
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

type conformanceRestrictedMethodCallRule struct {
	conformanceAbstractRule

	restrictions []restrictedMethodCallSpec
}

type restrictedMethodCallSpec struct {
	type_              rhino.JSType
	property           string
	restrictedCallType rhino.FunctionType
}

func newConformanceRules_RestrictedMethodCall(compiler AbstractCompiler, requirement *Requirement) (*conformanceRestrictedMethodCallRule, error) {
	r := &conformanceRestrictedMethodCallRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	registry := compiler.GetTypeRegistry()
	restrictions := make([]restrictedMethodCallSpec, 0, len(values))
	for _, value := range values {
		declName, restrictedDecl, ok := splitConformanceNameAndType(value)
		if !ok {
			return nil, errInvalidRequirementSpec("value should be in the form NAME:TYPE")
		}
		typeName, property, hasClassAndProp := splitConformanceClassAndProperty(declName)
		if !hasClassAndProp {
			return nil, errInvalidRequirementSpec("bad prop value")
		}

		evaluated, err := evaluateConformanceTypeString(compiler, restrictedDecl)
		if err != nil {
			return nil, err
		}
		restrictedCallType := evaluated.ToMaybeFunctionType()
		if restrictedCallType == nil {
			return nil, errInvalidRequirementSpec("invalid conformance type")
		}
		restrictions = append(restrictions, restrictedMethodCallSpec{
			type_:              registry.GetGlobalType(typeName),
			property:           property,
			restrictedCallType: restrictedCallType,
		})
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *conformanceRestrictedMethodCallRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	gp, ok := ast.AsGetProp(n)
	if !ok || !isConformanceCallTarget(n) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	for i := range r.restrictions {
		res := &r.restrictions[i]
		result := ConformanceRules_ConformanceResult_CONFORMANCE

		if gp.Property() == res.property {
			result = r.checkMethodCall(t, gp, gp.Object().Node(), res, false)
		} else if gp.Property() == "call" {
			if inner, ok := gp.Object().AsGetProp(); ok && inner.Property() == res.property {
				result = r.checkMethodCall(t, gp, inner.Object().Node(), res, true)
			}
		}

		if result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
			return result
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceRestrictedMethodCallRule) checkMethodCall(
	t *NodeTraversal,
	gp ast.GetProp,
	lhs rhino.Node,
	res *restrictedMethodCallSpec,
	isCallInvocation bool,
) *ConformanceRules_ConformanceResult {
	registry := t.GetCompiler().GetTypeRegistry()
	methodClassType := res.type_
	if methodClassType != nil && !lhs.IsNil() && lhs.GetJSType() != nil {
		targetType := lhs.GetJSType().RestrictByNotNullOrUndefined()
		n := gp.Node()
		if isLooseConformanceType(targetType) || isLooseObjectType(targetType, registry) {
			if r.reportLooseTypeViolations &&
				!validateConformanceCall(r.compiler, n.GetParent(), res.restrictedCallType, isCallInvocation) {
				return conformancePossibleViolationLooseTypes
			}
		} else if targetType.IsSubtypeOf(methodClassType) {
			if !validateConformanceCall(r.compiler, n.GetParent(), res.restrictedCallType, isCallInvocation) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

type conformanceRestrictedPropertyWriteRule struct {
	conformanceAbstractRule

	restrictions []restrictedPropertyWriteSpec
}

type restrictedPropertyWriteSpec struct {
	type_          rhino.JSType
	property       string
	restrictedType rhino.JSType
}

func newConformanceRules_RestrictedPropertyWrite(compiler AbstractCompiler, requirement *Requirement) (*conformanceRestrictedPropertyWriteRule, error) {
	r := &conformanceRestrictedPropertyWriteRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	registry := compiler.GetTypeRegistry()
	restrictions := make([]restrictedPropertyWriteSpec, 0, len(values))
	for _, value := range values {
		declName, restrictedDecl, ok := splitConformanceNameAndType(value)
		if !ok {
			return nil, errInvalidRequirementSpec("value should be in the form NAME:TYPE")
		}
		typeName, property, hasClassAndProp := splitConformanceClassAndProperty(declName)
		if !hasClassAndProp {
			return nil, errInvalidRequirementSpec("bad prop value")
		}
		restrictedType, err := evaluateConformanceTypeString(compiler, restrictedDecl)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, restrictedPropertyWriteSpec{
			type_:          registry.GetGlobalType(typeName),
			property:       property,
			restrictedType: restrictedType,
		})
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *conformanceRestrictedPropertyWriteRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if gp, ok := ast.AsGetProp(n); ok {
		if asgn, isAssign := ast.AsAssign(n.GetParent()); isAssign && asgn.Target().Is(n) {
			rhsType := asgn.Value().Node().GetJSType()
			targetType := gp.Object().Node().GetJSType()
			if rhsType != nil && targetType != nil {
				prop := gp.Property()
				var targetNotNullType rhino.JSType
				for i := range r.restrictions {
					res := &r.restrictions[i]
					if prop == res.property && res.restrictedType != nil && !rhsType.IsSubtypeOf(res.restrictedType) {
						if isLooseConformanceType(targetType) {
							if r.reportLooseTypeViolations {
								return conformancePossibleViolationLooseTypes
							}
						} else if res.type_ != nil {
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

type conformanceBannedCodePatternRule struct {
	conformanceAbstractRule

	restrictions []*TemplateAstMatcher
}

func newConformanceRules_BannedCodePattern(compiler AbstractCompiler, requirement *Requirement) (*conformanceBannedCodePatternRule, error) {
	r := &conformanceBannedCodePatternRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
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
		script, isScript := ast.AsScript(parseRoot)
		if !isScript {
			return nil, errInvalidRequirementSpec("invalid conformance template: " + value)
		}
		onlyItem, hasOne := script.OnlyItem()
		if !hasOne || !onlyItem.Node().IsFunction() {
			return nil, errInvalidRequirementSpec("invalid conformance template: " + value)
		}
		astMatcher := NewTemplateAstMatcher(compiler, onlyItem.Node(), r.typeMatchingStrategy)
		restrictions = append(restrictions, astMatcher)
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *conformanceBannedCodePatternRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	possibleViolation := false
	for _, matcher := range r.restrictions {
		if matcher.Matches(n) {
			if matcher.IsLooseMatch() {
				possibleViolation = true
			} else {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	if possibleViolation && r.reportLooseTypeViolations {
		return conformancePossibleViolationLooseTypes
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

var conformanceCustomRuleMu sync.RWMutex

// RegisterCustomConformanceRule registers a custom CheckConformance_Rule
// factory under javaClass (matching Requirement.java_class in ConformanceConfig).
func RegisterCustomConformanceRule(
	javaClass string,
	factory func(AbstractCompiler, *Requirement) (CheckConformance_Rule, error),
) {
	jutil.CheckNotNull(factory)
	conformanceCustomRuleMu.Lock()
	defer conformanceCustomRuleMu.Unlock()
	conformanceCustomRuleRegistry[javaClass] = factory
}

// UnregisterCustomConformanceRule removes a custom conformance rule factory
// previously registered under javaClass.
func UnregisterCustomConformanceRule(javaClass string) {
	conformanceCustomRuleMu.Lock()
	defer conformanceCustomRuleMu.Unlock()
	delete(conformanceCustomRuleRegistry, javaClass)
}

func newConformanceRules_CustomRule(compiler AbstractCompiler, requirement *Requirement) (CheckConformance_Rule, error) {
	if !requirement.HasJavaClass() {
		return nil, errInvalidRequirementSpec("missing java_class")
	}
	javaClass := requirement.GetJavaClass()
	conformanceCustomRuleMu.RLock()
	factory, ok := conformanceCustomRuleRegistry[javaClass]
	conformanceCustomRuleMu.RUnlock()
	if !ok {
		return nil, errInvalidRequirementSpec("Unknown conformance rule: " + javaClass)
	}
	return factory(compiler, requirement)
}

var conformanceCustomRuleRegistry = map[string]func(AbstractCompiler, *Requirement) (CheckConformance_Rule, error){
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
		return newConformanceBanForOfRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$RequireUseStrict": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRequireUseStrictRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanThrowOfNonErrorTypes": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanThrowOfNonErrorTypesRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanNullDeref": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanNullDerefRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanUnknownThis": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanUnknownThisRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanUnknownDirectThisPropsReferences": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanUnknownDirectThisPropsRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanUnknownTypedClassPropsReferences": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanUnknownTypedClassPropsRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanUnresolvedType": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanUnresolvedTypeRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$StrictBanUnresolvedType": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceStrictBanUnresolvedTypeRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanGlobalVars": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanGlobalVarsRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedEnhance": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedEnhance(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BannedModsRegex": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceRules_BannedModsRegex(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanCreateElement": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanCreateElementRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanCreateDom": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanCreateDomRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanElementSetAttribute": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanElementSetAttributeRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanSettingAttributes": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanSettingAttributesRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanExecCommand": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanExecCommandRule(c, r)
	},
	"com.google.javascript.jscomp.ConformanceRules$BanStaticThis": func(c AbstractCompiler, r *Requirement) (CheckConformance_Rule, error) {
		return newConformanceBanStaticThisRule(c, r)
	},
}

type conformanceBanForOfRule struct {
	conformanceAbstractRule
}

func newConformanceBanForOfRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanForOfRule, error) {
	r := &conformanceBanForOfRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceBanForOfRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if _, ok := ast.AsForOfLoop(n); ok {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

type conformanceRequireUseStrictRule struct {
	conformanceAbstractRule
}

func newConformanceRequireUseStrictRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceRequireUseStrictRule, error) {
	r := &conformanceRequireUseStrictRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	if len(requirement.GetValueList()) != 0 {
		return nil, errInvalidRequirementSpec("invalid value")
	}
	return r, nil
}

func (r *conformanceRequireUseStrictRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if _, ok := ast.AsScript(n); ok && !n.IsUseStrict() {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

type conformanceBanThrowOfNonErrorTypesRule struct {
	conformanceAbstractRule

	errorObjType rhino.JSType
}

func newConformanceBanThrowOfNonErrorTypesRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanThrowOfNonErrorTypesRule, error) {
	r := &conformanceBanThrowOfNonErrorTypesRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	r.errorObjType = compiler.GetTypeRegistry().GetGlobalType("Error")
	return r, nil
}

func (r *conformanceBanThrowOfNonErrorTypesRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
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

type conformanceBanNullDerefRule struct {
	conformanceAbstractRule
}

func newConformanceBanNullDerefRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanNullDerefRule, error) {
	r := &conformanceBanNullDerefRule{}
	if err := r.initTypeRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceBanNullDerefRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	var violation bool
	if get, ok := ast.AsNormalGet(n); ok {
		violation = r.isNullishDeref(get.Object().Node())
	} else if cn, ok := ast.AsCallOrNew(n); ok && !n.IsOptChainCall() {
		violation = r.isNullishDeref(cn.Callee().Node())
	} else if bin, ok := ast.AsBinary(n); ok && bin.Token() == rhino.Token_IN {
		violation = r.isNullishDeref(bin.Right().Node())
	}

	if violation {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBanNullDerefRule) isNullishDeref(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	type_ := n.GetJSType()
	return type_ != nil &&
		isKnownNodeType(n) &&
		!type_.IsAllType() &&
		(type_.IsNullable() || type_.IsVoidable()) &&
		!r.isAllowlistedType(n)
}

type conformanceBanUnknownThisRule struct {
	conformanceAbstractRule

	reports          map[rhino.Node]struct{}
	onGlobalScopeHit func(n rhino.Node)
}

func newConformanceBanUnknownThisRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanUnknownThisRule, error) {
	r := &conformanceBanUnknownThisRule{
		reports: make(map[rhino.Node]struct{}),
	}
	if err := r.initTypeRule(compiler, requirement, true, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceBanUnknownThisRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if !n.IsNil() && n.IsThis() {
		type_ := n.GetJSType()
		if type_ != nil && type_.IsUnknownType() && !r.isTypeImmediatelyTightened(n) {
			var root rhino.Node
			if t != nil {
				root = t.GetScopeRoot()
			}
			if _, reported := r.reports[root]; !reported {
				r.reports[root] = struct{}{}
				if r.onGlobalScopeHit != nil && root.IsRoot() {
					r.onGlobalScopeHit(n)
				}
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// conformanceBanUnknownDirectThisPropsRule bans unknown-typed `this.prop`
// references on a known `this` receiver.
type conformanceBanUnknownDirectThisPropsRule struct {
	conformanceAbstractRule
}

func newConformanceBanUnknownDirectThisPropsRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanUnknownDirectThisPropsRule, error) {
	r := &conformanceBanUnknownDirectThisPropsRule{}
	if err := r.initTypeRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceBanUnknownDirectThisPropsRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if gp, ok := ast.AsGetProp(n); ok &&
		isKnownThisNode(gp.Object().Node()) &&
		isUnknownNodeType(n) &&
		!isTemplateNodeType(n) &&
		r.isNodeUsed(n) &&
		!r.isTypeImmediatelyTightened(n) &&
		!isExplicitlyUnknownProp(gp) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func isKnownThisNode(n rhino.Node) bool {
	return !n.IsNil() && n.IsThis() && !isUnknownNodeType(n)
}

func isExplicitlyUnknownProp(gp ast.GetProp) bool {
	owner := rhino.ObjectType_Cast(gp.Object().Node().GetJSType())
	if owner == nil {
		return false
	}
	prop := owner.GetSlot(gp.Property())
	return prop != nil && !prop.IsTypeInferred()
}

// conformanceBanUnknownTypedClassPropsRule bans unknown-typed `instance.prop`
// references on class/interface instances.
type conformanceBanUnknownTypedClassPropsRule struct {
	conformanceAbstractRule
}

func newConformanceBanUnknownTypedClassPropsRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanUnknownTypedClassPropsRule, error) {
	r := &conformanceBanUnknownTypedClassPropsRule{}
	if err := r.initTypeRule(compiler, requirement, true, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceBanUnknownTypedClassPropsRule) checkNode(t *NodeTraversal, getprop rhino.Node) *ConformanceRules_ConformanceResult {
	if gp, ok := ast.AsGetProp(getprop); ok &&
		isUnknownNodeType(getprop) &&
		r.isNodeUsed(getprop) &&
		!r.isTypeImmediatelyTightened(getprop) &&
		r.isCheckablePropertySource(gp.Object().Node()) &&
		!isTemplateNodeType(getprop) &&
		!r.isDeclaredUnknown(gp) {
		propName := gp.Property()
		typeName := gp.Object().Node().GetJSType().String()
		return NewConformanceRules_ConformanceResultWithNote(
			ConformanceRules_ConformanceLevel_VIOLATION,
			"The property \""+propName+"\" on type \""+typeName+"\"")
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBanUnknownTypedClassPropsRule) isCheckablePropertySource(n rhino.Node) bool {
	return isKnownNodeType(n) &&
		!isTopNodeType(n) &&
		r.isClassType(n) &&
		!r.isNativeObjectType(n) &&
		!r.isAllowlistedType(n)
}

func (r *conformanceBanUnknownTypedClassPropsRule) isClassType(n rhino.Node) bool {
	if n.IsNil() || n.GetJSType() == nil {
		return false
	}
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

func (r *conformanceBanUnknownTypedClassPropsRule) isDeclaredUnknown(gp ast.GetProp) bool {
	objType := gp.Object().Node().GetJSType()
	if objType == nil {
		return false
	}
	targetType := objType.RestrictByNotNullOrUndefined().ToMaybeObjectType()
	if targetType == nil {
		return false
	}

	info := targetType.GetPropertyJSDocInfo(gp.Property())
	if info == nil || !info.HasType() {
		return false
	}

	typeExprNode := info.GetType().GetRoot()
	if typeExprNode.IsNil() {
		return false
	}
	if typeExprNode.GetToken() == rhino.Token_QMARK && !typeExprNode.HasChildren() {
		return true
	} else if typeExprNode.GetToken() == rhino.Token_PIPE {
		// JSDoc union type-AST node (`|`): check whether any branch is explicit `?`.
		for child := range typeExprNode.Children() {
			if child.GetToken() == rhino.Token_QMARK {
				return true
			}
		}
	}

	return false
}

// conformanceBanUnresolvedTypeRule bans accessing properties on unresolved
// forward-declared type names.
type conformanceBanUnresolvedTypeRule struct {
	conformanceAbstractRule
}

func newConformanceBanUnresolvedTypeRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanUnresolvedTypeRule, error) {
	r := &conformanceBanUnresolvedTypeRule{}
	if err := r.initTypeRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceBanUnresolvedTypeRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if gp, ok := ast.AsGetProp(n); ok {
		nonConformingPart, hasPart := nonConformingTypePart(gp.Object().Node().GetJSType())
		if hasPart && !r.isTypeImmediatelyTightened(n) {
			return NewConformanceRules_ConformanceResultWithNote(
				ConformanceRules_ConformanceLevel_VIOLATION,
				"Reference to type '"+nonConformingPart+"' never resolved.")
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func nonConformingTypePart(type_ rhino.JSType) (string, bool) {
	if type_ == nil {
		return "", false
	}
	if type_.IsUnionType() {
		var nonConformingParts []string
		for _, part := range type_.GetUnionMembers() {
			if nonConformingPart, ok := nonConformingTypePart(part); ok {
				nonConformingParts = append(nonConformingParts, nonConformingPart)
			}
		}
		if len(nonConformingParts) != 0 {
			return strings.Join(nonConformingParts, "|"), true
		}
	} else if type_.IsNoResolvedType() {
		if objType := type_.ToObjectType(); objType != nil {
			return objType.ReferenceName()
		}
	}
	return "", false
}

type conformanceStrictBanUnresolvedTypeRule struct {
	conformanceAbstractRule
}

func newConformanceStrictBanUnresolvedTypeRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceStrictBanUnresolvedTypeRule, error) {
	r := &conformanceStrictBanUnresolvedTypeRule{}
	if err := r.initTypeRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceStrictBanUnresolvedTypeRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if !n.IsNil() {
		if nonConformingPart, ok := nonConformingTypePart(n.GetJSType()); ok && !r.isTypeImmediatelyTightened(n) {
			return NewConformanceRules_ConformanceResultWithNote(
				ConformanceRules_ConformanceLevel_VIOLATION,
				"Reference to type '"+nonConformingPart+"' never resolved.")
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

type conformanceBanGlobalVarsRule struct {
	conformanceAbstractRule

	allowlistedNames map[string]struct{}
}

func newConformanceBanGlobalVarsRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanGlobalVarsRule, error) {
	r := &conformanceBanGlobalVarsRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
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

func (r *conformanceBanGlobalVarsRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if !n.IsNil() &&
		t != nil &&
		t.InGlobalScope() &&
		NodeUtil_IsDeclaration(n) &&
		!n.GetBooleanProp(rhino.Node_IS_NAMESPACE) &&
		!r.isAllowlisted(n) {
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

func (r *conformanceBanGlobalVarsRule) isAllowlisted(n rhino.Node) bool {
	if n.IsFromExterns() {
		return true
	}

	if fn, ok := ast.AsFunction(n); ok {
		name, hasName := fn.Name()
		return hasName && r.isAllowlistedName(name.Ident())
	}

	if NodeUtil_IsNameDeclaration(n) {
		allowlisted := true
		NodeUtil_VisitLhsNodesInNode(
			n,
			func(name rhino.Node) {
				if ident, ok := ast.AsName(name); !ok || !r.isAllowlistedName(ident.Ident()) {
					allowlisted = false
				}
			})
		return allowlisted
	}

	return false
}

func (r *conformanceBanGlobalVarsRule) isAllowlistedName(name string) bool {
	if _, ok := r.allowlistedNames[name]; ok {
		return true
	}
	// Names created by JSCompiler internals.
	return name == "$jscomp" ||
		strings.HasPrefix(name, "$jscomp$compprop") ||
		ClosureRewriteModule_IsModuleContent(name) ||
		ClosureRewriteModule_IsModuleExport(name) ||
		ScopedAliases_IsScopedAliases(name)
}

type conformanceBannedEnhanceRule struct {
	conformanceAbstractRule

	bannedEnhancedNamespaces map[string]struct{}
}

func newConformanceRules_BannedEnhance(compiler AbstractCompiler, requirement *Requirement) (*conformanceBannedEnhanceRule, error) {
	r := &conformanceBannedEnhanceRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
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

func (r *conformanceBannedEnhanceRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	docInfo := n.GetJSDocInfo()
	if docInfo == nil {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	enhance, hasEnhance := docInfo.Enhance()
	if !hasEnhance {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if _, ok := r.bannedEnhancedNamespaces[enhance]; ok {
		return NewConformanceRules_ConformanceResultWithNote(
			ConformanceRules_ConformanceLevel_VIOLATION, "The enhanced namespace \""+enhance+"\"")
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBannedEnhanceRule) GetPrecondition() CheckConformance_Precondition {
	return checkConformance_PreconditionBannedEnhanceScript
}

type conformanceBannedModsRegexRule struct {
	conformanceAbstractRule

	bannedModsRegex *jutil.Pattern
}

func newConformanceRules_BannedModsRegex(compiler AbstractCompiler, requirement *Requirement) (*conformanceBannedModsRegexRule, error) {
	r := &conformanceBannedModsRegexRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}

	bannedModsRegexList := requirement.GetValueList()
	if len(bannedModsRegexList) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}
	pattern, err := buildConformancePattern(bannedModsRegexList)
	if err != nil {
		return nil, err
	}
	r.bannedModsRegex = pattern
	return r, nil
}

func (r *conformanceBannedModsRegexRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	docInfo := n.GetJSDocInfo()
	if docInfo == nil {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if mods, ok := docInfo.Mods(); ok && r.bannedModsRegex.Matcher(mods).Find() {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBannedModsRegexRule) GetPrecondition() CheckConformance_Precondition {
	return checkConformance_PreconditionBannedModsRegexScript
}
