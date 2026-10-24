package rhino

// Nodes is the current store's records, passed by value: code that keeps a
// Nodes in a local or parameter reads a node's fields with one load (the
// slice's address and length are in registers), where a Node method first
// loads the current store's slice from a package variable.
//
// A Nodes is valid for one pass (the slice only moves between passes, see
// NodeStore.Settle). Nodes allocated during the pass beyond the slice, and
// nodes of the static store, are read through the slower path.
type Nodes struct {
	recs []nodeData
	base uint32
}

// CurrentNodes returns the current store's records.
func CurrentNodes() Nodes {
	return Nodes{recs: currentRecs, base: currentBase}
}

func (s Nodes) rec(n Node) *nodeData {
	if i := uint32(n) - s.base; i < uint32(len(s.recs)) {
		return &s.recs[i]
	}
	return n.rec()
}

// First is n.GetFirstChild().
func (s Nodes) First(n Node) Node { return Node(s.rec(n).first) }

// Next is n.GetNext().
func (s Nodes) Next(n Node) Node { return Node(s.rec(n).next) }

// Parent is n.GetParent().
func (s Nodes) Parent(n Node) Node { return Node(s.rec(n).parent) }

// Token is n.GetToken().
func (s Nodes) Token(n Node) Token { return s.rec(n).token }
