package jscomp

import (
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
		c.sourceMap.SetSourceFileMapping(c.AsSourceFileMapping())
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
	builder := NewCodePrinter_Builder(scriptNode)
	builder.SetCompilerOptions(c.options)
	builder.SetTagAsTypeSummary(c.options.ShouldGenerateTypedExterns())
	builder.SetTagAsStrict(firstOutput && c.options.ShouldEmitUseStrict())
	builder.SetLicenseTracker(licenseTracker)
	return builder.buildSourceAndMappingsInto(sm, gatherSourceMap, scratch), nil
}

// RenderChunksAndSourceMaps renders non-weak chunks and optional per-chunk
// source maps in parallel across a worker pool while preserving sequential
// dependency-order license deduplication across chunks.
func (c *Compiler) RenderChunksAndSourceMaps(
	chunks []*JSChunk,
	opts RenderChunkOptions,
) (rendered []RenderedChunk, err error) {
	c.InitCompilerOptionsIfTesting()

	activeChunks := make([]*JSChunk, 0, len(chunks))
	for _, m := range chunks {
		if m != nil && m.GetName() != JSChunk_WEAK_CHUNK_NAME {
			activeChunks = append(activeChunks, m)
		}
	}
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
				sm := c.GetSourceMap()
				if opts.WrapperPrefix != nil {
					if prefix := opts.WrapperPrefix(m); prefix != "" {
						sm.SetWrapperPrefix(prefix)
					}
				}
				name := chunkFileName(m)
				var mapBuf strings.Builder
				sm.AppendTo(&mapBuf, &name)
				smText = mapBuf.String()
			}
			rendered[j] = RenderedChunk{Chunk: m, Code: code, SourceMap: smText}
		}
		return rendered, nil
	}

	printed := make([]chunkPrinted, len(activeChunks))
	printErrs := make([]error, len(activeChunks))
	scratches := make([]codePrinterScratch, numWorkers)
	failures := parallelPrebuildIndexed(len(activeChunks), numWorkers, func(w, j int) {
		m := activeChunks[j]
		inputs := m.GetInputs()
		if len(inputs) == 0 {
			return
		}
		pc := &printed[j]
		pc.inputs = make([]chunkPrintedInput, len(inputs))
		for i, in := range inputs {
			scriptNode := in.GetAstRoot(c)
			pi := &pc.inputs[i]
			pi.root = scriptNode
			pi.rec.Reset(onlyScript)
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
			return nil, printErrs[j]
		}
		if f != nil {
			return nil, f
		}
	}

	cb := NewCompiler_CodeBuilder()
	for j, m := range activeChunks {
		if chunkTracker != nil {
			chunkTracker.SetCurrentChunkContext(m)
		}
		pc := &printed[j]
		var code string
		if len(pc.inputs) > 0 {
			cb.ResetAll()
			for i := range pc.inputs {
				pi := &pc.inputs[i]
				for _, file := range pi.rec.seenFiles {
					replay(file)
				}
				pi.startLine, pi.startCol, pi.emitted = c.appendSourceToCodeBuilderForChunk(
					cb, lt, i, pi.root, pi.sm.source)
				pi.sm.source = ""
				pi.root = rhino.NilNode
			}
			code = cb.String()
		}
		rendered[j] = RenderedChunk{Chunk: m, Code: code}
	}

	if gatherSourceMap {
		workerSM := make([]*SourceMap, numWorkers)
		for w := 0; w < numWorkers; w++ {
			workerSM[w] = c.newWorkerSourceMap()
		}
		mapFailures := parallelPrebuildIndexed(len(activeChunks), numWorkers, func(w, j int) {
			m := activeChunks[j]
			sm := workerSM[w]
			if j == len(activeChunks)-1 && c.sourceMap != nil {
				sm = c.sourceMap
			}
			sm.Reset()
			pc := &printed[j]
			for i := range pc.inputs {
				pi := &pc.inputs[i]
				if pi.emitted {
					sm.SetStartingPosition(pi.startLine, pi.startCol)
					for k := range pi.sm.mappings {
						mp := &pi.sm.mappings[k]
						sm.AddMappingOfNodePositions(mp.node, mp.start, mp.end)
					}
				}
				pi.sm.mappings = nil
			}
			if opts.WrapperPrefix != nil {
				if prefix := opts.WrapperPrefix(m); prefix != "" {
					sm.SetWrapperPrefix(prefix)
				}
			}
			name := chunkFileName(m)
			var mapBuf strings.Builder
			sm.AppendTo(&mapBuf, &name)
			rendered[j].SourceMap = mapBuf.String()
		})
		for _, f := range mapFailures {
			if f != nil {
				return nil, f
			}
		}
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
	activeChunks := make([]*JSChunk, 0, len(chunks))
	for _, m := range chunks {
		if m.GetName() != JSChunk_WEAK_CHUNK_NAME {
			activeChunks = append(activeChunks, m)
		}
	}
	if len(activeChunks) == 0 {
		return
	}

	numThreads := max(1, options.GetNumParallelThreads())
	gatherSourceMap := options.ShouldGatherSourceMapInfo()
	var workerSM []*SourceMap
	if gatherSourceMap {
		c.ResetAndIntitializeSourceMap()
		workerSM = make([]*SourceMap, numThreads)
		for w := 0; w < numThreads; w++ {
			workerSM[w] = c.newWorkerSourceMap()
		}
	}

	const batchSize = 256
	type chunkWriteTask struct {
		filename    string
		mapFilename string
		code        string
		wrapper     string
	}

	printed := make([]chunkPrinted, min(len(activeChunks), batchSize))
	writeTasks := make([]chunkWriteTask, min(len(activeChunks), batchSize))
	scratches := make([]codePrinterScratch, numThreads)

	for batchStart := 0; batchStart < len(activeChunks); batchStart += batchSize {
		batchEnd := min(batchStart+batchSize, len(activeChunks))
		batchLen := batchEnd - batchStart

		failures := parallelPrebuildIndexed(batchLen, numThreads, func(w, j int) {
			m := activeChunks[batchStart+j]
			inputs := m.GetInputs()
			nInputs := len(inputs)
			pc := &printed[j]
			if cap(pc.inputs) >= nInputs {
				pc.inputs = pc.inputs[:nInputs]
			} else {
				pc.inputs = make([]chunkPrintedInput, nInputs)
			}
			for i := 0; i < nInputs; i++ {
				scriptNode := inputs[i].GetAstRoot(c)
				pi := &pc.inputs[i]
				pi.root = scriptNode
				pi.rec.Reset(false)
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
					panic(printErr)
				}
			}
		})
		for _, f := range failures {
			if f != nil {
				panic(f)
			}
		}

		cb := NewCompiler_CodeBuilder()
		for j := 0; j < batchLen; j++ {
			m := activeChunks[batchStart+j]
			chunkFilename := hooks.ChunkFileName(m)
			hooks.EnsureDir(chunkFilename)
			var mapFilename string
			if gatherSourceMap {
				mapFilename = hooks.SourceMapPath(m)
			}
			mlicenseTracker.SetCurrentChunkContext(m)
			pc := &printed[j]
			var code string
			if len(pc.inputs) > 0 {
				cb.ResetAll()
				for i := range pc.inputs {
					pi := &pc.inputs[i]
					for _, file := range pi.rec.seenFiles {
						mlicenseTracker.trackLicensesForFileName(file)
					}
					pi.startLine, pi.startCol, pi.emitted = c.appendSourceToCodeBuilderForChunk(
						cb, mlicenseTracker, i, pi.root, pi.sm.source)
					pi.sm.source = ""
					pi.root = rhino.NilNode
				}
				code = cb.String()
			}
			baseName := jutil.NewFile(chunkFilename).GetName()
			wrapper := strings.ReplaceAll(hooks.WrapperTemplate(m.GetName()), "%basename%", baseName)
			writeTasks[j] = chunkWriteTask{
				filename:    chunkFilename,
				mapFilename: mapFilename,
				code:        code,
				wrapper:     wrapper,
			}
		}

		wFailures := parallelPrebuildIndexed(batchLen, numThreads, func(w, j int) {
			t := &writeTasks[j]
			defer func() {
				*t = chunkWriteTask{}
			}()
			if !gatherSourceMap {
				hooks.WriteChunkNoMap(t.filename, t.code, t.wrapper)
				return
			}
			globalIdx := batchStart + j
			sm := workerSM[w]
			if globalIdx == len(activeChunks)-1 && c.sourceMap != nil {
				sm = c.sourceMap
			}
			sm.Reset()
			pc := &printed[j]
			for i := range pc.inputs {
				pi := &pc.inputs[i]
				if pi.emitted {
					sm.SetStartingPosition(pi.startLine, pi.startCol)
					for k := range pi.sm.mappings {
						mp := &pi.sm.mappings[k]
						sm.AddMappingOfNodePositions(mp.node, mp.start, mp.end)
					}
				}
				if cap(pi.sm.mappings) > 16384 {
					pi.sm.mappings = nil
				} else {
					pi.sm.mappings = pi.sm.mappings[:0]
				}
			}
			hooks.WriteChunkWithMap(t.filename, t.mapFilename, sm, t.code, t.wrapper)
		})
		for _, f := range wFailures {
			if f != nil {
				panic(f)
			}
		}
	}
}
