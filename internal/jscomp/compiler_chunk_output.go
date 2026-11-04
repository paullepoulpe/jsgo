package jscomp

import (
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

func renderChunkSourceMapString(
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
	name := chunkFileName(m)
	var mapBuf strings.Builder
	sm.AppendTo(&mapBuf, name)
	return mapBuf.String()
}

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
	onParChunk func(globalIdx, batchIdx int, m *JSChunk, sm *SourceMap),
) error {
	const batchSize = 256
	chunkTracker, _ := lt.(*Compiler_ChunkGraphAwareLicenseTracker)
	batchCap := min(len(activeChunks), batchSize)
	scriptOffsets := make([]int, batchCap+1)
	printer := newParallelScriptPrinter(c, numWorkers, onlyScript, gatherSourceMap, batchCap)
	printer.postProcessInput = postProcessInput

	var workerSM []*SourceMap
	if gatherSourceMap && onParChunk != nil {
		workerSM = make([]*SourceMap, min(numWorkers, batchCap))
		for w := range workerSM {
			workerSM[w] = c.newWorkerSourceMap()
		}
	}

	cb := NewCompiler_CodeBuilder()
	for batchStart := 0; batchStart < len(activeChunks); batchStart += batchSize {
		batchLen := min(batchSize, len(activeChunks)-batchStart)
		totalScripts := 0
		for j := range batchLen {
			scriptOffsets[j] = totalScripts
			totalScripts += len(activeChunks[batchStart+j].GetInputs())
		}
		scriptOffsets[batchLen] = totalScripts

		slots := printer.ensureSlots(totalScripts)
		for j := range batchLen {
			m := activeChunks[batchStart+j]
			base := scriptOffsets[j]
			for i, in := range m.GetInputs() {
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
				onParChunk(globalIdx, j, m, sm)
			})
			for _, f := range parFailures {
				if f != nil {
					return f
				}
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
		totalInputs += len(m.GetInputs())
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

	replay, onlyScript, canParallel := compiler_licenseTrackerReplay(lt)
	rendered = make([]RenderedChunk, len(activeChunks))
	if !canParallel {
		for j, m := range activeChunks {
			inputs := m.GetInputs()
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
				smText = renderChunkSourceMapString(c.GetSourceMap(), m, opts.WrapperPrefix, chunkFileName)
				if opts.CompressSourceMap != nil {
					smBytes, err = opts.CompressSourceMap(m, smText)
					if err != nil {
						return nil, err
					}
				}
			}
			rendered[j] = RenderedChunk{Chunk: m, Code: code, SourceMap: smText, SourceMapBytes: smBytes}
		}
		return rendered, nil
	}

	var parErrs []error
	if gatherSourceMap && opts.CompressSourceMap != nil {
		parErrs = make([]error, len(activeChunks))
	}
	var onParChunk func(globalIdx, batchIdx int, m *JSChunk, sm *SourceMap)
	if gatherSourceMap {
		onParChunk = func(globalIdx, _ int, m *JSChunk, sm *SourceMap) {
			smText := renderChunkSourceMapString(sm, m, opts.WrapperPrefix, chunkFileName)
			rendered[globalIdx].SourceMap = smText
			if opts.CompressSourceMap != nil {
				b, cErr := opts.CompressSourceMap(m, smText)
				if cErr != nil {
					parErrs[globalIdx] = cErr
					return
				}
				rendered[globalIdx].SourceMapBytes = b
			}
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
		inputs := chunk.GetInputs()
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
	writeTasks := make([]chunkWriteTask, min(len(activeChunks), 256))

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
		func(_, j int, _ *JSChunk, sm *SourceMap) {
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
