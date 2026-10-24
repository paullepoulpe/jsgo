// Port of com/google/javascript/jscomp/SourceInfoCheck.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SourceInfoCheck is a simple pass to enforce some invariants about source
// information.
//
//   - All nodes in externs and code have a line number and non-zero length
//   - All JSTypeExpressions have a source name
//   - All nodes in a JSTYpeExpression have a source file attached (does not
//     verify lineno or charno)
//
// This pass does not check the content of the line numbers for things like
// "being monotonically increasing" because we often insert generated code
// with out-of-order line numbers. (e.g. imagine inlining a function from a.js
// into b.js, but having the source information point to a.js)
type SourceInfoCheck struct {
	compiler            AbstractCompiler
	requiresLineNumbers bool
}

// Java: SourceInfoCheck.MISSING_LINE_INFO
var sourceInfoCheck_MISSING_LINE_INFO = DiagnosticType_Error(
	"JSC_MISSING_LINE_INFO",
	"No source location information associated with {0}."+
		"\nMost likely a Node has been created without setting the source file"+
		" and line/column location.  Usually this is done using"+
		" Node.srcrefIfMissing and supplying a Node from the source AST.")

// Java: SourceInfoCheck.MISSING_LENGTH
var sourceInfoCheck_MISSING_LENGTH = DiagnosticType_Error(
	"JSC_MISSING_LENGTH",
	"Negative length associated with {0}.\n"+
		"Most likely a Node's source information was set incorrectly at parse time.")

// Java: SourceInfoCheck.MISSING_SOURCE_NAME
var sourceInfoCheck_MISSING_SOURCE_NAME = DiagnosticType_Error(
	"JSC_MISSING_SOURCE_NAME",
	"No source name associated with {0}.\n"+
		"Most likely a new type was created without setting the source name.")

// Java: SourceInfoCheck(AbstractCompiler)
func NewSourceInfoCheck(compiler AbstractCompiler) *SourceInfoCheck {
	s := &SourceInfoCheck{}
	s.requiresLineNumbers = false
	s.compiler = compiler
	return s
}

// Java: SourceInfoCheck.setCheckSubTree(Node)
func (s *SourceInfoCheck) SetCheckSubTree(root *rhino.Node) {
	s.requiresLineNumbers = true

	NodeTraversal_Traverse(s.compiler, root, s)
}

// Java: SourceInfoCheck.process(Node, Node)
func (s *SourceInfoCheck) Process(externs *rhino.Node, root *rhino.Node) {
	s.requiresLineNumbers = false

	NodeTraversal_Traverse(s.compiler, root, s)
}

// Java: SourceInfoCheck.shouldTraverse(NodeTraversal, Node, Node)
func (s *SourceInfoCheck) ShouldTraverse(unused *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	// Each JavaScript file is rooted in a script node, so we'll only
	// have line number information inside the script node.
	if n.IsScript() {
		s.requiresLineNumbers = true
	}
	return true
}

// Java: SourceInfoCheck.visit(NodeTraversal, Node, Node)
func (s *SourceInfoCheck) Visit(unused *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsScript() {
		s.requiresLineNumbers = false
	} else if s.requiresLineNumbers {
		if n.GetLineno() == -1 {
			// The tree version of the node is really the best diagnostic
			// info we have to offer here.
			s.compiler.Report(JSError_MakeWithNode(n, sourceInfoCheck_MISSING_LINE_INFO, n.ToStringTree()))
		} else if n.GetLength() < 0 {
			s.compiler.Report(JSError_MakeWithNode(n, sourceInfoCheck_MISSING_LENGTH, n.ToStringTree()))
		}
		info := n.GetJSDocInfo()
		if info != nil {
			s.CheckJSDoc(info)
		}
	}
}

// CheckJSDoc verifies that all type nodes have source files (if not an
// actual length).
//
// Java: SourceInfoCheck.checkJSDoc(JSDocInfo)
func (s *SourceInfoCheck) CheckJSDoc(info *rhino.JSDocInfo) {
	for it := info.TypeExprIter(); it.Next(); {
		expression := it.Expr()
		// Go-only: language-adapter — rhino.JSTypeExpression stores a Java
		// null source name as "".
		if expression.GetSourceName() == "" {
			s.compiler.Report(
				JSError_MakeWithNode(
					expression.GetRoot(), sourceInfoCheck_MISSING_SOURCE_NAME, "JSTypeExpression "+expression.String()))
		} else {
			NodeUtil_VisitPreOrder(expression.GetRoot(), NodeUtil_VisitorFunc(s.CheckSourceFile))
		}
	}
}

// Java: SourceInfoCheck.checkSourceFile(Node)
func (s *SourceInfoCheck) CheckSourceFile(n *rhino.Node) {
	sourceName := n.GetStaticSourceFile()
	if sourceName == nil {
		s.compiler.Report(JSError_MakeWithNode(n, sourceInfoCheck_MISSING_SOURCE_NAME, "type node "+n.ToStringTree()))
	}
}
