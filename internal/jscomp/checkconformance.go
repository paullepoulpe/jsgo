// Port of com/google/javascript/jscomp/CheckConformance.java
// @ closure-compiler v20260415 (shim; no conformance config supported).
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
//   • mergeRequirements (lines 418-500) + validateBehavior (lines 244-294)
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

// CheckConformance is a shim for Java's CheckConformance compiler pass.
// When no conformance config path is provided (the only supported mode
// today), Check always returns nil — matching Java's early-exit when
// this.categories.isEmpty() (line 169).
//
// Wire-in: call NewCheckConformance(path) in compiler.go and append
// its Check(script) result to the diagnostics slice. Because no config
// is ever loaded, the append is a no-op and introduces zero overhead.
type CheckConformance struct {
	// path is the source file path forwarded to any future diagnostics.
	path string
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
//   CheckConformance(AbstractCompiler)                        – line 151
//   CheckConformance(AbstractCompiler, configs, reportingMode) – line 158
//
// Both are collapsed here: without a ConformanceConfig proto type,
// the config argument is a plain file path (empty = no config).
//
// TODO: port the ImmutableList<ConformanceConfig> + ConformanceReportingMode
// constructor signature (lines 158-165) once proto types are available.
func NewCheckConformance(path string) *CheckConformance {
	return ccNew(path, "")
}

// ccNew is the internal constructor that also records an optional
// config path. Separated from NewCheckConformance so callers that
// eventually supply a config path can call ccNew directly without
// adding a package-level helper.
func ccNew(path, configPath string) *CheckConformance {
	return &CheckConformance{path: path, configPath: configPath}
}

// Check runs the conformance pass over script.
//
// Java's process() method (lines 167-172) delegates to
// NodeTraversal.traverseRoots only when categories is non-empty.
// Because this shim never builds any categories (no config), it
// mirrors that early-exit and always returns nil.
//
// When a config is eventually wired in, this method should:
//   1. call ccInitRules to build categories from the loaded config
//      (TODO: port initRules, lines 206-242)
//   2. traverse the script calling ccVisit for each node
//      (TODO: port shouldTraverse / visit, lines 175-204)
//
// TODO: port NodeTraversal.traverseRoots equivalent once categories
// are populated (Java lines 169-172).
func (c *CheckConformance) Check(script *rhino.Node) []diag.Diagnostic {
	// No conformance config → no categories → no diagnostics.
	// Mirrors: if (!this.categories.isEmpty()) { NodeTraversal.traverseRoots(…) }
	// (Java line 169).
	_ = script
	return nil
}
