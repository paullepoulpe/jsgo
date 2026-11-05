package jscomp

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: single-restore multi-locale Stage 3 (FINALIZATIONS) fan-out.
//
// In multi-stage builds with late localization, Stage 1 (checks) and Stage 2
// (optimizations) produce a single locale-independent TypedAST + CompilerState
// bundle that Stage 3 (FINALIZATIONS) specializes across N target locales.
// RunFinalizationsBatch deserializes that bundle once, snapshots the post-restore
// NodeStore and Compiler state in memory, and executes Stage 3 + post-compilation
// tasks for each locale sequentially while rewinding in memory between locales.

// FinalizationLocaleSpec configures one locale pass in Compiler.RunFinalizationsBatch.
type FinalizationLocaleSpec struct {
	// Locale sets CompilerOptions.SetLocale when non-empty.
	Locale string

	// MessageBundle sets CompilerOptions.SetMessageBundle when non-nil.
	MessageBundle MessageBundle

	// ConfigureOptions is an optional per-locale hook invoked on a fresh clone
	// of the post-RestoreState CompilerOptions after Locale and MessageBundle
	// have been applied. Mutations made here never leak to subsequent locales.
	ConfigureOptions func(opts *CompilerOptions)
}

type finalizationInputCheckpoint struct {
	sourceFile             *SourceFile
	inputId                *rhino.InputId
	scriptNode             rhino.Node
	isSyntheticExterns     bool
	isSyntheticTypeSummary bool
}

type finalizationChunkCheckpoint struct {
	name       string
	depIndices []int
	inputs     []finalizationInputCheckpoint
}

type savedErrorEntry struct {
	entry SortingErrorManager_ErrorWithLevel
	node  rhino.Node
}

type sortingErrorManagerCheckpoint struct {
	present            bool
	messages           []savedErrorEntry
	sorted             bool
	originalErrorCount int
	promotedErrorCount int
	warningCount       int
	typedPercent       float64
}

type finalizationsCheckpoint struct {
	nodeStore  *rhino.NodeStore
	nodeSnap   *rhino.NodeStoreSnapshot
	stateProto *JSCompilerStateProto

	baseOptions *CompilerOptions

	externAndJsRoot rhino.Node
	externsRoot     rhino.Node
	jsRoot          rhino.Node

	externs []finalizationInputCheckpoint
	chunks  []finalizationChunkCheckpoint

	externProperties    *OrderedSet[string]
	colorRegistry       *colors.ColorRegistry
	lastInjectedLibrary rhino.Node

	inputSourceMaps     map[string]*SourceMapInput
	inputSourceMapOrder []string

	astCompactNested bool
	errorSnap        sortingErrorManagerCheckpoint
}

// RunFinalizationsBatch restores savedState once and runs Stage3Passes +
// PerformPostCompilationTasks for each entry in locales, invoking onLocaleDone
// after each locale completes so the caller can render single-file JS,
// multi-chunk JS + SourceMaps, variable/property maps, and diagnostics.
//
// Prior to calling RunFinalizationsBatch, the Compiler must have been
// initialized with its externs, inputs/chunks, and base CompilerOptions via
// Init or InitChunks (just as required before Compiler.RestoreState).
func (c *Compiler) RunFinalizationsBatch(
	savedState io.Reader,
	locales []FinalizationLocaleSpec,
	onLocaleDone func(idx int, spec FinalizationLocaleSpec, c *Compiler) error,
) error {
	if c == nil {
		return errors.New("jscomp: nil Compiler in RunFinalizationsBatch")
	}
	if savedState == nil {
		return errors.New("jscomp: nil savedState reader in RunFinalizationsBatch")
	}

	// Ensure this batch deserializes into an isolated, empty NodeStore rather than
	// snapshotting dead nodes accumulated by earlier compilations in the process.
	if !c.astCompact.nested && c.externAndJsRoot.IsNil() {
		if cur := rhino.CurrentNodeStore(); cur.IsStatic() || cur.NodeCount() > 0 {
			fresh := rhino.NewNodeStore()
			rhino.SetCurrentNodeStore(fresh)
			c.nodeStore = fresh
		}
	}

	c.RestoreState(savedState)
	if c.HasErrors() {
		return fmt.Errorf("jscomp: compiler has %d error(s) after RestoreState", c.GetErrorCount())
	}
	if len(locales) == 0 {
		return nil
	}

	cp := c.captureFinalizationsCheckpoint()

	for idx, spec := range locales {
		if idx == 0 {
			c.options = cloneCompilerOptions(cp.baseOptions)
		} else {
			c.rewindToFinalizationsCheckpoint(cp)
		}

		c.prepareLocaleOptions(spec)
		c.Stage3Passes()
		c.PerformPostCompilationTasks()

		if onLocaleDone != nil {
			if err := onLocaleDone(idx, spec, c); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Compiler) captureFinalizationsCheckpoint() *finalizationsCheckpoint {
	store := rhino.CurrentNodeStore()
	nodeSnap := store.Snapshot()
	stateProto := c.GetCompilerStateProtoBuilder().Build()
	baseOptions := cloneCompilerOptions(c.options)

	externs := make([]finalizationInputCheckpoint, len(c.externs))
	for i, ext := range c.externs {
		var script rhino.Node
		if ext.ast != nil && ext.ast.isParsed() {
			script = ext.ast.root
		}
		externs[i] = finalizationInputCheckpoint{
			sourceFile:             ext.GetSourceFile(),
			inputId:                ext.GetInputId(),
			scriptNode:             script,
			isSyntheticExterns:     ext == c.syntheticExternsInput,
			isSyntheticTypeSummary: ext == c.syntheticTypeSummaryInput,
		}
	}

	rawChunks := c.GetChunks()
	chunkIndex := make(map[*JSChunk]int, len(rawChunks))
	for i, ch := range rawChunks {
		chunkIndex[ch] = i
	}
	chunks := make([]finalizationChunkCheckpoint, len(rawChunks))
	for i, ch := range rawChunks {
		deps := ch.GetDependencies()
		depIndices := make([]int, len(deps))
		for j, dep := range deps {
			depIndices[j] = chunkIndex[dep]
		}
		rawInputs := ch.GetInputs()
		inputs := make([]finalizationInputCheckpoint, len(rawInputs))
		for j, in := range rawInputs {
			var script rhino.Node
			if in.ast != nil && in.ast.isParsed() {
				script = in.ast.root
			}
			inputs[j] = finalizationInputCheckpoint{
				sourceFile: in.GetSourceFile(),
				inputId:    in.GetInputId(),
				scriptNode: script,
			}
		}
		chunks[i] = finalizationChunkCheckpoint{
			name:       ch.GetName(),
			depIndices: depIndices,
			inputs:     inputs,
		}
	}

	var externProps *OrderedSet[string]
	if c.externProperties != nil {
		ep := c.externProperties.Clone()
		externProps = &ep
	}

	var lastInjected rhino.Node
	if c.runtimeJsLibManager != nil {
		lastInjected = c.runtimeJsLibManager.GetLastInjectedLibrary()
	}

	return &finalizationsCheckpoint{
		nodeStore:           store,
		nodeSnap:            nodeSnap,
		stateProto:          stateProto,
		baseOptions:         baseOptions,
		externAndJsRoot:     c.externAndJsRoot,
		externsRoot:         c.externsRoot,
		jsRoot:              c.jsRoot,
		externs:             externs,
		chunks:              chunks,
		externProperties:    externProps,
		colorRegistry:       c.colorRegistry,
		lastInjectedLibrary: lastInjected,
		inputSourceMaps:     maps.Clone(c.inputSourceMaps),
		inputSourceMapOrder: slices.Clone(c.inputSourceMapOrder),
		astCompactNested:    c.astCompact.nested,
		errorSnap:           captureSortingErrorManager(c.errorManager),
	}
}

func (c *Compiler) rewindToFinalizationsCheckpoint(cp *finalizationsCheckpoint) {
	cp.nodeStore.RestoreSnapshot(cp.nodeSnap)
	c.nodeStore = cp.nodeStore
	c.options = cloneCompilerOptions(cp.baseOptions)

	// Recreate ChangeTracker and RuntimeJsLibManager so pass timelines and
	// interned runtime field caches start completely fresh for this locale.
	c.changeTracker = NewChangeTracker()
	c.changeTracker.AddChangeHandler(c.changeTracker.GetRecentChange())

	var nodeForCodeInsertion func() rhino.Node
	if c.options.GetRuntimeLibraryMode() == RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES {
		nodeForCodeInsertion = c.getSynthesizedExternsRoot
	} else {
		nodeForCodeInsertion = c.getNodeForRuntimeCodeInsertion
	}
	c.runtimeJsLibManager = RuntimeJsLibManager_Create(
		c.options.GetRuntimeLibraryMode(),
		c.loadResourceContents,
		c.changeTracker,
		nodeForCodeInsertion,
	)

	// Restore scalar/collection compiler fields from the post-RestoreState proto.
	c.RestoreFromState(cp.stateProto)

	// Restore externProperties, ColorRegistry, AST roots, externs, and chunkGraph.
	if cp.externProperties != nil {
		ep := cp.externProperties.Clone()
		c.SetExternProperties(&ep)
	} else {
		c.SetExternProperties(nil)
	}
	c.colorRegistry = cp.colorRegistry
	c.externAndJsRoot = cp.externAndJsRoot
	c.externsRoot = cp.externsRoot
	c.jsRoot = cp.jsRoot

	clear(c.inputsById)
	clear(c.externs)
	c.externs = c.externs[:0]
	c.syntheticExternsInput = nil
	c.syntheticTypeSummaryInput = nil

	for _, extSpec := range cp.externs {
		in := newRestoredCompilerInput(c, extSpec, true)
		c.putCompilerInput(in)
		if !extSpec.scriptNode.IsNil() {
			c.scriptNodeByFilename[extSpec.sourceFile.GetName()] = extSpec.scriptNode
		}
		c.externs = append(c.externs, in)
		if extSpec.isSyntheticExterns {
			c.syntheticExternsInput = in
		}
		if extSpec.isSyntheticTypeSummary {
			c.syntheticTypeSummaryInput = in
		}
	}

	chunks := make([]*JSChunk, len(cp.chunks))
	for i, chSpec := range cp.chunks {
		chunks[i] = NewJSChunk(chSpec.name)
	}
	for i, chSpec := range cp.chunks {
		for _, depIdx := range chSpec.depIndices {
			chunks[i].AddDependency(chunks[depIdx])
		}
	}
	c.chunkGraph = NewJSChunkGraph(chunks)

	for i, chSpec := range cp.chunks {
		for _, inSpec := range chSpec.inputs {
			in := newRestoredCompilerInput(c, inSpec, false)
			if !inSpec.scriptNode.IsNil() {
				c.scriptNodeByFilename[inSpec.sourceFile.GetName()] = inSpec.scriptNode
			}
			c.putCompilerInput(in)
			chunks[i].Add(in)
		}
	}

	c.runtimeJsLibManager.SetLastInjectedLibrary(cp.lastInjectedLibrary)

	// Clear transient per-run caches and pass state.
	c.passes = nil
	c.phaseOptimizer = nil
	c.rucGuarded = nil
	c.propIndex = nil
	c.inlineFnCache = nil
	c.astCompact = astCompaction{nested: cp.astCompactNested}
	c.transpilationNamespace = nil
	c.cachedAstAnalyzer = nil
	c.moduleMetadataMap = nil
	c.moduleMap = nil
	clear(c.moduleTypesByName)
	c.defineNames.Clear()
	clear(c.commentsPerFile)
	c.lastJsSource = ""
	c.hasLastJsSource = false
	c.externExports = ""
	c.hasExternExports = false
	c.currentPassIndex = -1
	c.tracker = nil

	clear(c.sourceMapOriginalSources)
	c.resolvedSourceMap = compiler_ResolvedSourceMap{}
	clear(c.inputSourceMaps)
	if len(cp.inputSourceMaps) > 0 {
		maps.Copy(c.inputSourceMaps, cp.inputSourceMaps)
	}
	c.inputSourceMapOrder = slices.Clone(cp.inputSourceMapOrder)
	c.sourceMap = nil

	restoreSortingErrorManager(c.errorManager, cp.errorSnap)
}

func (c *Compiler) prepareLocaleOptions(spec FinalizationLocaleSpec) {
	if spec.Locale != "" {
		c.options.SetLocale(spec.Locale)
	}
	if spec.MessageBundle != nil {
		c.options.SetMessageBundle(spec.MessageBundle)
	}
	if spec.ConfigureOptions != nil {
		spec.ConfigureOptions(c.options)
	}

	c.passes = nil
	c.initASTCompaction()
	if c.runtimeJsLibManager != nil {
		c.runtimeJsLibManager.mode = c.options.GetRuntimeLibraryMode()
	}
	c.ReconcileOptionsWithGuards()
	c.initWarningsGuard(c.options.GetWarningsGuard())

	c.sourceMap = nil
	c.InitBasedOnOptions()
	c.ResetAndInitializeSourceMap()

	c.tracker = nil
	c.MaybeSetTracker()
	if c.tracker != nil {
		c.tracker.UpdateAfterDeserialize(c.jsRoot)
	}
}

func newRestoredCompilerInput(c *Compiler, spec finalizationInputCheckpoint, isExtern bool) *CompilerInput {
	in := NewCompilerInputWithInputIdAndIsExtern(spec.sourceFile, spec.inputId, isExtern)
	in.SetCompiler(c)
	if !spec.scriptNode.IsNil() {
		in.ast.installDeserializedRoot(spec.scriptNode)
	}
	return in
}

func captureSortingErrorManager(em ErrorManager) sortingErrorManagerCheckpoint {
	if em == nil {
		return sortingErrorManagerCheckpoint{}
	}
	if ts, ok := em.(*ThreadSafeDelegatingErrorManager); ok && ts != nil {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		em = ts.delegated
	}
	holder, ok := em.(interface{ sortingErrorManager() *SortingErrorManager })
	if !ok || holder == nil {
		return sortingErrorManagerCheckpoint{}
	}
	sm := holder.sortingErrorManager()
	if sm == nil {
		return sortingErrorManagerCheckpoint{}
	}
	msgs := make([]savedErrorEntry, len(sm.messages))
	for i, m := range sm.messages {
		var n rhino.Node
		if m.Error != nil {
			n = m.Error.node
		}
		msgs[i] = savedErrorEntry{entry: m, node: n}
	}
	return sortingErrorManagerCheckpoint{
		present:            true,
		messages:           msgs,
		sorted:             sm.sorted,
		originalErrorCount: sm.originalErrorCount,
		promotedErrorCount: sm.promotedErrorCount,
		warningCount:       sm.warningCount,
		typedPercent:       sm.typedPercent,
	}
}

func restoreSortingErrorManager(em ErrorManager, snap sortingErrorManagerCheckpoint) {
	if em == nil {
		return
	}
	if ts, ok := em.(*ThreadSafeDelegatingErrorManager); ok && ts != nil {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		em = ts.delegated
	}
	if snap.present {
		if holder, ok := em.(interface{ sortingErrorManager() *SortingErrorManager }); ok && holder != nil {
			if sm := holder.sortingErrorManager(); sm != nil {
				clear(sm.seen)
				if sm.seen == nil {
					sm.seen = make(map[sortingErrorKey]struct{}, len(snap.messages))
				}
				sm.messages = sm.messages[:0]
				for _, saved := range snap.messages {
					if saved.entry.Error != nil {
						saved.entry.Error.node = saved.node
					}
					sm.messages = append(sm.messages, saved.entry)
					sm.seen[makeSortingErrorKey(saved.entry.Error, saved.entry.Level)] = struct{}{}
				}
				sm.sorted = snap.sorted
				sm.originalErrorCount = snap.originalErrorCount
				sm.promotedErrorCount = snap.promotedErrorCount
				sm.warningCount = snap.warningCount
				sm.typedPercent = snap.typedPercent
				return
			}
		}
	}
	if resetter, ok := em.(interface{ Reset() }); ok {
		resetter.Reset()
	}
}

func cloneCompilerOptions(src *CompilerOptions) *CompilerOptions {
	if src == nil {
		return nil
	}
	dst := *src

	if src.experimentalOutputFeatureSet != nil {
		exp := *src.experimentalOutputFeatureSet
		dst.experimentalOutputFeatureSet = &exp
	}
	if src.dependencyOptions != nil {
		dep := *src.dependencyOptions
		dep.entryPoints = slices.Clone(src.dependencyOptions.entryPoints)
		dst.dependencyOptions = &dep
	}
	if src.extraAnnotationNames != nil {
		c := src.extraAnnotationNames.Clone()
		dst.extraAnnotationNames = &c
	}
	if src.inputSourceMaps != nil {
		c := src.inputSourceMaps.Clone()
		dst.inputSourceMaps = &c
	}
	if dng, ok := src.nameGenerator.(*DefaultNameGenerator); ok && dng != nil {
		dst.nameGenerator = cloneDefaultNameGenerator(dng)
	}
	dst.stripTypes = src.stripTypes.Clone()
	dst.stripNameSuffixes = src.stripNameSuffixes.Clone()
	dst.stripNamePrefixes = src.stripNamePrefixes.Clone()
	if src.customPasses != nil {
		groups := make([]customPassGroup, len(src.customPasses.groups))
		for i, g := range src.customPasses.groups {
			groups[i] = customPassGroup{
				time:   g.time,
				passes: slices.Clone(g.passes),
			}
		}
		dst.customPasses = &customPassMap{groups: groups}
	}
	dst.defineReplacements = src.defineReplacements.Clone()
	if src.cssRenamingSkiplist != nil {
		c := src.cssRenamingSkiplist.Clone()
		dst.cssRenamingSkiplist = &c
	}
	if src.idGenerators != nil {
		c := src.idGenerators.Clone()
		dst.idGenerators = &c
	}
	dst.replaceStringsFunctionDescriptions = slices.Clone(src.replaceStringsFunctionDescriptions)
	if src.propertiesThatMustDisambiguate != nil {
		c := src.propertiesThatMustDisambiguate.Clone()
		dst.propertiesThatMustDisambiguate = &c
	}
	dst.moduleRoots = slices.Clone(src.moduleRoots)
	if src.injectPolyfillsNewerThan != nil {
		lang := *src.injectPolyfillsNewerThan
		dst.injectPolyfillsNewerThan = &lang
	}
	dst.forceLibraryInjection = slices.Clone(src.forceLibraryInjection)
	if src.unusedImportsToRemove != nil {
		c := src.unusedImportsToRemove.Clone()
		dst.unusedImportsToRemove = &c
	}
	dst.filesToPrintAfterEachPassRegexList = slices.Clone(src.filesToPrintAfterEachPassRegexList)
	dst.chunksToPrintAfterEachPassRegexList = slices.Clone(src.chunksToPrintAfterEachPassRegexList)
	dst.qnameUsesToPrintAfterEachPassList = slices.Clone(src.qnameUsesToPrintAfterEachPassList)
	if src.warningsGuard != nil {
		dst.warningsGuard = &ComposeWarningsGuard{
			WarningsGuardBase: src.warningsGuard.WarningsGuardBase,
			entries:           slices.Clone(src.warningsGuard.entries),
			nextSeq:           src.warningsGuard.nextSeq,
			demoteErrors:      src.warningsGuard.demoteErrors,
		}
	}
	dst.extraReportGenerators = slices.Clone(src.extraReportGenerators)
	dst.sourceMapLocationMappings = slices.Clone(src.sourceMapLocationMappings)
	dst.conformanceConfigs = slices.Clone(src.conformanceConfigs)
	if src.browserResolverPrefixReplacements != nil {
		c := src.browserResolverPrefixReplacements.Clone()
		dst.browserResolverPrefixReplacements = &c
	}
	dst.packageJsonEntryNames = slices.Clone(src.packageJsonEntryNames)
	dst.unknownDefinesToIgnore = slices.Clone(src.unknownDefinesToIgnore)
	dst.extraOptions = maps.Clone(src.extraOptions)
	dst.optionWarnings = slices.Clone(src.optionWarnings)

	return &dst
}

func cloneDefaultNameGenerator(src *DefaultNameGenerator) *DefaultNameGenerator {
	dst := &DefaultNameGenerator{
		reservedNames:   src.reservedNames,
		reservedMap:     src.reservedMap,
		prefix:          src.prefix,
		nameCount:       src.nameCount,
		hasFavored:      src.hasFavored,
		firstAllowed:    src.firstAllowed,
		nonFirstAllowed: src.nonFirstAllowed,
		priorities:      src.priorities,
	}
	dst.firstChars = dst.firstBuf[:len(src.firstChars)]
	for i, cp := range src.firstChars {
		dst.firstChars[i] = &dst.priorities[cp.order]
	}
	dst.nonFirstChars = dst.nonFirstBuf[:len(src.nonFirstChars)]
	for i, cp := range src.nonFirstChars {
		dst.nonFirstChars[i] = &dst.priorities[cp.order]
	}
	return dst
}
