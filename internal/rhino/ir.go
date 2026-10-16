// Port of com/google/javascript/rhino/IR.java @ closure-compiler v20260415

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1/GPL 2.0
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   John Lenz
 *
 * Alternatively, the contents of this file may be used under the terms of
 * the GNU General Public License Version 2 or later (the "GPL"), in which
 * case the provisions of the GPL are applicable instead of those above. If
 * you wish to allow use of your version of this file only under the terms of
 * the GPL and not to allow others to use your version of this file under the
 * MPL, indicate your decision by deleting the provisions above and replacing
 * them with the notice and other provisions required by the GPL. If you do
 * not delete the provisions above, a recipient may use your version of this
 * file under either the MPL or the GPL.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"iter"
	"math"
	"math/big"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IR is an AST construction helper class.
type IR struct{}

// Java: IR()
func NewIR() *IR { return &IR{} }

// Java: IR.empty()
func IR_Empty() *Node {
	return NewNode(Token_EMPTY)
}

// Java: IR.export(Node)
func IR_Export(declaration *Node) *Node {
	return NewNodeWithChild(Token_EXPORT, declaration)
}

// Java: IR.importNode(Node, Node, Node)
func IR_ImportNode(name *Node, importSpecs *Node, moduleIdentifier *Node) *Node {
	jutil.CheckStateMsg1(name.IsName() || name.IsEmpty(), "%s", name)
	jutil.CheckStateMsg1(
		importSpecs.IsImportSpec() || importSpecs.IsImportStar() || importSpecs.IsEmpty(),
		"%s", importSpecs)

	jutil.CheckStateMsg1(moduleIdentifier.IsStringLit(), "%s", moduleIdentifier)
	return NewNodeWith3Children(Token_IMPORT, name, importSpecs, moduleIdentifier)
}

// Java: IR.importStar(String)
func IR_ImportStar(name string) *Node {
	return Node_NewStringWithToken(Token_IMPORT_STAR, name)
}

// Java: IR.function(Node, Node, Node)
func IR_Function(name *Node, params *Node, body *Node) *Node {
	jutil.CheckState(name.IsName())
	jutil.CheckState(params.IsParamList())
	jutil.CheckState(body.IsBlock() || body.IsEmpty())
	return NewNodeWith3Children(Token_FUNCTION, name, params, body)
}

// Java: IR.arrowFunction(Node, Node, Node)
func IR_ArrowFunction(name *Node, params *Node, body *Node) *Node {
	jutil.CheckState(name.IsName())
	jutil.CheckState(params.IsParamList())
	jutil.CheckState(body.IsBlock() || IR_MayBeExpression(body))
	func_ := NewNodeWith3Children(Token_FUNCTION, name, params, body)
	func_.SetIsArrowFunction(true)
	return func_
}

// Java: IR.paramList(Node...)
func IR_ParamList(params ...*Node) *Node {
	paramList := NewNode(Token_PARAM_LIST)
	for _, param := range params {
		jutil.CheckState(param.IsName() || param.IsRest() || param.IsDefaultValue())
		paramList.AddChildToBack(param)
	}
	return paramList
}

// Java: IR.root(Node...)
func IR_Root(rootChildren ...*Node) *Node {
	root := NewNode(Token_ROOT)
	for _, child := range rootChildren {
		jutil.CheckState(child.GetToken() == Token_ROOT || child.GetToken() == Token_SCRIPT)
		root.AddChildToBack(child)
	}
	return root
}

// Java: IR.block()
func IR_Block() *Node {
	return NewNode(Token_BLOCK)
}

// Java: IR.block(Node)
func IR_BlockOfNode(stmt *Node) *Node {
	jutil.CheckStateMsg1(IR_MayBeStatement(stmt), "Block node cannot contain %s", stmt.GetToken())
	return NewNodeWithChild(Token_BLOCK, stmt)
}

// Java: IR.block(Node...)
func IR_BlockOfNodes(stmts ...*Node) *Node {
	block := IR_Block()
	for _, stmt := range stmts {
		jutil.CheckState(IR_MayBeStatement(stmt))
		block.AddChildToBack(stmt)
	}
	return block
}

// Java: IR.block(List<Node>)
func IR_BlockOfList(stmts []*Node) *Node {
	paramList := IR_Block()
	for _, stmt := range stmts {
		jutil.CheckState(IR_MayBeStatement(stmt))
		paramList.AddChildToBack(stmt)
	}
	return paramList
}

// Java: IR.blockUnchecked(Node)
func IR_BlockUnchecked(stmt *Node) *Node {
	return NewNodeWithChild(Token_BLOCK, stmt)
}

// Java: IR.script()
func IR_Script() *Node {
	// TODO(johnlenz): finish setting up the SCRIPT node
	return NewNode(Token_SCRIPT)
}

// Java: IR.script(Node...)
func IR_ScriptOfNodes(stmts ...*Node) *Node {
	block := IR_Script()
	for _, stmt := range stmts {
		jutil.CheckState(IR_MayBeStatementNoReturn(stmt))
		block.AddChildToBack(stmt)
	}
	return block
}

// Java: IR.script(List<Node>)
func IR_ScriptOfList(stmts []*Node) *Node {
	paramList := IR_Script()
	for _, stmt := range stmts {
		jutil.CheckState(IR_MayBeStatementNoReturn(stmt))
		paramList.AddChildToBack(stmt)
	}
	return paramList
}

// Java: IR.var(Node)
func IR_Var(lhs *Node) *Node {
	return IR_Declaration(lhs, Token_VAR)
}

// Java: IR.var(Node, Node)
func IR_VarWithValue(lhs *Node, value *Node) *Node {
	return IR_DeclarationWithValue(lhs, value, Token_VAR)
}

// Java: IR.let(Node)
func IR_Let(lhs *Node) *Node {
	return IR_Declaration(lhs, Token_LET)
}

// Java: IR.let(Node, Node)
func IR_LetWithValue(lhs *Node, value *Node) *Node {
	return IR_DeclarationWithValue(lhs, value, Token_LET)
}

// Java: IR.constNode(Node, Node)
func IR_ConstNode(lhs *Node, value *Node) *Node {
	return IR_DeclarationWithValue(lhs, value, Token_CONST)
}

// Java: IR.declaration(Node, Token)
func IR_Declaration(lhs *Node, type_ Token) *Node {
	jutil.CheckStateMsg1(lhs.IsName() || lhs.IsDestructuringPattern() || lhs.IsDestructuringLhs(), "%s", lhs)
	if lhs.IsDestructuringPattern() {
		lhs = NewNodeWithChild(Token_DESTRUCTURING_LHS, lhs)
	}
	return NewNodeWithChild(type_, lhs)
}

// Java: IR.declaration(Node, Node, Token)
func IR_DeclarationWithValue(lhs *Node, value *Node, type_ Token) *Node {
	if lhs.IsName() {
		jutil.CheckState(!lhs.HasChildren())
	} else {
		jutil.CheckState(lhs.IsArrayPattern() || lhs.IsObjectPattern())
		lhs = NewNodeWithChild(Token_DESTRUCTURING_LHS, lhs)
	}
	jutil.CheckStateMsg1(IR_MayBeExpression(value), "%s can't be an expression", value)

	lhs.AddChildToBack(value)
	return NewNodeWithChild(type_, lhs)
}

// Java: IR.returnNode()
func IR_ReturnNode() *Node {
	return NewNode(Token_RETURN)
}

// Java: IR.returnNode(Node)
func IR_ReturnNodeWithExpr(expr *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(expr))
	return NewNodeWithChild(Token_RETURN, expr)
}

// Java: IR.yieldNode(Node)
func IR_YieldNode(expr *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(expr))
	return NewNodeWithChild(Token_YIELD, expr)
}

// Java: IR.await(Node)
func IR_Await(expr *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(expr))
	return NewNodeWithChild(Token_AWAIT, expr)
}

// Java: IR.throwNode(Node)
func IR_ThrowNode(expr *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(expr))
	return NewNodeWithChild(Token_THROW, expr)
}

// Java: IR.exprResult(Node)
func IR_ExprResult(expr *Node) *Node {
	jutil.CheckStateMsg1(IR_MayBeExpression(expr), "%s", expr)
	return NewNodeWithChild(Token_EXPR_RESULT, expr)
}

// Java: IR.ifNode(Node, Node)
func IR_IfNode(cond *Node, then *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(cond))
	jutil.CheckState(then.IsBlock())
	return NewNodeWith2Children(Token_IF, cond, then)
}

// Java: IR.ifNode(Node, Node, Node)
func IR_IfNodeWithElse(cond *Node, then *Node, elseNode *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(cond))
	jutil.CheckState(then.IsBlock())
	jutil.CheckState(elseNode.IsBlock())
	return NewNodeWith3Children(Token_IF, cond, then, elseNode)
}

// Java: IR.doNode(Node, Node)
func IR_DoNode(body *Node, cond *Node) *Node {
	jutil.CheckState(body.IsBlock())
	jutil.CheckState(IR_MayBeExpression(cond))
	return NewNodeWith2Children(Token_DO, body, cond)
}

// Java: IR.whileNode(Node, Node)
func IR_WhileNode(cond *Node, body *Node) *Node {
	jutil.CheckState(body.IsBlock())
	jutil.CheckState(IR_MayBeExpression(cond))
	return NewNodeWith2Children(Token_WHILE, cond, body)
}

// Java: IR.forIn(Node, Node, Node)
func IR_ForIn(target *Node, cond *Node, body *Node) *Node {
	jutil.CheckState(target.IsVar() || IR_MayBeExpression(target))
	jutil.CheckState(IR_MayBeExpression(cond))
	jutil.CheckState(body.IsBlock())
	return NewNodeWith3Children(Token_FOR_IN, target, cond, body)
}

// Java: IR.forNode(Node, Node, Node, Node)
func IR_ForNode(init *Node, cond *Node, incr *Node, body *Node) *Node {
	jutil.CheckState(init.IsVar() || init.IsLet() || init.IsConst() || IR_MayBeExpressionOrEmpty(init))
	jutil.CheckState(IR_MayBeExpressionOrEmpty(cond))
	jutil.CheckState(IR_MayBeExpressionOrEmpty(incr))
	jutil.CheckState(body.IsBlock())
	r := NewNodeWith3Children(Token_FOR, init, cond, incr)
	r.AddChildToBack(body)
	return r
}

// Java: IR.switchNode(Node, Node...)
func IR_SwitchNode(cond *Node, cases ...*Node) *Node {
	jutil.CheckState(IR_MayBeExpression(cond))
	switchNode := NewNodeWithChild(Token_SWITCH, cond)
	switchBody := NewNode(Token_SWITCH_BODY)
	switchNode.AddChildToBack(switchBody)
	for _, caseNode := range cases {
		jutil.CheckState(caseNode.IsCase() || caseNode.IsDefaultCase())
		switchBody.AddChildToBack(caseNode)
	}
	return switchNode
}

// Java: IR.caseNode(Node, Node)
func IR_CaseNode(expr *Node, body *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(expr))
	jutil.CheckState(body.IsBlock())
	body.SetIsAddedBlock(true)
	return NewNodeWith2Children(Token_CASE, expr, body)
}

// Java: IR.defaultCase(Node)
func IR_DefaultCase(body *Node) *Node {
	jutil.CheckState(body.IsBlock())
	body.SetIsAddedBlock(true)
	return NewNodeWithChild(Token_DEFAULT_CASE, body)
}

// Java: IR.label(Node, Node)
func IR_Label(name *Node, stmt *Node) *Node {
	// TODO(johnlenz): additional validation here.
	jutil.CheckState(name.IsLabelName())
	jutil.CheckState(IR_MayBeStatement(stmt))
	return NewNodeWith2Children(Token_LABEL, name, stmt)
}

// Java: IR.labelName(String)
func IR_LabelName(name string) *Node {
	jutil.CheckState(name != "")
	return Node_NewStringWithToken(Token_LABEL_NAME, name)
}

// Java: IR.tryFinally(Node, Node)
func IR_TryFinally(tryBody *Node, finallyBody *Node) *Node {
	jutil.CheckState(tryBody.IsBlock())
	jutil.CheckState(finallyBody.IsBlock())
	catchBody := IR_Block().SrcrefIfMissing(tryBody)
	return NewNodeWith3Children(Token_TRY, tryBody, catchBody, finallyBody)
}

// Java: IR.tryCatch(Node, Node)
func IR_TryCatch(tryBody *Node, catchNode *Node) *Node {
	jutil.CheckState(tryBody.IsBlock())
	jutil.CheckState(catchNode.IsCatch())
	catchBody := IR_BlockUnchecked(catchNode).SrcrefIfMissing(catchNode)
	return NewNodeWith2Children(Token_TRY, tryBody, catchBody)
}

// Java: IR.tryCatchFinally(Node, Node, Node)
func IR_TryCatchFinally(tryBody *Node, catchNode *Node, finallyBody *Node) *Node {
	jutil.CheckState(finallyBody.IsBlock())
	tryNode := IR_TryCatch(tryBody, catchNode)
	tryNode.AddChildToBack(finallyBody)
	return tryNode
}

// Java: IR.catchNode(Node, Node)
func IR_CatchNode(expr *Node, body *Node) *Node {
	jutil.CheckState(expr.IsName())
	jutil.CheckState(body.IsBlock())
	return NewNodeWith2Children(Token_CATCH, expr, body)
}

// Java: IR.breakNode()
func IR_BreakNode() *Node {
	return NewNode(Token_BREAK)
}

// Java: IR.breakNode(Node)
func IR_BreakNodeWithLabel(name *Node) *Node {
	// TODO(johnlenz): additional validation here.
	jutil.CheckState(name.IsLabelName())
	return NewNodeWithChild(Token_BREAK, name)
}

// Java: IR.continueNode()
func IR_ContinueNode() *Node {
	return NewNode(Token_CONTINUE)
}

// Java: IR.continueNode(Node)
func IR_ContinueNodeWithLabel(name *Node) *Node {
	// TODO(johnlenz): additional validation here.
	jutil.CheckState(name.IsLabelName())
	return NewNodeWithChild(Token_CONTINUE, name)
}

// Java: IR.call(Node, Node...)
func IR_Call(target *Node, args ...*Node) *Node {
	call := NewNodeWithChild(Token_CALL, target)
	for _, arg := range args {
		jutil.CheckStateMsg1(IR_MayBeExpression(arg) || arg.IsSpread(), "%s", arg)
		call.AddChildToBack(arg)
	}
	return call
}

// Java: IR.startOptChainCall(Node, Node...)
func IR_StartOptChainCall(target *Node, args ...*Node) *Node {
	call := NewNodeWithChild(Token_OPTCHAIN_CALL, target)
	for _, arg := range args {
		jutil.CheckStateMsg1(IR_MayBeExpression(arg) || arg.IsSpread(), "%s", arg)
		call.AddChildToBack(arg)
	}
	call.SetIsOptionalChainStart(true)
	return call
}

// Java: IR.continueOptChainCall(Node, Node...)
func IR_ContinueOptChainCall(target *Node, args ...*Node) *Node {
	call := NewNodeWithChild(Token_OPTCHAIN_CALL, target)
	for _, arg := range args {
		jutil.CheckStateMsg1(IR_MayBeExpression(arg) || arg.IsSpread(), "%s", arg)
		call.AddChildToBack(arg)
	}
	call.SetIsOptionalChainStart(false)
	return call
}

// Java: IR.newNode(Node, Node...)
func IR_NewNode(target *Node, args ...*Node) *Node {
	newcall := NewNodeWithChild(Token_NEW, target)
	for _, arg := range args {
		jutil.CheckStateMsg1(IR_MayBeExpression(arg) || arg.IsSpread(), "%s", arg)
		newcall.AddChildToBack(arg)
	}
	return newcall
}

// Java: IR.name(String)
func IR_Name(name string) *Node {
	jutil.CheckStateMsg1(
		jutil.StringIndexOfChar(name, '.') == -1, "Invalid name '%s'. Did you mean to use NodeUtil.newQName?", name)

	return Node_NewStringWithToken(Token_NAME, name)
}

// Java: IR.startOptChainGetprop(Node, String)
func IR_StartOptChainGetprop(target *Node, prop string) *Node {
	jutil.CheckStateMsg1(IR_MayBeExpression(target), "%s", target)
	optChainGetProp := Node_NewStringWithToken(Token_OPTCHAIN_GETPROP, prop)
	optChainGetProp.AddChildToBack(target)
	optChainGetProp.SetIsOptionalChainStart(true)
	return optChainGetProp
}

// Java: IR.continueOptChainGetprop(Node, String)
func IR_ContinueOptChainGetprop(target *Node, prop string) *Node {
	jutil.CheckStateMsg1(IR_MayBeExpression(target), "%s", target)
	optChainGetProp := Node_NewStringWithToken(Token_OPTCHAIN_GETPROP, prop)
	optChainGetProp.AddChildToBack(target)
	optChainGetProp.SetIsOptionalChainStart(false)
	return optChainGetProp
}

// Java: IR.getprop(Node, String)
func IR_Getprop(target *Node, prop string) *Node {
	jutil.CheckState(IR_MayBeExpression(target))
	getprop := Node_NewStringWithToken(Token_GETPROP, prop)
	getprop.AddChildToBack(target)
	return getprop
}

// Java: IR.getprop(Node, String, String...)
func IR_GetpropWithMoreProps(target *Node, prop string, moreProps ...string) *Node {
	jutil.CheckState(IR_MayBeExpression(target))
	result := IR_Getprop(target, prop)
	for _, moreProp := range moreProps {
		result = IR_Getprop(result, moreProp)
	}
	return result
}

// Java: IR.startOptChainGetelem(Node, Node)
func IR_StartOptChainGetelem(target *Node, elem *Node) *Node {
	jutil.CheckStateMsg1(IR_MayBeExpression(target), "%s", target)
	jutil.CheckStateMsg1(IR_MayBeExpression(elem), "%s", elem)
	optChainGetElem := NewNodeWith2Children(Token_OPTCHAIN_GETELEM, target, elem)
	optChainGetElem.SetIsOptionalChainStart(true)
	return optChainGetElem
}

// Java: IR.continueOptChainGetelem(Node, Node)
func IR_ContinueOptChainGetelem(target *Node, elem *Node) *Node {
	jutil.CheckStateMsg1(IR_MayBeExpression(target), "%s", target)
	jutil.CheckStateMsg1(IR_MayBeExpression(elem), "%s", elem)
	optChainGetElem := NewNodeWith2Children(Token_OPTCHAIN_GETELEM, target, elem)
	optChainGetElem.SetIsOptionalChainStart(false)
	return optChainGetElem
}

// Java: IR.getelem(Node, Node)
func IR_Getelem(target *Node, elem *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(target))
	jutil.CheckState(IR_MayBeExpression(elem))
	return NewNodeWith2Children(Token_GETELEM, target, elem)
}

// Java: IR.delprop(Node)
func IR_Delprop(target *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(target))
	return NewNodeWithChild(Token_DELPROP, target)
}

// Java: IR.assign(Node, Node)
func IR_Assign(target *Node, expr *Node) *Node {
	jutil.CheckStateMsg1(target.IsValidAssignmentTarget(), "%s", target)
	jutil.CheckStateMsg1(IR_MayBeExpression(expr), "%s", expr)
	return NewNodeWith2Children(Token_ASSIGN, target, expr)
}

// Java: IR.hook(Node, Node, Node)
func IR_Hook(cond *Node, trueval *Node, falseval *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(cond))
	jutil.CheckState(IR_MayBeExpression(trueval))
	jutil.CheckState(IR_MayBeExpression(falseval))
	return NewNodeWith3Children(Token_HOOK, cond, trueval, falseval)
}

// Java: IR.in(Node, Node)
func IR_In(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_IN, expr1, expr2)
}

// Java: IR.comma(Node, Node)
func IR_Comma(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_COMMA, expr1, expr2)
}

// Java: IR.and(Node, Node)
func IR_And(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_AND, expr1, expr2)
}

// Java: IR.or(Node, Node)
func IR_Or(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_OR, expr1, expr2)
}

// Java: IR.coalesce(Node, Node)
func IR_Coalesce(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_COALESCE, expr1, expr2)
}

// Java: IR.not(Node)
func IR_Not(expr1 *Node) *Node {
	return IR_UnaryOp(Token_NOT, expr1)
}

// IR_Lt is "&lt;".
//
// Java: IR.lt(Node, Node)
func IR_Lt(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_LT, expr1, expr2)
}

// IR_Ge is "&gt;=".
//
// Java: IR.ge(Node, Node)
func IR_Ge(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_GE, expr1, expr2)
}

// IR_Eq is "==".
//
// Java: IR.eq(Node, Node)
func IR_Eq(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_EQ, expr1, expr2)
}

// IR_Ne is "!=".
//
// Java: IR.ne(Node, Node)
func IR_Ne(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_NE, expr1, expr2)
}

// IR_Sheq is "===".
//
// Java: IR.sheq(Node, Node)
func IR_Sheq(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_SHEQ, expr1, expr2)
}

// IR_Shne is "!==".
//
// Java: IR.shne(Node, Node)
func IR_Shne(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_SHNE, expr1, expr2)
}

// Java: IR.voidNode(Node)
func IR_VoidNode(expr1 *Node) *Node {
	return IR_UnaryOp(Token_VOID, expr1)
}

// Java: IR.neg(Node)
func IR_Neg(expr1 *Node) *Node {
	return IR_UnaryOp(Token_NEG, expr1)
}

// Java: IR.pos(Node)
func IR_Pos(expr1 *Node) *Node {
	return IR_UnaryOp(Token_POS, expr1)
}

// Java: IR.cast(Node, JSDocInfo)
func IR_Cast(expr1 *Node, jsdoc *JSDocInfo) *Node {
	op := IR_UnaryOp(Token_CAST, expr1)
	op.SetJSDocInfo(jsdoc)
	return op
}

// Java: IR.inc(Node, boolean)
func IR_Inc(exp *Node, isPost bool) *Node {
	op := IR_UnaryOp(Token_INC, exp)
	op.PutBooleanProp(Node_INCRDECR_PROP, isPost)
	return op
}

// Java: IR.dec(Node, boolean)
func IR_Dec(exp *Node, isPost bool) *Node {
	op := IR_UnaryOp(Token_DEC, exp)
	op.PutBooleanProp(Node_INCRDECR_PROP, isPost)
	return op
}

// Java: IR.add(Node, Node)
func IR_Add(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_ADD, expr1, expr2)
}

// Java: IR.sub(Node, Node)
func IR_Sub(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_SUB, expr1, expr2)
}

// IR_AssignOr is "||=".
//
// Java: IR.assignOr(Node, Node)
func IR_AssignOr(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_ASSIGN_OR, expr1, expr2)
}

// IR_AssignAnd is "&&=".
//
// Java: IR.assignAnd(Node, Node)
func IR_AssignAnd(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_ASSIGN_AND, expr1, expr2)
}

// IR_AssignCoalesce is "??=".
//
// Java: IR.assignCoalesce(Node, Node)
func IR_AssignCoalesce(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_ASSIGN_COALESCE, expr1, expr2)
}

// IR_BitwiseAnd is "&".
//
// Java: IR.bitwiseAnd(Node, Node)
func IR_BitwiseAnd(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_BITAND, expr1, expr2)
}

// IR_RightShift is ">>".
//
// Java: IR.rightShift(Node, Node)
func IR_RightShift(expr1 *Node, expr2 *Node) *Node {
	return IR_BinaryOp(Token_RSH, expr1, expr2)
}

// TODO(johnlenz): the rest of the ops

// literals

// Java: IR.objectlit(Node...)
func IR_Objectlit(propdefs ...*Node) *Node {
	objectlit := NewNode(Token_OBJECTLIT)
	for _, propdef := range propdefs {
		switch propdef.GetToken() {
		case Token_STRING_KEY,
			Token_MEMBER_FUNCTION_DEF,
			Token_GETTER_DEF,
			Token_SETTER_DEF,
			Token_OBJECT_SPREAD,
			Token_COMPUTED_PROP:
		default:
			panic(jutil.NewIllegalStateExceptionMsg("Unexpected OBJECTLIT child: " + propdef.String()))
		}

		objectlit.AddChildToBack(propdef)
	}
	return objectlit
}

// Java: IR.objectPattern(Node...)
func IR_ObjectPattern(keys ...*Node) *Node {
	objectPattern := NewNode(Token_OBJECT_PATTERN)
	for _, key := range keys {
		jutil.CheckState(key.IsStringKey() || key.IsComputedProp() || key.IsRest())
		objectPattern.AddChildToBack(key)
	}
	return objectPattern
}

// Java: IR.arrayPattern(Node...)
func IR_ArrayPattern(keys ...*Node) *Node {
	arrayPattern := NewNode(Token_ARRAY_PATTERN)
	for _, key := range keys {
		jutil.CheckState(key.IsRest() || key.IsValidAssignmentTarget())
		arrayPattern.AddChildToBack(key)
	}
	return arrayPattern
}

// Java: IR.computedProp(Node, Node)
func IR_ComputedProp(key *Node, value *Node) *Node {
	jutil.CheckStateMsg1(IR_MayBeExpression(key), "%s", key)
	jutil.CheckStateMsg1(IR_MayBeExpression(value), "%s", value)
	return NewNodeWith2Children(Token_COMPUTED_PROP, key, value)
}

// Java: IR.propdef(Node, Node)
func IR_Propdef(string_ *Node, value *Node) *Node {
	jutil.CheckState(string_.IsStringKey())
	jutil.CheckState(!string_.HasChildren())
	jutil.CheckState(IR_MayBeExpression(value))
	string_.AddChildToFront(value)
	return string_
}

// Java: IR.arraylit(Node...)
func IR_Arraylit(exprs ...*Node) *Node {
	return IR_ArraylitOfIterable(slices.Values(exprs))
}

// Java: IR.arraylit(Iterable<Node>)
func IR_ArraylitOfIterable(exprs iter.Seq[*Node]) *Node {
	arraylit := NewNode(Token_ARRAYLIT)
	for expr := range exprs {
		jutil.CheckStateMsg1(IR_MayBeExpressionOrEmpty(expr) || expr.IsSpread(), "%s", expr)
		arraylit.AddChildToBack(expr)
	}
	return arraylit
}

// Java: IR.regexp(Node)
func IR_Regexp(expr *Node) *Node {
	jutil.CheckState(expr.IsStringLit())
	return NewNodeWithChild(Token_REGEXP, expr)
}

// Java: IR.regexp(Node, Node)
func IR_RegexpWithFlags(expr *Node, flags *Node) *Node {
	jutil.CheckState(expr.IsStringLit())
	jutil.CheckState(flags.IsStringLit())
	return NewNodeWith2Children(Token_REGEXP, expr, flags)
}

// Java: IR.string(String)
func IR_String(s string) *Node {
	return Node_NewString(s)
}

// Java: IR.stringKey(String)
func IR_StringKey(s string) *Node {
	return Node_NewStringWithToken(Token_STRING_KEY, s)
}

// Java: IR.stringKey(String, Node)
func IR_StringKeyWithValue(s string, value *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(value) || value.IsDefaultValue() || value.IsObjectPattern())
	stringKey := IR_StringKey(s)
	stringKey.AddChildToFront(value)
	return stringKey
}

// Java: IR.quotedStringKey(String, Node)
func IR_QuotedStringKey(s string, value *Node) *Node {
	k := IR_StringKeyWithValue(s, value)
	k.PutBooleanProp(Node_QUOTED_PROP, true)
	return k
}

// Java: IR.templateLiteral()
func IR_TemplateLiteral() *Node {
	return NewNode(Token_TEMPLATELIT)
}

// Java: IR.templateLiteralString(String, String)
func IR_TemplateLiteralString(cooked *string, raw string) *Node {
	return Node_NewTemplateLitString(cooked, raw)
}

// Java: IR.templateLiteralSubstitution(Node)
func IR_TemplateLiteralSubstitution(child *Node) *Node {
	sub := NewNode(Token_TEMPLATELIT_SUB)
	sub.AddChildToBack(child)
	return sub
}

// Java: IR.iterRest(Node)
func IR_IterRest(target *Node) *Node {
	jutil.CheckStateMsg1(target.IsValidAssignmentTarget(), "%s", target)
	return NewNodeWithChild(Token_ITER_REST, target)
}

// Java: IR.objectRest(Node)
func IR_ObjectRest(target *Node) *Node {
	jutil.CheckStateMsg1(target.IsValidAssignmentTarget(), "%s", target)
	return NewNodeWithChild(Token_OBJECT_REST, target)
}

// Java: IR.iterSpread(Node)
func IR_IterSpread(expr *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(expr))
	return NewNodeWithChild(Token_ITER_SPREAD, expr)
}

// Java: IR.objectSpread(Node)
func IR_ObjectSpread(expr *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(expr))
	return NewNodeWithChild(Token_OBJECT_SPREAD, expr)
}

// Java: IR.superNode()
func IR_SuperNode() *Node {
	return NewNode(Token_SUPER)
}

// Java: IR.getterDef(String, Node)
func IR_GetterDef(name string, value *Node) *Node {
	jutil.CheckState(value.IsFunction())
	member := Node_NewStringWithToken(Token_GETTER_DEF, name)
	member.AddChildToFront(value)
	return member
}

// Java: IR.setterDef(String, Node)
func IR_SetterDef(name string, value *Node) *Node {
	jutil.CheckState(value.IsFunction())
	member := Node_NewStringWithToken(Token_SETTER_DEF, name)
	member.AddChildToFront(value)
	return member
}

// Java: IR.memberFieldDef(String, Node)
func IR_MemberFieldDef(name string, value *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(value))
	member := Node_NewStringWithToken(Token_MEMBER_FIELD_DEF, name)
	member.AddChildToFront(value)
	return member
}

// Java: IR.memberFunctionDef(String, Node)
func IR_MemberFunctionDef(name string, function *Node) *Node {
	jutil.CheckState(function.IsFunction())
	member := Node_NewStringWithToken(Token_MEMBER_FUNCTION_DEF, name)
	member.AddChildToBack(function)
	return member
}

// Java: IR.number(double)
func IR_Number(d float64) *Node {
	jutil.CheckStateMsg1(!math.IsNaN(d), "%s", d)
	jutil.CheckStateMsg1(base.JSCompDoubles_IsPositive(d), "%s", d)
	return Node_NewNumber(d)
}

// Java: IR.bigint(BigInteger)
func IR_Bigint(b *big.Int) *Node {
	jutil.CheckNotNull(b)
	jutil.CheckStateMsg1(b.Sign() >= 0, "%s", b)
	return Node_NewBigInt(b)
}

// Java: IR.thisNode()
func IR_ThisNode() *Node {
	return NewNode(Token_THIS)
}

// Java: IR.trueNode()
func IR_TrueNode() *Node {
	return NewNode(Token_TRUE)
}

// Java: IR.falseNode()
func IR_FalseNode() *Node {
	return NewNode(Token_FALSE)
}

// Java: IR.nullNode()
func IR_NullNode() *Node {
	return NewNode(Token_NULL)
}

// Java: IR.typeof(Node)
func IR_Typeof(expr *Node) *Node {
	return IR_UnaryOp(Token_TYPEOF, expr)
}

// Java: IR.importMeta()
func IR_ImportMeta() *Node {
	return NewNode(Token_IMPORT_META)
}

// helper methods

// Java: IR.binaryOp(Token, Node, Node)
func IR_BinaryOp(token Token, expr1 *Node, expr2 *Node) *Node {
	jutil.CheckStateMsg1(IR_MayBeExpression(expr1), "%s", expr1)
	jutil.CheckStateMsg1(IR_MayBeExpression(expr2), "%s", expr2)
	return NewNodeWith2Children(token, expr1, expr2)
}

// Java: IR.unaryOp(Token, Node)
func IR_UnaryOp(token Token, expr *Node) *Node {
	jutil.CheckState(IR_MayBeExpression(expr))
	return NewNodeWithChild(token, expr)
}

// Java: IR.mayBeExpressionOrEmpty(Node)
func IR_MayBeExpressionOrEmpty(n *Node) bool {
	return n.IsEmpty() || IR_MayBeExpression(n)
}

// NOTE: some nodes are neither statements nor expression nodes:
//   SCRIPT, LABEL_NAME, PARAM_LIST, CASE, DEFAULT_CASE, CATCH
//   GETTER_DEF, SETTER_DEF

// IR_MayBeStatementNoReturn: it isn't possible to always determine if a
// detached node is a expression, so make a best guess.
//
// Java: IR.mayBeStatementNoReturn(Node)
func IR_MayBeStatementNoReturn(n *Node) bool {
	switch n.GetToken() {
	case Token_EMPTY, Token_FUNCTION:
		// EMPTY and FUNCTION are used both in expression and statement
		// contexts
		return true
	case Token_BLOCK,
		Token_BREAK,
		Token_CLASS,
		Token_CONST,
		Token_CONTINUE,
		Token_DEBUGGER,
		Token_DO,
		Token_ENUM,
		Token_EXPR_RESULT,
		Token_FOR,
		Token_FOR_IN,
		Token_FOR_OF,
		Token_FOR_AWAIT_OF,
		Token_IF,
		Token_INTERFACE,
		Token_LABEL,
		Token_LET,
		Token_SWITCH,
		Token_THROW,
		Token_TRY,
		Token_VAR,
		Token_WHILE,
		Token_WITH:
		return true
	default:
		return false
	}
}

// IR_MayBeStatement: it isn't possible to always determine if a detached
// node is a expression, so make a best guess.
//
// Java: IR.mayBeStatement(Node)
func IR_MayBeStatement(n *Node) bool {
	if !IR_MayBeStatementNoReturn(n) {
		return n.IsReturn()
	}
	return true
}

// IR_MayBeExpression: it isn't possible to always determine if a detached
// node is a expression, so make a best guess.
//
// Java: IR.mayBeExpression(Node)
func IR_MayBeExpression(n *Node) bool {
	switch n.GetToken() {
	case Token_FUNCTION, Token_CLASS:
		// FUNCTION and CLASS are used both in expression and statement
		// contexts.
		return true
	case Token_ADD,
		Token_AND,
		Token_ARRAYLIT,
		Token_ASSIGN,
		Token_ASSIGN_BITOR,
		Token_ASSIGN_BITXOR,
		Token_ASSIGN_BITAND,
		Token_ASSIGN_LSH,
		Token_ASSIGN_RSH,
		Token_ASSIGN_URSH,
		Token_ASSIGN_ADD,
		Token_ASSIGN_SUB,
		Token_ASSIGN_MUL,
		Token_ASSIGN_EXPONENT,
		Token_ASSIGN_DIV,
		Token_ASSIGN_MOD,
		Token_ASSIGN_OR,
		Token_ASSIGN_AND,
		Token_ASSIGN_COALESCE,
		Token_AWAIT,
		Token_BIGINT,
		Token_BITAND,
		Token_BITOR,
		Token_BITNOT,
		Token_BITXOR,
		Token_CALL,
		Token_CAST,
		Token_COALESCE,
		Token_COMMA,
		Token_DEC,
		Token_DELPROP,
		Token_DIV,
		Token_DYNAMIC_IMPORT,
		Token_EQ,
		Token_EXPONENT,
		Token_FALSE,
		Token_GE,
		Token_GETPROP,
		Token_GETELEM,
		Token_GT,
		Token_HOOK,
		Token_IMPORT_META,
		Token_IN,
		Token_INC,
		Token_INSTANCEOF,
		Token_LE,
		Token_LSH,
		Token_LT,
		Token_MOD,
		Token_MUL,
		Token_NAME,
		Token_NE,
		Token_NEG,
		Token_NEW,
		Token_NEW_TARGET,
		Token_NOT,
		Token_NUMBER,
		Token_NULL,
		Token_OBJECTLIT,
		Token_OPTCHAIN_CALL,
		Token_OPTCHAIN_GETELEM,
		Token_OPTCHAIN_GETPROP,
		Token_OR,
		Token_POS,
		Token_REGEXP,
		Token_RSH,
		Token_SHEQ,
		Token_SHNE,
		Token_STRINGLIT,
		Token_SUB,
		Token_SUPER,
		Token_TEMPLATELIT,
		Token_TAGGED_TEMPLATELIT,
		Token_THIS,
		Token_TYPEOF,
		Token_TRUE,
		Token_URSH,
		Token_VOID,
		Token_YIELD:
		return true
	default:
		return false
	}
}
