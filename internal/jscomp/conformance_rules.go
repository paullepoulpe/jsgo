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

// conformanceRules_splitToList is Splitter.on(separator).splitToList(s):
// every piece between non-overlapping occurrences of the (ASCII) separator,
// empty pieces included.
//
// Go-only: language-adapter — Guava Splitter.
func conformanceRules_splitToList(separator string, s string) []string {
	return strings.Split(s, separator)
}

// Java: ConformanceRules.createTsAllowlist()
func ConformanceRules_CreateTsAllowlist() *ConformanceRules_AllowList {
	var result *ConformanceRules_AllowList
	jutil.Try(func() {
		result = NewConformanceRules_AllowList(
			[]string{}, []string{".*\\.closure\\.js", ".*\\.tsx\\.cl\\.js"})
	}, func(t jutil.Throwable) {
		panic(jutil.NewAssertionErrorWithCause(t))
	})
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
	r := &ConformanceRules_ConformanceResult{}
	r.level = level
	r.note = note
	return r
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

// ConformanceRules_ConformanceLevel is the possible check check results.
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

// ConformanceRules_BuildPattern returns the joined pattern, or nil when
// reqPatterns is empty. It panics with InvalidRequirementSpec (a checked
// exception) when a pattern is invalid.
//
// Java: ConformanceRules.buildPattern(List<String>)
func ConformanceRules_BuildPattern(reqPatterns []string) *jutil.Pattern {
	if len(reqPatterns) == 0 {
		return nil
	}

	// validate the patterns
	for _, reqPattern := range reqPatterns {
		jutil.Try(func() {
			jutil.Pattern_Compile(reqPattern)
		}, func(e *jutil.PatternSyntaxException) {
			panic(NewCheckConformance_InvalidRequirementSpecWithCause("invalid regex pattern", e))
		})
	}

	var pattern *jutil.Pattern
	jutil.Try(func() {
		jointRegExp := "(" + strings.Join(reqPatterns, "|") + ")"
		pattern = jutil.Pattern_Compile(jointRegExp)
	}, func(e *jutil.PatternSyntaxException) {
		panic(jutil.NewRuntimeExceptionWithMessageAndCause("bad joined regexp", e))
	})
	return pattern
}

// Java: ConformanceRules.AllowList
type ConformanceRules_AllowList struct {
	prefixes       []string
	regexp         *jutil.Pattern
	allowlistEntry *RequirementScopeEntry
}

// Java: ConformanceRules.AllowList(List<String>, List<String>)
func NewConformanceRules_AllowList(prefixes []string, regexps []string) *ConformanceRules_AllowList {
	a := &ConformanceRules_AllowList{}
	a.prefixes = slices.Clone(prefixes)
	a.regexp = ConformanceRules_BuildPattern(regexps)
	a.allowlistEntry = nil
	return a
}

// Java: ConformanceRules.AllowList(RequirementScopeEntry)
func NewConformanceRules_AllowListFromEntry(allowlistEntry *RequirementScopeEntry) *ConformanceRules_AllowList {
	a := &ConformanceRules_AllowList{}
	a.prefixes = slices.Clone(allowlistEntry.GetPrefixList())
	a.regexp = ConformanceRules_BuildPattern(allowlistEntry.GetRegexpList())
	a.allowlistEntry = allowlistEntry
	return a
}

// Matches returns true if the given path matches one of the prefixes or
// regexps, and false otherwise.
//
// Java: ConformanceRules.AllowList.matches(String)
func (a *ConformanceRules_AllowList) Matches(path string) bool {
	// If the path ends with .closure.js or .tsx.cl.js, it is probably a tsickle-generated file,
	// and there may be entries in the allow list for the TypeScript path
	var tsPath *string
	if strings.HasSuffix(path, ".closure.js") {
		p := path[:len(path)-len(".closure.js")] + ".ts"
		tsPath = &p
	} else if strings.HasSuffix(path, ".tsx.cl.js") {
		// TSX
		p := path[:len(path)-len(".cl.js")]
		tsPath = &p
	}
	if a.prefixes != nil {
		for _, prefix := range a.prefixes {
			if path != "" &&
				(strings.HasPrefix(path, prefix) || (tsPath != nil && strings.HasPrefix(*tsPath, prefix))) {
				return true
			}
		}
	}

	return a.regexp != nil &&
		(a.regexp.Matcher(path).Find() || (tsPath != nil && a.regexp.Matcher(*tsPath).Find()))
}

// ConformanceRules_AbstractRule is a conformance rule implementation to
// support things common to all rules such as allowlisting and reporting.
type ConformanceRules_AbstractRule interface {
	CheckConformance_Rule

	// TsIsAllowlisted is protected and overridable (called from the
	// constructor).
	TsIsAllowlisted() bool

	// CheckConformance returns whether the code represented by the Node
	// conforms to the rule.
	CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult

	// Report reports a conformance warning for the given node.
	Report(n rhino.Node, result *ConformanceRules_ConformanceResult, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior)

	GetConformanceRules_AbstractRuleBase() *ConformanceRules_AbstractRuleBase
}

// ConformanceRules_AbstractRuleBase holds AbstractRule's fields and base
// implementations.
type ConformanceRules_AbstractRuleBase struct {
	CheckConformance_RuleDefaults

	self ConformanceRules_AbstractRule

	compiler                  AbstractCompiler
	message                   string
	severity                  Requirement_Severity
	allowlists                []*ConformanceRules_AllowList
	onlyApplyTo               *ConformanceRules_AllowList
	reportLooseTypeViolations bool
	typeMatchingStrategy      TypeMatchingStrategy
	requirement               *Requirement
}

// InitConformanceRules_AbstractRuleBase is the AbstractRule constructor. It
// panics with InvalidRequirementSpec (a checked exception).
//
// Java: ConformanceRules.AbstractRule(AbstractCompiler, Requirement)
func (b *ConformanceRules_AbstractRuleBase) InitConformanceRules_AbstractRuleBase(
	self ConformanceRules_AbstractRule, compiler AbstractCompiler, requirement *Requirement) {
	b.self = self
	if !requirement.HasErrorMessage() {
		panic(NewCheckConformance_InvalidRequirementSpec("missing message"))
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
	allowlistsBuilder := []*ConformanceRules_AllowList{}
	for _, entry := range requirement.GetWhitelistEntryList() {
		allowlistsBuilder = append(allowlistsBuilder, NewConformanceRules_AllowListFromEntry(entry))
	}
	for _, entry := range requirement.GetAllowlistEntryList() {
		allowlistsBuilder = append(allowlistsBuilder, NewConformanceRules_AllowListFromEntry(entry))
	}

	if b.self.TsIsAllowlisted() {
		allowlistsBuilder = append(allowlistsBuilder, conformanceRules_ALL_TS_ALLOWLIST)
	}

	if requirement.GetWhitelistCount() > 0 || requirement.GetWhitelistRegexpCount() > 0 {
		allowlist := NewConformanceRules_AllowList(requirement.GetWhitelistList(), requirement.GetWhitelistRegexpList())
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	if requirement.GetAllowlistCount() > 0 || requirement.GetAllowlistRegexpCount() > 0 {
		allowlist := NewConformanceRules_AllowList(requirement.GetAllowlistList(), requirement.GetAllowlistRegexpList())
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	b.allowlists = allowlistsBuilder

	if requirement.GetOnlyApplyToCount() > 0 || requirement.GetOnlyApplyToRegexpCount() > 0 {
		b.onlyApplyTo =
			NewConformanceRules_AllowList(requirement.GetOnlyApplyToList(), requirement.GetOnlyApplyToRegexpList())
	} else {
		b.onlyApplyTo = nil
	}
	b.reportLooseTypeViolations = requirement.GetReportLooseTypeViolations()
	b.typeMatchingStrategy = ConformanceRules_AbstractRule_GetTypeMatchingStrategy(requirement)
	b.requirement = requirement
}

// GetConformanceRules_AbstractRuleBase gives access to the base fields
// through the interface.
func (b *ConformanceRules_AbstractRuleBase) GetConformanceRules_AbstractRuleBase() *ConformanceRules_AbstractRuleBase {
	return b
}

// Java: ConformanceRules.AbstractRule.tsIsAllowlisted()
func (b *ConformanceRules_AbstractRuleBase) TsIsAllowlisted() bool {
	// LINT.IfChange
	return false
	// LINT.ThenChange(//depot/google3/java/com/google/javascript/modules/librarydepsconformancechecker/conformance_checker.go)
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
func (b *ConformanceRules_AbstractRuleBase) FindAllowListForPath(path string) *ConformanceRules_AllowList {
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

// Java: ConformanceRules.AbstractRule.check(NodeTraversal, Node, LibraryLevelNonAllowlistedConformanceViolationsBehavior)
func (b *ConformanceRules_AbstractRuleBase) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	result := b.self.CheckConformance(t, n)
	if result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		b.self.Report(n, result, behavior)
	}
}

// Report reports a conformance warning for the given node.
//
// n is the node representing the violating code; result the result
// representing the confidence of the violation.
//
// Java: ConformanceRules.AbstractRule.report(Node, ConformanceResult, LibraryLevelNonAllowlistedConformanceViolationsBehavior)
func (b *ConformanceRules_AbstractRuleBase) Report(
	n rhino.Node,
	result *ConformanceRules_ConformanceResult,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	var msg *DiagnosticType
	if b.severity == Requirement_Severity_ERROR {
		// Always report findings that are errors, even if the types are too loose to be certain.
		// TODO(bangert): If this causes problems, add another severity category that only
		// errors when certain.
		msg = CheckConformance_CONFORMANCE_ERROR
	} else {
		if result.level == ConformanceRules_ConformanceLevel_VIOLATION {
			msg = CheckConformance_CONFORMANCE_VIOLATION
		} else {
			msg = CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION
		}
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
		allowlistEntry = allowlist.allowlistEntry // Optional.fromNullable
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

	// if the violation is not in a gencode or we're not in library level reporting mode,
	// determined by `shouldReport` above, then check the allowlists to decide whether to actually
	// report the violation or not.
	if shouldReport && !isAllowlisted {
		b.compiler.Report(err)
	}
}

// ConformanceRules_NoOp is a no-op rule that never reports any violations.
//
// This exists so that, if a requirement becomes obsolete but is extended by
// other requirements that can't all be simultaneously deleted, it can be
// changed to this rule, allowing it to be effectively removed without
// breaking downstream builds.
type ConformanceRules_NoOp struct {
	ConformanceRules_AbstractRuleBase
}

// Java: ConformanceRules.NoOp(AbstractCompiler, Requirement)
func NewConformanceRules_NoOp(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_NoOp {
	r := &ConformanceRules_NoOp{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.NoOp.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_NoOp) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_AbstractTypeRestrictionRule is the abstract class
// AbstractTypeRestrictionRule.
type ConformanceRules_AbstractTypeRestrictionRule interface {
	ConformanceRules_AbstractRule

	GetConformanceRules_AbstractTypeRestrictionRuleBase() *ConformanceRules_AbstractTypeRestrictionRuleBase
}

// ConformanceRules_AbstractTypeRestrictionRuleBase holds
// AbstractTypeRestrictionRule's fields and base implementations.
type ConformanceRules_AbstractTypeRestrictionRuleBase struct {
	ConformanceRules_AbstractRuleBase

	nativeObjectType   rhino.JSType
	allowlistedTypes   rhino.JSType
	assertionFunctions *CodingConvention_AssertionFunctionLookup
}

// InitConformanceRules_AbstractTypeRestrictionRuleBase is the
// AbstractTypeRestrictionRule constructor.
//
// Java: ConformanceRules.AbstractTypeRestrictionRule(AbstractCompiler, Requirement)
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) InitConformanceRules_AbstractTypeRestrictionRuleBase(
	self ConformanceRules_AbstractTypeRestrictionRule, compiler AbstractCompiler, requirement *Requirement) {
	b.InitConformanceRules_AbstractRuleBase(self, compiler, requirement)
	b.nativeObjectType = compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
	allowlistedTypeNames := requirement.GetValueList()
	b.allowlistedTypes = b.Union(allowlistedTypeNames)

	b.assertionFunctions =
		CodingConvention_AssertionFunctionLookup_Of(compiler.GetCodingConvention().GetAssertionFunctions())
}

// GetConformanceRules_AbstractTypeRestrictionRuleBase gives access to the
// base fields through the interface.
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) GetConformanceRules_AbstractTypeRestrictionRuleBase() *ConformanceRules_AbstractTypeRestrictionRuleBase {
	return b
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isAllowlistedType(Node)
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) IsAllowlistedType(n rhino.Node) bool {
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
		!ConformanceRules_AbstractTypeRestrictionRule_IsTemplateType(n) // TODO(johnlenz): Remove this restriction
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isNativeObjectType(Node)
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) IsNativeObjectType(n rhino.Node) bool {
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
	return (type_ == nil || type_.IsUnknownType())
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
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) Union(typeNames []string) rhino.JSType {
	registry := b.compiler.GetTypeRegistry()
	types := []rhino.JSType{}

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
	} else {
		return registry.CreateUnionTypeOfList(types)
	}
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isAssertionCall(Node)
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) IsAssertionCall(n rhino.Node) bool {
	if n.IsCall() && n.GetFirstChild().IsQualifiedName() {
		target := n.GetFirstChild()
		return b.assertionFunctions.LookupByCallee(target) != nil
	}
	return false
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isTypeImmediatelyTightened(Node)
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) IsTypeImmediatelyTightened(n rhino.Node) bool {
	return b.IsAssertionCall(n.GetParent()) ||
		n.GetParent().IsTypeOf() ||
		/* casted node */ n.GetJSTypeBeforeCast() != nil
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isUsed(Node)
func (b *ConformanceRules_AbstractTypeRestrictionRuleBase) IsUsed(n rhino.Node) bool {
	if n.GetParent().IsName() || NodeUtil_IsLhsByDestructuring(n) {
		return false
	}

	// Consider lvalues in assignment operations to be used iff the actual assignment
	// operation's result is used. e.g. for `a.b.c`:
	//     USED: `alert(x = a.b.c);`
	//   UNUSED: `x = a.b.c;`
	if NodeUtil_IsAssignmentOp(n.GetParent()) {
		return NodeUtil_IsExpressionResultUsed(n.GetParent())
	}

	return NodeUtil_IsExpressionResultUsed(n)
}

// ConformanceRules_InferredConstCheck checks that variables annotated as
// @const have an inferred type, if there is no type given explicitly.
type ConformanceRules_InferredConstCheck struct {
	ConformanceRules_AbstractRuleBase
}

// Java: ConformanceRules.InferredConstCheck(AbstractCompiler, Requirement)
func NewConformanceRules_InferredConstCheck(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_InferredConstCheck {
	r := &ConformanceRules_InferredConstCheck{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	return r
}

// Java: ConformanceRules.InferredConstCheck.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_InferredConstCheck) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	jsDoc := n.GetJSDocInfo()
	if jsDoc != nil && jsDoc.HasConstAnnotation() && jsDoc.GetType() == nil {
		if n.IsAssign() {
			n = n.GetFirstChild()
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
	ConformanceRules_AbstractRuleBase

	paths []string
}

// Java: ConformanceRules.BannedDependency(AbstractCompiler, Requirement)
func NewConformanceRules_BannedDependency(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BannedDependency {
	r := &ConformanceRules_BannedDependency{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	r.paths = requirement.GetValueList()
	if len(r.paths) == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}
	return r
}

// Java: ConformanceRules.BannedDependency.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedDependency) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsScript() {
		srcFile := *n.GetSourceFileName()
		for i := 0; i < len(r.paths); i++ {
			path := r.paths[i]
			if strings.HasPrefix(srcFile, path) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BannedDependencyRegex is the banned dependency via regex
// rule.
type ConformanceRules_BannedDependencyRegex struct {
	ConformanceRules_AbstractRuleBase

	pathRegexp *jutil.Pattern
}

// Java: ConformanceRules.BannedDependencyRegex(AbstractCompiler, Requirement)
func NewConformanceRules_BannedDependencyRegex(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BannedDependencyRegex {
	r := &ConformanceRules_BannedDependencyRegex{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	pathRegexpList := requirement.GetValueList()
	if len(pathRegexpList) == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value (no banned dependency regexps)"))
	}
	r.pathRegexp = ConformanceRules_BuildPattern(pathRegexpList)
	return r
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
	ConformanceRules_AbstractRuleBase

	// LINT.ThenChange(//depot/google3/javascript/typescript/compiler/conformance/rules/ban_third_party_script.ts)
	requirementType Requirement_Type
	qualifiedNames  []rhino.Node
	shortNames      *jutil.LinkedHashSet[string]
}

// Java: ConformanceRules.BannedName(AbstractCompiler, Requirement)
func NewConformanceRules_BannedName(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BannedName {
	r := &ConformanceRules_BannedName{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)
	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}
	r.requirementType = requirement.GetType()

	qualifiedBuilder := []rhino.Node{}
	shortBuilder := jutil.NewLinkedHashSet[string]()
	for _, name := range requirement.GetValueList() {
		qualifiedName := NodeUtil_NewQName(compiler, name)
		qualifiedBuilder = append(qualifiedBuilder, qualifiedName)
		shortBuilder.Add(qualifiedName.GetString())
	}
	r.qualifiedNames = qualifiedBuilder
	r.shortNames = shortBuilder
	return r
}

// Java: ConformanceRules.BannedName.IS_CANDIDATE_NODE
var conformanceRules_BannedName_IS_CANDIDATE_NODE CheckConformance_Precondition = NewCheckConformance_PreconditionFunc(
	func(n rhino.Node) bool {
		switch n.GetToken() {
		case rhino.Token_GETPROP:
			return n.GetFirstChild().IsQualifiedName()
		case rhino.Token_NAME:
			return n.GetString() != ""
		default:
			return false
		}
	})

// Java: ConformanceRules.BannedName.getPrecondition()
func (r *ConformanceRules_BannedName) GetPrecondition() CheckConformance_Precondition {
	return conformanceRules_BannedName_IS_CANDIDATE_NODE
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
	if !r.shortNames.Contains(n.GetString()) {
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
	rootName := NodeUtil_GetRootOfQualifiedName(n).GetQualifiedName()
	v := t.GetScope().GetVar(*rootName)
	// TODO(b/189382837): Turn the nullness check back into an assertion once the bug is fixed.
	return v != nil && v.IsGlobal()
}

// ConformanceRules_BannedProperty is the banned property rule.
type ConformanceRules_BannedProperty struct {
	ConformanceRules_AbstractRuleBase

	registry *rhino.JSTypeRegistry
	// Go-only: language-adapter — ImmutableSetMultimap<String, JSType>: the
	// values of a key are a set of JSTypes deduplicated with equals(), in
	// insertion order.
	props                   *jutil.LinkedHashMap[string, *jutil.ObjectLinkedHashSet[rhino.JSType]]
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
		if NodeUtil_IsLhsOfAssign(n) &&
			(NodeUtil_IsLiteralValue(n.GetNext() /* includeFunctions= */, false) ||
				NodeUtil_IsSomeCompileTimeConstStringValue(n.GetNext())) {
			return false
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
func NewConformanceRules_BannedProperty(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BannedProperty {
	r := &ConformanceRules_BannedProperty{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)

	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
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

	builder := jutil.NewLinkedHashMap[string, *jutil.ObjectLinkedHashSet[rhino.JSType]]()
	for _, value := range requirement.GetValueList() {
		typename := ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(value)
		property := ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(value)
		if typename == nil || property == nil {
			panic(NewCheckConformance_InvalidRequirementSpec("bad prop value"))
		}

		type_ := r.registry.GetGlobalType(*typename)
		if type_ == nil {
			type_ = r.registry.ResolveViaClosureNamespace(*typename)
		}

		// If type doesn't exist in the copmilation, it can't be a violation. Also, bottom types
		// match everything, which is almost surely not what the check is intended to check against.
		if type_ == nil || type_.IsUnknownType() || type_.IsEmptyType() {
			continue
		}

		types, ok := builder.Lookup(*property)
		if !ok {
			types = jutil.NewObjectLinkedHashSet[rhino.JSType]()
			builder.Put(*property, types)
		}
		types.Add(type_)
	}
	r.props = builder
	return r
}

// Java: ConformanceRules.BannedProperty.getPrecondition()
func (r *ConformanceRules_BannedProperty) GetPrecondition() CheckConformance_Precondition {
	return r.requirementPrecondition
}

// Java: ConformanceRules.BannedProperty.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedProperty) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	// ImmutableSetMultimap.get(null) is empty.
	name := r.ExtractName(n)
	var checkTypes *jutil.ObjectLinkedHashSet[rhino.JSType]
	if name != nil {
		checkTypes = r.props.Get(*name)
	}
	if checkTypes == nil || checkTypes.IsEmpty() {
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

	for checkType := range checkTypes.All() {
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
	switch n.GetToken() {
	case rhino.Token_GETELEM, rhino.Token_GETPROP:
		return n.GetFirstChild().GetJSType()
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
func (r *ConformanceRules_BannedProperty) ExtractName(n rhino.Node) *string {

	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_STRING_KEY:
		s := n.GetString()
		return &s
	case rhino.Token_GETELEM:
		string_ := n.GetSecondChild()
		if string_.IsStringLit() {
			s := string_.GetString()
			return &s
		}
		return nil
	case rhino.Token_COMPUTED_PROP:
		string_ := n.GetFirstChild()
		if string_.IsStringLit() {
			s := string_.GetString()
			return &s
		}
		return nil
	default:
		return nil
	}
}

// ConformanceRules_BannedStringRegex is the banned string rule.
type ConformanceRules_BannedStringRegex struct {
	ConformanceRules_AbstractRuleBase

	stringPattern *jutil.Pattern
}

// Java: ConformanceRules.BannedStringRegex(AbstractCompiler, Requirement)
func NewConformanceRules_BannedStringRegex(compiler AbstractCompiler, requirement *Requirement) *ConformanceRules_BannedStringRegex {
	r := &ConformanceRules_BannedStringRegex{}
	r.InitConformanceRules_AbstractRuleBase(r, compiler, requirement)

	if requirement.GetValueCount() == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}

	builder := []string{}
	for _, value := range requirement.GetValueList() {
		if jutil.StringTrim(value) == "" {
			panic(NewCheckConformance_InvalidRequirementSpec("empty strings or whitespace are not allowed"))
		}
		builder = append(builder, value)
	}

	values := builder
	if len(values) == 0 {
		panic(NewCheckConformance_InvalidRequirementSpec("missing value"))
	}
	stringRegex := ConformanceRules_BuildPattern(values)
	r.stringPattern = stringRegex
	return r
}

// Java: ConformanceRules.BannedStringRegex.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BannedStringRegex) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if node.IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if node.IsStringLit() {
		if r.stringPattern.Matcher(node.GetString()).Matches() {
			return ConformanceRules_ConformanceResult_VIOLATION
		}
	} else if node.IsTemplateLitString() {
		cookedString := node.GetCookedString()
		if cookedString != nil {
			if r.stringPattern.Matcher(*cookedString).Matches() {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		} else {
			jutil.CheckStateMsg1(
				ConformanceRules_IsTaggedTemplateLitString(node),
				"Found untagged template literal with invalid (uncookable) escape sequence: %s",
				node.GetRawString())
			rawString := node.GetRawString()
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
	return !node.GetGrandparent().IsNil() && node.GetGrandparent().IsTaggedTemplateLit()
}

// ConformanceRules_ConformanceUtil holds the static helpers of the private
// class ConformanceUtil.
type ConformanceRules_ConformanceUtil struct{}

// Java: ConformanceRules.ConformanceUtil.isCallTarget(Node)
func ConformanceRules_ConformanceUtil_IsCallTarget(n rhino.Node) bool {
	parent := n.GetParent()
	return (parent.IsCall() || parent.IsNew()) && parent.GetFirstChild() == n
}

// Java: ConformanceRules.ConformanceUtil.isLooseType(JSType)
func ConformanceRules_ConformanceUtil_IsLooseType(type_ rhino.JSType) bool {
	return type_.IsUnknownType() || type_.IsNoResolvedType() || type_.IsAllType()
}

// ConformanceRules_ConformanceUtil_EvaluateTypeString panics with
// InvalidRequirementSpec (a checked exception).
//
// Java: ConformanceRules.ConformanceUtil.evaluateTypeString(AbstractCompiler, String)
func ConformanceRules_ConformanceUtil_EvaluateTypeString(compiler AbstractCompiler, expression string) rhino.JSType {
	typeNodes := parsing.JsDocInfoParser_ParseTypeString(expression)
	if typeNodes.IsNil() {
		panic(NewCheckConformance_InvalidRequirementSpec("bad type expression"))
	}
	typeExpr := rhino.NewJSTypeExpression(typeNodes, "conformance")
	return compiler.GetTypeRegistry().EvaluateTypeExpressionInGlobalScope(typeExpr)
}

// ConformanceRules_ConformanceUtil_ValidateCall validates the parameters and
// the 'this' type, of a new or call.
//
// See TypeCheck#visitParameterList.
//
// Java: ConformanceRules.ConformanceUtil.validateCall(AbstractCompiler, Node, FunctionType, boolean)
func ConformanceRules_ConformanceUtil_ValidateCall(
	compiler AbstractCompiler,
	callOrNew rhino.Node,
	functionType rhino.FunctionType,
	isCallInvocation bool) bool {
	jutil.CheckState(callOrNew.IsCall() || callOrNew.IsNew())

	return ConformanceRules_ConformanceUtil_ValidateParameterList(compiler, callOrNew, functionType, isCallInvocation) &&
		ConformanceRules_ConformanceUtil_ValidateThis(callOrNew, functionType, isCallInvocation)
}

// Java: ConformanceRules.ConformanceUtil.validateThis(Node, FunctionType, boolean)
func ConformanceRules_ConformanceUtil_ValidateThis(
	callOrNew rhino.Node, functionType rhino.FunctionType, isCallInvocation bool) bool {

	if callOrNew.IsNew() {
		return true
	}

	thisType := functionType.GetTypeOfThis()
	if thisType == nil || thisType.IsUnknownType() {
		return true
	}

	var thisNode rhino.Node
	if isCallInvocation {
		thisNode = callOrNew.GetSecondChild()
	} else {
		thisNode = callOrNew.GetFirstFirstChild()
	}
	thisNodeType := thisNode.GetJSType().RestrictByNotNullOrUndefined()
	return thisNodeType.IsSubtypeOf(thisType)
}

// Java: ConformanceRules.ConformanceUtil.validateParameterList(AbstractCompiler, Node, FunctionType, boolean)
func ConformanceRules_ConformanceUtil_ValidateParameterList(
	compiler AbstractCompiler,
	callOrNew rhino.Node,
	functionType rhino.FunctionType,
	isCallInvocation bool) bool {
	arg := callOrNew.GetSecondChild() // skip the function name
	if isCallInvocation && !arg.IsNil() {
		arg = arg.GetNext()
	}

	// Get all the annotated types of the argument nodes
	argumentTypes := []rhino.JSType{}
	for ; !arg.IsNil(); arg = arg.GetNext() {
		argType := arg.GetJSType()
		if argType == nil {
			argType = compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		argumentTypes = append(argumentTypes, argType)
	}
	return functionType.AcceptsArguments(argumentTypes)
}

// ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName extracts
// the method name from a provided name.
//
// Java: ConformanceRules.ConformanceUtil.getPropertyFromDeclarationName(String)
func ConformanceRules_ConformanceUtil_GetPropertyFromDeclarationName(specName string) *string {
	parts := conformanceRules_splitToList(conformanceRules_ON_PROTOTYPE, specName)
	jutil.CheckState(len(parts) == 1 || len(parts) == 2)
	if len(parts) == 2 {
		return &parts[1]
	}
	return nil
}

// ConformanceRules_ConformanceUtil_GetClassFromDeclarationName extracts the
// class name from a provided name.
//
// Java: ConformanceRules.ConformanceUtil.getClassFromDeclarationName(String)
func ConformanceRules_ConformanceUtil_GetClassFromDeclarationName(specName string) *string {
	tmp := specName
	parts := conformanceRules_splitToList(conformanceRules_ON_PROTOTYPE, tmp)
	jutil.CheckState(len(parts) == 1 || len(parts) == 2)
	if len(parts) == 2 {
		return &parts[0]
	}
	return nil
}

// ConformanceRules_ConformanceUtil_RemoveTypeDecl panics with
// InvalidRequirementSpec (a checked exception).
//
// Java: ConformanceRules.ConformanceUtil.removeTypeDecl(String)
func ConformanceRules_ConformanceUtil_RemoveTypeDecl(specName string) string {
	index := jutil.StringIndexOfChar(specName, ':')
	if index < 1 {
		panic(NewCheckConformance_InvalidRequirementSpec("value should be in the form NAME:TYPE"))
	}
	return jutil.StringSubstring(specName, 0, index)
}

// Java: ConformanceRules.ConformanceUtil.getTypeFromValue(String)
func ConformanceRules_ConformanceUtil_GetTypeFromValue(specName string) *string {
	index := jutil.StringIndexOfChar(specName, ':')
	if index < 1 {
		return nil
	}
	s := jutil.StringSubstringFrom(specName, index+1)
	return &s
}

// Java: ConformanceRules.ConformanceUtil.isAnnotatedAsConst(JSDocInfo)
func ConformanceRules_ConformanceUtil_IsAnnotatedAsConst(jsdoc *rhino.JSDocInfo) bool {
	return jsdoc != nil && jsdoc.IsConstant()
}

// Java: ConformanceRules.ConformanceUtil.getValueFromGlobalName(Scope, GlobalNamespace, String)
//
// qualifiedName is nullable (Node.getQualifiedName() of a GETPROP that is
// not a qualified name): getOwnSlot(null) generates the namespace and finds
// no name.
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
		name := node.GetString()
		var_ := scope.GetVar(name)
		if var_ == nil {
			return nil
		}
		if !var_.IsConst() && !ConformanceRules_ConformanceUtil_IsAnnotatedAsConst(var_.GetJSDocInfo()) {
			return nil
		}
		initialValue := var_.GetInitialValue()
		return ConformanceRules_ConformanceUtil_InferStringValue(var_.GetScope().(*Scope), initialValue, globalNamespaceSupplier)
	case rhino.Token_GETPROP:
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
				NodeUtil_GetFirstPropMatchingKey(enumSource, node.GetString()),
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
	ConformanceRules_AbstractRuleBase

	restrictions []*ConformanceRules_RestrictedNameCall_Restriction
}

// Java: ConformanceRules.isLooseObject(JSType, JSTypeRegistry)
func ConformanceRules_IsLooseObject(type_ rhino.JSType, registry *rhino.JSTypeRegistry) bool {
	return type_.Equals(registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
}
