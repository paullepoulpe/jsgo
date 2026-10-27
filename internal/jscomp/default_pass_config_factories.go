// Port of com/google/javascript/jscomp/DefaultPassConfig.java @ closure-compiler v20261021

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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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

					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, jsRoot rhino.Node) {
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

					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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

					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
				return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) { compiler.ClearTypedScopeCreator() })
			}).
			Build()

	// Clears the top typed scope when we're done with it.
	c.clearTopTypedScopePass =
		PassFactory_Builder_().
			SetName("clearTopTypedScopePass").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(unused rhino.Node, srcRoot rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					var unreachable *CheckUnreachableCode
					if !options.Disables(DiagnosticGroups_CHECK_USELESS_CODE) {
						unreachable = NewCheckUnreachableCode(compiler)
						callbacks = append(callbacks, unreachable)
					}
					if !options.Disables(DiagnosticGroups_MISSING_RETURN) {
						mr := NewCheckMissingReturn(compiler)
						mr.sharedCfgSource = unreachable
						callbacks = append(callbacks, mr)
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
					return CompilerPassFunc(func(externs rhino.Node, jsRoot rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, jsRoot rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) { compiler.InitializeSyntheticCodeInput() })
				}).
			Build()

	c.removeSyntheticScript =
		PassFactory_Builder_().
			SetName("REMOVE_SYNTHETIC_SCRIPT").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) { compiler.RemoveSyntheticCodeInput() })
			}).
			Build()

	c.mergeSyntheticScript =
		PassFactory_Builder_().
			SetName("MERGE_SYNTHETIC_SCRIPT").
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) { compiler.MergeSyntheticCodeInput() })
			}).
			Build()

	// Replace locale values with stubs that can survive optimizations and be
	// replaced afterwards.
	c.protectLocaleData =
		PassFactory_Builder_().
			SetName("protectLocaleData").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
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
	compiler AbstractCompiler, prevVariableMap *VariableMap, externs rhino.Node, root rhino.Node) *VariableMap {
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
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		for _, pass := range passes {
			pass.Process(externs, root)
		}
	})
}

// Java: DefaultPassConfig.getAdditionalReplacements(CompilerOptions)
func defaultPassConfig_getAdditionalReplacements(options *CompilerOptions) *jutil.LinkedHashMap[string, rhino.Node] {
	additionalReplacements := jutil.NewLinkedHashMap[string, rhino.Node]()

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
