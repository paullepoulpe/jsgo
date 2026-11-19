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

// NodeData is the pointer-free 36-byte in-memory representation of an AST node.
// Because it contains zero Go pointers, blocks of NodeData allocated in chunked
// arrays are marked noscan by the Go runtime, eliminating GC traversal overhead.
type NodeData struct {
	Next         Node   // 0..3: next sibling (0 = none)
	First        Node   // 4..7: first child (0 = none)
	Token        uint8  // 8: AST token kind
	KindFlags    uint8  // 9: bits 0..2 kind; bits 3..7 internal flags
	Length16     uint16 // 10..11: 0..65534 inline; 0xFFFF = overflow in NodeStore.aux
	Parent       Node   // 12..15: parent node (0 = none)
	Previous     Node   // 16..19: previous sibling (circular linked list)
	Payload      uint32 // 20..23: kind-specific payload (StringID, numHandle, or inline ChangeTime/OrigName on plain kind=0 nodes)
	SourceFile   uint32 // 24..27: 1-based StaticSourceFile handle (0 = none)
	PropFlags    uint32 // 28..31: boolean & small-int property bitflags
	LinenoCharno int32  // 32..35: line (high 20 bits) and column (low 12 bits)
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

// Rec returns a pointer to the 36-byte NodeData record in the arena.
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

// Length returns the source code character length (0..65534 inline in Length16, or from lenOverflowTop).
func (n Node) Length() int32 {
	return getNodeLength(n, getLengthOverflowSlow)
}

func getNodeLength(n Node, slow func(Node) int32) int32 {
	if blk := nodeBlocks[n>>blockShift].Load(); blk != nil {
		if l16 := blk[n&blockMask].Length16; l16 != 0xFFFF {
			return int32(l16)
		}
		return slow(n)
	}
	return 0
}

//go:noinline
func getLengthOverflowSlow(n Node) int32 {
	return int32(lenOverflowTop[n>>lenOverflowTopShift].Load()[(n>>lenOverflowShift)&lenOverflowSubMask].Load()[n&lenOverflowMask])
}

// SetLength sets the source code character length (0..65534 inline in Length16, or in NodeStore.aux / lenOverflowPages).
func (n Node) SetLength(v int32) {
	setNodeLength(n, v, setLengthOverflowSlow)
}

func setNodeLength(n Node, v int32, slow func(Node, int32)) {
	if r := &nodeBlocks[n>>blockShift].Load()[n&blockMask]; uint32(v)|uint32(r.Length16) < 0xFFFF {
		r.Length16 = uint16(v)
	} else {
		slow(n, v)
	}
}

//go:noinline
func setLengthOverflowSlow(n Node, v int32) {
	rec := n.Rec()
	s := storeForNode(n)
	if uint32(v) < 0xFFFF {
		rec.Length16 = uint16(v)
		setLenOverflow(n, 0)
		s.setAuxLengthOverflow(uint32(n), 0)
		return
	}
	setLenOverflow(n, uint32(v))
	s.setAuxLengthOverflow(uint32(n), v)
	rec.Length16 = 0xFFFF
}

// SourceFile returns the source file identifier.
func (n Node) SourceFile() uint32 {
	return n.Rec().SourceFile
}

// SetSourceFile sets the source file identifier.
func (n Node) SetSourceFile(sf uint32) {
	n.Rec().SourceFile = sf
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
