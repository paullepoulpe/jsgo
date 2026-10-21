// Port of com/google/javascript/jscomp/DefaultPassConfig.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DefaultPassConfig holds the pass factories and meta-data for native
// JSCompiler passes.
//
// NOTE(johnlenz): this needs some non-trivial refactoring. The pass config
// should use as little state as possible. The recommended way for a pass to
// leave behind some state for a subsequent pass is through the compiler
// object. Any other state remaining here should only be used when the pass
// config is creating the list of checks and optimizations, not after passes
// have started executing.
//
// The general goal is for this class to be as minimal as possible. Option
// validation should occur before this class configures the compiler,
// business logic should live here (passes should not be inlined, etc).
//
// The PassFactory fields are initialized by NewDefaultPassConfig, in Java
// declaration order (initFields).
type DefaultPassConfig struct {
	PassConfigBase

	// A symbol table for registering references that get removed during
	// preprocessing.
	preprocessorSymbolTableFactory *PreprocessorSymbolTable_CachedInstanceFactory

	checkSideEffects                             *PassFactory
	stripSideEffectProtection                    *PassFactory
	suspiciousCode                               *PassFactory
	extraRequires                                *PassFactory
	checkMissingRequires                         *PassFactory
	checkJsDocAndEs6Modules                      *PassFactory
	generateExports                              *PassFactory
	generateIjs                                  *PassFactory
	removeExtraRequires                          *PassFactory
	exportTestFunctions                          *PassFactory
	gatherRawExports                             *PassFactory
	closurePrimitives                            *PassFactory
	closureProvidesRequires                      *PassFactory
	angularPass                                  *PassFactory
	replaceMessagesForChrome                     *PassFactory
	closureGoogScopeAliasesForIjs                *PassFactory
	closureGoogScopeAliases                      *PassFactory
	injectRuntimeLibraries                       *PassFactory
	removeWeakSources                            *PassFactory
	declaredGlobalExternsOnWindow                *PassFactory
	checkTypeImportCodeReferences                *PassFactory
	closureCheckModule                           *PassFactory
	closureRewriteModule                         *PassFactory
	checkClosureImports                          *PassFactory
	rewriteGoogJsImports                         *PassFactory
	closureReplaceGetCssName                     *PassFactory
	createSyntheticBlocks                        *PassFactory
	earlyPeepholeOptimizations                   *PassFactory
	earlyInlineVariables                         *PassFactory
	peepholeOptimizations                        *PassFactory
	peepholeOptimizationsOnceNormalized          *PassFactory
	peepholeOptimizationsOnceNonNormalized       *PassFactory
	latePeepholeOptimizations                    *PassFactory
	checkVars                                    *PassFactory
	inferConsts                                  *PassFactory
	checkRegExp                                  *PassFactory
	checkRegExpForOptimizations                  *PassFactory
	checkVariableReferences                      *PassFactory
	checkSuper                                   *PassFactory
	clearTypedScopeCreatorPass                   *PassFactory
	clearTopTypedScopePass                       *PassFactory
	inferTypes                                   *PassFactory
	inferJsDocInfo                               *PassFactory
	checkTypes                                   *PassFactory
	checkControlFlow                             *PassFactory
	checkAccessControls                          *PassFactory
	lintChecks                                   *PassFactory
	analyzerChecks                               *PassFactory
	checkRequiresAndProvidesSorted               *PassFactory
	checkStrictMode                              *PassFactory
	processTweaks                                *PassFactory
	processDefinesCheck                          *PassFactory
	processDefinesOptimize                       *PassFactory
	stripCode                                    *PassFactory
	checkConsts                                  *PassFactory
	rewriteCallerCodeLocation                    *PassFactory
	replaceToggles                               *PassFactory
	replaceIdGenerators                          *PassFactory
	replaceStrings                               *PassFactory
	closureCodeRemoval                           *PassFactory
	closureOptimizePrimitives                    *PassFactory
	rescopeGlobalSymbols                         *PassFactory
	convertChunksToESModules                     *PassFactory
	inlineAndCollapseProperties                  *PassFactory
	collapseObjectLiterals                       *PassFactory
	disambiguateProperties                       *PassFactory
	devirtualizeMethods                          *PassFactory
	optimizeCalls                                *PassFactory
	optimizeConstructors                         *PassFactory
	markPureFunctions                            *PassFactory
	inlineVariables                              *PassFactory
	inlineConstants                              *PassFactory
	inlineSimpleMethods                          *PassFactory
	deadAssignmentsElimination                   *PassFactory
	deadPropertyAssignmentElimination            *PassFactory
	inlineFunctions                              *PassFactory
	inlineProperties                             *PassFactory
	isolatePolyfills                             *PassFactory
	removeUnusedCode                             *PassFactory
	removeUnusedCodeOnce                         *PassFactory
	crossChunkCodeMotion                         *PassFactory
	crossChunkMethodMotion                       *PassFactory
	flowSensitiveInlineVariables                 *PassFactory
	coalesceVariableNames                        *PassFactory
	postNormalizePeephole                        *PassFactory
	collapseVariableDeclarations                 *PassFactory
	extractPrototypeMemberDeclarations           *PassFactory
	rewriteFunctionExpressions                   *PassFactory
	collapseAnonymousFunctions                   *PassFactory
	rewriteGlobalDeclarationsForTryCatchWrapping *PassFactory
	aliasStrings                                 *PassFactory
	ambiguateProperties                          *PassFactory
	markUnnormalized                             *PassFactory
	normalize                                    *PassFactory
	externExports                                *PassFactory
	denormalize                                  *PassFactory
	invertContextualRenaming                     *PassFactory
	renameProperties                             *PassFactory
	removePropertyRenamingCalls                  *PassFactory
	renameVars                                   *PassFactory
	renameLabels                                 *PassFactory
	convertToDottedProperties                    *PassFactory
	checkAstValidity                             *PassFactory
	varCheckValidity                             *PassFactory
	instrumentForCodeCoverage                    *PassFactory
	gatherExternPropertiesCheck                  *PassFactory
	gatherExternPropertiesOptimize               *PassFactory
	polymerPass                                  *PassFactory
	chromePass                                   *PassFactory
	j2clConstantHoisterPass                      *PassFactory
	j2clClinitPass                               *PassFactory
	j2clPropertyInlinerPass                      *PassFactory
	j2clPass                                     *PassFactory
	j2clUtilGetDefineRewriterPass                *PassFactory
	j2clAssertRemovalPass                        *PassFactory
	j2clSourceFileChecker                        *PassFactory
	j2clChecksPass                               *PassFactory
	checkConformance                             *PassFactory
	removeCastNodes                              *PassFactory
	typesToColors                                *PassFactory
	serializeTypedAst                            *PassFactory
	removeUnnecessarySyntheticExterns            *PassFactory
	optimizeToEs6                                *PassFactory
	whitespaceWrapGoogModules                    *PassFactory
	rewriteCommonJsModules                       *PassFactory
	rewriteScriptsToEs6Modules                   *PassFactory
	gatherModuleMetadataPass                     *PassFactory
	createModuleMapPass                          *PassFactory
	gatherGettersAndSetters                      *PassFactory
	addSyntheticScript                           *PassFactory
	removeSyntheticScript                        *PassFactory
	mergeSyntheticScript                         *PassFactory
	protectLocaleData                            *PassFactory
	substituteLocaleData                         *PassFactory
	transpileOnlyClosureUnaware                  *PassFactory
	transpileAndOptimizeClosureUnaware           *PassFactory
}

// For the --mark-as-compiled pass
//
// Java: DefaultPassConfig.COMPILED_CONSTANT_NAME
const defaultPassConfig_COMPILED_CONSTANT_NAME = "COMPILED"

// Constant name for Closure's locale
//
// Java: DefaultPassConfig.CLOSURE_LOCALE_CONSTANT_NAME
const defaultPassConfig_CLOSURE_LOCALE_CONSTANT_NAME = "goog.LOCALE"

// Miscellaneous errors.
//
// Java: DefaultPassConfig.GLOBAL_SYMBOL_NAMESPACE_PATTERN
var defaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN = jutil.Pattern_Compile("^[a-zA-Z0-9$_]+$")

// Java: DefaultPassConfig(CompilerOptions)
func NewDefaultPassConfig(options *CompilerOptions) *DefaultPassConfig {
	c := &DefaultPassConfig{}
	c.self = c
	c.InitPassConfigBase(options)
	c.initFields()
	return c
}

// initFields runs the field initializers of DefaultPassConfig, in Java
// declaration order (after super(options), as in Java).
//
// Java: DefaultPassConfig.<init> (field initializers)
func (c *DefaultPassConfig) initFields() {
	options := c.options

	c.preprocessorSymbolTableFactory = NewPreprocessorSymbolTable_CachedInstanceFactory()

	c.checkSideEffects =
		PassFactory_Builder_().
			SetName("checkSideEffects").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCheckSideEffects(
						compiler,
						options.GetCheckSuspiciousCode(),
						options.ShouldProtectHiddenSideEffects())
				}).
			Build()

	// Removes the "protector" functions that were added by CheckSideEffects.
	c.stripSideEffectProtection =
		PassFactory_Builder_().
			SetName(PassNames_STRIP_SIDE_EFFECT_PROTECTION).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewCheckSideEffects_StripProtection(compiler)
			}).
			Build()

	// Checks for code that is probably wrong (such as stray expressions).
	c.suspiciousCode =
		PassFactory_Builder_().
			SetName("suspiciousCode").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					var sharedCallbacks []NodeTraversal_Callback
					if options.GetCheckSuspiciousCode() {
						sharedCallbacks = append(sharedCallbacks, NewCheckSuspiciousCode())
						sharedCallbacks = append(sharedCallbacks, NewCheckDuplicateCase(compiler))
					}

					if options.Enables(DiagnosticGroups_GLOBAL_THIS) {
						sharedCallbacks = append(sharedCallbacks, NewCheckGlobalThis(compiler))
					}

					if options.Enables(DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT) {
						sharedCallbacks = append(sharedCallbacks, NewCheckDebuggerStatement(compiler))
					}

					return defaultPassConfig_combineChecks(compiler, sharedCallbacks)
				}).
			Build()

	// Checks that all goog.require()s are used.
	c.extraRequires =
		PassFactory_Builder_().
			SetName("checkExtraRequires").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCheckExtraRequires(compiler, options.GetUnusedImportsToRemove())
				}).
			Build()

	c.checkMissingRequires =
		PassFactory_Builder_().
			SetName("checkMissingRequires").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCheckMissingRequires(compiler, compiler.GetModuleMetadataMap())
				}).
			Build()

	// Verifies JSDoc annotations are used properly and checks for ES modules.
	c.checkJsDocAndEs6Modules =
		PassFactory_Builder_().
			SetName("checkJsDocAndEs6Modules").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return defaultPassConfig_combineChecks(
						compiler,
						[]NodeTraversal_Callback{NewCheckJSDoc(compiler), NewEs6CheckModule(compiler)})
				}).
			Build()

	// Generates exports for @export annotations.
	c.generateExports =
		PassFactory_Builder_().
			SetName(PassNames_GENERATE_EXPORTS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					convention := compiler.GetCodingConvention()
					pass :=
						NewGenerateExports(
							compiler,
							options.ShouldExportLocalPropertyDefinitions(),
							convention.GetExportSymbolFunction(),
							convention.GetExportPropertyFunction())
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						pass.Process(externs, root)
						compiler.AddExportedNames(pass.GetExportedVariableNames())
					})
				}).
			Build()

	c.generateIjs =
		PassFactory_Builder_().
			SetName("generateIjs").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewConvertToTypedInterface(compiler)
			}).
			Build()

	// Prunes unnecessary goog.requires and in .i.js files
	// (go/exclude-unnecessary-goog-requires-in-ijs)
	c.removeExtraRequires =
		PassFactory_Builder_().
			SetName("removeExtraRequires").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewExtraRequireRemover(compiler)
			}).
			Build()

	// Generates exports for functions associated with JsUnit.
	c.exportTestFunctions =
		PassFactory_Builder_().
			SetName(PassNames_EXPORT_TEST_FUNCTIONS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					convention := compiler.GetCodingConvention()
					if convention.GetExportSymbolFunction() != nil {
						return NewExportTestFunctions(
							compiler,
							convention.GetExportSymbolFunction(),
							convention.GetExportPropertyFunction())
					} else {
						return NewErrorPass(compiler, defaultPassConfig_GENERATE_EXPORTS_ERROR)
					}
				}).
			Build()

	// Raw exports processing pass.
	c.gatherRawExports =
		PassFactory_Builder_().
			SetName(PassNames_GATHER_RAW_EXPORTS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					pass := NewGatherRawExports(compiler)

					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						pass.Process(externs, root)
						compiler.AddExportedNames(pass.GetExportedVariableNames())
					})
				}).
			Build()

	// Closure pre-processing pass.
	c.closurePrimitives =
		PassFactory_Builder_().
			SetName("closurePrimitives").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
					return NewProcessClosurePrimitives(compiler)
				}).
			Build()

	// Closure provide/require rewriting pass.
	c.closureProvidesRequires =
		PassFactory_Builder_().
			SetName("closureProvidesRequires").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
					pass :=
						NewProcessClosureProvidesAndRequires(
							compiler, options.ShouldPreservesGoogProvidesAndRequires())
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						pass.Process(externs, root)
						compiler.AddExportedNames(pass.GetExportedVariableNames())
					})
				}).
			Build()

	// Process AngularJS-specific annotations.
	c.angularPass =
		PassFactory_Builder_().
			SetName(PassNames_ANGULAR_PASS).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewAngularPass(compiler)
			}).
			Build()

	c.replaceMessagesForChrome =
		PassFactory_Builder_().
			SetName(PassNames_REPLACE_MESSAGES).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewReplaceMessagesForChrome(
						compiler, NewGoogleJsMessageIdGenerator(options.GetTcProjectId()),
						/* allow messages with goog.getMsg */
					)
				}).
			Build()

	// Applies aliases and inlines goog.scope.
	c.closureGoogScopeAliasesForIjs =
		PassFactory_Builder_().
			SetName("closureGoogScopeAliasesForIjs").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return ScopedAliases_Builder_(compiler).Build()
			}).
			Build()

	// Applies aliases and inlines goog.scope, storing information about the
	// transformations performed.
	c.closureGoogScopeAliases =
		PassFactory_Builder_().
			SetName("closureGoogScopeAliases").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
					return ScopedAliases_Builder_(compiler).
						SetPreprocessorSymbolTable(c.preprocessorSymbolTableFactory.GetInstanceOrNull()).
						SetModuleMetadataMap(compiler.GetModuleMetadataMap()).
						SetInvalidModuleGetHandling(ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME).
						Build()
				}).
			Build()

	c.injectRuntimeLibraries =
		PassFactory_Builder_().
			SetName("InjectRuntimeLibraries").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewInjectRuntimeLibraries(
						compiler,
						jutil.NewLinkedHashSet(compiler.GetOptions().GetForceLibraryInjectionList()...))
				}).
			Build()

	c.removeWeakSources =
		PassFactory_Builder_().
			SetName("removeWeakSources").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewRemoveWeakSources(compiler)
			}).
			Build()

	c.declaredGlobalExternsOnWindow =
		PassFactory_Builder_().
			SetName(PassNames_DECLARED_GLOBAL_EXTERNS_ON_WINDOW).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewDeclaredGlobalExternsOnWindow(compiler)
			}).
			Build()

	c.checkTypeImportCodeReferences =
		PassFactory_Builder_().
			SetName("checkTypeImportCodeReferences").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewCheckTypeImportCodeReferences(compiler)
			}).
			Build()

	// Checks of correct usage of goog.module
	c.closureCheckModule =
		PassFactory_Builder_().
			SetName("closureCheckModule").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap())
				}).
			Build()

	// Rewrites goog.module
	c.closureRewriteModule =
		PassFactory_Builder_().
			SetName("closureRewriteModule").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
					return NewClosureRewriteModule(
						compiler,
						c.preprocessorSymbolTableFactory.GetInstanceOrNull(),
						compiler.GetTopScope())
				}).
			Build()

	// Checks goog.require, goog.forwardDeclare, goog.requireType, and
	// goog.module.get calls
	c.checkClosureImports =
		PassFactory_Builder_().
			SetName("checkGoogRequires").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCheckClosureImports(compiler, compiler.GetModuleMetadataMap())
				}).
			Build()

	// Rewrite imports for Closure Library's goog.js file to global goog
	// references.
	c.rewriteGoogJsImports =
		PassFactory_Builder_().
			SetName("rewriteGoogJsImports").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewRewriteGoogJsImports(
						compiler,
						RewriteGoogJsImports_Mode_LINT_AND_REWRITE,
						compiler.GetModuleMap())
				}).
			Build()

	// Processes goog.getCssName. The cssRenamingMap is used to lookup
	// replacement values for the classnames. If null, the raw class names are
	// inlined.
	c.closureReplaceGetCssName =
		PassFactory_Builder_().
			SetName("closureReplaceGetCssName").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, jsRoot *rhino.Node) {
						// Go-only: language-adapter — Optional<ImmutableSet.Builder<String>>:
						// nil is Optional.empty().
						var cssNames *jutil.LinkedHashSet[string]
						if options.ShouldGatherCssNames() {
							cssNames = jutil.NewLinkedHashSet[string]()
						}

						pass :=
							NewReplaceCssNames(
								compiler,
								options.GetCssRenamingMap(),
								ReplaceCssNames_CssNameCollectorFunc(func(cssName string) {
									if cssNames != nil {
										cssNames.Add(cssName)
									}
								}),
								options.GetCssRenamingSkiplist())
						pass.Process(externs, jsRoot)

						if cssNames != nil {
							compiler.SetCssNames(cssNames)
						} else {
							compiler.SetCssNames(nil)
						}
					})
				}).
			Build()

	// Creates synthetic blocks to prevent FoldConstants from moving code past
	// markers in the source.
	c.createSyntheticBlocks =
		PassFactory_Builder_().
			SetName("createSyntheticBlocks").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCreateSyntheticBlocks(
						compiler,
						*options.GetSyntheticBlockStartMarker(),
						*options.GetSyntheticBlockEndMarker())
				}).
			Build()

	c.earlyPeepholeOptimizations =
		PassFactory_Builder_().
			SetName("earlyPeepholeOptimizations").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					useTypesForOptimization :=
						compiler.GetOptions().ShouldUseTypesForLocalOptimization()
					var peepholeOptimizations []AbstractPeepholeOptimization
					peepholeOptimizations = append(peepholeOptimizations, NewPeepholeRemoveDeadCode())
					if compiler.GetOptions().GetJ2clPass().ShouldAddJ2clPasses() {
						peepholeOptimizations = append(peepholeOptimizations,
							NewJ2clEqualitySameRewriterPass(useTypesForOptimization))
					}
					return NewPeepholeOptimizationsPassOfList(
						compiler, "earlyPeepholeOptimizations", peepholeOptimizations)
				}).
			Build()

	c.earlyInlineVariables =
		PassFactory_Builder_().
			SetName("earlyInlineVariables").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					var mode InlineVariables_Mode
					if options.ShouldInlineVariables() {
						mode = InlineVariables_Mode_ALL
					} else if options.ShouldInlineLocalVariables() {
						mode = InlineVariables_Mode_LOCALS_ONLY
					} else {
						panic(jutil.NewIllegalStateExceptionMsg("No variable inlining option set."))
					}
					return NewInlineVariables(compiler, mode)
				}).
			Build()

	// Various peephole optimizations.
	c.peepholeOptimizations =
		PassFactory_Builder_().
			SetName(PassNames_PEEPHOLE_OPTIMIZATIONS).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return defaultPassConfig_createPeepholeOptimizationsPass(compiler, PassNames_PEEPHOLE_OPTIMIZATIONS, true)
				}).
			Build()

	// Various peephole optimizations.
	c.peepholeOptimizationsOnceNormalized =
		PassFactory_Builder_().
			SetName(PassNames_PEEPHOLE_OPTIMIZATIONS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return defaultPassConfig_createPeepholeOptimizationsPass(
						compiler,
						PassNames_PEEPHOLE_OPTIMIZATIONS,
						/* expectAstIsNormalized= */ true)
				}).
			Build()

	// Various peephole optimizations.
	c.peepholeOptimizationsOnceNonNormalized =
		PassFactory_Builder_().
			SetName(PassNames_PEEPHOLE_OPTIMIZATIONS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return defaultPassConfig_createPeepholeOptimizationsPass(
						compiler,
						PassNames_PEEPHOLE_OPTIMIZATIONS,
						/* expectAstIsNormalized= */ false)
				}).
			Build()

	// Same as peepholeOptimizations but aggressively merges code together
	c.latePeepholeOptimizations =
		PassFactory_Builder_().
			SetName("latePeepholeOptimizations").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					const late = true
					useTypesForOptimization :=
						options.ShouldUseTypesForLocalOptimization()
					return NewPeepholeOptimizationsPass(
						compiler,
						"latePeepholeOptimizations",
						NewStatementFusion(),
						NewPeepholeRemoveDeadCode(),
						NewPeepholeMinimizeConditions(late),
						NewPeepholeSubstituteAlternateSyntax(late),
						NewPeepholeReplaceKnownMethods(late, useTypesForOptimization),
						NewPeepholeFoldConstants(late, useTypesForOptimization))
				}).
			Build()

	// Checks that all variables are defined.
	c.checkVars =
		PassFactory_Builder_().SetName(PassNames_CHECK_VARS).SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
			return NewVarCheck(compiler)
		}).Build()

	// Infers constants.
	c.inferConsts =
		PassFactory_Builder_().
			SetName(PassNames_INFER_CONSTS).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewInferConsts(compiler)
			}).
			Build()

	// Checks for RegExp references.
	c.checkRegExp =
		PassFactory_Builder_().
			SetName(PassNames_CHECK_REG_EXP).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					pass :=
						NewCheckRegExp(
							compiler,
							options.GetAssumePropertiesAreStaticallyAnalyzable(),
							/* reportErrors= */ true)

					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						pass.Process(externs, root)
						compiler.SetHasRegExpGlobalReferences(pass.IsGlobalRegExpPropertiesUsed())
					})
				}).
			Build()

	// This pass is redundant if we've already run checkRegExp during this compilation
	c.checkRegExpForOptimizations =
		PassFactory_Builder_().
			SetName("checkRegExpForOptimizations").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					regExpCheck :=
						NewCheckRegExp(
							compiler,
							options.GetAssumePropertiesAreStaticallyAnalyzable(),
							/* reportErrors= */ false)

					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						regExpCheck.Process(externs, root)
						compiler.SetHasRegExpGlobalReferences(
							regExpCheck.IsGlobalRegExpPropertiesUsed())
					})
				}).
			Build()

	// Checks that references to variables look reasonable.
	c.checkVariableReferences =
		PassFactory_Builder_().
			SetName(PassNames_CHECK_VARIABLE_REFERENCES).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewVariableReferenceCheck(compiler)
			}).
			Build()

	c.checkSuper =
		PassFactory_Builder_().SetName("checkSuper").SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
			return NewCheckSuper(compiler)
		}).Build()

	// Clears the typed scope creator and all local typed scopes.
	c.clearTypedScopeCreatorPass =
		PassFactory_Builder_().
			SetName("clearTypedScopeCreatorPass").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) { compiler.ClearTypedScopeCreator() })
			}).
			Build()

	// Clears the top typed scope when we're done with it.
	c.clearTopTypedScopePass =
		PassFactory_Builder_().
			SetName("clearTopTypedScopePass").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						// clear these scopes which we don't need anymore so they can be garbage
						// collected
						compiler.SetTopScope(nil)
						for _, compilerInput := range compiler.GetInputsInOrder() {
							compilerInput.SetTypedScope(nil)
						}
					})
				}).
			Build()

	// Runs type inference.
	c.inferTypes =
		PassFactory_Builder_().
			SetName(PassNames_INFER_TYPES).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(unused *rhino.Node, srcRoot *rhino.Node) {
						if comp, ok := compiler.(*Compiler); ok {
							comp.evictPassArenaCaches()
						}
						globalRoot := srcRoot.GetParent()

						inferencePass :=
							NewTypeInferencePass(
								compiler,
								compiler.GetReverseAbstractInterpreter(),
								compiler.GetTypedScopeCreator().(*TypedScopeCreator))
						compiler.SetTypeCheckingHasRun(true)
						compiler.SetTopScope(inferencePass.InferAllScopes(globalRoot))
					})
				}).
			Build()

	c.inferJsDocInfo =
		PassFactory_Builder_().
			SetName("inferJsDocInfo").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewInferJSDocInfo(compiler)
			}).
			Build()

	// Checks type usage
	c.checkTypes =
		PassFactory_Builder_().
			SetName(PassNames_CHECK_TYPES).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						check :=
							NewTypeCheckWithScope(
								compiler,
								compiler.GetReverseAbstractInterpreter(),
								compiler.GetTypeRegistry(),
								compiler.GetTopScope(),
								compiler.GetTypedScopeCreator().(*TypedScopeCreator)).
								ReportUnknownTypes(
									options.Enables(
										DiagnosticGroup_ForType(TypeCheck_UNKNOWN_EXPR_TYPE))).
								ReportMissingProperties(
									!options.Disables(
										DiagnosticGroup_ForType(TypeCheck_INEXISTENT_PROPERTY)))
						check.Process(externs, root)
						compiler.GetErrorManager().SetTypedPercent(check.GetTypedPercent())
					})
				}).
			Build()

	// Checks possible execution paths of the program for problems: missing
	// return statements and dead code.
	c.checkControlFlow =
		PassFactory_Builder_().
			SetName("checkControlFlow").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					var callbacks []NodeTraversal_Callback
					if !options.Disables(DiagnosticGroups_CHECK_USELESS_CODE) {
						callbacks = append(callbacks, NewCheckUnreachableCode(compiler))
					}
					if !options.Disables(DiagnosticGroups_MISSING_RETURN) {
						callbacks = append(callbacks, NewCheckMissingReturn(compiler))
					}
					return defaultPassConfig_combineChecks(compiler, callbacks)
				}).
			Build()

	// Checks access controls. Depends on type-inference.
	c.checkAccessControls =
		PassFactory_Builder_().
			SetName("checkAccessControls").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewCheckAccessControls(compiler)
			}).
			Build()

	// Runs the single-file linter passes
	//
	// These is NOT the configuration for the standalone Linter binary. New
	// linter passes must also be added to LintPassConfig as well as this list.
	c.lintChecks =
		PassFactory_Builder_().
			SetName(PassNames_LINT_CHECKS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					callbacks := []NodeTraversal_Callback{
						NewCheckConstPrivateProperties(compiler),
						NewCheckConstantCaseNames(compiler),
						NewCheckEmptyStatements(compiler),
						NewCheckEnums(compiler),
						NewCheckEs6ModuleFileStructure(compiler),
						NewCheckEs6Modules(compiler),
						NewCheckNoMutatedEs6Exports(compiler),
						NewCheckGoogModuleTypeScriptName(compiler),
						NewCheckInterfaces(compiler),
						NewCheckJSDocStyle(compiler),
						NewCheckMissingSemicolon(compiler),
						NewCheckNullabilityModifiers(compiler),
						NewCheckPrimitiveAsObject(compiler),
						NewCheckPrototypeProperties(compiler),
						NewCheckUnusedPrivateProperties(compiler),
						NewCheckUnusedLabels(compiler),
						NewCheckUselessBlocks(compiler),
						NewCheckVar(compiler),
					}
					return defaultPassConfig_combineChecks(compiler, callbacks)
				}).
			Build()

	c.analyzerChecks =
		PassFactory_Builder_().
			SetName(PassNames_ANALYZER_CHECKS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					callbacks :=
						[]NodeTraversal_Callback{
							NewCheckArrayWithGoogObject(compiler),
							NewImplicitNullabilityCheck(compiler),
							NewCheckNestedNames(compiler)}

					return defaultPassConfig_combineChecks(compiler, callbacks)
				}).
			Build()

	c.checkRequiresAndProvidesSorted =
		PassFactory_Builder_().
			SetName("checkRequiresAndProvidesSorted").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return defaultPassConfig_combineChecks(
						compiler,
						[]NodeTraversal_Callback{
							NewCheckProvidesSorted(CheckProvidesSorted_Mode_COLLECT_AND_REPORT),
							NewCheckRequiresSorted(CheckRequiresSorted_Mode_COLLECT_AND_REPORT)})
				}).
			Build()

	// Checks that the code is ES5 strict compliant.
	c.checkStrictMode =
		PassFactory_Builder_().
			SetName("checkStrictMode").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					defaultLevel := CheckLevel_OFF
					if options.ExpectStrictModeInput() {
						defaultLevel = CheckLevel_ERROR
					}
					return NewStrictModeCheck(compiler, defaultLevel)
				}).
			Build()

	// Process goog.tweak.getTweak() calls.
	c.processTweaks =
		PassFactory_Builder_().
			SetName("processTweaks").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, jsRoot *rhino.Node) {
						NewProcessTweaks(compiler, options.GetTweakProcessing().ShouldStrip()).
							Process(externs, jsRoot)
					})
				}).
			Build()

	// Check @define-annotated constants.
	c.processDefinesCheck = c.createProcessDefines(ProcessDefines_Mode_CHECK)

	// Replace @define-annotated constants.
	c.processDefinesOptimize =
		c.createProcessDefines(ProcessDefines_Mode_OPTIMIZE)

	// Strips code for smaller compiled code. This is useful for removing debug
	// statements to prevent leaking them publicly.
	c.stripCode =
		PassFactory_Builder_().
			SetName("stripCode").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, jsRoot *rhino.Node) {
						options := compiler.GetOptions()
						pass :=
							NewStripCode(
								compiler,
								options.GetStripTypes(),
								options.GetStripNameSuffixes(),
								options.GetStripNamePrefixes(),
								options.GetTweakProcessing().ShouldStrip())

						pass.Process(externs, jsRoot)
					})
				}).
			// TODO(johnlenz): StripCode may be fooled by some newer features, like destructuring,
			// an).build();
			Build()

	// Checks that all constants are not modified
	c.checkConsts =
		PassFactory_Builder_().SetName("checkConsts").SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
			return NewConstCheck(compiler)
		}).Build()

	c.rewriteCallerCodeLocation =
		PassFactory_Builder_().
			SetName("rewriteCallerCodeLocation").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewRewriteCallerCodeLocation(compiler)
			}).
			Build()

	// Replaces goog.toggle calls with toggle lookups.
	c.replaceToggles =
		PassFactory_Builder_().
			SetName("replaceToggles").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewReplaceToggles(compiler)
			}).
			Build()

	// Generates unique ids.
	c.replaceIdGenerators =
		PassFactory_Builder_().
			SetName(PassNames_REPLACE_ID_GENERATORS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						pass :=
							NewReplaceIdGenerators(
								compiler,
								options.NeedsTranspilationOf(parser.FeatureSet_Feature_TEMPLATE_LITERALS),
								options.GetIdGenerators(),
								options.ShouldGeneratePseudoNames(),
								options.GetIdGeneratorsMapSerialized(),
								options.GetXidHashFunction())
						pass.Process(externs, root)
						compiler.SetIdGeneratorMap(pass.GetSerializedIdMappings())
					})
				}).
			Build()

	// Replace strings.
	c.replaceStrings =
		PassFactory_Builder_().
			SetName("replaceStrings").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						pass :=
							NewReplaceStrings(
								compiler,
								options.GetReplaceStringsPlaceholderToken(),
								options.GetReplaceStringsFunctionDescriptions())
						pass.Process(externs, root)
						compiler.SetStringMap(pass.GetStringMap())
					})
				}).
			Build()

	// Remove variables set to goog.abstractMethod.
	c.closureCodeRemoval =
		PassFactory_Builder_().
			SetName("closureCodeRemoval").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewClosureCodeRemoval(
						compiler,
						options.ShouldRemoveAbstractMethods(),
						options.ShouldRemoveClosureAsserts())
				}).
			Build()

	// Special case optimizations for closure functions.
	c.closureOptimizePrimitives =
		PassFactory_Builder_().
			SetName("closureOptimizePrimitives").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewClosureOptimizePrimitives(
						compiler, compiler.GetOptions().GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015))
				}).
			Build()

	// Puts global symbols into a single object.
	c.rescopeGlobalSymbols =
		PassFactory_Builder_().
			SetName("rescopeGlobalSymbols").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewRescopeGlobalSymbols(
						compiler,
						*options.GetRenamePrefixNamespace(),
						options.AssumeCrossChunkNamesForRenamePrefixNamespace(),
						options.GetOptimizeLocalAccessForGlobalSymbolNamespace())
				}).
			Build()

	// Converts cross chunk references into ES Module import and export
	// statements.
	c.convertChunksToESModules =
		PassFactory_Builder_().
			SetName("convertChunksToESModules").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewConvertChunksToESModules(compiler)
			}).
			Build()

	// Perform inlining of aliases and collapsing of qualified names to global
	// variables in order to improve later optimizations, such as
	// RemoveUnusedCode.
	c.inlineAndCollapseProperties =
		PassFactory_Builder_().
			SetName("inlineAndCollapseProperties").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return InlineAndCollapseProperties_Builder_(compiler).
						SetPropertyCollapseLevel(options.GetPropertyCollapseLevel()).
						SetChunkOutputType(options.GetChunkOutputType()).
						SetHaveModulesBeenRewritten(options.GetProcessCommonJSModules()).
						SetModuleResolutionMode(options.GetModuleResolutionMode()).
						Build()
				}).
			SetPreconditionCheck(
				func(options *CompilerOptions) *PassFactory_PreconditionResult {
					return NewPassFactory_PreconditionResult(
						options.GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_NONE ||
							options.GetAssumePropertiesAreStaticallyAnalyzable(),
						jutil.StringPtr("requires assumePropertiesAreStaticallyAnalyzable to be enabled"))
				}).
			Build()

	// Rewrite properties as variables.
	c.collapseObjectLiterals =
		PassFactory_Builder_().
			SetName(PassNames_COLLAPSE_OBJECT_LITERALS).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewInlineObjectLiterals(compiler, compiler.GetUniqueNameIdSupplier())
				}).
			Build()

	// Disambiguate property names based on type information.
	c.disambiguateProperties =
		PassFactory_Builder_().
			SetName(PassNames_DISAMBIGUATE_PROPERTIES).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewDisambiguateProperties(compiler, options.GetPropertiesThatMustDisambiguate())
				}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Rewrite instance methods as static methods, to make them easier to
	// inline.
	c.devirtualizeMethods =
		PassFactory_Builder_().
			SetName(PassNames_DEVIRTUALIZE_METHODS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return OptimizeCalls_Builder_().
						SetCompiler(compiler).
						SetConsiderExterns(false).
						AddPass(NewDevirtualizeMethods(compiler)).
						Build()
				}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Optimizes unused function arguments, unused return values, and inlines
	// constant parameters. Also runs RemoveUnusedCode.
	c.optimizeCalls =
		PassFactory_Builder_().
			SetName(PassNames_OPTIMIZE_CALLS).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return OptimizeCalls_Builder_().
						SetCompiler(compiler).
						SetConsiderExterns(false).
						// Remove unused return values.
						AddPass(NewOptimizeReturns(compiler).AsCallGraphCompilerPass()).
						// Remove all parameters that are constants or unused.
						AddPass(NewOptimizeParameters(compiler).AsCallGraphCompilerPass()).
						Build()
				}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Removes ECMAScript class constructors when an implicit constructor is
	// sufficient.
	c.optimizeConstructors =
		PassFactory_Builder_().
			SetName("optimizeConstructors").
			SetRunInFixedPointLoop(false).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return OptimizeCalls_Builder_().
						SetCompiler(compiler).
						SetConsiderExterns(false).
						// Remove redundant constructor definitions.
						AddPass(NewOptimizeConstructors(compiler).AsCallGraphCompilerPass()).
						Build()
				}).
			Build()

	// Look for function calls that are pure, and annotate them that way.
	c.markPureFunctions =
		PassFactory_Builder_().
			SetName("markPureFunctions").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewPureFunctionIdentifier_Driver(compiler)
			}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Inlines variables heuristically.
	c.inlineVariables =
		PassFactory_Builder_().
			SetName(PassNames_INLINE_VARIABLES).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					var mode InlineVariables_Mode
					if options.ShouldInlineVariables() {
						mode = InlineVariables_Mode_ALL
					} else if options.ShouldInlineLocalVariables() {
						mode = InlineVariables_Mode_LOCALS_ONLY
					} else {
						panic(jutil.NewIllegalStateExceptionMsg("No variable inlining option set."))
					}
					return NewInlineVariables(compiler, mode)
				}).
			Build()

	// Inlines variables that are marked as constants.
	c.inlineConstants =
		PassFactory_Builder_().
			SetName("inlineConstants").
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewInlineVariables(compiler, InlineVariables_Mode_CONSTANTS_ONLY)
				}).
			Build()

	// Inlines simple methods, like getters
	c.inlineSimpleMethods =
		PassFactory_Builder_().
			SetName("inlineSimpleMethods").
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewInlineSimpleMethods(compiler)
			}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Kills dead assignments.
	c.deadAssignmentsElimination =
		PassFactory_Builder_().
			SetName(PassNames_DEAD_ASSIGNMENT_ELIMINATION).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewDeadAssignmentsElimination(compiler)
			}).
			Build()

	// Kills dead property assignments.
	c.deadPropertyAssignmentElimination =
		PassFactory_Builder_().
			SetName("deadPropertyAssignmentElimination").
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewDeadPropertyAssignmentElimination(compiler)
			}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Inlines function calls.
	c.inlineFunctions =
		PassFactory_Builder_().
			SetName(PassNames_INLINE_FUNCTIONS).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewInlineFunctions(
						compiler,
						compiler.GetUniqueNameIdSupplier(),
						options.GetInlineFunctionsLevel(),
						options.AssumeStrictThis() || options.ExpectStrictModeInput(),
						options.AssumeClosuresOnlyCaptureReferences(),
						options.GetMaxFunctionSizeAfterInlining())
				}).
			Build()

	// Inlines constant properties.
	c.inlineProperties =
		PassFactory_Builder_().
			SetName(PassNames_INLINE_PROPERTIES).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewInlineProperties(compiler)
			}).
			Build()

	// Isolates injected polyfills & references from the global scope
	c.isolatePolyfills =
		PassFactory_Builder_().
			SetName("IsolatePolyfills").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewIsolatePolyfills(compiler)
			}).
			Build()

	c.removeUnusedCode =
		PassFactory_Builder_().
			SetName(PassNames_REMOVE_UNUSED_CODE).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewRemoveUnusedCode_Builder(compiler).
						RemoveLocalVars(options.ShouldRemoveUnusedLocalVariables()).
						RemoveGlobals(options.ShouldRemoveUnusedVariables()).
						PreserveFunctionExpressionNames(false).
						RemoveUnusedPrototypeProperties(
							options.ShouldRemoveUnusedPrototypeProperties()).
						RemoveUnusedThisProperties(options.IsRemoveUnusedClassProperties()).
						RemoveUnusedObjectDefinePropertiesDefinitions(
							options.IsRemoveUnusedClassProperties()).
						// If we are forcing injection of some library code, don't remove polyfills.
						// Otherwise, we might end up removing polyfills the user specifically asked
						// to include.
						RemoveUnusedPolyfills(
							len(options.GetForceLibraryInjectionList()) == 0 &&
								options.GetInjectPolyfillsNewerThan() == nil).
						AssumeGettersArePure(options.GetAssumeGettersArePure()).
						Build()
				}).
			Build()

	c.removeUnusedCodeOnce =
		c.removeUnusedCode.ToBuilder().SetRunInFixedPointLoop(false).Build()

	// Move global symbols to a deeper common chunk
	c.crossChunkCodeMotion =
		PassFactory_Builder_().
			SetName(PassNames_CROSS_CHUNK_CODE_MOTION).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCrossChunkCodeMotion(
						compiler,
						compiler.GetChunkGraph(),
						options.GetParentChunkCanSeeSymbolsDeclaredInChildren())
				}).
			Build()

	// Move methods to a deeper common chunk
	c.crossChunkMethodMotion =
		PassFactory_Builder_().
			SetName(PassNames_CROSS_CHUNK_METHOD_MOTION).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCrossChunkMethodMotion(
						compiler,
						compiler.GetCrossChunkIdGenerator(),
						/* canModifyExterns= */ false, // remove this
						options.GetCrossChunkCodeMotionNoStubMethods())
				}).
			Build()

	// A data-flow based variable inliner.
	c.flowSensitiveInlineVariables =
		PassFactory_Builder_().
			SetName(PassNames_FLOW_SENSITIVE_INLINE_VARIABLES).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewFlowSensitiveInlineVariables(compiler)
			}).
			Build()

	// Lowers block-scoped declarations to function-scoped declarations.

	// Uses register-allocation algorithms to use fewer variables.
	c.coalesceVariableNames =
		PassFactory_Builder_().
			SetName(PassNames_COALESCE_VARIABLE_NAMES).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCoalesceVariableNames(compiler, options.ShouldGeneratePseudoNames())
				}).
			Build()

	// Collapses assignment expressions (e.g., `x = 3; y = x;` becomes
	// `y = x = 3;`.
	c.postNormalizePeephole =
		PassFactory_Builder_().
			SetName(PassNames_POST_NORMALIZE_PEEPHOLE).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					var optimizations []AbstractPeepholeOptimization
					if options.ShouldCollapseVariableDeclarations() {
						optimizations = append(optimizations, NewExploitAssigns())
					}
					if options.ShouldOptimizeLetAndConst() {
						optimizations = append(optimizations,
							NewOptimizeLetAndConstPeephole(options.ShouldTreatGlobalScopeAsIsolated()))
					}
					return NewPeepholeOptimizationsPassOfList(
						compiler, PassNames_POST_NORMALIZE_PEEPHOLE, optimizations)
				}).
			Build()

	// Collapses variable declarations (e.g., `var x; var y;` becomes
	// `var x,y;`.
	c.collapseVariableDeclarations =
		PassFactory_Builder_().
			SetName(PassNames_COLLAPSE_VARIABLE_DECLARATIONS).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewCollapseVariableDeclarations(compiler)
			}).
			Build()

	// Extracts common sub-expressions.
	c.extractPrototypeMemberDeclarations =
		PassFactory_Builder_().
			SetName(PassNames_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					var pattern ExtractPrototypeMemberDeclarations_Pattern
					switch options.GetExtractPrototypeMemberDeclarationsMode() {
					case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP:
						pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP
					case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP:
						pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP
					case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE:
						pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE
					default:
						panic(jutil.NewIllegalStateExceptionMsg("unexpected"))
					}

					return NewExtractPrototypeMemberDeclarations(compiler, pattern)
				}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Rewrites common function definitions to be more compact.
	c.rewriteFunctionExpressions =
		PassFactory_Builder_().
			SetName(PassNames_REWRITE_FUNCTION_EXPRESSIONS).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewFunctionRewriter(compiler)
			}).
			Build()

	// Collapses functions to not use the VAR keyword.
	c.collapseAnonymousFunctions =
		PassFactory_Builder_().
			SetName(PassNames_COLLAPSE_ANONYMOUS_FUNCTIONS).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewCollapseAnonymousFunctions(compiler)
			}).
			Build()

	// Moves function declarations to the top to simulate actual hoisting and
	// rewrites block scope declarations to 'var'.
	c.rewriteGlobalDeclarationsForTryCatchWrapping =
		PassFactory_Builder_().
			SetName("rewriteGlobalDeclarationsForTryCatchWrapping").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewRewriteGlobalDeclarationsForTryCatchWrapping(compiler)
			}).
			Build()

	// Alias string literals with global variables, to reduce code size.
	c.aliasStrings =
		PassFactory_Builder_().
			SetName("aliasStrings").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewAliasStrings(
						compiler,
						compiler.GetChunkGraph(),
						options.ShouldOutputJsStringUsage(),
						options.GetAliasStringsMode())
				}).
			Build()

	// Renames properties so that the two properties that never appear on the
	// same object get the same name.
	c.ambiguateProperties =
		PassFactory_Builder_().
			SetName(PassNames_AMBIGUATE_PROPERTIES).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewAmbiguateProperties(
						compiler,
						options.GetPropertyReservedNamingFirstChars(),
						options.GetPropertyReservedNamingNonFirstChars(),
						compiler.GetExternProperties())
				}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Mark the point at which the normalized AST assumptions no longer hold.
	c.markUnnormalized =
		PassFactory_Builder_().
			SetName("markUnnormalized").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
					})
				}).
			Build()

	c.normalize =
		PassFactory_Builder_().
			SetName(PassNames_NORMALIZE).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return Normalize_CreateNormalizeForOptimizations(compiler)
			}).
			Build()

	c.externExports =
		PassFactory_Builder_().
			SetName(PassNames_EXTERN_EXPORTS).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewExternExportsPass(compiler)
			}).
			Build()

	// Denormalize the AST for code generation.
	c.denormalize =
		PassFactory_Builder_().
			SetName("denormalize").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewDenormalize(compiler, options.GetOutputFeatureSet())
				}).
			Build()

	// Inverting name normalization.
	c.invertContextualRenaming =
		PassFactory_Builder_().
			SetName("invertContextualRenaming").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return MakeDeclaredNamesUnique_GetContextualRenameInverter(compiler)
			}).
			Build()

	// Renames properties.
	c.renameProperties =
		PassFactory_Builder_().
			SetName("renameProperties").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					jutil.CheckState(options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED)
					prevPropertyMap := options.GetInputPropertyMap()
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						rprop :=
							NewRenameProperties(
								compiler,
								options.ShouldGeneratePseudoNames(),
								prevPropertyMap,
								options.GetPropertyReservedNamingFirstChars(),
								options.GetPropertyReservedNamingNonFirstChars(),
								options.GetNameGenerator())
						rprop.Process(externs, root)
						compiler.SetPropertyMap(rprop.GetPropertyMap())
					})
				}).
			SetPreconditionCheck(defaultPassConfig_requirePropertiesAreStaticallyAnalyzable).
			Build()

	// Inlines calls to property renaming functions.
	c.removePropertyRenamingCalls =
		PassFactory_Builder_().
			SetName("removePropertyRenamingCalls").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewRemovePropertyRenamingCalls(compiler)
			}).
			Build()

	// Renames variables.
	c.renameVars =
		PassFactory_Builder_().
			SetName("renameVars").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					prevVariableMap := options.GetInputVariableMap()
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						compiler.SetVariableMap(
							c.runVariableRenaming(compiler, prevVariableMap, externs, root))
					})
				}).
			Build()

	// Renames labels
	c.renameLabels =
		PassFactory_Builder_().SetName("renameLabels").SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
			return NewRenameLabels(compiler)
		}).Build()

	// Convert bracket access to dot access
	c.convertToDottedProperties =
		PassFactory_Builder_().
			SetName(PassNames_CONVERT_TO_DOTTED_PROPERTIES).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewConvertToDottedProperties(compiler)
			}).
			Build()

	c.checkAstValidity =
		PassFactory_Builder_().
			SetName("checkAstValidity").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewAstValidator(compiler)
			}).
			Build()

	// Checks that all variables are defined.
	c.varCheckValidity =
		PassFactory_Builder_().
			SetName("varCheckValidity").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass { return NewVarCheckWithValidityCheck(compiler, true) }).
			Build()

	c.instrumentForCodeCoverage =
		PassFactory_Builder_().
			SetName("instrumentForCodeCoverage").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					// TODO(johnlenz): make global instrumentation an option
					return NewCoverageInstrumentationPassWithInstrumentOption(
						compiler,
						CoverageInstrumentationPass_CoverageReach_CONDITIONAL,
						options.GetInstrumentForCoverageOption(),
						options.GetProductionInstrumentationArrayName())
				}).
			Build()

	c.gatherExternPropertiesCheck =
		c.createGatherExternProperties(GatherExternProperties_Mode_CHECK)

	c.gatherExternPropertiesOptimize =
		c.createGatherExternProperties(GatherExternProperties_Mode_OPTIMIZE)

	// Rewrites Polymer({})
	c.polymerPass =
		PassFactory_Builder_().SetName("polymerPass").SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
			return NewPolymerPass(compiler)
		}).Build()

	c.chromePass =
		PassFactory_Builder_().SetName("chromePass").SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
			return NewChromePass(compiler)
		}).Build()

	// Rewrites J2CL constructs to be more optimizable.
	c.j2clConstantHoisterPass =
		PassFactory_Builder_().
			SetName("j2clConstantHoisterPass").
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewJ2clConstantHoisterPass(compiler)
			}).
			Build()

	// Optimizes J2CL clinit methods.
	c.j2clClinitPass =
		PassFactory_Builder_().
			SetName("j2clClinitPass").
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					changedScopeNodes :=
						compiler.GetChangeTracker().GetChangedScopeNodesForPass("j2clClinitPass")
					return NewJ2clClinitPrunerPass(compiler, changedScopeNodes)
				}).
			Build()

	// Rewrites J2CL constructs to be more optimizable.
	c.j2clPropertyInlinerPass =
		PassFactory_Builder_().
			SetName("j2clPropertyInlinerPass").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewJ2clPropertyInlinerPass(compiler)
			}).
			Build()

	// Rewrites J2CL constructs to be more optimizable.
	c.j2clPass =
		PassFactory_Builder_().SetName("j2clPass").SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
			return NewJ2clPass(compiler)
		}).Build()

	// Rewrites J2CL constructs to be more optimizable.
	c.j2clUtilGetDefineRewriterPass =
		PassFactory_Builder_().
			SetName("j2clUtilGetDefineRewriterPass").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewJ2clUtilGetDefineRewriterPass(compiler)
			}).
			Build()

	c.j2clAssertRemovalPass =
		PassFactory_Builder_().
			SetName("j2clAssertRemovalPass").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewJ2clAssertRemovalPass(compiler)
			}).
			Build()

	c.j2clSourceFileChecker =
		PassFactory_Builder_().
			SetName("j2clSourceFileChecker").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewJ2clSourceFileChecker(compiler)
			}).
			Build()

	c.j2clChecksPass =
		PassFactory_Builder_().
			SetName("j2clChecksPass").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewJ2clChecksPass(compiler)
			}).
			Build()

	c.checkConformance =
		PassFactory_Builder_().
			SetName(PassNames_CHECK_CONFORMANCE).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					reportingMode := options.GetConformanceReportingMode()
					return NewCheckConformanceWithConfigs(
						compiler,
						options.GetConformanceConfigs(),
						&reportingMode)
				}).
			Build()

	c.removeCastNodes =
		PassFactory_Builder_().
			SetName("removeCastNodes").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewRemoveCastNodes(compiler)
			}).
			Build()

	// Convert types from type-checking to optimization "colors"
	c.typesToColors =
		PassFactory_Builder_().
			SetName("typesToColors").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, js *rhino.Node) {
						NewConvertTypesToColors(
							compiler,
							SerializationOptions_Builder_().
								SetIncludeDebugInfo(
									compiler.GetOptions().ShouldSerializeExtraDebugInfo()).
								Build()).
							Process(externs, js)

						compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
					})
				}).
			Build()

	// Emit a TypedAST into of the current compilation.
	c.serializeTypedAst =
		PassFactory_Builder_().
			SetName("serializeTypedAst").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return SerializeTypedAstPass_CreateFromPath(
						compiler,
						options.GetTypedAstOutputFile(),
						SerializationOptions_Builder_().
							SetIncludeDebugInfo(
								compiler.GetOptions().ShouldSerializeExtraDebugInfo()).
							// set the runtime libraries to serialize in the TypedAST proto
							SetRuntimeLibraries(
								compiler.GetRuntimeJsLibManager().GetInjectedLibraries()).
							Build())
				}).
			Build()

	c.removeUnnecessarySyntheticExterns =
		PassFactory_Builder_().
			SetName("removeUnnecessarySyntheticExterns").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewRemoveUnnecessarySyntheticExterns(compiler)
			}).
			Build()

	// Optimizations that output ES2015 features.
	c.optimizeToEs6 =
		PassFactory_Builder_().
			SetName("optimizeToEs6").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewSubstituteEs6Syntax(compiler)
			}).
			Build()

	// Rewrites goog.module in whitespace only mode
	c.whitespaceWrapGoogModules =
		PassFactory_Builder_().
			SetName("whitespaceWrapGoogModules").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewWhitespaceWrapGoogModules(compiler)
			}).
			Build()

	c.rewriteCommonJsModules =
		PassFactory_Builder_().
			SetName(PassNames_REWRITE_COMMON_JS_MODULES).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewProcessCommonJSModules(compiler)
			}).
			Build()

	c.rewriteScriptsToEs6Modules =
		PassFactory_Builder_().
			SetName(PassNames_REWRITE_SCRIPTS_TO_ES6_MODULES).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewEs6RewriteScriptsToModules(compiler)
			}).
			Build()

	c.gatherModuleMetadataPass =
		PassFactory_Builder_().
			SetName(PassNames_GATHER_MODULE_METADATA).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					// Force creation of the synthetic input so that we create metadata for it
					compiler.GetSynthesizedExternsInput()
					return NewGatherModuleMetadata(
						compiler,
						options.GetProcessCommonJSModules(),
						options.GetModuleResolutionMode())
				}).
			Build()

	c.createModuleMapPass =
		PassFactory_Builder_().
			SetName(PassNames_CREATE_MODULE_MAP).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap())
				}).
			// does not look at AST
			Build()

	c.gatherGettersAndSetters =
		PassFactory_Builder_().
			SetName(PassNames_GATHER_GETTERS_AND_SETTERS).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewGatherGetterAndSetterProperties(compiler)
			}).
			Build()

	// this pass is just adding script without looking at the AST so it should accept all features
	c.addSyntheticScript =
		PassFactory_Builder_().
			SetName("ADD_SYNTHETIC_SCRIPT").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, js *rhino.Node) { compiler.InitializeSyntheticCodeInput() })
				}).
			Build()

	c.removeSyntheticScript =
		PassFactory_Builder_().
			SetName("REMOVE_SYNTHETIC_SCRIPT").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return CompilerPassFunc(func(externs *rhino.Node, js *rhino.Node) { compiler.RemoveSyntheticCodeInput() })
			}).
			Build()

	c.mergeSyntheticScript =
		PassFactory_Builder_().
			SetName("MERGE_SYNTHETIC_SCRIPT").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return CompilerPassFunc(func(externs *rhino.Node, js *rhino.Node) { compiler.MergeSyntheticCodeInput() })
			}).
			Build()

	// Replace locale values with stubs that can survive optimizations and be
	// replaced afterwards.
	c.protectLocaleData =
		PassFactory_Builder_().
			SetName("protectLocaleData").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						NewLocaleDataPasses_ProtectGoogLocale(compiler).Process(externs, root)
					})
				}).
			Build()

	// Replace locale data stubs with the values for a locale.
	c.substituteLocaleData =
		PassFactory_Builder_().
			SetName("SubstituteLocaleData").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
						NewLocaleDataPasses_LocaleSubstitutions(
							compiler, compiler.GetOptions().GetLocale()).
							Process(externs, root)
					})
				}).
			Build()

	c.transpileOnlyClosureUnaware =
		PassFactory_Builder_().
			SetName("TranspileOnlyClosureUnaware").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewTranspileAndOptimizeClosureUnaware(
						compiler, NestedCompilerRunner_Mode_TRANSPILE_ONLY)
				}).
			Build()

	c.transpileAndOptimizeClosureUnaware =
		PassFactory_Builder_().
			SetName("TranspileAndOptimizeClosureUnaware").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewTranspileAndOptimizeClosureUnaware(
						compiler, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE)
				}).
			Build()
}

// Java: DefaultPassConfig.getPreprocessorSymbolTable()
func (c *DefaultPassConfig) GetPreprocessorSymbolTable() *PreprocessorSymbolTable {
	return c.preprocessorSymbolTableFactory.GetInstanceOrNull()
}

// Java: DefaultPassConfig.getTranspileOnlyPasses()
func (c *DefaultPassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	// Certain errors in block-scoped variable declarations will prevent correct transpilation
	passes.MaybeAdd(c.checkVariableReferences)
	passes.MaybeAdd(c.checkVars)

	passes.MaybeAdd(c.gatherModuleMetadataPass)
	passes.MaybeAdd(c.createModuleMapPass)

	if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		passes.MaybeAdd(c.rewriteGoogJsImports)
		switch options.GetEs6ModuleTranspilation() {
		case CompilerOptions_Es6ModuleTranspilation_COMPILE:
			TranspilationPasses_AddEs6ModulePass(passes, c.preprocessorSymbolTableFactory)
		case CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES:
			TranspilationPasses_AddEs6ModuleToCjsPass(passes)
		case CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS:
			TranspilationPasses_AddEs6RewriteImportPathPass(passes)
		case CompilerOptions_Es6ModuleTranspilation_NONE:
			// nothing
		}
	}

	passes.MaybeAdd(c.checkSuper)

	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)

	passes.MaybeAdd(c.transpileOnlyClosureUnaware)

	if options.NeedsTranspilationFrom(parser.FeatureSet_ES2015) && options.GetRewritePolyfills() {
		if options.GetIsolatePolyfills() {
			panic(jutil.NewIllegalStateExceptionMsg("Polyfill isolation cannot be used in transpileOnly mode"))
		}
		TranspilationPasses_AddRewritePolyfillPass(passes)
	} else if options.GetInjectPolyfillsNewerThan() != nil {
		TranspilationPasses_AddRewritePolyfillPass(passes)
	}

	passes.MaybeAdd(c.injectRuntimeLibraries)

	// Passes below this point may rely on normalization and must maintain normalization.
	passes.MaybeAdd(c.normalize)

	passes.MaybeAdd(c.gatherGettersAndSetters)

	TranspilationPasses_AddTranspilationPasses(passes, options)

	passes.MaybeAdd(c.markUnnormalized)
	// The transpilation passes may rely on normalize making all variables unique,
	// but we're doing only transpilation, so we want to put back the original variable names
	// wherever we can to meet user expectations.
	//
	// Also, if we don't do this, we could end up breaking runtime behavior that depends on specific
	// variable names.
	//
	// The primary concern is function parameter names, because some frameworks, like Angular,
	// do runtime injection of function call arguments based on the function parameter names.
	passes.MaybeAdd(c.invertContextualRenaming)

	// Es6ConvertSuper may add this so it needs to be removed for transpile-only mode.
	passes.MaybeAdd(c.removePropertyRenamingCalls)

	passes.AssertAllOneTimePasses()
	c.assertValidOrderForChecks(passes)
	return passes
}

// Java: DefaultPassConfig.getWhitespaceOnlyPasses()
func (c *DefaultPassConfig) GetWhitespaceOnlyPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.GetProcessCommonJSModules() {
		passes.MaybeAdd(c.rewriteCommonJsModules)
	} else if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		passes.MaybeAdd(c.rewriteScriptsToEs6Modules)
	}

	if options.ShouldWrapGoogModulesForWhitespaceOnly() {
		passes.MaybeAdd(c.whitespaceWrapGoogModules)
	}
	return passes
}

// Java: DefaultPassConfig.addModuleRewritingPasses(PassListBuilder, CompilerOptions)
func (c *DefaultPassConfig) addModuleRewritingPasses(checks *PassListBuilder, options *CompilerOptions) {
	if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		checks.MaybeAdd(c.rewriteGoogJsImports)
		TranspilationPasses_AddEs6ModulePass(checks, c.preprocessorSymbolTableFactory)
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(c.closureRewriteModule)
	}
}

// Java: DefaultPassConfig.getChecks()
func (c *DefaultPassConfig) GetChecks() *PassListBuilder {
	options := c.options
	checks := NewPassListBuilder(options)

	jutil.CheckStateMsg(
		!options.GetSkipNonTranspilationPasses(),
		"options.getSkipNonTranspilationPasses() cannot be mixed with PassConfig::getChecks. Call"+
			" PassConfig::getTranspileOnlyPasses instead.")

	if options.IsPropertyRenamingOnlyCompilationMode() {
		checks.MaybeAdd(c.addSyntheticScript)
		checks.MaybeAdd(c.gatherGettersAndSetters)
		checks.MaybeAdd(c.gatherModuleMetadataPass)
		checks.MaybeAdd(c.createModuleMapPass)
		checks.MaybeAdd(c.declaredGlobalExternsOnWindow)
		checks.MaybeAdd(c.checkSideEffects)
		checks.MaybeAdd(c.angularPass)
		checks.MaybeAdd(c.closureGoogScopeAliases)
		checks.MaybeAdd(c.closurePrimitives)
		c.addModuleRewritingPasses(checks, options)
		checks.MaybeAdd(c.clearTypedScopeCreatorPass)
		checks.MaybeAdd(c.clearTopTypedScopePass)
		checks.MaybeAdd(c.generateExports)
		checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_AFTER_STANDARD_CHECKS))
		checks.MaybeAdd(c.mergeSyntheticScript)
		checks.MaybeAdd(c.gatherExternPropertiesCheck)
		checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_SERIALIZATION))
		return checks
	}

	if options.ShouldGenerateTypedExterns() {
		checks.MaybeAdd(c.addSyntheticScript)
		checks.MaybeAdd(c.closureGoogScopeAliasesForIjs)
		checks.MaybeAdd(c.generateIjs)
		checks.MaybeAdd(c.removeExtraRequires)
		if options.ShouldWrapGoogModulesForWhitespaceOnly() {
			checks.MaybeAdd(c.whitespaceWrapGoogModules)
		}
		checks.MaybeAdd(c.removeSyntheticScript)
		return checks
	}

	// Run this pass before any pass tries to inject new runtime libraries
	checks.MaybeAdd(c.addSyntheticScript)

	checks.MaybeAdd(c.gatherGettersAndSetters)

	if options.GetLanguageIn().ToFeatureSet().ContainsFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT) &&
		!options.ShouldAllowDynamicImport() {
		checks.MaybeAdd(defaultPassConfig_forbidDynamicImportUsage)
	}

	checks.MaybeAdd(PassFactory_CreateEmptyPass("beforeStandardChecks"))

	if !options.GetProcessCommonJSModules() &&
		options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		checks.MaybeAdd(c.rewriteScriptsToEs6Modules)
	}

	// Run these passes after promoting scripts to modules, but before rewriting any other modules.
	checks.MaybeAdd(c.gatherModuleMetadataPass)
	checks.MaybeAdd(c.createModuleMapPass)

	if options.GetProcessCommonJSModules() {
		checks.MaybeAdd(c.rewriteCommonJsModules)
	}

	// Note: ChromePass can rewrite invalid @type annotations into valid ones, so should run before
	// JsDoc checks.
	if options.IsChromePassEnabled() {
		checks.MaybeAdd(c.chromePass)
	}

	// Verify JsDoc annotations and check ES modules
	checks.MaybeAdd(c.checkJsDocAndEs6Modules)

	checks.MaybeAdd(c.checkTypeImportCodeReferences)

	if options.Enables(DiagnosticGroups_LINT_CHECKS) {
		checks.MaybeAdd(c.lintChecks)
	}

	if options.GetClosurePass() && options.Enables(DiagnosticGroups_LINT_CHECKS) {
		checks.MaybeAdd(c.checkRequiresAndProvidesSorted)
	}

	if options.Enables(DiagnosticGroups_EXTRA_REQUIRE) {
		checks.MaybeAdd(c.extraRequires)
	}

	if options.Enables(DiagnosticGroups_MISSING_REQUIRE) {
		checks.MaybeAdd(c.checkMissingRequires)
	}

	checks.MaybeAdd(c.declaredGlobalExternsOnWindow)

	if !options.GetProcessCommonJSModules() {
		// TODO(ChadKillingsworth): move CommonJS module rewriting after VarCheck
		checks.MaybeAdd(c.checkVariableReferences)
		checks.MaybeAdd(c.checkVars)
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(c.checkClosureImports)
	}

	checks.MaybeAdd(c.checkStrictMode)

	if options.GetClosurePass() {
		checks.MaybeAdd(c.closureCheckModule)
	}

	checks.MaybeAdd(c.checkSuper)

	checks.MaybeAdd(c.checkSideEffects)

	if options.GetAngularPass() {
		checks.MaybeAdd(c.angularPass)
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(c.closureGoogScopeAliases)
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(c.closurePrimitives)
	}

	// TODO(b/141389184): Move this after the Polymer pass
	if options.ShouldRewriteModulesBeforeTypechecking() {
		c.addModuleRewritingPasses(checks, options)
	}

	// It's important that the PolymerPass run *after* the ClosurePrimitives and ChromePass rewrites
	// and *before* the suspicious code checks. This is enforced in the assertValidOrder method.
	if options.IsPolymerPassEnabled() {
		checks.MaybeAdd(c.polymerPass)
	}

	if options.GetProcessCommonJSModules() {
		// TODO(ChadKillingsworth): remove this branch.
		checks.MaybeAdd(c.checkVariableReferences)
		checks.MaybeAdd(c.checkVars)
	}

	if options.ShouldInferConsts() {
		checks.MaybeAdd(c.inferConsts)
	}

	if options.ShouldComputeFunctionSideEffects() {
		checks.MaybeAdd(c.checkRegExp)
	}

	checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_TYPE_CHECKING))

	if options.GetCheckTypes() || options.GetInferTypes() {
		checks.MaybeAdd(c.inferTypes)
		if options.GetCheckTypes() {
			checks.MaybeAdd(c.checkTypes)
		} else {
			checks.MaybeAdd(c.inferJsDocInfo)
		}
	}

	// Analyzer checks must be run after typechecking but before module rewriting.
	if options.Enables(DiagnosticGroups_ANALYZER_CHECKS) && options.IsTypecheckingEnabled() {
		checks.MaybeAdd(c.analyzerChecks)
	}

	// We assume that only clients who are going to re-compile, or do in-depth static analysis,
	// will need the typed scope creator after the compile job.
	if !options.PreservesDetailedSourceInfo() {
		checks.MaybeAdd(c.clearTypedScopeCreatorPass)
	}

	if options.ShouldRewriteModulesAfterTypechecking() {
		c.addModuleRewritingPasses(checks, options)
	}

	// Dynamic import rewriting must always occur after type checking
	if options.ShouldAllowDynamicImport() &&
		options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_DYNAMIC_IMPORT) {
		checks.MaybeAdd(defaultPassConfig_rewriteDynamicImports)
	}

	// We assume that only clients who are going to re-compile, or do in-depth static analysis,
	// will need the typed scope creator after the compile job.
	if !options.PreservesDetailedSourceInfo() {
		checks.MaybeAdd(c.clearTopTypedScopePass)
	}

	// CheckSuspiciousCode requires type information, so must run after the type checker.
	if options.GetCheckSuspiciousCode() ||
		options.Enables(DiagnosticGroups_GLOBAL_THIS) ||
		options.Enables(DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT) {
		checks.MaybeAdd(c.suspiciousCode)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		checks.MaybeAdd(c.j2clSourceFileChecker)
	}

	if !options.Disables(DiagnosticGroups_CHECK_USELESS_CODE) ||
		!options.Disables(DiagnosticGroups_MISSING_RETURN) {
		checks.MaybeAdd(c.checkControlFlow)
	}

	// CheckAccessControls only works if check types is on.
	if options.IsTypecheckingEnabled() &&
		(!options.Disables(DiagnosticGroups_ACCESS_CONTROLS) ||
			options.Enables(DiagnosticGroups_CONSTANT_PROPERTY)) {
		checks.MaybeAdd(c.checkAccessControls)
	}

	checks.MaybeAdd(c.checkConsts)

	checks.MaybeAdd(c.rewriteCallerCodeLocation)

	if len(options.GetConformanceConfigs()) != 0 {
		checks.MaybeAdd(c.checkConformance)
	}

	if options.GetTweakProcessing().IsOn() {
		checks.MaybeAdd(c.processTweaks)
	}

	checks.MaybeAdd(c.processDefinesCheck)

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		checks.MaybeAdd(c.j2clChecksPass)
	}

	if options.ShouldGenerateExports() {
		checks.MaybeAdd(c.generateExports)
	}

	checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_AFTER_STANDARD_CHECKS))

	checks.MaybeAdd(c.mergeSyntheticScript)

	// Create extern exports after the normalize because externExports depends on unique names.
	if options.GetExternExportsPath() != nil {
		checks.MaybeAdd(c.externExports)
	}

	if !options.IsChecksOnly() {
		checks.MaybeAdd(c.removeWeakSources)
	}

	// Gather property names in externs so they can be queried by the optimizing passes.
	// See b/180424427 for why this runs in stage 1 and not stage 2.
	checks.MaybeAdd(c.gatherExternPropertiesCheck)

	checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_SERIALIZATION))

	if options.GetTypedAstOutputFile() != nil {
		checks.MaybeAdd(c.serializeTypedAst)
	}

	// Validation should be after all pass additions
	checks.AssertAllOneTimePasses()
	c.assertValidOrderForChecks(checks)

	return checks
}

// Java: DefaultPassConfig.getOptimizations()
func (c *DefaultPassConfig) GetOptimizations() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.IsPropertyRenamingOnlyCompilationMode() {
		passes.MaybeAdd(c.removeUnnecessarySyntheticExterns)
		TranspilationPasses_AddTranspilationRuntimeLibraries(passes)
		passes.MaybeAdd(c.closureProvidesRequires)
		passes.MaybeAdd(c.processDefinesOptimize)
		passes.MaybeAdd(c.normalize)
		passes.MaybeAdd(c.gatherGettersAndSetters)
		TranspilationPasses_AddTranspilationPasses(passes, options)
		passes.MaybeAdd(c.gatherExternPropertiesOptimize)
		passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_STANDARD_OPTIMIZATIONS))
		passes.MaybeAdd(c.inlineAndCollapseProperties)
		passes.MaybeAdd(c.closureOptimizePrimitives)
		// If side-effects were protected, remove the protection now.
		passes.MaybeAdd(c.stripSideEffectProtection)
		return passes
	}

	if options.GetSkipNonTranspilationPasses() {
		// Reaching this if-condition means the 'getChecks()' phase has been skipped in favor of
		// 'getTranspileOnlyPasses'.
		return passes
	}

	passes.AddAll(c.getEarlyOptimizationPasses())
	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_OPTIMIZATIONS_HALFWAY_POINT))
	passes.AddAll(c.getLateOptimizationPasses())

	return passes
}

// Java: DefaultPassConfig.getFinalizations()
func (c *DefaultPassConfig) GetFinalizations() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.IsPropertyRenamingOnlyCompilationMode() {
		if options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() ||
			options.GetRenamePrefixNamespace() != nil {
			passes.MaybeAdd(c.rewriteGlobalDeclarationsForTryCatchWrapping)
		}

		passes.MaybeAdd(c.renameProperties)

		if options.GetRenamePrefixNamespace() != nil &&
			options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE {
			if !defaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN.Matches(*options.GetRenamePrefixNamespace()) {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"Illegal character in renamePrefixNamespace name: " +
						*options.GetRenamePrefixNamespace()))
			}
			passes.MaybeAdd(c.rescopeGlobalSymbols)
		}
		return passes
	}

	if options.DoLateLocalization() {
		if options.ShouldRunReplaceMessagesPass() {
			passes.MaybeAdd(c.getReplaceProtectedMessagesPass())
		}
		passes.MaybeAdd(c.substituteLocaleData)
		passes.MaybeAdd(c.peepholeOptimizationsOnceNormalized)
	}

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(c.flowSensitiveInlineVariables)

		// After inlining variable uses, some variables may be unused.
		// If we're doing late localization, the simple code removal pass runs added below will clean
		// those up. Otherwise, clean them up now.
		if !options.DoLateLocalization() && c.shouldRunRemoveUnusedCode() {
			passes.MaybeAdd(c.removeUnusedCodeOnce)
		}
	}

	if options.DoLateLocalization() {
		passes.AddAll(c.getPostL10nOptimizations())
	}

	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeChunkMotion"))

	if options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(c.crossChunkCodeMotion)
	}

	if options.ShouldRunCrossChunkMethodMotion() {
		passes.MaybeAdd(c.crossChunkMethodMotion)
	}

	passes.MaybeAdd(PassFactory_CreateEmptyPass("afterChunkMotion"))

	if options.GetOptimizeESClassConstructors() &&
		options.GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015) {
		passes.MaybeAdd(c.optimizeConstructors)
	}

	// Isolate injected polyfills from the global scope. Runs late in the optimization loop
	// to take advantage of property renaming & RemoveUnusedCode, as this pass will increase code
	// size by wrapping all potential polyfill usages.
	if options.GetIsolatePolyfills() {
		passes.MaybeAdd(c.isolatePolyfills)
	}

	if options.ShouldCollapseAnonymousFunctions() {
		// TODO: b/197349249 - Maybe we should just move this pass after denormalization. It seems
		// weird to convert from function expression to function declaration while we're still
		// supposed to be in a normalized state. But it requires testing as perhaps some optimizations
		// in that range will get affected if we skip this rewriting.
		passes.MaybeAdd(c.collapseAnonymousFunctions)
	}

	// Move functions before extracting prototype member declarations.
	if options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() ||
		// renamePrefixNamescape relies on rewriteGlobalDeclarationsForTryCatchWrapping
		// to preserve semantics.
		options.GetRenamePrefixNamespace() != nil {
		passes.MaybeAdd(c.rewriteGlobalDeclarationsForTryCatchWrapping)
	}

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS))
	// The mapped name anonymous function pass makes use of information that
	// the extract prototype member declarations pass removes so the former
	// happens before the latter.
	if options.GetExtractPrototypeMemberDeclarationsMode() !=
		CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF {
		passes.MaybeAdd(c.extractPrototypeMemberDeclarations)
	}

	if options.ShouldAmbiguateProperties() &&
		options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED &&
		options.IsTypecheckingEnabled() {
		passes.MaybeAdd(c.ambiguateProperties)
	}

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_RENAME_PROPERTIES))
	if options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED {
		passes.MaybeAdd(c.renameProperties)
	} else {
		passes.MaybeAdd(c.removePropertyRenamingCalls)
	}

	// Reserve global names added to the "windows" object.
	if options.ShouldReserveRawExports() {
		passes.MaybeAdd(c.gatherRawExports)
	}

	// This comes after property renaming because quoted property names must
	// not be renamed.
	if options.ShouldConvertToDottedProperties() {
		passes.MaybeAdd(c.convertToDottedProperties)
	}

	// Property renaming must happen before this pass runs since this
	// pass may convert dotted properties into quoted properties.  It
	// is beneficial to run before alias strings, alias keywords and
	// variable renaming.
	if options.ShouldRewriteFunctionExpressions() {
		passes.MaybeAdd(c.rewriteFunctionExpressions)
	}

	// This comes after converting quoted property accesses to dotted property
	// accesses in order to avoid aliasing property names.
	if options.GetAliasStringsMode() != CompilerOptions_AliasStringsMode_NONE {
		passes.MaybeAdd(c.aliasStrings)
	}

	if options.ShouldCoalesceVariableNames() {
		// Passes after this point can no longer depend on normalized AST
		// assumptions because the code is marked as un-normalized
		passes.MaybeAdd(c.coalesceVariableNames)

		// coalesceVariables creates identity assignments and more redundant code
		// that can be removed, rerun the peephole optimizations to clean them
		// up.
		if options.ShouldFoldConstants() {
			passes.MaybeAdd(c.peepholeOptimizationsOnceNonNormalized)
		}
	}

	// Passes after this point can no longer depend on normalized AST assumptions.
	passes.MaybeAdd(c.markUnnormalized)

	if options.ShouldCollapseVariableDeclarations() || options.ShouldOptimizeLetAndConst() {
		passes.MaybeAdd(c.postNormalizePeephole)
	}
	if options.ShouldCollapseVariableDeclarations() {
		passes.MaybeAdd(c.collapseVariableDeclarations)
	}

	// This pass works best after collapseVariableDeclarations.
	passes.MaybeAdd(c.denormalize)

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_VARIABLE_RENAMING))

	if options.GetVariableRenaming() != VariableRenamingPolicy_ALL {
		// If we're leaving some (or all) variables with their old names,
		// then we need to undo any of the markers we added for distinguishing
		// local variables ("x" -> "x$jscomp$1").
		passes.MaybeAdd(c.invertContextualRenaming)
	}

	if options.GetVariableRenaming() != VariableRenamingPolicy_OFF {
		passes.MaybeAdd(c.renameVars)
	}

	if options.ShouldRenameLabels() {
		passes.MaybeAdd(c.renameLabels)
	}

	if options.ShouldFoldConstants() {
		passes.MaybeAdd(c.latePeepholeOptimizations)
	}

	// If side-effects were protected, remove the protection now.
	// Note that when using precompiled libraries we always run this pass regardless of the
	// 'shouldProtectHiddenSideEffects' option: the library compilation may have run
	// with side effect protection enabled even if the binary disables it, so we assume we may
	// always need to strip side effect protection.
	if options.ShouldProtectHiddenSideEffects() || options.GetMergedPrecompiledLibraries() {
		passes.MaybeAdd(c.stripSideEffectProtection)
	}

	if options.GetRenamePrefixNamespace() != nil &&
		options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE {
		if !defaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN.Matches(*options.GetRenamePrefixNamespace()) {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Illegal character in renamePrefixNamespace name: " +
					*options.GetRenamePrefixNamespace()))
		}
		passes.MaybeAdd(c.rescopeGlobalSymbols)
	}

	// Raise to ES2015, if allowed
	if options.GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015) {
		passes.MaybeAdd(c.optimizeToEs6)
	}

	// Must run after all non-safety-check passes as the optimizations do not support modules.
	if options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_ES_MODULES {
		passes.MaybeAdd(c.convertChunksToESModules)
	}

	// Safety checks.  These should always be the last passes.
	passes.MaybeAdd(c.checkAstValidity)
	passes.MaybeAdd(c.varCheckValidity)
	return passes
}

// Java: DefaultPassConfig.getEarlyOptimizationLoopPasses()
func (c *DefaultPassConfig) getEarlyOptimizationLoopPasses() *PassListBuilder {
	options := c.options
	earlyLoopPasses := NewPassListBuilder(options)

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		earlyLoopPasses.MaybeAdd(c.inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		earlyLoopPasses.MaybeAdd(c.inlineConstants)
	}

	if options.GetCollapseObjectLiterals() {
		earlyLoopPasses.MaybeAdd(c.collapseObjectLiterals)
	}

	if c.shouldRunRemoveUnusedCode() {
		earlyLoopPasses.MaybeAdd(c.removeUnusedCode)
	}

	if options.ShouldFoldConstants() {
		earlyLoopPasses.MaybeAdd(c.peepholeOptimizations)
	}

	earlyLoopPasses.AssertAllLoopablePasses()
	return earlyLoopPasses
}

// getPostL10nOptimizations adds optimization passes that need to run after
// late localization has been done.
//
// Once we've replaced message references with string constants and
// `goog.LOCALE` with a constant value, we need to re-run some optimizations,
// so they can throw away code for locales that aren't relevant and perform
// constant folding on the message strings we've now inserted.
//
// Java: DefaultPassConfig.getPostL10nOptimizations()
func (c *DefaultPassConfig) getPostL10nOptimizations() *PassListBuilder {
	options := c.options
	// Localization replaced lots of function calls with constants to get statements like these.
	//
	// `goog.LOCALE = 'es-419';`
	// `const MSG_GREETING = 'Hola';`
	// `const MSG_GREETING_WITH_NAME = 'Hola, ' + person.getName();`
	//
	// Before calling this method we should also have run `flowSensitiveInlineVariables`
	// To make optimization opportunities for peepholeOptimizations.
	//
	// For example, `flowSensitiveInlineVariables` will change this
	// ```
	// var x = 'localized version of message';
	// x = x + '&nbsp;';
	// ```
	// to this
	// ```
	// var x = 'localized version of message' + '&nbsp;';
	// ```
	// Which constant folding can then turn into this
	// ```
	// var x = 'localized version of message&nbsp;';
	// ```
	// Now we should have unblocked a lot of potential optimizations,
	// so do an optimization loop to perform those.
	// This loop is similar to the one created by getMainOptimizationLoop().
	// These should be in the same order as those, but only optimizations we expect to need
	// doing to clean up after localization are included.
	loopPasses := NewPassListBuilder(options)

	if options.ShouldOptimizeCalls() {
		loopPasses.MaybeAdd(c.optimizeCalls)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		loopPasses.MaybeAdd(c.j2clConstantHoisterPass)
		loopPasses.MaybeAdd(c.j2clClinitPass)
	}

	// It is important that inlineVariables and peepholeOptimizations run after inlineFunctions,
	// because inlineFunctions relies on them to clean up patterns it introduces. This affects our
	// size-based loop-termination heuristic.
	if options.GetInlineFunctionsLevel() != CompilerOptions_Reach_NONE {
		loopPasses.MaybeAdd(c.inlineFunctions)
	}

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		loopPasses.MaybeAdd(c.inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		loopPasses.MaybeAdd(c.inlineConstants)
	}

	if c.shouldRunRemoveUnusedCode() {
		loopPasses.MaybeAdd(c.removeUnusedCode)
	}

	if options.ShouldFoldConstants() {
		loopPasses.MaybeAdd(c.peepholeOptimizations)
	}

	loopPasses.AssertAllLoopablePasses()
	return loopPasses
}

// getEarlyOptimizationPasses: these are the passes run in the first half of
// optimizations, which consists of transpilation and some early
// optimization passes.
//
// Java: DefaultPassConfig.getEarlyOptimizationPasses()
func (c *DefaultPassConfig) getEarlyOptimizationPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)
	// At this point all checks have been done.
	if options.ShouldExportTestFunctions() {
		passes.MaybeAdd(c.exportTestFunctions)
	}

	if options.GetMergedPrecompiledLibraries() {
		// Weak sources aren't removed at the library level
		passes.MaybeAdd(c.removeWeakSources)

		// it would be safe to always recompute side effects even if not using precompiled libraries
		// (the else case) but it's unnecessary so skip it to improve build times.
		passes.MaybeAdd(c.checkRegExpForOptimizations)

		// This runs during getChecks(), so only needs to be run here if using precompiled .typedasts
		if options.GetJ2clPass().ShouldAddJ2clPasses() {
			passes.MaybeAdd(c.j2clSourceFileChecker)
		}
	} else {
		c.addNonTypedAstNormalizationPasses(passes)
	}

	// Remove synthetic extern declarations of names that are now defined in source
	// This is expected to do nothing when in a monolithic build
	passes.MaybeAdd(c.removeUnnecessarySyntheticExterns)

	if options.GetSyntheticBlockStartMarker() != nil {
		// This pass must run before the first fold constants pass.
		passes.MaybeAdd(c.createSyntheticBlocks)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(c.j2clPass)
	}

	passes.MaybeAdd(c.transpileAndOptimizeClosureUnaware)

	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)

	if options.GetRewritePolyfills() ||
		options.GetIsolatePolyfills() ||
		options.GetInjectPolyfillsNewerThan() != nil {
		TranspilationPasses_AddRewritePolyfillPass(passes)
	}

	passes.MaybeAdd(c.injectRuntimeLibraries)

	if options.GetClosurePass() {
		passes.MaybeAdd(c.closureProvidesRequires)
	}

	if options.ShouldRunReplaceMessagesForChrome() {
		passes.MaybeAdd(c.replaceMessagesForChrome)
	} else if options.ShouldRunReplaceMessagesPass() {
		if options.DoLateLocalization() {
			// With late localization we protect the messages from mangling by optimizations now,
			// then actually replace them after optimizations.
			// The purpose of doing this is to separate localization from optimization so we can
			// optimize just once for all locales.
			passes.MaybeAdd(c.getProtectMessagesPass())
		} else {
			// TODO(bradfordcsmith): At the moment we expect the optimized output may be slightly
			// smaller if you replace messages before optimizing, but if we can change that, it would
			// be good to drop this early replacement entirely.
			passes.MaybeAdd(c.getFullReplaceMessagesPass())
		}
	}

	if options.DoLateLocalization() {
		passes.MaybeAdd(c.protectLocaleData)
	}

	// Replace 'goog.getCssName' before processing defines
	if options.GetClosurePass() && !options.ShouldPreserveGoogLibraryPrimitives() {
		passes.MaybeAdd(c.closureReplaceGetCssName)
	}

	// Defines in code always need to be processed.
	passes.MaybeAdd(c.processDefinesOptimize)
	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_EARLY_OPTIMIZATIONS_TRANSPILATION))

	passes.MaybeAdd(c.normalize)

	passes.MaybeAdd(c.gatherGettersAndSetters)

	// TODO(b/329447979): Add an early removeUnusedCode pass here
	TranspilationPasses_AddTranspilationPasses(passes, options)

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(c.j2clUtilGetDefineRewriterPass)
	}

	if options.GetInstrumentForCoverageOption() != CompilerOptions_InstrumentOption_NONE {
		passes.MaybeAdd(c.instrumentForCodeCoverage)
	}

	passes.MaybeAdd(c.gatherExternPropertiesOptimize)

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_STANDARD_OPTIMIZATIONS))

	// Abstract method removal works best on minimally modified code, and also
	// only needs to run once.
	if options.GetClosurePass() &&
		(options.ShouldRemoveAbstractMethods() || options.ShouldRemoveClosureAsserts()) {
		passes.MaybeAdd(c.closureCodeRemoval)
	}

	if options.ShouldRemoveJ2clAsserts() {
		passes.MaybeAdd(c.j2clAssertRemovalPass)
	}

	passes.MaybeAdd(c.replaceToggles)
	passes.MaybeAdd(c.inlineAndCollapseProperties)

	if options.GetTweakProcessing().ShouldStrip() ||
		!options.GetStripTypes().IsEmpty() ||
		!options.GetStripNameSuffixes().IsEmpty() ||
		!options.GetStripNamePrefixes().IsEmpty() {
		passes.MaybeAdd(c.stripCode)
	}

	// Ideally this pass would run before transpilation which would allow it to be simplified.
	// It needs to run after `inlineAndCollapseProperties` in order to identify idGenerator calls.
	if options.ShouldReplaceIdGenerators() {
		passes.MaybeAdd(c.replaceIdGenerators)
	}

	// Inline getters/setters in J2CL classes so that Object.defineProperties() calls (resulting
	// from desugaring) don't block class stripping.
	if options.GetJ2clPass().ShouldAddJ2clPasses() &&
		options.GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_ALL {
		// Relies on collapseProperties-triggered aggressive alias inlining.
		passes.MaybeAdd(c.j2clPropertyInlinerPass)
	}

	if options.ShouldInferConsts() {
		passes.MaybeAdd(c.inferConsts)
	}

	// A marker pass to allow {@code ExtraPassConfig} passes to order themselves before
	// RemoveUnusedCode.
	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_OBFUSCATION_PASS_MARKER))

	// Running RemoveUnusedCode before disambiguate properties allows disambiguate properties to be
	// more effective if code that would prevent disambiguation can be removed.
	// TODO(b/66971163): Rename options since we're not actually using smartNameRemoval here now.
	if options.GetSmartNameRemoval() {

		// These passes remove code that is dead because of define flags.
		// If the dead code is weakly typed, running these passes before property
		// disambiguation results in more code removal.
		// The passes are one-time on purpose. (The later runs are loopable.)
		if options.ShouldFoldConstants() &&
			(options.ShouldInlineVariables() || options.ShouldInlineLocalVariables()) {
			passes.MaybeAdd(c.earlyInlineVariables)
			passes.MaybeAdd(c.earlyPeepholeOptimizations)
		}

		passes.MaybeAdd(c.removeUnusedCodeOnce)
	}

	// Property disambiguation should only run once and needs to be done
	// soon after type checking, both so that it can make use of type
	// information and so that other passes can take advantage of the renamed
	// properties.
	if options.ShouldDisambiguateProperties() && options.IsTypecheckingEnabled() {
		passes.MaybeAdd(c.disambiguateProperties)
	}

	if options.ShouldComputeFunctionSideEffects() {
		passes.MaybeAdd(c.markPureFunctions)
	}

	passes.AssertAllOneTimePasses()
	return passes
}

// getLateOptimizationPasses: these are the passes run in the second half of
// optimizations, which consists of the early optimization loop and the main
// optimization loop.
//
// Java: DefaultPassConfig.getLateOptimizationPasses()
func (c *DefaultPassConfig) getLateOptimizationPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)
	if options.GetSmartNameRemoval() {
		// Place one-time marker passes around this loop to prevent the addition of a looping pass
		// above or below from accidentally becoming part of the loop.
		passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_EARLY_OPTIMIZATION_LOOP))
		passes.AddAll(c.getEarlyOptimizationLoopPasses())
		// TODO(): Remove this early loop or rename the option that enables it
		// to something more appropriate.
		passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_AFTER_EARLY_OPTIMIZATION_LOOP))
	}

	// This needs to come after the inline constants pass, which is run within
	// the code removing passes.
	if options.GetClosurePass() {
		passes.MaybeAdd(c.closureOptimizePrimitives)
	}

	// ReplaceStrings runs after CollapseProperties in order to simplify
	// pulling in values of constants defined in enums structures. It also runs
	// after disambiguate properties and smart name removal so that it can
	// correctly identify logging types and can replace references to string
	// expressions.
	if len(options.GetReplaceStringsFunctionDescriptions()) != 0 {
		passes.MaybeAdd(c.replaceStrings)
	}

	// TODO(user): This forces a first crack at crossChunkCodeMotion
	// before devirtualization. Once certain functions are devirtualized,
	// it confuses crossChunkCodeMotion ability to recognized that
	// it is recursive.

	// TODO(user): This is meant for a temporary quick win.
	// In the future, we might want to improve our analysis in
	// CrossChunkCodeMotion so we don't need to do this.
	if options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(c.crossChunkCodeMotion)
	}

	// Method devirtualization benefits from property disambiguation so
	// it should run after that pass but before passes that do
	// optimizations based on global names (like cross-chunk code motion
	// and inline functions).  RemoveUnusedCode does better if run before
	// this pass.
	if options.ShouldDevirtualizeMethods() {
		passes.MaybeAdd(c.devirtualizeMethods)
	}

	passes.MaybeAdd(c.getCustomPasses(CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP))

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_MAIN_OPTIMIZATIONS))

	// Because FlowSensitiveInlineVariables does not operate on the global scope due to compilation
	// time, we need to run it once before InlineFunctions so that we don't miss inlining
	// opportunities when a function will be inlined into the global scope.
	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(c.flowSensitiveInlineVariables)
	}

	passes.AddAll(c.getMainOptimizationLoop())
	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_AFTER_MAIN_OPTIMIZATIONS))

	// Some optimizations belong outside the loop because running them more
	// than once would either have no benefit or be incorrect.
	passes.MaybeAdd(c.getCustomPasses(CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP))

	c.assertValidOrderForOptimizations(passes)
	return passes
}

// getMainOptimizationLoop creates the passes for the main optimization
// loop.
//
// Java: DefaultPassConfig.getMainOptimizationLoop()
func (c *DefaultPassConfig) getMainOptimizationLoop() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)
	if options.ShouldInlineGetters() {
		passes.MaybeAdd(c.inlineSimpleMethods)
	}

	if options.ShouldInlineProperties() && options.IsTypecheckingEnabled() {
		passes.MaybeAdd(c.inlineProperties)
	}

	if options.ShouldRunDeadPropertyAssignmentElimination() {
		passes.MaybeAdd(c.deadPropertyAssignmentElimination)
	}

	if options.ShouldOptimizeCalls() {
		passes.MaybeAdd(c.optimizeCalls)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(c.j2clConstantHoisterPass)
		passes.MaybeAdd(c.j2clClinitPass)
	}

	// It is important that inlineVariables and peepholeOptimizations run after inlineFunctions,
	// because inlineFunctions relies on them to clean up patterns it introduces. This affects our
	// size-based loop-termination heuristic.
	if options.GetInlineFunctionsLevel() != CompilerOptions_Reach_NONE {
		passes.MaybeAdd(c.inlineFunctions)
	}

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(c.inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		passes.MaybeAdd(c.inlineConstants)
	}

	if options.ShouldRunDeadAssignmentElimination() {
		passes.MaybeAdd(c.deadAssignmentsElimination)
	}

	if options.GetCollapseObjectLiterals() {
		passes.MaybeAdd(c.collapseObjectLiterals)
	}

	if c.shouldRunRemoveUnusedCode() {
		passes.MaybeAdd(c.removeUnusedCode)
	}

	if options.ShouldFoldConstants() {
		passes.MaybeAdd(c.peepholeOptimizations)
	}

	passes.AssertAllLoopablePasses()
	return passes
}

// addNonTypedAstNormalizationPasses is for use in builds that run
// getOptimizations() on the result of an AST directly from getChecks(),
// that has not gone through TypedAST serialization/deserialization.
//
// Java: DefaultPassConfig.addNonTypedAstNormalizationPasses(PassListBuilder)
func (c *DefaultPassConfig) addNonTypedAstNormalizationPasses(passes *PassListBuilder) {
	passes.MaybeAdd(c.removeCastNodes)
	passes.MaybeAdd(c.typesToColors)
}

// Java: DefaultPassConfig.shouldRunRemoveUnusedCode()
func (c *DefaultPassConfig) shouldRunRemoveUnusedCode() bool {
	options := c.options
	return options.ShouldRemoveUnusedVariables() ||
		options.ShouldRemoveUnusedLocalVariables() ||
		options.ShouldRemoveUnusedPrototypeProperties() ||
		options.IsRemoveUnusedClassProperties() ||
		options.GetRewritePolyfills()
}

// assertValidOrderForChecks: certain checks and rewriting passes need to
// run in a particular order. For example, the PolymerPass will not work
// correctly unless it runs after the goog.provide() processing. This
// enforces those constraints.
//
// checks is the list of check passes.
//
// Java: DefaultPassConfig.assertValidOrderForChecks(PassListBuilder)
func (c *DefaultPassConfig) assertValidOrderForChecks(checks *PassListBuilder) {
	checks.AssertPassOrder(
		c.declaredGlobalExternsOnWindow,
		c.checkVars,
		"declaredGlobalExternsOnWindow must happen before VarCheck, which adds synthetic externs")
	checks.AssertPassOrder(
		c.chromePass,
		c.checkJsDocAndEs6Modules,
		"The ChromePass must run before after JsDoc and Es6 module checking.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		c.processDefinesCheck,
		"Must rewrite goog.module before processing @define's, so that @defines in modules work.")
	checks.AssertPassOrder(
		c.closurePrimitives, c.polymerPass, "The Polymer pass must run after goog.provide processing.")
	checks.AssertPassOrder(
		c.chromePass, c.polymerPass, "The Polymer pass must run after ChromePass processing.")
	checks.AssertPassOrder(
		c.polymerPass, c.suspiciousCode, "The Polymer pass must run before suspiciousCode processing.")
	checks.AssertPassOrder(
		c.addSyntheticScript,
		c.gatherModuleMetadataPass,
		"Cannot add a synthetic script node after module metadata creation.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		c.removeSyntheticScript,
		"Synthetic script node should be removed only after module rewriting.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		c.rewriteCallerCodeLocation,
		"ClosureRewriteModule must happen before RewriteCallerCodeLocation, so that exported"+
			" functions and call sites are rewritten correctly.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		TranspilationPasses_GetEs6RewriteDestructuring(
			Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS),
		"RewriteCallerCodeLocation must happen before Es6RewriteDestructuring, because we need"+
			" ReWriteCallerCodeLocation to run before default parameters get rewritten.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		TranspilationPasses_GetEs6RewriteDestructuring(
			Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST),
		"RewriteCallerCodeLocation must happen before Es6RewriteDestructuring, because we need"+
			" ReWriteCallerCodeLocation to run before default parameters get rewritten.")

	if checks.Contains(c.closureGoogScopeAliases) {
		jutil.CheckStateMsg(
			checks.Contains(c.checkVariableReferences),
			"goog.scope processing requires variable checking")
	}
	// TODO(lharker): add this back once fixing the ProcessCommonJSModules VarCheck branch
	// checks.assertPassOrder(
	//     checkVariableReferences,
	//     closureGoogScopeAliases,
	//     "Variable checking must happen before goog.scope processing.");

	checks.AssertPassOrder(
		c.gatherModuleMetadataPass,
		c.closureCheckModule,
		"Need to gather module metadata before checking closure modules.")

	checks.AssertPassOrder(
		c.gatherModuleMetadataPass,
		c.createModuleMapPass,
		"Need to gather module metadata before scanning modules.")

	checks.AssertPassOrder(
		c.createModuleMapPass,
		c.rewriteCommonJsModules,
		"Need to gather module information before rewriting CommonJS modules.")

	checks.AssertPassOrder(
		c.rewriteScriptsToEs6Modules,
		c.gatherModuleMetadataPass,
		"Need to gather module information after rewriting scripts to modules.")

	checks.AssertPassOrder(
		c.gatherModuleMetadataPass,
		c.checkMissingRequires,
		"Need to gather module information before checking for missing requires.")
}

// assertValidOrderForOptimizations: certain optimizations need to run in a
// particular order. For example, OptimizeCalls must run before
// RemoveSuperMethodsPass, because the former can invalidate assumptions in
// the latter. This enforces those constraints.
//
// optimizations is the list of optimization passes.
//
// Java: DefaultPassConfig.assertValidOrderForOptimizations(PassListBuilder)
func (c *DefaultPassConfig) assertValidOrderForOptimizations(optimizations *PassListBuilder) {
	optimizations.AssertPassOrder(
		c.j2clPass,
		TranspilationPasses_rewriteGenerators,
		"J2CL normalization should be done before generator re-writing.")

	optimizations.AssertPassOrder(
		TranspilationPasses_rewritePolyfills,
		c.processDefinesOptimize,
		"Polyfill injection must be done before processDefines as some polyfills reference "+
			"goog.defines.")
	optimizations.AssertPassOrder(
		TranspilationPasses_injectTranspilationRuntimeLibraries,
		c.processDefinesOptimize,
		"Runtime library injection must be done before processDefines some runtime libraries "+
			"reference goog.defines.")

	optimizations.AssertPassOrder(
		c.processDefinesOptimize,
		c.j2clUtilGetDefineRewriterPass,
		"J2CL define re-writing should be done after processDefines since it relies on "+
			"Compiler#getDefineNames to have been populated by it.")

	optimizations.AssertPassOrder(
		c.removeUnusedCode,
		c.isolatePolyfills,
		"Polyfill isolation should be done after RemovedUnusedCode. Otherwise unused polyfill"+
			" removal will not find any polyfill usages and will delete all polyfills.")

	optimizations.AssertPassOrder(
		TranspilationPasses_instrumentAsyncContext,
		TranspilationPasses_rewriteAsyncIteration,
		"AsyncContext should be instrumentated before await and/or yield is transpiled away")

	optimizations.AssertPassOrder(
		TranspilationPasses_instrumentAsyncContext,
		TranspilationPasses_rewriteAsyncFunctions,
		"AsyncContext should be instrumentated before await and/or yield is transpiled away")

	optimizations.AssertPassOrder(
		TranspilationPasses_instrumentAsyncContext,
		TranspilationPasses_rewriteGenerators,
		"AsyncContext should be instrumentated before await and/or yield is transpiled away")
}

// Java: DefaultPassConfig.GENERATE_EXPORTS_ERROR
var defaultPassConfig_GENERATE_EXPORTS_ERROR = DiagnosticType_Error(
	"JSC_GENERATE_EXPORTS_ERROR",
	"Exports can only be generated if export symbol/property functions are set.")

// getFullReplaceMessagesPass returns the form of `replaceMessages` that
// does the replacement all at once, and which must run before
// optimizations.
//
// Java: DefaultPassConfig.getFullReplaceMessagesPass()
func (c *DefaultPassConfig) getFullReplaceMessagesPass() *PassFactory {
	options := c.options
	return PassFactory_Builder_().
		SetName(PassNames_REPLACE_MESSAGES).
		SetInternalFactory(
			func(compiler AbstractCompiler) CompilerPass {
				return NewReplaceMessages(
					compiler,
					options.GetMessageBundle(),
					/* allow messages with goog.getMsg */
					options.GetStrictMessageReplacement()).
					GetFullReplacementPass()
			}).
		Build()
}

// getProtectMessagesPass returns the pass that protects messages from
// mangling, so they can be found and replaced after optimizations.
//
// Java: DefaultPassConfig.getProtectMessagesPass()
func (c *DefaultPassConfig) getProtectMessagesPass() *PassFactory {
	options := c.options
	return PassFactory_Builder_().
		SetName("protectMessages").
		SetInternalFactory(
			func(compiler AbstractCompiler) CompilerPass {
				return NewReplaceMessages(
					compiler,
					options.GetMessageBundle(),
					/* allow messages with goog.getMsg */
					options.GetStrictMessageReplacement()).
					GetMsgProtectionPass()
			}).
		Build()
}

// getReplaceProtectedMessagesPass returns the pass that identifies
// protected messages and completes their replacement.
//
// Java: DefaultPassConfig.getReplaceProtectedMessagesPass()
func (c *DefaultPassConfig) getReplaceProtectedMessagesPass() *PassFactory {
	options := c.options
	return PassFactory_Builder_().
		SetName("replaceProtectedMessages").
		SetInternalFactory(
			func(compiler AbstractCompiler) CompilerPass {
				return NewReplaceMessages(
					compiler,
					options.GetMessageBundle(),
					/* allow messages with goog.getMsg */
					options.GetStrictMessageReplacement()).
					GetReplacementCompletionPass()
			}).
		Build()
}

// defaultPassConfig_createPeepholeOptimizationsPass: various peephole
// optimizations.
//
// Java: DefaultPassConfig.createPeepholeOptimizationsPass(AbstractCompiler, String, boolean)
func defaultPassConfig_createPeepholeOptimizationsPass(
	compiler AbstractCompiler, passName string, expectAstIsNormalized bool) CompilerPass {
	jutil.CheckArgumentMsg1(
		expectAstIsNormalized == compiler.GetLifeCycleStage().IsNormalized(),
		"%s",
		compiler.GetLifeCycleStage())
	const late = false
	useTypesForOptimization :=
		compiler.GetOptions().ShouldUseTypesForLocalOptimization()
	var optimizations []AbstractPeepholeOptimization
	if expectAstIsNormalized {
		// MinimizeExitPoints requires the AST to be normalized.
		optimizations = append(optimizations, NewMinimizeExitPoints())
	}
	optimizations = append(optimizations, NewPeepholeMinimizeConditions(late))
	optimizations = append(optimizations, NewPeepholeSubstituteAlternateSyntax(late))
	optimizations = append(optimizations, NewPeepholeReplaceKnownMethods(late, useTypesForOptimization))
	optimizations = append(optimizations, NewPeepholeRemoveDeadCode())
	if compiler.GetOptions().GetJ2clPass().ShouldAddJ2clPasses() {
		optimizations = append(optimizations, NewJ2clEqualitySameRewriterPass(useTypesForOptimization))
		optimizations = append(optimizations, NewJ2clStringValueOfRewriterPass())
		optimizations = append(optimizations, NewJ2clUndefinedChecksRewriterPass())
	}
	optimizations = append(optimizations, NewPeepholeFoldConstants(late, useTypesForOptimization))
	optimizations = append(optimizations, NewPeepholeCollectPropertyAssignments())
	return NewPeepholeOptimizationsPassOfList(compiler, passName, optimizations)
}

// defaultPassConfig_combineChecks executes the given callbacks with a
// CombinedCompilerPass.
//
// Java: DefaultPassConfig.combineChecks(AbstractCompiler, List<NodeTraversal.Callback>)
func defaultPassConfig_combineChecks(
	compiler AbstractCompiler, callbacks []NodeTraversal_Callback) CompilerPass {
	jutil.CheckArgument(len(callbacks) != 0)
	return NewCombinedCompilerPassOfList(compiler, callbacks)
}

// Java: DefaultPassConfig.createProcessDefines(ProcessDefines.Mode)
func (c *DefaultPassConfig) createProcessDefines(mode ProcessDefines_Mode) *PassFactory {
	options := c.options
	return PassFactory_Builder_().
		SetName("processDefines_" + mode.String()).
		SetInternalFactory(
			func(compiler AbstractCompiler) CompilerPass {
				return NewProcessDefines_Builder(compiler).
					PutReplacements(defaultPassConfig_getAdditionalReplacements(options)).
					PutReplacements(options.GetDefineReplacements()).
					SetMode(mode).
					SetRecognizeClosureDefines(compiler.GetOptions().GetClosurePass()).
					SetEnableZonesDefineName(options.GetEnableZonesDefineName()).
					SetZoneInputPattern(options.GetZoneInputPattern()).
					SetUnknownDefinesToIgnore(options.GetUnknownDefinesToIgnore()).
					Build()
			}).
		Build()
}

// Java: DefaultPassConfig.runVariableRenaming(AbstractCompiler, VariableMap, Node, Node)
func (c *DefaultPassConfig) runVariableRenaming(
	compiler AbstractCompiler, prevVariableMap *VariableMap, externs *rhino.Node, root *rhino.Node) *VariableMap {
	options := c.options
	reservedChars := jutil.NewLinkedHashSet[uint16]()
	reservedNames := jutil.NewLinkedHashSet[string]()
	if options.GetRenamePrefixNamespace() != nil {
		// don't use the prefix name as a global symbol.
		reservedNames.Add(*options.GetRenamePrefixNamespace())
	}
	for name := range compiler.GetExportedNames().All() {
		reservedNames.Add(name)
	}
	for name := range parsing.ParserRunner_GetReservedVars().All() {
		reservedNames.Add(name)
	}
	rn :=
		NewRenameVars(
			compiler,
			options.GetRenamePrefix(),
			options.GetVariableRenaming() == VariableRenamingPolicy_LOCAL,
			options.ShouldGeneratePseudoNames(),
			options.ShouldPreferStableNames(),
			prevVariableMap,
			reservedChars,
			reservedNames,
			options.GetNameGenerator())
	rn.Process(externs, root)
	return rn.GetVariableMap()
}

// createGatherExternProperties: extern property names gathering pass.
//
// Java: DefaultPassConfig.createGatherExternProperties(GatherExternProperties.Mode)
func (c *DefaultPassConfig) createGatherExternProperties(mode GatherExternProperties_Mode) *PassFactory {
	return PassFactory_Builder_().
		SetName("gatherExternProperties").
		SetInternalFactory(func(compiler AbstractCompiler) CompilerPass { return NewGatherExternProperties(compiler, mode) }).
		Build()
}

// getCustomPasses runs custom passes that are designated to run at a
// particular time.
//
// TODO(nickreid): Deprecate this API
//
// Java: DefaultPassConfig.getCustomPasses(CustomPassExecutionTime)
func (c *DefaultPassConfig) getCustomPasses(executionTime CustomPassExecutionTime) *PassFactory {
	options := c.options
	return PassFactory_Builder_().
		SetName("runCustomPasses").
		SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
			return defaultPassConfig_runInSerial(options.GetCustomPassesAt(executionTime))
		}).
		Build()
}

// defaultPassConfig_runInSerial creates a compiler pass that runs the given
// passes in serial.
//
// Java: DefaultPassConfig.runInSerial(Collection<CompilerPass>)
func defaultPassConfig_runInSerial(passes []CompilerPass) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		for _, pass := range passes {
			pass.Process(externs, root)
		}
	})
}

// Java: DefaultPassConfig.getAdditionalReplacements(CompilerOptions)
func defaultPassConfig_getAdditionalReplacements(options *CompilerOptions) *jutil.LinkedHashMap[string, *rhino.Node] {
	additionalReplacements := jutil.NewLinkedHashMap[string, *rhino.Node]()

	if options.ShouldMarkAsCompiled() || options.GetClosurePass() {
		additionalReplacements.Put(defaultPassConfig_COMPILED_CONSTANT_NAME, rhino.IR_TrueNode())
	}

	if options.GetClosurePass() && options.GetLocale() != nil && !options.DoLateLocalization() {
		additionalReplacements.Put(defaultPassConfig_CLOSURE_LOCALE_CONSTANT_NAME, rhino.IR_String(*options.GetLocale()))
	}

	return additionalReplacements
}

// Rewrites ES modules import paths to be browser compliant
//
// Java: DefaultPassConfig.forbidDynamicImportUsage
var defaultPassConfig_forbidDynamicImportUsage = PassFactory_Builder_().
	SetName("FORBID_DYNAMIC_IMPORT").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewForbidDynamicImportUsage(compiler)
	}).
	Build()

// Rewrites Dynamic Import Expressions
//
// Java: DefaultPassConfig.rewriteDynamicImports
var defaultPassConfig_rewriteDynamicImports = PassFactory_Builder_().
	SetName("REWRITE_DYNAMIC_IMPORT").
	SetInternalFactory(
		func(compiler AbstractCompiler) CompilerPass {
			return NewRewriteDynamicImports(
				compiler,
				compiler.GetOptions().GetDynamicImportAlias(),
				compiler.GetOptions().GetChunkOutputType())
		}).
	Build()

// Java: DefaultPassConfig.requirePropertiesAreStaticallyAnalyzable(CompilerOptions)
func defaultPassConfig_requirePropertiesAreStaticallyAnalyzable(
	options *CompilerOptions) *PassFactory_PreconditionResult {
	return NewPassFactory_PreconditionResult(
		options.GetAssumePropertiesAreStaticallyAnalyzable(),
		jutil.StringPtr("requires assumePropertiesAreStaticallyAnalyzable to be enabled"))
}
