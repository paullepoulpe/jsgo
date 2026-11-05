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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var allTSAllowlist = createTsAllowlist()

func createTsAllowlist() *conformanceAllowList {
	result, err := newConformanceAllowList(
		[]string{}, []string{".*\\.closure\\.js", ".*\\.tsx\\.cl\\.js"})
	if err != nil {
		panic(jutil.NewAssertionError(err.Error()))
	}
	return result
}

// ConformanceRules_ConformanceResult represents the outcome of checking a node
// against a conformance rule, optionally carrying a violation-specific note.
type ConformanceRules_ConformanceResult struct {
	level ConformanceRules_ConformanceLevel
	note  string
}

func NewConformanceRules_ConformanceResult(level ConformanceRules_ConformanceLevel) *ConformanceRules_ConformanceResult {
	return NewConformanceRules_ConformanceResultWithNote(level, "")
}

func NewConformanceRules_ConformanceResultWithNote(level ConformanceRules_ConformanceLevel, note string) *ConformanceRules_ConformanceResult {
	return &ConformanceRules_ConformanceResult{level: level, note: note}
}

// IsViolation reports whether r represents a definite or possible conformance violation.
func (r *ConformanceRules_ConformanceResult) IsViolation() bool {
	return r != nil && r.level != ConformanceRules_ConformanceLevel_CONFORMANCE
}

// GetLevel returns the conformance level of r.
func (r *ConformanceRules_ConformanceResult) GetLevel() ConformanceRules_ConformanceLevel {
	if r == nil {
		return ConformanceRules_ConformanceLevel_CONFORMANCE
	}
	return r.level
}

// GetNote returns the optional violation-specific note attached to r.
func (r *ConformanceRules_ConformanceResult) GetNote() string {
	if r == nil {
		return ""
	}
	return r.note
}

var ConformanceRules_ConformanceResult_CONFORMANCE = NewConformanceRules_ConformanceResult(ConformanceRules_ConformanceLevel_CONFORMANCE)

var ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION = NewConformanceRules_ConformanceResult(ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION)

var conformancePossibleViolationLooseTypes = NewConformanceRules_ConformanceResultWithNote(
	ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION,
	"The type information available for this expression is too loose "+
		"to ensure conformance.")

var ConformanceRules_ConformanceResult_VIOLATION = NewConformanceRules_ConformanceResult(ConformanceRules_ConformanceLevel_VIOLATION)

// ConformanceRules_ConformanceLevel enumerates the possible check outcomes.
type ConformanceRules_ConformanceLevel int32

const (
	// Nothing interesting detected.
	ConformanceRules_ConformanceLevel_CONFORMANCE ConformanceRules_ConformanceLevel = iota
	// In the optionally typed world of the Closure Compiler type system
	// it is possible that patterns match with looser types than the target pattern.
	ConformanceRules_ConformanceLevel_POSSIBLE_VIOLATION
	// Definitely a violation.
	ConformanceRules_ConformanceLevel_VIOLATION
)

var conformanceLevelNames = [...]string{"CONFORMANCE", "POSSIBLE_VIOLATION", "VIOLATION"}

func (l ConformanceRules_ConformanceLevel) String() string {
	return conformanceLevelNames[l]
}

func buildConformancePattern(reqPatterns []string) (*jutil.Pattern, error) {
	if len(reqPatterns) == 0 {
		return nil, nil
	}

	for _, reqPattern := range reqPatterns {
		if !jutil.IsValidPattern(reqPattern) {
			return nil, errInvalidRequirementSpec("invalid regex pattern")
		}
	}

	if len(reqPatterns) == 1 {
		return jutil.Pattern_Compile(reqPatterns[0]), nil
	}
	return jutil.Pattern_Compile("(" + strings.Join(reqPatterns, "|") + ")"), nil
}

type conformanceAllowList struct {
	prefixes       []string
	regexp         *jutil.Pattern
	allowlistEntry *RequirementScopeEntry
}

func newConformanceAllowList(prefixes []string, regexps []string) (*conformanceAllowList, error) {
	regexp, err := buildConformancePattern(regexps)
	if err != nil {
		return nil, err
	}
	return &conformanceAllowList{
		prefixes: slices.Clone(prefixes),
		regexp:   regexp,
	}, nil
}

func newConformanceAllowListFromEntry(allowlistEntry *RequirementScopeEntry) (*conformanceAllowList, error) {
	regexp, err := buildConformancePattern(allowlistEntry.GetRegexpList())
	if err != nil {
		return nil, err
	}
	return &conformanceAllowList{
		prefixes:       slices.Clone(allowlistEntry.GetPrefixList()),
		regexp:         regexp,
		allowlistEntry: allowlistEntry,
	}, nil
}

func splitConformanceTsPath(path string) (tsPath string, hasTsPath bool) {
	if stem, ok := strings.CutSuffix(path, ".closure.js"); ok {
		return stem + ".ts", true
	}
	if strings.HasSuffix(path, ".tsx.cl.js") {
		return strings.TrimSuffix(path, ".cl.js"), true
	}
	return "", false
}

// matches reports whether path (or its original `.ts`/`.tsx` counterpart for
// tsickle-generated files) matches one of the allowlist prefixes or regexps.
func (a *conformanceAllowList) matches(path string) bool {
	tsPath, hasTsPath := splitConformanceTsPath(path)
	return a.matchesWithTsPath(path, tsPath, hasTsPath)
}

func (a *conformanceAllowList) matchesWithTsPath(path, tsPath string, hasTsPath bool) bool {
	if len(a.prefixes) != 0 && path != "" {
		for _, prefix := range a.prefixes {
			if strings.HasPrefix(path, prefix) || (hasTsPath && strings.HasPrefix(tsPath, prefix)) {
				return true
			}
		}
	}

	return a.regexp != nil &&
		(a.regexp.Find(path) || (hasTsPath && a.regexp.Find(tsPath)))
}

func isDefaultConformanceErrorManager(em ErrorManager) bool {
	switch m := em.(type) {
	case *SortingErrorManager, *PrintStreamErrorManager, *LoggerErrorManager, *BlackHoleErrorManager:
		return true
	case *ThreadSafeDelegatingErrorManager:
		return m.hasDefaultConformanceReporting()
	default:
		return false
	}
}

func hasDefaultConformanceReporting(compiler AbstractCompiler) bool {
	if compiler == nil || compiler.GetConformanceViolationInterceptor() != nil {
		return false
	}
	return isDefaultConformanceErrorManager(compiler.GetErrorManager())
}

// ConformanceRules_AbstractRule holds the common rule fields and helpers shared
// by built-in, type-restriction, and embedder-registered custom conformance rules.
//
// Java: ConformanceRules.AbstractRule
type ConformanceRules_AbstractRule struct {
	compiler                  AbstractCompiler
	message                   string
	severity                  Requirement_Severity
	allowlists                []*conformanceAllowList
	combinedPrefixes          []string
	combinedRegexps           []*jutil.Pattern
	onlyApplyTo               *conformanceAllowList
	reportLooseTypeViolations bool
	typeMatchingStrategy      TypeMatchingStrategy
	requirement               *Requirement
	checkConformance          func(*NodeTraversal, rhino.Node) *ConformanceRules_ConformanceResult

	// Cached path normalization for the most recently inspected source file.
	cachedRawPath       string
	cachedRawTsPath     string
	cachedNormPath      string
	cachedNormTsPath    string
	hasCachedPath       bool
	hasCachedRawTsPath  bool
	hasCachedNormTsPath bool

	// Type-restriction rule fields (populated by initTypeRule).
	nativeObjectType   rhino.JSType
	allowlistedTypes   rhino.JSType
	assertionFunctions *CodingConvention_AssertionFunctionLookup
}

type conformanceAbstractRule = ConformanceRules_AbstractRule

func (b *ConformanceRules_AbstractRule) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := b.checkConformance(t, n); result.IsViolation() {
		b.Report(n, result, behavior)
	}
}

func (*ConformanceRules_AbstractRule) GetPrecondition() CheckConformance_Precondition {
	return CheckConformance_Precondition_CHECK_ALL
}

// GetCompiler returns the compiler associated with this rule.
func (b *ConformanceRules_AbstractRule) GetCompiler() AbstractCompiler { return b.compiler }

// GetRequirement returns the Requirement proto backing this rule.
func (b *ConformanceRules_AbstractRule) GetRequirement() *Requirement { return b.requirement }

// GetMessage returns the formatted diagnostic message for this rule.
func (b *ConformanceRules_AbstractRule) GetMessage() string { return b.message }

// GetSeverity returns the effective severity of this rule.
func (b *ConformanceRules_AbstractRule) GetSeverity() Requirement_Severity { return b.severity }

// InitRule initializes the common AbstractRule state and node-check callback.
func (b *ConformanceRules_AbstractRule) InitRule(
	compiler AbstractCompiler,
	requirement *Requirement,
	tsIsAllowlisted bool,
	check func(*NodeTraversal, rhino.Node) *ConformanceRules_ConformanceResult,
) error {
	return b.initRule(compiler, requirement, tsIsAllowlisted, check)
}

// InitTypeRule initializes AbstractTypeRestrictionRule state and node-check callback.
func (b *ConformanceRules_AbstractRule) InitTypeRule(
	compiler AbstractCompiler,
	requirement *Requirement,
	tsIsAllowlisted bool,
	check func(*NodeTraversal, rhino.Node) *ConformanceRules_ConformanceResult,
) error {
	return b.initTypeRule(compiler, requirement, tsIsAllowlisted, check)
}

// initRule initializes the common AbstractRule state and check callback.
func (b *conformanceAbstractRule) initRule(
	compiler AbstractCompiler,
	requirement *Requirement,
	tsIsAllowlisted bool,
	check func(*NodeTraversal, rhino.Node) *ConformanceRules_ConformanceResult,
) error {
	b.checkConformance = check
	if !requirement.HasErrorMessage() {
		return errInvalidRequirementSpec("missing message")
	}
	b.compiler = compiler
	message := requirement.GetErrorMessage()
	if len(requirement.GetConfigFileList()) > 0 {
		message +=
			"\n  defined in " + strings.Join(requirement.GetConfigFileList(), "\n  extended by ")
	}
	b.message = message
	if requirement.GetSeverity() == Requirement_Severity_UNSPECIFIED {
		b.severity = Requirement_Severity_WARNING
	} else {
		b.severity = requirement.GetSeverity()
	}

	var allowlistsBuilder []*conformanceAllowList
	for _, entry := range requirement.GetWhitelistEntryList() {
		allowlist, err := newConformanceAllowListFromEntry(entry)
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	for _, entry := range requirement.GetAllowlistEntryList() {
		allowlist, err := newConformanceAllowListFromEntry(entry)
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}

	if tsIsAllowlisted {
		allowlistsBuilder = append(allowlistsBuilder, allTSAllowlist)
	}

	if len(requirement.GetWhitelistList()) > 0 || len(requirement.GetWhitelistRegexpList()) > 0 {
		allowlist, err := newConformanceAllowList(requirement.GetWhitelistList(), requirement.GetWhitelistRegexpList())
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	if len(requirement.GetAllowlistList()) > 0 || len(requirement.GetAllowlistRegexpList()) > 0 {
		allowlist, err := newConformanceAllowList(requirement.GetAllowlistList(), requirement.GetAllowlistRegexpList())
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	b.allowlists = allowlistsBuilder
	var combinedPrefixes []string
	var combinedRegexps []*jutil.Pattern
	for _, al := range allowlistsBuilder {
		combinedPrefixes = append(combinedPrefixes, al.prefixes...)
		if al.regexp != nil {
			combinedRegexps = append(combinedRegexps, al.regexp)
		}
	}
	b.combinedPrefixes = combinedPrefixes
	b.combinedRegexps = combinedRegexps

	if len(requirement.GetOnlyApplyToList()) > 0 || len(requirement.GetOnlyApplyToRegexpList()) > 0 {
		onlyApplyTo, err := newConformanceAllowList(requirement.GetOnlyApplyToList(), requirement.GetOnlyApplyToRegexpList())
		if err != nil {
			return err
		}
		b.onlyApplyTo = onlyApplyTo
	} else {
		b.onlyApplyTo = nil
	}
	b.reportLooseTypeViolations = requirement.GetReportLooseTypeViolations()
	b.typeMatchingStrategy = conformanceTypeMatchingStrategy(requirement)
	b.requirement = requirement
	return nil
}

// initTypeRule initializes AbstractTypeRestrictionRule state and check callback.
func (b *conformanceAbstractRule) initTypeRule(
	compiler AbstractCompiler,
	requirement *Requirement,
	tsIsAllowlisted bool,
	check func(*NodeTraversal, rhino.Node) *ConformanceRules_ConformanceResult,
) error {
	if err := b.initRule(compiler, requirement, tsIsAllowlisted, check); err != nil {
		return err
	}
	b.nativeObjectType = compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
	b.allowlistedTypes = b.unionTypes(requirement.GetValueList())
	b.assertionFunctions =
		CodingConvention_AssertionFunctionLookup_Of(compiler.GetCodingConvention().GetAssertionFunctions())
	return nil
}

func conformanceTypeMatchingStrategy(requirement *Requirement) TypeMatchingStrategy {
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

func computeConformanceNormPath(compiler AbstractCompiler, rawPath string) string {
	if compiler != nil {
		if opts := compiler.GetOptions(); opts != nil {
			if pathRegex := opts.GetConformanceRemoveRegexFromPath(); pathRegex != nil {
				return pathRegex.Matcher(rawPath).ReplaceFirst("")
			}
		}
	}
	return rawPath
}

func (b *conformanceAbstractRule) ensureCachedPaths(rawPath string) {
	if b.hasCachedPath && b.cachedRawPath == rawPath {
		return
	}
	rawTsPath, hasRawTsPath := splitConformanceTsPath(rawPath)
	normPath := computeConformanceNormPath(b.compiler, rawPath)
	normTsPath, hasNormTsPath := rawTsPath, hasRawTsPath
	if normPath != rawPath {
		normTsPath, hasNormTsPath = splitConformanceTsPath(normPath)
	}
	b.cachedRawPath = rawPath
	b.cachedRawTsPath = rawTsPath
	b.cachedNormPath = normPath
	b.cachedNormTsPath = normTsPath
	b.hasCachedRawTsPath = hasRawTsPath
	b.hasCachedNormTsPath = hasNormTsPath
	b.hasCachedPath = true
}

func (b *conformanceAbstractRule) matchesAnyAllowlist(normPath, normTsPath string, hasNormTsPath bool) bool {
	if len(b.combinedPrefixes) != 0 && normPath != "" {
		for _, prefix := range b.combinedPrefixes {
			if strings.HasPrefix(normPath, prefix) || (hasNormTsPath && strings.HasPrefix(normTsPath, prefix)) {
				return true
			}
		}
	}
	for _, re := range b.combinedRegexps {
		if re.Find(normPath) || (hasNormTsPath && re.Find(normTsPath)) {
			return true
		}
	}
	return false
}

func (b *conformanceAbstractRule) isExcludedOrAllowlisted(
	rawPath, rawTsPath string,
	hasRawTsPath bool,
	normPath, normTsPath string,
	hasNormTsPath bool,
	hasPath bool,
) bool {
	if b.onlyApplyTo != nil && (!hasPath || !b.onlyApplyTo.matchesWithTsPath(rawPath, rawTsPath, hasRawTsPath)) {
		return true
	}
	return hasPath && b.matchesAnyAllowlist(normPath, normTsPath, hasNormTsPath)
}

func (b *conformanceAbstractRule) findAllowListForPath(rawPath string) *conformanceAllowList {
	b.ensureCachedPaths(rawPath)
	for _, allowlist := range b.allowlists {
		if allowlist.matchesWithTsPath(b.cachedNormPath, b.cachedNormTsPath, b.hasCachedNormTsPath) {
			return allowlist
		}
	}
	return nil
}

// Report reports a conformance finding for node n if it is not allowlisted.
func (b *conformanceAbstractRule) Report(
	n rhino.Node,
	result *ConformanceRules_ConformanceResult,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if hasDefaultConformanceReporting(b.compiler) {
		if behavior == ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_RECORD_ONLY {
			return
		}
		sourcePath, hasPath := NodeUtil_SourceName(n)
		if hasPath {
			b.ensureCachedPaths(sourcePath)
			if b.isExcludedOrAllowlisted(
				sourcePath, b.cachedRawTsPath, b.hasCachedRawTsPath,
				b.cachedNormPath, b.cachedNormTsPath, b.hasCachedNormTsPath,
				true,
			) {
				return
			}
		} else if b.onlyApplyTo != nil {
			return
		}
		var msg *DiagnosticType
		if b.severity == Requirement_Severity_ERROR {
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
		b.compiler.Report(JSError_MakeWithRequirement(b.requirement, n, msg, b.message, separator, result.note))
		return
	}

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

	var allowlist *conformanceAllowList
	sourcePath, hasPath := NodeUtil_SourceName(n)
	if hasPath {
		allowlist = b.findAllowListForPath(sourcePath)
	}
	isAllowlisted :=
		allowlist != nil || (b.onlyApplyTo != nil && (!hasPath || !b.onlyApplyTo.matchesWithTsPath(sourcePath, b.cachedRawTsPath, b.hasCachedRawTsPath)))
	var allowlistEntry *RequirementScopeEntry
	if allowlist != nil {
		allowlistEntry = allowlist.allowlistEntry
	}
	var shouldReport bool
	if interceptor := b.compiler.GetConformanceViolationInterceptor(); interceptor != nil {
		shouldReport = interceptor(
			b.requirement,
			allowlistEntry,
			err,
			behavior,
			isAllowlisted)
	} else {
		shouldReport = b.compiler.
			GetErrorManager().
			// returns true even if the violation is allowlisted
			ShouldReportConformanceViolation(
				b.requirement,
				allowlistEntry,
				err,
				behavior,
				isAllowlisted)
	}

	if shouldReport && !isAllowlisted {
		b.compiler.Report(err)
	}
}

func (b *conformanceAbstractRule) isAllowlistedType(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	if b.allowlistedTypes != nil && n.GetJSType() != nil {
		targetType := n.GetJSType().RestrictByNotNullOrUndefined()
		if targetType.IsSubtypeOf(b.allowlistedTypes) {
			return true
		}
	}
	return false
}

func isKnownNodeType(n rhino.Node) bool {
	return !isUnknownNodeType(n) &&
		!isBottomNodeType(n) &&
		!isTemplateNodeType(n)
}

func (b *conformanceAbstractRule) isNativeObjectType(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	type_ := n.GetJSType()
	return type_ != nil && type_.RestrictByNotNullOrUndefined().Equals(b.nativeObjectType)
}

func isTopNodeType(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	type_ := n.GetJSType()
	return type_ != nil && type_.IsAllType()
}

func isUnknownNodeType(n rhino.Node) bool {
	if n.IsNil() {
		return true
	}
	type_ := n.GetJSType()
	return type_ == nil || type_.IsUnknownType()
}

func isTemplateNodeType(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	type_ := n.GetJSType()
	if type_ == nil {
		return false
	}
	if type_.IsUnionType() {
		for _, member := range type_.GetUnionMembers() {
			if member.IsTemplateType() {
				return true
			}
		}
	}
	return type_.IsTemplateType()
}

func isBottomNodeType(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	type_ := n.GetJSType()
	return type_ != nil && type_.RestrictByNotNullOrUndefined().IsEmptyType()
}

func (b *conformanceAbstractRule) unionTypes(typeNames []string) rhino.JSType {
	registry := b.compiler.GetTypeRegistry()
	var types []rhino.JSType

	for _, typeName := range typeNames {
		type_ := registry.GetGlobalType(typeName)
		if type_ == nil {
			// Types exported from a goog.module (e.g. WizEvent) are resolved via Closure namespace
			// rather than global type registry slots.
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

func (b *conformanceAbstractRule) isAssertionCall(n rhino.Node) bool {
	if call, ok := ast.AsCall(n); ok && b.assertionFunctions != nil {
		callee := call.Callee().Node()
		if callee.IsQualifiedName() {
			return b.assertionFunctions.LookupByCallee(callee) != nil
		}
	}
	return false
}

func (b *conformanceAbstractRule) isTypeImmediatelyTightened(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	if n.GetJSTypeBeforeCast() != nil {
		return true
	}
	parent := n.GetParent()
	if parent.IsNil() {
		return false
	}
	if u, ok := ast.AsUnary(parent); ok && u.Token() == rhino.Token_TYPEOF {
		return true
	}
	return b.isAssertionCall(parent)
}

func (b *conformanceAbstractRule) isNodeUsed(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	parent := n.GetParent()
	if parent.IsNil() || parent.IsName() || NodeUtil_IsLhsByDestructuring(n) {
		return false
	}

	// Consider lvalues in assignment operations to be used iff the actual assignment
	// operation's result is used. e.g. for `a.b.c`:
	//     USED: `alert(x = a.b.c);`
	//   UNUSED: `x = a.b.c;`
	if asgn, ok := ast.AsAssignment(parent); ok {
		if asgn.Node().GetParent().IsNil() {
			return false
		}
		return NodeUtil_IsExpressionResultUsed(parent)
	}

	return NodeUtil_IsExpressionResultUsed(n)
}

func isConformanceCallTarget(n rhino.Node) bool {
	role := ast.RoleOf(n)
	return role == ast.RoleCallCallee || role == ast.RoleNewCallee
}

func isLooseConformanceType(type_ rhino.JSType) bool {
	return type_ == nil || type_.IsUnknownType() || type_.IsNoResolvedType() || type_.IsAllType()
}

func evaluateConformanceTypeString(compiler AbstractCompiler, expression string) (rhino.JSType, error) {
	typeNodes := parsing.JsDocInfoParser_ParseTypeString(expression)
	if typeNodes.IsNil() {
		return nil, errInvalidRequirementSpec("bad type expression")
	}
	typeExpr := rhino.NewJSTypeExpression(typeNodes, "conformance")
	return compiler.GetTypeRegistry().EvaluateTypeExpressionInGlobalScope(typeExpr), nil
}

// validateConformanceCall validates the parameters and receiver `this` type of
// a call or `new` invocation against functionType.
func validateConformanceCall(
	compiler AbstractCompiler,
	callOrNew rhino.Node,
	functionType rhino.FunctionType,
	isCallInvocation bool) bool {
	cn, ok := ast.AsCallOrNew(callOrNew)
	if !ok || callOrNew.IsOptChainCall() {
		return true
	}

	return validateConformanceParameterList(compiler, cn, functionType, isCallInvocation) &&
		validateConformanceCallThis(cn, functionType, isCallInvocation)
}

func validateConformanceCallThis(
	callOrNew ast.CallOrNew, functionType rhino.FunctionType, isCallInvocation bool) bool {
	if callOrNew.Node().IsOptChainCall() {
		return false
	}
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
	} else if get, ok := ast.AsNormalGet(callOrNew.Callee().Node()); ok {
		thisNode = get.Object().Node()
	}
	if thisNode.IsNil() || thisNode.GetJSType() == nil {
		return false
	}
	thisNodeType := thisNode.GetJSType().RestrictByNotNullOrUndefined()
	return thisNodeType.IsSubtypeOf(thisType)
}

func validateConformanceParameterList(
	compiler AbstractCompiler,
	callOrNew ast.CallOrNew,
	functionType rhino.FunctionType,
	isCallInvocation bool) bool {
	var argumentTypes []rhino.JSType
	skippedThisArg := false
	for arg := range callOrNew.Args() {
		if isCallInvocation && !skippedThisArg {
			skippedThisArg = true
			continue
		}
		argType := arg.Node().GetJSType()
		if argType == nil {
			argType = compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		argumentTypes = append(argumentTypes, argType)
	}
	return functionType.AcceptsArguments(argumentTypes)
}

func splitConformanceClassAndProperty(specName string) (className, propName string, ok bool) {
	className, propName, ok = strings.Cut(specName, ".prototype.")
	if !ok {
		return "", "", false
	}
	if strings.Contains(propName, ".prototype.") {
		panic(jutil.NewIllegalStateExceptionMsg("unexpected multiple .prototype. segments: " + specName))
	}
	if className == "" || propName == "" {
		return "", "", false
	}
	return className, propName, true
}

// splitConformanceNameAndType splits a "NAME:TYPE" requirement value at the
// first colon, requiring a non-empty name prefix.
func splitConformanceNameAndType(specName string) (name, typeDecl string, ok bool) {
	name, typeDecl, ok = strings.Cut(specName, ":")
	if !ok || name == "" {
		return "", "", false
	}
	return name, typeDecl, true
}

func isAnnotatedAsConst(jsdoc *rhino.JSDocInfo) bool {
	return jsdoc != nil && jsdoc.IsConstant()
}

// ConformanceRules_ConformanceUtil_ValueFromGlobalName resolves the constant string
// value assigned to an optional global qualified name without pointer boxing.
func ConformanceRules_ConformanceUtil_ValueFromGlobalName(
	scope *Scope, gns *GlobalNamespace, qualifiedName base.OptString) (string, bool) {
	qname, hasQName := qualifiedName.Get()
	return conformanceValueFromGlobalName(scope, gns, qname, hasQName)
}

type conformanceVisitSet struct {
	inline   [8]rhino.Node
	n        int
	overflow []rhino.Node
}

const conformanceMaxInferenceDepth = 64

func (s *conformanceVisitSet) add(node rhino.Node) bool {
	if node.IsNil() {
		return true
	}
	if s.n+len(s.overflow) >= conformanceMaxInferenceDepth {
		return false
	}
	for i := range s.n {
		if s.inline[i] == node {
			return false
		}
	}
	for _, prev := range s.overflow {
		if prev == node {
			return false
		}
	}
	if s.n < len(s.inline) {
		s.inline[s.n] = node
		s.n++
	} else {
		s.overflow = append(s.overflow, node)
	}
	return true
}

func conformanceValueFromGlobalName(
	scope *Scope, gns *GlobalNamespace, qualifiedName string, hasQualifiedName bool) (string, bool) {
	return conformanceValueFromGlobalNameVisited(scope, gns, qualifiedName, hasQualifiedName, conformanceVisitSet{})
}

func conformanceValueFromGlobalNameVisited(
	scope *Scope, gns *GlobalNamespace, qualifiedName string, hasQualifiedName bool, visited conformanceVisitSet) (string, bool) {
	if gns == nil {
		return "", false
	}
	var gn *GlobalNamespace_Name
	if !hasQualifiedName {
		gns.EnsureGenerated()
	} else if slot := gns.GetOwnSlot(qualifiedName); slot != nil {
		gn = slot.(*GlobalNamespace_Name)
	}
	if gn == nil ||
		gn.GetGlobalSets() != 1 ||
		gn.GetTotalSets() != 1 ||
		!isAnnotatedAsConst(gn.GetJSDocInfo()) {
		return "", false
	}
	decl := gn.GetDeclaration()
	if decl == nil || !visited.add(decl.GetNode()) {
		return "", false
	}

	return inferConformanceStringValueVisited(
		scope, NodeUtil_GetRValueOfLValue(decl.GetNode()), func() *GlobalNamespace { return gns }, visited)
}

// ConformanceRules_ConformanceUtil_InferStringValueOk infers a constant string
// value for node without pointer boxing.
func ConformanceRules_ConformanceUtil_InferStringValueOk(
	scope *Scope, node rhino.Node, globalNamespaceSupplier func() *GlobalNamespace) (string, bool) {
	return inferConformanceStringValue(scope, node, globalNamespaceSupplier)
}

func inferConformanceStringValue(
	scope *Scope, node rhino.Node, globalNamespaceSupplier func() *GlobalNamespace) (string, bool) {
	return inferConformanceStringValueVisited(scope, node, globalNamespaceSupplier, conformanceVisitSet{})
}

func inferConformanceStringValueVisited(
	scope *Scope, node rhino.Node, globalNamespaceSupplier func() *GlobalNamespace, visited conformanceVisitSet) (string, bool) {
	if node.IsNil() {
		return "", false
	}

	switch node.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT, rhino.Token_TEMPLATELIT_STRING:
		return NodeUtil_GetStringValue(node)
	case rhino.Token_NAME:
		if scope == nil || !visited.add(node) {
			return "", false
		}
		var_ := scope.GetVarByID(node.GetStringID())
		if var_ == nil {
			return "", false
		}
		if !var_.IsConst() && !isAnnotatedAsConst(var_.GetJSDocInfo()) {
			return "", false
		}
		if !visited.add(var_.GetNameNode()) {
			return "", false
		}
		initialValue := var_.GetInitialValue()
		return inferConformanceStringValueVisited(var_.GetScope().(*Scope), initialValue, globalNamespaceSupplier, visited)
	case rhino.Token_GETPROP:
		if !visited.add(node) {
			return "", false
		}
		gp := ast.MustGetProp(node)
		type_ := node.GetJSType()
		if type_ == nil {
			return "", false
		}

		// The JS style guide requires enums to be effectively immutable and all enum items should
		// be statically known. See go/js-style#features-objects-enums.
		if type_.IsEnumElementType() {
			enumSource := type_.ToMaybeEnumElementType().GetEnumType().GetSource()
			if enumSource.IsNil() || (!enumSource.IsObjectLit() && !enumSource.IsClassMembers()) {
				return "", false
			}
			return inferConformanceStringValueVisited(
				nil,
				NodeUtil_GetFirstPropMatchingKey(enumSource, gp.Property()),
				globalNamespaceSupplier,
				visited)
		} else if type_.IsString() && globalNamespaceSupplier != nil {
			gns := globalNamespaceSupplier()
			if gns == nil {
				return "", false
			}
			qname, hasQName := node.GetQualifiedNameString()
			return conformanceValueFromGlobalNameVisited(scope, gns, qname, hasQName, visited)
		}
		return "", false
	}
	return "", false
}

func isXidType(type_ rhino.JSType) bool {
	if type_ == nil {
		return false
	}
	enumElTy := type_.ToMaybeEnumElementType()
	if enumElTy == nil {
		return false
	}
	refName, ok := enumElTy.GetEnumType().ReferenceName()
	return ok && refName == "enum{xid.String}"
}

func isEventHandlerAttrName(attr string) bool {
	return attr != "on" && strings.HasPrefix(attr, "on")
}

func isLooseObjectType(type_ rhino.JSType, registry *rhino.JSTypeRegistry) bool {
	return type_ != nil && type_.Equals(registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
}
