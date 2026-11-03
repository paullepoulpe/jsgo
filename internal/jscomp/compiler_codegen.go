// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261028

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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ToSourceOfJSChunk converts the parse tree for a chunk back to JS code.
//
// Consider using toSource(JSChunk, LicenseTracker) for better license
// handling. This call will emit all license text attached to all direct
// inputs to the chunk, which can be very inefficient.
//
// Java: Compiler.toSource(JSChunk)
func (c *Compiler) ToSourceOfJSChunk(chunk *JSChunk) string {
	return c.ToSourceOfLicenseTrackerAndJSChunk(
		// Backwards-compatible implementation of license tracking.
		NewCompiler_ScriptNodeLicensesOnlyTracker(c), chunk)
}

// ToSourceOfLicenseTrackerAndJSChunk converts the parse tree for a chunk
// back to JS code, using the given License Tracker to determine which
// licenses should be emitted before the source code.
//
// licenseTracker: The license tracker implementation to use.
// ChunkGraphAwareLicenseTracker is a suitable implementation when this
// method is being called on each chunk in the chunk graph in dependency
// order.
//
// Java: Compiler.toSource(LicenseTracker, JSChunk)
func (c *Compiler) ToSourceOfLicenseTrackerAndJSChunk(licenseTracker CodePrinter_LicenseTracker, chunk *JSChunk) string {
	return Compiler_RunInCompilerThread(c,
		func() string {
			inputs := chunk.GetInputs()
			numInputs := len(inputs)
			if numInputs == 0 {
				return ""
			}
			cb := NewCompiler_CodeBuilder()
			if numInputs > 1 && c.options != nil && c.options.GetNumParallelThreads() > 1 {
				replay, onlyScript, canParallel := compiler_licenseTrackerReplay(licenseTracker)
				if canParallel {
					scripts := make([]rhino.Node, numInputs)
					for i := 0; i < numInputs; i++ {
						scriptNode := inputs[i].GetAstRoot(c)
						if scriptNode.IsNil() {
							panic(jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + chunk.GetName()))
						}
						scripts[i] = scriptNode
					}
					batchSize := 64
					if numInputs <= 64 {
						batchSize = max(c.options.GetNumParallelThreads(), 4)
					}
					type scriptPrinted struct {
						sm  CodePrinter_SourceAndMappings
						rec compiler_ParallelScriptLicenseRecorder
					}
					printed := make([]scriptPrinted, min(numInputs, batchSize))
					numThreads := c.options.GetNumParallelThreads()
					scratches := make([]codePrinterScratch, numThreads)
					for batchStart := 0; batchStart < numInputs; batchStart += batchSize {
						batchEnd := min(batchStart+batchSize, numInputs)
						batchLen := batchEnd - batchStart
						failures := runParallelIndexedTasks(batchLen, numThreads, func(w, j int) {
							idx := batchStart + j
							printed[j].rec.Reset(onlyScript)
							c.toSourceAndMappingsIntoWithScratch(scripts[idx], idx == 0, &printed[j].rec, &printed[j].sm, &scratches[w])
						})
						for _, f := range failures {
							if f != nil {
								panic(f)
							}
						}
						for j := 0; j < batchLen; j++ {
							idx := batchStart + j
							printed[j].rec.Replay(replay)
							c.appendSourceAndMappingsToCodeBuilder(cb, licenseTracker, idx, scripts[idx], &printed[j].sm)
						}
					}
					return cb.String()
				}
			}
			for i := 0; i < numInputs; i++ {
				scriptNode := inputs[i].GetAstRoot(c)
				if scriptNode.IsNil() {
					panic(jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + chunk.GetName()))
				}
				c.ToSourceOfCodeBuilder(cb, licenseTracker, i, scriptNode)
			}
			return cb.String()
		})
}

// ToSourceOfCodeBuilder writes out JS code from a root node. If printing
// input delimiters, this method will attach a comment to the start of the
// text indicating which input the output derived from. If there were any
// preserve annotations within the root's source, they will also be printed
// in a block comment at the beginning of the output.
//
// The LicenseTracker provided determines which licenses attached to the
// nodes visited should be included in the output. When building a single
// JS bundle, consider using the SingleBinaryLicenseTracker implementation.
//
// Java: Compiler.toSource(CodeBuilder, LicenseTracker, int, Node)
func (c *Compiler) ToSourceOfCodeBuilder(
	cb *Compiler_CodeBuilder,
	licenseTracker CodePrinter_LicenseTracker,
	inputSeqNum int,
	root rhino.Node) {
	c.runInCompilerThread(func() {
		var sourceAndMappings CodePrinter_SourceAndMappings
		c.toSourceAndMappingsInto(root, inputSeqNum == 0, licenseTracker, &sourceAndMappings)
		c.appendSourceAndMappingsToCodeBuilder(cb, licenseTracker, inputSeqNum, root, &sourceAndMappings)
	})
}

func resetSourceAndMappings(sm *CodePrinter_SourceAndMappings) {
	if cap(sm.mappings) > 16384 {
		sm.mappings = nil
	} else {
		sm.mappings = sm.mappings[:0]
	}
}

func (c *Compiler) appendSourceAndMappingsToCodeBuilder(
	cb *Compiler_CodeBuilder,
	licenseTracker CodePrinter_LicenseTracker,
	inputSeqNum int,
	root rhino.Node,
	sourceAndMappings *CodePrinter_SourceAndMappings) {
	startLine, startCol, emitted := c.appendSourceAndMetricsToCodeBuilderForChunk(
		cb, licenseTracker, inputSeqNum, root, sourceAndMappings.source, sourceAndMappings)
	if !emitted {
		return
	}

	// If there is a valid source map, then indicate to it that the current
	// root node's mappings are offset by the given string builder buffer.
	// This offset is a result of licenses being added to the output buffer.
	if c.options.ShouldGatherSourceMapInfo() {
		c.sourceMap.SetStartingPosition(startLine, startCol)
		for i := range sourceAndMappings.mappings {
			m := &sourceAndMappings.mappings[i]
			c.sourceMap.AddMappingOfNodePositions(m.node, m.start, m.end)
		}
		resetSourceAndMappings(sourceAndMappings)
	}
}

// appendSourceAndResolvedMappingsToCodeBuilder is
// appendSourceAndMappingsToCodeBuilder with mappings already resolved by
// SourceMap.resolveMappingOfNodePositions (Go-only, fast mode).
func (c *Compiler) appendSourceAndResolvedMappingsToCodeBuilder(
	cb *Compiler_CodeBuilder,
	licenseTracker CodePrinter_LicenseTracker,
	inputSeqNum int,
	root rhino.Node,
	sourceAndMappings *CodePrinter_SourceAndMappings,
	resolved []sourceMap_resolvedMapping) {
	startLine, startCol, emitted := c.appendSourceAndMetricsToCodeBuilderForChunk(
		cb, licenseTracker, inputSeqNum, root, sourceAndMappings.source, sourceAndMappings)
	if emitted {
		c.sourceMap.SetStartingPosition(startLine, startCol)
		for i := range resolved {
			c.sourceMap.addResolvedMapping(&resolved[i])
		}
	}
	resetSourceAndMappings(sourceAndMappings)
}

func (c *Compiler) appendSourceToCodeBuilderForChunk(
	cb *Compiler_CodeBuilder,
	licenseTracker CodePrinter_LicenseTracker,
	inputSeqNum int,
	root rhino.Node,
	code string,
) (startLine int, startCol int, emitted bool) {
	return c.appendSourceAndMetricsToCodeBuilderForChunk(cb, licenseTracker, inputSeqNum, root, code, nil)
}

func (c *Compiler) appendSourceAndMetricsToCodeBuilderForChunk(
	cb *Compiler_CodeBuilder,
	licenseTracker CodePrinter_LicenseTracker,
	inputSeqNum int,
	root rhino.Node,
	code string,
	sm *CodePrinter_SourceAndMappings,
) (startLine int, startCol int, emitted bool) {
	if c.options.ShouldPrintInputDelimiter() {
		if (cb.GetLength() > 0) && !cb.EndsWith("\n") {
			cb.Append("\n") // Make sure that the label starts on a new line
		}
		jutil.CheckState(root.IsScript())

		delimiter := c.options.GetInputDelimiter()

		inputName := root.GetInputId().GetIdName()
		sourceName := root.GetSourceFileName()
		jutil.CheckState(sourceName != nil)
		jutil.CheckState(*sourceName != "")

		delimiter = strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(delimiter, "%name%", inputName),
				"%num%", strconv.Itoa(inputSeqNum)),
			"%n%", "\n")

		cb.Append(delimiter).Append("\n")
	}

	// Check whether there is any license information that should be emitted.
	for _, license := range licenseTracker.EmitLicenses() {
		cb.Append("/*\n").Append(license).Append("*/\n")
	}

	// Check whether there's any actual code to emit.
	if code == "" {
		// Nothing to do.
		return 0, 0, false
	}

	startLine = cb.GetLineIndex()
	startCol = cb.GetColumnIndex()

	if sm != nil && sm.hasMetrics {
		cb.appendWithMetrics(code, sm.utf16Len, sm.lineIndex, sm.lineLength)
	} else {
		cb.Append(code)
	}

	// In order to avoid parse ambiguity when files are concatenated
	// together, all files should end in a semi-colon. Do a quick
	// heuristic check if there's an obvious semi-colon already there.
	// ';' and '\n' are ASCII and never appear as trailing bytes of multi-byte UTF-8/WTF-8.
	hasSemiColon := strings.HasSuffix(code, ";") || strings.HasSuffix(code, ";\n")
	if !hasSemiColon {
		cb.Append(";")
	}
	return startLine, startCol, true
}

// ToSourceOfNode generates JavaScript source code for an AST, doesn't
// generate source map info.
//
// Java: Compiler.toSource(Node)
func (c *Compiler) ToSourceOfNode(n rhino.Node) string {
	c.InitCompilerOptionsIfTesting()
	if !n.IsNil() && !n.IsRoot() && !n.IsScript() {
		return c.toSourceWithoutMappings(n, false, nil)
	}
	// We don't know whether this node is part of a single JS binary, or part of a module, or
	// something else. Use the legacy license tracking implementation to match existing expectations
	// about what licenses are included in the resulting JS.
	lt := NewCompiler_ScriptNodeLicensesOnlyTracker(c)
	code := c.toSourceWithoutMappings(n, false, lt)
	licenses := lt.EmitLicenses()
	if len(licenses) == 0 {
		return code
	}
	var sb strings.Builder
	for _, license := range licenses {
		sb.WriteString("/*\n")
		sb.WriteString(license)
		sb.WriteString("*/\n")
	}
	sb.WriteString(code)
	return sb.String()
}

var emptyLicenseSet = &OrderedSet[string]{}

func resetOrderedSet[K comparable](s *OrderedSet[K]) {
	clear(s.m.entries)
	s.m.entries = s.m.entries[:0]
	clear(s.m.index)
	s.m.holes = 0
}

func addAllOrderedSet[K comparable](dst, src *OrderedSet[K]) {
	for k := range src.All() {
		dst.Add(k)
	}
}

// licenseTracker maintains two sets of licenses without virtual self-dispatch:
// the "currently seen set" emitted on the next call to EmitLicenses, and the
// "globally seen set" of every license ever emitted. When onlyScript is true,
// only licenses on SCRIPT/ROOT nodes are tracked (ScriptNodeLicensesOnlyTracker);
// otherwise only licenses on non-SCRIPT/non-ROOT nodes are tracked
// (SingleBinaryLicenseTracker).
type licenseTracker struct {
	compiler     AbstractCompiler
	globalSeen   map[string]struct{}
	current      OrderedSet[string]
	lastSeenFile string
	onlyScript   bool
}

// Compiler_ScriptNodeLicensesOnlyTracker emits all licenses attached to script
// nodes, regardless of whether the content of that script is used in the output.
//
// Java: Compiler.ScriptNodeLicensesOnlyTracker
type Compiler_ScriptNodeLicensesOnlyTracker = licenseTracker

// Java: Compiler.ScriptNodeLicensesOnlyTracker(AbstractCompiler)
func NewCompiler_ScriptNodeLicensesOnlyTracker(compiler AbstractCompiler) *Compiler_ScriptNodeLicensesOnlyTracker {
	return &licenseTracker{compiler: compiler, onlyScript: true}
}

// Compiler_SingleBinaryLicenseTracker retains only licenses for which there are
// useful non-SCRIPT/non-ROOT nodes in a binary.
//
// Java: Compiler.SingleBinaryLicenseTracker
type Compiler_SingleBinaryLicenseTracker = licenseTracker

// Java: Compiler.SingleBinaryLicenseTracker(AbstractCompiler)
func NewCompiler_SingleBinaryLicenseTracker(compiler AbstractCompiler) *Compiler_SingleBinaryLicenseTracker {
	return newCompiler_SingleBinaryLicenseTracker(compiler)
}

func newCompiler_SingleBinaryLicenseTracker(compiler AbstractCompiler) *Compiler_SingleBinaryLicenseTracker {
	return &licenseTracker{compiler: compiler}
}

func (t *licenseTracker) addLicense(lic string) {
	if _, ok := t.globalSeen[lic]; ok {
		return
	}
	t.current.Add(lic)
}

// Java: Compiler.SeenSetLicenseTracker.trackLicensesForNode(Node)
func (t *licenseTracker) TrackLicensesForNode(node rhino.Node) {
	isRootOrScript := node.IsRoot() || node.IsScript()
	if t.onlyScript != isRootOrScript {
		return
	}
	file := node.GetSourceFileName()
	if file == nil {
		return
	}
	t.trackLicensesForFileName(*file)
}

func (t *licenseTracker) trackLicensesForFileName(file string) {
	if file == t.lastSeenFile {
		return
	}
	t.lastSeenFile = file
	license, ok := compiler_getLicenseForFileName(t.compiler, file)
	if !ok {
		return
	}
	t.addLicense(license)
}

// Java: Compiler.SeenSetLicenseTracker.emitLicenses()
func (t *licenseTracker) EmitLicenses() []string {
	if t.current.IsEmpty() {
		return nil
	}
	rv := t.current.ToSlice()
	if t.globalSeen == nil {
		t.globalSeen = make(map[string]struct{}, len(rv))
	}
	for _, lic := range rv {
		t.globalSeen[lic] = struct{}{}
	}
	resetOrderedSet(&t.current)
	return rv
}

// Compiler_ChunkGraphAwareLicenseTracker is an implementation of a license
// tracker to be used when licenses can be de-duped among Chunks (because
// licenses might be loaded by a chunk's transitive dependencies. This
// tracker has a specific lifecycle to function correctly:
//
//   - all chunks must be processed in dependency-order
//   - setCurrentChunkContext must be called before any nodes in a Chunk are
//     processed
//   - emitLicenses can either be called once per chunk, or once per input
//     file in the chunk. Each call to emitLicenses will push all
//     newly-seen licenses into tracker-internal storage, and those licenses
//     won't be emitted again until setCurrentChunkContext resets this
//     internal state.
//
// Java: Compiler.ChunkGraphAwareLicenseTracker
type Compiler_ChunkGraphAwareLicenseTracker struct {
	compiler           AbstractCompiler
	licensesFromChunks map[*JSChunk]*OrderedSet[string]
	cumulativeLicenses map[*JSChunk]*OrderedSet[string]

	haveInitializedCurrentChunkLicenses bool
	currentChunkLicencesInTDeps         OrderedSet[string]
	currentChunkTDepsRef                *OrderedSet[string]

	lastSeenFile             string
	licensesNewInCurrentFile OrderedSet[string]

	// Re-initialized in SetCurrentChunkContext and stored into licensesFromChunks.
	licensesNewInCurrentChunk *OrderedSet[string]
	log                       diagnostic.LogFile
	currentChunk              *JSChunk
}

// Java: Compiler.ChunkGraphAwareLicenseTracker(AbstractCompiler)
func NewCompiler_ChunkGraphAwareLicenseTracker(compiler AbstractCompiler) *Compiler_ChunkGraphAwareLicenseTracker {
	t := &Compiler_ChunkGraphAwareLicenseTracker{
		licensesFromChunks:        make(map[*JSChunk]*OrderedSet[string], 16),
		cumulativeLicenses:        make(map[*JSChunk]*OrderedSet[string], 16),
		licensesNewInCurrentChunk: emptyLicenseSet,
		log:                       diagnostic.LogFile_CreateNoOp(),
		compiler:                  compiler,
	}
	t.currentChunkTDepsRef = &t.currentChunkLicencesInTDeps
	return t
}

// Java: Compiler.ChunkGraphAwareLicenseTracker.setLogFile(LogFile)
func (t *Compiler_ChunkGraphAwareLicenseTracker) SetLogFile(file diagnostic.LogFile) {
	t.log = file
}

// Java: Compiler.ChunkGraphAwareLicenseTracker.setCurrentChunkContext(JSChunk)
func (t *Compiler_ChunkGraphAwareLicenseTracker) SetCurrentChunkContext(chunk *JSChunk) {
	if t.log.IsLogging() {
		t.log.LogFormat("Initializing licenses from deps for chunk %s", chunk)
	}
	t.currentChunk = chunk
	if !t.currentChunkLicencesInTDeps.IsEmpty() {
		resetOrderedSet(&t.currentChunkLicencesInTDeps)
	}
	t.currentChunkTDepsRef = emptyLicenseSet
	if !t.licensesNewInCurrentFile.IsEmpty() {
		resetOrderedSet(&t.licensesNewInCurrentFile)
	}
	t.haveInitializedCurrentChunkLicenses = false
	if _, visited := t.licensesFromChunks[chunk]; visited {
		panic(jutil.NewIllegalStateExceptionMsg("Visiting a chunk more than once is not allowed."))
	}
	t.licensesNewInCurrentChunk = emptyLicenseSet
	t.licensesFromChunks[chunk] = emptyLicenseSet
}

func (t *Compiler_ChunkGraphAwareLicenseTracker) mergeDepLicenses(deps []*JSChunk, dst *OrderedSet[string]) (*OrderedSet[string], bool) {
	if len(deps) == 0 {
		return emptyLicenseSet, false
	}
	if len(deps) == 1 {
		return t.ensureCumulativeLicenses(deps[0]), false
	}
	var shared *OrderedSet[string]
	needMerge := false
	for _, dep := range deps {
		cumDep := t.ensureCumulativeLicenses(dep)
		if cumDep.IsEmpty() {
			continue
		}
		if shared == nil {
			shared = cumDep
		} else if shared != cumDep {
			needMerge = true
		}
	}
	if !needMerge {
		if shared == nil {
			return emptyLicenseSet, false
		}
		return shared, false
	}
	if dst == nil {
		dst = &OrderedSet[string]{}
	}
	for _, dep := range deps {
		if cumDep := t.cumulativeLicenses[dep]; !cumDep.IsEmpty() {
			addAllOrderedSet(dst, cumDep)
		}
	}
	return dst, true
}

func (t *Compiler_ChunkGraphAwareLicenseTracker) ensureCumulativeLicenses(chunk *JSChunk) *OrderedSet[string] {
	if cum, ok := t.cumulativeLicenses[chunk]; ok {
		return cum
	}
	chunkNew, visited := t.licensesFromChunks[chunk]
	if !visited {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Chunk aware license analysis error - analysis out of order."))
	}
	depsCum, merged := t.mergeDepLicenses(chunk.deps, nil)
	var res *OrderedSet[string]
	switch {
	case chunkNew.IsEmpty():
		res = depsCum
	case depsCum.IsEmpty():
		res = chunkNew
	case merged:
		addAllOrderedSet(depsCum, chunkNew)
		res = depsCum
	default:
		res = &OrderedSet[string]{}
		addAllOrderedSet(res, depsCum)
		addAllOrderedSet(res, chunkNew)
	}
	t.cumulativeLicenses[chunk] = res
	return res
}

// Java: Compiler.ChunkGraphAwareLicenseTracker.trackLicensesForNode(Node)
func (t *Compiler_ChunkGraphAwareLicenseTracker) TrackLicensesForNode(node rhino.Node) {
	if node.IsRoot() || node.IsScript() {
		// Neither of these node kinds on their own emit code, so skip their analysis.
		return
	}
	sourceFile := node.GetSourceFileName()
	if sourceFile == nil {
		return
	}
	t.trackLicensesForFileName(*sourceFile)
}

func (t *Compiler_ChunkGraphAwareLicenseTracker) trackLicensesForFileName(sourceFile string) {
	if t.lastSeenFile == sourceFile {
		// We just saw this file during the previous call to trackLicensesForNode
		return
	}
	t.lastSeenFile = sourceFile
	// Lazily initialize the set of licenses provided by the chunk's transitive dependencies
	if !t.haveInitializedCurrentChunkLicenses {
		if t.log.IsLogging() {
			for depChunk := range t.currentChunk.GetAllDependencies().All() {
				t.log.LogFormat("Chunk %s depends on chunk %s", t.currentChunk, depChunk)
				depSet, visited := t.licensesFromChunks[depChunk]
				if !visited {
					panic(jutil.NewIllegalStateExceptionMsg(
						"Chunk aware license analysis error - analysis out of order."))
				}
				if !depSet.IsEmpty() {
					addAllOrderedSet(&t.currentChunkLicencesInTDeps, depSet)
				}
			}
			t.currentChunkTDepsRef = &t.currentChunkLicencesInTDeps
		} else {
			t.currentChunkTDepsRef, _ = t.mergeDepLicenses(t.currentChunk.deps, &t.currentChunkLicencesInTDeps)
		}
		t.haveInitializedCurrentChunkLicenses = true
	}
	license, ok := compiler_getLicenseForFileName(t.compiler, sourceFile)
	if !ok {
		return
	}
	if t.currentChunkTDepsRef.Contains(license) {
		return
	}
	if t.licensesNewInCurrentChunk.Contains(license) {
		return
	}
	newlyAddedLicense := t.licensesNewInCurrentFile.Add(license)
	if t.log.IsLogging() {
		if !newlyAddedLicense {
			t.log.LogFormat(
				"Chunk %s already depends on license\n==================\n%s\n==================\n\n",
				t.currentChunk, license)
		} else {
			t.log.LogFormat(
				"Chunk %s has new license:\n==================\n%s\n==================\n\n",
				t.currentChunk, license)
		}
	}
}

type compiler_ParallelScriptLicenseRecorder struct {
	onlyScript   bool
	lastSeenFile string
	hasLastSeen  bool
	seenFiles    []string
}

func (r *compiler_ParallelScriptLicenseRecorder) Reset(onlyScript bool) {
	r.onlyScript = onlyScript
	r.lastSeenFile = ""
	r.hasLastSeen = false
	r.seenFiles = r.seenFiles[:0]
}

func (r *compiler_ParallelScriptLicenseRecorder) Replay(replay func(string)) {
	for _, file := range r.seenFiles {
		replay(file)
	}
}

func (r *compiler_ParallelScriptLicenseRecorder) TrackLicensesForNode(node rhino.Node) {
	isRootOrScript := node.IsRoot() || node.IsScript()
	if r.onlyScript != isRootOrScript {
		return
	}
	file := node.GetSourceFileName()
	if file == nil {
		return
	}
	if r.hasLastSeen && *file == r.lastSeenFile {
		return
	}
	r.lastSeenFile = *file
	r.hasLastSeen = true
	r.seenFiles = append(r.seenFiles, *file)
}

func (r *compiler_ParallelScriptLicenseRecorder) EmitLicenses() []string {
	return nil
}

func compiler_licenseTrackerReplay(lt CodePrinter_LicenseTracker) (replay func(string), onlyScript bool, ok bool) {
	switch t := lt.(type) {
	case *Compiler_ChunkGraphAwareLicenseTracker:
		return t.trackLicensesForFileName, false, true
	case *licenseTracker:
		return t.trackLicensesForFileName, t.onlyScript, true
	default:
		return nil, false, false
	}
}

// Java: Compiler.ChunkGraphAwareLicenseTracker.emitLicenses()
func (t *Compiler_ChunkGraphAwareLicenseTracker) EmitLicenses() []string {
	if t.licensesNewInCurrentFile.IsEmpty() {
		return nil
	}
	rv := t.licensesNewInCurrentFile.ToSlice()
	if t.licensesNewInCurrentChunk == emptyLicenseSet {
		t.licensesNewInCurrentChunk = NewOrderedSet(rv...)
		t.licensesFromChunks[t.currentChunk] = t.licensesNewInCurrentChunk
	} else {
		t.licensesNewInCurrentChunk.AddAll(rv...)
	}
	resetOrderedSet(&t.licensesNewInCurrentFile)
	return rv
}

func (c *Compiler) acquirePrintScratch() (*codePrinterScratch, bool) {
	s := c.ensureSync()
	if s.printScratchBusy.CompareAndSwap(false, true) {
		return &s.printScratch, true
	}
	return nil, false
}

func (c *Compiler) toSourceAndMappingsInto(
	n rhino.Node, firstOutput bool, licenseTracker CodePrinter_LicenseTracker, sm *CodePrinter_SourceAndMappings) {
	scratch, owned := c.acquirePrintScratch()
	if owned {
		defer c.sync.printScratchBusy.Store(false)
	}
	c.toSourceAndMappingsIntoWithScratch(n, firstOutput, licenseTracker, sm, scratch)
}

func (c *Compiler) printRootWithScratch(
	n rhino.Node,
	firstOutput bool,
	gatherSourceMap bool,
	licenseTracker CodePrinter_LicenseTracker,
	sm *CodePrinter_SourceAndMappings,
	scratch *codePrinterScratch,
) string {
	builder := CodePrinter_Builder{root: n}
	builder.SetCompilerOptions(c.options)
	builder.SetTagAsTypeSummary(c.options.ShouldGenerateTypedExterns())
	builder.SetTagAsStrict(firstOutput && c.options.ShouldEmitUseStrict())
	builder.SetLicenseTracker(licenseTracker)
	return builder.buildSourceAndMappingsInto(sm, gatherSourceMap, scratch)
}

func (c *Compiler) toSourceAndMappingsIntoWithScratch(
	n rhino.Node,
	firstOutput bool,
	licenseTracker CodePrinter_LicenseTracker,
	sm *CodePrinter_SourceAndMappings,
	scratch *codePrinterScratch) {
	c.printRootWithScratch(n, firstOutput, true, licenseTracker, sm, scratch)
}

func (c *Compiler) toSourceWithoutMappings(
	n rhino.Node, firstOutput bool, licenseTracker CodePrinter_LicenseTracker) string {
	scratch, owned := c.acquirePrintScratch()
	if owned {
		defer c.sync.printScratchBusy.Store(false)
	}
	return c.printRootWithScratch(n, firstOutput, false, licenseTracker, nil, scratch)
}

func compiler_getLicenseForFileName(compiler AbstractCompiler, fileName string) (string, bool) {
	licenseRoot := compiler.GetScriptNode(fileName)
	if licenseRoot.IsNil() {
		return "", false
	}
	docInfo := licenseRoot.GetJSDocInfo()
	if docInfo == nil {
		return "", false
	}
	lic := docInfo.GetLicense()
	if lic == nil {
		return "", false
	}
	return *lic, true
}

// Compiler_CodeBuilder stores a buffer of text to which more can be
// appended. This is just like a StringBuilder except that we also track the
// number of lines.
//
// Java: Compiler.CodeBuilder
type Compiler_CodeBuilder struct {
	sb        strings.Builder
	utf16Len  int
	lineCount int
	colCount  int
}

// Java: Compiler.CodeBuilder()
func NewCompiler_CodeBuilder() *Compiler_CodeBuilder {
	return &Compiler_CodeBuilder{}
}

// Reset removes all text, but leaves the line count unchanged.
//
// Java: Compiler.CodeBuilder.reset()
func (b *Compiler_CodeBuilder) Reset() {
	b.sb.Reset()
	b.utf16Len = 0
}

// ResetAll removes all text and resets line and column counts to zero.
func (b *Compiler_CodeBuilder) ResetAll() {
	b.sb.Reset()
	b.utf16Len = 0
	b.lineCount = 0
	b.colCount = 0
}

func (b *Compiler_CodeBuilder) appendWithMetrics(str string, utf16Len, newlines, colLen int) {
	if len(str) == 0 {
		return
	}
	b.sb.WriteString(str)
	b.utf16Len += utf16Len
	if newlines == 0 {
		b.colCount += colLen
	} else {
		b.lineCount += newlines
		b.colCount = colLen
	}
}

// Append appends the given string to the text buffer.
//
// Java: Compiler.CodeBuilder.append(String)
func (b *Compiler_CodeBuilder) Append(str string) *Compiler_CodeBuilder {
	if len(str) == 0 {
		return b
	}
	b.sb.WriteString(str)

	// Adjust the line and column information for the new text.
	//
	// Go-only: performance — Java loops over str.indexOf('\n', index + 1),
	// which here would convert all of str to UTF-16 once per line (quadratic
	// in the size of a script). '\n' never occurs inside a multi-byte WTF-8
	// sequence, so counting and finding it on the bytes gives the same line
	// count, and the UTF-16 length after the last '\n' is
	// str.length() - (lastIndex + 1).
	lastIndex := strings.LastIndexByte(str, '\n')
	if jutil.IsASCII(str) {
		b.utf16Len += len(str)
		if lastIndex == -1 {
			b.colCount += len(str)
		} else {
			b.lineCount += 1 + strings.Count(str[:lastIndex], "\n")
			b.colCount = len(str) - (lastIndex + 1)
		}
		return b
	}
	if lastIndex == -1 {
		// No new lines, append the new characters added.
		n := jutil.StringLength(str)
		b.utf16Len += n
		b.colCount += n
	} else {
		tailLen := jutil.StringLength(str[lastIndex+1:])
		b.utf16Len += jutil.StringLength(str[:lastIndex]) + 1 + tailLen
		b.lineCount += 1 + strings.Count(str[:lastIndex], "\n")
		b.colCount = tailLen
	}

	return b
}

// String returns all text in the text buffer.
//
// Java: Compiler.CodeBuilder.toString()
func (b *Compiler_CodeBuilder) String() string {
	return b.sb.String()
}

// GetLength returns the length of the text buffer.
//
// Java: Compiler.CodeBuilder.getLength()
func (b *Compiler_CodeBuilder) GetLength() int {
	return b.utf16Len
}

// GetLineIndex returns the (zero-based) index of the last line in the text
// buffer.
//
// Java: Compiler.CodeBuilder.getLineIndex()
func (b *Compiler_CodeBuilder) GetLineIndex() int {
	return b.lineCount
}

// GetColumnIndex returns the (zero-based) index of the last column in the
// text buffer.
//
// Java: Compiler.CodeBuilder.getColumnIndex()
func (b *Compiler_CodeBuilder) GetColumnIndex() int {
	return b.colCount
}

// EndsWith determines whether the text ends with the given suffix.
//
// Java: Compiler.CodeBuilder.endsWith(String)
func (b *Compiler_CodeBuilder) EndsWith(suffix string) bool {
	return b.utf16Len > jutil.StringLength(suffix) && strings.HasSuffix(b.sb.String(), suffix)
}

// ------------------------------------------------------------------------
// Optimizations
// ------------------------------------------------------------------------

// Java: Compiler.performTranspilationAndOptimizations(SegmentOfCompilationToRun)
func (c *Compiler) PerformTranspilationAndOptimizations(segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun) {
	jutil.CheckState(c.options.ShouldOptimize())
	// getOptimizations() also includes transpilation passes
	allOptimizationPassesToRun := c.GetPassConfig().GetOptimizations().Build()
	if len(allOptimizationPassesToRun) == 0 {
		return
	}

	// Get the subset of optimization passes to run in the current segment of optimizations.
	optimizationPassesToRunInCurrentSegment := c.getOptimizationPassesToRunInCurrentSegment(
		segmentOfCompilationToRun, allOptimizationPassesToRun)

	// note whether any script gets transpiled
	c.markTranspiledFiles()

	c.phaseOptimizer = c.createPhaseOptimizer()
	c.phaseOptimizer.Consume(optimizationPassesToRunInCurrentSegment)
	c.beginASTCompaction()
	c.phaseOptimizer.Process(c.externsRoot, c.jsRoot)
	c.maybeCompactAST("end of optimizations", nil, nil, nil)
	c.astCompact.allowed = false
	c.phaseOptimizer = nil
}
