// Port of com/google/javascript/jscomp/AbstractCompiler.java @ closure-compiler v20261021

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
	"reflect"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AbstractCompiler is an abstract compiler, to help remove the circular
// dependency of passes on JSCompiler.
//
// This is an abstract class, so that we can make the methods
// package-private.
//
// Note: Please avoid adding any new local state or concrete implementations
// to this class, as there will be concrete implementations of
// AbstractCompiler where inheriting a default implementation is likely
// incorrect. Local State should instate be declared in any concrete
// implementations of AbstractCompiler that need it.
//
// The Go interface declares every abstract and overridable method;
// implementations embed AbstractCompilerBase (PORTING.md §4).
//
// Go-only: language-adapter — Java Iterable/List/ImmutableList results are
// []T; Set/ImmutableSet values are *jutil.LinkedHashSet[T] (insertion
// order); java.util.function.Supplier<T> is func() T; Class<?> is
// reflect.Type; a Java null String result is a nil *string.
type AbstractCompiler interface {
	SourceExcerptProvider
	CompilerInputProvider

	// GetAbstractCompilerBase gives access to the fields of the base class.
	GetAbstractCompilerBase() *AbstractCompilerBase

	// BeforePass will be called before each pass runs.
	//
	// Java: AbstractCompiler.beforePass(String)
	BeforePass(passName string)

	// AfterPass will be called after each pass finishes.
	//
	// Java: AbstractCompiler.afterPass(String)
	AfterPass(passName string)

	// TODO(nicksantos): Decide if all of these are really necessary.
	// Many of them are just accessors that should be passed to the
	// CompilerPass's constructor.

	// Java: AbstractCompiler.getTypedAstDeserializer(SourceFile)
	GetTypedAstDeserializer(file *SourceFile) func() rhino.Node

	// GetSourceFileByName looks up a source file by name. May return nil.
	//
	// Java: AbstractCompiler.getSourceFileByName(String)
	GetSourceFileByName(sourceName string) *SourceFile

	// Java: AbstractCompiler.getScriptNode(String)
	GetScriptNode(filename string) rhino.Node

	// GetChunkGraph gets the chunk graph.
	//
	// Java: AbstractCompiler.getChunkGraph()
	GetChunkGraph() *JSChunkGraph

	// GetInputsInOrder gets the inputs in the order in which they are being
	// processed. Only for use by AbstractCompilerRunner.
	//
	// Java: AbstractCompiler.getInputsInOrder()
	GetInputsInOrder() []*CompilerInput

	// GetNumberOfInputs gets the total number of inputs.
	//
	// This can be useful as a guide for the initial allocated size for data
	// structures.
	//
	// Java: AbstractCompiler.getNumberOfInputs()
	GetNumberOfInputs() int

	//
	// Intermediate state and results produced and needed by particular passes.
	// TODO(rluble): move these into the general structure for keeping state between pass runs.
	//

	// AddExportedNames adds exported names to keep track.
	//
	// Java: AbstractCompiler.addExportedNames(Set<String>)
	AddExportedNames(exportedVariableNames *jutil.LinkedHashSet[string])

	// GetExportedNames gets the names that have been exported.
	//
	// Java: AbstractCompiler.getExportedNames()
	GetExportedNames() *jutil.LinkedHashSet[string]

	// SetDefineNames adds @define names to keep track.
	//
	// Java: AbstractCompiler.setDefineNames(Collection<String>)
	SetDefineNames(defineNames []string)

	// Java: AbstractCompiler.getDefineNames()
	GetDefineNames() *jutil.LinkedHashSet[string]

	// SetVariableMap sets the variable renaming map
	//
	// Java: AbstractCompiler.setVariableMap(VariableMap)
	SetVariableMap(variableMap *VariableMap)

	// SetPropertyMap sets the property renaming map
	//
	// Java: AbstractCompiler.setPropertyMap(VariableMap)
	SetPropertyMap(propertyMap *VariableMap)

	// SetStringMap sets the string replacement map
	//
	// Java: AbstractCompiler.setStringMap(VariableMap)
	SetStringMap(stringMap *VariableMap)

	// SetCssNames sets the css names found during compilation.
	//
	// Java: AbstractCompiler.setCssNames(Set<String>)
	SetCssNames(newCssNames *jutil.LinkedHashSet[string])

	// SetInstrumentationMapping sets the mapping for instrumentation
	// parameter encoding.
	//
	// Java: AbstractCompiler.setInstrumentationMapping(VariableMap)
	SetInstrumentationMapping(instrumentationMapping *VariableMap)

	// SetIdGeneratorMap sets the id generator for cross-chunk motion.
	//
	// Java: AbstractCompiler.setIdGeneratorMap(String)
	SetIdGeneratorMap(serializedIdMappings string)

	// GetTranspiledFiles gets whether any file needed to transpile any
	// feature
	//
	// Java: AbstractCompiler.getTranspiledFiles()
	GetTranspiledFiles() bool

	// GetCrossChunkIdGenerator gets the id generator for cross-chunk motion.
	//
	// Java: AbstractCompiler.getCrossChunkIdGenerator()
	GetCrossChunkIdGenerator() *IdGenerator

	// SetAnonymousFunctionNameMap sets the naming map for anonymous
	// functions
	//
	// Java: AbstractCompiler.setAnonymousFunctionNameMap(VariableMap)
	SetAnonymousFunctionNameMap(functionMap *VariableMap)

	//
	// End of intermediate state needed by passes.
	//

	// SetTypeCheckingHasRun sets whether the typechecking passes have run.
	//
	// Java: AbstractCompiler.setTypeCheckingHasRun(boolean)
	SetTypeCheckingHasRun(hasRun bool)

	// HasTypeCheckingRun returns whether the typechecking passes have run
	//
	// Java: AbstractCompiler.hasTypeCheckingRun()
	HasTypeCheckingRun() bool

	// HasOptimizationColors reports whether the AST has been annotated with
	// optimization colors.
	//
	// Java: AbstractCompiler.hasOptimizationColors()
	HasOptimizationColors() bool

	// IsTypeRegistryCleared returns `true` when type checking has run, but
	// the type registry has been cleared.
	//
	// See also `clearJSTypeRegistry()`.
	//
	// Java: AbstractCompiler.isTypeRegistryCleared()
	IsTypeRegistryCleared() bool

	// GetTypeRegistry gets a central registry of type information from the
	// compiled JS.
	//
	// Java: AbstractCompiler.getTypeRegistry()
	GetTypeRegistry() *rhino.JSTypeRegistry

	// Java: AbstractCompiler.clearJSTypeRegistry()
	ClearJSTypeRegistry()

	// GetColorRegistry gets a central registry of colors from deserialized
	// JS types.
	//
	// Java: AbstractCompiler.getColorRegistry()
	GetColorRegistry() *colors.ColorRegistry

	// SetColorRegistry sets the color registry
	//
	// Java: AbstractCompiler.setColorRegistry(ColorRegistry)
	SetColorRegistry(registry *colors.ColorRegistry)

	// Java: AbstractCompiler.forwardDeclareType(String)
	ForwardDeclareType(typeName string)

	// GetTypedScopeCreator gets a memoized scope creator with type
	// information.
	//
	// Java: AbstractCompiler.getTypedScopeCreator()
	GetTypedScopeCreator() ScopeCreator

	// Java: AbstractCompiler.clearTypedScopeCreator()
	ClearTypedScopeCreator()

	// GetTopScope gets the top scope.
	//
	// Java: AbstractCompiler.getTopScope()
	GetTopScope() *TypedScope

	// SetTopScope sets the top scope.
	//
	// Java: AbstractCompiler.setTopScope(TypedScope)
	SetTopScope(x *TypedScope)

	// GetTranspilationNamespace returns a scope containing only externs and
	// synthetic code or other code in the first script.
	//
	// Intended for transpilation passes to look up types when synthesizing
	// new code.
	//
	// Java: AbstractCompiler.getTranspilationNamespace()
	GetTranspilationNamespace() rhino.StaticScope

	// Report reports an error or warning.
	//
	// Java: AbstractCompiler.report(JSError)
	Report(error *JSError)

	// ThrowInternalError reports an internal error.
	//
	// Java: AbstractCompiler.throwInternalError(String, Throwable)
	ThrowInternalError(msg string, cause jutil.Throwable)

	// GetCodingConvention gets the current coding convention.
	//
	// Java: AbstractCompiler.getCodingConvention()
	GetCodingConvention() CodingConvention

	// ReportChangeToEnclosingScope: passes that make modifications in a
	// scope that is different than the Compiler.currentScope use this (eg,
	// InlineVariables and many others)
	//
	// Java: AbstractCompiler.reportChangeToEnclosingScope(Node)
	ReportChangeToEnclosingScope(n rhino.Node)

	// ReportChangeToChangeScope: mark modifications in a scope that is
	// different than the Compiler.currentScope use this (eg, InlineVariables
	// and many others)
	//
	// Java: AbstractCompiler.reportChangeToChangeScope(Node)
	ReportChangeToChangeScope(changeScopeRoot rhino.Node)

	// ReportFunctionDeleted marks a specific function node as known to be
	// deleted. Is part of having accurate change tracking which is necessary
	// to streamline optimizations.
	//
	// Java: AbstractCompiler.reportFunctionDeleted(Node)
	ReportFunctionDeleted(node rhino.Node)

	// ReportDisambiguatePropertiesSummary is used by
	// `DisambiguateProperties` to record a one-line summary of the work it
	// accomplished, if any.
	//
	// This information will be included in the tracer mode output, if that
	// is enabled.
	//
	// Java: AbstractCompiler.reportDisambiguatePropertiesSummary(Supplier<String>)
	ReportDisambiguatePropertiesSummary(summarySupplier func() string)

	// ReportAmbiguatePropertiesSummary is used by `AmbiguateProperties` to
	// record a one-line summary of the work it accomplished, if any.
	//
	// This information will be included in the tracer mode output, if that
	// is enabled.
	//
	// Java: AbstractCompiler.reportAmbiguatePropertiesSummary(Supplier<String>)
	ReportAmbiguatePropertiesSummary(summarySupplier func() string)

	// GetNodeForCodeInsertion gets a suitable SCRIPT node to serve as a
	// parent for code insertion. If chunk contains any inputs, the returned
	// node will be the SCRIPT node corresponding to its first input. If
	// chunk is empty, on the other hand, then the returned node will be the
	// first SCRIPT node in a non-empty chunk that chunk depends on (the
	// deepest one possible).
	//
	// chunk is a chunk. If nil, will return the first SCRIPT node of all
	// chunks. The result is a SCRIPT node (never nil).
	//
	// Java: AbstractCompiler.getNodeForCodeInsertion(JSChunk)
	GetNodeForCodeInsertion(chunk *JSChunk) rhino.Node

	// Java: AbstractCompiler.getTypeValidator()
	GetTypeValidator() *TypeValidator

	// GetTypeMismatches gets the central registry of type violations.
	//
	// Java: AbstractCompiler.getTypeMismatches()
	GetTypeMismatches() []*TypeMismatch

	// Java: AbstractCompiler.setExternExports(String)
	SetExternExports(externExports string)

	// ParseSyntheticCode parses code for injecting, and associate it with a
	// given source file.
	//
	// Java: AbstractCompiler.parseSyntheticCode(String, String)
	ParseSyntheticCode(filename string, code string) rhino.Node

	// ParseTestCode parses code for testing.
	//
	// Java: AbstractCompiler.parseTestCode(String)
	ParseTestCode(code string) rhino.Node

	// ParseTestCodeOfList parses code for testing.
	//
	// Java: AbstractCompiler.parseTestCode(ImmutableList<String>)
	ParseTestCodeOfList(code []string) rhino.Node

	// ToSource prints a node to source code.
	//
	// Java: AbstractCompiler.toSource()
	ToSource() string

	// ToSourceOfNode prints a node to source code.
	//
	// Java: AbstractCompiler.toSource(Node)
	ToSourceOfNode(root rhino.Node) string

	// PreferRegexParser reports whether to prefer using JsFileRegexParser
	// when calculating DependencyInfo as opposed to triggering an AST-based
	// parse.
	//
	// This is primarily for performance reasons. The two should produce the
	// same results except in some edge cases.
	//
	// Java: AbstractCompiler.preferRegexParser()
	PreferRegexParser() bool

	// GetDefaultErrorReporter gets a default error reporter for injecting
	// into Rhino.
	//
	// Java: AbstractCompiler.getDefaultErrorReporter()
	GetDefaultErrorReporter() rhino.ErrorReporter

	// GetReverseAbstractInterpreter gets an interpreter for type analysis.
	//
	// Java: AbstractCompiler.getReverseAbstractInterpreter()
	GetReverseAbstractInterpreter() ReverseAbstractInterpreter

	// GetLifeCycleStage returns the current life-cycle stage of the AST
	// we're working on.
	//
	// Java: AbstractCompiler.getLifeCycleStage()
	GetLifeCycleStage() AbstractCompiler_LifeCycleStage

	// InitRuntimeLibraryTypedAsts deserializes runtime libraries from a
	// TypedAST packaged as a JAR resource and reconcile their Colors with
	// the current inputs.
	//
	// This method must be called anywhere that Colors are reconciled for
	// application to the AST. Otherwise Color information won't be
	// consistent. `colorPoolBuilder` must be the same builder as used for
	// the other inputs, and the caller retains ownership.
	//
	// colorPoolBuilder, if present (non-nil), includes inferred optimization
	// colors on the deserialized ASTs. If absent, does not include colors.
	//
	// Java: AbstractCompiler.initRuntimeLibraryTypedAsts(Optional<ColorPool.Builder>)
	InitRuntimeLibraryTypedAsts(colorPoolBuilder *ColorPool_Builder)

	// GetUniqueIdSupplier generates unique String Ids when requested via a
	// compiler instance.
	//
	// This supplier provides Ids that are deterministic and unique across
	// all input files given to the compiler. The generated ID format is:
	// uniqueId = "fileHashCode$counterForThisFile"
	//
	// Java: AbstractCompiler.getUniqueIdSupplier()
	GetUniqueIdSupplier() *UniqueIdSupplier

	// GetUniqueNameIdSupplier generates unique ids.
	//
	// Deprecated: because the generated names during transpilation are not
	// unique across all input files. Use the new supplier by calling
	// GetUniqueIdSupplier().
	//
	// Java: AbstractCompiler.getUniqueNameIdSupplier()
	GetUniqueNameIdSupplier() func() string

	// HasHaltingErrors reports whether any errors have been encountered that
	// should stop the compilation process.
	//
	// Java: AbstractCompiler.hasHaltingErrors()
	HasHaltingErrors() bool

	// Java: AbstractCompiler.getChangeTracker()
	GetChangeTracker() *ChangeTracker

	// AddIndexProvider registers a provider for some type of index.
	//
	// Go-only: language-adapter — IndexProvider<?> is IndexProvider[any]:
	// a Go interface method cannot take a type parameter.
	//
	// Java: AbstractCompiler.addIndexProvider(IndexProvider<?>)
	AddIndexProvider(indexProvider IndexProvider[any])

	// GetIndex returns, from a provider, the desired index of type T,
	// otherwise nil if no provider is registered for the given type.
	//
	// Go-only: language-adapter — <T> T getIndex(Class<T>) takes the
	// reflect.Type and returns the index as any.
	//
	// Java: AbstractCompiler.getIndex(Class<T>)
	GetIndex(type_ reflect.Type) any

	// GetJsRoot returns the root of the source tree, ignoring externs
	//
	// Java: AbstractCompiler.getJsRoot()
	GetJsRoot() rhino.Node

	// HasScopeChanged: true iff a function changed since the last time a
	// pass was run
	//
	// Java: AbstractCompiler.hasScopeChanged(Node)
	HasScopeChanged(n rhino.Node) bool

	// GetParserConfig returns the parser configuration for the specified
	// context.
	//
	// Java: AbstractCompiler.getParserConfig(ConfigContext)
	GetParserConfig(context AbstractCompiler_ConfigContext) *parsing.Config

	// GetErrorManager gets the error manager.
	//
	// Java: AbstractCompiler.getErrorManager()
	GetErrorManager() ErrorManager

	// SetLifeCycleStage sets the current life-cycle state.
	//
	// Java: AbstractCompiler.setLifeCycleStage(LifeCycleStage)
	SetLifeCycleStage(stage AbstractCompiler_LifeCycleStage)

	// AreNodesEqualForInlining: are the nodes equal for the purpose of
	// inlining? If type aware optimizations are on, type equality is
	// checked.
	//
	// Java: AbstractCompiler.areNodesEqualForInlining(Node, Node)
	AreNodesEqualForInlining(n1 rhino.Node, n2 rhino.Node) bool

	// SetHasRegExpGlobalReferences sets if RegExp global properties are
	// used.
	//
	// references is whether there are references to the RegExp global
	// object properties.
	//
	// Java: AbstractCompiler.setHasRegExpGlobalReferences(boolean)
	SetHasRegExpGlobalReferences(references bool)

	// HasRegExpGlobalReferences reports whether the AST contains references
	// to the RegExp global object properties.
	//
	// Java: AbstractCompiler.hasRegExpGlobalReferences()
	HasRegExpGlobalReferences() bool

	// GetErrorLevel returns the error level the given error object will be
	// reported at (nil when no guard decides).
	//
	// Java: AbstractCompiler.getErrorLevel(JSError)
	GetErrorLevel(error *JSError) *CheckLevel

	// Process runs a given compiler-pass by calling its process() method.
	//
	// pass is the pass to be run.
	//
	// Java: AbstractCompiler.process(CompilerPass)
	Process(pass CompilerPass)

	// GetRoot returns the root node of the AST, which includes both externs
	// and source.
	//
	// Java: AbstractCompiler.getRoot()
	GetRoot() rhino.Node

	// Java: AbstractCompiler.getOptions()
	GetOptions() *CompilerOptions

	// GetAllowableFeatures returns the set of language features currently
	// allowed to exist in the AST
	//
	// At the start of compilation, all features in the languageIn are
	// allowed. At the end of compilation, only features in the languageOut
	// should be allowed.
	//
	// Transpilation passes should mark the transpiled features as "not
	// allowed" to prevent future passes from accidentally introducing new
	// uses of those features. For example, if transpiling away ES2017
	// features, the RewriteAsyncFunctions pass will mark async functions as
	// no longer allowed.
	//
	// Java: AbstractCompiler.getAllowableFeatures()
	GetAllowableFeatures() *parser.FeatureSet

	// SetAllowableFeatures sets the feature set allowed in the AST. See
	// GetAllowableFeatures()
	//
	// In most cases, callers of this method should only be shrinking the
	// allowable feature set and not adding to it. One known exception is the
	// ConvertChunksToESModules pass.
	//
	// Java: AbstractCompiler.setAllowableFeatures(FeatureSet)
	SetAllowableFeatures(featureSet *parser.FeatureSet)

	// MarkFeatureNotAllowed marks feature as no longer allowed in the AST.
	// See GetAllowableFeatures()
	//
	// Java: AbstractCompiler.markFeatureNotAllowed(FeatureSet.Feature)
	MarkFeatureNotAllowed(feature parser.FeatureSet_Feature)

	// MarkFeatureSetNotAllowed marks features as no longer allowed in the
	// AST. See GetAllowableFeatures()
	//
	// Java: AbstractCompiler.markFeatureSetNotAllowed(FeatureSet)
	MarkFeatureSetNotAllowed(featureSet *parser.FeatureSet)

	// GetSynthesizedExternsInput returns a CompilerInput that can be
	// modified to add additional extern definitions to the beginning of the
	// externs AST
	//
	// Java: AbstractCompiler.getSynthesizedExternsInput()
	GetSynthesizedExternsInput() *CompilerInput

	// GetSynthesizedTypeSummaryInput returns a CompilerInput that can be
	// modified to add additional type summary definitions to the beginning
	// of the externs AST (note that the externs AST includes type summary
	// input).
	//
	// Java: AbstractCompiler.getSynthesizedTypeSummaryInput()
	GetSynthesizedTypeSummaryInput() *CompilerInput

	// GetLastPassName gets the last pass name set by setProgress (nil when
	// there is none).
	//
	// Java: AbstractCompiler.getLastPassName()
	GetLastPassName() *string

	// GetRuntimeJsLibManager returns a class responsible for managing
	// injection of runtime libraries under the js/ subdirectory.
	//
	// Java: AbstractCompiler.getRuntimeJsLibManager()
	GetRuntimeJsLibManager() *RuntimeJsLibManager

	// SetExternProperties sets the names of the properties defined in
	// externs.
	//
	// externProperties is the set of property names defined in externs.
	//
	// Java: AbstractCompiler.setExternProperties(ImmutableSet<String>)
	SetExternProperties(externProperties *jutil.LinkedHashSet[string])

	// GetExternProperties gets the names of the properties defined in
	// externs or nil if GatherExternProperties pass was not run yet.
	//
	// Java: AbstractCompiler.getExternProperties()
	GetExternProperties() *jutil.LinkedHashSet[string]

	// AddInputSourceMap adds a SourceMapInput for the given sourceFileName,
	// to be used for error reporting and source map combining.
	//
	// Java: AbstractCompiler.addInputSourceMap(String, SourceMapInput)
	AddInputSourceMap(name string, sourceMap *SourceMapInput)

	// Java: AbstractCompiler.getBase64SourceMapContents(String)
	GetBase64SourceMapContents(sourceFileName string) *string

	// Java: AbstractCompiler.addComments(String, List<Comment>)
	AddComments(filename string, comments []*parser.Comment)

	// GetAccessorSummary returns a summary an entry for every property name
	// found in the AST with a getter and / or setter defined.
	//
	// Property names for which there are no getters or setters will not be
	// in the map.
	//
	// Java: AbstractCompiler.getAccessorSummary()
	GetAccessorSummary() *AccessorSummary

	// SetAccessorSummary sets the summary of properties with getters and
	// setters.
	//
	// Java: AbstractCompiler.setAccessorSummary(AccessorSummary)
	SetAccessorSummary(summary *AccessorSummary)

	// GetComments returns all the comments from the given file.
	//
	// Java: AbstractCompiler.getComments(String)
	GetComments(filename string) []*parser.Comment

	// GetModuleLoader gets the module loader.
	//
	// Java: AbstractCompiler.getModuleLoader()
	GetModuleLoader() *ModuleLoader

	// GetModuleTypeByName looks up the type of a module from its name (nil
	// when unknown).
	//
	// Java: AbstractCompiler.getModuleTypeByName(String)
	GetModuleTypeByName(moduleName string) *CompilerInput_ModuleType

	// SetRunJ2clPasses sets whether J2CL passes should run
	//
	// Java: AbstractCompiler.setRunJ2clPasses(boolean)
	SetRunJ2clPasses(runJ2clPasses bool)

	// RunJ2clPasses reports whether J2CL passes should run
	//
	// Java: AbstractCompiler.runJ2clPasses()
	RunJ2clPasses() bool

	// CreateAstFactory returns a new AstFactory that will add type
	// information to the nodes it creates if and only if type checking has
	// already happened and types have not been converted into colors.
	//
	// Note that the AstFactory will /not/ add colors to the AST if types
	// have been converted into colors. The AstFactory does not understand
	// colors, although color support could certainly be added if it proves
	// useful.
	//
	// Java: AbstractCompiler.createAstFactory()
	CreateAstFactory() *AstFactory

	// CreateAstFactoryWithoutTypes returns a new AstFactory that will not
	// add type information, regardless of whether type checking has already
	// happened.
	//
	// Java: AbstractCompiler.createAstFactoryWithoutTypes()
	CreateAstFactoryWithoutTypes() *AstFactory

	// GetAstAnalyzer returns a new AstAnalyzer configured correctly to
	// answer questions about Nodes in the AST currently being compiled.
	//
	// Java: AbstractCompiler.getAstAnalyzer()
	GetAstAnalyzer() *AstAnalyzer

	// Java: AbstractCompiler.createDefaultExpressionDecomposer()
	CreateDefaultExpressionDecomposer() *ExpressionDecomposer

	// Java: AbstractCompiler.createExpressionDecomposer(Supplier<String>, ImmutableSet<String>, Scope)
	CreateExpressionDecomposer(
		uniqueNameIdSupplier func() string,
		knownConstantFunctions *jutil.LinkedHashSet[string],
		scope *Scope) *ExpressionDecomposer

	// Java: AbstractCompiler.getModuleMetadataMap()
	GetModuleMetadataMap() *ModuleMetadataMap

	// Java: AbstractCompiler.setModuleMetadataMap(ModuleMetadataMap)
	SetModuleMetadataMap(moduleMetadataMap *ModuleMetadataMap)

	// Java: AbstractCompiler.getModuleMap()
	GetModuleMap() *ModuleMap

	// Java: AbstractCompiler.setModuleMap(ModuleMap)
	SetModuleMap(moduleMap *ModuleMap)

	// Java: AbstractCompiler.isDebugLoggingEnabled()
	IsDebugLoggingEnabled() bool

	// Java: AbstractCompiler.getDebugLogFilterList()
	GetDebugLogFilterList() []string

	// CreateOrReopenLog provides logging access to a file with the specified
	// name.
	//
	// Java: AbstractCompiler.createOrReopenLog(Class<?>, String, String...)
	CreateOrReopenLog(owner reflect.Type, firstNamePart string, restNameParts ...string) diagnostic.LogFile

	// CreateOrReopenIndexedLog provides logging access to a file with the
	// specified name, differentiated by the index of the current pass.
	//
	// Indexing helps in separating logs from different pass loops. The
	// filename pattern is
	// "[debug_log_directory]/[owner_name]/([name_part[i]]/){0,n-1}[pass_index]_[name_part[n]]".
	//
	// Java: AbstractCompiler.createOrReopenIndexedLog(Class<?>, String, String...)
	CreateOrReopenIndexedLog(owner reflect.Type, firstNamePart string, restNameParts ...string) diagnostic.LogFile

	// GetSyntheticCodeInputId returns the InputId of the synthetic code
	// input (even if it is not initialized yet).
	//
	// Java: AbstractCompiler.getSyntheticCodeInputId()
	GetSyntheticCodeInputId() *rhino.InputId

	// InitializeSyntheticCodeInput adds a synthetic script to the front of
	// the AST
	//
	// Useful to allow inserting code into the global scope, earlier than any
	// of the user-provided code, in the case that the first user-provided
	// input is a module.
	//
	// Java: AbstractCompiler.initializeSyntheticCodeInput()
	InitializeSyntheticCodeInput()

	// RemoveSyntheticCodeInput removes the script added by
	// InitializeSyntheticCodeInput
	//
	// Java: AbstractCompiler.removeSyntheticCodeInput()
	RemoveSyntheticCodeInput()

	// MergeSyntheticCodeInput merges all code in the script added by
	// InitializeSyntheticCodeInput into the first non-synthetic script. Will
	// crash if the first non-synthetic script is a module and module
	// rewriting has not occurred.
	//
	// Java: AbstractCompiler.mergeSyntheticCodeInput()
	MergeSyntheticCodeInput()
}

// AbstractCompilerBase holds the (field-less) state and the concrete
// methods of AbstractCompiler. Implementations embed it and call
// InitAbstractCompilerBase from their constructor.
type AbstractCompilerBase struct {
	self AbstractCompiler
}

// InitAbstractCompilerBase is the Go form of AbstractCompiler's (implicit)
// constructor: it records the most-derived object for virtual calls.
func (b *AbstractCompilerBase) InitAbstractCompilerBase(self AbstractCompiler) {
	b.self = self
}

// GetAbstractCompilerBase gives access to the fields of the base class.
func (b *AbstractCompilerBase) GetAbstractCompilerBase() *AbstractCompilerBase { return b }

// Java: AbstractCompiler.READ_ERROR
var AbstractCompiler_READ_ERROR = DiagnosticType_Error("JSC_READ_ERROR", "Cannot read file {0}: {1}")

// Java: AbstractCompiler.FILL_FILE_SUFFIX
const abstractCompiler_FILL_FILE_SUFFIX = "$fillFile"

// AbstractCompiler_CreateFillFileName: empty chunks get an empty "fill"
// file, so that we can move code into an empty chunk.
//
// Java: AbstractCompiler.createFillFileName(String)
func AbstractCompiler_CreateFillFileName(chunkName string) string {
	return chunkName + abstractCompiler_FILL_FILE_SUFFIX
}

// AbstractCompiler_IsFillFileName returns whether a file name was created by
// AbstractCompiler_CreateFillFileName.
//
// Java: AbstractCompiler.isFillFileName(String)
func AbstractCompiler_IsFillFileName(fileName string) bool {
	return strings.HasSuffix(fileName, abstractCompiler_FILL_FILE_SUFFIX)
}

// InitRuntimeLibraryTypedAsts is the base implementation, which is not
// J2CL compatible.
//
// Java: AbstractCompiler.initRuntimeLibraryTypedAsts(Optional<ColorPool.Builder>)
func (b *AbstractCompilerBase) InitRuntimeLibraryTypedAsts(colorPoolBuilder *ColorPool_Builder) {
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Implementation in Compiler.java is not J2CL compatible."))
}

// AbstractCompiler_ConfigContext represents the different contexts for
// which the compiler could have distinct configurations.
type AbstractCompiler_ConfigContext int32

const (
	// Normal JavaScript.
	AbstractCompiler_ConfigContext_DEFAULT AbstractCompiler_ConfigContext = iota

	// Externs files.
	AbstractCompiler_ConfigContext_EXTERNS
)

var abstractCompiler_ConfigContext_names = [...]string{"DEFAULT", "EXTERNS"}

// AbstractCompiler_ConfigContext_values is ConfigContext.values().
func AbstractCompiler_ConfigContext_values() []AbstractCompiler_ConfigContext {
	return []AbstractCompiler_ConfigContext{
		AbstractCompiler_ConfigContext_DEFAULT, AbstractCompiler_ConfigContext_EXTERNS,
	}
}

// String is Enum.name().
func (e AbstractCompiler_ConfigContext) String() string {
	return abstractCompiler_ConfigContext_names[e]
}

// Ordinal is Enum.ordinal().
func (e AbstractCompiler_ConfigContext) Ordinal() int { return int(e) }

// AbstractCompiler_LifeCycleStage is what point in optimizations we're in.
// For use by compiler passes
type AbstractCompiler_LifeCycleStage int32

const (
	AbstractCompiler_LifeCycleStage_RAW AbstractCompiler_LifeCycleStage = iota

	// See constraints put on the AST by either running the
	// ConvertTypesToColors.java pass /or/ by having created the AST via
	// serialization/TypedAstDeserializer.java
	// NORMALIZED implies this constraint
	AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC

	// See constraints put on the tree by Normalize.java
	AbstractCompiler_LifeCycleStage_NORMALIZED

	// The normalize pass has put constraints on the tree,
	// but variables and properties have been renamed so
	// coding conventions no longer apply.
	AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED
)

var abstractCompiler_LifeCycleStage_names = [...]string{
	"RAW", "COLORS_AND_SIMPLIFIED_JSDOC", "NORMALIZED", "NORMALIZED_OBFUSCATED",
}

// AbstractCompiler_LifeCycleStage_values is LifeCycleStage.values().
func AbstractCompiler_LifeCycleStage_values() []AbstractCompiler_LifeCycleStage {
	return []AbstractCompiler_LifeCycleStage{
		AbstractCompiler_LifeCycleStage_RAW,
		AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC,
		AbstractCompiler_LifeCycleStage_NORMALIZED,
		AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED,
	}
}

// String is Enum.name().
func (s AbstractCompiler_LifeCycleStage) String() string {
	return abstractCompiler_LifeCycleStage_names[s]
}

// Ordinal is Enum.ordinal().
func (s AbstractCompiler_LifeCycleStage) Ordinal() int { return int(s) }

// Java: AbstractCompiler.LifeCycleStage.isNormalized()
func (s AbstractCompiler_LifeCycleStage) IsNormalized() bool {
	return s == AbstractCompiler_LifeCycleStage_NORMALIZED ||
		s == AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED
}

// Java: AbstractCompiler.LifeCycleStage.isNormalizedUnobfuscated()
func (s AbstractCompiler_LifeCycleStage) IsNormalizedUnobfuscated() bool {
	return s == AbstractCompiler_LifeCycleStage_NORMALIZED
}

// Java: AbstractCompiler.LifeCycleStage.isNormalizedObfuscated()
func (s AbstractCompiler_LifeCycleStage) IsNormalizedObfuscated() bool {
	return s == AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED
}

// Java: AbstractCompiler.LifeCycleStage.hasColorAndSimplifiedJSDoc()
func (s AbstractCompiler_LifeCycleStage) HasColorAndSimplifiedJSDoc() bool {
	return s != AbstractCompiler_LifeCycleStage_RAW
}

// Java: AbstractCompiler.RUNTIME_LIB_DIR
const AbstractCompiler_RUNTIME_LIB_DIR = RuntimeJsLibManager_RUNTIME_LIB_DIR
