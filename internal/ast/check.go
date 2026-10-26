package ast

import (
	"fmt"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Check walks the tree under n and calls report for every node that does
// not follow the grammar (grammar.txt): a child of a kind its role does
// not allow, a missing child, or an extra one. Nodes whose token has no
// kind in the grammar (JSDoc type expressions, for example) are not
// checked, nor is what is under a role of type Node.
func Check(n rhino.Node, report func(n rhino.Node, problem string)) {
	initSpecs()
	if k := specByToken[n.GetToken()]; k != nil {
		checkKind(n, k, report)
	}
}

type childMode int

const (
	required childMode = iota
	optional
	emptyOr
	list
)

type childSpec struct {
	role    string
	typ     string // a kind, a category, or Node
	mode    childMode
	allowed *allowedKinds // typ, compiled by initSpecs (nil for Node)
}

type kindSpec struct {
	name     string
	context  bool
	tokens   []rhino.Token
	children []childSpec
}

// allowedKinds gives, for each token, the kind a node with that token has
// in a place of some type (nil: not allowed there).
type allowedKinds [256]*kindSpec

var (
	specsOnce   sync.Once
	specByToken allowedKinds // kinds recognized by their token alone
)

func initSpecs() {
	specsOnce.Do(func() {
		byName := map[string]*kindSpec{}
		for i := range grammarKinds {
			k := &grammarKinds[i]
			byName[k.name] = k
			if !k.context {
				for _, t := range k.tokens {
					specByToken[t] = k
				}
			}
		}
		types := map[string]*allowedKinds{}
		allow := func(a *allowedKinds, k *kindSpec) {
			for _, t := range k.tokens {
				a[t] = k
			}
		}
		// order-independent: builds a lookup table keyed by name
		for name, k := range byName {
			a := &allowedKinds{}
			allow(a, k)
			types[name] = a
		}
		// order-independent: builds a lookup table keyed by category
		for c, kinds := range grammarCategories {
			a := &allowedKinds{}
			for _, kn := range kinds {
				allow(a, byName[kn])
			}
			types[c] = a
		}
		for i := range grammarKinds {
			for j := range grammarKinds[i].children {
				ch := &grammarKinds[i].children[j]
				ch.allowed = types[ch.typ] // nil for Node
			}
		}
	})
}

func checkChild(c rhino.Node, s *childSpec, parent *kindSpec, report func(rhino.Node, string)) {
	if s.allowed == nil {
		return // Node: anything, not checked further
	}
	k := s.allowed[c.GetToken()]
	if k == nil {
		report(c, fmt.Sprintf("%s.%s is a %s, want %s", parent.name, s.role, c.GetToken(), s.typ))
		return
	}
	checkKind(c, k, report)
}

func checkKind(n rhino.Node, k *kindSpec, report func(rhino.Node, string)) {
	c := n.GetFirstChild()
	for i := range k.children {
		s := &k.children[i]
		switch s.mode {
		case required:
			if c == rhino.NilNode {
				report(n, fmt.Sprintf("%s has no %s", k.name, s.role))
				return
			}
			checkChild(c, s, k, report)
		case emptyOr:
			if c == rhino.NilNode {
				report(n, fmt.Sprintf("%s has no %s (not even an EMPTY node)", k.name, s.role))
				return
			}
			if !c.IsEmpty() {
				checkChild(c, s, k, report)
			}
		case optional:
			if c == rhino.NilNode {
				return
			}
			checkChild(c, s, k, report)
		case list:
			for ; c != rhino.NilNode; c = c.GetNext() {
				checkChild(c, s, k, report)
			}
			return
		}
		c = c.GetNext()
	}
	if c != rhino.NilNode {
		report(c, fmt.Sprintf("%s has an extra child %s", k.name, c.GetToken()))
	}
}
