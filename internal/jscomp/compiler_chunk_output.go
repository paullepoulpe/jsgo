package jscomp

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RenderedChunk holds the rendered JavaScript source and optional source map
// for a single JSChunk.
type RenderedChunk struct {
	Chunk     *JSChunk
	Code      string
	SourceMap string
}

// RenderChunkOptions configures optional license tracking, source-map naming,
// wrapper prefixing, and worker parallelism for Compiler.RenderChunksAndSourceMaps.
type RenderChunkOptions struct {
	// LicenseTracker tracks emitted licenses across chunks in dependency order.
	// When nil, NewCompiler_ChunkGraphAwareLicenseTracker(c) is used.
	LicenseTracker CodePrinter_LicenseTracker

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

type chunkPrintedInput struct {
	root      rhino.Node
	sm        CodePrinter_SourceAndMappings
	rec       compiler_ParallelScriptLicenseRecorder
	startLine int
	startCol  int
	emitted   bool
}

type chunkPrinted struct {
	inputs []chunkPrintedInput
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

func (c *Compiler) printChunkInput(
	chunkName string,
	scriptNode rhino.Node,
	firstOutput bool,
	gatherSourceMap bool,
	licenseTracker CodePrinter_LicenseTracker,
	sm *CodePrinter_SourceAndMappings,
	scratch *codePrinterScratch,
) (string, error) {
	if scriptNode.IsNil() {
		return "", jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + chunkName)
	}
	return c.printRootWithScratch(scriptNode, firstOutput, gatherSourceMap, licenseTracker, sm, scratch), nil
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

func populateChunkSourceMap(sm *SourceMap, pc *chunkPrinted) {
	sm.Reset()
	for i := range pc.inputs {
		pi := &pc.inputs[i]
		if pi.emitted {
			sm.SetStartingPosition(pi.startLine, pi.startCol)
			for k := range pi.sm.mappings {
				mp := &pi.sm.mappings[k]
				sm.AddMappingOfNodePositions(mp.node, mp.start, mp.end)
			}
		}
		resetSourceAndMappings(&pi.sm)
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
	sm.AppendTo(&mapBuf, &name)
	return mapBuf.String()
}

func (c *Compiler) renderChunksBatched(
	activeChunks []*JSChunk,
	numWorkers int,
	lt CodePrinter_LicenseTracker,
	onlyScript bool,
	replay func(string),
	gatherSourceMap bool,
	onSeqChunk func(globalIdx, batchIdx int, m *JSChunk, code string),
	onParChunk func(globalIdx, batchIdx int, m *JSChunk, sm *SourceMap),
) error {
	const batchSize = 256
	chunkTracker, _ := lt.(*Compiler_ChunkGraphAwareLicenseTracker)
	batchCap := min(len(activeChunks), batchSize)
	printed := make([]chunkPrinted, batchCap)
	printErrs := make([]error, batchCap)
	scratches := make([]codePrinterScratch, numWorkers)

	var workerSM []*SourceMap
	if gatherSourceMap && onParChunk != nil {
		workerSM = make([]*SourceMap, numWorkers)
		for w := range workerSM {
			workerSM[w] = c.newWorkerSourceMap()
		}
	}

	cb := NewCompiler_CodeBuilder()
	for batchStart := 0; batchStart < len(activeChunks); batchStart += batchSize {
		batchLen := min(batchSize, len(activeChunks)-batchStart)
		clear(printErrs[:batchLen])

		failures := runParallelIndexedTasks(batchLen, numWorkers, func(w, j int) {
			m := activeChunks[batchStart+j]
			inputs := m.GetInputs()
			nInputs := len(inputs)
			pc := &printed[j]
			if cap(pc.inputs) >= nInputs {
				pc.inputs = pc.inputs[:nInputs]
			} else {
				pc.inputs = make([]chunkPrintedInput, nInputs)
			}
			for i, in := range inputs {
				scriptNode := in.GetAstRoot(c)
				pi := &pc.inputs[i]
				pi.root = scriptNode
				pi.rec.Reset(onlyScript)
				pi.sm.mappings = pi.sm.mappings[:0]
				pi.startLine = 0
				pi.startCol = 0
				pi.emitted = false
				var printErr error
				if gatherSourceMap {
					_, printErr = c.printChunkInput(m.GetName(), scriptNode, i == 0, true, &pi.rec, &pi.sm, &scratches[w])
				} else {
					pi.sm.source, printErr = c.printChunkInput(m.GetName(), scriptNode, i == 0, false, &pi.rec, nil, &scratches[w])
				}
				if printErr != nil {
					printErrs[j] = printErr
					return
				}
			}
		})
		for j, f := range failures {
			if printErrs[j] != nil {
				return printErrs[j]
			}
			if f != nil {
				return f
			}
		}

		for j := range batchLen {
			globalIdx := batchStart + j
			m := activeChunks[globalIdx]
			if chunkTracker != nil {
				chunkTracker.SetCurrentChunkContext(m)
			}
			pc := &printed[j]
			var code string
			if len(pc.inputs) > 0 {
				cb.ResetAll()
				for i := range pc.inputs {
					pi := &pc.inputs[i]
					pi.rec.Replay(replay)
					pi.startLine, pi.startCol, pi.emitted = c.appendSourceToCodeBuilderForChunk(
						cb, lt, i, pi.root, pi.sm.source)
					pi.sm.source = ""
					pi.root = rhino.NilNode
				}
				code = cb.String()
			}
			onSeqChunk(globalIdx, j, m, code)
		}

		if onParChunk != nil {
			parFailures := runParallelIndexedTasks(batchLen, numWorkers, func(w, j int) {
				globalIdx := batchStart + j
				m := activeChunks[globalIdx]
				var sm *SourceMap
				if gatherSourceMap {
					sm = workerSM[w]
					if globalIdx == len(activeChunks)-1 && c.sourceMap != nil {
						sm = c.sourceMap
					}
					populateChunkSourceMap(sm, &printed[j])
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
	numWorkers = max(1, min(numWorkers, len(activeChunks)))

	lt := opts.LicenseTracker
	if lt == nil {
		lt = NewCompiler_ChunkGraphAwareLicenseTracker(c)
	}
	chunkTracker, _ := lt.(*Compiler_ChunkGraphAwareLicenseTracker)

	gatherSourceMap := opts.GenerateSourceMaps || (c.options != nil && c.options.ShouldGatherSourceMapInfo())
	if gatherSourceMap {
		c.ensureSourceMapInitialized()
		c.ResetAndIntitializeSourceMap()
	}

	chunkFileName := opts.ChunkFileName
	if chunkFileName == nil {
		chunkFileName = func(m *JSChunk) string { return m.GetName() + ".js" }
	}

	replay, onlyScript, canParallel := compiler_licenseTrackerReplay(lt)
	rendered = make([]RenderedChunk, len(activeChunks))
	if !canParallel {
		for j, m := range activeChunks {
			for _, in := range m.GetInputs() {
				if in.GetAstRoot(c).IsNil() {
					return nil, jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + m.GetName())
				}
			}
			if gatherSourceMap {
				c.ResetAndIntitializeSourceMap()
			}
			if chunkTracker != nil {
				chunkTracker.SetCurrentChunkContext(m)
			}
			code := c.ToSourceOfLicenseTrackerAndJSChunk(lt, m)
			var smText string
			if gatherSourceMap {
				smText = renderChunkSourceMapString(c.GetSourceMap(), m, opts.WrapperPrefix, chunkFileName)
			}
			rendered[j] = RenderedChunk{Chunk: m, Code: code, SourceMap: smText}
		}
		return rendered, nil
	}

	var onParChunk func(globalIdx, batchIdx int, m *JSChunk, sm *SourceMap)
	if gatherSourceMap {
		onParChunk = func(globalIdx, _ int, m *JSChunk, sm *SourceMap) {
			rendered[globalIdx].SourceMap = renderChunkSourceMapString(sm, m, opts.WrapperPrefix, chunkFileName)
		}
	}
	err = c.renderChunksBatched(
		activeChunks,
		numWorkers,
		lt,
		onlyScript,
		replay,
		gatherSourceMap,
		func(globalIdx, _ int, m *JSChunk, code string) {
			rendered[globalIdx] = RenderedChunk{Chunk: m, Code: code}
		},
		onParChunk,
	)
	if err != nil {
		return nil, err
	}
	return rendered, nil
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
		c.ResetAndIntitializeSourceMap()
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
		gatherSourceMap,
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
