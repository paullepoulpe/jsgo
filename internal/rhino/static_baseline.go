package rhino

import (
	"math"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/colors"
)

type staticBaselineSnapshot struct {
	captured bool

	next             uint32
	blocksLen        int
	tcsLen           int
	tcNonNil         []bool
	origNamesLen     int
	origNonNil       []bool
	auxNonNil        [auxShards]bool
	typesCount       uint32
	colorsCount      uint32
	sourceFilesCount uint32
	bigIntsCount     uint32

	numberNext uint32
}

var (
	staticBaselineMu sync.Mutex
	staticBaseline   staticBaselineSnapshot
)

// HasStaticBaseline reports whether CaptureStaticBaseline has recorded a
// post-init baseline snapshot.
func HasStaticBaseline() bool {
	staticBaselineMu.Lock()
	ok := staticBaseline.captured
	staticBaselineMu.Unlock()
	return ok
}

// CaptureStaticBaseline snapshots the current post-init state of package-level
// rhino tables (staticStore and numberTable) so subsequent calls to
// RewindStaticBaseline can reclaim per-request template nodes, side-table
// entries, and interned float64 literals in long-lived worker processes.
func CaptureStaticBaseline() {
	staticBaselineMu.Lock()
	defer staticBaselineMu.Unlock()
	captureStaticBaselineLocked()
}

func captureStaticBaselineLocked() {
	s := staticStore
	s.mu.Lock()
	snap := staticBaselineSnapshot{
		captured:         true,
		next:             s.next.Load(),
		blocksLen:        len(s.blocks),
		tcsLen:           len(s.tcs),
		tcNonNil:         make([]bool, len(s.tcs)),
		origNamesLen:     len(s.origNames),
		origNonNil:       make([]bool, len(s.origNames)),
		typesCount:       s.types.count.Load(),
		colorsCount:      s.colors.count.Load(),
		sourceFilesCount: s.sourceFiles.count.Load(),
		bigIntsCount:     s.bigInts.count.Load(),
	}
	for i, c := range s.tcs {
		snap.tcNonNil[i] = c != nil
	}
	for i, c := range s.origNames {
		snap.origNonNil[i] = c != nil
	}
	for i := range s.aux {
		sh := &s.aux[i]
		sh.mu.RLock()
		snap.auxNonNil[i] = len(sh.m) > 0
		sh.mu.RUnlock()
	}
	s.mu.Unlock()

	numberTable.mu.Lock()
	snap.numberNext = numberTable.next
	numberTable.mu.Unlock()

	staticBaseline = snap
}

// RewindStaticBaseline resets the current per-compilation NodeStore and rewinds
// package-level rhino tables (staticStore and numberTable) back to the snapshot
// recorded by CaptureStaticBaseline. If CaptureStaticBaseline has not yet been
// called, RewindStaticBaseline resets the current NodeStore, captures the
// current static state as the baseline, and returns.
//
// Callers must invoke RewindStaticBaseline only between compilations when no
// compilation is actively mutating or reading per-request staticStore nodes.
func RewindStaticBaseline() {
	ResetNodeStore()
	staticBaselineMu.Lock()
	defer staticBaselineMu.Unlock()
	if !staticBaseline.captured {
		captureStaticBaselineLocked()
		return
	}
	snap := staticBaseline

	rewindStaticStore(&snap)
	numberTable.rewind(snap.numberNext)
}

func rewindStaticStore(snap *staticBaselineSnapshot) {
	s := staticStore
	s.mu.Lock()
	defer s.mu.Unlock()

	first := s.base >> blockShift
	curNext := s.next.Load()
	if curNext > snap.next {
		// Unregister and release blocks allocated after baseline.
		for k := snap.blocksLen; k < len(s.blocks); k++ {
			if b := s.blocks[k]; b != nil {
				nodeBlocks[first+uint32(k)].CompareAndSwap(b, nil)
				s.blocks[k] = nil
			}
		}
		s.blocks = s.blocks[:snap.blocksLen]

		// Zero newly populated nodeData slots in the partial tail baseline block
		// so future allocations receive zeroed records.
		if snap.blocksLen > 0 {
			lastK := snap.blocksLen - 1
			blkBase := s.base + uint32(lastK)<<blockShift
			if snap.next < blkBase+blockLen && s.blocks[lastK] != nil {
				startOff := snap.next - blkBase
				endOff := min(uint32(blockLen), curNext-blkBase)
				clear(s.blocks[lastK][startOff:endOff])
			}
		}

		// Rewind tc columns.
		if len(s.tcs) > snap.tcsLen {
			tcField.unregister(first+uint32(snap.tcsLen), s.tcs[snap.tcsLen:])
			clear(s.tcs[snap.tcsLen:])
			s.tcs = s.tcs[:snap.tcsLen]
		}
		for k := range snap.tcsLen {
			if !snap.tcNonNil[k] && s.tcs[k] != nil {
				tcField.dir[first+uint32(k)].CompareAndSwap(s.tcs[k], nil)
				s.tcs[k] = nil
			}
		}
		if snap.blocksLen > 0 {
			lastK := snap.blocksLen - 1
			blkBase := s.base + uint32(lastK)<<blockShift
			if lastK < len(s.tcs) && s.tcs[lastK] != nil && snap.next < blkBase+blockLen {
				startOff := snap.next - blkBase
				endOff := min(uint32(blockLen), curNext-blkBase)
				clear(s.tcs[lastK][startOff:endOff])
			}
		}

		// Rewind origName paged columns.
		if len(s.origNames) > snap.origNamesLen {
			origNameField.unregister(first+uint32(snap.origNamesLen), s.origNames[snap.origNamesLen:])
			clear(s.origNames[snap.origNamesLen:])
			s.origNames = s.origNames[:snap.origNamesLen]
		}
		for k := range snap.origNamesLen {
			if !snap.origNonNil[k] && s.origNames[k] != nil {
				origNameField.dir[first+uint32(k)].CompareAndSwap(s.origNames[k], nil)
				s.origNames[k] = nil
			}
		}
		if snap.blocksLen > 0 {
			lastK := snap.blocksLen - 1
			blkBase := s.base + uint32(lastK)<<blockShift
			if lastK < len(s.origNames) && s.origNames[lastK] != nil && snap.next < blkBase+blockLen {
				col := s.origNames[lastK]
				startOff := snap.next - blkBase
				endOff := min(uint32(blockLen), curNext-blkBase)
				firstDropPage := (startOff + columnPageMask) >> columnPageShift
				lastUsedPage := (endOff - 1) >> columnPageShift
				for p := firstDropPage; p <= lastUsedPage && p < pagesPerBlock; p++ {
					col[p].Store(nil)
				}
				if startOff&columnPageMask != 0 {
					bp := startOff >> columnPageShift
					if pg := col[bp].Load(); pg != nil {
						pageEnd := min(uint32(columnPageLen), endOff-(bp<<columnPageShift))
						clear(pg[startOff&columnPageMask : pageEnd])
					}
				}
			}
		}

		// Rewind cold aux shards; reallocate map if it grew large so bucket memory is freed.
		for i := range s.aux {
			sh := &s.aux[i]
			sh.mu.Lock()
			if len(sh.m) > 0 {
				if !snap.auxNonNil[i] {
					sh.m = nil
				} else {
					prevLen := len(sh.m)
					for id := range sh.m { // order-independent: deleting keys >= snap.next from aux shard
						if id >= snap.next {
							delete(sh.m, id)
						}
					}
					if prevLen > len(sh.m)*2+32 {
						compact := make(map[uint32]nodeAux, len(sh.m))
						for id, v := range sh.m { // order-independent: copying surviving baseline aux entries
							compact[id] = v
						}
						sh.m = compact
					}
				}
			}
			sh.mu.Unlock()
		}

		s.next.Store(snap.next)
	}

	s.types.rewind(snap.typesCount, func(h uint32, v *JSTypeBase) {
		v.nodeHandle.CompareAndSwap(h, 0)
	})
	s.colors.rewind(snap.colorsCount, func(h uint32, v *colors.Color) {
		if storeID, ch := v.LoadNodeHandle(); storeID == s.id && ch == h {
			v.StoreNodeHandle(0, 0)
		}
	})
	s.sourceFiles.rewind(snap.sourceFilesCount)
	s.bigInts.rewind(snap.bigIntsCount)
}

func (t *handleTable[P]) rewind(baseCount uint32, clearOwner func(uint32, *P)) {
	curCount := t.count.Load()
	if curCount <= baseCount {
		return
	}
	lastCurBlk := min(curCount>>tableBlockShift, uint32(tableMaxBlocks-1))
	firstBlk := (baseCount + 1) >> tableBlockShift
	for k := firstBlk; k <= lastCurBlk; k++ {
		var blk *handleBlock[P]
		if k == 0 {
			blk = t.blk0.Load()
		} else if sub := t.dir[k>>tableSubShift].Load(); sub != nil {
			blk = sub[k&tableSubMask].Load()
		}
		if blk == nil {
			continue
		}
		start := uint32(0)
		if k == 0 {
			start = 1
		}
		if k == (baseCount >> tableBlockShift) {
			start = (baseCount & tableBlockMask) + 1
		}
		end := uint32(tableBlockLen)
		if k == lastCurBlk && (curCount>>tableBlockShift) == lastCurBlk {
			end = (curCount & tableBlockMask) + 1
		}
		for off := start; off < end; off++ {
			if v := blk[off].Load(); v != nil {
				if clearOwner != nil {
					clearOwner(k<<tableBlockShift|off, v)
				}
				blk[off].Store(nil)
			}
		}
	}

	if baseCount == 0 {
		t.blk0.Store(nil)
		lastTop := lastCurBlk >> tableSubShift
		for top := uint32(0); top <= lastTop && top < tableDirLen; top++ {
			t.dir[top].Store(nil)
		}
	} else {
		lastBaseBlk := baseCount >> tableBlockShift
		lastBaseTop := lastBaseBlk >> tableSubShift
		if sub := t.dir[lastBaseTop].Load(); sub != nil {
			for subIdx := (lastBaseBlk & tableSubMask) + 1; subIdx < tableSubLen; subIdx++ {
				sub[subIdx].Store(nil)
			}
		}
		lastTop := lastCurBlk >> tableSubShift
		for top := lastBaseTop + 1; top <= lastTop && top < tableDirLen; top++ {
			t.dir[top].Store(nil)
		}
	}
	t.count.Store(baseCount)
}

func (t *storeTable[T]) rewind(baseCount uint32) {
	curCount := t.count.Load()
	if curCount <= baseCount {
		return
	}
	lastCurChunk := (curCount - 1) >> storeChunkShift
	lastTop := min(lastCurChunk>>storeDirShift, uint32(storeTopLen-1))
	if baseCount == 0 {
		for top := uint32(0); top <= lastTop; top++ {
			t.dir[top].Store(nil)
		}
		t.count.Store(0)
		return
	}

	lastBaseChunk := (baseCount - 1) >> storeChunkShift
	lastBaseTop := lastBaseChunk >> storeDirShift
	lastBaseMid := lastBaseChunk & storeDirMask
	if baseCount&storeChunkMask != 0 {
		if mid := t.dir[lastBaseTop].Load(); mid != nil {
			if chunk := mid[lastBaseMid].Load(); chunk != nil {
				end := uint32(storeChunkLen)
				if lastCurChunk == lastBaseChunk {
					end = ((curCount - 1) & storeChunkMask) + 1
				}
				clear(chunk[baseCount&storeChunkMask : end])
			}
		}
	}
	if mid := t.dir[lastBaseTop].Load(); mid != nil {
		for m := lastBaseMid + 1; m < storeDirLen; m++ {
			mid[m].Store(nil)
		}
	}
	for top := lastBaseTop + 1; top <= lastTop; top++ {
		t.dir[top].Store(nil)
	}
	t.count.Store(baseCount)
}

func (t *wideNumbers) rewind(baseNext uint32) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.next <= baseNext {
		return
	}
	const chunkShift = 12
	const chunkLen = 1 << chunkShift
	const chunkMask = chunkLen - 1

	if baseNext == 0 {
		t.index = nil
		lastChunk := (t.next - 1) >> chunkShift
		for c := uint32(0); c <= lastChunk && int(c) < len(t.chunks); c++ {
			t.chunks[c].Store(nil)
		}
		t.next = 0
		return
	}

	prevMapLen := len(t.index)
	if prevMapLen > int(baseNext)*2+64 {
		rebuilt := make(map[uint64]uint32, baseNext)
		for i := uint32(0); i < baseNext; i++ {
			rebuilt[math.Float64bits(t.at(i))] = i
		}
		t.index = rebuilt
	} else {
		for i := baseNext; i < t.next; i++ {
			delete(t.index, math.Float64bits(t.at(i)))
		}
	}

	lastBaseChunk := (baseNext - 1) >> chunkShift
	lastCurChunk := (t.next - 1) >> chunkShift
	if baseNext&chunkMask != 0 {
		if c := t.chunks[lastBaseChunk].Load(); c != nil {
			end := uint32(chunkLen)
			if lastCurChunk == lastBaseChunk {
				end = ((t.next - 1) & chunkMask) + 1
			}
			clear(c[baseNext&chunkMask : end])
		}
	}
	for c := lastBaseChunk + 1; c <= lastCurChunk && int(c) < len(t.chunks); c++ {
		t.chunks[c].Store(nil)
	}
	t.next = baseNext
}
