// Port of com/google/javascript/jscomp/TranspilationPasses.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TranspilationPasses provides a single place to manage transpilation
// passes.
//
// Go-only: language-adapter — the Java class only has static members (and a
// private constructor), which are package-level functions and variables here.

// Java: TranspilationPasses.addEs6ModulePass(PassListBuilder, PreprocessorSymbolTable.CachedInstanceFactory)
func TranspilationPasses_AddEs6ModulePass(
	passes *PassListBuilder,
	preprocessorTableFactory *PreprocessorSymbolTable_CachedInstanceFactory) {
	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName("es6RewriteModule").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					preprocessorTableFactory.MaybeInitialize(compiler)
					return NewEs6RewriteModulesWithChunkOutputType(
						compiler,
						compiler.GetModuleMetadataMap(),
						compiler.GetModuleMap(),
						preprocessorTableFactory.GetInstanceOrNull(),
						compiler.GetTopScope(),
						compiler.GetOptions().GetChunkOutputType())
				}).
			Build())
}

// Java: TranspilationPasses.addTranspilationRuntimeLibraries(PassListBuilder)
func TranspilationPasses_AddTranspilationRuntimeLibraries(passes *PassListBuilder) {
	// Inject runtime libraries needed for the transpilation we will have to do. Should run before
	// typechecking.
	passes.MaybeAdd(TranspilationPasses_injectTranspilationRuntimeLibraries)
}

// Java: TranspilationPasses.addEs6ModuleToCjsPass(PassListBuilder)
func TranspilationPasses_AddEs6ModuleToCjsPass(passes *PassListBuilder) {
	passes.MaybeAdd(TranspilationPasses_es6RewriteModuleToCjs)
}

// Java: TranspilationPasses.addEs6RewriteImportPathPass(PassListBuilder)
func TranspilationPasses_AddEs6RewriteImportPathPass(passes *PassListBuilder) {
	passes.MaybeAdd(TranspilationPasses_es6RelativizeImportPaths)
}

// TranspilationPasses_AddTranspilationPasses adds transpilation passes.
//
// Passes added in this method either use TranspilationPasses.processTranspile
// or early-exit by checking their feature in the script's featureset. So they
// get run if the feature they're responsible for removing exists in the
// script.
//
// Java: TranspilationPasses.addTranspilationPasses(PassListBuilder, CompilerOptions)
func TranspilationPasses_AddTranspilationPasses(passes *PassListBuilder, options *CompilerOptions) {

	// bag of small, local, independent rewritings being done in a single traversal
	passes.MaybeAdd(TranspilationPasses_peepholeTranspilationsPasses)

	passes.MaybeAdd(TranspilationPasses_CreateUnifiedFeatureRemovalPass("featureRemovalPasses", options))

	// Es6NormalizeClasses is always run so that future passes can make assumptions about classes.
	// Always does the following:
	// -  Extracts classes defined in expression contexts (e.g. `foo(class {})` and
	//    `class C extends class {} {}`).
	// -  Removes all static initialization (i.e. removes static initialization blocks and moves
	//    static field initializers into a generated static method).
	// -  Provides a single name for classes (i.e. removes the inner class name).
	// -  Moves computed field expressions into named variables.
	// If transpilation is needed for Feature.PUBLIC_CLASS_FIELDS:
	// -  Rewrites public fields as assignments in the constructor.
	// -  Removes all field declarations (public and static).
	passes.MaybeAdd(TranspilationPasses_es6NormalizeClasses)

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT) {
		passes.MaybeAdd(TranspilationPasses_rewriteLogicalAssignmentOperatorsPass)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_OPTIONAL_CHAINING) {
		passes.MaybeAdd(TranspilationPasses_rewriteOptionalChainingOperator)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_NULL_COALESCE_OP) {
		passes.MaybeAdd(TranspilationPasses_rewriteNullishCoalesceOperator)
	}

	// NOTE: This needs to be _before_ await and yield are transpiled away.
	if options.GetInstrumentAsyncContext() {
		passes.MaybeAdd(TranspilationPasses_instrumentAsyncContext)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_FOR_AWAIT_OF) ||
		options.NeedsTranspilationOf(parser.FeatureSet_Feature_ASYNC_GENERATORS) {
		passes.MaybeAdd(TranspilationPasses_rewriteAsyncIteration)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD) ||
		options.NeedsTranspilationOf(parser.FeatureSet_Feature_OBJECT_PATTERN_REST) {
		passes.MaybeAdd(TranspilationPasses_rewriteObjectSpread)
		if !options.NeedsTranspilationOf(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING) &&
			options.NeedsTranspilationOf(parser.FeatureSet_Feature_OBJECT_PATTERN_REST) {
			passes.MaybeAdd(
				TranspilationPasses_GetEs6RewriteDestructuring(Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST))
		}
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_ASYNC_FUNCTIONS) {
		passes.MaybeAdd(TranspilationPasses_rewriteAsyncFunctions)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_EXPONENT_OP) {
		passes.MaybeAdd(TranspilationPasses_rewriteExponentialOperator)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_CLASSES) {
		passes.MaybeAdd(TranspilationPasses_es6ConvertSuper)
	}

	if options.NeedsTranspilationFrom(
		parser.FeatureSet_BARE_MINIMUM.WithFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)) {
		passes.MaybeAdd(
			TranspilationPasses_GetEs6RewriteDestructuring(Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS))
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_ARROW_FUNCTIONS) {
		passes.MaybeAdd(TranspilationPasses_es6RewriteArrowFunction)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_CLASSES) {
		passes.MaybeAdd(TranspilationPasses_es6RewriteClass)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_REST_PARAMETERS) {
		passes.MaybeAdd(TranspilationPasses_es6RewriteRestParameters)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS) {
		passes.MaybeAdd(TranspilationPasses_es6RewriteSpreadExpressions)
	}

	if options.NeedsTranspilationFrom(
		parser.FeatureSet_BARE_MINIMUM.WithFeatures(
			parser.FeatureSet_Feature_COMPUTED_PROPERTIES, parser.FeatureSet_Feature_MEMBER_DECLARATIONS, parser.FeatureSet_Feature_TEMPLATE_LITERALS)) {
		passes.MaybeAdd(TranspilationPasses_lateConvertEs6ToEs3)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_FOR_OF) {
		passes.MaybeAdd(TranspilationPasses_es6ForOf)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION) {
		passes.MaybeAdd(TranspilationPasses_rewriteBlockScopedFunctionDeclaration)
	}

	if options.NeedsTranspilationFrom(
		parser.FeatureSet_BARE_MINIMUM.WithFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS, parser.FeatureSet_Feature_CONST_DECLARATIONS)) {
		passes.MaybeAdd(TranspilationPasses_rewriteBlockScopedDeclaration)
	}

	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_GENERATORS) {
		passes.MaybeAdd(TranspilationPasses_rewriteGenerators)
	}
	// This pass must run at the end of all transpiler passes. It validates that only supported
	// features remain in the compiler's featureSet
	passes.MaybeAdd(
		TranspilationPasses_CreatePostTranspileUnsupportedFeaturesRemovedCheck(
			"postTranspileUnsupportedFeaturesRemovedCheck"))
}

// Java: TranspilationPasses.peepholeTranspilationsPasses
var TranspilationPasses_peepholeTranspilationsPasses = PassFactory_Builder_().
	SetName("peepholeTranspilationsPasses").
	SetInternalFactory(
		func(compiler AbstractCompiler) CompilerPass {
			peepholeTranspilations := []AbstractPeepholeTranspilation{}
			peepholeTranspilations = append(peepholeTranspilations,
				NewReportUntranspilableFeatures(
					compiler,
					compiler.GetOptions().GetBrowserFeaturesetYearObject(),
					compiler.GetOptions().GetOutputFeatureSet()))
			if compiler.GetOptions().NeedsTranspilationOf(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING) {
				peepholeTranspilations = append(peepholeTranspilations, NewRewriteCatchWithNoBinding(compiler))
			}
			if compiler.
				GetOptions().
				NeedsTranspilationOf(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES) {
				peepholeTranspilations = append(peepholeTranspilations, NewEs6NormalizeShorthandProperties(compiler))
			}
			if compiler.GetOptions().NeedsTranspilationOf(parser.FeatureSet_Feature_NEW_TARGET) {
				peepholeTranspilations = append(peepholeTranspilations, NewRewriteNewDotTarget(compiler))
			}
			return PeepholeTranspilationsPass_Create(compiler, peepholeTranspilations)
		}).
	Build()

// TranspilationPasses_AddRewritePolyfillPass adds the pass to inject ES2015
// polyfills, which goes after the late ES2015 passes.
//
// Java: TranspilationPasses.addRewritePolyfillPass(PassListBuilder)
func TranspilationPasses_AddRewritePolyfillPass(passes *PassListBuilder) {
	passes.MaybeAdd(TranspilationPasses_rewritePolyfills)
}

// Rewrites ES6 modules
//
// Java: TranspilationPasses.es6RewriteModuleToCjs
var TranspilationPasses_es6RewriteModuleToCjs = PassFactory_Builder_().
	SetName("es6RewriteModuleToCjs").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RewriteModulesToCommonJsModules(compiler)
	}).
	Build()

// Rewrites ES6 modules import paths to be browser compliant
//
// Java: TranspilationPasses.es6RelativizeImportPaths
var TranspilationPasses_es6RelativizeImportPaths = PassFactory_Builder_().
	SetName("es6RelativizeImportPaths").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RelativizeImportPaths(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteAsyncFunctions
var TranspilationPasses_rewriteAsyncFunctions = PassFactory_Builder_().
	SetName("rewriteAsyncFunctions").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return RewriteAsyncFunctions_Create(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteAsyncIteration
var TranspilationPasses_rewriteAsyncIteration = PassFactory_Builder_().
	SetName("rewriteAsyncIteration").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return RewriteAsyncIteration_Create(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteObjectSpread
var TranspilationPasses_rewriteObjectSpread = PassFactory_Builder_().
	SetName("rewriteObjectSpread").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewRewriteObjectSpread(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteExponentialOperator
var TranspilationPasses_rewriteExponentialOperator = PassFactory_Builder_().
	SetName("rewriteExponentialOperator").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs7RewriteExponentialOperator(compiler)
	}).
	Build()

// Java: TranspilationPasses.es6NormalizeClasses
var TranspilationPasses_es6NormalizeClasses = PassFactory_Builder_().
	SetName(PassNames_ES6_NORMALIZE_CLASSES).
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6NormalizeClasses(compiler)
	}).
	Build()

// Java: TranspilationPasses.es6RewriteClass
var TranspilationPasses_es6RewriteClass = PassFactory_Builder_().
	SetName("Es6RewriteClass").
	SetInternalFactory(
		func(compiler AbstractCompiler) CompilerPass {
			return NewEs6RewriteClass(
				compiler, compiler.GetOptions().GetEs6SubclassTranspilation())
		}).
	Build()

// Java: TranspilationPasses.getEs6RewriteDestructuring(ObjectDestructuringRewriteMode)
func TranspilationPasses_GetEs6RewriteDestructuring(rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode) *PassFactory {
	return PassFactory_Builder_().
		SetName("Es6RewriteDestructuring").
		SetInternalFactory(
			func(compiler AbstractCompiler) CompilerPass {
				return NewEs6RewriteDestructuring_Builder(compiler).
					SetDestructuringRewriteMode(rewriteMode).
					Build()
			}).
		Build()
}

// Java: TranspilationPasses.es6RewriteArrowFunction
var TranspilationPasses_es6RewriteArrowFunction = PassFactory_Builder_().
	SetName("Es6RewriteArrowFunction").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RewriteArrowFunction(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewritePolyfills
var TranspilationPasses_rewritePolyfills = PassFactory_Builder_().
	SetName("RewritePolyfills").
	SetInternalFactory(
		func(compiler AbstractCompiler) CompilerPass {
			return NewRewritePolyfills(
				compiler,
				compiler.GetOptions().GetRewritePolyfills(),
				compiler.GetOptions().GetIsolatePolyfills(),
				compiler.GetOptions().GetInjectPolyfillsNewerThan())
		}).
	Build()

// Java: TranspilationPasses.instrumentAsyncContext
var TranspilationPasses_instrumentAsyncContext = PassFactory_Builder_().
	SetName("instrumentAsyncContext").
	SetInternalFactory(
		func(compiler AbstractCompiler) CompilerPass {
			return NewInstrumentAsyncContext(
				compiler,
				compiler.
					GetOptions().
					GetOutputFeatureSet().
					ContainsFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS))
		}).
	Build()

// Java: TranspilationPasses.es6ConvertSuper
var TranspilationPasses_es6ConvertSuper = PassFactory_Builder_().
	SetName("es6ConvertSuper").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6ConvertSuper(compiler)
	}).
	Build()

// Injects runtime library code needed for transpiled ES2015+ code.
//
// Java: TranspilationPasses.injectTranspilationRuntimeLibraries
var TranspilationPasses_injectTranspilationRuntimeLibraries = PassFactory_Builder_().
	SetName("es6InjectRuntimeLibraries").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewInjectTranspilationRuntimeLibraries(compiler)
	}).
	Build()

// Transpiles REST parameters.
//
// Java: TranspilationPasses.es6RewriteRestParameters
var TranspilationPasses_es6RewriteRestParameters = PassFactory_Builder_().
	SetName("es6RewriteRestParameters").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RewriteRestParameters(compiler)
	}).
	Build()

// Transpiles SPREAD in both array literals, new and function calls.
//
// Java: TranspilationPasses.es6RewriteSpreadExpressions
var TranspilationPasses_es6RewriteSpreadExpressions = PassFactory_Builder_().
	SetName("es6RewriteSpreadExpressions").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RewriteSpreadExpressions(compiler)
	}).
	Build()

// Does the main ES2015 to ES3 conversion. There are a few other passes which
// run before this one, to convert constructs which are not converted by this
// pass.
//
// Java: TranspilationPasses.lateConvertEs6ToEs3
var TranspilationPasses_lateConvertEs6ToEs3 = PassFactory_Builder_().
	SetName("lateConvertEs6").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewLateEs6ToEs3Converter(compiler)
	}).
	Build()

// Java: TranspilationPasses.es6ForOf
var TranspilationPasses_es6ForOf = PassFactory_Builder_().SetName("es6ForOf").SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
	return NewEs6ForOfConverter(compiler)
}).Build()

// Java: TranspilationPasses.rewriteBlockScopedFunctionDeclaration
var TranspilationPasses_rewriteBlockScopedFunctionDeclaration = PassFactory_Builder_().
	SetName("Es6RewriteBlockScopedFunctionDeclaration").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RewriteBlockScopedFunctionDeclaration(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteBlockScopedDeclaration
var TranspilationPasses_rewriteBlockScopedDeclaration = PassFactory_Builder_().
	SetName("Es6RewriteBlockScopedDeclaration").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RewriteBlockScopedDeclaration(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteGenerators
var TranspilationPasses_rewriteGenerators = PassFactory_Builder_().
	SetName("rewriteGenerators").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RewriteGenerators(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteLogicalAssignmentOperatorsPass
var TranspilationPasses_rewriteLogicalAssignmentOperatorsPass = PassFactory_Builder_().
	SetName("rewriteLogicalAssignmentOperatorsPass").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewRewriteLogicalAssignmentOperatorsPass(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteOptionalChainingOperator
var TranspilationPasses_rewriteOptionalChainingOperator = PassFactory_Builder_().
	SetName("rewriteOptionalChainingOperator").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewRewriteOptionalChainingOperator(compiler)
	}).
	Build()

// Java: TranspilationPasses.rewriteNullishCoalesceOperator
var TranspilationPasses_rewriteNullishCoalesceOperator = PassFactory_Builder_().
	SetName("rewriteNullishCoalesceOperator").
	SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
		return NewRewriteNullishCoalesceOperator(compiler)
	}).
	Build()

// TranspilationPasses_DoesScriptHaveAnyOfTheseFeatures returns true if the
// script's featureSet contains any feature from the given featureSet.
//
// The script's featureSet gets accurately maintained during transpile. It can
// be relied upon to check if a feature exists in the AST or not.
//
// Java: TranspilationPasses.doesScriptHaveAnyOfTheseFeatures(Node, FeatureSet)
func TranspilationPasses_DoesScriptHaveAnyOfTheseFeatures(script *rhino.Node, featureSet *parser.FeatureSet) bool {
	features := NodeUtil_GetFeatureSetOfScript(script)
	return features != nil && features.ContainsAtLeastOneOf(featureSet)
}

// TranspilationPasses_ProcessTranspile runs the given transpilation callbacks
// on every source script if the given featuresToRunFor are unsupported in the
// output language and exist in that script.
//
// compiler is an AbstractCompiler; combinedRoot is the combined root for all
// JS files; featuresToRunFor are the features for which these callbacks run;
// callbacks are the callbacks that should be invoked if a file has
// output-unsupported features.
//
// Deprecated: Please use a regular NodeTraversal object directly, using
// `shouldTraverse` to skip SCRIPT node if desired.
//
// Java: TranspilationPasses.processTranspile(AbstractCompiler, Node, FeatureSet, NodeTraversal.Callback...)
func TranspilationPasses_ProcessTranspile(
	compiler AbstractCompiler,
	combinedRoot *rhino.Node,
	featuresToRunFor *parser.FeatureSet,
	callbacks ...NodeTraversal_Callback) {
	languageOutFeatures := compiler.GetOptions().GetOutputFeatureSet()
	if languageOutFeatures.Contains(featuresToRunFor) {
		// all featuresToRunFor are supported by languageOut and don't need to get transpiled.
		return
	}

	for singleRoot := combinedRoot.GetFirstChild(); singleRoot != nil; singleRoot = singleRoot.GetNext() {
		jutil.CheckState(singleRoot.IsScript())
		// only run the callbacks if any feature from the given featureSet exists in the script
		if TranspilationPasses_DoesScriptHaveAnyOfTheseFeatures(singleRoot, featuresToRunFor) {
			for _, callback := range callbacks {
				NodeTraversal_Traverse(compiler, singleRoot, callback)
			}
		}
	}
}

// TranspilationPasses_MaybeMarkFeatureAsTranspiledAway removes the given
// features from the FEATURE_SET prop of all scripts under root. Also removes
// from the compiler's featureset.
//
// Java: TranspilationPasses.maybeMarkFeatureAsTranspiledAway(AbstractCompiler, Node, Feature)
func TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(
	compiler AbstractCompiler, root *rhino.Node, feature parser.FeatureSet_Feature) {
	// We don't bother to do this if the compiler has halting errors, which avoids unnecessary
	// warnings from AstValidator warning that the features are still there.
	if !compiler.HasHaltingErrors() {
		compiler.MarkFeatureNotAllowed(feature)
		NodeUtil_RemoveFeatureFromAllScripts(root, feature, compiler)
	}
}

// TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway removes the given
// features from the FEATURE_SET prop of all scripts under root. Also removes
// from the compiler's featureset.
//
// Java: TranspilationPasses.maybeMarkFeaturesAsTranspiledAway(AbstractCompiler, Node, FeatureSet)
func TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(
	compiler AbstractCompiler, root *rhino.Node, transpiledFeatures *parser.FeatureSet) {
	// We don't bother to do this if the compiler has halting errors, which avoids unnecessary
	// warnings from AstValidator warning that the features are still there.
	if !compiler.HasHaltingErrors() {
		// remove the features from the compiler's featureSet and remove the features from every
		// script's featureset
		NodeUtil_RemoveFeaturesFromAllScripts(root, transpiledFeatures, compiler)
	}
}

// Java: TranspilationPasses.postTranspileCheckUnsupportedFeaturesRemoved(AbstractCompiler)
func TranspilationPasses_PostTranspileCheckUnsupportedFeaturesRemoved(compiler AbstractCompiler) {
	outputFeatures := compiler.GetOptions().GetOutputFeatureSet()
	currentFeatures := compiler.GetAllowableFeatures()
	// features modules, importMeta and Dynamic module import may exist in the output even though
	// unsupported
	if compiler.GetOptions().GetChunkOutputType() == CompilerOptions_ChunkOutputType_ES_MODULES {
		currentFeatures = currentFeatures.WithoutFeatures(
			parser.FeatureSet_Feature_MODULES,
			parser.FeatureSet_Feature_IMPORT_META,
			parser.FeatureSet_Feature_DYNAMIC_IMPORT,
		)
	}

	if outputFeatures.GetFeatures().IsEmpty() {
		// In some cases (e.g. targets built using `gen_closurized_js`), the outputFeatures is not
		// set. Only when set, confirm that the output featureSet is respected by JSCompiler.
		return
	}

	if !outputFeatures.Contains(currentFeatures) {
		// Confirm that the output featureSet is respected by JSCompiler.
		diff := currentFeatures.Without(outputFeatures)
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unsupported feature(s) leaked to output code:" + diff.GetFeatures().String()))
	}
}

// TranspilationPasses_CreateFeatureRemovalPass returns a pass that just
// removes features from the source scripts' FeatureSet and the compiler's
// featureset.
//
// Doing this indicates that the AST no longer contains uses of the features,
// or that they are no longer of concern for some other reason.
//
// Java: TranspilationPasses.createFeatureRemovalPass(String, FeatureSet)
func TranspilationPasses_CreateFeatureRemovalPass(
	passName string, featuresToRemove *parser.FeatureSet) *PassFactory {
	return PassFactory_Builder_().
		SetName(passName).
		SetInternalFactory(
			func(compiler AbstractCompiler) CompilerPass {
				return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
					TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(compiler, root, featuresToRemove)
				})
			}).
		Build()
}

// TranspilationPasses_CreatePostTranspileUnsupportedFeaturesRemovedCheck
// returns a pass that just checks that post-transpile only supported features
// exist in the code.
//
// Java: TranspilationPasses.createPostTranspileUnsupportedFeaturesRemovedCheck(String)
func TranspilationPasses_CreatePostTranspileUnsupportedFeaturesRemovedCheck(passName string) *PassFactory {
	return PassFactory_Builder_().
		SetName(passName).
		SetInternalFactory(
			func(compiler AbstractCompiler) CompilerPass {
				return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
					TranspilationPasses_PostTranspileCheckUnsupportedFeaturesRemoved(compiler)
				})
			}).
		Build()
}

// TranspilationPasses_CreateUnifiedFeatureRemovalPass creates a single pass
// to mark all transpiled features as removed from the source scripts'
// FeatureSet and the compiler's featureset.
//
// Doing this indicates that the AST no longer contains uses of the features,
// or that they are no longer of concern for some other reason.
//
// Java: TranspilationPasses.createUnifiedFeatureRemovalPass(String, CompilerOptions)
func TranspilationPasses_CreateUnifiedFeatureRemovalPass(
	passName string, options *CompilerOptions) *PassFactory {

	featuresToMarkRemoved := parser.FeatureSet_BARE_MINIMUM

	// The compiler doesn't transpile regex features.
	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_REGEXP_FLAG_D) {
		featuresToMarkRemoved = featuresToMarkRemoved.With(parser.FeatureSet_Feature_REGEXP_FLAG_D)
	}
	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_BIGINT) {
		featuresToMarkRemoved = featuresToMarkRemoved.With(parser.FeatureSet_Feature_BIGINT)
	}
	if options.NeedsTranspilationOf(parser.FeatureSet_Feature_NUMERIC_SEPARATOR) {
		// Numeric separators are flagged as present by the parser,
		// but never actually represented in the AST.
		// The only thing we need to do is mark them as not present in the AST.
		featuresToMarkRemoved = featuresToMarkRemoved.With(parser.FeatureSet_Feature_NUMERIC_SEPARATOR)
	}
	if options.GetChunkOutputType() != CompilerOptions_ChunkOutputType_ES_MODULES {
		// Default output mode of JSCompiler is a script, unless chunkOutputType is set to
		// `ES_MODULES` where each output chunk is an ES module.
		featuresToMarkRemoved = featuresToMarkRemoved.With(parser.FeatureSet_Feature_MODULES)
		// Since import.meta cannot be transpiled, it is passed-through when the output format
		// is a module. Otherwise it must be marked removed.
		featuresToMarkRemoved = featuresToMarkRemoved.With(parser.FeatureSet_Feature_IMPORT_META)
		// Dynamic imports are preserved for open source output only when the chunk output type is
		// ES_MODULES
		featuresToMarkRemoved = featuresToMarkRemoved.With(parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	}

	if options.NeedsTranspilationFrom(
		parser.FeatureSet_BARE_MINIMUM.WithFeatures(
			parser.FeatureSet_Feature_BINARY_LITERALS,
			parser.FeatureSet_Feature_OCTAL_LITERALS,
			parser.FeatureSet_Feature_REGEXP_FLAG_U,
			parser.FeatureSet_Feature_REGEXP_FLAG_Y)) {
		// Binary and octal literals are effectively transpiled by the parser.
		// There's no transpilation we can do for the new regexp flags.
		featuresToMarkRemoved =
			featuresToMarkRemoved.WithFeatures(
				parser.FeatureSet_Feature_BINARY_LITERALS,
				parser.FeatureSet_Feature_OCTAL_LITERALS,
				parser.FeatureSet_Feature_REGEXP_FLAG_U,
				parser.FeatureSet_Feature_REGEXP_FLAG_Y)
	}

	if options.NeedsTranspilationFrom(parser.FeatureSet_ES5) {
		// this means we're transpiling to ES3 output
		featuresToMarkRemoved =
			featuresToMarkRemoved.WithFeatures(
				// TODO(b/354075108): Stop tracking these 2 features. These don't get transpiled and
				// are not supported in ES3. But they get possibly inlined / renamed by the optimizer.
				parser.FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS, // does not get transpiled
				parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES,      // does not get transpiled
				// GETTERs and SETTERs get reported in lateConvertEs6ToEs3 for ES3 output. If
				// we're here it means that GETTERs and SETTERs don't exist.
				parser.FeatureSet_Feature_GETTER,
				parser.FeatureSet_Feature_SETTER)
	}

	// these ES5 features are transpiled away unconditionally regardless of output level.
	featuresToMarkRemoved =
		featuresToMarkRemoved.WithFeatures(
			parser.FeatureSet_Feature_STRING_CONTINUATION, // transpiled away during parsing
			parser.FeatureSet_Feature_TRAILING_COMMA,      // transpiled away during Normalization
		)

	// ES_NEXT_RUNTIME and ES_UNSTABLE_RUNTIME are not "real" features. They only exist to ensure
	// that FeatureSet.ES_NEXT and FeatureSet.ES_UNSTABLE have distinct version names as some
	// polyfills target "es_next" and "es_unstable" for their fromLang.
	featuresToMarkRemoved =
		featuresToMarkRemoved.WithFeatures(parser.FeatureSet_Feature_ES_NEXT_RUNTIME, parser.FeatureSet_Feature_ES_UNSTABLE_RUNTIME)

	return TranspilationPasses_CreateFeatureRemovalPass(passName, featuresToMarkRemoved)
}
