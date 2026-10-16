// Port of com/google/javascript/jscomp/CompilerOptions.java @ closure-compiler v20260415

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
	"fmt"
	"reflect"
	"slices"
	"strings"

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
	outputFeatureSet            *parser.FeatureSet
	experimentalForceTranspiles []CompilerOptions_ExperimentalForceTranspile
	languageOutIsDefaultStrict  *bool
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
	// More aggressive function inlining
	assumeClosuresOnlyCaptureReferences bool
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
	renamePrefixNamespaceAssumeCrossChunkNames bool
	collapsePropertiesLevel                    CompilerOptions_PropertyCollapseLevel
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
	defineReplacements *jutil.LinkedHashMap[string, any]
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
	extraReportGenerators *jutil.ArrayList[SortingErrorManager_ErrorReportGenerator]
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

// CompilerOptions_Reach is a common enum for compiler passes that can run either globally or locally.
type CompilerOptions_Reach int32

const (
	CompilerOptions_Reach_ALL CompilerOptions_Reach = iota
	CompilerOptions_Reach_LOCAL_ONLY
	CompilerOptions_Reach_NONE
)

var compilerOptions_Reach_names = [...]string{"ALL", "LOCAL_ONLY", "NONE"}

// CompilerOptions_Reach_values is Reach.values().
func CompilerOptions_Reach_values() []CompilerOptions_Reach {
	return []CompilerOptions_Reach{CompilerOptions_Reach_ALL, CompilerOptions_Reach_LOCAL_ONLY, CompilerOptions_Reach_NONE}
}

// String is Enum.name().
func (e CompilerOptions_Reach) String() string { return compilerOptions_Reach_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_Reach) Ordinal() int { return int(e) }

// Java: CompilerOptions.Reach.isOn()
func (e CompilerOptions_Reach) IsOn() bool {
	return e != CompilerOptions_Reach_NONE
}

// Java: CompilerOptions.Reach.includesGlobals()
func (e CompilerOptions_Reach) IncludesGlobals() bool {
	return e == CompilerOptions_Reach_ALL
}

// CompilerOptions_PropertyCollapseLevel is CompilerOptions.PropertyCollapseLevel.
type CompilerOptions_PropertyCollapseLevel int32

const (
	CompilerOptions_PropertyCollapseLevel_ALL CompilerOptions_PropertyCollapseLevel = iota
	CompilerOptions_PropertyCollapseLevel_NONE
	CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT
)

var compilerOptions_PropertyCollapseLevel_names = [...]string{"ALL", "NONE", "MODULE_EXPORT"}

// CompilerOptions_PropertyCollapseLevel_values is PropertyCollapseLevel.values().
func CompilerOptions_PropertyCollapseLevel_values() []CompilerOptions_PropertyCollapseLevel {
	return []CompilerOptions_PropertyCollapseLevel{CompilerOptions_PropertyCollapseLevel_ALL, CompilerOptions_PropertyCollapseLevel_NONE, CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT}
}

// String is Enum.name().
func (e CompilerOptions_PropertyCollapseLevel) String() string {
	return compilerOptions_PropertyCollapseLevel_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_PropertyCollapseLevel) Ordinal() int { return int(e) }

// CompilerOptions_BrowserFeaturesetYear represents browser feature set year to use for compilation. It tells the
// JSCompiler to output code that works on the releases of major browsers that
// were current as of January 1 of the given year, without including
// transpilation or other workarounds for browsers older than that
type CompilerOptions_BrowserFeaturesetYear int32

const (
	CompilerOptions_BrowserFeaturesetYear_YEAR_2012 CompilerOptions_BrowserFeaturesetYear = iota
	CompilerOptions_BrowserFeaturesetYear_YEAR_2018
	CompilerOptions_BrowserFeaturesetYear_YEAR_2019
	CompilerOptions_BrowserFeaturesetYear_YEAR_2020
	CompilerOptions_BrowserFeaturesetYear_YEAR_2021
	CompilerOptions_BrowserFeaturesetYear_YEAR_2022
	CompilerOptions_BrowserFeaturesetYear_YEAR_2023
	CompilerOptions_BrowserFeaturesetYear_YEAR_2024
	CompilerOptions_BrowserFeaturesetYear_YEAR_2025
	CompilerOptions_BrowserFeaturesetYear_YEAR_2026
)

var compilerOptions_BrowserFeaturesetYear_names = [...]string{"YEAR_2012", "YEAR_2018", "YEAR_2019", "YEAR_2020", "YEAR_2021", "YEAR_2022", "YEAR_2023", "YEAR_2024", "YEAR_2025", "YEAR_2026"}

// CompilerOptions_BrowserFeaturesetYear_values is BrowserFeaturesetYear.values().
func CompilerOptions_BrowserFeaturesetYear_values() []CompilerOptions_BrowserFeaturesetYear {
	return []CompilerOptions_BrowserFeaturesetYear{CompilerOptions_BrowserFeaturesetYear_YEAR_2012, CompilerOptions_BrowserFeaturesetYear_YEAR_2018, CompilerOptions_BrowserFeaturesetYear_YEAR_2019, CompilerOptions_BrowserFeaturesetYear_YEAR_2020, CompilerOptions_BrowserFeaturesetYear_YEAR_2021, CompilerOptions_BrowserFeaturesetYear_YEAR_2022, CompilerOptions_BrowserFeaturesetYear_YEAR_2023, CompilerOptions_BrowserFeaturesetYear_YEAR_2024, CompilerOptions_BrowserFeaturesetYear_YEAR_2025, CompilerOptions_BrowserFeaturesetYear_YEAR_2026}
}

// String is Enum.name().
func (e CompilerOptions_BrowserFeaturesetYear) String() string {
	return compilerOptions_BrowserFeaturesetYear_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_BrowserFeaturesetYear) Ordinal() int { return int(e) }

var compilerOptions_BrowserFeaturesetYear_ptrs = func() []CompilerOptions_BrowserFeaturesetYear {
	v := CompilerOptions_BrowserFeaturesetYear_values()
	return v
}()

// Ptr returns the canonical pointer to e (Go-only: language-adapter — a
// nullable Java BrowserFeaturesetYear reference is a *CompilerOptions_BrowserFeaturesetYear; nil is null).
func (e CompilerOptions_BrowserFeaturesetYear) Ptr() *CompilerOptions_BrowserFeaturesetYear {
	return &compilerOptions_BrowserFeaturesetYear_ptrs[e]
}

// Enum fields, indexed by ordinal (Java: BrowserFeaturesetYear(int, FeatureSet)).
var compilerOptions_BrowserFeaturesetYear_year = [...]int{2012, 2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025, 2026}

var compilerOptions_BrowserFeaturesetYear_featureSet = [...]*parser.FeatureSet{
	CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT.ToFeatureSet(),
	CompilerOptions_LanguageMode_ECMASCRIPT_2016.ToFeatureSet(),
	parser.FeatureSet_BROWSER_2019,
	parser.FeatureSet_BROWSER_2020,
	parser.FeatureSet_BROWSER_2021,
	parser.FeatureSet_BROWSER_2022,
	parser.FeatureSet_BROWSER_2023,
	parser.FeatureSet_BROWSER_2024,
	parser.FeatureSet_BROWSER_2025,
	parser.FeatureSet_BROWSER_2026,
}

// Java: CompilerOptions.BrowserFeaturesetYear.YEAR_MAP
var compilerOptions_BrowserFeaturesetYear_YEAR_MAP = func() *jutil.LinkedHashMap[int, CompilerOptions_BrowserFeaturesetYear] {
	m := jutil.NewLinkedHashMap[int, CompilerOptions_BrowserFeaturesetYear]()
	// go/keep-sorted start
	m.Put(2012, CompilerOptions_BrowserFeaturesetYear_YEAR_2012)
	m.Put(2018, CompilerOptions_BrowserFeaturesetYear_YEAR_2018)
	m.Put(2019, CompilerOptions_BrowserFeaturesetYear_YEAR_2019)
	m.Put(2020, CompilerOptions_BrowserFeaturesetYear_YEAR_2020)
	m.Put(2021, CompilerOptions_BrowserFeaturesetYear_YEAR_2021)
	m.Put(2022, CompilerOptions_BrowserFeaturesetYear_YEAR_2022)
	m.Put(2023, CompilerOptions_BrowserFeaturesetYear_YEAR_2023)
	m.Put(2024, CompilerOptions_BrowserFeaturesetYear_YEAR_2024)
	m.Put(2025, CompilerOptions_BrowserFeaturesetYear_YEAR_2025)
	m.Put(2026, CompilerOptions_BrowserFeaturesetYear_YEAR_2026)
	// go/keep-sorted end
	return m
}()

// Java: CompilerOptions.BrowserFeaturesetYear.from(int)
func CompilerOptions_BrowserFeaturesetYear_From(year int) CompilerOptions_BrowserFeaturesetYear {
	jutil.CheckStateMsg1(
		compilerOptions_BrowserFeaturesetYear_YEAR_MAP.ContainsKey(year),
		"Illegal browser_featureset_year=%s. We support values 2012, or 2018..2026 only",
		year)
	return compilerOptions_BrowserFeaturesetYear_YEAR_MAP.Get(year)
}

// Java: CompilerOptions.BrowserFeaturesetYear.setDependentValuesFromYear(CompilerOptions)
func (e CompilerOptions_BrowserFeaturesetYear) SetDependentValuesFromYear(options *CompilerOptions) {
	year := compilerOptions_BrowserFeaturesetYear_year[e]
	options.SetOutputFeatureSet(compilerOptions_BrowserFeaturesetYear_featureSet[e])
	// All values targeted by browser featureset year are default strict.
	defaultStrict := true
	options.languageOutIsDefaultStrict = &defaultStrict
	options.SetDefineToNumberLiteral("goog.FEATURESET_YEAR", year)
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES5", year > 2012)
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES6", year >= 2018)
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES2020", year >= 2021)
}

// Java: CompilerOptions.BrowserFeaturesetYear.getFeatureSet()
func (e CompilerOptions_BrowserFeaturesetYear) GetFeatureSet() *parser.FeatureSet {
	return compilerOptions_BrowserFeaturesetYear_featureSet[e]
}

// Java: CompilerOptions.BrowserFeaturesetYear.getYear()
func (e CompilerOptions_BrowserFeaturesetYear) GetYear() int {
	return compilerOptions_BrowserFeaturesetYear_year[e]
}

// CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor returns the
// minimum browser featureset year required for the given feature, or nil if
// the feature is in a language level higher than that corresponding to the
// most recent BrowserFeaturesetYear.
//
// Java: CompilerOptions.BrowserFeaturesetYear.minimumRequiredFor(FeatureSet.Feature)
func CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(feature parser.FeatureSet_Feature) *CompilerOptions_BrowserFeaturesetYear {
	// Depends on YEAR_TO_FEATURE_SET being created with keys in increasing order
	for year := range compilerOptions_BrowserFeaturesetYear_YEAR_MAP.Keys() {
		yearObject := compilerOptions_BrowserFeaturesetYear_YEAR_MAP.Get(year)
		if yearObject.GetFeatureSet().ContainsFeature(feature) {
			return yearObject.Ptr()
		}
	}
	return nil
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

// CompilerOptions_IncrementalCheckMode is CompilerOptions.IncrementalCheckMode.
type CompilerOptions_IncrementalCheckMode int32

const (
	// Normal mode
	CompilerOptions_IncrementalCheckMode_OFF CompilerOptions_IncrementalCheckMode = iota
	// The compiler should generate an output file that represents the type-only
	// interface of the code being compiled. This is useful for incremental type
	// checking.
	CompilerOptions_IncrementalCheckMode_GENERATE_IJS
)

var compilerOptions_IncrementalCheckMode_names = [...]string{"OFF", "GENERATE_IJS"}

// CompilerOptions_IncrementalCheckMode_values is IncrementalCheckMode.values().
func CompilerOptions_IncrementalCheckMode_values() []CompilerOptions_IncrementalCheckMode {
	return []CompilerOptions_IncrementalCheckMode{CompilerOptions_IncrementalCheckMode_OFF, CompilerOptions_IncrementalCheckMode_GENERATE_IJS}
}

// String is Enum.name().
func (e CompilerOptions_IncrementalCheckMode) String() string {
	return compilerOptions_IncrementalCheckMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_IncrementalCheckMode) Ordinal() int { return int(e) }

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

// Java: CompilerOptions.UNLIMITED_FUN_SIZE_AFTER_INLINING
const CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING = -1

// CompilerOptions_ExtractPrototypeMemberDeclarationsMode is CompilerOptions.ExtractPrototypeMemberDeclarationsMode.
type CompilerOptions_ExtractPrototypeMemberDeclarationsMode int32

const (
	CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF CompilerOptions_ExtractPrototypeMemberDeclarationsMode = iota
	CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP
	CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP
	CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE
)

var compilerOptions_ExtractPrototypeMemberDeclarationsMode_names = [...]string{"OFF", "USE_GLOBAL_TEMP", "USE_CHUNK_TEMP", "USE_IIFE"}

// CompilerOptions_ExtractPrototypeMemberDeclarationsMode_values is ExtractPrototypeMemberDeclarationsMode.values().
func CompilerOptions_ExtractPrototypeMemberDeclarationsMode_values() []CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return []CompilerOptions_ExtractPrototypeMemberDeclarationsMode{CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF, CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP, CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP, CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE}
}

// String is Enum.name().
func (e CompilerOptions_ExtractPrototypeMemberDeclarationsMode) String() string {
	return compilerOptions_ExtractPrototypeMemberDeclarationsMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_ExtractPrototypeMemberDeclarationsMode) Ordinal() int { return int(e) }

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

// CompilerOptions_OutputJs is what type of JS file should be output by this compilation
type CompilerOptions_OutputJs int32

const (
	// Don't output anything.
	CompilerOptions_OutputJs_NONE CompilerOptions_OutputJs = iota
	// Output a "sentinel" file containing just a comment.
	CompilerOptions_OutputJs_SENTINEL
	// Output the compiled JS.
	CompilerOptions_OutputJs_NORMAL
)

var compilerOptions_OutputJs_names = [...]string{"NONE", "SENTINEL", "NORMAL"}

// CompilerOptions_OutputJs_values is OutputJs.values().
func CompilerOptions_OutputJs_values() []CompilerOptions_OutputJs {
	return []CompilerOptions_OutputJs{CompilerOptions_OutputJs_NONE, CompilerOptions_OutputJs_SENTINEL, CompilerOptions_OutputJs_NORMAL}
}

// String is Enum.name().
func (e CompilerOptions_OutputJs) String() string { return compilerOptions_OutputJs_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_OutputJs) Ordinal() int { return int(e) }

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
func (o *CompilerOptions) GetExtraReportGenerators() *jutil.ArrayList[SortingErrorManager_ErrorReportGenerator] {
	return o.extraReportGenerators
}

// Java: CompilerOptions.addReportGenerator(SortingErrorManager.ErrorReportGenerator)
func (o *CompilerOptions) AddReportGenerator(generator SortingErrorManager_ErrorReportGenerator) {
	o.extraReportGenerators.Add(generator)
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

// CompilerOptions_ConformanceReportingMode is whether to respect the library-level non-allowlisted conformance
// behavior specified in the conformance config.
type CompilerOptions_ConformanceReportingMode int32

const (
	// Used when running conformance checks at the binary level. It ignores the
	// library-level conformance behavior specified in the conformance config and
	// always report the conformance violations. For example, setting this allows
	// the compiler to ignore the "RECORD_ONLY" library-level behavior specified in
	// the config and unconditionally report the non-allowlisted violations when it
	// is running at the binary level. (default)
	CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG CompilerOptions_ConformanceReportingMode = iota
	// Used when running conformance checks at the library level. It respects the
	// library-level conformance behavior specified in the conformance config and
	// only reports the conformance violations according to the behavior. For
	// example, setting this respects the "RECORD_ONLY" library-level behavior
	// specified in the Boq Web configs and only records the violations when it is
	// running at the library level, and also respects the "REPORT_AS_BUILD_ERROR"
	// behavior specified in the global conformance configs to report errors for
	// go/jscp at the library level.
	CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
)

var compilerOptions_ConformanceReportingMode_names = [...]string{"IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG", "RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG"}

// CompilerOptions_ConformanceReportingMode_values is ConformanceReportingMode.values().
func CompilerOptions_ConformanceReportingMode_values() []CompilerOptions_ConformanceReportingMode {
	return []CompilerOptions_ConformanceReportingMode{CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG, CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG}
}

// String is Enum.name().
func (e CompilerOptions_ConformanceReportingMode) String() string {
	return compilerOptions_ConformanceReportingMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_ConformanceReportingMode) Ordinal() int { return int(e) }

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
	o.experimentalForceTranspiles = []CompilerOptions_ExperimentalForceTranspile{}
	o.inferConsts = true
	o.incrementalCheckMode = CompilerOptions_IncrementalCheckMode_OFF
	o.parseJsDocDocumentation = parsing.Config_JsDocParsing_TYPES_ONLY
	o.dependencyOptions = DependencyOptions_None()
	o.numParallelThreads = 1
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
	o.extraReportGenerators = jutil.NewArrayList[SortingErrorManager_ErrorReportGenerator]()
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
	o.deadAssignmentElimination = false
	o.deadPropertyAssignmentElimination = base.Tri_UNKNOWN
	o.inlineConstantVars = false
	o.inlineFunctionsLevel = CompilerOptions_Reach_NONE
	o.maxFunctionSizeAfterInlining = CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING
	o.assumeStrictThis = false
	o.assumeClosuresOnlyCaptureReferences = false
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
	o.defineReplacements = jutil.NewLinkedHashMap[string, any]()
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

// Returns the map of define replacements.
//
// Java: CompilerOptions.getDefineReplacements()
func (o *CompilerOptions) GetDefineReplacements() *jutil.LinkedHashMap[string, *rhino.Node] {
	m := jutil.NewLinkedHashMap[string, *rhino.Node]()
	for name, value := range o.defineReplacements.All() {
		switch v := value.(type) {
		case bool:
			m.Put(name, NodeUtil_BooleanNode(v))
		case int32:
			m.Put(name, NodeUtil_NumberNode(float64(v), nil))
		case float64:
			m.Put(name, NodeUtil_NumberNode(v, nil))
		case string:
			m.Put(name, rhino.IR_String(v))
		default:
			panic(jutil.NewIllegalStateExceptionMsg(jutil.ToJavaString(value)))
		}
	}
	return m
}

// Sets the value of the @define variable in JS to a boolean literal.
//
// Java: CompilerOptions.setDefineToBooleanLiteral(String, boolean)
func (o *CompilerOptions) SetDefineToBooleanLiteral(defineName string, value bool) {
	o.defineReplacements.Put(defineName, value)
}

// Sets the value of the @define variable in JS to a String literal.
//
// Java: CompilerOptions.setDefineToStringLiteral(String, String)
func (o *CompilerOptions) SetDefineToStringLiteral(defineName string, value string) {
	o.defineReplacements.Put(defineName, value)
}

// Sets the value of the @define variable in JS to a number literal.
//
// Java: CompilerOptions.setDefineToNumberLiteral(String, int)
func (o *CompilerOptions) SetDefineToNumberLiteral(defineName string, value int) {
	o.defineReplacements.Put(defineName, int32(value)) // boxed as Integer
}

// Sets the value of the @define variable in JS to a number literal.
//
// Java: CompilerOptions.setDefineToDoubleLiteral(String, double)
func (o *CompilerOptions) SetDefineToDoubleLiteral(defineName string, value float64) {
	o.defineReplacements.Put(defineName, value)
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

// CompilerOptions_ExperimentalForceTranspile are options to force transpile specific features for performance
// experiments.
//
// Deprecated.
type CompilerOptions_ExperimentalForceTranspile int32

const (
	// Transpile performance sensitive features away, preserving ASYNC/AWAIT for
	// debuggability
	//
	// Deprecated: Use setLanguageOut(LanguageMode.ECMASCRIPT5), or
	// setLangaugeOut(LanguageMode.ECMASCRIPT_2017) instead.
	CompilerOptions_ExperimentalForceTranspile_PERFORMANT_WITH_ASYNC_STACKS CompilerOptions_ExperimentalForceTranspile = iota
)

var compilerOptions_ExperimentalForceTranspile_names = [...]string{"PERFORMANT_WITH_ASYNC_STACKS"}

// CompilerOptions_ExperimentalForceTranspile_values is ExperimentalForceTranspile.values().
func CompilerOptions_ExperimentalForceTranspile_values() []CompilerOptions_ExperimentalForceTranspile {
	return []CompilerOptions_ExperimentalForceTranspile{CompilerOptions_ExperimentalForceTranspile_PERFORMANT_WITH_ASYNC_STACKS}
}

// String is Enum.name().
func (e CompilerOptions_ExperimentalForceTranspile) String() string {
	return compilerOptions_ExperimentalForceTranspile_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_ExperimentalForceTranspile) Ordinal() int { return int(e) }

// Java: CompilerOptions.setExperimentalForceTranspiles(ExperimentalForceTranspile...)
func (o *CompilerOptions) SetExperimentalForceTranspiles(experimentalForceTranspile ...CompilerOptions_ExperimentalForceTranspile) {
	if len(experimentalForceTranspile) > 1 {
		jutil.CheckState(
			!slices.Contains(o.experimentalForceTranspiles,
				CompilerOptions_ExperimentalForceTranspile_PERFORMANT_WITH_ASYNC_STACKS))
	}
	o.experimentalForceTranspiles = append([]CompilerOptions_ExperimentalForceTranspile{}, experimentalForceTranspile...)
}

// Java: CompilerOptions.getExperimentalForceTranspiles()
func (o *CompilerOptions) GetExperimentalForceTranspiles() []CompilerOptions_ExperimentalForceTranspile {
	return o.experimentalForceTranspiles
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

// @return Whether assumeClosuresOnlyCaptureReferences is set.
//
// Java: CompilerOptions.assumeClosuresOnlyCaptureReferences()
func (o *CompilerOptions) AssumeClosuresOnlyCaptureReferences() bool {
	return o.assumeClosuresOnlyCaptureReferences
}

// Whether to assume closures capture only what they reference. This allows more aggressive
// function inlining.
//
// Java: CompilerOptions.setAssumeClosuresOnlyCaptureReferences(boolean)
func (o *CompilerOptions) SetAssumeClosuresOnlyCaptureReferences(enable bool) {
	o.assumeClosuresOnlyCaptureReferences = enable
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

// Ensures that functions labelled with `@requireInlining` JSDoc are, in fact, actually inlined by
// the optimizer.
//
// This may only be set when using advanced optimizations, as otherwise the compiler is unable
// to do most inlining. The compiler will throw an exception during initialization if this is set
// to true, but advanced optimizations are not being used.
//
// Java: CompilerOptions.setValidateRequiredInlinings(boolean)
func (o *CompilerOptions) SetValidateRequiredInlinings(validateRequireInliningAnnotation bool) {
	o.validateRequireInliningAnnotation = base.Tri_ForBoolean(validateRequireInliningAnnotation)
}

// Java: CompilerOptions.getShouldValidateRequiredInlinings()
func (o *CompilerOptions) GetShouldValidateRequiredInlinings() base.Tri {
	return o.validateRequireInliningAnnotation
}

// Java: CompilerOptions.setPreserveDetailedSourceInfo(boolean)
func (o *CompilerOptions) SetPreserveDetailedSourceInfo(preserveDetailedSourceInfo bool) {
	o.preserveDetailedSourceInfo = preserveDetailedSourceInfo
}

// Java: CompilerOptions.preservesDetailedSourceInfo()
func (o *CompilerOptions) PreservesDetailedSourceInfo() bool {
	return o.preserveDetailedSourceInfo
}

// Java: CompilerOptions.setPreserveNonJSDocComments(boolean)
func (o *CompilerOptions) SetPreserveNonJSDocComments(preserveNonJSDocComments bool) {
	o.preserveNonJSDocComments = preserveNonJSDocComments
}

// Java: CompilerOptions.getPreserveNonJSDocComments()
func (o *CompilerOptions) GetPreserveNonJSDocComments() bool {
	return o.preserveNonJSDocComments
}

// Java: CompilerOptions.setContinueAfterErrors(boolean)
func (o *CompilerOptions) SetContinueAfterErrors(continueAfterErrors bool) {
	o.continueAfterErrors = continueAfterErrors
}

// Java: CompilerOptions.canContinueAfterErrors()
func (o *CompilerOptions) CanContinueAfterErrors() bool {
	return o.continueAfterErrors
}

// Enables or disables the parsing of JSDoc documentation, and optionally also the preservation of
// all whitespace and formatting within a JSDoc comment. By default, whitespace is collapsed for
// all comments except {@literal @license} and {@literal @preserve} blocks,
//
// Java: CompilerOptions.setParseJsDocDocumentation(Config.JsDocParsing)
func (o *CompilerOptions) SetParseJsDocDocumentation(parseJsDocDocumentation parsing.Config_JsDocParsing) {
	o.parseJsDocDocumentation = parseJsDocDocumentation
}

// Checks JSDoc documentation will be parsed.
//
// @return True when JSDoc documentation will be parsed, false if not.
//
// Java: CompilerOptions.isParseJsDocDocumentation()
func (o *CompilerOptions) IsParseJsDocDocumentation() parsing.Config_JsDocParsing {
	return o.parseJsDocDocumentation
}

// Skip all passes (other than transpilation, if requested). Don't inject any runtime libraries
// (unless explicitly requested) or do any checks/optimizations (this is useful for per-file
// transpilation).
//
// Java: CompilerOptions.setSkipNonTranspilationPasses(boolean)
func (o *CompilerOptions) SetSkipNonTranspilationPasses(skipNonTranspilationPasses bool) {
	o.skipNonTranspilationPasses = skipNonTranspilationPasses
}

// Java: CompilerOptions.getSkipNonTranspilationPasses()
func (o *CompilerOptions) GetSkipNonTranspilationPasses() bool {
	return o.skipNonTranspilationPasses
}

// Java: CompilerOptions.setDevMode(DevMode)
func (o *CompilerOptions) SetDevMode(devMode CompilerOptions_DevMode) {
	o.devMode = devMode
}

// Java: CompilerOptions.getDevMode()
func (o *CompilerOptions) GetDevMode() CompilerOptions_DevMode {
	return o.devMode
}

// Java: CompilerOptions.setCheckDeterminism(boolean)
func (o *CompilerOptions) SetCheckDeterminism(checkDeterminism bool) {
	o.checkDeterminism = checkDeterminism
}

// Java: CompilerOptions.getCheckDeterminism()
func (o *CompilerOptions) GetCheckDeterminism() bool {
	return o.checkDeterminism
}

// Java: CompilerOptions.setMessageBundle(MessageBundle)
func (o *CompilerOptions) SetMessageBundle(messageBundle MessageBundle) {
	o.messageBundle = messageBundle
}

// Java: CompilerOptions.getMessageBundle()
func (o *CompilerOptions) GetMessageBundle() MessageBundle {
	return o.messageBundle
}

// Java: CompilerOptions.setCheckSymbols(boolean)
func (o *CompilerOptions) SetCheckSymbols(checkSymbols bool) {
	o.checkSymbols = checkSymbols
}

// Java: CompilerOptions.getCheckSymbols()
func (o *CompilerOptions) GetCheckSymbols() bool {
	return o.checkSymbols
}

// Java: CompilerOptions.setCheckSuspiciousCode(boolean)
func (o *CompilerOptions) SetCheckSuspiciousCode(checkSuspiciousCode bool) {
	o.checkSuspiciousCode = checkSuspiciousCode
}

// Java: CompilerOptions.getCheckSuspiciousCode()
func (o *CompilerOptions) GetCheckSuspiciousCode() bool {
	return o.checkSuspiciousCode
}

// Java: CompilerOptions.setCheckTypes(boolean)
func (o *CompilerOptions) SetCheckTypes(checkTypes bool) {
	o.checkTypes = checkTypes
}

// Java: CompilerOptions.getCheckTypes()
func (o *CompilerOptions) GetCheckTypes() bool {
	return o.checkTypes
}

// Java: CompilerOptions.setFoldConstants(boolean)
func (o *CompilerOptions) SetFoldConstants(foldConstants bool) {
	o.foldConstants = foldConstants
}

// Java: CompilerOptions.shouldFoldConstants()
func (o *CompilerOptions) ShouldFoldConstants() bool {
	return o.foldConstants
}

// Java: CompilerOptions.setDeadAssignmentElimination(boolean)
func (o *CompilerOptions) SetDeadAssignmentElimination(deadAssignmentElimination bool) {
	o.deadAssignmentElimination = deadAssignmentElimination
}

// Java: CompilerOptions.shouldRunDeadAssignmentElimination()
func (o *CompilerOptions) ShouldRunDeadAssignmentElimination() bool {
	return o.deadAssignmentElimination
}

// Java: CompilerOptions.setDeadPropertyAssignmentElimination(boolean)
func (o *CompilerOptions) SetDeadPropertyAssignmentElimination(deadPropertyAssignmentElimination bool) {
	o.deadPropertyAssignmentElimination = base.Tri_ForBoolean(deadPropertyAssignmentElimination)
}

// Java: CompilerOptions.shouldRunDeadPropertyAssignmentElimination()
func (o *CompilerOptions) ShouldRunDeadPropertyAssignmentElimination() bool {
	// By default, run dead property assignment elimination when both
	//   1. regular dead assignment elimiantion is enabled.
	//   2. the polymer pass is disabled.
	// This is because Polymer relies heavily on getters and setters. Those are very incompatible
	// eliminating seemingly "dead" property assignments. See b/28912819.
	return o.deadPropertyAssignmentElimination.ToBoolean(
		o.deadAssignmentElimination && !o.polymerPass)
}

// Java: CompilerOptions.setInlineConstantVars(boolean)
func (o *CompilerOptions) SetInlineConstantVars(inlineConstantVars bool) {
	o.inlineConstantVars = inlineConstantVars
}

// Java: CompilerOptions.shouldInlineConstantVars()
func (o *CompilerOptions) ShouldInlineConstantVars() bool {
	return o.inlineConstantVars
}

// Java: CompilerOptions.setCrossChunkCodeMotion(boolean)
func (o *CompilerOptions) SetCrossChunkCodeMotion(crossChunkCodeMotion bool) {
	o.crossChunkCodeMotion = crossChunkCodeMotion
}

// Java: CompilerOptions.setCrossChunkCodeMotionNoStubMethods(boolean)
func (o *CompilerOptions) SetCrossChunkCodeMotionNoStubMethods(crossChunkCodeMotionNoStubMethods bool) {
	o.crossChunkCodeMotionNoStubMethods = crossChunkCodeMotionNoStubMethods
}

// Java: CompilerOptions.getCrossChunkCodeMotionNoStubMethods()
func (o *CompilerOptions) GetCrossChunkCodeMotionNoStubMethods() bool {
	return o.crossChunkCodeMotionNoStubMethods
}

// Java: CompilerOptions.setParentChunkCanSeeSymbolsDeclaredInChildren(boolean)
func (o *CompilerOptions) SetParentChunkCanSeeSymbolsDeclaredInChildren(parentChunkCanSeeSymbolsDeclaredInChildren bool) {
	o.parentChunkCanSeeSymbolsDeclaredInChildren = parentChunkCanSeeSymbolsDeclaredInChildren
}

// Java: CompilerOptions.getParentChunkCanSeeSymbolsDeclaredInChildren()
func (o *CompilerOptions) GetParentChunkCanSeeSymbolsDeclaredInChildren() bool {
	return o.parentChunkCanSeeSymbolsDeclaredInChildren
}

// Java: CompilerOptions.setCrossChunkMethodMotion(boolean)
func (o *CompilerOptions) SetCrossChunkMethodMotion(crossChunkMethodMotion bool) {
	o.crossChunkMethodMotion = crossChunkMethodMotion
}

// Java: CompilerOptions.getCrossChunkMethodMotion()
func (o *CompilerOptions) GetCrossChunkMethodMotion() bool {
	return o.crossChunkMethodMotion
}

// Java: CompilerOptions.setCoalesceVariableNames(boolean)
func (o *CompilerOptions) SetCoalesceVariableNames(coalesceVariableNames bool) {
	o.coalesceVariableNames = coalesceVariableNames
}

// Java: CompilerOptions.shouldCoalesceVariableNames()
func (o *CompilerOptions) ShouldCoalesceVariableNames() bool {
	return o.coalesceVariableNames
}

// Java: CompilerOptions.setInlineLocalVariables(boolean)
func (o *CompilerOptions) SetInlineLocalVariables(inlineLocalVariables bool) {
	o.inlineLocalVariables = inlineLocalVariables
}

// Java: CompilerOptions.shouldInlineLocalVariables()
func (o *CompilerOptions) ShouldInlineLocalVariables() bool {
	return o.inlineLocalVariables
}

// Java: CompilerOptions.setFlowSensitiveInlineVariables(boolean)
func (o *CompilerOptions) SetFlowSensitiveInlineVariables(enabled bool) {
	o.flowSensitiveInlineVariables = enabled
}

// Java: CompilerOptions.getFlowSensitiveInlineVariables()
func (o *CompilerOptions) GetFlowSensitiveInlineVariables() bool {
	return o.flowSensitiveInlineVariables
}

// Java: CompilerOptions.setSmartNameRemoval(boolean)
func (o *CompilerOptions) SetSmartNameRemoval(smartNameRemoval bool) {
	// TODO(bradfordcsmith): Remove the smart name removal option.
	o.smartNameRemoval = smartNameRemoval
	if smartNameRemoval {
		// To get the effect this option used to have we need to enable these options.
		// Don't disable them here if they were set explicitly, though.
		o.removeUnusedVars = true
		o.removeUnusedPrototypeProperties = true
	}
}

// Java: CompilerOptions.getSmartNameRemoval()
func (o *CompilerOptions) GetSmartNameRemoval() bool {
	return o.smartNameRemoval
}

// Java: CompilerOptions.setExtractPrototypeMemberDeclarations(boolean)
func (o *CompilerOptions) SetExtractPrototypeMemberDeclarations(enabled bool) {
	if enabled {
		o.extractPrototypeMemberDeclarations = CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP
	} else {
		o.extractPrototypeMemberDeclarations = CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF
	}
}

// USE_IIFE is currently unused. Consider removing support for it and
// deleting this setter.
//
// Java: CompilerOptions.setExtractPrototypeMemberDeclarations(ExtractPrototypeMemberDeclarationsMode)
func (o *CompilerOptions) SetExtractPrototypeMemberDeclarationsOfMode(mode CompilerOptions_ExtractPrototypeMemberDeclarationsMode) {
	o.extractPrototypeMemberDeclarations = mode
}

// Java: CompilerOptions.getExtractPrototypeMemberDeclarationsMode()
func (o *CompilerOptions) GetExtractPrototypeMemberDeclarationsMode() CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return o.extractPrototypeMemberDeclarations
}

// Java: CompilerOptions.setRemoveUnusedPrototypeProperties(boolean)
func (o *CompilerOptions) SetRemoveUnusedPrototypeProperties(enabled bool) {
	o.removeUnusedPrototypeProperties = enabled
	// InlineSimpleMethods makes similar assumptions to
	// RemoveUnusedCode, so they are enabled together.
	o.inlineGetters = enabled
}

// Java: CompilerOptions.shouldRemoveUnusedPrototypeProperties()
func (o *CompilerOptions) ShouldRemoveUnusedPrototypeProperties() bool {
	return o.removeUnusedPrototypeProperties
}

// Java: CompilerOptions.shouldInlineGetters()
func (o *CompilerOptions) ShouldInlineGetters() bool {
	return o.inlineGetters
}

// Java: CompilerOptions.setCollapseVariableDeclarations(boolean)
func (o *CompilerOptions) SetCollapseVariableDeclarations(enabled bool) {
	o.collapseVariableDeclarations = enabled
}

// Java: CompilerOptions.shouldCollapseVariableDeclarations()
func (o *CompilerOptions) ShouldCollapseVariableDeclarations() bool {
	return o.collapseVariableDeclarations
}

// Java: CompilerOptions.setCollapseAnonymousFunctions(boolean)
func (o *CompilerOptions) SetCollapseAnonymousFunctions(enabled bool) {
	o.collapseAnonymousFunctions = enabled
}

// Java: CompilerOptions.shouldCollapseAnonymousFunctions()
func (o *CompilerOptions) ShouldCollapseAnonymousFunctions() bool {
	return o.collapseAnonymousFunctions
}

// Java: CompilerOptions.setAliasStringsMode(AliasStringsMode)
func (o *CompilerOptions) SetAliasStringsMode(aliasStringsMode CompilerOptions_AliasStringsMode) {
	o.aliasStringsMode = aliasStringsMode
}

// Java: CompilerOptions.getAliasStringsMode()
func (o *CompilerOptions) GetAliasStringsMode() CompilerOptions_AliasStringsMode {
	return o.aliasStringsMode
}

// Java: CompilerOptions.setOutputJsStringUsage(boolean)
func (o *CompilerOptions) SetOutputJsStringUsage(outputJsStringUsage bool) {
	o.outputJsStringUsage = outputJsStringUsage
}

// Java: CompilerOptions.shouldOutputJsStringUsage()
func (o *CompilerOptions) ShouldOutputJsStringUsage() bool {
	return o.outputJsStringUsage
}

// Java: CompilerOptions.setConvertToDottedProperties(boolean)
func (o *CompilerOptions) SetConvertToDottedProperties(convertToDottedProperties bool) {
	o.convertToDottedProperties = convertToDottedProperties
}

// Java: CompilerOptions.shouldConvertToDottedProperties()
func (o *CompilerOptions) ShouldConvertToDottedProperties() bool {
	return o.convertToDottedProperties
}

// Java: CompilerOptions.setUseTypesForLocalOptimization(boolean)
func (o *CompilerOptions) SetUseTypesForLocalOptimization(useTypesForLocalOptimization bool) {
	o.useTypesForLocalOptimization = useTypesForLocalOptimization
}

// Java: CompilerOptions.shouldUseTypesForLocalOptimization()
func (o *CompilerOptions) ShouldUseTypesForLocalOptimization() bool {
	return o.useTypesForLocalOptimization
}

// Java: CompilerOptions.setUseTypesForOptimization(boolean)
func (o *CompilerOptions) SetUseTypesForOptimization(useTypesForOptimization bool) {
	if useTypesForOptimization {
		o.disambiguateProperties = true
		o.ambiguateProperties = true
		o.inlineProperties = true
		o.useTypesForLocalOptimization = true
	}
}

// Whether any type-based optimizations are enabled
//
// Java: CompilerOptions.requiresTypesForOptimization()
func (o *CompilerOptions) RequiresTypesForOptimization() bool {
	// Add new type-based optimization options to this list, or types may be erased from the AST
	// before the optimziations run.
	return o.disambiguateProperties ||
		o.ambiguateProperties ||
		o.inlineProperties ||
		o.useTypesForLocalOptimization
}

// Java: CompilerOptions.setRewriteFunctionExpressions(boolean)
func (o *CompilerOptions) SetRewriteFunctionExpressions(rewriteFunctionExpressions bool) {
	o.rewriteFunctionExpressions = rewriteFunctionExpressions
}

// Java: CompilerOptions.shouldRewriteFunctionExpressions()
func (o *CompilerOptions) ShouldRewriteFunctionExpressions() bool {
	return o.rewriteFunctionExpressions
}

// Java: CompilerOptions.setOptimizeCalls(boolean)
func (o *CompilerOptions) SetOptimizeCalls(optimizeCalls bool) {
	o.optimizeCalls = optimizeCalls
}

// Java: CompilerOptions.shouldOptimizeCalls()
func (o *CompilerOptions) ShouldOptimizeCalls() bool {
	return o.optimizeCalls
}

// Java: CompilerOptions.getOptimizeESClassConstructors()
func (o *CompilerOptions) GetOptimizeESClassConstructors() bool {
	return o.optimizeESClassConstructors
}

// Java: CompilerOptions.setOptimizeESClassConstructors(boolean)
func (o *CompilerOptions) SetOptimizeESClassConstructors(optimizeESClassConstructors bool) {
	o.optimizeESClassConstructors = optimizeESClassConstructors
}

// Java: CompilerOptions.setVariableRenaming(VariableRenamingPolicy)
func (o *CompilerOptions) SetVariableRenaming(variableRenaming VariableRenamingPolicy) {
	o.variableRenaming = variableRenaming
}

// Java: CompilerOptions.getVariableRenaming()
func (o *CompilerOptions) GetVariableRenaming() VariableRenamingPolicy {
	return o.variableRenaming
}

// Java: CompilerOptions.setPropertyRenaming(PropertyRenamingPolicy)
func (o *CompilerOptions) SetPropertyRenaming(propertyRenaming PropertyRenamingPolicy) {
	o.propertyRenaming = propertyRenaming
}

// Java: CompilerOptions.getPropertyRenaming()
func (o *CompilerOptions) GetPropertyRenaming() PropertyRenamingPolicy {
	return o.propertyRenaming
}

// Java: CompilerOptions.setLabelRenaming(boolean)
func (o *CompilerOptions) SetLabelRenaming(labelRenaming bool) {
	o.labelRenaming = labelRenaming
}

// Java: CompilerOptions.shouldRenameLabels()
func (o *CompilerOptions) ShouldRenameLabels() bool {
	return o.labelRenaming
}

// Java: CompilerOptions.setReserveRawExports(boolean)
func (o *CompilerOptions) SetReserveRawExports(reserveRawExports bool) {
	o.reserveRawExports = reserveRawExports
}

// Java: CompilerOptions.shouldReserveRawExports()
func (o *CompilerOptions) ShouldReserveRawExports() bool {
	return o.reserveRawExports
}

// Java: CompilerOptions.setPreferStableNames(boolean)
func (o *CompilerOptions) SetPreferStableNames(preferStableNames bool) {
	o.preferStableNames = preferStableNames
}

// Java: CompilerOptions.shouldPreferStableNames()
func (o *CompilerOptions) ShouldPreferStableNames() bool {
	return o.preferStableNames
}

// Java: CompilerOptions.setGeneratePseudoNames(boolean)
func (o *CompilerOptions) SetGeneratePseudoNames(generatePseudoNames bool) {
	o.generatePseudoNames = generatePseudoNames
}

// Java: CompilerOptions.shouldGeneratePseudoNames()
func (o *CompilerOptions) ShouldGeneratePseudoNames() bool {
	return o.generatePseudoNames
}

// Java: CompilerOptions.setPropertyRenamingOnlyCompilationMode(boolean)
func (o *CompilerOptions) SetPropertyRenamingOnlyCompilationMode(propertyRenamingOnlyCompilationMode bool) {
	o.propertyRenamingOnlyCompilationMode = propertyRenamingOnlyCompilationMode
}

// Java: CompilerOptions.isPropertyRenamingOnlyCompilationMode()
func (o *CompilerOptions) IsPropertyRenamingOnlyCompilationMode() bool {
	return o.propertyRenamingOnlyCompilationMode
}

// Java: CompilerOptions.setRenamePrefix(String)
func (o *CompilerOptions) SetRenamePrefix(renamePrefix *string) {
	o.renamePrefix = renamePrefix
}

// Java: CompilerOptions.getRenamePrefix()
func (o *CompilerOptions) GetRenamePrefix() *string {
	return o.renamePrefix
}

// Java: CompilerOptions.getRenamePrefixNamespace()
func (o *CompilerOptions) GetRenamePrefixNamespace() *string {
	return o.renamePrefixNamespace
}

// Java: CompilerOptions.setRenamePrefixNamespace(String)
func (o *CompilerOptions) SetRenamePrefixNamespace(renamePrefixNamespace *string) {
	o.renamePrefixNamespace = renamePrefixNamespace
}

// Java: CompilerOptions.setCollapsePropertiesLevel(PropertyCollapseLevel)
func (o *CompilerOptions) SetCollapsePropertiesLevel(level CompilerOptions_PropertyCollapseLevel) {
	o.collapsePropertiesLevel = level
}

// Java: CompilerOptions.setCollapseProperties(boolean)
func (o *CompilerOptions) SetCollapseProperties(fullyCollapse bool) {
	if fullyCollapse {
		o.collapsePropertiesLevel = CompilerOptions_PropertyCollapseLevel_ALL
	} else {
		o.collapsePropertiesLevel = CompilerOptions_PropertyCollapseLevel_NONE
	}
}

// Java: CompilerOptions.getCollapsePropertiesLevel()
func (o *CompilerOptions) GetCollapsePropertiesLevel() CompilerOptions_PropertyCollapseLevel {
	return o.collapsePropertiesLevel
}

// Java: CompilerOptions.setDevirtualizeMethods(boolean)
func (o *CompilerOptions) SetDevirtualizeMethods(devirtualizeMethods bool) {
	o.devirtualizeMethods = devirtualizeMethods
}

// Java: CompilerOptions.shouldDevirtualizeMethods()
func (o *CompilerOptions) ShouldDevirtualizeMethods() bool {
	return o.devirtualizeMethods
}

// Java: CompilerOptions.setComputeFunctionSideEffects(boolean)
func (o *CompilerOptions) SetComputeFunctionSideEffects(computeFunctionSideEffects bool) {
	o.computeFunctionSideEffects = computeFunctionSideEffects
}

// Java: CompilerOptions.shouldComputeFunctionSideEffects()
func (o *CompilerOptions) ShouldComputeFunctionSideEffects() bool {
	return o.computeFunctionSideEffects
}

// Java: CompilerOptions.setDisambiguateProperties(boolean)
func (o *CompilerOptions) SetDisambiguateProperties(disambiguateProperties bool) {
	o.disambiguateProperties = disambiguateProperties
}

// Java: CompilerOptions.shouldDisambiguateProperties()
func (o *CompilerOptions) ShouldDisambiguateProperties() bool {
	return o.disambiguateProperties
}

// Java: CompilerOptions.setAmbiguateProperties(boolean)
func (o *CompilerOptions) SetAmbiguateProperties(ambiguateProperties bool) {
	o.ambiguateProperties = ambiguateProperties
}

// Java: CompilerOptions.shouldAmbiguateProperties()
func (o *CompilerOptions) ShouldAmbiguateProperties() bool {
	return o.ambiguateProperties
}

// Java: CompilerOptions.setInputVariableMap(VariableMap)
func (o *CompilerOptions) SetInputVariableMap(inputVariableMap *VariableMap) {
	o.inputVariableMap = inputVariableMap
}

// Java: CompilerOptions.getInputVariableMap()
func (o *CompilerOptions) GetInputVariableMap() *VariableMap {
	return o.inputVariableMap
}

// Java: CompilerOptions.setInputPropertyMap(VariableMap)
func (o *CompilerOptions) SetInputPropertyMap(inputPropertyMap *VariableMap) {
	o.inputPropertyMap = inputPropertyMap
}

// Java: CompilerOptions.getInputPropertyMap()
func (o *CompilerOptions) GetInputPropertyMap() *VariableMap {
	return o.inputPropertyMap
}

// Java: CompilerOptions.setExportTestFunctions(boolean)
func (o *CompilerOptions) SetExportTestFunctions(exportTestFunctions bool) {
	o.exportTestFunctions = exportTestFunctions
}

// Java: CompilerOptions.shouldExportTestFunctions()
func (o *CompilerOptions) ShouldExportTestFunctions() bool {
	return o.exportTestFunctions
}

// Java: CompilerOptions.setSyntheticBlockStartMarker(String)
func (o *CompilerOptions) SetSyntheticBlockStartMarker(syntheticBlockStartMarker *string) {
	o.syntheticBlockStartMarker = syntheticBlockStartMarker
}

// Java: CompilerOptions.setSyntheticBlockEndMarker(String)
func (o *CompilerOptions) SetSyntheticBlockEndMarker(syntheticBlockEndMarker *string) {
	o.syntheticBlockEndMarker = syntheticBlockEndMarker
}

// Java: CompilerOptions.getSyntheticBlockStartMarker()
func (o *CompilerOptions) GetSyntheticBlockStartMarker() *string {
	return o.syntheticBlockStartMarker
}

// Java: CompilerOptions.getSyntheticBlockEndMarker()
func (o *CompilerOptions) GetSyntheticBlockEndMarker() *string {
	return o.syntheticBlockEndMarker
}

// Java: CompilerOptions.setLocale(String)
func (o *CompilerOptions) SetLocale(locale *string) {
	o.locale = locale
}

// Java: CompilerOptions.getLocale()
func (o *CompilerOptions) GetLocale() *string {
	return o.locale
}

// Java: CompilerOptions.setDoLateLocalization(boolean)
func (o *CompilerOptions) SetDoLateLocalization(doLateLocalization bool) {
	o.doLateLocalization = doLateLocalization
}

// Java: CompilerOptions.doLateLocalization()
func (o *CompilerOptions) DoLateLocalization() bool {
	return o.doLateLocalization
}

// Should we run any form of the `ReplaceMessages` pass?
//
// Java: CompilerOptions.shouldRunReplaceMessagesPass()
func (o *CompilerOptions) ShouldRunReplaceMessagesPass() bool {
	return !o.ShouldRunReplaceMessagesForChrome() && o.messageBundle != nil
}

// Java: CompilerOptions.setMarkAsCompiled(boolean)
func (o *CompilerOptions) SetMarkAsCompiled(markAsCompiled bool) {
	o.markAsCompiled = markAsCompiled
}

// Java: CompilerOptions.shouldMarkAsCompiled()
func (o *CompilerOptions) ShouldMarkAsCompiled() bool {
	return o.markAsCompiled
}

// Java: CompilerOptions.setClosurePass(boolean)
func (o *CompilerOptions) SetClosurePass(closurePass bool) {
	o.closurePass = closurePass
}

// Java: CompilerOptions.getClosurePass()
func (o *CompilerOptions) GetClosurePass() bool {
	return o.closurePass
}

// Preserve closure primitives.
//
// For now, this only preserves goog.provide(), goog.require() and goog.module() calls.
//
// Java: CompilerOptions.setPreserveClosurePrimitives(boolean)
func (o *CompilerOptions) SetPreserveClosurePrimitives(preserveClosurePrimitives bool) {
	o.preserveClosurePrimitives = preserveClosurePrimitives
}

// Java: CompilerOptions.shouldPreservesGoogProvidesAndRequires()
func (o *CompilerOptions) ShouldPreservesGoogProvidesAndRequires() bool {
	return o.preserveClosurePrimitives
}

// Java: CompilerOptions.shouldPreserveGoogModule()
func (o *CompilerOptions) ShouldPreserveGoogModule() bool {
	return o.preserveClosurePrimitives
}

// Do not process goog. intrinsics, such as goog.getCssName().
//
// Java: CompilerOptions.shouldPreserveGoogLibraryPrimitives()
func (o *CompilerOptions) ShouldPreserveGoogLibraryPrimitives() bool {
	return o.preserveClosurePrimitives
}

// Java: CompilerOptions.setPreserveTypeAnnotations(boolean)
func (o *CompilerOptions) SetPreserveTypeAnnotations(preserveTypeAnnotations bool) {
	o.preserveTypeAnnotations = preserveTypeAnnotations
}

// Java: CompilerOptions.shouldPreserveTypeAnnotations()
func (o *CompilerOptions) ShouldPreserveTypeAnnotations() bool {
	return o.preserveTypeAnnotations
}

// Java: CompilerOptions.setGentsMode(boolean)
func (o *CompilerOptions) SetGentsMode(gentsMode bool) {
	o.gentsMode = gentsMode
}

// Java: CompilerOptions.getGentsMode()
func (o *CompilerOptions) GetGentsMode() bool {
	return o.gentsMode
}

// Java: CompilerOptions.setGatherCssNames(boolean)
func (o *CompilerOptions) SetGatherCssNames(gatherCssNames bool) {
	o.gatherCssNames = gatherCssNames
}

// Java: CompilerOptions.shouldGatherCssNames()
func (o *CompilerOptions) ShouldGatherCssNames() bool {
	return o.gatherCssNames
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.setStripTypes(Set<String>)
func (o *CompilerOptions) SetStripTypes(stripTypes *jutil.LinkedHashSet[string]) {
	o.stripTypes = jutil.NewLinkedHashSet(stripTypes.ToSlice()...)
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.getStripTypes()
func (o *CompilerOptions) GetStripTypes() *jutil.LinkedHashSet[string] {
	return o.stripTypes
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.setStripNameSuffixes(Set<String>)
func (o *CompilerOptions) SetStripNameSuffixes(stripNameSuffixes *jutil.LinkedHashSet[string]) {
	o.stripNameSuffixes = jutil.NewLinkedHashSet(stripNameSuffixes.ToSlice()...)
}

// Java: CompilerOptions.getStripNameSuffixes()
func (o *CompilerOptions) GetStripNameSuffixes() *jutil.LinkedHashSet[string] {
	return o.stripNameSuffixes
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.setStripNamePrefixes(Set<String>)
func (o *CompilerOptions) SetStripNamePrefixes(stripNamePrefixes *jutil.LinkedHashSet[string]) {
	o.stripNamePrefixes = jutil.NewLinkedHashSet(stripNamePrefixes.ToSlice()...)
}

// Java: CompilerOptions.getStripNamePrefixes()
func (o *CompilerOptions) GetStripNamePrefixes() *jutil.LinkedHashSet[string] {
	return o.stripNamePrefixes
}

// Java: CompilerOptions.addCustomPass(CustomPassExecutionTime, CompilerPass)
func (o *CompilerOptions) AddCustomPass(time CustomPassExecutionTime, customPass CompilerPass) {
	if o.customPasses == nil {
		o.customPasses = jutil.NewSetMultimap[CustomPassExecutionTime, CompilerPass]()
	}
	o.customPasses.Put(time, customPass)
}

// Java: CompilerOptions.getCustomPassesAt(CustomPassExecutionTime)
func (o *CompilerOptions) GetCustomPassesAt(executionTime CustomPassExecutionTime) []CompilerPass {
	if o.customPasses == nil {
		return []CompilerPass{}
	}
	return o.customPasses.Get(executionTime).ToSlice()
}

// Java: CompilerOptions.setDefineReplacements(Map<String, Object>)
func (o *CompilerOptions) SetDefineReplacements(defineReplacements *jutil.LinkedHashMap[string, any]) {
	o.defineReplacements.Clear()
	for k, v := range defineReplacements.All() {
		o.defineReplacements.Put(k, v)
	}
}

// Java: CompilerOptions.getEnableZonesDefineName()
func (o *CompilerOptions) GetEnableZonesDefineName() *string {
	return o.enableZonesDefineName
}

// Java: CompilerOptions.setEnableZonesDefineName(String)
func (o *CompilerOptions) SetEnableZonesDefineName(enableZonesDefineName *string) {
	o.enableZonesDefineName = enableZonesDefineName
}

// Java: CompilerOptions.getZoneInputPattern()
func (o *CompilerOptions) GetZoneInputPattern() *jutil.Pattern {
	return o.zoneInputPattern
}

// Java: CompilerOptions.setZoneInputPattern(Pattern)
func (o *CompilerOptions) SetZoneInputPattern(zoneInputPattern *jutil.Pattern) {
	o.zoneInputPattern = zoneInputPattern
}

// Java: CompilerOptions.setRewriteGlobalDeclarationsForTryCatchWrapping(boolean)
func (o *CompilerOptions) SetRewriteGlobalDeclarationsForTryCatchWrapping(rewrite bool) {
	o.rewriteGlobalDeclarationsForTryCatchWrapping = rewrite
}

// Java: CompilerOptions.shouldRewriteGlobalDeclarationsForTryCatchWrapping()
func (o *CompilerOptions) ShouldRewriteGlobalDeclarationsForTryCatchWrapping() bool {
	return o.rewriteGlobalDeclarationsForTryCatchWrapping
}

// Java: CompilerOptions.setCssRenamingMap(CssRenamingMap)
func (o *CompilerOptions) SetCssRenamingMap(cssRenamingMap CssRenamingMap) {
	o.cssRenamingMap = cssRenamingMap
}

// Java: CompilerOptions.getCssRenamingMap()
func (o *CompilerOptions) GetCssRenamingMap() CssRenamingMap {
	return o.cssRenamingMap
}

// Java: CompilerOptions.setCssRenamingSkiplist(Set<String>)
func (o *CompilerOptions) SetCssRenamingSkiplist(skiplist *jutil.LinkedHashSet[string]) {
	o.cssRenamingSkiplist = jutil.NewLinkedHashSet(skiplist.ToSlice()...)
}

// Java: CompilerOptions.getCssRenamingSkiplist()
func (o *CompilerOptions) GetCssRenamingSkiplist() *jutil.LinkedHashSet[string] {
	return o.cssRenamingSkiplist
}

// Java: CompilerOptions.setReplaceStringsFunctionDescriptions(List<String>)
func (o *CompilerOptions) SetReplaceStringsFunctionDescriptions(replaceStringsFunctionDescriptions []string) {
	o.replaceStringsFunctionDescriptions = append([]string{}, replaceStringsFunctionDescriptions...)
}

// Java: CompilerOptions.getReplaceStringsFunctionDescriptions()
func (o *CompilerOptions) GetReplaceStringsFunctionDescriptions() []string {
	return o.replaceStringsFunctionDescriptions
}

// Java: CompilerOptions.setReplaceStringsPlaceholderToken(String)
func (o *CompilerOptions) SetReplaceStringsPlaceholderToken(replaceStringsPlaceholderToken string) {
	o.replaceStringsPlaceholderToken = replaceStringsPlaceholderToken
}

// Java: CompilerOptions.getReplaceStringsPlaceholderToken()
func (o *CompilerOptions) GetReplaceStringsPlaceholderToken() string {
	return o.replaceStringsPlaceholderToken
}

// Java: CompilerOptions.setPrettyPrint(boolean)
func (o *CompilerOptions) SetPrettyPrint(prettyPrint bool) {
	o.prettyPrint = prettyPrint
}

// Java: CompilerOptions.isPrettyPrint()
func (o *CompilerOptions) IsPrettyPrint() bool {
	return o.prettyPrint
}

// Java: CompilerOptions.setLineBreak(boolean)
func (o *CompilerOptions) SetLineBreak(lineBreak bool) {
	o.lineBreak = lineBreak
}

// Java: CompilerOptions.shouldAddLineBreak()
func (o *CompilerOptions) ShouldAddLineBreak() bool {
	return o.lineBreak
}

// Java: CompilerOptions.setPrintInputDelimiter(boolean)
func (o *CompilerOptions) SetPrintInputDelimiter(printInputDelimiter bool) {
	o.printInputDelimiter = printInputDelimiter
}

// Java: CompilerOptions.shouldPrintInputDelimiter()
func (o *CompilerOptions) ShouldPrintInputDelimiter() bool {
	return o.printInputDelimiter
}

// Java: CompilerOptions.setInputDelimiter(String)
func (o *CompilerOptions) SetInputDelimiter(inputDelimiter string) {
	o.inputDelimiter = inputDelimiter
}

// Java: CompilerOptions.getInputDelimiter()
func (o *CompilerOptions) GetInputDelimiter() string {
	return o.inputDelimiter
}

// Java: CompilerOptions.setDebugLogDirectory(Path)
func (o *CompilerOptions) SetDebugLogDirectory(dir *jutil.Path) {
	o.debugLogDirectory = dir
}

// Java: CompilerOptions.getDebugLogDirectory()
func (o *CompilerOptions) GetDebugLogDirectory() *jutil.Path {
	return o.debugLogDirectory
}

// Java: CompilerOptions.setDebugLogFilter(String)
func (o *CompilerOptions) SetDebugLogFilter(filter *string) {
	o.debugLogFilter = filter
}

// Java: CompilerOptions.getDebugLogFilter()
func (o *CompilerOptions) GetDebugLogFilter() *string {
	return o.debugLogFilter
}

// Java: CompilerOptions.shouldSerializeExtraDebugInfo()
func (o *CompilerOptions) ShouldSerializeExtraDebugInfo() bool {
	// NOTE: this is tied to debug logging (via checking getDebugLogDirectory()) for convenience,
	// and it's possible to split up the two options if this begins to trigger memory issues.
	return o.serializeExtraDebugInfo || o.GetDebugLogDirectory() != nil
}

// Java: CompilerOptions.setSerializeExtraDebugInfo(boolean)
func (o *CompilerOptions) SetSerializeExtraDebugInfo(serializeExtraDebugInfo bool) {
	o.serializeExtraDebugInfo = serializeExtraDebugInfo
}

// Java: CompilerOptions.setQuoteKeywordProperties(boolean)
func (o *CompilerOptions) SetQuoteKeywordProperties(quoteKeywordProperties bool) {
	o.quoteKeywordProperties = quoteKeywordProperties
}

// Java: CompilerOptions.shouldQuoteKeywordProperties()
func (o *CompilerOptions) ShouldQuoteKeywordProperties() bool {
	// Never quote properties in .i.js files
	if o.incrementalCheckMode == CompilerOptions_IncrementalCheckMode_GENERATE_IJS {
		return false
	}
	return o.quoteKeywordProperties || parser.FeatureSet_ES3.Contains(o.GetOutputFeatureSet())
}

// Java: CompilerOptions.setErrorFormat(ErrorFormat)
func (o *CompilerOptions) SetErrorFormat(errorFormat ErrorFormat) {
	o.errorFormat = errorFormat
}

// Java: CompilerOptions.getErrorFormat()
func (o *CompilerOptions) GetErrorFormat() ErrorFormat {
	return o.errorFormat
}

// Java: CompilerOptions.setWarningsGuard(ComposeWarningsGuard)
func (o *CompilerOptions) SetWarningsGuard(warningsGuard *ComposeWarningsGuard) {
	o.warningsGuard = warningsGuard
}

// Java: CompilerOptions.setLineLengthThreshold(int)
func (o *CompilerOptions) SetLineLengthThreshold(lineLengthThreshold int) {
	o.lineLengthThreshold = lineLengthThreshold
}

// Java: CompilerOptions.getLineLengthThreshold()
func (o *CompilerOptions) GetLineLengthThreshold() int {
	return o.lineLengthThreshold
}

// Java: CompilerOptions.setUseOriginalNamesInOutput(boolean)
func (o *CompilerOptions) SetUseOriginalNamesInOutput(useOriginalNamesInOutput bool) {
	o.useOriginalNamesInOutput = useOriginalNamesInOutput
}

// Java: CompilerOptions.getUseOriginalNamesInOutput()
func (o *CompilerOptions) GetUseOriginalNamesInOutput() bool {
	return o.useOriginalNamesInOutput
}

// Java: CompilerOptions.setExternExportsPath(String)
func (o *CompilerOptions) SetExternExportsPath(externExportsPath *string) {
	o.externExportsPath = externExportsPath
}

// Java: CompilerOptions.getExternExportsPath()
func (o *CompilerOptions) GetExternExportsPath() *string {
	return o.externExportsPath
}

// Java: CompilerOptions.setSourceMapOutputPath(String)
func (o *CompilerOptions) SetSourceMapOutputPath(sourceMapOutputPath *string) {
	o.sourceMapOutputPath = sourceMapOutputPath
}

// Java: CompilerOptions.setApplyInputSourceMaps(boolean)
func (o *CompilerOptions) SetApplyInputSourceMaps(applyInputSourceMaps bool) {
	o.applyInputSourceMaps = applyInputSourceMaps
}

// Java: CompilerOptions.getApplyInputSourceMaps()
func (o *CompilerOptions) GetApplyInputSourceMaps() bool {
	return o.applyInputSourceMaps
}

// Java: CompilerOptions.setResolveSourceMapAnnotations(boolean)
func (o *CompilerOptions) SetResolveSourceMapAnnotations(resolveSourceMapAnnotations bool) {
	o.resolveSourceMapAnnotations = resolveSourceMapAnnotations
}

// Java: CompilerOptions.getResolveSourceMapAnnotations()
func (o *CompilerOptions) GetResolveSourceMapAnnotations() bool {
	return o.resolveSourceMapAnnotations
}

// Java: CompilerOptions.setSourceMapIncludeSourcesContent(boolean)
func (o *CompilerOptions) SetSourceMapIncludeSourcesContent(sourceMapIncludeSourcesContent bool) {
	o.sourceMapIncludeSourcesContent = sourceMapIncludeSourcesContent
}

// Java: CompilerOptions.getSourceMapIncludeSourcesContent()
func (o *CompilerOptions) GetSourceMapIncludeSourcesContent() bool {
	return o.sourceMapIncludeSourcesContent
}

// Java: CompilerOptions.setParseInlineSourceMaps(boolean)
func (o *CompilerOptions) SetParseInlineSourceMaps(parseInlineSourceMaps bool) {
	o.parseInlineSourceMaps = parseInlineSourceMaps
}

// Java: CompilerOptions.getParseInlineSourceMaps()
func (o *CompilerOptions) GetParseInlineSourceMaps() bool {
	return o.parseInlineSourceMaps
}

// Java: CompilerOptions.setSourceMapDetailLevel(SourceMap.DetailLevel)
func (o *CompilerOptions) SetSourceMapDetailLevel(sourceMapDetailLevel SourceMap_DetailLevel) {
	o.sourceMapDetailLevel = sourceMapDetailLevel
}

// Java: CompilerOptions.getSourceMapDetailLevel()
func (o *CompilerOptions) GetSourceMapDetailLevel() SourceMap_DetailLevel {
	return o.sourceMapDetailLevel
}

// Java: CompilerOptions.setSourceMapFormat(SourceMap.Format)
func (o *CompilerOptions) SetSourceMapFormat(sourceMapFormat SourceMap_Format) {
	o.sourceMapFormat = sourceMapFormat
}

// Java: CompilerOptions.getSourceMapFormat()
func (o *CompilerOptions) GetSourceMapFormat() SourceMap_Format {
	return o.sourceMapFormat
}

// Java: CompilerOptions.setSourceMapLocationMappings(List<? extends SourceMap.LocationMapping>)
func (o *CompilerOptions) SetSourceMapLocationMappings(sourceMapLocationMappings []SourceMap_LocationMapping) {
	o.sourceMapLocationMappings = append([]SourceMap_LocationMapping{}, sourceMapLocationMappings...)
}

// Java: CompilerOptions.getSourceMapLocationMappings()
func (o *CompilerOptions) GetSourceMapLocationMappings() []SourceMap_LocationMapping {
	return o.sourceMapLocationMappings
}

// Rewrites CommonJS modules so that modules can be concatenated together, by renaming all globals
// to avoid conflicting with other modules.
//
// Java: CompilerOptions.setProcessCommonJSModules(boolean)
func (o *CompilerOptions) SetProcessCommonJSModules(processCommonJSModules bool) {
	o.processCommonJSModules = processCommonJSModules
}

// Java: CompilerOptions.getProcessCommonJSModules()
func (o *CompilerOptions) GetProcessCommonJSModules() bool {
	return o.processCommonJSModules
}

// CompilerOptions_Es6ModuleTranspilation is how ES modules should be transformed.
type CompilerOptions_Es6ModuleTranspilation int32

const (
	// Do not touch any Es6 module feature.
	CompilerOptions_Es6ModuleTranspilation_NONE CompilerOptions_Es6ModuleTranspilation = iota
	// Rewrite import paths to resolved, relative paths only.
	CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS
	// Rewrite to common js like modules for bundling.
	CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES
	// Compile ES modules.
	CompilerOptions_Es6ModuleTranspilation_COMPILE
)

var compilerOptions_Es6ModuleTranspilation_names = [...]string{"NONE", "RELATIVIZE_IMPORT_PATHS", "TO_COMMON_JS_LIKE_MODULES", "COMPILE"}

// CompilerOptions_Es6ModuleTranspilation_values is Es6ModuleTranspilation.values().
func CompilerOptions_Es6ModuleTranspilation_values() []CompilerOptions_Es6ModuleTranspilation {
	return []CompilerOptions_Es6ModuleTranspilation{CompilerOptions_Es6ModuleTranspilation_NONE, CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS, CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES, CompilerOptions_Es6ModuleTranspilation_COMPILE}
}

// String is Enum.name().
func (e CompilerOptions_Es6ModuleTranspilation) String() string {
	return compilerOptions_Es6ModuleTranspilation_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_Es6ModuleTranspilation) Ordinal() int { return int(e) }

// Java: CompilerOptions.setEs6ModuleTranspilation(Es6ModuleTranspilation)
func (o *CompilerOptions) SetEs6ModuleTranspilation(value CompilerOptions_Es6ModuleTranspilation) {
	o.es6ModuleTranspilation = value
}

// Java: CompilerOptions.getEs6ModuleTranspilation()
func (o *CompilerOptions) GetEs6ModuleTranspilation() CompilerOptions_Es6ModuleTranspilation {
	return o.es6ModuleTranspilation
}

// Sets a path prefix for CommonJS modules (maps to #setModuleRoots(List)).
//
// Java: CompilerOptions.setCommonJSModulePathPrefix(String)
func (o *CompilerOptions) SetCommonJSModulePathPrefix(commonJSModulePathPrefix string) {
	o.SetModuleRoots([]string{commonJSModulePathPrefix})
}

// Sets the module roots.
//
// Java: CompilerOptions.setModuleRoots(List<String>)
func (o *CompilerOptions) SetModuleRoots(moduleRoots []string) {
	o.moduleRoots = append([]string{}, moduleRoots...)
}

// Java: CompilerOptions.getModuleRoots()
func (o *CompilerOptions) GetModuleRoots() []string {
	return o.moduleRoots
}

// Sets whether to rewrite polyfills.
//
// Java: CompilerOptions.setRewritePolyfills(boolean)
func (o *CompilerOptions) SetRewritePolyfills(rewritePolyfills bool) {
	o.rewritePolyfills = rewritePolyfills
}

// Java: CompilerOptions.getRewritePolyfills()
func (o *CompilerOptions) GetRewritePolyfills() bool {
	return o.rewritePolyfills
}

// Sets whether to isolate polyfills from the global scope.
//
// Java: CompilerOptions.setIsolatePolyfills(boolean)
func (o *CompilerOptions) SetIsolatePolyfills(isolatePolyfills bool) {
	o.isolatePolyfills = isolatePolyfills
	if o.isolatePolyfills {
		o.SetDefineToBooleanLiteral("$jscomp.ISOLATE_POLYFILLS", isolatePolyfills)
	}
}

// Java: CompilerOptions.getIsolatePolyfills()
func (o *CompilerOptions) GetIsolatePolyfills() bool {
	return o.isolatePolyfills
}

// Java: CompilerOptions.setInjectPolyfillsNewerThan(LanguageMode)
func (o *CompilerOptions) SetInjectPolyfillsNewerThan(injectPolyfillsNewerThan *CompilerOptions_LanguageMode) {
	o.injectPolyfillsNewerThan = injectPolyfillsNewerThan
}

// Java: CompilerOptions.getInjectPolyfillsNewerThan()
func (o *CompilerOptions) GetInjectPolyfillsNewerThan() *CompilerOptions_LanguageMode {
	return o.injectPolyfillsNewerThan
}

// Sets whether to transpile async functions and generators for AsyncContext.
//
// Java: CompilerOptions.setInstrumentAsyncContext(boolean)
func (o *CompilerOptions) SetInstrumentAsyncContext(instrumentAsyncContext bool) {
	o.instrumentAsyncContext = instrumentAsyncContext
	o.SetDefineToBooleanLiteral("$jscomp.INSTRUMENT_ASYNC_CONTEXT", instrumentAsyncContext)
}

// Java: CompilerOptions.getInstrumentAsyncContext()
func (o *CompilerOptions) GetInstrumentAsyncContext() bool {
	return o.instrumentAsyncContext
}

// CompilerOptions_Es6SubclassTranspilation configures the transpiled output for ES6 classes when targeting ES5.
//
// This is needed for transpilation of ES6 classes that extend some unknown
// superclass, in case that superclass is still a native ES6 class at runtime
// (not transpiled to ES5).
//
// The default is CONCISE_UNSAFE, which assumes that if a superclass in an
// extends clause cannot be statically analyzed, it's safe to assume it is an
// ES5 class.
//
// The "safe" behavior relies on $jscomp.construct, which is our wrapper for
// Reflect.construct. This is correct, but has the disadvantage of pulling in
// more helper utilities into the compiled output and of the transpiled output
// being slightly larger.
//
// TODO: b/36789413 - always enable the "safe" transpilation.
type CompilerOptions_Es6SubclassTranspilation int32

const (
	CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE CompilerOptions_Es6SubclassTranspilation = iota
	CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT
)

var compilerOptions_Es6SubclassTranspilation_names = [...]string{"CONCISE_UNSAFE", "SAFE_REFLECT_CONSTRUCT"}

// CompilerOptions_Es6SubclassTranspilation_values is Es6SubclassTranspilation.values().
func CompilerOptions_Es6SubclassTranspilation_values() []CompilerOptions_Es6SubclassTranspilation {
	return []CompilerOptions_Es6SubclassTranspilation{CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE, CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT}
}

// String is Enum.name().
func (e CompilerOptions_Es6SubclassTranspilation) String() string {
	return compilerOptions_Es6SubclassTranspilation_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_Es6SubclassTranspilation) Ordinal() int { return int(e) }

// Java: CompilerOptions.setEs6SubclassTranspilation(Es6SubclassTranspilation)
func (o *CompilerOptions) SetEs6SubclassTranspilation(es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) {
	o.es6SubclassTranspilation = es6SubclassTranspilation
}

// Java: CompilerOptions.getEs6SubclassTranspilation()
func (o *CompilerOptions) GetEs6SubclassTranspilation() CompilerOptions_Es6SubclassTranspilation {
	return o.es6SubclassTranspilation
}

// Sets list of libraries to always inject, even if not needed.
//
// Java: CompilerOptions.setForceLibraryInjection(Iterable<String>)
func (o *CompilerOptions) SetForceLibraryInjection(libraries []string) {
	o.forceLibraryInjection = append([]string{}, libraries...)
}

// Java: CompilerOptions.getForceLibraryInjectionList()
func (o *CompilerOptions) GetForceLibraryInjectionList() []string {
	return o.forceLibraryInjection
}

// Sets the set of libraries to never inject, even if required.
//
// Java: CompilerOptions.setPreventLibraryInjection(boolean)
func (o *CompilerOptions) SetPreventLibraryInjection(preventLibraryInjection bool) {
	if preventLibraryInjection {
		o.runtimeLibraryMode = RuntimeJsLibManager_RuntimeLibraryMode_NO_OP
	} else {
		o.runtimeLibraryMode = RuntimeJsLibManager_RuntimeLibraryMode_INJECT
	}
}

// Java: CompilerOptions.setRuntimeLibraryMode(RuntimeLibraryMode)
func (o *CompilerOptions) SetRuntimeLibraryMode(runtimeLibraryMode RuntimeJsLibManager_RuntimeLibraryMode) {
	o.runtimeLibraryMode = runtimeLibraryMode
}

// Java: CompilerOptions.getRuntimeLibraryMode()
func (o *CompilerOptions) GetRuntimeLibraryMode() RuntimeJsLibManager_RuntimeLibraryMode {
	return o.runtimeLibraryMode
}

// Java: CompilerOptions.setUnusedImportsToRemove(ImmutableSet<String>)
func (o *CompilerOptions) SetUnusedImportsToRemove(unusedImportsToRemove *jutil.LinkedHashSet[string]) {
	o.unusedImportsToRemove = unusedImportsToRemove
}

// Java: CompilerOptions.getUnusedImportsToRemove()
func (o *CompilerOptions) GetUnusedImportsToRemove() *jutil.LinkedHashSet[string] {
	return o.unusedImportsToRemove
}

// Java: CompilerOptions.setInstrumentForCoverageOption(InstrumentOption)
func (o *CompilerOptions) SetInstrumentForCoverageOption(instrumentForCoverageOption CompilerOptions_InstrumentOption) {
	o.instrumentForCoverageOption = instrumentForCoverageOption
}

// Java: CompilerOptions.getInstrumentForCoverageOption()
func (o *CompilerOptions) GetInstrumentForCoverageOption() CompilerOptions_InstrumentOption {
	return o.instrumentForCoverageOption
}

// Sets the name for the global array which is used by PRODUCTION instrumentation. The array is
// declared during the instrumentation pass with the name provided through this setter.
//
// Java: CompilerOptions.setProductionInstrumentationArrayName(String)
func (o *CompilerOptions) SetProductionInstrumentationArrayName(productionInstrumentationArrayName string) {
	o.productionInstrumentationArrayName = productionInstrumentationArrayName
}

// Java: CompilerOptions.getProductionInstrumentationArrayName()
func (o *CompilerOptions) GetProductionInstrumentationArrayName() string {
	return o.productionInstrumentationArrayName
}

// Java: CompilerOptions.getConformanceConfigs()
func (o *CompilerOptions) GetConformanceConfigs() []*ConformanceConfig {
	return o.conformanceConfigs
}

// Both enable and configure conformance checks, if non-null.
//
// @deprecated See go/binary-level-conformance-deprecated.
//
// Java: CompilerOptions.setConformanceConfig(ConformanceConfig)
func (o *CompilerOptions) SetConformanceConfig(conformanceConfig *ConformanceConfig) {
	o.SetConformanceConfigs([]*ConformanceConfig{jutil.CheckNotNull(conformanceConfig)})
}

// Both enable and configure conformance checks, if non-null.
//
// @deprecated See go/binary-level-conformance-deprecated.
//
// Java: CompilerOptions.setConformanceConfigs(List<ConformanceConfig>)
func (o *CompilerOptions) SetConformanceConfigs(configs []*ConformanceConfig) {
	for _, c := range configs {
		jutil.CheckNotNull(c) // ImmutableList.copyOf rejects null elements.
	}
	o.conformanceConfigs = append([]*ConformanceConfig{}, configs...)
}

// Java: CompilerOptions.clearConformanceConfigs()
func (o *CompilerOptions) ClearConformanceConfigs() {
	o.conformanceConfigs = []*ConformanceConfig{}
}

// Whether the output should contain a 'use strict' directive.
//
// Java: CompilerOptions.shouldEmitUseStrict()
func (o *CompilerOptions) ShouldEmitUseStrict() bool {
	// Fall back to the language in's strictness if there is no output language explicitly set
	// for backwards compatibility.
	if o.emitUseStrict != nil {
		return *o.emitUseStrict
	}
	if o.languageOutIsDefaultStrict != nil {
		return *o.languageOutIsDefaultStrict
	}
	return o.languageIn.IsDefaultStrict()
}

// Java: CompilerOptions.setEmitUseStrict(boolean)
func (o *CompilerOptions) SetEmitUseStrict(emitUseStrict bool) *CompilerOptions {
	o.emitUseStrict = &emitUseStrict
	return o
}

// Java: CompilerOptions.getModuleResolutionMode()
func (o *CompilerOptions) GetModuleResolutionMode() ModuleLoader_ResolutionMode {
	return o.moduleResolutionMode
}

// Java: CompilerOptions.setModuleResolutionMode(ResolutionMode)
func (o *CompilerOptions) SetModuleResolutionMode(moduleResolutionMode ModuleLoader_ResolutionMode) {
	o.moduleResolutionMode = moduleResolutionMode
}

// Java: CompilerOptions.getBrowserResolverPrefixReplacements()
func (o *CompilerOptions) GetBrowserResolverPrefixReplacements() *jutil.LinkedHashMap[string, string] {
	return o.browserResolverPrefixReplacements
}

// Java: CompilerOptions.setBrowserResolverPrefixReplacements(ImmutableMap<String, String>)
func (o *CompilerOptions) SetBrowserResolverPrefixReplacements(browserResolverPrefixReplacements *jutil.LinkedHashMap[string, string]) {
	o.browserResolverPrefixReplacements = browserResolverPrefixReplacements
}

// Java: CompilerOptions.setPathEscaper(ModuleLoader.PathEscaper)
func (o *CompilerOptions) SetPathEscaper(pathEscaper ModuleLoader_PathEscaper) {
	o.pathEscaper = pathEscaper
}

// Java: CompilerOptions.getPathEscaper()
func (o *CompilerOptions) GetPathEscaper() ModuleLoader_PathEscaper {
	return o.pathEscaper
}

// Java: CompilerOptions.getPackageJsonEntryNames()
func (o *CompilerOptions) GetPackageJsonEntryNames() []string {
	return o.packageJsonEntryNames
}

// Java: CompilerOptions.setPackageJsonEntryNames(List<String>)
func (o *CompilerOptions) SetPackageJsonEntryNames(names []string) {
	o.packageJsonEntryNames = names
}

// Java: CompilerOptions.setUseSizeHeuristicToStopOptimizationLoop(boolean)
func (o *CompilerOptions) SetUseSizeHeuristicToStopOptimizationLoop(mayStopEarly bool) {
	o.useSizeHeuristicToStopOptimizationLoop = mayStopEarly
}

// Java: CompilerOptions.shouldUseSizeHeuristicToStopOptimizationLoop()
func (o *CompilerOptions) ShouldUseSizeHeuristicToStopOptimizationLoop() bool {
	return o.useSizeHeuristicToStopOptimizationLoop
}

// Java: CompilerOptions.setMaxOptimizationLoopIterations(int)
func (o *CompilerOptions) SetMaxOptimizationLoopIterations(maxIterations int) {
	o.optimizationLoopMaxIterations = maxIterations
}

// Java: CompilerOptions.getMaxOptimizationLoopIterations()
func (o *CompilerOptions) GetMaxOptimizationLoopIterations() int {
	return o.optimizationLoopMaxIterations
}

// Java: CompilerOptions.getChunkOutputType()
func (o *CompilerOptions) GetChunkOutputType() CompilerOptions_ChunkOutputType {
	return o.chunkOutputType
}

// Java: CompilerOptions.setChunkOutputType(ChunkOutputType)
func (o *CompilerOptions) SetChunkOutputType(chunkOutputType CompilerOptions_ChunkOutputType) {
	o.chunkOutputType = chunkOutputType
}

// Java: CompilerOptions.setStrictMessageReplacement(boolean)
func (o *CompilerOptions) SetStrictMessageReplacement(strictMessageReplacement bool) {
	o.strictMessageReplacement = strictMessageReplacement
}

// Java: CompilerOptions.getStrictMessageReplacement()
func (o *CompilerOptions) GetStrictMessageReplacement() bool {
	return o.strictMessageReplacement
}

// Java: CompilerOptions.toString()
func (o *CompilerOptions) String() string {
	h := newCompilerOptions_toStringHelper("CompilerOptions")
	h.add("aliasStringsMode", o.GetAliasStringsMode())
	h.addUnconditional("ambiguateProperties", o.ambiguateProperties)
	h.addUnconditional("angularPass", o.angularPass)
	h.addUnconditional("assumeClosuresOnlyCaptureReferences", o.assumeClosuresOnlyCaptureReferences)
	h.addUnconditional("assumeGettersArePure", o.assumeGettersArePure)
	h.addUnconditional("assumePropertiesAreStaticallyAnalyzable", o.assumePropertiesAreStaticallyAnalyzable)
	h.addUnconditional("assumeStrictThis", o.AssumeStrictThis())
	h.add("browserResolverPrefixReplacements", o.browserResolverPrefixReplacements)
	h.addUnconditional("checkDeterminism", o.GetCheckDeterminism())
	h.addUnconditional("checkSuspiciousCode", o.checkSuspiciousCode)
	h.addUnconditional("checkSymbols", o.checkSymbols)
	h.addUnconditional("checkTypes", o.checkTypes)
	h.addUnconditional("checksOnly", o.checksOnly)
	h.add("chunksToPrintAfterEachPassRegexList", o.chunksToPrintAfterEachPassRegexList)
	h.addUnconditional("closurePass", o.closurePass)
	h.addUnconditional("coalesceVariableNames", o.coalesceVariableNames)
	h.add("codingConvention", o.GetCodingConvention())
	h.addUnconditional("collapseAnonymousFunctions", o.collapseAnonymousFunctions)
	h.addUnconditional("collapseObjectLiterals", o.collapseObjectLiterals)
	h.add("collapseProperties", o.collapsePropertiesLevel)
	h.addUnconditional("collapseVariableDeclarations", o.collapseVariableDeclarations)
	h.addUnconditional("colorizeErrorOutput", o.ShouldColorizeErrorOutput())
	h.addUnconditional("computeFunctionSideEffects", o.computeFunctionSideEffects)
	h.add("conformanceConfigs", o.GetConformanceConfigs())
	h.add("conformanceRemoveRegexFromPath", compilerOptions_optionalOf(o.conformanceRemoveRegexFromPath))
	h.add("conformanceReportingMode", o.conformanceReportingMode)
	h.addUnconditional("continueAfterErrors", o.CanContinueAfterErrors())
	h.addUnconditional("convertToDottedProperties", o.convertToDottedProperties)
	h.addUnconditional("crossChunkCodeMotion", o.crossChunkCodeMotion)
	h.addUnconditional("crossChunkCodeMotionNoStubMethods", o.crossChunkCodeMotionNoStubMethods)
	h.addUnconditional("crossChunkMethodMotion", o.crossChunkMethodMotion)
	h.add("cssRenamingMap", o.cssRenamingMap)
	h.add("cssRenamingSkiplist", o.cssRenamingSkiplist)
	h.add("customPasses", o.customPasses)
	h.addUnconditional("deadAssignmentElimination", o.deadAssignmentElimination)
	h.add("debugLogDirectory", o.debugLogDirectory)
	h.add("defineReplacements", o.GetDefineReplacements())
	h.add("dependencyOptions", o.GetDependencyOptions())
	h.add("devMode", o.devMode)
	h.addUnconditional("devirtualizeMethods", o.devirtualizeMethods)
	h.addUnconditional("disambiguateProperties", o.disambiguateProperties)
	h.add("emitUseStrict", compilerOptions_optionalOf(o.emitUseStrict))
	h.addUnconditional("enableModuleRewriting", o.enableModuleRewriting)
	h.add("environment", o.GetEnvironment())
	h.add("errorFormat", o.errorFormat)
	h.add("errorHandler", o.errorHandler)
	h.add("es6ModuleTranspilation", o.es6ModuleTranspilation)
	h.add("es6SubclassTranspilation", o.es6SubclassTranspilation)
	h.addUnconditional("exportLocalPropertyDefinitions", o.exportLocalPropertyDefinitions)
	h.addUnconditional("exportTestFunctions", o.exportTestFunctions)
	h.add("externExportsPath", o.externExportsPath)
	h.add("extraAnnotationNames", o.extraAnnotationNames)
	h.add("extractPrototypeMemberDeclarations", o.extractPrototypeMemberDeclarations)
	h.add("filesToPrintAfterEachPassRegexList", o.filesToPrintAfterEachPassRegexList)
	h.addUnconditional("flowSensitiveInlineVariables", o.flowSensitiveInlineVariables)
	h.addUnconditional("foldConstants", o.foldConstants)
	h.add("forceLibraryInjection", o.forceLibraryInjection)
	h.addUnconditional("gatherCssNames", o.gatherCssNames)
	h.addUnconditional("generateExports", o.generateExports)
	h.addUnconditional("generatePseudoNames", o.generatePseudoNames)
	h.addUnconditional("generateTypedExterns", o.ShouldGenerateTypedExterns())
	h.add("idGenerators", o.idGenerators)
	h.add("idGeneratorsMapSerialized", o.idGeneratorsMapSerialized)
	h.add("incrementalCheckMode", o.incrementalCheckMode)
	h.addUnconditional("inferConsts", o.inferConsts)
	h.addUnconditional("inferTypes", o.inferTypes)
	h.addUnconditional("inlineConstantVars", o.inlineConstantVars)
	h.add("inlineFunctionsLevel", o.inlineFunctionsLevel)
	h.addUnconditional("inlineGetters", o.inlineGetters)
	h.addUnconditional("inlineLocalVariables", o.inlineLocalVariables)
	h.addUnconditional("inlineProperties", o.inlineProperties)
	h.addUnconditional("inlineVariables", o.inlineVariables)
	h.add("inputDelimiter", o.inputDelimiter)
	h.add("inputPropertyMap", o.inputPropertyMap)
	h.add("inputSourceMaps", o.inputSourceMaps)
	h.add("inputVariableMap", o.inputVariableMap)
	h.addUnconditional("instrumentForCoverageOnly", o.instrumentForCoverageOnly)
	h.add("instrumentForCoverageOption", o.instrumentForCoverageOption.String())
	h.addUnconditional("isolatePolyfills", o.isolatePolyfills)
	h.addUnconditional("j2clMinifierEnabled", o.j2clMinifierEnabled)
	h.add("j2clMinifierPruningManifest", o.j2clMinifierPruningManifest)
	h.add("j2clPassMode", o.j2clPassMode)
	h.addUnconditional("labelRenaming", o.labelRenaming)
	h.add("languageIn", o.GetLanguageIn())
	h.add("languageOutIsDefaultStrict", compilerOptions_optionalOf(o.languageOutIsDefaultStrict))
	h.addUnconditional("lineBreak", o.lineBreak)
	h.addUnconditional("lineLengthThreshold", o.lineLengthThreshold)
	h.add("locale", o.locale)
	h.addUnconditional("markAsCompiled", o.markAsCompiled)
	h.addUnconditional("maxFunctionSizeAfterInlining", o.maxFunctionSizeAfterInlining)
	h.add("messageBundle", o.messageBundle)
	h.add("moduleRoots", o.moduleRoots)
	h.add("nameGenerator", o.nameGenerator)
	h.addUnconditional("numParallelThreads", o.numParallelThreads)
	h.addUnconditional("optimizeCalls", o.optimizeCalls)
	h.addUnconditional("optimizeESClassConstructors", o.optimizeESClassConstructors)
	h.add("outputCharset", o.outputCharset)
	h.add("outputFeatureSet", compilerOptions_optionalOf(o.outputFeatureSet))
	h.add("outputJs", o.outputJs)
	h.addUnconditional("outputJsStringUsage", o.outputJsStringUsage)
	h.addUnconditional("parentChunkCanSeeSymbolsDeclaredInChildren", o.parentChunkCanSeeSymbolsDeclaredInChildren)
	h.add("parseJsDocDocumentation", o.IsParseJsDocDocumentation())
	h.add("pathEscaper", o.pathEscaper)
	h.addUnconditional("polymerPass", o.polymerPass)
	h.addUnconditional("preferSingleQuotes", o.preferSingleQuotes)
	h.addUnconditional("preferStableNames", o.preferStableNames)
	h.addUnconditional("preserveDetailedSourceInfo", o.PreservesDetailedSourceInfo())
	h.addUnconditional("preserveGoogProvidesAndRequires", o.preserveClosurePrimitives)
	h.addUnconditional("preserveNonJSDocComments", o.GetPreserveNonJSDocComments())
	h.addUnconditional("preserveTypeAnnotations", o.preserveTypeAnnotations)
	h.addUnconditional("prettyPrint", o.prettyPrint)
	h.addUnconditional("printConfig", o.printConfig)
	h.addUnconditional("printInputDelimiter", o.printInputDelimiter)
	h.addUnconditional("printSourceAfterEachPass", o.printSourceAfterEachPass)
	h.addUnconditional("processCommonJSModules", o.processCommonJSModules)
	h.add("productionInstrumentationArrayName", o.productionInstrumentationArrayName)
	h.add("propertiesThatMustDisambiguate", o.propertiesThatMustDisambiguate)
	h.add("propertyRenaming", o.propertyRenaming)
	h.addUnconditional("propertyRenamingOnlyCompilationMode", o.propertyRenamingOnlyCompilationMode)
	h.addUnconditional("protectHiddenSideEffects", o.protectHiddenSideEffects)
	h.add("qnameUsesToPrintAfterEachPassRegexList", o.qnameUsesToPrintAfterEachPassList)
	h.addUnconditional("quoteKeywordProperties", o.quoteKeywordProperties)
	h.addUnconditional("removeAbstractMethods", o.removeAbstractMethods)
	h.addUnconditional("removeClosureAsserts", o.removeClosureAsserts)
	h.addUnconditional("removeJ2clAsserts", o.removeJ2clAsserts)
	h.addUnconditional("removeUnusedClassProperties", o.removeUnusedClassProperties)
	h.addUnconditional("removeUnusedLocalVars", o.removeUnusedLocalVars)
	h.addUnconditional("removeUnusedPrototypeProperties", o.removeUnusedPrototypeProperties)
	h.addUnconditional("removeUnusedVars", o.removeUnusedVars)
	h.add("renamePrefix", o.renamePrefix)
	h.add("renamePrefixNamespace", o.renamePrefixNamespace)
	h.addUnconditional("renamePrefixNamespaceAssumeCrossChunkNames", o.renamePrefixNamespaceAssumeCrossChunkNames)
	h.addUnconditional("replaceIdGenerators", o.replaceIdGenerators)
	h.addUnconditional("replaceMessagesWithChromeI18n", o.replaceMessagesWithChromeI18n)
	h.add("replaceStringsFunctionDescriptions", o.replaceStringsFunctionDescriptions)
	h.add("replaceStringsPlaceholderToken", o.replaceStringsPlaceholderToken)
	h.addUnconditional("reserveRawExports", o.reserveRawExports)
	h.addUnconditional("rewriteFunctionExpressions", o.rewriteFunctionExpressions)
	h.addUnconditional("rewriteGlobalDeclarationsForTryCatchWrapping", o.rewriteGlobalDeclarationsForTryCatchWrapping)
	h.addUnconditional("rewriteModulesBeforeTypechecking", o.rewriteModulesBeforeTypechecking)
	h.addUnconditional("rewritePolyfills", o.rewritePolyfills)
	h.add("runtimeLibraryMode", o.runtimeLibraryMode)
	h.addUnconditional("skipNonTranspilationPasses", o.skipNonTranspilationPasses)
	h.addUnconditional("smartNameRemoval", o.smartNameRemoval)
	h.add("sourceMapDetailLevel", o.sourceMapDetailLevel)
	h.add("sourceMapFormat", o.sourceMapFormat)
	h.add("sourceMapLocationMappings", o.sourceMapLocationMappings)
	h.add("sourceMapOutputPath", o.sourceMapOutputPath)
	h.addUnconditional("strictMessageReplacement", o.strictMessageReplacement)
	h.add("stripNamePrefixes", o.stripNamePrefixes)
	h.add("stripNameSuffixes", o.stripNameSuffixes)
	h.add("stripTypes", o.stripTypes)
	h.addUnconditional("summaryDetailLevel", o.summaryDetailLevel)
	h.add("syntheticBlockEndMarker", o.syntheticBlockEndMarker)
	h.add("syntheticBlockStartMarker", o.syntheticBlockStartMarker)
	h.add("tcProjectId", o.tcProjectId)
	h.add("tracer", o.tracer)
	h.addUnconditional("trustedStrings", o.trustedStrings)
	h.add("tweakProcessing", o.GetTweakProcessing())
	h.add("unusedImportsToRemove", o.unusedImportsToRemove)
	h.addUnconditional("useTypesForLocalOptimization", o.useTypesForLocalOptimization)
	h.add("variableRenaming", o.variableRenaming)
	h.add("warningsGuard", o.GetWarningsGuard())
	h.addUnconditional("wrapGoogModulesForWhitespaceOnly", o.wrapGoogModulesForWhitespaceOnly)
	return h.String()
}

// CompilerOptions_InstrumentOption is an option to determine what level of code instrumentation is performed, if any
type CompilerOptions_InstrumentOption int32

const (
	// No coverage instrumentation is performed
	CompilerOptions_InstrumentOption_NONE CompilerOptions_InstrumentOption = iota
	// Collect coverage for every executable statement.
	CompilerOptions_InstrumentOption_LINE_ONLY
	// Collect coverage for control-flow branches.
	CompilerOptions_InstrumentOption_BRANCH_ONLY
	// Collect coverage for functions where code is compiled for production.
	CompilerOptions_InstrumentOption_PRODUCTION
)

var compilerOptions_InstrumentOption_names = [...]string{"NONE", "LINE_ONLY", "BRANCH_ONLY", "PRODUCTION"}

// CompilerOptions_InstrumentOption_values is InstrumentOption.values().
func CompilerOptions_InstrumentOption_values() []CompilerOptions_InstrumentOption {
	return []CompilerOptions_InstrumentOption{CompilerOptions_InstrumentOption_NONE, CompilerOptions_InstrumentOption_LINE_ONLY, CompilerOptions_InstrumentOption_BRANCH_ONLY, CompilerOptions_InstrumentOption_PRODUCTION}
}

// String is Enum.name().
func (e CompilerOptions_InstrumentOption) String() string {
	return compilerOptions_InstrumentOption_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_InstrumentOption) Ordinal() int { return int(e) }

var compilerOptions_InstrumentOption_ptrs = func() []CompilerOptions_InstrumentOption {
	v := CompilerOptions_InstrumentOption_values()
	return v
}()

// Ptr returns the canonical pointer to e (Go-only: language-adapter — a
// nullable Java InstrumentOption reference is a *CompilerOptions_InstrumentOption; nil is null).
func (e CompilerOptions_InstrumentOption) Ptr() *CompilerOptions_InstrumentOption {
	return &compilerOptions_InstrumentOption_ptrs[e]
}

// Java: CompilerOptions.InstrumentOption.fromString(String)
func CompilerOptions_InstrumentOption_FromString(value *string) *CompilerOptions_InstrumentOption {
	if value == nil {
		return nil
	}
	switch *value {
	case "NONE":
		return CompilerOptions_InstrumentOption_NONE.Ptr()
	case "LINE":
		return CompilerOptions_InstrumentOption_LINE_ONLY.Ptr()
	case "BRANCH":
		return CompilerOptions_InstrumentOption_BRANCH_ONLY.Ptr()
	case "PRODUCTION":
		return CompilerOptions_InstrumentOption_PRODUCTION.Ptr()
	default:
		return nil
	}
}

// CompilerOptions_ChunkOutputType is the format for compiler output when multiple chunks are requested.
type CompilerOptions_ChunkOutputType int32

const (
	// Original format - each chunk is a simple script
	CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE CompilerOptions_ChunkOutputType = iota
	// Each chunk is an ES module with dependencies and symbol access set up using
	// import and export statements.
	CompilerOptions_ChunkOutputType_ES_MODULES
)

var compilerOptions_ChunkOutputType_names = [...]string{"GLOBAL_NAMESPACE", "ES_MODULES"}

// CompilerOptions_ChunkOutputType_values is ChunkOutputType.values().
func CompilerOptions_ChunkOutputType_values() []CompilerOptions_ChunkOutputType {
	return []CompilerOptions_ChunkOutputType{CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE, CompilerOptions_ChunkOutputType_ES_MODULES}
}

// String is Enum.name().
func (e CompilerOptions_ChunkOutputType) String() string {
	return compilerOptions_ChunkOutputType_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_ChunkOutputType) Ordinal() int { return int(e) }

// CompilerOptions_LanguageMode is a language mode applies to the whole compilation job. As a result, the
// compiler does not support mixed strict and non-strict in the same compilation
// job. Therefore, the 'use strict' directive is ignored when the language mode
// is not strict.
type CompilerOptions_LanguageMode int32

const (
	// 90's JavaScript
	CompilerOptions_LanguageMode_ECMASCRIPT3 CompilerOptions_LanguageMode = iota
	// Traditional JavaScript
	CompilerOptions_LanguageMode_ECMASCRIPT5
	// Nitpicky, traditional JavaScript
	CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT
	// ECMAScript standard approved in 2015.
	CompilerOptions_LanguageMode_ECMASCRIPT_2015
	// ECMAScript standard approved in 2016. Adds the exponent operator (**).
	CompilerOptions_LanguageMode_ECMASCRIPT_2016
	// ECMAScript standard approved in 2017. Adds async/await and other syntax
	CompilerOptions_LanguageMode_ECMASCRIPT_2017
	// ECMAScript standard approved in 2018. Adds "..." in object literals/patterns.
	CompilerOptions_LanguageMode_ECMASCRIPT_2018
	// ECMAScript standard approved in 2019. Adds catch blocks with no error binding.
	CompilerOptions_LanguageMode_ECMASCRIPT_2019
	// ECMAScript standard approved in 2020.
	CompilerOptions_LanguageMode_ECMASCRIPT_2020
	// ECMAScript standard approved in 2021.
	CompilerOptions_LanguageMode_ECMASCRIPT_2021
	// ECMAScript standard approved in 2022.
	CompilerOptions_LanguageMode_ECMASCRIPT_2022
	// ECMAScript features from the upcoming standard.
	CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	// Use stable features.
	CompilerOptions_LanguageMode_STABLE
	// For languageOut only. The same language mode as the input.
	CompilerOptions_LanguageMode_NO_TRANSPILE
	// For testing only. Features that can be parsed/checked but cannot necessarily
	// be understood by the rest of the compiler yet.
	CompilerOptions_LanguageMode_UNSTABLE
	// For testing only. Features that can be parsed but cannot necessarily be
	// understood by the rest of the compiler yet.
	CompilerOptions_LanguageMode_UNSUPPORTED
)

var compilerOptions_LanguageMode_names = [...]string{"ECMASCRIPT3", "ECMASCRIPT5", "ECMASCRIPT5_STRICT", "ECMASCRIPT_2015", "ECMASCRIPT_2016", "ECMASCRIPT_2017", "ECMASCRIPT_2018", "ECMASCRIPT_2019", "ECMASCRIPT_2020", "ECMASCRIPT_2021", "ECMASCRIPT_2022", "ECMASCRIPT_NEXT", "STABLE", "NO_TRANSPILE", "UNSTABLE", "UNSUPPORTED"}

// CompilerOptions_LanguageMode_values is LanguageMode.values().
func CompilerOptions_LanguageMode_values() []CompilerOptions_LanguageMode {
	return []CompilerOptions_LanguageMode{CompilerOptions_LanguageMode_ECMASCRIPT3, CompilerOptions_LanguageMode_ECMASCRIPT5, CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT, CompilerOptions_LanguageMode_ECMASCRIPT_2015, CompilerOptions_LanguageMode_ECMASCRIPT_2016, CompilerOptions_LanguageMode_ECMASCRIPT_2017, CompilerOptions_LanguageMode_ECMASCRIPT_2018, CompilerOptions_LanguageMode_ECMASCRIPT_2019, CompilerOptions_LanguageMode_ECMASCRIPT_2020, CompilerOptions_LanguageMode_ECMASCRIPT_2021, CompilerOptions_LanguageMode_ECMASCRIPT_2022, CompilerOptions_LanguageMode_ECMASCRIPT_NEXT, CompilerOptions_LanguageMode_STABLE, CompilerOptions_LanguageMode_NO_TRANSPILE, CompilerOptions_LanguageMode_UNSTABLE, CompilerOptions_LanguageMode_UNSUPPORTED}
}

// String is Enum.name().
func (e CompilerOptions_LanguageMode) String() string { return compilerOptions_LanguageMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_LanguageMode) Ordinal() int { return int(e) }

var compilerOptions_LanguageMode_ptrs = func() []CompilerOptions_LanguageMode {
	v := CompilerOptions_LanguageMode_values()
	return v
}()

// Ptr returns the canonical pointer to e (Go-only: language-adapter — a
// nullable Java LanguageMode reference is a *CompilerOptions_LanguageMode; nil is null).
func (e CompilerOptions_LanguageMode) Ptr() *CompilerOptions_LanguageMode {
	return &compilerOptions_LanguageMode_ptrs[e]
}

// CompilerOptions_LanguageMode_STABLE_IN is the default input language
// level of JSCompiler.
//
// Java: CompilerOptions.LanguageMode.STABLE_IN
const CompilerOptions_LanguageMode_STABLE_IN = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT

// CompilerOptions_LanguageMode_STABLE_OUT is the default output language
// level of JSCompiler.
//
// Java: CompilerOptions.LanguageMode.STABLE_OUT
const CompilerOptions_LanguageMode_STABLE_OUT = CompilerOptions_LanguageMode_ECMASCRIPT5

// IsDefaultStrict reports whether this language mode defaults to strict
// mode
//
// Java: CompilerOptions.LanguageMode.isDefaultStrict()
func (e CompilerOptions_LanguageMode) IsDefaultStrict() bool {
	switch e {
	case CompilerOptions_LanguageMode_ECMASCRIPT3, CompilerOptions_LanguageMode_ECMASCRIPT5:
		return false
	default:
		return true
	}
}

// CompilerOptions_LanguageMode_ValidCommandLineNames returns a list of valid
// names used to select a `LanguageMode` on the command line.
//
// Java: CompilerOptions.LanguageMode.validCommandLineNames()
func CompilerOptions_LanguageMode_ValidCommandLineNames() []string {
	names := []string{}
	for _, mode := range CompilerOptions_LanguageMode_values() {
		if mode != CompilerOptions_LanguageMode_UNSUPPORTED {
			names = append(names, mode.String())
			if strings.HasPrefix(mode.String(), "ECMASCRIPT") {
				names = append(names, strings.Replace(mode.String(), "ECMASCRIPT", "ES", 1))
			}
		}
	}
	names = append(names, "ECMASCRIPT6")
	names = append(names, "ES6")
	names = append(names, "ECMASCRIPT6_STRICT")
	names = append(names, "ES6_STRICT")
	return names
}

// Java: CompilerOptions.LanguageMode.fromString(String)
func CompilerOptions_LanguageMode_FromString(value *string) *CompilerOptions_LanguageMode {
	if value == nil {
		return nil
	}
	// Trim spaces, disregard case, and allow abbreviation of ECMASCRIPT for convenience.
	canonicalizedName := jutil.StringToUpperCaseASCII(jutil.StringTrim(*value))
	if strings.HasPrefix(canonicalizedName, "ES") {
		canonicalizedName = "ECMASCRIPT" + canonicalizedName[2:]
	}

	if canonicalizedName == "ECMASCRIPT6" ||
		canonicalizedName == "ECMASCRIPT6_STRICT" {
		return CompilerOptions_LanguageMode_ECMASCRIPT_2015.Ptr()
	}

	// LanguageMode.valueOf(canonicalizedName), or null for an unknown name.
	for _, mode := range CompilerOptions_LanguageMode_values() {
		if mode.String() == canonicalizedName {
			return mode.Ptr()
		}
	}
	return nil // unknown name.
}

// Java: CompilerOptions.LanguageMode.toFeatureSet()
func (e CompilerOptions_LanguageMode) ToFeatureSet() *parser.FeatureSet {
	switch e {
	case CompilerOptions_LanguageMode_ECMASCRIPT3:
		return parser.FeatureSet_ES3
	case CompilerOptions_LanguageMode_ECMASCRIPT5, CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT:
		return parser.FeatureSet_ES5
	case CompilerOptions_LanguageMode_ECMASCRIPT_2015:
		return parser.FeatureSet_ES2015_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2016:
		return parser.FeatureSet_ES2016_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2017:
		return parser.FeatureSet_ES2017_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2018:
		return parser.FeatureSet_ES2018_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2019:
		return parser.FeatureSet_ES2019_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2020:
		return parser.FeatureSet_ES2020_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2021:
		return parser.FeatureSet_ES2021_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2022:
		return parser.FeatureSet_ES2022_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_NEXT:
		return parser.FeatureSet_ES_NEXT
	case CompilerOptions_LanguageMode_NO_TRANSPILE, CompilerOptions_LanguageMode_UNSTABLE:
		return parser.FeatureSet_ES_UNSTABLE
	case CompilerOptions_LanguageMode_UNSUPPORTED:
		return parser.FeatureSet_ES_UNSUPPORTED
	case CompilerOptions_LanguageMode_STABLE:
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"STABLE has different feature sets for language in and out. " +
				"Use STABLE_IN or STABLE_OUT."))
	}
	panic(jutil.NewIllegalStateException())
}

// CompilerOptions_DevMode is when to do the extra validity checks
type CompilerOptions_DevMode int32

const (
	// Don't do any extra checks.
	CompilerOptions_DevMode_OFF CompilerOptions_DevMode = iota
	// After the initial parse
	CompilerOptions_DevMode_START
	// At the start and at the end of all optimizations.
	CompilerOptions_DevMode_START_AND_END
	// After every pass
	CompilerOptions_DevMode_EVERY_PASS
)

var compilerOptions_DevMode_names = [...]string{"OFF", "START", "START_AND_END", "EVERY_PASS"}

// CompilerOptions_DevMode_values is DevMode.values().
func CompilerOptions_DevMode_values() []CompilerOptions_DevMode {
	return []CompilerOptions_DevMode{CompilerOptions_DevMode_OFF, CompilerOptions_DevMode_START, CompilerOptions_DevMode_START_AND_END, CompilerOptions_DevMode_EVERY_PASS}
}

// String is Enum.name().
func (e CompilerOptions_DevMode) String() string { return compilerOptions_DevMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_DevMode) Ordinal() int { return int(e) }

// CompilerOptions_TracerMode is how much tracing we want to do
type CompilerOptions_TracerMode int32

const (
	// Collect all timing and size metrics. Very slow.
	CompilerOptions_TracerMode_ALL CompilerOptions_TracerMode = iota
	// Collect all timing and size metrics, except gzipped size. Slow.
	CompilerOptions_TracerMode_RAW_SIZE
	// For size data, don't serialize the AST, just count the number of nodes.
	CompilerOptions_TracerMode_AST_SIZE
	// Collect timing metrics only.
	CompilerOptions_TracerMode_TIMING_ONLY
	// Collect no timing and size metrics.
	CompilerOptions_TracerMode_OFF
)

var compilerOptions_TracerMode_names = [...]string{"ALL", "RAW_SIZE", "AST_SIZE", "TIMING_ONLY", "OFF"}

// CompilerOptions_TracerMode_values is TracerMode.values().
func CompilerOptions_TracerMode_values() []CompilerOptions_TracerMode {
	return []CompilerOptions_TracerMode{CompilerOptions_TracerMode_ALL, CompilerOptions_TracerMode_RAW_SIZE, CompilerOptions_TracerMode_AST_SIZE, CompilerOptions_TracerMode_TIMING_ONLY, CompilerOptions_TracerMode_OFF}
}

// String is Enum.name().
func (e CompilerOptions_TracerMode) String() string { return compilerOptions_TracerMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_TracerMode) Ordinal() int { return int(e) }

// Java: CompilerOptions.TracerMode.isOn()
func (e CompilerOptions_TracerMode) IsOn() bool {
	return e != CompilerOptions_TracerMode_OFF
}

// CompilerOptions_TweakProcessing is the option for the ProcessTweaks pass
type CompilerOptions_TweakProcessing int32

const (
	// Do not run the ProcessTweaks pass.
	CompilerOptions_TweakProcessing_OFF CompilerOptions_TweakProcessing = iota
	// Run the pass, but do not strip out the calls.
	CompilerOptions_TweakProcessing_CHECK
	// Strip out all calls to goog.tweak.*.
	CompilerOptions_TweakProcessing_STRIP
)

var compilerOptions_TweakProcessing_names = [...]string{"OFF", "CHECK", "STRIP"}

// CompilerOptions_TweakProcessing_values is TweakProcessing.values().
func CompilerOptions_TweakProcessing_values() []CompilerOptions_TweakProcessing {
	return []CompilerOptions_TweakProcessing{CompilerOptions_TweakProcessing_OFF, CompilerOptions_TweakProcessing_CHECK, CompilerOptions_TweakProcessing_STRIP}
}

// String is Enum.name().
func (e CompilerOptions_TweakProcessing) String() string {
	return compilerOptions_TweakProcessing_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_TweakProcessing) Ordinal() int { return int(e) }

// Java: CompilerOptions.TweakProcessing.isOn()
func (e CompilerOptions_TweakProcessing) IsOn() bool {
	return e != CompilerOptions_TweakProcessing_OFF
}

// Java: CompilerOptions.TweakProcessing.shouldStrip()
func (e CompilerOptions_TweakProcessing) ShouldStrip() bool {
	return e == CompilerOptions_TweakProcessing_STRIP
}

// CompilerOptions_IsolationMode is what kind of isolation is going to be used
type CompilerOptions_IsolationMode int32

const (
	// output does not include additional isolation.
	CompilerOptions_IsolationMode_NONE CompilerOptions_IsolationMode = iota
	// The output should be wrapped in an IIFE to isolate global variables.
	CompilerOptions_IsolationMode_IIFE
)

var compilerOptions_IsolationMode_names = [...]string{"NONE", "IIFE"}

// CompilerOptions_IsolationMode_values is IsolationMode.values().
func CompilerOptions_IsolationMode_values() []CompilerOptions_IsolationMode {
	return []CompilerOptions_IsolationMode{CompilerOptions_IsolationMode_NONE, CompilerOptions_IsolationMode_IIFE}
}

// String is Enum.name().
func (e CompilerOptions_IsolationMode) String() string { return compilerOptions_IsolationMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_IsolationMode) Ordinal() int { return int(e) }

// CompilerOptions_SegmentOfCompilationToRun is what segment of the compilation to run.
type CompilerOptions_SegmentOfCompilationToRun int32

const (
	// all JSC compilation in a single action.
	CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION CompilerOptions_SegmentOfCompilationToRun = iota
	// all checks
	CompilerOptions_SegmentOfCompilationToRun_CHECKS
	// first half of optimizations
	CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF
	// second half of optimizations
	CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF
	// all optimizations
	CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS
	// all optimizations and finalizations in one action
	CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_AND_FINALIZATIONS
	// all finalizations
	CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS
)

var compilerOptions_SegmentOfCompilationToRun_names = [...]string{"ENTIRE_COMPILATION", "CHECKS", "OPTIMIZATIONS_FIRST_HALF", "OPTIMIZATIONS_SECOND_HALF", "OPTIMIZATIONS", "OPTIMIZATIONS_AND_FINALIZATIONS", "FINALIZATIONS"}

// CompilerOptions_SegmentOfCompilationToRun_values is SegmentOfCompilationToRun.values().
func CompilerOptions_SegmentOfCompilationToRun_values() []CompilerOptions_SegmentOfCompilationToRun {
	return []CompilerOptions_SegmentOfCompilationToRun{CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION, CompilerOptions_SegmentOfCompilationToRun_CHECKS, CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF, CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF, CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS, CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_AND_FINALIZATIONS, CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS}
}

// String is Enum.name().
func (e CompilerOptions_SegmentOfCompilationToRun) String() string {
	return compilerOptions_SegmentOfCompilationToRun_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_SegmentOfCompilationToRun) Ordinal() int { return int(e) }

// CompilerOptions_AliasStringsMode is a mode enum used to indicate the alias strings policy for the AliasStrings pass
type CompilerOptions_AliasStringsMode int32

const (
	// Do not alias string literals.
	CompilerOptions_AliasStringsMode_NONE CompilerOptions_AliasStringsMode = iota
	// Alias all string literals with a length greater than 100 characters.
	CompilerOptions_AliasStringsMode_LARGE
	// Alias all string literals where it may improve code size
	CompilerOptions_AliasStringsMode_ALL
	// Alias all string regardless of code size
	CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE
)

var compilerOptions_AliasStringsMode_names = [...]string{"NONE", "LARGE", "ALL", "ALL_AGGRESSIVE"}

// CompilerOptions_AliasStringsMode_values is AliasStringsMode.values().
func CompilerOptions_AliasStringsMode_values() []CompilerOptions_AliasStringsMode {
	return []CompilerOptions_AliasStringsMode{CompilerOptions_AliasStringsMode_NONE, CompilerOptions_AliasStringsMode_LARGE, CompilerOptions_AliasStringsMode_ALL, CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE}
}

// String is Enum.name().
func (e CompilerOptions_AliasStringsMode) String() string {
	return compilerOptions_AliasStringsMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_AliasStringsMode) Ordinal() int { return int(e) }

// CompilerOptions_Environment is an environment specifies the built-in externs that are loaded for a given compilation.
type CompilerOptions_Environment int32

const (
	// Hand crafted externs that have traditionally been the default externs.
	CompilerOptions_Environment_BROWSER CompilerOptions_Environment = iota
	// Only language externs are loaded.
	CompilerOptions_Environment_CUSTOM
)

var compilerOptions_Environment_names = [...]string{"BROWSER", "CUSTOM"}

// CompilerOptions_Environment_values is Environment.values().
func CompilerOptions_Environment_values() []CompilerOptions_Environment {
	return []CompilerOptions_Environment{CompilerOptions_Environment_BROWSER, CompilerOptions_Environment_CUSTOM}
}

// String is Enum.name().
func (e CompilerOptions_Environment) String() string { return compilerOptions_Environment_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_Environment) Ordinal() int { return int(e) }

// CompilerOptions_JsonStreamMode is whether standard input or standard output should be an array of JSON encoded files
type CompilerOptions_JsonStreamMode int32

const (
	// stdin/out are both single files.
	CompilerOptions_JsonStreamMode_NONE CompilerOptions_JsonStreamMode = iota
	// stdin is a json stream.
	CompilerOptions_JsonStreamMode_IN
	// stdout is a json stream.
	CompilerOptions_JsonStreamMode_OUT
	// stdin and stdout are both json streams.
	CompilerOptions_JsonStreamMode_BOTH
)

var compilerOptions_JsonStreamMode_names = [...]string{"NONE", "IN", "OUT", "BOTH"}

// CompilerOptions_JsonStreamMode_values is JsonStreamMode.values().
func CompilerOptions_JsonStreamMode_values() []CompilerOptions_JsonStreamMode {
	return []CompilerOptions_JsonStreamMode{CompilerOptions_JsonStreamMode_NONE, CompilerOptions_JsonStreamMode_IN, CompilerOptions_JsonStreamMode_OUT, CompilerOptions_JsonStreamMode_BOTH}
}

// String is Enum.name().
func (e CompilerOptions_JsonStreamMode) String() string {
	return compilerOptions_JsonStreamMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_JsonStreamMode) Ordinal() int { return int(e) }

// CompilerOptions_J2clPassMode is a mode enum used to indicate whether J2clPass should be enabled,
// disabled, or enabled automatically if there is any J2cl source file (i.e. in
// the AUTO mode).
type CompilerOptions_J2clPassMode int32

const (
	// J2clPass is disabled.
	CompilerOptions_J2clPassMode_OFF CompilerOptions_J2clPassMode = iota
	// It auto-detects whether there are J2cl generated file. If yes, execute J2clPass.
	CompilerOptions_J2clPassMode_AUTO
)

var compilerOptions_J2clPassMode_names = [...]string{"OFF", "AUTO"}

// CompilerOptions_J2clPassMode_values is J2clPassMode.values().
func CompilerOptions_J2clPassMode_values() []CompilerOptions_J2clPassMode {
	return []CompilerOptions_J2clPassMode{CompilerOptions_J2clPassMode_OFF, CompilerOptions_J2clPassMode_AUTO}
}

// String is Enum.name().
func (e CompilerOptions_J2clPassMode) String() string { return compilerOptions_J2clPassMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_J2clPassMode) Ordinal() int { return int(e) }

// Java: CompilerOptions.J2clPassMode.shouldAddJ2clPasses()
func (e CompilerOptions_J2clPassMode) ShouldAddJ2clPasses() bool {
	return e == CompilerOptions_J2clPassMode_AUTO
}

// Java: CompilerOptions.expectStrictModeInput()
func (o *CompilerOptions) ExpectStrictModeInput() bool {
	if o.isStrictModeInput != nil {
		return *o.isStrictModeInput
	}
	return o.GetLanguageIn().IsDefaultStrict()
}

// Java: CompilerOptions.setStrictModeInput(boolean)
func (o *CompilerOptions) SetStrictModeInput(isStrictModeInput bool) *CompilerOptions {
	o.isStrictModeInput = &isStrictModeInput
	return o
}

// Java: CompilerOptions.getPropertyReservedNamingFirstChars()
func (o *CompilerOptions) GetPropertyReservedNamingFirstChars() *jutil.LinkedHashSet[uint16] {
	if o.polymerPass {
		return compilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS
	} else if o.angularPass {
		return compilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS
	}
	return jutil.NewLinkedHashSet[uint16]()
}

// Java: CompilerOptions.getPropertyReservedNamingNonFirstChars()
func (o *CompilerOptions) GetPropertyReservedNamingNonFirstChars() *jutil.LinkedHashSet[uint16] {
	if o.polymerPass {
		return compilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS
	}
	return jutil.NewLinkedHashSet[uint16]()
}

// Java: CompilerOptions.shouldOptimize()
func (o *CompilerOptions) ShouldOptimize() bool {
	return !o.skipNonTranspilationPasses &&
		!o.checksOnly &&
		!o.ShouldGenerateTypedExterns() &&
		!o.instrumentForCoverageOnly
}

// compilerOptions_toStringHelper is MoreObjects.toStringHelper(this).omitNullValues().
//
// Go-only: language-adapter — Guava's ToStringHelper.
type compilerOptions_toStringHelper struct {
	sb    strings.Builder
	first bool
}

func newCompilerOptions_toStringHelper(className string) *compilerOptions_toStringHelper {
	h := &compilerOptions_toStringHelper{first: true}
	h.sb.WriteString(className)
	h.sb.WriteByte('{')
	return h
}

// add is ToStringHelper.add(String, Object): null values are omitted.
func (h *compilerOptions_toStringHelper) add(name string, value any) {
	if jutil.IsNil(value) {
		return
	}
	h.addUnconditional(name, value)
}

// addUnconditional is ToStringHelper.add(String, boolean/int).
func (h *compilerOptions_toStringHelper) addUnconditional(name string, value any) {
	if !h.first {
		h.sb.WriteString(", ")
	}
	h.first = false
	h.sb.WriteString(name)
	h.sb.WriteByte('=')
	h.sb.WriteString(compilerOptions_toJavaString(value))
}

func (h *compilerOptions_toStringHelper) String() string {
	return h.sb.String() + "}"
}

// compilerOptions_optional is a present or absent Guava Optional.
type compilerOptions_optional struct {
	value any // nil when absent
}

// compilerOptions_optionalOf wraps a nil-able pointer as an Optional.
func compilerOptions_optionalOf[T any](p *T) compilerOptions_optional {
	if p == nil {
		return compilerOptions_optional{}
	}
	return compilerOptions_optional{value: p}
}

// compilerOptions_toJavaString is String.valueOf(Object) for the option
// values.
//
// Go-only: language-adapter — slices are printed like Java lists, and a
// value whose Java class does not override toString() prints as
// Object.toString() (class name, '@', identity hash).
func compilerOptions_toJavaString(value any) string {
	switch v := value.(type) {
	case compilerOptions_optional:
		if v.value == nil {
			return "Optional.absent()"
		}
		return "Optional.of(" + compilerOptions_toJavaString(v.value) + ")"
	case *bool:
		return jutil.ToJavaString(*v)
	case string, bool, int, int32, int64, float64:
		return jutil.ToJavaString(v)
	case fmt.Stringer:
		return v.String()
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Slice {
		parts := make([]string, rv.Len())
		for i := range parts {
			parts[i] = compilerOptions_toJavaString(rv.Index(i).Interface())
		}
		return "[" + strings.Join(parts, ", ") + "]"
	}
	t := rv.Type()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return jutil.ObjectToString(
		"com.google.javascript.jscomp."+strings.ReplaceAll(t.Name(), "_", "$"), value)
}
