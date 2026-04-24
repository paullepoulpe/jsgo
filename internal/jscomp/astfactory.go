// Port of com/google/javascript/jscomp/AstFactory.java
// @ closure-compiler v20260415 (subset; node-creation factories).
//
// AstFactory creates typed AST nodes. Three modes: JSTYPE, COLOR, NONE.
// Source-reference information is never set here — caller responsibility.
// Type-system-heavy helpers (createQName, etc.) are stubs pending
// LifeCycleStage / RuntimeJsLibManager ports.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// astTypeMode mirrors AstFactory.TypeMode in Java.
type astTypeMode int

const (
	astTypeModeNone   astTypeMode = iota
	astTypeModeJSType             // full JSType (JSTypeRegistry)
	astTypeModeColor              // optimisation-only Color (TODO)
)

// AstFactory creates AST nodes optionally annotated with type information.
// Mirrors com.google.javascript.jscomp.AstFactory. Use factory functions.
type AstFactory struct {
	registry    *jstype.Registry
	unknownType jstype.JSType
	typeMode    astTypeMode
}

// NewAstFactoryWithoutTypes mirrors AstFactory.createFactoryWithoutTypes.
func NewAstFactoryWithoutTypes() *AstFactory {
	return &AstFactory{typeMode: astTypeModeNone}
}

// NewAstFactoryWithTypes mirrors AstFactory.createFactoryWithTypes.
func NewAstFactoryWithTypes(reg *jstype.Registry) *AstFactory {
	return &AstFactory{registry: reg, unknownType: jstype.Unknown, typeMode: astTypeModeJSType}
}

// IsAddingTypes mirrors AstFactory.isAddingTypes.
func (af *AstFactory) IsAddingTypes() bool { return af.typeMode == astTypeModeJSType }

// IsAddingColors mirrors AstFactory.isAddingColors.
func (af *AstFactory) IsAddingColors() bool { return af.typeMode == astTypeModeColor }

// setType attaches type to n in JSTYPE mode. Color mode is a TODO.
func (af *AstFactory) setType(n *rhino.Node, t jstype.JSType) {
	if af.typeMode == astTypeModeJSType && t != nil {
		n.SetJSType(t)
	}
}

// nodeJSType extracts the JSType stored on n, or nil if absent.
func nodeJSType(n *rhino.Node) jstype.JSType {
	t, _ := n.GetJSType().(jstype.JSType)
	return t
}

// CreateEmpty mirrors AstFactory.createEmpty. Uses ERROR as the zero
// token since rhino does not yet define EMPTY.
func (af *AstFactory) CreateEmpty() *rhino.Node { return rhino.New(rhino.ERROR) }

// CreateBlock mirrors AstFactory.createBlock.
func (af *AstFactory) CreateBlock(stmts ...*rhino.Node) *rhino.Node {
	n := rhino.New(rhino.BLOCK)
	for _, s := range stmts {
		n.AddChildToBack(s)
	}
	return n
}

// CreateIf mirrors AstFactory.createIf(cond, then).
func (af *AstFactory) CreateIf(cond, then *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.IF)
	n.AddChildToBack(cond)
	n.AddChildToBack(then)
	return n
}

// CreateIfElse mirrors AstFactory.createIf(cond, then, elseNode).
func (af *AstFactory) CreateIfElse(cond, then, elseNode *rhino.Node) *rhino.Node {
	n := af.CreateIf(cond, then)
	n.AddChildToBack(elseNode)
	return n
}

// CreateReturn mirrors AstFactory.createReturn.
func (af *AstFactory) CreateReturn(value *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.RETURN)
	if value != nil {
		n.AddChildToBack(value)
	}
	return n
}

// CreateThrow mirrors AstFactory.createThrow.
func (af *AstFactory) CreateThrow(expr *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.THROW)
	n.AddChildToBack(expr)
	return n
}

// CreateExprResult mirrors AstFactory.exprResult; sets srcref.
func (af *AstFactory) CreateExprResult(expr *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.EXPR_RESULT)
	n.AddChildToBack(expr)
	n.SrcrefIfMissing(expr)
	return n
}

// CreateString mirrors AstFactory.createString.
func (af *AstFactory) CreateString(value string) *rhino.Node {
	n := rhino.NewString(rhino.STRINGLIT, value)
	af.setType(n, jstype.String)
	return n
}

// CreateNumber mirrors AstFactory.createNumber.
func (af *AstFactory) CreateNumber(value float64) *rhino.Node {
	n := rhino.NewNumber(value)
	af.setType(n, jstype.Number)
	return n
}

// CreateBoolean mirrors AstFactory.createBoolean.
func (af *AstFactory) CreateBoolean(value bool) *rhino.Node {
	tok := rhino.TRUE
	if !value {
		tok = rhino.FALSE
	}
	n := rhino.New(tok)
	af.setType(n, jstype.Boolean)
	return n
}

// CreateNull mirrors AstFactory.createNull.
func (af *AstFactory) CreateNull() *rhino.Node {
	n := rhino.New(rhino.NULL)
	af.setType(n, jstype.Null)
	return n
}

// CreateThis mirrors AstFactory.createThis.
func (af *AstFactory) CreateThis(thisType jstype.JSType) *rhino.Node {
	n := rhino.New(rhino.THIS)
	af.setType(n, thisType)
	return n
}

// CreateAssign mirrors AstFactory.createAssign(Node, Node).
func (af *AstFactory) CreateAssign(lhs, rhs *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.ASSIGN)
	n.AddChildToBack(lhs)
	n.AddChildToBack(rhs)
	af.setType(n, nodeJSType(rhs))
	return n
}

// CreateName mirrors AstFactory.createName(String, Type).
func (af *AstFactory) CreateName(name string, t jstype.JSType) *rhino.Node {
	n := rhino.NewString(rhino.NAME, name)
	af.setType(n, t)
	return n
}

// CreateCall mirrors AstFactory.createCall.
func (af *AstFactory) CreateCall(callee *rhino.Node, resultType jstype.JSType, args ...*rhino.Node) *rhino.Node {
	n := rhino.New(rhino.CALL)
	n.AddChildToBack(callee)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	af.setType(n, resultType)
	return n
}

// CreateGetProp mirrors AstFactory.createGetProp(Node, String, Type).
func (af *AstFactory) CreateGetProp(receiver *rhino.Node, propertyName string, t jstype.JSType) *rhino.Node {
	n := rhino.NewString(rhino.GETPROP, propertyName)
	n.AddChildToFront(receiver)
	af.setType(n, t)
	return n
}
