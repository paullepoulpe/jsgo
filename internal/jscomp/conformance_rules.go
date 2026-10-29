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
	"regexp"
	"regexp/syntax"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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

func conformanceRules_isValidPattern(regex string) bool {
	var sb strings.Builder
	rs := []rune(regex)
	n := len(rs)
	inClass := false
	classStart := 0
	for i := 0; i < n; i++ {
		c := rs[i]
		switch {
		case c == '\\':
			if i+1 >= n {
				return false
			}
			i++
			e := rs[i]
			switch {
			case e == 'Q':
				end := -1
				for j := i + 1; j+1 < n; j++ {
					if rs[j] == '\\' && rs[j+1] == 'E' {
						end = j
						break
					}
				}
				var lit string
				if end == -1 {
					lit = string(rs[i+1:])
					i = n
				} else {
					lit = string(rs[i+1 : end])
					i = end + 1
				}
				if inClass {
					for _, r := range lit {
						sb.WriteString(regexp.QuoteMeta(string(r)))
					}
				} else {
					sb.WriteString(regexp.QuoteMeta(lit))
				}
			case e == 'E':
			case e == 's', e == 'S', e == 'v', e == 'h', e == 'V', e == 'H',
				e == 'e', e == 'a',
				e == '0', e >= '1' && e <= '9', e == 'k',
				e == 'b', e == 'B', e == 'z', e == 'Z', e == 'G', e == 'A', e == 'R', e == 'X',
				e == 'p', e == 'P', e == 'x', e == 'u', e == 'c', e == 'N':
				sb.WriteByte('a')
			case e == 'd', e == 'D', e == 'w', e == 'W', e == 't', e == 'n', e == 'r', e == 'f':
				sb.WriteByte('\\')
				sb.WriteRune(e)
			case (e >= 'a' && e <= 'z') || (e >= 'A' && e <= 'Z'):
				return false
			default:
				sb.WriteString(regexp.QuoteMeta(string(e)))
			}
		case inClass:
			if c == ']' {
				if i == classStart || (i == classStart+1 && rs[classStart] == '^') {
					sb.WriteString(`\]`)
				} else {
					inClass = false
					sb.WriteByte(']')
				}
			} else {
				sb.WriteRune(c)
			}
		case c == '[':
			inClass = true
			classStart = i + 1
			sb.WriteByte('[')
		default:
			sb.WriteRune(c)
		}
	}
	if inClass {
		return false
	}
	if _, err := syntax.Parse(sb.String(), syntax.Perl); err != nil {
		if se, ok := err.(*syntax.Error); ok {
			switch se.Code {
			case syntax.ErrInvalidRepeatSize, syntax.ErrInvalidRepeatOp, syntax.ErrLarge,
				syntax.ErrNestingDepth, syntax.ErrInternalError:
				return true
			}
		}
		return false
	}
	return true
}

// Java: ConformanceRules.buildPattern(List<String>)
func conformanceRules_buildPattern(reqPatterns []string) (*jutil.Pattern, error) {
	if len(reqPatterns) == 0 {
		return nil, nil
	}

	// validate the patterns
	for _, reqPattern := range reqPatterns {
		if !conformanceRules_isValidPattern(reqPattern) {
			return nil, errInvalidRequirementSpec("invalid regex pattern")
		}
	}

	if len(reqPatterns) == 1 {
		return jutil.Pattern_Compile(reqPatterns[0]), nil
	}
	return jutil.Pattern_Compile("(" + strings.Join(reqPatterns, "|") + ")"), nil
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
		(a.regexp.Find(path) || (hasTsPath && a.regexp.Find(tsPath)))
}

// ConformanceRules_AbstractRule holds the common rule fields and helpers shared
// by AbstractRule and AbstractTypeRestrictionRule subclasses, without virtual
// self-dispatch pointers.
type ConformanceRules_AbstractRule struct {
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

// Java: CheckConformance.Rule.getPrecondition() (default method)
func (*ConformanceRules_AbstractRule) GetPrecondition() CheckConformance_Precondition {
	return CheckConformance_Precondition_CHECK_ALL
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

	if len(requirement.GetWhitelistList()) > 0 || len(requirement.GetWhitelistRegexpList()) > 0 {
		allowlist, err := newConformanceRules_AllowList(requirement.GetWhitelistList(), requirement.GetWhitelistRegexpList())
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	if len(requirement.GetAllowlistList()) > 0 || len(requirement.GetAllowlistRegexpList()) > 0 {
		allowlist, err := newConformanceRules_AllowList(requirement.GetAllowlistList(), requirement.GetAllowlistRegexpList())
		if err != nil {
			return err
		}
		allowlistsBuilder = append(allowlistsBuilder, allowlist)
	}
	b.allowlists = allowlistsBuilder

	if len(requirement.GetOnlyApplyToList()) > 0 || len(requirement.GetOnlyApplyToRegexpList()) > 0 {
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
		!(allowlist == nil && (b.onlyApplyTo == nil || (path != nil && b.onlyApplyTo.Matches(*path))))
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
func conformanceRules_AbstractTypeRestrictionRule_isKnown(n rhino.Node) bool {
	return !conformanceRules_AbstractTypeRestrictionRule_isUnknown(n) &&
		!conformanceRules_AbstractTypeRestrictionRule_isBottom(n) &&
		!conformanceRules_AbstractTypeRestrictionRule_isTemplateType(n)
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isNativeObjectType(Node)
func (b *ConformanceRules_AbstractRule) IsNativeObjectType(n rhino.Node) bool {
	type_ := n.GetJSType().RestrictByNotNullOrUndefined()
	return type_.Equals(b.nativeObjectType)
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isTop(Node)
func conformanceRules_AbstractTypeRestrictionRule_isTop(n rhino.Node) bool {
	type_ := n.GetJSType()
	return type_ != nil && type_.IsAllType()
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isUnknown(Node)
func conformanceRules_AbstractTypeRestrictionRule_isUnknown(n rhino.Node) bool {
	type_ := n.GetJSType()
	return type_ == nil || type_.IsUnknownType()
}

// Java: ConformanceRules.AbstractTypeRestrictionRule.isTemplateType(Node)
func conformanceRules_AbstractTypeRestrictionRule_isTemplateType(n rhino.Node) bool {
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
func conformanceRules_AbstractTypeRestrictionRule_isBottom(n rhino.Node) bool {
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

// Java: ConformanceRules.isLooseObject(JSType, JSTypeRegistry)
func ConformanceRules_IsLooseObject(type_ rhino.JSType, registry *rhino.JSTypeRegistry) bool {
	return type_.Equals(registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
}
