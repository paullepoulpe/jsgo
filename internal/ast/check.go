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
	role string
	typ  string // a kind, a category, or Node
	mode childMode
}

type kindSpec struct {
	name     string
	context  bool
	tokens   []rhino.Token
	children []childSpec
}

var (
	specsOnce   sync.Once
	specByName  map[string]*kindSpec
	specByToken map[rhino.Token]*kindSpec            // kinds recognized by their token
	catByToken  map[string]map[rhino.Token]*kindSpec // category -> token -> kind
)

func initSpecs() {
	specsOnce.Do(func() {
		specByName = map[string]*kindSpec{}
		specByToken = map[rhino.Token]*kindSpec{}
		for i := range grammarKinds {
			k := &grammarKinds[i]
			specByName[k.name] = k
			if !k.context {
				for _, t := range k.tokens {
					specByToken[t] = k
				}
			}
		}
		catByToken = map[string]map[rhino.Token]*kindSpec{}
		for c, kinds := range grammarCategories {
			m := map[rhino.Token]*kindSpec{}
			for _, kn := range kinds {
				for _, t := range specByName[kn].tokens {
					m[t] = specByName[kn]
				}
			}
			catByToken[c] = m
		}
	})
}

// kindFor returns the kind n has as a child of type typ (nil: n is not a
// typ; typ Node gives nil too, and the caller does not look further).
func kindFor(n rhino.Node, typ string) (k *kindSpec, ok bool) {
	if typ == "Node" {
		return nil, true
	}
	if m, isCat := catByToken[typ]; isCat {
		k = m[n.GetToken()]
		return k, k != nil
	}
	k = specByName[typ]
	for _, t := range k.tokens {
		if t == n.GetToken() {
			return k, true
		}
	}
	return nil, false
}

func checkChild(c rhino.Node, s childSpec, parent *kindSpec, report func(rhino.Node, string)) {
	k, ok := kindFor(c, s.typ)
	if !ok {
		report(c, fmt.Sprintf("%s.%s is a %s, want %s", parent.name, s.role, c.GetToken(), s.typ))
		return
	}
	if k != nil {
		checkKind(c, k, report)
	}
}

func checkKind(n rhino.Node, k *kindSpec, report func(rhino.Node, string)) {
	c := n.GetFirstChild()
	for _, s := range k.children {
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
