// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261104

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package nodestore

// Node is a 32-bit handle into the chunked arena of AST node records.
// NilNode is 0 (Java null Node).
type Node uint32

// NilNode is the nil node handle.
const NilNode Node = 0

// NodeData is the pointer-free 32-byte in-memory representation of an AST node.
// Because it contains zero Go pointers, blocks of NodeData allocated in chunked
// arrays are marked noscan by the Go runtime, eliminating GC traversal overhead.
// At 32 bytes (exact half of a 64-byte L1D cache line), records never straddle
// cache lines and block indexing uses a single << 5 shift.
type NodeData struct {
	Next         Node   // 0..3: next sibling (0 = none)
	First        Node   // 4..7: first child (0 = none)
	Token        uint8  // 8: AST token kind
	KindFlags    uint8  // 9: bits 0..2 kind; bits 3..7 internal flags
	SourceFile16 uint16 // 10..11: 0 = none, 1..65534 = inline StaticSourceFile handle, 0xFFFF = overflow
	Parent       Node   // 12..15: parent node (0 = none)
	Previous     Node   // 16..19: previous sibling (circular linked list)
	Payload      uint32 // 20..23: kind-specific payload (StringID, numHandle, or inline ChangeTime/OrigName on plain kind=0 nodes)
	PropFlags    uint32 // 24..27: boolean & small-int property bitflags
	LinenoCharno int32  // 28..31: line (high 20 bits) and column (low 12 bits)
}

// IsNil reports whether n is the nil node (handle 0).
func (n Node) IsNil() bool {
	return n == 0
}

// IsJavaNil_ implements javalang.JavaNilChecker for reflection-free nil checks.
func (n Node) IsJavaNil_() bool {
	return n == 0
}

// ID returns the 32-bit node handle as an integer.
func (n Node) ID() uint32 {
	return uint32(n)
}

// HashCode returns a 32-bit hash of the node handle matching Java Object.hashCode().
func (n Node) HashCode() int32 {
	return int32(uint32(n) * 0x9E3779B1)
}

// Equals reports whether o is a Node with the same handle.
func (n Node) Equals(o any) bool {
	m, ok := o.(Node)
	return ok && m == n
}

// Rec returns a pointer to the 32-byte NodeData record in the arena.
// Record lookup requires two memory dereferences without mutexes; an unmapped
// or nil handle faults on the nil *nodeBlock dereference (recovered as NPE).
func (n Node) Rec() *NodeData {
	return &nodeBlocks[n>>blockShift].Load()[n&blockMask]
}

// Next returns the next sibling node.
func (n Node) Next() Node {
	if n == 0 {
		return 0
	}
	return n.Rec().Next
}

// SetNext sets the next sibling node.
func (n Node) SetNext(next Node) {
	n.Rec().Next = next
}

// FirstChild returns the first child node.
func (n Node) FirstChild() Node {
	if n == 0 {
		return 0
	}
	return n.Rec().First
}

// SetFirstChild sets the first child node.
func (n Node) SetFirstChild(first Node) {
	n.Rec().First = first
}

// Parent returns the parent node.
func (n Node) Parent() Node {
	if n == 0 {
		return 0
	}
	return n.Rec().Parent
}

// SetParent sets the parent node.
func (n Node) SetParent(parent Node) {
	n.Rec().Parent = parent
}

// Previous returns the previous sibling node.
func (n Node) Previous() Node {
	if n == 0 {
		return 0
	}
	return n.Rec().Previous
}

// SetPrevious sets the previous sibling node.
func (n Node) SetPrevious(prev Node) {
	n.Rec().Previous = prev
}

// Token returns the node's token kind.
func (n Node) Token() uint8 {
	return n.Rec().Token
}

// SetToken sets the node's token kind.
func (n Node) SetToken(tok uint8) {
	n.Rec().Token = tok
}

// BlockOf returns the *NodeBlock containing non-nil Node n (panicking with runtime.Error on unmapped/nil handle just like n.Rec()).
func BlockOf(n Node) *NodeBlock {
	b := nodeBlocks[n>>blockShift].Load()
	_ = &b[0] // preserve immediate nil-deref runtime.Error if unmapped/NilNode without reading b[0]
	return b
}

// FirstNextTok reads First, Next, and Token directly from b without a nodeBlocks directory load.
// Caller guarantees uint32(n)>>BlockShift corresponds to b.
func (b *NodeBlock) FirstNextTok(n Node) (first, next Node, tok uint8) {
	rec := &b[n&blockMask]
	return rec.First, rec.Next, rec.Token
}

// FirstNextTok loads First, Next, and Token in a single Rec() lookup.
func (n Node) FirstNextTok() (first, next Node, tok uint8) {
	rec := n.Rec()
	return rec.First, rec.Next, rec.Token
}

// HasColdProps reports whether n currently has entries in the cold property list.
func (n Node) HasColdProps() bool {
	return n.Rec().KindFlags&kindFlagHasColdProps != 0
}

// Lineno returns the 1-based source line number, or -1 if unknown.
func (n Node) Lineno() int {
	lc := n.Rec().LinenoCharno
	if lc < 0 {
		return -1
	}
	return int(lc >> 12)
}

// Charno returns the 0-based source character offset within the line, or -1 if unknown.
func (n Node) Charno() int {
	lc := n.Rec().LinenoCharno
	if lc < 0 {
		return -1
	}
	return int(lc & 0xFFF)
}

// SetLinenoCharno sets the line and character offset, storing -1 when either is negative.
func (n Node) SetLinenoCharno(lineno, charno int) {
	rec := n.Rec()
	if (lineno | charno) < 0 {
		rec.LinenoCharno = -1
		return
	}
	if charno > 0xFFF {
		charno = 0xFFF
	}
	rec.LinenoCharno = int32((lineno << 12) | charno)
}

// HasNonNegativeLength reports whether n.Length() >= 0 without reading len16Top
// when no active node has a negative length. Panics with runtime.Error on NilNode
// to match rhino.Node.GetLength().
func (n Node) HasNonNegativeLength() bool {
	return hasNonNegativeLength(n, hasNonNegativeLengthSlow)
}

func hasNonNegativeLength(n Node, slow func(Node) bool) bool {
	if n != 0 && negLenCount.Load() == 0 {
		return true
	}
	return slow(n)
}

//go:noinline
func hasNonNegativeLengthSlow(n Node) bool {
	if n == 0 {
		_ = (*NodeData)(nil).Token
	}
	return n.Length() >= 0
}

// Length returns the source code character length (0..65534 in len16Top, or from lenOverflowTop).
func (n Node) Length() int32 {
	return getNodeLength(n, getLengthOverflowSlow)
}

func getNodeLength(n Node, slow func(Node) int32) int32 {
	if l16 := len16Top[n>>len16TopShift].Load()[uint8(n>>len16PageShift)].Load()[n&len16PageMask]; l16 != 0xFFFF {
		return int32(l16)
	}
	return slow(n)
}

//go:noinline
func getLengthOverflowSlow(n Node) int32 {
	if p := loadLenOverflowPage(uint32(n) >> lenOverflowShift); p != nil {
		return int32(p[uint32(n)&lenOverflowMask])
	}
	return 0
}

// SetLength sets the source code character length (0..65534 lazily in len16Top, or in lenOverflowTop / NodeStore.aux).
func (n Node) SetLength(v int32) {
	setNodeLength(n, v, setLengthSlow)
}

func setNodeLength(n Node, v int32, slow func(Node, int32)) {
	if p := &len16WriteTop[n>>len16TopShift].Load()[uint8(n>>len16PageShift)].Load()[n&len16PageMask]; uint32(v)|uint32(*p) < 0xFFFF {
		*p = uint16(v)
	} else {
		slow(n, v)
	}
}

//go:noinline
func setLengthSlow(n Node, v int32) {
	_ = n.Rec().Token
	idx := uint32(n) >> len16PageShift
	if v == 0 && loadLen16Page(idx) == &emptyLen16Page {
		return
	}
	pg := ensureLen16Page(n, idx)
	off := uint32(n) & len16PageMask
	s := storeForNode(n)
	oldNeg := pg[off] == 0xFFFF && getLengthOverflowSlow(n) < 0
	if !oldNeg && v < 0 {
		s.negLenNodes.Add(1)
		negLenCount.Add(1)
	}
	if uint32(v) < 0xFFFF {
		if pg[off] == 0xFFFF {
			setLenOverflow(n, 0)
			s.setAuxLengthOverflow(uint32(n), 0)
		}
		pg[off] = uint16(v)
	} else {
		setLenOverflow(n, uint32(v))
		s.setAuxLengthOverflow(uint32(n), v)
		pg[off] = 0xFFFF
	}
	if oldNeg && v >= 0 {
		s.negLenNodes.Add(^uint32(0))
		negLenCount.Add(^uint32(0))
	}
}

// SourceFile returns the source file identifier (0..65534 inline in SourceFile16, or from sfOverflowTop).
func (n Node) SourceFile() uint32 {
	return getNodeSourceFile(n, getSourceFileOverflowSlow)
}

func getNodeSourceFile(n Node, slow func(Node) uint32) uint32 {
	if sf16 := n.Rec().SourceFile16; sf16 != 0xFFFF {
		return uint32(sf16)
	}
	return slow(n)
}

//go:noinline
func getSourceFileOverflowSlow(n Node) uint32 {
	if p := loadSFOverflowPage(uint32(n) >> sfOverflowShift); p != nil {
		return p[uint32(n)&sfOverflowMask]
	}
	return 0
}

// SetSourceFile sets the source file identifier (0..65534 inline in SourceFile16, or in sfOverflowTop).
func (n Node) SetSourceFile(sf uint32) {
	setNodeSourceFile(n, sf, setSourceFileOverflowSlow)
}

func setNodeSourceFile(n Node, sf uint32, slow func(Node, *NodeData, uint32)) {
	rec := n.Rec()
	if sf|uint32(rec.SourceFile16) < 0xFFFF {
		rec.SourceFile16 = uint16(sf)
		return
	}
	slow(n, rec, sf)
}

//go:noinline
func setSourceFileOverflowSlow(n Node, rec *NodeData, sf uint32) {
	if sf < 0xFFFF {
		setSFOverflow(n, 0)
		rec.SourceFile16 = uint16(sf)
		return
	}
	setSFOverflow(n, sf)
	rec.SourceFile16 = 0xFFFF
}

// Payload returns the raw payload word.
func (n Node) Payload() uint32 {
	return n.Rec().Payload
}

// SetPayload sets the raw payload word, preserving any inline original name in origNameField
// and any inline plain-node ChangeTime in changeTimeField before mutating Payload.
func (n Node) SetPayload(p uint32) {
	rec := n.Rec()
	if p != rec.Payload {
		if rec.KindFlags&kindFlagInlineOrigName != 0 {
			n.materializeOrigName(rec, rec.Payload)
		} else if rec.KindFlags&(kindFlagsMask|kindFlagHasChangeTime) == kindFlagHasChangeTime {
			n.migrateInlineChangeTime(rec)
		}
	}
	rec.Payload = p
}

// SetStringPayload sets the interned StringID payload on a StringNode, materializing
// any inline original name into origNameField before mutating Payload.
func (n Node) SetStringPayload(id StringID) {
	setStringPayload(n, id, materializeOrigNameSlow)
}

func setStringPayload(n Node, id StringID, slow func(Node, *NodeData, uint32)) {
	rec := n.Rec()
	old := rec.Payload
	if uint32(id) == old {
		return
	}
	if rec.KindFlags&kindFlagInlineOrigName != 0 {
		slow(n, rec, old)
	}
	rec.Payload = uint32(id)
}

//go:noinline
func materializeOrigNameSlow(n Node, rec *NodeData, old uint32) {
	if rec.KindFlags&kindFlagsMask == 0 {
		rec.Payload = 0
	}
	if old != 0 {
		origNameField.set(n, old)
		rec.KindFlags = (rec.KindFlags &^ kindFlagInlineOrigName) | kindFlagHasExtOrigName
	} else {
		rec.KindFlags &^= kindFlagInlineOrigName
	}
}

func (n Node) materializeOrigName(rec *NodeData, old uint32) {
	materializeOrigNameSlow(n, rec, old)
}

//go:noinline
func (n Node) migrateInlineChangeTime(rec *NodeData) {
	if rec.Payload != 0 {
		s := storeForNode(n)
		if s.getAuxChangeTime(uint32(n)) == 0 {
			s.setAuxChangeTime(uint32(n), int32(rec.Payload))
			rec.Payload = 0
		}
	}
}

const (
	kindFlagHasColdProps      uint8 = 1 << 3
	kindFlagHasExtOrigName    uint8 = 1 << 4 // origNameField has non-zero entry for n
	kindFlagHasChangeTime     uint8 = 1 << 5
	kindFlagHasJSDocOrCast    uint8 = 1 << 6 // aux has jsDoc != nil || typeBeforeCast != nil
	kindFlagHasJSDoc          uint8 = kindFlagHasJSDocOrCast
	kindFlagHasTypeBeforeCast uint8 = kindFlagHasJSDocOrCast
	kindFlagInlineOrigName    uint8 = 1 << 7 // kind==3: OrigNameID == Payload (str); kind==0 && bit5==0: OrigNameID == Payload
	kindFlagAnyOrigName       uint8 = kindFlagInlineOrigName | kindFlagHasExtOrigName
	kindFlagInternalMask      uint8 = kindFlagInlineOrigName | kindFlagHasJSDocOrCast | kindFlagHasChangeTime | kindFlagHasExtOrigName | kindFlagHasColdProps
	kindFlagsMask             uint8 = ^kindFlagInternalMask // 0x07
	kindFlagClonedMask        uint8 = kindFlagHasColdProps | kindFlagHasJSDocOrCast | kindFlagHasChangeTime
)

// PropFlags returns the inline property bitmask.
func (n Node) PropFlags() int32 {
	return int32(n.Rec().PropFlags)
}

// SetPropFlags sets the inline property bitmask.
func (n Node) SetPropFlags(v int32) {
	n.Rec().PropFlags = uint32(v)
}

// KindFlags returns the kind flags (low 3 bits, masking internal ColdProps, TypeBeforeCast, ChangeTime, JSDoc, and InlineOrigName flags).
func (n Node) KindFlags() uint8 {
	return n.Rec().KindFlags & kindFlagsMask
}

// SetKindFlags sets the kind flags (low 3 bits), preserving internal flags and migrating inline OrigName or ChangeTime when needed.
func (n Node) SetKindFlags(flags uint8) {
	rec := n.Rec()
	newKind := flags & kindFlagsMask
	oldKind := rec.KindFlags & kindFlagsMask
	if oldKind != newKind {
		if rec.KindFlags&kindFlagInlineOrigName != 0 {
			n.materializeOrigName(rec, rec.Payload)
		} else if oldKind == 0 && rec.KindFlags&kindFlagHasChangeTime != 0 {
			n.migrateInlineChangeTime(rec)
		}
	}
	rec.KindFlags = (rec.KindFlags & kindFlagInternalMask) | newKind
}

// CheckNotNull asserts that n is not NilNode, panicking otherwise.
func CheckNotNull(n Node) Node {
	if n == 0 {
		panic("NullPointerException: node is null")
	}
	return n
}

// CheckNotNullMsg asserts that n is not NilNode with a formatted message.
func CheckNotNullMsg(n Node, msg string) Node {
	if n == 0 {
		panic("NullPointerException: " + msg)
	}
	return n
}
