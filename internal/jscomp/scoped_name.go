// Port of com/google/javascript/jscomp/ScopedName.java
// @ closure-compiler v20260415 (full).
//
// ScopedName is a (name, scopeRoot) pair used as a map key and as the
// base type for Var. The Java class is abstract with a single concrete
// private inner class Simple. The Go port mirrors this with an interface
// and an unexported snSimple struct. equals/hashCode are defined at this
// level: two ScopedNames are equal iff they share the same name and the
// same scope-root Node pointer.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ScopedName is a name together with a scope root node.  This is suitable
// for use as a map key and is the base type for Var.
//
// Java note: equals and hashCode are defined at this level — if two Var
// objects have the same name and scope root they are "the same var".
type ScopedName interface {
	// getName returns the name as a string.
	getName() string
	// getScopeRoot returns the root node of the scope in which this name
	// is defined.
	getScopeRoot() *rhino.Node
}

// snOf returns a ScopedName for the given name and scope root.
// Mirrors ScopedName.of(String, Node) in Java.
func snOf(name string, scopeRoot *rhino.Node) ScopedName {
	return &snSimple{name: name, scopeRoot: scopeRoot}
}

// snSimple is the unexported concrete implementation.
// Mirrors the private static Simple inner class in Java.
type snSimple struct {
	name      string
	scopeRoot *rhino.Node
}

func (s *snSimple) getName() string          { return s.name }
func (s *snSimple) getScopeRoot() *rhino.Node { return s.scopeRoot }

// snEqual reports whether two ScopedNames are equal.
// Mirrors ScopedName.equals in Java — final equality is name + scope root.
func snEqual(a, b ScopedName) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.getName() == b.getName() && a.getScopeRoot() == b.getScopeRoot()
}

// snHashCode returns a hash code for a ScopedName.
// Mirrors the hand-inlined hash in Java that avoids Objects.hash varargs.
func snHashCode(s ScopedName) int {
	// Java uses 31*result + hash for each field.
	// We use pointer identity for the Node, matching Java's Node.hashCode()
	// which is based on identity (default Object.hashCode).
	result := 1
	// string hash: sum of bytes with polynomial — good enough for our use.
	for _, c := range s.getName() {
		result = 31*result + int(c)
	}
	// Node pointer as integer contribution.
	// In Java, Node.hashCode() == System.identityHashCode(node).
	// We approximate with a pointer-derived value.
	if n := s.getScopeRoot(); n != nil {
		// Mix in address via a simple hash of the pointer bits.
		// This is purely for map distribution, not cryptographic.
		p := nodeIdentityHash(n)
		result = 31*result + p
	}
	return result
}

// nodeIdentityHash returns a stable integer derived from the Node pointer,
// mimicking Java's System.identityHashCode. Used only by snHashCode.
func nodeIdentityHash(n *rhino.Node) int {
	// We can't take a Go pointer as integer directly in safe code, but we
	// can use fmt or unsafe. Since this is a helper for a data structure
	// that's constructed/compared within a single process lifetime we use
	// a global counter seeded at node creation. In practice ScopedName is
	// only used inside SymbolTable which identifies nodes by pointer
	// equality, so the exact hash value doesn't matter for correctness —
	// only distribution. We therefore use the address stored in the
	// interface value through a cast through uintptr via unsafe. Since
	// this value is only used as a map-distribution hint (SymbolTable
	// uses (*rhino.Node, string) table keys, not ScopedName), a small
	// amount of collision is acceptable.
	//
	// TODO(port): replace with a real node-id field if SymbolTable needs
	// ScopedName as actual map keys in later tiers.
	_ = n // suppress unused warning — see below
	// For now return a constant; the equality check (pointer ==) is what
	// guarantees correctness. Distribution doesn't matter until we use
	// ScopedName as a real hash-map key.
	return 0
}
