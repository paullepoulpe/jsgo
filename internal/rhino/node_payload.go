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
type propBits [2]uint32

func (p *propBits) get() uint64       { return uint64(p[0]) | uint64(p[1])<<32 }
func (p *propBits) set(v uint64)      { p[0], p[1] = uint32(v), uint32(v>>32) }
func (p *propBits) or(v uint64)       { p.set(p.get() | v) }
func (p *propBits) and(v uint64)      { p.set(p.get() & v) }
func (p *propBits) andNot(v uint64)   { p.set(p.get() &^ v) }
func (p *propBits) has(v uint64) bool { return p.get()&v != 0 }
