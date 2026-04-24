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

// createGetPropWithUnknownType mirrors AstFactory.createGetPropWithUnknownType.
func (af *AstFactory) createGetPropWithUnknownType(receiver *rhino.Node, propertyName string) *rhino.Node {
	return af.CreateGetProp(receiver, propertyName, af.unknownType)
}

// createGetPropWithoutColor mirrors AstFactory.createGetPropWithoutColor.
func (af *AstFactory) createGetPropWithoutColor(receiver *rhino.Node, propertyName string) *rhino.Node {
	return af.CreateGetProp(receiver, propertyName, nil)
}

// createGetElem mirrors AstFactory.createGetElem.
func (af *AstFactory) createGetElem(receiver, key *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.GETELEM)
	n.AddChildToBack(receiver)
	n.AddChildToBack(key)
	af.setType(n, af.unknownType)
	return n
}

// CreateFor mirrors AstFactory.createFor(Node, Node, Node, Node).
func (af *AstFactory) CreateFor(init, cond, incr, body *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.FOR)
	n.AddChildToBack(init)
	n.AddChildToBack(cond)
	n.AddChildToBack(incr)
	n.AddChildToBack(body)
	return n
}

// createBreak mirrors AstFactory.createBreak.
func (af *AstFactory) createBreak() *rhino.Node {
	return rhino.New(rhino.BREAK)
}

// createCatch mirrors AstFactory.createCatch.
func (af *AstFactory) createCatch(error, block *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.CATCH)
	n.AddChildToBack(error)
	n.AddChildToBack(block)
	return n
}

// createTryFinally mirrors AstFactory.createTryFinally.
func (af *AstFactory) createTryFinally(tryBlock, finallyBlock *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.TRY)
	n.AddChildToBack(tryBlock)
	n.AddChildToBack(rhino.New(rhino.BLOCK)) // empty catch list
	n.AddChildToBack(finallyBlock)
	return n
}

// createTryCatchFinally mirrors AstFactory.createTryCatchFinally.
func (af *AstFactory) createTryCatchFinally(tryBlock, catchNode, finallyBlock *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.TRY)
	n.AddChildToBack(tryBlock)
	catches := rhino.New(rhino.BLOCK)
	catches.AddChildToBack(catchNode)
	n.AddChildToBack(catches)
	n.AddChildToBack(finallyBlock)
	return n
}

// createVoid mirrors AstFactory.createVoid.
func (af *AstFactory) createVoid(child *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.VOID)
	n.AddChildToBack(child)
	af.setType(n, jstype.Void)
	return n
}

// createUndefinedValue mirrors AstFactory.createUndefinedValue.
// Java: `return createVoid(createNumber(0))`.
func (af *AstFactory) createUndefinedValue() *rhino.Node {
	return af.createVoid(af.CreateNumber(0))
}

// createNot mirrors AstFactory.createNot.
func (af *AstFactory) createNot(child *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.NOT)
	n.AddChildToBack(child)
	af.setType(n, jstype.Boolean)
	return n
}

// createNameWithUnknownType mirrors AstFactory.createNameWithUnknownType.
func (af *AstFactory) createNameWithUnknownType(name string) *rhino.Node {
	return af.CreateName(name, af.unknownType)
}

// createAssignStatement mirrors AstFactory.createAssignStatement.
func (af *AstFactory) createAssignStatement(lhs, rhs *rhino.Node) *rhino.Node {
	return af.CreateExprResult(af.CreateAssign(lhs, rhs))
}

// createObjectLit mirrors AstFactory.createObjectLit.
func (af *AstFactory) createObjectLit(elements ...*rhino.Node) *rhino.Node {
	n := rhino.New(rhino.OBJECTLIT)
	for _, e := range elements {
		n.AddChildToBack(e)
	}
	af.setType(n, af.unknownType)
	return n
}

// createArraylit mirrors AstFactory.createArraylit.
func (af *AstFactory) createArraylit(elements ...*rhino.Node) *rhino.Node {
	n := rhino.New(rhino.ARRAYLIT)
	for _, e := range elements {
		n.AddChildToBack(e)
	}
	af.setType(n, af.unknownType)
	return n
}

// createAnd mirrors AstFactory.createAnd.
func (af *AstFactory) createAnd(left, right *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.AND)
	n.AddChildToBack(left)
	n.AddChildToBack(right)
	// AND's type is the union of both operands for JSType mode.
	af.setType(n, af.unknownType)
	return n
}

// createOr mirrors AstFactory.createOr.
func (af *AstFactory) createOr(left, right *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.OR)
	n.AddChildToBack(left)
	n.AddChildToBack(right)
	af.setType(n, af.unknownType)
	return n
}

// createAdd mirrors AstFactory.createAdd.
func (af *AstFactory) createAdd(left, right *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.ADD)
	n.AddChildToBack(left)
	n.AddChildToBack(right)
	af.setType(n, jstype.Number)
	return n
}

// createSub mirrors AstFactory.createSub.
func (af *AstFactory) createSub(left, right *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.SUB)
	n.AddChildToBack(left)
	n.AddChildToBack(right)
	af.setType(n, jstype.Number)
	return n
}

// createInc mirrors AstFactory.createInc.
func (af *AstFactory) createInc(operand *rhino.Node, _ bool) *rhino.Node {
	n := rhino.New(rhino.INC)
	n.AddChildToBack(operand)
	af.setType(n, jstype.Number)
	return n
}

// createLessThan mirrors AstFactory.createLessThan.
func (af *AstFactory) createLessThan(left, right *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.LT)
	n.AddChildToBack(left)
	n.AddChildToBack(right)
	af.setType(n, jstype.Boolean)
	return n
}

// createSheq mirrors AstFactory.createSheq.
func (af *AstFactory) createSheq(expr1, expr2 *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.SHEQ)
	n.AddChildToBack(expr1)
	n.AddChildToBack(expr2)
	af.setType(n, jstype.Boolean)
	return n
}

// createEq mirrors AstFactory.createEq.
func (af *AstFactory) createEq(expr1, expr2 *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.EQ)
	n.AddChildToBack(expr1)
	n.AddChildToBack(expr2)
	af.setType(n, jstype.Boolean)
	return n
}

// createNe mirrors AstFactory.createNe.
func (af *AstFactory) createNe(expr1, expr2 *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.NE)
	n.AddChildToBack(expr1)
	n.AddChildToBack(expr2)
	af.setType(n, jstype.Boolean)
	return n
}

// createHook mirrors AstFactory.createHook.
func (af *AstFactory) createHook(condition, expr1, expr2 *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.HOOK)
	n.AddChildToBack(condition)
	n.AddChildToBack(expr1)
	n.AddChildToBack(expr2)
	af.setType(n, af.unknownType)
	return n
}

// createNewNode mirrors AstFactory.createNewNode.
func (af *AstFactory) createNewNode(target *rhino.Node, args ...*rhino.Node) *rhino.Node {
	n := rhino.New(rhino.NEW)
	n.AddChildToBack(target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	af.setType(n, af.unknownType)
	return n
}

// createCallWithUnknownType mirrors AstFactory.createCallWithUnknownType.
func (af *AstFactory) createCallWithUnknownType(callee *rhino.Node, args ...*rhino.Node) *rhino.Node {
	return af.CreateCall(callee, af.unknownType, args...)
}

// createLabel mirrors AstFactory.createLabel.
func (af *AstFactory) createLabel(label, stmt *rhino.Node) *rhino.Node {
	n := rhino.New(rhino.LABEL)
	n.AddChildToBack(label)
	n.AddChildToBack(stmt)
	return n
}

// createEmptyFunction mirrors AstFactory.createEmptyFunction(Type).
func (af *AstFactory) createEmptyFunction(t jstype.JSType) *rhino.Node {
	name := rhino.NewString(rhino.NAME, "")
	params := rhino.New(rhino.PARAM_LIST)
	body := rhino.New(rhino.BLOCK)
	n := rhino.New(rhino.FUNCTION)
	n.AddChildToBack(name)
	n.AddChildToBack(params)
	n.AddChildToBack(body)
	af.setType(n, t)
	return n
}

// createSingleLetNameDeclaration mirrors
// AstFactory.createSingleLetNameDeclaration(String).
func (af *AstFactory) createSingleLetNameDeclaration(variableName string) *rhino.Node {
	nameNode := af.CreateName(variableName, jstype.Void)
	n := rhino.New(rhino.LET)
	n.AddChildToBack(nameNode)
	return n
}

// createSingleVarNameDeclaration mirrors
// AstFactory.createSingleVarNameDeclaration(String).
func (af *AstFactory) createSingleVarNameDeclaration(variableName string) *rhino.Node {
	nameNode := af.CreateName(variableName, jstype.Void)
	n := rhino.New(rhino.VAR)
	n.AddChildToBack(nameNode)
	return n
}

// createParamList mirrors AstFactory.createParamList(String...).
func (af *AstFactory) createParamList(parameterNames ...string) *rhino.Node {
	n := rhino.New(rhino.PARAM_LIST)
	for _, pname := range parameterNames {
		p := af.CreateName(pname, af.unknownType)
		n.AddChildToBack(p)
	}
	return n
}

// isAddingTypes mirrors AstFactory.isAddingTypes (lowercase alias).
func (af *AstFactory) isAddingTypes() bool { return af.IsAddingTypes() }

// isAddingColors mirrors AstFactory.isAddingColors (lowercase alias).
func (af *AstFactory) isAddingColors() bool { return af.IsAddingColors() }

// exprResult mirrors AstFactory.exprResult (lowercase Java name).
func (af *AstFactory) exprResult(expr *rhino.Node) *rhino.Node {
	return af.CreateExprResult(expr)
}

// createEmpty mirrors AstFactory.createEmpty (lowercase Java alias).
func (af *AstFactory) createEmpty() *rhino.Node {
	return af.CreateEmpty()
}

// createBlock mirrors AstFactory.createBlock (lowercase Java alias).
func (af *AstFactory) createBlock(stmts ...*rhino.Node) *rhino.Node {
	return af.CreateBlock(stmts...)
}

// createIf mirrors AstFactory.createIf (lowercase Java alias).
func (af *AstFactory) createIf(cond, then *rhino.Node) *rhino.Node {
	return af.CreateIf(cond, then)
}

// createReturn mirrors AstFactory.createReturn (lowercase Java alias).
func (af *AstFactory) createReturn(value *rhino.Node) *rhino.Node {
	return af.CreateReturn(value)
}

// createThrow mirrors AstFactory.createThrow (lowercase Java alias).
func (af *AstFactory) createThrow(expr *rhino.Node) *rhino.Node {
	return af.CreateThrow(expr)
}

// createString mirrors AstFactory.createString (lowercase Java alias).
func (af *AstFactory) createString(value string) *rhino.Node {
	return af.CreateString(value)
}

// createNumber mirrors AstFactory.createNumber (lowercase Java alias).
func (af *AstFactory) createNumber(value float64) *rhino.Node {
	return af.CreateNumber(value)
}

// createBoolean mirrors AstFactory.createBoolean (lowercase Java alias).
func (af *AstFactory) createBoolean(value bool) *rhino.Node {
	return af.CreateBoolean(value)
}

// createNull mirrors AstFactory.createNull (lowercase Java alias).
func (af *AstFactory) createNull() *rhino.Node {
	return af.CreateNull()
}

// createThis mirrors AstFactory.createThis (lowercase Java alias).
func (af *AstFactory) createThis(thisType jstype.JSType) *rhino.Node {
	return af.CreateThis(thisType)
}

// createAssign mirrors AstFactory.createAssign (lowercase Java alias).
func (af *AstFactory) createAssign(lhs, rhs *rhino.Node) *rhino.Node {
	return af.CreateAssign(lhs, rhs)
}

// createName mirrors AstFactory.createName (lowercase Java alias).
func (af *AstFactory) createName(name string, t jstype.JSType) *rhino.Node {
	return af.CreateName(name, t)
}

// createCall mirrors AstFactory.createCall (lowercase Java alias).
func (af *AstFactory) createCall(callee *rhino.Node, resultType jstype.JSType, args ...*rhino.Node) *rhino.Node {
	return af.CreateCall(callee, resultType, args...)
}

// createGetProp mirrors AstFactory.createGetProp (lowercase Java alias).
func (af *AstFactory) createGetProp(receiver *rhino.Node, propertyName string, t jstype.JSType) *rhino.Node {
	return af.CreateGetProp(receiver, propertyName, t)
}
