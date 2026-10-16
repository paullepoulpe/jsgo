// Port of com/google/javascript/jscomp/AbstractCompiler.java
// @ closure-compiler v20260415

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/legacy/internal/diag"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
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

	// runJ2clPasses mirrors Compiler.runJ2clPasses.
	runJ2clPasses bool

	// typeRegistryCleared mirrors Compiler.typeRegistryCleared.
	typeRegistryCleared bool

	// exportedNames tracks names added via addExportedNames.
	exportedNames map[string]bool

	// defineNames tracks names added via setDefineNames.
	defineNames map[string]bool

	// externProperties mirrors Compiler.externProperties.
	externProperties map[string]bool

	// accessorSummary mirrors Compiler.accessorSummary (opaque until AccessorSummary is ported).
	accessorSummary interface{}

	// moduleMetadataMap mirrors Compiler.moduleMetadataMap (opaque until ModuleMetadataMap is ported).
	moduleMetadataMap interface{}

	// moduleMap mirrors Compiler.moduleMap (opaque until ModuleMap is ported).
	moduleMap interface{}

	// crossChunkIdGenerator mirrors Compiler.crossChunkIdGenerator (opaque until IdGenerator is ported).
	crossChunkIdGenerator interface{}

	// inputSourceMaps mirrors Compiler.inputSourceMaps.
	inputSourceMaps map[string]interface{}

	// commentsPerFile mirrors Compiler.commentsPerFile.
	commentsPerFile map[string]interface{}

	// variableMap mirrors Compiler.variableMap.
	variableMap interface{}

	// propertyMap mirrors Compiler.propertyMap.
	propertyMap interface{}

	// stringMap mirrors Compiler.stringMap.
	stringMap interface{}

	// cssNames mirrors Compiler.cssNames.
	cssNames map[string]bool

	// instrumentationMapping mirrors Compiler.instrumentationMapping.
	instrumentationMapping interface{}

	// idGeneratorMap mirrors Compiler.idGeneratorMap.
	idGeneratorMap string

	// anonymousFunctionNameMap mirrors Compiler.anonymousFunctionNameMap.
	anonymousFunctionNameMap interface{}

	// transpiledFiles mirrors Compiler.transpiledFiles.
	transpiledFiles bool

	// disambiguatePropertiesSummary mirrors Compiler.disambiguatePropertiesSummary.
	disambiguatePropertiesSummary func() string

	// ambiguatePropertiesSummary mirrors Compiler.ambiguatePropertiesSummary.
	ambiguatePropertiesSummary func() string

	// syntheticCodeInput mirrors Compiler.syntheticCodeInput.
	syntheticCodeInput *CompilerInput

	// scriptsByFilename mirrors Compiler.scriptNodeByFilename.
	scriptsByFilename map[string]*rhino.Node

	// numberOfInputs mirrors Compiler.getNumberOfInputs.
	numberOfInputs int

	// externExports mirrors Compiler.externExports.
	externExports string

	// debugLogEnabled mirrors Compiler.isDebugLoggingEnabled.
	debugLogEnabled bool

	// debugLogFilterList mirrors Compiler.getDebugLogFilterList.
	debugLogFilterList []string

	// allowableFeatures mirrors Compiler.allowableFeatures (opaque until FeatureSet is ported).
	allowableFeatures interface{}
}

// Go-only: language-adapter
// NewAbstractCompiler returns an AbstractCompiler with default options.
func NewAbstractCompiler(opts *CompilerOptions) *AbstractCompiler {
	if opts == nil {
		opts = NewCompilerOptions()
	}
	return &AbstractCompiler{
		options:         opts,
		exportedNames:   make(map[string]bool),
		defineNames:     make(map[string]bool),
		inputSourceMaps: make(map[string]interface{}),
		commentsPerFile: make(map[string]interface{}),
		scriptsByFilename: make(map[string]*rhino.Node),
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
func (c *AbstractCompiler) AfterPass(passName string) {
	c.passName = passName
}

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

// Go-only: language-adapter
// Diagnostics returns accumulated diagnostics. Java's Compiler uses
// getErrorManager().generateReport() or iterates errors; Go uses a slice.
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

// ReportChangeToEnclosingScope notifies the change tracker that the enclosing
// scope of n was modified. Mirrors Java's AbstractCompiler.reportChangeToEnclosingScope.
// The Java Compiler implementation delegates to changeTracker; here we record
// the notification without a full ChangeTracker port.
func (c *AbstractCompiler) ReportChangeToEnclosingScope(n *rhino.Node) {
	// No-op: ChangeTracker not yet ported. Java Compiler delegates to
	// changeTracker.reportChangeToEnclosingScope(n).
	_ = n
}

// ReportChangeToChangeScope notifies the change tracker that a change-scope
// root was modified. Mirrors Java's AbstractCompiler.reportChangeToChangeScope.
func (c *AbstractCompiler) ReportChangeToChangeScope(n *rhino.Node) {
	// No-op: ChangeTracker not yet ported. Java Compiler delegates to
	// changeTracker.reportChangeToChangeScope(n).
	_ = n
}

// ReportFunctionDeleted marks a function node as deleted for change tracking.
// Mirrors Java's AbstractCompiler.reportFunctionDeleted.
func (c *AbstractCompiler) ReportFunctionDeleted(n *rhino.Node) {
	// No-op: ChangeTracker not yet ported. Java Compiler delegates to
	// changeTracker.reportFunctionDeleted(n).
	_ = n
}

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
	return c.typeRegistryCleared
}

// Go-only: language-adapter
// SetTypeRegistryCleared records that the JS type registry has been cleared.
// Mirrors Java's Compiler.clearJSTypeRegistry() which sets typeRegistryCleared.
func (c *AbstractCompiler) SetTypeRegistryCleared(cleared bool) {
	c.typeRegistryCleared = cleared
}

// ---------------------------------------------------------------------------
// Source access
// ---------------------------------------------------------------------------

// GetScriptNode returns the SCRIPT AST node for the named file.
// Mirrors Java's AbstractCompiler.getScriptNode.
// Returns the registered script node for the filename, or nil if not found.
func (c *AbstractCompiler) GetScriptNode(filename string) *rhino.Node {
	if c.scriptsByFilename == nil {
		return nil
	}
	return c.scriptsByFilename[filename]
}

// Go-only: language-adapter
// RegisterScriptNode registers a SCRIPT-rooted node for the given filename.
// Mirrors Java's Compiler.scriptNodeByFilename map population during parse.
func (c *AbstractCompiler) RegisterScriptNode(filename string, n *rhino.Node) {
	if c.scriptsByFilename == nil {
		c.scriptsByFilename = make(map[string]*rhino.Node)
	}
	c.scriptsByFilename[filename] = n
}

// GetSourceFileByName looks up a source file by name.
// Mirrors Java's AbstractCompiler.getSourceFileByName.
// Returns nil until SourceFile is ported; source names are tracked via scriptsByFilename.
// TODO(port): Compiler.getSourceFileByName(String) — requires SourceFile type port
func (c *AbstractCompiler) GetSourceFileByName(sourceName string) interface{} {
	// SourceFile not yet ported; check scriptsByFilename as a proxy.
	_ = sourceName
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
// TODO(port): Compiler.parseSyntheticCode(String filename, String code)
func (c *AbstractCompiler) ParseSyntheticCode(filename, code string) *rhino.Node {
	// Not yet implemented; requires ParserRunner → AST construction outside JsAst.
	_ = filename
	_ = code
	return nil
}

// ToSource prints the AST as source. Mirrors Java's AbstractCompiler.toSource.
// TODO(port): Compiler.toSource() — requires CodePrinter port
func (c *AbstractCompiler) ToSource() string {
	// CodePrinter not yet ported. Java delegates to toSource(Node) with jsRoot.
	_ = c.options
	return ""
}

// ToSourceNode prints a single node to source. Mirrors Java's
// AbstractCompiler.toSource(Node).
// TODO(port): Compiler.toSource(Node root) — requires CodePrinter port
func (c *AbstractCompiler) ToSourceNode(root *rhino.Node) string {
	// Not yet implemented; requires CodePrinter which is not yet ported.
	_ = root
	return ""
}

// ---------------------------------------------------------------------------
// Module / chunk information
// ---------------------------------------------------------------------------

// GetNumberOfInputs returns the total number of compiler inputs.
// Mirrors Java's AbstractCompiler.getNumberOfInputs.
func (c *AbstractCompiler) GetNumberOfInputs() int {
	return c.numberOfInputs
}

// Go-only: language-adapter
// SetNumberOfInputs stores the total input count.
// Mirrors Java's Compiler field numberOfInputs populated during init.
func (c *AbstractCompiler) SetNumberOfInputs(n int) {
	c.numberOfInputs = n
}

// GetModuleTypeByName looks up the module type for the named module.
// Mirrors Java's AbstractCompiler.getModuleTypeByName.
// TODO(port): Compiler.getModuleTypeByName(String moduleName) — requires moduleTypesByName map
func (c *AbstractCompiler) GetModuleTypeByName(moduleName string) ModuleType {
	return ModuleTypeNone
}

// ---------------------------------------------------------------------------
// Exported / define names tracking
// ---------------------------------------------------------------------------

// AddExportedNames records exported variable names. Mirrors Java's
// AbstractCompiler.addExportedNames. Note: Java uses Set<String> while Go uses
// map[string]bool for idiomatic membership testing.
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
// Extern exports
// ---------------------------------------------------------------------------

// Go-only: language-adapter
// GetExternExports returns the extern exports string.
// Mirrors Java's Compiler.externExports field. Java exposes this via
// Compiler.toSource() rather than a direct getter.
func (c *AbstractCompiler) GetExternExports() string {
	return c.externExports
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
	return c.debugLogEnabled
}

// Go-only: language-adapter
// SetDebugLoggingEnabled enables or disables debug logging.
// Mirrors Java's Compiler.isDebugLoggingEnabled (controlled by options.getDebugLogDirectory).
func (c *AbstractCompiler) SetDebugLoggingEnabled(enabled bool) {
	c.debugLogEnabled = enabled
}

// GetDebugLogFilterList returns the debug-log filter list.
// Mirrors Java's AbstractCompiler.getDebugLogFilterList.
func (c *AbstractCompiler) GetDebugLogFilterList() []string {
	return c.debugLogFilterList
}

// Go-only: language-adapter
// SetDebugLogFilterList stores the debug-log filter list.
// Mirrors Java's Compiler.getDebugLogFilterList (driven by options.getDebugLogDirectory).
func (c *AbstractCompiler) SetDebugLogFilterList(list []string) {
	c.debugLogFilterList = list
}

// ---------------------------------------------------------------------------
// Synthetic code input lifecycle
// ---------------------------------------------------------------------------

// InitializeSyntheticCodeInput adds a synthetic script to the front of the AST.
// Mirrors Java's AbstractCompiler.initializeSyntheticCodeInput.
// TODO(port): Compiler.initializeSyntheticCodeInput() — requires full AST root management
func (c *AbstractCompiler) InitializeSyntheticCodeInput() {
	if c.syntheticCodeInput == nil {
		c.syntheticCodeInput = NewCompilerInput(c.GetSyntheticCodeInputId(), "")
	}
}

// RemoveSyntheticCodeInput removes the script added by InitializeSyntheticCodeInput.
// Mirrors Java's AbstractCompiler.removeSyntheticCodeInput.
// TODO(port): Compiler.removeSyntheticCodeInput() — requires full AST root management
func (c *AbstractCompiler) RemoveSyntheticCodeInput() {
	c.syntheticCodeInput = nil
}

// MergeSyntheticCodeInput merges the synthetic code input into the first non-synthetic
// script. Mirrors Java's AbstractCompiler.mergeSyntheticCodeInput.
// TODO(port): Compiler.mergeSyntheticCodeInput() — requires full AST root management
func (c *AbstractCompiler) MergeSyntheticCodeInput() {
	c.syntheticCodeInput = nil
}

// ---------------------------------------------------------------------------
// J2CL pass control
// ---------------------------------------------------------------------------

// SetRunJ2clPasses sets whether J2CL passes should run.
// Mirrors Java's AbstractCompiler.setRunJ2clPasses.
func (c *AbstractCompiler) SetRunJ2clPasses(run bool) {
	c.runJ2clPasses = run
}

// RunJ2clPasses reports whether J2CL passes should run.
// Mirrors Java's AbstractCompiler.runJ2clPasses.
func (c *AbstractCompiler) RunJ2clPasses() bool {
	return c.runJ2clPasses
}

// ---------------------------------------------------------------------------
// Externs properties
// ---------------------------------------------------------------------------

// SetExternProperties stores the set of property names defined in externs.
// Mirrors Java's AbstractCompiler.setExternProperties.
func (c *AbstractCompiler) SetExternProperties(props map[string]bool) {
	c.externProperties = props
}

// GetExternProperties returns the set of property names defined in externs,
// or nil if GatherExternProperties has not run. Mirrors Java's
// AbstractCompiler.getExternProperties.
func (c *AbstractCompiler) GetExternProperties() map[string]bool {
	return c.externProperties
}

// ---------------------------------------------------------------------------
// Variable / property / string maps
// ---------------------------------------------------------------------------

// SetVariableMap stores the variable renaming map. Mirrors Java's
// AbstractCompiler.setVariableMap.
func (c *AbstractCompiler) SetVariableMap(m interface{}) {
	c.variableMap = m
}

// Go-only: language-adapter
// GetVariableMap returns the variable renaming map. Java accesses the VariableMap
// via Result.variableRenaming; Go exposes it directly.
func (c *AbstractCompiler) GetVariableMap() interface{} {
	return c.variableMap
}

// SetPropertyMap stores the property renaming map. Mirrors Java's
// AbstractCompiler.setPropertyMap.
func (c *AbstractCompiler) SetPropertyMap(m interface{}) {
	c.propertyMap = m
}

// Go-only: language-adapter
// GetPropertyMap returns the property renaming map. Java accesses the VariableMap
// via Result.propertyRenaming; Go exposes it directly.
func (c *AbstractCompiler) GetPropertyMap() interface{} {
	return c.propertyMap
}

// SetStringMap stores the string replacement map. Mirrors Java's
// AbstractCompiler.setStringMap.
func (c *AbstractCompiler) SetStringMap(m interface{}) {
	c.stringMap = m
}

// Go-only: language-adapter
// GetStringMap returns the string replacement map.
// Java accesses this via Result.stringMap; Go exposes it directly.
func (c *AbstractCompiler) GetStringMap() interface{} {
	return c.stringMap
}

// SetCssNames stores the CSS names found during compilation. Mirrors Java's
// AbstractCompiler.setCssNames.
func (c *AbstractCompiler) SetCssNames(names map[string]bool) {
	c.cssNames = names
}

// Go-only: language-adapter
// GetCssNames returns the CSS names found during compilation.
// Java accesses this via Result.cssNames; Go exposes it directly.
func (c *AbstractCompiler) GetCssNames() map[string]bool {
	return c.cssNames
}

// SetInstrumentationMapping stores the instrumentation parameter encoding map.
// Mirrors Java's AbstractCompiler.setInstrumentationMapping.
func (c *AbstractCompiler) SetInstrumentationMapping(m interface{}) {
	c.instrumentationMapping = m
}

// SetIdGeneratorMap stores the ID generator map for cross-chunk motion.
// Mirrors Java's AbstractCompiler.setIdGeneratorMap.
func (c *AbstractCompiler) SetIdGeneratorMap(serialized string) {
	c.idGeneratorMap = serialized
}

// Go-only: language-adapter
// GetIdGeneratorMap returns the serialized ID generator map.
// Java accesses this via Result.idRenaming; Go exposes it directly.
func (c *AbstractCompiler) GetIdGeneratorMap() string {
	return c.idGeneratorMap
}

// GetTranspiledFiles reports whether any file needed to transpile any feature.
// Mirrors Java's AbstractCompiler.getTranspiledFiles.
func (c *AbstractCompiler) GetTranspiledFiles() bool {
	return c.transpiledFiles
}

// Go-only: language-adapter
// SetTranspiledFiles stores whether any file needed transpilation.
// Java sets this field internally during transpilation passes.
func (c *AbstractCompiler) SetTranspiledFiles(v bool) {
	c.transpiledFiles = v
}

// SetAnonymousFunctionNameMap stores the naming map for anonymous functions.
// Mirrors Java's AbstractCompiler.setAnonymousFunctionNameMap.
func (c *AbstractCompiler) SetAnonymousFunctionNameMap(m interface{}) {
	c.anonymousFunctionNameMap = m
}

// Go-only: language-adapter
// GetAnonymousFunctionNameMap returns the anonymous function naming map.
// Java accesses this via Result.functionNames; Go exposes it directly.
func (c *AbstractCompiler) GetAnonymousFunctionNameMap() interface{} {
	return c.anonymousFunctionNameMap
}

// ---------------------------------------------------------------------------
// Source map support
// ---------------------------------------------------------------------------

// AddInputSourceMap adds a SourceMapInput for the given source file.
// Mirrors Java's AbstractCompiler.addInputSourceMap.
func (c *AbstractCompiler) AddInputSourceMap(name string, sourceMap interface{}) {
	if c.inputSourceMaps == nil {
		c.inputSourceMaps = make(map[string]interface{})
	}
	c.inputSourceMaps[name] = sourceMap
}

// Go-only: language-adapter
// GetInputSourceMap returns the SourceMapInput for the given source file, or nil.
// Java accesses inputSourceMaps directly in Compiler; Go needs an explicit getter.
func (c *AbstractCompiler) GetInputSourceMap(name string) interface{} {
	if c.inputSourceMaps == nil {
		return nil
	}
	return c.inputSourceMaps[name]
}

// GetBase64SourceMapContents returns the base64-encoded source map for a file.
// Mirrors Java's AbstractCompiler.getBase64SourceMapContents.
// TODO(port): Compiler.getBase64SourceMapContents(String sourceFileName) — requires SourceMap port
func (c *AbstractCompiler) GetBase64SourceMapContents(sourceFileName string) string {
	// Not yet implemented; requires SourceMap / SourceMapInput types.
	_ = sourceFileName
	return ""
}

// ---------------------------------------------------------------------------
// Comments
// ---------------------------------------------------------------------------

// AddComments stores the comments from a given file.
// Mirrors Java's AbstractCompiler.addComments.
func (c *AbstractCompiler) AddComments(filename string, comments interface{}) {
	if c.commentsPerFile == nil {
		c.commentsPerFile = make(map[string]interface{})
	}
	c.commentsPerFile[filename] = comments
}

// GetComments returns all comments from the given file.
// Mirrors Java's AbstractCompiler.getComments.
func (c *AbstractCompiler) GetComments(filename string) interface{} {
	if c.commentsPerFile == nil {
		return nil
	}
	return c.commentsPerFile[filename]
}

// ---------------------------------------------------------------------------
// Accessor summary
// ---------------------------------------------------------------------------

// SetAccessorSummary stores the summary of properties with getters and setters.
// Mirrors Java's AbstractCompiler.setAccessorSummary.
func (c *AbstractCompiler) SetAccessorSummary(summary interface{}) {
	c.accessorSummary = summary
}

// GetAccessorSummary returns the summary of properties with getters and setters.
// Mirrors Java's AbstractCompiler.getAccessorSummary.
func (c *AbstractCompiler) GetAccessorSummary() interface{} {
	return c.accessorSummary
}

// ---------------------------------------------------------------------------
// Module metadata / module map
// ---------------------------------------------------------------------------

// GetModuleMetadataMap returns the module metadata map.
// Mirrors Java's AbstractCompiler.getModuleMetadataMap.
func (c *AbstractCompiler) GetModuleMetadataMap() interface{} {
	return c.moduleMetadataMap
}

// SetModuleMetadataMap stores the module metadata map.
// Mirrors Java's AbstractCompiler.setModuleMetadataMap.
func (c *AbstractCompiler) SetModuleMetadataMap(m interface{}) {
	c.moduleMetadataMap = m
}

// GetModuleMap returns the module map.
// Mirrors Java's AbstractCompiler.getModuleMap.
func (c *AbstractCompiler) GetModuleMap() interface{} {
	return c.moduleMap
}

// SetModuleMap stores the module map.
// Mirrors Java's AbstractCompiler.setModuleMap.
func (c *AbstractCompiler) SetModuleMap(m interface{}) {
	c.moduleMap = m
}

// ---------------------------------------------------------------------------
// Cross-chunk ID generator
// ---------------------------------------------------------------------------

// GetCrossChunkIdGenerator returns the ID generator for cross-chunk motion.
// Mirrors Java's AbstractCompiler.getCrossChunkIdGenerator.
// TODO(port): Compiler.getCrossChunkIdGenerator() — requires IdGenerator port
func (c *AbstractCompiler) GetCrossChunkIdGenerator() interface{} {
	return c.crossChunkIdGenerator
}

// Go-only: language-adapter
// SetCrossChunkIdGenerator stores the cross-chunk ID generator.
// Java initialises this inside Compiler(); Go needs an explicit setter.
func (c *AbstractCompiler) SetCrossChunkIdGenerator(gen interface{}) {
	c.crossChunkIdGenerator = gen
}

// ---------------------------------------------------------------------------
// Error level
// ---------------------------------------------------------------------------

// GetErrorLevel returns the error level the given JSError will be reported at.
// Mirrors Java's AbstractCompiler.getErrorLevel(JSError).
// This stub returns the default level from the diagnostic type.
func (c *AbstractCompiler) GetErrorLevel(code string, defaultLevel diag.Level) diag.Level {
	return defaultLevel
}

// ---------------------------------------------------------------------------
// Extern exports
// ---------------------------------------------------------------------------

// SetExternExports stores the extern exports string.
// Mirrors Java's AbstractCompiler.setExternExports.
func (c *AbstractCompiler) SetExternExports(externExports string) {
	c.externExports = externExports
}

// ---------------------------------------------------------------------------
// Disambiguation / ambiguation property summary
// ---------------------------------------------------------------------------

// ReportDisambiguatePropertiesSummary is called by DisambiguateProperties to
// record a summary. Mirrors Java's AbstractCompiler.reportDisambiguatePropertiesSummary.
func (c *AbstractCompiler) ReportDisambiguatePropertiesSummary(summary func() string) {
	c.disambiguatePropertiesSummary = summary
}

// Go-only: language-adapter
// GetDisambiguatePropertiesSummary returns the recorded disambiguation summary.
// Java logs this via Tracer; Go exposes it as a string getter.
func (c *AbstractCompiler) GetDisambiguatePropertiesSummary() string {
	if c.disambiguatePropertiesSummary == nil {
		return ""
	}
	return c.disambiguatePropertiesSummary()
}

// ReportAmbiguatePropertiesSummary is called by AmbiguateProperties to record
// a summary. Mirrors Java's AbstractCompiler.reportAmbiguatePropertiesSummary.
func (c *AbstractCompiler) ReportAmbiguatePropertiesSummary(summary func() string) {
	c.ambiguatePropertiesSummary = summary
}

// Go-only: language-adapter
// GetAmbiguatePropertiesSummary returns the recorded ambiguate summary.
// Java logs this via Tracer; Go exposes it as a string getter.
func (c *AbstractCompiler) GetAmbiguatePropertiesSummary() string {
	if c.ambiguatePropertiesSummary == nil {
		return ""
	}
	return c.ambiguatePropertiesSummary()
}

// ---------------------------------------------------------------------------
// Feature set management
// ---------------------------------------------------------------------------

// GetAllowableFeatures returns the feature set currently allowed in the AST.
// Mirrors Java's AbstractCompiler.getAllowableFeatures.
// TODO(port): Compiler.getAllowableFeatures() — requires FeatureSet port
func (c *AbstractCompiler) GetAllowableFeatures() interface{} {
	return c.allowableFeatures
}

// SetAllowableFeatures sets the feature set allowed in the AST.
// Mirrors Java's AbstractCompiler.setAllowableFeatures.
// TODO(port): Compiler.setAllowableFeatures(FeatureSet featureSet) — requires FeatureSet port
func (c *AbstractCompiler) SetAllowableFeatures(featureSet interface{}) {
	c.allowableFeatures = featureSet
}

// MarkFeatureNotAllowed marks a feature as no longer allowed in the AST.
// Mirrors Java's AbstractCompiler.markFeatureNotAllowed.
// TODO(port): Compiler.markFeatureNotAllowed(Feature feature) — requires FeatureSet port
func (c *AbstractCompiler) MarkFeatureNotAllowed(feature interface{}) {
	// No-op until FeatureSet is ported. Java Compiler removes the feature from allowableFeatures.
	_ = feature
}

// MarkFeatureSetNotAllowed marks a feature set as no longer allowed in the AST.
// Mirrors Java's AbstractCompiler.markFeatureSetNotAllowed.
// TODO(port): Compiler.markFeatureSetNotAllowed(FeatureSet featureSet) — requires FeatureSet port
func (c *AbstractCompiler) MarkFeatureSetNotAllowed(featureSet interface{}) {
	// No-op until FeatureSet is ported. Java Compiler removes features from allowableFeatures.
	_ = featureSet
}

// ---------------------------------------------------------------------------
// AstFactory / ExpressionDecomposer / AstAnalyzer
// ---------------------------------------------------------------------------

// CreateAstFactory returns an AstFactory configured for the current
// compilation state. Mirrors Java's AbstractCompiler.createAstFactory.
// TODO(port): Compiler.createAstFactory() — requires AstFactory port
func (c *AbstractCompiler) CreateAstFactory() interface{} {
	// AstFactory not yet ported. Java checks hasTypeCheckingRun() to decide
	// whether to inject type information.
	withTypes := c.HasTypeCheckingRun()
	_ = withTypes
	return nil
}

// CreateAstFactoryWithoutTypes returns an AstFactory that will not add type
// information. Mirrors Java's AbstractCompiler.createAstFactoryWithoutTypes.
// TODO(port): Compiler.createAstFactoryWithoutTypes() — requires AstFactory port
func (c *AbstractCompiler) CreateAstFactoryWithoutTypes() interface{} {
	// AstFactory not yet ported. Java creates AstFactory.createFactoryWithoutTypes().
	_ = c.options
	return nil
}

// GetAstAnalyzer returns an AstAnalyzer for the current compilation.
// Mirrors Java's AbstractCompiler.getAstAnalyzer.
// TODO(port): Compiler.getAstAnalyzer() — requires AstAnalyzer port
func (c *AbstractCompiler) GetAstAnalyzer() interface{} {
	// AstAnalyzer not yet ported. Java creates new AstAnalyzer(this, options).
	_ = c.options
	return nil
}

// ---------------------------------------------------------------------------
// SyntheticCodeInputId
// ---------------------------------------------------------------------------

// GetSyntheticCodeInputId returns the InputId of the synthetic code input.
// Mirrors Java's AbstractCompiler.getSyntheticCodeInputId.
func (c *AbstractCompiler) GetSyntheticCodeInputId() string {
	return " [synthetic] "
}

// ---------------------------------------------------------------------------
// Change tracking
// ---------------------------------------------------------------------------

// HasScopeChanged reports whether a function changed since the last pass run.
// Mirrors Java's AbstractCompiler.hasScopeChanged.
// Returns true conservatively: ChangeTracker not yet ported. Java Compiler
// delegates to changeTracker.hasScopeChanged(n).
func (c *AbstractCompiler) HasScopeChanged(n *rhino.Node) bool {
	// Conservative default: assume all scopes may have changed.
	_ = n
	return true
}

// ForwardDeclareType records a forward-declared type name.
// Mirrors Java's AbstractCompiler.forwardDeclareType.
// TODO(port): Compiler.forwardDeclareType(String typeName) — requires forwardDeclaredTypes set
func (c *AbstractCompiler) ForwardDeclareType(typeName string) {
	// No-op: forwardDeclaredTypes not yet wired. Java Compiler adds to
	// forwardDeclaredTypes set.
	_ = typeName
}

// ---------------------------------------------------------------------------
// Runtime library direction
// ---------------------------------------------------------------------------

// GetRuntimeLibDir returns the runtime library directory path.
// Mirrors Java's AbstractCompiler.RUNTIME_LIB_DIR.
const RuntimeLibDir = "js"
