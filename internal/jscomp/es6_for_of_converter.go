// Port of com/google/javascript/jscomp/Es6ForOfConverter.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6ForOfConverter converts ES6 "for of" loops to ES5.
type Es6ForOfConverter struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler   AbstractCompiler
	namer      *DefaultNameGenerator
	astFactory *AstFactory
	namespace  rhino.StaticScope
}

// Java: Es6ForOfConverter.transpiledFeatures
var Es6ForOfConverter_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_FOR_OF)

const (
	// Java: Es6ForOfConverter.ITER_BASE
	Es6ForOfConverter_ITER_BASE = "$jscomp$iter$"

	// Java: Es6ForOfConverter.ITER_RESULT
	Es6ForOfConverter_ITER_RESULT = "$jscomp$key$"
)

// Java: Es6ForOfConverter(AbstractCompiler)
func NewEs6ForOfConverter(compiler AbstractCompiler) *Es6ForOfConverter {
	e := &Es6ForOfConverter{}
	e.compiler = compiler
	e.namer = NewDefaultNameGenerator()
	e.astFactory = compiler.CreateAstFactory()
	e.namespace = compiler.GetTranspilationNamespace()
	return e
}

// Java: Es6ForOfConverter.process(Node, Node)
func (e *Es6ForOfConverter) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, Es6ForOfConverter_transpiledFeatures, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6ForOfConverter_transpiledFeatures)
}

// Java: Es6ForOfConverter.visit(NodeTraversal, Node, Node)
func (e *Es6ForOfConverter) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsForOf() {
		e.VisitForOf(t, n)
	}
}

// TODO(lharker): break up this method
//
// Java: Es6ForOfConverter.visitForOf(NodeTraversal, Node)
func (e *Es6ForOfConverter) VisitForOf(t *NodeTraversal, node rhino.Node) {
	// `var v` or `let v` or v any valid lhs
	variable := node.RemoveFirstChild()
	iterable := node.RemoveFirstChild()
	body := node.RemoveFirstChild()
	varJSDocInfo := variable.GetJSDocInfo()
	// `$jscomp$iter$0`
	iterName :=
		e.astFactory.CreateName(
			Es6ForOfConverter_ITER_BASE+e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput()),
			AstFactory_Type_OfColorId(colors.StandardColors_ITERATOR_ID))
	iterName.MakeNonIndexable()
	// `$jscomp$iter$0.next()`
	getNext :=
		e.astFactory.CreateCallWithUnknownType(
			e.astFactory.CreateGetPropWithUnknownType(iterName.CloneTree(), "next"))
	// generate a unique iterator result name for every for-of loop getting rewritten to avoid
	// conflicts
	iteratorResultName :=
		Es6ForOfConverter_ITER_RESULT + e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput()) + "$"

	if NodeUtil_IsNameDeclaration(variable) {
		iteratorResultName += variable.GetFirstChild().GetString()
	} else if variable.IsName() {
		iteratorResultName += variable.GetString()
	} else {
		// give arbitrary lhs expressions an arbitrary name
		iteratorResultName += e.namer.GenerateNextName()
	}
	// `$jscomp$key$extraName`
	iterResult := e.astFactory.CreateNameWithUnknownType(iteratorResultName)
	iterResult.MakeNonIndexable()

	// `$jscomp.makeIterator(iterable)`
	callMakeIterator :=
		e.astFactory.
			CreateJSCompMakeIteratorCall(iterable, e.namespace).
			SrcrefTreeIfMissing(iterable)
	// `var $jscomp$iter$0 = $jscomp.makeIterator(iterable)`
	initIter := rhino.IR_VarWithValue(iterName.CloneTree(), callMakeIterator).SrcrefTreeIfMissing(iterable)
	// var $jscomp$key$extraName = $jscomp$iter$0.next();
	initIterResult :=
		rhino.IR_VarWithValue(iterResult.CloneTree(), getNext.CloneTree()).SrcrefTreeIfMissing(iterable)

	// !$jscomp$key$extraName.done
	cond :=
		e.astFactory.CreateNot(
			e.astFactory.CreateGetProp(iterResult.CloneTree(), "done", AstFactory_Type_OfColor(colors.StandardColors_BOOLEAN)))
	// $jscomp$key$extraName = $jscomp$iter$0.next()
	incr := e.astFactory.CreateAssign(iterResult.CloneTree(), getNext.CloneTree())

	var declarationOrAssign rhino.Node
	if !NodeUtil_IsNameDeclaration(variable) {
		// e.g. `for(a.b of []) {}`
		declarationOrAssign =
			e.astFactory.CreateAssign(
				variable.CloneTree().SetJSDocInfo(nil),
				e.astFactory.CreateGetProp(iterResult.CloneTree(), "value", AstFactory_Type_(variable)))
		declarationOrAssign.SetJSDocInfo(varJSDocInfo)
		declarationOrAssign = rhino.IR_ExprResult(declarationOrAssign)
	} else {
		// `for(let a of []) {}` or `for(const a of []) {}`
		type_ := AstFactory_Type_(variable.GetFirstChild())
		declarationType := variable.GetToken() // i.e. VAR, CONST, or LET.
		jutil.CheckStateMsg(
			declarationType != rhino.Token_VAR,
			"var initializers must've gotten moved out of the loop during normalize")
		declarationOrAssign =
			e.astFactory.CreateSingleNameDeclaration(
				declarationType,
				variable.GetFirstChild().GetString(),
				e.astFactory.CreateGetProp(iterResult.CloneTree(), "value", type_))
		if variable.GetFirstChild().GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
			// if the original name was const, then the new name should be too
			// e.g. `for(let CID of []) {}` where `CID` was originally marked constant by coding
			// convention
			declarationOrAssign.GetFirstChild().PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
		declarationOrAssign.SetJSDocInfo(varJSDocInfo)
	}
	newBody := rhino.IR_BlockOfNodes(declarationOrAssign, body).Srcref(body)
	empty := e.astFactory.CreateEmpty()
	newFor := rhino.IR_ForNode(empty, cond, incr, newBody).SrcrefTreeIfMissing(node)

	// Build finally block:
	// (0, $jscomp.iteratorClose)($jscomp$iter$0, $jscomp$key$extraName);
	callIteratorClose :=
		e.astFactory.
			CreateJscompIteratorCloseCall(
				iterName.CloneTree(), iterResult.CloneTree(), e.namespace).
			SrcrefTreeIfMissing(node)
	finallyBlock := e.astFactory.CreateBlock(e.astFactory.ExprResult(callIteratorClose))

	// Check if the for loop has a parent that is a label i.e. `loop1: for(...of ...)`
	labelNames := []rhino.Node{}
	insertionPoint := node
	for !insertionPoint.GetParent().IsNil() && insertionPoint.GetParent().IsLabel() {
		insertionPoint = insertionPoint.GetParent()
		labelNames = append(labelNames, insertionPoint.GetFirstChild().CloneNode())
	}

	innerLoop := newFor
	for _, labelName := range labelNames {
		innerLoop = e.astFactory.CreateLabel(labelName, innerLoop)
	}
	tryBlock := e.astFactory.CreateBlock(innerLoop)
	tryFinally := e.astFactory.CreateTryFinally(tryBlock, finallyBlock).SrcrefTreeIfMissing(node)

	insertionPoint.ReplaceWith(tryFinally)

	initIter.InsertBefore(tryFinally)
	initIterResult.InsertAfter(initIter)
	e.compiler.ReportChangeToEnclosingScope(tryFinally)
}
