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
	"github.com/paullepoulpe/jsgo/internal/ast"
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

var es6ForOfVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FOR_OF)
	return s
}()

func (*Es6ForOfConverter) HandledTokens() ast.TokenSet {
	return es6ForOfVisitTokens
}

func (*Es6ForOfConverter) visitTokens() *ast.TokenSet {
	return &es6ForOfVisitTokens
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
	return &Es6ForOfConverter{
		compiler:   compiler,
		namer:      NewDefaultNameGenerator(),
		astFactory: compiler.CreateAstFactory(),
		namespace:  compiler.GetTranspilationNamespace(),
	}
}

// Java: Es6ForOfConverter.process(Node, Node)
func (e *Es6ForOfConverter) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, Es6ForOfConverter_transpiledFeatures, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6ForOfConverter_transpiledFeatures)
}

// Java: Es6ForOfConverter.visit(NodeTraversal, Node, Node)
func (e *Es6ForOfConverter) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if forOf, ok := ast.AsForOf(n); ok {
		e.visitForOf(t, forOf)
	}
}

func (e *Es6ForOfConverter) visitForOf(t *NodeTraversal, forOf ast.ForOf) {
	node := forOf.Node()
	variable := node.RemoveFirstChild()
	iterable := node.RemoveFirstChild()
	body := node.RemoveFirstChild()
	varJSDocInfo := variable.GetJSDocInfo()

	iterName :=
		e.astFactory.CreateName(
			Es6ForOfConverter_ITER_BASE+e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput()),
			AstFactory_Type_OfColorId(colors.StandardColors_ITERATOR_ID))
	iterName.MakeNonIndexable()

	getNext :=
		e.astFactory.CreateCallWithUnknownType(
			e.astFactory.CreateGetPropWithUnknownType(iterName.CloneTree(), "next"))

	iteratorResultName :=
		Es6ForOfConverter_ITER_RESULT + e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput()) + "$"

	nameDecl, isNameDecl := ast.AsNameDeclaration(variable)
	if isNameDecl {
		iteratorResultName += nameDecl.MustFirstDeclarator().MustName().Ident()
	} else if nm, ok := ast.AsName(variable); ok {
		iteratorResultName += nm.Ident()
	} else {
		iteratorResultName += e.namer.GenerateNextName()
	}

	iterResult := e.astFactory.CreateNameWithUnknownType(iteratorResultName)
	iterResult.MakeNonIndexable()

	callMakeIterator :=
		e.astFactory.
			CreateJSCompMakeIteratorCall(iterable, e.namespace).
			SrcrefTreeIfMissing(iterable)
	initIter := rhino.IR_VarWithValue(iterName.CloneTree(), callMakeIterator).SrcrefTreeIfMissing(iterable)
	initIterResult :=
		rhino.IR_VarWithValue(iterResult.CloneTree(), getNext.CloneTree()).SrcrefTreeIfMissing(iterable)

	cond :=
		e.astFactory.CreateNot(
			e.astFactory.CreateGetProp(iterResult.CloneTree(), "done", AstFactory_Type_OfColor(colors.StandardColors_BOOLEAN)))
	incr := e.astFactory.CreateAssign(iterResult.CloneTree(), getNext.CloneTree())

	var declarationOrAssign rhino.Node
	if !isNameDecl {
		declarationOrAssign =
			e.astFactory.CreateAssign(
				variable.CloneTree().SetJSDocInfo(nil),
				e.astFactory.CreateGetProp(iterResult.CloneTree(), "value", AstFactory_Type_(variable)))
		declarationOrAssign.SetJSDocInfo(varJSDocInfo)
		declarationOrAssign = rhino.IR_ExprResult(declarationOrAssign)
	} else {
		declName := nameDecl.MustFirstDeclarator().MustName()
		type_ := AstFactory_Type_(declName.Node())
		declarationType := nameDecl.Token()
		jutil.CheckStateMsg(
			declarationType != rhino.Token_VAR,
			"var initializers must've gotten moved out of the loop during normalize")
		declarationOrAssign =
			e.astFactory.CreateSingleNameDeclaration(
				declarationType,
				declName.Ident(),
				e.astFactory.CreateGetProp(iterResult.CloneTree(), "value", type_))
		if declName.IsConstantName() {
			ast.MustNameDeclaration(declarationOrAssign).MustFirstDeclarator().MustName().SetConstantName(true)
		}
		declarationOrAssign.SetJSDocInfo(varJSDocInfo)
	}
	newBody := rhino.IR_BlockOfNodes(declarationOrAssign, body).Srcref(body)
	empty := e.astFactory.CreateEmpty()
	newFor := rhino.IR_ForNode(empty, cond, incr, newBody).SrcrefTreeIfMissing(node)

	callIteratorClose :=
		e.astFactory.
			CreateJscompIteratorCloseCall(
				iterName.CloneTree(), iterResult.CloneTree(), e.namespace).
			SrcrefTreeIfMissing(node)
	finallyBlock := e.astFactory.CreateBlock(e.astFactory.ExprResult(callIteratorClose))

	var labelNames []rhino.Node
	insertionPoint := node
	for {
		lbl, ok := ast.AsLabel(insertionPoint.GetParent())
		if !ok {
			break
		}
		insertionPoint = lbl.Node()
		labelNames = append(labelNames, lbl.Name().Node().CloneNode())
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
