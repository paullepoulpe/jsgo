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

// TODO(b/329447979): Run this pass only for scripts containing Feature.SUPER and not for scripts
// containing just Feature.CLASSES?
//
// Java: Es6ConvertSuper.FEATURES_TO_RUN_FOR
var es6ConvertSuper_FEATURES_TO_RUN_FOR = parser.FeatureSet_BARE_MINIMUM.WithFeatures(parser.FeatureSet_Feature_SUPER, parser.FeatureSet_Feature_CLASSES)

// Java: Es6ConvertSuper(AbstractCompiler)
func NewEs6ConvertSuper(compiler AbstractCompiler) *Es6ConvertSuper {
	e := &Es6ConvertSuper{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.constructorSynthesizer = NewSynthesizeExplicitConstructors(compiler)
	return e
}

// Java: Es6ConvertSuper.visit(NodeTraversal, Node, Node)
func (e *Es6ConvertSuper) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsClass() {
		e.constructorSynthesizer.SynthesizeClassConstructorIfMissing(t, n)
	} else if n.IsSuper() {
		e.VisitSuper(n, parent)
	}
}

// Java: Es6ConvertSuper.isInterface(Node)
func (e *Es6ConvertSuper) IsInterface(classNode rhino.Node) bool {
	classJsDocInfo := NodeUtil_GetBestJSDocInfo(classNode)
	return classJsDocInfo != nil && classJsDocInfo.IsInterface()
}

// Java: Es6ConvertSuper.visitSuper(Node, Node)
func (e *Es6ConvertSuper) VisitSuper(node rhino.Node, parent rhino.Node) {
	jutil.CheckState(node.IsSuper())
	exprRoot := node

	if exprRoot.GetParent().IsGetElem() || exprRoot.GetParent().IsGetProp() {
		exprRoot = exprRoot.GetParent()
	}

	enclosingMemberDef := NodeUtil_GetEnclosingNode(
		exprRoot,
		func(n rhino.Node) bool {
			switch n.GetToken() {
			case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP:
				return true
			default:
				return false
			}
		})

	if enclosingMemberDef.IsNil() || enclosingMemberDef.GetParent().IsObjectLit() {
		e.compiler.Report(
			JSError_MakeWithNode(node, TranspilationUtil_CANNOT_CONVERT, "super is not supported in object literal methods."))
		return
	}

	if parent.IsCall() {
		// super(...)
		e.VisitSuperCall(node, parent, enclosingMemberDef)
	} else if parent.IsGetProp() || parent.IsGetElem() {
		if parent.GetFirstChild() == node {
			if parent.GetParent().IsCall() && NodeUtil_IsInvocationTarget(parent) {
				// super.something(...) or super['something'](..)
				e.VisitSuperPropertyCall(node, parent, enclosingMemberDef)
			} else {
				// super.something or super['something']
				e.VisitSuperPropertyAccess(node, parent, enclosingMemberDef)
			}
		} else {
			// super.something used in some other way
			e.compiler.Report(
				JSError_MakeWithNode(
					node,
					TranspilationUtil_CANNOT_CONVERT_YET,
					"Only calls to super or to a method of super are supported."))
		}
	} else if parent.IsNew() {
		panic(jutil.NewIllegalStateExceptionMsg("This should never happen. Did Es6SuperCheck fail to run?"))
	} else {
		// some other use of super we don't support yet
		e.compiler.Report(
			JSError_MakeWithNode(
				node,
				TranspilationUtil_CANNOT_CONVERT_YET,
				"Only calls to super or to a method of super are supported."))
	}
}

// Java: Es6ConvertSuper.visitSuperCall(Node, Node, Node)
func (e *Es6ConvertSuper) VisitSuperCall(node rhino.Node, parent rhino.Node, enclosingMemberDef rhino.Node) {
	jutil.CheckStateMsg1(parent.IsCall(), "%s", parent)
	jutil.CheckStateMsg1(node.IsSuper(), "%s", node)

	clazz := NodeUtil_GetEnclosingClass(node)
	superName := es6ConvertSuper_secondChildOfClazz(clazz)
	if !superName.IsQualifiedName() {
		// This will be reported as an error in Es6ToEs3Converter.
		return
	}

	if NodeUtil_IsEs6ConstructorMemberFunctionDef(enclosingMemberDef) {
		// Calls to super() constructors will be transpiled by Es6ConvertSuperConstructorCalls later.
		if node.IsFromExterns() || e.IsInterface(clazz) {
			// If a class is defined in an externs file or as an interface, it's only a stub, not an
			// implementation that should be instantiated.
			// A call to super() shouldn't actually exist for these cases and is problematic to
			// transpile, so just drop it.
			enclosingStatement := NodeUtil_GetEnclosingStatement(node)
			enclosingStatementParent := enclosingStatement.GetParent()
			enclosingStatement.Detach()
			e.compiler.ReportChangeToEnclosingScope(enclosingStatementParent)
		}
		// Calls to super() constructors will be transpiled by Es6ConvertSuperConstructorCalls
		// later.
		return
	} else {
		// super can only be directly called in a constructor
		panic(jutil.NewIllegalStateExceptionMsg("This should never happen. Did Es6SuperCheck fail to run?"))
	}
}

// Java: Es6ConvertSuper.visitSuperPropertyCall(Node, Node, Node)
func (e *Es6ConvertSuper) VisitSuperPropertyCall(node rhino.Node, parent rhino.Node, enclosingMemberDef rhino.Node) {
	jutil.CheckStateMsg1(parent.IsGetProp() || parent.IsGetElem(), "%s", parent)
	jutil.CheckStateMsg1(node.IsSuper(), "%s", node)
	grandparent := parent.GetParent()
	jutil.CheckState(grandparent.IsCall())

	clazz := NodeUtil_GetEnclosingClass(node)
	superName := es6ConvertSuper_secondChildOfClazz(clazz)
	if !superName.IsQualifiedName() {
		// This will be reported as an error in Es6ToEs3Converter.
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
		thisNode.MakeNonIndexable() // no direct correlation with original source
		thisNode.InsertAfter(callTarget)
		grandparent.SrcrefTreeIfMissing(parent)

		// We added a `this` reference to the call. We need to make sure this member is not collapsed.
		e.EnsureNoCollapseOnStaticMemberDef(enclosingMemberDef)
	} else {
		// Replace super node to give
		// super.method(...) -> SuperClass.prototype.method(...)
		expandedSuper := e.astFactory.CreatePrototypeAccess(superName.CloneTree()).SrcrefTree(node)
		expandedSuper.SetOriginalName(jutil.StringPtr("super"))
		node.ReplaceWith(expandedSuper)
		// Set the 'this' object correctly for the call
		// SuperClass.prototype.method(...) -> SuperClass.prototype.method.call(this, ...)
		callTarget = e.astFactory.CreateGetPropWithUnknownType(callTarget.Detach(), "call")
		grandparent.AddChildToFront(callTarget)
		thisNode := e.astFactory.CreateThisForEs6Class(clazz)
		thisNode.MakeNonIndexable() // no direct correlation with original source
		thisNode.InsertAfter(callTarget)
		grandparent.PutBooleanProp(rhino.Node_FREE_CALL, false)
		grandparent.SrcrefTreeIfMissing(parent)
	}
	e.compiler.ReportChangeToEnclosingScope(grandparent)
}

// Java: Es6ConvertSuper.visitSuperPropertyAccess(Node, Node, Node)
func (e *Es6ConvertSuper) VisitSuperPropertyAccess(node rhino.Node, parent rhino.Node, enclosingMemberDef rhino.Node) {
	jutil.CheckStateMsg1(parent.IsGetProp() || parent.IsGetElem(), "%s", parent)
	jutil.CheckStateMsg1(node.IsSuper(), "%s", node)
	grandparent := parent.GetParent()

	if NodeUtil_IsLValue(parent) {
		// We don't support assigning to a super property
		e.compiler.Report(JSError_MakeWithNode(parent, TranspilationUtil_CANNOT_CONVERT_YET, "assigning to a super property"))
		return
	}

	clazz := NodeUtil_GetEnclosingClass(node)
	superName := es6ConvertSuper_secondChildOfClazz(clazz)
	if !superName.IsQualifiedName() {
		// This will be reported as an error in Es6ToEs3Converter.
		return
	}

	var newProp rhino.Node
	if enclosingMemberDef.IsStaticMember() {
		// super.prop -> SuperClass.prop
		newProp = superName.CloneTree().SrcrefTree(node)
	} else {
		// super.prop -> SuperClass.prototype.prop
		newProp = e.astFactory.CreatePrototypeAccess(superName.CloneTree()).SrcrefTree(node)
	}

	if parent.IsGetProp() &&
		e.compiler.GetAccessorSummary().GetKind(parent.GetString()).HasGetter() {
		// Rewrites
		//   super.x
		// as
		//   Reflect.get(super, JSCompiler_renameProperty('x', ParentClass), this)
		superPlaceholder := rhino.IR_Name("sp") // No need for type info as this will be replaced.

		reflectGet := e.astFactory.CreateGetPropWithUnknownType(
			e.astFactory.CreateConstantName("Reflect", AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
			"get")
		var propertyArg rhino.Node
		if parent.IsGetProp() {
			propertyArg = e.astFactory.CreateCall(
				e.astFactory.CreateName(
					NodeUtil_JSC_PROPERTY_NAME_FN, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
				AstFactory_Type_OfColor(colors.StandardColors_STRING),
				e.astFactory.CreateString(parent.GetString()),
				superName.CloneTree().SrcrefTree(node))
		} else {
			// For getElem, we just use the accessor value in the Reflect.get call.
			propertyArg = node.GetNext().Detach()
		}
		reflectGetCall := e.astFactory.
			CreateCall(
				reflectGet,
				AstFactory_Type_(parent),
				superPlaceholder,
				propertyArg,
				e.astFactory.CreateThis(AstFactory_Type_(clazz))).
			SrcrefTreeIfMissing(parent)

		// Replace the getProp/getElem with the Reflect.get call.
		parent.ReplaceWith(reflectGetCall)
		// Swap in the `n` (super) where it goes in the Reflect.get call.
		superPlaceholder.ReplaceWith(node.Detach())

		// We added a `this` reference to the call. We need to make sure this member is not collapsed.
		e.EnsureNoCollapseOnStaticMemberDef(enclosingMemberDef)
	}

	node.ReplaceWith(newProp)
	newProp.SetOriginalName(jutil.StringPtr("super"))

	e.compiler.ReportChangeToEnclosingScope(grandparent)
}

// es6ConvertSuper_secondChildOfClazz is clazz.getSecondChild(), the
// superclass of the class that encloses a super.
//
// Go-only: language-adapter — the JVM's "helpful" NullPointerException
// message (jutil.NewHelpfulNullPointerException). A super in an object
// literal's method has no enclosing class, so clazz is null there and Java
// dereferences it; the INTERNAL COMPILER ERROR report that follows prints
// the exception's message, which the JVM derives from the bytecode of this
// call. A Go nil dereference has none.
func es6ConvertSuper_secondChildOfClazz(clazz rhino.Node) rhino.Node {
	if clazz.IsNil() {
		panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "com.google.javascript.rhino.Node.getSecondChild()" because "clazz" is null`))
	}
	return clazz.GetSecondChild()
}

// Java: Es6ConvertSuper.ensureNoCollapseOnStaticMemberDef(Node)
func (e *Es6ConvertSuper) EnsureNoCollapseOnStaticMemberDef(memberDef rhino.Node) {
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
