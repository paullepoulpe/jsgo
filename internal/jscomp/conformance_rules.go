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
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ConformanceRules holds the standard conformance rules. See
// third_party/java_src/jscomp/java/com/google/javascript/jscomp/conformance.proto
//
// Java: ConformanceRules() is private (non-instantiable); the class has only
// static members.
type ConformanceRules struct{}

// Java: ConformanceRules.ALL_TS_ALLOWLIST
var conformanceRules_ALL_TS_ALLOWLIST = ConformanceRules_CreateTsAllowlist()

// Java: ConformanceRules.ON_PROTOTYPE (Splitter.on(".prototype."))
const conformanceRules_ON_PROTOTYPE = ".prototype."

// Java: ConformanceRules.ON_DOT (Splitter.on("."))
const conformanceRules_ON_DOT = "."

func conformanceRules_splitToList(separator string, s string) []string {
	return strings.Split(s, separator)
}

// Java: ConformanceRules.createTsAllowlist()
func ConformanceRules_CreateTsAllowlist() *ConformanceRules_AllowList {
	result, err := newConformanceRules_AllowList(
		[]string{}, []string{".*\\.closure\\.js", ".*\\.tsx\\.cl\\.js"})
	if err != nil {
		panic(jutil.NewAssertionError(err.Error()))
	}
	return result
}

// ConformanceRules_ConformanceResult: classes extending AbstractRule must
// return ConformanceResult from their checkConformance implementation. For
// simple rules, the constants CONFORMANCE, POSSIBLE_VIOLATION, VIOLATION are
// sufficient. However, for some rules additional clarification specific to
// the violation instance is helpful, for that, an instance of this class can
// be created to associate a note with the violation.
type ConformanceRules_ConformanceResult struct {
	level ConformanceRules_ConformanceLevel
	note  string
}

// Java: ConformanceRules.ConformanceResult(ConformanceLevel)
func NewConformanceRules_ConformanceResult(level ConformanceRules_ConformanceLevel) *ConformanceRules_ConformanceResult {
	return NewConformanceRules_ConformanceResultWithNote(level, "")
}

// Java: ConformanceRules.ConformanceResult(ConformanceLevel, String)
func NewConformanceRules_ConformanceResultWithNote(level ConformanceRules_ConformanceLevel, note string) *ConformanceRules_ConformanceResult {
	return &ConformanceRules_ConformanceResult{level: level, note: note}
}

// For CONFORMANCE rules that don't generate notes:
//
// Java: ConformanceRules.ConformanceResult.CONFORMANCE
var ConformanceRules_ConformanceResult_CONFORMANCE = NewConformanceRules_ConformanceResult(ConformanceRules_ConformanceLevel_CONFORMANCE)

// Java: ConformanceRules.ConformanceResult.POSSIBLE_VIOLATION
var ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION = NewConformanceRules_ConformanceResult(ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION)

// Java: ConformanceRules.ConformanceResult.POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
var conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES = NewConformanceRules_ConformanceResultWithNote(
	ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION,
	"The type information available for this expression is too loose "+
		"to ensure conformance.")

// Java: ConformanceRules.ConformanceResult.VIOLATION
var ConformanceRules_ConformanceResult_VIOLATION = NewConformanceRules_ConformanceResult(ConformanceRules_ConformanceLevel_VIOLATION)

// ConformanceRules_ConformanceLevel is the possible check results.
type ConformanceRules_ConformanceLevel int32

const (
	// Nothing interesting detected.
	ConformanceRules_ConformanceLevel_CONFORMANCE ConformanceRules_ConformanceLevel = iota
	// In the optionally typed world of the Closure Compiler type system
	// it is possible that detect patterns that match with looser types
	// that the target pattern.
	ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION
	// Definitely a violation.
	ConformanceRules_ConformanceLevel_VIOLATION
)

var conformanceRules_ConformanceLevel_names = [...]string{"CONFORMANCE", "POSSIBLE_VIOLATION", "VIOLATION"}

// String is Enum.name().
func (l ConformanceRules_ConformanceLevel) String() string {
	return conformanceRules_ConformanceLevel_names[l]
}

// Ordinal is Enum.ordinal().
func (l ConformanceRules_ConformanceLevel) Ordinal() int { return int(l) }

func compileConformancePattern(regex string) (p *jutil.Pattern, synErr *jutil.PatternSyntaxException) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*jutil.PatternSyntaxException); ok {
				synErr = e
				return
			}
			panic(r)
		}
	}()
	return jutil.Pattern_Compile(regex), nil
}

// Java: ConformanceRules.buildPattern(List<String>)
func conformanceRules_buildPattern(reqPatterns []string) (*jutil.Pattern, error) {
	if len(reqPatterns) == 0 {
		return nil, nil
	}

	// validate the patterns
	for _, reqPattern := range reqPatterns {
		if _, err := compileConformancePattern(reqPattern); err != nil {
			return nil, errInvalidRequirementSpec("invalid regex pattern")
		}
	}

	jointRegExp := "(" + strings.Join(reqPatterns, "|") + ")"
	pattern, err := compileConformancePattern(jointRegExp)
	if err != nil {
		panic(jutil.NewRuntimeExceptionWithMessageAndCause("bad joined regexp", err))
	}
	return pattern, nil
}

// Java: ConformanceRules.AllowList
type ConformanceRules_AllowList struct {
	prefixes       []string
	regexp         *jutil.Pattern
	allowlistEntry *RequirementScopeEntry
}

// Java: ConformanceRules.AllowList(List<String>, List<String>)
func newConformanceRules_AllowList(prefixes []string, regexps []string) (*ConformanceRules_AllowList, error) {
	regexp, err := conformanceRules_buildPattern(regexps)
	if err != nil {
		return nil, err
	}
	return &ConformanceRules_AllowList{
		prefixes: slices.Clone(prefixes),
		regexp:   regexp,
	}, nil
}

// Java: ConformanceRules.AllowList(RequirementScopeEntry)
func newConformanceRules_AllowListFromEntry(allowlistEntry *RequirementScopeEntry) (*ConformanceRules_AllowList, error) {
	regexp, err := conformanceRules_buildPattern(allowlistEntry.GetRegexpList())
	if err != nil {
		return nil, err
	}
	return &ConformanceRules_AllowList{
		prefixes:       slices.Clone(allowlistEntry.GetPrefixList()),
		regexp:         regexp,
		allowlistEntry: allowlistEntry,
	}, nil
}

// Matches returns true if the given path matches one of the prefixes or
// regexps, and false otherwise.
//
// Java: ConformanceRules.AllowList.matches(String)
func (a *ConformanceRules_AllowList) Matches(path string) bool {
	// If the path ends with .closure.js or .tsx.cl.js, it is probably a tsickle-generated file,
	// and there may be entries in the allow list for the TypeScript path
	var tsPath string
	hasTsPath := false
	if strings.HasSuffix(path, ".closure.js") {
		tsPath = path[:len(path)-len(".closure.js")] + ".ts"
		hasTsPath = true
	} else if strings.HasSuffix(path, ".tsx.cl.js") {
		// TSX
		tsPath = path[:len(path)-len(".cl.js")]
		hasTsPath = true
	}
	if a.prefixes != nil && path != "" {
		for _, prefix := range a.prefixes {
			if strings.HasPrefix(path, prefix) || (hasTsPath && strings.HasPrefix(tsPath, prefix)) {
				return true
			}
		}
	}

	return a.regexp != nil &&
		(a.regexp.Matcher(path).Find() || (hasTsPath && a.regexp.Matcher(tsPath).Find()))
}

// ConformanceRules_AbstractRule holds the common rule fields and helpers shared
// by AbstractRule and AbstractTypeRestrictionRule subclasses, without virtual
// self-dispatch pointers.
type ConformanceRules_AbstractRule struct {
	CheckConformance_RuleDefaults

	compiler                  AbstractCompiler
	message                   string
	severity                  Requirement_Severity
	allowlists                []*ConformanceRules_AllowList
	onlyApplyTo               *ConformanceRules_AllowList
	reportLooseTypeViolations bool
	typeMatchingStrategy      TypeMatchingStrategy
	requirement               *Requirement

	// Type-restriction rule fields (populated by initTypeRestrictionRule).
	nativeObjectType   rhino.JSType
	allowlistedTypes   rhino.JSType
	assertionFunctions *CodingConvention_AssertionFunctionLookup
}

// initAbstractRule initializes the common AbstractRule state.
//
// Java: ConformanceRules.AbstractRule(AbstractCompiler, Requirement)
func (b *ConformanceRules_AbstractRule) initAbstractRule(
	compiler AbstractCompiler, requirement *Requirement, tsIsAllowlisted bool) error {
	if !requirement.HasErrorMessage() {
		return errInvalidRequirementSpec("missing message")
	}
	b.compiler = compiler
	message := requirement.GetErrorMessage()
	if requirement.GetConfigFileCount() > 0 {
		message +=
			"\n  defined in " + strings.Join(requirement.GetConfigFileList(), "\n  extended by ")
	}
	b.message = message
	if requirement.GetSeverity() == Requirement_Severity_UNSPECIFIED {
		b.severity = Requirement_Severity_WARNING
	} else {
		b.severity = requirement.GetSeverity()
	}

	// build allowlists
	var allowlistsBuilder []*ConformanceRules_AllowList
	for _, entry := range requirement.GetWhitelistEntryList() {
		allowlist, err := newConformanceRules_AllowListFromEntry(entry)
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	for _, entry := range requirement.GetAllowlistEntryList() {
		allowlist, err := newConformanceRules_AllowListFromEntry(entry)
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}

	if tsIsAllowlisted {
		allowlistsBuilder = append(allowlistsBuilder, conformanceRules_ALL_TS_ALLOWLIST)
	}

	if requirement.GetWhitelistCount() > 0 || requirement.GetWhitelistRegexpCount() > 0 {
		allowlist, err := newConformanceRules_AllowList(requirement.GetWhitelistList(), requirement.GetWhitelistRegexpList())
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	if requirement.GetAllowlistCount() > 0 || requirement.GetAllowlistRegexpCount() > 0 {
		allowlist, err := newConformanceRules_AllowList(requirement.GetAllowlistList(), requirement.GetAllowlistRegexpList())
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	b.allowlists = allowlistsBuilder

	if requirement.GetOnlyApplyToCount() > 0 || requirement.GetOnlyApplyToRegexpCount() > 0 {
		onlyApplyTo, err := newConformanceRules_AllowList(requirement.GetOnlyApplyToList(), requirement.GetOnlyApplyToRegexpList())
		if err != nil {
			return err
		}
		b.onlyApplyTo = onlyApplyTo
	} else {
		b.onlyApplyTo = nil
	}
	b.reportLooseTypeViolations = requirement.GetReportLooseTypeViolations()
	b.typeMatchingStrategy = ConformanceRules_AbstractRule_GetTypeMatchingStrategy(requirement)
	b.requirement = requirement
	return nil
}

// initTypeRestrictionRule initializes AbstractTypeRestrictionRule state.
//
// Java: ConformanceRules.AbstractTypeRestrictionRule(AbstractCompiler, Requirement)
func (b *ConformanceRules_AbstractRule) initTypeRestrictionRule(
	compiler AbstractCompiler, requirement *Requirement, tsIsAllowlisted bool) error {
	if err := b.initAbstractRule(compiler, requirement, tsIsAllowlisted); err != nil {
		return err
	}
	b.nativeObjectType = compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
	b.allowlistedTypes = b.Union(requirement.GetValueList())
	b.assertionFunctions =
		CodingConvention_AssertionFunctionLookup_Of(compiler.GetCodingConvention().GetAssertionFunctions())
	return nil
}

// Java: ConformanceRules.AbstractRule.getTypeMatchingStrategy(Requirement)
func ConformanceRules_AbstractRule_GetTypeMatchingStrategy(requirement *Requirement) TypeMatchingStrategy {
	switch requirement.GetTypeMatchingStrategy() {
	case Requirement_TypeMatchingStrategy_LOOSE:
		return TypeMatchingStrategy_LOOSE
	case Requirement_TypeMatchingStrategy_STRICT_NULLABILITY:
		return TypeMatchingStrategy_STRICT_NULLABILITY
	case Requirement_TypeMatchingStrategy_SUBTYPES:
		return TypeMatchingStrategy_SUBTYPES
	case Requirement_TypeMatchingStrategy_EXACT:
		return TypeMatchingStrategy_EXACT
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unknown TypeMatchingStrategy"))
	}
}

// FindAllowListForPath returns the first AllowList entry that matches the
// given path, and null otherwise.
//
// Java: ConformanceRules.AbstractRule.findAllowListForPath(String)
func (b *ConformanceRules_AbstractRule) FindAllowListForPath(path string) *ConformanceRules_AllowList {
	pathRegex := b.compiler.GetOptions().GetConformanceRemoveRegexFromPath()
	if pathRegex != nil {
		path = pathRegex.Matcher(path).ReplaceFirst("")
	}

	for _, allowlist := range b.allowlists {
		if allowlist.Matches(path) {
			return allowlist
		}
	}
	return nil
}

// Report reports a conformance warning for the given node.
//
// Java: ConformanceRules.AbstractRule.report(Node, ConformanceResult, LibraryLevelNonAllowlistedConformanceViolationsBehavior)
func (b *ConformanceRules_AbstractRule) Report(
	n rhino.Node,
	result *ConformanceRules_ConformanceResult,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	var msg *DiagnosticType
	if b.severity == Requirement_Severity_ERROR {
		// Always report findings that are errors, even if the types are too loose to be certain.
		msg = CheckConformance_CONFORMANCE_ERROR
	} else if result.level == ConformanceRules_ConformanceLevel_VIOLATION {
		msg = CheckConformance_CONFORMANCE_VIOLATION
	} else {
		msg = CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION
	}
	separator := "\n"
	if result.note == "" {
		separator = ""
	}
	err := JSError_MakeWithRequirement(b.requirement, n, msg, b.message, separator, result.note)

	path := NodeUtil_GetSourceName(n)
	var allowlist *ConformanceRules_AllowList
	if path != nil {
		allowlist = b.FindAllowListForPath(*path)
	}
	isAllowlisted :=
		!(allowlist == nil && (b.onlyApplyTo == nil || b.onlyApplyTo.Matches(*path)))
	var allowlistEntry *RequirementScopeEntry
	if allowlist != nil {
		allowlistEntry = allowlist.allowlistEntry
	}
	shouldReport :=
		b.compiler.
			GetErrorManager().
			// returns true even if the violation is allowlisted
			ShouldReportConformanceViolation(
				b.requirement,
				allowlistEntry,
				err,
				behavior,
				isAllowlisted)

	if shouldReport && !isAllowlisted {
		b.compiler.Report(err)
	}
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isAllowlistedType(Node)
func (b *ConformanceRules_AbstractRule) IsAllowlistedType(n rhino.Node) bool {
	if b.allowlistedTypes != nil && n.GetJSType() != nil {
		targetType := n.GetJSType().RestrictByNotNullOrUndefined()
		if targetType.IsSubtypeOf(b.allowlistedTypes) {
			return true
		}
	}
	return false
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isKnown(Node)
func ConformanceRules_AbstractTypeRestrictionRule_IsKnown(n rhino.Node) bool {
	return !ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n) &&
		!ConformanceRules_AbstractTypeRestrictionRule_IsBottom(n) &&
		!ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(n)
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isNativeObjectType(Node)
func (b *ConformanceRules_AbstractRule) IsNativeObjectType(n rhino.Node) bool {
	type_ := n.GetJSType().RestrictByNotNullOrUndefined()
	return type_.Equals(b.nativeObjectType)
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isTop(Node)
func ConformanceRules_AbstractTypeRestrictionRule_IsTop(n rhino.Node) bool {
	type_ := n.GetJSType()
	return type_ != nil && type_.IsAllType()
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isUnknown(Node)
func ConformanceRules_AbstractTypeRestrictionRule_IsUnknown(n rhino.Node) bool {
	type_ := n.GetJSType()
	return type_ == nil || type_.IsUnknownType()
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isTemplateType(Node)
func ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(n rhino.Node) bool {
	type_ := n.GetJSType()
	if type_.IsUnionType() {
		for _, member := range type_.GetUnionMembers() {
			if member.IsTemplateType() {
				return true
			}
		}
	}
	return type_.IsTemplateType()
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isBottom(Node)
func ConformanceRules_AbstractTypeRestrictionRule_IsBottom(n rhino.Node) bool {
	type_ := n.GetJSType().RestrictByNotNullOrUndefined()
	return type_.IsEmptyType()
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.union(List<String>)
func (b *ConformanceRules_AbstractRule) Union(typeNames []string) rhino.JSType {
	registry := b.compiler.GetTypeRegistry()
	var types []rhino.JSType

	for _, typeName := range typeNames {
		type_ := registry.GetGlobalType(typeName)
		if type_ == nil {
			// For a few types like `google3.javascript.apps.wiz.events.wiz_event.WizEvent`, we need
			// to resolve via closure namespace. This is because the class type WizEvent is exported
			// from a goog.module, so is not a global type name in the global type registry.
			type_ = registry.ResolveViaClosureNamespace(typeName)
		}
		if type_ != nil {
			types = append(types, type_)
		}
	}
	if len(types) == 0 {
		return nil
	}
	return registry.CreateUnionTypeOfList(types)
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isAssertionCall(Node)
func (b *ConformanceRules_AbstractRule) IsAssertionCall(n rhino.Node) bool {
	if call, ok := ast.AsCall(n); ok {
		callee := call.Callee().Node()
		if callee.IsQualifiedName() {
			return b.assertionFunctions.LookupByCallee(callee) != nil
		}
	}
	return false
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isTypeImmediatelyTightened(Node)
func (b *ConformanceRules_AbstractRule) IsTypeImmediatelyTightened(n rhino.Node) bool {
	parent := n.GetParent()
	return b.IsAssertionCall(parent) ||
		parent.IsTypeOf() ||
		/* casted node */ n.GetJSTypeBeforeCast() != nil
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isUsed(Node)
func (b *ConformanceRules_AbstractRule) IsUsed(n rhino.Node) bool {
	parent := n.GetParent()
	if parent.IsName() || NodeUtil_IsLhsByDestructuring(n) {
		return false
	}

	// Consider lvalues in assignment operations to be used iff the actual assignment
	// operation's result is used. e.g. for `a.b.c`:
	//     USED: `alert(x = a.b.c);`
	//   UNUSED: `x = a.b.c;`
	if NodeUtil_IsAssignmentOp(parent) {
		return NodeUtil_IsExpressionResultUsed(parent)
	}

	return NodeUtil_IsExpressionResultUsed(n)
}

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
		if r.pathRegexp != nil && r.pathRegexp.Matcher(srcFile).Find() {
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

	typeSets := make(map[string]*jutil.ObjectLinkedHashSet[rhino.JSType])
	var propOrder []string
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

		types, exists := typeSets[property]
		if !exists {
			types = jutil.NewObjectLinkedHashSet[rhino.JSType]()
			typeSets[property] = types
			propOrder = append(propOrder, property)
		}
		types.Add(type_)
	}
	props := make(map[string][]rhino.JSType, len(propOrder))
	for _, prop := range propOrder {
		props[prop] = typeSets[prop].ToSlice()
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
		if r.stringPattern.Matcher(s.Value()).Matches() {
			return ConformanceRules_ConformanceResult_VIOLATION
		}
	} else if tls, ok := ast.AsTemplateLitString(node); ok {
		if cookedString, ok := tls.Cooked(); ok {
			if r.stringPattern.Matcher(cookedString).Matches() {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		} else {
			rawString := tls.Raw()
			jutil.CheckStateMsg1(
				ConformanceRules_IsTaggedTemplateLitString(node),
				"Found untagged template literal with invalid (uncookable) escape sequence: %s",
				rawString)
			if r.stringPattern.Matcher(rawString).Matches() {
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

// ConformanceRules_ConformanceUtil holds the static helpers of the private
// class ConformanceUtil.
type ConformanceRules_ConformanceUtil struct{}

// Java: ConformanceRules.ConformanceUtil.isCallTarget(Node)
func ConformanceRules_ConformanceUtil_IsCallTarget(n rhino.Node) bool {
	role := ast.RoleOf(n)
	return role == ast.RoleCallCallee || role == ast.RoleNewCallee
}

// Java: ConformanceRules.ConformanceUtil.isLooseType(JSType)
func ConformanceRules_ConformanceUtil_IsLooseType(type_ rhino.JSType) bool {
	return type_.IsUnknownType() || type_.IsNoResolvedType() || type_.IsAllType()
}

// Java: ConformanceRules.ConformanceUtil.evaluateTypeString(AbstractCompiler, String)
func conformanceRules_evaluateTypeString(compiler AbstractCompiler, expression string) (rhino.JSType, error) {
	typeNodes := parsing.JsDocInfoParser_ParseTypeString(expression)
	if typeNodes.IsNil() {
		return nil, errInvalidRequirementSpec("bad type expression")
	}
	typeExpr := rhino.NewJSTypeExpression(typeNodes, "conformance")
	return compiler.GetTypeRegistry().EvaluateTypeExpressionInGlobalScope(typeExpr), nil
}

// ConformanceRules_ConformanceUtil_ValidateCall validates the parameters and
// the 'this' type, of a new or call.
//
// Java: ConformanceRules.ConformanceUtil.validateCall(AbstractCompiler, Node, FunctionType, boolean)
func ConformanceRules_ConformanceUtil_ValidateCall(
	compiler AbstractCompiler,
	callOrNew rhino.Node,
	functionType rhino.FunctionType,
	isCallInvocation bool) bool {
	cn := ast.MustCallOrNew(callOrNew)

	return conformanceRules_validateParameterList(compiler, cn, functionType, isCallInvocation) &&
		conformanceRules_validateThis(cn, functionType, isCallInvocation)
}

// Java: ConformanceRules.ConformanceUtil.validateThis(Node, FunctionType, boolean)
func conformanceRules_validateThis(
	callOrNew ast.CallOrNew, functionType rhino.FunctionType, isCallInvocation bool) bool {
	if _, isNew := callOrNew.AsNew(); isNew {
		return true
	}

	thisType := functionType.GetTypeOfThis()
	if thisType == nil || thisType.IsUnknownType() {
		return true
	}

	var thisNode rhino.Node
	if isCallInvocation {
		if firstArg, ok := callOrNew.FirstArg(); ok {
			thisNode = firstArg.Node()
		}
	} else if gp, ok := ast.AsGetProp(callOrNew.Callee().Node()); ok {
		thisNode = gp.Object().Node()
	} else {
		thisNode = callOrNew.Callee().Node().GetFirstChild()
	}
	thisNodeType := thisNode.GetJSType().RestrictByNotNullOrUndefined()
	return thisNodeType.IsSubtypeOf(thisType)
}

// Java: ConformanceRules.ConformanceUtil.validateParameterList(AbstractCompiler, Node, FunctionType, boolean)
func conformanceRules_validateParameterList(
	compiler AbstractCompiler,
	callOrNew ast.CallOrNew,
	functionType rhino.FunctionType,
	isCallInvocation bool) bool {
	var arg rhino.Node
	if firstArg, ok := callOrNew.FirstArg(); ok {
		arg = firstArg.Node()
	}
	if isCallInvocation && !arg.IsNil() {
		arg = arg.GetNext()
	}

	var argumentTypes []rhino.JSType
	for ; !arg.IsNil(); arg = arg.GetNext() {
		argType := arg.GetJSType()
		if argType == nil {
			argType = compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		argumentTypes = append(argumentTypes, argType)
	}
	return functionType.AcceptsArguments(argumentTypes)
}

func conformanceRules_splitClassAndProperty(specName string) (className string, propName string, ok bool) {
	parts := conformanceRules_splitToList(conformanceRules_ON_PROTOTYPE, specName)
	jutil.CheckState(len(parts) == 1 || len(parts) == 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}
	return "", "", false
}

// Java: ConformanceRules.ConformanceUtil.removeTypeDecl(String)
func conformanceRules_removeTypeDecl(specName string) (string, error) {
	index := strings.IndexByte(specName, ':')
	if index < 1 {
		return "", errInvalidRequirementSpec("value should be in the form NAME:TYPE")
	}
	return specName[:index], nil
}

// Java: ConformanceRules.ConformanceUtil.getTypeFromValue(String)
func conformanceRules_getTypeFromValue(specName string) (string, bool) {
	index := strings.IndexByte(specName, ':')
	if index < 1 {
		return "", false
	}
	return specName[index+1:], true
}

// Java: ConformanceRules.ConformanceUtil.isAnnotatedAsConst(JSDocInfo)
func ConformanceRules_ConformanceUtil_IsAnnotatedAsConst(jsdoc *rhino.JSDocInfo) bool {
	return jsdoc != nil && jsdoc.IsConstant()
}

// Java: ConformanceRules.ConformanceUtil.getValueFromGlobalName(Scope, GlobalNamespace, String)
func ConformanceRules_ConformanceUtil_GetValueFromGlobalName(
	scope *Scope, gns *GlobalNamespace, qualifiedName *string) *string {
	var gn *GlobalNamespace_Name
	if qualifiedName == nil {
		gns.EnsureGenerated()
	} else if slot := gns.GetOwnSlot(*qualifiedName); slot != nil {
		gn = slot.(*GlobalNamespace_Name)
	}
	if gn == nil ||
		gn.GetGlobalSets() != 1 ||
		gn.GetTotalSets() != 1 ||
		!ConformanceRules_ConformanceUtil_IsAnnotatedAsConst(gn.GetJSDocInfo()) {
		return nil
	}

	return ConformanceRules_ConformanceUtil_InferStringValue(
		scope, NodeUtil_GetRValueOfLValue(gn.GetDeclaration().GetNode()), func() *GlobalNamespace { return gns })
}

// Java: ConformanceRules.ConformanceUtil.inferStringValue(Scope, Node, Supplier<GlobalNamespace>)
func ConformanceRules_ConformanceUtil_InferStringValue(
	scope *Scope, node rhino.Node, globalNamespaceSupplier func() *GlobalNamespace) *string {
	if node.IsNil() {
		return nil
	}

	switch node.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT, rhino.Token_TEMPLATELIT_STRING:
		return NodeUtil_GetStringValue(node)
	case rhino.Token_NAME:
		if scope == nil {
			return nil
		}
		var_ := scope.GetVarByID(node.GetStringID())
		if var_ == nil {
			return nil
		}
		if !var_.IsConst() && !ConformanceRules_ConformanceUtil_IsAnnotatedAsConst(var_.GetJSDocInfo()) {
			return nil
		}
		initialValue := var_.GetInitialValue()
		return ConformanceRules_ConformanceUtil_InferStringValue(var_.GetScope().(*Scope), initialValue, globalNamespaceSupplier)
	case rhino.Token_GETPROP:
		gp := ast.MustGetProp(node)
		type_ := node.GetJSType()
		if type_ == nil {
			return nil
		}

		// The JS style guide requires enums to be effectively immutable and all enum items should
		// be statically known. See go/js-style#features-objects-enums.
		if type_.IsEnumElementType() {
			enumSource := type_.ToMaybeEnumElementType().GetEnumType().GetSource()
			if enumSource.IsNil() {
				return nil
			}
			if !enumSource.IsObjectLit() && !enumSource.IsClassMembers() {
				return nil
			}
			return ConformanceRules_ConformanceUtil_InferStringValue(
				nil,
				NodeUtil_GetFirstPropMatchingKey(enumSource, gp.Property()),
				globalNamespaceSupplier)
		} else if type_.IsString() {
			gns := globalNamespaceSupplier()
			if gns == nil {
				return nil
			}
			return ConformanceRules_ConformanceUtil_GetValueFromGlobalName(scope, gns, node.GetQualifiedName())
		}
		return nil
	default:
		// Do nothing
	}
	return nil
}

// Java: ConformanceRules.ConformanceUtil.isXid(JSType)
func ConformanceRules_ConformanceUtil_IsXid(type_ rhino.JSType) bool {
	if type_ == nil {
		return false
	}
	enumElTy := type_.ToMaybeEnumElementType()
	if enumElTy != nil &&
		*enumElTy.GetEnumType().GetReferenceName() == "enum{xid.String}" {
		return true
	}
	return false
}

// Java: ConformanceRules.ConformanceUtil.isEventHandlerAttrName(String)
func ConformanceRules_ConformanceUtil_IsEventHandlerAttrName(attr string) bool {
	return attr != "on" && strings.HasPrefix(attr, "on")
}

// ConformanceRules_RestrictedNameCall is the restricted name call rule.
type ConformanceRules_RestrictedNameCall struct {
	ConformanceRules_AbstractRule

	restrictions []ConformanceRules_RestrictedNameCall_Restriction
}

// Java: ConformanceRules.isLooseObject(JSType, JSTypeRegistry)
func ConformanceRules_IsLooseObject(type_ rhino.JSType, registry *rhino.JSTypeRegistry) bool {
	return type_.Equals(registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
}
