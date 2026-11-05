// Port of com/google/javascript/jscomp/CheckConformance.java @ closure-compiler v20261028

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
	"math/bits"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type flatConformanceRule struct {
	rule               CheckConformance_Rule
	absRule            *ConformanceRules_AbstractRule
	cat                *CheckConformance_Category
	behavior           ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior
	catIdx             int
	needPrecond        bool
	canFilterPerScript bool
	isBannedName       bool
}

// CheckConformance provides a framework for checking code against a set of
// user configured conformance requirements Requirement. The requirements are
// specified by the ConformanceConfig proto, which allows for both standard
// checks (forbidden properties, variables, or dependencies) and allow for
// more complex checks using requirements of CUSTOM type.
//
// The requirements are translated into rules CheckConformance.Rule which are
// then checked by this pass.
//
// Conformance violations are both reported as compiler errors, and are also
// reported separately to the {cI gue@link ErrorManager}
type CheckConformance struct {
	compiler           AbstractCompiler
	categories         []*CheckConformance_Category
	activeTokens       ast.TokenSet
	hasAllTokens       bool
	flatRules          []flatConformanceRule
	useMask64          bool
	tokenMask0         [256]uint64
	tokenMasks         [256][]uint64
	scriptSourceFileID uint32
	scriptSkipMask0    uint64
	scriptSkipMask     []uint64
}

var (
	_ nodeTraversal_visitTokens                         = (*CheckConformance)(nil)
	_ nodeTraversal_onlyScriptAndFunctionShouldTraverse = (*CheckConformance)(nil)
)

func (*CheckConformance) shouldTraverseOnlyScriptAndFunction() bool {
	return true
}

func (c *CheckConformance) visitTokens() *ast.TokenSet {
	if c.hasAllTokens {
		return nil
	}
	return &c.activeTokens
}

// Java: CheckConformance.CONFORMANCE_ERROR
var CheckConformance_CONFORMANCE_ERROR = DiagnosticType_Error("JSC_CONFORMANCE_ERROR", "Violation: {0}{1}{2}")

// Java: CheckConformance.CONFORMANCE_VIOLATION
var CheckConformance_CONFORMANCE_VIOLATION = DiagnosticType_Warning("JSC_CONFORMANCE_VIOLATION", "Violation: {0}{1}{2}")

// Java: CheckConformance.CONFORMANCE_POSSIBLE_VIOLATION
var CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION = DiagnosticType_Warning("JSC_CONFORMANCE_POSSIBLE_VIOLATION", "Possible violation: {0}{1}{2}")

// Java: CheckConformance.INVALID_REQUIREMENT_SPEC
var CheckConformance_INVALID_REQUIREMENT_SPEC = DiagnosticType_Error(
	"JSC_INVALID_REQUIREMENT_SPEC",
	"Invalid requirement. Reason: {0}\nRequirement spec:\n{1}")

// CheckConformance_Rule is a rule that can be checked for conformance.
type CheckConformance_Rule interface {
	// GetPrecondition returns a precondition for this rule.
	//
	// This method will only be called once (per rule) during the creation of
	// the CheckConformance pass. Therefore, the return must be constant.
	//
	// Returning null means that there is no precondition. This is convenient,
	// but can be a major performance hit.
	GetPrecondition() CheckConformance_Precondition

	// Check performs conformance check.
	Check(t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior)
}

// CheckConformance_Precondition is a condition that must be true for a rule
// to possibly match a node.
//
// Instances are used as keys to group rules with common preconditions.
// Grouping allows shared computation to be done only once per node, which is
// a substantial performance improvement.
type CheckConformance_Precondition interface {
	ShouldCheck(n rhino.Node) bool
}

// checkConformance_NodePrecondition enumerates built-in node preconditions by
// distinct scalar value so rules group identically to Java without closure
// pointer boxing.
type checkConformance_NodePrecondition uint8

const (
	checkConformance_PreconditionCheckAll checkConformance_NodePrecondition = iota
	checkConformance_PreconditionBannedNameCandidate
	checkConformance_PreconditionBannedEnhanceScript
	checkConformance_PreconditionBannedModsRegexScript
)

func (p checkConformance_NodePrecondition) ShouldCheck(n rhino.Node) bool {
	switch p {
	case checkConformance_PreconditionCheckAll:
		return true
	case checkConformance_PreconditionBannedNameCandidate:
		switch n.GetToken() {
		case rhino.Token_GETPROP:
			return n.GetFirstChild().IsQualifiedName()
		case rhino.Token_NAME:
			return n.GetStringID() != strpool.EmptyID
		default:
			return false
		}
	case checkConformance_PreconditionBannedEnhanceScript, checkConformance_PreconditionBannedModsRegexScript:
		return n.IsScript()
	default:
		panic(jutil.NewIllegalStateException())
	}
}

// Java: CheckConformance.Precondition.CHECK_ALL
var CheckConformance_Precondition_CHECK_ALL CheckConformance_Precondition = checkConformance_PreconditionCheckAll

type categorizedRule struct {
	rule     CheckConformance_Rule
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior
}

type checkConformance_CategoryKind uint8

const (
	checkConformance_KindCheckAll checkConformance_CategoryKind = iota
	checkConformance_KindNodePrecondition
	checkConformance_KindBannedProperty
	checkConformance_KindCustom
)

// Java: CheckConformance.Category
type CheckConformance_Category struct {
	precondition CheckConformance_Precondition
	bannedProp   ConformanceRules_BannedProperty_RequirementPrecondition
	nodePrecond  checkConformance_NodePrecondition
	kind         checkConformance_CategoryKind
	rules        []categorizedRule
}

func (cat *CheckConformance_Category) shouldCheck(n rhino.Node) bool {
	switch cat.kind {
	case checkConformance_KindCheckAll:
		return true
	case checkConformance_KindNodePrecondition:
		return cat.nodePrecond.ShouldCheck(n)
	case checkConformance_KindBannedProperty:
		return cat.bannedProp.ShouldCheck(n)
	default:
		return cat.precondition.ShouldCheck(n)
	}
}

// Java: CheckConformance(AbstractCompiler)
func NewCheckConformance(compiler AbstractCompiler) *CheckConformance {
	return NewCheckConformanceWithConfigs(compiler, []*ConformanceConfig{}, nil)
}

// NewCheckConformanceWithConfigs: configs are the rules to check.
// reportingMode is a nullable enum (nil is Java null).
//
// Java: CheckConformance(AbstractCompiler, ImmutableList<ConformanceConfig>, ConformanceReportingMode)
func NewCheckConformanceWithConfigs(
	compiler AbstractCompiler,
	configs []*ConformanceConfig,
	reportingMode *CompilerOptions_ConformanceReportingMode) *CheckConformance {
	c := &CheckConformance{compiler: compiler}
	// Initialize the map of functions to inspect for renaming candidates.
	c.categories = c.initRules(compiler, configs, reportingMode)
	c.buildDispatchTables()
	return c
}

// Java: CheckConformance.process(Node, Node)
func (c *CheckConformance) Process(externs rhino.Node, root rhino.Node) {
	if len(c.categories) != 0 {
		NodeTraversal_TraverseRoots(c.compiler, c, externs, root)
	}
}

// Java: CheckConformance.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckConformance) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Don't inspect extern files, *.tsmes.closure.js, weak sources, or @closureUnaware ASTs.
	if !n.IsScript() {
		return true
	}
	if !c.IsScriptOfInterest(t.GetInput().GetSourceFile()) {
		return false
	}
	if strings.HasSuffix(t.GetSourceName(), "tsmes.closure.js") {
		return false
	}
	c.updateScriptFilter(n)
	return true
}

// Java: CheckConformance.isScriptOfInterest(SourceFile)
func (c *CheckConformance) IsScriptOfInterest(sf *SourceFile) bool {
	return !sf.IsWeak() && !sf.IsExtern()
}

// Java: CheckConformance.visit(NodeTraversal, Node, Node)
func (c *CheckConformance) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	tok := uint8(n.GetToken())
	if c.useMask64 {
		mask := c.tokenMask0[tok]
		if mask == 0 {
			return
		}
		if c.scriptSkipMask0 != 0 && n.SourceFileID() == c.scriptSourceFileID {
			mask &^= c.scriptSkipMask0
			if mask == 0 {
				return
			}
		}
		var testedCats, passedCats uint64
		for mask != 0 {
			idx := bits.TrailingZeros64(mask)
			mask &= mask - 1
			c.runFlatRule(&c.flatRules[idx], t, n, &testedCats, &passedCats)
		}
		return
	}
	wordMasks := c.tokenMasks[tok]
	if len(wordMasks) == 0 {
		return
	}
	applySkip := len(c.scriptSkipMask) != 0 && n.SourceFileID() == c.scriptSourceFileID
	var testedCats, passedCats uint64
	for w, mask := range wordMasks {
		if applySkip {
			mask &^= c.scriptSkipMask[w]
		}
		base := w << 6
		for mask != 0 {
			idx := base + bits.TrailingZeros64(mask)
			mask &= mask - 1
			c.runFlatRule(&c.flatRules[idx], t, n, &testedCats, &passedCats)
		}
	}
}

func (c *CheckConformance) runFlatRule(
	entry *flatConformanceRule,
	t *NodeTraversal,
	n rhino.Node,
	testedCats, passedCats *uint64,
) {
	if entry.needPrecond {
		if entry.catIdx < 64 {
			catBit := uint64(1) << uint(entry.catIdx)
			if *testedCats&catBit == 0 {
				*testedCats |= catBit
				if entry.cat.shouldCheck(n) {
					*passedCats |= catBit
				}
			}
			if *passedCats&catBit == 0 {
				return
			}
		} else if !entry.cat.shouldCheck(n) {
			return
		}
	}
	if entry.absRule != nil {
		if result := entry.absRule.checkConformance(t, n); result.IsViolation() {
			entry.absRule.Report(n, result, entry.behavior)
		}
		return
	}
	entry.rule.Check(t, n, entry.behavior)
}

func (c *CheckConformance) updateScriptFilter(script rhino.Node) {
	c.scriptSkipMask0 = 0
	if len(c.scriptSkipMask) != 0 {
		clear(c.scriptSkipMask)
	}
	sfID := script.SourceFileID()
	c.scriptSourceFileID = sfID
	if sfID == 0 || len(c.flatRules) == 0 || !hasDefaultConformanceReporting(c.compiler) {
		return
	}
	rawPath, hasPath := script.SourceFileName()
	if !hasPath {
		return
	}
	rawTsPath, hasRawTsPath := splitConformanceTsPath(rawPath)
	normPath := computeConformanceNormPath(c.compiler, rawPath)
	normTsPath, hasNormTsPath := rawTsPath, hasRawTsPath
	if normPath != rawPath {
		normTsPath, hasNormTsPath = splitConformanceTsPath(normPath)
	}
	hasCustomBannedNameEval := CustomBannedNameEvaluator != nil
	for i := range c.flatRules {
		entry := &c.flatRules[i]
		if !entry.canFilterPerScript || (entry.isBannedName && hasCustomBannedNameEval) {
			continue
		}
		if entry.behavior == ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_RECORD_ONLY ||
			entry.absRule.isExcludedOrAllowlisted(rawPath, rawTsPath, hasRawTsPath, normPath, normTsPath, hasNormTsPath, true) {
			if c.useMask64 {
				c.scriptSkipMask0 |= uint64(1) << uint(i)
			} else {
				c.scriptSkipMask[i>>6] |= uint64(1) << uint(i&63)
			}
		}
	}
}

func categoryNeedsRuntimePrecond(cat *CheckConformance_Category) bool {
	switch cat.kind {
	case checkConformance_KindCheckAll:
		return false
	case checkConformance_KindNodePrecondition:
		switch cat.nodePrecond {
		case checkConformance_PreconditionBannedEnhanceScript, checkConformance_PreconditionBannedModsRegexScript:
			return false
		default:
			return true
		}
	case checkConformance_KindBannedProperty:
		return cat.bannedProp != ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY
	default:
		return true
	}
}

func inspectConformanceRule(rule CheckConformance_Rule) (
	tokens ast.TokenSet,
	allTokens bool,
	absRule *ConformanceRules_AbstractRule,
	canFilterPerScript bool,
	isBannedName bool,
) {
	switch r := rule.(type) {
	case *conformanceNoOpRule:
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceInferredConstRule:
		return tokens, true, &r.conformanceAbstractRule, true, false
	case *conformanceBannedDependencyRule:
		tokens.Add(rhino.Token_SCRIPT)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBannedDependencyRegexRule:
		tokens.Add(rhino.Token_SCRIPT)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBannedNameRule:
		tokens.Add(rhino.Token_GETPROP)
		tokens.Add(rhino.Token_NAME)
		return tokens, false, &r.conformanceAbstractRule, true, true
	case *ConformanceRules_BannedProperty:
		if len(r.props) != 0 {
			tokens.Add(rhino.Token_STRING_KEY)
			tokens.Add(rhino.Token_GETPROP)
			tokens.Add(rhino.Token_GETELEM)
			tokens.Add(rhino.Token_COMPUTED_PROP)
		}
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBannedStringRegexRule:
		tokens.Add(rhino.Token_STRINGLIT)
		tokens.Add(rhino.Token_TEMPLATELIT_STRING)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceRestrictedNameCallRule:
		tokens.Add(rhino.Token_NAME)
		tokens.Add(rhino.Token_GETPROP)
		tokens.Add(rhino.Token_THIS)
		tokens.Add(rhino.Token_SUPER)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceRestrictedMethodCallRule:
		tokens.Add(rhino.Token_GETPROP)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceRestrictedPropertyWriteRule:
		tokens.Add(rhino.Token_GETPROP)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBannedCodePatternRule:
		for _, m := range r.restrictions {
			start := m.templateStart
			if start.IsNil() {
				return tokens, true, &r.conformanceAbstractRule, true, false
			}
			switch start.GetToken() {
			case templateTypeParam, templateStringLiteral:
				return tokens, true, &r.conformanceAbstractRule, true, false
			case templateLocalName:
				tokens.Add(rhino.Token_NAME)
			default:
				tokens.Add(start.GetToken())
			}
		}
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanForOfRule:
		tokens.Add(rhino.Token_FOR_OF)
		tokens.Add(rhino.Token_FOR_AWAIT_OF)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceRequireUseStrictRule:
		tokens.Add(rhino.Token_SCRIPT)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanThrowOfNonErrorTypesRule:
		if r.errorObjType != nil {
			tokens.Add(rhino.Token_THROW)
		}
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanNullDerefRule:
		tokens.Add(rhino.Token_GETPROP)
		tokens.Add(rhino.Token_GETELEM)
		tokens.Add(rhino.Token_CALL)
		tokens.Add(rhino.Token_NEW)
		tokens.Add(rhino.Token_IN)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanUnknownThisRule:
		// BanUnknownThis deduplicates reports per scope root (including global scope),
		// so keep it unfiltered across scripts while restricting it to THIS nodes.
		tokens.Add(rhino.Token_THIS)
		return tokens, false, &r.conformanceAbstractRule, false, false
	case *conformanceBanUnknownDirectThisPropsRule:
		tokens.Add(rhino.Token_GETPROP)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanUnknownTypedClassPropsRule:
		tokens.Add(rhino.Token_GETPROP)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanUnresolvedTypeRule:
		tokens.Add(rhino.Token_GETPROP)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceStrictBanUnresolvedTypeRule:
		return tokens, true, &r.conformanceAbstractRule, true, false
	case *conformanceBanGlobalVarsRule:
		tokens.Add(rhino.Token_VAR)
		tokens.Add(rhino.Token_LET)
		tokens.Add(rhino.Token_CONST)
		tokens.Add(rhino.Token_FUNCTION)
		tokens.Add(rhino.Token_CLASS)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBannedEnhanceRule:
		tokens.Add(rhino.Token_SCRIPT)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBannedModsRegexRule:
		tokens.Add(rhino.Token_SCRIPT)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanCreateElementRule:
		tokens.Add(rhino.Token_CALL)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanCreateDomRule:
		tokens.Add(rhino.Token_CALL)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanElementSetAttributeRule:
		tokens.Add(rhino.Token_CALL)
		tokens.Add(rhino.Token_GETELEM)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanSettingAttributesRule:
		tokens.Add(rhino.Token_CALL)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanExecCommandRule:
		tokens.Add(rhino.Token_CALL)
		return tokens, false, &r.conformanceAbstractRule, true, false
	case *conformanceBanStaticThisRule:
		tokens.Add(rhino.Token_THIS)
		return tokens, false, &r.conformanceAbstractRule, true, false
	default:
		switch p := rule.GetPrecondition().(type) {
		case checkConformance_NodePrecondition:
			switch p {
			case checkConformance_PreconditionBannedNameCandidate:
				tokens.Add(rhino.Token_GETPROP)
				tokens.Add(rhino.Token_NAME)
				return tokens, false, nil, false, false
			case checkConformance_PreconditionBannedEnhanceScript, checkConformance_PreconditionBannedModsRegexScript:
				tokens.Add(rhino.Token_SCRIPT)
				return tokens, false, nil, false, false
			}
		case ConformanceRules_BannedProperty_RequirementPrecondition:
			if p == ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY {
				tokens.Add(rhino.Token_STRING_KEY)
				tokens.Add(rhino.Token_GETPROP)
				tokens.Add(rhino.Token_GETELEM)
				tokens.Add(rhino.Token_COMPUTED_PROP)
				return tokens, false, nil, false, false
			}
		}
		return tokens, true, nil, false, false
	}
}

func (c *CheckConformance) buildDispatchTables() {
	c.activeTokens = ast.TokenSet{}
	c.hasAllTokens = false
	totalRules := 0
	for _, cat := range c.categories {
		totalRules += len(cat.rules)
	}
	if totalRules == 0 {
		c.useMask64 = true
		return
	}

	c.flatRules = make([]flatConformanceRule, 0, totalRules)
	c.useMask64 = totalRules <= 64
	numWords := (totalRules + 63) >> 6
	var maskBacking []uint64
	if !c.useMask64 {
		c.scriptSkipMask = make([]uint64, numWords)
		maskBacking = make([]uint64, 256*numWords)
	}

	for catIdx := len(c.categories) - 1; catIdx >= 0; catIdx-- {
		cat := c.categories[catIdx]
		needPrecond := categoryNeedsRuntimePrecond(cat)
		rules := cat.rules
		for rIdx := len(rules) - 1; rIdx >= 0; rIdx-- {
			entry := rules[rIdx]
			tokens, allTokens, absRule, canFilter, isBannedName := inspectConformanceRule(entry.rule)
			flatIdx := len(c.flatRules)
			c.flatRules = append(c.flatRules, flatConformanceRule{
				rule:               entry.rule,
				absRule:            absRule,
				cat:                cat,
				behavior:           entry.behavior,
				catIdx:             catIdx,
				needPrecond:        needPrecond,
				canFilterPerScript: canFilter,
				isBannedName:       isBannedName,
			})
			if allTokens {
				c.hasAllTokens = true
				if c.useMask64 {
					bit := uint64(1) << uint(flatIdx)
					for tok := range 256 {
						c.tokenMask0[tok] |= bit
					}
				} else {
					w, bit := flatIdx>>6, uint64(1)<<uint(flatIdx&63)
					for tok := range 256 {
						if c.tokenMasks[tok] == nil {
							c.tokenMasks[tok] = maskBacking[tok*numWords : (tok+1)*numWords]
						}
						c.tokenMasks[tok][w] |= bit
					}
				}
				continue
			}
			for wIdx, word := range tokens {
				c.activeTokens[wIdx] |= word
				baseTok := wIdx << 6
				for word != 0 {
					tok := baseTok + bits.TrailingZeros64(word)
					word &= word - 1
					if c.useMask64 {
						c.tokenMask0[tok] |= uint64(1) << uint(flatIdx)
					} else {
						if c.tokenMasks[tok] == nil {
							c.tokenMasks[tok] = maskBacking[tok*numWords : (tok+1)*numWords]
						}
						c.tokenMasks[tok][flatIdx>>6] |= uint64(1) << uint(flatIdx&63)
					}
				}
			}
		}
	}
}

func (c *CheckConformance) newCategory(rule CheckConformance_Rule, precondition CheckConformance_Precondition) *CheckConformance_Category {
	cat := &CheckConformance_Category{precondition: precondition, kind: checkConformance_KindCustom}
	switch p := precondition.(type) {
	case nil:
		c.hasAllTokens = true
		cat.kind = checkConformance_KindCheckAll
	case checkConformance_NodePrecondition:
		switch p {
		case checkConformance_PreconditionCheckAll:
			c.hasAllTokens = true
			cat.kind = checkConformance_KindCheckAll
		case checkConformance_PreconditionBannedNameCandidate:
			c.activeTokens.Add(rhino.Token_GETPROP)
			c.activeTokens.Add(rhino.Token_NAME)
			cat.kind = checkConformance_KindNodePrecondition
			cat.nodePrecond = p
		case checkConformance_PreconditionBannedEnhanceScript, checkConformance_PreconditionBannedModsRegexScript:
			c.activeTokens.Add(rhino.Token_SCRIPT)
			cat.kind = checkConformance_KindNodePrecondition
			cat.nodePrecond = p
		default:
			c.hasAllTokens = true
		}
	case ConformanceRules_BannedProperty_RequirementPrecondition:
		if p == ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY {
			c.addBannedPropertyTokens()
		} else if _, isBannedProp := rule.(*ConformanceRules_BannedProperty); isBannedProp {
			c.addBannedPropertyTokens()
		} else {
			c.hasAllTokens = true
		}
		cat.kind = checkConformance_KindBannedProperty
		cat.bannedProp = p
	default:
		c.hasAllTokens = true
	}
	return cat
}

func (c *CheckConformance) addBannedPropertyTokens() {
	c.activeTokens.Add(rhino.Token_STRING_KEY)
	c.activeTokens.Add(rhino.Token_GETPROP)
	c.activeTokens.Add(rhino.Token_GETELEM)
	c.activeTokens.Add(rhino.Token_COMPUTED_PROP)
}

type mergedRequirement struct {
	requirement *Requirement
	behavior    ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior
}

// initRules builds the data structures needed by this pass from the provided
// configurations.
//
// Java: CheckConformance.initRules(AbstractCompiler, ImmutableList<ConformanceConfig>, ConformanceReportingMode)
func (c *CheckConformance) initRules(
	compiler AbstractCompiler,
	configs []*ConformanceConfig,
	reportingMode *CompilerOptions_ConformanceReportingMode) []*CheckConformance_Category {
	isLibraryLevelReportingMode :=
		reportingMode != nil &&
			*reportingMode == CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	if isLibraryLevelReportingMode {
		if !CheckConformance_ValidateBehaviorSettingOfConfigs(compiler, configs, *reportingMode) {
			return nil
		}
	}

	mergedRequirements := c.mergeRequirements(compiler, configs, isLibraryLevelReportingMode)

	var categories []*CheckConformance_Category
	categoryIndex := make(map[CheckConformance_Precondition]int)
	for _, entry := range mergedRequirements {
		rule := CheckConformance_InitRule(compiler, entry.requirement)
		if rule == nil {
			continue
		}
		precondition := rule.GetPrecondition()
		idx, ok := categoryIndex[precondition]
		if !ok {
			idx = len(categories)
			categoryIndex[precondition] = idx
			categories = append(categories, c.newCategory(rule, precondition))
		} else {
			// Ensure activeTokens/hasAllTokens still updates if multiple rules share a precondition
			// (e.g. a non-BannedProperty custom rule reusing RequirementPrecondition).
			if _, isBannedPropPrecond := precondition.(ConformanceRules_BannedProperty_RequirementPrecondition); isBannedPropPrecond {
				if _, isBannedProp := rule.(*ConformanceRules_BannedProperty); !isBannedProp &&
					precondition != ConformanceRules_BannedProperty_RequirementPrecondition_BANNED_PROPERTY {
					c.hasAllTokens = true
				}
			}
		}
		categories[idx].rules = append(categories[idx].rules, categorizedRule{rule: rule, behavior: entry.behavior})
	}
	return categories
}

// Java: CheckConformance.validateBehaviorSettingOfConfigs(AbstractCompiler, ImmutableList<ConformanceConfig>, ConformanceReportingMode)
func CheckConformance_ValidateBehaviorSettingOfConfigs(
	compiler AbstractCompiler,
	configs []*ConformanceConfig,
	reportingMode CompilerOptions_ConformanceReportingMode) bool {
	// only validate if the behavior matters (i.e. if the compiler is running in library-level
	// conformance mode)
	jutil.CheckState(
		reportingMode == CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG)
	baseRequirementBehaviors := make(map[string]ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior)
	for _, config := range configs {
		if !config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior() {
			// nothing to validate
			continue
		}
		behavior := config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior()
		for _, requirement := range config.GetRequirementList() {
			if requirement.HasRuleId() {
				baseRequirementBehaviors[requirement.GetRuleId()] = behavior
			}
		}
	}

	for _, config := range configs {
		if !config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior() {
			// nothing to validate
			continue
		}
		for _, requirement := range config.GetRequirementList() {
			if requirement.HasExtends() {
				extendingBehavior := config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior()
				baseRuleId := requirement.GetExtends()
				if baseBehavior, ok := baseRequirementBehaviors[baseRuleId]; ok {
					if baseBehavior != extendingBehavior {
						CheckConformance_ReportInvalidRequirement(
							compiler,
							requirement,
							"extending rule's config may not specify a different value of"+
								" 'library_level_non_allowlisted_conformance_violations_behavior' than the"+
								" base rule's config. Skipping all conformance checks.")
						return false
					}
				}
			}
		}
	}
	return true
}

const checkConformance_EXTENDABLE_FIELDS_STR = "[config_file, error_message, extends, only_apply_to, only_apply_to_regexp, whitelist, whitelist_regexp, allowlist, allowlist_regexp, value]"

// checkConformance_countNonExtendableFields counts set fields on r that are not
// in CheckConformance.EXTENDABLE_FIELDS (matching iterating r.getAllFields().keySet()).
func checkConformance_countNonExtendableFields(r *Requirement) int {
	n := 0
	if r.HasType() {
		n++
	}
	if r.HasJavaClass() {
		n++
	}
	if r.HasRuleId() {
		n++
	}
	if r.HasReportLooseTypeViolations() {
		n++
	}
	if r.HasSeverity() {
		n++
	}
	if r.HasTypeMatchingStrategy() {
		n++
	}
	if len(r.GetWhitelistEntryList()) > 0 {
		n++
	}
	if r.HasAllowExtendingValue() {
		n++
	}
	if len(r.GetAllowlistEntryList()) > 0 {
		n++
	}
	return n
}

type rootRequirementEntry struct {
	builder  *Requirement
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior
}

// mergeRequirements gets requirements from all configs, validates them, and
// merges any extending requirements into their respective requirements from
// which they're extending. This means merging the allowlists/whitelists of
// requirements with 'extends' equal to 'rule_id' of other rule.
//
// Java: CheckConformance.mergeRequirements(AbstractCompiler, List<ConformanceConfig>, boolean)
func (c *CheckConformance) MergeRequirements(
	compiler AbstractCompiler,
	configs []*ConformanceConfig,
	isLibraryLevelConformanceReportingMode bool,
) []*Requirement {
	merged := c.mergeRequirements(compiler, configs, isLibraryLevelConformanceReportingMode)
	out := make([]*Requirement, len(merged))
	for i, entry := range merged {
		out[i] = entry.requirement
	}
	return out
}

func (c *CheckConformance) mergeRequirements(
	compiler AbstractCompiler,
	configs []*ConformanceConfig,
	isLibraryLevelConformanceReportingMode bool,
) []mergedRequirement {
	// Root requirements (i.e. requirements that have no 'extends' field).
	var rootRequirements []rootRequirementEntry
	rootIndex := make(map[*Requirement]int)
	putRoot := func(builder *Requirement, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
		if idx, ok := rootIndex[builder]; ok {
			rootRequirements[idx].behavior = behavior
			return
		}
		rootIndex[builder] = len(rootRequirements)
		rootRequirements = append(rootRequirements, rootRequirementEntry{builder: builder, behavior: behavior})
	}
	// Requirements that are extendable (i.e. have a 'rule_id' field).
	extendable := make(map[string]*Requirement)

	// 1. Process the root requirements and add them to the rootRequirements list, and process the
	// extendable requirements and add them to the extendable map.
	for _, config := range configs {
		for _, requirement := range config.GetRequirementList() {
			builder := requirement.Clone()
			if requirement.HasRuleId() {
				if !checkConformance_validateExtendableRequirement(compiler, requirement, extendable) {
					continue
				}
				// has a rule_id, so it's extendable
				extendable[requirement.GetRuleId()] = builder
			}
			hasBehavior := config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior()
			behavior := config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior()
			if !requirement.HasExtends() {
				// does not extend anything, so it's a root requirement. May or may not have a rule_id.
				if isLibraryLevelConformanceReportingMode && hasBehavior {
					putRoot(builder, behavior)
				} else {
					putRoot(builder, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED)
				}
			}
			// Enforce that requirements used in library-level conformance checks must have their
			// rule_id or extends fields set.
			libraryLevelConformanceEnabled :=
				isLibraryLevelConformanceReportingMode && hasBehavior &&
					behavior == ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_RECORD_ONLY
			if !requirement.HasRuleId() &&
				!requirement.HasExtends() &&
				libraryLevelConformanceEnabled {
				CheckConformance_ReportInvalidRequirement(
					compiler,
					requirement,
					"Library-level conformance requirements require their rule_id or extends field set.")
			}
		}
	}

	// 2. Process extending requirements and merge them into their respective extended requirements.
	for _, config := range configs {
		for _, requirement := range config.GetRequirementList() {
			if requirement.HasExtends() {
				extendableRequirement := extendable[requirement.GetExtends()]
				if !CheckConformance_ValidateExtendingRequirement(compiler, requirement, extendableRequirement) {
					continue
				}
				if config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior() {
					// overwrite the behavior of the extendable requirement with the behavior of the
					// extending requirement.
					putRoot(
						extendableRequirement,
						config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior())
				}
				CheckConformance_MergeExtendingRequirementIntoExtended(extendableRequirement, requirement)
			}
		}
	}

	// 3. Remove duplicates from the allowlists/whitelists of the merged root requirements.
	cleanedUpRootRequirements := make([]mergedRequirement, 0, len(rootRequirements))
	for _, entry := range rootRequirements {
		CheckConformance_RemoveDuplicates(entry.builder)
		requirement := entry.builder.Build()
		behavior := ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED
		if isLibraryLevelConformanceReportingMode {
			behavior = entry.behavior
			if compiler != nil {
				compiler.RecordLibraryLevelConformanceBehavior(requirement, behavior)
			}
		}
		cleanedUpRootRequirements = append(cleanedUpRootRequirements, mergedRequirement{
			requirement: requirement,
			behavior:    behavior,
		})
	}

	return cleanedUpRootRequirements
}

// CheckConformance_MergeExtendingRequirementIntoExtended merges the
// allowlists/whitelists of the extending requirement into the extended
// requirement.
//
// Java: CheckConformance.mergeExtendingRequirementIntoExtended(Requirement.Builder, Requirement)
func CheckConformance_MergeExtendingRequirementIntoExtended(
	extendedRequirement *Requirement, extendingRequirement *Requirement) {
	jutil.CheckStateMsg(
		extendingRequirement.GetExtends() == extendedRequirement.GetRuleId(),
		"Attempting to merge an requirement with the wrong extended requirement")
	extendedRequirement.
		AddAllWhitelist(extendingRequirement.GetWhitelistList()).
		AddAllWhitelistRegexp(extendingRequirement.GetWhitelistRegexpList()).
		AddAllAllowlist(extendingRequirement.GetAllowlistList()).
		AddAllAllowlistRegexp(extendingRequirement.GetAllowlistRegexpList()).
		AddAllOnlyApplyTo(extendingRequirement.GetOnlyApplyToList()).
		AddAllOnlyApplyToRegexp(extendingRequirement.GetOnlyApplyToRegexpList()).
		AddAllWhitelistEntry(extendingRequirement.GetWhitelistEntryList()).
		AddAllAllowlistEntry(extendingRequirement.GetAllowlistEntryList()).
		AddAllValue(extendingRequirement.GetValueList()).
		AddAllConfigFile(extendingRequirement.GetConfigFileList())
	if extendingRequirement.HasErrorMessage() && extendedRequirement.HasErrorMessage() {
		extendedRequirement.SetErrorMessage(
			extendedRequirement.GetErrorMessage() + " " + extendingRequirement.GetErrorMessage())
	}
}

// Java: CheckConformance.validateExtendableRequirement(AbstractCompiler, Requirement, Map<String, Requirement.Builder>)
func checkConformance_validateExtendableRequirement(
	compiler AbstractCompiler,
	requirement *Requirement,
	extendable map[string]*Requirement) bool {
	jutil.CheckStateMsg(requirement.HasRuleId(), "Extendable requirement must have a rule_id")

	if requirement.GetRuleId() == "" {
		CheckConformance_ReportInvalidRequirement(compiler, requirement, "empty rule_id")
		return false
	}
	if _, exists := extendable[requirement.GetRuleId()]; exists {
		CheckConformance_ReportInvalidRequirement(
			compiler,
			requirement,
			"two requirements with the same rule_id: "+requirement.GetRuleId())
		return false
	}
	return true
}

// CheckConformance_ValidateExtendingRequirement validates an extending
// requirement is correctly extending an extendable requirement.
//
// Java: CheckConformance.validateExtendingRequirement(AbstractCompiler, Requirement, Requirement.Builder)
func CheckConformance_ValidateExtendingRequirement(
	compiler AbstractCompiler,
	extendingRequirement *Requirement,
	extendableRequirement *Requirement) bool {
	jutil.CheckStateMsg(extendingRequirement.HasExtends(), "Extending requirement must have an extends")
	if extendableRequirement == nil {
		CheckConformance_ReportInvalidRequirement(
			compiler,
			extendingRequirement,
			"no extendable requirement with rule_id: "+extendingRequirement.GetExtends())
		return false
	}
	for range checkConformance_countNonExtendableFields(extendingRequirement) {
		CheckConformance_ReportInvalidRequirement(
			compiler, extendingRequirement, "extending rules allow only "+checkConformance_EXTENDABLE_FIELDS_STR)
	}
	if len(extendingRequirement.GetValueList()) > 0 &&
		!extendableRequirement.GetAllowExtendingValue() {
		CheckConformance_ReportInvalidRequirement(
			compiler,
			extendingRequirement,
			"extending rule may not specify 'value' if extendable rule does not allow it")
	}
	return true
}

func checkConformance_dedupStrings(in []string) []string {
	if len(in) <= 1 {
		return in
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// Java: CheckConformance.removeDuplicates(Requirement.Builder)
func CheckConformance_RemoveDuplicates(requirement *Requirement) {
	list1 := checkConformance_dedupStrings(requirement.GetWhitelistList())
	requirement.ClearWhitelist().AddAllWhitelist(list1)

	allowlist := checkConformance_dedupStrings(requirement.GetAllowlistList())
	requirement.ClearAllowlist().AddAllAllowlist(allowlist)

	list2 := checkConformance_dedupStrings(requirement.GetWhitelistRegexpList())
	requirement.ClearWhitelistRegexp().AddAllWhitelistRegexp(list2)

	allowlistRegexp := checkConformance_dedupStrings(requirement.GetAllowlistRegexpList())
	requirement.ClearAllowlistRegexp().AddAllAllowlistRegexp(allowlistRegexp)

	list3 := checkConformance_dedupStrings(requirement.GetOnlyApplyToList())
	requirement.ClearOnlyApplyTo().AddAllOnlyApplyTo(list3)

	list4 := checkConformance_dedupStrings(requirement.GetOnlyApplyToRegexpList())
	requirement.ClearOnlyApplyToRegexp().AddAllOnlyApplyToRegexp(list4)
}

type errInvalidRequirementSpec string

func (e errInvalidRequirementSpec) Error() string { return string(e) }

// CheckConformance_InitRule translates a requirement specified in the
// conformance config file into a rule instance.
//
// Java: CheckConformance.initRule(AbstractCompiler, Requirement)
func CheckConformance_InitRule(compiler AbstractCompiler, requirement *Requirement) CheckConformance_Rule {
	rule, err := checkConformance_initRule(compiler, requirement)
	if err != nil {
		CheckConformance_ReportInvalidRequirement(compiler, requirement, err.Error())
		return nil
	}
	return rule
}

func checkConformance_initRule(compiler AbstractCompiler, requirement *Requirement) (CheckConformance_Rule, error) {
	switch requirement.GetType() {
	case Requirement_Type_CUSTOM:
		return newConformanceRules_CustomRule(compiler, requirement)
	case Requirement_Type_NO_OP:
		return newConformanceRules_NoOp(compiler, requirement)
	case Requirement_Type_BANNED_CODE_PATTERN:
		return newConformanceRules_BannedCodePattern(compiler, requirement)
	case Requirement_Type_BANNED_DEPENDENCY:
		return newConformanceRules_BannedDependency(compiler, requirement)
	case Requirement_Type_BANNED_DEPENDENCY_REGEX:
		return newConformanceRules_BannedDependencyRegex(compiler, requirement)
	case Requirement_Type_BANNED_ENHANCE:
		return newConformanceRules_BannedEnhance(compiler, requirement)
	case Requirement_Type_BANNED_MODS_REGEX:
		return newConformanceRules_BannedModsRegex(compiler, requirement)
	case Requirement_Type_BANNED_NAME,
		Requirement_Type_BANNED_NAME_CALL,
		Requirement_Type_BANNED_SCRIPT_IMPORT:
		return newConformanceRules_BannedName(compiler, requirement)
	case Requirement_Type_BANNED_PROPERTY,
		Requirement_Type_BANNED_PROPERTY_READ,
		Requirement_Type_BANNED_PROPERTY_WRITE,
		Requirement_Type_BANNED_PROPERTY_NON_CONSTANT_WRITE,
		Requirement_Type_BANNED_PROPERTY_CALL:
		return newConformanceRules_BannedProperty(compiler, requirement)
	case Requirement_Type_RESTRICTED_NAME_CALL:
		return newConformanceRules_RestrictedNameCall(compiler, requirement)
	case Requirement_Type_RESTRICTED_METHOD_CALL:
		return newConformanceRules_RestrictedMethodCall(compiler, requirement)
	case Requirement_Type_RESTRICTED_PROPERTY_WRITE:
		return newConformanceRules_RestrictedPropertyWrite(compiler, requirement)
	case Requirement_Type_BANNED_STRING_REGEX:
		return newConformanceRules_BannedStringRegex(compiler, requirement)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(requirement.GetType().String()))
	}
}

// Java: CheckConformance.reportInvalidRequirement(AbstractCompiler, Requirement, String)
func CheckConformance_ReportInvalidRequirement(
	compiler AbstractCompiler, requirement *Requirement, reason string) {
	compiler.Report(
		JSError_Make(
			// TextFormat.printer().printToString(requirement) is Message.toString().
			CheckConformance_INVALID_REQUIREMENT_SPEC, reason, requirement.String()))
}
