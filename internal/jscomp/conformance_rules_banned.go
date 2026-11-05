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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type conformanceNoOpRule struct {
	conformanceAbstractRule
}

func newConformanceRules_NoOp(compiler AbstractCompiler, requirement *Requirement) (*conformanceNoOpRule, error) {
	r := &conformanceNoOpRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceNoOpRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// conformanceInferredConstRule checks that variables annotated `@const` without
// an explicit type annotation do not infer to the unknown type.
type conformanceInferredConstRule struct {
	conformanceAbstractRule
}

func newConformanceRules_InferredConstCheck(compiler AbstractCompiler, requirement *Requirement) (*conformanceInferredConstRule, error) {
	r := &conformanceInferredConstRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceInferredConstRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	jsDoc := n.GetJSDocInfo()
	if jsDoc != nil && jsDoc.HasConstAnnotation() && jsDoc.GetType() == nil {
		if asgn, ok := ast.AsAssign(n); ok {
			n = asgn.Target().Node()
		}
		type_ := n.GetJSType()
		if type_ != nil && type_.IsUnknownType() && !NodeUtil_IsNamespaceDecl(n) {
			return ConformanceRules_ConformanceResult_VIOLATION
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

type conformanceBannedDependencyRule struct {
	conformanceAbstractRule

	paths []string
}

func newConformanceRules_BannedDependency(compiler AbstractCompiler, requirement *Requirement) (*conformanceBannedDependencyRule, error) {
	r := &conformanceBannedDependencyRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	r.paths = requirement.GetValueList()
	if len(r.paths) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}
	return r, nil
}

func (r *conformanceBannedDependencyRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if _, ok := ast.AsScript(n); ok {
		srcFile, hasSrc := n.SourceFileName()
		if !hasSrc {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		for _, path := range r.paths {
			if strings.HasPrefix(srcFile, path) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

type conformanceBannedDependencyRegexRule struct {
	conformanceAbstractRule

	pathRegexp *jutil.Pattern
}

func newConformanceRules_BannedDependencyRegex(compiler AbstractCompiler, requirement *Requirement) (*conformanceBannedDependencyRegexRule, error) {
	r := &conformanceBannedDependencyRegexRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	pathRegexpList := requirement.GetValueList()
	if len(pathRegexpList) == 0 {
		return nil, errInvalidRequirementSpec("missing value (no banned dependency regexps)")
	}
	pathRegexp, err := buildConformancePattern(pathRegexpList)
	if err != nil {
		return nil, err
	}
	r.pathRegexp = pathRegexp
	return r, nil
}

func (r *conformanceBannedDependencyRegexRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if _, ok := ast.AsScript(n); ok {
		srcFile, hasSrc := n.SourceFileName()
		if !hasSrc {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		if r.pathRegexp != nil && r.pathRegexp.Find(srcFile) {
			return ConformanceRules_ConformanceResult_VIOLATION
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// CustomBannedNameEvaluator is an optional embedder hook invoked when a
// BannedName rule (BANNED_NAME, BANNED_NAME_CALL, or BANNED_SCRIPT_IMPORT)
// matches a global qualified name node. When handled is true, res is returned
// directly; otherwise the default VIOLATION result is returned.
var CustomBannedNameEvaluator func(reqType Requirement_Type, n rhino.Node) (res *ConformanceRules_ConformanceResult, handled bool)

type conformanceBannedNameRule struct {
	conformanceAbstractRule

	// LINT.ThenChange(//depot/google3/javascript/typescript/compiler/conformance/rules/ban_third_party_script.ts)
	requirementType Requirement_Type
	qualifiedNames  []rhino.Node
	shortNames      map[uint32]struct{}
}

func newConformanceRules_BannedName(compiler AbstractCompiler, requirement *Requirement) (*conformanceBannedNameRule, error) {
	r := &conformanceBannedNameRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}
	r.requirementType = requirement.GetType()

	qualifiedBuilder := make([]rhino.Node, 0, len(values))
	shortBuilder := make(map[uint32]struct{}, len(values))
	for _, name := range values {
		qualifiedName := NodeUtil_NewQName(compiler, name)
		qualifiedBuilder = append(qualifiedBuilder, qualifiedName)
		shortBuilder[qualifiedName.GetStringID()] = struct{}{}
	}
	r.qualifiedNames = qualifiedBuilder
	r.shortNames = shortBuilder
	return r, nil
}

func (r *conformanceBannedNameRule) GetPrecondition() CheckConformance_Precondition {
	return checkConformance_PreconditionBannedNameCandidate
}

func (r *conformanceBannedNameRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if r.requirementType == Requirement_Type_BANNED_NAME_CALL && !isConformanceCallTarget(n) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// Fast short-name prefilter before inspecting scopes or full qualified chains.
	if _, ok := r.shortNames[n.GetStringID()]; !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if NodeUtil_IsInSyntheticScript(n) || !isRootOfQualifiedNameGlobal(t, n) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	for i := len(r.qualifiedNames) - 1; i >= 0; i-- {
		if n.MatchesQualifiedNameNode(r.qualifiedNames[i]) {
			if CustomBannedNameEvaluator != nil {
				if res, handled := CustomBannedNameEvaluator(r.requirementType, n); handled {
					return res
				}
			}
			return ConformanceRules_ConformanceResult_VIOLATION
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func isRootOfQualifiedNameGlobal(t *NodeTraversal, n rhino.Node) bool {
	if n.IsNil() || t == nil || t.GetScope() == nil {
		return false
	}
	rootNode := NodeUtil_GetRootOfQualifiedName(n)
	if rootNode.IsNil() {
		return false
	}
	var v *Var
	if name, ok := ast.AsName(rootNode); ok {
		v = t.GetScope().GetVarByID(name.Node().GetStringID())
	} else {
		rootName, ok := rootNode.GetQualifiedNameString()
		if !ok {
			return false
		}
		v = t.GetScope().GetVar(rootName)
	}
	return v != nil && v.IsGlobal()
}

// ConformanceRules_BannedProperty implements BANNED_PROPERTY,
// BANNED_PROPERTY_READ, BANNED_PROPERTY_WRITE, BANNED_PROPERTY_NON_CONSTANT_WRITE,
// and BANNED_PROPERTY_CALL requirements.
type ConformanceRules_BannedProperty struct {
	conformanceAbstractRule

	registry                *rhino.JSTypeRegistry
	props                   map[string][]rhino.JSType
	requirementPrecondition ConformanceRules_BannedProperty_RequirementPrecondition
}

// ConformanceRules_BannedProperty_RequirementPrecondition groups banned-property
// rules sharing the same syntactic candidate filter.
type ConformanceRules_BannedProperty_RequirementPrecondition int32

const (
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY ConformanceRules_BannedProperty_RequirementPrecondition = iota
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL
)

var bannedPropertyPreconditionNames = [...]string{
	"BANNED_PROPERTY",
	"BANNED_PROPERTY_WRITE",
	"BANNED_PROPERTY_NON_CONSTANT_WRITE",
	"BANNED_PROPERTY_READ",
	"BANNED_PROPERTY_CALL",
}

func (p ConformanceRules_BannedProperty_RequirementPrecondition) String() string {
	return bannedPropertyPreconditionNames[p]
}

func (p ConformanceRules_BannedProperty_RequirementPrecondition) ShouldCheck(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	switch p {
	case ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY:
		switch n.GetToken() {
		case rhino.Token_STRING_KEY,
			rhino.Token_GETPROP,
			rhino.Token_GETELEM,
			rhino.Token_OPTCHAIN_GETPROP,
			rhino.Token_OPTCHAIN_GETELEM,
			rhino.Token_COMPUTED_PROP:
			return true
		default:
			return false
		}
	case ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE:
		return NodeUtil_IsLValue(n)
	case ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE:
		if !NodeUtil_IsLValue(n) {
			return false
		}
		if asgn, ok := ast.AsAssign(n.GetParent()); ok && asgn.Target().Is(n) {
			rhs := asgn.Value().Node()
			if NodeUtil_IsLiteralValue(rhs, false) || NodeUtil_IsSomeCompileTimeConstStringValue(rhs) {
				return false
			}
		}
		return true
	case ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ:
		return !NodeUtil_IsLValue(n) && !n.GetParent().IsNil() && NodeUtil_IsExpressionResultUsed(n)
	case ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL:
		return isConformanceCallTarget(n)
	}
	panic(jutil.NewIllegalStateException())
}

func newConformanceRules_BannedProperty(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedProperty, error) {
	r := &ConformanceRules_BannedProperty{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}

	if len(requirement.GetValueList()) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	switch requirement.GetType() {
	case Requirement_Type_BANNED_PROPERTY:
		r.requirementPrecondition = ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY
	case Requirement_Type_BANNED_PROPERTY_READ:
		r.requirementPrecondition = ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ
	case Requirement_Type_BANNED_PROPERTY_WRITE:
		r.requirementPrecondition = ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE
	case Requirement_Type_BANNED_PROPERTY_NON_CONSTANT_WRITE:
		r.requirementPrecondition =
			ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE
	case Requirement_Type_BANNED_PROPERTY_CALL:
		r.requirementPrecondition = ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL
	default:
		panic(jutil.NewAssertionError(requirement.GetType().String()))
	}

	r.registry = compiler.GetTypeRegistry()

	props := make(map[string][]rhino.JSType)
	for _, value := range requirement.GetValueList() {
		typename, property, ok := splitConformanceClassAndProperty(value)
		if !ok {
			return nil, errInvalidRequirementSpec("bad prop value")
		}

		type_ := r.registry.GetGlobalType(typename)
		if type_ == nil {
			type_ = r.registry.ResolveViaClosureNamespace(typename)
		}

		// If type doesn't exist in the compilation, it can't be a violation. Also, bottom types
		// match everything, which is almost surely not what the check is intended to check against.
		if type_ == nil || type_.IsUnknownType() || type_.IsEmptyType() {
			continue
		}

		if !slices.ContainsFunc(props[property], func(existing rhino.JSType) bool { return type_.Equals(existing) }) {
			props[property] = append(props[property], type_)
		}
	}
	r.props = props
	return r, nil
}

func (r *ConformanceRules_BannedProperty) GetPrecondition() CheckConformance_Precondition {
	return r.requirementPrecondition
}

func (r *ConformanceRules_BannedProperty) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	name, ok := r.extractName(n)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	checkTypes := r.props[name]
	if len(checkTypes) == 0 {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// Avoid union/type operations until we know the property name is actually tracked.
	foundType := r.extractType(n)
	if foundType == nil {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	foundType = foundType.RestrictByNotNullOrUndefined()

	for _, checkType := range checkTypes {
		result := r.matchTypes(foundType, checkType)
		if result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
			return result
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *ConformanceRules_BannedProperty) matchTypes(foundType, checkType rhino.JSType) *ConformanceRules_ConformanceResult {
	if foundObj := foundType.ToMaybeObjectType(); foundObj != nil {
		if foundObj.IsFunctionPrototypeType() {
			ownerFun := foundObj.GetOwnerFunction()
			if ownerFun != nil && ownerFun.IsConstructor() {
				foundType = ownerFun.GetInstanceType()
			}
		} else if foundObj.IsTemplatizedType() {
			foundType = foundObj.GetRawType()
		}
	}

	if foundType.IsUnknownType() ||
		foundType.IsTemplateType() ||
		foundType.IsEmptyType() ||
		foundType.IsAllType() ||
		isLooseObjectType(foundType, r.registry) {
		if r.reportLooseTypeViolations {
			return conformancePossibleViolationLooseTypes
		}
	} else if foundType.IsSubtypeOf(checkType) {
		return ConformanceRules_ConformanceResult_VIOLATION
	} else if checkType.IsSubtypeWithoutStructuralTyping(foundType) {
		if r.matchesPrototype(checkType, foundType) {
			return ConformanceRules_ConformanceResult_VIOLATION
		} else if r.reportLooseTypeViolations {
			// Access of a banned property through a super class may be a violation.
			return conformancePossibleViolationLooseTypes
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *ConformanceRules_BannedProperty) matchesPrototype(type_, maybePrototype rhino.JSType) bool {
	methodClassObjectType := type_.ToMaybeObjectType()
	return methodClassObjectType != nil &&
		methodClassObjectType.GetImplicitPrototype() != nil &&
		methodClassObjectType.GetImplicitPrototype().Equals(maybePrototype)
}

func (r *ConformanceRules_BannedProperty) extractType(n rhino.Node) rhino.JSType {
	if n.IsNil() {
		return nil
	}
	if get, ok := ast.AsNormalOrOptChainGet(n); ok {
		return get.Object().Node().GetJSType()
	}
	switch n.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_COMPUTED_PROP:
		parent := n.GetParent()
		if parent.IsNil() {
			return nil
		}
		switch parent.GetToken() {
		case rhino.Token_OBJECT_PATTERN, rhino.Token_OBJECTLIT:
			return parent.GetJSType()
		case rhino.Token_CLASS_MEMBERS:
			return nil
		default:
			panic(jutil.NewAssertionError("unexpected parent for property node: " + parent.GetToken().String()))
		}
	default:
		return nil
	}
}

func (r *ConformanceRules_BannedProperty) extractName(n rhino.Node) (string, bool) {
	if get, ok := ast.AsNormalOrOptChainGet(n); ok {
		return get.StaticPropertyName()
	}
	if sk, ok := ast.AsStringKey(n); ok {
		return sk.Key(), true
	}
	if cp, ok := ast.AsComputedProp(n); ok {
		if s, ok := cp.Key().AsStringLit(); ok {
			return s.Value(), true
		}
		return "", false
	}
	return "", false
}

type conformanceBannedStringRegexRule struct {
	conformanceAbstractRule

	stringPattern *jutil.Pattern
}

func newConformanceRules_BannedStringRegex(compiler AbstractCompiler, requirement *Requirement) (*conformanceBannedStringRegexRule, error) {
	r := &conformanceBannedStringRegexRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	for _, value := range values {
		if jutil.StringTrim(value) == "" {
			return nil, errInvalidRequirementSpec("empty strings or whitespace are not allowed")
		}
	}

	stringRegex, err := buildConformancePattern(values)
	if err != nil {
		return nil, err
	}
	r.stringPattern = stringRegex
	return r, nil
}

func (r *conformanceBannedStringRegexRule) checkNode(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if node.IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if s, ok := ast.AsStringLit(node); ok {
		if r.stringPattern.Matches(s.Value()) {
			return ConformanceRules_ConformanceResult_VIOLATION
		}
	} else if tls, ok := ast.AsTemplateLitString(node); ok {
		if cookedString, ok := tls.Cooked(); ok {
			if r.stringPattern.Matches(cookedString) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		} else {
			rawString := tls.Raw()
			if !isTaggedTemplateLitString(tls) {
				panic(jutil.NewIllegalStateExceptionMsg(
					"Found untagged template literal with invalid (uncookable) escape sequence: " + rawString))
			}
			if r.stringPattern.Matches(rawString) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func isTaggedTemplateLitString(tls ast.TemplateLitString) bool {
	parent := tls.Node().GetParent()
	if parent.IsNil() {
		return false
	}
	gp := parent.GetParent()
	return !gp.IsNil() && gp.IsTaggedTemplateLit()
}

type conformanceRestrictedNameCallRule struct {
	conformanceAbstractRule

	restrictions []restrictedNameCallSpec
}
