// Port of com/google/javascript/jscomp/Es6ConvertSuper.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// Es6ConvertSuper converts `super.method()` calls and adds constructors to
// any classes that lack them.
//
// This has to run before the main Es6RewriteClass pass. The super()
// constructor calls are not converted here, but rather in
// Es6ConvertSuperConstructorCalls, which runs later.
type Es6ConvertSuper struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler               AbstractCompiler
	astFactory             *AstFactory
	constructorSynthesizer *SynthesizeExplicitConstructors
}

var es6ConvertSuper_FEATURES_TO_RUN_FOR = parser.FeatureSet_BARE_MINIMUM.WithFeatures(
	parser.FeatureSet_Feature_SUPER,
	parser.FeatureSet_Feature_CLASSES,
)

// Java: Es6ConvertSuper(AbstractCompiler)
func NewEs6ConvertSuper(compiler AbstractCompiler) *Es6ConvertSuper {
	return &Es6ConvertSuper{
		compiler:               compiler,
		astFactory:             compiler.CreateAstFactory(),
		constructorSynthesizer: NewSynthesizeExplicitConstructors(compiler),
	}
}

// Java: Es6ConvertSuper.visit(NodeTraversal, Node, Node)
func (e *Es6ConvertSuper) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsClass(n); ok {
		e.constructorSynthesizer.SynthesizeClassConstructorIfMissing(t, n)
	} else if n.IsSuper() {
		e.visitSuper(n, parent)
	}
}

func (e *Es6ConvertSuper) isInterface(classNode rhino.Node) bool {
	classJsDocInfo := NodeUtil_GetBestJSDocInfo(classNode)
	return classJsDocInfo != nil && classJsDocInfo.IsInterface()
}

func (e *Es6ConvertSuper) visitSuper(node rhino.Node, parent rhino.Node) {
	exprRoot := node
	if _, ok := ast.AsNormalGet(exprRoot.GetParent()); ok {
		exprRoot = exprRoot.GetParent()
	}

	enclosingMemberDef := NodeUtil_GetEnclosingNode(
		exprRoot,
		func(n rhino.Node) bool {
			_, isMethod := ast.AsMethodDef(n)
			_, isCompProp := ast.AsComputedProp(n)
			return isMethod || isCompProp
		})

	if enclosingMemberDef.IsNil() || enclosingMemberDef.GetParent().IsObjectLit() {
		e.compiler.Report(
			JSError_MakeWithNode(node, TranspilationUtil_CANNOT_CONVERT, "super is not supported in object literal methods."))
		return
	}

	if _, isCall := ast.AsCall(parent); isCall {
		e.visitSuperCall(node, enclosingMemberDef)
	} else if propGet, isGet := ast.AsNormalGet(parent); isGet {
		if propGet.Object().Node() == node {
			if _, isParentCall := ast.AsCall(parent.GetParent()); isParentCall && NodeUtil_IsInvocationTarget(parent) {
				e.visitSuperPropertyCall(node, parent, enclosingMemberDef)
			} else {
				e.visitSuperPropertyAccess(node, parent, enclosingMemberDef)
			}
		} else {
			e.compiler.Report(
				JSError_MakeWithNode(
					node,
					TranspilationUtil_CANNOT_CONVERT_YET,
					"Only calls to super or to a method of super are supported."))
		}
	} else if _, isNew := ast.AsNew(parent); isNew {
		panic(jutil.NewIllegalStateExceptionMsg("This should never happen. Did Es6SuperCheck fail to run?"))
	} else {
		e.compiler.Report(
			JSError_MakeWithNode(
				node,
				TranspilationUtil_CANNOT_CONVERT_YET,
				"Only calls to super or to a method of super are supported."))
	}
}

func (e *Es6ConvertSuper) visitSuperCall(node rhino.Node, enclosingMemberDef rhino.Node) {
	clazz := NodeUtil_GetEnclosingClass(node)
	superName := es6ConvertSuper_secondChildOfClazz(clazz)
	if !superName.IsQualifiedName() {
		// This will be reported as an error in Es6ToEs3Converter.
		return
	}

	if mfd, ok := ast.AsMemberFunctionDef(enclosingMemberDef); ok && mfd.IsEs6Constructor() {
		// Calls to super() constructors will be transpiled by Es6ConvertSuperConstructorCalls later.
		if node.IsFromExterns() || e.isInterface(clazz) {
			enclosingStatement := NodeUtil_GetEnclosingStatement(node)
			enclosingStatementParent := enclosingStatement.GetParent()
			enclosingStatement.Detach()
			e.compiler.ReportChangeToEnclosingScope(enclosingStatementParent)
		}
		return
	}
	panic(jutil.NewIllegalStateExceptionMsg("This should never happen. Did Es6SuperCheck fail to run?"))
}

func (e *Es6ConvertSuper) visitSuperPropertyCall(node rhino.Node, parent rhino.Node, enclosingMemberDef rhino.Node) {
	call := ast.MustCall(parent.GetParent())
	grandparent := call.Node()

	clazz := NodeUtil_GetEnclosingClass(node)
	superName := es6ConvertSuper_secondChildOfClazz(clazz)
	if !superName.IsQualifiedName() {
		return
	}

	callTarget := parent
	if enclosingMemberDef.IsStaticMember() {
		expandedSuper := superName.CloneTree().SrcrefTree(node)
		expandedSuper.SetOriginalName(jutil.StringPtr("super"))
		node.ReplaceWith(expandedSuper)
		callTarget = e.astFactory.CreateGetPropWithUnknownType(callTarget.Detach(), "call")
		grandparent.AddChildToFront(callTarget)
		thisNode := e.astFactory.CreateThis(AstFactory_Type_(clazz))
		thisNode.MakeNonIndexable()
		thisNode.InsertAfter(callTarget)
		grandparent.SrcrefTreeIfMissing(parent)

		e.ensureNoCollapseOnStaticMemberDef(enclosingMemberDef)
	} else {
		expandedSuper := e.astFactory.CreatePrototypeAccess(superName.CloneTree()).SrcrefTree(node)
		expandedSuper.SetOriginalName(jutil.StringPtr("super"))
		node.ReplaceWith(expandedSuper)
		callTarget = e.astFactory.CreateGetPropWithUnknownType(callTarget.Detach(), "call")
		grandparent.AddChildToFront(callTarget)
		thisNode := e.astFactory.CreateThisForEs6Class(clazz)
		thisNode.MakeNonIndexable()
		thisNode.InsertAfter(callTarget)
		call.SetFreeCall(false)
		grandparent.SrcrefTreeIfMissing(parent)
	}
	e.compiler.ReportChangeToEnclosingScope(grandparent)
}

func (e *Es6ConvertSuper) visitSuperPropertyAccess(node rhino.Node, parent rhino.Node, enclosingMemberDef rhino.Node) {
	grandparent := parent.GetParent()

	if NodeUtil_IsLValue(parent) {
		e.compiler.Report(JSError_MakeWithNode(parent, TranspilationUtil_CANNOT_CONVERT_YET, "assigning to a super property"))
		return
	}

	clazz := NodeUtil_GetEnclosingClass(node)
	superName := es6ConvertSuper_secondChildOfClazz(clazz)
	if !superName.IsQualifiedName() {
		return
	}

	var newProp rhino.Node
	if enclosingMemberDef.IsStaticMember() {
		newProp = superName.CloneTree().SrcrefTree(node)
	} else {
		newProp = e.astFactory.CreatePrototypeAccess(superName.CloneTree()).SrcrefTree(node)
	}

	if gp, ok := ast.AsGetProp(parent); ok &&
		e.compiler.GetAccessorSummary().GetKind(gp.Property()).HasGetter() {
		superPlaceholder := rhino.IR_Name("sp")

		reflectGet := e.astFactory.CreateGetPropWithUnknownType(
			e.astFactory.CreateConstantName("Reflect", AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
			"get")
		propertyArg := e.astFactory.CreateCall(
			e.astFactory.CreateName(
				NodeUtil_JSC_PROPERTY_NAME_FN, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
			AstFactory_Type_OfColor(colors.StandardColors_STRING),
			e.astFactory.CreateString(gp.Property()),
			superName.CloneTree().SrcrefTree(node))
		reflectGetCall := e.astFactory.
			CreateCall(
				reflectGet,
				AstFactory_Type_(parent),
				superPlaceholder,
				propertyArg,
				e.astFactory.CreateThis(AstFactory_Type_(clazz))).
			SrcrefTreeIfMissing(parent)

		parent.ReplaceWith(reflectGetCall)
		superPlaceholder.ReplaceWith(node.Detach())

		e.ensureNoCollapseOnStaticMemberDef(enclosingMemberDef)
	}

	node.ReplaceWith(newProp)
	newProp.SetOriginalName(jutil.StringPtr("super"))

	e.compiler.ReportChangeToEnclosingScope(grandparent)
}

// es6ConvertSuper_secondChildOfClazz preserves the JVM's helpful NPE message
// when `clazz` is null (e.g. a `super` expression with no enclosing class).
func es6ConvertSuper_secondChildOfClazz(clazz rhino.Node) rhino.Node {
	if clazz.IsNil() {
		panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "com.google.javascript.rhino.Node.getSecondChild()" because "clazz" is null`))
	}
	return ast.MustClass(clazz).ExtendsNode()
}

func (e *Es6ConvertSuper) ensureNoCollapseOnStaticMemberDef(memberDef rhino.Node) {
	if !memberDef.IsStaticMember() {
		return
	}
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(memberDef.GetJSDocInfo())
	if builder.RecordNoCollapse() {
		memberDef.SetJSDocInfo(builder.Build())
	}
}

// Java: Es6ConvertSuper.process(Node, Node)
func (e *Es6ConvertSuper) Process(externs rhino.Node, root rhino.Node) {
	// Might need to synthesize constructors for ambient classes in .d.ts externs
	TranspilationPasses_ProcessTranspile(e.compiler, externs, es6ConvertSuper_FEATURES_TO_RUN_FOR, e)
	TranspilationPasses_ProcessTranspile(e.compiler, root, es6ConvertSuper_FEATURES_TO_RUN_FOR, e)
}
