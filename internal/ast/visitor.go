package ast

import "github.com/paullepoulpe/jsgo/internal/rhino"

// TokenSet is a set of rhino tokens.
type TokenSet [4]uint64

// Add adds t to s.
func (s *TokenSet) Add(t rhino.Token) { s[t>>6] |= 1 << (t & 63) }

// Has reports whether t is in s.
func (s *TokenSet) Has(t rhino.Token) bool { return s[t>>6]&(1<<(t&63)) != 0 }

// mustError is the panic value of a Must conversion of node to kind. Its
// message is built only when printed, so the conversions stay small enough
// to inline.
type mustError struct {
	kind string
	node rhino.Node
}

func (e mustError) Error() string {
	if e.node == rhino.NilNode {
		return "ast.Must" + e.kind + ": the nil node"
	}
	return "ast.Must" + e.kind + ": a " + e.node.GetToken().String() + " node"
}

// tokenSetOf returns the set of ts.
func tokenSetOf(ts ...rhino.Token) TokenSet {
	var s TokenSet
	for _, t := range ts {
		s.Add(t)
	}
	return s
}
