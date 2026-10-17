// Port of com/google/javascript/jscomp/ijs/ProcessConstJsdocCallback.java @ closure-compiler v20261015

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ProcessConstJsdocCallback is a callback that calls the abstract method on
// every "inferrable const". This is a constant declaration for which there is
// no declared type and an RHS is present. This is useful for giving warnings
// like the CONSTANT_WITHOUT_EXPLICIT_TYPE diagnostic.
//
// As a side effect, this callback also populates the given FileInfo (assumed
// empty) with all of the declarations found throughout the compilation.
type ProcessConstJsdocCallback interface {
	NodeTraversal_AbstractPostOrderCallback

	GetProcessConstJsdocCallbackBase() *ProcessConstJsdocCallbackBase

	// Java: ProcessConstJsdocCallback.processConstWithRhs(NodeTraversal, Node)
	ProcessConstWithRhs(t *NodeTraversal, lhs *rhino.Node)
}

// ProcessConstJsdocCallbackBase holds the fields and base methods of the
// abstract class ProcessConstJsdocCallback; subclasses embed it, define
// ProcessConstWithRhs and hand themselves to InitProcessConstJsdocCallbackBase.
type ProcessConstJsdocCallbackBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	self ProcessConstJsdocCallback

	currentFile *FileInfo
}

var processConstJsdocCallback_GOOG_DEFINE = rhino.QualifiedName_Of("goog.define")
var processConstJsdocCallback_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")
var processConstJsdocCallback_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")
var processConstJsdocCallback_CJS_REQUIRE = rhino.QualifiedName_Of("require")

// Java: ProcessConstJsdocCallback(FileInfo)
func (b *ProcessConstJsdocCallbackBase) InitProcessConstJsdocCallbackBase(self ProcessConstJsdocCallback, currentFile *FileInfo) {
	b.self = self
	b.currentFile = currentFile
}

// GetProcessConstJsdocCallbackBase gives access to the base through the
// interface.
func (b *ProcessConstJsdocCallbackBase) GetProcessConstJsdocCallbackBase() *ProcessConstJsdocCallbackBase {
	return b
}

// Java: ProcessConstJsdocCallback.visit(NodeTraversal, Node, Node)
func (b *ProcessConstJsdocCallbackBase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		if NodeUtil_IsStatementParent(parent) {
			b.currentFile.RecordNameDeclaration(n.GetFirstChild())
		}
	case rhino.Token_MEMBER_FIELD_DEF:
		if NodeUtil_GetRValueOfLValue(n) != nil {
			b.ProcessDeclarationWithRhs(t, n)
		}
		b.currentFile.RecordMemberFieldDef(n)
	case rhino.Token_FUNCTION:
		if NodeUtil_IsStatementParent(parent) {
			b.currentFile.RecordNameDeclaration(n.GetFirstChild())
		} else if ClassUtil_IsClassMethod(n) && ClassUtil_HasNamedClass(n) {
			b.currentFile.RecordMethod(n)
		}
	case rhino.Token_EXPR_RESULT:
		expr := n.GetFirstChild()
		switch expr.GetToken() {
		case rhino.Token_CALL:
			callee := expr.GetFirstChild()
			if processConstJsdocCallback_GOOG_PROVIDE.Matches(callee) {
				b.currentFile.MarkProvided(expr.GetLastChild().GetString())
			} else if processConstJsdocCallback_GOOG_REQUIRE.Matches(callee) || processConstJsdocCallback_CJS_REQUIRE.Matches(callee) {
				b.currentFile.RecordImport(expr.GetLastChild().GetString())
			} else if processConstJsdocCallback_GOOG_DEFINE.Matches(callee) {
				b.currentFile.RecordDefine(expr)
			}
		case rhino.Token_ASSIGN:
			b.RecordDeclaration(t, expr.GetFirstChild(), expr.GetLastChild())
		case rhino.Token_GETPROP:
			b.currentFile.RecordNameDeclaration(expr)
		case rhino.Token_GETELEM:
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg("Unexpected declaration: " + expr.String()))
		}
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		jutil.CheckStateMsg1(n.HasOneChild(), "%s", n)
		b.RecordDeclaration(t, n.GetFirstChild(), n.GetFirstChild().GetLastChild())
	case rhino.Token_STRING_KEY:
		if parent.IsObjectLit() && n.HasOneChild() {
			b.ProcessDeclarationWithRhs(t, n)
			b.currentFile.RecordStringKeyDeclaration(n)
		}
	default:
	}
}

// Java: ProcessConstJsdocCallback.recordDeclaration(NodeTraversal, Node, Node)
func (b *ProcessConstJsdocCallbackBase) RecordDeclaration(t *NodeTraversal, lhs *rhino.Node, rhs *rhino.Node) {
	if rhs != nil &&
		rhs.IsCall() &&
		processConstJsdocCallback_GOOG_DEFINE.Matches(rhs.GetFirstChild()) &&
		lhs.IsQualifiedName() {
		b.currentFile.RecordDefine(rhs)
	} else if lhs.IsGetElem() {
		return
	} else {
		b.RecordNameDeclaration(lhs, rhs)
		if !lhs.IsDestructuringLhs() && rhs != nil {
			b.ProcessDeclarationWithRhs(t, lhs)
		}
	}
}

// Java: ProcessConstJsdocCallback.recordNameDeclaration(Node, Node)
func (b *ProcessConstJsdocCallbackBase) RecordNameDeclaration(lhs *rhino.Node, rhs *rhino.Node) {
	jutil.CheckArgument(NodeUtil_IsNameDeclaration(lhs.GetParent()) || lhs.GetParent().IsAssign())
	isImport := PotentialDeclaration_IsImportRhs(rhs)
	isAlias := PotentialDeclaration_IsAliasDeclaration(lhs, rhs)
	NodeUtil_VisitLhsNodesInNode(
		lhs.GetParent(),
		func(name *rhino.Node) {
			if isAlias || isImport {
				b.currentFile.RecordAliasDeclaration(name)
			} else {
				b.currentFile.RecordNameDeclaration(name)
			}
		})
}

// Java: ProcessConstJsdocCallback.processDeclarationWithRhs(NodeTraversal, Node)
func (b *ProcessConstJsdocCallbackBase) ProcessDeclarationWithRhs(t *NodeTraversal, lhs *rhino.Node) {
	jutil.CheckArgumentMsg1(lhs.IsQualifiedName() || lhs.IsStringKey() || lhs.IsMemberFieldDef(), "%s", lhs)
	jutil.CheckStateMsg1(NodeUtil_GetRValueOfLValue(lhs) != nil, "%s", lhs)
	if PotentialDeclaration_IsConstToBeInferred(lhs) {
		b.self.ProcessConstWithRhs(t, lhs)
	}
}
