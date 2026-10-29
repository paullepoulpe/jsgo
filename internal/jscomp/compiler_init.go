// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"bytes"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// createMessageFormatter creates a message formatter instance
// corresponding to the value of CompilerOptions.
//
// Java: Compiler.createMessageFormatter()
func (c *Compiler) createMessageFormatter() MessageFormatter {
	colorize := c.options.ShouldColorizeErrorOutput()
	return c.options.GetErrorFormat().ToFormatter(c, colorize)
}

// Java: Compiler.initExperimentalOutputFeatureSetOptions(CompilerOptions)
func (c *Compiler) initExperimentalOutputFeatureSetOptions(options *CompilerOptions) {
	featureSet := options.GetOutputFeatureSet()
	if jutil.ObjectsEquals(featureSet, parser.FeatureSet_ES5) || jutil.ObjectsEquals(featureSet, parser.FeatureSet_ES3) {
		// do not change anything if language_out is configured to be ES5 or ES3.
		return
	}

	if options.GetExperimentalOutputFeatureSet() != nil {
		switch *options.GetExperimentalOutputFeatureSet() {
		case CompilerOptions_ExperimentalOutputFeatureSet_ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS:
			options.SetOutputFeatureSet(
				parser.FeatureSet_ES5.WithFeatures(
					parser.FeatureSet_Feature_ASYNC_FUNCTIONS,
					parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION,
					parser.FeatureSet_Feature_CONST_DECLARATIONS,
					parser.FeatureSet_Feature_GENERATORS,
					parser.FeatureSet_Feature_LET_DECLARATIONS))
		case CompilerOptions_ExperimentalOutputFeatureSet_BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD:
			options.SetOutputFeatureSet(
				parser.FeatureSet_BROWSER_2019.WithoutFeatures(
					parser.FeatureSet_Feature_CLASSES,
					parser.FeatureSet_Feature_CLASS_GETTER_SETTER,
					parser.FeatureSet_Feature_NEW_TARGET,
					parser.FeatureSet_Feature_SUPER,
					parser.FeatureSet_Feature_SPREAD_EXPRESSIONS))
		}
	}
}

// InitOptions initializes the compiler options. It's called as part of a
// normal compile() job. Public for the callers that are not doing a normal
// compile() job.
//
// Java: Compiler.initOptions(CompilerOptions)
func (c *Compiler) InitOptions(options *CompilerOptions) {
	c.options = options
	c.initASTCompaction()
	c.allowableFeatures = options.GetLanguageIn().ToFeatureSet()
	c.initExperimentalOutputFeatureSetOptions(options)
	if c.errorManager == nil {
		if c.outStream == nil {
			c.SetErrorManager(NewLoggerErrorManagerWithFormatter(c.createMessageFormatter(), Compiler_logger))
		} else {
			// ImmutableSet.builder(): duplicates (by identity) are dropped.
			var builder []SortingErrorManager_ErrorReportGenerator
			add := func(g SortingErrorManager_ErrorReportGenerator) {
				for _, e := range builder {
					if jutil.ObjectsEquals(e, g) {
						return
					}
				}
				builder = append(builder, g)
			}
			add(NewPrintStreamErrorReportGenerator(
				c.createMessageFormatter(), c.outStream, options.GetSummaryDetailLevel()))
			for _, g := range options.GetExtraReportGenerators() {
				add(g)
			}
			c.SetErrorManager(NewSortingErrorManager(builder))
		}
	}

	if options.GetSkipNonTranspilationPasses() {
		options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_NO_OP)
	}

	c.moduleLoader = moduleLoader_NewEmpty()

	c.ReconcileOptionsWithGuards()
	switch options.GetShouldValidateRequiredInlinings() {
	case base.Tri_TRUE:
		sufficientOptimizations := InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(options)
		jutil.CheckStateMsg1(
			sufficientOptimizations.ShouldValidate(),
			"shouldValidateRequiredInlinings cannot be set unless CompilerOptions enables"+
				" sufficient optimization, but found: %s",
			sufficientOptimizations.OptMessage())
	case base.Tri_FALSE:
		// do nothing.
	case base.Tri_UNKNOWN:
		options.SetValidateRequiredInlinings(
			InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(options).
				ShouldValidate())
	}

	// TODO(johnlenz): generally, the compiler should not be changing the options object
	// provided by the user.  This should be handled a different way.

	// Turn off type-based optimizations when type checking is off
	if !options.IsTypecheckingEnabled() {
		options.SetUseTypesForLocalOptimization(false)
		options.SetUseTypesForOptimization(false)
	}

	if options.AssumeForwardDeclaredForMissingTypes() {
		c.forwardDeclaredTypes = compilerForwardDeclaredSet{all: true}
	}

	c.initWarningsGuard(options.GetWarningsGuard())
	if c.IsDebugLoggingEnabled() {
		// If debug logs are requested, then we'll always log the configuration for convenience.
		options.SetPrintConfig(true)
	}

	// If pruning many unused inputs, it's possible that it's faster to regex-parse everything and
	// AST-parse only the non-pruned inputs, rather than just AST-parse everything once.
	c.preferRegexParser = c.preferRegexParser || options.GetDependencyOptions().ShouldPrune()

	if options.GetMergedPrecompiledLibraries() &&
		options.GetDependencyOptions().NeedsManagement() {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Using precompiled libraries (i.e. TypedAST) is incompatible with flags that " +
				"automatically order/prune dependencies: " +
				jutil.ToJavaString(options.GetDependencyOptions()) +
				""))
	}

	if options.GetMergedPrecompiledLibraries() && len(options.GetConformanceConfigs()) != 0 {
		var conformanceConfigFiles []string
		seenConfigs := make(map[string]struct{})
		for _, conformanceConfig := range options.GetConformanceConfigs() {
			for _, requirement := range conformanceConfig.GetRequirementList() {
				for _, cf := range requirement.GetConfigFileList() {
					if _, ok := seenConfigs[cf]; !ok {
						seenConfigs[cf] = struct{}{}
						conformanceConfigFiles = append(conformanceConfigFiles, cf)
					}
				}
			}
		}
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"use_precompiled_libraries mode does not support passing JS conformance configs: this" +
				" mode skips all check passes including conformance\n" +
				"Unsupported conformance configs: [" +
				strings.Join(conformanceConfigFiles, ", ") + "]"))
	}

	if options.ShouldRunDeadAssignmentElimination() &&
		!(options.ShouldRemoveUnusedVariables() || options.ShouldRemoveUnusedLocalVariables()) {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Invalid flag combination: enabling dead assignment elimination also requires unused" +
				" variable removal. Either disable dead assignment elimination or enable unused" +
				" variable removal."))
	}
	if options.ShouldRunDeadPropertyAssignmentElimination() && options.IsPolymerPassEnabled() {
		// The Polymer source is usually not included in the compilation, but it creates
		// getters/setters for many properties in compiled code. Dead property property assignment
		// elimination is only safe when it knows about getters/setters.
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Invalid flag combination: dead assignment elimination is incompatible with" +
				" setPolymerPass(true)."))
	}
}

// Java: Compiler.printConfig()
func (c *Compiler) PrintConfig() {
	if c.IsDebugLoggingEnabled() {
		// Log to separate files for convenience.
		c.logToFile("externs.log", func() string { return jutil.CollectionToString(c.externs) })
		c.logToFile(
			"inputs.json",
			func() string {
				return jutil.NewGsonBuilder().SetPrettyPrinting().Create().ToJson(c.chunkGraph.ToJson())
			})
		c.logToFile("options.log", func() string { return c.options.String() })
		c.logToFile("warningsGuard.log", func() string { return jutil.ToJavaString(c.warningsGuard) })
	} else {
		// Debug logging is not enabled, so use stderr
		err := jutil.System_Err()
		err.Println("==== Externs ====")
		err.Println(jutil.CollectionToString(c.externs))
		err.Println("==== Inputs ====")
		err.Println(jutil.NewGsonBuilder().SetPrettyPrinting().Create().ToJson(c.chunkGraph.ToJson()))
		err.Println("==== CompilerOptions ====")
		err.Println(c.options.String())
		err.Println("==== WarningsGuard ====")
		err.Println(jutil.ToJavaString(c.warningsGuard))
	}
}

// Java: Compiler.maybeLogPrunedInputs()
func (c *Compiler) maybeLogPrunedInputs() {
	if c.IsDebugLoggingEnabled() {
		c.logToFile(
			"inputs_pruned.json",
			func() string {
				return jutil.NewGsonBuilder().SetPrettyPrinting().Create().ToJson(c.chunkGraph.ToJson())
			})
	}
}

// Java: Compiler.logToFile(String, Supplier<String>)
func (c *Compiler) logToFile(logFileName string, logStringSupplier func() string) {
	log := c.CreateOrReopenLog("Compiler", logFileName)
	defer log.Close()
	log.LogSupplier(logStringSupplier)
}

// Java: Compiler.initWarningsGuard(WarningsGuard)
func (c *Compiler) initWarningsGuard(warningsGuard WarningsGuard) {
	guards := []WarningsGuard{}
	guards = append(guards,
		NewJ2clSuppressWarningsGuard(),
		NewClosureUnawareCodeWarningsGuard(),
		NewSuppressDocWarningsGuard(c, DiagnosticGroups_GetRegisteredGroups()),
		jutil.CheckNotNull(warningsGuard))

	c.warningsGuard = NewComposeWarningsGuardFromList(guards)
}

// ReconcileOptionsWithGuards: when the CompilerOptions and its
// WarningsGuard overlap, reconcile any discrepancies.
//
// Java: Compiler.reconcileOptionsWithGuards()
func (c *Compiler) ReconcileOptionsWithGuards() {
	options := c.options
	// DiagnosticGroups override the plain checkTypes option.
	if options.Enables(DiagnosticGroups_CHECK_TYPES) {
		options.SetCheckTypes(true)
	} else if options.Disables(DiagnosticGroups_CHECK_TYPES) {
		options.SetCheckTypes(false)
	} else if !options.GetCheckTypes() {
		// If DiagnosticGroups did not override the plain checkTypes
		// option, and checkTypes is disabled, then turn off the
		// parser type warnings.
		options.SetWarningLevel(
			DiagnosticGroup_ForType(RhinoErrorReporter_TYPE_PARSE_ERROR), CheckLevel_OFF)
	}

	if !options.GetCheckTypes() {
		options.SetWarningLevel(DiagnosticGroups_BOUNDED_GENERICS, CheckLevel_OFF)
	}

	// All passes must run the variable check. This synthesizes
	// variables later so that the compiler doesn't crash. It also
	// checks the externs file for validity. If you don't want to warn
	// about missing variable declarations, we shut that specific
	// error off.
	if !options.GetCheckSymbols() && !options.Enables(DiagnosticGroups_CHECK_VARIABLES) {
		options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	}

	// If we're in transpile-only mode, we don't need to do checks for suspicious var usages.
	// Since we still have to run VariableReferenceCheck before transpilation to check block-scoped
	// variables, though, we disable the unnecessary warnings it produces relating to vars here.
	if options.GetSkipNonTranspilationPasses() &&
		!options.Enables(DiagnosticGroups_CHECK_VARIABLES) {
		options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	}

	// If we're in transpile-only mode, we don't need to check for missing requires unless the user
	// explicitly enables missing-provide checks.
	if options.GetSkipNonTranspilationPasses() &&
		!options.Enables(DiagnosticGroups_MISSING_PROVIDE) {
		options.SetWarningLevel(DiagnosticGroups_MISSING_PROVIDE, CheckLevel_OFF)
	}

	// The default warning level of ARTIFICIAL_FUNCTION_PURITY_VALIDATION varies with whether or not
	// property disambiguation is enabled: if yes, it defaults to ERROR; if no, it defaults to OFF.
	// This is because without property disambiguation, it's often impossible for the compiler to
	// enforce function purity; and we expect projects that are code-size-sensitive enough to care
	// to have disambiguation already enabled.
	// We still want to allow individual projects to explicitly disable or enable this error.
	if options.ShouldDisambiguateProperties() &&
		options.
			GetWarningsGuard().
			MustRunChecks(DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION) == base.Tri_UNKNOWN {
		options.SetWarningLevel(
			DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION, CheckLevel_ERROR)
	}
}

func (c *Compiler) checkMissingTypedAstFile(file *SourceFile) {
	jutil.CheckStateMsg1(
		strings.HasPrefix(file.GetName(), compiler_SYNTHETIC_FILE_NAME_PREFIX) ||
			AbstractCompiler_IsFillFileName(file.GetName()),
		"TypedAST filesystem initialized, but missing requested file: %s",
		file)
}

func (c *Compiler) deserializeTypedAstRoot(file *SourceFile) (rhino.Node, bool) {
	if c.typedAstFilesystem == nil {
		return rhino.NilNode, false
	}
	entry, ok := c.typedAstFilesystem.Remove(file)
	if !ok {
		c.checkMissingTypedAstFile(file)
		return rhino.NilNode, false
	}
	return entry.DeserializeNew(), true
}

// Java: Compiler.getTypedAstDeserializer(SourceFile)
func (c *Compiler) GetTypedAstDeserializer(file *SourceFile) func() rhino.Node {
	if c.typedAstFilesystem == nil {
		return nil
	}
	entry, ok := c.typedAstFilesystem.Remove(file)
	if !ok {
		c.checkMissingTypedAstFile(file)
		return nil
	}
	return entry.DeserializeNew
}

func (c *Compiler) takeTypedAstDeserializer(file *SourceFile) (*ScriptNodeDeserializer, func() rhino.Node) {
	if c.typedAstFilesystem == nil {
		return nil, nil
	}
	entry, ok := c.typedAstFilesystem.Remove(file)
	if !ok {
		c.checkMissingTypedAstFile(file)
		return nil, nil
	}
	if entry.sd != nil {
		return entry.sd, typedAstDeserializer_scriptSentinel
	}
	return nil, entry.DeserializeNew
}

// InitWithTypedAstFilesystem initializes a compiler with deserialized state
// from the given TypedAst.List
//
// This method initializes all the state needed to run `.stage2Passes()` or
// do any similar sort of optimization work.
//
// typedAstListStream: a gzipped, binary-serialized TypedAst.List proto
//
// Java: Compiler.initWithTypedAstFilesystem(List<SourceFile>, List<SourceFile>, CompilerOptions, InputStream)
func (c *Compiler) InitWithTypedAstFilesystem(
	externs []*SourceFile,
	sources []*SourceFile,
	options *CompilerOptions,
	typedAstListStream io.Reader) {

	capHint := len(externs) + len(sources)
	requiredOrder := make([]*SourceFile, 0, capHint)
	requiredSet := make(map[*SourceFile]struct{}, capHint)
	for _, f := range externs {
		if _, dup := requiredSet[f]; !dup {
			requiredSet[f] = struct{}{}
			requiredOrder = append(requiredOrder, f)
		}
	}
	for _, f := range sources {
		if _, dup := requiredSet[f]; !dup {
			requiredSet[f] = struct{}{}
			requiredOrder = append(requiredOrder, f)
		}
	}

	options.SetMergedPrecompiledLibraries(true)

	c.InitOptions(options)
	c.Init(externs, sources, options)
	c.mergeAndDeserializeTypedAsts(requiredOrder, requiredSet, typedAstListStream, options)
}

// InitChunksWithTypedAstFilesystem initializes a compiler with deserialized
// state from the given TypedAst.List
//
// This method initializes all the state needed to run `.stage2Passes()` or
// do any similar sort of optimization work.
//
// typedAstListStream: a gzipped, binary-serialized TypedAst.List proto
//
// Java: Compiler.initChunksWithTypedAstFilesystem(List<SourceFile>, List<JSChunk>, CompilerOptions, InputStream)
func (c *Compiler) InitChunksWithTypedAstFilesystem(
	externs []*SourceFile,
	chunks []*JSChunk,
	options *CompilerOptions,
	typedAstListStream io.Reader) {

	capHint := len(externs)
	for _, chunk := range chunks {
		capHint += len(chunk.GetInputs())
	}
	requiredOrder := make([]*SourceFile, 0, capHint)
	requiredSet := make(map[*SourceFile]struct{}, capHint)
	for _, f := range externs {
		if _, dup := requiredSet[f]; !dup {
			requiredSet[f] = struct{}{}
			requiredOrder = append(requiredOrder, f)
		}
	}
	for _, chunk := range chunks {
		for _, input := range chunk.GetInputs() {
			f := input.GetSourceFile()
			if _, dup := requiredSet[f]; !dup {
				requiredSet[f] = struct{}{}
				requiredOrder = append(requiredOrder, f)
			}
		}
	}

	options.SetMergedPrecompiledLibraries(true)

	c.InitOptions(options)
	c.InitChunks(externs, chunks, options)
	c.mergeAndDeserializeTypedAsts(requiredOrder, requiredSet, typedAstListStream, options)
}

// Java: Compiler.mergeAndDeserializeTypedAsts(ImmutableSet<SourceFile>, InputStream, CompilerOptions)
func (c *Compiler) mergeAndDeserializeTypedAsts(
	requiredOrder []*SourceFile,
	requiredSet map[*SourceFile]struct{},
	typedAstListStream io.Reader,
	options *CompilerOptions) {
	jutil.CheckState(c.typedAstFilesystem == nil)
	c.MaybeSetTracker()

	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
	// To speed up builds that don't run type-based optimizations, skip type deserialization
	deserializeTypes := options.RequiresTypesForOptimization()

	astData := Compiler_RunInCompilerThread(c,
		func() *TypedAstDeserializer_DeserializedAst {
			tracer := c.NewTracer("deserializeTypedAst")
			defer c.StopTracer(tracer, "deserializeTypedAst")
			return typedAstDeserializer_deserializeFullAstFromSet(
				c,
				c.SYNTHETIC_EXTERNS_FILE,
				requiredOrder,
				requiredSet,
				typedAstListStream,
				deserializeTypes,
				options.GetResolveSourceMapAnnotations(),
				options.GetParseInlineSourceMaps())
		})

	c.typedAstFilesystem = astData.filesystem
	c.externProperties = astData.GetExternProperties()
	c.colorRegistry = astData.GetColorRegistry()
	c.SetTypeCheckingHasRun(deserializeTypes)

	for library := range astData.GetRuntimeLibraries().All() {
		c.runtimeJsLibManager.EnsureLibraryInjected(library, false)
	}

	c.GetSynthesizedExternsInput() // Force lazy creation.

	// TODO(lharker): refactor things to avoid this misleading 'parse' call. It's not actually
	// parsing JS source code, just doing more deserialization & building the Rhino AST.
	c.runInCompilerThread(c.parseForCompilationInternal)
}

// compiler_getResourceAsStream is Compiler.class.getResourceAsStream(path):
// nil (Java null) when there is no such resource.
//
// Go-only: language-adapter — classpath resources are embedded in package
// resources.
func compiler_getResourceAsStream(path string) io.Reader {
	data := resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.Compiler", path)
	if data == nil {
		return nil
	}
	return bytes.NewReader(data)
}

// Java: Compiler.initRuntimeLibraryTypedAsts(Optional<ColorPool.Builder>)
func (c *Compiler) InitRuntimeLibraryTypedAsts(colorPoolBuilder *ColorPool_Builder) {
	jutil.CheckState(c.runtimeLibraryTypedAsts == nil)

	path := strings.Join([]string{
		"/runtime_libs.typedast"}, "")

	var stream io.Reader
	if data := resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.Compiler", path); data != nil {
		stream = &runtimeLibsTypedAstStream{data: data}
	} else {
		stream = compiler_getResourceAsStream(path)
	}

	astData := TypedAstDeserializer_DeserializeRuntimeLibraries(
		c,
		c.SYNTHETIC_EXTERNS_FILE,
		colorPoolBuilder,
		stream,
		c.GetOptions().GetResolveSourceMapAnnotations(),
		c.GetOptions().GetParseInlineSourceMaps())

	// Re-index the runtime libraries by file name rather than SourceFile object
	fs := astData.filesystem
	runtimeLibraryTypedAsts := make(map[string]*ScriptNodeDeserializer, len(fs))
	for key, entry := range fs { // order-independent: each SourceFile in filePoolBuilder has a unique filename
		if entry.sd != nil {
			runtimeLibraryTypedAsts[key.GetName()] = entry.sd
		}
	}

	c.runtimeLibraryTypedAsts = runtimeLibraryTypedAsts
}

// Init initializes the instance state needed for a compile job.
//
// Java: Compiler.init(List<SourceFile>, List<SourceFile>, CompilerOptions)
func (c *Compiler) Init(externs []*SourceFile, sources []*SourceFile, options *CompilerOptions) {
	chunk := NewJSChunk(JSChunk_STRONG_CHUNK_NAME)
	for _, source := range sources {
		chunk.Add(NewCompilerInput(source))
	}

	chunks := make([]*JSChunk, 0, 1)
	chunks = append(chunks, chunk)
	c.InitChunks(externs, chunks, options)
	c.addFilesToSourceMap(sources)
}

// InitChunks initializes the instance state needed for a compile job if
// the sources are in chunks.
//
// Java: Compiler.initChunks(List<SourceFile>, List<JSChunk>, CompilerOptions)
func (c *Compiler) InitChunks(externs []*SourceFile, chunks []*JSChunk, options *CompilerOptions) {
	c.InitOptions(options)

	c.checkFirstChunk(chunks)

	clear(c.externs)
	c.externs = c.externs[:0]
	for _, file := range externs {
		c.externs = append(c.externs, NewCompilerInputWithIsExtern(file /* isExtern= */, true))
	}

	// Generate the chunk graph, and report any errors in the chunk specification as errors.
	chunkGraph, depErr := tryNewJSChunkGraph(chunks)
	if depErr != nil {
		// problems with the chunk format.  Report as an error.  The
		// message gives all details.
		c.Report(
			JSError_Make(
				Compiler_CHUNK_DEPENDENCY_ERROR, depErr.GetChunk().GetName(), depErr.GetDependentChunk().GetName()))
		return
	}
	c.chunkGraph = chunkGraph

	// Creating the chunk graph can move weak source around, and end up with empty chunks.
	c.fillEmptyChunks(c.GetChunks())

	c.commentsPerFile = make(map[string][]parser.Comment)
	c.InitBasedOnOptions()

	c.InitInputsByIdMap()

	c.initAST()

	if c.IsDebugLoggingEnabled() || options.ShouldPrintConfig() {
		// Always print configuration logs when debug logging is enabled.
		// NOTE: initChunks() is called by every execution path through the compiler code that
		// actually performs a compilation action. That's what makes this a good place to put this
		// logging.
		c.PrintConfig()
	}

	// If debug logging is enabled, write out the module / chunk graph in GraphViz format.
	// This graph is often too big to reasonably render.
	// Using gvpr(1) is recommended to extract the parts of the graph that are of interest.
	// `dot -Tpng graph_file.dot > graph_file.png` will render an image.
	func() {
		chunkGraphLog := c.CreateOrReopenLog("Compiler", "chunk_graph.dot")
		defer chunkGraphLog.Close()
		if chunkGraphLog.IsLogging() {
			chunkGraphLog.LogString(DotFormatter_ToDotOfGraphvizGraph(c.chunkGraph.ToGraphvizGraph()))
		}
	}()

	var nodeForCodeInsertion func() rhino.Node
	if options.GetRuntimeLibraryMode() == RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES {
		nodeForCodeInsertion = c.getSynthesizedExternsRoot
	} else {
		nodeForCodeInsertion = c.getNodeForRuntimeCodeInsertion
	}
	c.runtimeJsLibManager = RuntimeJsLibManager_Create(
		options.GetRuntimeLibraryMode(),
		c.loadResourceContents,
		c.changeTracker,
		nodeForCodeInsertion)
}

// Java: Compiler.getRuntimeJsLibManager()
func (c *Compiler) GetRuntimeJsLibManager() *RuntimeJsLibManager {
	return c.runtimeJsLibManager
}

func (c *Compiler) putInputSourceMap(name string, sm *SourceMapInput) {
	if _, exists := c.inputSourceMaps[name]; !exists {
		c.inputSourceMapOrder = append(c.inputSourceMapOrder, name)
	}
	c.inputSourceMaps[name] = sm
}

// InitBasedOnOptions does any initialization that is dependent on the
// compiler options.
//
// Java: Compiler.initBasedOnOptions()
func (c *Compiler) InitBasedOnOptions() {
	for k, v := range c.options.GetInputSourceMaps().All() {
		c.putInputSourceMap(k, v)
	}
	// Create the source map if necessary.
	if c.options.ShouldGatherSourceMapInfo() {
		c.sourceMap = c.options.GetSourceMapFormat().GetInstance()
		c.sourceMap.SetPrefixMappings(c.options.GetSourceMapLocationMappings())
		// Go-only: fast mode encodes the mappings on workers (byte-identical).
		if c.options.isFastMode() && c.options.GetNumParallelThreads() > 1 && c.sourceMap.generatorV3 != nil {
			c.sourceMap.generatorV3.SetParallelWorkers(c.options.GetNumParallelThreads())
		}
		if c.options.GetApplyInputSourceMaps() {
			c.sourceMap.SetCompilerSourceMapping(c)
			if c.options.GetSourceMapIncludeSourcesContent() {
				for _, name := range c.inputSourceMapOrder {
					c.addSourceMapSourceFiles(c.inputSourceMaps[name])
				}
			}
		}
	}
}

// Java: Compiler.EMPTY_CHUNK_LIST_ERROR
var compiler_EMPTY_CHUNK_LIST_ERROR = DiagnosticType_Error("JSC_EMPTY_CHUNK_LIST_ERROR", "At least one chunk must be provided")

// Java: Compiler.EMPTY_ROOT_CHUNK_ERROR
var compiler_EMPTY_ROOT_CHUNK_ERROR = DiagnosticType_Error(
	"JSC_EMPTY_ROOT_CHUNK_ERROR",
	"Root chunk ''{0}'' must contain at least one source code input")

// checkFirstChunk verifies that at least one chunk has been provided and
// that the first one has at least one source code input.
//
// Java: Compiler.checkFirstChunk(List<JSChunk>)
func (c *Compiler) checkFirstChunk(chunks []*JSChunk) {
	if len(chunks) == 0 {
		c.Report(JSError_Make(compiler_EMPTY_CHUNK_LIST_ERROR))
	} else if len(chunks[0].GetInputs()) == 0 && len(chunks) > 1 {
		// The root chunk may only be empty if there is exactly 1 chunk.
		c.Report(JSError_Make(compiler_EMPTY_ROOT_CHUNK_ERROR, chunks[0].GetName()))
	}
}

// Compiler_JoinPathParts creates an OS specific path string from parts
//
// Java: Compiler.joinPathParts(String...)
func Compiler_JoinPathParts(pathParts ...string) string {
	return compiler_pathJoiner(pathParts)
}

// fillEmptyChunks fills any empty chunks with a place holder file. It makes
// any cross chunk motion easier.
//
// Java: Compiler.fillEmptyChunks(Iterable<JSChunk>)
func (c *Compiler) fillEmptyChunks(chunks []*JSChunk) {
	for _, chunk := range chunks {
		if chunk.GetName() != JSChunk_WEAK_CHUNK_NAME && len(chunk.GetInputs()) == 0 {
			input := NewCompilerInput(SourceFile_FromCode(AbstractCompiler_CreateFillFileName(chunk.GetName()), ""))
			input.SetCompiler(c)
			chunk.Add(input)
		}
	}
}

// RebuildInputsFromChunks rebuilds the internal input map by iterating
// over all chunks. This is necessary if inputs have been added to or
// removed from a chunk after an init or initChunks call.
//
// Java: Compiler.rebuildInputsFromChunks()
func (c *Compiler) RebuildInputsFromChunks() {
	c.InitInputsByIdMap()
}

// Java: Compiler.DUPLICATE_INPUT
var Compiler_DUPLICATE_INPUT = DiagnosticType_Error("JSC_DUPLICATE_INPUT", "Duplicate input: {0}")

// Java: Compiler.DUPLICATE_EXTERN_INPUT
var Compiler_DUPLICATE_EXTERN_INPUT = DiagnosticType_Error("JSC_DUPLICATE_EXTERN_INPUT", "Duplicate extern input: {0}")

// InitInputsByIdMap creates a map to make looking up an input by name fast.
// Also checks for duplicate inputs.
//
// Java: Compiler.initInputsByIdMap()
func (c *Compiler) InitInputsByIdMap() {
	clear(c.inputsById)
	for _, input := range c.externs {
		previous := c.putCompilerInput(input)
		if previous != nil {
			c.Report(JSError_Make(Compiler_DUPLICATE_EXTERN_INPUT, input.GetName()))
		}
	}

	for _, input := range c.chunkGraph.GetAllInputs() {
		previous := c.putCompilerInput(input)
		if previous != nil {
			c.Report(JSError_Make(Compiler_DUPLICATE_INPUT, input.GetName()))
		}
	}
}

// initAST sets up the skeleton of the AST (the externs and root).
//
// Java: Compiler.initAST()
func (c *Compiler) initAST() {
	c.jsRoot = rhino.IR_Root()
	c.externsRoot = rhino.IR_Root()
	c.externAndJsRoot = rhino.IR_Root(c.externsRoot, c.jsRoot)
}

// Java: Compiler.isDebugLoggingEnabled()
func (c *Compiler) IsDebugLoggingEnabled() bool {
	return c.GetOptions().GetDebugLogDirectory() != nil
}

// Java: Compiler.getDebugLogFilterList()
func (c *Compiler) GetDebugLogFilterList() []string {
	if c.GetOptions().GetDebugLogFilter() == nil {
		return []string{}
	}
	var out []string
	for _, s := range strings.Split(*c.GetOptions().GetDebugLogFilter(), ",") {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// Java: Compiler.createOrReopenLog(Class<?>, String, String...)
func (c *Compiler) CreateOrReopenLog(owner string, firstNamePart string, restNameParts ...string) diagnostic.LogFile {
	if !c.IsDebugLoggingEnabled() {
		return diagnostic.LogFile_CreateNoOp()
	}

	dir := c.GetOptions().GetDebugLogDirectory()
	relativeParts := jutil.Path_Of(firstNamePart, restNameParts...)
	file := dir.ResolveString(owner).Resolve(relativeParts)

	// If a filter list for log file names was provided, only create a log file if any
	// of the filter strings matches.
	filters := c.GetDebugLogFilterList()
	if len(filters) == 0 {
		return diagnostic.LogFile_CreateOrReopen(file)
	}

	for _, filter := range filters {
		if strings.Contains(file.String(), filter) {
			return diagnostic.LogFile_CreateOrReopen(file)
		}
	}
	return diagnostic.LogFile_CreateNoOp()
}

// Java: Compiler.createOrReopenIndexedLog(Class<?>, String, String...)
func (c *Compiler) CreateOrReopenIndexedLog(owner string, firstNamePart string, restNameParts ...string) diagnostic.LogFile {
	jutil.CheckStateMsg1(c.currentPassIndex >= 0, "%s", c.currentPassIndex)

	index := strconv.Itoa(c.currentPassIndex)
	for len(index) < 3 {
		index = "0" + index
	}
	length := len(restNameParts)
	if length == 0 {
		firstNamePart = index + "_" + firstNamePart
	} else {
		restNameParts[length-1] = index + "_" + restNameParts[length-1]
	}

	return c.CreateOrReopenLog(owner, firstNamePart, restNameParts...)
}

// Java: Compiler.maybePrintSourceAfterEachPass(String)
func (c *Compiler) maybePrintSourceAfterEachPass(passName string) {
	if !c.options.ShouldPrintSourceAfterEachPass() {
		return
	}

	currentJsSource := c.GetCurrentJsSource()
	if c.hasLastJsSource && currentJsSource == c.lastJsSource {
		return
	}

	if c.IsDebugLoggingEnabled() {
		func() {
			log := c.CreateOrReopenIndexedLog("Compiler", "source_after_pass", passName)
			defer log.Close()
			log.LogString(currentJsSource)
		}()
	} else {
		err := jutil.System_Err()
		err.PrintlnEmpty()
		if c.hasDebugMessage {
			err.Println("// DEBUG: " + c.debugMessage)
		}
		err.Println("// " + passName + " yields:")
		err.Println("// ************************************")
		err.Println(currentJsSource)
		err.Println("// ************************************")
	}

	c.lastJsSource = currentJsSource
	c.hasLastJsSource = true
}

// Java: Compiler.getCurrentJsSource()
func (c *Compiler) GetCurrentJsSource() string {
	c.ResetAndIntitializeSourceMap()

	fileNameRegexList := c.options.GetFilesToPrintAfterEachPassRegexList()
	chunkNameRegexList := c.options.GetChunksToPrintAfterEachPassRegexList()
	var qnameOrder []string
	var qnameSet map[string]struct{}
	if rawQnames := c.options.GetQnameUsesToPrintAfterEachPassList(); len(rawQnames) > 0 {
		qnameSet = make(map[string]struct{}, len(rawQnames))
		for _, q := range rawQnames {
			if _, ok := qnameSet[q]; !ok {
				qnameSet[q] = struct{}{}
				qnameOrder = append(qnameOrder, q)
			}
		}
	}
	var builder strings.Builder

	if len(fileNameRegexList) == 0 && len(chunkNameRegexList) == 0 && len(qnameOrder) == 0 {
		return c.ToSource()
	}
	if len(fileNameRegexList) != 0 {
		rhino.CheckNotNull(c.externsRoot)
		rhino.CheckNotNull(c.jsRoot)
		for _, r := range []rhino.Node{c.externsRoot, c.jsRoot} {
			for fileNode := r.GetFirstChild(); !fileNode.IsNil(); fileNode = fileNode.GetNext() {
				fileName := fileNode.GetSourceFileName()
				for _, regex := range fileNameRegexList {
					if jutil.Pattern_Compile(regex).Matches(*fileName) {
						source := "// " + *fileName + "\n" + c.ToSourceOfNode(fileNode)
						builder.WriteString(source)
						break
					}
				}
			}
		}
		if builder.Len() == 0 {
			builder.WriteString("// No files matched any of: ")
			builder.WriteString(jutil.CollectionToString(fileNameRegexList))
		}
	}
	if len(chunkNameRegexList) != 0 {
		chunks := jutil.CheckNotNull(c.GetChunks())
		unmatchedChunkNames := []string{}
		// As we are emitting all chunks at once in dependency order, we can be smarter about
		// deciding which licenses to emit where in the source.
		lt := NewCompiler_ChunkGraphAwareLicenseTracker(c)
		for _, jsChunk := range chunks {
			chunkName := jsChunk.GetName()
			lt.SetCurrentChunkContext(jsChunk)
			for _, regex := range chunkNameRegexList {
				if jutil.Pattern_Compile(regex).Matches(chunkName) {
					source := "// module '" + chunkName + "'\n" + c.ToSourceOfLicenseTrackerAndJSChunk(lt, jsChunk)
					builder.WriteString(source)
					break
				}
			}
			unmatchedChunkNames = append(unmatchedChunkNames, chunkName)
		}
		if builder.Len() == 0 {
			builder.WriteString("// No chunks were matched:\n")
			for _, chunkName := range unmatchedChunkNames {
				builder.WriteString("// ")
				builder.WriteString(chunkName)
				builder.WriteString("\n")
			}
		}
	}
	if len(qnameOrder) != 0 {
		// Keep track of how a qualified name we're interested in gets renamed.
		// Note that a qname can actually be renamed multiple times during compilation,
		// so we need a multimap here.
		originalToNewQNameMap := make(map[string][]string)
		// Get a stream of top level statement nodes that contain at least one reference to one of the
		// qnames we are interested in. Take note of relevant renamings while we're at it.
		//
		// Go-only: language-adapter — Java's lazy Stream.filter runs when
		// the stream is consumed (after the header below is written), so
		// the filter is applied in the loop at the end.
		statements := c.getTopLevelStatements(c.jsRoot)
		statementFilter := func(statement rhino.Node) bool {
			// filter to just those statements containing interesting qnames
			return NodeUtil_Has(
				statement,
				func(n rhino.Node) bool {
					rhino.CheckNotNull(n)
					if !n.IsQualifiedName() {
						return false
					}
					qname := n.GetQualifiedName()
					if _, ok := qnameSet[*qname]; ok {
						return true
					}
					originalQname := n.GetOriginalQualifiedName()
					if originalQname != nil {
						if _, ok := qnameSet[*originalQname]; ok {
							existing := originalToNewQNameMap[*originalQname]
							if !slices.Contains(existing, *qname) {
								originalToNewQNameMap[*originalQname] = append(existing, *qname)
							}
							return true
						}
					}
					return false
				},
				func(node rhino.Node) bool { return true })
		}
		builder.WriteString("//\n")
		builder.WriteString("// closure-compiler: Printing all of the top-level statements\n")
		builder.WriteString("// that contain references to these qualified names.\n")
		builder.WriteString("//\n")
		for _, qname := range qnameOrder {
			builder.WriteString("// '")
			builder.WriteString(qname)
			builder.WriteString("'\n")
			for _, newName := range originalToNewQNameMap[qname] {
				builder.WriteString("// '")
				builder.WriteString(newName)
				builder.WriteString("' (originally '")
				builder.WriteString(qname)
				builder.WriteString("')\n")
			}
		}
		builder.WriteString("//\n")
		for _, statement := range statements {
			if !statementFilter(statement) {
				continue
			}
			builder.WriteString("// ")
			builder.WriteString(statement.GetLocation())
			builder.WriteByte('\n')
			builder.WriteString(c.ToSourceOfNode(statement))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

// Java: Compiler.getTopLevelStatements(Node)
func (c *Compiler) getTopLevelStatements(root rhino.Node) []rhino.Node {
	cb := &compiler_topLevelStatementsCallback{}
	NodeTraversal_Traverse(c, root, cb)
	return cb.builder
}

// compiler_topLevelStatementsCallback is the anonymous
// AbstractPreOrderCallback of getTopLevelStatements: it collects the
// outermost statements.
type compiler_topLevelStatementsCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	builder []rhino.Node
}

// Java: Compiler.getTopLevelStatements(...).new AbstractPreOrderCallback() {...}.shouldTraverse(NodeTraversal, Node, Node)
func (cb *compiler_topLevelStatementsCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if NodeUtil_IsStatement(n) {
		cb.builder = append(cb.builder, n)
		return false
	} else {
		return true
	}
}

// Java: Compiler.initializeModuleLoader()
func (c *Compiler) InitializeModuleLoader() {

	var moduleResolverFactory ModuleLoader_ModuleResolverFactory
	switch c.options.GetModuleResolutionMode() {
	case ModuleLoader_ResolutionMode_BROWSER:
		moduleResolverFactory = BrowserModuleResolver_FACTORY
	case ModuleLoader_ResolutionMode_NODE:
		// processJsonInputs requires a module loader to already be defined
		// so we redefine it afterwards with the package.json inputs
		moduleResolverFactory = NewNodeModuleResolver_FactoryWithEntries(c.ProcessJsonInputs(c.chunkGraph.GetAllInputs()))
	case ModuleLoader_ResolutionMode_WEBPACK:
		moduleResolverFactory = NewWebpackModuleResolver_Factory(c.inputPathByWebpackId)
	case ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES:
		moduleResolverFactory = NewBrowserWithTransformedPrefixesModuleResolver_Factory(
			c.options.GetBrowserResolverPrefixReplacements())
	}

	inputs := c.chunkGraph.GetAllInputs()
	dependencyInfos := make([]DependencyInfo, len(inputs))
	for i, in := range inputs {
		dependencyInfos[i] = in
	}
	c.moduleLoader = ModuleLoader_Builder_().
		SetModuleRoots(c.options.GetModuleRoots()).
		SetInputs(dependencyInfos).
		SetFactory(moduleResolverFactory).
		SetPathResolver(ModuleLoader_PathResolver_RELATIVE).
		SetPathEscaper(c.options.GetPathEscaper()).
		Build()
}
