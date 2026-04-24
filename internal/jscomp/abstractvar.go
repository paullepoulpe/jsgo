// Port of com/google/javascript/jscomp/AbstractVar.java
// @ closure-compiler v20260415
//
// AbstractVar<S,V> is the generic base class for Var (used with Scope) and
// TypedVar (used with TypedScope). It extends ScopedName and carries the name,
// the name AST node, the enclosing scope, a declaration index, the compiler
// input, and optional implicit-goog-namespace metadata.
//
// Java inheritance chain:
//   ScopedName ← AbstractVar<S,V>
//                ← Var   (S=Scope,  V=Var)
//                ← TypedVar (S=TypedScope, V=TypedVar)
//
// Go port strategy:
//   AbstractVarBase is a struct that holds all fields declared in Java's
//   AbstractVar. Both UntypedVar (port of Var.java) and GoTypedVar (port of
//   TypedVar.java) embed AbstractVarBase.
//
//   The existing TypedScope in scope.go does not use AbstractVarBase; it uses
//   plain JSType maps. This file defines the hierarchy for callers that need
//   to work with the raw Var / TypedVar types (e.g. VarCheck, SymbolTable).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AbstractVarBase holds all fields from Java's AbstractVar<S,V>.
// Embed this in UntypedVar and GoTypedVar to mirror Java's inheritance.
type AbstractVarBase struct {
	name     string
	nameNode *rhino.Node
	index    int
	// implicitGoogNamespace: non-nil iff this is a goog.provide / goog.module namespace.
	implicitGoogNamespace *rhino.Node
}

// GetName returns the variable's identifier name.
// Mirrors AbstractVar.getName() in Java.
func (v *AbstractVarBase) GetName() string { return v.name }

// GetNode returns the NAME AST node that declares this variable, or nil for
// implicit variables.
// Mirrors AbstractVar.getNode() in Java.
func (v *AbstractVarBase) GetNode() *rhino.Node { return v.nameNode }

// GetNameNode is an alias for GetNode (matches Java's getNameNode()).
func (v *AbstractVarBase) GetNameNode() *rhino.Node { return v.nameNode }

// GetParentNode returns the parent of the NAME node, or nil if the name node
// is nil.
// Mirrors AbstractVar.getParentNode() in Java.
func (v *AbstractVarBase) GetParentNode() *rhino.Node {
	if v.nameNode == nil {
		return nil
	}
	return v.nameNode.Parent()
}

// GetIndex returns the declaration index within the scope.
// Mirrors AbstractVar.getIndex() in Java.
func (v *AbstractVarBase) GetIndex() int { return v.index }

// IsImplicitGoogNamespace reports whether this variable was synthesised from
// a goog.provide or legacy goog.module statement rather than an explicit
// declaration.
// Mirrors AbstractVar.isImplicitGoogNamespace() in Java.
func (v *AbstractVarBase) IsImplicitGoogNamespace() bool {
	return v.implicitGoogNamespace != nil
}

// IsBleedingFunction reports whether this is a named function expression whose
// name bleeds into the enclosing scope.
// Mirrors AbstractVar.isBleedingFunction() in Java.
func (v *AbstractVarBase) IsBleedingFunction() bool {
	parent := v.GetParentNode()
	if parent == nil {
		return false
	}
	return parent.Token() == rhino.FUNCTION
}

// declarationTypeOf returns the declaration token (VAR, LET, CONST, FUNCTION,
// CLASS, CATCH, IMPORT, PARAM_LIST) for a name node, walking up the ancestor
// chain. Returns -1 for implicit variables.
// Mirrors AbstractVar.declarationType() in Java.
func declarationTypeOf(nameNode *rhino.Node) rhino.Token {
	if nameNode == nil {
		return -1
	}
	for cur := nameNode; cur != nil; cur = cur.Parent() {
		switch cur.Token() {
		case rhino.VAR, rhino.LET, rhino.CONST,
			rhino.FUNCTION, rhino.CLASS, rhino.CATCH,
			rhino.IMPORT, rhino.PARAM_LIST:
			return cur.Token()
		}
	}
	return -1
}

// AbstractVarIface is the interface satisfied by both UntypedVar and GoTypedVar.
// It mirrors the StaticSlot / StaticRef contract from Java.
type AbstractVarIface interface {
	GetName() string
	GetNode() *rhino.Node
	GetParentNode() *rhino.Node
	GetIndex() int
	IsBleedingFunction() bool
	IsImplicitGoogNamespace() bool
}
