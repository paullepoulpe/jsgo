package sourcemap

import (
	"io"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: a parallel encoder for the "mappings" field, used when
// SetParallelWorkers(n > 1) was called (jsgo's fast mode). The output is
// byte-identical to the sequential LineMapper's.
//
// Phase 1 runs the mapping traversal once, sequentially, as PrepMappings
// does (computing lineCount), and also:
//   - assigns the source and name ids in first-use order, as WriteEntry would;
//   - simulates the LineMapper's relative-value state;
//   - every parallelChunkMappings mappings, records a checkpoint: the reader
//     position, the traversal's stack and position, and that LineMapper state.
//
// Phase 2 re-runs the traversal from each checkpoint to the next on worker
// goroutines, each with a LineMapper started from the checkpoint's state and
// using the id tables of phase 1 (read-only), writing to a private buffer;
// the buffers are concatenated in order.
//
// Generated lines only grow during the traversal and entries are written
// only for lines below lineCount, so the checkpoints taken below lineCount
// carry exactly the state the sequential LineMapper has there; the chunk that
// reaches lineCount runs to the end.

var parallelChunkMappings = 1 << 15

// SetParallelWorkers sets the number of goroutines AppendTo may use to
// encode the mappings (1, the default, encodes sequentially).
func (g *SourceMapGeneratorV3) SetParallelWorkers(n int) {
	g.parallelWorkers = n
}

type lineMapperState struct {
	previousLine, previousColumn             int
	previousSourceFileId, previousSourceLine int
	previousSourceColumn, previousNameId     int
}

type mappingCheckpoint struct {
	reader     sourceMapGeneratorV3_mappingReader
	stack      []SourceMapGeneratorV3_Mapping
	line, col  int
	state      lineMapperState
	startIndex int // index of the next mapping to read
}

// phase1Visitor is the UsedMappingCheck plus the id assignment and the
// simulated LineMapper state.
type phase1Visitor struct {
	g       *SourceMapGeneratorV3
	maxLine int
	st      lineMapperState
}

func (p *phase1Visitor) Visit(m *SourceMapGeneratorV3_Mapping, line, col, nextLine, nextCol int) {
	if m != nil {
		if endLine := int(m.endLine); endLine > p.maxLine {
			p.maxLine = endLine
		}
	}
	// LineMapper.Visit, assuming line < lineCount (see the file comment).
	if p.st.previousLine != line {
		p.st.previousColumn = 0
	}
	if line != nextLine || col != nextCol {
		p.st.previousColumn = col
		if m != nil {
			g := p.g
			sourceId := g.sourceIdxToFileId[m.sourceFileIdx]
			if sourceId < 0 {
				sourceId = int32(len(g.orderedSourceFiles))
				g.orderedSourceFiles = append(g.orderedSourceFiles, g.addedSourceFiles[m.sourceFileIdx])
				g.sourceIdxToFileId[m.sourceFileIdx] = sourceId
			}
			p.st.previousSourceFileId = int(sourceId)
			p.st.previousSourceLine = int(m.originalLine)
			p.st.previousSourceColumn = int(m.originalCol)
			if m.originalNameIdx >= 0 {
				nameId := g.symbolIdxToNameId[m.originalNameIdx]
				if nameId < 0 {
					nameId = int32(len(g.orderedSymbolNames))
					g.orderedSymbolNames = append(g.orderedSymbolNames, g.addedSymbolNames[m.originalNameIdx])
					g.symbolIdxToNameId[m.originalNameIdx] = nameId
				}
				p.st.previousNameId = int(nameId)
			}
		}
		p.st.previousLine = line
	}
}

// resetIdTables prepares the id tables as AppendLineMappings does.
func (g *SourceMapGeneratorV3) resetIdTables() {
	nSources := len(g.addedSourceFiles)
	g.sourceIdxToFileId = make([]int32, nSources)
	for i := range g.sourceIdxToFileId {
		g.sourceIdxToFileId[i] = -1
	}
	clear(g.orderedSourceFiles)
	g.orderedSourceFiles = g.orderedSourceFiles[:0]
	nSymbols := len(g.addedSymbolNames)
	g.symbolIdxToNameId = make([]int32, nSymbols)
	for i := range g.symbolIdxToNameId {
		g.symbolIdxToNameId[i] = -1
	}
	clear(g.orderedSymbolNames)
	g.orderedSymbolNames = g.orderedSymbolNames[:0]
}

// prepMappingsParallel is PrepMappings plus phase 1; it returns the raw max
// line (before the prefix adjustment) and the checkpoints.
func (g *SourceMapGeneratorV3) prepMappingsParallel() (int, []mappingCheckpoint) {
	g.resetIdTables()
	v := &phase1Visitor{g: g, st: lineMapperState{previousLine: -1}}
	t := SourceMapGeneratorV3_MappingTraversal{outer: g}
	var cps []mappingCheckpoint
	r := sourceMapGeneratorV3_mappingReader{store: &g.mappings}
	idx := 0
	for {
		if idx%parallelChunkMappings == 0 {
			cps = append(cps, mappingCheckpoint{
				reader:     r,
				stack:      append([]SourceMapGeneratorV3_Mapping(nil), t.stack...),
				line:       t.line,
				col:        t.col,
				state:      v.st,
				startIndex: idx,
			})
		}
		if !r.next(&t.cur) {
			break
		}
		idx++
		t.step(v)
	}
	t.finish(v)
	return v.maxLine, cps
}

// step is the body of Traverse's loop for the mapping in t.cur.
func (t *SourceMapGeneratorV3_MappingTraversal) step(v SourceMapGeneratorV3_MappingVisitor) {
	m := &t.cur
	for len(t.stack) > 0 && !t.IsOverlapped(&t.stack[len(t.stack)-1], m) {
		t.closing = t.stack[len(t.stack)-1]
		t.stack = t.stack[:len(t.stack)-1]
		t.MaybeVisit(v, &t.closing)
	}
	var parent *SourceMapGeneratorV3_Mapping
	if len(t.stack) > 0 {
		parent = &t.stack[len(t.stack)-1]
	}
	t.MaybeVisitParent(v, parent, m)
	t.stack = append(t.stack, *m)
}

// finish closes the mappings left on the stack, as the end of Traverse does.
func (t *SourceMapGeneratorV3_MappingTraversal) finish(v SourceMapGeneratorV3_MappingVisitor) {
	for len(t.stack) > 0 {
		t.closing = t.stack[len(t.stack)-1]
		t.stack = t.stack[:len(t.stack)-1]
		t.MaybeVisit(v, &t.closing)
	}
}

// chunkWriter collects a chunk's bytes.
type chunkWriter struct{ b []byte }

func (w *chunkWriter) WriteString(s string) (int, error) {
	w.b = append(w.b, s...)
	return len(s), nil
}

func (w *chunkWriter) Write(p []byte) (int, error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

// appendLineMappingsParallel writes the "mappings" value (quotes included),
// byte-identical to LineMapper.AppendLineMappings.
func (g *SourceMapGeneratorV3) appendLineMappingsParallel(out io.StringWriter, maxLine int, cps []mappingCheckpoint) {
	// Usable chunk starts: checkpoints below lineCount (the first one always).
	n := 1
	for n < len(cps) && cps[n].line < maxLine {
		n++
	}
	ends := make([]int, n) // number of mappings each chunk reads; -1 = to the end
	for i := 0; i < n; i++ {
		if i+1 < n {
			ends[i] = cps[i+1].startIndex - cps[i].startIndex
		} else {
			ends[i] = -1
		}
	}
	bufs := make([][]byte, n)
	workers := min(g.parallelWorkers, n)
	var next atomic.Int64
	var wg sync.WaitGroup
	panics := make([]any, n)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := int(next.Add(1) - 1)
				if i >= n {
					return
				}
				func() {
					defer func() {
						if r := recover(); r != nil {
							panics[i] = r
						}
					}()
					bufs[i] = g.encodeChunk(&cps[i], ends[i], maxLine)
				}()
			}
		}()
	}
	wg.Wait()
	for _, p := range panics {
		if p != nil {
			panic(p)
		}
	}
	appendString(out, "\"")
	for _, b := range bufs {
		appendBytes(out, b)
	}
	appendString(out, ";\"")
}

func (g *SourceMapGeneratorV3) encodeChunk(cp *mappingCheckpoint, count int, maxLine int) []byte {
	w := &chunkWriter{b: make([]byte, 0, 64<<10)}
	l := &SourceMapGeneratorV3_LineMapper{outer: g, out: w, maxLine: maxLine, frozenIds: true}
	l.previousLine = cp.state.previousLine
	l.previousColumn = cp.state.previousColumn
	l.previousSourceFileId = cp.state.previousSourceFileId
	l.previousSourceLine = cp.state.previousSourceLine
	l.previousSourceColumn = cp.state.previousSourceColumn
	l.previousNameId = cp.state.previousNameId
	l.buffering = true
	t := SourceMapGeneratorV3_MappingTraversal{outer: g, line: cp.line, col: cp.col}
	t.stack = append(t.stack, cp.stack...)
	r := cp.reader
	for i := 0; count < 0 || i < count; i++ {
		if !r.next(&t.cur) {
			break
		}
		t.step(l)
	}
	if count < 0 {
		t.finish(l)
	}
	l.flush()
	return w.b
}

// appendToParallel is AppendTo with the parallel mappings encoder.
func (g *SourceMapGeneratorV3) appendToParallel(out io.StringWriter, name *string) {
	rawMax, cps := g.prepMappingsParallel()
	maxLine := jutil.IntAdd(jutil.IntAdd(rawMax, g.prefixPosition.GetLine()), 1)

	appendString(out, "{\n")
	SourceMapGeneratorV3_AppendFirstField(out, "version", "3")
	if name != nil {
		SourceMapGeneratorV3_AppendField(out, "file", SourceMapGeneratorV3_EscapeString(*name))
	}
	SourceMapGeneratorV3_AppendField(out, "lineCount", strconv.Itoa(maxLine))
	if g.sourceRootPath != "" {
		SourceMapGeneratorV3_AppendField(out, "sourceRoot", SourceMapGeneratorV3_EscapeString(g.sourceRootPath))
	}
	SourceMapGeneratorV3_AppendFieldStart(out, "mappings")
	g.appendLineMappingsParallel(out, maxLine, cps)
	SourceMapGeneratorV3_AppendFieldEnd(out)
	g.appendToTail(out)
}
