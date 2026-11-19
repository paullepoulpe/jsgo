package protobuf

import (
	"encoding/binary"
	"io"
	"math"
	"sync"
	"sync/atomic"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
)

// ---------------------------------------------------------------------------
// TemplateStringValue & AstNode (com.google.javascript.jscomp.serialization.*)
// ---------------------------------------------------------------------------

type TemplateStringValue struct {
	rawStringPointer_    int32
	cookedStringPointer_ int32
	cachedSize           int32
}

type TemplateStringValue_Builder = TemplateStringValue

var templateStringValueDefault = &TemplateStringValue{}

func TemplateStringValue_NewBuilder() *TemplateStringValue_Builder {
	return &TemplateStringValue{}
}

func (m *TemplateStringValue) GetRawStringPointer() int32    { return m.rawStringPointer_ }
func (m *TemplateStringValue) GetCookedStringPointer() int32 { return m.cookedStringPointer_ }

func (m *TemplateStringValue) SetRawStringPointer(v int32) *TemplateStringValue {
	m.rawStringPointer_ = v
	m.cachedSize = 0
	return m
}

func (m *TemplateStringValue) SetCookedStringPointer(v int32) *TemplateStringValue {
	m.cookedStringPointer_ = v
	m.cachedSize = 0
	return m
}

func (m *TemplateStringValue) Build() *TemplateStringValue {
	m.cachedSize = 0
	return m
}

func (m *TemplateStringValue) serializedSize() int {
	if m.cachedSize > 0 {
		return int(m.cachedSize - 1)
	}
	size := 0
	if m.rawStringPointer_ != 0 {
		size += 1 + protoVarintLen(uint64(uint32(m.rawStringPointer_)))
	}
	if m.cookedStringPointer_ != 0 {
		size += 1 + protoVarintLen(uint64(int64(m.cookedStringPointer_)))
	}
	if size < math.MaxInt32 {
		m.cachedSize = int32(size + 1)
	}
	return size
}

func (m *TemplateStringValue) appendTo(b []byte) []byte {
	if m.rawStringPointer_ != 0 {
		b = protoAppendTag(b, 1, protoWireVarint)
		b = protoAppendUint32(b, m.rawStringPointer_)
	}
	if m.cookedStringPointer_ != 0 {
		b = protoAppendTag(b, 2, protoWireVarint)
		b = protoAppendInt32(b, m.cookedStringPointer_)
	}
	return b
}

func (m *TemplateStringValue) mergeFromReader(r *protoReader) {
	for {
		tag := r.readTag()
		switch tag {
		case 0:
			return
		case 8:
			m.rawStringPointer_ = r.readUint32()
		case 16:
			m.cookedStringPointer_ = r.readInt32()
		default:
			r.skipField(tag)
		}
	}
}

type AstNode_ValueCase int32

const (
	AstNode_ValueCase_VALUE_NOT_SET         AstNode_ValueCase = 0
	AstNode_ValueCase_STRING_VALUE_POINTER  AstNode_ValueCase = 3
	AstNode_ValueCase_DOUBLE_VALUE          AstNode_ValueCase = 4
	AstNode_ValueCase_TEMPLATE_STRING_VALUE AstNode_ValueCase = 8
)

const (
	astNodeFlagValueMask   uint8 = 0x0F
	astNodeFlagHasType     uint8 = 1 << 4
	astNodeFlagHasTemplate uint8 = 1 << 5
	astNodeFlagBuilt       uint8 = 1 << 6
)

const astAuxOverflowSentinel uint16 = math.MaxUint16

func packTemplateStringBits(raw, cooked int32) uint64 {
	return uint64(uint32(raw)) | (uint64(uint32(cooked)) << 32)
}

func unpackTemplateStringBits(bits uint64) (raw, cooked int32) {
	return int32(uint32(bits)), int32(uint32(bits >> 32))
}

func templateStringSerializedSize(raw, cooked int32) int {
	size := 0
	if raw != 0 {
		size += 1 + protoVarintLen(uint64(uint32(raw)))
	}
	if cooked != 0 {
		size += 1 + protoVarintLen(uint64(int64(cooked)))
	}
	return size
}

// astNodeRec is a 36-byte pointer-free (100% noscan) record for deserialized
// AstNodes stored in astNodeSlab.recs. Using 32-bit aligned scalar fields
// eliminates both alignment padding and all GC pointer scanning across millions
// of deserialized AST nodes.
type astNodeRec struct {
	boolPropsLo    uint32
	boolPropsHi    uint32
	childOff       uint32
	childCount     uint32
	value32        uint32
	relativeLine   int32
	relativeColumn int32
	type_          int32
	auxIdx         uint16
	kind           uint8
	flags          uint8
}

func (r *astNodeRec) valueCase() AstNode_ValueCase {
	return AstNode_ValueCase(r.flags & astNodeFlagValueMask)
}

func (r *astNodeRec) setValueCase(vc AstNode_ValueCase, hasTemplate bool) {
	f := (r.flags &^ (astNodeFlagValueMask | astNodeFlagHasTemplate)) | (uint8(vc) & astNodeFlagValueMask)
	if hasTemplate {
		f |= astNodeFlagHasTemplate
	}
	r.flags = f
}

func (r *astNodeRec) boolProps() int64 {
	return int64(uint64(r.boolPropsLo) | (uint64(r.boolPropsHi) << 32))
}

// astNodeAux holds cold/sparse AstNode fields (present on <5% of AST nodes).
// 12 bytes, 100% pointer-free (noscan).
type astNodeAux struct {
	originalNamePointer int32
	sourceFile          int32
	jsdocIdx            uint32 // 0 = none; 1..N = 1-based index in slab.jsdocs
}

// astNodeData backs standalone/builder AstNode instances (when slab == nil) or
// serves as the single shared slab descriptor embedded in astNodeSlab (when
// slab != nil). Pointer fields are packed at the head so Go's GC skips the
// scalar tail. Unpinned single-parent children are linked intrusively via
// firstChild/lastChild/nextSibling without allocating []*AstNode slices.
type astNodeData struct {
	slab                 *astNodeSlab
	bundle               *astBuilderBundle
	selfNode             *AstNode
	firstChild           *astNodeData
	lastChild            *astNodeData
	nextSibling          *astNodeData
	child_               []*AstNode
	jsdoc_               *OptimizationJsdoc
	valueBits_           uint64
	booleanProperties_   int64
	relativeLine_        int32
	relativeColumn_      int32
	originalNamePointer_ int32
	type_                int32
	sourceFile_          int32
	cachedSize           int32
	childCount_          uint32
	linkState_           uint32
	kind_                uint8
	flags_               uint8
}

func (d *astNodeData) valueCase() AstNode_ValueCase {
	return AstNode_ValueCase(d.flags_ & astNodeFlagValueMask)
}

func (d *astNodeData) setValueCase(vc AstNode_ValueCase, hasTemplate bool) {
	f := (d.flags_ &^ (astNodeFlagValueMask | astNodeFlagHasTemplate)) | (uint8(vc) & astNodeFlagValueMask)
	if hasTemplate {
		f |= astNodeFlagHasTemplate
	}
	d.flags_ = f
}

func (d *astNodeData) isZero() bool {
	return d.slab == nil && d.childCount_ == 0 && len(d.child_) == 0 && d.jsdoc_ == nil &&
		d.valueBits_ == 0 && d.booleanProperties_ == 0 &&
		d.relativeLine_ == 0 && d.relativeColumn_ == 0 &&
		d.originalNamePointer_ == 0 && d.type_ == 0 && d.sourceFile_ == 0 &&
		d.kind_ == 0 && d.flags_ == 0
}

func pinAstData(d *astNodeData) {
	if atomic.LoadUint32(&d.linkState_)&astLinkFlagPinned == 0 {
		atomic.OrUint32(&d.linkState_, astLinkFlagPinned)
	}
}

func (d *astNodeData) ensureChildSlice() {
	if d.child_ != nil {
		return
	}
	n := int(d.childCount_)
	if n == 0 {
		d.child_ = make([]*AstNode, 0, 4)
		return
	}
	ch := make([]*AstNode, n, max(n, 4))
	i := 0
	for c := d.firstChild; c != nil; c = c.nextSibling {
		ch[i] = c.selfNode
		i++
	}
	d.child_ = ch
	d.firstChild = nil
	d.lastChild = nil
}

// AstNode is a 16-byte handle pointing either to a shared astNodeSlab (for
// deserialized ASTs, indexing pointer-free astNodeRecs by idx) or to a mutable
// builder astNodeData (when d.slab == nil).
type AstNode struct {
	d   *astNodeData
	idx uint32
}

const (
	astLinkFlagLinked uint32 = 1 << 0
	astLinkFlagPinned uint32 = 1 << 1

	astBuilderChunkSize         = 256
	astBuilderMaxFreeBundles    = 262144
	astBuilderMaxPooledChildren = 4096
)

type astBuilderBundle struct {
	node     AstNode
	data     astNodeData
	poolNext *astBuilderBundle
}

type astBuilderAllocBatch struct {
	idx  uint32
	done uint32
	ptrs [astBuilderChunkSize]*astBuilderBundle
}

type astBuilderChunk struct {
	items [astBuilderChunkSize]astBuilderBundle
}

type astDataChunk struct {
	items [astBuilderChunkSize]astNodeData
}

var (
	astBuilderBatchPtr   atomic.Pointer[astBuilderAllocBatch]
	astBuilderPoolMu     sync.Mutex
	astBuilderCurBatch   *astBuilderAllocBatch
	astBuilderBatchFree  []*astBuilderAllocBatch
	astBuilderRetired    []*astBuilderAllocBatch
	astBuilderFree       []*astBuilderBundle

	astDataChunkMu  sync.Mutex
	astDataChunkCur *astDataChunk
	astDataChunkPos int
)

func allocStandaloneAstData(self *AstNode) *astNodeData {
	astDataChunkMu.Lock()
	if astDataChunkCur == nil || astDataChunkPos >= astBuilderChunkSize {
		astDataChunkCur = &astDataChunk{}
		astDataChunkPos = 0
	}
	nd := &astDataChunkCur.items[astDataChunkPos]
	astDataChunkPos++
	astDataChunkMu.Unlock()
	nd.selfNode = self
	return nd
}

//go:noinline
func newAstNodeBuilderSlow() *AstNode {
	astBuilderPoolMu.Lock()
	if batch := astBuilderBatchPtr.Load(); batch != nil {
		if i := atomic.AddUint32(&batch.idx, 1) - 1; i < astBuilderChunkSize {
			b := batch.ptrs[i]
			atomic.AddUint32(&batch.done, 1)
			astBuilderPoolMu.Unlock()
			return &b.node
		}
	}
	if old := astBuilderCurBatch; old != nil {
		astBuilderCurBatch = nil
		astBuilderBatchPtr.Store(nil)
		if atomic.LoadUint32(&old.done) == astBuilderChunkSize {
			clear(old.ptrs[:])
			if len(astBuilderBatchFree) < 16 {
				astBuilderBatchFree = append(astBuilderBatchFree, old)
			}
		} else {
			astBuilderRetired = append(astBuilderRetired, old)
		}
	}
	if len(astBuilderRetired) > 0 {
		w := 0
		for _, rb := range astBuilderRetired {
			if atomic.LoadUint32(&rb.done) == astBuilderChunkSize {
				clear(rb.ptrs[:])
				if len(astBuilderBatchFree) < 16 {
					astBuilderBatchFree = append(astBuilderBatchFree, rb)
				}
			} else {
				astBuilderRetired[w] = rb
				w++
			}
		}
		clear(astBuilderRetired[w:])
		astBuilderRetired = astBuilderRetired[:w]
	}
	var nb *astBuilderAllocBatch
	if nbf := len(astBuilderBatchFree); nbf > 0 {
		nb = astBuilderBatchFree[nbf-1]
		astBuilderBatchFree[nbf-1] = nil
		astBuilderBatchFree = astBuilderBatchFree[:nbf-1]
	} else {
		nb = &astBuilderAllocBatch{idx: astBuilderChunkSize}
	}
	nFree := len(astBuilderFree)
	if nFree >= astBuilderChunkSize {
		start := nFree - astBuilderChunkSize
		copy(nb.ptrs[:], astBuilderFree[start:])
		clear(astBuilderFree[start:])
		astBuilderFree = astBuilderFree[:start]
	} else {
		copied := 0
		if nFree > 0 {
			copied = copy(nb.ptrs[:], astBuilderFree)
			clear(astBuilderFree)
			astBuilderFree = astBuilderFree[:0]
		}
		nc := &astBuilderChunk{}
		need := astBuilderChunkSize - copied
		for i := 0; i < need; i++ {
			b := &nc.items[i]
			b.node.d = &b.data
			b.data.bundle = b
			b.data.selfNode = &b.node
			nb.ptrs[copied+i] = b
		}
		for i := need; i < astBuilderChunkSize; i++ {
			b := &nc.items[i]
			b.node.d = &b.data
			b.data.bundle = b
			b.data.selfNode = &b.node
			astBuilderFree = append(astBuilderFree, b)
		}
	}
	ret := nb.ptrs[0]
	atomic.StoreUint32(&nb.done, 1)
	atomic.StoreUint32(&nb.idx, 1)
	astBuilderCurBatch = nb
	astBuilderBatchPtr.Store(nb)
	astBuilderPoolMu.Unlock()
	return &ret.node
}

func (d *astNodeData) resetBuilderFields() {
	if d.firstChild != nil {
		d.firstChild = nil
		d.lastChild = nil
	}
	if d.nextSibling != nil {
		d.nextSibling = nil
	}
	if len(d.child_) > 0 {
		clear(d.child_)
		if cap(d.child_) <= astBuilderMaxPooledChildren {
			d.child_ = d.child_[:0]
		} else {
			d.child_ = nil
		}
	} else if cap(d.child_) > astBuilderMaxPooledChildren {
		d.child_ = nil
	}
	if d.jsdoc_ != nil {
		d.jsdoc_ = nil
	}
	d.valueBits_ = 0
	d.booleanProperties_ = 0
	d.relativeLine_ = 0
	d.relativeColumn_ = 0
	d.originalNamePointer_ = 0
	d.type_ = 0
	d.sourceFile_ = 0
	d.cachedSize = 0
	d.childCount_ = 0
	d.kind_ = 0
	d.flags_ = 0
	atomic.StoreUint32(&d.linkState_, 0)
}

func recycleAstBuilderTree(root *astBuilderBundle) {
	if atomic.LoadUint32(&root.data.linkState_) != 0 {
		return
	}
	root.poolNext = nil
	workHead := root
	var freeHead *astBuilderBundle
	for workHead != nil {
		b := workHead
		workHead = b.poolNext
		b.poolNext = nil
		d := &b.data
		for c := d.firstChild; c != nil; {
			next := c.nextSibling
			if cb := c.bundle; cb != nil && c == &cb.data && atomic.LoadUint32(&c.linkState_) == astLinkFlagLinked {
				cb.poolNext = workHead
				workHead = cb
			}
			c = next
		}
		d.firstChild = nil
		d.lastChild = nil
		d.nextSibling = nil
		if children := d.child_; children != nil {
			for i, ch := range children {
				children[i] = nil
				if ch != nil {
					if cd := ch.d; cd != nil && cd.bundle != nil && ch == &cd.bundle.node && atomic.LoadUint32(&cd.linkState_) == astLinkFlagLinked {
						cd.bundle.poolNext = workHead
						workHead = cd.bundle
					}
				}
			}
			if cap(children) > 0 && cap(children) <= astBuilderMaxPooledChildren {
				d.child_ = children[:0]
			} else {
				d.child_ = nil
			}
		}
		if d.jsdoc_ != nil {
			d.jsdoc_ = nil
		}
		d.valueBits_ = 0
		d.booleanProperties_ = 0
		d.relativeLine_ = 0
		d.relativeColumn_ = 0
		d.originalNamePointer_ = 0
		d.type_ = 0
		d.sourceFile_ = 0
		d.cachedSize = 0
		d.childCount_ = 0
		d.kind_ = 0
		d.flags_ = 0
		atomic.StoreUint32(&d.linkState_, 0)
		if b.node.d != d {
			b.node.d = d
		}
		if d.selfNode != &b.node {
			d.selfNode = &b.node
		}
		b.node.idx = 0
		b.poolNext = freeHead
		freeHead = b
	}
	if freeHead == nil {
		return
	}
	astBuilderPoolMu.Lock()
	for freeHead != nil {
		b := freeHead
		freeHead = b.poolNext
		b.poolNext = nil
		if len(astBuilderFree) < astBuilderMaxFreeBundles {
			astBuilderFree = append(astBuilderFree, b)
		}
	}
	astBuilderPoolMu.Unlock()
}

const (
	astSlabMaxPooledRecs     = 262144
	astSlabMaxPooledChildren = 524288
)

type astNodeSlab struct {
	shared        astNodeData
	root          AstNode
	recs          []astNodeRec
	childIdxs     []uint32
	vals64        []uint64
	tsvs          []TemplateStringValue
	aux           []astNodeAux
	jsdocs        []OptimizationJsdoc
	scratch       []uint32
	handles       []AstNode
	cachedSizes   []uint32
	auxOverflow   map[uint32]astNodeAux
	extCache      any
	handlesInit   int
	handlesMu     sync.Mutex
	handlesReady  uint32
	extCacheReady uint32
	sizesOnce     sync.Once
	inPool        bool
	mutated       bool
}

var astNodeSlabPool = sync.Pool{
	New: func() any {
		s := &astNodeSlab{}
		s.shared.slab = s
		s.shared.flags_ = astNodeFlagBuilt
		s.root = AstNode{d: &s.shared, idx: 0}
		return s
	},
}

func newAstNodeSlab(remBytes int) *astNodeSlab {
	s := astNodeSlabPool.Get().(*astNodeSlab)
	s.inPool = false
	s.shared.slab = s
	s.shared.flags_ = astNodeFlagBuilt
	s.root = AstNode{d: &s.shared, idx: 0}
	if remBytes <= 0 {
		if cap(s.recs) < 4 {
			s.recs = make([]astNodeRec, 0, 4)
		} else {
			s.recs = s.recs[:0]
		}
		if s.childIdxs != nil {
			s.childIdxs = s.childIdxs[:0]
		}
		return s
	}
	// Closure TypedAst averages ~13-15 wire bytes per AstNode; remBytes/12 + 4
	// fits typical scripts in a single noscan slice without ~2x overallocation.
	expectedNodes := min(max(remBytes/12+4, 4), astSlabMaxPooledRecs)
	if cap(s.recs) < expectedNodes {
		newCap := expectedNodes
		if cap(s.recs) > 0 {
			newCap = max(expectedNodes, min(cap(s.recs)*2, astSlabMaxPooledRecs))
		}
		s.recs = make([]astNodeRec, 0, newCap)
	} else {
		s.recs = s.recs[:0]
	}
	if expectedNodes > 1 && cap(s.childIdxs) < expectedNodes-1 {
		newChildCap := expectedNodes - 1
		if cap(s.childIdxs) > 0 {
			newChildCap = max(expectedNodes-1, min(cap(s.childIdxs)*2, astSlabMaxPooledChildren))
		}
		s.childIdxs = make([]uint32, 0, newChildCap)
	} else if s.childIdxs != nil {
		s.childIdxs = s.childIdxs[:0]
	}
	return s
}

func (s *astNodeSlab) allocRec() uint32 {
	idx := uint32(len(s.recs))
	if len(s.recs) == cap(s.recs) {
		s.growRecs()
	}
	s.recs = s.recs[:idx+1]
	s.recs[idx] = astNodeRec{}
	return idx
}

//go:noinline
func (s *astNodeSlab) growRecs() {
	oldCap := cap(s.recs)
	newCap := oldCap + max(oldCap/2, 64)
	next := make([]astNodeRec, len(s.recs), newCap)
	copy(next, s.recs)
	s.recs = next
}

func (s *astNodeSlab) appendChildIdxs(ids []uint32) uint32 {
	off := uint32(len(s.childIdxs))
	need := len(s.childIdxs) + len(ids)
	if need > cap(s.childIdxs) {
		s.growChildIdxs(need)
	}
	s.childIdxs = append(s.childIdxs, ids...)
	return off
}

//go:noinline
func (s *astNodeSlab) growChildIdxs(need int) {
	oldCap := cap(s.childIdxs)
	newCap := max(oldCap+max(oldCap/2, 64), need)
	next := make([]uint32, len(s.childIdxs), newCap)
	copy(next, s.childIdxs)
	s.childIdxs = next
}

func (s *astNodeSlab) setRecDoubleBits(cur *astNodeRec, bits uint64) {
	if cur.valueCase() == AstNode_ValueCase_DOUBLE_VALUE {
		s.vals64[cur.value32] = bits
		return
	}
	if s.vals64 == nil {
		s.vals64 = make([]uint64, 0, min(max(cap(s.recs)/16, 8), 8192))
	}
	cur.value32 = uint32(len(s.vals64))
	s.vals64 = append(s.vals64, bits)
}

func (s *astNodeSlab) mergeRecTemplateString(cur *astNodeRec, r *protoReader) {
	var tv TemplateStringValue
	reuse := cur.valueCase() == AstNode_ValueCase_TEMPLATE_STRING_VALUE && cur.flags&astNodeFlagHasTemplate != 0
	if reuse {
		tv = s.tsvs[cur.value32]
		tv.cachedSize = 0
	}
	tv.mergeFromReader(r)
	if reuse {
		s.tsvs[cur.value32] = tv
		return
	}
	if s.tsvs == nil {
		s.tsvs = make([]TemplateStringValue, 0, min(max(cap(s.recs)/32, 4), 2048))
	}
	cur.value32 = uint32(len(s.tsvs))
	s.tsvs = append(s.tsvs, tv)
}

func (s *astNodeSlab) getAux(idx uint32, auxIdx uint16) astNodeAux {
	if auxIdx == 0 {
		return astNodeAux{}
	}
	if auxIdx < astAuxOverflowSentinel {
		return s.aux[auxIdx-1]
	}
	return s.auxOverflow[idx]
}

func (s *astNodeSlab) putAux(idx uint32, cur *astNodeRec, a astNodeAux) {
	if cur.auxIdx == 0 {
		if len(s.aux) < int(astAuxOverflowSentinel-1) {
			if s.aux == nil {
				s.aux = make([]astNodeAux, 0, min(max(cap(s.recs)/16, 4), 8192))
			}
			s.aux = append(s.aux, a)
			cur.auxIdx = uint16(len(s.aux))
			return
		}
		if s.auxOverflow == nil {
			s.auxOverflow = make(map[uint32]astNodeAux)
		}
		cur.auxIdx = astAuxOverflowSentinel
		s.auxOverflow[idx] = a
		return
	}
	if cur.auxIdx < astAuxOverflowSentinel {
		s.aux[cur.auxIdx-1] = a
		return
	}
	s.auxOverflow[idx] = a
}

func (s *astNodeSlab) setAuxOrigName(idx uint32, cur *astNodeRec, val int32) {
	if val == 0 && cur.auxIdx == 0 {
		return
	}
	a := s.getAux(idx, cur.auxIdx)
	a.originalNamePointer = val
	s.putAux(idx, cur, a)
}

func (s *astNodeSlab) setAuxSourceFile(idx uint32, cur *astNodeRec, val int32) {
	if val == 0 && cur.auxIdx == 0 {
		return
	}
	a := s.getAux(idx, cur.auxIdx)
	a.sourceFile = val
	s.putAux(idx, cur, a)
}

func (s *astNodeSlab) mergeAuxJsdoc(idx uint32, cur *astNodeRec, r *protoReader) {
	a := s.getAux(idx, cur.auxIdx)
	if a.jsdocIdx == 0 {
		if s.jsdocs == nil {
			s.jsdocs = make([]OptimizationJsdoc, 0, min(max(cap(s.recs)/128, 4), 512))
		}
		n := len(s.jsdocs)
		if n < cap(s.jsdocs) {
			s.jsdocs = s.jsdocs[:n+1]
			j := &s.jsdocs[n]
			var savedKind []int32
			if cap(j.kind_) <= 16 {
				savedKind = j.kind_[:0]
			}
			*j = OptimizationJsdoc{kind_: savedKind}
		} else {
			s.jsdocs = append(s.jsdocs, OptimizationJsdoc{})
		}
		a.jsdocIdx = uint32(len(s.jsdocs))
		s.putAux(idx, cur, a)
	} else {
		s.jsdocs[a.jsdocIdx-1].cachedSize = 0
	}
	s.jsdocs[a.jsdocIdx-1].mergeFromReader(r)
}

func (s *astNodeSlab) ensureHandlesLocked() {
	if atomic.LoadUint32(&s.handlesReady) != 0 {
		return
	}
	n := len(s.recs) - 1
	if n <= 0 {
		s.handles = s.handles[:0]
		atomic.StoreUint32(&s.handlesReady, 1)
		return
	}
	if cap(s.handles) < n {
		newCap := max(n, max(cap(s.recs)-1, cap(s.handles)*2))
		s.handles = make([]AstNode, n, newCap)
		s.handlesInit = 0
	} else {
		s.handles = s.handles[:n]
	}
	if s.handlesInit < n {
		shared := &s.shared
		h := s.handles
		for i := s.handlesInit; i < n; i++ {
			h[i] = AstNode{d: shared, idx: uint32(i + 1)}
		}
		s.handlesInit = n
	}
	atomic.StoreUint32(&s.handlesReady, 1)
}

//go:noinline
func (s *astNodeSlab) initHandlesSlow(idx uint32) *AstNode {
	if idx == 0 {
		return &s.root
	}
	s.handlesMu.Lock()
	s.ensureHandlesLocked()
	h := &s.handles[idx-1]
	s.handlesMu.Unlock()
	return h
}

func (s *astNodeSlab) getHandle(idx uint32) *AstNode {
	if idx != 0 && atomic.LoadUint32(&s.handlesReady) != 0 {
		return &s.handles[idx-1]
	}
	return s.initHandlesSlow(idx)
}

func (s *astNodeSlab) childAt(idx uint32, i int) *AstNode {
	rec := &s.recs[idx]
	cIdx := s.childIdxs[rec.childOff : rec.childOff+rec.childCount][i]
	if atomic.LoadUint32(&s.handlesReady) != 0 {
		return &s.handles[cIdx-1]
	}
	return s.initHandlesSlow(cIdx)
}

// SlabNodeCount returns the total number of AstNode records in m's backing slab,
// or 0 if m is nil or not backed by an unmutated astNodeSlab.
func (m *AstNode) SlabNodeCount() int {
	if m == nil || m.d == nil || m.d.slab == nil {
		return 0
	}
	return len(m.d.slab.recs)
}

// SlabCache returns the external wrapper cache stored on m's backing slab.
func (m *AstNode) SlabCache() any {
	if m == nil || m.d == nil || m.d.slab == nil {
		return nil
	}
	s := m.d.slab
	if atomic.LoadUint32(&s.extCacheReady) != 0 {
		return s.extCache
	}
	s.handlesMu.Lock()
	c := s.extCache
	s.handlesMu.Unlock()
	return c
}

// SetSlabCache stores an external wrapper cache on m's backing slab.
func (m *AstNode) SetSlabCache(c any) {
	if m == nil || m.d == nil || m.d.slab == nil {
		return
	}
	s := m.d.slab
	s.handlesMu.Lock()
	s.extCache = c
	if c != nil {
		atomic.StoreUint32(&s.extCacheReady, 1)
	} else {
		atomic.StoreUint32(&s.extCacheReady, 0)
	}
	s.handlesMu.Unlock()
}

// SlabChildFast resolves child i on the lock-free fast path when m's slab
// external wrapper cache has already been initialized for the current parse.
func (m *AstNode) SlabChildFast(i int) (uint32, any) {
	if d := m.d; d != nil {
		if s := d.slab; s != nil && atomic.LoadUint32(&s.extCacheReady) != 0 {
			rec := &s.recs[m.idx]
			cIdx := s.childIdxs[rec.childOff : rec.childOff+rec.childCount][i]
			return cIdx, s.extCache
		}
	}
	return 0, nil
}

// SlabChildInit initializes m's slab handles and external wrapper cache on first
// GetChild access and returns the child's 0-based slab index plus the cache value.
func (m *AstNode) SlabChildInit(i int, initFn func(prev any, root *AstNode, handles []AstNode) any) (uint32, any, bool) {
	if m == nil || m.d == nil {
		return 0, nil, false
	}
	s := m.d.slab
	if s == nil {
		return 0, nil, false
	}
	rec := &s.recs[m.idx]
	cIdx := s.childIdxs[rec.childOff : rec.childOff+rec.childCount][i]
	if atomic.LoadUint32(&s.extCacheReady) != 0 {
		return cIdx, s.extCache, true
	}
	s.handlesMu.Lock()
	s.ensureHandlesLocked()
	if atomic.LoadUint32(&s.extCacheReady) == 0 {
		s.extCache = initFn(s.extCache, &s.root, s.handles)
		atomic.StoreUint32(&s.extCacheReady, 1)
	}
	c := s.extCache
	s.handlesMu.Unlock()
	return cIdx, c, true
}

// ReleaseSlab clears m's backing astNodeSlab and returns it to astNodeSlabPool
// when m is backed by an unmutated slab within pooled capacity bounds.
func (m *AstNode) ReleaseSlab() {
	if m == nil || m.d == nil {
		return
	}
	s := m.d.slab
	if s == nil || s.mutated || s.inPool {
		return
	}
	if m != &s.root {
		*m = AstNode{}
	}
	if cap(s.recs) > astSlabMaxPooledRecs || cap(s.childIdxs) > astSlabMaxPooledChildren {
		return
	}
	s.recs = s.recs[:0]
	s.childIdxs = s.childIdxs[:0]
	s.vals64 = s.vals64[:0]
	s.tsvs = s.tsvs[:0]
	s.aux = s.aux[:0]
	for i := range s.jsdocs {
		var savedKind []int32
		if cap(s.jsdocs[i].kind_) <= 16 {
			savedKind = s.jsdocs[i].kind_[:0]
		}
		s.jsdocs[i] = OptimizationJsdoc{kind_: savedKind}
	}
	s.jsdocs = s.jsdocs[:0]
	if cap(s.scratch) <= 4096 {
		s.scratch = s.scratch[:0]
	} else {
		s.scratch = nil
	}
	if cap(s.cachedSizes) <= astSlabMaxPooledRecs {
		s.cachedSizes = s.cachedSizes[:0]
	} else {
		s.cachedSizes = nil
	}
	if len(s.auxOverflow) > 0 {
		clear(s.auxOverflow)
	}
	if cap(s.handles) <= astSlabMaxPooledRecs {
		s.handles = s.handles[:0]
	} else {
		s.handles = nil
		s.handlesInit = 0
		s.extCache = nil
	}
	atomic.StoreUint32(&s.handlesReady, 0)
	atomic.StoreUint32(&s.extCacheReady, 0)
	s.sizesOnce = sync.Once{}
	s.shared = astNodeData{slab: s, flags_: astNodeFlagBuilt}
	s.root = AstNode{d: &s.shared, idx: 0}
	s.inPool = true
	astNodeSlabPool.Put(s)
}

func (s *astNodeSlab) materializeData(idx uint32) *astNodeData {
	rec := &s.recs[idx]
	nd := &astNodeData{
		booleanProperties_: rec.boolProps(),
		relativeLine_:      rec.relativeLine,
		relativeColumn_:    rec.relativeColumn,
		type_:              rec.type_,
		kind_:              rec.kind,
		flags_:             rec.flags &^ astNodeFlagBuilt,
	}
	switch rec.valueCase() {
	case AstNode_ValueCase_STRING_VALUE_POINTER:
		nd.valueBits_ = uint64(rec.value32)
	case AstNode_ValueCase_DOUBLE_VALUE:
		nd.valueBits_ = s.vals64[rec.value32]
	case AstNode_ValueCase_TEMPLATE_STRING_VALUE:
		if rec.flags&astNodeFlagHasTemplate != 0 {
			tv := s.tsvs[rec.value32]
			nd.valueBits_ = packTemplateStringBits(tv.rawStringPointer_, tv.cookedStringPointer_)
		}
	}
	if rec.auxIdx != 0 {
		a := s.getAux(idx, rec.auxIdx)
		nd.originalNamePointer_ = a.originalNamePointer
		nd.sourceFile_ = a.sourceFile
		if a.jsdocIdx != 0 {
			nd.jsdoc_ = &s.jsdocs[a.jsdocIdx-1]
		}
	}
	if rec.childCount > 0 {
		sub := s.childIdxs[rec.childOff : rec.childOff+rec.childCount]
		children := make([]*AstNode, len(sub))
		for i, cIdx := range sub {
			children[i] = s.getHandle(cIdx)
		}
		nd.child_ = children
	}
	s.mutated = true
	return nd
}

type AstNode_Builder = AstNode

func AstNode_NewBuilder() *AstNode_Builder {
	if batch := astBuilderBatchPtr.Load(); batch != nil {
		if i := atomic.AddUint32(&batch.idx, 1) - 1; i < astBuilderChunkSize {
			b := batch.ptrs[i]
			atomic.AddUint32(&batch.done, 1)
			return &b.node
		}
	}
	return newAstNodeBuilderSlow()
}

func AstNode_ParseFrom(data []byte) (*AstNode, error) {
	rd := protoReader{b: data, limit: math.MaxInt32}
	return parseAstNodeFromReader(&rd)
}

func AstNode_ParseFromCodedInput(in *CodedInputStream, _ *ExtensionRegistry) (n *AstNode, err error) {
	if err := in.materializeBytes(); err != nil {
		return nil, err
	}
	return parseAstNodeFromReader(&in.rd)
}

func parseAstNodeFromReader(rd *protoReader) (n *AstNode, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if ex, ok := rec.(error); ok {
				err = ex
				return
			}
			panic(rec)
		}
	}()
	slab := newAstNodeSlab(len(rd.b) - rd.pos)
	rootIdx := slab.allocRec()
	slab.mergeRec(rootIdx, rd)
	if cap(slab.scratch) <= 4096 {
		slab.scratch = slab.scratch[:0]
	} else {
		slab.scratch = nil
	}
	return &slab.root, nil
}

func (m *AstNode) mutData() *astNodeData {
	if d := m.d; d != nil && d.flags_&astNodeFlagBuilt == 0 {
		return d
	}
	return m.mutDataSlow()
}

//go:noinline
func (m *AstNode) mutDataSlow() *astNodeData {
	d := m.d
	if d == nil {
		nd := allocStandaloneAstData(m)
		m.d = nd
		m.idx = 0
		return nd
	}
	if s := d.slab; s != nil {
		nd := s.materializeData(m.idx)
		nd.selfNode = m
		m.d = nd
		m.idx = 0
		return nd
	}
	// Copy-on-write when mutating a builder/message copy after Build().
	pinAstData(d)
	cp := allocStandaloneAstData(m)
	cp.jsdoc_ = d.jsdoc_
	cp.valueBits_ = d.valueBits_
	cp.booleanProperties_ = d.booleanProperties_
	cp.relativeLine_ = d.relativeLine_
	cp.relativeColumn_ = d.relativeColumn_
	cp.originalNamePointer_ = d.originalNamePointer_
	cp.type_ = d.type_
	cp.sourceFile_ = d.sourceFile_
	cp.kind_ = d.kind_
	cp.flags_ = d.flags_ &^ astNodeFlagBuilt
	if d.child_ != nil {
		if len(d.child_) > 0 {
			ch := make([]*AstNode, len(d.child_))
			for i, c := range d.child_ {
				ch[i] = c
				if c != nil && c.d != nil && c.d.slab == nil {
					pinAstData(c.d)
				}
			}
			cp.child_ = ch
		}
	} else if d.childCount_ > 0 {
		ch := make([]*AstNode, d.childCount_)
		i := 0
		for c := d.firstChild; c != nil; c = c.nextSibling {
			pinAstData(c)
			ch[i] = c.selfNode
			i++
		}
		cp.child_ = ch
	}
	m.d = cp
	m.idx = 0
	return cp
}

func (m *AstNode) GetKind() NodeKind {
	d := m.d
	if d == nil {
		return NodeKind_NODE_KIND_UNSPECIFIED
	}
	if s := d.slab; s != nil {
		return NodeKind(s.recs[m.idx].kind)
	}
	return NodeKind(d.kind_)
}

func (m *AstNode) GetChildList() []*AstNode {
	d := m.d
	if d == nil {
		return nil
	}
	if s := d.slab; s != nil {
		rec := &s.recs[m.idx]
		if rec.childCount == 0 {
			return nil
		}
		sub := s.childIdxs[rec.childOff : rec.childOff+rec.childCount]
		out := make([]*AstNode, len(sub))
		for i, cIdx := range sub {
			out[i] = s.getHandle(cIdx)
		}
		return out
	}
	pinAstData(d)
	if d.child_ != nil {
		for _, c := range d.child_ {
			if c != nil && c.d != nil && c.d.slab == nil {
				pinAstData(c.d)
			}
		}
		return d.child_
	}
	if d.childCount_ == 0 {
		return nil
	}
	out := make([]*AstNode, d.childCount_)
	i := 0
	for c := d.firstChild; c != nil; c = c.nextSibling {
		pinAstData(c)
		out[i] = c.selfNode
		i++
	}
	if d.flags_&astNodeFlagBuilt == 0 {
		d.child_ = out
		d.firstChild = nil
		d.lastChild = nil
	}
	return out
}

func (m *AstNode) GetChildCount() int {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		return int(s.recs[m.idx].childCount)
	}
	if d.child_ != nil {
		return len(d.child_)
	}
	return int(d.childCount_)
}

//go:noinline
func (d *astNodeData) builderChildAt(i int) *AstNode {
	pinAstData(d)
	if d.child_ != nil {
		c := d.child_[i]
		if c != nil && c.d != nil && c.d.slab == nil {
			pinAstData(c.d)
		}
		return c
	}
	if uint(i) >= uint(d.childCount_) {
		return d.child_[i]
	}
	if d.flags_&astNodeFlagBuilt == 0 {
		d.ensureChildSlice()
		c := d.child_[i]
		if c != nil && c.d != nil && c.d.slab == nil {
			pinAstData(c.d)
		}
		return c
	}
	c := d.firstChild
	for j := 0; j < i; j++ {
		c = c.nextSibling
	}
	pinAstData(c)
	return c.selfNode
}

func (m *AstNode) GetChild(i int) *AstNode {
	if s := m.d.slab; s != nil {
		return s.childAt(m.idx, i)
	}
	return m.d.builderChildAt(i)
}

func (m *AstNode) GetValueCase() AstNode_ValueCase {
	d := m.d
	if d == nil {
		return AstNode_ValueCase_VALUE_NOT_SET
	}
	if s := d.slab; s != nil {
		return s.recs[m.idx].valueCase()
	}
	return d.valueCase()
}

func (m *AstNode) GetStringValuePointer() int32 {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		rec := &s.recs[m.idx]
		if rec.valueCase() == AstNode_ValueCase_STRING_VALUE_POINTER {
			return int32(rec.value32)
		}
		return 0
	}
	if d.valueCase() == AstNode_ValueCase_STRING_VALUE_POINTER {
		return int32(d.valueBits_)
	}
	return 0
}

func (m *AstNode) GetDoubleValue() float64 {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		rec := &s.recs[m.idx]
		if rec.valueCase() == AstNode_ValueCase_DOUBLE_VALUE {
			return math.Float64frombits(s.vals64[rec.value32])
		}
		return 0
	}
	if d.valueCase() == AstNode_ValueCase_DOUBLE_VALUE {
		return math.Float64frombits(d.valueBits_)
	}
	return 0
}

//go:noinline
func (d *astNodeData) getTemplateStringSlow() *TemplateStringValue {
	if d.valueCase() == AstNode_ValueCase_TEMPLATE_STRING_VALUE && d.flags_&astNodeFlagHasTemplate != 0 {
		raw, cooked := unpackTemplateStringBits(d.valueBits_)
		if raw == 0 && cooked == 0 {
			return templateStringValueDefault
		}
		return &TemplateStringValue{
			rawStringPointer_:    raw,
			cookedStringPointer_: cooked,
		}
	}
	return templateStringValueDefault
}

func (m *AstNode) GetTemplateStringValue() *TemplateStringValue {
	d := m.d
	if d == nil {
		return templateStringValueDefault
	}
	if s := d.slab; s != nil {
		rec := &s.recs[m.idx]
		if rec.valueCase() == AstNode_ValueCase_TEMPLATE_STRING_VALUE && rec.flags&astNodeFlagHasTemplate != 0 {
			return &s.tsvs[rec.value32]
		}
		return templateStringValueDefault
	}
	return d.getTemplateStringSlow()
}

func (m *AstNode) GetRelativeLine() int32 {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		return s.recs[m.idx].relativeLine
	}
	return d.relativeLine_
}

func (m *AstNode) GetRelativeColumn() int32 {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		return s.recs[m.idx].relativeColumn
	}
	return d.relativeColumn_
}

func (m *AstNode) HasJsdoc() bool {
	d := m.d
	if d == nil {
		return false
	}
	if s := d.slab; s != nil {
		rec := &s.recs[m.idx]
		if rec.auxIdx == 0 {
			return false
		}
		return s.getAux(m.idx, rec.auxIdx).jsdocIdx != 0
	}
	return d.jsdoc_ != nil
}

//go:noinline
func (s *astNodeSlab) getJsdocSlow(idx uint32, auxIdx uint16) *OptimizationJsdoc {
	if jIdx := s.getAux(idx, auxIdx).jsdocIdx; jIdx != 0 {
		return &s.jsdocs[jIdx-1]
	}
	return OptimizationJsdoc_GetDefaultInstance()
}

func (m *AstNode) GetJsdoc() *OptimizationJsdoc {
	d := m.d
	if d != nil {
		if s := d.slab; s != nil {
			if auxIdx := s.recs[m.idx].auxIdx; auxIdx != 0 {
				return s.getJsdocSlow(m.idx, auxIdx)
			}
			return OptimizationJsdoc_GetDefaultInstance()
		}
		if d.jsdoc_ != nil {
			return d.jsdoc_
		}
	}
	return OptimizationJsdoc_GetDefaultInstance()
}

func (m *AstNode) GetOriginalNamePointer() int32 {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		rec := &s.recs[m.idx]
		if rec.auxIdx == 0 {
			return 0
		}
		return s.getAux(m.idx, rec.auxIdx).originalNamePointer
	}
	return d.originalNamePointer_
}

func (m *AstNode) GetBooleanProperties() int64 {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		return s.recs[m.idx].boolProps()
	}
	return d.booleanProperties_
}

func (m *AstNode) HasType() bool {
	d := m.d
	if d == nil {
		return false
	}
	if s := d.slab; s != nil {
		return s.recs[m.idx].flags&astNodeFlagHasType != 0
	}
	return d.flags_&astNodeFlagHasType != 0
}

func (m *AstNode) GetType() int32 {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		return s.recs[m.idx].type_
	}
	return d.type_
}

func (m *AstNode) GetSourceFile() int32 {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		rec := &s.recs[m.idx]
		if rec.auxIdx == 0 {
			return 0
		}
		return s.getAux(m.idx, rec.auxIdx).sourceFile
	}
	return d.sourceFile_
}

func (m *AstNode) SetKind(v NodeKind) *AstNode {
	d := m.mutData()
	if uint32(v) > uint32(NodeKind_UNRECOGNIZED) {
		d.kind_ = uint8(NodeKind_UNRECOGNIZED)
	} else {
		d.kind_ = uint8(v)
	}
	d.cachedSize = 0
	return m
}

func (m *AstNode) AddChild(v *AstNode) *AstNode {
	if d := m.d; d != nil && d.flags_&astNodeFlagBuilt == 0 && d.child_ == nil && v != nil {
		if vd := v.d; vd != nil && vd.slab == nil && vd != d && vd.selfNode == v && atomic.CompareAndSwapUint32(&vd.linkState_, 0, astLinkFlagLinked) {
			vd.nextSibling = nil
			if d.lastChild == nil {
				d.firstChild = vd
			} else {
				d.lastChild.nextSibling = vd
			}
			d.lastChild = vd
			d.childCount_++
			d.cachedSize = 0
			return m
		}
	}
	return m.addChildSlow(v)
}

//go:noinline
func (m *AstNode) addChildSlow(v *AstNode) *AstNode {
	d := m.mutData()
	if d.child_ == nil && v != nil {
		if vd := v.d; vd != nil && vd.slab == nil && vd != d && vd.selfNode == v && atomic.CompareAndSwapUint32(&vd.linkState_, 0, astLinkFlagLinked) {
			vd.nextSibling = nil
			if d.lastChild == nil {
				d.firstChild = vd
			} else {
				d.lastChild.nextSibling = vd
			}
			d.lastChild = vd
			d.childCount_++
			d.cachedSize = 0
			return m
		}
	}
	if v != nil {
		if vd := v.d; vd != nil && vd.slab == nil {
			if vd == d || d.bundle == nil || atomic.LoadUint32(&d.linkState_)&astLinkFlagPinned != 0 || !atomic.CompareAndSwapUint32(&vd.linkState_, 0, astLinkFlagLinked) {
				pinAstData(vd)
			}
		}
	}
	d.ensureChildSlice()
	d.child_ = append(d.child_, v)
	d.cachedSize = 0
	return m
}

func (m *AstNode) ResetForPool() {
	d := m.d
	if d != nil && d.slab == nil {
		d.resetBuilderFields()
		if d.selfNode != m && (d.bundle == nil || m == &d.bundle.node) {
			d.selfNode = m
		}
		m.idx = 0
		return
	}
	*m = AstNode{}
}

func (m *AstNode) SetStringValuePointer(v int32) *AstNode {
	d := m.mutData()
	d.valueBits_ = uint64(uint32(v))
	d.setValueCase(AstNode_ValueCase_STRING_VALUE_POINTER, false)
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetDoubleValue(v float64) *AstNode {
	d := m.mutData()
	d.valueBits_ = math.Float64bits(v)
	d.setValueCase(AstNode_ValueCase_DOUBLE_VALUE, false)
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetTemplateStringValue(v *TemplateStringValue) *AstNode {
	d := m.mutData()
	if v != nil {
		d.valueBits_ = packTemplateStringBits(v.rawStringPointer_, v.cookedStringPointer_)
		d.setValueCase(AstNode_ValueCase_TEMPLATE_STRING_VALUE, true)
	} else {
		d.valueBits_ = 0
		d.setValueCase(AstNode_ValueCase_TEMPLATE_STRING_VALUE, false)
	}
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetRelativeLine(v int32) *AstNode {
	d := m.mutData()
	d.relativeLine_ = v
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetRelativeColumn(v int32) *AstNode {
	d := m.mutData()
	d.relativeColumn_ = v
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetJsdoc(v *OptimizationJsdoc) *AstNode {
	d := m.mutData()
	d.jsdoc_ = v
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetOriginalNamePointer(v int32) *AstNode {
	d := m.mutData()
	d.originalNamePointer_ = v
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetBooleanProperties(v int64) *AstNode {
	d := m.mutData()
	d.booleanProperties_ = v
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetType(v int32) *AstNode {
	d := m.mutData()
	d.type_ = v
	d.flags_ |= astNodeFlagHasType
	d.cachedSize = 0
	return m
}

func (m *AstNode) SetSourceFile(v int32) *AstNode {
	d := m.mutData()
	d.sourceFile_ = v
	d.cachedSize = 0
	return m
}

func (m *AstNode) Build() *AstNode {
	if d := m.d; d != nil && d.slab == nil {
		d.cachedSize = 0
		d.flags_ |= astNodeFlagBuilt
	}
	return m
}

func (s *astNodeSlab) recScalarSerializedSize(idx uint32, rec *astNodeRec) int {
	size := 0
	if k := rec.kind; k != 0 {
		if k < uint8(NodeKind_UNRECOGNIZED) {
			size += 3
		} else {
			size += 11
		}
	}
	vc := rec.valueCase()
	switch vc {
	case AstNode_ValueCase_STRING_VALUE_POINTER:
		size += 1 + protoVarintLen(uint64(rec.value32))
	case AstNode_ValueCase_DOUBLE_VALUE:
		size += 1 + 8
	}
	if rec.relativeLine != 0 {
		size += 1 + protoVarintLen(uint64(uint32((rec.relativeLine<<1)^(rec.relativeLine>>31))))
	}
	if rec.relativeColumn != 0 {
		size += 1 + protoVarintLen(uint64(uint32((rec.relativeColumn<<1)^(rec.relativeColumn>>31))))
	}
	var a astNodeAux
	if rec.auxIdx != 0 {
		a = s.getAux(idx, rec.auxIdx)
		if a.jsdocIdx != 0 {
			js := s.jsdocs[a.jsdocIdx-1].serializedSize()
			size += 1 + protoVarintLen(uint64(js)) + js
		}
	}
	if vc == AstNode_ValueCase_TEMPLATE_STRING_VALUE && rec.flags&astNodeFlagHasTemplate != 0 {
		tv := s.tsvs[rec.value32]
		ts := templateStringSerializedSize(tv.rawStringPointer_, tv.cookedStringPointer_)
		size += 1 + protoVarintLen(uint64(ts)) + ts
	}
	if a.originalNamePointer != 0 {
		size += 1 + protoVarintLen(uint64(uint32(a.originalNamePointer)))
	}
	if bp := uint64(rec.boolProps()); bp != 0 {
		size += 1 + protoVarintLen(bp)
	}
	if rec.flags&astNodeFlagHasType != 0 {
		size += 1 + protoVarintLen(uint64(uint32(rec.type_)))
	}
	if a.sourceFile != 0 {
		size += 2 + protoVarintLen(uint64(uint32(a.sourceFile)))
	}
	return size
}

func (s *astNodeSlab) initCachedSizes() {
	n := len(s.recs)
	var sizes []uint32
	if cap(s.cachedSizes) >= n {
		sizes = s.cachedSizes[:n]
	} else {
		sizes = make([]uint32, n)
	}
	for i := n - 1; i >= 0; i-- {
		rec := &s.recs[i]
		sz := s.recScalarSerializedSize(uint32(i), rec)
		if rec.childCount > 0 {
			sub := s.childIdxs[rec.childOff : rec.childOff+rec.childCount]
			for _, cIdx := range sub {
				cs := int(sizes[cIdx])
				sz += 1 + protoVarintLen(uint64(cs)) + cs
			}
		}
		sizes[i] = uint32(sz)
	}
	s.cachedSizes = sizes
}

func (s *astNodeSlab) recSerializedSize(idx uint32) int {
	if s.mutated {
		if idx == 0 && s.root.d != &s.shared {
			return s.root.serializedSize()
		}
		if idx > 0 && atomic.LoadUint32(&s.handlesReady) != 0 {
			if h := &s.handles[idx-1]; h.d != &s.shared {
				return h.serializedSize()
			}
		}
		rec := &s.recs[idx]
		sz := s.recScalarSerializedSize(idx, rec)
		if rec.childCount > 0 {
			sub := s.childIdxs[rec.childOff : rec.childOff+rec.childCount]
			for _, cIdx := range sub {
				cs := s.recSerializedSize(cIdx)
				sz += 1 + protoVarintLen(uint64(cs)) + cs
			}
		}
		return sz
	}
	s.sizesOnce.Do(s.initCachedSizes)
	return int(s.cachedSizes[idx])
}

func (d *astNodeData) dataSerializedSize() int {
	if cs := atomic.LoadInt32(&d.cachedSize); cs > 0 {
		return int(cs - 1)
	}
	size := 0
	if k := d.kind_; k != 0 {
		if k < uint8(NodeKind_UNRECOGNIZED) {
			size += 3
		} else {
			size += 11
		}
	}
	if d.child_ == nil {
		for c := d.firstChild; c != nil; c = c.nextSibling {
			s := c.dataSerializedSize()
			size += 1 + protoVarintLen(uint64(s)) + s
		}
	} else {
		for _, c := range d.child_ {
			s := c.serializedSize()
			size += 1 + protoVarintLen(uint64(s)) + s
		}
	}
	vc := d.valueCase()
	switch vc {
	case AstNode_ValueCase_STRING_VALUE_POINTER:
		size += 1 + protoVarintLen(uint64(uint32(d.valueBits_)))
	case AstNode_ValueCase_DOUBLE_VALUE:
		size += 1 + 8
	}
	if d.relativeLine_ != 0 {
		size += 1 + protoVarintLen(uint64(uint32((d.relativeLine_<<1)^(d.relativeLine_>>31))))
	}
	if d.relativeColumn_ != 0 {
		size += 1 + protoVarintLen(uint64(uint32((d.relativeColumn_<<1)^(d.relativeColumn_>>31))))
	}
	if d.jsdoc_ != nil {
		s := d.jsdoc_.serializedSize()
		size += 1 + protoVarintLen(uint64(s)) + s
	}
	if vc == AstNode_ValueCase_TEMPLATE_STRING_VALUE && d.flags_&astNodeFlagHasTemplate != 0 {
		raw, cooked := unpackTemplateStringBits(d.valueBits_)
		s := templateStringSerializedSize(raw, cooked)
		size += 1 + protoVarintLen(uint64(s)) + s
	}
	if d.originalNamePointer_ != 0 {
		size += 1 + protoVarintLen(uint64(uint32(d.originalNamePointer_)))
	}
	if d.booleanProperties_ != 0 {
		size += 1 + protoVarintLen(uint64(d.booleanProperties_))
	}
	if d.flags_&astNodeFlagHasType != 0 {
		size += 1 + protoVarintLen(uint64(uint32(d.type_)))
	}
	if d.sourceFile_ != 0 {
		size += 2 + protoVarintLen(uint64(uint32(d.sourceFile_)))
	}
	if size < math.MaxInt32 {
		atomic.StoreInt32(&d.cachedSize, int32(size+1))
	}
	return size
}

func (m *AstNode) serializedSize() int {
	d := m.d
	if d == nil {
		return 0
	}
	if s := d.slab; s != nil {
		return s.recSerializedSize(m.idx)
	}
	return d.dataSerializedSize()
}

func putTag1Varint(dst []byte, pos int, tag byte, v uint64) int {
	if v < 0x80 {
		_ = dst[pos+1]
		dst[pos] = tag
		dst[pos+1] = byte(v)
		return pos + 2
	}
	dst[pos] = tag
	for pos++; v >= 0x80; pos++ {
		dst[pos] = byte(v) | 0x80
		v >>= 7
	}
	dst[pos] = byte(v)
	return pos + 1
}

func (s *astNodeSlab) putRecTo(idx uint32, dst []byte, pos int) int {
	if s.mutated {
		if idx == 0 && s.root.d != &s.shared {
			return s.root.putTo(dst, pos)
		}
		if idx > 0 && atomic.LoadUint32(&s.handlesReady) != 0 {
			if h := &s.handles[idx-1]; h.d != &s.shared {
				return h.putTo(dst, pos)
			}
		}
	}
	rec := &s.recs[idx]
	if k := rec.kind; k != 0 {
		if k < uint8(NodeKind_UNRECOGNIZED) {
			kn := uint32(nodeKind_numbers[k])
			_ = dst[pos+2]
			dst[pos] = 0x08
			dst[pos+1] = byte(kn) | 0x80
			dst[pos+2] = byte(kn >> 7)
			pos += 3
		} else {
			pos = putTag1Varint(dst, pos, 0x08, math.MaxUint64)
		}
	}
	if rec.childCount > 0 {
		sub := s.childIdxs[rec.childOff : rec.childOff+rec.childCount]
		if !s.mutated {
			sizes := s.cachedSizes
			for _, cIdx := range sub {
				pos = putTag1Varint(dst, pos, 0x12, uint64(sizes[cIdx]))
				pos = s.putRecTo(cIdx, dst, pos)
			}
		} else {
			for _, cIdx := range sub {
				pos = putTag1Varint(dst, pos, 0x12, uint64(s.recSerializedSize(cIdx)))
				pos = s.putRecTo(cIdx, dst, pos)
			}
		}
	}
	vc := rec.valueCase()
	switch vc {
	case AstNode_ValueCase_STRING_VALUE_POINTER:
		pos = putTag1Varint(dst, pos, 0x18, uint64(rec.value32))
	case AstNode_ValueCase_DOUBLE_VALUE:
		dst[pos] = 0x21
		binary.LittleEndian.PutUint64(dst[pos+1:], s.vals64[rec.value32])
		pos += 9
	}
	if rec.relativeLine != 0 {
		pos = putTag1Varint(dst, pos, 0x28, uint64(uint32((rec.relativeLine<<1)^(rec.relativeLine>>31))))
	}
	if rec.relativeColumn != 0 {
		pos = putTag1Varint(dst, pos, 0x30, uint64(uint32((rec.relativeColumn<<1)^(rec.relativeColumn>>31))))
	}
	var a astNodeAux
	if rec.auxIdx != 0 {
		a = s.getAux(idx, rec.auxIdx)
		if a.jsdocIdx != 0 {
			jsdoc := &s.jsdocs[a.jsdocIdx-1]
			js := jsdoc.serializedSize()
			pos = putTag1Varint(dst, pos, 0x3A, uint64(js))
			pos = len(jsdoc.appendTo(dst[:pos : pos+js]))
		}
	}
	if vc == AstNode_ValueCase_TEMPLATE_STRING_VALUE && rec.flags&astNodeFlagHasTemplate != 0 {
		tv := s.tsvs[rec.value32]
		raw, cooked := tv.rawStringPointer_, tv.cookedStringPointer_
		ts := templateStringSerializedSize(raw, cooked)
		pos = putTag1Varint(dst, pos, 0x42, uint64(ts))
		if raw != 0 {
			pos = putTag1Varint(dst, pos, 0x08, uint64(uint32(raw)))
		}
		if cooked != 0 {
			pos = putTag1Varint(dst, pos, 0x10, uint64(int64(cooked)))
		}
	}
	if a.originalNamePointer != 0 {
		pos = putTag1Varint(dst, pos, 0x48, uint64(uint32(a.originalNamePointer)))
	}
	if bp := uint64(rec.boolProps()); bp != 0 {
		pos = putTag1Varint(dst, pos, 0x70, bp)
	}
	if rec.flags&astNodeFlagHasType != 0 {
		pos = putTag1Varint(dst, pos, 0x78, uint64(uint32(rec.type_)))
	}
	if a.sourceFile != 0 {
		dst[pos] = 0x80
		dst[pos+1] = 0x01
		pos += 2 + protoPutVarint(dst[pos+2:], uint64(uint32(a.sourceFile)))
	}
	return pos
}

func (d *astNodeData) putDataTo(dst []byte, pos int) int {
	if k := d.kind_; k != 0 {
		if k < uint8(NodeKind_UNRECOGNIZED) {
			kn := uint32(nodeKind_numbers[k])
			_ = dst[pos+2]
			dst[pos] = 0x08
			dst[pos+1] = byte(kn) | 0x80
			dst[pos+2] = byte(kn >> 7)
			pos += 3
		} else {
			pos = putTag1Varint(dst, pos, 0x08, math.MaxUint64)
		}
	}
	if d.child_ == nil {
		for c := d.firstChild; c != nil; c = c.nextSibling {
			cs := int(atomic.LoadInt32(&c.cachedSize)) - 1
			if cs < 0 {
				cs = c.dataSerializedSize()
			}
			pos = putTag1Varint(dst, pos, 0x12, uint64(cs))
			pos = c.putDataTo(dst, pos)
		}
	} else {
		for _, c := range d.child_ {
			pos = putTag1Varint(dst, pos, 0x12, uint64(c.serializedSize()))
			pos = c.putTo(dst, pos)
		}
	}
	vc := d.valueCase()
	switch vc {
	case AstNode_ValueCase_STRING_VALUE_POINTER:
		pos = putTag1Varint(dst, pos, 0x18, uint64(uint32(d.valueBits_)))
	case AstNode_ValueCase_DOUBLE_VALUE:
		dst[pos] = 0x21
		binary.LittleEndian.PutUint64(dst[pos+1:], d.valueBits_)
		pos += 9
	}
	if d.relativeLine_ != 0 {
		pos = putTag1Varint(dst, pos, 0x28, uint64(uint32((d.relativeLine_<<1)^(d.relativeLine_>>31))))
	}
	if d.relativeColumn_ != 0 {
		pos = putTag1Varint(dst, pos, 0x30, uint64(uint32((d.relativeColumn_<<1)^(d.relativeColumn_>>31))))
	}
	if d.jsdoc_ != nil {
		js := d.jsdoc_.serializedSize()
		pos = putTag1Varint(dst, pos, 0x3A, uint64(js))
		pos = len(d.jsdoc_.appendTo(dst[:pos : pos+js]))
	}
	if vc == AstNode_ValueCase_TEMPLATE_STRING_VALUE && d.flags_&astNodeFlagHasTemplate != 0 {
		raw, cooked := unpackTemplateStringBits(d.valueBits_)
		s := templateStringSerializedSize(raw, cooked)
		pos = putTag1Varint(dst, pos, 0x42, uint64(s))
		if raw != 0 {
			pos = putTag1Varint(dst, pos, 0x08, uint64(uint32(raw)))
		}
		if cooked != 0 {
			pos = putTag1Varint(dst, pos, 0x10, uint64(int64(cooked)))
		}
	}
	if d.originalNamePointer_ != 0 {
		pos = putTag1Varint(dst, pos, 0x48, uint64(uint32(d.originalNamePointer_)))
	}
	if d.booleanProperties_ != 0 {
		pos = putTag1Varint(dst, pos, 0x70, uint64(d.booleanProperties_))
	}
	if d.flags_&astNodeFlagHasType != 0 {
		pos = putTag1Varint(dst, pos, 0x78, uint64(uint32(d.type_)))
	}
	if d.sourceFile_ != 0 {
		dst[pos] = 0x80
		dst[pos+1] = 0x01
		pos += 2 + protoPutVarint(dst[pos+2:], uint64(uint32(d.sourceFile_)))
	}
	return pos
}

func (m *AstNode) putTo(dst []byte, pos int) int {
	d := m.d
	if d == nil {
		return pos
	}
	if s := d.slab; s != nil {
		return s.putRecTo(m.idx, dst, pos)
	}
	return d.putDataTo(dst, pos)
}

func (m *AstNode) appendTo(b []byte) []byte {
	sz := m.serializedSize()
	if sz == 0 {
		return b
	}
	off := len(b)
	if cap(b)-off < sz {
		next := make([]byte, off+sz)
		copy(next, b)
		b = next
	} else {
		b = b[:off+sz]
	}
	end := m.putTo(b, off)
	return b[:end]
}

func (m *AstNode) mergeFromReader(r *protoReader) {
	if m.d == nil || (m.d.slab == nil && m.d.isZero()) {
		slab := newAstNodeSlab(len(r.b) - r.pos)
		rootIdx := slab.allocRec()
		slab.mergeRec(rootIdx, r)
		if cap(slab.scratch) <= 4096 {
			slab.scratch = slab.scratch[:0]
		} else {
			slab.scratch = nil
		}
		*m = slab.root
		return
	}
	d := m.mutData()
	slab := newAstNodeSlab(len(r.b) - r.pos)
	d.mergeIntoBuilder(r, slab)
}

func (d *astNodeData) mergeIntoBuilder(r *protoReader, slab *astNodeSlab) {
	d.cachedSize = 0
	var childIdxs []uint32
	for {
		tag := r.readTag()
		switch tag {
		case 0:
			slab.scratch = nil
			if len(childIdxs) > 0 {
				d.ensureChildSlice()
				for _, cIdx := range childIdxs {
					d.child_ = append(d.child_, slab.getHandle(cIdx))
				}
			}
			return
		case 8:
			d.kind_ = uint8(nodeKind_fromNumber(r.readInt32()))
		case 18:
			data := r.readMessageBytes()
			r.checkRecursionLimit()
			cIdx := slab.allocRec()
			origB, origPos := r.b, r.pos
			r.b, r.pos = data, 0
			r.depth++
			slab.mergeRec(cIdx, r)
			r.depth--
			r.b, r.pos = origB, origPos
			childIdxs = append(childIdxs, cIdx)
		case 24:
			d.valueBits_ = uint64(uint32(r.readUint32()))
			d.setValueCase(AstNode_ValueCase_STRING_VALUE_POINTER, false)
		case 33:
			d.valueBits_ = r.readFixed64()
			d.setValueCase(AstNode_ValueCase_DOUBLE_VALUE, false)
		case 40:
			d.relativeLine_ = r.readSint32()
		case 48:
			d.relativeColumn_ = r.readSint32()
		case 58:
			data := r.readMessageBytes()
			r.checkRecursionLimit()
			if d.jsdoc_ == nil {
				d.jsdoc_ = &OptimizationJsdoc{}
			} else {
				cp := *d.jsdoc_
				cp.cachedSize = 0
				d.jsdoc_ = &cp
			}
			origB, origPos := r.b, r.pos
			r.b, r.pos = data, 0
			r.depth++
			d.jsdoc_.mergeFromReader(r)
			r.depth--
			r.b, r.pos = origB, origPos
		case 66:
			data := r.readMessageBytes()
			r.checkRecursionLimit()
			var tv TemplateStringValue
			if d.valueCase() == AstNode_ValueCase_TEMPLATE_STRING_VALUE && d.flags_&astNodeFlagHasTemplate != 0 {
				tv.rawStringPointer_, tv.cookedStringPointer_ = unpackTemplateStringBits(d.valueBits_)
			}
			origB, origPos := r.b, r.pos
			r.b, r.pos = data, 0
			r.depth++
			tv.mergeFromReader(r)
			r.depth--
			r.b, r.pos = origB, origPos
			d.valueBits_ = packTemplateStringBits(tv.rawStringPointer_, tv.cookedStringPointer_)
			d.setValueCase(AstNode_ValueCase_TEMPLATE_STRING_VALUE, true)
		case 72:
			d.originalNamePointer_ = r.readUint32()
		case 112:
			d.booleanProperties_ = r.readInt64()
		case 120:
			d.type_ = r.readUint32()
			d.flags_ |= astNodeFlagHasType
		case 128:
			d.sourceFile_ = r.readUint32()
		default:
			r.skipField(tag)
		}
	}
}

func (s *astNodeSlab) mergeRec(idx uint32, r *protoReader) {
	cur := s.recs[idx]
	var localChildren [16]uint32
	nLocal := 0
	scratchStart := -1
	if cur.childCount > 0 {
		scratchStart = len(s.scratch)
		existing := s.childIdxs[cur.childOff : cur.childOff+cur.childCount]
		s.scratch = append(s.scratch, existing...)
	}
	for {
		tag := r.readTag()
		switch tag {
		case 0:
			if scratchStart >= 0 {
				ids := s.scratch[scratchStart:]
				cur.childOff = s.appendChildIdxs(ids)
				cur.childCount = uint32(len(ids))
				s.scratch = s.scratch[:scratchStart]
			} else if nLocal > 0 {
				cur.childOff = s.appendChildIdxs(localChildren[:nLocal])
				cur.childCount = uint32(nLocal)
			}
			s.recs[idx] = cur
			return
		case 8:
			cur.kind = uint8(nodeKind_fromNumber(r.readInt32()))
		case 18:
			data := r.readMessageBytes()
			r.checkRecursionLimit()
			cIdx := s.allocRec()
			origB, origPos := r.b, r.pos
			r.b, r.pos = data, 0
			r.depth++
			s.mergeRec(cIdx, r)
			r.depth--
			r.b, r.pos = origB, origPos
			if scratchStart >= 0 {
				s.scratch = append(s.scratch, cIdx)
			} else if nLocal < len(localChildren) {
				localChildren[nLocal] = cIdx
				nLocal++
			} else {
				scratchStart = len(s.scratch)
				s.scratch = append(s.scratch, localChildren[:]...)
				s.scratch = append(s.scratch, cIdx)
			}
		case 24:
			cur.value32 = uint32(r.readUint32())
			cur.setValueCase(AstNode_ValueCase_STRING_VALUE_POINTER, false)
		case 33:
			s.setRecDoubleBits(&cur, r.readFixed64())
			cur.setValueCase(AstNode_ValueCase_DOUBLE_VALUE, false)
		case 40:
			cur.relativeLine = r.readSint32()
		case 48:
			cur.relativeColumn = r.readSint32()
		case 58:
			data := r.readMessageBytes()
			r.checkRecursionLimit()
			origB, origPos := r.b, r.pos
			r.b, r.pos = data, 0
			r.depth++
			s.mergeAuxJsdoc(idx, &cur, r)
			r.depth--
			r.b, r.pos = origB, origPos
		case 66:
			data := r.readMessageBytes()
			r.checkRecursionLimit()
			origB, origPos := r.b, r.pos
			r.b, r.pos = data, 0
			r.depth++
			s.mergeRecTemplateString(&cur, r)
			r.depth--
			r.b, r.pos = origB, origPos
			cur.setValueCase(AstNode_ValueCase_TEMPLATE_STRING_VALUE, true)
		case 72:
			s.setAuxOrigName(idx, &cur, r.readUint32())
		case 112:
			bp := uint64(r.readInt64())
			cur.boolPropsLo = uint32(bp)
			cur.boolPropsHi = uint32(bp >> 32)
		case 120:
			cur.type_ = r.readUint32()
			cur.flags |= astNodeFlagHasType
		case 128:
			s.setAuxSourceFile(idx, &cur, r.readUint32())
		default:
			r.skipField(tag)
		}
	}
}

func (m *AstNode) ToByteArray() []byte {
	if m != nil {
		if d := m.d; d != nil && d.slab == nil {
			pinAstData(d)
		}
	}
	return m.appendTo(make([]byte, 0, m.serializedSize()))
}

const (
	wireSlabSize    = 2 * 1024 * 1024
	wireSlabMaxItem = 512 * 1024
)

var (
	wireSlabMu   sync.Mutex
	wireSlabBuf  []byte
	wireSlabTail []byte
)

func allocWireBytes(sz int) []byte {
	if sz <= 0 {
		return nil
	}
	if sz > wireSlabMaxItem {
		return make([]byte, 0, sz)
	}
	wireSlabMu.Lock()
	if tailRem := cap(wireSlabTail) - len(wireSlabTail); tailRem >= sz {
		off := len(wireSlabTail)
		end := off + sz
		wireSlabTail = wireSlabTail[:end]
		sub := wireSlabTail[off:off:end]
		if cap(wireSlabTail)-end < 256 {
			wireSlabTail = nil
		}
		wireSlabMu.Unlock()
		return sub
	}
	rem := cap(wireSlabBuf) - len(wireSlabBuf)
	if rem < sz {
		if rem >= 256 && rem > cap(wireSlabTail)-len(wireSlabTail) {
			wireSlabTail = wireSlabBuf
		}
		wireSlabBuf = make([]byte, 0, wireSlabSize)
	}
	off := len(wireSlabBuf)
	end := off + sz
	wireSlabBuf = wireSlabBuf[:end]
	sub := wireSlabBuf[off:off:end]
	wireSlabMu.Unlock()
	return sub
}

func (m *AstNode) ToByteString() ByteString {
	if m == nil {
		return ByteString_EMPTY
	}
	sz := m.serializedSize()
	var raw []byte
	if sz > 0 {
		raw = allocWireBytes(sz)[:sz:sz]
		end := m.putTo(raw, 0)
		raw = raw[:end:end]
	}
	if d := m.d; d != nil && d.bundle != nil &&
		NodeKind(d.kind_) == NodeKind_SOURCE_FILE &&
		m == &d.bundle.node &&
		atomic.LoadUint32(&d.linkState_) == 0 {
		recycleAstBuilderTree(d.bundle)
	}
	if len(raw) == 0 {
		return ByteString_EMPTY
	}
	return ByteString{data: raw}
}

// ---------------------------------------------------------------------------
// TypedAst, TypedAst.List, LazyAst, StringPoolProto, ExternsSummary Builders
// ---------------------------------------------------------------------------

const TypedAst_List_TYPED_ASTS_FIELD_NUMBER int32 = 1

type TypedAst_Builder = TypedAst

func TypedAst_NewBuilder() *TypedAst_Builder { return &TypedAst{} }

func (m *TypedAst) Clear() *TypedAst {
	*m = TypedAst{}
	return m
}

func (m *TypedAst) Build() *TypedAst {
	cp := *m
	cp.lazySlab = nil
	cp.lazyStep = 0
	cp.cachedSize = 0
	return &cp
}

func (m *TypedAst) AddExternAst(v *LazyAst) *TypedAst {
	m.externAst_ = append(m.externAst_, v)
	m.cachedSize = 0
	return m
}

func (m *TypedAst) AddCodeAst(v *LazyAst) *TypedAst {
	m.codeAst_ = append(m.codeAst_, v)
	m.cachedSize = 0
	return m
}

func (m *TypedAst) SetTypePool(v *TypePool) *TypedAst {
	m.typePool_ = v
	m.cachedSize = 0
	return m
}

func (m *TypedAst) SetStringPool(v *StringPoolProto) *TypedAst {
	m.stringPool_ = v
	m.cachedSize = 0
	return m
}

func (m *TypedAst) SetSourceFilePool(v *SourceFilePool) *TypedAst {
	m.sourceFilePool_ = v
	m.cachedSize = 0
	return m
}

func (m *TypedAst) SetExternsSummary(v *ExternsSummary) *TypedAst {
	m.externsSummary_ = v
	m.cachedSize = 0
	return m
}

func (m *TypedAst) AddAllRuntimeLibraryToInject(v []string) *TypedAst {
	m.runtimeLibraryToInject_ = append(m.runtimeLibraryToInject_, v...)
	m.cachedSize = 0
	return m
}

func (m *SourceFilePool) appendTo(b []byte) []byte {
	for _, x := range m.sourceFile_ {
		b = protoAppendTag(b, 1, protoWireBytes)
		b = protoAppendVarint(b, uint64(x.serializedSize()))
		b = x.appendTo(b)
	}
	return b
}

func (m *TypedAst) appendTo(b []byte) []byte {
	if m.typePool_ != nil {
		b = protoAppendTag(b, 1, protoWireBytes)
		b = protoAppendVarint(b, uint64(m.typePool_.serializedSize()))
		b = m.typePool_.appendTo(b)
	}
	if m.stringPool_ != nil {
		b = protoAppendTag(b, 2, protoWireBytes)
		b = protoAppendVarint(b, uint64(m.stringPool_.serializedSize()))
		b = m.stringPool_.appendTo(b)
	}
	for _, x := range m.externAst_ {
		b = protoAppendTag(b, 3, protoWireBytes)
		b = protoAppendVarint(b, uint64(x.serializedSize()))
		b = x.appendTo(b)
	}
	for _, x := range m.codeAst_ {
		b = protoAppendTag(b, 4, protoWireBytes)
		b = protoAppendVarint(b, uint64(x.serializedSize()))
		b = x.appendTo(b)
	}
	if m.sourceFilePool_ != nil {
		b = protoAppendTag(b, 5, protoWireBytes)
		b = protoAppendVarint(b, uint64(m.sourceFilePool_.serializedSize()))
		b = m.sourceFilePool_.appendTo(b)
	}
	if m.externsSummary_ != nil {
		b = protoAppendTag(b, 6, protoWireBytes)
		b = protoAppendVarint(b, uint64(m.externsSummary_.serializedSize()))
		b = m.externsSummary_.appendTo(b)
	}
	for _, x := range m.runtimeLibraryToInject_ {
		b = protoAppendTag(b, 7, protoWireBytes)
		b = protoAppendString(b, x)
	}
	return b
}

func (m *TypedAst) ToByteArray() []byte {
	return m.appendTo(make([]byte, 0, m.serializedSize()))
}

type TypedAst_List struct {
	typedAsts_ []*TypedAst
}

type TypedAst_List_Builder = TypedAst_List

func TypedAst_List_NewBuilder() *TypedAst_List_Builder { return &TypedAst_List{} }

func (l *TypedAst_List) GetTypedAstsList() []*TypedAst { return l.typedAsts_ }

func (l *TypedAst_List) AddTypedAsts(v *TypedAst) *TypedAst_List {
	l.typedAsts_ = append(l.typedAsts_, v)
	return l
}

func (l *TypedAst_List) Build() *TypedAst_List { return l }

func (l *TypedAst_List) serializedSize() int {
	size := 0
	for _, ast := range l.typedAsts_ {
		s := ast.serializedSize()
		size += 1 + protoVarintLen(uint64(s)) + s
	}
	return size
}

func (l *TypedAst_List) appendTo(b []byte) []byte {
	for _, ast := range l.typedAsts_ {
		b = protoAppendTag(b, 1, protoWireBytes)
		b = protoAppendVarint(b, uint64(ast.serializedSize()))
		b = ast.appendTo(b)
	}
	return b
}

func (l *TypedAst_List) ToByteArray() []byte {
	return l.appendTo(make([]byte, 0, l.serializedSize()))
}

func (l *TypedAst_List) WriteTo(out io.Writer) (int64, error) {
	return streamProtoMessageTo(l, out, false)
}

func (l *TypedAst_List) WriteToStream(out io.Writer) error {
	_, err := l.WriteTo(out)
	return err
}

type LazyAst_Builder = LazyAst

func LazyAst_NewBuilder() *LazyAst_Builder { return &LazyAst{} }

func (m *LazyAst) GetScriptByteString() ByteString {
	return ByteString{data: m.script_}
}
func (m *LazyAst) SetScript(v ByteString) *LazyAst {
	m.script_ = v.bytesImmutable()
	m.cachedSize = 0
	return m
}
func (m *LazyAst) SetScriptBytes(v []byte) *LazyAst {
	m.script_ = v
	m.cachedSize = 0
	return m
}
func (m *LazyAst) SetSourceFile(v int32) *LazyAst {
	m.sourceFile_ = v
	m.cachedSize = 0
	return m
}
func (m *LazyAst) SetSourceMappingUrl(v string) *LazyAst {
	m.sourceMappingUrl_ = v
	m.cachedSize = 0
	return m
}
func (m *LazyAst) Build() *LazyAst {
	m.cachedSize = 0
	return m
}

type StringPoolProto_Builder = StringPoolProto

func StringPoolProto_NewBuilder() *StringPoolProto_Builder { return &StringPoolProto{} }

func (m *StringPoolProto) addWireString(v []byte) {
	entryLen := 1 + protoVarintLen(uint64(len(v))) + len(v)
	nc := len(m.wireChunks_)
	if nc == 0 || cap(m.wireChunks_[nc-1])-len(m.wireChunks_[nc-1]) < entryLen {
		chunkCap := min(4096<<min(nc, 6), 262144)
		if entryLen > chunkCap {
			chunkCap = entryLen
		}
		m.wireChunks_ = append(m.wireChunks_, make([]byte, 0, chunkCap))
		nc++
	}
	b := m.wireChunks_[nc-1]
	b = append(b, 0x0A)
	b = protoAppendVarint(b, uint64(len(v)))
	b = append(b, v...)
	m.wireChunks_[nc-1] = b
	m.wireStringsSize_ += entryLen
	m.wireStringsCnt_++
	m.cachedSize = 0
}

func (m *StringPoolProto) GetStringsByteStringList() []ByteString {
	if len(m.wireChunks_) == 0 {
		return nil
	}
	out := make([]ByteString, 0, m.wireStringsCnt_)
	for _, c := range m.wireChunks_ {
		pos := 0
		for pos < len(c) {
			pos++ // skip field 1 tag (0x0A)
			sz, _ := protoDecodeSliceVarint64(c, &pos)
			end := pos + int(sz)
			out = append(out, ByteString{data: c[pos:end:end]})
			pos = end
		}
	}
	return out
}
func (m *StringPoolProto) SetMaxLength(v int32) *StringPoolProto {
	m.maxLength_ = v
	m.cachedSize = 0
	return m
}
func (m *StringPoolProto) AddStrings(v ByteString) *StringPoolProto {
	m.addWireString(v.data)
	return m
}
func (m *StringPoolProto) AddStringsBytes(v []byte) *StringPoolProto {
	m.addWireString(v)
	return m
}
func (m *StringPoolProto) appendTo(b []byte) []byte {
	for _, c := range m.wireChunks_ {
		b = append(b, c...)
	}
	if m.maxLength_ != 0 {
		b = protoAppendTag(b, 2, protoWireVarint)
		b = protoAppendUint32(b, m.maxLength_)
	}
	return b
}
func (m *StringPoolProto) Build() *StringPoolProto {
	if n := len(m.wireChunks_); n > 0 {
		last := m.wireChunks_[n-1]
		m.wireChunks_[n-1] = last[:len(last):len(last)]
		m.wireChunks_ = m.wireChunks_[:n:n]
	}
	m.cachedSize = 0
	return m
}

type ExternsSummary_Builder = ExternsSummary

func ExternsSummary_NewBuilder() *ExternsSummary_Builder { return &ExternsSummary{} }

func (m *ExternsSummary) AddPropNamePtr(v int32) *ExternsSummary {
	m.propNamePtr_ = append(m.propNamePtr_, v)
	m.cachedSize = 0
	return m
}
func (m *ExternsSummary) Build() *ExternsSummary {
	m.cachedSize = 0
	return m
}

// ---------------------------------------------------------------------------
// SourceFilePool & SourceFileProto Builders
// ---------------------------------------------------------------------------

type SourceFilePool_Builder = SourceFilePool

func SourceFilePool_NewBuilder() *SourceFilePool_Builder { return &SourceFilePool{} }

func (m *SourceFilePool) AddSourceFile(v *SourceFileProto) *SourceFilePool {
	m.sourceFile_ = append(m.sourceFile_, v)
	m.cachedSize = 0
	return m
}
func (m *SourceFilePool) AddAllSourceFile(v []*SourceFileProto) *SourceFilePool {
	m.sourceFile_ = append(m.sourceFile_, v...)
	m.cachedSize = 0
	return m
}
func (m *SourceFilePool) Build() *SourceFilePool {
	m.cachedSize = 0
	return m
}

type SourceFileProto_Builder = SourceFileProto

func SourceFileProto_NewBuilder() *SourceFileProto_Builder { return &SourceFileProto{} }

func (m *SourceFileProto) SetFilename(v string) *SourceFileProto {
	m.filename_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) SetPreloadedContents(v string) *SourceFileProto {
	m.preloadedContents_ = v
	m.loaderCase_ = 2
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) SetFileOnDisk(v *SourceFileProto_FileOnDisk) *SourceFileProto {
	m.fileOnDisk_ = v
	m.loaderCase_ = 3
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) SetZipEntry(v *SourceFileProto_ZipEntryOnDisk) *SourceFileProto {
	m.zipEntry_ = v
	m.loaderCase_ = 4
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) SetStubFile(v bool) *SourceFileProto {
	m.stubFile_ = v
	m.loaderCase_ = 9
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) SetSourceKind(v SourceFileProto_SourceKind) *SourceFileProto {
	m.sourceKind_ = v.GetNumber()
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) SetIsClosureUnawareCode(v bool) *SourceFileProto {
	m.isClosureUnawareCode_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) SetNumBytesPlusOne(v int32) *SourceFileProto {
	m.numBytesPlusOne_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) SetNumLinesPlusOne(v int32) *SourceFileProto {
	m.numLinesPlusOne_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto) Build() *SourceFileProto {
	m.cachedSize = 0
	return m
}

type SourceFileProto_FileOnDisk_Builder = SourceFileProto_FileOnDisk

func SourceFileProto_FileOnDisk_NewBuilder() *SourceFileProto_FileOnDisk_Builder {
	return &SourceFileProto_FileOnDisk{}
}
func (m *SourceFileProto_FileOnDisk) SetActualPath(v string) *SourceFileProto_FileOnDisk {
	m.actualPath_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto_FileOnDisk) SetCharset(v string) *SourceFileProto_FileOnDisk {
	m.charset_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto_FileOnDisk) Build() *SourceFileProto_FileOnDisk {
	m.cachedSize = 0
	return m
}

type SourceFileProto_ZipEntryOnDisk_Builder = SourceFileProto_ZipEntryOnDisk

func SourceFileProto_ZipEntryOnDisk_NewBuilder() *SourceFileProto_ZipEntryOnDisk_Builder {
	return &SourceFileProto_ZipEntryOnDisk{}
}
func (m *SourceFileProto_ZipEntryOnDisk) SetZipPath(v string) *SourceFileProto_ZipEntryOnDisk {
	m.zipPath_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto_ZipEntryOnDisk) SetEntryName(v string) *SourceFileProto_ZipEntryOnDisk {
	m.entryName_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto_ZipEntryOnDisk) SetCharset(v string) *SourceFileProto_ZipEntryOnDisk {
	m.charset_ = v
	m.cachedSize = 0
	return m
}
func (m *SourceFileProto_ZipEntryOnDisk) Build() *SourceFileProto_ZipEntryOnDisk {
	m.cachedSize = 0
	return m
}

// ---------------------------------------------------------------------------
// TypePool, TypeProto, ObjectTypeProto, UnionTypeProto, OptimizationJsdoc Builders
// ---------------------------------------------------------------------------

type TypePool_Builder = TypePool

func TypePool_NewBuilder() *TypePool_Builder { return &TypePool{} }

func (m *TypePool) AddType(v *TypeProto) *TypePool {
	m.type_ = append(m.type_, v)
	m.cachedSize = 0
	return m
}
func (m *TypePool) AddDisambiguationEdges(v *SubtypingEdge) *TypePool {
	m.disambiguationEdges_ = append(m.disambiguationEdges_, v)
	m.cachedSize = 0
	return m
}
func (m *TypePool) AddDisambiguationEdgesBuilder() *SubtypingEdge {
	e := m.ensureArena().allocEdge()
	m.disambiguationEdges_ = append(m.disambiguationEdges_, e)
	m.cachedSize = 0
	return e
}
func (m *TypePool) GetDebugInfoBuilder() *TypePool_DebugInfo {
	if m.debugInfo_ == nil {
		m.debugInfo_ = &TypePool_DebugInfo{}
	}
	m.cachedSize = 0
	return m.debugInfo_
}
func (m *TypePool) Build() *TypePool {
	m.cachedSize = 0
	return m
}

func (m *TypePool_DebugInfo) AddMismatchBuilder() *TypePool_DebugInfo_Mismatch {
	mm := &TypePool_DebugInfo_Mismatch{}
	m.mismatch_ = append(m.mismatch_, mm)
	m.cachedSize = 0
	return mm
}
func (m *TypePool_DebugInfo) Build() *TypePool_DebugInfo {
	m.cachedSize = 0
	return m
}

func (m *TypePool_DebugInfo_Mismatch) SetSourceRef(v string) *TypePool_DebugInfo_Mismatch {
	m.sourceRef_ = v
	m.cachedSize = 0
	return m
}
func (m *TypePool_DebugInfo_Mismatch) AddAllInvolvedColor(v []int32) *TypePool_DebugInfo_Mismatch {
	m.involvedColor_ = append(m.involvedColor_, v...)
	m.cachedSize = 0
	return m
}
func (m *TypePool_DebugInfo_Mismatch) Build() *TypePool_DebugInfo_Mismatch {
	m.cachedSize = 0
	return m
}

type SubtypingEdge_Builder = SubtypingEdge

func SubtypingEdge_NewBuilder() *SubtypingEdge_Builder { return &SubtypingEdge{} }

func (m *SubtypingEdge) SetSubtype(v int32) *SubtypingEdge {
	m.subtype_ = v
	m.cachedSize = 0
	return m
}
func (m *SubtypingEdge) SetSupertype(v int32) *SubtypingEdge {
	m.supertype_ = v
	m.cachedSize = 0
	return m
}
func (m *SubtypingEdge) Build() *SubtypingEdge {
	m.cachedSize = 0
	return m
}

type UnionTypeProto_Builder = UnionTypeProto

func UnionTypeProto_NewBuilder() *UnionTypeProto_Builder { return &UnionTypeProto{} }

func (m *UnionTypeProto) AddUnionMember(v int32) *UnionTypeProto {
	m.unionMember_ = append(m.unionMember_, v)
	m.cachedSize = 0
	return m
}
func (m *UnionTypeProto) AddAllUnionMember(v []int32) *UnionTypeProto {
	m.unionMember_ = append(m.unionMember_, v...)
	m.cachedSize = 0
	return m
}
func (m *UnionTypeProto) Build() *UnionTypeProto {
	m.cachedSize = 0
	return m
}

type ObjectTypeProto_Builder = ObjectTypeProto

func ObjectTypeProto_NewBuilder() *ObjectTypeProto_Builder { return &ObjectTypeProto{} }

func (m *ObjectTypeProto) GetUuidByteString() ByteString { return ByteString{data: m.uuid_} }

func (m *ObjectTypeProto) SetIsInvalidating(v bool) *ObjectTypeProto {
	if v {
		m.flags_ |= objTypeFlagIsInvalidating
	} else {
		m.flags_ &^= objTypeFlagIsInvalidating
	}
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) SetUuid(v ByteString) *ObjectTypeProto {
	m.uuid_ = v.bytesImmutable()
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) SetUuidBytes(v []byte) *ObjectTypeProto {
	m.uuid_ = v
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) SetPropertiesKeepOriginalName(v bool) *ObjectTypeProto {
	if v {
		m.flags_ |= objTypeFlagPropertiesKeepOriginalName
	} else {
		m.flags_ &^= objTypeFlagPropertiesKeepOriginalName
	}
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) AddAllInstanceType(v []int32) *ObjectTypeProto {
	m.instanceType_ = append(m.instanceType_, v...)
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) AddAllPrototype(v []int32) *ObjectTypeProto {
	m.prototype_ = append(m.prototype_, v...)
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) SetMarkedConstructor(v bool) *ObjectTypeProto {
	if v {
		m.flags_ |= objTypeFlagMarkedConstructor
	} else {
		m.flags_ &^= objTypeFlagMarkedConstructor
	}
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) AddAllOwnProperty(v []int32) *ObjectTypeProto {
	m.ownProperty_ = append(m.ownProperty_, v...)
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) SetClosureAssert(v bool) *ObjectTypeProto {
	if v {
		m.flags_ |= objTypeFlagClosureAssert
	} else {
		m.flags_ &^= objTypeFlagClosureAssert
	}
	m.cachedSize = 0
	return m
}
func (m *ObjectTypeProto) Build() *ObjectTypeProto {
	m.cachedSize = 0
	return m
}

type TypeProto_Builder = TypeProto

func TypeProto_NewBuilder() *TypeProto_Builder { return &TypeProto{} }

func (m *TypeProto) SetObject(v *ObjectTypeProto) *TypeProto {
	m.object_ = v
	m.kindCase_ = 2
	m.cachedSize = 0
	return m
}
func (m *TypeProto) GetObjectBuilder() *ObjectTypeProto {
	if m.kindCase_ != 2 || m.object_ == nil {
		m.object_ = &ObjectTypeProto{}
		m.kindCase_ = 2
	}
	m.cachedSize = 0
	return m.object_
}
func (m *TypeProto) SetUnion(v *UnionTypeProto) *TypeProto {
	m.union_ = v
	m.kindCase_ = 4
	m.cachedSize = 0
	return m
}
func (m *TypeProto) GetUnionBuilder() *UnionTypeProto {
	if m.kindCase_ != 4 || m.union_ == nil {
		m.union_ = &UnionTypeProto{}
		m.kindCase_ = 4
	}
	m.cachedSize = 0
	return m.union_
}
func (m *TypeProto) Build() *TypeProto {
	m.cachedSize = 0
	return m
}

type OptimizationJsdoc_Builder = OptimizationJsdoc

func OptimizationJsdoc_NewBuilder() *OptimizationJsdoc_Builder { return &OptimizationJsdoc{} }

func (m *OptimizationJsdoc) GetKindValueList() []int32 {
	if m.kind_ != nil {
		return m.kind_
	}
	if m.kindInlineLen_ == 0 {
		return nil
	}
	out := make([]int32, m.kindInlineLen_)
	for i := range out {
		out[i] = int32(m.kindInline_[i])
	}
	return out
}
func (m *OptimizationJsdoc) GetKindList() []JsdocTag {
	n := m.GetKindCount()
	out := make([]JsdocTag, n)
	for i := 0; i < n; i++ {
		if t, ok := JsdocTag_ForNumber(m.kindAt(i)); ok {
			out[i] = t
		} else {
			out[i] = JsdocTag_UNRECOGNIZED
		}
	}
	return out
}
func (m *OptimizationJsdoc) AddKind(v JsdocTag) *OptimizationJsdoc {
	m.addKindNum(v.GetNumber())
	m.cachedSize = 0
	return m
}
func (m *OptimizationJsdoc) SetLicenseTextPointer(v int32) *OptimizationJsdoc {
	m.licenseTextPointer_ = v
	m.cachedSize = 0
	return m
}
func (m *OptimizationJsdoc) SetDescriptionPointer(v int32) *OptimizationJsdoc {
	m.descriptionPointer_ = v
	m.cachedSize = 0
	return m
}
func (m *OptimizationJsdoc) SetMeaningPointer(v int32) *OptimizationJsdoc {
	m.meaningPointer_ = v
	m.cachedSize = 0
	return m
}
func (m *OptimizationJsdoc) SetAlternateMessageIdPointer(v int32) *OptimizationJsdoc {
	m.alternateMessageIdPointer_ = v
	m.cachedSize = 0
	return m
}
func (m *OptimizationJsdoc) SetClosureUnawareModePointer(v int32) *OptimizationJsdoc {
	m.closureUnawareModePointer_ = v
	m.cachedSize = 0
	return m
}
func (m *OptimizationJsdoc) Build() *OptimizationJsdoc {
	m.cachedSize = 0
	return m
}

// Requirement.toBuilder() helper used by CheckConformance & ConformanceAllowlister.
func (m *Requirement) ToBuilder() *Requirement {
	cp := *m
	cp.whitelist_ = append([]string(nil), m.whitelist_...)
	cp.whitelistRegexp_ = append([]string(nil), m.whitelistRegexp_...)
	cp.allowlist_ = append([]string(nil), m.allowlist_...)
	cp.allowlistRegexp_ = append([]string(nil), m.allowlistRegexp_...)
	cp.onlyApplyTo_ = append([]string(nil), m.onlyApplyTo_...)
	cp.onlyApplyToRegexp_ = append([]string(nil), m.onlyApplyToRegexp_...)
	cp.value_ = append([]string(nil), m.value_...)
	cp.configFile_ = append([]string(nil), m.configFile_...)
	cp.whitelistEntry_ = append([]*RequirementScopeEntry(nil), m.whitelistEntry_...)
	cp.allowlistEntry_ = append([]*RequirementScopeEntry(nil), m.allowlistEntry_...)
	cp.cachedSize = 0
	return &cp
}

func (m *Requirement) Clear() *Requirement {
	*m = Requirement{}
	return m
}

func (m *Requirement) ClearWhitelistEntry() *Requirement {
	m.whitelistEntry_ = nil
	m.cachedSize = 0
	return m
}

func (m *Requirement) ClearAllowlistEntry() *Requirement {
	m.allowlistEntry_ = nil
	m.cachedSize = 0
	return m
}

func (m *Requirement) PopulatedFieldNames() []string {
	var out []string
	if m.HasErrorMessage() {
		out = append(out, "error_message")
	}
	if len(m.whitelist_) > 0 {
		out = append(out, "whitelist")
	}
	if len(m.whitelistRegexp_) > 0 {
		out = append(out, "whitelist_regexp")
	}
	if len(m.onlyApplyTo_) > 0 {
		out = append(out, "only_apply_to")
	}
	if len(m.onlyApplyToRegexp_) > 0 {
		out = append(out, "only_apply_to_regexp")
	}
	if m.HasType() {
		out = append(out, "type")
	}
	if len(m.value_) > 0 {
		out = append(out, "value")
	}
	if m.HasJavaClass() {
		out = append(out, "java_class")
	}
	if m.HasRuleId() {
		out = append(out, "rule_id")
	}
	if m.HasExtends() {
		out = append(out, "extends")
	}
	if m.HasReportLooseTypeViolations() {
		out = append(out, "report_loose_type_violations")
	}
	if m.HasSeverity() {
		out = append(out, "severity")
	}
	if m.HasTypeMatchingStrategy() {
		out = append(out, "type_matching_strategy")
	}
	if len(m.whitelistEntry_) > 0 {
		out = append(out, "whitelist_entry")
	}
	if m.HasAllowExtendingValue() {
		out = append(out, "allow_extending_value")
	}
	if len(m.configFile_) > 0 {
		out = append(out, "config_file")
	}
	if len(m.allowlist_) > 0 {
		out = append(out, "allowlist")
	}
	if len(m.allowlistRegexp_) > 0 {
		out = append(out, "allowlist_regexp")
	}
	if len(m.allowlistEntry_) > 0 {
		out = append(out, "allowlist_entry")
	}
	return out
}

// ---------------------------------------------------------------------------
// Generic ProtoMessage serialization & delimited I/O helpers
// ---------------------------------------------------------------------------

type ProtoSerializable interface {
	serializedSize() int
	appendTo(b []byte) []byte
}

const protoStreamBufSize = 64 * 1024

var protoStreamBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, protoStreamBufSize)
		return &b
	},
}

type protoStreamWritable interface {
	ProtoSerializable
	writeStreamTo(w *protoStreamWriter)
}

type protoStreamWriter struct {
	out     io.Writer
	buf     []byte
	written int64
	err     error
}

func (w *protoStreamWriter) flush() {
	if w.err != nil || len(w.buf) == 0 {
		return
	}
	n, err := w.out.Write(w.buf)
	w.written += int64(n)
	w.buf = w.buf[:0]
	if err != nil {
		w.err = err
	}
}

func (w *protoStreamWriter) ensureSpace(need int) {
	if cap(w.buf)-len(w.buf) < need {
		w.flush()
	}
}

func (w *protoStreamWriter) writeRawBytes(p []byte) {
	if w.err != nil || len(p) == 0 {
		return
	}
	if len(p) <= cap(w.buf)-len(w.buf) {
		w.buf = append(w.buf, p...)
		return
	}
	w.flush()
	if w.err != nil {
		return
	}
	if len(p) >= protoStreamBufSize/2 {
		n, err := w.out.Write(p)
		w.written += int64(n)
		if err != nil {
			w.err = err
		}
		return
	}
	w.buf = append(w.buf, p...)
}

func (w *protoStreamWriter) writeBytesField(fieldNum int, p []byte) {
	w.ensureSpace(20)
	if w.err != nil {
		return
	}
	w.buf = protoAppendTag(w.buf, fieldNum, protoWireBytes)
	w.buf = protoAppendVarint(w.buf, uint64(len(p)))
	w.writeRawBytes(p)
}

func (w *protoStreamWriter) writeStringField(fieldNum int, s string) {
	if w.err != nil {
		return
	}
	if utf8.ValidString(s) {
		w.ensureSpace(20)
		if w.err != nil {
			return
		}
		w.buf = protoAppendTag(w.buf, fieldNum, protoWireBytes)
		w.buf = protoAppendVarint(w.buf, uint64(len(s)))
		if len(s) <= cap(w.buf)-len(w.buf) {
			w.buf = append(w.buf, s...)
			return
		}
		for len(s) > 0 && w.err == nil {
			avail := cap(w.buf) - len(w.buf)
			if avail == 0 {
				w.flush()
				if w.err != nil {
					return
				}
				avail = cap(w.buf)
			}
			take := min(len(s), avail)
			w.buf = append(w.buf, s[:take]...)
			s = s[take:]
		}
		return
	}
	w.writeBytesField(fieldNum, javanio.StandardCharsets_UTF_8.Encode(s))
}

func (w *protoStreamWriter) writeBoolField(fieldNum int, v bool) {
	w.ensureSpace(12)
	if w.err != nil {
		return
	}
	w.buf = protoAppendTag(w.buf, fieldNum, protoWireVarint)
	w.buf = protoAppendBool(w.buf, v)
}

func (w *protoStreamWriter) writeInt32Field(fieldNum int, v int32) {
	w.ensureSpace(20)
	if w.err != nil {
		return
	}
	w.buf = protoAppendTag(w.buf, fieldNum, protoWireVarint)
	w.buf = protoAppendInt32(w.buf, v)
}

func (w *protoStreamWriter) writeUint32Field(fieldNum int, v int32) {
	w.ensureSpace(16)
	if w.err != nil {
		return
	}
	w.buf = protoAppendTag(w.buf, fieldNum, protoWireVarint)
	w.buf = protoAppendUint32(w.buf, v)
}

func (w *protoStreamWriter) writePackedInt32Field(fieldNum int, vals []int32) {
	if w.err != nil || len(vals) == 0 {
		return
	}
	packedLen := 0
	for _, x := range vals {
		packedLen += protoVarintLen(uint64(int64(x)))
	}
	w.ensureSpace(20)
	if w.err != nil {
		return
	}
	w.buf = protoAppendTag(w.buf, fieldNum, protoWireBytes)
	w.buf = protoAppendVarint(w.buf, uint64(packedLen))
	for _, x := range vals {
		w.ensureSpace(12)
		if w.err != nil {
			return
		}
		w.buf = protoAppendInt32(w.buf, x)
	}
}

func (w *protoStreamWriter) writeSubMessage(fieldNum int, sub ProtoSerializable) {
	if w.err != nil {
		return
	}
	sz := sub.serializedSize()
	w.ensureSpace(20)
	if w.err != nil {
		return
	}
	w.buf = protoAppendTag(w.buf, fieldNum, protoWireBytes)
	w.buf = protoAppendVarint(w.buf, uint64(sz))
	if sz == 0 {
		return
	}
	if sw, ok := sub.(protoStreamWritable); ok {
		sw.writeStreamTo(w)
		return
	}
	if sz <= cap(w.buf) {
		w.ensureSpace(sz)
		if w.err != nil {
			return
		}
		w.buf = sub.appendTo(w.buf)
		return
	}
	w.writeRawBytes(MessageToByteArray(sub))
}

func (l *TypedAst_List) writeStreamTo(w *protoStreamWriter) {
	for _, ast := range l.typedAsts_ {
		w.writeSubMessage(1, ast)
		if w.err != nil {
			return
		}
	}
}

func (m *TypedAst) writeStreamTo(w *protoStreamWriter) {
	if m.typePool_ != nil {
		w.writeSubMessage(1, m.typePool_)
	}
	if m.stringPool_ != nil {
		w.writeSubMessage(2, m.stringPool_)
	}
	for _, x := range m.externAst_ {
		w.writeSubMessage(3, x)
		if w.err != nil {
			return
		}
	}
	for _, x := range m.codeAst_ {
		w.writeSubMessage(4, x)
		if w.err != nil {
			return
		}
	}
	if m.sourceFilePool_ != nil {
		w.writeSubMessage(5, m.sourceFilePool_)
	}
	if m.externsSummary_ != nil {
		w.writeSubMessage(6, m.externsSummary_)
	}
	for _, x := range m.runtimeLibraryToInject_ {
		w.writeStringField(7, x)
		if w.err != nil {
			return
		}
	}
}

func (m *TypePool) writeStreamTo(w *protoStreamWriter) {
	for _, x := range m.type_ {
		w.writeSubMessage(1, x)
		if w.err != nil {
			return
		}
	}
	for _, x := range m.disambiguationEdges_ {
		w.writeSubMessage(2, x)
		if w.err != nil {
			return
		}
	}
	if m.debugInfo_ != nil {
		w.writeSubMessage(3, m.debugInfo_)
	}
	if len(m.unknownFields) > 0 {
		w.writeRawBytes(protoCanonicalUnknownFields(m.unknownFields))
	}
}

func (m *StringPoolProto) writeStreamTo(w *protoStreamWriter) {
	for _, c := range m.wireChunks_ {
		w.writeRawBytes(c)
		if w.err != nil {
			return
		}
	}
	if m.maxLength_ != 0 {
		w.writeUint32Field(2, m.maxLength_)
	}
}

func (m *LazyAst) writeStreamTo(w *protoStreamWriter) {
	if len(m.script_) != 0 {
		w.writeBytesField(1, m.script_)
	}
	if m.sourceFile_ != 0 {
		w.writeUint32Field(2, m.sourceFile_)
	}
	if m.sourceMappingUrl_ != "" {
		w.writeStringField(3, m.sourceMappingUrl_)
	}
}

func (m *SourceFilePool) writeStreamTo(w *protoStreamWriter) {
	for _, x := range m.sourceFile_ {
		w.writeSubMessage(1, x)
		if w.err != nil {
			return
		}
	}
}

func (m *ExternsSummary) writeStreamTo(w *protoStreamWriter) {
	w.writePackedInt32Field(1, m.propNamePtr_)
}

func (m *ChunkProto) writeStreamTo(w *protoStreamWriter) {
	if m.name_ != "" {
		w.writeStringField(1, m.name_)
	}
	w.writePackedInt32Field(2, m.dependencies_)
	for _, x := range m.inputIds_ {
		w.writeStringField(3, x)
		if w.err != nil {
			return
		}
	}
}

func (m *AccessorSummaryProto) writeStreamTo(w *protoStreamWriter) {
	if m.assumeAlwaysGetterAndSetter_ {
		w.writeBoolField(1, true)
	}
	for _, x := range m.accessors_ {
		w.writeSubMessage(2, x)
		if w.err != nil {
			return
		}
	}
}

func (m *JSCompilerStateProto) writeStreamTo(w *protoStreamWriter) {
	w.writePackedInt32Field(1, m.allowableFeatures_)
	if m.typeCheckingHasRun_ {
		w.writeBoolField(2, true)
	}
	if m.hasRegExpGlobalReferences_ {
		w.writeBoolField(3, true)
	}
	if m.lifeCycleStage_ != 0 {
		w.writeInt32Field(4, m.lifeCycleStage_)
	}
	if m.mergedPrecompiledLibraries_ {
		w.writeBoolField(5, true)
	}
	for _, x := range m.chunks_ {
		w.writeSubMessage(6, x)
		if w.err != nil {
			return
		}
	}
	if m.uniqueNameId_ != 0 {
		w.writeInt32Field(7, m.uniqueNameId_)
	}
	for _, x := range m.uniqueIdSupplier_ {
		w.writeSubMessage(8, x)
		if w.err != nil {
			return
		}
	}
	for _, x := range m.exportedNames_ {
		w.writeStringField(9, x)
		if w.err != nil {
			return
		}
	}
	if m.hasCssNames_ {
		w.writeBoolField(10, true)
	}
	for _, x := range m.cssNames_ {
		w.writeStringField(11, x)
		if w.err != nil {
			return
		}
	}
	if m.bitField0_&(1<<0) != 0 {
		w.writeStringField(12, m.idGeneratorMap_)
	}
	if m.transpiledFiles_ {
		w.writeBoolField(13, true)
	}
	if m.idGeneratorCurrentId_ != 0 {
		w.writeInt32Field(14, m.idGeneratorCurrentId_)
	}
	if m.runJ2ClPasses_ {
		w.writeBoolField(15, true)
	}
	for _, x := range m.externs_ {
		w.writeStringField(16, x)
		if w.err != nil {
			return
		}
	}
	for _, x := range m.injectedLibraries_ {
		w.writeStringField(17, x)
		if w.err != nil {
			return
		}
	}
	if m.lastInjectedLibraryIndexInFirstScript_ != 0 {
		w.writeInt32Field(18, m.lastInjectedLibraryIndexInFirstScript_)
	}
	if m.accessorSummary_ != nil {
		w.writeSubMessage(19, m.accessorSummary_)
	}
	if m.hasStringMap_ {
		w.writeBoolField(20, true)
	}
	for _, x := range m.stringMap_ {
		w.writeSubMessage(21, x)
		if w.err != nil {
			return
		}
	}
	if m.hasInstrumentationMapping_ {
		w.writeBoolField(22, true)
	}
	for _, x := range m.instrumentationMapping_ {
		w.writeSubMessage(23, x)
		if w.err != nil {
			return
		}
	}
	for _, x := range m.conformanceViolations_ {
		w.writeSubMessage(24, x)
		if w.err != nil {
			return
		}
	}
	for _, x := range m.unattributedConformanceViolations_ {
		w.writeSubMessage(25, x)
		if w.err != nil {
			return
		}
	}
}

func streamProtoMessageTo(m ProtoSerializable, out io.Writer, delimited bool) (int64, error) {
	if m == nil {
		if delimited {
			n, err := out.Write([]byte{0})
			return int64(n), err
		}
		return 0, nil
	}
	bp := protoStreamBufPool.Get().(*[]byte)
	w := protoStreamWriter{
		out: out,
		buf: (*bp)[:0],
	}
	if delimited {
		sz := m.serializedSize()
		w.buf = protoAppendVarint(w.buf, uint64(sz))
		if sz > 0 {
			if sw, ok := m.(protoStreamWritable); ok {
				sw.writeStreamTo(&w)
			} else if sz <= cap(w.buf)-len(w.buf) {
				w.buf = m.appendTo(w.buf)
			} else {
				w.writeRawBytes(m.appendTo(make([]byte, 0, sz)))
			}
		}
	} else if sw, ok := m.(protoStreamWritable); ok {
		sw.writeStreamTo(&w)
	} else {
		sz := m.serializedSize()
		if sz > 0 {
			if sz <= cap(w.buf) {
				w.buf = m.appendTo(w.buf)
			} else {
				w.writeRawBytes(m.appendTo(make([]byte, 0, sz)))
			}
		}
	}
	w.flush()
	if cap(w.buf) == protoStreamBufSize {
		*bp = w.buf[:0]
		protoStreamBufPool.Put(bp)
	}
	return w.written, w.err
}

func panicProtoWriteErr(err error) {
	if th, ok := err.(javalang.Throwable); ok {
		panic(th)
	}
	panic(javaio.NewIOException(err.Error()))
}

func MessageToByteArray(m ProtoSerializable) []byte {
	if m == nil {
		return nil
	}
	sz := m.serializedSize()
	return m.appendTo(make([]byte, 0, sz))
}

func MessageToByteString(m ProtoSerializable) ByteString {
	if ast, ok := m.(*AstNode); ok {
		return ast.ToByteString()
	}
	return ByteString{data: MessageToByteArray(m)}
}

func WriteMessageTo(m ProtoSerializable, out io.Writer) {
	if _, err := streamProtoMessageTo(m, out, false); err != nil {
		panicProtoWriteErr(err)
	}
}

func WriteMessageDelimitedTo(m ProtoSerializable, out io.Writer) {
	if _, err := streamProtoMessageTo(m, out, true); err != nil {
		panicProtoWriteErr(err)
	}
}

func ReadDelimitedBytes(in io.Reader) ([]byte, bool) {
	br := protoByteReader(in)
	sz64, ok, th := protoReadStreamVarint64(br, true)
	if th != nil {
		panic(th)
	}
	if !ok {
		return nil, false
	}
	if sz64 > math.MaxInt32 {
		panic(protoNegativeSize())
	}
	data := make([]byte, int(sz64))
	if int(sz64) > 0 {
		if th := protoReadFull(in, data); th != nil {
			panic(th)
		}
	}
	return data, true
}

func JSCompilerStateProto_ParseDelimitedFrom(in io.Reader) *JSCompilerStateProto {
	data, ok := ReadDelimitedBytes(in)
	if !ok {
		return nil
	}
	return JSCompilerStateProto_ParseFrom(data)
}

// ---------------------------------------------------------------------------
// CompilerState sub-message Builders & Oneof / Count helpers
// ---------------------------------------------------------------------------

type ChunkProto_Builder = ChunkProto

func ChunkProto_NewBuilder() *ChunkProto_Builder { return &ChunkProto{} }

func (m *ChunkProto) Clear() *ChunkProto {
	*m = ChunkProto{}
	return m
}
func (m *ChunkProto) Build() *ChunkProto {
	m.cachedSize = 0
	return m
}
func (m *ChunkProto) SetName(v string) *ChunkProto {
	m.name_ = v
	m.cachedSize = 0
	return m
}
func (m *ChunkProto) AddDependencies(v int32) *ChunkProto {
	m.dependencies_ = append(m.dependencies_, v)
	m.cachedSize = 0
	return m
}
func (m *ChunkProto) AddInputIds(v string) *ChunkProto {
	m.inputIds_ = append(m.inputIds_, v)
	m.cachedSize = 0
	return m
}

type UniqueIdProto_Builder = UniqueIdProto

func UniqueIdProto_NewBuilder() *UniqueIdProto_Builder { return &UniqueIdProto{} }

func (m *UniqueIdProto) Clear() *UniqueIdProto {
	*m = UniqueIdProto{}
	return m
}
func (m *UniqueIdProto) Build() *UniqueIdProto {
	m.cachedSize = 0
	return m
}
func (m *UniqueIdProto) SetHash(v int32) *UniqueIdProto {
	m.hash_ = v
	m.cachedSize = 0
	return m
}
func (m *UniqueIdProto) SetCounter(v int32) *UniqueIdProto {
	m.counter_ = v
	m.cachedSize = 0
	return m
}

type AccessorSummaryEntryProto_Builder = AccessorSummaryEntryProto

func AccessorSummaryEntryProto_NewBuilder() *AccessorSummaryEntryProto_Builder {
	return &AccessorSummaryEntryProto{}
}
func (m *AccessorSummaryEntryProto) Clear() *AccessorSummaryEntryProto {
	*m = AccessorSummaryEntryProto{}
	return m
}
func (m *AccessorSummaryEntryProto) Build() *AccessorSummaryEntryProto {
	m.cachedSize = 0
	return m
}
func (m *AccessorSummaryEntryProto) SetName(v string) *AccessorSummaryEntryProto {
	m.name_ = v
	m.cachedSize = 0
	return m
}
func (m *AccessorSummaryEntryProto) SetKind(v PropertyAccessKindProto) *AccessorSummaryEntryProto {
	m.kind_ = v.GetNumber()
	m.cachedSize = 0
	return m
}

type AccessorSummaryProto_Builder = AccessorSummaryProto

func AccessorSummaryProto_NewBuilder() *AccessorSummaryProto_Builder {
	return &AccessorSummaryProto{}
}
func (m *AccessorSummaryProto) Clear() *AccessorSummaryProto {
	*m = AccessorSummaryProto{}
	return m
}
func (m *AccessorSummaryProto) Build() *AccessorSummaryProto {
	m.cachedSize = 0
	return m
}
func (m *AccessorSummaryProto) SetAssumeAlwaysGetterAndSetter(v bool) *AccessorSummaryProto {
	m.assumeAlwaysGetterAndSetter_ = v
	m.cachedSize = 0
	return m
}
func (m *AccessorSummaryProto) AddAccessors(v *AccessorSummaryEntryProto) *AccessorSummaryProto {
	m.accessors_ = append(m.accessors_, v)
	m.cachedSize = 0
	return m
}

type VariableMapEntryProto_Builder = VariableMapEntryProto

func VariableMapEntryProto_NewBuilder() *VariableMapEntryProto_Builder {
	return &VariableMapEntryProto{}
}
func (m *VariableMapEntryProto) Clear() *VariableMapEntryProto {
	*m = VariableMapEntryProto{}
	return m
}
func (m *VariableMapEntryProto) Build() *VariableMapEntryProto {
	m.cachedSize = 0
	return m
}
func (m *VariableMapEntryProto) SetOriginalName(v string) *VariableMapEntryProto {
	m.originalName_ = v
	m.cachedSize = 0
	return m
}
func (m *VariableMapEntryProto) SetNewName(v string) *VariableMapEntryProto {
	m.newName_ = v
	m.cachedSize = 0
	return m
}

func (m *JSCompilerStateProto) GetChunks(index int) *ChunkProto {
	return m.chunks_[index]
}

type TypeProto_KindCase int32

const (
	TypeProto_KindCase_OBJECT       TypeProto_KindCase = 0
	TypeProto_KindCase_UNION        TypeProto_KindCase = 1
	TypeProto_KindCase_KIND_NOT_SET TypeProto_KindCase = 2
)

var typeProtoKindCaseNames = [...]string{"OBJECT", "UNION", "KIND_NOT_SET"}
var typeProtoKindCaseNums = [...]int32{2, 4, 0}

func (c TypeProto_KindCase) String() string {
	if int(c) >= 0 && int(c) < len(typeProtoKindCaseNames) {
		return typeProtoKindCaseNames[c]
	}
	return "KIND_NOT_SET"
}
func (c TypeProto_KindCase) GetNumber() int32 {
	if int(c) >= 0 && int(c) < len(typeProtoKindCaseNums) {
		return typeProtoKindCaseNums[c]
	}
	return 0
}

func (m *TypeProto) GetKindCase() TypeProto_KindCase {
	switch m.kindCase_ {
	case 2:
		return TypeProto_KindCase_OBJECT
	case 4:
		return TypeProto_KindCase_UNION
	default:
		return TypeProto_KindCase_KIND_NOT_SET
	}
}

var sourceFileLoaderCaseNames = [...]string{
	"PRELOADED_CONTENTS",
	"FILE_ON_DISK",
	"ZIP_ENTRY",
	"STUB_FILE",
	"LOADER_NOT_SET",
}
var sourceFileLoaderCaseNums = [...]int32{2, 3, 4, 9, 0}

func (c SourceFileProto_LoaderCase) String() string {
	if int(c) >= 0 && int(c) < len(sourceFileLoaderCaseNames) {
		return sourceFileLoaderCaseNames[c]
	}
	return "LOADER_NOT_SET"
}
func (c SourceFileProto_LoaderCase) GetNumber() int32 {
	if int(c) >= 0 && int(c) < len(sourceFileLoaderCaseNums) {
		return sourceFileLoaderCaseNums[c]
	}
	return 0
}

func (m *Mapping_OriginalMapping) HasIdentifier() bool {
	return m.Identifier != ""
}

func (m *ObjectTypeProto) GetOwnPropertyCount() int   { return len(m.ownProperty_) }
func (m *ObjectTypeProto) GetOwnProperty(i int) int32 { return m.ownProperty_[i] }
func (m *OptimizationJsdoc) GetKindCount() int {
	if m.kind_ != nil {
		return len(m.kind_)
	}
	return int(m.kindInlineLen_)
}
func (m *UnionTypeProto) GetUnionMemberCount() int    { return len(m.unionMember_) }
func (m *TypePool) GetTypeCount() int                 { return len(m.type_) }
func (m *Requirement) GetWhitelistCount() int         { return len(m.whitelist_) }
func (m *Requirement) GetWhitelistRegexpCount() int   { return len(m.whitelistRegexp_) }
func (m *Requirement) GetAllowlistCount() int         { return len(m.allowlist_) }
func (m *Requirement) GetAllowlistRegexpCount() int   { return len(m.allowlistRegexp_) }
func (m *Requirement) GetOnlyApplyToCount() int       { return len(m.onlyApplyTo_) }
func (m *Requirement) GetOnlyApplyToRegexpCount() int { return len(m.onlyApplyToRegexp_) }
func (m *Requirement) GetValueCount() int             { return len(m.value_) }
func (m *Requirement) GetConfigFileCount() int        { return len(m.configFile_) }
