// Port of com/google/javascript/jscomp/AbstractCompiler.java
// @ closure-compiler v20260415

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// LifeCycleStage represents what point in optimizations the compiler is at.
// Mirrors Java's AbstractCompiler.LifeCycleStage enum.
type LifeCycleStage int

const (
	// LifeCycleRaw is the initial state: no normalization, no color.
	LifeCycleRaw LifeCycleStage = iota
	// LifeCycleColorsAndSimplifiedJSDoc after ConvertTypesToColors.
	LifeCycleColorsAndSimplifiedJSDoc
	// LifeCycleNormalized after Normalize.java.
	LifeCycleNormalized
	// LifeCycleNormalizedObfuscated after variable renaming.
	LifeCycleNormalizedObfuscated
)

func (s LifeCycleStage) IsNormalized() bool {
	return s == LifeCycleNormalized || s == LifeCycleNormalizedObfuscated
}

func (s LifeCycleStage) IsNormalizedUnobfuscated() bool {
	return s == LifeCycleNormalized
}

func (s LifeCycleStage) IsNormalizedObfuscated() bool {
	return s == LifeCycleNormalizedObfuscated
}

func (s LifeCycleStage) HasColorAndSimplifiedJSDoc() bool {
	return s != LifeCycleRaw
}

// ConfigContext distinguishes parse configuration contexts.
// Mirrors Java's AbstractCompiler.ConfigContext enum.
type ConfigContext int

const (
	ConfigContextDefault ConfigContext = iota
	ConfigContextExterns
)

// ReadError is a compile-time error emitted when a source file cannot be read.
// Mirrors Java's AbstractCompiler.READ_ERROR DiagnosticType.
var ReadError = diag.DiagnosticTypeError("JSC_READ_ERROR", "Cannot read file {0}: {1}")

// AbstractCompiler is the base type for the compiler. Passes receive a
// *AbstractCompiler to access compiler state without depending on the full
// Compiler implementation. Mirrors Java's AbstractCompiler abstract class.
//
// This Go version is a struct rather than an interface so that passes can
// embed it and override only the methods they need. Methods that haven't been
// wired to a real implementation yet return zero values or no-op to keep the
// pass pipeline exercisable.
//
// Local state fields added here are shared across all passes that hold a
// *AbstractCompiler. Note that Java's AbstractCompiler is deliberately
// state-free (all state is in Compiler); here we carry a minimum set so that
// the Go pass pipeline can run end-to-end in tests.
type AbstractCompiler struct {
	options     *CompilerOptions
	lifeCycle   LifeCycleStage
	diagnostics []diag.Diagnostic
	passName    string

	// typeCheckingHasRun tracks whether the typechecking passes have run.
	typeCheckingHasRun bool

	// hasRegExpGlobalReferences mirrors
	// AbstractCompiler.hasRegExpGlobalReferences.
	hasRegExpGlobalReferences bool

	// exportedNames tracks names added via addExportedNames.
	exportedNames map[string]bool

	// defineNames tracks names added via setDefineNames.
	defineNames map[string]bool
}

// NewAbstractCompiler returns an AbstractCompiler with default options.
func NewAbstractCompiler(opts *CompilerOptions) *AbstractCompiler {
	if opts == nil {
		opts = NewCompilerOptions()
	}
	return &AbstractCompiler{
		options:       opts,
		exportedNames: make(map[string]bool),
		defineNames:   make(map[string]bool),
	}
}

// ---------------------------------------------------------------------------
// Options / lifecycle
// ---------------------------------------------------------------------------

// GetOptions returns the compiler options. Mirrors Java's AbstractCompiler.getOptions.
func (c *AbstractCompiler) GetOptions() *CompilerOptions {
	return c.options
}

// GetLifeCycleStage returns the current AST life-cycle stage.
// Mirrors Java's AbstractCompiler.getLifeCycleStage.
func (c *AbstractCompiler) GetLifeCycleStage() LifeCycleStage {
	return c.lifeCycle
}

// SetLifeCycleStage updates the current AST life-cycle stage.
// Mirrors Java's AbstractCompiler.setLifeCycleStage.
func (c *AbstractCompiler) SetLifeCycleStage(s LifeCycleStage) {
	c.lifeCycle = s
}

// ---------------------------------------------------------------------------
// Pass lifecycle hooks
// ---------------------------------------------------------------------------

// BeforePass is called before each pass runs.
// Mirrors Java's AbstractCompiler.beforePass.
func (c *AbstractCompiler) BeforePass(passName string) {
	c.passName = passName
}

// AfterPass is called after each pass finishes.
// Mirrors Java's AbstractCompiler.afterPass.
func (c *AbstractCompiler) AfterPass(passName string) {}

// GetLastPassName returns the name of the most-recently started pass.
// Mirrors Java's AbstractCompiler.getLastPassName.
func (c *AbstractCompiler) GetLastPassName() string {
	return c.passName
}

// ---------------------------------------------------------------------------
// Diagnostics / reporting
// ---------------------------------------------------------------------------

// Report records a diagnostic. Mirrors Java's AbstractCompiler.report.
func (c *AbstractCompiler) Report(d diag.Diagnostic) {
	c.diagnostics = append(c.diagnostics, d)
}

// Diagnostics returns accumulated diagnostics.
func (c *AbstractCompiler) Diagnostics() []diag.Diagnostic {
	return c.diagnostics
}

// HasHaltingErrors reports whether any error-level diagnostics have been
// accumulated. Mirrors Java's AbstractCompiler.hasHaltingErrors.
func (c *AbstractCompiler) HasHaltingErrors() bool {
	for _, d := range c.diagnostics {
		if d.Level == diag.Error {
			return true
		}
	}
	return false
}

// ReportChangeToEnclosingScope is a no-op shim. Mirrors Java's
// AbstractCompiler.reportChangeToEnclosingScope.
func (c *AbstractCompiler) ReportChangeToEnclosingScope(n *rhino.Node) {}

// ReportChangeToChangeScope is a no-op shim. Mirrors Java's
// AbstractCompiler.reportChangeToChangeScope.
func (c *AbstractCompiler) ReportChangeToChangeScope(n *rhino.Node) {}

// ReportFunctionDeleted is a no-op shim. Mirrors Java's
// AbstractCompiler.reportFunctionDeleted.
func (c *AbstractCompiler) ReportFunctionDeleted(n *rhino.Node) {}

// ThrowInternalError panics with the given message. Mirrors Java's
// AbstractCompiler.throwInternalError.
func (c *AbstractCompiler) ThrowInternalError(msg string, cause error) {
	if cause != nil {
		panic(fmt.Sprintf("Internal compiler error: %s: %v", msg, cause))
	}
	panic(fmt.Sprintf("Internal compiler error: %s", msg))
}

// ---------------------------------------------------------------------------
// Type-checking state
// ---------------------------------------------------------------------------

// HasTypeCheckingRun reports whether type checking has run.
// Mirrors Java's AbstractCompiler.hasTypeCheckingRun.
func (c *AbstractCompiler) HasTypeCheckingRun() bool {
	return c.typeCheckingHasRun || c.lifeCycle >= LifeCycleColorsAndSimplifiedJSDoc
}

// SetTypeCheckingHasRun sets whether type checking has completed.
// Mirrors Java's AbstractCompiler.setTypeCheckingHasRun.
func (c *AbstractCompiler) SetTypeCheckingHasRun(hasRun bool) {
	c.typeCheckingHasRun = hasRun
	if hasRun && c.lifeCycle < LifeCycleColorsAndSimplifiedJSDoc {
		c.lifeCycle = LifeCycleColorsAndSimplifiedJSDoc
	}
}

// HasOptimizationColors reports whether the AST has been annotated with
// optimization colors. Mirrors Java's AbstractCompiler.hasOptimizationColors.
func (c *AbstractCompiler) HasOptimizationColors() bool {
	return c.lifeCycle >= LifeCycleColorsAndSimplifiedJSDoc
}

// IsTypeRegistryCleared returns true when type checking has run but the
// type registry has been cleared. Mirrors Java's
// AbstractCompiler.isTypeRegistryCleared.
func (c *AbstractCompiler) IsTypeRegistryCleared() bool {
	return false // stub: type registry clearing not yet implemented
}

// ---------------------------------------------------------------------------
// Source access
// ---------------------------------------------------------------------------

// GetScriptNode returns the SCRIPT AST node for the named file.
// Mirrors Java's AbstractCompiler.getScriptNode.
// This stub always returns nil; the real implementation is in Compiler.
func (c *AbstractCompiler) GetScriptNode(filename string) *rhino.Node {
	return nil
}

// GetSourceFileByName looks up a source file by name.
// Mirrors Java's AbstractCompiler.getSourceFileByName.
func (c *AbstractCompiler) GetSourceFileByName(sourceName string) interface{} {
	return nil
}

// PreferRegexParser reports whether the regex-based dependency parser is
// preferred over a full AST-based parse. Mirrors Java's
// AbstractCompiler.preferRegexParser.
func (c *AbstractCompiler) PreferRegexParser() bool {
	return true
}

// ParseSyntheticCode parses code for injecting. Mirrors Java's
// AbstractCompiler.parseSyntheticCode.
func (c *AbstractCompiler) ParseSyntheticCode(filename, code string) *rhino.Node {
	return nil
}

// ToSource prints the AST as source. Mirrors Java's AbstractCompiler.toSource.
func (c *AbstractCompiler) ToSource() string {
	return ""
}

// ToSourceNode prints a single node to source. Mirrors Java's
// AbstractCompiler.toSource(Node).
func (c *AbstractCompiler) ToSourceNode(root *rhino.Node) string {
	return ""
}

// ---------------------------------------------------------------------------
// Module / chunk information
// ---------------------------------------------------------------------------

// GetNumberOfInputs returns the total number of compiler inputs.
// Mirrors Java's AbstractCompiler.getNumberOfInputs.
func (c *AbstractCompiler) GetNumberOfInputs() int {
	return 0
}

// GetModuleTypeByName looks up the module type for the named module.
// Mirrors Java's AbstractCompiler.getModuleTypeByName.
func (c *AbstractCompiler) GetModuleTypeByName(moduleName string) ModuleType {
	return ModuleTypeNone
}

// ---------------------------------------------------------------------------
// Exported / define names tracking
// ---------------------------------------------------------------------------

// AddExportedNames records exported variable names. Mirrors Java's
// AbstractCompiler.addExportedNames.
func (c *AbstractCompiler) AddExportedNames(names map[string]bool) {
	for n := range names {
		c.exportedNames[n] = true
	}
}

// GetExportedNames returns the set of names that have been exported.
// Mirrors Java's AbstractCompiler.getExportedNames.
func (c *AbstractCompiler) GetExportedNames() map[string]bool {
	return c.exportedNames
}

// SetDefineNames stores the @define names. Mirrors Java's
// AbstractCompiler.setDefineNames.
func (c *AbstractCompiler) SetDefineNames(names []string) {
	c.defineNames = make(map[string]bool, len(names))
	for _, n := range names {
		c.defineNames[n] = true
	}
}

// GetDefineNames returns the registered @define names. Mirrors Java's
// AbstractCompiler.getDefineNames.
func (c *AbstractCompiler) GetDefineNames() map[string]bool {
	return c.defineNames
}

// ---------------------------------------------------------------------------
// RegExp global references
// ---------------------------------------------------------------------------

// SetHasRegExpGlobalReferences records whether the AST references
// RegExp global properties. Mirrors Java's
// AbstractCompiler.setHasRegExpGlobalReferences.
func (c *AbstractCompiler) SetHasRegExpGlobalReferences(references bool) {
	c.hasRegExpGlobalReferences = references
}

// HasRegExpGlobalReferences reports whether the AST contains references to
// RegExp global object properties. Mirrors Java's
// AbstractCompiler.hasRegExpGlobalReferences.
func (c *AbstractCompiler) HasRegExpGlobalReferences() bool {
	return c.hasRegExpGlobalReferences
}

// ---------------------------------------------------------------------------
// AreNodesEqualForInlining
// ---------------------------------------------------------------------------

// AreNodesEqualForInlining reports whether two nodes are equal for the
// purpose of inlining. Mirrors Java's AbstractCompiler.areNodesEqualForInlining.
// This stub compares by pointer identity.
func (c *AbstractCompiler) AreNodesEqualForInlining(n1, n2 *rhino.Node) bool {
	return n1 == n2
}

// ---------------------------------------------------------------------------
// Process (run a pass)
// ---------------------------------------------------------------------------

// Process runs a compiler pass. Mirrors Java's AbstractCompiler.process.
func (c *AbstractCompiler) Process(pass CompilerPass) {
	pass.Process(nil, nil)
}

// ---------------------------------------------------------------------------
// Fill-file helpers
// ---------------------------------------------------------------------------

// createFillFileName mirrors Java's AbstractCompiler.createFillFileName.
func createFillFileName(chunkName string) string {
	return chunkName + "$fillFile"
}

// IsFillFileName mirrors Java's AbstractCompiler.isFillFileName.
func IsFillFileName(fileName string) bool {
	const suffix = "$fillFile"
	return len(fileName) >= len(suffix) && fileName[len(fileName)-len(suffix):] == suffix
}

// ---------------------------------------------------------------------------
// Debug logging
// ---------------------------------------------------------------------------

// IsDebugLoggingEnabled reports whether debug logging is enabled.
// Mirrors Java's AbstractCompiler.isDebugLoggingEnabled.
func (c *AbstractCompiler) IsDebugLoggingEnabled() bool {
	return false
}

// GetDebugLogFilterList returns the debug-log filter list.
// Mirrors Java's AbstractCompiler.getDebugLogFilterList.
func (c *AbstractCompiler) GetDebugLogFilterList() []string {
	return nil
}

// ---------------------------------------------------------------------------
// Synthetic code input lifecycle
// ---------------------------------------------------------------------------

// InitializeSyntheticCodeInput adds a synthetic script to the front of the AST.
// Mirrors Java's AbstractCompiler.initializeSyntheticCodeInput.
func (c *AbstractCompiler) InitializeSyntheticCodeInput() {}

// RemoveSyntheticCodeInput removes the script added by InitializeSyntheticCodeInput.
// Mirrors Java's AbstractCompiler.removeSyntheticCodeInput.
func (c *AbstractCompiler) RemoveSyntheticCodeInput() {}

// MergeSyntheticCodeInput merges the synthetic code input into the first non-synthetic
// script. Mirrors Java's AbstractCompiler.mergeSyntheticCodeInput.
func (c *AbstractCompiler) MergeSyntheticCodeInput() {}

// ---------------------------------------------------------------------------
// J2CL pass control
// ---------------------------------------------------------------------------

// SetRunJ2clPasses sets whether J2CL passes should run.
// Mirrors Java's AbstractCompiler.setRunJ2clPasses.
func (c *AbstractCompiler) SetRunJ2clPasses(run bool) {}

// RunJ2clPasses reports whether J2CL passes should run.
// Mirrors Java's AbstractCompiler.runJ2clPasses.
func (c *AbstractCompiler) RunJ2clPasses() bool { return false }

// ---------------------------------------------------------------------------
// Externs properties
// ---------------------------------------------------------------------------

// SetExternProperties stores the set of property names defined in externs.
// Mirrors Java's AbstractCompiler.setExternProperties.
func (c *AbstractCompiler) SetExternProperties(props map[string]bool) {}

// GetExternProperties returns the set of property names defined in externs,
// or nil if GatherExternProperties has not run. Mirrors Java's
// AbstractCompiler.getExternProperties.
func (c *AbstractCompiler) GetExternProperties() map[string]bool { return nil }
