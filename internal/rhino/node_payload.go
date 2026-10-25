package rhino

import (
	"math"
	"sync"
	"sync/atomic"
)

// Go-only: the leaf payload, original name and property bits of a node
// record (nodeData).
//
// The payload is one 32-bit value whose meaning depends on the node's kind:
//
//	node_kindNode                      CHANGE_TIME (0: none)
//	node_kindString                    the interned string
//	node_kindNumber                    the number (numberPayload)
//	node_kindBigInt                    the BigInt's id
//	node_kindTemplateLiteralSubstring  the interned raw string (the cooked
//	                                   string is in the ext box when it is not
//	                                   the raw one)
//
// The original name (the name before renaming, kept for source maps) of a
// string node is usually its own string: kindFlagOrigNameIsString says so.
// Any other original name is in the origName column (node_tables.go).

// --- numbers ---

// A number that is an int31 (and not -0) is stored in the payload itself,
// shifted left by one; any other number is stored in numberTable, and the
// payload is its index shifted left by one, plus one.

func numberPayload(x float64) uint32 {
	if i := int32(x); float64(i) == x && i >= -1<<30 && i < 1<<30 && !(x == 0 && math.Signbit(x)) {
		return uint32(i) << 1
	}
	return numberTable.add(x)<<1 | 1
}

func payloadNumber(p uint32) float64 {
	if p&1 == 0 {
		return float64(int32(p) >> 1)
	}
	return numberTable.at(p >> 1)
}

// numberTable holds the numbers that do not fit in a payload, for the life
// of the process (few numbers in a program are not small integers).
var numberTable wideNumbers

type wideNumbers struct {
	mu     sync.Mutex
	index  map[uint64]uint32 // float64 bits -> position
	next   uint32
	chunks [1 << 12]atomic.Pointer[[1 << 12]float64]
}

func (t *wideNumbers) add(x float64) uint32 {
	bits := math.Float64bits(x)
	t.mu.Lock()
	defer t.mu.Unlock()
	if i, ok := t.index[bits]; ok {
		return i
	}
	if t.index == nil {
		t.index = map[uint64]uint32{}
	}
	i := t.next
	if int(i>>12) >= len(t.chunks) {
		panic("rhino: too many distinct non-integer numbers")
	}
	c := t.chunks[i>>12].Load()
	if c == nil {
		c = new([1 << 12]float64)
		t.chunks[i>>12].Store(c)
	}
	c[i&(1<<12-1)] = x
	t.index[bits] = i
	t.next++
	return i
}

func (t *wideNumbers) at(i uint32) float64 {
	return t.chunks[i>>12].Load()[i&(1<<12-1)]
}

// --- original names ---

// node_originalNameID returns the id of n's original name (0: none).
func (n Node) node_originalNameID() rhinoStringID {
	r := n.rec()
	if r.kindFlags&kindFlagOrigNameIsString != 0 {
		return rhinoStringID(r.payload)
	}
	return rhinoStringID(origNameField.get(n))
}

// node_setOriginalNameID sets n's original name (0: none).
func (n Node) node_setOriginalNameID(id rhinoStringID) {
	r := n.rec()
	if id != 0 && r.kindFlags&kindMask == uint8(node_kindString) && rhinoStringID(r.payload) == id {
		r.kindFlags |= kindFlagOrigNameIsString
		id = 0
	} else {
		r.kindFlags &^= kindFlagOrigNameIsString
	}
	if id != 0 || origNameField.get(n) != 0 {
		origNameField.set(n, &n.store().origNames, uint32(id))
	}
}

// setStringPayload sets the string of a string node; its original name
// stays what it was (if it was the node's string, it is now kept apart).
func (n Node) setStringPayload(id rhinoStringID) {
	r := n.rec()
	if r.kindFlags&kindFlagOrigNameIsString != 0 && rhinoStringID(r.payload) != id {
		orig := rhinoStringID(r.payload)
		r.payload = uint32(id)
		n.node_setOriginalNameID(orig)
		return
	}
	r.payload = uint32(id)
}

// --- template literal strings ---

// node_cookedID returns the id of the cooked string of a template literal
// string (0: none, an invalid escape in a tagged template).
func (n Node) node_cookedID() rhinoStringID {
	if box := n.extBoxOrNil(); box != nil {
		if box.cookedNil {
			return 0
		}
		if box.cookedID != 0 {
			return box.cookedID
		}
	}
	return rhinoStringID(n.rec().payload)
}

// initTemplateLitPayload sets the raw and cooked strings of a new template
// literal string (cookedID 0: no cooked string).
func (n Node) initTemplateLitPayload(cookedID, rawID rhinoStringID) {
	n.rec().payload = uint32(rawID)
	switch {
	case cookedID == 0:
		n.ensureExtBox().cookedNil = true
	case cookedID != rawID:
		n.ensureExtBox().cookedID = cookedID
	}
}

// --- property bits ---

// propBits are the 64 property bits of a node (see the propFlag and
// nodeProp constants), kept as two 32-bit words so that the record needs
// only 4-byte alignment. The methods compile to single 64-bit operations.
// --- property bits ---

// A node's property bits have a logical 64-bit layout, which the code reads
// and writes with Node.props and Node.setProps:
//
//	bits 0..39   the boolean properties (propBit_*), the common ones first
//	bit  40      propFlagHasJSDoc
//	bits 41..42  ConstantVarFlags
//	bits 43..46  SideEffectFlags
//	bits 47..62  the source-file id (nodePropSourceFile*)
//
// The record keeps 32 of them, in a propWord. The other boolean properties
// (bits propHotCount..39) are rare, so they live in the node's aux entry
// (nodeAux.cold), and the word only says whether there are any.
type propWord uint32

// The bits of a propWord:
//
//	bits 0..7    the first propHotCount boolean properties
//	bit  8       propFlagHasJSDoc
//	bits 9..14   ConstantVarFlags, SideEffectFlags
//	bit  15      propWordHasCold: nodeAux.cold is not zero
//	bits 16..31  the source-file id
const (
	propWordHotMask propWord = 1<<propHotCount - 1
	propWordJSDoc   propWord = 1 << 8
	propWordHasCold propWord = 1 << 15

	propLogicalMidShift = 40 // logical bits 40..46 are word bits 8..14
	propWordMidShift    = 8
	propMidMask         = 0x7F
	propWordSFShift     = 16
)

// props returns n's property bits in the logical layout.
func (n Node) props() uint64 {
	p := n.rec().propFlags
	v := uint64(p&propWordHotMask) |
		uint64(p>>propWordMidShift&propMidMask)<<propLogicalMidShift |
		uint64(p>>propWordSFShift)<<nodePropSourceFileShift
	if p&propWordHasCold != 0 {
		v |= uint64(n.coldProps()) << propHotCount
	}
	return v
}

// setProps sets n's property bits from the logical layout.
func (n Node) setProps(v uint64) {
	r := n.rec()
	cold := uint32(v>>propHotCount) & (1<<(propBitCount-propHotCount) - 1)
	p := propWord(v)&propWordHotMask |
		propWord(v>>propLogicalMidShift&propMidMask)<<propWordMidShift |
		propWord(v>>nodePropSourceFileShift)<<propWordSFShift
	if cold != 0 {
		p |= propWordHasCold
	}
	hadCold := r.propFlags&propWordHasCold != 0
	r.propFlags = p
	if cold != 0 || hadCold {
		n.setColdProps(cold)
	}
}

func (n Node) orProps(v uint64)     { n.setProps(n.props() | v) }
func (n Node) andProps(v uint64)    { n.setProps(n.props() & v) }
func (n Node) andNotProps(v uint64) { n.setProps(n.props() &^ v) }
