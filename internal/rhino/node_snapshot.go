package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/colors"
)

// Go-only: in-memory NodeStore snapshot and rewind.
//
// NodeStoreSnapshot captures a complete, immutable in-memory checkpoint of a
// non-static NodeStore — including node records, JSType/Color columns, paged
// original-name columns, cold auxiliary boxes and unshared property chains,
// JSDocInfo type-expression node handles, and the store's type/color handle
// tables. RestoreSnapshot rewinds the store back to that exact checkpoint
// in place, reusing already-allocated block, column, page, and shard-map
// storage across repeated rewinds (even if AST compaction ran in between).

type snapNodeBlock struct {
	present bool
	data    []nodeData
}

type snapNodeColumn struct {
	present bool
	data    []uint32
}

type snapColumnPage struct {
	idx  uint16
	data columnPage
}

type snapPagedColumn struct {
	pages []snapColumnPage
}

type snapAuxEntry struct {
	id  uint32
	aux nodeAux
}

type snapJSDocExpr struct {
	expr *JSTypeExpression
	root Node
}

type snapJSDocTransform struct {
	info  *JSDocInfo
	nodes []Node
}

// NodeStoreSnapshot is a reusable in-memory checkpoint of a NodeStore.
type NodeStoreSnapshot struct {
	storeID uint32
	base    uint32
	next    uint32

	blocks    []snapNodeBlock
	tcs       []snapNodeColumn
	origNames []*snapPagedColumn
	aux       [auxShards][]snapAuxEntry

	extCount          int
	unsharedPropCount int

	jsDocExprs      []snapJSDocExpr
	jsDocTransforms []snapJSDocTransform

	types  []*JSTypeBase
	colors []*colors.Color
}

func usedSlotsInBlock(blockIdx, totalNodes int) int {
	rem := totalNodes - blockIdx*blockLen
	if rem <= 0 {
		return 0
	}
	if rem >= blockLen {
		return blockLen
	}
	return rem
}

// Snapshot captures an immutable in-memory checkpoint of s.
// Call between passes only (when no concurrent goroutine is mutating s).
func (s *NodeStore) Snapshot() *NodeStoreSnapshot {
	if s == nil || s.static {
		panic("rhino.NodeStore.Snapshot: cannot snapshot nil or static store")
	}

	s.mu.Lock()
	base := s.base
	next := s.next.Load()
	totalNodes := 0
	if base != 0 && next > base {
		totalNodes = int(next - base)
	}

	blocks := make([]snapNodeBlock, len(s.blocks))
	for i, src := range s.blocks {
		if src != nil {
			used := usedSlotsInBlock(i, totalNodes)
			data := make([]nodeData, used)
			copy(data, src[:used])
			blocks[i] = snapNodeBlock{present: true, data: data}
		}
	}

	var tcs []snapNodeColumn
	if len(s.tcs) > 0 {
		tcs = make([]snapNodeColumn, len(s.tcs))
		for i, src := range s.tcs {
			if src != nil {
				used := usedSlotsInBlock(i, totalNodes)
				data := make([]uint32, used)
				copy(data, src[:used])
				tcs[i] = snapNodeColumn{present: true, data: data}
			}
		}
	}

	var origNames []*snapPagedColumn
	if len(s.origNames) > 0 {
		origNames = make([]*snapPagedColumn, len(s.origNames))
		for i, src := range s.origNames {
			if src == nil {
				continue
			}
			var pages []snapColumnPage
			for pIdx := range src {
				if p := src[pIdx].Load(); p != nil {
					pages = append(pages, snapColumnPage{
						idx:  uint16(pIdx),
						data: *p,
					})
				}
			}
			origNames[i] = &snapPagedColumn{pages: pages}
		}
	}

	types := snapshotHandleTable(&s.types)
	colorsSlice := snapshotHandleTable(&s.colors)
	s.mu.Unlock()

	snap := &NodeStoreSnapshot{
		storeID:   s.id,
		base:      base,
		next:      next,
		blocks:    blocks,
		tcs:       tcs,
		origNames: origNames,
		types:     types,
		colors:    colorsSlice,
	}

	// First count ext boxes and unshared prop items across shards so the snapshot
	// stores them in two contiguous backing slices instead of thousands of small allocs.
	extCount := 0
	unsharedPropCount := 0
	for sh := range s.aux {
		shard := &s.aux[sh]
		shard.mu.RLock()
		// order-independent: counting ext boxes and unshared prop list items across map entries
		for _, a := range shard.m {
			if a.ext != nil {
				extCount++
				for cur := a.ext.propListHead; cur != nil && !cur.shared; cur = cur.next {
					unsharedPropCount++
				}
			}
		}
		shard.mu.RUnlock()
	}
	snap.extCount = extCount
	snap.unsharedPropCount = unsharedPropCount

	var extPool []nodeExtBox
	if extCount > 0 {
		extPool = make([]nodeExtBox, extCount)
	}
	var propPool []Node_PropListItem
	if unsharedPropCount > 0 {
		propPool = make([]Node_PropListItem, unsharedPropCount)
	}
	extIdx := 0
	propIdx := 0

	var seenDocs map[*JSDocInfo]struct{}
	var seenExprs map[*JSTypeExpression]struct{}
	recordDoc := func(info *JSDocInfo) {
		if info == nil {
			return
		}
		if seenDocs == nil {
			seenDocs = make(map[*JSDocInfo]struct{})
		}
		if _, done := seenDocs[info]; done {
			return
		}
		seenDocs[info] = struct{}{}

		for it := info.TypeExprIter(); it.Next(); {
			e := it.Expr()
			if e == nil || e == JSTypeExpression_IMPLICIT_TEMPLATE_BOUND {
				continue
			}
			if !e.root.IsNil() && uint32(e.root)&staticBit != 0 {
				continue
			}
			if seenExprs == nil {
				seenExprs = make(map[*JSTypeExpression]struct{})
			}
			if _, exists := seenExprs[e]; !exists {
				seenExprs[e] = struct{}{}
				snap.jsDocExprs = append(snap.jsDocExprs, snapJSDocExpr{
					expr: e,
					root: e.root,
				})
			}
		}

		if transforms := info.GetTypeTransformationEntries(); len(transforms) > 0 {
			nodes := make([]Node, len(transforms))
			for idx := range transforms {
				nodes[idx] = transforms[idx].TypeNode
			}
			snap.jsDocTransforms = append(snap.jsDocTransforms, snapJSDocTransform{
				info:  info,
				nodes: nodes,
			})
		}
	}

	for sh := range s.aux {
		shard := &s.aux[sh]
		shard.mu.RLock()
		if len(shard.m) > 0 {
			entries := make([]snapAuxEntry, 0, len(shard.m))
			// order-independent: copying shard map entries into snapshot; unique ids and deduplicated JSDoc pointers have no order dependence
			for id, a := range shard.m {
				copied := nodeAux{jsDoc: a.jsDoc}
				recordDoc(a.jsDoc)
				if a.ext != nil {
					boxCopy := &extPool[extIdx]
					extIdx++
					*boxCopy = *a.ext
					boxCopy.propListHead, propIdx = clonePropListUnsharedInto(a.ext.propListHead, propPool, propIdx)
					copied.ext = boxCopy
					recordDoc(a.ext.jsDocInfo)
					for item := a.ext.propListHead; item != nil; item = item.next {
						if item.isObject {
							if doc, ok := item.objectValue.(*JSDocInfo); ok {
								recordDoc(doc)
							}
						}
					}
				}
				entries = append(entries, snapAuxEntry{id: id, aux: copied})
			}
			snap.aux[sh] = entries
		}
		shard.mu.RUnlock()
	}

	return snap
}

// RestoreSnapshot rewinds s to the state captured in snap.
// snap remains immutable and may be restored into s repeatedly.
// Call between passes only (when no concurrent goroutine is accessing s).
func (s *NodeStore) RestoreSnapshot(snap *NodeStoreSnapshot) {
	if s == nil || s.static {
		panic("rhino.NodeStore.RestoreSnapshot: cannot restore into nil or static store")
	}
	if snap == nil {
		panic("rhino.NodeStore.RestoreSnapshot: nil snapshot")
	}
	if s.id != snap.storeID {
		panic("rhino.NodeStore.RestoreSnapshot: snapshot belongs to a different NodeStore")
	}

	wasCurrent := current == s
	if s.base != 0 {
		raiseNextBlock(s.blockEnd())
		if wasCurrent {
			s.unregister()
		}
	}

	s.mu.Lock()
	sameBase := s.base == snap.base && s.base != 0
	prevTotalNodes := 0
	if sameBase {
		if prevNext := s.next.Load(); prevNext > s.base {
			prevTotalNodes = int(prevNext - s.base)
		}
	}

	// Restore node record blocks, reusing existing 576 KiB block allocations
	// and touching only slots that were live in snap or dirtied since snap.
	s.blocks = resizePtrSlice(s.blocks, len(snap.blocks))
	for i, sb := range snap.blocks {
		if !sb.present {
			s.blocks[i] = nil
			continue
		}
		dst := s.blocks[i]
		if dst == nil {
			dst = new(nodeBlock)
			copy(dst[:len(sb.data)], sb.data)
			s.blocks[i] = dst
			continue
		}
		copy(dst[:len(sb.data)], sb.data)
		dirtyEnd := blockLen
		if sameBase {
			dirtyEnd = usedSlotsInBlock(i, prevTotalNodes)
		}
		if dirtyEnd > len(sb.data) {
			clear(dst[len(sb.data):dirtyEnd])
		}
	}

	// Restore JSType/Color columns, reusing existing 64 KiB column allocations.
	s.tcs = resizePtrSlice(s.tcs, len(snap.tcs))
	for i, sc := range snap.tcs {
		if !sc.present {
			s.tcs[i] = nil
			continue
		}
		dst := s.tcs[i]
		if dst == nil {
			dst = new(nodeColumn)
			copy(dst[:len(sc.data)], sc.data)
			s.tcs[i] = dst
			continue
		}
		copy(dst[:len(sc.data)], sc.data)
		dirtyEnd := blockLen
		if sameBase {
			dirtyEnd = usedSlotsInBlock(i, prevTotalNodes)
		}
		if dirtyEnd > len(sc.data) {
			clear(dst[len(sc.data):dirtyEnd])
		}
	}

	// Restore paged original-name columns, reusing allocated pagedColumns and columnPages.
	s.origNames = resizePtrSlice(s.origNames, len(snap.origNames))
	for i, src := range snap.origNames {
		if src == nil {
			s.origNames[i] = nil
			continue
		}
		dst := s.origNames[i]
		if dst == nil {
			dst = new(pagedColumn)
			s.origNames[i] = dst
		}
		spIdx := 0
		for pIdx := range dst {
			if spIdx < len(src.pages) && int(src.pages[spIdx].idx) == pIdx {
				p := dst[pIdx].Load()
				if p == nil {
					p = new(columnPage)
					*p = src.pages[spIdx].data
					dst[pIdx].Store(p)
				} else {
					*p = src.pages[spIdx].data
				}
				spIdx++
			} else if dst[pIdx].Load() != nil {
				dst[pIdx].Store(nil)
			}
		}
	}

	restoreHandleTable(&s.types, snap.types)
	restoreHandleTable(&s.colors, snap.colors)

	if snap.base == 0 && wasCurrent {
		s.setBase(nextBlock.Load() << blockShift)
	} else {
		s.base = snap.base
		s.next.Store(snap.next)
	}
	s.mu.Unlock()

	// Restore cached node handles on types and colors so future SetJSType / SetColor
	// lookups hit their original table slots directly.
	var lastReg *JSTypeRegistry
	for h := 1; h < len(snap.types); h++ {
		if t := snap.types[h]; t != nil {
			t.nodeHandle.Store(uint32(h))
			if reg := t.GetRegistry(); reg != nil && reg != lastReg {
				reg.noteNodeStore(s)
				lastReg = reg
			}
		}
	}
	for h := 1; h < len(snap.colors); h++ {
		if c := snap.colors[h]; c != nil {
			c.StoreNodeHandle(s.id, uint32(h))
		}
	}

	// Restore cold auxiliary shards using two contiguous pool allocations for all
	// ext boxes and unshared property list items so subsequent passes never mutate snap.
	var extPool []nodeExtBox
	if snap.extCount > 0 {
		extPool = make([]nodeExtBox, snap.extCount)
	}
	var propPool []Node_PropListItem
	if snap.unsharedPropCount > 0 {
		propPool = make([]Node_PropListItem, snap.unsharedPropCount)
	}
	extIdx := 0
	propIdx := 0

	for sh := range s.aux {
		shard := &s.aux[sh]
		entries := snap.aux[sh]
		shard.mu.Lock()
		if len(entries) == 0 {
			if len(shard.m) > 0 {
				clear(shard.m)
			}
		} else {
			if shard.m == nil {
				shard.m = make(map[uint32]nodeAux, len(entries))
			} else {
				clear(shard.m)
			}
			for _, e := range entries {
				restored := nodeAux{jsDoc: e.aux.jsDoc}
				if e.aux.ext != nil {
					boxCopy := &extPool[extIdx]
					extIdx++
					*boxCopy = *e.aux.ext
					boxCopy.propListHead, propIdx = clonePropListUnsharedInto(e.aux.ext.propListHead, propPool, propIdx)
					restored.ext = boxCopy
				}
				shard.m[e.id] = restored
			}
		}
		shard.mu.Unlock()
	}

	// Restore JSDocInfo type-expression and type-transformation Node handles
	// in case Compactor forwarded them during the rewound run.
	for _, je := range snap.jsDocExprs {
		je.expr.root = je.root
	}
	for _, jt := range snap.jsDocTransforms {
		transforms := jt.info.GetTypeTransformationEntries()
		for idx := range jt.nodes {
			if idx < len(transforms) {
				transforms[idx].TypeNode = jt.nodes[idx]
			}
		}
	}

	if wasCurrent {
		s.register()
	}
}

// clonePropListUnsharedInto deep-copies the unshared prefix of head into pool
// starting at poolIdx, stopping as soon as it reaches a shared tail item.
func clonePropListUnsharedInto(head *Node_PropListItem, pool []Node_PropListItem, poolIdx int) (*Node_PropListItem, int) {
	if head == nil || head.shared {
		return head, poolIdx
	}
	headCopy := &pool[poolIdx]
	poolIdx++
	*headCopy = *head
	dst := headCopy
	for cur := head.next; cur != nil && !cur.shared; cur = cur.next {
		itemCopy := &pool[poolIdx]
		poolIdx++
		*itemCopy = *cur
		dst.next = itemCopy
		dst = itemCopy
	}
	return headCopy, poolIdx
}

func resizePtrSlice[T any](s []*T, n int) []*T {
	if len(s) > n {
		clear(s[n:])
		return s[:n]
	}
	if len(s) < n {
		if cap(s) >= n {
			return s[:n]
		}
		grown := make([]*T, n)
		copy(grown, s)
		return grown
	}
	return s
}

func snapshotHandleTable[T any](t *handleTable[T]) []T {
	next := int(t.next.Load())
	if next <= 1 {
		return nil
	}
	out := make([]T, next)
	for h := 1; h < next; {
		blkIdx := h >> tableBlockShift
		off := h & tableBlockMask
		end := min(next, (blkIdx+1)<<tableBlockShift)
		if src := t.blocks[blkIdx].Load(); src != nil {
			copy(out[h:end], src[off:off+(end-h)])
		}
		h = end
	}
	return out
}

func restoreHandleTable[T any](t *handleTable[T], snapEntries []T) {
	snapNext := uint32(1)
	if len(snapEntries) > 1 {
		snapNext = uint32(len(snapEntries))
	}
	curNext := int(t.next.Load())
	curBlocks := 0
	if curNext > 1 {
		curBlocks = ((curNext - 1) >> tableBlockShift) + 1
	}
	snapBlocks := 0
	if snapNext > 1 {
		snapBlocks = int(((snapNext - 1) >> tableBlockShift) + 1)
	}

	for k := snapBlocks; k < curBlocks; k++ {
		if k < len(t.blocks) {
			t.blocks[k].Store(nil)
		}
	}

	for k := range snapBlocks {
		start := k << tableBlockShift
		end := min(int(snapNext), start+tableBlockLen)
		count := end - start

		dst := t.blocks[k].Load()
		if dst == nil {
			dst = new([tableBlockLen]T)
			copy(dst[:count], snapEntries[start:end])
			t.blocks[k].Store(dst)
			continue
		}
		copy(dst[:count], snapEntries[start:end])
		if count < tableBlockLen {
			dirtyEnd := min(tableBlockLen, max(count, curNext-start))
			if dirtyEnd > count {
				clear(dst[count:dirtyEnd])
			}
		}
	}

	t.next.Store(snapNext)
}
