package ast

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Accessors that the grammar cannot describe on its own.

// Catch returns the catch clause of x; ok is false if x has none.
//
// (The tree keeps the catch clause in a BLOCK, the second child of the TRY,
// which is empty when there is no catch.)
func (x Try) Catch() (Catch, bool) {
	return x.catches().Catch()
}

// SetCatch sets the catch clause of x to v (which must not be in a tree).
func (x Try) SetCatch(v Catch) {
	x.catches().SetCatch(v)
}

// NewTry returns try body, without catch or finally yet (use SetCatch and
// SetFinally: a try needs at least one of them).
func NewTry(body Block) Try {
	n := rhino.NewNode(rhino.Token_TRY)
	n.AddChildToBack(body.node)
	n.AddChildToBack(rhino.NewNode(rhino.Token_BLOCK))
	return Try{n}
}

// NewFunction returns function name(params) body; an empty name makes an
// anonymous function.
func NewFunction(name string, params ParamList, body FunctionBody) Function {
	return Function{rhino.IR_Function(rhino.IR_Name(name), params.node, body.node)}
}

// Name returns the name of x; ok is false if x is anonymous.
//
// (An anonymous function has a NAME child with an empty string.)
func (x Function) Name() (Name, bool) {
	n := x.name()
	return n, n.Ident() != ""
}

// RemoveName makes x anonymous.
//
// (The NAME child stays, with an empty string.)
func (x Function) RemoveName() { x.name().SetIdent("") }

// IsArrow reports whether x is an arrow function.
func (x Function) IsArrow() bool { return x.node.IsArrowFunction() }

// SetArrow makes x an arrow function, or a normal one.
func (x Function) SetArrow(arrow bool) { x.node.SetIsArrowFunction(arrow) }

// Value returns the number.
func (x Number) Value() float64 { return x.node.GetDouble() }

// Body returns the body of the loop x.
func (x Loop) Body() Block {
	if x.node.IsDo() {
		return Block{x.node.GetFirstChild()}
	}
	return Block{x.node.GetLastChild()}
}

// NewName returns a new Name.
func NewName(ident string) Name { return Name{rhino.IR_Name(ident)} }

// NewStringLit returns a new StringLit.
func NewStringLit(value string) StringLit { return StringLit{rhino.IR_String(value)} }

// NewGetProp returns object.property.
func NewGetProp(object Expr, property string) GetProp {
	return GetProp{rhino.IR_Getprop(object.node, property)}
}

// NewLabelName returns a new LabelName.
func NewLabelName(ident string) LabelName { return LabelName{rhino.IR_LabelName(ident)} }

// Raw returns the raw text of x, as written in the source.
func (x TemplateLitString) Raw() string { return x.node.GetRawString() }

// Cooked returns the cooked text of x (escapes processed); ok is false if
// it has none (an invalid escape in a tagged template).
func (x TemplateLitString) Cooked() (s string, ok bool) {
	if p := x.node.GetCookedString(); p != nil {
		return *p, true
	}
	return "", false
}
