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

func (c *DefaultPassConfig) makeNoopPass(name string) *PassFactory {
	return NewPassFactoryBuilder().
		SetName(name).
		SetInternalFactory(func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		}).
		Build()
}

func (c *DefaultPassConfig) makeConditionalPass(name string, cond func(*CompilerOptions) bool) *PassFactory {
	return NewPassFactoryBuilder().
		SetName(name).
		SetCondition(cond).
		SetInternalFactory(func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		}).
		Build()
}

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
// Currently a no-op — ordering assertions are added as passes mature.
func (c *DefaultPassConfig) assertValidOrderForChecks(checks *PassListBuilder) {}

// assertValidOrderForOptimizations validates that the optimization passes are
// in a consistent order. Mirrors Java's DefaultPassConfig.assertValidOrderForOptimizations.
// Currently a no-op — ordering assertions are added as passes mature.
func (c *DefaultPassConfig) assertValidOrderForOptimizations(optimizations *PassListBuilder) {}
