// Port of com/google/javascript/jscomp/TypedVar.java
// @ closure-compiler v20260415
//
// TypedVar extends AbstractVar<TypedScope, TypedVar> and adds a JSType slot,
// an "inferred" flag, and two TypeInference flags (markedEscaped /
// markedAssignedExactlyOnce).
//
// Java class: public class TypedVar extends AbstractVar<TypedScope, TypedVar>
//             implements StaticTypedSlot, StaticTypedRef
//
// Note: the existing TypedScope in scope.go uses plain JSType maps rather
// than TypedVar instances. GoTypedVar is provided for callers (e.g. VarCheck,
// SymbolTable) that need to work at the TypedVar level.
//
// equals/hashCode: inherited from ScopedName in Java — two TypedVars with the
// same name and scope root are equal, regardless of type information.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino/jstype"
)

// GoTypedVar is the Go port of com/google/javascript/jscomp/TypedVar.java.
// It holds a typed variable declared in a TypedScope.
//
// The name prefix "Go" avoids collision with the existing TypedScope-internal
// representation; a future refactor may merge the two.
type GoTypedVar struct {
	AbstractVarBase

	// scope is the enclosing TypedScope.
	scope *TypedScope

	// jsType is the declared or inferred JSType of this variable.
	jsType jstype.JSType

	// typeInferred is true when jsType was inferred rather than declared.
	// Mirrors TypedVar.typeInferred in Java.
	typeInferred bool

	// markedEscaped is set by TypeInference when the variable escapes.
	// Mirrors TypedVar.markedEscaped in Java.
	markedEscaped bool

	// markedAssignedExactlyOnce is set by TypeInference when the assignment
	// pattern is detected. Mirrors TypedVar.markedAssignedExactlyOnce in Java.
	markedAssignedExactlyOnce bool
}

// newGoTypedVar creates a GoTypedVar.
// Mirrors the package-private TypedVar constructor in Java.
func newGoTypedVar(
	inferred bool,
	name string,
	nameNode *rhino.Node,
	t jstype.JSType,
	scope *TypedScope,
	index int,
) *GoTypedVar {
	return &GoTypedVar{
		AbstractVarBase: AbstractVarBase{
			name:     name,
			nameNode: nameNode,
			index:    index,
		},
		scope:        scope,
		jsType:       t,
		typeInferred: inferred,
	}
}

// GetScope returns the enclosing TypedScope.
// Mirrors AbstractVar.getScope() in Java.
func (v *GoTypedVar) GetScope() *TypedScope { return v.scope }

// GetType returns the JSType of this variable.
// Mirrors TypedVar.getType() / StaticTypedSlot.getType() in Java.
func (v *GoTypedVar) GetType() jstype.JSType { return v.jsType }

// SetType updates the JSType.
// Mirrors TypedVar.setType(JSType) in Java.
func (v *GoTypedVar) SetType(t jstype.JSType) { v.jsType = t }

// IsTypeInferred reports whether the type was inferred.
// Mirrors TypedVar.isTypeInferred() in Java.
func (v *GoTypedVar) IsTypeInferred() bool { return v.typeInferred }

// MarkEscaped marks this variable as having escaped its declaring scope.
// Mirrors TypedVar.markEscaped() in Java.
func (v *GoTypedVar) MarkEscaped() { v.markedEscaped = true }

// IsMarkedEscaped reports whether MarkEscaped has been called.
// Mirrors TypedVar.isMarkedEscaped() in Java.
func (v *GoTypedVar) IsMarkedEscaped() bool { return v.markedEscaped }

// MarkAssignedExactlyOnce marks this variable as assigned exactly once.
// Mirrors TypedVar.markAssignedExactlyOnce() in Java.
func (v *GoTypedVar) MarkAssignedExactlyOnce() { v.markedAssignedExactlyOnce = true }

// IsMarkedAssignedExactlyOnce reports whether MarkAssignedExactlyOnce has
// been called. Mirrors TypedVar.isMarkedAssignedExactlyOnce() in Java.
func (v *GoTypedVar) IsMarkedAssignedExactlyOnce() bool { return v.markedAssignedExactlyOnce }

// GetInputName returns the name of the compiler input that contains this
// variable's declaration. Returns "<non-file>" when no input is available.
// Mirrors TypedVar.getInputName() in Java.
func (v *GoTypedVar) GetInputName() string {
	// CompilerInput is not yet ported; return a placeholder.
	return "<non-file>"
}

// String returns a human-readable description.
// Mirrors TypedVar.toString() in Java.
func (v *GoTypedVar) String() string {
	return fmt.Sprintf("Var %s{%v}", v.name, v.jsType)
}
