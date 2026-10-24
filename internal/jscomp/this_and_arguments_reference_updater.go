// Port of com/google/javascript/jscomp/ThisAndArgumentsReferenceUpdater.java @ closure-compiler v20261021

/*
 * Copyright 2024 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ThisAndArgumentsReferenceUpdater rewrites references to `this` and
// `arguments` in a single function to instead refer to local variables.
//
// When systematically transpiling away all arrow functions, an instance is
// generated for each arrow function in order of *decreasing* depth. This
// isn't too inefficient, because instances don't traverse into non-arrow
// functions and all nested functions will already have been "de-arrowed".
// This class is also used for one-off cases (e.g. instrumenting generators
// for AsyncContext), in which case there may be nested arrow functions that
// need to be recursed into.
type ThisAndArgumentsReferenceUpdater struct {
	compiler   AbstractCompiler
	context    *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext
	astFactory *AstFactory
}

// The name of the vars that capture 'this' and 'arguments' for converting
// arrow functions. Note that these names can be reused (once per scope)
// because declarations in nested scopes will shadow one another, which
// results in the intended behaviour.
const (
	// Java: ThisAndArgumentsReferenceUpdater.ARGUMENTS_VAR
	ThisAndArgumentsReferenceUpdater_ARGUMENTS_VAR = "$jscomp$arguments"
	// Java: ThisAndArgumentsReferenceUpdater.THIS_VAR
	ThisAndArgumentsReferenceUpdater_THIS_VAR = "$jscomp$this"
)

// Java: ThisAndArgumentsReferenceUpdater(AbstractCompiler, ThisAndArgumentsContext, AstFactory)
func NewThisAndArgumentsReferenceUpdater(
	compiler AbstractCompiler, context *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext, astFactory *AstFactory) *ThisAndArgumentsReferenceUpdater {
	u := &ThisAndArgumentsReferenceUpdater{}
	u.compiler = compiler
	u.context = context
	u.astFactory = astFactory
	return u
}

// Java: ThisAndArgumentsReferenceUpdater.visit(NodeTraversal, Node, Node)
func (u *ThisAndArgumentsReferenceUpdater) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsThis() {
		u.context.SetNeedsThisVarWithType(AstFactory_Type_(n))

		name := u.astFactory.CreateName(
			ThisAndArgumentsReferenceUpdater_THIS_VAR+"$"+u.context.uniqueId, u.context.GetThisType()).Srcref(n)
		name.MakeNonIndexable()
		if u.compiler.GetOptions().PreservesDetailedSourceInfo() {
			originalName := "this"
			name.SetOriginalName(&originalName)
		}

		n.ReplaceWith(name)
	} else if n.IsName() && n.GetString() == "arguments" {
		u.context.SetNeedsArgumentsVar()

		name := u.astFactory.
			CreateName(ThisAndArgumentsReferenceUpdater_ARGUMENTS_VAR+"$"+u.context.uniqueId, AstFactory_Type_(n)).
			Srcref(n)
		if u.compiler.GetOptions().PreservesDetailedSourceInfo() {
			originalName := "arguments"
			name.SetOriginalName(&originalName)
		}

		n.ReplaceWith(name)
	}
}

// Java: ThisAndArgumentsReferenceUpdater.shouldTraverse(NodeTraversal, Node, Node)
func (u *ThisAndArgumentsReferenceUpdater) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsMemberFieldDef() {
		return false
	}
	if !parent.IsNil() && parent.IsComputedFieldDef() && n != parent.GetFirstChild() {
		return false
	}
	return !n.IsFunction() || n.IsArrowFunction()
}

// ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext accumulates
// information about a scope in which `this` and `arguments` are consistent.
//
// Instances are maintained in a DFS stack during traversal of
// Es6RewriteArrowFunction. They can't be immutable because a context isn't
// fully defined by a single node (`super()` makes this hard).
type ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext struct {
	scopeBody          rhino.Node
	isConstructor      bool
	lastSuperStatement rhino.Node // Last statement in the body that refers to super().
	// An object of ThisAndArgumentsContext exists for each Script and each Function. Store a
	// uniqueId string based on the script node's filename's hash. This gets used to generate a
	// unique name when rewriting `this` to `$jscomp$this`.
	uniqueId string

	needsThisVar bool
	thisType     AstFactory_Type

	needsArgumentsVar bool
}

// Java: ThisAndArgumentsContext(Node, boolean, String)
func NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
	scopeBody rhino.Node, isConstructor bool, uniqueId string) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	c := &ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext{
		lastSuperStatement: rhino.NilNode,
		uniqueId:           "",
		needsThisVar:       false,
		needsArgumentsVar:  false,
	}
	c.scopeBody = scopeBody
	c.isConstructor = isConstructor
	c.uniqueId = uniqueId
	return c
}

// Java: ThisAndArgumentsContext.getThisType()
func (c *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) GetThisType() AstFactory_Type {
	return c.thisType
}

// Java: ThisAndArgumentsContext.setNeedsThisVarWithType(AstFactory.Type)
func (c *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) SetNeedsThisVarWithType(type_ AstFactory_Type) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	c.thisType = type_
	c.needsThisVar = true
	return c
}

// Java: ThisAndArgumentsContext.setNeedsArgumentsVar()
func (c *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) SetNeedsArgumentsVar() *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	c.needsArgumentsVar = true
	return c
}

// Java: ThisAndArgumentsContext.addVarDeclarations(AbstractCompiler, AstFactory, Node)
func (c *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) AddVarDeclarations(
	compiler AbstractCompiler, astFactory *AstFactory, script rhino.Node) {
	if c.needsThisVar {
		name := astFactory.CreateName(ThisAndArgumentsReferenceUpdater_THIS_VAR+"$"+c.uniqueId, c.GetThisType())
		thisVar := rhino.IR_ConstNode(name, astFactory.CreateThis(c.GetThisType()))
		NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature_CONST_DECLARATIONS, compiler)
		thisVar.SrcrefTreeIfMissing(c.scopeBody)
		ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_MakeTreeNonIndexable(thisVar)

		if c.lastSuperStatement.IsNil() {
			ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_InsertVarDeclaration(thisVar, c.scopeBody)
		} else {
			// Not safe to reference `this` until after super() has been called.
			// TODO(bradfordcsmith): Some complex cases still aren't covered, like
			//     if (...) { super(); arrow function } else { super(); }
			thisVar.InsertAfter(c.lastSuperStatement)
		}
		compiler.ReportChangeToEnclosingScope(thisVar)
	}

	if c.needsArgumentsVar {
		argumentsVar :=
			astFactory.CreateArgumentsAliasDeclaration(ThisAndArgumentsReferenceUpdater_ARGUMENTS_VAR + "$" + c.uniqueId)
		NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature_CONST_DECLARATIONS, compiler)

		ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_InsertVarDeclaration(argumentsVar, c.scopeBody)

		argumentsVar.SrcrefTreeIfMissing(c.scopeBody)
		compiler.ReportChangeToEnclosingScope(argumentsVar)
	}
}

// Java: ThisAndArgumentsContext.insertVarDeclaration(Node, Node)
func ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_InsertVarDeclaration(varDeclaration rhino.Node, scopeBody rhino.Node) {
	if scopeBody.GetParent().IsFunction() {
		// for functions, we must find the correct insertion point to preserve normalization
		insertBeforePoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(scopeBody)
		if !insertBeforePoint.IsNil() {
			varDeclaration.InsertBefore(insertBeforePoint)
		} else {
			// functionBody only contains hoisted function declarations
			scopeBody.AddChildToBack(varDeclaration)
		}
	} else {
		scopeBody.AddChildToFront(varDeclaration)
	}
}

// Java: ThisAndArgumentsContext.makeTreeNonIndexable(Node)
func ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_MakeTreeNonIndexable(n rhino.Node) {
	n.MakeNonIndexable()
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext_MakeTreeNonIndexable(child)
	}
}
