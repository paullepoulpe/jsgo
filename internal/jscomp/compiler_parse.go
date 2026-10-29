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
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ------------------------------------------------------------------------
// Parsing
// ------------------------------------------------------------------------

// ParseInputs parses the externs and main inputs.
//
// Returns a synthetic root node whose two children are the externs root
// and the main root
//
// Java: Compiler.parseInputs()
func (c *Compiler) ParseInputs() (result rhino.Node) {
	devMode := c.options.GetDevMode() != CompilerOptions_DevMode_OFF

	// If old roots exist (we are parsing a second time), detach each of the
	// individual file parse trees.
	c.externsRoot.DetachChildren()
	c.jsRoot.DetachChildren()
	c.scriptNodeByFilename.Clear()

	tracer := c.NewTracer(PassNames_PARSE_INPUTS)
	c.BeforePass(PassNames_PARSE_INPUTS)

	defer func() {
		c.AfterPass(PassNames_PARSE_INPUTS)
		c.StopTracer(tracer, PassNames_PARSE_INPUTS)
	}()
	// Parse externs sources.
	if c.options.GetNumParallelThreads() > 1 {
		NewPrebuildAst(c, c.options.GetNumParallelThreads()).Prebuild(slices.Values(c.externs))
	}
	for _, input := range c.externs {
		n := rhino.CheckNotNull(input.GetAstRoot(c))
		if c.HasErrors() {
			return rhino.NilNode
		}
		c.externsRoot.AddChildToBack(n)
		c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), n)
	}

	if c.options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) ||
		c.options.GetProcessCommonJSModules() {
		c.InitializeModuleLoader()
	} else {
		// Use an empty module loader if we're not actually dealing with modules.
		c.moduleLoader = moduleLoader_NewEmpty()
	}

	if c.options.GetDependencyOptions().NeedsManagement() {
		c.findModulesFromEntryPoints(
			c.options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES),
			c.options.GetProcessCommonJSModules())
	} else if c.options.NeedsTranspilationFrom(parser.FeatureSet_ES2015_MODULES) ||
		c.options.GetProcessCommonJSModules() {
		if c.options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
			c.parsePotentialModules(c.chunkGraph.GetAllInputs())
		}

		// Build a map of module identifiers for any input which provides no namespace.
		// These files could be imported modules which have no exports, but do have side effects.
		inputModuleIdentifiers := jutil.NewLinkedHashMap[string, *CompilerInput]()
		for _, input := range c.chunkGraph.GetAllInputs() {
			if len(input.GetKnownProvides()) == 0 {
				modPath := c.moduleLoader.Resolve(input.GetSourceFile().GetName())
				inputModuleIdentifiers.Put(modPath.ToModuleName(), input)
			}
		}

		// Find out if any input attempted to import a module that had no exports.
		// In this case we must force module rewriting to occur on the imported file
		inputsToRewrite := jutil.NewLinkedHashMap[string, *CompilerInput]()
		for _, input := range c.chunkGraph.GetAllInputs() {
			for _, require := range input.GetKnownRequiredSymbols() {
				if inputModuleIdentifiers.ContainsKey(require) &&
					!inputsToRewrite.ContainsKey(require) {
					inputsToRewrite.Put(require, inputModuleIdentifiers.Get(require))
				}
			}
		}

		for input := range inputsToRewrite.Values() {
			input.SetJsModuleType(CompilerInput_ModuleType_IMPORTED_SCRIPT)
		}
	}

	if c.moduleLoader != nil {
		c.moduleLoader.SetErrorHandler(c.AsErrorHandler())
	}

	c.OrderInputs()
	c.markClosureUnawareCode()

	// If in IDE mode, we ignore the error and keep going.
	if c.HasErrors() {
		return rhino.NilNode
	}

	// Build the AST.
	allInputs := c.chunkGraph.GetAllInputs()
	if c.options.GetNumParallelThreads() > 1 {
		NewPrebuildAst(c, c.options.GetNumParallelThreads()).Prebuild(slices.Values(allInputs))
	}

	needSourceInfo := c.options.ShouldGatherSourceMapInfo() ||
		c.options.GetExternExportsPath() != nil ||
		len(c.options.GetReplaceStringsFunctionDescriptions()) != 0
	annotatedInParallel := false
	if needSourceInfo && !devMode && c.options.GetDevMode() == CompilerOptions_DevMode_OFF &&
		c.options.GetNumParallelThreads() > 1 && len(allInputs) > 1 && !c.HasErrors() {
		sia := SourceInformationAnnotator_Create()
		failures := parallelPrebuild(len(allInputs), c.options.GetNumParallelThreads(), func(i int) {
			if root := allInputs[i].GetAstRoot(c); !root.IsNil() {
				sia.AnnotateTree(root)
			}
		})
		for _, f := range failures {
			if f != nil {
				panic(f)
			}
		}
		annotatedInParallel = true
	}

	for _, input := range allInputs {
		n := rhino.CheckNotNull(input.GetAstRoot(c))
		if devMode {
			c.runValidityCheck()
			if c.HasErrors() {
				return rhino.NilNode
			}
		}

		if needSourceInfo && !annotatedInParallel {
			// Annotate the nodes in the tree with information from the
			// input file. This information is used to construct the SourceMap.
			if c.options.GetDevMode() == CompilerOptions_DevMode_OFF {
				SourceInformationAnnotator_Create().AnnotateTree(n)
			} else {
				sia := SourceInformationAnnotator_CreateWithAnnotationChecks(input.GetName())
				NodeTraversal_Traverse(c, n, sia)
			}
		}

		if NodeUtil_IsFromTypeSummary(n) {
			input.SetIsExtern()
			c.externsRoot.AddChildToBack(n)
		} else {
			c.jsRoot.AddChildToBack(n)
		}
		c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), n)
	}

	if c.HasErrors() {
		return rhino.NilNode
	}
	// Save on memory. Any future calls to "createInputConsideringTypedAstFilesystem" will throw.
	// TODO(lharker): do we actually need the synthetic externs file in stage 2?
	if c.typedAstFilesystem != nil {
		// If anything (besides the synthetic externs) is in the typedAstFilesystem, that indicates
		// it never should have been added in the first place: every CompilerInput should now be
		// initialized, a process which removes it from the typedAstFilesystem
		jutil.CheckStateMsg1(
			c.typedAstFilesystem.Size() == 1 || c.typedAstFilesystem.IsEmpty(),
			"%s",
			c.typedAstFilesystem.Size())
	}
	return c.externAndJsRoot
}

// PreferRegexParser hints to other parts of the compiler as to whether
// every input will eventually require a full AST parse.
//
// NOTE: this is just intended as a hint towards improving performance.
// Returning true or false shouldn't affect any non-performance behavior.
//
// Java: Compiler.preferRegexParser()
func (c *Compiler) PreferRegexParser() bool {
	return c.preferRegexParser
}

// SetPreferRegexParser artificially sets the return value of
// `preferRegexParser`
//
// Non-public and not exposed on AbstractCompiler, as for now we want to
// maintain control of when the regex vs. AST parser is used.
//
// Java: Compiler.setPreferRegexParser(boolean)
func (c *Compiler) SetPreferRegexParser(preferRegexParser bool) {
	c.preferRegexParser = preferRegexParser
}

// Java: Compiler.orderInputsWithLargeStack()
func (c *Compiler) OrderInputsWithLargeStack() {
	c.runInCompilerThread(func() {
		tracer := c.NewTracer("orderInputsWithLargeStack")
		defer c.StopTracer(tracer, "orderInputsWithLargeStack")
		c.OrderInputs()
	})
}

// Java: Compiler.orderInputs()
func (c *Compiler) OrderInputs() {
	c.maybeDoThreadedParsing()

	// Before dependency pruning, save a copy of the original inputs to use for externs hoisting.
	originalInputs := append([]*CompilerInput(nil), c.chunkGraph.GetAllInputs()...)

	if c.tracker != nil {
		c.tracker.RecordPrePruningInputCount(slices.Values(originalInputs), slices.Values(c.externs))
	}

	// Externs must be marked before dependency management since it needs to know what is an extern.
	c.markExterns(originalInputs)

	// Check if the sources need to be re-ordered.
	staleInputs := false
	if c.options.GetDependencyOptions().NeedsManagement() {
		jutil.Try(func() {
			jutil.Try(func() {
				result :=
					c.chunkGraph.ManageDependencies(c, c.options.GetDependencyOptions())
				staleInputs = true
				c.maybeLogPruningAnalysis(result)
			}, func(e *SortedDependencies_MissingProvideException) {
				c.Report(JSError_Make(Compiler_MISSING_ENTRY_ERROR, e.GetMessage()))
			})
		}, func(e *JSChunkGraph_MissingChunkException) {
			c.Report(JSError_Make(Compiler_MISSING_CHUNK_ERROR, e.GetMessage()))
		})
	}
	c.HoistExterns(originalInputs)

	// Manage dependencies may move weak sources around, and end up with empty chunks.
	c.fillEmptyChunks(c.GetChunks())
	c.HoistNoCompileFiles()

	if staleInputs {
		c.repartitionInputs()
	}
	c.maybeLogPrunedInputs()
}

// findModulesFromEntryPoints finds modules by recursively traversing
// dependencies starting with the entry points.
//
// Triggers either a regex parse or full AST-based parse of each file in
// order to gather dependency information.
//
// If the dependency mode is set to PRUNE_LEGACY, inputs which the regex
// parse does not identify as ES modules and which do not contain any
// provide statements are considered to be additional entry points.
//
// Java: Compiler.findModulesFromEntryPoints(boolean, boolean)
func (c *Compiler) findModulesFromEntryPoints(
	supportEs6Modules bool, supportCommonJSModules bool) {
	c.maybeDoThreadedParsing()
	entryPoints := []*CompilerInput{}
	inputsByProvide := jutil.NewLinkedHashMap[string, *CompilerInput]()
	inputsByIdentifier := jutil.NewLinkedHashMap[string, *CompilerInput]()
	for _, input := range c.chunkGraph.GetAllInputs() {
		var provides []string
		for _, p := range input.GetProvides() {
			if !strings.HasPrefix(p, "module$") {
				provides = append(provides, p)
			}
		}
		if !c.options.GetDependencyOptions().ShouldDropMoochers() && len(provides) == 0 {
			entryPoints = append(entryPoints, input)
		}
		inputsByIdentifier.Put(
			ModuleIdentifier_ForFile(input.GetPath().String()).String(), input)
		for _, provide := range provides {
			inputsByProvide.Put(provide, input)
		}
	}
	for _, moduleIdentifier := range c.options.GetDependencyOptions().EntryPoints() {
		input := inputsByProvide.Get(moduleIdentifier.String())
		if input == nil {
			input = inputsByIdentifier.Get(moduleIdentifier.String())
		}
		if input != nil {
			entryPoints = append(entryPoints, input)
		}
	}

	workingInputSet := jutil.NewLinkedHashSet[*CompilerInput]()
	for _, input := range c.chunkGraph.GetAllInputs() {
		workingInputSet.Add(input)
	}
	for _, entryPoint := range entryPoints {
		c.findModulesFromInput(
			entryPoint,
			/* wasImportedByModule= */ false,
			workingInputSet,
			inputsByIdentifier,
			inputsByProvide,
			supportEs6Modules,
			supportCommonJSModules)
	}
}

// findModulesFromInput traverses an input's dependencies to find
// additional modules.
//
// Java: Compiler.findModulesFromInput(CompilerInput, boolean, Set<CompilerInput>, Map<String, CompilerInput>, Map<String, CompilerInput>, boolean, boolean)
func (c *Compiler) findModulesFromInput(
	input *CompilerInput,
	wasImportedByModule bool,
	inputs *jutil.LinkedHashSet[*CompilerInput],
	inputsByIdentifier *jutil.LinkedHashMap[string, *CompilerInput],
	inputsByProvide *jutil.LinkedHashMap[string, *CompilerInput],
	supportEs6Modules bool,
	supportCommonJSModules bool) {
	if !inputs.Remove(input) {
		// It's possible for a module to be included as both a script
		// and a module in the same compilation. In these cases, it should
		// be forced to be a module.
		if wasImportedByModule && input.GetJsModuleType() == CompilerInput_ModuleType_NONE {
			input.SetJsModuleType(CompilerInput_ModuleType_IMPORTED_SCRIPT)
		}
		return
	}

	findDeps := NewFindModuleDependencies(
		c, supportEs6Modules, supportCommonJSModules, c.inputPathByWebpackId)
	findDeps.Process(rhino.CheckNotNull(input.GetAstRoot(c)))

	// If this input was imported by another module, it is itself a module
	// so we force it to be detected as such.
	if wasImportedByModule && input.GetJsModuleType() == CompilerInput_ModuleType_NONE {
		input.SetJsModuleType(CompilerInput_ModuleType_IMPORTED_SCRIPT)
	}
	c.moduleTypesByName.Put(input.GetPath().ToModuleName(), input.GetJsModuleType())

	var allDeps []string
	allDeps = append(allDeps, input.GetRequiredSymbols()...)
	allDeps = append(allDeps, input.GetDynamicRequires()...)
	allDeps = append(allDeps, input.GetTypeRequires()...)
	for _, requiredNamespace := range allDeps {
		var requiredInput *CompilerInput
		requiredByModuleImport := false
		if inputsByProvide.ContainsKey(requiredNamespace) {
			requiredInput = inputsByProvide.Get(requiredNamespace)
		} else if inputsByIdentifier.ContainsKey(requiredNamespace) {
			requiredByModuleImport = true
			requiredInput = inputsByIdentifier.Get(requiredNamespace)
		}

		if requiredInput != nil {
			c.findModulesFromInput(
				requiredInput,
				requiredByModuleImport,
				inputs,
				inputsByIdentifier,
				inputsByProvide,
				supportEs6Modules,
				supportCommonJSModules)
		}
	}
}

// HoistExterns hoists inputs with the @externs annotation into the externs
// list.
//
// Java: Compiler.hoistExterns(ImmutableList<CompilerInput>)
func (c *Compiler) HoistExterns(originalInputs []*CompilerInput) {
	staleInputs := false

	for _, input := range originalInputs {
		if c.hoistIfExtern(input) {
			staleInputs = true
		}
	}
	if staleInputs {
		c.repartitionInputs()
	}
}

// hoistIfExtern hoists a compiler input to externs if it contains the
// @externs annotation. Return whether or not the given input was hoisted.
//
// Java: Compiler.hoistIfExtern(CompilerInput)
func (c *Compiler) hoistIfExtern(input *CompilerInput) bool {
	if input.GetHasExternsAnnotation() {
		// If the input file is explicitly marked as an externs file, then move it out of the main
		// JS root and put it with the other externs.
		root := input.GetAstRoot(c)
		c.externsRoot.AddChildToBack(root)
		c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)

		chunk := input.GetChunk()
		if chunk != nil {
			chunk.Remove(input)
		}

		c.externs = append(c.externs, input)
		return true
	}
	return false
}

// markClosureUnawareCode marks source files annotated with
// `@closureUnaware` with the closureUnawareCode bit on SourceFile.
//
// Java: Compiler.markClosureUnawareCode()
func (c *Compiler) markClosureUnawareCode() {
	for _, input := range c.chunkGraph.GetAllInputs() {
		root := input.GetAstRoot(c)
		info := root.GetJSDocInfo()
		if info != nil && info.IsClosureUnawareCode() {
			input.GetSourceFile().MarkAsClosureUnawareCode()
		}
	}
}

// markExterns marks inputs with the @externs annotation as an Extern
// source file type. This is so that externs marking can be done before
// dependency management, and externs hoisting done after dependency
// management.
//
// Java: Compiler.markExterns(ImmutableList<CompilerInput>)
func (c *Compiler) markExterns(originalInputs []*CompilerInput) {
	for _, input := range originalInputs {
		if input.GetHasExternsAnnotation() {
			input.SetIsExtern()
		}
	}
}

// HoistNoCompileFiles hoists inputs with the @nocompile annotation out of
// the inputs.
//
// Java: Compiler.hoistNoCompileFiles()
func (c *Compiler) HoistNoCompileFiles() {
	staleInputs := false
	c.maybeDoThreadedParsing()
	// Iterate a copy because hoisting modifies what we're iterating over.
	for _, input := range append([]*CompilerInput(nil), c.chunkGraph.GetAllInputs()...) {
		if input.GetHasNoCompileAnnotation() {
			input.GetChunk().Remove(input)
			staleInputs = true
		}
	}

	if staleInputs {
		c.repartitionInputs()
	}
}

// Java: Compiler.maybeDoThreadedParsing()
func (c *Compiler) maybeDoThreadedParsing() {
	if c.options.GetNumParallelThreads() > 1 {
		NewPrebuildDependencyInfo(c.options.GetNumParallelThreads()).
			Prebuild(slices.Values(c.chunkGraph.GetAllInputs()))
	}
}

// Java: Compiler.repartitionInputs()
func (c *Compiler) repartitionInputs() {
	c.fillEmptyChunks(c.GetChunks())
	c.RebuildInputsFromChunks()
}

// ProcessJsonInputs transforms JSON files to a module export that closure
// compiler can process and keeps track of any "main" entries in
// package.json files.
//
// Java: Compiler.processJsonInputs(Iterable<CompilerInput>)
func (c *Compiler) ProcessJsonInputs(inputsToProcess []*CompilerInput) *jutil.LinkedHashMap[string, string] {
	rewriteJson := NewRewriteJsonToModule(c)
	for _, input := range inputsToProcess {
		if !strings.HasSuffix(input.GetSourceFile().GetName(), ".json") {
			continue
		}

		input.SetCompiler(c)
		readFailed := false
		jutil.Try(func() {
			// JSON objects need wrapped in parens to parse properly
			input.GetSourceFile().SetCodeDeprecated("(" + input.GetSourceFile().GetCode() + ")")
		}, func(e jutil.IIOException) {
			readFailed = true
		})
		if readFailed {
			continue
		}

		root := rhino.CheckNotNull(input.GetAstRoot(c))
		input.SetJsModuleType(CompilerInput_ModuleType_JSON)
		rewriteJson.Process(rhino.NilNode, root)
	}
	return rewriteJson.GetPackageJsonMainEntries()
}

// Java: Compiler.parsePotentialModules(Iterable<CompilerInput>)
func (c *Compiler) parsePotentialModules(inputsToProcess []*CompilerInput) {
	filteredInputs := []*CompilerInput{}
	for _, input := range inputsToProcess {
		// Only process files that are detected as ES6 modules
		if !c.options.GetDependencyOptions().ShouldPrune() ||
			!JsFileRegexParser_IsSupported() ||
			input.IsEs6Module() {
			filteredInputs = append(filteredInputs, input)
		}
	}
	if c.options.GetNumParallelThreads() > 1 {
		NewPrebuildAst(c, c.options.GetNumParallelThreads()).Prebuild(slices.Values(filteredInputs))
	}
	for _, input := range filteredInputs {
		input.SetCompiler(c)
		// Call getRequires to force regex-based dependency parsing to happen.
		input.GetRequires()
		input.SetJsModuleType(CompilerInput_ModuleType_ES6)
	}
}

// NewCompilerOptions allows subclasses to override the default
// CompileOptions object.
//
// Java: Compiler.newCompilerOptions()
func (c *Compiler) NewCompilerOptions() *CompilerOptions {
	return NewCompilerOptions()
}

// Java: Compiler.initCompilerOptionsIfTesting()
func (c *Compiler) InitCompilerOptionsIfTesting() {
	if c.options == nil {
		// initialization for tests that don't initialize the compiler
		// by the normal mechanisms.
		c.InitOptions(c.NewCompilerOptions())
		// allow tests to use unsupported features.
		c.options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	}
}

// Java: Compiler.parseSyntheticCode(String, String)
func (c *Compiler) ParseSyntheticCode(fileName string, js string) rhino.Node {
	c.InitCompilerOptionsIfTesting()
	source := SourceFile_FromCode(compiler_SYNTHETIC_FILE_NAME_PREFIX+fileName+"] ", js)
	c.addFilesToSourceMap([]*SourceFile{source})
	return c.parseCodeHelper(source)
}

// Java: Compiler.parseTestCode(String)
func (c *Compiler) ParseTestCode(js string) rhino.Node {
	c.InitCompilerOptionsIfTesting()
	c.InitBasedOnOptions()
	// NOTE: The file name used here is important, because at least Es6RewriteClassTest expects
	// it to match the default name used by the CompilerTestCase class for the source file
	// created when its `test("source string", "expected string")` method is called.
	// Also, some test cases rely on generated variable names that are based partially on the file
	// name or error messages that include the file name.
	return c.parseCodeHelper(SourceFile_FromCode("testcode", js))
}

// Java: Compiler.parseTestCode(ImmutableList<String>)
func (c *Compiler) ParseTestCodeOfList(code []string) rhino.Node {
	c.InitCompilerOptionsIfTesting()
	c.InitBasedOnOptions()
	srcs := make([]*SourceFile, len(code))
	for index, value := range code {
		srcs[index] = SourceFile_FromCode("testcode"+strconv.Itoa(index), value)
	}
	return c.parseCodeHelperOfList(srcs)
}

// Java: Compiler.parseCodeHelper(SourceFile)
func (c *Compiler) parseCodeHelper(src *SourceFile) rhino.Node {
	input := NewCompilerInput(src)
	c.putCompilerInput(input)
	root := input.GetAstRoot(c)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)
	return rhino.CheckNotNull(root)
}

// Java: Compiler.parseCodeHelper(List<SourceFile>)
func (c *Compiler) parseCodeHelperOfList(srcs []*SourceFile) rhino.Node {
	root := rhino.IR_Root()
	for _, src := range srcs {
		root.AddChildToBack(c.parseCodeHelper(src))
	}
	return root
}

// Java: Compiler.getDefaultErrorReporter()
func (c *Compiler) GetDefaultErrorReporter() rhino.ErrorReporter {
	return c.oldErrorReporter
}

// ------------------------------------------------------------------------
// Convert back to source code
// ------------------------------------------------------------------------

// ToSource converts the main parse tree back to JS code.
//
// Java: Compiler.toSource()
func (c *Compiler) ToSource() string {
	return Compiler_RunInCompilerThread(c,
		func() string {
			tracer := c.NewTracer("toSource")
			defer c.StopTracer(tracer, "toSource")
			cb := NewCompiler_CodeBuilder()
			// We are emitting all the sources at once, so use the SingleBinaryLicenseTracker
			// to de-dupe seen licenses across all inputs.
			lt := newCompiler_SingleBinaryLicenseTracker(c)
			if !c.jsRoot.IsNil() {
				var scripts []rhino.Node
				if c.options.ShouldPrintExterns() {
					for scriptNode := c.externsRoot.GetFirstChild(); !scriptNode.IsNil(); scriptNode = scriptNode.GetNext() {
						scripts = append(scripts, scriptNode)
					}
				}
				for scriptNode := c.jsRoot.GetFirstChild(); !scriptNode.IsNil(); scriptNode = scriptNode.GetNext() {
					scripts = append(scripts, scriptNode)
				}
				if c.options.GetNumParallelThreads() > 1 && len(scripts) > 1 {
					batchSize := 64
					if len(scripts) <= 64 {
						batchSize = max(c.options.GetNumParallelThreads(), 4)
					}
					// Go-only: fast mode prints up to 256 scripts per batch (better
					// balance between scripts of very different sizes) and resolves
					// the source-map mappings on the workers, leaving only the
					// generator's appends to the sequential merge. Same output.
					fastResolve := c.options.isFastMode() && c.options.ShouldGatherSourceMapInfo() && c.sourceMap != nil
					var resolvers []*SourceMap
					if c.options.isFastMode() {
						batchSize = min(len(scripts), 256)
					}
					if fastResolve {
						resolvers = make([]*SourceMap, c.options.GetNumParallelThreads())
					}
					type scriptPrinted struct {
						sm       CodePrinter_SourceAndMappings
						licenses []string
						lt       *Compiler_SingleBinaryLicenseTracker
						resolved []sourceMap_resolvedMapping
					}
					printed := make([]scriptPrinted, min(len(scripts), batchSize))
					numThreads := c.options.GetNumParallelThreads()
					scratches := make([]codePrinterScratch, numThreads)
					for batchStart := 0; batchStart < len(scripts); batchStart += batchSize {
						batchEnd := min(batchStart+batchSize, len(scripts))
						batchLen := batchEnd - batchStart
						failures := parallelPrebuildIndexed(batchLen, numThreads, func(w, j int) {
							idx := batchStart + j
							localLT := printed[j].lt
							if localLT == nil {
								localLT = newCompiler_SingleBinaryLicenseTracker(c)
								printed[j].lt = localLT
							} else {
								localLT.current.clear()
								localLT.lastSeenFile = ""
							}
							c.toSourceAndMappingsIntoWithScratch(scripts[idx], idx == 0, localLT, &printed[j].sm, &scratches[w])
							if fastResolve {
								if resolvers[w] == nil {
									resolvers[w] = c.newWorkerSourceMap()
								}
								rs := resolvers[w]
								out := printed[j].resolved[:0]
								for k := range printed[j].sm.mappings {
									m := &printed[j].sm.mappings[k]
									if r, ok := rs.resolveMappingOfNodePositions(m.node, m.start, m.end); ok {
										out = append(out, r)
									}
								}
								printed[j].resolved = out
							}
							if localLT.current.isEmpty() {
								printed[j].licenses = printed[j].licenses[:0]
							} else {
								printed[j].licenses = append(printed[j].licenses[:0], localLT.current.order...)
							}
						})
						for _, f := range failures {
							if f != nil {
								panic(f)
							}
						}
						for j := 0; j < batchLen; j++ {
							idx := batchStart + j
							for _, lic := range printed[j].licenses {
								lt.addLicense(lic)
							}
							printed[j].licenses = printed[j].licenses[:0]
							if fastResolve {
								c.appendSourceAndResolvedMappingsToCodeBuilder(cb, lt, idx, scripts[idx], &printed[j].sm, printed[j].resolved)
								printed[j].resolved = printed[j].resolved[:0]
								continue
							}
							c.appendSourceAndMappingsToCodeBuilder(cb, lt, idx, scripts[idx], &printed[j].sm)
						}
					}
				} else {
					for i, scriptNode := range scripts {
						c.ToSourceOfCodeBuilder(cb, lt, i, scriptNode)
					}
				}
			}
			return cb.String()
		})
}
