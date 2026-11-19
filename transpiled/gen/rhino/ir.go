/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
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
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/IR.java (generated mechanically by j2go).

package rhino

import (
	"math"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.IR
type IR struct {
	idHash_ atomic.Int32
}

func (this_ *IR) BindSelf_(s any) {
	_ = s
}
func (this_ *IR) IsJavaNil_() bool           { return this_ == nil }
func (*IR) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.IR", "IR" }
func (this_ *IR) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IR) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.IR", "com.google.javascript.rhino.IR", func(v_ any) bool {
		return j2gort.InstanceOf[*IR](v_)
	})
}

func (this_ *IR) initFields_IR() {
	_ = this_
}

// Java constructor: IR()
func NewIR() *IR {
	this_ := &IR{}
	InitIR(this_)
	return this_
}

func InitIR(this_ *IR) {
	this_.initFields_IR()
}

// Java static method: IR.empty()
func IR_Empty() Node {
	return NewNode_1(StaticGet_Token_EMPTY())
}

// Java static method: IR.export(com.google.javascript.rhino.Node declaration)
func IR_Export(declaration Node) Node {
	_ = declaration
	return NewNode_2(StaticGet_Token_EXPORT(), declaration)
}

// Java static method: IR.importNode(com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node importSpecs, com.google.javascript.rhino.Node moduleIdentifier)
func IR_ImportNode(name Node, importSpecs Node, moduleIdentifier Node) Node {
	_ = name
	_ = importSpecs
	_ = moduleIdentifier
	if !(name.IsName() || name.IsEmpty()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(name))
	}
	if !((importSpecs.IsImportSpec() || importSpecs.IsImportStar()) || importSpecs.IsEmpty()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(importSpecs))
	}
	if !(moduleIdentifier.IsStringLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(moduleIdentifier))
	}
	return NewNode_4(StaticGet_Token_IMPORT(), name, importSpecs, moduleIdentifier)
}

// Java static method: IR.importStar(java.lang.String name)
func IR_ImportStar(name string) Node {
	_ = name
	return Node_NewString_2(StaticGet_Token_IMPORT_STAR(), name)
}

// Java static method: IR.function(com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node params, com.google.javascript.rhino.Node body)
func IR_Function(name Node, params Node, body Node) Node {
	_ = name
	_ = params
	_ = body
	if !(name.IsName()) {
		j2gort.CheckState(false)
	}
	if !(params.IsParamList()) {
		j2gort.CheckState(false)
	}
	if !(body.IsBlock() || body.IsEmpty()) {
		j2gort.CheckState(false)
	}
	return NewNode_4(StaticGet_Token_FUNCTION(), name, params, body)
}

// Java static method: IR.arrowFunction(com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node params, com.google.javascript.rhino.Node body)
func IR_ArrowFunction(name Node, params Node, body Node) Node {
	_ = name
	_ = params
	_ = body
	if !(name.IsName()) {
		j2gort.CheckState(false)
	}
	if !(params.IsParamList()) {
		j2gort.CheckState(false)
	}
	if !(body.IsBlock() || IR_MayBeExpression(body)) {
		j2gort.CheckState(false)
	}
	var func_ Node = NewNode_4(StaticGet_Token_FUNCTION(), name, params, body)
	_ = func_
	func_.SetIsArrowFunction(true)
	return func_
}

// Java static method: IR.paramList(com.google.javascript.rhino.Node[] params)
func IR_ParamList(params []Node) Node {
	_ = params
	var paramList Node = NewNode_1(StaticGet_Token_PARAM_LIST())
	_ = paramList
	for _, el_1 := range j2gort.CheckNotNullSlice(params) {
		param := el_1
		_ = param
		if !((param.IsName() || param.IsRest()) || param.IsDefaultValue()) {
			j2gort.CheckState(false)
		}
		paramList.AddChildToBack(param)
	}
	return paramList
}

// Java static method: IR.root(com.google.javascript.rhino.Node[] rootChildren)
func IR_Root(rootChildren []Node) Node {
	_ = rootChildren
	var root Node = NewNode_1(StaticGet_Token_ROOT())
	_ = root
	for _, el_2 := range j2gort.CheckNotNullSlice(rootChildren) {
		child := el_2
		_ = child
		if !((child.RawToken_() == uint8(89)) || (child.RawToken_() == uint8(94))) {
			j2gort.CheckState(false)
		}
		root.AddChildToBack(child)
	}
	return root
}

// Java static method: IR.block()
func IR_Block_0() Node {
	return NewNode_1(StaticGet_Token_BLOCK())
}

// Java static method: IR.block(com.google.javascript.rhino.Node stmt)
func IR_Block_1_Node(stmt Node) Node {
	_ = stmt
	if !(IR_MayBeStatement(stmt)) {
		j2gort.CheckStateFmt(false, "Block node cannot contain %s", j2gort.BoxNullableVal(stmt.GetToken()))
	}
	return NewNode_2(StaticGet_Token_BLOCK(), stmt)
}

// Java static method: IR.block(com.google.javascript.rhino.Node[] stmts)
func IR_Block_1_Node_2(stmts []Node) Node {
	_ = stmts
	var block Node = IR_Block_0()
	_ = block
	for _, el_3 := range j2gort.CheckNotNullSlice(stmts) {
		stmt := el_3
		_ = stmt
		if !(IR_MayBeStatement(stmt)) {
			j2gort.CheckState(false)
		}
		block.AddChildToBack(stmt)
	}
	return block
}

// Java static method: IR.block(java.util.List<com.google.javascript.rhino.Node> stmts)
func IR_Block_1_Ext_java_util_List(stmts j2gort.Ext_java_util_List) Node {
	_ = stmts
	var paramList Node = IR_Block_0()
	_ = paramList
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(stmts))
		for cur_4.HasNext() {
			stmt := Node(cur_4.NextNode())
			_ = stmt
			if !(IR_MayBeStatement(stmt)) {
				j2gort.CheckState(false)
			}
			paramList.AddChildToBack(stmt)
		}
	}
	return paramList
}

// Java static method: IR.blockUnchecked(com.google.javascript.rhino.Node stmt)
func IR_BlockUnchecked(stmt Node) Node {
	_ = stmt
	return NewNode_2(StaticGet_Token_BLOCK(), stmt)
}

// Java static method: IR.script()
func IR_Script_0() Node {
	return NewNode_1(StaticGet_Token_SCRIPT())
}

// Java static method: IR.script(com.google.javascript.rhino.Node[] stmts)
func IR_Script_1_Node(stmts []Node) Node {
	_ = stmts
	var block Node = IR_Script_0()
	_ = block
	for _, el_5 := range j2gort.CheckNotNullSlice(stmts) {
		stmt := el_5
		_ = stmt
		if !(IR_MayBeStatementNoReturn(stmt)) {
			j2gort.CheckState(false)
		}
		block.AddChildToBack(stmt)
	}
	return block
}

// Java static method: IR.script(java.util.List<com.google.javascript.rhino.Node> stmts)
func IR_Script_1_Ext_java_util_List(stmts j2gort.Ext_java_util_List) Node {
	_ = stmts
	var paramList Node = IR_Script_0()
	_ = paramList
	{
		cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(stmts))
		for cur_6.HasNext() {
			stmt := Node(cur_6.NextNode())
			_ = stmt
			if !(IR_MayBeStatementNoReturn(stmt)) {
				j2gort.CheckState(false)
			}
			paramList.AddChildToBack(stmt)
		}
	}
	return paramList
}

// Java static method: IR.var(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node value)
func IR_Var_2(lhs Node, value Node) Node {
	_ = lhs
	_ = value
	return IR_Declaration_3(lhs, value, StaticGet_Token_VAR())
}

// Java static method: IR.var(com.google.javascript.rhino.Node lhs)
func IR_Var_1(lhs Node) Node {
	_ = lhs
	return IR_Declaration_2(lhs, StaticGet_Token_VAR())
}

// Java static method: IR.let(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node value)
func IR_Let_2(lhs Node, value Node) Node {
	_ = lhs
	_ = value
	return IR_Declaration_3(lhs, value, StaticGet_Token_LET())
}

// Java static method: IR.let(com.google.javascript.rhino.Node lhs)
func IR_Let_1(lhs Node) Node {
	_ = lhs
	return IR_Declaration_2(lhs, StaticGet_Token_LET())
}

// Java static method: IR.constNode(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node value)
func IR_ConstNode(lhs Node, value Node) Node {
	_ = lhs
	_ = value
	return IR_Declaration_3(lhs, value, StaticGet_Token_CONST())
}

// Java static method: IR.declaration(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Token type)
func IR_Declaration_2(lhs Node, type_ Token) Node {
	_ = lhs
	_ = type_
	if !((lhs.IsName() || lhs.IsDestructuringPattern()) || lhs.IsDestructuringLhs()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(lhs))
	}
	if lhs.IsDestructuringPattern() {
		lhs = NewNode_2(StaticGet_Token_DESTRUCTURING_LHS(), lhs)
	}
	return NewNode_2(type_, lhs)
}

// Java static method: IR.declaration(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node value, com.google.javascript.rhino.Token type)
func IR_Declaration_3(lhs Node, value Node, type_ Token) Node {
	_ = lhs
	_ = value
	_ = type_
	if lhs.IsName() {
		if !(!(lhs.HasChildren())) {
			j2gort.CheckState(false)
		}
	} else {
		if !(lhs.IsArrayPattern() || lhs.IsObjectPattern()) {
			j2gort.CheckState(false)
		}
		lhs = NewNode_2(StaticGet_Token_DESTRUCTURING_LHS(), lhs)
	}
	if !(IR_MayBeExpression(value)) {
		j2gort.CheckStateFmt(false, "%s can't be an expression", j2gort.BoxNode(value))
	}
	lhs.AddChildToBack(value)
	return NewNode_2(type_, lhs)
}

// Java static method: IR.returnNode()
func IR_ReturnNode_0() Node {
	return NewNode_1(StaticGet_Token_RETURN())
}

// Java static method: IR.returnNode(com.google.javascript.rhino.Node expr)
func IR_ReturnNode_1(expr Node) Node {
	_ = expr
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_RETURN(), expr)
}

// Java static method: IR.yieldNode(com.google.javascript.rhino.Node expr)
func IR_YieldNode(expr Node) Node {
	_ = expr
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_YIELD(), expr)
}

// Java static method: IR.await(com.google.javascript.rhino.Node expr)
func IR_Await(expr Node) Node {
	_ = expr
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_AWAIT(), expr)
}

// Java static method: IR.throwNode(com.google.javascript.rhino.Node expr)
func IR_ThrowNode(expr Node) Node {
	_ = expr
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_THROW(), expr)
}

// Java static method: IR.exprResult(com.google.javascript.rhino.Node expr)
func IR_ExprResult(expr Node) Node {
	_ = expr
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(expr))
	}
	return NewNode_2(StaticGet_Token_EXPR_RESULT(), expr)
}

// Java static method: IR.ifNode(com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node then)
func IR_IfNode_2(cond Node, then Node) Node {
	_ = cond
	_ = then
	if !(IR_MayBeExpression(cond)) {
		j2gort.CheckState(false)
	}
	if !(then.IsBlock()) {
		j2gort.CheckState(false)
	}
	return NewNode_3(StaticGet_Token_IF(), cond, then)
}

// Java static method: IR.ifNode(com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node then, com.google.javascript.rhino.Node elseNode)
func IR_IfNode_3(cond Node, then Node, elseNode Node) Node {
	_ = cond
	_ = then
	_ = elseNode
	if !(IR_MayBeExpression(cond)) {
		j2gort.CheckState(false)
	}
	if !(then.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(elseNode.IsBlock()) {
		j2gort.CheckState(false)
	}
	return NewNode_4(StaticGet_Token_IF(), cond, then, elseNode)
}

// Java static method: IR.doNode(com.google.javascript.rhino.Node body, com.google.javascript.rhino.Node cond)
func IR_DoNode(body Node, cond Node) Node {
	_ = body
	_ = cond
	if !(body.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpression(cond)) {
		j2gort.CheckState(false)
	}
	return NewNode_3(StaticGet_Token_DO(), body, cond)
}

// Java static method: IR.whileNode(com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node body)
func IR_WhileNode(cond Node, body Node) Node {
	_ = cond
	_ = body
	if !(body.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpression(cond)) {
		j2gort.CheckState(false)
	}
	return NewNode_3(StaticGet_Token_WHILE(), cond, body)
}

// Java static method: IR.forIn(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node body)
func IR_ForIn(target Node, cond Node, body Node) Node {
	_ = target
	_ = cond
	_ = body
	if !(target.IsVar() || IR_MayBeExpression(target)) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpression(cond)) {
		j2gort.CheckState(false)
	}
	if !(body.IsBlock()) {
		j2gort.CheckState(false)
	}
	return NewNode_4(StaticGet_Token_FOR_IN(), target, cond, body)
}

// Java static method: IR.forNode(com.google.javascript.rhino.Node init, com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node incr, com.google.javascript.rhino.Node body)
func IR_ForNode(init Node, cond Node, incr Node, body Node) Node {
	_ = init
	_ = cond
	_ = incr
	_ = body
	if !(((init.IsVar() || init.IsLet()) || init.IsConst()) || IR_MayBeExpressionOrEmpty(init)) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpressionOrEmpty(cond)) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpressionOrEmpty(incr)) {
		j2gort.CheckState(false)
	}
	if !(body.IsBlock()) {
		j2gort.CheckState(false)
	}
	var r Node = NewNode_4(StaticGet_Token_FOR(), init, cond, incr)
	_ = r
	r.AddChildToBack(body)
	return r
}

// Java static method: IR.switchNode(com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node[] cases)
func IR_SwitchNode(cond Node, cases []Node) Node {
	_ = cond
	_ = cases
	if !(IR_MayBeExpression(cond)) {
		j2gort.CheckState(false)
	}
	var switchNode Node = NewNode_2(StaticGet_Token_SWITCH(), cond)
	_ = switchNode
	var switchBody Node = NewNode_1(StaticGet_Token_SWITCH_BODY())
	_ = switchBody
	switchNode.AddChildToBack(switchBody)
	for _, el_7 := range j2gort.CheckNotNullSlice(cases) {
		caseNode := el_7
		_ = caseNode
		if !(caseNode.IsCase() || caseNode.IsDefaultCase()) {
			j2gort.CheckState(false)
		}
		switchBody.AddChildToBack(caseNode)
	}
	return switchNode
}

// Java static method: IR.caseNode(com.google.javascript.rhino.Node expr, com.google.javascript.rhino.Node body)
func IR_CaseNode(expr Node, body Node) Node {
	_ = expr
	_ = body
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckState(false)
	}
	if !(body.IsBlock()) {
		j2gort.CheckState(false)
	}
	body.SetIsAddedBlock(true)
	return NewNode_3(StaticGet_Token_CASE(), expr, body)
}

// Java static method: IR.defaultCase(com.google.javascript.rhino.Node body)
func IR_DefaultCase(body Node) Node {
	_ = body
	if !(body.IsBlock()) {
		j2gort.CheckState(false)
	}
	body.SetIsAddedBlock(true)
	return NewNode_2(StaticGet_Token_DEFAULT_CASE(), body)
}

// Java static method: IR.label(com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node stmt)
func IR_Label(name Node, stmt Node) Node {
	_ = name
	_ = stmt
	if !(name.IsLabelName()) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeStatement(stmt)) {
		j2gort.CheckState(false)
	}
	return NewNode_3(StaticGet_Token_LABEL(), name, stmt)
}

// Java static method: IR.labelName(java.lang.String name)
func IR_LabelName(name string) Node {
	_ = name
	if !(!(j2gort.CheckNotNullString(name) == "")) {
		j2gort.CheckState(false)
	}
	return Node_NewString_2(StaticGet_Token_LABEL_NAME(), name)
}

// Java static method: IR.tryFinally(com.google.javascript.rhino.Node tryBody, com.google.javascript.rhino.Node finallyBody)
func IR_TryFinally(tryBody Node, finallyBody Node) Node {
	_ = tryBody
	_ = finallyBody
	if !(tryBody.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(finallyBody.IsBlock()) {
		j2gort.CheckState(false)
	}
	var catchBody Node = IR_Block_0().SrcrefIfMissing(tryBody)
	_ = catchBody
	return NewNode_4(StaticGet_Token_TRY(), tryBody, catchBody, finallyBody)
}

// Java static method: IR.tryCatch(com.google.javascript.rhino.Node tryBody, com.google.javascript.rhino.Node catchNode)
func IR_TryCatch(tryBody Node, catchNode Node) Node {
	_ = tryBody
	_ = catchNode
	if !(tryBody.IsBlock()) {
		j2gort.CheckState(false)
	}
	if !(catchNode.IsCatch()) {
		j2gort.CheckState(false)
	}
	var catchBody Node = IR_BlockUnchecked(catchNode).SrcrefIfMissing(catchNode)
	_ = catchBody
	return NewNode_3(StaticGet_Token_TRY(), tryBody, catchBody)
}

// Java static method: IR.tryCatchFinally(com.google.javascript.rhino.Node tryBody, com.google.javascript.rhino.Node catchNode, com.google.javascript.rhino.Node finallyBody)
func IR_TryCatchFinally(tryBody Node, catchNode Node, finallyBody Node) Node {
	_ = tryBody
	_ = catchNode
	_ = finallyBody
	if !(finallyBody.IsBlock()) {
		j2gort.CheckState(false)
	}
	var tryNode Node = IR_TryCatch(tryBody, catchNode)
	_ = tryNode
	tryNode.AddChildToBack(finallyBody)
	return tryNode
}

// Java static method: IR.catchNode(com.google.javascript.rhino.Node expr, com.google.javascript.rhino.Node body)
func IR_CatchNode(expr Node, body Node) Node {
	_ = expr
	_ = body
	if !(expr.IsName()) {
		j2gort.CheckState(false)
	}
	if !(body.IsBlock()) {
		j2gort.CheckState(false)
	}
	return NewNode_3(StaticGet_Token_CATCH(), expr, body)
}

// Java static method: IR.breakNode()
func IR_BreakNode_0() Node {
	return NewNode_1(StaticGet_Token_BREAK())
}

// Java static method: IR.breakNode(com.google.javascript.rhino.Node name)
func IR_BreakNode_1(name Node) Node {
	_ = name
	if !(name.IsLabelName()) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_BREAK(), name)
}

// Java static method: IR.continueNode()
func IR_ContinueNode_0() Node {
	return NewNode_1(StaticGet_Token_CONTINUE())
}

// Java static method: IR.continueNode(com.google.javascript.rhino.Node name)
func IR_ContinueNode_1(name Node) Node {
	_ = name
	if !(name.IsLabelName()) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_CONTINUE(), name)
}

// Java static method: IR.call(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node[] args)
func IR_Call(target Node, args []Node) Node {
	_ = target
	_ = args
	var call Node = NewNode_2(StaticGet_Token_CALL(), target)
	_ = call
	for _, el_8 := range j2gort.CheckNotNullSlice(args) {
		arg := el_8
		_ = arg
		if !(IR_MayBeExpression(arg) || arg.IsSpread()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(arg))
		}
		call.AddChildToBack(arg)
	}
	return call
}

// Java static method: IR.startOptChainCall(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node[] args)
func IR_StartOptChainCall(target Node, args []Node) Node {
	_ = target
	_ = args
	var call Node = NewNode_2(StaticGet_Token_OPTCHAIN_CALL(), target)
	_ = call
	for _, el_9 := range j2gort.CheckNotNullSlice(args) {
		arg := el_9
		_ = arg
		if !(IR_MayBeExpression(arg) || arg.IsSpread()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(arg))
		}
		call.AddChildToBack(arg)
	}
	call.SetIsOptionalChainStart(true)
	return call
}

// Java static method: IR.continueOptChainCall(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node[] args)
func IR_ContinueOptChainCall(target Node, args []Node) Node {
	_ = target
	_ = args
	var call Node = NewNode_2(StaticGet_Token_OPTCHAIN_CALL(), target)
	_ = call
	for _, el_10 := range j2gort.CheckNotNullSlice(args) {
		arg := el_10
		_ = arg
		if !(IR_MayBeExpression(arg) || arg.IsSpread()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(arg))
		}
		call.AddChildToBack(arg)
	}
	call.SetIsOptionalChainStart(false)
	return call
}

// Java static method: IR.newNode(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node[] args)
func IR_NewNode(target Node, args []Node) Node {
	_ = target
	_ = args
	var newcall Node = NewNode_2(StaticGet_Token_NEW(), target)
	_ = newcall
	for _, el_11 := range j2gort.CheckNotNullSlice(args) {
		arg := el_11
		_ = arg
		if !(IR_MayBeExpression(arg) || arg.IsSpread()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(arg))
		}
		newcall.AddChildToBack(arg)
	}
	return newcall
}

// Java static method: IR.name(java.lang.String name)
func IR_Name(name string) Node {
	_ = name
	if !(j2gort.JStrIndexOfCp(name, int32(uint16(46))) == int32(-1)) {
		j2gort.CheckStateFmt(false, "Invalid name '%s'. Did you mean to use NodeUtil.newQName?", j2gort.BoxString(name))
	}
	return Node_NewString_2(StaticGet_Token_NAME(), name)
}

// Java static method: IR.startOptChainGetprop(com.google.javascript.rhino.Node target, java.lang.String prop)
func IR_StartOptChainGetprop(target Node, prop string) Node {
	_ = target
	_ = prop
	if !(IR_MayBeExpression(target)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(target))
	}
	var optChainGetProp Node = Node_NewString_2(StaticGet_Token_OPTCHAIN_GETPROP(), prop)
	_ = optChainGetProp
	optChainGetProp.AddChildToBack(target)
	optChainGetProp.SetIsOptionalChainStart(true)
	return optChainGetProp
}

// Java static method: IR.continueOptChainGetprop(com.google.javascript.rhino.Node target, java.lang.String prop)
func IR_ContinueOptChainGetprop(target Node, prop string) Node {
	_ = target
	_ = prop
	if !(IR_MayBeExpression(target)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(target))
	}
	var optChainGetProp Node = Node_NewString_2(StaticGet_Token_OPTCHAIN_GETPROP(), prop)
	_ = optChainGetProp
	optChainGetProp.AddChildToBack(target)
	optChainGetProp.SetIsOptionalChainStart(false)
	return optChainGetProp
}

// Java static method: IR.getprop(com.google.javascript.rhino.Node target, java.lang.String prop)
func IR_Getprop_2(target Node, prop string) Node {
	_ = target
	_ = prop
	if !(IR_MayBeExpression(target)) {
		j2gort.CheckState(false)
	}
	var getprop Node = Node_NewString_2(StaticGet_Token_GETPROP(), prop)
	_ = getprop
	getprop.AddChildToBack(target)
	return getprop
}

// Java static method: IR.getprop(com.google.javascript.rhino.Node target, java.lang.String prop, java.lang.String[] moreProps)
func IR_Getprop_3(target Node, prop string, moreProps []string) Node {
	_ = target
	_ = prop
	_ = moreProps
	if !(IR_MayBeExpression(target)) {
		j2gort.CheckState(false)
	}
	var result Node = IR_Getprop_2(target, prop)
	_ = result
	for _, el_12 := range j2gort.CheckNotNullSlice(moreProps) {
		moreProp := el_12
		_ = moreProp
		result = IR_Getprop_2(result, moreProp)
	}
	return result
}

// Java static method: IR.startOptChainGetelem(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node elem)
func IR_StartOptChainGetelem(target Node, elem Node) Node {
	_ = target
	_ = elem
	if !(IR_MayBeExpression(target)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(target))
	}
	if !(IR_MayBeExpression(elem)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(elem))
	}
	var optChainGetElem Node = NewNode_3(StaticGet_Token_OPTCHAIN_GETELEM(), target, elem)
	_ = optChainGetElem
	optChainGetElem.SetIsOptionalChainStart(true)
	return optChainGetElem
}

// Java static method: IR.continueOptChainGetelem(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node elem)
func IR_ContinueOptChainGetelem(target Node, elem Node) Node {
	_ = target
	_ = elem
	if !(IR_MayBeExpression(target)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(target))
	}
	if !(IR_MayBeExpression(elem)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(elem))
	}
	var optChainGetElem Node = NewNode_3(StaticGet_Token_OPTCHAIN_GETELEM(), target, elem)
	_ = optChainGetElem
	optChainGetElem.SetIsOptionalChainStart(false)
	return optChainGetElem
}

// Java static method: IR.getelem(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node elem)
func IR_Getelem(target Node, elem Node) Node {
	_ = target
	_ = elem
	if !(IR_MayBeExpression(target)) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpression(elem)) {
		j2gort.CheckState(false)
	}
	return NewNode_3(StaticGet_Token_GETELEM(), target, elem)
}

// Java static method: IR.delprop(com.google.javascript.rhino.Node target)
func IR_Delprop(target Node) Node {
	_ = target
	if !(IR_MayBeExpression(target)) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_DELPROP(), target)
}

// Java static method: IR.assign(com.google.javascript.rhino.Node target, com.google.javascript.rhino.Node expr)
func IR_Assign(target Node, expr Node) Node {
	_ = target
	_ = expr
	if !(target.IsValidAssignmentTarget()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(target))
	}
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(expr))
	}
	return NewNode_3(StaticGet_Token_ASSIGN(), target, expr)
}

// Java static method: IR.hook(com.google.javascript.rhino.Node cond, com.google.javascript.rhino.Node trueval, com.google.javascript.rhino.Node falseval)
func IR_Hook(cond Node, trueval Node, falseval Node) Node {
	_ = cond
	_ = trueval
	_ = falseval
	if !(IR_MayBeExpression(cond)) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpression(trueval)) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpression(falseval)) {
		j2gort.CheckState(false)
	}
	return NewNode_4(StaticGet_Token_HOOK(), cond, trueval, falseval)
}

// Java static method: IR.in(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_In(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_IN(), expr1, expr2)
}

// Java static method: IR.comma(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Comma(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_COMMA(), expr1, expr2)
}

// Java static method: IR.and(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_And(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_AND(), expr1, expr2)
}

// Java static method: IR.or(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Or(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_OR(), expr1, expr2)
}

// Java static method: IR.coalesce(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Coalesce(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_COALESCE(), expr1, expr2)
}

// Java static method: IR.not(com.google.javascript.rhino.Node expr1)
func IR_Not(expr1 Node) Node {
	_ = expr1
	return IR_UnaryOp(StaticGet_Token_NOT(), expr1)
}

// Java static method: IR.lt(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Lt(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_LT(), expr1, expr2)
}

// Java static method: IR.ge(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Ge(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_GE(), expr1, expr2)
}

// Java static method: IR.eq(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Eq(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_EQ(), expr1, expr2)
}

// Java static method: IR.ne(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Ne(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_NE(), expr1, expr2)
}

// Java static method: IR.sheq(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Sheq(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_SHEQ(), expr1, expr2)
}

// Java static method: IR.shne(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Shne(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_SHNE(), expr1, expr2)
}

// Java static method: IR.voidNode(com.google.javascript.rhino.Node expr1)
func IR_VoidNode(expr1 Node) Node {
	_ = expr1
	return IR_UnaryOp(StaticGet_Token_VOID(), expr1)
}

// Java static method: IR.neg(com.google.javascript.rhino.Node expr1)
func IR_Neg(expr1 Node) Node {
	_ = expr1
	return IR_UnaryOp(StaticGet_Token_NEG(), expr1)
}

// Java static method: IR.pos(com.google.javascript.rhino.Node expr1)
func IR_Pos(expr1 Node) Node {
	_ = expr1
	return IR_UnaryOp(StaticGet_Token_POS(), expr1)
}

// Java static method: IR.cast(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.JSDocInfo jsdoc)
func IR_Cast(expr1 Node, jsdoc *JSDocInfo) Node {
	_ = expr1
	_ = jsdoc
	var op Node = IR_UnaryOp(StaticGet_Token_CAST(), expr1)
	_ = op
	op.SetJSDocInfo(jsdoc)
	return op
}

// Java static method: IR.inc(com.google.javascript.rhino.Node exp, boolean isPost)
func IR_Inc(exp Node, isPost bool) Node {
	_ = exp
	_ = isPost
	var op Node = IR_UnaryOp(StaticGet_Token_INC(), exp)
	_ = op
	op.PutBooleanProp(StaticGet_Node_INCRDECR_PROP(), isPost)
	return op
}

// Java static method: IR.dec(com.google.javascript.rhino.Node exp, boolean isPost)
func IR_Dec(exp Node, isPost bool) Node {
	_ = exp
	_ = isPost
	var op Node = IR_UnaryOp(StaticGet_Token_DEC(), exp)
	_ = op
	op.PutBooleanProp(StaticGet_Node_INCRDECR_PROP(), isPost)
	return op
}

// Java static method: IR.add(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Add(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_ADD(), expr1, expr2)
}

// Java static method: IR.sub(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_Sub(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_SUB(), expr1, expr2)
}

// Java static method: IR.assignOr(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_AssignOr(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_ASSIGN_OR(), expr1, expr2)
}

// Java static method: IR.assignAnd(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_AssignAnd(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_ASSIGN_AND(), expr1, expr2)
}

// Java static method: IR.assignCoalesce(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_AssignCoalesce(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_ASSIGN_COALESCE(), expr1, expr2)
}

// Java static method: IR.bitwiseAnd(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_BitwiseAnd(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_BITAND(), expr1, expr2)
}

// Java static method: IR.rightShift(com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_RightShift(expr1 Node, expr2 Node) Node {
	_ = expr1
	_ = expr2
	return IR_BinaryOp(StaticGet_Token_RSH(), expr1, expr2)
}

// Java static method: IR.objectlit(com.google.javascript.rhino.Node[] propdefs)
func IR_Objectlit(propdefs []Node) Node {
	_ = propdefs
	var objectlit Node = NewNode_1(StaticGet_Token_OBJECTLIT())
	_ = objectlit
	for _, el_13 := range j2gort.CheckNotNullSlice(propdefs) {
		propdef := el_13
		_ = propdef
		switch propdef.TokenSwitchOrd_() {
		case int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(128) /* Token.OBJECT_SPREAD */, int32(129) /* Token.COMPUTED_PROP */ :
		default:
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected OBJECTLIT child: "+j2gort.StrVal(propdef))))
		}
		objectlit.AddChildToBack(propdef)
	}
	return objectlit
}

// Java static method: IR.objectPattern(com.google.javascript.rhino.Node[] keys)
func IR_ObjectPattern(keys []Node) Node {
	_ = keys
	var objectPattern Node = NewNode_1(StaticGet_Token_OBJECT_PATTERN())
	_ = objectPattern
	for _, el_14 := range j2gort.CheckNotNullSlice(keys) {
		key := el_14
		_ = key
		if !((key.IsStringKey() || key.IsComputedProp()) || key.IsRest()) {
			j2gort.CheckState(false)
		}
		objectPattern.AddChildToBack(key)
	}
	return objectPattern
}

// Java static method: IR.arrayPattern(com.google.javascript.rhino.Node[] keys)
func IR_ArrayPattern(keys []Node) Node {
	_ = keys
	var arrayPattern Node = NewNode_1(StaticGet_Token_ARRAY_PATTERN())
	_ = arrayPattern
	for _, el_15 := range j2gort.CheckNotNullSlice(keys) {
		key := el_15
		_ = key
		if !(key.IsRest() || key.IsValidAssignmentTarget()) {
			j2gort.CheckState(false)
		}
		arrayPattern.AddChildToBack(key)
	}
	return arrayPattern
}

// Java static method: IR.computedProp(com.google.javascript.rhino.Node key, com.google.javascript.rhino.Node value)
func IR_ComputedProp(key Node, value Node) Node {
	_ = key
	_ = value
	if !(IR_MayBeExpression(key)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(key))
	}
	if !(IR_MayBeExpression(value)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(value))
	}
	return NewNode_3(StaticGet_Token_COMPUTED_PROP(), key, value)
}

// Java static method: IR.propdef(com.google.javascript.rhino.Node string, com.google.javascript.rhino.Node value)
func IR_Propdef(string_ Node, value Node) Node {
	_ = string_
	_ = value
	if !(string_.IsStringKey()) {
		j2gort.CheckState(false)
	}
	if !(!(string_.HasChildren())) {
		j2gort.CheckState(false)
	}
	if !(IR_MayBeExpression(value)) {
		j2gort.CheckState(false)
	}
	string_.AddChildToFront(value)
	return string_
}

// Java static method: IR.arraylit(com.google.javascript.rhino.Node[] exprs)
func IR_Arraylit_1_Node(exprs []Node) Node {
	_ = exprs
	return IR_Arraylit_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Arrays_AsList(func() []any {
		arr_16 := exprs
		if arr_16 == nil {
			return nil
		}
		out_17 := make([]any, len(arr_16))
		for i_, el_18 := range arr_16 {
			out_17[i_] = j2gort.BoxNode(el_18)
		}
		return out_17
	}())))
}

// Java static method: IR.arraylit(java.lang.Iterable<com.google.javascript.rhino.Node> exprs)
func IR_Arraylit_1_Ext_java_lang_Iterable(exprs j2gort.Ext_java_lang_Iterable) Node {
	_ = exprs
	var arraylit Node = NewNode_1(StaticGet_Token_ARRAYLIT())
	_ = arraylit
	{
		cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNull(exprs))
		for cur_19.HasNext() {
			expr := Node(cur_19.NextNode())
			_ = expr
			if !(IR_MayBeExpressionOrEmpty(expr) || expr.IsSpread()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(expr))
			}
			arraylit.AddChildToBack(expr)
		}
	}
	return arraylit
}

// Java static method: IR.regexp(com.google.javascript.rhino.Node expr)
func IR_Regexp_1(expr Node) Node {
	_ = expr
	if !(expr.IsStringLit()) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_REGEXP(), expr)
}

// Java static method: IR.regexp(com.google.javascript.rhino.Node expr, com.google.javascript.rhino.Node flags)
func IR_Regexp_2(expr Node, flags Node) Node {
	_ = expr
	_ = flags
	if !(expr.IsStringLit()) {
		j2gort.CheckState(false)
	}
	if !(flags.IsStringLit()) {
		j2gort.CheckState(false)
	}
	return NewNode_3(StaticGet_Token_REGEXP(), expr, flags)
}

// Java static method: IR.string(java.lang.String s)
func IR_String(s string) Node {
	_ = s
	return Node_NewString_1(s)
}

// Java static method: IR.stringKey(java.lang.String s)
func IR_StringKey_1(s string) Node {
	_ = s
	return Node_NewString_2(StaticGet_Token_STRING_KEY(), s)
}

// Java static method: IR.stringKey(java.lang.String s, com.google.javascript.rhino.Node value)
func IR_StringKey_2(s string, value Node) Node {
	_ = s
	_ = value
	if !((IR_MayBeExpression(value) || value.IsDefaultValue()) || value.IsObjectPattern()) {
		j2gort.CheckState(false)
	}
	var stringKey Node = IR_StringKey_1(s)
	_ = stringKey
	stringKey.AddChildToFront(value)
	return stringKey
}

// Java static method: IR.quotedStringKey(java.lang.String s, com.google.javascript.rhino.Node value)
func IR_QuotedStringKey(s string, value Node) Node {
	_ = s
	_ = value
	var k Node = IR_StringKey_2(s, value)
	_ = k
	k.PutBooleanProp(StaticGet_Node_QUOTED_PROP(), true)
	return k
}

// Java static method: IR.templateLiteral()
func IR_TemplateLiteral() Node {
	return NewNode_1(StaticGet_Token_TEMPLATELIT())
}

// Java static method: IR.templateLiteralString(java.lang.@org.jspecify.annotations.Nullable String cooked, java.lang.String raw)
func IR_TemplateLiteralString(cooked string, raw string) Node {
	_ = cooked
	_ = raw
	return Node_NewTemplateLitString_2(cooked, raw)
}

// Java static method: IR.templateLiteralSubstitution(com.google.javascript.rhino.Node child)
func IR_TemplateLiteralSubstitution(child Node) Node {
	_ = child
	var sub Node = NewNode_1(StaticGet_Token_TEMPLATELIT_SUB())
	_ = sub
	sub.AddChildToBack(child)
	return sub
}

// Java static method: IR.iterRest(com.google.javascript.rhino.Node target)
func IR_IterRest(target Node) Node {
	_ = target
	if !(target.IsValidAssignmentTarget()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(target))
	}
	return NewNode_2(StaticGet_Token_ITER_REST(), target)
}

// Java static method: IR.objectRest(com.google.javascript.rhino.Node target)
func IR_ObjectRest(target Node) Node {
	_ = target
	if !(target.IsValidAssignmentTarget()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(target))
	}
	return NewNode_2(StaticGet_Token_OBJECT_REST(), target)
}

// Java static method: IR.iterSpread(com.google.javascript.rhino.Node expr)
func IR_IterSpread(expr Node) Node {
	_ = expr
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_ITER_SPREAD(), expr)
}

// Java static method: IR.objectSpread(com.google.javascript.rhino.Node expr)
func IR_ObjectSpread(expr Node) Node {
	_ = expr
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckState(false)
	}
	return NewNode_2(StaticGet_Token_OBJECT_SPREAD(), expr)
}

// Java static method: IR.superNode()
func IR_SuperNode() Node {
	return NewNode_1(StaticGet_Token_SUPER())
}

// Java static method: IR.getterDef(java.lang.String name, com.google.javascript.rhino.Node value)
func IR_GetterDef(name string, value Node) Node {
	_ = name
	_ = value
	if !(value.IsFunction()) {
		j2gort.CheckState(false)
	}
	var member Node = Node_NewString_2(StaticGet_Token_GETTER_DEF(), name)
	_ = member
	member.AddChildToFront(value)
	return member
}

// Java static method: IR.setterDef(java.lang.String name, com.google.javascript.rhino.Node value)
func IR_SetterDef(name string, value Node) Node {
	_ = name
	_ = value
	if !(value.IsFunction()) {
		j2gort.CheckState(false)
	}
	var member Node = Node_NewString_2(StaticGet_Token_SETTER_DEF(), name)
	_ = member
	member.AddChildToFront(value)
	return member
}

// Java static method: IR.memberFieldDef(java.lang.String name, com.google.javascript.rhino.Node value)
func IR_MemberFieldDef(name string, value Node) Node {
	_ = name
	_ = value
	if !(IR_MayBeExpression(value)) {
		j2gort.CheckState(false)
	}
	var member Node = Node_NewString_2(StaticGet_Token_MEMBER_FIELD_DEF(), name)
	_ = member
	member.AddChildToFront(value)
	return member
}

// Java static method: IR.memberFunctionDef(java.lang.String name, com.google.javascript.rhino.Node function)
func IR_MemberFunctionDef(name string, function Node) Node {
	_ = name
	_ = function
	if !(function.IsFunction()) {
		j2gort.CheckState(false)
	}
	var member Node = Node_NewString_2(StaticGet_Token_MEMBER_FUNCTION_DEF(), name)
	_ = member
	member.AddChildToBack(function)
	return member
}

// Java static method: IR.number(double d)
func IR_Number(d float64) Node {
	_ = d
	if !(!(math.IsNaN(d))) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfDouble(d))
	}
	if !(pkg_base.JSCompDoubles_IsPositive(d)) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfDouble(d))
	}
	return Node_NewNumber(d)
}

// Java static method: IR.bigint(java.math.BigInteger b)
func IR_Bigint(b *j2gort.Ext_java_math_BigInteger) Node {
	_ = b
	j2gort.CheckNotNullPtr(b)
	if !(b.Signum() >= int32(0)) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(b))
	}
	return Node_NewBigInt(b)
}

// Java static method: IR.thisNode()
func IR_ThisNode() Node {
	return NewNode_1(StaticGet_Token_THIS())
}

// Java static method: IR.trueNode()
func IR_TrueNode() Node {
	return NewNode_1(StaticGet_Token_TRUE())
}

// Java static method: IR.falseNode()
func IR_FalseNode() Node {
	return NewNode_1(StaticGet_Token_FALSE())
}

// Java static method: IR.nullNode()
func IR_NullNode() Node {
	return NewNode_1(StaticGet_Token_NULL())
}

// Java static method: IR.typeof(com.google.javascript.rhino.Node expr)
func IR_Typeof(expr Node) Node {
	_ = expr
	return IR_UnaryOp(StaticGet_Token_TYPEOF(), expr)
}

// Java static method: IR.importMeta()
func IR_ImportMeta() Node {
	return NewNode_1(StaticGet_Token_IMPORT_META())
}

// Java static method: IR.newTarget()
func IR_NewTarget() Node {
	return NewNode_1(StaticGet_Token_NEW_TARGET())
}

// Java static method: IR.binaryOp(com.google.javascript.rhino.Token token, com.google.javascript.rhino.Node expr1, com.google.javascript.rhino.Node expr2)
func IR_BinaryOp(token Token, expr1 Node, expr2 Node) Node {
	_ = token
	_ = expr1
	_ = expr2
	if !(IR_MayBeExpression(expr1)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(expr1))
	}
	if !(IR_MayBeExpression(expr2)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(expr2))
	}
	return NewNode_3(token, expr1, expr2)
}

// Java static method: IR.unaryOp(com.google.javascript.rhino.Token token, com.google.javascript.rhino.Node expr)
func IR_UnaryOp(token Token, expr Node) Node {
	_ = token
	_ = expr
	if !(IR_MayBeExpression(expr)) {
		j2gort.CheckState(false)
	}
	return NewNode_2(token, expr)
}

// Java static method: IR.mayBeExpressionOrEmpty(com.google.javascript.rhino.Node n)
func IR_MayBeExpressionOrEmpty(n Node) bool {
	_ = n
	return (n.IsEmpty() || IR_MayBeExpression(n))
}

// Java static method: IR.mayBeStatementNoReturn(com.google.javascript.rhino.Node n)
func IR_MayBeStatementNoReturn(n Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(88) /* Token.EMPTY */, int32(73) /* Token.FUNCTION */ :
		return true
	case int32(90) /* Token.BLOCK */, int32(82) /* Token.BREAK */, int32(105) /* Token.CLASS */, int32(97) /* Token.CONST */, int32(83) /* Token.CONTINUE */, int32(98) /* Token.DEBUGGER */, int32(79) /* Token.DO */, int32(167) /* Token.ENUM */, int32(93) /* Token.EXPR_RESULT */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(74) /* Token.IF */, int32(164) /* Token.INTERFACE */, int32(92) /* Token.LABEL */, int32(111) /* Token.LET */, int32(75) /* Token.SWITCH */, int32(43) /* Token.THROW */, int32(48) /* Token.TRY */, int32(84) /* Token.VAR */, int32(78) /* Token.WHILE */, int32(85) /* Token.WITH */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: IR.mayBeStatement(com.google.javascript.rhino.Node n)
func IR_MayBeStatement(n Node) bool {
	_ = n
	if !(IR_MayBeStatementNoReturn(n)) {
		return n.IsReturn()
	}
	return true
}

// Java static method: IR.mayBeExpression(com.google.javascript.rhino.Node n)
func IR_MayBeExpression(n Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */ :
		return true
	case int32(13) /* Token.ADD */, int32(69) /* Token.AND */, int32(46) /* Token.ARRAYLIT */, int32(51) /* Token.ASSIGN */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(63) /* Token.ASSIGN_EXPONENT */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */, int32(115) /* Token.AWAIT */, int32(34) /* Token.BIGINT */, int32(3) /* Token.BITAND */, int32(1) /* Token.BITOR */, int32(20) /* Token.BITNOT */, int32(2) /* Token.BITXOR */, int32(28) /* Token.CALL */, int32(101) /* Token.CAST */, int32(70) /* Token.COALESCE */, int32(50) /* Token.COMMA */, int32(72) /* Token.DEC */, int32(24) /* Token.DELPROP */, int32(16) /* Token.DIV */, int32(124) /* Token.DYNAMIC_IMPORT */, int32(4) /* Token.EQ */, int32(18) /* Token.EXPONENT */, int32(38) /* Token.FALSE */, int32(9) /* Token.GE */, int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */, int32(8) /* Token.GT */, int32(67) /* Token.HOOK */, int32(136) /* Token.IMPORT_META */, int32(44) /* Token.IN */, int32(71) /* Token.INC */, int32(45) /* Token.INSTANCEOF */, int32(7) /* Token.LE */, int32(10) /* Token.LSH */, int32(6) /* Token.LT */, int32(17) /* Token.MOD */, int32(15) /* Token.MUL */, int32(32) /* Token.NAME */, int32(5) /* Token.NE */, int32(22) /* Token.NEG */, int32(23) /* Token.NEW */, int32(135) /* Token.NEW_TARGET */, int32(19) /* Token.NOT */, int32(33) /* Token.NUMBER */, int32(36) /* Token.NULL */, int32(47) /* Token.OBJECTLIT */, int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(68) /* Token.OR */, int32(21) /* Token.POS */, int32(42) /* Token.REGEXP */, int32(11) /* Token.RSH */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(35) /* Token.STRINGLIT */, int32(14) /* Token.SUB */, int32(110) /* Token.SUPER */, int32(131) /* Token.TEMPLATELIT */, int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(37) /* Token.THIS */, int32(25) /* Token.TYPEOF */, int32(39) /* Token.TRUE */, int32(12) /* Token.URSH */, int32(87) /* Token.VOID */, int32(114) /* Token.YIELD */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}
