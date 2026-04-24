// Port of com/google/javascript/jscomp/Var.java
// @ closure-compiler v20260415
//
// Var is the concrete variable type used by the untyped Scope. It extends
// AbstractVar<Scope, Var> in Java.
//
// Java class: public class Var extends AbstractVar<Scope, Var>
//
// Go port: UntypedVar embeds AbstractVarBase (which holds all AbstractVar
// fields) and adds a pointer back to the enclosing UntypedScope. The name
// "UntypedVar" avoids a collision with any future use of "Var" as a short
// alias in callers.
//
// Constants ARGUMENTS and EXPORTS mirror the Java static fields.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	// VarArguments is the special "arguments" implicit variable name.
	// Mirrors Var.ARGUMENTS in Java.
	VarArguments = "arguments"

	// VarExports is the special "exports" implicit variable name used in
	// goog.module scopes.
	// Mirrors Var.EXPORTS in Java.
	VarExports = "exports"
)

// UntypedVar is the Go port of com/google/javascript/jscomp/Var.java.
// It holds a variable declared in an UntypedScope.
type UntypedVar struct {
	AbstractVarBase
	scope *UntypedScope
}

// newUntypedVar creates a variable declared at nameNode in the given scope.
// Mirrors the package-private Var constructor in Java.
func newUntypedVar(name string, nameNode *rhino.Node, scope *UntypedScope, index int) *UntypedVar {
	return &UntypedVar{
		AbstractVarBase: AbstractVarBase{
			name:     name,
			nameNode: nameNode,
			index:    index,
		},
		scope: scope,
	}
}

// createImplicitGoogNamespace creates a Var that represents an implicit
// goog.provide or legacy goog.module namespace (no real AST nameNode).
// Mirrors Var.createImplicitGoogNamespace in Java.
func createUntypedImplicitGoogNamespace(name string, scope *UntypedScope, definition *rhino.Node) *UntypedVar {
	return &UntypedVar{
		AbstractVarBase: AbstractVarBase{
			name:                  name,
			nameNode:              nil,
			index:                 -1,
			implicitGoogNamespace: definition,
		},
		scope: scope,
	}
}

// GetScope returns the UntypedScope that contains this variable.
// Mirrors AbstractVar.getScope() in Java.
func (v *UntypedVar) GetScope() *UntypedScope { return v.scope }

// IsGlobal reports whether the enclosing scope is the global scope.
// Mirrors AbstractVar.isGlobal() in Java.
func (v *UntypedVar) IsGlobal() bool {
	return v.scope != nil && v.scope.IsGlobal()
}

// IsLocal reports whether the enclosing scope is a local (non-global) scope.
// Mirrors AbstractVar.isLocal() in Java.
func (v *UntypedVar) IsLocal() bool {
	return v.scope != nil && v.scope.IsLocal()
}

// IsVar reports whether this variable was declared with the `var` keyword.
// Mirrors AbstractVar.isVar() in Java.
func (v *UntypedVar) IsVar() bool {
	return declarationTypeOf(v.nameNode) == rhino.VAR
}

// IsLet reports whether this variable was declared with `let`.
func (v *UntypedVar) IsLet() bool {
	return declarationTypeOf(v.nameNode) == rhino.LET
}

// IsConst reports whether this variable was declared with `const`.
func (v *UntypedVar) IsConst() bool {
	return declarationTypeOf(v.nameNode) == rhino.CONST
}

// IsParam reports whether this variable is a function parameter.
func (v *UntypedVar) IsParam() bool {
	return declarationTypeOf(v.nameNode) == rhino.PARAM_LIST
}

// String returns a human-readable description.
// Mirrors Var.toString() in Java.
func (v *UntypedVar) String() string {
	return fmt.Sprintf("Var %s @ %v", v.name, v.nameNode)
}
