package jscomp

import (
	"io"
	"maps"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func (m *customPassMap) add(time CustomPassExecutionTime, pass CompilerPass) {
	for i := range m.groups {
		if m.groups[i].time == time {
			m.groups[i].passes = append(m.groups[i].passes, pass)
			return
		}
	}
	m.groups = append(m.groups, customPassGroup{
		time:   time,
		passes: []CompilerPass{pass},
	})
}

// AddCustomPass registers a custom CompilerPass to run at executionTime.
//
// Java: CompilerOptions.addCustomPass(CustomPassExecutionTime, CompilerPass)
func (o *CompilerOptions) AddCustomPass(executionTime CustomPassExecutionTime, customPass CompilerPass) {
	jutil.CheckNotNull(customPass)
	if o.customPasses == nil {
		o.customPasses = &customPassMap{}
	}
	o.customPasses.add(executionTime, customPass)
}

// SetNameGenerator sets the shared NameGenerator used for variable and property renaming.
//
// Java: CompilerOptions.setNameGenerator(NameGenerator)
func (o *CompilerOptions) SetNameGenerator(nameGenerator NameGenerator) {
	o.nameGenerator = jutil.CheckNotNull(nameGenerator)
}

// SetCssRenamingMap sets the map used in renaming CSS class names.
//
// Java: CompilerOptions.setCssRenamingMap(CssRenamingMap)
func (o *CompilerOptions) SetCssRenamingMap(cssRenamingMap CssRenamingMap) {
	o.cssRenamingMap = cssRenamingMap
}

// SetXidHashFunction sets the hash function used for xid generation.
//
// Java: CompilerOptions.setXidHashFunction(Xid.HashFunction)
func (o *CompilerOptions) SetXidHashFunction(xidHashFunction Xid_HashFunction) {
	o.xidHashFunction = xidHashFunction
}

// SetChunkIdHashFunction sets the hash function used for generating chunk IDs.
//
// Java: CompilerOptions.setChunkIdHashFunction(Xid.HashFunction)
func (o *CompilerOptions) SetChunkIdHashFunction(chunkIdHashFunction Xid_HashFunction) {
	o.chunkIdHashFunction = chunkIdHashFunction
}

// SetStateCompressionWrapper sets an optional stream wrapper applied when
// writing saved compilation state.
//
// Java: CompilerOptions.setStateCompressionWrapper(Function<OutputStream, OutputStream>)
func (o *CompilerOptions) SetStateCompressionWrapper(wrapper func(io.Writer) io.Writer) {
	o.stateCompressionWrapper = wrapper
}

// SetStateDecompressionWrapper sets an optional stream wrapper applied when
// reading saved compilation state.
//
// Java: CompilerOptions.setStateDecompressionWrapper(Function<InputStream, InputStream>)
func (o *CompilerOptions) SetStateDecompressionWrapper(wrapper func(io.Reader) io.Reader) {
	o.stateDecompressionWrapper = wrapper
}

// SetTypedAstOutputStream configures an in-memory output stream to receive
// serialized TypedAST output during checks (wired to
// SerializeTypedAstPass_CreateFromOutputStream).
func (o *CompilerOptions) SetTypedAstOutputStream(out io.Writer) {
	o.typedAstOutputStream = out
}

// GetTypedAstOutputStream returns the configured in-memory TypedAST output stream, or nil.
func (o *CompilerOptions) GetTypedAstOutputStream() io.Writer {
	return o.typedAstOutputStream
}

// HasTypedAstOutput reports whether either a file path or an output stream is
// configured for TypedAST serialization.
func (o *CompilerOptions) HasTypedAstOutput() bool {
	return o.typedAstOutputFile != nil || o.typedAstOutputStream != nil
}

// SetLicenseFilter installs an optional predicate invoked on each file's
// @license/@preserve comment text before emission; returning false suppresses it.
func (o *CompilerOptions) SetLicenseFilter(filter func(string) bool) {
	o.licenseFilter = filter
}

// GetLicenseFilter returns the configured license filter predicate, or nil.
func (o *CompilerOptions) GetLicenseFilter() func(string) bool {
	return o.licenseFilter
}

// SetReconcileOptionsHook installs an optional callback invoked at the end of
// Compiler.ReconcileOptionsWithGuards.
func (o *CompilerOptions) SetReconcileOptionsHook(hook func(*CompilerOptions)) {
	o.reconcileOptionsHook = hook
}

// GetReconcileOptionsHook returns the configured reconcile-options callback, or nil.
func (o *CompilerOptions) GetReconcileOptionsHook() func(*CompilerOptions) {
	return o.reconcileOptionsHook
}

// SetChecksDecorator registers an embedder callback invoked at the end of
// DefaultPassConfig.GetChecks().
func SetChecksDecorator(o *CompilerOptions, fn func(*DefaultPassConfig, *PassListBuilder)) {
	o.checksDecorator = fn
}

// SetOptimizationsDecorator registers an embedder callback invoked at the end of
// DefaultPassConfig.GetOptimizations().
func SetOptimizationsDecorator(o *CompilerOptions, fn func(*DefaultPassConfig, *PassListBuilder)) {
	o.optimizationsDecorator = fn
}

// SetFinalizationsDecorator registers an embedder callback invoked at the end of
// DefaultPassConfig.GetFinalizations().
func SetFinalizationsDecorator(o *CompilerOptions, fn func(*DefaultPassConfig, *PassListBuilder)) {
	o.finalizationsDecorator = fn
}

// SetExtraOption stores an embedder-defined key-value pair on CompilerOptions.
func SetExtraOption(o *CompilerOptions, key string, val any) {
	if o.extraOptions == nil {
		o.extraOptions = make(map[string]any)
	}
	o.extraOptions[key] = val
}

// ExtraOption retrieves an embedder-defined option value from CompilerOptions.
func ExtraOption(o *CompilerOptions, key string) (any, bool) {
	if o == nil || o.extraOptions == nil {
		return nil, false
	}
	v, ok := o.extraOptions[key]
	return v, ok
}

// ExtraOptionsSnapshot returns a shallow copy of all embedder-defined extra
// options attached to CompilerOptions, or nil if none have been set.
func ExtraOptionsSnapshot(o *CompilerOptions) map[string]any {
	if o == nil || len(o.extraOptions) == 0 {
		return nil
	}
	return maps.Clone(o.extraOptions)
}

// AddOptionWarning queues a deferred option warning emitted as
// StrictWarningsGuard_UNRAISABLE_WARNING when input parsing begins.
func AddOptionWarning(o *CompilerOptions, msg string) {
	o.optionWarnings = append(o.optionWarnings, msg)
}

// LanguageOutIsDefaultStrict returns the tri-state languageOutIsDefaultStrict
// setting on CompilerOptions.
func LanguageOutIsDefaultStrict(o *CompilerOptions) base.Tri {
	return o.languageOutIsDefaultStrict
}

// OptimizationLoopMaxIterations returns the configured maximum optimization
// loop iterations on CompilerOptions (0 means default).
func OptimizationLoopMaxIterations(o *CompilerOptions) int {
	return o.optimizationLoopMaxIterations
}

// PreserveClosurePrimitives reports whether preserveClosurePrimitives is enabled.
func PreserveClosurePrimitives(o *CompilerOptions) bool {
	return o.preserveClosurePrimitives
}

// EnableModuleRewriting reports whether module rewriting is enabled.
func EnableModuleRewriting(o *CompilerOptions) bool {
	return o.enableModuleRewriting
}

// ReplaceMessagesWithChromeI18n returns whether Chrome i18n message replacement
// is enabled together with its optional Translation Console project ID.
func ReplaceMessagesWithChromeI18n(o *CompilerOptions) (bool, base.OptString) {
	return o.replaceMessagesWithChromeI18n, o.tcProjectId
}

// ClearConformanceConfigs clears all configured conformance configs.
//
// Java: CompilerOptions.clearConformanceConfigs()
func (o *CompilerOptions) ClearConformanceConfigs() {
	o.conformanceConfigs = []*ConformanceConfig{}
}

// SetExportTestFunctions enables or disables exporting test functions.
//
// Java: CompilerOptions.setExportTestFunctions(boolean)
func (o *CompilerOptions) SetExportTestFunctions(export bool) {
	o.exportTestFunctions = export
}

// SetIdGeneratorsOfMap sets the id generators map from an explicit RenamingMap map.
//
// Java: CompilerOptions.setIdGenerators(Map<String, RenamingMap>)
func (o *CompilerOptions) SetIdGeneratorsOfMap(idGenerators *OrderedStringMap[RenamingMap]) {
	jutil.CheckNotNull(idGenerators)
	cloned := idGenerators.Clone()
	o.idGenerators = &cloned
}

// SetInferConst enables or disables InferConsts.
//
// Java: CompilerOptions.setInferConst(boolean)
func (o *CompilerOptions) SetInferConst(infer bool) {
	o.inferConsts = infer
}

// SetJ2clMinifierEnabled enables or disables the J2CL minifier.
//
// Java: CompilerOptions.setJ2clMinifierEnabled(boolean)
func (o *CompilerOptions) SetJ2clMinifierEnabled(enabled bool) {
	o.j2clMinifierEnabled = enabled
}

// SetOutputJsStringUsage enables or disables JS string usage reporting.
//
// Java: CompilerOptions.setOutputJsStringUsage(boolean)
func (o *CompilerOptions) SetOutputJsStringUsage(output bool) {
	o.outputJsStringUsage = output
}

// SetParentChunkCanSeeSymbolsDeclaredInChildren configures cross-chunk symbol visibility.
//
// Java: CompilerOptions.setParentChunkCanSeeSymbolsDeclaredInChildren(boolean)
func (o *CompilerOptions) SetParentChunkCanSeeSymbolsDeclaredInChildren(canSee bool) {
	o.parentChunkCanSeeSymbolsDeclaredInChildren = canSee
}

// SetPrintExterns configures whether output printing emits externs instead of JS sources.
//
// Java: CompilerOptions.setPrintExterns(boolean)
func (o *CompilerOptions) SetPrintExterns(printExterns bool) {
	o.printExterns = printExterns
}

// SetPropertyRenamingOnlyCompilationMode enables or disables property-renaming-only mode.
//
// Java: CompilerOptions.setPropertyRenamingOnlyCompilationMode(boolean)
func (o *CompilerOptions) SetPropertyRenamingOnlyCompilationMode(enabled bool) {
	o.propertyRenamingOnlyCompilationMode = enabled
}

// SetReplaceStringsFunctionDescriptions sets the ReplaceStrings function descriptors.
//
// Java: CompilerOptions.setReplaceStringsFunctionDescriptions(List<String>)
func (o *CompilerOptions) SetReplaceStringsFunctionDescriptions(descriptions []string) {
	jutil.CheckNotNull(descriptions)
	o.replaceStringsFunctionDescriptions = append([]string{}, descriptions...)
}

// SetReplaceStringsPlaceholderToken sets the ReplaceStrings placeholder token.
//
// Java: CompilerOptions.setReplaceStringsPlaceholderToken(String)
func (o *CompilerOptions) SetReplaceStringsPlaceholderToken(token string) {
	o.replaceStringsPlaceholderToken = token
}

// SetSerializeExtraDebugInfo configures whether TypedAST serialization includes debug info.
//
// Java: CompilerOptions.setSerializeExtraDebugInfo(boolean)
func (o *CompilerOptions) SetSerializeExtraDebugInfo(serialize bool) {
	o.serializeExtraDebugInfo = serialize
}

// SetStrictMessageReplacement enables or disables strict message replacement.
//
// Java: CompilerOptions.setStrictMessageReplacement(boolean)
func (o *CompilerOptions) SetStrictMessageReplacement(strict bool) {
	o.strictMessageReplacement = strict
}

// JenkinsHash32WithSeed computes Bob Jenkins' 32-bit lookup2 hash over data with
// the given initial seed (matching GoogleJsMessageIdGenerator.FP.hash32).
func JenkinsHash32WithSeed(data []byte, seed int32) int32 {
	return int32(googleJsMessageIdGenerator_FP_hash32(data, uint32(seed)))
}
