// Port of com/google/javascript/jscomp/CompilerOptions.java @ closure-compiler v20261021

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
	"io"
	"os"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

var defaultConformanceRemoveRegexFromPath = jutil.Pattern_Compile(
	"^((.*/)?google3/)?(/?(blaze|bazel)-out/[^/]+/(bin|genfiles)/)?")

// CompilerOptions holds the compiler options.
//
// Go-only: language-adapter — the fields are declared together here (Java
// interleaves them with the accessors) and keep their Java declaration
// order; their initializers run first in NewCompilerOptions, as Java's
// field initializers run before the constructor body. Value types: a Java
// String that may be null is a *string; Optional<Boolean> is a *bool and
// Optional<FeatureSet>/Optional<Pattern> are nil-able pointers (nil is
// absent); ImmutableSet/Set are *jutil.LinkedHashSet, ImmutableMap/Map are
// *jutil.LinkedHashMap and ImmutableList/List are slices; the @define
// replacement values are bool, int32 (Integer), float64 (Double) or string.
type CompilerOptions struct {
	// Should the compiled output start with "'use strict';"?
	emitUseStrict *bool
	// The JavaScript language version accepted.
	languageIn CompilerOptions_LanguageMode
	// The JavaScript features that are allowed to be in the output.
	outputFeatureSet             *parser.FeatureSet
	experimentalOutputFeatureSet *CompilerOptions_ExperimentalOutputFeatureSet
	languageOutIsDefaultStrict   *bool
	// The builtin set of externs to be used
	environment CompilerOptions_Environment
	// Represents BrowserFeaturesetYear to use for compilation
	browserFeaturesetYear *CompilerOptions_BrowserFeaturesetYear
	// Instrument code for the purpose of collecting coverage data - restrict to coverage pass only,
	// and skip all other passes.
	instrumentForCoverageOnly  bool
	typedAstOutputFile         *jutil.Path
	mergedPrecompiledLibraries bool
	// Whether to infer consts. This should not be configurable by external clients. This is a
	// transitional flag for a new type of const analysis.
	//
	// TODO(nicksantos): Remove this option.
	inferConsts bool
	// Whether the compiler should assume that a function's "this" value never needs coercion (for
	// example in non-strict "null" or "undefined" will be coerced to the global "this" and primitives
	// to objects).
	assumeStrictThis           bool
	preserveDetailedSourceInfo bool
	preserveNonJSDocComments   bool
	continueAfterErrors        bool
	incrementalCheckMode       CompilerOptions_IncrementalCheckMode
	parseJsDocDocumentation    parsing.Config_JsDocParsing
	printExterns               bool
	// Even if checkTypes is disabled, clients such as IDEs might want to still infer types.
	inferTypes bool
	// Configures the compiler to skip as many passes as possible. If transpilation is requested, it
	// will be run, but all others passes will be skipped.
	skipNonTranspilationPasses bool
	// Configures the compiler to run expensive validity checks after every pass. Only intended for
	// internal development.
	devMode CompilerOptions_DevMode
	// Configures the compiler to log a hash code of the AST after every pass. Only intended for
	// internal development.
	checkDeterminism bool
	// --------------------------------
	// Input Options
	// --------------------------------
	dependencyOptions *DependencyOptions
	// Returns localized replacement for MSG_* variables
	messageBundle MessageBundle
	// Whether we should report an error if a message is absent from a bundle.
	strictMessageReplacement bool
	// --------------------------------
	// Checks
	// --------------------------------
	// Checks that all symbols are defined
	// TODO(tbreisacher): Remove this and deprecate the corresponding setter.
	checkSymbols bool
	// Checks for suspicious statements that have no effect
	checkSuspiciousCode bool
	// Checks types on expressions
	checkTypes bool
	// A set of extra annotation names which are accepted and silently ignored when encountered in a
	// source file. Defaults to null which has the same effect as specifying an empty set.
	extraAnnotationNames *jutil.LinkedHashSet[string]
	// TODO(bradfordcsmith): Investigate how can we use multi-threads as default.
	numParallelThreads int
	// --------------------------------
	// Optimizations
	// --------------------------------
	// Folds constants (e.g. (2 + 3) to 5)
	foldConstants bool
	// Remove assignments to values that can not be referenced
	deadAssignmentElimination bool
	// Remove assignments to properties that cannot be referenced
	deadPropertyAssignmentElimination base.Tri
	// Inlines constants (symbols that are all CAPS)
	inlineConstantVars bool
	// For projects that want to avoid the creation of giant functions after inlining.
	maxFunctionSizeAfterInlining int
	// Inlines properties
	inlineProperties bool
	// Move code to a deeper chunk
	crossChunkCodeMotion bool
	// Don't generate stub functions when moving methods deeper.
	//
	// Note, switching on this option may break existing code that depends on enumerating prototype
	// methods for mixin behavior, such as goog.object.extend, since the prototype assignments will be
	// removed from the parent chunk and moved to a later chunk.
	crossChunkCodeMotionNoStubMethods bool
	// Whether when chunk B depends on chunk A and chunk B declares a symbol, this symbol can be seen
	// in A after B has been loaded. This is often true, but may not be true when loading code using
	// nested eval.
	parentChunkCanSeeSymbolsDeclaredInChildren bool
	// Merge two variables together as one.
	coalesceVariableNames bool
	// Lower block-scoped declarations to function-scoped declarations.
	optimizeLetAndConst bool
	// Assume the global scope is lexically isolated (e.g. via IIFE or Module).
	assumeGlobalScopeIsIsolated bool
	// Move methods to a deeper chunk
	crossChunkMethodMotion bool
	// Inlines trivial getters
	inlineGetters bool
	// Inlines variables
	inlineVariables bool
	// Inlines variables
	inlineLocalVariables bool
	// TODO(user): This is temporary. Once flow sensitive inlining is stable
	// Remove this.
	flowSensitiveInlineVariables bool
	// Removes code associated with unused global names
	smartNameRemoval bool
	// Extracts common prototype member declarations
	extractPrototypeMemberDeclarations CompilerOptions_ExtractPrototypeMemberDeclarationsMode
	// Removes unused member prototypes
	removeUnusedPrototypeProperties bool
	// Removes unused member properties
	removeUnusedClassProperties bool
	// Removes unused variables
	removeUnusedVars bool
	// Removes unused variables in local scope.
	removeUnusedLocalVars bool
	// Collapses multiple variable declarations into one
	collapseVariableDeclarations bool
	// Collapses anonymous function declarations into named function declarations
	collapseAnonymousFunctions bool
	// Aliases string literals to global instances, to reduce code size.
	aliasStringsMode CompilerOptions_AliasStringsMode
	// Print string usage as part of the compilation log.
	outputJsStringUsage bool
	// Converts quoted property accesses to dot syntax (a['b'] &rarr; a.b)
	convertToDottedProperties bool
	// Reduces the size of common function expressions.
	rewriteFunctionExpressions bool
	// Remove unused function arguments, remove unused return values, and inlines constant parameters.
	optimizeCalls bool
	// Removes trivial constructors where ES class implicit constructors are sufficient.
	optimizeESClassConstructors bool
	// Use type information to enable additional optimization opportunities.
	useTypesForLocalOptimization           bool
	useSizeHeuristicToStopOptimizationLoop bool
	// Do up to this many iterations of the optimization loop. Setting this field to some small
	// number, say 3 or 4, allows a large project to build faster, but sacrifice some code size.
	optimizationLoopMaxIterations int
	// --------------------------------
	// Renaming
	// --------------------------------
	// Controls which variables get renamed.
	variableRenaming VariableRenamingPolicy
	// Controls which properties get renamed.
	propertyRenaming PropertyRenamingPolicy
	// Controls if property renaming only compilation mode is needed.
	propertyRenamingOnlyCompilationMode bool
	// Controls label renaming.
	labelRenaming bool
	// Reserve property names on the global this object.
	reserveRawExports bool
	// Use a renaming heuristic with better stability across source changes. With this option each
	// symbol is more likely to receive the same name between builds. The cost may be a slight
	// increase in code size.
	preferStableNames bool
	// Generate pseudo names for variables and properties for debugging purposes.
	generatePseudoNames bool
	// Specifies a prefix for all globals
	renamePrefix *string
	// Specifies the name of an object that will be used to store all non-extern globals.
	renamePrefixNamespace *string
	// Used by tests of the RescopeGlobalSymbols pass to avoid having declare 2 chunks in simple
	// cases.
	renamePrefixNamespaceAssumeCrossChunkNames  bool
	optimizeLocalAccessForGlobalSymbolNamespace CompilerOptions_OptimizeLocalAccess
	collapsePropertiesLevel                     CompilerOptions_PropertyCollapseLevel
	// Split object literals into individual variables when possible.
	collapseObjectLiterals bool
	// Devirtualize prototype method by rewriting them to be static calls that take the this pointer
	// as their first argument
	devirtualizeMethods bool
	// Use @nosideeffects annotations, function bodies and name graph to determine if calls have side
	// effects.
	computeFunctionSideEffects bool
	// Rename properties to disambiguate between unrelated fields based on type information.
	disambiguateProperties bool
	// Rename unrelated properties to the same name to reduce code size.
	ambiguateProperties bool
	// Input sourcemap files, indexed by the JS files they refer to
	inputSourceMaps *jutil.LinkedHashMap[string, *SourceMapInput]
	// Input variable renaming map.
	//
	// During renaming, the compiler uses this map and the inputPropertyMap to try to preserve
	// renaming mappings from a previous compilation. The application is delta encoding: keeping the
	// diff between consecutive versions of one's code small. The compiler does NOT guarantee to
	// respect these maps; projects should not use these maps to prevent renaming or to select
	// particular names. Point questioners to this post:
	// http://closuretools.blogspot.com/2011/01/property-by-any-other-name-part-3.html
	inputVariableMap *VariableMap
	// Input property renaming map.
	inputPropertyMap *VariableMap
	// Whether to export test functions.
	exportTestFunctions bool
	// Shared name generator
	nameGenerator NameGenerator
	// --------------------------------
	// Special-purpose alterations
	// --------------------------------
	// Replace UI strings with chrome.i18n.getMessage calls. Used by Chrome extensions/apps.
	replaceMessagesWithChromeI18n bool
	tcProjectId                   *string
	// A CodingConvention to use during the compile.
	codingConvention          CodingConvention
	syntheticBlockStartMarker *string
	syntheticBlockEndMarker   *string
	// Compiling locale
	locale *string
	// If true, then perform localization passes as late as possible.
	//
	// At the moment this only affects when `ReplaceMessages` does the work of inserting the
	// localized form of declared messages from the message bundle.
	//
	// TODO(johnlenz): Use this option to control late substitution of other locale-specific code.
	doLateLocalization bool
	// Sets the special "COMPILED" value to true
	markAsCompiled bool
	// Processes goog.provide() and goog.require() calls
	closurePass bool
	// Do not strip goog.provide()/goog.require() calls from the code.
	preserveClosurePrimitives bool
	// Processes AngularJS-specific annotations
	angularPass bool
	// Processes Polymer code
	polymerPass bool
	// Processes cr.* functions
	chromePass bool
	// Processes the output of J2CL
	j2clPassMode                CompilerOptions_J2clPassMode
	j2clMinifierEnabled         bool
	j2clMinifierPruningManifest *string
	// Remove goog.abstractMethod assignments and @abstract methods.
	removeAbstractMethods bool
	// Remove goog.asserts calls.
	removeClosureAsserts bool
	// Remove J2CL assert calls.
	removeJ2clAsserts bool
	// Gather CSS names (requires closurePass)
	gatherCssNames bool
	// Names of types to strip
	stripTypes *jutil.LinkedHashSet[string]
	// Name suffixes that determine which variables and properties to strip
	stripNameSuffixes *jutil.LinkedHashSet[string]
	// Name prefixes that determine which variables and properties to strip
	stripNamePrefixes *jutil.LinkedHashSet[string]
	// Custom passes
	customPasses *customPassMap
	// Replacements for @defines. Will be Boolean, Numbers, or Strings
	defineReplacements *jutil.LinkedHashMap[string, defineReplacement]
	// What kind of processing to do for goog.tweak functions.
	tweakProcessing CompilerOptions_TweakProcessing
	// Move top-level function declarations to the top
	rewriteGlobalDeclarationsForTryCatchWrapping bool
	checksOnly                                   bool
	outputJs                                     CompilerOptions_OutputJs
	generateExports                              bool
	exportLocalPropertyDefinitions               bool
	// Map used in the renaming of CSS class names.
	cssRenamingMap CssRenamingMap
	// Skiplist used in the renaming of CSS class names.
	cssRenamingSkiplist *jutil.LinkedHashSet[string]
	// Replace id generators
	replaceIdGenerators bool
	// true by default for legacy reasons.
	// Id generators to replace.
	idGenerators *jutil.LinkedHashMap[string, RenamingMap]
	// Hash function to use for xid generation.
	xidHashFunction Xid_HashFunction
	// Hash function to use for generating chunk ID.
	chunkIdHashFunction Xid_HashFunction
	// A previous map of ids (serialized to a string by a previous compile). This will be used as a
	// hint during the ReplaceIdGenerators pass, which will attempt to reuse the same ids.
	idGeneratorsMapSerialized *string
	// Configuration strings
	replaceStringsFunctionDescriptions []string
	replaceStringsPlaceholderToken     string
	// List of properties that we report invalidation errors for.
	propertiesThatMustDisambiguate *jutil.LinkedHashSet[string]
	// Whether to enforce that the @requireInlining annotation works.
	validateRequireInliningAnnotation base.Tri
	// Rewrite CommonJS modules so that they can be concatenated together.
	processCommonJSModules bool
	// CommonJS module prefix.
	moduleRoots []string
	// Inject polyfills
	rewritePolyfills bool
	// Isolates injected polyfills from the global scope.
	isolatePolyfills bool
	// Configures polyfill injection of all polyfills newer than the given language mode, regardless
	// of whether they are actually referenced in the code being compiled.
	//
	// How is this different from --inject_library? --inject_library takes a specific JS file
	// defined in the jscomp/js package, and injects that and all of its dependencies. This option
	// takes a language mode.
	//
	// For example, the common use of `--inject_library` is passing `--inject_library=es6_runtime`.
	// es6_runtime.js is a file that contains both all the ES6 polyfills and also ES6 transpilation
	// utilities (library code referenced by transpiled code). So
	//
	// <ul>
	//   <li>es6_runtime includes transpilation utilities, not just polyfills.
	//       `--inject_polyfills_newer_than=ES5` will not add transpilation utilities, just the
	//       polyfills,
	//   <li>if someone forgot to add a new polyfill to es6_runtime.js, then `--inject_library` will
	//       not add it, but `--inject_polyfills_newer_than=ES5` will add it.
	// </ul>
	injectPolyfillsNewerThan *CompilerOptions_LanguageMode
	// Whether to emit the safe but slower/more verbose transpilation of class super, when the
	// compiler cannot statically determine whether the superclass is an ES6 class.
	es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation
	// Whether to instrument reentrant functions for AsyncContext.
	instrumentAsyncContext bool
	// Runtime libraries to always inject.
	forceLibraryInjection []string
	// Runtime libraries to never inject.
	runtimeLibraryMode                   RuntimeJsLibManager_RuntimeLibraryMode
	assumeForwardDeclaredForMissingTypes bool
	// A Set of goog.requires to be removed. If null, ALL of the unused goog.requires will be counted
	// as a candidate to be removed.
	unusedImportsToRemove *jutil.LinkedHashSet[string]
	// --------------------------------
	// Output options
	// --------------------------------
	// Do not strip closure-style type annotations from code.
	preserveTypeAnnotations bool
	// To distinguish between gents and non-gents mode so that we can turn off checking the sanity of
	// the source location of comments, and also provide a different mode for comment printing between
	// those two.
	gentsMode bool
	// Output in pretty indented format
	prettyPrint bool
	// Line break the output a bit more aggressively
	lineBreak bool
	// Prints a separator comment before each JS script
	printInputDelimiter bool
	// The string to use as the separator for printInputDelimiter
	inputDelimiter string
	// A directory into which human readable debug log files can be written.
	//
	// null indicates that no such files should be written.
	debugLogDirectory *jutil.Path
	// A comma separated list of strings used to filter which debug logs are written. If a filter is
	// specified, only log files with a path name that contains at least one of the filter Strings are
	// written.
	//
	// null indicates that no such files should be written.
	debugLogFilter          *string
	serializeExtraDebugInfo bool
	// Go-only: language-adapter — Function<OutputStream, OutputStream> and
	// Function<InputStream, InputStream>; nil is Java null.
	stateCompressionWrapper   func(io.Writer) io.Writer
	stateDecompressionWrapper func(io.Reader) io.Reader
	// Whether to write keyword properties as foo['class'] instead of foo.class; needed for IE8.
	quoteKeywordProperties bool
	preferSingleQuotes     bool
	trustedStrings         bool
	// Should only be used when debugging compiler bugs.
	printSourceAfterEachPass bool
	// Used to narrow down the printed source when overall input size is large. If these are both
	// empty the entire source is printed.
	filesToPrintAfterEachPassRegexList  []string
	chunksToPrintAfterEachPassRegexList []string
	qnameUsesToPrintAfterEachPassList   []string
	tracer                              CompilerOptions_TracerMode
	tracerOutput                        *jutil.Path
	colorizeErrorOutput                 bool
	errorFormat                         ErrorFormat
	warningsGuard                       *ComposeWarningsGuard
	summaryDetailLevel                  int
	lineLengthThreshold                 int
	// Whether to use the original names of nodes in the code output. This option is only really
	// useful when using the compiler to print code meant to check in to source.
	useOriginalNamesInOutput bool
	// --------------------------------
	// Special Output Options
	// --------------------------------
	// The output path for the created externs file.
	externExportsPath     *string
	extraReportGenerators []SortingErrorManager_ErrorReportGenerator
	// --------------------------------
	// Debugging Options
	// --------------------------------
	// The output path for the source map.
	sourceMapOutputPath *string
	// Should we gather information needed to generate a source map even when sourceMapOutputPath
	// indicates otherwise?
	//
	// This is necessary in cases where we're doing a partial compilation, because the source map
	// output path may not be known until the final stage, but we must still gather information in the
	// early stages.
	shouldAlwaysGatherSourceMapInfo bool
	// The detail level for the generated source map.
	sourceMapDetailLevel SourceMap_DetailLevel
	// The source map file format
	sourceMapFormat SourceMap_Format
	// Whether to parse inline source maps.
	parseInlineSourceMaps bool
	// Whether to apply input source maps to the output, i.e. map back to original inputs from input
	// files that have source maps applied to them.
	applyInputSourceMaps bool
	// Whether to resolve source mapping annotations. Cannot do this in an appengine or js environment
	// since we don't have access to the filesystem.
	resolveSourceMapAnnotations bool
	sourceMapLocationMappings   []SourceMap_LocationMapping
	// Whether to include full file contents in the source map.
	sourceMapIncludeSourcesContent bool
	// Charset to use when generating code. If null, then output ASCII.
	outputCharset *jutil.Charset
	// When set, assume that apparently side-effect free code is meaningful.
	protectHiddenSideEffects bool
	// Ignore the possibility that getter invocations (gets) can have side-effects and that the
	// results of gets can be side-effected by local state mutations.
	//
	// When false, it doesn't necessarily mean that all gets are considered side-effectful
	// or side-effected. Gets that can be proven to be pure may still be considered as such.
	//
	// Recall that object-spread is capable of triggering getters. Since the syntax doesn't
	// explicitly specify a property, it is essentailly impossible to prove it has no side-effects
	// without this assumption.
	assumeGettersArePure bool
	// When true, the compiler assumes that all non-quoted properties are only referenced by
	// their non-quoted names (not reflectively), or otherwise are defined in the externs. This is
	// basically the set of assumptions that PropertyRenamingPolicy.ALL_UNQUOTED makes, but
	// has repercussions beyond just property renaming.
	//
	// When false, the compiler does not attempt to do any property optimizations.
	//
	// NOTE: lharker - I added this in 2025 specifically for @closureUnaware support. Previously
	// this was always implicitly true. I've tried to update various parts of the compiler to check
	// this field, rather than assuming property definitions are known, but it's quite likely there
	// are still some silent assumptions somewhere. So this is still experimental and needs more
	// validation on real code.
	assumePropertiesAreStaticallyAnalyzable bool
	// Assume that static (class-side) inheritance is not being used and that static methods will not
	// be referenced via `this` or through subclasses.
	//
	// When true, the compiler is free to make unsafe (breaking) optimizations to code that
	// depends on static inheritance. These optimizations represent a substantial code-size reduction
	// for older projects and therefore cannot be unilaterally disabled. true was the
	// long-standing implicit assumption before static inheritance came about in ES2015.
	//
	// Example of what may break if this flag is true:
	//
	// <pre>{@code
	// class Parent {
	//   static method() { }
	// }
	//
	// class Child extends Parent { }
	//
	// Child.method(); // `method` will not be defined.
	// }</pre>
	assumeStaticInheritanceIsNotUsed bool
	// Handler for compiler warnings and errors.
	errorHandler                       ErrorHandler
	instrumentForCoverageOption        CompilerOptions_InstrumentOption
	productionInstrumentationArrayName string
	// List of conformance configs to use in CheckConformance.
	conformanceConfigs       []*ConformanceConfig
	conformanceReportingMode CompilerOptions_ConformanceReportingMode
	// Remove the first match of this regex from any paths when checking conformance allowlists.
	//
	// You can use this to make absolute paths relative to the root of your source tree. This is
	// useful to work around CI and build systems that use absolute paths.
	conformanceRemoveRegexFromPath *jutil.Pattern
	// For use in CompilationLevel#WHITESPACE_ONLY mode, when using goog.module.
	wrapGoogModulesForWhitespaceOnly bool
	// Print all configuration options to stderr after the compiler is initialized.
	printConfig bool
	// Are the input files written for strict mode?
	isStrictModeInput *bool
	// Store four unique states:
	//  - run module rewriting before typechecking
	//  - run module rewriting after typechecking
	//  - don't run module rewriting, but if it's reenabled, run it before typechecking
	//  - don't run module rewriting, but if it's reenabled, run it after typechecking
	rewriteModulesBeforeTypechecking bool
	enableModuleRewriting            bool
	// Which algorithm to use for locating ES6 and CommonJS modules
	moduleResolutionMode ModuleLoader_ResolutionMode
	// Map of prefix replacements for use when moduleResolutionMode is {@link
	// ResolutionMode#BROWSER_WITH_TRANSFORMED_PREFIXES}.
	browserResolverPrefixReplacements *jutil.LinkedHashMap[string, string]
	pathEscaper                       ModuleLoader_PathEscaper
	// Which entries to look for in package.json files when processing modules
	packageJsonEntryNames  []string
	allowDynamicImport     bool
	dynamicImportAlias     *string
	chunkOutputType        CompilerOptions_ChunkOutputType
	unknownDefinesToIgnore []string
	inlineFunctionsLevel   CompilerOptions_Reach
	enableZonesDefineName  *string
	zoneInputPattern       *jutil.Pattern
	es6ModuleTranspilation CompilerOptions_Es6ModuleTranspilation
	// Go-only: fast mode (Go-to-Go deterministic across thread counts, relaxed Java map order)
	fastMode bool
}

// Initializes compiler options. All options are disabled by default.
//
// Command-line frontends to the compiler should set these properties like a builder.
//
// Java: CompilerOptions()
func NewCompilerOptions() *CompilerOptions {
	o := &CompilerOptions{}
	// Field initializers.
	o.experimentalOutputFeatureSet = nil
	o.optimizeLocalAccessForGlobalSymbolNamespace = CompilerOptions_OptimizeLocalAccess_DISABLED
	o.inferConsts = true
	o.incrementalCheckMode = CompilerOptions_IncrementalCheckMode_OFF
	o.parseJsDocDocumentation = parsing.Config_JsDocParsing_TYPES_ONLY
	o.dependencyOptions = DependencyOptions_None()
	o.numParallelThreads = 1
	o.fastMode = os.Getenv("JSGO_FAST_MODE") == "1"
	o.useSizeHeuristicToStopOptimizationLoop = true
	o.j2clMinifierEnabled = true
	o.removeJ2clAsserts = true
	o.replaceIdGenerators = true
	o.moduleRoots = []string{ModuleLoader_DEFAULT_FILENAME_PREFIX}
	o.es6SubclassTranspilation = CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE
	o.forceLibraryInjection = []string{}
	o.runtimeLibraryMode = RuntimeJsLibManager_RuntimeLibraryMode_INJECT
	o.inputDelimiter = "// Input %num%"
	o.filesToPrintAfterEachPassRegexList = []string{}
	o.chunksToPrintAfterEachPassRegexList = []string{}
	o.qnameUsesToPrintAfterEachPassList = []string{}
	o.warningsGuard = NewComposeWarningsGuard()
	o.summaryDetailLevel = 1
	o.lineLengthThreshold = CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD
	o.extraReportGenerators = []SortingErrorManager_ErrorReportGenerator{}
	o.sourceMapDetailLevel = SourceMap_DetailLevel_ALL
	o.sourceMapFormat = SourceMap_Format_DEFAULT
	o.parseInlineSourceMaps = true
	o.resolveSourceMapAnnotations = true
	o.sourceMapLocationMappings = []SourceMap_LocationMapping{}
	o.assumeGettersArePure = true
	o.assumePropertiesAreStaticallyAnalyzable = true
	o.assumeStaticInheritanceIsNotUsed = true
	o.conformanceConfigs = []*ConformanceConfig{}
	o.conformanceReportingMode = CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
	o.conformanceRemoveRegexFromPath = defaultConformanceRemoveRegexFromPath
	o.wrapGoogModulesForWhitespaceOnly = true
	o.enableModuleRewriting = true
	o.unknownDefinesToIgnore = []string{}
	o.es6ModuleTranspilation = CompilerOptions_Es6ModuleTranspilation_COMPILE
	// Constructor body.
	// Accepted language
	o.languageIn = CompilerOptions_LanguageMode_STABLE_IN

	// Which environment to use
	o.environment = CompilerOptions_Environment_BROWSER
	o.browserResolverPrefixReplacements = jutil.NewLinkedHashMap[string, string]()

	// Modules
	o.moduleResolutionMode = ModuleLoader_ResolutionMode_BROWSER
	o.packageJsonEntryNames = []string{"browser", "module", "main"}
	o.pathEscaper = ModuleLoader_PathEscaper_ESCAPE
	o.rewriteModulesBeforeTypechecking = false
	o.enableModuleRewriting = true

	// Checks
	o.skipNonTranspilationPasses = false
	o.devMode = CompilerOptions_DevMode_OFF
	o.checkDeterminism = false
	o.checkSymbols = false
	o.checkSuspiciousCode = false
	o.checkTypes = false
	o.computeFunctionSideEffects = false
	o.extraAnnotationNames = nil

	// Optimizations
	o.foldConstants = false
	o.coalesceVariableNames = false
	o.optimizeLetAndConst = true
	o.assumeGlobalScopeIsIsolated = false
	o.deadAssignmentElimination = false
	o.deadPropertyAssignmentElimination = base.Tri_UNKNOWN
	o.inlineConstantVars = false
	o.inlineFunctionsLevel = CompilerOptions_Reach_NONE
	o.maxFunctionSizeAfterInlining = CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING
	o.assumeStrictThis = false
	o.inlineProperties = false
	o.crossChunkCodeMotion = false
	o.parentChunkCanSeeSymbolsDeclaredInChildren = false
	o.crossChunkMethodMotion = false
	o.inlineGetters = false
	o.inlineVariables = false
	o.inlineLocalVariables = false
	o.smartNameRemoval = false
	o.extractPrototypeMemberDeclarations = CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF
	o.removeUnusedPrototypeProperties = false
	o.removeUnusedClassProperties = false
	o.removeUnusedVars = false
	o.removeUnusedLocalVars = false
	o.collapseVariableDeclarations = false
	o.collapseAnonymousFunctions = false
	o.aliasStringsMode = CompilerOptions_AliasStringsMode_NONE
	o.outputJsStringUsage = false
	o.convertToDottedProperties = false
	o.rewriteFunctionExpressions = false
	o.validateRequireInliningAnnotation = base.Tri_UNKNOWN

	// Renaming
	o.variableRenaming = VariableRenamingPolicy_OFF
	o.propertyRenaming = PropertyRenamingPolicy_OFF
	o.propertyRenamingOnlyCompilationMode = false
	o.labelRenaming = false
	o.generatePseudoNames = false
	o.preferStableNames = false
	o.renamePrefix = nil
	o.collapsePropertiesLevel = CompilerOptions_PropertyCollapseLevel_NONE
	o.collapseObjectLiterals = false
	o.devirtualizeMethods = false
	o.disambiguateProperties = false
	o.ambiguateProperties = false
	o.exportTestFunctions = false
	o.nameGenerator = NewDefaultNameGenerator()

	// Alterations

	o.syntheticBlockStartMarker = nil
	o.syntheticBlockEndMarker = nil
	o.locale = nil
	o.doLateLocalization = false
	o.markAsCompiled = false
	o.closurePass = false
	o.preserveClosurePrimitives = false
	o.angularPass = false
	o.polymerPass = false
	o.j2clPassMode = CompilerOptions_J2clPassMode_AUTO
	o.j2clMinifierEnabled = true
	o.removeAbstractMethods = false
	o.removeClosureAsserts = false
	o.stripTypes = jutil.NewLinkedHashSet[string]()
	o.stripNameSuffixes = jutil.NewLinkedHashSet[string]()
	o.stripNamePrefixes = jutil.NewLinkedHashSet[string]()
	o.customPasses = nil
	o.defineReplacements = jutil.NewLinkedHashMap[string, defineReplacement]()
	o.tweakProcessing = CompilerOptions_TweakProcessing_OFF
	o.rewriteGlobalDeclarationsForTryCatchWrapping = false
	o.checksOnly = false
	o.outputJs = CompilerOptions_OutputJs_NORMAL
	o.generateExports = true
	o.exportLocalPropertyDefinitions = true
	o.cssRenamingMap = nil
	o.cssRenamingSkiplist = nil
	o.idGenerators = jutil.NewLinkedHashMap[string, RenamingMap]()
	o.replaceStringsFunctionDescriptions = []string{}
	o.replaceStringsPlaceholderToken = ""
	o.propertiesThatMustDisambiguate = jutil.NewLinkedHashSet[string]()
	o.inputSourceMaps = jutil.NewLinkedHashMap[string, *SourceMapInput]()

	o.instrumentForCoverageOption = CompilerOptions_InstrumentOption_NONE
	o.productionInstrumentationArrayName = ""

	// Output
	o.preserveTypeAnnotations = false
	o.gentsMode = false
	o.printInputDelimiter = false
	o.prettyPrint = false
	o.lineBreak = false
	o.tracer = CompilerOptions_TracerMode_OFF
	o.colorizeErrorOutput = false
	o.errorFormat = ErrorFormat_FULL
	o.chunkOutputType = CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE

	// Debugging
	o.errorHandler = nil
	o.printSourceAfterEachPass = false
	o.strictMessageReplacement = false
	return o
}

// Java: CompilerOptions.UNLIMITED_FUN_SIZE_AFTER_INLINING
const CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING = -1

// Java: CompilerOptions.setSkipTranspilationAndCrash(boolean)
func (o *CompilerOptions) SetSkipTranspilationAndCrash(value bool) {
}

// Skip all possible passes, to make the compiler as fast as possible.
//
// Java: CompilerOptions.skipAllCompilerPasses()
func (o *CompilerOptions) SkipAllCompilerPasses() {
	o.SetSkipNonTranspilationPasses(true)
}

// Sets the level of parallelism for compilation passes that can exploit multi-threading.
//
// Some compiler passes may take advantage of multi-threading, for example, parsing inputs.
// This sets the level of parallelism. The compiler will not start more than this number of
// threads.
//
// @param parallelism up to this number of parallel threads may be created.
//
// Java: CompilerOptions.setNumParallelThreads(int)
func (o *CompilerOptions) SetNumParallelThreads(parallelism int) {
	o.numParallelThreads = parallelism
}

// Java: CompilerOptions.getNumParallelThreads()
func (o *CompilerOptions) GetNumParallelThreads() int {
	return o.numParallelThreads
}

// SetFastMode enables or disables Go-to-Go deterministic fast mode.
func (o *CompilerOptions) SetFastMode(fastMode bool) {
	o.fastMode = fastMode
}

// isFastMode returns whether Go-to-Go deterministic fast mode is enabled.
func (o *CompilerOptions) isFastMode() bool {
	return o != nil && o.fastMode
}

// CompilerOptions_IsFastMode returns whether Go-to-Go deterministic fast mode is enabled.
func CompilerOptions_IsFastMode(o *CompilerOptions) bool {
	return o != nil && o.fastMode
}

// Java: CompilerOptions.setCodingConvention(CodingConvention)
func (o *CompilerOptions) SetCodingConvention(codingConvention CodingConvention) {
	o.codingConvention = codingConvention
}

// Java: CompilerOptions.getCodingConvention()
func (o *CompilerOptions) GetCodingConvention() CodingConvention {
	return o.codingConvention
}

// Should the compiler print its configuration options to stderr when they are initialized?
//
// Default false.
//
// Java: CompilerOptions.setPrintConfig(boolean)
func (o *CompilerOptions) SetPrintConfig(printConfig bool) {
	o.printConfig = printConfig
}

// Java: CompilerOptions.shouldPrintConfig()
func (o *CompilerOptions) ShouldPrintConfig() bool {
	return o.printConfig
}

// Java: CompilerOptions.getBrowserFeaturesetYearObject()
func (o *CompilerOptions) GetBrowserFeaturesetYearObject() *CompilerOptions_BrowserFeaturesetYear {
	return o.browserFeaturesetYear
}

// Java: CompilerOptions.setBrowserFeaturesetYear(int)
func (o *CompilerOptions) SetBrowserFeaturesetYear(year int) {
	o.SetBrowserFeaturesetYearOfBrowserFeaturesetYear(CompilerOptions_BrowserFeaturesetYear_From(year))
}

// Java: CompilerOptions.setBrowserFeaturesetYear(BrowserFeaturesetYear)
func (o *CompilerOptions) SetBrowserFeaturesetYearOfBrowserFeaturesetYear(year CompilerOptions_BrowserFeaturesetYear) {
	o.browserFeaturesetYear = year.Ptr()
	o.browserFeaturesetYear.SetDependentValuesFromYear(o)
}

// Sets ECMAScript version to use.
//
// Java: CompilerOptions.setLanguage(LanguageMode)
func (o *CompilerOptions) SetLanguage(language CompilerOptions_LanguageMode) {
	jutil.CheckState(language != CompilerOptions_LanguageMode_NO_TRANSPILE)
	o.SetLanguageIn(language)
	o.SetLanguageOut(language)
}

// Sets ECMAScript version to use for the input. If you are not transpiling from one version to
// another, use #setLanguage instead.
//
// Java: CompilerOptions.setLanguageIn(LanguageMode)
func (o *CompilerOptions) SetLanguageIn(languageIn CompilerOptions_LanguageMode) {
	jutil.CheckState(languageIn != CompilerOptions_LanguageMode_NO_TRANSPILE)
	if languageIn == CompilerOptions_LanguageMode_STABLE {
		o.languageIn = CompilerOptions_LanguageMode_STABLE_IN
	} else {
		o.languageIn = languageIn
	}
}

// Java: CompilerOptions.getLanguageIn()
func (o *CompilerOptions) GetLanguageIn() CompilerOptions_LanguageMode {
	return o.languageIn
}

// Sets ECMAScript version to use for the output.
//
// If you are not transpiling from one version to another, use #setLanguage instead.
//
// If you you need something more fine grained (e.g. "ES2017 without modules") use
// #setOutputFeatureSet.
//
// Java: CompilerOptions.setLanguageOut(LanguageMode)
func (o *CompilerOptions) SetLanguageOut(languageOut CompilerOptions_LanguageMode) {
	if languageOut == CompilerOptions_LanguageMode_NO_TRANSPILE {
		o.languageOutIsDefaultStrict = nil
		o.outputFeatureSet = nil
	} else {
		if languageOut == CompilerOptions_LanguageMode_STABLE {
			languageOut = CompilerOptions_LanguageMode_STABLE_OUT
		}
		isDefaultStrict := languageOut.IsDefaultStrict()
		o.languageOutIsDefaultStrict = &isDefaultStrict
		o.SetOutputFeatureSet(languageOut.ToFeatureSet())
	}
}

// Sets the features that allowed to appear in the output. Any feature in the input that is not in
// this output must be transpiled away.
//
// Note: this is a package private API since not every FeatureSet value can be properly output
// by the compiler without crashing. Both the `setBrowserFeaturesetYear` and `setLanguageOut` APIs
// are supported alternatives.
//
// Java: CompilerOptions.setOutputFeatureSet(FeatureSet)
func (o *CompilerOptions) SetOutputFeatureSet(featureSet *parser.FeatureSet) {
	o.outputFeatureSet = jutil.CheckNotNull(featureSet) // Optional.of
}

// Java: CompilerOptions.legacySetOutputFeatureSet(FeatureSet)
func (o *CompilerOptions) LegacySetOutputFeatureSet(featureSet *parser.FeatureSet) {
	o.SetOutputFeatureSet(featureSet)
}

// Gets the set of features that can appear in the output.
//
// Java: CompilerOptions.getOutputFeatureSet()
func (o *CompilerOptions) GetOutputFeatureSet() *parser.FeatureSet {
	if o.outputFeatureSet != nil {
		return o.outputFeatureSet
	}

	// Backwards compatibility for those that predate language out.
	return o.languageIn.ToFeatureSet()
}

// SetTranspilePublicClassFields disables emitting public class fields
// natively, forcing them to be transpiled into constructor assignments.
//
// Java: CompilerOptions.setTranspilePublicClassFields(boolean)
func (o *CompilerOptions) SetTranspilePublicClassFields(transpile bool) {
	if transpile {
		o.SetOutputFeatureSet(o.GetOutputFeatureSet().WithoutFeatures(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS))
	} else {
		o.SetOutputFeatureSet(o.GetOutputFeatureSet().With(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS))
	}
}

// Java: CompilerOptions.getTranspilePublicClassFields()
func (o *CompilerOptions) GetTranspilePublicClassFields() bool {
	return !o.GetOutputFeatureSet().Has(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
}

// SetExperimentalOutputFeatureSet: a nil experimentalOutputFeatureSet is
// Java null (Optional.fromNullable).
//
// Java: CompilerOptions.setExperimentalOutputFeatureSet(ExperimentalOutputFeatureSet)
func (o *CompilerOptions) SetExperimentalOutputFeatureSet(experimentalOutputFeatureSet *CompilerOptions_ExperimentalOutputFeatureSet) {
	if experimentalOutputFeatureSet == nil {
		o.experimentalOutputFeatureSet = nil
	} else {
		v := *experimentalOutputFeatureSet
		o.experimentalOutputFeatureSet = &v
	}
}

// GetExperimentalOutputFeatureSet: nil is Optional.absent().
//
// Java: CompilerOptions.getExperimentalOutputFeatureSet()
func (o *CompilerOptions) GetExperimentalOutputFeatureSet() *CompilerOptions_ExperimentalOutputFeatureSet {
	return o.experimentalOutputFeatureSet
}

// Java: CompilerOptions.needsTranspilationFrom(FeatureSet)
func (o *CompilerOptions) NeedsTranspilationFrom(languageLevel *parser.FeatureSet) bool {
	// TODO(johnplaisted): This isn't really accurate. This should instead be the *parsed* language,
	// not the *input* language.
	return o.GetLanguageIn().ToFeatureSet().Contains(languageLevel) &&
		!o.GetOutputFeatureSet().Contains(languageLevel)
}

// Java: CompilerOptions.needsTranspilationOf(FeatureSet.Feature)
func (o *CompilerOptions) NeedsTranspilationOf(feature parser.FeatureSet_Feature) bool {
	return o.GetLanguageIn().ToFeatureSet().Has(feature) && !o.GetOutputFeatureSet().Has(feature)
}

// Set which set of builtin externs to use.
//
// Java: CompilerOptions.setEnvironment(Environment)
func (o *CompilerOptions) SetEnvironment(environment CompilerOptions_Environment) {
	o.environment = environment
}

// Java: CompilerOptions.getEnvironment()
func (o *CompilerOptions) GetEnvironment() CompilerOptions_Environment {
	return o.environment
}

// Sets the dependency management options.
//
// Java: CompilerOptions.setDependencyOptions(DependencyOptions)
func (o *CompilerOptions) SetDependencyOptions(dependencyOptions *DependencyOptions) {
	o.dependencyOptions = dependencyOptions
}

// Java: CompilerOptions.getDependencyOptions()
func (o *CompilerOptions) GetDependencyOptions() *DependencyOptions {
	return o.dependencyOptions
}

// Java: CompilerOptions.setMergedPrecompiledLibraries(boolean)
func (o *CompilerOptions) SetMergedPrecompiledLibraries(mergedPrecompiledLibraries bool) {
	o.mergedPrecompiledLibraries = mergedPrecompiledLibraries
}

// Java: CompilerOptions.getMergedPrecompiledLibraries()
func (o *CompilerOptions) GetMergedPrecompiledLibraries() bool {
	return o.mergedPrecompiledLibraries
}

// Java: CompilerOptions.setWrapGoogModulesForWhitespaceOnly(boolean)
func (o *CompilerOptions) SetWrapGoogModulesForWhitespaceOnly(enable bool) {
	o.wrapGoogModulesForWhitespaceOnly = enable
}

// Java: CompilerOptions.shouldWrapGoogModulesForWhitespaceOnly()
func (o *CompilerOptions) ShouldWrapGoogModulesForWhitespaceOnly() bool {
	return o.wrapGoogModulesForWhitespaceOnly
}

// Whether to enable the bad module rewriting before typechecking that we want to get rid of
//
// Java: CompilerOptions.setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf(boolean)
func (o *CompilerOptions) SetBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf(b bool) {
	o.rewriteModulesBeforeTypechecking = b
}

// Java: CompilerOptions.shouldRewriteModulesBeforeTypechecking()
func (o *CompilerOptions) ShouldRewriteModulesBeforeTypechecking() bool {
	return o.enableModuleRewriting &&
		(o.rewriteModulesBeforeTypechecking || o.processCommonJSModules)
}

// Experimental option to disable all Closure and ES module and goog.provide rewriting
//
// Use at your own risk - disabling module rewriting is not fully tested yet.
//
// Java: CompilerOptions.setEnableModuleRewriting(boolean)
func (o *CompilerOptions) SetEnableModuleRewriting(enable bool) {
	o.enableModuleRewriting = enable
}

// Java: CompilerOptions.shouldRewriteModulesAfterTypechecking()
func (o *CompilerOptions) ShouldRewriteModulesAfterTypechecking() bool {
	return o.enableModuleRewriting && !o.rewriteModulesBeforeTypechecking
}

// Java: CompilerOptions.shouldRewriteModules()
func (o *CompilerOptions) ShouldRewriteModules() bool {
	return o.enableModuleRewriting
}

// Whether to enable support for dynamic import expressions
//
// Java: CompilerOptions.setAllowDynamicImport(boolean)
func (o *CompilerOptions) SetAllowDynamicImport(value bool) {
	o.allowDynamicImport = value
}

// Java: CompilerOptions.shouldAllowDynamicImport()
func (o *CompilerOptions) ShouldAllowDynamicImport() bool {
	return o.allowDynamicImport
}

// Set the alias name for dynamic import expressions
//
// Java: CompilerOptions.getDynamicImportAlias()
func (o *CompilerOptions) GetDynamicImportAlias() *string {
	return o.dynamicImportAlias
}

// Set the alias name for dynamic import expressions
//
// Java: CompilerOptions.setDynamicImportAlias(String)
func (o *CompilerOptions) SetDynamicImportAlias(value *string) {
	o.dynamicImportAlias = value
}

// Java: CompilerOptions.shouldAliasDynamicImport()
func (o *CompilerOptions) ShouldAliasDynamicImport() bool {
	return o.dynamicImportAlias != nil
}

// Java: CompilerOptions.shouldRunCrossChunkCodeMotion()
func (o *CompilerOptions) ShouldRunCrossChunkCodeMotion() bool {
	return o.crossChunkCodeMotion
}

// Java: CompilerOptions.shouldRunCrossChunkMethodMotion()
func (o *CompilerOptions) ShouldRunCrossChunkMethodMotion() bool {
	return o.crossChunkMethodMotion
}

// Flattens multi-level property names (e.g. a$b = x)
//
// @deprecated use getPropertyCollapseLevel
//
// Java: CompilerOptions.shouldCollapseProperties()
func (o *CompilerOptions) ShouldCollapseProperties() bool {
	return o.collapsePropertiesLevel != CompilerOptions_PropertyCollapseLevel_NONE
}

// Java: CompilerOptions.getPropertyCollapseLevel()
func (o *CompilerOptions) GetPropertyCollapseLevel() CompilerOptions_PropertyCollapseLevel {
	return o.collapsePropertiesLevel
}

// Java: CompilerOptions.setCollapseObjectLiterals(boolean)
func (o *CompilerOptions) SetCollapseObjectLiterals(enabled bool) {
	o.collapseObjectLiterals = enabled
}

// Java: CompilerOptions.getCollapseObjectLiterals()
func (o *CompilerOptions) GetCollapseObjectLiterals() bool {
	return o.collapseObjectLiterals
}

// @return Whether to attempt to remove unused class properties
//
// Java: CompilerOptions.isRemoveUnusedClassProperties()
func (o *CompilerOptions) IsRemoveUnusedClassProperties() bool {
	return o.removeUnusedClassProperties
}

// @param removeUnusedClassProperties Whether to attempt to remove unused class properties
//
// Java: CompilerOptions.setRemoveUnusedClassProperties(boolean)
func (o *CompilerOptions) SetRemoveUnusedClassProperties(removeUnusedClassProperties bool) {
	o.removeUnusedClassProperties = removeUnusedClassProperties
}

// Set the function inlining policy for the compiler.
//
// Java: CompilerOptions.setInlineFunctions(Reach)
func (o *CompilerOptions) SetInlineFunctions(reach CompilerOptions_Reach) {
	o.inlineFunctionsLevel = reach
}

// Get the function inlining policy for the compiler.
//
// Java: CompilerOptions.getInlineFunctionsLevel()
func (o *CompilerOptions) GetInlineFunctionsLevel() CompilerOptions_Reach {
	return o.inlineFunctionsLevel
}

// Java: CompilerOptions.setMaxFunctionSizeAfterInlining(int)
func (o *CompilerOptions) SetMaxFunctionSizeAfterInlining(funAstSize int) {
	jutil.CheckArgument(funAstSize > 0)
	o.maxFunctionSizeAfterInlining = funAstSize
}

// Java: CompilerOptions.getMaxFunctionSizeAfterInlining()
func (o *CompilerOptions) GetMaxFunctionSizeAfterInlining() int {
	return o.maxFunctionSizeAfterInlining
}

// Java: CompilerOptions.setInlineVariables(boolean)
func (o *CompilerOptions) SetInlineVariables(inlineVariables bool) {
	o.inlineVariables = inlineVariables
}

// Java: CompilerOptions.shouldInlineVariables()
func (o *CompilerOptions) ShouldInlineVariables() bool {
	return o.inlineVariables
}

// Set the variable inlining policy for the compiler.
//
// Java: CompilerOptions.setInlineVariables(Reach)
func (o *CompilerOptions) SetInlineVariablesOfReach(reach CompilerOptions_Reach) {
	switch reach {
	case CompilerOptions_Reach_ALL:
		o.inlineVariables = true
		o.inlineLocalVariables = true
	case CompilerOptions_Reach_LOCAL_ONLY:
		o.inlineVariables = false
		o.inlineLocalVariables = true
	case CompilerOptions_Reach_NONE:
		o.inlineVariables = false
		o.inlineLocalVariables = false
	}
}

// Set the function inlining policy for the compiler.
//
// Java: CompilerOptions.setInlineProperties(boolean)
func (o *CompilerOptions) SetInlineProperties(enable bool) {
	o.inlineProperties = enable
}

// Java: CompilerOptions.shouldInlineProperties()
func (o *CompilerOptions) ShouldInlineProperties() bool {
	return o.inlineProperties
}

// Set the variable removal policy for the compiler.
//
// Java: CompilerOptions.setRemoveUnusedVariables(Reach)
func (o *CompilerOptions) SetRemoveUnusedVariables(reach CompilerOptions_Reach) {
	switch reach {
	case CompilerOptions_Reach_ALL:
		o.removeUnusedVars = true
		o.removeUnusedLocalVars = true
	case CompilerOptions_Reach_LOCAL_ONLY:
		o.removeUnusedVars = false
		o.removeUnusedLocalVars = true
	case CompilerOptions_Reach_NONE:
		o.removeUnusedVars = false
		o.removeUnusedLocalVars = false
	}
}

// Java: CompilerOptions.shouldRemoveUnusedVariables()
func (o *CompilerOptions) ShouldRemoveUnusedVariables() bool {
	return o.removeUnusedVars
}

// Java: CompilerOptions.shouldRemoveUnusedLocalVariables()
func (o *CompilerOptions) ShouldRemoveUnusedLocalVariables() bool {
	return o.removeUnusedLocalVars
}

// Java: CompilerOptions.setRemoveAbstractMethods(boolean)
func (o *CompilerOptions) SetRemoveAbstractMethods(remove bool) {
	o.removeAbstractMethods = remove
}

// Java: CompilerOptions.shouldRemoveAbstractMethods()
func (o *CompilerOptions) ShouldRemoveAbstractMethods() bool {
	return o.removeAbstractMethods
}

// Java: CompilerOptions.setRemoveClosureAsserts(boolean)
func (o *CompilerOptions) SetRemoveClosureAsserts(remove bool) {
	o.removeClosureAsserts = remove
}

// Java: CompilerOptions.shouldRemoveClosureAsserts()
func (o *CompilerOptions) ShouldRemoveClosureAsserts() bool {
	return o.removeClosureAsserts
}

// Java: CompilerOptions.setRemoveJ2clAsserts(boolean)
func (o *CompilerOptions) SetRemoveJ2clAsserts(remove bool) {
	o.removeJ2clAsserts = remove
}

// Java: CompilerOptions.shouldRemoveJ2clAsserts()
func (o *CompilerOptions) ShouldRemoveJ2clAsserts() bool {
	return o.removeJ2clAsserts
}

// Sets how goog.tweak calls are processed.
//
// Java: CompilerOptions.setTweakProcessing(TweakProcessing)
func (o *CompilerOptions) SetTweakProcessing(tweakProcessing CompilerOptions_TweakProcessing) {
	o.tweakProcessing = tweakProcessing
}

// Java: CompilerOptions.getTweakProcessing()
func (o *CompilerOptions) GetTweakProcessing() CompilerOptions_TweakProcessing {
	return o.tweakProcessing
}
