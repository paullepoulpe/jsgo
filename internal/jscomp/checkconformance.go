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
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// Stub types for proto-dependent infrastructure not yet ported.
// ---------------------------------------------------------------------------

// ccRequirement is a minimal stub for the Requirement proto type used by
// mergeRequirements and initRule. The Java Requirement is a protobuf message
// with fields: type, rule_id, extends, error_message, allowlist, etc.
// TODO: replace with the generated proto type once ConformanceConfig.proto is ported.
type ccRequirement struct {
	ruleID  string
	extends string
}

// ccConformanceConfig is a minimal stub for the ConformanceConfig proto.
// Java: com/google/javascript/jscomp/ConformanceConfig.proto
// TODO: replace with the generated proto type.
type ccConformanceConfig struct {
	requirements []*ccRequirement
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
// Go-only: constructor shim; Java uses the class constructor directly.
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
// Builds the category list from the merged requirements of all configs.
// Returns an empty list when configs is nil/empty, matching Java's behaviour
// when the ImmutableList is empty.
//
// TODO: port mergeRequirements, validateBehaviorSettingOfConfigs, initRule once
// ConformanceConfig proto types are available.
func initRules(configs []*ccConformanceConfig) []*Category {
	if len(configs) == 0 {
		return nil
	}
	// TODO: port mergeRequirements (lines 418-500) and initRule (lines 622-653).
	return nil
}

// mergeRequirements mirrors Java's mergeRequirements(AbstractCompiler,
//
//	List<ConformanceConfig>, boolean) (lines 418-500).
//
// Collects root requirements from all configs, merges extending requirements
// into their extended counterparts (allowlist/whitelist merging), and deduplicates.
// Returns the final flat list of Requirement objects to process.
//
// TODO: implement once Requirement proto type is available.
func mergeRequirements(configs []*ccConformanceConfig) []*ccRequirement {
	var result []*ccRequirement
	for _, cfg := range configs {
		for _, req := range cfg.requirements {
			if req.extends == "" {
				result = append(result, req)
			}
		}
	}
	// TODO: merge extending requirements into extended roots (lines 465-483).
	// TODO: remove duplicates from allowlists (lines 486-498).
	return result
}

// validateBehaviorSettingOfConfigs mirrors Java's
// validateBehaviorSettingOfConfigs(AbstractCompiler, ImmutableList<ConformanceConfig>,
//
//	ConformanceReportingMode) (lines 244-294).
//
// Validates that extending rules do not specify a different
// library_level_non_allowlisted_conformance_violations_behavior than their
// base rule. Returns false and reports an error if a mismatch is found.
//
// TODO: implement once ConformanceConfig proto is available.
func validateBehaviorSettingOfConfigs(configs []*ccConformanceConfig) bool {
	// TODO: implement validation logic (Java lines 255-293).
	_ = configs
	return true
}

// validateExtendableRequirement mirrors Java's
// validateExtendableRequirement(AbstractCompiler, Requirement,
//
//	Map<String, Requirement.Builder>) (lines 533-551).
//
// Validates that an extendable requirement has a unique, non-empty rule_id.
// Returns false when invalid.
//
// TODO: implement once Requirement proto type is available.
func validateExtendableRequirement(req *ccRequirement, extendable map[string]*ccRequirement) bool {
	if req.ruleID == "" {
		return false
	}
	if _, dup := extendable[req.ruleID]; dup {
		return false
	}
	return true
}

// validateExtendingRequirement mirrors Java's
// validateExtendingRequirement(AbstractCompiler, Requirement,
//
//	Requirement.Builder) (lines 566-592).
//
// Validates that an extending requirement correctly targets an existing
// extendable requirement and only uses EXTENDABLE_FIELDS.
// Returns false when invalid.
//
// TODO: implement field-set validation once Requirement proto is available.
func validateExtendingRequirement(extending, extendable *ccRequirement) bool {
	if extendable == nil {
		return false
	}
	return true
}

// removeDuplicates mirrors Java's removeDuplicates(Requirement.Builder) (lines 594-613).
// Deduplicates allowlist/whitelist/onlyApplyTo entries in the merged requirement.
//
// TODO: implement once Requirement.Builder proto is available.
func removeDuplicates(req *ccRequirement) {
	// TODO: deduplicate whitelist, allowlist, whitelist_regexp, allowlist_regexp,
	// only_apply_to, only_apply_to_regexp fields (Java lines 595-612).
	_ = req
}

// mergeExtendingRequirementIntoExtended mirrors Java's
// mergeExtendingRequirementIntoExtended(Requirement.Builder, Requirement) (lines 508-524).
// Merges the allowlists/whitelists from the extending requirement into the extended one.
//
// TODO: implement once Requirement.Builder proto is available.
func mergeExtendingRequirementIntoExtended(extended, extending *ccRequirement) {
	// TODO: merge whitelist, allowlist, whitelist_regexp, allowlist_regexp,
	// only_apply_to, only_apply_to_regexp, value, config_file fields.
	_ = extended
	_ = extending
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
