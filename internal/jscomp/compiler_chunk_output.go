package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ParallelChunkOutputCallbacks supplies the CLI-specific path expansion,
// wrapper lookup, and stream writing hooks to Compiler.OutputChunksParallel.
type ParallelChunkOutputCallbacks struct {
	ChunkFileName     func(m *JSChunk) string
	SourceMapPath     func(m *JSChunk) *string
	WrapperTemplate   func(chunkName string) string
	EnsureDir         func(filename *string)
	WriteChunkWithMap func(filename string, mapFilename string, sm *SourceMap, code string, wrapper string)
	WriteChunkNoMap   func(filename string, code string, wrapper string)
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

	numThreads := options.GetNumParallelThreads()
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
	type printedInput struct {
		root      rhino.Node
		sm        CodePrinter_SourceAndMappings
		rec       compiler_ParallelScriptLicenseRecorder
		startLine int
		startCol  int
		emitted   bool
	}
	type printedChunk struct {
		inputs []printedInput
		single [2]printedInput
	}
	type chunkWriteTask struct {
		filename    string
		mapFilename string
		code        string
		wrapper     string
	}

	printed := make([]printedChunk, min(len(activeChunks), batchSize))
	writeBufA := make([]chunkWriteTask, min(len(activeChunks), batchSize))
	var writeBufB []chunkWriteTask
	if !gatherSourceMap {
		writeBufB = make([]chunkWriteTask, min(len(activeChunks), batchSize))
	}

	var prevWriteDone chan jutil.Throwable
	waitPrevWrite := func() {
		if prevWriteDone != nil {
			if err := <-prevWriteDone; err != nil {
				panic(err)
			}
			prevWriteDone = nil
		}
	}
	defer waitPrevWrite()

	useBufB := false
	scratches := make([]codePrinterScratch, numThreads)
	for batchStart := 0; batchStart < len(activeChunks); batchStart += batchSize {
		batchEnd := min(batchStart+batchSize, len(activeChunks))
		batchLen := batchEnd - batchStart

		failures := parallelPrebuildIndexed(batchLen, numThreads, func(w, j int) {
			m := activeChunks[batchStart+j]
			inputs := m.GetInputs()
			nInputs := len(inputs)
			pc := &printed[j]
			if nInputs == 0 {
				pc.inputs = nil
				return
			}
			if nInputs <= len(pc.single) {
				pc.inputs = pc.single[:nInputs]
			} else if cap(pc.inputs) >= nInputs {
				pc.inputs = pc.inputs[:nInputs]
			} else {
				pc.inputs = make([]printedInput, nInputs)
			}
			for i := 0; i < nInputs; i++ {
				scriptNode := inputs[i].GetAstRoot(c)
				if scriptNode.IsNil() {
					panic(jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + m.GetName()))
				}
				pi := &pc.inputs[i]
				pi.root = scriptNode
				pi.rec.Reset(false)
				pi.sm.mappings = pi.sm.mappings[:0]
				pi.startLine = 0
				pi.startCol = 0
				pi.emitted = false
				c.toSourceAndMappingsIntoWithScratch(scriptNode, i == 0, &pi.rec, &pi.sm, &scratches[w])
			}
		})
		for _, f := range failures {
			if f != nil {
				waitPrevWrite()
				panic(f)
			}
		}

		waitPrevWrite()
		writeTasks := writeBufA[:batchLen]
		if !gatherSourceMap {
			if useBufB {
				writeTasks = writeBufB[:batchLen]
			}
			useBufB = !useBufB
		}

		cb := NewCompiler_CodeBuilder()
		for j := 0; j < batchLen; j++ {
			m := activeChunks[batchStart+j]
			chunkFilename := hooks.ChunkFileName(m)
			hooks.EnsureDir(&chunkFilename)
			var mapFilename string
			if gatherSourceMap {
				if mp := hooks.SourceMapPath(m); mp != nil {
					mapFilename = *mp
				}
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
					if gatherSourceMap {
						pi.startLine, pi.startCol, pi.emitted = c.appendSourceToCodeBuilderForChunk(
							cb, mlicenseTracker, i, pi.root, pi.sm.source)
						pi.sm.source = ""
					} else {
						c.appendSourceAndMappingsToCodeBuilder(cb, mlicenseTracker, i, pi.root, &pi.sm)
					}
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

		if gatherSourceMap {
			wFailures := parallelPrebuildIndexed(batchLen, numThreads, func(w, j int) {
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
				t := &writeTasks[j]
				hooks.WriteChunkWithMap(t.filename, t.mapFilename, sm, t.code, t.wrapper)
				t.code = ""
				t.wrapper = ""
				t.filename = ""
				t.mapFilename = ""
			})
			for _, f := range wFailures {
				if f != nil {
					panic(f)
				}
			}
			continue
		}

		done := make(chan jutil.Throwable, 1)
		prevWriteDone = done
		go func(tasks []chunkWriteTask) {
			wFailures := parallelPrebuild(len(tasks), numThreads, func(j int) {
				t := &tasks[j]
				hooks.WriteChunkNoMap(t.filename, t.code, t.wrapper)
				t.code = ""
				t.wrapper = ""
				t.filename = ""
			})
			for _, f := range wFailures {
				if f != nil {
					done <- f
					return
				}
			}
			done <- nil
		}(writeTasks)
	}
}
