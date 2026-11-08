package jscomp

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// RenderedChunk holds the rendered JavaScript source and optional source map
// (and optional worker-pool transformed/compressed source map bytes) for a
// single JSChunk.
type RenderedChunk struct {
	Chunk          *JSChunk
	Code           string
	SourceMap      string
	SourceMapBytes []byte
}

// RenderChunkOptions configures optional license tracking, source-map naming,
// wrapper prefixing, per-input post-processing, source-map compression, and
// worker parallelism for Compiler.RenderChunksAndSourceMaps.
type RenderChunkOptions struct {
	// LicenseTracker tracks emitted licenses across chunks in dependency order.
	// When nil, NewCompiler_ChunkGraphAwareLicenseTracker(c) is used unless
	// PerChunkScriptLicenses is true.
	LicenseTracker CodePrinter_LicenseTracker

	// PerChunkScriptLicenses, when true and LicenseTracker is nil, uses a fresh
	// ScriptNodeLicensesOnlyTracker per chunk (deduplicating script-level
	// licenses within each chunk, not across chunks), matching
	// AbstractCommandLineRunner.createJsonFileFromChunk.
	PerChunkScriptLicenses bool

	// GenerateSourceMaps requests per-chunk source map rendering even when
	// CompilerOptions.SetSourceMapOutputPath was not set. When false, source maps
	// are rendered whenever c.GetOptions().ShouldGatherSourceMapInfo() is true.
	GenerateSourceMaps bool

	// ChunkFileName returns the output filename recorded in each chunk's source
	// map ("file" field). When nil, chunk.GetName() + ".js" is used.
	ChunkFileName func(*JSChunk) string

	// WrapperPrefix returns an optional prefix prepended before the chunk's JS
	// code so the chunk's source map offsets are adjusted accordingly.
	WrapperPrefix func(*JSChunk) string

	// PostProcessInput, when non-nil, transforms each printed CompilerInput's
	// generated JavaScript code in parallel before chunk concatenation.
	PostProcessInput func(input *CompilerInput, code string) string

	// CompressSourceMap, when non-nil and source maps are generated, transforms
	// or compresses each chunk's rendered source map JSON on the worker pool and
	// stores the result in RenderedChunk.SourceMapBytes.
	CompressSourceMap func(chunk *JSChunk, sourceMap string) ([]byte, error)

	// CompressSourceMapBytes, when non-nil and source maps are generated, is a
	// zero-copy alternative to CompressSourceMap that receives a worker-owned
	// []byte buffer valid for the duration of the callback and stores the result
	// in RenderedChunk.SourceMapBytes.
	CompressSourceMapBytes func(chunk *JSChunk, sourceMap []byte) ([]byte, error)

	// OmitSourceMapString, when true alongside CompressSourceMap or
	// CompressSourceMapBytes, omits populating RenderedChunk.SourceMap so
	// multi-thousand-chunk builds do not retain uncompressed source map JSON
	// strings in memory.
	OmitSourceMapString bool

	// NumWorkers overrides the worker pool size. When <= 0, uses
	// c.GetOptions().GetNumParallelThreads() if > 1, or runtime.GOMAXPROCS(0).
	NumWorkers int
}

// ParallelChunkOutputCallbacks supplies the CLI-specific path expansion,
// wrapper lookup, and stream writing hooks to Compiler.OutputChunksParallel.
type ParallelChunkOutputCallbacks struct {
	ChunkFileName     func(m *JSChunk) string
	SourceMapPath     func(m *JSChunk) string
	WrapperTemplate   func(chunkName string) string
	EnsureDir         func(filename string)
	WriteChunkWithMap func(filename string, mapFilename string, sm *SourceMap, code string, wrapper string)
	WriteChunkNoMap   func(filename string, code string, wrapper string)
}

func (c *Compiler) ensureSourceMapInitialized() {
	if c.sourceMap != nil {
		return
	}
	c.sourceMap = c.options.GetSourceMapFormat().GetInstance()
	c.sourceMap.SetPrefixMappings(c.options.GetSourceMapLocationMappings())
	if c.options.GetApplyInputSourceMaps() {
		c.sourceMap.SetCompilerSourceMapping(c)
	}
}

func filterActiveChunks(chunks []*JSChunk) []*JSChunk {
	active := make([]*JSChunk, 0, len(chunks))
	for _, m := range chunks {
		if m != nil && m.GetName() != JSChunk_WEAK_CHUNK_NAME {
			active = append(active, m)
		}
	}
	return active
}

func populateChunkSourceMap(sm *SourceMap, scripts []parallelPrintedScript) {
	sm.Reset()
	for i := range scripts {
		addPrintedScriptMappings(sm, &scripts[i])
	}
}

func prepareChunkSourceMapOutput(
	sm *SourceMap,
	m *JSChunk,
	wrapperPrefix func(*JSChunk) string,
	chunkFileName func(*JSChunk) string,
) string {
	if wrapperPrefix != nil {
		if prefix := wrapperPrefix(m); prefix != "" {
			sm.SetWrapperPrefix(prefix)
		}
	}
	return chunkFileName(m)
}

func renderChunkSourceMapString(
	sm *SourceMap,
	m *JSChunk,
	wrapperPrefix func(*JSChunk) string,
	chunkFileName func(*JSChunk) string,
) string {
	name := prepareChunkSourceMapOutput(sm, m, wrapperPrefix, chunkFileName)
	var mapBuf strings.Builder
	if est := sm.EstimatedByteSize(); est > 0 {
		mapBuf.Grow(est)
	}
	sm.AppendTo(&mapBuf, name)
	return mapBuf.String()
}

func renderChunkSourceMapBytesInto(
	buf *bytes.Buffer,
	sm *SourceMap,
	m *JSChunk,
	wrapperPrefix func(*JSChunk) string,
	chunkFileName func(*JSChunk) string,
) []byte {
	name := prepareChunkSourceMapOutput(sm, m, wrapperPrefix, chunkFileName)
	buf.Reset()
	if est := sm.EstimatedByteSize(); est > buf.Cap() {
		buf.Grow(est - buf.Len())
	}
	sm.AppendTo(buf, name)
	return buf.Bytes()
}

const (
	chunkRenderMaxBatchChunks  = 64
	chunkRenderMaxBatchScripts = 256
)

func (c *Compiler) renderChunksBatched(
	activeChunks []*JSChunk,
	numWorkers int,
	lt CodePrinter_LicenseTracker,
	onlyScript bool,
	replay func(string),
	resetLTPerChunk func(),
	gatherSourceMap bool,
	postProcessInput func(*CompilerInput, string) string,
	onSeqChunk func(globalIdx, batchIdx int, m *JSChunk, code string),
	onParChunk func(workerIdx, globalIdx, batchIdx int, m *JSChunk, sm *SourceMap),
) error {
	chunkTracker, _ := lt.(*Compiler_ChunkGraphAwareLicenseTracker)
	totalInputs := 0
	for _, m := range activeChunks {
		totalInputs += len(c.chunkInputsForOutput(m))
	}
	scriptOffsets := make([]int, chunkRenderMaxBatchChunks+1)
	printer := newParallelScriptPrinter(
		c,
		numWorkers,
		onlyScript,
		gatherSourceMap,
		min(max(totalInputs, 1), chunkRenderMaxBatchScripts),
	)
	printer.postProcessInput = postProcessInput

	var workerSM []*SourceMap
	if gatherSourceMap && onParChunk != nil {
		workerSM = make([]*SourceMap, max(1, min(numWorkers, min(len(activeChunks), chunkRenderMaxBatchChunks))))
		for w := range workerSM {
			workerSM[w] = c.newWorkerSourceMap()
		}
	}

	cb := NewCompiler_CodeBuilder()
	for chunkIdx := 0; chunkIdx < len(activeChunks); {
		firstChunk := activeChunks[chunkIdx]
		if inputs := c.chunkInputsForOutput(firstChunk); len(inputs) > chunkRenderMaxBatchScripts {
			// Stream large chunks (such as a multi-thousand-input base chunk) in
			// bounded script sub-batches so per-script codeBytes and mapping slices
			// never accumulate beyond chunkRenderMaxBatchScripts at once.
			for _, in := range inputs {
				if in.GetAstRoot(c).IsNil() {
					return jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + firstChunk.GetName())
				}
			}
			if resetLTPerChunk != nil {
				resetLTPerChunk()
			}
			if chunkTracker != nil {
				chunkTracker.SetCurrentChunkContext(firstChunk)
			}
			cb.ResetAll()
			var sm *SourceMap
			if gatherSourceMap && onParChunk != nil {
				sm = workerSM[0]
				if chunkIdx == len(activeChunks)-1 && c.sourceMap != nil {
					sm = c.sourceMap
				}
				sm.Reset()
			}
			for inStart := 0; inStart < len(inputs); inStart += chunkRenderMaxBatchScripts {
				subLen := min(chunkRenderMaxBatchScripts, len(inputs)-inStart)
				slots := printer.ensureSlots(subLen)
				for k := range subLen {
					idx := inStart + k
					if err := printer.setSlot(k, firstChunk.GetName(), inputs[idx], idx == 0); err != nil {
						return err
					}
				}
				if err := printer.printBatch(subLen); err != nil {
					return err
				}
				for k := range subLen {
					c.appendPrintedScriptToCodeBuilder(cb, lt, replay, inStart+k, &slots[k])
					if sm != nil {
						addPrintedScriptMappings(sm, &slots[k])
					} else {
						resetSourceAndMappings(&slots[k].sm)
					}
				}
			}
			code := cb.String()
			onSeqChunk(chunkIdx, 0, firstChunk, code)
			if onParChunk != nil {
				parFailures := runParallelIndexedTasks(1, 1, func(w, _ int) {
					onParChunk(w, chunkIdx, 0, firstChunk, sm)
				})
				if len(parFailures) > 0 && parFailures[0] != nil {
					return parFailures[0]
				}
			}
			chunkIdx++
			continue
		}

		batchStart := chunkIdx
		batchLen := 0
		totalScripts := 0
		for batchStart+batchLen < len(activeChunks) && batchLen < chunkRenderMaxBatchChunks {
			nIn := len(c.chunkInputsForOutput(activeChunks[batchStart+batchLen]))
			if batchLen > 0 && totalScripts+nIn > chunkRenderMaxBatchScripts {
				break
			}
			scriptOffsets[batchLen] = totalScripts
			totalScripts += nIn
			batchLen++
		}
		scriptOffsets[batchLen] = totalScripts
		chunkIdx += batchLen

		slots := printer.ensureSlots(totalScripts)
		for j := range batchLen {
			m := activeChunks[batchStart+j]
			base := scriptOffsets[j]
			for i, in := range c.chunkInputsForOutput(m) {
				if err := printer.setSlot(base+i, m.GetName(), in, i == 0); err != nil {
					return err
				}
			}
		}
		if err := printer.printBatch(totalScripts); err != nil {
			return err
		}

		for j := range batchLen {
			globalIdx := batchStart + j
			m := activeChunks[globalIdx]
			if resetLTPerChunk != nil {
				resetLTPerChunk()
			}
			if chunkTracker != nil {
				chunkTracker.SetCurrentChunkContext(m)
			}
			chunkSlots := slots[scriptOffsets[j]:scriptOffsets[j+1]]
			var code string
			if len(chunkSlots) > 0 {
				cb.ResetAll()
				for i := range chunkSlots {
					c.appendPrintedScriptToCodeBuilder(cb, lt, replay, i, &chunkSlots[i])
				}
				code = cb.String()
			}
			onSeqChunk(globalIdx, j, m, code)
		}

		if onParChunk != nil {
			parFailures := runParallelIndexedTasks(batchLen, min(numWorkers, batchLen), func(w, j int) {
				globalIdx := batchStart + j
				m := activeChunks[globalIdx]
				var sm *SourceMap
				if gatherSourceMap {
					sm = workerSM[w]
					if globalIdx == len(activeChunks)-1 && c.sourceMap != nil {
						sm = c.sourceMap
					}
					populateChunkSourceMap(sm, slots[scriptOffsets[j]:scriptOffsets[j+1]])
				}
				onParChunk(w, globalIdx, j, m, sm)
			})
			for _, f := range parFailures {
				if f != nil {
					return f
				}
			}
		} else {
			for k := range totalScripts {
				resetSourceAndMappings(&slots[k].sm)
			}
		}
	}
	return nil
}

// RenderChunksAndSourceMaps renders non-weak chunks and optional per-chunk
// source maps in parallel across a worker pool while preserving sequential
// dependency-order license deduplication across chunks.
func (c *Compiler) RenderChunksAndSourceMaps(
	chunks []*JSChunk,
	opts RenderChunkOptions,
) (rendered []RenderedChunk, err error) {
	c.InitCompilerOptionsIfTesting()

	activeChunks := filterActiveChunks(chunks)
	if len(activeChunks) == 0 {
		return []RenderedChunk{}, nil
	}

	numWorkers := opts.NumWorkers
	if numWorkers <= 0 {
		if c.options != nil && c.options.GetNumParallelThreads() > 1 {
			numWorkers = c.options.GetNumParallelThreads()
		} else {
			numWorkers = runtime.GOMAXPROCS(0)
		}
	}
	totalInputs := 0
	for _, m := range activeChunks {
		totalInputs += len(c.chunkInputsForOutput(m))
	}
	numWorkers = max(1, min(numWorkers, max(len(activeChunks), totalInputs)))

	var resetLTPerChunk func()
	lt := opts.LicenseTracker
	if lt == nil {
		if opts.PerChunkScriptLicenses {
			scriptLT := NewCompiler_ScriptNodeLicensesOnlyTracker(c)
			lt = scriptLT
			resetLTPerChunk = func() {
				clear(scriptLT.globalSeen)
				resetOrderedSet(&scriptLT.current)
				scriptLT.lastSeenFile = ""
			}
		} else {
			lt = NewCompiler_ChunkGraphAwareLicenseTracker(c)
		}
	}
	chunkTracker, _ := lt.(*Compiler_ChunkGraphAwareLicenseTracker)

	gatherSourceMap := opts.GenerateSourceMaps || (c.options != nil && c.options.ShouldGatherSourceMapInfo())
	if gatherSourceMap {
		c.ensureSourceMapInitialized()
		c.ResetAndInitializeSourceMap()
	}

	chunkFileName := opts.ChunkFileName
	if chunkFileName == nil {
		chunkFileName = func(m *JSChunk) string { return m.GetName() + ".js" }
	}

	hasCompress := opts.CompressSourceMap != nil || opts.CompressSourceMapBytes != nil
	omitString := opts.OmitSourceMapString && hasCompress
	renderOneMap := func(sm *SourceMap, m *JSChunk, mapBuf *bytes.Buffer) (string, []byte, error) {
		if opts.CompressSourceMapBytes != nil && omitString {
			raw := renderChunkSourceMapBytesInto(mapBuf, sm, m, opts.WrapperPrefix, chunkFileName)
			b, cErr := opts.CompressSourceMapBytes(m, raw)
			if mapBuf.Cap() > 16<<20 {
				*mapBuf = bytes.Buffer{}
			}
			return "", b, cErr
		}
		smText := renderChunkSourceMapString(sm, m, opts.WrapperPrefix, chunkFileName)
		var smBytes []byte
		var cErr error
		if opts.CompressSourceMapBytes != nil {
			smBytes, cErr = opts.CompressSourceMapBytes(m, []byte(smText))
		} else if opts.CompressSourceMap != nil {
			smBytes, cErr = opts.CompressSourceMap(m, smText)
		}
		if omitString {
			smText = ""
		}
		return smText, smBytes, cErr
	}

	replay, onlyScript, canParallel := compiler_licenseTrackerReplay(lt)
	rendered = make([]RenderedChunk, len(activeChunks))
	if !canParallel {
		var seqMapBuf bytes.Buffer
		for j, m := range activeChunks {
			inputs := c.chunkInputsForOutput(m)
			for _, in := range inputs {
				if in.GetAstRoot(c).IsNil() {
					return nil, jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + m.GetName())
				}
			}
			if gatherSourceMap {
				c.ResetAndInitializeSourceMap()
			}
			if resetLTPerChunk != nil {
				resetLTPerChunk()
			}
			if chunkTracker != nil {
				chunkTracker.SetCurrentChunkContext(m)
			}
			var code string
			if opts.PostProcessInput != nil {
				code = c.toSourceChunkSequentialWithPostProcess(lt, m, gatherSourceMap, opts.PostProcessInput)
			} else {
				code = c.ToSourceOfLicenseTrackerAndJSChunk(lt, m)
			}
			var smText string
			var smBytes []byte
			if gatherSourceMap {
				smText, smBytes, err = renderOneMap(c.GetSourceMap(), m, &seqMapBuf)
				if err != nil {
					return nil, err
				}
			}
			rendered[j] = RenderedChunk{Chunk: m, Code: code, SourceMap: smText, SourceMapBytes: smBytes}
		}
		return rendered, nil
	}

	var parErrs []error
	if gatherSourceMap && hasCompress {
		parErrs = make([]error, len(activeChunks))
	}
	var workerMapBufs []bytes.Buffer
	if gatherSourceMap && opts.CompressSourceMapBytes != nil && omitString {
		workerMapBufs = make([]bytes.Buffer, max(1, numWorkers))
	}
	var onParChunk func(workerIdx, globalIdx, batchIdx int, m *JSChunk, sm *SourceMap)
	if gatherSourceMap {
		onParChunk = func(w, globalIdx, _ int, m *JSChunk, sm *SourceMap) {
			var mapBuf *bytes.Buffer
			if len(workerMapBufs) > 0 {
				mapBuf = &workerMapBufs[w]
			}
			smText, smBytes, cErr := renderOneMap(sm, m, mapBuf)
			if cErr != nil {
				parErrs[globalIdx] = cErr
				return
			}
			rendered[globalIdx].SourceMap = smText
			rendered[globalIdx].SourceMapBytes = smBytes
		}
	}
	err = c.renderChunksBatched(
		activeChunks,
		numWorkers,
		lt,
		onlyScript,
		replay,
		resetLTPerChunk,
		gatherSourceMap,
		opts.PostProcessInput,
		func(globalIdx, _ int, m *JSChunk, code string) {
			rendered[globalIdx] = RenderedChunk{Chunk: m, Code: code}
		},
		onParChunk,
	)
	if err != nil {
		return nil, err
	}
	for _, pErr := range parErrs {
		if pErr != nil {
			return nil, pErr
		}
	}
	return rendered, nil
}

func (c *Compiler) toSourceChunkSequentialWithPostProcess(
	lt CodePrinter_LicenseTracker,
	chunk *JSChunk,
	gatherSourceMap bool,
	postProcess func(*CompilerInput, string) string,
) string {
	return Compiler_RunInCompilerThread(c, func() string {
		inputs := c.chunkInputsForOutput(chunk)
		if len(inputs) == 0 {
			return ""
		}
		cb := NewCompiler_CodeBuilder()
		scratch, owned := c.acquirePrintScratch()
		if owned {
			defer c.sync.printScratchBusy.Store(false)
		}
		sm := CodePrinter_SourceAndMappings{bytesOnly: true}
		for i, in := range inputs {
			scriptNode := in.GetAstRoot(c)
			c.printRootWithScratch(scriptNode, i == 0, gatherSourceMap, lt, &sm, scratch)
			orig := string(sm.codeBytes)
			if out := postProcess(in, orig); out != orig {
				sm.codeBytes = append(sm.codeBytes[:0], out...)
				sm.hasMetrics = false
			}
			startLine, startCol, emitted := c.appendSourceAndMetricsToCodeBuilderForChunk(cb, lt, i, scriptNode, &sm)
			if emitted && gatherSourceMap {
				addMappingsToSourceMap(c.sourceMap, startLine, startCol, &sm)
			} else {
				resetSourceAndMappings(&sm)
			}
		}
		return cb.String()
	})
}

// OutputChunksParallel renders non-weak chunks in bounded batches across
// options.GetNumParallelThreads() workers, preserving sequential license
// deduplication order across chunks and emitting per-chunk JS and optional V3
// source maps in parallel.
func (c *Compiler) OutputChunksParallel(
	chunks []*JSChunk,
	options *CompilerOptions,
	mlicenseTracker *Compiler_ChunkGraphAwareLicenseTracker,
	hooks ParallelChunkOutputCallbacks,
) {
	activeChunks := filterActiveChunks(chunks)
	if len(activeChunks) == 0 {
		return
	}

	numThreads := max(1, options.GetNumParallelThreads())
	gatherSourceMap := options.ShouldGatherSourceMapInfo()
	if gatherSourceMap {
		c.ResetAndInitializeSourceMap()
	}

	type chunkWriteTask struct {
		filename    string
		mapFilename string
		code        string
		wrapper     string
	}
	writeTasks := make([]chunkWriteTask, min(len(activeChunks), chunkRenderMaxBatchChunks))

	err := c.renderChunksBatched(
		activeChunks,
		numThreads,
		mlicenseTracker,
		false,
		mlicenseTracker.trackLicensesForFileName,
		nil,
		gatherSourceMap,
		nil,
		func(_, j int, m *JSChunk, code string) {
			chunkFilename := hooks.ChunkFileName(m)
			hooks.EnsureDir(chunkFilename)
			var mapFilename string
			if gatherSourceMap {
				mapFilename = hooks.SourceMapPath(m)
			}
			wrapper := strings.ReplaceAll(hooks.WrapperTemplate(m.GetName()), "%basename%", filepath.Base(chunkFilename))
			writeTasks[j] = chunkWriteTask{
				filename:    chunkFilename,
				mapFilename: mapFilename,
				code:        code,
				wrapper:     wrapper,
			}
		},
		func(_, _, j int, _ *JSChunk, sm *SourceMap) {
			t := &writeTasks[j]
			defer func() {
				*t = chunkWriteTask{}
			}()
			if sm == nil {
				hooks.WriteChunkNoMap(t.filename, t.code, t.wrapper)
				return
			}
			hooks.WriteChunkWithMap(t.filename, t.mapFilename, sm, t.code, t.wrapper)
		},
	)
	if err != nil {
		panic(err)
	}
}
