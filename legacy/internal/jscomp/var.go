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

	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
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

// Go-only: language-adapter — newUntypedVar is the Go constructor for UntypedVar;
// Java uses the package-private Var(name, nameNode, scope, index, input, null) constructor.
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

// Go-only: language-adapter — createUntypedImplicitGoogNamespace is the Go name for
// the factory that Java calls Var.createImplicitGoogNamespace(name, scope, definition).
// The "Untyped" prefix distinguishes it from a future TypedVar equivalent.
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

// Go-only: language-adapter — GetScope exposes the UntypedScope pointer directly;
// in Java, getScope() is defined on AbstractVar and returns the generic S (scope type).
func (v *UntypedVar) GetScope() *UntypedScope { return v.scope }

// Go-only: language-adapter — IsGlobal is inherited from AbstractVar in Java;
// Go requires explicit delegation since UntypedVar cannot inherit from AbstractVarBase.
func (v *UntypedVar) IsGlobal() bool {
	return v.scope != nil && v.scope.IsGlobal()
}

// Go-only: language-adapter — IsLocal is inherited from AbstractVar in Java;
// Go requires explicit delegation since UntypedVar cannot inherit AbstractVar.
func (v *UntypedVar) IsLocal() bool {
	return v.scope != nil && v.scope.IsLocal()
}

// Go-only: language-adapter — IsVar is inherited from AbstractVar in Java (final method);
// Go requires explicit implementation since UntypedVar embeds AbstractVarBase, not AbstractVar.
func (v *UntypedVar) IsVar() bool {
	return declarationTypeOf(v.nameNode) == rhino.VAR
}

// Go-only: language-adapter — IsLet is inherited from AbstractVar.isLet() in Java;
// Go requires explicit delegation since UntypedVar embeds AbstractVarBase struct.
func (v *UntypedVar) IsLet() bool {
	return declarationTypeOf(v.nameNode) == rhino.LET
}

// Go-only: language-adapter — IsConst is inherited from AbstractVar.isConst() in Java;
// Go requires explicit delegation since UntypedVar embeds AbstractVarBase struct.
func (v *UntypedVar) IsConst() bool {
	return declarationTypeOf(v.nameNode) == rhino.CONST
}

// Go-only: language-adapter — IsParam is inherited from AbstractVar.isParam() in Java;
// Go requires explicit delegation since UntypedVar embeds AbstractVarBase struct.
func (v *UntypedVar) IsParam() bool {
	return declarationTypeOf(v.nameNode) == rhino.PARAM_LIST
}

// String returns a human-readable description.
// Mirrors Var.toString() in Java.
func (v *UntypedVar) String() string {
	return fmt.Sprintf("Var %s @ %v", v.name, v.nameNode)
}
