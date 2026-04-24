// Port of com/google/javascript/jscomp/CompilerOptions.java
// @ closure-compiler v20260415 (narrow; option knobs the current passes consult).
//
// Java's CompilerOptions is a ~3800-line bag of flags covering every
// command-line switch (language in/out, externs handling, pass toggles,
// warning levels, …). This port carries the knobs jsgo uses today plus
// stubs for the Java public API so type-checking passes can call the
// same method names they use in Java. Defaults match Java's defaults
// for each option.
//
// Warning-level knobs (setWarningLevel / isOn / levelFor) plus a
// selection of common boolean flags (checkTypes, checkSuspiciousCode,
// skipNonTranspilationPasses, …) are implemented. Everything else is
// a stub returning an appropriate zero value.

package jscomp

// LanguageMode mirrors CompilerOptions.LanguageMode in Java.
// Only the subset jsgo cares about today.
type LanguageMode int

const (
	// LanguageModeEcmascript2015 is the default input language.
	LanguageModeEcmascript2015 LanguageMode = iota
	// LanguageModeEcmascriptNext accepts all current ES features.
	LanguageModeEcmascriptNext
	// LanguageModeEcmascript5 is ES5 output / input.
	LanguageModeEcmascript5
)

// CompilerOptions holds the subset of --jscomp_* knobs jsgo honours.
// Additional fields (module resolution, path renaming, …) land as later
// tiers plumb them through the Compiler shim.
type CompilerOptions struct {
	levels                    map[*DiagnosticGroup]CheckLevel
	checkTypes                bool
	checkSuspiciousCode       bool
	checkSymbols              bool
	inferConst                bool
	skipNonTranspilationPasses bool
	languageIn                LanguageMode
	languageOut               LanguageMode
	wrapGoogModulesForWhitespace bool
	inlineVariables           bool
	inlineFunctions           bool
	collapseObjectLiterals    bool
	removeUnusedClassProperties bool
	printConfig               bool
	numParallelThreads        int
	assumeGettersArePure      bool
	assumeStaticInheritanceIsNotUsed bool
}

// NewCompilerOptions returns options pre-populated with Java's defaults.
// Mirrors CompilerOptions() constructor's defaults.
func NewCompilerOptions() *CompilerOptions {
	o := &CompilerOptions{
		levels:      map[*DiagnosticGroup]CheckLevel{},
		checkTypes:  true, // --jscomp_error=checkTypes by default
		checkSymbols: true,
		languageIn:  LanguageModeEcmascriptNext,
		languageOut: LanguageModeEcmascript5,
	}
	o.levels[GroupCheckTypes] = LevelError
	o.levels[GroupStrictCheckTypes] = LevelOff
	return o
}

// SetWarningLevel mirrors CompilerOptions.setWarningLevel(DiagnosticGroup, CheckLevel).
func (o *CompilerOptions) SetWarningLevel(g *DiagnosticGroup, lvl CheckLevel) {
	if o == nil {
		return
	}
	if o.levels == nil {
		o.levels = map[*DiagnosticGroup]CheckLevel{}
	}
	o.levels[g] = lvl
}

// resetWarningsGuard mirrors CompilerOptions.resetWarningsGuard — clears
// all per-group overrides, reverting to defaults.
func (o *CompilerOptions) resetWarningsGuard() {
	if o == nil {
		return
	}
	o.levels = map[*DiagnosticGroup]CheckLevel{}
}

// Go-only: LevelFor returns the current level configured for `g`, or LevelOff.
// Java uses DiagnosticGroup.getCheckLevel(warnings) via a WarningsGuard chain.
func (o *CompilerOptions) LevelFor(g *DiagnosticGroup) CheckLevel {
	if o == nil || o.levels == nil {
		return LevelOff
	}
	return o.levels[g]
}

// isOn mirrors CompilerOptions.isOn — reports whether the group is at
// WARNING or ERROR level. Mirrors the Java isOn(DiagnosticGroup) helper.
func (o *CompilerOptions) isOn(g *DiagnosticGroup) bool {
	return o.LevelFor(g) != LevelOff
}

// IsOn is the exported Go alias for isOn, used by passes and tests.
func (o *CompilerOptions) IsOn(g *DiagnosticGroup) bool {
	return o.isOn(g)
}

// setCheckTypes mirrors CompilerOptions.setCheckTypes(boolean).
// Enables / disables the full type-checking pass.
func (o *CompilerOptions) setCheckTypes(checkTypes bool) {
	if o == nil {
		return
	}
	o.checkTypes = checkTypes
}

// getCheckTypes mirrors CompilerOptions.getCheckTypes().
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

// shouldCheckSymbols mirrors CompilerOptions.checkSymbols.
func (o *CompilerOptions) shouldCheckSymbols() bool {
	if o == nil {
		return false
	}
	return o.checkSymbols
}

// setInferConst mirrors CompilerOptions.setInferConst(boolean).
func (o *CompilerOptions) setInferConst(value bool) {
	if o == nil {
		return
	}
	o.inferConst = value
}

// shouldInferConsts mirrors CompilerOptions.shouldInferConsts().
func (o *CompilerOptions) shouldInferConsts() bool {
	if o == nil {
		return false
	}
	return o.inferConst
}

// setLanguageIn mirrors CompilerOptions.setLanguageIn(LanguageMode).
func (o *CompilerOptions) setLanguageIn(lang LanguageMode) {
	if o == nil {
		return
	}
	o.languageIn = lang
}

// getLanguageIn mirrors CompilerOptions.getLanguageIn().
func (o *CompilerOptions) getLanguageIn() LanguageMode {
	if o == nil {
		return LanguageModeEcmascriptNext
	}
	return o.languageIn
}

// setLanguageOut mirrors CompilerOptions.setLanguageOut(LanguageMode).
func (o *CompilerOptions) setLanguageOut(lang LanguageMode) {
	if o == nil {
		return
	}
	o.languageOut = lang
}

// getLanguageOut mirrors CompilerOptions.getLanguageOut() (returns output mode).
func (o *CompilerOptions) getLanguageOut() LanguageMode {
	if o == nil {
		return LanguageModeEcmascript5
	}
	return o.languageOut
}

// getSkipNonTranspilationPasses mirrors
// CompilerOptions.getSkipNonTranspilationPasses().
func (o *CompilerOptions) getSkipNonTranspilationPasses() bool {
	if o == nil {
		return false
	}
	return o.skipNonTranspilationPasses
}

// shouldWrapGoogModulesForWhitespaceOnly mirrors
// CompilerOptions.shouldWrapGoogModulesForWhitespaceOnly().
func (o *CompilerOptions) shouldWrapGoogModulesForWhitespaceOnly() bool {
	if o == nil {
		return false
	}
	return o.wrapGoogModulesForWhitespace
}

// setWrapGoogModulesForWhitespaceOnly mirrors
// CompilerOptions.setWrapGoogModulesForWhitespaceOnly(boolean).
func (o *CompilerOptions) setWrapGoogModulesForWhitespaceOnly(enable bool) {
	if o == nil {
		return
	}
	o.wrapGoogModulesForWhitespace = enable
}

// setInlineVariables mirrors CompilerOptions.setInlineVariables(boolean).
func (o *CompilerOptions) setInlineVariables(enable bool) {
	if o == nil {
		return
	}
	o.inlineVariables = enable
}

// shouldInlineVariables mirrors CompilerOptions.shouldInlineVariables().
func (o *CompilerOptions) shouldInlineVariables() bool {
	if o == nil {
		return false
	}
	return o.inlineVariables
}

// setInlineFunctions mirrors CompilerOptions.setInlineFunctions(Reach).
// Simplified to a boolean (no Reach enum in jsgo yet).
func (o *CompilerOptions) setInlineFunctions(enable bool) {
	if o == nil {
		return
	}
	o.inlineFunctions = enable
}

// setCollapseObjectLiterals mirrors
// CompilerOptions.setCollapseObjectLiterals(boolean).
func (o *CompilerOptions) setCollapseObjectLiterals(enable bool) {
	if o == nil {
		return
	}
	o.collapseObjectLiterals = enable
}

// getCollapseObjectLiterals mirrors
// CompilerOptions.getCollapseObjectLiterals().
func (o *CompilerOptions) getCollapseObjectLiterals() bool {
	if o == nil {
		return false
	}
	return o.collapseObjectLiterals
}

// isRemoveUnusedClassProperties mirrors
// CompilerOptions.isRemoveUnusedClassProperties().
func (o *CompilerOptions) isRemoveUnusedClassProperties() bool {
	if o == nil {
		return false
	}
	return o.removeUnusedClassProperties
}

// setRemoveUnusedClassProperties mirrors
// CompilerOptions.setRemoveUnusedClassProperties(boolean).
func (o *CompilerOptions) setRemoveUnusedClassProperties(enable bool) {
	if o == nil {
		return
	}
	o.removeUnusedClassProperties = enable
}

// skipAllCompilerPasses mirrors CompilerOptions.skipAllCompilerPasses().
// Disables all optional passes (turns off checkTypes etc.).
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
}

// setPrintConfig mirrors CompilerOptions.setPrintConfig(boolean).
func (o *CompilerOptions) setPrintConfig(printConfig bool) {
	if o == nil {
		return
	}
	o.printConfig = printConfig
}

// setNumParallelThreads mirrors CompilerOptions.setNumParallelThreads(int).
func (o *CompilerOptions) setNumParallelThreads(parallelism int) {
	if o == nil {
		return
	}
	o.numParallelThreads = parallelism
}

// getNumParallelThreads mirrors CompilerOptions.getNumParallelThreads().
func (o *CompilerOptions) getNumParallelThreads() int {
	if o == nil {
		return 1
	}
	if o.numParallelThreads <= 0 {
		return 1
	}
	return o.numParallelThreads
}

// setAssumeGettersArePure mirrors
// CompilerOptions.setAssumeGettersArePure(boolean).
func (o *CompilerOptions) setAssumeGettersArePure(value bool) {
	if o == nil {
		return
	}
	o.assumeGettersArePure = value
}

// getAssumeGettersArePure mirrors CompilerOptions.getAssumeGettersArePure().
func (o *CompilerOptions) getAssumeGettersArePure() bool {
	if o == nil {
		return false
	}
	return o.assumeGettersArePure
}

// setAssumeStaticInheritanceIsNotUsed mirrors
// CompilerOptions.setAssumeStaticInheritanceIsNotUsed(boolean).
func (o *CompilerOptions) setAssumeStaticInheritanceIsNotUsed(value bool) {
	if o == nil {
		return
	}
	o.assumeStaticInheritanceIsNotUsed = value
}

// getAssumeStaticInheritanceIsNotUsed mirrors
// CompilerOptions.getAssumeStaticInheritanceIsNotUsed().
func (o *CompilerOptions) getAssumeStaticInheritanceIsNotUsed() bool {
	if o == nil {
		return false
	}
	return o.assumeStaticInheritanceIsNotUsed
}
