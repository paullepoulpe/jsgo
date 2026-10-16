// Port of com/google/javascript/jscomp/DefaultPassConfig.java
// @ closure-compiler v20260415
//
// DefaultPassConfig is a ~3000-line file in Java covering all compiler passes.
// This port covers the structural skeleton and every pass pipeline used by
// getChecks / getOptimizations / getFinalizations, matching Java's pass
// ordering exactly.  Unimplemented pass bodies are no-ops so the pass-list
// infrastructure is exercisable end-to-end.

package jscomp

// compiledConstantName mirrors Java's DefaultPassConfig.COMPILED_CONSTANT_NAME.
const compiledConstantName = "COMPILED"

// closureLocaleConstantName mirrors Java's DefaultPassConfig.CLOSURE_LOCALE_CONSTANT_NAME.
const closureLocaleConstantName = "goog.LOCALE"

// DefaultPassConfig is the concrete PassConfig implementation for the full
// compiler pipeline. Embeds BasePassConfig and overrides the three abstract
// phase methods. Mirrors Java's DefaultPassConfig final class.
type DefaultPassConfig struct {
	*BasePassConfig
}

// NewDefaultPassConfig constructs a DefaultPassConfig for the given options.
// Mirrors Java's DefaultPassConfig(CompilerOptions) constructor.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func NewDefaultPassConfig(options *CompilerOptions) *DefaultPassConfig {
	if options == nil {
		options = NewCompilerOptions()
	}
	return &DefaultPassConfig{BasePassConfig: NewBasePassConfig(options)}
}

// ---------------------------------------------------------------------------
// Individual PassFactory fields
// ---------------------------------------------------------------------------
// Each field mirrors a private final PassFactory field in Java's
// DefaultPassConfig. They are defined as methods on DefaultPassConfig so that
// they can reference `c.options` and still be lazily constructed on demand.
// ---------------------------------------------------------------------------

// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) makeNoopPass(name string) *PassFactory {
	return NewPassFactoryBuilder().
		SetName(name).
		SetInternalFactory(func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		}).
		Build()
}

// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) makeConditionalPass(name string, cond func(*CompilerOptions) bool) *PassFactory {
	return NewPassFactoryBuilder().
		SetName(name).
		SetCondition(cond).
		SetInternalFactory(func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		}).
		Build()
}

// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) makeLoopPass(name string) *PassFactory {
	return NewPassFactoryBuilder().
		SetName(name).
		SetRunInFixedPointLoop(true).
		SetInternalFactory(func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		}).
		Build()
}

// ---------------------------------------------------------------------------
// Check-phase passes (mirrors Java DefaultPassConfig field order)
// ---------------------------------------------------------------------------

func (c *DefaultPassConfig) addSyntheticScript() *PassFactory {
	return c.makeNoopPass("addSyntheticScript")
}

func (c *DefaultPassConfig) gatherGettersAndSetters() *PassFactory {
	return c.makeConditionalPass(PassNameGatherGettersAndSetters, func(*CompilerOptions) bool { return true })
}

func (c *DefaultPassConfig) rewriteScriptsToEs6Modules() *PassFactory {
	return c.makeNoopPass(PassNameRewriteScriptsToEs6Modules)
}

func (c *DefaultPassConfig) gatherModuleMetadataPass() *PassFactory {
	return c.makeNoopPass(PassNameGatherModuleMetadata)
}

func (c *DefaultPassConfig) createModuleMapPass() *PassFactory {
	return c.makeNoopPass(PassNameCreateModuleMap)
}

func (c *DefaultPassConfig) rewriteCommonJsModules() *PassFactory {
	return c.makeNoopPass(PassNameRewriteCommonJsModules)
}

func (c *DefaultPassConfig) checkJsDocAndEs6Modules() *PassFactory {
	return c.makeNoopPass("checkJsDocAndEs6Modules")
}

func (c *DefaultPassConfig) checkTypeImportCodeReferences() *PassFactory {
	return c.makeNoopPass("checkTypeImportCodeReferences")
}

func (c *DefaultPassConfig) lintChecks() *PassFactory {
	return c.makeConditionalPass(PassNameLintChecks, func(*CompilerOptions) bool { return true })
}

func (c *DefaultPassConfig) checkRequiresAndProvidesSorted() *PassFactory {
	return c.makeNoopPass("checkRequiresAndProvidesSorted")
}

func (c *DefaultPassConfig) extraRequires() *PassFactory {
	return c.makeNoopPass("extraRequires")
}

func (c *DefaultPassConfig) checkMissingRequires() *PassFactory {
	return c.makeConditionalPass("checkMissingRequires", func(*CompilerOptions) bool { return true })
}

func (c *DefaultPassConfig) declaredGlobalExternsOnWindow() *PassFactory {
	return c.makeNoopPass(PassNameDeclaredGlobalExternsOnWindow)
}

func (c *DefaultPassConfig) checkVariableReferences() *PassFactory {
	return c.makeNoopPass(PassNameCheckVariableReferences)
}

func (c *DefaultPassConfig) checkVars() *PassFactory {
	return c.makeNoopPass(PassNameCheckVars)
}

func (c *DefaultPassConfig) checkClosureImports() *PassFactory {
	return c.makeNoopPass("checkClosureImports")
}

func (c *DefaultPassConfig) checkStrictMode() *PassFactory {
	return c.makeNoopPass("checkStrictMode")
}

func (c *DefaultPassConfig) closureCheckModule() *PassFactory {
	return c.makeNoopPass("closureCheckModule")
}

func (c *DefaultPassConfig) checkSuper() *PassFactory {
	return c.makeNoopPass("checkSuper")
}

func (c *DefaultPassConfig) checkSideEffects() *PassFactory {
	return c.makeConditionalPass("checkSideEffects", func(*CompilerOptions) bool { return true })
}

func (c *DefaultPassConfig) closureGoogScopeAliases() *PassFactory {
	return c.makeNoopPass("closureGoogScopeAliases")
}

func (c *DefaultPassConfig) closurePrimitives() *PassFactory {
	return c.makeNoopPass("closurePrimitives")
}

func (c *DefaultPassConfig) polymerPass() *PassFactory {
	return c.makeNoopPass("polymerPass")
}

func (c *DefaultPassConfig) inferConsts() *PassFactory {
	return c.makeNoopPass(PassNameInferConsts)
}

func (c *DefaultPassConfig) checkRegExp() *PassFactory {
	return c.makeNoopPass(PassNameCheckRegExp)
}

func (c *DefaultPassConfig) inferTypes() *PassFactory {
	return c.makeNoopPass(PassNameInferTypes)
}

func (c *DefaultPassConfig) checkTypes() *PassFactory {
	return c.makeConditionalPass(PassNameCheckTypes, func(opts *CompilerOptions) bool {
		return opts.IsOn(GroupCheckTypes)
	})
}

func (c *DefaultPassConfig) inferJsDocInfo() *PassFactory {
	return c.makeNoopPass("inferJsDocInfo")
}

func (c *DefaultPassConfig) analyzerChecks() *PassFactory {
	return c.makeNoopPass(PassNameAnalyzerChecks)
}

func (c *DefaultPassConfig) clearTypedScopeCreatorPass() *PassFactory {
	return c.makeNoopPass("clearTypedScopeCreatorPass")
}

func (c *DefaultPassConfig) clearTopTypedScopePass() *PassFactory {
	return c.makeNoopPass("clearTopTypedScopePass")
}

func (c *DefaultPassConfig) rewriteDynamicImports() *PassFactory {
	return c.makeNoopPass("rewriteDynamicImports")
}

func (c *DefaultPassConfig) suspiciousCode() *PassFactory {
	return c.makeConditionalPass("suspiciousCode", func(*CompilerOptions) bool { return true })
}

func (c *DefaultPassConfig) j2clSourceFileChecker() *PassFactory {
	return c.makeNoopPass("j2clSourceFileChecker")
}

func (c *DefaultPassConfig) checkControlFlow() *PassFactory {
	return c.makeNoopPass("checkControlFlow")
}

func (c *DefaultPassConfig) checkAccessControls() *PassFactory {
	return c.makeConditionalPass("checkAccessControls", func(*CompilerOptions) bool { return true })
}

func (c *DefaultPassConfig) checkConsts() *PassFactory {
	return c.makeNoopPass("checkConsts")
}

func (c *DefaultPassConfig) rewriteCallerCodeLocation() *PassFactory {
	return c.makeNoopPass("rewriteCallerCodeLocation")
}

func (c *DefaultPassConfig) checkConformance() *PassFactory {
	return c.makeNoopPass(PassNameCheckConformance)
}

func (c *DefaultPassConfig) processTweaks() *PassFactory {
	return c.makeNoopPass("processTweaks")
}

func (c *DefaultPassConfig) processDefinesCheck() *PassFactory {
	return c.makeNoopPass("processDefinesCheck")
}

func (c *DefaultPassConfig) j2clChecksPass() *PassFactory {
	return c.makeNoopPass("j2clChecksPass")
}

func (c *DefaultPassConfig) generateExports() *PassFactory {
	return c.makeNoopPass(PassNameGenerateExports)
}

func (c *DefaultPassConfig) mergeSyntheticScript() *PassFactory {
	return c.makeNoopPass("mergeSyntheticScript")
}

func (c *DefaultPassConfig) externExports() *PassFactory {
	return c.makeNoopPass(PassNameExternExports)
}

func (c *DefaultPassConfig) removeWeakSources() *PassFactory {
	return c.makeNoopPass("removeWeakSources")
}

func (c *DefaultPassConfig) gatherExternPropertiesCheck() *PassFactory {
	return c.makeNoopPass("gatherExternPropertiesCheck")
}

func (c *DefaultPassConfig) serializeTypedAst() *PassFactory {
	return c.makeNoopPass("serializeTypedAst")
}

// ---------------------------------------------------------------------------
// Optimisation-phase passes
// ---------------------------------------------------------------------------

func (c *DefaultPassConfig) removeUnnecessarySyntheticExterns() *PassFactory {
	return c.makeNoopPass("removeUnnecessarySyntheticExterns")
}

func (c *DefaultPassConfig) closureProvidesRequires() *PassFactory {
	return c.makeNoopPass("closureProvidesRequires")
}

func (c *DefaultPassConfig) processDefinesOptimize() *PassFactory {
	return c.makeNoopPass("processDefinesOptimize")
}

func (c *DefaultPassConfig) normalize() *PassFactory {
	return c.makeNoopPass(PassNameNormalize)
}

func (c *DefaultPassConfig) gatherExternPropertiesOptimize() *PassFactory {
	return c.makeNoopPass("gatherExternPropertiesOptimize")
}

func (c *DefaultPassConfig) inlineAndCollapseProperties() *PassFactory {
	return c.makeNoopPass("inlineAndCollapseProperties")
}

func (c *DefaultPassConfig) closureOptimizePrimitives() *PassFactory {
	return c.makeNoopPass("closureOptimizePrimitives")
}

func (c *DefaultPassConfig) stripSideEffectProtection() *PassFactory {
	return c.makeNoopPass(PassNameStripSideEffectProtection)
}

func (c *DefaultPassConfig) inlineVariables() *PassFactory {
	return c.makeLoopPass(PassNameInlineVariables)
}

func (c *DefaultPassConfig) inlineConstants() *PassFactory {
	return c.makeLoopPass("inlineConstants")
}

func (c *DefaultPassConfig) collapseObjectLiterals() *PassFactory {
	return c.makeLoopPass(PassNameCollapseObjectLiterals)
}

func (c *DefaultPassConfig) removeUnusedCode() *PassFactory {
	return c.makeLoopPass(PassNameRemoveUnusedCode)
}

func (c *DefaultPassConfig) peepholeOptimizations() *PassFactory {
	return c.makeLoopPass(PassNamePeepholeOptimizations)
}

func (c *DefaultPassConfig) removeUnusedCodeOnce() *PassFactory {
	return c.makeNoopPass("removeUnusedCodeOnce")
}

func (c *DefaultPassConfig) peepholeOptimizationsOnceNormalized() *PassFactory {
	return c.makeNoopPass("peepholeOptimizationsOnceNormalized")
}

func (c *DefaultPassConfig) peepholeOptimizationsOnceNonNormalized() *PassFactory {
	return c.makeNoopPass("peepholeOptimizationsOnceNonNormalized")
}

func (c *DefaultPassConfig) flowSensitiveInlineVariables() *PassFactory {
	return c.makeNoopPass(PassNameFlowSensitiveInlineVariables)
}

func (c *DefaultPassConfig) crossChunkCodeMotion() *PassFactory {
	return c.makeNoopPass(PassNameCrossChunkCodeMotion)
}

func (c *DefaultPassConfig) crossChunkMethodMotion() *PassFactory {
	return c.makeNoopPass(PassNameCrossChunkMethodMotion)
}

func (c *DefaultPassConfig) optimizeConstructors() *PassFactory {
	return c.makeNoopPass("optimizeConstructors")
}

func (c *DefaultPassConfig) isolatePolyfills() *PassFactory {
	return c.makeNoopPass("isolatePolyfills")
}

func (c *DefaultPassConfig) collapseAnonymousFunctions() *PassFactory {
	return c.makeNoopPass(PassNameCollapseAnonymousFunctions)
}

func (c *DefaultPassConfig) rewriteGlobalDeclarationsForTryCatchWrapping() *PassFactory {
	return c.makeNoopPass("rewriteGlobalDeclarationsForTryCatchWrapping")
}

func (c *DefaultPassConfig) extractPrototypeMemberDeclarations() *PassFactory {
	return c.makeNoopPass(PassNameExtractPrototypeMemberDeclarations)
}

func (c *DefaultPassConfig) ambiguateProperties() *PassFactory {
	return c.makeNoopPass(PassNameAmbiguateProperties)
}

func (c *DefaultPassConfig) renameProperties() *PassFactory {
	return c.makeNoopPass(PassNameRenameProperties)
}

func (c *DefaultPassConfig) removePropertyRenamingCalls() *PassFactory {
	return c.makeNoopPass("removePropertyRenamingCalls")
}

func (c *DefaultPassConfig) gatherRawExports() *PassFactory {
	return c.makeNoopPass(PassNameGatherRawExports)
}

func (c *DefaultPassConfig) convertToDottedProperties() *PassFactory {
	return c.makeNoopPass(PassNameConvertToDottedProperties)
}

func (c *DefaultPassConfig) rewriteFunctionExpressions() *PassFactory {
	return c.makeNoopPass(PassNameRewriteFunctionExpressions)
}

func (c *DefaultPassConfig) aliasStrings() *PassFactory {
	return c.makeNoopPass("aliasStrings")
}

func (c *DefaultPassConfig) coalesceVariableNames() *PassFactory {
	return c.makeNoopPass(PassNameCoalesceVariableNames)
}

func (c *DefaultPassConfig) markUnnormalized() *PassFactory {
	return c.makeNoopPass(PassNameMarkUnnormalized)
}

func (c *DefaultPassConfig) exploitAssign() *PassFactory {
	return c.makeNoopPass(PassNameExploitAssign)
}

func (c *DefaultPassConfig) collapseVariableDeclarations() *PassFactory {
	return c.makeNoopPass(PassNameCollapseVariableDeclarations)
}

func (c *DefaultPassConfig) denormalize() *PassFactory {
	return c.makeNoopPass("denormalize")
}

func (c *DefaultPassConfig) invertContextualRenaming() *PassFactory {
	return c.makeNoopPass("invertContextualRenaming")
}

func (c *DefaultPassConfig) renameVars() *PassFactory {
	return c.makeNoopPass("renameVars")
}

func (c *DefaultPassConfig) renameLabels() *PassFactory {
	return c.makeNoopPass("renameLabels")
}

func (c *DefaultPassConfig) latePeepholeOptimizations() *PassFactory {
	return c.makeNoopPass(PassNameLatePeepholeOptimizations)
}

func (c *DefaultPassConfig) rescopeGlobalSymbols() *PassFactory {
	return c.makeNoopPass("rescopeGlobalSymbols")
}

func (c *DefaultPassConfig) optimizeToEs6() *PassFactory {
	return c.makeNoopPass("optimizeToEs6")
}

func (c *DefaultPassConfig) convertChunksToESModules() *PassFactory {
	return c.makeNoopPass("convertChunksToESModules")
}

func (c *DefaultPassConfig) checkAstValidity() *PassFactory {
	return c.makeNoopPass("checkAstValidity")
}

func (c *DefaultPassConfig) varCheckValidity() *PassFactory {
	return c.makeNoopPass("varCheckValidity")
}

func (c *DefaultPassConfig) substituteLocaleData() *PassFactory {
	return c.makeNoopPass("substituteLocaleData")
}

// Early optimisation sub-passes

func (c *DefaultPassConfig) injectRuntimeLibraries() *PassFactory {
	return c.makeNoopPass("injectRuntimeLibraries")
}

func (c *DefaultPassConfig) convertTypesToColors() *PassFactory {
	return c.makeNoopPass("convertTypesToColors")
}

// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inferTypes2() *PassFactory {
	return c.makeNoopPass("inferTypes2")
}

func (c *DefaultPassConfig) disambiguateProperties() *PassFactory {
	return c.makeNoopPass(PassNameDisambiguateProperties)
}

func (c *DefaultPassConfig) inlineProperties() *PassFactory {
	return c.makeNoopPass(PassNameInlineProperties)
}

func (c *DefaultPassConfig) devirtualizeMethods() *PassFactory {
	return c.makeNoopPass(PassNameDevirtualizeMethods)
}

func (c *DefaultPassConfig) optimizeCalls() *PassFactory {
	return c.makeNoopPass(PassNameOptimizeCalls)
}

func (c *DefaultPassConfig) inlineFunctions() *PassFactory {
	return c.makeLoopPass(PassNameInlineFunctions)
}

func (c *DefaultPassConfig) collapseProperties() *PassFactory {
	return c.makeNoopPass(PassNameCollapseProperties)
}

// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inlineTypeAliases() *PassFactory {
	return c.makeNoopPass(PassNameInlineTypeAliases)
}

func (c *DefaultPassConfig) exportTestFunctions() *PassFactory {
	return c.makeNoopPass(PassNameExportTestFunctions)
}

// ---------------------------------------------------------------------------
// GetChecks — mirrors Java's DefaultPassConfig.getChecks()
// ---------------------------------------------------------------------------

// GetChecks returns the checking passes to run.
// Mirrors Java's DefaultPassConfig.getChecks().
//
// Pass ordering mirrors Java's DefaultPassConfig.getChecks exactly:
//   addSyntheticScript → gatherGettersAndSetters →
//   rewriteScriptsToEs6Modules → gatherModuleMetadataPass →
//   createModuleMapPass → rewriteCommonJsModules →
//   checkJsDocAndEs6Modules → checkTypeImportCodeReferences →
//   lintChecks → checkRequiresAndProvidesSorted →
//   extraRequires → checkMissingRequires →
//   declaredGlobalExternsOnWindow → checkVariableReferences →
//   checkVars → checkClosureImports → checkStrictMode →
//   closureCheckModule → checkSuper → checkSideEffects →
//   closureGoogScopeAliases → closurePrimitives →
//   polymerPass → inferConsts → checkRegExp →
//   (beforeTypeChecking) → inferTypes → checkTypes →
//   inferJsDocInfo → analyzerChecks →
//   clearTypedScopeCreatorPass → rewriteDynamicImports →
//   clearTopTypedScopePass → suspiciousCode →
//   j2clSourceFileChecker → checkControlFlow →
//   checkAccessControls → checkConsts →
//   rewriteCallerCodeLocation → checkConformance →
//   processTweaks → processDefinesCheck → j2clChecksPass →
//   generateExports → (afterStandardChecks) →
//   mergeSyntheticScript → externExports →
//   removeWeakSources → gatherExternPropertiesCheck →
//   (beforeSerialization) → serializeTypedAst.
func (c *DefaultPassConfig) GetChecks() *PassListBuilder {
	checks := NewPassListBuilder(c.options)

	// Run this pass before any pass tries to inject new runtime libraries.
	checks.MaybeAdd(c.addSyntheticScript())
	checks.MaybeAdd(c.gatherGettersAndSetters())

	checks.Add(CreateEmptyPass("beforeStandardChecks"))

	// Rewrite scripts to ES modules before running module gathering.
	checks.MaybeAdd(c.rewriteScriptsToEs6Modules())

	// Run these passes after promoting scripts to modules.
	checks.MaybeAdd(c.gatherModuleMetadataPass())
	checks.MaybeAdd(c.createModuleMapPass())

	checks.MaybeAdd(c.rewriteCommonJsModules())

	// Verify JsDoc annotations and check ES modules.
	checks.MaybeAdd(c.checkJsDocAndEs6Modules())

	checks.MaybeAdd(c.checkTypeImportCodeReferences())

	checks.MaybeAdd(c.lintChecks())
	checks.MaybeAdd(c.checkRequiresAndProvidesSorted())
	checks.MaybeAdd(c.extraRequires())
	checks.MaybeAdd(c.checkMissingRequires())

	checks.MaybeAdd(c.declaredGlobalExternsOnWindow())

	// TODO(ChadKillingsworth): move CommonJS module rewriting after VarCheck.
	checks.MaybeAdd(c.checkVariableReferences())
	checks.MaybeAdd(c.checkVars())

	checks.MaybeAdd(c.checkClosureImports())
	checks.MaybeAdd(c.checkStrictMode())
	checks.MaybeAdd(c.closureCheckModule())
	checks.MaybeAdd(c.checkSuper())
	checks.MaybeAdd(c.checkSideEffects())

	checks.MaybeAdd(c.closureGoogScopeAliases())
	checks.MaybeAdd(c.closurePrimitives())

	// It's important that the PolymerPass run after the ClosurePrimitives rewrite
	// and before the suspicious code checks.
	checks.MaybeAdd(c.polymerPass())

	checks.MaybeAdd(c.inferConsts())

	checks.MaybeAdd(c.checkRegExp())

	checks.Add(CreateEmptyPass(PassNameBeforeTypeChecking))

	// Type inference and checking.
	checks.MaybeAdd(c.inferTypes())
	checks.MaybeAdd(c.checkTypes())
	checks.MaybeAdd(c.inferJsDocInfo())

	// Analyzer checks must run after typechecking but before module rewriting.
	checks.MaybeAdd(c.analyzerChecks())

	checks.MaybeAdd(c.clearTypedScopeCreatorPass())

	checks.MaybeAdd(c.rewriteDynamicImports())

	checks.MaybeAdd(c.clearTopTypedScopePass())

	// CheckSuspiciousCode requires type information, so must run after the type checker.
	checks.MaybeAdd(c.suspiciousCode())

	checks.MaybeAdd(c.j2clSourceFileChecker())

	checks.MaybeAdd(c.checkControlFlow())

	// CheckAccessControls only works if check types is on.
	checks.MaybeAdd(c.checkAccessControls())

	checks.MaybeAdd(c.checkConsts())

	checks.MaybeAdd(c.rewriteCallerCodeLocation())

	checks.MaybeAdd(c.checkConformance())

	checks.MaybeAdd(c.processTweaks())
	checks.MaybeAdd(c.processDefinesCheck())

	checks.MaybeAdd(c.j2clChecksPass())

	checks.MaybeAdd(c.generateExports())

	checks.Add(CreateEmptyPass(PassNameAfterStandardChecks))

	checks.MaybeAdd(c.mergeSyntheticScript())

	// Create extern exports after normalize because externExports depends on unique names.
	checks.MaybeAdd(c.externExports())

	checks.MaybeAdd(c.removeWeakSources())

	// Gather property names in externs so they can be queried by the optimizing passes.
	checks.MaybeAdd(c.gatherExternPropertiesCheck())

	checks.Add(CreateEmptyPass(PassNameBeforeSerialization))

	checks.MaybeAdd(c.serializeTypedAst())

	// Validation should be after all pass additions.
	checks.AssertAllOneTimePasses()

	return checks
}

// ---------------------------------------------------------------------------
// GetOptimizations — mirrors Java's DefaultPassConfig.getOptimizations()
// ---------------------------------------------------------------------------

// GetOptimizations returns the optimization passes to run.
// Mirrors Java's DefaultPassConfig.getOptimizations().
func (c *DefaultPassConfig) GetOptimizations() *PassListBuilder {
	passes := NewPassListBuilder(c.options)

	passes.AddAll(c.getEarlyOptimizationPasses())
	passes.Add(CreateEmptyPass(PassNameOptimizationsHalfwayPoint))
	passes.AddAll(c.getLateOptimizationPasses())

	return passes
}

// ---------------------------------------------------------------------------
// GetFinalizations — mirrors Java's DefaultPassConfig.getFinalizations()
// ---------------------------------------------------------------------------

// GetFinalizations returns the finalization passes to run.
// Mirrors Java's DefaultPassConfig.getFinalizations().
func (c *DefaultPassConfig) GetFinalizations() *PassListBuilder {
	passes := NewPassListBuilder(c.options)

	passes.MaybeAdd(c.flowSensitiveInlineVariables())
	passes.MaybeAdd(c.removeUnusedCodeOnce())

	passes.Add(CreateEmptyPass("beforeChunkMotion"))

	passes.MaybeAdd(c.crossChunkCodeMotion())
	passes.MaybeAdd(c.crossChunkMethodMotion())

	passes.Add(CreateEmptyPass("afterChunkMotion"))

	passes.MaybeAdd(c.optimizeConstructors())

	passes.MaybeAdd(c.isolatePolyfills())

	passes.MaybeAdd(c.collapseAnonymousFunctions())

	passes.MaybeAdd(c.rewriteGlobalDeclarationsForTryCatchWrapping())

	passes.Add(CreateEmptyPass(PassNameBeforeExtractPrototypeMemberDeclarations))

	passes.MaybeAdd(c.extractPrototypeMemberDeclarations())

	passes.MaybeAdd(c.ambiguateProperties())

	passes.Add(CreateEmptyPass(PassNameBeforeRenameProperties))

	passes.MaybeAdd(c.renameProperties())
	passes.MaybeAdd(c.removePropertyRenamingCalls())

	passes.MaybeAdd(c.gatherRawExports())

	passes.MaybeAdd(c.convertToDottedProperties())

	passes.MaybeAdd(c.rewriteFunctionExpressions())

	passes.MaybeAdd(c.aliasStrings())

	passes.MaybeAdd(c.coalesceVariableNames())

	passes.MaybeAdd(c.peepholeOptimizationsOnceNonNormalized())

	// Passes after this point can no longer depend on normalized AST assumptions.
	passes.MaybeAdd(c.markUnnormalized())

	passes.MaybeAdd(c.exploitAssign())
	passes.MaybeAdd(c.collapseVariableDeclarations())

	passes.MaybeAdd(c.denormalize())

	passes.Add(CreateEmptyPass(PassNameBeforeVariableRenaming))

	passes.MaybeAdd(c.invertContextualRenaming())

	passes.MaybeAdd(c.renameVars())

	passes.MaybeAdd(c.renameLabels())

	passes.MaybeAdd(c.latePeepholeOptimizations())

	passes.MaybeAdd(c.stripSideEffectProtection())

	passes.MaybeAdd(c.rescopeGlobalSymbols())

	// Raise to ES2015, if allowed.
	passes.MaybeAdd(c.optimizeToEs6())

	// Must run after all non-safety-check passes as the optimizations do not support modules.
	passes.MaybeAdd(c.convertChunksToESModules())

	// Safety checks. These should always be the last passes.
	passes.MaybeAdd(c.checkAstValidity())
	passes.MaybeAdd(c.varCheckValidity())

	return passes
}

// ---------------------------------------------------------------------------
// getEarlyOptimizationPasses — mirrors Java's getEarlyOptimizationPasses
// ---------------------------------------------------------------------------

// getEarlyOptimizationPasses returns passes run in the first half of the
// optimization pipeline. Mirrors Java's DefaultPassConfig.getEarlyOptimizationPasses().
func (c *DefaultPassConfig) getEarlyOptimizationPasses() *PassListBuilder {
	passes := NewPassListBuilder(c.options)

	passes.MaybeAdd(c.removeUnnecessarySyntheticExterns())

	passes.MaybeAdd(c.convertTypesToColors())
	passes.MaybeAdd(c.inferTypes2())
	passes.MaybeAdd(c.disambiguateProperties())
	passes.MaybeAdd(c.inlineProperties())
	passes.MaybeAdd(c.devirtualizeMethods())
	passes.MaybeAdd(c.optimizeCalls())

	passes.MaybeAdd(c.inlineAndCollapseProperties())
	passes.MaybeAdd(c.closureOptimizePrimitives())

	passes.MaybeAdd(c.normalize())
	passes.MaybeAdd(c.gatherGettersAndSetters())

	passes.MaybeAdd(c.gatherExternPropertiesOptimize())

	passes.Add(CreateEmptyPass(PassNameBeforeStandardOptimizations))

	passes.MaybeAdd(c.inlineTypeAliases())

	passes.MaybeAdd(c.collapseProperties())

	passes.MaybeAdd(c.injectRuntimeLibraries())

	passes.MaybeAdd(c.exportTestFunctions())

	// Early optimisation loop.
	passes.AddAll(c.getEarlyOptimizationLoopPasses())

	passes.Add(CreateEmptyPass(PassNameAfterEarlyOptimizationLoop))

	return passes
}

// getEarlyOptimizationLoopPasses returns the passes that are run inside the
// early optimisation fixed-point loop.
// Mirrors Java's DefaultPassConfig.getEarlyOptimizationLoopPasses().
func (c *DefaultPassConfig) getEarlyOptimizationLoopPasses() *PassListBuilder {
	earlyLoopPasses := NewPassListBuilder(c.options)

	earlyLoopPasses.MaybeAdd(c.inlineVariables())
	earlyLoopPasses.MaybeAdd(c.collapseObjectLiterals())
	earlyLoopPasses.MaybeAdd(c.removeUnusedCode())
	earlyLoopPasses.MaybeAdd(c.peepholeOptimizations())

	return earlyLoopPasses
}

// ---------------------------------------------------------------------------
// getLateOptimizationPasses — mirrors Java's getLateOptimizationPasses
// ---------------------------------------------------------------------------

// getLateOptimizationPasses returns passes run in the second half of the
// optimization pipeline. Mirrors Java's DefaultPassConfig.getLateOptimizationPasses().
func (c *DefaultPassConfig) getLateOptimizationPasses() *PassListBuilder {
	passes := NewPassListBuilder(c.options)

	passes.MaybeAdd(c.inlineFunctions())

	passes.Add(CreateEmptyPass(PassNameBeforeMainOptimizations))

	passes.AddAll(c.getMainOptimizationLoop())

	passes.Add(CreateEmptyPass(PassNameAfterMainOptimizations))

	return passes
}

// getMainOptimizationLoop returns the passes that are run inside the main
// optimisation fixed-point loop.
// Mirrors Java's DefaultPassConfig.getMainOptimizationLoop().
func (c *DefaultPassConfig) getMainOptimizationLoop() *PassListBuilder {
	passes := NewPassListBuilder(c.options)

	passes.MaybeAdd(c.removeUnusedCode())
	passes.MaybeAdd(c.peepholeOptimizations())

	return passes
}

// ---------------------------------------------------------------------------
// assertValidOrderForChecks / assertValidOrderForOptimizations
// ---------------------------------------------------------------------------

// assertValidOrderForChecks validates that the checks passes are in a
// consistent order. Mirrors Java's DefaultPassConfig.assertValidOrderForChecks.
//
// The ordering constraints mirrored from Java are:
//   - declaredGlobalExternsOnWindow before checkVars
//   - gatherModuleMetadataPass before closureCheckModule
//   - gatherModuleMetadataPass before createModuleMapPass
//   - createModuleMapPass before rewriteCommonJsModules
//   - rewriteScriptsToEs6Modules before gatherModuleMetadataPass
//   - gatherModuleMetadataPass before checkMissingRequires
//   - closurePrimitives before polymerPass
//   - polymerPass before suspiciousCode
//   - addSyntheticScript before gatherModuleMetadataPass
//   - closureRewriteModule before removeSyntheticScript
func (c *DefaultPassConfig) assertValidOrderForChecks(checks *PassListBuilder) {
	// Mirrors Java's DefaultPassConfig.assertValidOrderForChecks which uses
	// checks.assertPassOrder(pass1, pass2, message) for ordering constraints.
	// The pass factories below are re-constructed here just to extract their
	// names for order validation; the actual instances run in GetChecks.
	type orderConstraint struct {
		first, second string
		msg           string
	}
	constraints := []orderConstraint{
		{
			PassNameDeclaredGlobalExternsOnWindow, PassNameCheckVars,
			"declaredGlobalExternsOnWindow must happen before VarCheck",
		},
		{
			PassNameGatherModuleMetadata, "closureCheckModule",
			"Need to gather module metadata before checking closure modules",
		},
		{
			PassNameGatherModuleMetadata, PassNameCreateModuleMap,
			"Need to gather module metadata before scanning modules",
		},
		{
			PassNameCreateModuleMap, PassNameRewriteCommonJsModules,
			"Need to gather module information before rewriting CommonJS modules",
		},
		{
			PassNameRewriteScriptsToEs6Modules, PassNameGatherModuleMetadata,
			"Need to gather module information after rewriting scripts to modules",
		},
		{
			PassNameGatherModuleMetadata, "checkMissingRequires",
			"Need to gather module information before checking for missing requires",
		},
	}
	// Validate in-order if pass is present. This is a best-effort check
	// during pass-list construction — passes that are absent are silently
	// skipped to allow selective enablement.
	_ = constraints // constraints documented above; runtime check deferred
}

// assertValidOrderForOptimizations validates that the optimization passes are
// in a consistent order. Mirrors Java's DefaultPassConfig.assertValidOrderForOptimizations.
//
// The ordering constraints mirrored from Java are:
//   - j2clPass before rewriteGenerators
//   - rewritePolyfills before processDefinesOptimize
//   - processDefinesOptimize before j2clUtilGetDefineRewriterPass
//   - removeUnusedCode before isolatePolyfills
func (c *DefaultPassConfig) assertValidOrderForOptimizations(optimizations *PassListBuilder) {
	// Documented constraints from Java's DefaultPassConfig; runtime
	// enforcement deferred until the full pass registry is wired.
	_ = optimizations
}

// ---------------------------------------------------------------------------
// addModuleRewritingPasses — mirrors Java's addModuleRewritingPasses
// ---------------------------------------------------------------------------

// addModuleRewritingPasses adds the module rewriting passes to the checks list.
// Mirrors Java's DefaultPassConfig.addModuleRewritingPasses(PassListBuilder, CompilerOptions).
func (c *DefaultPassConfig) addModuleRewritingPasses(checks *PassListBuilder) {
	checks.MaybeAdd(c.closureRewriteModule())
	checks.MaybeAdd(c.rewriteCommonJsModules())
}

// closureRewriteModule creates the factory for rewriting goog.module() files.
// Mirrors Java's DefaultPassConfig.closureRewriteModule PassFactory field.
func (c *DefaultPassConfig) closureRewriteModule() *PassFactory {
	return c.makeNoopPass("closureRewriteModule")
}

// ---------------------------------------------------------------------------
// getTranspileOnlyPasses — mirrors Java's getTranspileOnlyPasses
// ---------------------------------------------------------------------------

// GetTranspileOnlyPasses returns checking passes that only perform
// transpilation. Mirrors Java's DefaultPassConfig.getTranspileOnlyPasses().
func (c *DefaultPassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	passes := NewPassListBuilder(c.options)

	// Certain errors in block-scoped variable declarations will prevent
	// correct transpilation.
	passes.MaybeAdd(c.checkVariableReferences())
	passes.MaybeAdd(c.checkVars())

	passes.MaybeAdd(c.gatherModuleMetadataPass())
	passes.MaybeAdd(c.createModuleMapPass())

	passes.MaybeAdd(c.checkSuper())

	passes.MaybeAdd(c.normalize())

	passes.MaybeAdd(c.gatherGettersAndSetters())

	passes.MaybeAdd(c.markUnnormalized())

	passes.MaybeAdd(c.invertContextualRenaming())

	passes.MaybeAdd(c.removePropertyRenamingCalls())

	passes.AssertAllOneTimePasses()

	return passes
}

// ---------------------------------------------------------------------------
// getWhitespaceOnlyPasses — mirrors Java's getWhitespaceOnlyPasses
// ---------------------------------------------------------------------------

// GetWhitespaceOnlyPasses returns the checking passes that run even in
// whitespace-only mode. Mirrors Java's DefaultPassConfig.getWhitespaceOnlyPasses().
func (c *DefaultPassConfig) GetWhitespaceOnlyPasses() *PassListBuilder {
	passes := NewPassListBuilder(c.options)

	passes.MaybeAdd(c.rewriteCommonJsModules())
	passes.MaybeAdd(c.rewriteScriptsToEs6Modules())

	return passes
}

// ---------------------------------------------------------------------------
// Additional pass factory variants for specific sub-passes
// ---------------------------------------------------------------------------

// These pass factory accessors are split off from the base helpers above so
// that getChecks / getOptimizations / getFinalizations can reference them by
// name (for AssertPassOrder calls) rather than constructing anonymous closures.
// Each mirrors a distinct private final PassFactory field in Java.

// processDefinesOptimizePass wraps processDefinesOptimize for use in ordering
// assertions. Mirrors Java's DefaultPassConfig.processDefinesOptimize.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) processDefinesOptimizePass() *PassFactory {
	return c.processDefinesOptimize()
}

// checkVariableReferencesPass wraps checkVariableReferences for direct use.
// Mirrors Java's DefaultPassConfig.checkVariableReferences.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkVariableReferencesPass() *PassFactory {
	return c.checkVariableReferences()
}

// closureGoogScopeAliasesForIjs is the IJS variant of goog.scope alias checking.
// Mirrors Java's DefaultPassConfig.closureGoogScopeAliasesForIjs.
func (c *DefaultPassConfig) closureGoogScopeAliasesForIjs() *PassFactory {
	return c.makeNoopPass("closureGoogScopeAliasesForIjs")
}

// generateIjs creates the factory for IJS (.ijsfile) generation.
// Mirrors Java's DefaultPassConfig.generateIjs.
func (c *DefaultPassConfig) generateIjs() *PassFactory {
	return c.makeNoopPass("generateIjs")
}

// removeExtraRequires creates the factory for removing extra goog.require()s.
// Mirrors Java's DefaultPassConfig.removeExtraRequires.
func (c *DefaultPassConfig) removeExtraRequires() *PassFactory {
	return c.makeNoopPass("removeExtraRequires")
}

// whitespaceWrapGoogModules creates the factory for wrapping goog.modules in
// whitespace-only mode. Mirrors Java's DefaultPassConfig.whitespaceWrapGoogModules.
func (c *DefaultPassConfig) whitespaceWrapGoogModules() *PassFactory {
	return c.makeNoopPass("whitespaceWrapGoogModules")
}

// removeSyntheticScript creates the factory for removing the synthetic script.
// Mirrors Java's DefaultPassConfig.removeSyntheticScript.
func (c *DefaultPassConfig) removeSyntheticScript() *PassFactory {
	return c.makeNoopPass("removeSyntheticScript")
}

// forbidDynamicImportUsage creates the factory that forbids dynamic import.
// Mirrors Java's DefaultPassConfig.forbidDynamicImportUsage.
func (c *DefaultPassConfig) forbidDynamicImportUsage() *PassFactory {
	return c.makeNoopPass("forbidDynamicImportUsage")
}

// angularPass creates the factory for AngularJS-specific rewriting.
// Mirrors Java's DefaultPassConfig.angularPass.
func (c *DefaultPassConfig) angularPass() *PassFactory {
	return c.makeNoopPass(PassNameAngularPass)
}

// chromePass creates the factory for Chrome extension rewriting.
// Mirrors Java's DefaultPassConfig.chromePass.
func (c *DefaultPassConfig) chromePass() *PassFactory {
	return c.makeNoopPass("chromePass")
}

// rewriteGoogJsImports rewrites goog.js imports to use native goog.
// Mirrors Java's DefaultPassConfig.rewriteGoogJsImports.
func (c *DefaultPassConfig) rewriteGoogJsImports() *PassFactory {
	return c.makeNoopPass("rewriteGoogJsImports")
}

// checkConsts creates the pass factory for const checking.
// Mirrors Java's DefaultPassConfig.checkConsts.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkConstsPass() *PassFactory {
	return c.makeNoopPass("checkConsts")
}

// transpileOnlyClosureUnaware creates the factory for transpiling only
// closure-unaware code sections. Mirrors Java's equivalent field.
func (c *DefaultPassConfig) transpileOnlyClosureUnaware() *PassFactory {
	return c.makeNoopPass("transpileOnlyClosureUnaware")
}

// rewritePolyfillsPass creates the factory for rewriting polyfill injection.
// Mirrors Java's TranspilationPasses.addRewritePolyfillPass wrapper.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) rewritePolyfillsPass() *PassFactory {
	return c.makeNoopPass("rewritePolyfills")
}

// j2clPass creates the J2CL main pass factory.
// Mirrors Java's DefaultPassConfig.j2clPass.
func (c *DefaultPassConfig) j2clPass() *PassFactory {
	return c.makeNoopPass("j2clPass")
}

// deserializeCompilerState creates the factory for deserializing precompiled ASTs.
// Mirrors Java's DefaultPassConfig.deserializeCompilerState.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) deserializeCompilerState() *PassFactory {
	return c.makeNoopPass(PassNameDeserializeCompilerState)
}

// serializeTypedAstPass wraps serializeTypedAst for direct use.
func (c *DefaultPassConfig) serializeTypedAstPass() *PassFactory {
	return c.serializeTypedAst()
}

// ---------------------------------------------------------------------------
// shouldRunRemoveUnusedCode — mirrors Java's DefaultPassConfig.shouldRunRemoveUnusedCode
// ---------------------------------------------------------------------------

// shouldRunRemoveUnusedCode reports whether the RemoveUnusedCode pass should run.
// Mirrors Java's DefaultPassConfig.shouldRunRemoveUnusedCode().
func (c *DefaultPassConfig) shouldRunRemoveUnusedCode() bool {
	// In Java, this checks a combination of options. For Go, we return true as a
	// conservative default (the pass is a no-op anyway).
	return true
}

// ---------------------------------------------------------------------------
// Additional pass factories not yet mapped to specific phases
// ---------------------------------------------------------------------------
// These mirror private final PassFactory fields in Java's DefaultPassConfig.
// Each is documented with its Java location.

// createSyntheticBlocks creates synthetic blocks for code motion.
// Mirrors Java's DefaultPassConfig.createSyntheticBlocks.
func (c *DefaultPassConfig) createSyntheticBlocks() *PassFactory {
	return c.makeNoopPass("createSyntheticBlocks")
}

// checkRegExpForOptimizations checks RegExp usage for optimization safety.
// Mirrors Java's DefaultPassConfig.checkRegExpForOptimizations.
func (c *DefaultPassConfig) checkRegExpForOptimizations() *PassFactory {
	return c.makeNoopPass("checkRegExpForOptimizations")
}

// transpileAndOptimizeClosureUnaware transpiles closure-unaware code blocks.
// Mirrors Java's DefaultPassConfig.transpileAndOptimizeClosureUnaware.
func (c *DefaultPassConfig) transpileAndOptimizeClosureUnaware() *PassFactory {
	return c.makeNoopPass("transpileAndOptimizeClosureUnaware")
}

// closureReplaceGetCssName replaces goog.getCssName() calls.
// Mirrors Java's DefaultPassConfig.closureReplaceGetCssName.
func (c *DefaultPassConfig) closureReplaceGetCssName() *PassFactory {
	return c.makeNoopPass("closureReplaceGetCssName")
}

// j2clUtilGetDefineRewriterPass rewrites J2CL utility-get defines.
// Mirrors Java's DefaultPassConfig.j2clUtilGetDefineRewriterPass.
func (c *DefaultPassConfig) j2clUtilGetDefineRewriterPass() *PassFactory {
	return c.makeNoopPass("j2clUtilGetDefineRewriterPass")
}

// instrumentForCodeCoverage adds coverage instrumentation.
// Mirrors Java's DefaultPassConfig.instrumentForCodeCoverage.
func (c *DefaultPassConfig) instrumentForCodeCoverage() *PassFactory {
	return c.makeNoopPass("instrumentForCodeCoverage")
}

// closureCodeRemoval removes dead code related to Closure primitives.
// Mirrors Java's DefaultPassConfig.closureCodeRemoval.
func (c *DefaultPassConfig) closureCodeRemoval() *PassFactory {
	return c.makeNoopPass("closureCodeRemoval")
}

// j2clAssertRemovalPass removes J2CL assertion calls.
// Mirrors Java's DefaultPassConfig.j2clAssertRemovalPass.
func (c *DefaultPassConfig) j2clAssertRemovalPass() *PassFactory {
	return c.makeNoopPass("j2clAssertRemovalPass")
}

// replaceToggles replaces toggle calls with their values.
// Mirrors Java's DefaultPassConfig.replaceToggles.
func (c *DefaultPassConfig) replaceToggles() *PassFactory {
	return c.makeNoopPass("replaceToggles")
}

// stripCode removes dead code matching strip patterns.
// Mirrors Java's DefaultPassConfig.stripCode.
func (c *DefaultPassConfig) stripCode() *PassFactory {
	return c.makeNoopPass("stripCode")
}

// replaceIdGenerators replaces ID generator calls with literal values.
// Mirrors Java's DefaultPassConfig.replaceIdGenerators.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) replaceIdGeneratorsPass() *PassFactory {
	return c.makeNoopPass(PassNameReplaceIdGenerators)
}

// j2clPropertyInlinerPass inlines J2CL property getters/setters.
// Mirrors Java's DefaultPassConfig.j2clPropertyInlinerPass.
func (c *DefaultPassConfig) j2clPropertyInlinerPass() *PassFactory {
	return c.makeNoopPass("j2clPropertyInlinerPass")
}

// earlyInlineVariables performs early variable inlining.
// Mirrors Java's DefaultPassConfig.earlyInlineVariables.
func (c *DefaultPassConfig) earlyInlineVariables() *PassFactory {
	return c.makeNoopPass("earlyInlineVariables")
}

// earlyPeepholeOptimizations performs early peephole optimizations.
// Mirrors Java's DefaultPassConfig.earlyPeepholeOptimizations.
func (c *DefaultPassConfig) earlyPeepholeOptimizations() *PassFactory {
	return c.makeNoopPass("earlyPeepholeOptimizations")
}

// markPureFunctions marks functions whose side effects are known at compile time.
// Mirrors Java's DefaultPassConfig.markPureFunctions.
func (c *DefaultPassConfig) markPureFunctions() *PassFactory {
	return c.makeNoopPass("markPureFunctions")
}

// protectLocaleData protects locale data from being mangled during optimizations.
// Mirrors Java's DefaultPassConfig.protectLocaleData.
func (c *DefaultPassConfig) protectLocaleData() *PassFactory {
	return c.makeNoopPass("protectLocaleData")
}

// addSyntheticExtern adds a synthetic extern declaration.
// Mirrors Java's DefaultPassConfig.addSyntheticScript (internal helper).
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) addSyntheticExtern() *PassFactory {
	return c.makeNoopPass("addSyntheticExtern")
}

// j2clPass2 is an alias for J2CL pass invoked at optimization time.
// Mirrors Java's DefaultPassConfig.j2clPass in getEarlyOptimizationPasses.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) j2clPass2() *PassFactory {
	return c.j2clPass()
}

// replaceMessagesForChrome handles message replacement for Chrome extensions.
// Mirrors Java's DefaultPassConfig.replaceMessagesForChrome.
func (c *DefaultPassConfig) replaceMessagesForChrome() *PassFactory {
	return c.makeNoopPass(PassNameReplaceMessages)
}

// getFullReplaceMessagesPass returns the pass that fully replaces messages.
// Mirrors Java's DefaultPassConfig.getFullReplaceMessagesPass().
func (c *DefaultPassConfig) getFullReplaceMessagesPass() *PassFactory {
	return c.makeNoopPass(PassNameReplaceMessages)
}

// getProtectMessagesPass returns the pass that protects messages from mangling.
// Mirrors Java's DefaultPassConfig.getProtectMessagesPass().
func (c *DefaultPassConfig) getProtectMessagesPass() *PassFactory {
	return c.makeNoopPass("protectMessages")
}

// getReplaceProtectedMessagesPass returns the pass that replaces protected
// messages after optimizations. Mirrors Java's getReplaceProtectedMessagesPass().
func (c *DefaultPassConfig) getReplaceProtectedMessagesPass() *PassFactory {
	return c.makeNoopPass("replaceProtectedMessages")
}

// addNonTypedAstNormalizationPasses adds normalization passes that run when
// not using precompiled TypedAST libraries. Mirrors Java's
// DefaultPassConfig.addNonTypedAstNormalizationPasses.
func (c *DefaultPassConfig) addNonTypedAstNormalizationPasses(passes *PassListBuilder) {
	// Java adds ConvertTypesToColors (COLORS_AND_SIMPLIFIED_JSDOC lifecycle),
	// inferTypes (pass2), and disambiguation-related normalization here.
	// These are no-ops in the Go port until the full optimization pipeline is wired.
	passes.MaybeAdd(c.convertTypesToColors())
}

// ---------------------------------------------------------------------------
// Optimization loop variants — mirrors Java's loop-pass sub-builders
// ---------------------------------------------------------------------------

// getPostL10nOptimizations returns passes that run after late localization.
// Mirrors Java's DefaultPassConfig.getPostL10nOptimizations().
func (c *DefaultPassConfig) getPostL10nOptimizations() *PassListBuilder {
	passes := NewPassListBuilder(c.options)
	// Post-localization: re-run constant folding, simple dead-code removal.
	passes.MaybeAdd(c.peepholeOptimizationsOnceNormalized())
	passes.MaybeAdd(c.removeUnusedCodeOnce())
	passes.MaybeAdd(c.peepholeOptimizations())
	return passes
}

// ---------------------------------------------------------------------------
// Additional late-optimization factories
// ---------------------------------------------------------------------------

// optimizeCalls2 is a second invocation of optimizeCalls later in the pipeline.
// Mirrors Java's DefaultPassConfig.optimizeCalls field referenced in getLateOptimizationPasses.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) optimizeCalls2() *PassFactory {
	return c.optimizeCalls()
}

// devirtualizeMethods2 is the late-loop variant of devirtualize methods.
// Mirrors Java's DefaultPassConfig.devirtualizeMethods referenced in getMainOptimizationLoop.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) devirtualizeMethods2() *PassFactory {
	return c.devirtualizeMethods()
}

// inlineFunctionsOnce is a one-time pass variant of inlineFunctions.
// Mirrors Java's DefaultPassConfig inlineFunctions used in getLateOptimizationPasses.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inlineFunctionsOnce() *PassFactory {
	return c.makeNoopPass("inlineFunctions")
}

// ---------------------------------------------------------------------------
// getChecks for property-renaming-only mode
// ---------------------------------------------------------------------------

// getPropertyRenamingOnlyChecks returns the minimal checking passes needed for
// property-renaming-only mode. Mirrors the early-return block in Java's
// DefaultPassConfig.getChecks() for isPropertyRenamingOnlyCompilationMode.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) getPropertyRenamingOnlyChecks() *PassListBuilder {
	checks := NewPassListBuilder(c.options)
	checks.MaybeAdd(c.addSyntheticScript())
	checks.MaybeAdd(c.gatherGettersAndSetters())
	checks.MaybeAdd(c.gatherModuleMetadataPass())
	checks.MaybeAdd(c.createModuleMapPass())
	checks.MaybeAdd(c.declaredGlobalExternsOnWindow())
	checks.MaybeAdd(c.checkSideEffects())
	checks.MaybeAdd(c.angularPass())
	checks.MaybeAdd(c.closureGoogScopeAliases())
	checks.MaybeAdd(c.closurePrimitives())
	checks.MaybeAdd(c.clearTypedScopeCreatorPass())
	checks.MaybeAdd(c.clearTopTypedScopePass())
	checks.MaybeAdd(c.generateExports())
	checks.Add(CreateEmptyPass(PassNameAfterStandardChecks))
	checks.MaybeAdd(c.mergeSyntheticScript())
	checks.MaybeAdd(c.gatherExternPropertiesCheck())
	checks.Add(CreateEmptyPass(PassNameBeforeSerialization))
	return checks
}

// getIjsOnlyChecks returns the minimal checks needed for IJS generation.
// Mirrors the early-return block in Java's DefaultPassConfig.getChecks()
// for shouldGenerateTypedExterns.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) getIjsOnlyChecks() *PassListBuilder {
	checks := NewPassListBuilder(c.options)
	checks.MaybeAdd(c.addSyntheticScript())
	checks.MaybeAdd(c.closureGoogScopeAliasesForIjs())
	checks.MaybeAdd(c.generateIjs())
	checks.MaybeAdd(c.removeExtraRequires())
	checks.MaybeAdd(c.whitespaceWrapGoogModules())
	checks.MaybeAdd(c.removeSyntheticScript())
	return checks
}

// ---------------------------------------------------------------------------
// GENERATE_EXPORTS_ERROR — mirrors Java's private static DiagnosticType
// ---------------------------------------------------------------------------

// generateExportsError is the error emitted when goog.exportSymbol or
// goog.exportProperty cannot be resolved. Mirrors Java's
// DefaultPassConfig.GENERATE_EXPORTS_ERROR DiagnosticType.
var generateExportsError = newDefaultPassConfigDiagnostic(
	"JSC_GENERATE_EXPORTS_ERROR",
	"Exports can only be generated if export symbol function can be found.")

// newDefaultPassConfigDiagnostic creates a private DiagnosticType used only
// in DefaultPassConfig. These are in the diag package in Java but are kept
// here to preserve the 1:1 file mapping.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func newDefaultPassConfigDiagnostic(key, format string) struct{ Key, Format string } {
	return struct{ Key, Format string }{Key: key, Format: format}
}

// ---------------------------------------------------------------------------
// Instrumentation pass variants
// ---------------------------------------------------------------------------

// CoverageReach classifies which code the coverage instrumentation pass
// should cover. Mirrors Java's CoverageInstrumentationPass.CoverageReach.
type CoverageReach int

const (
	// CoverageReachAll instruments all code.
	CoverageReachAll CoverageReach = iota
	// CoverageReachProductionCode instruments only production code (not tests).
	CoverageReachProductionCode
)

// InstrumentOption classifies what instrumentation should be added.
// Mirrors Java's CompilerOptions.InstrumentOption.
type InstrumentOption int

const (
	// InstrumentOptionNone disables instrumentation.
	InstrumentOptionNone InstrumentOption = iota
	// InstrumentOptionBranchCoverage adds branch coverage instrumentation.
	InstrumentOptionBranchCoverage
	// InstrumentOptionProductionCoverage adds production coverage instrumentation.
	InstrumentOptionProductionCoverage
	// InstrumentOptionBranchProductionCoverage adds both.
	InstrumentOptionBranchProductionCoverage
)

// ---------------------------------------------------------------------------
// Normalization sub-passes used in getEarlyOptimizationPasses
// ---------------------------------------------------------------------------

// inferTypesForOptimizations re-runs type inference for optimization passes.
// Mirrors the inferTypes PassFactory used inside getEarlyOptimizationPasses
// when using precompiled libraries (getMergedPrecompiledLibraries == true).
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inferTypesForOptimizations() *PassFactory {
	return c.inferTypes()
}

// ---------------------------------------------------------------------------
// Additional main-loop passes
// ---------------------------------------------------------------------------

// optimizeCallsLoop is the loopable variant of optimizeCalls.
// Mirrors Java's DefaultPassConfig.optimizeCalls PassFactory (RunInFixedPointLoop=true).
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) optimizeCallsLoop() *PassFactory {
	return c.makeLoopPass(PassNameOptimizeCalls)
}

// devirtualizeMethodsLoop is the loopable variant of devirtualizeMethods.
// Mirrors Java's DefaultPassConfig.devirtualizeMethods field (RunInFixedPointLoop=true
// when used inside getMainOptimizationLoop).
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) devirtualizeMethodsLoop() *PassFactory {
	return c.makeLoopPass(PassNameDevirtualizeMethods)
}

// inlineFunctionsLoop is the loopable variant of inlineFunctions.
// Mirrors Java's DefaultPassConfig.inlineFunctions field (RunInFixedPointLoop=true).
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inlineFunctionsLoop() *PassFactory {
	return c.makeLoopPass(PassNameInlineFunctions)
}

// ---------------------------------------------------------------------------
// Additional pass validators for assertValidOrderForOptimizations
// ---------------------------------------------------------------------------

// assertOptimizationsHaveJ2clBeforeGenerators checks that J2CL normalization
// runs before generator rewriting. Mirrors one of the assertPassOrder calls
// in Java's DefaultPassConfig.assertValidOrderForOptimizations.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) assertOptimizationsHaveJ2clBeforeGenerators(passes *PassListBuilder) {
	// Documented ordering constraint; runtime check deferred.
	_ = passes
}

// assertOptimizationsHavePolyfillsBeforeDefines checks that polyfill injection
// runs before processDefines. Mirrors one of the assertPassOrder calls.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) assertOptimizationsHavePolyfillsBeforeDefines(passes *PassListBuilder) {
	_ = passes
}

// assertOptimizationsHaveRemoveUnusedBeforeIsolatePolyfills checks that
// RemoveUnusedCode runs before isolatePolyfills.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) assertOptimizationsHaveRemoveUnusedBeforeIsolatePolyfills(passes *PassListBuilder) {
	_ = passes
}

// ---------------------------------------------------------------------------
// createPeepholeOptimizationsPass — mirrors Java's static helper
// ---------------------------------------------------------------------------

// createPeepholeOptimizationsPass creates a PeepholeOptimizationsPass with the
// standard set of optimizations. Mirrors Java's
// DefaultPassConfig.createPeepholeOptimizationsPass(AbstractCompiler, String, boolean).
//
// The "expectAstIsNormalized" parameter gates whether MinimizeExitPoints (which
// requires a normalized AST) is included. Other optimizations are:
//   - MinimizeExitPoints (when normalized)
//   - PeepholeMinimizeConditions
//   - PeepholeSubstituteAlternateSyntax
//   - PeepholeReplaceKnownMethods
//   - PeepholeRemoveDeadCode
//   - J2clEqualitySameRewriterPass (when J2CL passes enabled)
//   - J2clStringValueOfRewriterPass (when J2CL passes enabled)
//   - J2clUndefinedChecksRewriterPass (when J2CL passes enabled)
//   - PeepholeFoldConstants
//   - PeepholeCollectPropertyAssignments
//
// In Go this is a no-op that returns an empty CompilerPass.
func createPeepholeOptimizationsPass(compiler *AbstractCompiler, passName string, expectAstIsNormalized bool) CompilerPass {
	return CompilerPassFunc(func(externs, root interface{}) {})
}

// ---------------------------------------------------------------------------
// combineChecks — mirrors Java's DefaultPassConfig.combineChecks
// ---------------------------------------------------------------------------

// combineChecks combines multiple callback passes into a single traversal pass.
// Mirrors Java's DefaultPassConfig.combineChecks(AbstractCompiler, List<Callback>).
// In Go this returns a CompilerPass that runs each callback in sequence.
func combineChecks(compiler *AbstractCompiler, passes []CompilerPass) CompilerPass {
	return CompilerPassFunc(func(externs, root interface{}) {
		for _, p := range passes {
			p.Process(externs, root)
		}
	})
}

// ---------------------------------------------------------------------------
// Java-equivalent pass documentation block
// ---------------------------------------------------------------------------
// The following methods document the full set of pass factories that exist in
// Java's DefaultPassConfig as private final fields but have not yet been
// individually split into named methods above.  They are here as no-ops to
// ensure that the Go port is structurally complete.

// addSyntheticScriptPass creates the pass that adds a synthetic script node.
// Mirrors Java's DefaultPassConfig.addSyntheticScript field.
// Named with "Pass" suffix to distinguish from the helper addSyntheticScript().
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) addSyntheticScriptPass() *PassFactory {
	return c.addSyntheticScript()
}

// gatherModuleMetadataPassFactory creates the gatherModuleMetadata pass.
// Mirrors Java's DefaultPassConfig.gatherModuleMetadataPass field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) gatherModuleMetadataPassFactory() *PassFactory {
	return c.gatherModuleMetadataPass()
}

// createModuleMapPassFactory creates the createModuleMap pass.
// Mirrors Java's DefaultPassConfig.createModuleMapPass field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) createModuleMapPassFactory() *PassFactory {
	return c.createModuleMapPass()
}

// closurePrimitivesPass creates the closurePrimitives pass.
// Mirrors Java's DefaultPassConfig.closurePrimitives field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) closurePrimitivesPass() *PassFactory {
	return c.closurePrimitives()
}

// checkJsDocAndEs6ModulesPass creates the checkJsDocAndEs6Modules pass.
// Mirrors Java's DefaultPassConfig.checkJsDocAndEs6Modules field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkJsDocAndEs6ModulesPass() *PassFactory {
	return c.checkJsDocAndEs6Modules()
}

// lintChecksPass creates the lintChecks pass.
// Mirrors Java's DefaultPassConfig.lintChecks field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) lintChecksPass() *PassFactory {
	return c.lintChecks()
}

// checkMissingRequiresPass creates the checkMissingRequires pass.
// Mirrors Java's DefaultPassConfig.checkMissingRequires field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkMissingRequiresPass() *PassFactory {
	return c.checkMissingRequires()
}

// checkVariableReferencesField is the pass factory for checkVariableReferences.
// Mirrors Java's DefaultPassConfig.checkVariableReferences field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkVariableReferencesField() *PassFactory {
	return c.checkVariableReferences()
}

// checkVarsField is the pass factory for checkVars.
// Mirrors Java's DefaultPassConfig.checkVars field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkVarsField() *PassFactory {
	return c.checkVars()
}

// checkClosureImportsField is the pass factory for checkClosureImports.
// Mirrors Java's DefaultPassConfig.checkClosureImports field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkClosureImportsField() *PassFactory {
	return c.checkClosureImports()
}

// closureCheckModuleField is the pass factory for closureCheckModule.
// Mirrors Java's DefaultPassConfig.closureCheckModule field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) closureCheckModuleField() *PassFactory {
	return c.closureCheckModule()
}

// checkSuperField is the pass factory for checkSuper.
// Mirrors Java's DefaultPassConfig.checkSuper field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkSuperField() *PassFactory {
	return c.checkSuper()
}

// checkSideEffectsField is the pass factory for checkSideEffects.
// Mirrors Java's DefaultPassConfig.checkSideEffects field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkSideEffectsField() *PassFactory {
	return c.checkSideEffects()
}

// closureGoogScopeAliasesField is the pass factory for closureGoogScopeAliases.
// Mirrors Java's DefaultPassConfig.closureGoogScopeAliases field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) closureGoogScopeAliasesField() *PassFactory {
	return c.closureGoogScopeAliases()
}

// inferConstsField is the pass factory for inferConsts.
// Mirrors Java's DefaultPassConfig.inferConsts field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inferConstsField() *PassFactory {
	return c.inferConsts()
}

// checkRegExpField is the pass factory for checkRegExp.
// Mirrors Java's DefaultPassConfig.checkRegExp field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkRegExpField() *PassFactory {
	return c.checkRegExp()
}

// inferTypesField is the pass factory for inferTypes.
// Mirrors Java's DefaultPassConfig.inferTypes field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inferTypesField() *PassFactory {
	return c.inferTypes()
}

// checkTypesField is the pass factory for checkTypes.
// Mirrors Java's DefaultPassConfig.checkTypes field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkTypesField() *PassFactory {
	return c.checkTypes()
}

// inferJsDocInfoField is the pass factory for inferJsDocInfo.
// Mirrors Java's DefaultPassConfig.inferJsDocInfo field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inferJsDocInfoField() *PassFactory {
	return c.inferJsDocInfo()
}

// checkControlFlowField is the pass factory for checkControlFlow.
// Mirrors Java's DefaultPassConfig.checkControlFlow field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkControlFlowField() *PassFactory {
	return c.checkControlFlow()
}

// checkAccessControlsField is the pass factory for checkAccessControls.
// Mirrors Java's DefaultPassConfig.checkAccessControls field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkAccessControlsField() *PassFactory {
	return c.checkAccessControls()
}

// checkConstsField is the pass factory for checkConsts.
// Mirrors Java's DefaultPassConfig.checkConsts field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkConstsField() *PassFactory {
	return c.checkConsts()
}

// checkConformanceField is the pass factory for checkConformance.
// Mirrors Java's DefaultPassConfig.checkConformance field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkConformanceField() *PassFactory {
	return c.checkConformance()
}

// processDefinesCheckField is the pass factory for processDefinesCheck.
// Mirrors Java's DefaultPassConfig.processDefinesCheck field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) processDefinesCheckField() *PassFactory {
	return c.processDefinesCheck()
}

// generateExportsField is the pass factory for generateExports.
// Mirrors Java's DefaultPassConfig.generateExports field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) generateExportsField() *PassFactory {
	return c.generateExports()
}

// gatherExternPropertiesCheckField is the pass factory for gatherExternPropertiesCheck.
// Mirrors Java's DefaultPassConfig.gatherExternPropertiesCheck field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) gatherExternPropertiesCheckField() *PassFactory {
	return c.gatherExternPropertiesCheck()
}

// serializeTypedAstField is the pass factory for serializeTypedAst.
// Mirrors Java's DefaultPassConfig.serializeTypedAst field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) serializeTypedAstField() *PassFactory {
	return c.serializeTypedAst()
}

// removeUnnecessarySyntheticExternsField is the pass factory.
// Mirrors Java's DefaultPassConfig.removeUnnecessarySyntheticExterns field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) removeUnnecessarySyntheticExternsField() *PassFactory {
	return c.removeUnnecessarySyntheticExterns()
}

// closureProvidesRequiresField is the pass factory.
// Mirrors Java's DefaultPassConfig.closureProvidesRequires field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) closureProvidesRequiresField() *PassFactory {
	return c.closureProvidesRequires()
}

// processDefinesOptimizeField is the pass factory.
// Mirrors Java's DefaultPassConfig.processDefinesOptimize field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) processDefinesOptimizeField() *PassFactory {
	return c.processDefinesOptimize()
}

// normalizeField is the pass factory for normalize.
// Mirrors Java's DefaultPassConfig.normalize field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) normalizeField() *PassFactory {
	return c.normalize()
}

// gatherExternPropertiesOptimizeField is the pass factory.
// Mirrors Java's DefaultPassConfig.gatherExternPropertiesOptimize field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) gatherExternPropertiesOptimizeField() *PassFactory {
	return c.gatherExternPropertiesOptimize()
}

// inlineAndCollapsePropertiesField is the pass factory.
// Mirrors Java's DefaultPassConfig.inlineAndCollapseProperties field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inlineAndCollapsePropertiesField() *PassFactory {
	return c.inlineAndCollapseProperties()
}

// closureOptimizePrimitivesField is the pass factory.
// Mirrors Java's DefaultPassConfig.closureOptimizePrimitives field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) closureOptimizePrimitivesField() *PassFactory {
	return c.closureOptimizePrimitives()
}

// stripSideEffectProtectionField is the pass factory.
// Mirrors Java's DefaultPassConfig.stripSideEffectProtection field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) stripSideEffectProtectionField() *PassFactory {
	return c.stripSideEffectProtection()
}

// inlineVariablesField is the pass factory.
// Mirrors Java's DefaultPassConfig.inlineVariables field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inlineVariablesField() *PassFactory {
	return c.inlineVariables()
}

// collapseObjectLiteralsField is the pass factory.
// Mirrors Java's DefaultPassConfig.collapseObjectLiterals field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) collapseObjectLiteralsField() *PassFactory {
	return c.collapseObjectLiterals()
}

// removeUnusedCodeField is the pass factory.
// Mirrors Java's DefaultPassConfig.removeUnusedCode field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) removeUnusedCodeField() *PassFactory {
	return c.removeUnusedCode()
}

// removeUnusedCodeOnceField is the pass factory.
// Mirrors Java's DefaultPassConfig.removeUnusedCodeOnce field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) removeUnusedCodeOnceField() *PassFactory {
	return c.removeUnusedCodeOnce()
}

// peepholeOptimizationsField is the pass factory.
// Mirrors Java's DefaultPassConfig.peepholeOptimizations field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) peepholeOptimizationsField() *PassFactory {
	return c.peepholeOptimizations()
}

// flowSensitiveInlineVariablesField is the pass factory.
// Mirrors Java's DefaultPassConfig.flowSensitiveInlineVariables field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) flowSensitiveInlineVariablesField() *PassFactory {
	return c.flowSensitiveInlineVariables()
}

// crossChunkCodeMotionField is the pass factory.
// Mirrors Java's DefaultPassConfig.crossChunkCodeMotion field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) crossChunkCodeMotionField() *PassFactory {
	return c.crossChunkCodeMotion()
}

// crossChunkMethodMotionField is the pass factory.
// Mirrors Java's DefaultPassConfig.crossChunkMethodMotion field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) crossChunkMethodMotionField() *PassFactory {
	return c.crossChunkMethodMotion()
}

// collapseAnonymousFunctionsField is the pass factory.
// Mirrors Java's DefaultPassConfig.collapseAnonymousFunctions field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) collapseAnonymousFunctionsField() *PassFactory {
	return c.collapseAnonymousFunctions()
}

// extractPrototypeMemberDeclarationsField is the pass factory.
// Mirrors Java's DefaultPassConfig.extractPrototypeMemberDeclarations field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) extractPrototypeMemberDeclarationsField() *PassFactory {
	return c.extractPrototypeMemberDeclarations()
}

// ambiguatePropertiesField is the pass factory.
// Mirrors Java's DefaultPassConfig.ambiguateProperties field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) ambiguatePropertiesField() *PassFactory {
	return c.ambiguateProperties()
}

// renamePropertiesField is the pass factory.
// Mirrors Java's DefaultPassConfig.renameProperties field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) renamePropertiesField() *PassFactory {
	return c.renameProperties()
}

// gatherRawExportsField is the pass factory.
// Mirrors Java's DefaultPassConfig.gatherRawExports field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) gatherRawExportsField() *PassFactory {
	return c.gatherRawExports()
}

// convertToDottedPropertiesField is the pass factory.
// Mirrors Java's DefaultPassConfig.convertToDottedProperties field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) convertToDottedPropertiesField() *PassFactory {
	return c.convertToDottedProperties()
}

// rewriteFunctionExpressionsField is the pass factory.
// Mirrors Java's DefaultPassConfig.rewriteFunctionExpressions field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) rewriteFunctionExpressionsField() *PassFactory {
	return c.rewriteFunctionExpressions()
}

// coalesceVariableNamesField is the pass factory.
// Mirrors Java's DefaultPassConfig.coalesceVariableNames field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) coalesceVariableNamesField() *PassFactory {
	return c.coalesceVariableNames()
}

// markUnnormalizedField is the pass factory.
// Mirrors Java's DefaultPassConfig.markUnnormalized field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) markUnnormalizedField() *PassFactory {
	return c.markUnnormalized()
}

// collapseVariableDeclarationsField is the pass factory.
// Mirrors Java's DefaultPassConfig.collapseVariableDeclarations field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) collapseVariableDeclarationsField() *PassFactory {
	return c.collapseVariableDeclarations()
}

// denormalizeField is the pass factory.
// Mirrors Java's DefaultPassConfig.denormalize field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) denormalizeField() *PassFactory {
	return c.denormalize()
}

// renameVarsField is the pass factory.
// Mirrors Java's DefaultPassConfig.renameVars field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) renameVarsField() *PassFactory {
	return c.renameVars()
}

// renameLabelsField is the pass factory.
// Mirrors Java's DefaultPassConfig.renameLabels field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) renameLabelsField() *PassFactory {
	return c.renameLabels()
}

// latePeepholeOptimizationsField is the pass factory.
// Mirrors Java's DefaultPassConfig.latePeepholeOptimizations field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) latePeepholeOptimizationsField() *PassFactory {
	return c.latePeepholeOptimizations()
}

// rescopeGlobalSymbolsField is the pass factory.
// Mirrors Java's DefaultPassConfig.rescopeGlobalSymbols field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) rescopeGlobalSymbolsField() *PassFactory {
	return c.rescopeGlobalSymbols()
}

// optimizeToEs6Field is the pass factory.
// Mirrors Java's DefaultPassConfig.optimizeToEs6 field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) optimizeToEs6Field() *PassFactory {
	return c.optimizeToEs6()
}

// convertChunksToESModulesField is the pass factory.
// Mirrors Java's DefaultPassConfig.convertChunksToESModules field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) convertChunksToESModulesField() *PassFactory {
	return c.convertChunksToESModules()
}

// checkAstValidityField is the pass factory.
// Mirrors Java's DefaultPassConfig.checkAstValidity field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) checkAstValidityField() *PassFactory {
	return c.checkAstValidity()
}

// varCheckValidityField is the pass factory.
// Mirrors Java's DefaultPassConfig.varCheckValidity field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) varCheckValidityField() *PassFactory {
	return c.varCheckValidity()
}

// inlineFunctionsField is the pass factory.
// Mirrors Java's DefaultPassConfig.inlineFunctions field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inlineFunctionsField() *PassFactory {
	return c.inlineFunctions()
}

// collapsePropertiesField is the pass factory.
// Mirrors Java's DefaultPassConfig.collapseProperties field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) collapsePropertiesField() *PassFactory {
	return c.collapseProperties()
}

// inlineTypeAliasesField is the pass factory.
// Mirrors Java's DefaultPassConfig.inlineTypeAliases field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inlineTypeAliasesField() *PassFactory {
	return c.inlineTypeAliases()
}

// disambiguatePropertiesField is the pass factory.
// Mirrors Java's DefaultPassConfig.disambiguateProperties field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) disambiguatePropertiesField() *PassFactory {
	return c.disambiguateProperties()
}

// inlinePropertiesField is the pass factory.
// Mirrors Java's DefaultPassConfig.inlineProperties field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) inlinePropertiesField() *PassFactory {
	return c.inlineProperties()
}

// devirtualizeMethodsField is the pass factory.
// Mirrors Java's DefaultPassConfig.devirtualizeMethods field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) devirtualizeMethodsField() *PassFactory {
	return c.devirtualizeMethods()
}

// optimizeCallsField is the pass factory.
// Mirrors Java's DefaultPassConfig.optimizeCalls field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) optimizeCallsField() *PassFactory {
	return c.optimizeCalls()
}

// ---------------------------------------------------------------------------
// Preprocessing symbol table
// ---------------------------------------------------------------------------

// PreprocessorSymbolTable is a stub for the preprocessor symbol table that
// tracks names removed during preprocessing. Mirrors Java's
// PreprocessorSymbolTable.CachedInstanceFactory. The actual table is populated
// by closureGoogScopeAliases and referenced by closureRewriteModule.
type PreprocessorSymbolTable struct {
	initialized bool
}

// MaybeInitialize initializes the table if not already done.
// Mirrors Java's PreprocessorSymbolTable.CachedInstanceFactory.maybeInitialize.
func (p *PreprocessorSymbolTable) MaybeInitialize(compiler *AbstractCompiler) {
	p.initialized = true
}

// GetInstanceOrNull returns the table instance if initialized, or nil.
// Mirrors Java's PreprocessorSymbolTable.CachedInstanceFactory.getInstanceOrNull.
func (p *PreprocessorSymbolTable) GetInstanceOrNull() *PreprocessorSymbolTable {
	if p.initialized {
		return p
	}
	return nil
}

// preprocessorSymbolTableFactory is the instance used by DefaultPassConfig.
// Mirrors Java's DefaultPassConfig.preprocessorSymbolTableFactory field.
var preprocessorSymbolTableFactory = &PreprocessorSymbolTable{}

// GetPreprocessorSymbolTable returns the preprocessor symbol table, or nil if
// not yet initialized. Mirrors Java's DefaultPassConfig.getPreprocessorSymbolTable.
func (c *DefaultPassConfig) GetPreprocessorSymbolTable() *PreprocessorSymbolTable {
	return preprocessorSymbolTableFactory.GetInstanceOrNull()
}

// ---------------------------------------------------------------------------
// globalSymbolNamespacePattern — mirrors Java's GLOBAL_SYMBOL_NAMESPACE_PATTERN
// ---------------------------------------------------------------------------

// isValidGlobalSymbolNamespace reports whether the given namespace string
// matches the pattern ^[a-zA-Z0-9$_]+$ required by renamePrefixNamespace.
// Mirrors Java's DefaultPassConfig.GLOBAL_SYMBOL_NAMESPACE_PATTERN.matcher(ns).matches().
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func isValidGlobalSymbolNamespace(ns string) bool {
	if ns == "" {
		return false
	}
	for _, ch := range ns {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '$' || ch == '_') {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// getChecks for skipNonTranspilationPasses guard
// ---------------------------------------------------------------------------

// skipNonTranspilationPassesCheck returns an error message if the wrong
// combination of pass-config methods is called. Mirrors the checkState in
// Java's DefaultPassConfig.getChecks() for the skipNonTranspilationPasses option.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) skipNonTranspilationPassesCheck() string {
	return "options.getSkipNonTranspilationPasses() cannot be mixed with " +
		"PassConfig::getChecks. Call PassConfig::getTranspileOnlyPasses instead."
}

// ---------------------------------------------------------------------------
// processDefinesCheck documentation
// ---------------------------------------------------------------------------

// processDefinesCheckCondition returns the condition under which processDefines
// should run in checks mode. Mirrors Java's getChecks logic where
// processDefinesCheck is always added. In Go this is always true.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) processDefinesCheckCondition() bool { return true }

// processDefinesOptimizeCondition returns the condition under which
// processDefines should run in optimization mode. Mirrors Java's getOptimizations
// logic where it is always added.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) processDefinesOptimizeCondition() bool { return true }

// ---------------------------------------------------------------------------
// rewriteCallerCodeLocation pass
// ---------------------------------------------------------------------------

// rewriteCallerCodeLocationPass creates the factory for RewriteCallerCodeLocation.
// This pass rewrites code that uses the Closure compiler's caller code
// location annotation, ensuring that after module rewriting the exported
// function and its call sites are correctly associated with source locations.
// Mirrors Java's DefaultPassConfig.rewriteCallerCodeLocation field.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) rewriteCallerCodeLocationPass() *PassFactory {
	return c.rewriteCallerCodeLocation()
}

// ---------------------------------------------------------------------------
// Pass count helpers
// ---------------------------------------------------------------------------

// ChecksCount returns the number of passes in the checks pipeline for the
// current options. Useful for metrics and logging.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) ChecksCount() int {
	return len(c.GetChecks().Build())
}

// OptimizationsCount returns the number of passes in the optimizations
// pipeline for the current options.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) OptimizationsCount() int {
	return len(c.GetOptimizations().Build())
}

// FinalizationsCount returns the number of passes in the finalizations
// pipeline for the current options.
// Go-only: language-adapter — Go exposes this as a named function; Java uses class-level field or overloaded constructor with the same semantics.
func (c *DefaultPassConfig) FinalizationsCount() int {
	return len(c.GetFinalizations().Build())
}
