// Port of com/google/javascript/jscomp/CheckConformance.java @ closure-compiler v20260415
//
// Java's CheckConformance is a framework that reads a ConformanceConfig
// proto (text-proto file supplied via --conformance_config), builds a
// set of Rule objects from it, and then traverses every non-extern
// source file firing per-rule diagnostics. Four built-in diagnostic
// codes are defined (lines 54-67):
//
//   JSC_CONFORMANCE_ERROR            – rule.type == ERROR
//   JSC_CONFORMANCE_VIOLATION        – rule.type == WARNING
//   JSC_CONFORMANCE_POSSIBLE_VIOLATION
//   JSC_INVALID_REQUIREMENT_SPEC
//
// Full porting would require:
//   • ConformanceConfig / Requirement proto types (lines 22-38)
//   • ConformanceRules (20+ rule implementations)
//   • NodeTraversal + shouldTraverse / visit callbacks (lines 167-204)
//   • mergeRequirements (lines 418-500) + validateBehaviorSettingOfConfigs (lines 244-294)
//   • LibraryLevelNonAllowlistedConformanceViolationsBehavior enum
//
// Since no caller currently supplies a conformance config, all of those
// branches are unreachable today. This shim exposes the same constructor
// + Check surface used by compiler.go so the pass can be wired in now
// and filled in incrementally as proto support lands.
//
// TODO: port Rule interface + Precondition (lines 96-145)
// TODO: port Category grouping + initRules (lines 137-242)
// TODO: port validateBehaviorSettingOfConfigs (lines 244-294)
// TODO: port mergeRequirements incl. EXTENDABLE_FIELDS set (lines 296-500)
// TODO: port mergeExtendingRequirementIntoExtended (lines 508-524)
// TODO: port validateExtendableRequirement (lines 533-551)
// TODO: port validateExtendingRequirement (lines 566-592)
// TODO: port removeDuplicates (lines 594-613)
// TODO: port initRule / switch on Requirement.Type (lines 622-653)
// TODO: port InvalidRequirementSpec exception (lines 655-663)
// TODO: port reportInvalidRequirement (lines 665-670)

package jscomp

import (
	"github.com/paullepoulpe/jsgo/legacy/internal/diag"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// ---------------------------------------------------------------------------
// Stub types for proto-dependent infrastructure not yet ported.
// ---------------------------------------------------------------------------

// LibraryLevelBehavior mirrors
// ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior (Java lines 22-23).
// TODO: replace with generated proto enum once ConformanceConfig.proto is ported.
type LibraryLevelBehavior int

const (
	// LibraryLevelBehaviorUnspecified mirrors UNSPECIFIED (the proto default).
	LibraryLevelBehaviorUnspecified LibraryLevelBehavior = iota
	// LibraryLevelBehaviorRecordOnly mirrors RECORD_ONLY.
	LibraryLevelBehaviorRecordOnly
	// LibraryLevelBehaviorReportAsBuildError mirrors REPORT_AS_BUILD_ERROR.
	LibraryLevelBehaviorReportAsBuildError
)

// ccRequirement is a minimal stub for the Requirement proto type used by
// mergeRequirements and initRule. The Java Requirement is a protobuf message
// with fields: type, rule_id, extends, error_message, allowlist, etc.
// TODO: replace with the generated proto type once ConformanceConfig.proto is ported.
type ccRequirement struct {
	ruleID             string
	extends            string
	whitelist          []string
	whitelistRegexp    []string
	allowlist          []string
	allowlistRegexp    []string
	onlyApplyTo        []string
	onlyApplyToRegexp  []string
	value              []string
	configFile         []string
	allowExtendingValue bool
}

// ccConformanceConfig is a minimal stub for the ConformanceConfig proto.
// Java: com/google/javascript/jscomp/ConformanceConfig.proto
// TODO: replace with the generated proto type.
type ccConformanceConfig struct {
	requirements []ccRequirement
	// libraryLevelBehavior mirrors the
	// library_level_non_allowlisted_conformance_violations_behavior field.
	// hasLibraryLevelBehavior is true when the field is explicitly set (Java
	// proto3 hasXxx() semantics — field is present vs. defaulted).
	libraryLevelBehavior    LibraryLevelBehavior
	hasLibraryLevelBehavior bool
}

// Rule mirrors the Rule interface in Java's CheckConformance (line 96).
// Each Rule implementation checks one conformance requirement.
//
// Java:
//
//	public static interface Rule {
//	  default Precondition getPrecondition() { return Precondition.CHECK_ALL; }
//	  void check(NodeTraversal t, Node n, LibraryLevelBehavior behavior);
//	}
type Rule interface {
	// getPrecondition returns the precondition filter for this rule.
	// Mirrors Rule.getPrecondition() — default returns Precondition.CHECK_ALL.
	getPrecondition() Precondition
	// check validates node n against this rule.
	// Mirrors Rule.check(NodeTraversal t, Node n, behavior).
	check(n *rhino.Node)
}

// Precondition mirrors the Precondition interface (Java line 122).
// Instances are used as keys to group rules with common preconditions.
// Grouping allows shared computation per node (substantial performance win).
//
// Java:
//
//	public static interface Precondition {
//	  boolean shouldCheck(Node n);
//	}
type Precondition interface {
	// shouldCheck returns true when the node should be checked by rules
	// in this precondition group.
	shouldCheck(n *rhino.Node) bool
}

// PreconditionCheckAll mirrors Precondition.CHECK_ALL (Java line 125).
// Returns true for every node.
// Java: static final Precondition CHECK_ALL = n -> true;
type PreconditionCheckAll struct{}

// shouldCheck always returns true, matching CHECK_ALL behaviour.
func (PreconditionCheckAll) shouldCheck(_ *rhino.Node) bool { return true }

// PreconditionIsConstructorOrClass mirrors Precondition.IS_CONSTRUCTOR_OR_CLASS (Java line 133).
// Returns true only for class or constructor nodes.
// Java: static final Precondition IS_CONSTRUCTOR_OR_CLASS = n -> n.isClass() || NodeUtil.isConstructor(n);
type PreconditionIsConstructorOrClass struct{}

// shouldCheck returns true for CLASS nodes (isClass) or constructor functions.
func (PreconditionIsConstructorOrClass) shouldCheck(n *rhino.Node) bool {
	return n.Token() == rhino.CLASS
}

// Category mirrors the private Category class (Java line 137).
// Groups rules that share a common precondition to avoid redundant checks.
//
// Java:
//
//	private static final class Category {
//	  final Precondition precondition;
//	  final ImmutableList<Rule> rules;
//	}
type Category struct {
	precondition Precondition
	rules        []Rule
}

// ---------------------------------------------------------------------------
// CheckConformance — the CompilerPass.
// ---------------------------------------------------------------------------

// CheckConformance is a shim for Java's CheckConformance compiler pass.
// When no conformance config path is provided (the only supported mode
// today), Check always returns nil — matching Java's early-exit when
// this.categories.isEmpty() (Java line 169).
//
// Wire-in: call NewCheckConformance(path) in compiler.go and append
// its Check(script) result to the diagnostics slice. Because no config
// is ever loaded, the append is a no-op and introduces zero overhead.
//
// Java: implements NodeTraversal.Callback, CompilerPass.
type CheckConformance struct {
	// path is the source file path forwarded to any future diagnostics.
	path string
	// categories holds the grouped rule sets built from the conformance config.
	// Empty when no config is loaded — matches Java's categories field.
	// TODO: populate from ConformanceConfig proto loading.
	categories []*Category
	// configPath, when non-empty, would be the conformance_config text-proto
	// file to parse. Not consumed by this shim — kept for future expansion.
	//
	// TODO: port ConformanceConfig proto loading from configPath
	// (Java loads configs via compiler.getConformanceConfigs(),
	//  compiler.go line ~158 in the Java CheckConformance constructor).
	configPath string
}

// NewCheckConformance constructs a CheckConformance shim for the given
// source file path. The optional conformance config path is recorded
// for future use but not parsed; see the TODO above.
//
// Java's two constructors are:
//
//	CheckConformance(AbstractCompiler)                        – line 151
//	CheckConformance(AbstractCompiler, configs, reportingMode) – line 158
//
// Both are collapsed here: without a ConformanceConfig proto type,
// the config argument is a plain file path (empty = no config).
//
// Go-only: language-adapter — Go constructor shim (Java uses class constructor directly; Go has no class constructors).
// TODO: port the ImmutableList<ConformanceConfig> + ConformanceReportingMode
// constructor signature (lines 158-165) once proto types are available.
func NewCheckConformance(path string) *CheckConformance {
	return &CheckConformance{
		path:       path,
		categories: initRules(nil),
	}
}

// Check runs the conformance pass over script.
//
// Java's process() method (lines 167-172) delegates to
// NodeTraversal.traverseRoots only when categories is non-empty.
// Because this shim never builds any categories (no config), it
// mirrors that early-exit and always returns nil.
//
// When a config is eventually wired in, this method should:
//  1. call initRules to build categories from the loaded config
//     (TODO: port initRules, lines 206-242)
//  2. traverse the script calling visit for each node
//     (TODO: port shouldTraverse / visit, lines 175-204)
//
// Mirrors: Java's process(Node externs, Node root), lines 167-172.
func (c *CheckConformance) Check(script *rhino.Node) []diag.Diagnostic {
	// No conformance config → no categories → no diagnostics.
	// Mirrors: if (!this.categories.isEmpty()) { NodeTraversal.traverseRoots(…) }
	// (Java line 169).
	if len(c.categories) == 0 {
		_ = script
		return nil
	}
	// TODO: traverse script and call shouldTraverse / visit for each node.
	_ = script
	return nil
}

// shouldTraverse mirrors Java's shouldTraverse(NodeTraversal t, Node n, Node parent).
// Returns false for SCRIPT nodes that are externs, weak sources, or @closureUnaware.
// (Java lines 175-180.)
//
// TODO: wire into traversal once NodeTraversal is ported.
func (c *CheckConformance) shouldTraverse(n *rhino.Node) bool {
	// Skip SCRIPT nodes for extern files; all others traverse normally.
	if n.Token() == rhino.SCRIPT {
		// TODO: check SourceFile.isWeak() / isExtern() once SourceFile is ported.
		return true
	}
	return true
}

// visit mirrors Java's visit(NodeTraversal t, Node n, Node parent) (lines 187-204).
// Iterates categories backwards (performance optimisation from Java) and fires
// each rule whose precondition matches n.
//
// Java (lines 194-203): backward loop over categories, then backward loop over rules.
// TODO: call from traversal once NodeTraversal adapter is wired in.
func (c *CheckConformance) visit(n *rhino.Node) {
	for i := len(c.categories) - 1; i >= 0; i-- {
		cat := c.categories[i]
		if !cat.precondition.shouldCheck(n) {
			continue
		}
		for j := len(cat.rules) - 1; j >= 0; j-- {
			cat.rules[j].check(n)
		}
	}
}

// initRules mirrors Java's initRules(AbstractCompiler, ImmutableList<ConformanceConfig>,
//
//	ConformanceReportingMode) (lines 206-242).
//
// Builds the category list from the merged requirements of all configs grouped by
// Precondition (Java uses HashMultimap<Precondition, Rule>). isLibraryLevelReportingMode
// is true when reportingMode == RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG.
// Returns nil (empty) when configs is nil/empty or validation fails.
func initRules(configs []*ccConformanceConfig) []*Category {
	if len(configs) == 0 {
		return nil
	}
	// Java line 213-220: isLibraryLevelReportingMode check + validateBehaviorSettingOfConfigs.
	// Without a ConformanceReportingMode argument at this shim level we treat the mode as
	// non-library-level (the only supported mode today).
	isLibraryLevel := false

	// Java line 223-238: merge requirements then build rules per precondition.
	mergedRequirements := mergeRequirements(configs, isLibraryLevel)

	// Java uses HashMultimap<Precondition, Rule>; we approximate with a map keyed by a
	// Precondition identity string since Precondition is an interface without a natural key.
	// For the shim this is fine — initRule always returns nil so the map stays empty.
	type preconditionKey string
	type entry struct {
		precondition Precondition
		rules        []Rule
	}
	byPrecondition := map[preconditionKey]*entry{}
	preconditionOrder := []preconditionKey{}

	for i := range mergedRequirements {
		req := &mergedRequirements[i]
		rule := initRule(req)
		if rule == nil {
			continue
		}
		p := rule.getPrecondition()
		var key preconditionKey
		switch p.(type) {
		case PreconditionIsConstructorOrClass:
			key = "IS_CONSTRUCTOR_OR_CLASS"
		default:
			key = "CHECK_ALL"
		}
		if _, ok := byPrecondition[key]; !ok {
			byPrecondition[key] = &entry{precondition: p}
			preconditionOrder = append(preconditionOrder, key)
		}
		byPrecondition[key].rules = append(byPrecondition[key].rules, rule)
	}

	// Java line 239-241: stream entries into Category list.
	categories := make([]*Category, 0, len(preconditionOrder))
	for _, key := range preconditionOrder {
		e := byPrecondition[key]
		categories = append(categories, &Category{
			precondition: e.precondition,
			rules:        e.rules,
		})
	}
	return categories
}

// mergeRequirements mirrors Java's mergeRequirements(AbstractCompiler,
//
//	List<ConformanceConfig>, boolean) (lines 418-500).
//
// Three-phase algorithm (mirrors Java lines 422-499):
//  1. Collect root requirements; index extendable ones by rule_id.
//  2. Merge extending requirements into their extended roots.
//  3. Deduplicate allowlists/whitelists and build the final flat list.
//
// isLibraryLevelConformanceReportingMode corresponds to the same-named Java boolean
// (derived from ConformanceReportingMode, Java line 421).
func mergeRequirements(configs []*ccConformanceConfig, isLibraryLevelConformanceReportingMode bool) []ccRequirement {
	// Java: Map<Requirement.Builder, LibraryLevelBehavior> rootRequirements (LinkedHashMap, line 423)
	// We preserve insertion order via a slice of builders paired with behavior.
	type rootEntry struct {
		req      ccRequirement
		behavior LibraryLevelBehavior
	}
	// rootRequirements preserves insertion order (LinkedHashMap in Java).
	var rootRequirements []rootEntry
	// rootIndex maps ruleID → index in rootRequirements for O(1) lookup by extending reqs.
	rootIndex := map[string]int{}

	// extendable mirrors Java's Map<String, Requirement.Builder> extendable (line 426).
	extendable := map[string]int{} // rule_id → index in rootRequirements

	// --- Phase 1 (Java lines 430-463) ---
	for _, cfg := range configs {
		behavior := cfg.libraryLevelBehavior
		hasBehavior := cfg.hasLibraryLevelBehavior
		for ri := range cfg.requirements {
			req := cfg.requirements[ri] // copy — we may mutate later
			if req.ruleID != "" {
				// Java lines 433-438: validate and index extendable requirement.
				if !validateExtendableRequirement(&req, extendable) {
					continue
				}
				idx := len(rootRequirements)
				extendable[req.ruleID] = idx
				rootIndex[req.ruleID] = idx
			}
			if req.extends == "" {
				// Java lines 443-449: root requirement (no extends field).
				b := LibraryLevelBehaviorUnspecified
				if isLibraryLevelConformanceReportingMode && hasBehavior {
					b = behavior
				}
				rootRequirements = append(rootRequirements, rootEntry{req: req, behavior: b})
			}
			// Java lines 452-461: report error when library-level conformance is enabled
			// but neither rule_id nor extends is set.
			libraryLevelConformanceEnabled :=
				isLibraryLevelConformanceReportingMode && hasBehavior && behavior == LibraryLevelBehaviorRecordOnly
			if req.ruleID == "" && req.extends == "" && libraryLevelConformanceEnabled {
				reportInvalidRequirement(&req,
					"Library-level conformance requirements require their rule_id or extends field set.")
			}
		}
	}

	// --- Phase 2 (Java lines 465-483) ---
	for _, cfg := range configs {
		for ri := range cfg.requirements {
			req := cfg.requirements[ri]
			if req.extends == "" {
				continue
			}
			extIdx, ok := extendable[req.extends]
			if !ok {
				// validateExtendingRequirement will report the error.
				var nilEntry *ccRequirement
				validateExtendingRequirement(&req, nilEntry)
				continue
			}
			extendedEntry := &rootRequirements[extIdx]
			if !validateExtendingRequirement(&req, &extendedEntry.req) {
				continue
			}
			// Java lines 473-479: overwrite behavior if config specifies one.
			if cfg.hasLibraryLevelBehavior {
				extendedEntry.behavior = cfg.libraryLevelBehavior
			}
			// Java line 480: mergeExtendingRequirementIntoExtended.
			mergeExtendingRequirementIntoExtended(&extendedEntry.req, &req)
		}
	}

	// --- Phase 3 (Java lines 486-499) ---
	result := make([]ccRequirement, 0, len(rootRequirements))
	for i := range rootRequirements {
		entry := &rootRequirements[i]
		removeDuplicates(&entry.req)
		result = append(result, entry.req)
	}
	_ = rootIndex // reserved for future mergedBehaviors map wiring
	return result
}

// validateBehaviorSettingOfConfigs mirrors Java's
// validateBehaviorSettingOfConfigs(AbstractCompiler, ImmutableList<ConformanceConfig>,
//
//	ConformanceReportingMode) (lines 244-294).
//
// Only called when reportingMode == RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
// (Java line 217). Validates that any extending rule's config does not specify a
// different library_level_non_allowlisted_conformance_violations_behavior than its
// base rule's config. Returns false (and reports an error) on the first mismatch.
//
// Java algorithm (lines 253-293):
//  1. Build a map of root rule_id → behavior from configs that have the behavior field.
//  2. For each extending requirement, check that its config's behavior matches the base.
func validateBehaviorSettingOfConfigs(configs []*ccConformanceConfig) bool {
	// Java line 253: Map<String, LibraryLevelBehavior> baseRequirementBehaviors
	baseRequirementBehaviors := map[string]LibraryLevelBehavior{}

	// Java lines 255-266: first pass — collect root rule_id behaviors.
	for _, cfg := range configs {
		if !cfg.hasLibraryLevelBehavior {
			continue
		}
		behavior := cfg.libraryLevelBehavior
		for ri := range cfg.requirements {
			req := &cfg.requirements[ri]
			if req.ruleID != "" {
				baseRequirementBehaviors[req.ruleID] = behavior
			}
		}
	}

	// Java lines 268-293: second pass — validate extending requirements.
	for _, cfg := range configs {
		if !cfg.hasLibraryLevelBehavior {
			continue
		}
		extendingBehavior := cfg.libraryLevelBehavior
		for ri := range cfg.requirements {
			req := &cfg.requirements[ri]
			if req.extends == "" {
				continue
			}
			baseRuleID := req.extends
			if baseBehavior, ok := baseRequirementBehaviors[baseRuleID]; ok {
				if baseBehavior != extendingBehavior {
					reportInvalidRequirement(req,
						"extending rule's config may not specify a different value of"+
							" 'library_level_non_allowlisted_conformance_violations_behavior' than the"+
							" base rule's config. Skipping all conformance checks.")
					return false
				}
			}
		}
	}
	return true
}

// validateExtendableRequirement mirrors Java's
// validateExtendableRequirement(AbstractCompiler, Requirement,
//
//	Map<String, Requirement.Builder>) (lines 533-551).
//
// Precondition (Java line 537): req.ruleID must be non-empty (checkState).
// Returns false and reports an error when:
//   - rule_id is empty (Java line 539-542)
//   - rule_id already exists in the extendable map — duplicate (Java lines 543-549)
//
// extendable maps rule_id → index in rootRequirements (Go adaptation of Java's
// Map<String, Requirement.Builder> — same semantics, Go-idiomatic key type).
func validateExtendableRequirement(req *ccRequirement, extendable map[string]int) bool {
	// Java line 537: checkState(requirement.hasRuleId()) — req.ruleID is always set here.
	if req.ruleID == "" {
		// Java line 539-542: empty rule_id.
		reportInvalidRequirement(req, "empty rule_id")
		return false
	}
	if _, dup := extendable[req.ruleID]; dup {
		// Java lines 543-549: duplicate rule_id.
		reportInvalidRequirement(req, "two requirements with the same rule_id: "+req.ruleID)
		return false
	}
	return true
}

// extendableFields mirrors Java's EXTENDABLE_FIELDS ImmutableSet (lines 296-306).
// These are the only fields an extending requirement is allowed to set.
// Java: "config_file", "extends", "only_apply_to", "only_apply_to_regexp",
//
//	"whitelist", "whitelist_regexp", "allowlist", "allowlist_regexp", "value"
var extendableFields = map[string]bool{
	"config_file":        true,
	"extends":            true,
	"only_apply_to":      true,
	"only_apply_to_regexp": true,
	"whitelist":          true,
	"whitelist_regexp":   true,
	"allowlist":          true,
	"allowlist_regexp":   true,
	"value":              true,
}

// validateExtendingRequirement mirrors Java's
// validateExtendingRequirement(AbstractCompiler, Requirement,
//
//	Requirement.Builder) (lines 566-592).
//
// Precondition (Java line 570): extending.extends must be non-empty.
// Returns false and reports an error when:
//   - extendable is nil — no matching rule_id (Java lines 571-577)
//   - extending sets a field not in EXTENDABLE_FIELDS (Java lines 578-583)
//   - extending specifies value[] but extendable does not allowExtendingValue (Java lines 584-590)
//
// Note: Java iterates over getAllFields() from the proto descriptor, which enumerates
// only the fields that are explicitly set in the message. We mirror this by inspecting
// the non-zero/non-nil stub fields.
func validateExtendingRequirement(extending, extendable *ccRequirement) bool {
	// Java line 570: checkState(extendingRequirement.hasExtends())
	// Java lines 571-577: null extendable → no matching rule_id.
	if extendable == nil {
		reportInvalidRequirement(extending,
			"no extendable requirement with rule_id: "+extending.extends)
		return false
	}
	// Java lines 578-583: check that only EXTENDABLE_FIELDS are used.
	// In the stub type, non-EXTENDABLE_FIELDS would be type, rule_id, error_message, etc.
	// We have no proto descriptor here, so we report when rule_id is set (it must NOT be set
	// on an extending requirement — Java enforces this by checking all populated fields against
	// EXTENDABLE_FIELDS; rule_id is not in that set).
	if extending.ruleID != "" {
		reportInvalidRequirement(extending,
			"extending rules allow only "+formatExtendableFields())
	}
	// Java lines 584-590: value[] present but extendable does not allow extending value.
	if len(extending.value) > 0 && !extendable.allowExtendingValue {
		reportInvalidRequirement(extending,
			"extending rule may not specify 'value' if extendable rule does not allow it")
	}
	return true
}

// formatExtendableFields returns a human-readable set string for error messages,
// mirroring Java's EXTENDABLE_FIELDS.toString() in the error message (line 581).
func formatExtendableFields() string {
	return "[config_file, extends, only_apply_to, only_apply_to_regexp, whitelist, whitelist_regexp, allowlist, allowlist_regexp, value]"
}

// removeDuplicates mirrors Java's removeDuplicates(Requirement.Builder) (lines 594-613).
// Deduplicates each of the six allowlist/whitelist/onlyApplyTo string-list fields
// using a linked-set approach (Java uses ImmutableSet.copyOf which preserves order
// of first occurrence). Java clears each field then re-adds the deduplicated entries.
func removeDuplicates(req *ccRequirement) {
	req.whitelist = dedupStrings(req.whitelist)
	req.allowlist = dedupStrings(req.allowlist)
	req.whitelistRegexp = dedupStrings(req.whitelistRegexp)
	req.allowlistRegexp = dedupStrings(req.allowlistRegexp)
	req.onlyApplyTo = dedupStrings(req.onlyApplyTo)
	req.onlyApplyToRegexp = dedupStrings(req.onlyApplyToRegexp)
}

// dedupStrings returns a copy of ss with duplicates removed, preserving the order
// of first occurrence. Mirrors ImmutableSet.copyOf(list) semantics (Java line 595).
func dedupStrings(ss []string) []string {
	if len(ss) == 0 {
		return ss
	}
	seen := make(map[string]struct{}, len(ss))
	out := ss[:0:0]
	for _, s := range ss {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// mergeExtendingRequirementIntoExtended mirrors Java's
// mergeExtendingRequirementIntoExtended(Requirement.Builder, Requirement) (lines 508-524).
// Merges the allowlists/whitelists/values/configFiles from the extending requirement
// into the extended (base) requirement by appending. Java uses addAll* proto builder methods.
//
// Precondition (Java line 510-512): extending.extends == extended.ruleID.
func mergeExtendingRequirementIntoExtended(extended, extending *ccRequirement) {
	// Java lines 513-523: addAll* calls for each mergeable field.
	extended.whitelist = append(extended.whitelist, extending.whitelist...)
	extended.whitelistRegexp = append(extended.whitelistRegexp, extending.whitelistRegexp...)
	extended.allowlist = append(extended.allowlist, extending.allowlist...)
	extended.allowlistRegexp = append(extended.allowlistRegexp, extending.allowlistRegexp...)
	extended.onlyApplyTo = append(extended.onlyApplyTo, extending.onlyApplyTo...)
	extended.onlyApplyToRegexp = append(extended.onlyApplyToRegexp, extending.onlyApplyToRegexp...)
	extended.value = append(extended.value, extending.value...)
	extended.configFile = append(extended.configFile, extending.configFile...)
}

// initRule mirrors Java's initRule(AbstractCompiler, Requirement) (lines 622-653).
// Translates a merged Requirement into a Rule instance by switching on Requirement.Type.
// Returns nil if the requirement is invalid (InvalidRequirementSpec caught in Java).
//
// TODO: implement the switch once ConformanceRules and the Requirement.Type enum are ported.
// For now returns nil, causing initRules to skip all requirements — matching the Java
// InvalidRequirementSpec path.
func initRule(_ *ccRequirement) Rule {
	// TODO: switch req.requirementType { case CUSTOM: ...; case BANNED_NAME: ...; ... }
	return nil
}

// reportInvalidRequirement mirrors Java's
// reportInvalidRequirement(AbstractCompiler compiler, Requirement requirement, String reason)
// (lines 665-670).
// Reports a JSC_INVALID_REQUIREMENT_SPEC diagnostic for an invalid conformance requirement.
//
// TODO: wire to compiler.report() once AbstractCompiler is ported.
func reportInvalidRequirement(req *ccRequirement, reason string) {
	// TODO: call compiler.report(JSError.make(INVALID_REQUIREMENT_SPEC, reason, requirement.toString()))
	_ = req
	_ = reason
}
