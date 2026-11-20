/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/OptionalChainRewriter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.OptionalChainRewriter
type OptionalChainRewriter struct {
	compiler          AbstractCompiler                        `gson:"compiler"`
	astFactory        *AstFactory                             `gson:"astFactory"`
	tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator `gson:"tmpVarNameCreator"`
	scope             *Scope                                  `gson:"scope"`
	deletesToDelete   j2gort.Deque[pkg_rhino.Node]            `gson:"deletesToDelete"`
	chainParent       pkg_rhino.Node                          `gson:"chainParent"`
	wholeChain        pkg_rhino.Node                          `gson:"wholeChain"`
	insertionPoint    pkg_rhino.Node                          `gson:"insertionPoint"`
	idHash_           atomic.Int32
}

func (this_ *OptionalChainRewriter) BindSelf_(s any) {
	_ = s
}
func (this_ *OptionalChainRewriter) IsJavaNil_() bool { return this_ == nil }
func (*OptionalChainRewriter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptionalChainRewriter", "OptionalChainRewriter"
}
func (this_ *OptionalChainRewriter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptionalChainRewriter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptionalChainRewriter", "com.google.javascript.jscomp.OptionalChainRewriter", func(v_ any) bool {
		return j2gort.InstanceOf[*OptionalChainRewriter](v_)
	})
}

func (this_ *OptionalChainRewriter) initFields_OptionalChainRewriter() {
	_ = this_
}

// Java constructor: OptionalChainRewriter(com.google.javascript.jscomp.OptionalChainRewriter.Builder builder, com.google.javascript.rhino.Node wholeChain)
func NewOptionalChainRewriter(builder *OptionalChainRewriter_Builder, wholeChain pkg_rhino.Node) *OptionalChainRewriter {
	this_ := &OptionalChainRewriter{}
	InitOptionalChainRewriter(this_, builder, wholeChain)
	return this_
}

func InitOptionalChainRewriter(this_ *OptionalChainRewriter, builder *OptionalChainRewriter_Builder, wholeChain pkg_rhino.Node) {
	_ = builder
	_ = wholeChain
	this_.initFields_OptionalChainRewriter()
	if !(NodeUtil_IsEndOfFullOptChain(wholeChain)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(wholeChain))
	}
	this_.compiler = builder.compiler
	this_.astFactory = builder.astFactory
	this_.tmpVarNameCreator = j2gort.CheckNotNullIface(builder.tmpVarNameCreator)
	this_.scope = builder.scope
	this_.wholeChain = wholeChain
	this_.chainParent = j2gort.CheckNotNullNode(wholeChain.GetParent())
	var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(wholeChain)
	_ = enclosingStatement
	if (enclosingStatement.GetPrevious() != 0) && enclosingStatement.GetPrevious().IsLabelName() {
		this_.insertionPoint = enclosingStatement.GetParent()
	} else {
		this_.insertionPoint = enclosingStatement
	}
	this_.deletesToDelete = j2gort.Deque[pkg_rhino.Node]{}
}

// Java static method: OptionalChainRewriter.builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func OptionalChainRewriter_S_Builder(compiler AbstractCompiler) *OptionalChainRewriter_Builder {
	_ = compiler
	return NewOptionalChainRewriter_Builder(compiler)
}

// Java method: OptionalChainRewriter.rewrite()
func (this_ *OptionalChainRewriter) Rewrite() {
	if !(NodeUtil_IsOptChainNode(this_.wholeChain)) {
		j2gort.CheckStateFmt(false, "already rewritten: %s", j2gort.BoxNode(this_.wholeChain))
	}
	var startNodeStack j2gort.Deque[pkg_rhino.Node]
	_ = startNodeStack
	var subchainEnd pkg_rhino.Node = this_.wholeChain
	_ = subchainEnd
	for NodeUtil_IsOptChainNode(subchainEnd) {
		var subchainStart pkg_rhino.Node = NodeUtil_GetStartOfOptChainSegment(subchainEnd)
		_ = subchainStart
		startNodeStack.Push(j2gort.CheckNotNullNode(subchainStart))
		subchainEnd = subchainStart.GetFirstChild()
	}
	if !(!(startNodeStack.IsEmpty())) {
		j2gort.CheckState(false)
	}
	var optChainReplacement pkg_rhino.Node = this_.priv_OptionalChainRewriter_rewriteInitialSegment_4152(startNodeStack.Pop(), this_.wholeChain)
	_ = optChainReplacement
	for !(startNodeStack.IsEmpty()) {
		this_.priv_OptionalChainRewriter_rewriteInitialSegment_4152(startNodeStack.Pop(), this_.wholeChain)
	}
	if ((this_.chainParent.IsCall() && !(this_.chainParent.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL()))) && optChainReplacement.IsFirstChildOf(this_.chainParent)) && NodeUtil_IsNormalGet(this_.wholeChain) {
		var thisValue pkg_rhino.Node = this_.wholeChain.GetFirstChild()
		_ = thisValue
		var tmpThisNode pkg_rhino.Node = this_.GetSubExprNameNode(thisValue)
		_ = tmpThisNode
		optChainReplacement.Detach()
		this_.chainParent.AddChildToFront(tmpThisNode)
		var dotCallNode pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(optChainReplacement, "call").SrcrefTreeIfMissing(optChainReplacement)
		_ = dotCallNode
		this_.chainParent.AddChildToFront(dotCallNode)
	}
	this_.compiler.ReportChangeToEnclosingScope(this_.chainParent)
	if this_.chainParent.IsDelProp() {
		for !(this_.deletesToDelete.IsEmpty()) {
			var delete_ pkg_rhino.Node = this_.deletesToDelete.Pop()
			_ = delete_
			if !(delete_.GetFirstChild().IsHook()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(delete_))
			}
			var hook pkg_rhino.Node = delete_.GetFirstChild()
			_ = hook
			hook.Detach()
			delete_.ReplaceWith(hook)
			this_.compiler.ReportChangeToEnclosingScope(hook)
		}
	}
	var enclosingScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(this_.insertionPoint)
	_ = enclosingScript
	NodeUtil_AddFeatureToScript(enclosingScript, pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
}

// Java method: OptionalChainRewriter.rewriteInitialSegment(com.google.javascript.rhino.Node fullChainStart, com.google.javascript.rhino.Node fullChainEnd)
func (this_ *OptionalChainRewriter) priv_OptionalChainRewriter_rewriteInitialSegment_4152(fullChainStart pkg_rhino.Node, fullChainEnd pkg_rhino.Node) pkg_rhino.Node {
	_ = fullChainStart
	_ = fullChainEnd
	var receiverNode pkg_rhino.Node = fullChainStart.GetFirstChild()
	_ = receiverNode
	var initialChainEnd pkg_rhino.Node = NodeUtil_GetEndOfOptChainSegment(fullChainStart)
	_ = initialChainEnd
	var isBeingDeleted bool = fullChainEnd.GetParent().IsDelProp()
	_ = isBeingDeleted
	if isBeingDeleted {
		this_.deletesToDelete.AddLast(j2gort.CheckNotNullNode(fullChainEnd.GetParent()))
	}
	if !(!(NodeUtil_IsOptChainNode(receiverNode))) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(receiverNode))
	}
	NodeUtil_ConvertToNonOptionalChainSegment(initialChainEnd)
	var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = placeholder
	fullChainEnd.ReplaceWith(placeholder)
	if (NodeUtil_IsNormalGet(receiverNode) && fullChainStart.IsCall()) && !(fullChainStart.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL())) {
		var thisValue pkg_rhino.Node = receiverNode.GetFirstChild()
		_ = thisValue
		var tmpThisNode pkg_rhino.Node = this_.GetSubExprNameNode(thisValue)
		_ = tmpThisNode
		var tmpReceiverNode pkg_rhino.Node = this_.GetSubExprNameNode(receiverNode)
		_ = tmpReceiverNode
		receiverNode = fullChainStart.RemoveFirstChild()
		fullChainStart.AddChildToFront(tmpThisNode)
		fullChainStart.AddChildToFront(this_.astFactory.CreateGetPropWithUnknownType(tmpReceiverNode, "call").SrcrefTreeIfMissing(receiverNode))
	} else {
		var tmpReceiverNode_2 pkg_rhino.Node = this_.GetSubExprNameNode(receiverNode)
		_ = tmpReceiverNode_2
		receiverNode = fullChainStart.GetFirstChild()
		receiverNode.ReplaceWith(tmpReceiverNode_2)
	}
	var optChainReplacement pkg_rhino.Node = this_.astFactory.CreateHook(this_.astFactory.CreateEq(receiverNode, this_.astFactory.CreateNull()), func() pkg_rhino.Node {
		if isBeingDeleted {
			return this_.astFactory.CreateBoolean(true)
		}
		return this_.astFactory.CreateUndefinedValue()
	}(), func() pkg_rhino.Node {
		if isBeingDeleted {
			return this_.astFactory.CreateDelProp(fullChainEnd)
		}
		return fullChainEnd
	}()).SrcrefTreeIfMissing(fullChainEnd)
	_ = optChainReplacement
	placeholder.ReplaceWith(optChainReplacement)
	return optChainReplacement
}

// Java method: OptionalChainRewriter.getSubExprNameNode(com.google.javascript.rhino.Node subExpr)
func (this_ *OptionalChainRewriter) GetSubExprNameNode(subExpr pkg_rhino.Node) pkg_rhino.Node {
	_ = subExpr
	var tempVarName string = this_.DeclareTempVarName(subExpr)
	_ = tempVarName
	var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = placeholder
	subExpr.ReplaceWith(placeholder)
	subExpr.RemoveDirectEvalProp()
	var replacement pkg_rhino.Node = this_.astFactory.CreateAssign_2_string_Node(tempVarName, subExpr).SrcrefTreeIfMissing(subExpr)
	_ = replacement
	placeholder.ReplaceWith(replacement)
	return replacement.GetFirstChild().CloneNode_0()
}

// Java method: OptionalChainRewriter.declareTempVarName(com.google.javascript.rhino.Node valueNode)
func (this_ *OptionalChainRewriter) DeclareTempVarName(valueNode pkg_rhino.Node) string {
	_ = valueNode
	var tempVarName string = this_.tmpVarNameCreator.CreateTmpVarName()
	_ = tempVarName
	var declarationStatement pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_1(tempVarName).SrcrefTree(valueNode)
	_ = declarationStatement
	declarationStatement.GetFirstChild().SetInferredConstantVar(true)
	declarationStatement.InsertBefore(this_.insertionPoint)
	this_.compiler.ReportChangeToEnclosingScope(declarationStatement)
	if this_.scope != nil {
		this_.scope.Declare(tempVarName, declarationStatement.GetFirstChild(), nil)
	}
	return tempVarName
}

// Java interface: com.google.javascript.jscomp.OptionalChainRewriter.TmpVarNameCreator
type OptionalChainRewriter_TmpVarNameCreator interface {
	Implements_jscomp_OptionalChainRewriter_TmpVarNameCreator()
	IsJavaNil_() bool
	CreateTmpVarName() string
}

func ToIface_OptionalChainRewriter_TmpVarNameCreator[T interface {
	OptionalChainRewriter_TmpVarNameCreator
	comparable
}](v T) OptionalChainRewriter_TmpVarNameCreator {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_OptionalChainRewriter_TmpVarNameCreator[T OptionalChainRewriter_TmpVarNameCreator](v T) OptionalChainRewriter_TmpVarNameCreator {
	var vi OptionalChainRewriter_TmpVarNameCreator = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptionalChainRewriter.TmpVarNameCreator", "com.google.javascript.jscomp.OptionalChainRewriter$TmpVarNameCreator", func(v_ any) bool {
		return j2gort.InstanceOf[OptionalChainRewriter_TmpVarNameCreator](v_)
	})
}

type lambda_OptionalChainRewriter_TmpVarNameCreator struct {
	fn func() string
}

func (l *lambda_OptionalChainRewriter_TmpVarNameCreator) IsJavaNil_() bool { return l == nil }
func NewLambda_OptionalChainRewriter_TmpVarNameCreator(fn func() string) OptionalChainRewriter_TmpVarNameCreator {
	if fn == nil {
		return nil
	}
	return &lambda_OptionalChainRewriter_TmpVarNameCreator{fn: fn}
}

func (*lambda_OptionalChainRewriter_TmpVarNameCreator) Implements_jscomp_OptionalChainRewriter_TmpVarNameCreator() {
}
func (l *lambda_OptionalChainRewriter_TmpVarNameCreator) CreateTmpVarName() string {
	return l.fn()
}

// Java class: com.google.javascript.jscomp.OptionalChainRewriter.Builder
type OptionalChainRewriter_Builder struct {
	compiler          AbstractCompiler                        `gson:"compiler"`
	astFactory        *AstFactory                             `gson:"astFactory"`
	tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator `gson:"tmpVarNameCreator"`
	scope             *Scope                                  `gson:"scope"`
	idHash_           atomic.Int32
}

func (this_ *OptionalChainRewriter_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *OptionalChainRewriter_Builder) IsJavaNil_() bool { return this_ == nil }
func (*OptionalChainRewriter_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptionalChainRewriter$Builder", "Builder"
}
func (this_ *OptionalChainRewriter_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptionalChainRewriter_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptionalChainRewriter.Builder", "com.google.javascript.jscomp.OptionalChainRewriter$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*OptionalChainRewriter_Builder](v_)
	})
}

func (this_ *OptionalChainRewriter_Builder) initFields_OptionalChainRewriter_Builder() {
	_ = this_
}

// Java constructor: Builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewOptionalChainRewriter_Builder(compiler AbstractCompiler) *OptionalChainRewriter_Builder {
	this_ := &OptionalChainRewriter_Builder{}
	InitOptionalChainRewriter_Builder(this_, compiler)
	return this_
}

func InitOptionalChainRewriter_Builder(this_ *OptionalChainRewriter_Builder, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_OptionalChainRewriter_Builder()
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: Builder.setTmpVarNameCreator(com.google.javascript.jscomp.OptionalChainRewriter.TmpVarNameCreator tmpVarNameCreator)
func (this_ *OptionalChainRewriter_Builder) SetTmpVarNameCreator(tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator) *OptionalChainRewriter_Builder {
	_ = tmpVarNameCreator
	this_.tmpVarNameCreator = j2gort.CheckNotNullIface(tmpVarNameCreator)
	return this_
}

// Java method: Builder.setScope(com.google.javascript.jscomp.Scope scope)
func (this_ *OptionalChainRewriter_Builder) SetScope(scope *Scope) *OptionalChainRewriter_Builder {
	_ = scope
	this_.scope = j2gort.CheckNotNullPtr(scope)
	return this_
}

// Java method: Builder.build(com.google.javascript.rhino.Node wholeChain)
func (this_ *OptionalChainRewriter_Builder) Build(wholeChain pkg_rhino.Node) *OptionalChainRewriter {
	_ = wholeChain
	return NewOptionalChainRewriter(this_, wholeChain)
}
