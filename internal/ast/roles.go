package ast

import "github.com/paullepoulpe/jsgo/internal/rhino"

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
