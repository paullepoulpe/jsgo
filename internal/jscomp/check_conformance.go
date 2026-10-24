// Port of com/google/javascript/jscomp/CheckConformance.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	compiler   AbstractCompiler
	categories []*CheckConformance_Category

	// Map of root requirements to their behavior specified in their configs, or of their extending
	// requirement's config. It is only populated if both of the following conditions are met:
	//
	//  1. the conformance reporting mode is RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG,
	//     i.e. the conformance checks are being run during the CheckJS action.
	//  2. the requirement has a behavior that is different from the default UNSPECIFIED behavior
	//     (i.e. RECORD_ONLY or REPORT_AS_BUILD_ERROR).
	//
	// The behavior of each library-level requirement gets stored in this map, and gets passed into
	// the check() method of its corresponding rule to determine whether to record or report the
	// conformance violations.
	//
	// Go-only: language-adapter — Requirement keys compare with the proto's
	// equals()/hashCode() (ObjectLinkedHashMap).
	mergedBehaviors *jutil.ObjectLinkedHashMap[*Requirement, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]

	// Map of rules to their behavior. Populated from the mergedBehaviors map. The behavior of each
	// rule is passed into the check() method of the rule to determine whether to record or report the
	// conformance violations.
	ruleToBehavior *jutil.LinkedHashMap[CheckConformance_Rule, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]
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

// CheckConformance_RuleDefaults holds Rule's default methods; implementations
// embed it.
//
// Go-only: language-adapter — Java interface default methods.
type CheckConformance_RuleDefaults struct{}

// Java: CheckConformance.Rule.getPrecondition() (default method)
func (CheckConformance_RuleDefaults) GetPrecondition() CheckConformance_Precondition {
	return CheckConformance_Precondition_CHECK_ALL
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

// CheckConformance_PreconditionFunc is the adapter for Precondition lambdas.
// Lambdas are used as map keys (by identity), so the adapter is used through
// a pointer: a Go func value is not hashable.
//
// Go-only: language-adapter
type CheckConformance_PreconditionFunc func(n rhino.Node) bool

// ShouldCheck calls f.
func (f *CheckConformance_PreconditionFunc) ShouldCheck(n rhino.Node) bool { return (*f)(n) }

// NewCheckConformance_PreconditionFunc allocates a Precondition lambda
// instance (one per Java lambda instance).
//
// Go-only: language-adapter
func NewCheckConformance_PreconditionFunc(f func(n rhino.Node) bool) *CheckConformance_PreconditionFunc {
	p := CheckConformance_PreconditionFunc(f)
	return &p
}

// checkConformance_checkAll is the anonymous Precondition of CHECK_ALL.
type checkConformance_checkAll struct {
	// Go-only: language-adapter — distinct identity (a map key).
	_ byte
}

// Java: CheckConformance.Precondition.CHECK_ALL.new Precondition() {...}.shouldCheck(Node)
func (*checkConformance_checkAll) ShouldCheck(n rhino.Node) bool {
	return true
}

// Java: CheckConformance.Precondition.CHECK_ALL
var CheckConformance_Precondition_CHECK_ALL CheckConformance_Precondition = &checkConformance_checkAll{}

// Java: CheckConformance.Precondition.IS_CONSTRUCTOR_OR_CLASS
var CheckConformance_Precondition_IS_CONSTRUCTOR_OR_CLASS CheckConformance_Precondition = NewCheckConformance_PreconditionFunc(
	func(n rhino.Node) bool { return n.IsClass() || NodeUtil_IsConstructor(n) })

// Java: CheckConformance.Category
type CheckConformance_Category struct {
	precondition CheckConformance_Precondition
	rules        []CheckConformance_Rule
}

// Java: CheckConformance.Category(Precondition, ImmutableList<Rule>)
func NewCheckConformance_Category(precondition CheckConformance_Precondition, rules []CheckConformance_Rule) *CheckConformance_Category {
	c := &CheckConformance_Category{}
	c.precondition = precondition
	c.rules = rules
	return c
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
	c := &CheckConformance{}
	c.mergedBehaviors = jutil.NewObjectLinkedHashMap[*Requirement, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]()
	c.ruleToBehavior = jutil.NewLinkedHashMap[CheckConformance_Rule, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]()
	c.compiler = compiler
	// Initialize the map of functions to inspect for renaming candidates.
	c.categories = c.InitRules(compiler, configs, reportingMode)
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
	return !n.IsScript() ||
		(c.IsScriptOfInterest(t.GetInput().GetSourceFile()) &&
			!strings.HasSuffix(*t.GetSourceName(), "tsmes.closure.js"))
}

// Java: CheckConformance.isScriptOfInterest(SourceFile)
func (c *CheckConformance) IsScriptOfInterest(sf *SourceFile) bool {
	return !sf.IsWeak() && !sf.IsExtern()
}

// Java: CheckConformance.visit(NodeTraversal, Node, Node)
func (c *CheckConformance) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	/*
	 * Use counted loops and backward iteration for performance.
	 *
	 * <p>These loops are run a huge number of times. The overhead of enhanced-for loops and even
	 * calling size() can add seconds of build time to large projects.
	 */
	for cat := len(c.categories) - 1; cat >= 0; cat-- {
		category := c.categories[cat]
		if category.precondition.ShouldCheck(n) {
			for r := len(category.rules) - 1; r >= 0; r-- {
				rule := category.rules[r]
				behavior := c.ruleToBehavior.GetOrDefault(rule, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED)
				rule.Check(t, n, behavior)
			}
		}
	}
}

// InitRules builds the data structures need by this pass from the provided
// configurations.
//
// Java: CheckConformance.initRules(AbstractCompiler, ImmutableList<ConformanceConfig>, ConformanceReportingMode)
func (c *CheckConformance) InitRules(
	compiler AbstractCompiler,
	configs []*ConformanceConfig,
	reportingMode *CompilerOptions_ConformanceReportingMode) []*CheckConformance_Category {
	// Go-only: language-adapter — HashMultimap<Precondition, Rule>: both keys
	// and values hash by identity, so Java's iteration order is not
	// deterministic; this keeps insertion order.
	builder := jutil.NewLinkedHashMap[CheckConformance_Precondition, *jutil.LinkedHashSet[CheckConformance_Rule]]()

	isLibraryLevelReportingMode :=
		reportingMode != nil &&
			*reportingMode == CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	if isLibraryLevelReportingMode {
		if !CheckConformance_ValidateBehaviorSettingOfConfigs(compiler, configs, *reportingMode) {
			return []*CheckConformance_Category{}
		}
	}

	mergedRequirements := c.MergeRequirements(compiler, configs, isLibraryLevelReportingMode)
	for _, requirement := range mergedRequirements {
		behavior := ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED
		rule := CheckConformance_InitRule(compiler, requirement)
		if rule != nil {
			precondition := rule.GetPrecondition()
			rules, ok := builder.Lookup(precondition)
			if !ok {
				rules = jutil.NewLinkedHashSet[CheckConformance_Rule]()
				builder.Put(precondition, rules)
			}
			rules.Add(rule)
			if isLibraryLevelReportingMode && c.mergedBehaviors.ContainsKey(requirement) {
				behavior = c.mergedBehaviors.Get(requirement)
				// checkNotNull(behavior, "The mergeBehaviors method inserted a null behavior in the
				// mergedBehaviors map"): a Go enum value is never null.
			}
			c.ruleToBehavior.Put(rule, behavior)
		}
	}
	categories := []*CheckConformance_Category{}
	for precondition, rules := range builder.All() {
		categories = append(categories, NewCheckConformance_Category(precondition, rules.ToSlice()))
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
	baseRequirementBehaviors := jutil.NewLinkedHashMap[string, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]()
	for _, config := range configs {
		if !config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior() {
			// nothing to validate
			continue
		}
		behavior := config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior()
		for _, requirement := range config.GetRequirementList() {
			if requirement.HasRuleId() {
				baseRequirementBehaviors.Put(requirement.GetRuleId(), behavior)
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
				if baseRequirementBehaviors.ContainsKey(baseRuleId) {
					baseBehavior := baseRequirementBehaviors.Get(baseRuleId)
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

// Java: CheckConformance.EXTENDABLE_FIELDS
var checkConformance_EXTENDABLE_FIELDS = jutil.NewLinkedHashSet(
	"config_file",
	"extends",
	"only_apply_to",
	"only_apply_to_regexp",
	"whitelist",
	"whitelist_regexp",
	"allowlist",
	"allowlist_regexp",
	"value")

// MergeRequirements gets requirements from all configs, validates them, and
// merges any extending requirements into their respective requirements from
// which they're extending. This means merging the allowlists/whitelists of
// requirements with 'extends' equal to 'rule_id' of other rule.
//
// The requirements inheritance can not have a chain of more than 1 level.
// This means that a requirement that extends another requirement can not
// itself be extended. This is enforced by this pass and reported as an error
// if a requirement has both 'extends' and 'rule_id' fields set. Given that,
// we can safely define the following types of requirements:
//
// Root requirements are simply the requirements that do not extend any other
// requirement. They do not have an `extends` field. They may or may not have
// a 'rule_id' field. For example, without a 'rule_id' field:
//
//	requirement: {
//	  type: BANNED_NAME
//	  value: 'eval'
//	  error_message: 'eval is not allowed'
//	}
//
// Root requirements may have a 'rule_id' field. For example:
//
//	requirement: {
//	  rule_id: 'eval'
//	  type: BANNED_NAME
//	  value: 'eval'
//	  error_message: 'eval is not allowed'
//	}
//
// Extendable root requirements are those root requirements that have a
// 'rule_id' field, and therefore _can_ be extended. They may or may not
// getting extended by another requirement. For example:
//
//	requirement: {
//	  rule_id: 'extendable'
//	  type: BANNED_NAME
//	  value: 'eval'
//	  error_message: 'eval is not allowed'
//	}
//
// Extended root requirements are extendable root requirements that are
// extended by at least one other requirement. For example:
//
//	requirement: {
//	  rule_id: 'extended'
//	  type: BANNED_NAME
//	  value: 'eval'
//	  error_message: 'eval is not allowed'
//	}
//
//	requirement: {
//	  extends: 'extended'
//	}
//
// Extending requirements are leaf-requirements that have an 'extends' field.
// For example:
//
//	requirement: {
//	  extends: 'eval'
//
// The following diagram illustrates the relationship between the different
// types of requirements that can be defined.
//
//	        +---------------------------+    +------------------------------+
//	        |    Root Requirements      |    |  Extending Requirements      |
//	        |                           |    | (leaf, with 'extends' field) |
//	        +---------------------------+    +------------------------------+
//	           /                        \
//	          /                          \
//	+-------------------------------+      +-----------------------------------+
//	|Extendable Root Requirements   |      | Non-Extendable Root Requirements  |
//	|     (with rule_id)            |      |        (without rule_id)          |
//	+-------------------------------+      +-----------------------------------+
//	          |                      \
//	          |                       \
//	+----------------------------+    +--------------------------------+
//	| Extended Root Requirements |    | Non-Extended Root Requirements |
//	+----------------------------+    +--------------------------------+
//
// Here's the algorithm:
//
//  1. Process the root requirements and add them to the rootRequirements
//     list, and process the extendable requirements and add them to the
//     extendable map.
//  2. Process extending-leaf requirements and merge them into their
//     respective extended root requirements.
//  3. Remove duplicates from the allowlists/whitelists of the merged extended
//     root requirements.
//
// Returns a list of root requirements of all configs after merging.
//
// Java: CheckConformance.mergeRequirements(AbstractCompiler, List<ConformanceConfig>, boolean)
func (c *CheckConformance) MergeRequirements(
	compiler AbstractCompiler,
	configs []*ConformanceConfig,
	isLibraryLevelConformanceReportingMode bool) []*Requirement {
	// Root requirements (i.e. requirements that have no 'extends' field).
	rootRequirements := jutil.NewLinkedHashMap[*Requirement_Builder, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]()
	// Requirements that are extendable (i.e. have a 'rule_id' field).
	extendable := jutil.NewLinkedHashMap[string, *Requirement_Builder]()

	// 1. Process the root requirements and add them to the rootRequirements list, and process the
	// extendable requirements and add them to the extendable map.
	for _, config := range configs {
		for _, requirement := range config.GetRequirementList() {
			builder := requirement.ToBuilder()
			if requirement.HasRuleId() {
				if !CheckConformance_ValidateExtendableRequirement(compiler, requirement, extendable) {
					continue
				}
				// has a rule_id, so it's extendable
				extendable.Put(requirement.GetRuleId(), builder)
			}
			hasBehavior := config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior()
			behavior := config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior()
			if !requirement.HasExtends() {
				// does not extend anything, so it's a root requirement. May or may not have a rule_id.
				if !isLibraryLevelConformanceReportingMode {
					rootRequirements.Put(builder, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED)
				} else {
					if hasBehavior {
						rootRequirements.Put(builder, behavior)
					} else {
						rootRequirements.Put(builder, ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED)
					}
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
				extendableRequirement := extendable.Get(requirement.GetExtends())
				if !CheckConformance_ValidateExtendingRequirement(compiler, requirement, extendableRequirement) {
					continue
				}
				if config.HasLibraryLevelNonAllowlistedConformanceViolationsBehavior() {
					// overwrite the behavior of the extendable requirement with the behavior of the
					// extending requirement.
					rootRequirements.Put(
						extendableRequirement,
						config.GetLibraryLevelNonAllowlistedConformanceViolationsBehavior())
				}
				CheckConformance_MergeExtendingRequirementIntoExtended(extendableRequirement, requirement)
			}
		}
	}

	// 3. Remove duplicates from the allowlists/whitelists of the merged root requirements.
	cleanedUpRootRequirements := make([]*Requirement, 0, rootRequirements.Size())
	for builder, behavior := range rootRequirements.All() {
		CheckConformance_RemoveDuplicates(builder)
		requirement := builder.Build()
		if isLibraryLevelConformanceReportingMode &&
			behavior != ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_UNSPECIFIED {
			c.mergedBehaviors.Put(requirement, behavior)
		}
		cleanedUpRootRequirements = append(cleanedUpRootRequirements, requirement)
	}

	return cleanedUpRootRequirements
}

// CheckConformance_MergeExtendingRequirementIntoExtended merges the
// allowlists/whitelists of the extending requirement into the extended
// requirement.
//
// extendedRequirement is the extended requirement to merge the extending
// requirement into; extendingRequirement the extending requirement to merge
// into the extended requirement.
//
// Java: CheckConformance.mergeExtendingRequirementIntoExtended(Requirement.Builder, Requirement)
func CheckConformance_MergeExtendingRequirementIntoExtended(
	extendedRequirement *Requirement_Builder, extendingRequirement *Requirement) {
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
}

// CheckConformance_ValidateExtendableRequirement validates a extendable
// requirement.
//
// A valid extendable requirement is a requirement that has a unique,
// non-empty rule_id.
//
// Returns true if the requirement is valid, false otherwise.
//
// Java: CheckConformance.validateExtendableRequirement(AbstractCompiler, Requirement, Map<String, Requirement.Builder>)
func CheckConformance_ValidateExtendableRequirement(
	compiler AbstractCompiler,
	requirement *Requirement,
	extendable *jutil.LinkedHashMap[string, *Requirement_Builder]) bool {
	jutil.CheckStateMsg(requirement.HasRuleId(), "Extendable requirement must have a rule_id")

	if requirement.GetRuleId() == "" {
		CheckConformance_ReportInvalidRequirement(compiler, requirement, "empty rule_id")
		return false
	}
	if extendable.ContainsKey(requirement.GetRuleId()) {
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
// A valid extending requirement is a requirement that:
//
//  1. has a non-empty extends field, and
//  2. only extends an existing extendable requirement
//  3. does not specify any fields that are not in EXTENDABLE_FIELDS.
//  4. does not specify a value if the extendable requirement does not allow
//     extending values.
//
// Returns true if the requirement is valid, false otherwise.
//
// Java: CheckConformance.validateExtendingRequirement(AbstractCompiler, Requirement, Requirement.Builder)
func CheckConformance_ValidateExtendingRequirement(
	compiler AbstractCompiler,
	extendingRequirement *Requirement,
	extendableRequirement *Requirement_Builder) bool {
	jutil.CheckStateMsg(extendingRequirement.HasExtends(), "Extending requirement must have an extends")
	if extendableRequirement == nil {
		CheckConformance_ReportInvalidRequirement(
			compiler,
			extendingRequirement,
			"no extendable requirement with rule_id: "+extendingRequirement.GetExtends())
		return false
	}
	for field := range extendingRequirement.GetAllFields().Keys() {
		if !checkConformance_EXTENDABLE_FIELDS.Contains(field.GetName()) {
			CheckConformance_ReportInvalidRequirement(
				compiler, extendingRequirement, "extending rules allow only "+checkConformance_EXTENDABLE_FIELDS.String())
		}
	}
	if extendingRequirement.GetValueCount() > 0 &&
		!extendableRequirement.GetAllowExtendingValue() {
		CheckConformance_ReportInvalidRequirement(
			compiler,
			extendingRequirement,
			"extending rule may not specify 'value' if extendable rule does not allow it")
	}
	return true
}

// Java: CheckConformance.removeDuplicates(Requirement.Builder)
func CheckConformance_RemoveDuplicates(requirement *Requirement_Builder) {
	list1 := jutil.NewLinkedHashSet(requirement.GetWhitelistList()...).ToSlice()
	requirement.ClearWhitelist().AddAllWhitelist(list1)

	allowlist := jutil.NewLinkedHashSet(requirement.GetAllowlistList()...).ToSlice()
	requirement.ClearAllowlist().AddAllAllowlist(allowlist)

	list2 := jutil.NewLinkedHashSet(requirement.GetWhitelistRegexpList()...).ToSlice()
	requirement.ClearWhitelistRegexp().AddAllWhitelistRegexp(list2)

	allowlistRegexp := jutil.NewLinkedHashSet(requirement.GetAllowlistRegexpList()...).ToSlice()
	requirement.ClearAllowlistRegexp().AddAllAllowlistRegexp(allowlistRegexp)

	list3 := jutil.NewLinkedHashSet(requirement.GetOnlyApplyToList()...).ToSlice()
	requirement.ClearOnlyApplyTo().AddAllOnlyApplyTo(list3)

	list4 := jutil.NewLinkedHashSet(requirement.GetOnlyApplyToRegexpList()...).ToSlice()
	requirement.ClearOnlyApplyToRegexp().AddAllOnlyApplyToRegexp(list4)
}

// CheckConformance_InitRule translates a requirement specified in the
// conformance config file into a rule instance.
//
// compiler is the compiler to use for the rule; requirement the requirement
// to initialize the rule for. Returns the initialized rule, or null if the
// requirement is invalid.
//
// Java: CheckConformance.initRule(AbstractCompiler, Requirement)
func CheckConformance_InitRule(compiler AbstractCompiler, requirement *Requirement) CheckConformance_Rule {
	var result CheckConformance_Rule
	jutil.Try(func() {
		switch requirement.GetType() {
		case Requirement_Type_CUSTOM:
			result = NewConformanceRules_CustomRuleProxy(compiler, requirement)
		case Requirement_Type_NO_OP:
			result = NewConformanceRules_NoOp(compiler, requirement)
		case Requirement_Type_BANNED_CODE_PATTERN:
			result = NewConformanceRules_BannedCodePattern(compiler, requirement)
		case Requirement_Type_BANNED_DEPENDENCY:
			result = NewConformanceRules_BannedDependency(compiler, requirement)
		case Requirement_Type_BANNED_DEPENDENCY_REGEX:
			result = NewConformanceRules_BannedDependencyRegex(compiler, requirement)
		case Requirement_Type_BANNED_ENHANCE:
			result = NewConformanceRules_BannedEnhance(compiler, requirement)
		case Requirement_Type_BANNED_MODS_REGEX:
			result = NewConformanceRules_BannedModsRegex(compiler, requirement)
		case Requirement_Type_BANNED_NAME,
			Requirement_Type_BANNED_NAME_CALL:
			result = NewConformanceRules_BannedName(compiler, requirement)
		case Requirement_Type_BANNED_PROPERTY,
			Requirement_Type_BANNED_PROPERTY_READ,
			Requirement_Type_BANNED_PROPERTY_WRITE,
			Requirement_Type_BANNED_PROPERTY_NON_CONSTANT_WRITE,
			Requirement_Type_BANNED_PROPERTY_CALL:
			result = NewConformanceRules_BannedProperty(compiler, requirement)
		case Requirement_Type_RESTRICTED_NAME_CALL:
			result = NewConformanceRules_RestrictedNameCall(compiler, requirement)
		case Requirement_Type_RESTRICTED_METHOD_CALL:
			result = NewConformanceRules_RestrictedMethodCall(compiler, requirement)
		case Requirement_Type_RESTRICTED_PROPERTY_WRITE:
			result = NewConformanceRules_RestrictedPropertyWrite(compiler, requirement)
		case Requirement_Type_BANNED_STRING_REGEX:
			result = NewConformanceRules_BannedStringRegex(compiler, requirement)
		default:
			// A Java switch expression over all enum constants: an unknown value
			// throws IncompatibleClassChangeError.
			panic(jutil.NewIllegalStateExceptionMsg(requirement.GetType().String()))
		}
	}, func(e ICheckConformance_InvalidRequirementSpec) {
		CheckConformance_ReportInvalidRequirement(compiler, requirement, e.GetMessage())
		result = nil
	})
	return result
}

// CheckConformance_InvalidRequirementSpec is a checked exception.
type CheckConformance_InvalidRequirementSpec struct {
	jutil.JavaException
}

func (*CheckConformance_InvalidRequirementSpec) invalidRequirementSpec() {}

// ICheckConformance_InvalidRequirementSpec matches InvalidRequirementSpec.
type ICheckConformance_InvalidRequirementSpec interface {
	jutil.IException
	invalidRequirementSpec()
}

// Java: CheckConformance.InvalidRequirementSpec(String)
func NewCheckConformance_InvalidRequirementSpec(message string) *CheckConformance_InvalidRequirementSpec {
	e := &CheckConformance_InvalidRequirementSpec{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.CheckConformance$InvalidRequirementSpec", jutil.StringPtr(message), nil)
	return e
}

// Java: CheckConformance.InvalidRequirementSpec(String, Throwable)
func NewCheckConformance_InvalidRequirementSpecWithCause(message string, cause jutil.Throwable) *CheckConformance_InvalidRequirementSpec {
	e := &CheckConformance_InvalidRequirementSpec{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.CheckConformance$InvalidRequirementSpec", jutil.StringPtr(message), cause)
	return e
}

// Java: CheckConformance.reportInvalidRequirement(AbstractCompiler, Requirement, String)
func CheckConformance_ReportInvalidRequirement(
	compiler AbstractCompiler, requirement *Requirement, reason string) {
	compiler.Report(
		JSError_Make(
			// TextFormat.printer().printToString(requirement) is Message.toString().
			CheckConformance_INVALID_REQUIREMENT_SPEC, reason, requirement.String()))
}
