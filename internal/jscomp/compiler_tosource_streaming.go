package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// toSourceStreaming is ToSource's multi-script path (Go-only). Workers print
// the scripts in index order while this goroutine merges the finished
// scripts, in order, as soon as they complete: the same calls in the same
// order as the sequential path, so the output is identical. In fast mode
// (resolveOnWorkers) the workers also resolve their script's source-map
// mappings with a per-worker SourceMap; that can report input-source-map
// errors from a worker, so compat keeps resolving in the merge. At most
// `window` scripts are printed but not yet merged, so the printed text and
// mappings of the whole program are never held at once, and the sequential
// merge overlaps with printing instead of following it.
func (c *Compiler) toSourceStreaming(
	cb *Compiler_CodeBuilder, lt *Compiler_SingleBinaryLicenseTracker, scripts []rhino.Node, resolveOnWorkers bool) {
	numThreads := c.options.GetNumParallelThreads()
	resolve := resolveOnWorkers && c.options.ShouldGatherSourceMapInfo() && c.sourceMap != nil
	window := min(len(scripts), 2*numThreads)

	type slot struct {
		sm       CodePrinter_SourceAndMappings
		licenses []string
		lt       *Compiler_SingleBinaryLicenseTracker
		resolved []sourceMap_resolvedMapping
		failure  jutil.Throwable
		done     chan struct{}
	}
	slots := make([]slot, window)
	for i := range slots {
		slots[i].done = make(chan struct{}, 1)
	}
	// One token per slot: a worker takes a token before claiming the next
	// script and the merge returns it once that script is appended. The
	// claimed, unmerged scripts are therefore always a run of at most
	// `window` consecutive indices, so script i owns slots[i%window].
	tokens := make(chan struct{}, window)
	for range window {
		tokens <- struct{}{}
	}
	stop := make(chan struct{})
	var next atomic.Int64
	var wg sync.WaitGroup
	for w := range min(numThreads, len(scripts)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var scratch codePrinterScratch
			var rs *SourceMap
			for {
				select {
				case <-tokens:
				case <-stop:
					return
				}
				idx := int(next.Add(1) - 1)
				if idx >= len(scripts) {
					return
				}
				s := &slots[idx%window]
				s.failure = prebuildRunIndexedTask(func(int, int) {
					if s.lt == nil {
						s.lt = newCompiler_SingleBinaryLicenseTracker(c)
					} else {
						s.lt.current.clear()
						s.lt.lastSeenFile = ""
					}
					c.toSourceAndMappingsIntoWithScratch(scripts[idx], idx == 0, s.lt, &s.sm, &scratch)
					if resolve {
						if rs == nil {
							rs = c.newWorkerSourceMap()
						}
						out := s.resolved[:0]
						if cap(out) < len(s.sm.mappings) {
							out = make([]sourceMap_resolvedMapping, 0, len(s.sm.mappings))
						}
						for k := range s.sm.mappings {
							m := &s.sm.mappings[k]
							if r, ok := rs.resolveMappingOfNodePositions(m.node, m.start, m.end); ok {
								out = append(out, r)
							}
						}
						s.resolved = out
					}
					if s.lt.current.isEmpty() {
						s.licenses = s.licenses[:0]
					} else {
						s.licenses = append(s.licenses[:0], s.lt.current.order...)
					}
				}, w, idx)
				s.done <- struct{}{}
			}
		}()
	}

	var failure jutil.Throwable
	for idx := range scripts {
		s := &slots[idx%window]
		<-s.done
		if s.failure != nil {
			failure = s.failure
			break
		}
		for _, lic := range s.licenses {
			lt.addLicense(lic)
		}
		s.licenses = s.licenses[:0]
		if resolve {
			c.appendSourceAndResolvedMappingsToCodeBuilder(cb, lt, idx, scripts[idx], &s.sm, s.resolved)
			if cap(s.resolved) > 1<<16 {
				s.resolved = nil
			} else {
				s.resolved = s.resolved[:0]
			}
		} else {
			c.appendSourceAndMappingsToCodeBuilder(cb, lt, idx, scripts[idx], &s.sm)
		}
		s.sm.source = ""
		tokens <- struct{}{}
	}
	close(stop)
	wg.Wait()
	if failure != nil {
		panic(failure)
	}
}
