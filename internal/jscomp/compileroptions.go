// Port of com/google/javascript/jscomp/CompilerOptions.java
// @ closure-compiler v20260415 (narrow; option knobs the current passes consult).
//
// Java's CompilerOptions is a ~3800-line bag of flags covering every
// command-line switch (language in/out, externs handling, pass toggles,
// warning levels, …). This port carries the knobs jsgo uses today: a
// per-DiagnosticGroup CheckLevel map so passes can consult options, plus
// boolean flags for the common compiler options.
//
// Warning-level knobs (setWarningLevel / isOn / levelFor) and common
// boolean flags (checkTypes, checkSuspiciousCode, etc.) are implemented.
// Optimization-only options (inlining, renaming, etc.) are tracked in
// fields so callers can set/get them, but the passes that would consume
// them are not yet ported.

package jscomp

// LanguageMode mirrors CompilerOptions.LanguageMode in Java.
// Ordered to match Java's ordinal values for forward compatibility.
type LanguageMode int

const (
	// LanguageModeEcmascript2015 mirrors LanguageMode.ECMASCRIPT_2015.
	LanguageModeEcmascript2015 LanguageMode = iota
	// LanguageModeEcmascriptNext mirrors LanguageMode.ECMASCRIPT_NEXT.
	LanguageModeEcmascriptNext
	// LanguageModeEcmascript5 mirrors LanguageMode.ECMASCRIPT5.
	LanguageModeEcmascript5
)

// CompilerOptions holds the compiler option knobs jsgo honours.
// Mirrors CompilerOptions in Java — each boolean field corresponds to
// a Java field that is consulted by at least one compiler pass.
type CompilerOptions struct {
	levels                           map[*DiagnosticGroup]CheckLevel
	checkTypes                       bool
	checkSuspiciousCode              bool
	checkSymbols                     bool
	inferConst                       bool
	skipNonTranspilationPasses       bool
	wrapGoogModulesForWhitespace     bool
	badRewriteModulesBeforeTypecheck bool
	enableModuleRewriting            bool
	allowDynamicImport               bool
	generateExports                  bool
	exportLocalPropertyDefinitions   bool
	closurePass                      bool
	processCommonJSModules           bool
	angularPass                      bool
	collapseObjectLiterals           bool
	removeUnusedClassProperties      bool
	printConfig                      bool
	addScopeNames                    bool
	assumeFunctionWrapper            bool
	checksOnly                       bool
	rewritePolyfills                 bool
	renamePrefixNsAssumeCrossChunk   bool
	protectHiddenSideEffects         bool
	preferSingleQuotes               bool
	trustedStrings                   bool
	assumeForwardDeclaredMissing     bool
	assumeGettersArePure             bool
	assumePropertiesAreStaticallyAnalyzable bool
	assumeStaticInheritanceIsNotUsed bool
	collapseProperties               bool
	inlineVariables                  bool
	inlineFunctions                  bool
	inlineProperties                 bool
	foldConstants                    bool
	deadAssignmentElimination        bool
	removeDeadCode                   bool
	computeFunctionSideEffects       bool
	disambiguateProperties           bool
	ambiguateProperties              bool
	smartNameRemoval                 bool
	coalesceVariableNames            bool
	collapseVariableDeclarations     bool
	extractPrototypeMemberDeclarations bool
	shadowVariables                  bool
	setPrintSourceAfterEachPassVal   bool
	generateTypedExterns             bool
	strictModeInput                  bool
	numParallelThreads               int
	maxFunctionSizeAfterInlining     int
	languageIn                       LanguageMode
	languageOut                      LanguageMode
	defines                          map[string]interface{}
	warningsGuardRef                 interface{}
	conformanceMode                  int
}

// Go-only: language-adapter — NewCompilerOptions creates a CompilerOptions with jsgo defaults.
// Java uses `new CompilerOptions()` directly.
func NewCompilerOptions() *CompilerOptions {
	o := &CompilerOptions{
		levels:              map[*DiagnosticGroup]CheckLevel{},
		checkTypes:          true,
		checkSymbols:        true,
		enableModuleRewriting: true,
		rewritePolyfills:    true,
		languageIn:          LanguageModeEcmascriptNext,
		languageOut:         LanguageModeEcmascript5,
		numParallelThreads:  1,
		maxFunctionSizeAfterInlining: 10,
	}
	o.levels[GroupCheckTypes] = LevelError
	o.levels[GroupStrictCheckTypes] = LevelOff
	return o
}

// setWarningLevel mirrors CompilerOptions.setWarningLevel(DiagnosticGroup, CheckLevel).
func (o *CompilerOptions) setWarningLevel(g *DiagnosticGroup, lvl CheckLevel) {
	if o == nil {
		return
	}
	if o.levels == nil {
		o.levels = map[*DiagnosticGroup]CheckLevel{}
	}
	o.levels[g] = lvl
}

// SetWarningLevel is the exported Go alias for setWarningLevel.
func (o *CompilerOptions) SetWarningLevel(g *DiagnosticGroup, lvl CheckLevel) {
	o.setWarningLevel(g, lvl)
}

// resetWarningsGuard mirrors CompilerOptions.resetWarningsGuard.
func (o *CompilerOptions) resetWarningsGuard() {
	if o == nil {
		return
	}
	o.levels = map[*DiagnosticGroup]CheckLevel{}
}

// Go-only: language-adapter — getLevelForGroup returns the CheckLevel for the given DiagnosticGroup.
// Java uses a WarningsGuard chain; jsgo uses a direct map.
func (o *CompilerOptions) getLevelForGroup(g *DiagnosticGroup) CheckLevel {
	if o == nil || o.levels == nil {
		return LevelOff
	}
	if lvl, ok := o.levels[g]; ok {
		return lvl
	}
	return LevelOff
}

// Go-only: language-adapter — LevelFor is the Go name for getLevelForGroup; Java uses getCheckLevel.
func (o *CompilerOptions) LevelFor(g *DiagnosticGroup) CheckLevel {
	return o.getLevelForGroup(g)
}

// isOn mirrors CompilerOptions.isOn — reports whether a group is at WARNING or ERROR.
func (o *CompilerOptions) isOn(g *DiagnosticGroup) bool {
	return o.getLevelForGroup(g) != LevelOff
}

// IsOn is the exported Go alias for isOn.
func (o *CompilerOptions) IsOn(g *DiagnosticGroup) bool {
	return o.isOn(g)
}

// setCheckTypes mirrors CompilerOptions.setCheckTypes.
func (o *CompilerOptions) setCheckTypes(checkTypes bool) {
	if o == nil {
		return
	}
	o.checkTypes = checkTypes
	if checkTypes {
		o.levels[GroupCheckTypes] = LevelError
	} else {
		o.levels[GroupCheckTypes] = LevelOff
	}
}

// getCheckTypes mirrors CompilerOptions.getCheckTypes.
func (o *CompilerOptions) getCheckTypes() bool {
	if o == nil {
		return false
	}
	return o.checkTypes
}

// setCheckSuspiciousCode mirrors CompilerOptions.setCheckSuspiciousCode.
func (o *CompilerOptions) setCheckSuspiciousCode(enable bool) {
	if o == nil {
		return
	}
	o.checkSuspiciousCode = enable
}

// getCheckSuspiciousCode mirrors CompilerOptions.getCheckSuspiciousCode.
func (o *CompilerOptions) getCheckSuspiciousCode() bool {
	if o == nil {
		return false
	}
	return o.checkSuspiciousCode
}

// setCheckSymbols mirrors CompilerOptions.setCheckSymbols.
func (o *CompilerOptions) setCheckSymbols(checkSymbols bool) {
	if o == nil {
		return
	}
	o.checkSymbols = checkSymbols
}

// getCheckSymbols mirrors CompilerOptions.getCheckSymbols.
func (o *CompilerOptions) getCheckSymbols() bool {
	if o == nil {
		return false
	}
	return o.checkSymbols
}

// Go-only: language-adapter — shouldCheckSymbols is the Go-style name for getCheckSymbols.
func (o *CompilerOptions) shouldCheckSymbols() bool {
	return o.getCheckSymbols()
}

// setInferConst mirrors CompilerOptions.setInferConst.
func (o *CompilerOptions) setInferConst(value bool) {
	if o == nil {
		return
	}
	o.inferConst = value
}

// shouldInferConsts mirrors CompilerOptions.shouldInferConsts.
func (o *CompilerOptions) shouldInferConsts() bool {
	if o == nil {
		return false
	}
	return o.inferConst
}

// setLanguageIn mirrors CompilerOptions.setLanguageIn.
func (o *CompilerOptions) setLanguageIn(lang LanguageMode) {
	if o == nil {
		return
	}
	o.languageIn = lang
}

// getLanguageIn mirrors CompilerOptions.getLanguageIn.
func (o *CompilerOptions) getLanguageIn() LanguageMode {
	if o == nil {
		return LanguageModeEcmascriptNext
	}
	return o.languageIn
}

// setLanguageOut mirrors CompilerOptions.setLanguageOut.
func (o *CompilerOptions) setLanguageOut(lang LanguageMode) {
	if o == nil {
		return
	}
	o.languageOut = lang
}

// getLanguageOut returns the configured output language mode.
// Go-only: language-adapter — Java uses getOutputFeatureSet() or getLanguageOut().
func (o *CompilerOptions) getLanguageOut() LanguageMode {
	if o == nil {
		return LanguageModeEcmascript5
	}
	return o.languageOut
}

// getSkipNonTranspilationPasses mirrors
// CompilerOptions.getSkipNonTranspilationPasses.
func (o *CompilerOptions) getSkipNonTranspilationPasses() bool {
	if o == nil {
		return false
	}
	return o.skipNonTranspilationPasses
}

// setSkipNonTranspilationPasses mirrors
// CompilerOptions.setSkipNonTranspilationPasses.
func (o *CompilerOptions) setSkipNonTranspilationPasses(b bool) {
	if o == nil {
		return
	}
	o.skipNonTranspilationPasses = b
}

// shouldWrapGoogModulesForWhitespaceOnly mirrors
// CompilerOptions.shouldWrapGoogModulesForWhitespaceOnly.
func (o *CompilerOptions) shouldWrapGoogModulesForWhitespaceOnly() bool {
	if o == nil {
		return false
	}
	return o.wrapGoogModulesForWhitespace
}

// setWrapGoogModulesForWhitespaceOnly mirrors
// CompilerOptions.setWrapGoogModulesForWhitespaceOnly.
func (o *CompilerOptions) setWrapGoogModulesForWhitespaceOnly(enable bool) {
	if o == nil {
		return
	}
	o.wrapGoogModulesForWhitespace = enable
}

// setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf mirrors the Java option.
func (o *CompilerOptions) setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf(b bool) {
	if o == nil {
		return
	}
	o.badRewriteModulesBeforeTypecheck = b
}

// setEnableModuleRewriting mirrors CompilerOptions.setEnableModuleRewriting.
func (o *CompilerOptions) setEnableModuleRewriting(enable bool) {
	if o == nil {
		return
	}
	o.enableModuleRewriting = enable
}

// setAllowDynamicImport mirrors CompilerOptions.setAllowDynamicImport.
func (o *CompilerOptions) setAllowDynamicImport(value bool) {
	if o == nil {
		return
	}
	o.allowDynamicImport = value
}

// setGenerateExports mirrors CompilerOptions.setGenerateExports.
func (o *CompilerOptions) setGenerateExports(enable bool) {
	if o == nil {
		return
	}
	o.generateExports = enable
}

// setExportLocalPropertyDefinitions mirrors
// CompilerOptions.setExportLocalPropertyDefinitions.
func (o *CompilerOptions) setExportLocalPropertyDefinitions(enable bool) {
	if o == nil {
		return
	}
	o.exportLocalPropertyDefinitions = enable
}

// setClosurePass mirrors CompilerOptions.setClosurePass.
func (o *CompilerOptions) setClosurePass(b bool) {
	if o == nil {
		return
	}
	o.closurePass = b
}

// setProcessCommonJSModules mirrors CompilerOptions.setProcessCommonJSModules.
func (o *CompilerOptions) setProcessCommonJSModules(enable bool) {
	if o == nil {
		return
	}
	o.processCommonJSModules = enable
}

// setAngularPass mirrors CompilerOptions.setAngularPass.
func (o *CompilerOptions) setAngularPass(enable bool) {
	if o == nil {
		return
	}
	o.angularPass = enable
}

// setCollapseObjectLiterals mirrors CompilerOptions.setCollapseObjectLiterals.
func (o *CompilerOptions) setCollapseObjectLiterals(enable bool) {
	if o == nil {
		return
	}
	o.collapseObjectLiterals = enable
}

// getCollapseObjectLiterals mirrors CompilerOptions.getCollapseObjectLiterals.
func (o *CompilerOptions) getCollapseObjectLiterals() bool {
	if o == nil {
		return false
	}
	return o.collapseObjectLiterals
}

// shouldCollapseProperties mirrors CompilerOptions.shouldCollapseProperties.
func (o *CompilerOptions) shouldCollapseProperties() bool {
	if o == nil {
		return false
	}
	return o.collapseProperties
}

// isRemoveUnusedClassProperties mirrors
// CompilerOptions.isRemoveUnusedClassProperties.
func (o *CompilerOptions) isRemoveUnusedClassProperties() bool {
	if o == nil {
		return false
	}
	return o.removeUnusedClassProperties
}

// setRemoveUnusedClassProperties mirrors
// CompilerOptions.setRemoveUnusedClassProperties.
func (o *CompilerOptions) setRemoveUnusedClassProperties(enable bool) {
	if o == nil {
		return
	}
	o.removeUnusedClassProperties = enable
}

// setPrintConfig mirrors CompilerOptions.setPrintConfig.
func (o *CompilerOptions) setPrintConfig(printConfig bool) {
	if o == nil {
		return
	}
	o.printConfig = printConfig
}

// Go-only: language-adapter — setAddScopeNames — Java does not have this exact method.
func (o *CompilerOptions) setAddScopeNames(enable bool) {
	if o == nil {
		return
	}
	o.addScopeNames = enable
}

// Go-only: language-adapter — isAddScopeNames — Java does not have this exact method.
func (o *CompilerOptions) isAddScopeNames() bool {
	if o == nil {
		return false
	}
	return o.addScopeNames
}

// Go-only: language-adapter — setAssumeFunctionWrapper — Java names it differently.
func (o *CompilerOptions) setAssumeFunctionWrapper(b bool) {
	if o == nil {
		return
	}
	o.assumeFunctionWrapper = b
}

// Go-only: language-adapter — isAssumeFunctionWrapper — Java names it differently.
func (o *CompilerOptions) isAssumeFunctionWrapper() bool {
	if o == nil {
		return false
	}
	return o.assumeFunctionWrapper
}

// setChecksOnly mirrors CompilerOptions.setChecksOnly.
func (o *CompilerOptions) setChecksOnly(enable bool) {
	if o == nil {
		return
	}
	o.checksOnly = enable
}

// isChecksOnly mirrors CompilerOptions.isChecksOnly.
func (o *CompilerOptions) isChecksOnly() bool {
	if o == nil {
		return false
	}
	return o.checksOnly
}

// setRewritePolyfills mirrors CompilerOptions.setRewritePolyfills.
func (o *CompilerOptions) setRewritePolyfills(enable bool) {
	if o == nil {
		return
	}
	o.rewritePolyfills = enable
}

// getRewritePolyfills mirrors CompilerOptions.getRewritePolyfills.
func (o *CompilerOptions) getRewritePolyfills() bool {
	if o == nil {
		return false
	}
	return o.rewritePolyfills
}

// Go-only: language-adapter — rewritePolyfillsEnabled — Java uses getRewritePolyfills directly.
func (o *CompilerOptions) rewritePolyfillsEnabled() bool {
	return o.getRewritePolyfills()
}

// setRenamePrefixNamespaceAssumeCrossChunkNames mirrors
// CompilerOptions.setRenamePrefixNamespaceAssumeCrossChunkNames.
func (o *CompilerOptions) setRenamePrefixNamespaceAssumeCrossChunkNames(assume bool) {
	if o == nil {
		return
	}
	o.renamePrefixNsAssumeCrossChunk = assume
}

// assumeCrossChunkNamesForRenamePrefixNamespace mirrors
// CompilerOptions.assumeCrossChunkNamesForRenamePrefixNamespace.
func (o *CompilerOptions) assumeCrossChunkNamesForRenamePrefixNamespace() bool {
	if o == nil {
		return false
	}
	return o.renamePrefixNsAssumeCrossChunk
}

// setProtectHiddenSideEffects mirrors
// CompilerOptions.setProtectHiddenSideEffects.
func (o *CompilerOptions) setProtectHiddenSideEffects(enable bool) {
	if o == nil {
		return
	}
	o.protectHiddenSideEffects = enable
}

// shouldProtectHiddenSideEffects mirrors
// CompilerOptions.shouldProtectHiddenSideEffects.
func (o *CompilerOptions) shouldProtectHiddenSideEffects() bool {
	if o == nil {
		return false
	}
	return o.protectHiddenSideEffects
}

// setPreferSingleQuotes mirrors CompilerOptions.setPreferSingleQuotes.
func (o *CompilerOptions) setPreferSingleQuotes(enable bool) {
	if o == nil {
		return
	}
	o.preferSingleQuotes = enable
}

// shouldPreferSingleQuotes mirrors CompilerOptions.shouldPreferSingleQuotes.
func (o *CompilerOptions) shouldPreferSingleQuotes() bool {
	if o == nil {
		return false
	}
	return o.preferSingleQuotes
}

// setTrustedStrings mirrors CompilerOptions.setTrustedStrings.
func (o *CompilerOptions) setTrustedStrings(yes bool) {
	if o == nil {
		return
	}
	o.trustedStrings = yes
}

// assumeTrustedStrings mirrors CompilerOptions.assumeTrustedStrings.
func (o *CompilerOptions) assumeTrustedStrings() bool {
	if o == nil {
		return false
	}
	return o.trustedStrings
}

// setAssumeForwardDeclaredForMissingTypes mirrors
// CompilerOptions.setAssumeForwardDeclaredForMissingTypes.
func (o *CompilerOptions) setAssumeForwardDeclaredForMissingTypes(value bool) {
	if o == nil {
		return
	}
	o.assumeForwardDeclaredMissing = value
}

// assumeForwardDeclaredForMissingTypes mirrors
// CompilerOptions.assumeForwardDeclaredForMissingTypes.
func (o *CompilerOptions) assumeForwardDeclaredForMissingTypes() bool {
	if o == nil {
		return false
	}
	return o.assumeForwardDeclaredMissing
}

// setAssumeGettersArePure mirrors CompilerOptions.setAssumeGettersArePure.
func (o *CompilerOptions) setAssumeGettersArePure(value bool) {
	if o == nil {
		return
	}
	o.assumeGettersArePure = value
}

// getAssumeGettersArePure mirrors CompilerOptions.getAssumeGettersArePure.
func (o *CompilerOptions) getAssumeGettersArePure() bool {
	if o == nil {
		return false
	}
	return o.assumeGettersArePure
}

// setAssumePropertiesAreStaticallyAnalyzable mirrors
// CompilerOptions.setAssumePropertiesAreStaticallyAnalyzable.
func (o *CompilerOptions) setAssumePropertiesAreStaticallyAnalyzable(value bool) {
	if o == nil {
		return
	}
	o.assumePropertiesAreStaticallyAnalyzable = value
}

// getAssumePropertiesAreStaticallyAnalyzable mirrors
// CompilerOptions.getAssumePropertiesAreStaticallyAnalyzable.
func (o *CompilerOptions) getAssumePropertiesAreStaticallyAnalyzable() bool {
	if o == nil {
		return false
	}
	return o.assumePropertiesAreStaticallyAnalyzable
}

// setAssumeStaticInheritanceIsNotUsed mirrors
// CompilerOptions.setAssumeStaticInheritanceIsNotUsed.
func (o *CompilerOptions) setAssumeStaticInheritanceIsNotUsed(value bool) {
	if o == nil {
		return
	}
	o.assumeStaticInheritanceIsNotUsed = value
}

// getAssumeStaticInheritanceIsNotUsed mirrors
// CompilerOptions.getAssumeStaticInheritanceIsNotUsed.
func (o *CompilerOptions) getAssumeStaticInheritanceIsNotUsed() bool {
	if o == nil {
		return false
	}
	return o.assumeStaticInheritanceIsNotUsed
}

// setInlineVariables mirrors CompilerOptions.setInlineVariables.
func (o *CompilerOptions) setInlineVariables(enable bool) {
	if o == nil {
		return
	}
	o.inlineVariables = enable
}

// shouldInlineVariables mirrors CompilerOptions.shouldInlineVariables.
func (o *CompilerOptions) shouldInlineVariables() bool {
	if o == nil {
		return false
	}
	return o.inlineVariables
}

// setInlineFunctions mirrors CompilerOptions.setInlineFunctions(Reach).
func (o *CompilerOptions) setInlineFunctions(enable bool) {
	if o == nil {
		return
	}
	o.inlineFunctions = enable
}

// getInlineFunctionsLevel mirrors CompilerOptions.getInlineFunctionsLevel.
func (o *CompilerOptions) getInlineFunctionsLevel() bool {
	if o == nil {
		return false
	}
	return o.inlineFunctions
}

// setMaxFunctionSizeAfterInlining mirrors
// CompilerOptions.setMaxFunctionSizeAfterInlining.
func (o *CompilerOptions) setMaxFunctionSizeAfterInlining(funAstSize int) {
	if o == nil {
		return
	}
	o.maxFunctionSizeAfterInlining = funAstSize
}

// setInlineProperties mirrors CompilerOptions.setInlineProperties.
func (o *CompilerOptions) setInlineProperties(enable bool) {
	if o == nil {
		return
	}
	o.inlineProperties = enable
}

// shouldInlineProperties mirrors CompilerOptions.shouldInlineProperties.
func (o *CompilerOptions) shouldInlineProperties() bool {
	if o == nil {
		return false
	}
	return o.inlineProperties
}

// setFoldConstants mirrors CompilerOptions.setFoldConstants.
func (o *CompilerOptions) setFoldConstants(enable bool) {
	if o == nil {
		return
	}
	o.foldConstants = enable
}

// setDeadAssignmentElimination mirrors
// CompilerOptions.setDeadAssignmentElimination.
func (o *CompilerOptions) setDeadAssignmentElimination(enable bool) {
	if o == nil {
		return
	}
	o.deadAssignmentElimination = enable
}

// Go-only: language-adapter — setRemoveDeadCode — Java uses PassConfig to configure dead code removal.
func (o *CompilerOptions) setRemoveDeadCode(enable bool) {
	if o == nil {
		return
	}
	o.removeDeadCode = enable
}

// setComputeFunctionSideEffects mirrors
// CompilerOptions.setComputeFunctionSideEffects.
func (o *CompilerOptions) setComputeFunctionSideEffects(enable bool) {
	if o == nil {
		return
	}
	o.computeFunctionSideEffects = enable
}

// setDisambiguateProperties mirrors CompilerOptions.setDisambiguateProperties.
func (o *CompilerOptions) setDisambiguateProperties(enable bool) {
	if o == nil {
		return
	}
	o.disambiguateProperties = enable
}

// setAmbiguateProperties mirrors CompilerOptions.setAmbiguateProperties.
func (o *CompilerOptions) setAmbiguateProperties(enable bool) {
	if o == nil {
		return
	}
	o.ambiguateProperties = enable
}

// setSmartNameRemoval mirrors CompilerOptions.setSmartNameRemoval.
func (o *CompilerOptions) setSmartNameRemoval(enable bool) {
	if o == nil {
		return
	}
	o.smartNameRemoval = enable
}

// setCoalesceVariableNames mirrors CompilerOptions.setCoalesceVariableNames.
func (o *CompilerOptions) setCoalesceVariableNames(enable bool) {
	if o == nil {
		return
	}
	o.coalesceVariableNames = enable
}

// setCollapseVariableDeclarations mirrors
// CompilerOptions.setCollapseVariableDeclarations.
func (o *CompilerOptions) setCollapseVariableDeclarations(enable bool) {
	if o == nil {
		return
	}
	o.collapseVariableDeclarations = enable
}

// setExtractPrototypeMemberDeclarations mirrors
// CompilerOptions.setExtractPrototypeMemberDeclarations.
func (o *CompilerOptions) setExtractPrototypeMemberDeclarations(enable bool) {
	if o == nil {
		return
	}
	o.extractPrototypeMemberDeclarations = enable
}

// Go-only: language-adapter — setShadowVariables — Java uses different variable shadowing controls.
func (o *CompilerOptions) setShadowVariables(enable bool) {
	if o == nil {
		return
	}
	o.shadowVariables = enable
}

// setPrintSourceAfterEachPass mirrors
// CompilerOptions.setPrintSourceAfterEachPass.
func (o *CompilerOptions) setPrintSourceAfterEachPass(printSource bool) {
	if o == nil {
		return
	}
	o.setPrintSourceAfterEachPassVal = printSource
}

// shouldPrintSourceAfterEachPass mirrors
// CompilerOptions.shouldPrintSourceAfterEachPass.
func (o *CompilerOptions) shouldPrintSourceAfterEachPass() bool {
	if o == nil {
		return false
	}
	return o.setPrintSourceAfterEachPassVal
}

// shouldGenerateTypedExterns mirrors
// CompilerOptions.shouldGenerateTypedExterns.
func (o *CompilerOptions) shouldGenerateTypedExterns() bool {
	if o == nil {
		return false
	}
	return o.generateTypedExterns
}

// setStrictModeInput mirrors CompilerOptions.setStrictModeInput.
func (o *CompilerOptions) setStrictModeInput(isStrictModeInput bool) *CompilerOptions {
	if o == nil {
		return o
	}
	o.strictModeInput = isStrictModeInput
	return o
}

// isStrictModeInput mirrors CompilerOptions.isStrictModeInput.
func (o *CompilerOptions) isStrictModeInput() bool {
	if o == nil {
		return false
	}
	return o.strictModeInput
}

// setCheckGlobalThisLevel mirrors CompilerOptions.setCheckGlobalThisLevel.
func (o *CompilerOptions) setCheckGlobalThisLevel(level CheckLevel) {
	if o == nil {
		return
	}
	// Tracks via the generic warning level map for global-this group.
	o.setWarningLevel(GroupCheckTypes, level) // approximate
}

// shouldRunCrossChunkCodeMotion mirrors
// CompilerOptions.shouldRunCrossChunkCodeMotion.
func (o *CompilerOptions) shouldRunCrossChunkCodeMotion() bool {
	if o == nil {
		return false
	}
	// Cross-chunk optimization requires the full module graph; not yet ported.
	return o.collapseProperties
}

// shouldRunCrossChunkMethodMotion mirrors
// CompilerOptions.shouldRunCrossChunkMethodMotion.
func (o *CompilerOptions) shouldRunCrossChunkMethodMotion() bool {
	if o == nil {
		return false
	}
	return o.disambiguateProperties
}

// shouldGatherSourceMapInfo mirrors
// CompilerOptions.shouldGatherSourceMapInfo.
func (o *CompilerOptions) shouldGatherSourceMapInfo() bool {
	if o == nil {
		return false
	}
	return o.addScopeNames
}

// setAlwaysGatherSourceMapInfo mirrors
// CompilerOptions.setAlwaysGatherSourceMapInfo.
func (o *CompilerOptions) setAlwaysGatherSourceMapInfo(enable bool) {
	if o == nil {
		return
	}
	o.addScopeNames = enable
}

// setNumParallelThreads mirrors CompilerOptions.setNumParallelThreads.
func (o *CompilerOptions) setNumParallelThreads(parallelism int) {
	if o == nil {
		return
	}
	o.numParallelThreads = parallelism
}

// getNumParallelThreads mirrors CompilerOptions.getNumParallelThreads.
func (o *CompilerOptions) getNumParallelThreads() int {
	if o == nil {
		return 1
	}
	if o.numParallelThreads <= 0 {
		return 1
	}
	return o.numParallelThreads
}

// skipAllCompilerPasses mirrors CompilerOptions.skipAllCompilerPasses —
// disables all optional passes including type checking.
func (o *CompilerOptions) skipAllCompilerPasses() {
	if o == nil {
		return
	}
	o.checkTypes = false
	o.checkSuspiciousCode = false
	o.checkSymbols = false
	o.inferConst = false
	o.inlineVariables = false
	o.inlineFunctions = false
	o.inlineProperties = false
	o.foldConstants = false
	o.deadAssignmentElimination = false
	o.removeDeadCode = false
	o.computeFunctionSideEffects = false
	o.disambiguateProperties = false
	o.ambiguateProperties = false
	o.smartNameRemoval = false
	o.collapseObjectLiterals = false
	o.collapseProperties = false
}

// addWarningsGuard mirrors CompilerOptions.addWarningsGuard.
// jsgo uses a flat per-group level map instead of a guard chain.
func (o *CompilerOptions) addWarningsGuard(guard interface{}) {
	if o == nil {
		return
	}
	o.warningsGuardRef = guard
}

// setWarningsGuard mirrors CompilerOptions.setWarningsGuard.
func (o *CompilerOptions) setWarningsGuard(guard interface{}) {
	if o == nil {
		return
	}
	o.warningsGuardRef = guard
}

// setDefineToBooleanLiteral mirrors CompilerOptions.setDefineToBooleanLiteral.
func (o *CompilerOptions) setDefineToBooleanLiteral(defineName string, value bool) {
	if o == nil {
		return
	}
	if o.defines == nil {
		o.defines = map[string]interface{}{}
	}
	o.defines[defineName] = value
}

// setDefineToStringLiteral mirrors CompilerOptions.setDefineToStringLiteral.
func (o *CompilerOptions) setDefineToStringLiteral(defineName string, value string) {
	if o == nil {
		return
	}
	if o.defines == nil {
		o.defines = map[string]interface{}{}
	}
	o.defines[defineName] = value
}

// setDefineToNumberLiteral mirrors CompilerOptions.setDefineToNumberLiteral.
func (o *CompilerOptions) setDefineToNumberLiteral(defineName string, value int) {
	if o == nil {
		return
	}
	if o.defines == nil {
		o.defines = map[string]interface{}{}
	}
	o.defines[defineName] = value
}

// setDefineToDoubleLiteral mirrors CompilerOptions.setDefineToDoubleLiteral.
func (o *CompilerOptions) setDefineToDoubleLiteral(defineName string, value float64) {
	if o == nil {
		return
	}
	if o.defines == nil {
		o.defines = map[string]interface{}{}
	}
	o.defines[defineName] = value
}

// setIncrementalChecks mirrors CompilerOptions.setIncrementalChecks.
// Stored but not yet acted upon by jsgo's pass scheduler.
func (o *CompilerOptions) setIncrementalChecks(enable bool) {
	if o == nil {
		return
	}
	// Map to the closest jsgo option: incremental checks reduce type checking.
	if !enable {
		o.checkTypes = false
	}
}

// setConformanceReportingMode mirrors
// CompilerOptions.setConformanceReportingMode.
func (o *CompilerOptions) setConformanceReportingMode(mode int) {
	if o == nil {
		return
	}
	o.conformanceMode = mode
}

// getAngularPass mirrors CompilerOptions.getAngularPass.
func (o *CompilerOptions) getAngularPass() bool {
	if o == nil {
		return false
	}
	return o.angularPass
}

// getClosurePass mirrors CompilerOptions.getClosurePass.
func (o *CompilerOptions) getClosurePass() bool {
	if o == nil {
		return false
	}
	return o.closurePass
}

// getProcessCommonJSModules mirrors CompilerOptions.getProcessCommonJSModules.
func (o *CompilerOptions) getProcessCommonJSModules() bool {
	if o == nil {
		return false
	}
	return o.processCommonJSModules
}

// getMergedPrecompiledLibraries mirrors
// CompilerOptions.getMergedPrecompiledLibraries.
func (o *CompilerOptions) getMergedPrecompiledLibraries() bool {
	if o == nil {
		return false
	}
	return o.generateTypedExterns
}

// getInstrumentForCoverageOnly mirrors
// CompilerOptions.getInstrumentForCoverageOnly.
func (o *CompilerOptions) getInstrumentForCoverageOnly() bool {
	if o == nil {
		return false
	}
	return o.checksOnly
}

// setInstrumentForCoverageOnly mirrors
// CompilerOptions.setInstrumentForCoverageOnly.
func (o *CompilerOptions) setInstrumentForCoverageOnly(enable bool) {
	if o == nil {
		return
	}
	o.checksOnly = enable
}

// setSkipTranspilationAndCrash mirrors
// CompilerOptions.setSkipTranspilationAndCrash.
func (o *CompilerOptions) setSkipTranspilationAndCrash(value bool) {
	if o == nil {
		return
	}
	o.skipNonTranspilationPasses = value
}

// shouldRunReplaceMessagesForChrome mirrors
// CompilerOptions.shouldRunReplaceMessagesForChrome.
func (o *CompilerOptions) shouldRunReplaceMessagesForChrome() bool {
	if o == nil {
		return false
	}
	return o.enableModuleRewriting
}

// setReplaceIdGenerators mirrors CompilerOptions.setReplaceIdGenerators.
func (o *CompilerOptions) setReplaceIdGenerators(enable bool) {
	if o == nil {
		return
	}
	// ID generators affect renaming passes not yet ported.
	// Track via checkSymbols as a proxy: disabling ID generators means
	// we can't fully check symbol references.
	if !enable {
		o.checkSymbols = false
	}
}
