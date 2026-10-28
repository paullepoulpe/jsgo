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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	customPasses *jutil.SetMultimap[CustomPassExecutionTime, CompilerPass]
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

// The number of characters after which we insert a line break in the code
//
// Java: CompilerOptions.DEFAULT_LINE_LENGTH_THRESHOLD
const CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD = 500

// Java: CompilerOptions.POLYMER_PROPERTY_RESERVED_FIRST_CHARS
var compilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS = jutil.NewLinkedHashSet(
	jutil.StringToUTF16("ABCDEFGHIJKLMNOPQRSTUVWXYZ$")...)

// Java: CompilerOptions.POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS
var compilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS = jutil.NewLinkedHashSet[uint16]('_', '$')

// Java: CompilerOptions.ANGULAR_PROPERTY_RESERVED_FIRST_CHARS
var compilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS = jutil.NewLinkedHashSet[uint16]('$')

// Java: CompilerOptions.getAngularPropertyReservedFirstChars()
func CompilerOptions_GetAngularPropertyReservedFirstChars() *jutil.LinkedHashSet[uint16] {
	return compilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS
}

// Java: CompilerOptions.getPolymerPropertyReservedFirstChars()
func CompilerOptions_GetPolymerPropertyReservedFirstChars() *jutil.LinkedHashSet[uint16] {
	return compilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS
}

// Java: CompilerOptions.shouldRunCrossChunkCodeMotion()
func (o *CompilerOptions) ShouldRunCrossChunkCodeMotion() bool {
	return o.crossChunkCodeMotion
}

// Java: CompilerOptions.shouldRunCrossChunkMethodMotion()
func (o *CompilerOptions) ShouldRunCrossChunkMethodMotion() bool {
	return o.crossChunkMethodMotion
}

// Java: CompilerOptions.getSourceMapOutputPath()
func (o *CompilerOptions) GetSourceMapOutputPath() *string {
	return o.sourceMapOutputPath
}

// Java: CompilerOptions.shouldGatherSourceMapInfo()
func (o *CompilerOptions) ShouldGatherSourceMapInfo() bool {
	return o.shouldAlwaysGatherSourceMapInfo || !(o.sourceMapOutputPath == nil || *o.sourceMapOutputPath == "")
}

// Java: CompilerOptions.setAlwaysGatherSourceMapInfo(boolean)
func (o *CompilerOptions) SetAlwaysGatherSourceMapInfo(shouldAlwaysGatherSourceMapInfo bool) {
	o.shouldAlwaysGatherSourceMapInfo = shouldAlwaysGatherSourceMapInfo
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

// Java: CompilerOptions.setInstrumentForCoverageOnly(boolean)
func (o *CompilerOptions) SetInstrumentForCoverageOnly(instrumentForCoverageOnly bool) {
	o.instrumentForCoverageOnly = instrumentForCoverageOnly
}

// Java: CompilerOptions.getInstrumentForCoverageOnly()
func (o *CompilerOptions) GetInstrumentForCoverageOnly() bool {
	return o.instrumentForCoverageOnly
}

// Sets file to output in-progress TypedAST format to. DO NOT USE!
//
// The "TypedAST format" is currently a gzipped TypedAst proto but this is not stable.
//
// Java: CompilerOptions.setTypedAstOutputFile(Path)
func (o *CompilerOptions) SetTypedAstOutputFile(file *jutil.Path) {
	o.typedAstOutputFile = file
}

// Java: CompilerOptions.getTypedAstOutputFile()
func (o *CompilerOptions) GetTypedAstOutputFile() *jutil.Path {
	return o.typedAstOutputFile
}

// Java: CompilerOptions.setMergedPrecompiledLibraries(boolean)
func (o *CompilerOptions) SetMergedPrecompiledLibraries(mergedPrecompiledLibraries bool) {
	o.mergedPrecompiledLibraries = mergedPrecompiledLibraries
}

// Java: CompilerOptions.getMergedPrecompiledLibraries()
func (o *CompilerOptions) GetMergedPrecompiledLibraries() bool {
	return o.mergedPrecompiledLibraries
}

// Java: CompilerOptions.setSkipTranspilationAndCrash(boolean)
func (o *CompilerOptions) SetSkipTranspilationAndCrash(value bool) {
}

// Sets the input sourcemap files, indexed by the JS files they refer to.
//
// @param inputSourceMaps the collection of input sourcemap files
//
// Java: CompilerOptions.setInputSourceMaps(ImmutableMap<String, SourceMapInput>)
func (o *CompilerOptions) SetInputSourceMaps(inputSourceMaps *jutil.LinkedHashMap[string, *SourceMapInput]) {
	o.inputSourceMaps = inputSourceMaps
}

// Java: CompilerOptions.getInputSourceMaps()
func (o *CompilerOptions) GetInputSourceMaps() *jutil.LinkedHashMap[string, *SourceMapInput] {
	return o.inputSourceMaps
}

// TODO(tbreisacher): Remove this method after ctemplate issues are solved.
//
// Java: CompilerOptions.setInferConst(boolean)
func (o *CompilerOptions) SetInferConst(value bool) {
	o.inferConsts = value
}

// Java: CompilerOptions.shouldInferConsts()
func (o *CompilerOptions) ShouldInferConsts() bool {
	return o.inferConsts
}

// Java: CompilerOptions.setIncrementalChecks(IncrementalCheckMode)
func (o *CompilerOptions) SetIncrementalChecks(value CompilerOptions_IncrementalCheckMode) {
	o.incrementalCheckMode = value
	switch value {
	case CompilerOptions_IncrementalCheckMode_OFF:
	case CompilerOptions_IncrementalCheckMode_GENERATE_IJS:
		o.SetPreserveTypeAnnotations(true)
		o.SetOutputJs(CompilerOptions_OutputJs_NORMAL)
	}
}

// Java: CompilerOptions.shouldGenerateTypedExterns()
func (o *CompilerOptions) ShouldGenerateTypedExterns() bool {
	return o.incrementalCheckMode == CompilerOptions_IncrementalCheckMode_GENERATE_IJS
}

// Java: CompilerOptions.setPrintExterns(boolean)
func (o *CompilerOptions) SetPrintExterns(printExterns bool) {
	o.printExterns = printExterns
}

// Java: CompilerOptions.shouldPrintExterns()
func (o *CompilerOptions) ShouldPrintExterns() bool {
	return o.printExterns || o.incrementalCheckMode == CompilerOptions_IncrementalCheckMode_GENERATE_IJS
}

// Deprecated. Please use setWarningLevel(DiagnosticGroups.GLOBAL_THIS, level) instead.
//
// Java: CompilerOptions.setCheckGlobalThisLevel(CheckLevel)
func (o *CompilerOptions) SetCheckGlobalThisLevel(level CheckLevel) {
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

// Java: CompilerOptions.UNLIMITED_FUN_SIZE_AFTER_INLINING
const CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING = -1

// Useful for tests to avoid having to declare two chunks
//
// Java: CompilerOptions.setRenamePrefixNamespaceAssumeCrossChunkNames(boolean)
func (o *CompilerOptions) SetRenamePrefixNamespaceAssumeCrossChunkNames(assume bool) {
	o.renamePrefixNamespaceAssumeCrossChunkNames = assume
}

// Java: CompilerOptions.assumeCrossChunkNamesForRenamePrefixNamespace()
func (o *CompilerOptions) AssumeCrossChunkNamesForRenamePrefixNamespace() bool {
	return o.renamePrefixNamespaceAssumeCrossChunkNames
}

// SetOptimizeLocalAccessForGlobalSymbolNamespace: optimize access to global
// symbols.
//
// Java: CompilerOptions.setOptimizeLocalAccessForGlobalSymbolNamespace(OptimizeLocalAccess)
func (o *CompilerOptions) SetOptimizeLocalAccessForGlobalSymbolNamespace(optimize CompilerOptions_OptimizeLocalAccess) {
	o.optimizeLocalAccessForGlobalSymbolNamespace = optimize
}

// Java: CompilerOptions.getOptimizeLocalAccessForGlobalSymbolNamespace()
func (o *CompilerOptions) GetOptimizeLocalAccessForGlobalSymbolNamespace() CompilerOptions_OptimizeLocalAccess {
	return o.optimizeLocalAccessForGlobalSymbolNamespace
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

// Java: CompilerOptions.setNameGenerator(NameGenerator)
func (o *CompilerOptions) SetNameGenerator(nameGenerator NameGenerator) {
	o.nameGenerator = nameGenerator
}

// Java: CompilerOptions.getNameGenerator()
func (o *CompilerOptions) GetNameGenerator() NameGenerator {
	return o.nameGenerator
}

// Java: CompilerOptions.setReplaceMessagesWithChromeI18n(boolean, String)
func (o *CompilerOptions) SetReplaceMessagesWithChromeI18n(replaceMessagesWithChromeI18n bool, tcProjectId *string) {
	if replaceMessagesWithChromeI18n && o.messageBundle != nil {
		if _, isEmpty := o.messageBundle.(*EmptyMessageBundle); !isEmpty {
			panic(jutil.NewRuntimeException(
				"When replacing messages with" +
					" chrome.i18n.getMessage, a message bundle should not be specified."))
		}
	}

	o.replaceMessagesWithChromeI18n = replaceMessagesWithChromeI18n
	o.tcProjectId = tcProjectId
}

// Java: CompilerOptions.getTcProjectId()
func (o *CompilerOptions) GetTcProjectId() *string {
	return o.tcProjectId
}

// Should we run the pass that does replacement of the chrome-specific `chrome.i18n.getMessage()`
// translatable message definitions?
//
// This form of l10n is incompatible with our standard `goog.getMsg()` messages.
//
// Java: CompilerOptions.shouldRunReplaceMessagesForChrome()
func (o *CompilerOptions) ShouldRunReplaceMessagesForChrome() bool {
	if o.replaceMessagesWithChromeI18n {
		_, isEmpty := o.messageBundle.(*EmptyMessageBundle)
		jutil.CheckStateMsg(
			o.messageBundle == nil || isEmpty,
			"When replacing messages with chrome.i18n.getMessage, a message bundle should not be"+
				" specified.")
		jutil.CheckStateMsg(
			!o.doLateLocalization, "Late localization is not supported for chrome.i18n.getMessage")
		return true
	} else {
		return false
	}
}

// If true, considers all missing types to be forward declared (useful for partial
// compilation).
//
// Java: CompilerOptions.setAssumeForwardDeclaredForMissingTypes(boolean)
func (o *CompilerOptions) SetAssumeForwardDeclaredForMissingTypes(assumeForwardDeclaredForMissingTypes bool) {
	o.assumeForwardDeclaredForMissingTypes = assumeForwardDeclaredForMissingTypes
}

// Java: CompilerOptions.assumeForwardDeclaredForMissingTypes()
func (o *CompilerOptions) AssumeForwardDeclaredForMissingTypes() bool {
	return o.assumeForwardDeclaredForMissingTypes
}

// Normally, when there are an equal number of single and double quotes in a string, the compiler
// will use double quotes. Set this to true to prefer single quotes.
//
// Java: CompilerOptions.setPreferSingleQuotes(boolean)
func (o *CompilerOptions) SetPreferSingleQuotes(enabled bool) {
	o.preferSingleQuotes = enabled
}

// Java: CompilerOptions.shouldPreferSingleQuotes()
func (o *CompilerOptions) ShouldPreferSingleQuotes() bool {
	return o.preferSingleQuotes
}

// Some people want to put arbitrary user input into strings, which are then run through the
// compiler. These scripts are then put into HTML. By default, we assume strings are untrusted. If
// the compiler is run from the command-line, we assume that strings are trusted.
//
// Java: CompilerOptions.setTrustedStrings(boolean)
func (o *CompilerOptions) SetTrustedStrings(yes bool) {
	o.trustedStrings = yes
}

// Java: CompilerOptions.assumeTrustedStrings()
func (o *CompilerOptions) AssumeTrustedStrings() bool {
	return o.trustedStrings
}

// Java: CompilerOptions.setPrintSourceAfterEachPass(boolean)
func (o *CompilerOptions) SetPrintSourceAfterEachPass(printSource bool) {
	o.printSourceAfterEachPass = printSource
}

// Java: CompilerOptions.shouldPrintSourceAfterEachPass()
func (o *CompilerOptions) ShouldPrintSourceAfterEachPass() bool {
	return o.printSourceAfterEachPass
}

// Java: CompilerOptions.setFilesToPrintAfterEachPassRegexList(List<String>)
func (o *CompilerOptions) SetFilesToPrintAfterEachPassRegexList(filePathRegexList []string) {
	o.filesToPrintAfterEachPassRegexList = filePathRegexList
}

// Java: CompilerOptions.setChunksToPrintAfterEachPassRegexList(List<String>)
func (o *CompilerOptions) SetChunksToPrintAfterEachPassRegexList(chunkPathRegexList []string) {
	o.chunksToPrintAfterEachPassRegexList = chunkPathRegexList
}

// Java: CompilerOptions.setQnameUsesToPrintAfterEachPassList(List<String>)
func (o *CompilerOptions) SetQnameUsesToPrintAfterEachPassList(qnameRegexList []string) {
	o.qnameUsesToPrintAfterEachPassList = qnameRegexList
}

// Java: CompilerOptions.getFilesToPrintAfterEachPassRegexList()
func (o *CompilerOptions) GetFilesToPrintAfterEachPassRegexList() []string {
	return o.filesToPrintAfterEachPassRegexList
}

// Java: CompilerOptions.getChunksToPrintAfterEachPassRegexList()
func (o *CompilerOptions) GetChunksToPrintAfterEachPassRegexList() []string {
	return o.chunksToPrintAfterEachPassRegexList
}

// Java: CompilerOptions.getQnameUsesToPrintAfterEachPassList()
func (o *CompilerOptions) GetQnameUsesToPrintAfterEachPassList() []string {
	return o.qnameUsesToPrintAfterEachPassList
}

// Java: CompilerOptions.getTracerMode()
func (o *CompilerOptions) GetTracerMode() CompilerOptions_TracerMode {
	return o.tracer
}

// NOTE: Timing information will not be printed if compiler.disableThreads() is called!
//
// Java: CompilerOptions.setTracerMode(TracerMode)
func (o *CompilerOptions) SetTracerMode(mode CompilerOptions_TracerMode) {
	o.tracer = mode
}

// Java: CompilerOptions.getTracerOutput()
func (o *CompilerOptions) GetTracerOutput() *jutil.Path {
	return o.tracerOutput
}

// Java: CompilerOptions.setTracerOutput(Path)
func (o *CompilerOptions) SetTracerOutput(out *jutil.Path) {
	o.tracerOutput = out
}

// Java: CompilerOptions.getExtraReportGenerators()
func (o *CompilerOptions) GetExtraReportGenerators() []SortingErrorManager_ErrorReportGenerator {
	return o.extraReportGenerators
}

// Java: CompilerOptions.addReportGenerator(SortingErrorManager.ErrorReportGenerator)
func (o *CompilerOptions) AddReportGenerator(generator SortingErrorManager_ErrorReportGenerator) {
	o.extraReportGenerators = append(o.extraReportGenerators, generator)
}

// When enabled, assume that apparently side-effect free code is meaningful.
//
// Java: CompilerOptions.setProtectHiddenSideEffects(boolean)
func (o *CompilerOptions) SetProtectHiddenSideEffects(enable bool) {
	o.protectHiddenSideEffects = enable
}

// Whether or not the compiler should wrap apparently side-effect free code to prevent it from
// being removed
//
// Java: CompilerOptions.shouldProtectHiddenSideEffects()
func (o *CompilerOptions) ShouldProtectHiddenSideEffects() bool {
	// TODO(lharker): would it make sense to make 'checksOnly' a tristate instead? either we're
	// running only checks for refactoring; running checks to generate a .typedast; or running
	// checks as part of a fully optimizing build?
	return o.protectHiddenSideEffects && (!o.checksOnly || o.GetTypedAstOutputFile() != nil)
}

// Java: CompilerOptions.setAssumeGettersArePure(boolean)
func (o *CompilerOptions) SetAssumeGettersArePure(x bool) {
	o.assumeGettersArePure = x
}

// Java: CompilerOptions.getAssumeGettersArePure()
func (o *CompilerOptions) GetAssumeGettersArePure() bool {
	return o.assumeGettersArePure
}

// Java: CompilerOptions.setAssumePropertiesAreStaticallyAnalyzable(boolean)
func (o *CompilerOptions) SetAssumePropertiesAreStaticallyAnalyzable(x bool) {
	o.assumePropertiesAreStaticallyAnalyzable = x
}

// Java: CompilerOptions.getAssumePropertiesAreStaticallyAnalyzable()
func (o *CompilerOptions) GetAssumePropertiesAreStaticallyAnalyzable() bool {
	return o.assumePropertiesAreStaticallyAnalyzable
}

// Java: CompilerOptions.setAssumeStaticInheritanceIsNotUsed(boolean)
func (o *CompilerOptions) SetAssumeStaticInheritanceIsNotUsed(x bool) {
	o.assumeStaticInheritanceIsNotUsed = x
}

// Java: CompilerOptions.getAssumeStaticInheritanceIsNotUsed()
func (o *CompilerOptions) GetAssumeStaticInheritanceIsNotUsed() bool {
	return o.assumeStaticInheritanceIsNotUsed
}

// Java: CompilerOptions.getConformanceReportingMode()
func (o *CompilerOptions) GetConformanceReportingMode() CompilerOptions_ConformanceReportingMode {
	return o.conformanceReportingMode
}

// Java: CompilerOptions.setConformanceReportingMode(ConformanceReportingMode)
func (o *CompilerOptions) SetConformanceReportingMode(mode CompilerOptions_ConformanceReportingMode) {
	o.conformanceReportingMode = mode
}

// Java: CompilerOptions.setConformanceRemoveRegexFromPath(Optional<Pattern>)
func (o *CompilerOptions) SetConformanceRemoveRegexFromPath(pattern *jutil.Pattern) {
	o.conformanceRemoveRegexFromPath = pattern
}

// Java: CompilerOptions.getConformanceRemoveRegexFromPath()
func (o *CompilerOptions) GetConformanceRemoveRegexFromPath() *jutil.Pattern {
	return o.conformanceRemoveRegexFromPath
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
	o.conformanceRemoveRegexFromPath =
		jutil.Pattern_Compile("^((.*/)?google3/)?(/?(blaze|bazel)-out/[^/]+/(bin|genfiles)/)?")
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

type defineReplacementKind uint8

const (
	defineReplacementBool defineReplacementKind = iota
	defineReplacementInt
	defineReplacementDouble
	defineReplacementString
)

type defineReplacement struct {
	kind    defineReplacementKind
	boolVal bool
	intVal  int32
	numVal  float64
	strVal  string
}

func (d defineReplacement) toNode() rhino.Node {
	switch d.kind {
	case defineReplacementBool:
		return NodeUtil_BooleanNode(d.boolVal)
	case defineReplacementInt:
		return NodeUtil_NumberNode(float64(d.intVal), rhino.NilNode)
	case defineReplacementDouble:
		return NodeUtil_NumberNode(d.numVal, rhino.NilNode)
	case defineReplacementString:
		return rhino.IR_String(d.strVal)
	default:
		panic(jutil.NewIllegalStateException())
	}
}

// Returns the map of define replacements.
//
// Java: CompilerOptions.getDefineReplacements()
func (o *CompilerOptions) GetDefineReplacements() *jutil.LinkedHashMap[string, rhino.Node] {
	m := jutil.NewLinkedHashMap[string, rhino.Node]()
	for name, value := range o.defineReplacements.All() {
		m.Put(name, value.toNode())
	}
	return m
}

// Sets the value of the @define variable in JS to a boolean literal.
//
// Java: CompilerOptions.setDefineToBooleanLiteral(String, boolean)
func (o *CompilerOptions) SetDefineToBooleanLiteral(defineName string, value bool) {
	o.defineReplacements.Put(defineName, defineReplacement{kind: defineReplacementBool, boolVal: value})
}

// Sets the value of the @define variable in JS to a String literal.
//
// Java: CompilerOptions.setDefineToStringLiteral(String, String)
func (o *CompilerOptions) SetDefineToStringLiteral(defineName string, value string) {
	o.defineReplacements.Put(defineName, defineReplacement{kind: defineReplacementString, strVal: value})
}

// Sets the value of the @define variable in JS to a number literal.
//
// Java: CompilerOptions.setDefineToNumberLiteral(String, int)
func (o *CompilerOptions) SetDefineToNumberLiteral(defineName string, value int) {
	o.defineReplacements.Put(defineName, defineReplacement{kind: defineReplacementInt, intVal: int32(value)})
}

// Sets the value of the @define variable in JS to a number literal.
//
// Java: CompilerOptions.setDefineToDoubleLiteral(String, double)
func (o *CompilerOptions) SetDefineToDoubleLiteral(defineName string, value float64) {
	o.defineReplacements.Put(defineName, defineReplacement{kind: defineReplacementDouble, numVal: value})
}

// Skip all possible passes, to make the compiler as fast as possible.
//
// Java: CompilerOptions.skipAllCompilerPasses()
func (o *CompilerOptions) SkipAllCompilerPasses() {
	o.SetSkipNonTranspilationPasses(true)
}

// Whether the warnings guard in this Options object enables the given group of warnings.
//
// Java: CompilerOptions.enables(DiagnosticGroup)
func (o *CompilerOptions) Enables(group *DiagnosticGroup) bool {
	return o.warningsGuard.MustRunChecks(group) == base.Tri_TRUE
}

// Whether the warnings guard in this Options object disables the given group of warnings.
//
// Java: CompilerOptions.disables(DiagnosticGroup)
func (o *CompilerOptions) Disables(group *DiagnosticGroup) bool {
	return o.warningsGuard.MustRunChecks(group) == base.Tri_FALSE
}

// Java: CompilerOptions.setUnknownDefinesToIgnore(ImmutableList<String>)
func (o *CompilerOptions) SetUnknownDefinesToIgnore(unknownDefinesToIgnore []string) {
	o.unknownDefinesToIgnore = unknownDefinesToIgnore
}

// Java: CompilerOptions.addUnknownDefinesToIgnore(ImmutableList<String>)
func (o *CompilerOptions) AddUnknownDefinesToIgnore(unknownDefinesToIgnore []string) {
	combined := make([]string, 0, len(o.unknownDefinesToIgnore)+len(unknownDefinesToIgnore))
	combined = append(combined, o.unknownDefinesToIgnore...)
	combined = append(combined, unknownDefinesToIgnore...)
	o.unknownDefinesToIgnore = combined
}

// Java: CompilerOptions.getUnknownDefinesToIgnore()
func (o *CompilerOptions) GetUnknownDefinesToIgnore() []string {
	return o.unknownDefinesToIgnore
}

// Configure the given type of warning to the given level.
//
// Java: CompilerOptions.setWarningLevel(DiagnosticGroup, CheckLevel)
func (o *CompilerOptions) SetWarningLevel(type_ *DiagnosticGroup, level CheckLevel) {
	o.AddWarningsGuard(NewDiagnosticGroupWarningsGuard(type_, level))
}

// Java: CompilerOptions.getWarningsGuard()
func (o *CompilerOptions) GetWarningsGuard() WarningsGuard {
	return o.warningsGuard
}

// Reset the warnings guard.
//
// Java: CompilerOptions.resetWarningsGuard()
func (o *CompilerOptions) ResetWarningsGuard() {
	o.warningsGuard = NewComposeWarningsGuard()
}

// Add a guard to the set of warnings guards.
//
// Java: CompilerOptions.addWarningsGuard(WarningsGuard)
func (o *CompilerOptions) AddWarningsGuard(guard WarningsGuard) {
	o.warningsGuard.AddGuard(guard)
}

// Sets the variable and property renaming policies for the compiler, in a way that clears
// warnings about the renaming policy being uninitialized from flags.
//
// Java: CompilerOptions.setRenamingPolicy(VariableRenamingPolicy, PropertyRenamingPolicy)
func (o *CompilerOptions) SetRenamingPolicy(newVariablePolicy VariableRenamingPolicy, newPropertyPolicy PropertyRenamingPolicy) {
	o.variableRenaming = newVariablePolicy
	o.propertyRenaming = newPropertyPolicy
}

// @param replaceIdGenerators the replaceIdGenerators to set
//
// Java: CompilerOptions.setReplaceIdGenerators(boolean)
func (o *CompilerOptions) SetReplaceIdGenerators(replaceIdGenerators bool) {
	o.replaceIdGenerators = replaceIdGenerators
}

// Java: CompilerOptions.shouldReplaceIdGenerators()
func (o *CompilerOptions) ShouldReplaceIdGenerators() bool {
	return o.replaceIdGenerators
}

// Sets the id generators to replace.
//
// Java: CompilerOptions.setIdGenerators(Set<String>)
func (o *CompilerOptions) SetIdGenerators(idGenerators *jutil.LinkedHashSet[string]) {
	builder := jutil.NewLinkedHashMap[string, RenamingMap]()
	for name := range idGenerators.All() {
		builder.Put(name, RenamingMap(RenamingToken_INCONSISTENT))
	}
	o.idGenerators = builder
}

// Sets the id generators to replace.
//
// Java: CompilerOptions.setIdGenerators(Map<String, RenamingMap>)
func (o *CompilerOptions) SetIdGeneratorsOfMap(idGenerators *jutil.LinkedHashMap[string, RenamingMap]) {
	o.idGenerators = jutil.NewLinkedHashMapFrom(idGenerators)
}

// Java: CompilerOptions.getIdGenerators()
func (o *CompilerOptions) GetIdGenerators() *jutil.LinkedHashMap[string, RenamingMap] {
	return o.idGenerators
}

// A previous map of ids (serialized to a string by a previous compile). This will be used as a
// hint during the ReplaceIdGenerators pass, which will attempt to reuse the same ids.
//
// Java: CompilerOptions.setIdGeneratorsMap(String)
func (o *CompilerOptions) SetIdGeneratorsMap(previousMappings *string) {
	o.idGeneratorsMapSerialized = previousMappings
}

// Java: CompilerOptions.getIdGeneratorsMapSerialized()
func (o *CompilerOptions) GetIdGeneratorsMapSerialized() *string {
	return o.idGeneratorsMapSerialized
}

// Sets the hash function to use for Xid
//
// Java: CompilerOptions.setXidHashFunction(Xid.HashFunction)
func (o *CompilerOptions) SetXidHashFunction(xidHashFunction Xid_HashFunction) {
	o.xidHashFunction = xidHashFunction
}

// Java: CompilerOptions.getXidHashFunction()
func (o *CompilerOptions) GetXidHashFunction() Xid_HashFunction {
	return o.xidHashFunction
}

// Sets the hash function to use for chunk ID generation
//
// Java: CompilerOptions.setChunkIdHashFunction(Xid.HashFunction)
func (o *CompilerOptions) SetChunkIdHashFunction(chunkIdHashFunction Xid_HashFunction) {
	o.chunkIdHashFunction = chunkIdHashFunction
}

// Java: CompilerOptions.getChunkIdHashFunction()
func (o *CompilerOptions) GetChunkIdHashFunction() Xid_HashFunction {
	return o.chunkIdHashFunction
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

// Sets the functions whose debug strings to replace.
//
// Java: CompilerOptions.setReplaceStringsConfiguration(String, List<String>)
func (o *CompilerOptions) SetReplaceStringsConfiguration(placeholderToken string, functionDescriptors []string) {
	o.replaceStringsPlaceholderToken = placeholderToken
	o.replaceStringsFunctionDescriptions = append([]string{}, functionDescriptors...)
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

// Java: CompilerOptions.setColorizeErrorOutput(boolean)
func (o *CompilerOptions) SetColorizeErrorOutput(colorizeErrorOutput bool) {
	o.colorizeErrorOutput = colorizeErrorOutput
}

// Java: CompilerOptions.shouldColorizeErrorOutput()
func (o *CompilerOptions) ShouldColorizeErrorOutput() bool {
	return o.colorizeErrorOutput
}

// Java: CompilerOptions.setChecksOnly(boolean)
func (o *CompilerOptions) SetChecksOnly(checksOnly bool) {
	o.checksOnly = checksOnly
}

// Java: CompilerOptions.isChecksOnly()
func (o *CompilerOptions) IsChecksOnly() bool {
	return o.checksOnly
}

// Java: CompilerOptions.setOutputJs(OutputJs)
func (o *CompilerOptions) SetOutputJs(outputJs CompilerOptions_OutputJs) {
	o.outputJs = outputJs
}

// Java: CompilerOptions.getOutputJs()
func (o *CompilerOptions) GetOutputJs() CompilerOptions_OutputJs {
	return o.outputJs
}

// Java: CompilerOptions.setGenerateExports(boolean)
func (o *CompilerOptions) SetGenerateExports(generateExports bool) {
	o.generateExports = generateExports
}

// Java: CompilerOptions.shouldGenerateExports()
func (o *CompilerOptions) ShouldGenerateExports() bool {
	return o.generateExports
}

// Java: CompilerOptions.setExportLocalPropertyDefinitions(boolean)
func (o *CompilerOptions) SetExportLocalPropertyDefinitions(export bool) {
	o.exportLocalPropertyDefinitions = export
}

// Java: CompilerOptions.shouldExportLocalPropertyDefinitions()
func (o *CompilerOptions) ShouldExportLocalPropertyDefinitions() bool {
	return o.exportLocalPropertyDefinitions
}

// Java: CompilerOptions.setAngularPass(boolean)
func (o *CompilerOptions) SetAngularPass(angularPass bool) {
	o.angularPass = angularPass
}

// Java: CompilerOptions.getAngularPass()
func (o *CompilerOptions) GetAngularPass() bool {
	return o.angularPass
}

// Java: CompilerOptions.setPolymerVersion(Integer)
func (o *CompilerOptions) SetPolymerVersion(polymerVersion *int32) {
	jutil.CheckArgumentMsg1(
		polymerVersion == nil || *polymerVersion == 1 || *polymerVersion == 2,
		"Invalid Polymer version: (%s)",
		polymerVersion)
	// Internally we model the polymerVersion as a boolean. if True, run the PolymerPass; otherwise
	// don't. The public API is an integer for historical reasons & in case we need to modify it
	// in the future.
	o.polymerPass = polymerVersion != nil
}

// Java: CompilerOptions.isPolymerPassEnabled()
func (o *CompilerOptions) IsPolymerPassEnabled() bool {
	return o.polymerPass
}

// Java: CompilerOptions.setChromePass(boolean)
func (o *CompilerOptions) SetChromePass(chromePass bool) {
	o.chromePass = chromePass
}

// Java: CompilerOptions.isChromePassEnabled()
func (o *CompilerOptions) IsChromePassEnabled() bool {
	return o.chromePass
}

// Java: CompilerOptions.setJ2clPass(J2clPassMode)
func (o *CompilerOptions) SetJ2clPass(j2clPassMode CompilerOptions_J2clPassMode) {
	o.j2clPassMode = j2clPassMode
}

// Java: CompilerOptions.getJ2clPass()
func (o *CompilerOptions) GetJ2clPass() CompilerOptions_J2clPassMode {
	return o.j2clPassMode
}

// Java: CompilerOptions.setJ2clMinifierEnabled(boolean)
func (o *CompilerOptions) SetJ2clMinifierEnabled(enabled bool) {
	o.j2clMinifierEnabled = enabled
}

// Java: CompilerOptions.isJ2clMinifierEnabled()
func (o *CompilerOptions) IsJ2clMinifierEnabled() bool {
	return o.j2clMinifierEnabled
}

// Java: CompilerOptions.setJ2clMinifierPruningManifest(String)
func (o *CompilerOptions) SetJ2clMinifierPruningManifest(j2clMinifierPruningManifest *string) {
	o.j2clMinifierPruningManifest = j2clMinifierPruningManifest
}

// Java: CompilerOptions.getJ2clMinifierPruningManifest()
func (o *CompilerOptions) GetJ2clMinifierPruningManifest() *string {
	return o.j2clMinifierPruningManifest
}

// Java: CompilerOptions.setCodingConvention(CodingConvention)
func (o *CompilerOptions) SetCodingConvention(codingConvention CodingConvention) {
	o.codingConvention = codingConvention
}

// Java: CompilerOptions.getCodingConvention()
func (o *CompilerOptions) GetCodingConvention() CodingConvention {
	return o.codingConvention
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

// Controls how detailed the compilation summary is. Values: 0 (never print summary), 1 (print
// summary only if there are errors or warnings), 2 (print summary if type checking is on), 3
// (always print summary). The default level is 1
//
// Java: CompilerOptions.setSummaryDetailLevel(int)
func (o *CompilerOptions) SetSummaryDetailLevel(summaryDetailLevel int) {
	o.summaryDetailLevel = summaryDetailLevel
}

// Java: CompilerOptions.getSummaryDetailLevel()
func (o *CompilerOptions) GetSummaryDetailLevel() int {
	return o.summaryDetailLevel
}

// Java: CompilerOptions.setExtraAnnotationNames(Iterable<String>)
func (o *CompilerOptions) SetExtraAnnotationNames(extraAnnotationNames []string) {
	o.extraAnnotationNames = jutil.NewLinkedHashSet(extraAnnotationNames...)
}

// Java: CompilerOptions.getExtraAnnotationNames()
func (o *CompilerOptions) GetExtraAnnotationNames() *jutil.LinkedHashSet[string] {
	return o.extraAnnotationNames
}

// Sets the output charset.
//
// Java: CompilerOptions.setOutputCharset(Charset)
func (o *CompilerOptions) SetOutputCharset(charset *jutil.Charset) {
	o.outputCharset = charset
}

// Gets the output charset.
//
// Java: CompilerOptions.getOutputCharset()
func (o *CompilerOptions) GetOutputCharset() *jutil.Charset {
	return o.outputCharset
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

// Set a custom handler for warnings and errors.
//
// This is mostly used for piping the warnings and errors to a file behind the scenes.
//
// If you want to filter warnings and errors, you should use a WarningsGuard.
//
// If you want to change how warnings and errors are reported to the user, you should set a
// ErrorManager on the Compiler. An ErrorManager is intended to summarize the errors for a single
// compile job.
//
// Java: CompilerOptions.setErrorHandler(ErrorHandler)
func (o *CompilerOptions) SetErrorHandler(handler ErrorHandler) {
	o.errorHandler = handler
}

// Java: CompilerOptions.getErrorHandler()
func (o *CompilerOptions) GetErrorHandler() ErrorHandler {
	return o.errorHandler
}

// If true, enables type inference. If checkTypes is enabled, this flag has no effect.
//
// Java: CompilerOptions.setInferTypes(boolean)
func (o *CompilerOptions) SetInferTypes(enable bool) {
	o.inferTypes = enable
}

// Gets the inferTypes flag. Note that if checkTypes is enabled, this flag is ignored when
// configuring the compiler.
//
// Java: CompilerOptions.getInferTypes()
func (o *CompilerOptions) GetInferTypes() bool {
	return o.inferTypes
}

// @deprecated This is a no-op.
//
// Java: CompilerOptions.setNewTypeInference(boolean)
func (o *CompilerOptions) SetNewTypeInference(enable bool) {
}

// @return true if either typechecker is ON.
//
// Java: CompilerOptions.isTypecheckingEnabled()
func (o *CompilerOptions) IsTypecheckingEnabled() bool {
	return o.checkTypes
}

// @return Whether assumeStrictThis is set.
//
// Java: CompilerOptions.assumeStrictThis()
func (o *CompilerOptions) AssumeStrictThis() bool {
	return o.assumeStrictThis
}

// If true, enables enables additional optimizations.
//
// Java: CompilerOptions.setAssumeStrictThis(boolean)
func (o *CompilerOptions) SetAssumeStrictThis(enable bool) {
	o.assumeStrictThis = enable
}

// Sets the list of properties that we report property invalidation errors for.
//
// Java: CompilerOptions.setPropertiesThatMustDisambiguate(Set<String>)
func (o *CompilerOptions) SetPropertiesThatMustDisambiguate(names *jutil.LinkedHashSet[string]) {
	o.propertiesThatMustDisambiguate = jutil.NewLinkedHashSet(names.ToSlice()...)
}

// Java: CompilerOptions.getPropertiesThatMustDisambiguate()
func (o *CompilerOptions) GetPropertiesThatMustDisambiguate() *jutil.LinkedHashSet[string] {
	return o.propertiesThatMustDisambiguate
}
