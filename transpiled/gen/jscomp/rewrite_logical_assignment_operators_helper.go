/*
 * Copyright 2021 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/RewriteLogicalAssignmentOperatorsHelper.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteLogicalAssignmentOperatorsHelper
type RewriteLogicalAssignmentOperatorsHelper struct {
	compiler AbstractCompiler `gson:"compiler"`
	astFactory *AstFactory `gson:"astFactory"`
	uniqueIdSupplier *UniqueIdSupplier `gson:"uniqueIdSupplier"`
	idHash_ atomic.Int32
}

func (this_ *RewriteLogicalAssignmentOperatorsHelper) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteLogicalAssignmentOperatorsHelper) IsJavaNil_() bool { return this_ == nil }
func (*RewriteLogicalAssignmentOperatorsHelper) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.RewriteLogicalAssignmentOperatorsHelper", "RewriteLogicalAssignmentOperatorsHelper" }
func (this_ *RewriteLogicalAssignmentOperatorsHelper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteLogicalAssignmentOperatorsHelper) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteLogicalAssignmentOperatorsHelper", "com.google.javascript.jscomp.RewriteLogicalAssignmentOperatorsHelper", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteLogicalAssignmentOperatorsHelper](v_)
	})
}

func StaticGet_RewriteLogicalAssignmentOperatorsHelper_TEMP_VAR_NAME_PREFIX() string { return "$jscomp$logical$assign$tmp" }
func StaticGet_RewriteLogicalAssignmentOperatorsHelper_TEMP_INDEX_VAR_NAME_PREFIX() string { return "$jscomp$logical$assign$tmpindex" }
func (this_ *RewriteLogicalAssignmentOperatorsHelper) initFields_RewriteLogicalAssignmentOperatorsHelper() {
	_ = this_
}

// Java constructor: RewriteLogicalAssignmentOperatorsHelper(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.AstFactory astFactory, com.google.javascript.jscomp.UniqueIdSupplier uniqueIdSupplier)
func NewRewriteLogicalAssignmentOperatorsHelper(compiler AbstractCompiler, astFactory *AstFactory, uniqueIdSupplier *UniqueIdSupplier) *RewriteLogicalAssignmentOperatorsHelper {
	this_ := &RewriteLogicalAssignmentOperatorsHelper{}
	InitRewriteLogicalAssignmentOperatorsHelper(this_, compiler, astFactory, uniqueIdSupplier)
	return this_
}

func InitRewriteLogicalAssignmentOperatorsHelper(this_ *RewriteLogicalAssignmentOperatorsHelper, compiler AbstractCompiler, astFactory *AstFactory, uniqueIdSupplier *UniqueIdSupplier) {
	_ = compiler
	_ = astFactory
	_ = uniqueIdSupplier
	this_.initFields_RewriteLogicalAssignmentOperatorsHelper()
	this_.compiler = compiler
	this_.astFactory = astFactory
	this_.uniqueIdSupplier = uniqueIdSupplier
}

// Java method: RewriteLogicalAssignmentOperatorsHelper.visitLogicalAssignmentOperator(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node logicalAssignment)
func (this_ *RewriteLogicalAssignmentOperatorsHelper) VisitLogicalAssignmentOperator(t *NodeTraversal, logicalAssignment pkg_rhino.Node) {
	_ = t
	_ = logicalAssignment
	var enclosingFunction pkg_rhino.Node = NodeUtil_GetEnclosingFunction(logicalAssignment)
	_ = enclosingFunction
	if (((enclosingFunction != 0) && !(NodeUtil_GetFunctionBody(enclosingFunction).IsBlock()))) {
		var returnValue pkg_rhino.Node = NodeUtil_GetFunctionBody(enclosingFunction)
		_ = returnValue
		var body pkg_rhino.Node = pkg_rhino.IR_Block_1_Node(pkg_rhino.IR_ReturnNode_1(returnValue.Detach())).SrcrefTreeIfMissing(returnValue)
		_ = body
		enclosingFunction.AddChildToBack(body)
		this_.compiler.ReportChangeToEnclosingScope(body)
	}
	var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(logicalAssignment)
	_ = enclosingStatement
	for (enclosingStatement.GetParent().IsLabel()) {
		enclosingStatement = enclosingStatement.GetParent()
	}
	var left pkg_rhino.Node = logicalAssignment.RemoveFirstChild()
	_ = left
	var right pkg_rhino.Node = logicalAssignment.GetLastChild().Detach()
	_ = right
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	for (left.IsCast()) {
		left = left.RemoveFirstChild()
	}
	if (left.IsName()) {
		replacement = this_.HandleLHSName(logicalAssignment, left, right)
	} else {
		if !((left.IsGetProp() || left.IsGetElem())) { j2gort.CheckStateMsg(false, j2gort.BoxNode(left)) }
		replacement = this_.HandleLHSPropertyReference(t, logicalAssignment, left, right, enclosingStatement)
	}
	logicalAssignment.ReplaceWith(replacement)
	this_.compiler.ReportChangeToEnclosingScope(enclosingStatement)
}

// Java method: RewriteLogicalAssignmentOperatorsHelper.handleLHSName(com.google.javascript.rhino.Node logicalAssignment, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right)
func (this_ *RewriteLogicalAssignmentOperatorsHelper) HandleLHSName(logicalAssignment pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node) pkg_rhino.Node {
	_ = logicalAssignment
	_ = left
	_ = right
	var assignToRHS pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(left, right).Srcref(right)
	_ = assignToRHS
	return pkg_rhino.NewNode_3(NodeUtil_GetOpFromAssignmentOp(logicalAssignment), left.CloneNode_0(), assignToRHS).CopyTypeFrom(logicalAssignment).Srcref(logicalAssignment)
}

// Java method: RewriteLogicalAssignmentOperatorsHelper.handleLHSPropertyReference(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node logicalAssignment, com.google.javascript.rhino.Node left, com.google.javascript.rhino.Node right, com.google.javascript.rhino.Node enclosingStatement)
func (this_ *RewriteLogicalAssignmentOperatorsHelper) HandleLHSPropertyReference(t *NodeTraversal, logicalAssignment pkg_rhino.Node, left pkg_rhino.Node, right pkg_rhino.Node, enclosingStatement pkg_rhino.Node) pkg_rhino.Node {
	_ = t
	_ = logicalAssignment
	_ = left
	_ = right
	_ = enclosingStatement
	var input *CompilerInput = t.GetInput()
	_ = input
	var uniqueId string = this_.uniqueIdSupplier.GetUniqueId(input)
	_ = uniqueId
	var tempVarName string = ("$jscomp$logical$assign$tmp" + j2gort.StrValString(uniqueId))
	_ = tempVarName
	var assignToRHS pkg_rhino.Node = pkg_rhino.Node(0)
	_ = assignToRHS
	var newLHS pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newLHS
	var objectNode pkg_rhino.Node = left.RemoveFirstChild()
	_ = objectNode
	var let pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_1(tempVarName).SrcrefTree(logicalAssignment)
	_ = let
	let.InsertBefore(enclosingStatement)
	var tempName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(tempVarName, AstFactory_S_Type_1_Node(objectNode)).Srcref(objectNode)
	_ = tempName
	var assignTemp pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(tempName, objectNode).Srcref(objectNode)
	_ = assignTemp
	if (left.IsGetProp()) {
		var propertyName string = left.GetString()
		_ = propertyName
		var tempProp pkg_rhino.Node = this_.astFactory.CreateGetProp(tempName.CloneNode_0(), propertyName, AstFactory_S_Type_1_Node(left)).Srcref(right)
		_ = tempProp
		assignToRHS = this_.astFactory.CreateAssign_2_Node_Node(tempProp, right).Srcref(right)
		newLHS = this_.astFactory.CreateGetProp(assignTemp, propertyName, AstFactory_S_Type_1_Node(left)).Srcref(left)
	} else {
		if !(left.IsGetElem()) { j2gort.CheckStateMsg(false, j2gort.BoxNode(left)) }
		var tempIndexVarName string = ("$jscomp$logical$assign$tmpindex" + j2gort.StrValString(uniqueId))
		_ = tempIndexVarName
		var indexExprNode pkg_rhino.Node = left.GetLastChild().Detach()
		_ = indexExprNode
		var letIndex pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_1(tempIndexVarName).SrcrefTree(logicalAssignment)
		_ = letIndex
		letIndex.InsertBefore(enclosingStatement)
		var tempIndexName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(tempIndexVarName, AstFactory_S_Type_1_Node(indexExprNode)).Srcref(indexExprNode)
		_ = tempIndexName
		var assignTempIndex pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(tempIndexName, indexExprNode).Srcref(indexExprNode)
		_ = assignTempIndex
		var tempElem pkg_rhino.Node = this_.astFactory.CreateGetElem(tempName.CloneNode_0(), tempIndexName.CloneNode_0()).CopyTypeFrom(left).Srcref(right)
		_ = tempElem
		assignToRHS = this_.astFactory.CreateAssign_2_Node_Node(tempElem, right).Srcref(right)
		newLHS = this_.astFactory.CreateGetElem(assignTemp, assignTempIndex).CopyTypeFrom(left).Srcref(left)
	}
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
	return pkg_rhino.NewNode_3(NodeUtil_GetOpFromAssignmentOp(logicalAssignment), newLHS, assignToRHS).CopyTypeFrom(logicalAssignment).Srcref(logicalAssignment)
}

