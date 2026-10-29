// Package ast is a typed view of the JavaScript AST.
//
// The compiler stores the AST as rhino.Nodes: every node has a token, a
// list of children and generic accessors (GetFirstChild, GetSecondChild,
// GetString...), so what a child means is a matter of convention: the
// second child of an IF is its then-block. This package gives each node
// kind a Go type whose methods name those children:
//
//	ifs, ok := ast.AsIf(n)  // ok is false if n is not an if statement
//	cond := ifs.Cond()      // an ast.Expr
//	if els, ok := ifs.Else(); ok {
//		...                 // els is an ast.Block
//	}
//
// A category (Expr, Statement, Target...) is the set of kinds that can
// appear in one place; As methods convert between kinds and categories,
// checking the token when the conversion can fail.
//
// RoleOf answers the question from the child's side: which child of its
// parent a node is (ast.RoleOf(n) == ast.RoleAssignTarget).
//
// The types are views: each holds the rhino.Node it describes (4 bytes),
// so they cost nothing to create or pass, and Node() goes back to the
// rhino API for anything this package does not cover.
//
// Accessors and setters find a child by its position, as the rhino API
// does: SetCond replaces the child where the cond of x is. In the middle of
// an edit that removed an earlier child, they reach the wrong one, so read
// every child you need before the first edit:
//
//	then, els := ifs.Then(), ifs.Else() // before the cond is detached
//
// The kinds, their children and the categories are listed in grammar.txt;
// nodes.gen.go is generated from it (go generate ./internal/ast). Check
// verifies a tree against the grammar.
package ast

//go:generate go run ../../tools/astgen -grammar grammar.txt -out nodes.gen.go

import (
	"iter"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReadAhead iterates over a list from its first element, as returned by a
// First<X> accessor:
//
//	for member := range ast.ReadAhead(obj.FirstMember()) {
//		member.Detach() // fine: the next member was already read
//	}
//
// It reads the next element before the loop body runs, so the body may
// detach, move or replace the current element. The list iterators
// (obj.Members(), block.Statements(), ...) read it after the body, like a
// for loop over GetNext; use them when the body changes the elements after
// the current one instead.
func ReadAhead[T interface{ Next() (T, bool) }](first T, ok bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		x, more := first, ok
		for more {
			next, nextOk := x.Next()
			if !yield(x) {
				return
			}
			x, more = next, nextOk
		}
	}
}

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

// listMustError is the panic value of a MustFirst or MustOnly list accessor.
// Its message is built only when printed, so the accessors stay small enough
// to inline.
type listMustError struct {
	method string
	node   rhino.Node
}

func (e listMustError) Error() string {
	if e.node == rhino.NilNode {
		return "ast." + e.method + ": the nil node"
	}
	return "ast." + e.method + ": " + e.node.GetToken().String() + " has " + strconv.Itoa(e.node.GetChildCount()) + " children"
}

// tokenSetOf returns the set of ts.
func tokenSetOf(ts ...rhino.Token) TokenSet {
	var s TokenSet
	for _, t := range ts {
		s.Add(t)
	}
	return s
}

// Role is the place of a node in its parent: which child of which kind it
// is. There is one Role per child of each kind of the grammar
// (RoleAssignTarget, RoleIfThen, RoleCallArgs...), and RoleOf finds a
// node's:
//
//	if ast.RoleOf(n) == ast.RoleAssignTarget { ... } // n is assigned to
//
// instead of reading the parent as a kind and comparing one of its
// children with n. A placeholder EMPTY node has the role it stands for
// (RoleForCond for the EMPTY of for (;;)).
type Role uint16

// roleRow gives the roles of the children of a node of one token: the
// fixed children in order, then the list for every child after them.
type roleRow struct {
	fixed   [maxFixedRoles]Role
	nfixed  uint8
	context bool // the token is also a context kind's (TryCatches): ask contextRole
	list    Role
}

// maxFixedRoles is the most children before a list (or in all) of a kind.
const maxFixedRoles = 4

// RoleOf returns the role of n in its parent; RoleNone if n has no parent
// or its parent's kind is not in the grammar.
func RoleOf(n rhino.Node) Role {
	return RoleIn(n, n.GetParent())
}

// RoleIn returns the role of n in its parent p (where p == n.GetParent());
// RoleNone if p is NilNode or p's kind is not in the grammar. Prefer RoleIn
// over RoleOf when the caller already holds p.
func RoleIn(n, p rhino.Node) Role {
	if p == rhino.NilNode {
		return RoleNone
	}
	row := &roleRows[p.GetToken()]
	if row.context {
		if r := contextRole(p); r != RoleNone {
			return r
		}
	}
	if row.nfixed == 0 {
		return row.list
	}
	c := p.GetFirstChild()
	for i := 0; i < int(row.nfixed); i++ {
		if c == n {
			return row.fixed[i]
		}
		c = c.GetNext()
	}
	return row.list
}

// Kind returns the name of the parent's kind of r ("Assign" for
// RoleAssignTarget), or "" for RoleNone.
func (r Role) Kind() string { return roleNames[r][0] }

// Name returns the role name of r, as the grammar writes it ("target" for
// RoleAssignTarget), or "" for RoleNone. Roles of different kinds can share
// a name: RoleOf(n).Name() == "body" asks whether n is the body of its
// parent, whatever that is.
func (r Role) Name() string { return roleNames[r][1] }

// String returns r as Kind.name ("Assign.target"), or "none".
func (r Role) String() string {
	if r == RoleNone {
		return "none"
	}
	return r.Kind() + "." + r.Name()
}
