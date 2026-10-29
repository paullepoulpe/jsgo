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

// ConformanceRules_NoOp is a no-op rule that never reports any violations.
type ConformanceRules_NoOp struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.NoOp(AbstractCompiler, Requirement)
func newConformanceRules_NoOp(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_NoOp, error) {
	r := &ConformanceRules_NoOp{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_NoOp) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.NoOp.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_NoOp) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_InferredConstCheck checks that variables annotated as
// @const have an inferred type, if there is no type given explicitly.
type ConformanceRules_InferredConstCheck struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.InferredConstCheck(AbstractCompiler, Requirement)
func newConformanceRules_InferredConstCheck(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_InferredConstCheck, error) {
	r := &ConformanceRules_InferredConstCheck{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_InferredConstCheck) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.InferredConstCheck.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_InferredConstCheck) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
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

// ConformanceRules_BannedDependency is the banned dependency rule.
type ConformanceRules_BannedDependency struct {
	ConformanceRules_AbstractRule

	paths []string
}

// Java: ConformanceRules.BannedDependency(AbstractCompiler, Requirement)
func newConformanceRules_BannedDependency(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedDependency, error) {
	r := &ConformanceRules_BannedDependency{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	r.paths = requirement.GetValueList()
	if len(r.paths) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}
	return r, nil
}

func (r *ConformanceRules_BannedDependency) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BannedDependency.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedDependency) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsScript() {
		srcFile := *n.GetSourceFileName()
		for i := 0; i < len(r.paths); i++ {
			if strings.HasPrefix(srcFile, r.paths[i]) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BannedDependencyRegex is the banned dependency via regex
// rule.
type ConformanceRules_BannedDependencyRegex struct {
	ConformanceRules_AbstractRule

	pathRegexp *jutil.Pattern
}

// Java: ConformanceRules.BannedDependencyRegex(AbstractCompiler, Requirement)
func newConformanceRules_BannedDependencyRegex(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedDependencyRegex, error) {
	r := &ConformanceRules_BannedDependencyRegex{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	pathRegexpList := requirement.GetValueList()
	if len(pathRegexpList) == 0 {
		return nil, errInvalidRequirementSpec("missing value (no banned dependency regexps)")
	}
	pathRegexp, err := conformanceRules_buildPattern(pathRegexpList)
	if err != nil {
		return nil, err
	}
	r.pathRegexp = pathRegexp
	return r, nil
}

func (r *ConformanceRules_BannedDependencyRegex) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BannedDependencyRegex.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedDependencyRegex) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsScript() {
		srcFile := *n.GetSourceFileName()
		if r.pathRegexp != nil && r.pathRegexp.Find(srcFile) {
			return ConformanceRules_ConformanceResult_VIOLATION
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BannedName is the banned name rule.
type ConformanceRules_BannedName struct {
	ConformanceRules_AbstractRule

	// LINT.ThenChange(//depot/google3/javascript/typescript/compiler/conformance/rules/ban_third_party_script.ts)
	requirementType Requirement_Type
	qualifiedNames  []rhino.Node
	shortNames      map[uint32]struct{}
}

// Java: ConformanceRules.BannedName(AbstractCompiler, Requirement)
func newConformanceRules_BannedName(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedName, error) {
	r := &ConformanceRules_BannedName{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
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

// Java: ConformanceRules.BannedName.getPrecondition()
func (r *ConformanceRules_BannedName) GetPrecondition() CheckConformance_Precondition {
	return checkConformance_PreconditionBannedNameCandidate
}

func (r *ConformanceRules_BannedName) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BannedName.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedName) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if r.requirementType == Requirement_Type_BANNED_NAME_CALL &&
		!ConformanceRules_ConformanceUtil_IsCallTarget(n) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	/*
	 * Efficiently filter out nearly all candidates.
	 *
	 * <p>Ideally we could use a hashset containing the qualified names, but it turns out that
	 * creating wrapper a object for every Node is more expensive than the loop below.
	 */
	if _, ok := r.shortNames[n.GetStringID()]; !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	/*
	 * Defer expensive infrequent checks.
	 *
	 * <p>These could be a precondition, but they don't usually need to be checked. Since they're
	 * kind of expensive, it's cheaper overall to defer them.
	 */
	if NodeUtil_IsInSyntheticScript(n) || !ConformanceRules_BannedName_IsRootOfQualifiedNameGlobal(t, n) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	for i := len(r.qualifiedNames) - 1; i >= 0; i-- {
		if n.MatchesQualifiedNameNode(r.qualifiedNames[i]) {
			return ConformanceRules_ConformanceResult_VIOLATION
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BannedName.isRootOfQualifiedNameGlobal(NodeTraversal, Node)
func ConformanceRules_BannedName_IsRootOfQualifiedNameGlobal(t *NodeTraversal, n rhino.Node) bool {
	rootNode := NodeUtil_GetRootOfQualifiedName(n)
	var v *Var
	if rootNode.IsName() {
		v = t.GetScope().GetVarByID(rootNode.GetStringID())
	} else {
		rootName := rootNode.GetQualifiedName()
		v = t.GetScope().GetVar(*rootName)
	}
	return v != nil && v.IsGlobal()
}

// ConformanceRules_BannedProperty is the banned property rule.
type ConformanceRules_BannedProperty struct {
	ConformanceRules_AbstractRule

	registry                *rhino.JSTypeRegistry
	props                   map[string][]rhino.JSType
	requirementPrecondition ConformanceRules_BannedProperty_RequirementPrecondition
}

// ConformanceRules_BannedProperty_RequirementPrecondition is the enum
// BannedProperty.RequirementPrecondition (a Precondition).
type ConformanceRules_BannedProperty_RequirementPrecondition int32

const (
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY ConformanceRules_BannedProperty_RequirementPrecondition = iota
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_WRITE
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_NON_CONSTANT_WRITE
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ
	ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL
)

var conformanceRules_BannedProperty_RequirementPrecondition_names = [...]string{
	"BANNED_PROPERTY",
	"BANNED_PROPERTY_WRITE",
	"BANNED_PROPERTY_NON_CONSTANT_WRITE",
	"BANNED_PROPERTY_READ",
	"BANNED_PROPERTY_CALL",
}

// String is Enum.name().
func (p ConformanceRules_BannedProperty_RequirementPrecondition) String() string {
	return conformanceRules_BannedProperty_RequirementPrecondition_names[p]
}

// Ordinal is Enum.ordinal().
func (p ConformanceRules_BannedProperty_RequirementPrecondition) Ordinal() int { return int(p) }

// ShouldCheck is the per-constant body of shouldCheck.
//
// Java: ConformanceRules.BannedProperty.RequirementPrecondition.shouldCheck(Node)
func (p ConformanceRules_BannedProperty_RequirementPrecondition) ShouldCheck(n rhino.Node) bool {
	switch p {
	case ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY:
		switch n.GetToken() {
		case rhino.Token_STRING_KEY, rhino.Token_GETPROP, rhino.Token_GETELEM, rhino.Token_COMPUTED_PROP:
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
		if asgn, ok := ast.AsAssign(n.GetParent()); ok && asgn.Target().Node() == n {
			rhs := asgn.Value().Node()
			if NodeUtil_IsLiteralValue(rhs, false) || NodeUtil_IsSomeCompileTimeConstStringValue(rhs) {
				return false
			}
		}
		return true
	case ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_READ:
		return !NodeUtil_IsLValue(n) && NodeUtil_IsExpressionResultUsed(n)
	case ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY_CALL:
		return ConformanceRules_ConformanceUtil_IsCallTarget(n)
	}
	panic(jutil.NewIllegalStateException())
}

// Java: ConformanceRules.BannedProperty(AbstractCompiler, Requirement)
func newConformanceRules_BannedProperty(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedProperty, error) {
	r := &ConformanceRules_BannedProperty{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}

	if requirement.GetValueCount() == 0 {
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
		typename, property, ok := conformanceRules_splitClassAndProperty(value)
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

// Java: ConformanceRules.BannedProperty.getPrecondition()
func (r *ConformanceRules_BannedProperty) GetPrecondition() CheckConformance_Precondition {
	return r.requirementPrecondition
}

func (r *ConformanceRules_BannedProperty) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BannedProperty.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedProperty) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	name, ok := r.ExtractName(n)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	checkTypes := r.props[name]
	if len(checkTypes) == 0 {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	/*
	 * Avoid type operations when possible.
	 *
	 * <p>checkTypes is almost always empty, and operations on unions can be expensive.
	 */
	foundType := r.ExtractType(n)
	if foundType == nil {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	foundType = foundType.RestrictByNotNullOrUndefined()

	for _, checkType := range checkTypes {
		result := r.MatchTypes(foundType, checkType)
		if result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
			return result
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BannedProperty.matchTypes(JSType, JSType)
func (r *ConformanceRules_BannedProperty) MatchTypes(foundType rhino.JSType, checkType rhino.JSType) *ConformanceRules_ConformanceResult {
	foundObj := foundType.ToMaybeObjectType()
	if foundObj != nil {
		if foundObj.IsFunctionPrototypeType() {
			ownerFun := foundObj.GetOwnerFunction()
			if ownerFun.IsConstructor() {
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
		ConformanceRules_IsLooseObject(foundType, r.registry) {
		if r.reportLooseTypeViolations {
			return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
		}
	} else if foundType.IsSubtypeOf(checkType) {
		return ConformanceRules_ConformanceResult_VIOLATION
	} else if checkType.IsSubtypeWithoutStructuralTyping(foundType) {
		if r.MatchesPrototype(checkType, foundType) {
			return ConformanceRules_ConformanceResult_VIOLATION
		} else if r.reportLooseTypeViolations {
			// Access of a banned property through a super class may be a violation
			return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BannedProperty.matchesPrototype(JSType, JSType)
func (r *ConformanceRules_BannedProperty) MatchesPrototype(type_ rhino.JSType, maybePrototype rhino.JSType) bool {
	methodClassObjectType := type_.ToMaybeObjectType()
	if methodClassObjectType != nil {
		if methodClassObjectType.GetImplicitPrototype().Equals(maybePrototype) {
			return true
		}
	}
	return false
}

// Java: ConformanceRules.BannedProperty.extractType(Node)
func (r *ConformanceRules_BannedProperty) ExtractType(n rhino.Node) rhino.JSType {
	if get, ok := ast.AsNormalGet(n); ok {
		return get.Object().Node().GetJSType()
	}
	switch n.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_COMPUTED_PROP:
		parent := n.GetParent()
		switch parent.GetToken() {
		case rhino.Token_OBJECT_PATTERN, rhino.Token_OBJECTLIT:
			return parent.GetJSType()
		case rhino.Token_CLASS_MEMBERS:
			return nil
		default:
			panic(jutil.NewAssertionError("null"))
		}
	default:
		return nil
	}
}

// Java: ConformanceRules.BannedProperty.extractName(Node)
func (r *ConformanceRules_BannedProperty) ExtractName(n rhino.Node) (string, bool) {
	if get, ok := ast.AsNormalGet(n); ok {
		return get.StaticPropertyName()
	}
	if sk, ok := ast.AsStringKey(n); ok {
		return sk.Key(), true
	}
	if cp, ok := ast.AsComputedProp(n); ok {
		if s, ok := ast.AsStringLit(cp.Key().Node()); ok {
			return s.Value(), true
		}
		return "", false
	}
	return "", false
}

// ConformanceRules_BannedStringRegex is the banned string rule.
type ConformanceRules_BannedStringRegex struct {
	ConformanceRules_AbstractRule

	stringPattern *jutil.Pattern
}

// Java: ConformanceRules.BannedStringRegex(AbstractCompiler, Requirement)
func newConformanceRules_BannedStringRegex(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BannedStringRegex, error) {
	r := &ConformanceRules_BannedStringRegex{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
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

	stringRegex, err := conformanceRules_buildPattern(values)
	if err != nil {
		return nil, err
	}
	r.stringPattern = stringRegex
	return r, nil
}

func (r *ConformanceRules_BannedStringRegex) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BannedStringRegex.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedStringRegex) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
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
			jutil.CheckStateMsg1(
				ConformanceRules_IsTaggedTemplateLitString(node),
				"Found untagged template literal with invalid (uncookable) escape sequence: %s",
				rawString)
			if r.stringPattern.Matches(rawString) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_IsTaggedTemplateLitString: is this template literal string
// a tagged template literal string?
//
// Java: ConformanceRules.isTaggedTemplateLitString(Node)
func ConformanceRules_IsTaggedTemplateLitString(node rhino.Node) bool {
	jutil.CheckState(node.IsTemplateLitString())
	gp := node.GetGrandparent()
	return !gp.IsNil() && gp.IsTaggedTemplateLit()
}

// ConformanceRules_RestrictedNameCall is the restricted name call rule.
type ConformanceRules_RestrictedNameCall struct {
	ConformanceRules_AbstractRule

	restrictions []ConformanceRules_RestrictedNameCall_Restriction
}
