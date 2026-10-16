// Port of com/google/javascript/jscomp/ScopeCreator.java
// @ closure-compiler v20260415
//
// ScopeCreator is a small interface (48 lines in Java) that declares how
// objects capable of creating scopes from the parse tree behave. It contains
// one abstract method createScope(Node, AbstractScope) and one constant:
// ASSERT_NO_SCOPES_CREATED.
//
// Java interface: interface ScopeCreator
//
// Go port strategy:
//   UntypedScopeCreator is the interface that mirrors Java's ScopeCreator for
//   the untyped (Scope/Var) hierarchy. The TypedScopeCreator in scope.go
//   already has its own shape; this file provides the untyped counterpart.
//
//   AssertNoScopesCreated is a sentinel implementation that panics if called,
//   mirroring ScopeCreator.ASSERT_NO_SCOPES_CREATED in Java.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// UntypedScopeCreator is the Go port of com/google/javascript/jscomp/ScopeCreator.java.
// It is satisfied by any object that can create UntypedScopes from AST nodes.
//
// Mirrors: interface ScopeCreator in Java.
type UntypedScopeCreator interface {
	// CreateUntypedScope creates a new UntypedScope for the given AST node.
	// parent is nil for the global scope, non-nil for all child scopes.
	//
	// Mirrors ScopeCreator.createScope(Node, AbstractScope) in Java.
	CreateUntypedScope(n *rhino.Node, parent *UntypedScope) *UntypedScope
}

// assertNoScopesCreated is the Go port of ScopeCreator.ASSERT_NO_SCOPES_CREATED.
// It panics when CreateUntypedScope is called, signalling that the traversal
// must not create any scopes.
type assertNoScopesCreated struct{}

// AssertNoUntypedScopesCreated is the singleton sentinel that mirrors
// ScopeCreator.ASSERT_NO_SCOPES_CREATED in Java. Pass it to any traversal
// that must not create scopes.
var AssertNoUntypedScopesCreated UntypedScopeCreator = assertNoScopesCreated{}

// Go-only: language-adapter — CreateUntypedScope is the Go name for Java's
// createScope; the "Untyped" prefix distinguishes it from the typed scope
// creation path. Java uses a single createScope(Node, AbstractScope) signature.
func (assertNoScopesCreated) CreateUntypedScope(_ *rhino.Node, _ *UntypedScope) *UntypedScope {
	panic("ScopeCreator.ASSERT_NO_SCOPES_CREATED: no scopes may be created by this traversal")
}

// Go-only: language-adapter — SimpleScopeCreator is a lightweight Go-only scope
// creator used internally; Java uses SyntacticScopeCreator (a separate class)
// for the full scope-building logic. The split exists because Go can't implement
// the ScopeCreator interface with a stateful SyntacticScopeCreator directly.
type SimpleScopeCreator struct{}

// Go-only: language-adapter — CreateUntypedScope implements ScopeCreator.createScope
// for UntypedScope; the "Untyped" prefix is a Go-only distinction since Java's
// createScope returns AbstractScope<?> and the type is inferred at the call site.
func (c SimpleScopeCreator) CreateUntypedScope(n *rhino.Node, parent *UntypedScope) *UntypedScope {
	var scope *UntypedScope
	if parent == nil {
		scope = CreateGlobalUntypedScope(n)
	} else {
		scope = CreateChildUntypedScope(parent, n)
	}
	walkAndDeclare(n, scope)
	return scope
}

// Go-only: language-adapter — walkAndDeclare implements the AST traversal for
// SimpleScopeCreator; Java inlines this logic in SyntacticScopeCreator.declareLHS.
func walkAndDeclare(n *rhino.Node, scope *UntypedScope) {
	if n == nil {
		return
	}
	for child := n.FirstChild(); child != nil; child = child.Next() {
		switch child.Token() {
		case rhino.VAR, rhino.LET, rhino.CONST:
			declareNameDeclaration(child, scope)
		case rhino.FUNCTION:
			if nameNode := child.FirstChild(); nameNode != nil && nameNode.String_() != "" {
				scope.Declare(nameNode.String_(), nameNode)
			}
			// Do not recurse into function bodies.
		case rhino.CLASS:
			kids := child.Children()
			if len(kids) > 0 && kids[0].String_() != "" {
				scope.Declare(kids[0].String_(), kids[0])
			}
			// Do not recurse into class bodies.
		case rhino.BLOCK, rhino.IF, rhino.WHILE, rhino.DO_WHILE,
			rhino.FOR, rhino.FOR_IN, rhino.FOR_OF,
			rhino.SWITCH, rhino.CASE, rhino.DEFAULT_CASE,
			rhino.TRY, rhino.CATCH, rhino.FINALLY, rhino.LABEL:
			walkAndDeclare(child, scope)
		}
	}
}

// Go-only: language-adapter — declareNameDeclaration handles VAR/LET/CONST
// children; Java's SyntacticScopeCreator.declareLHS handles this inline.
func declareNameDeclaration(decl *rhino.Node, scope *UntypedScope) {
	for child := decl.FirstChild(); child != nil; child = child.Next() {
		if child.Token() == rhino.NAME && child.String_() != "" {
			scope.Declare(child.String_(), child)
		}
	}
}
