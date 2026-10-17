// Port of com/google/javascript/jscomp/SubstituteEs6Syntax.java @ closure-compiler v20261015

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SubstituteEs6Syntax is an optimization that does peephole optimizations
// of ES6 code.
type SubstituteEs6Syntax struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler                       AbstractCompiler
	objectLiteralShorthandWasAdded bool
}

// Java: SubstituteEs6Syntax(AbstractCompiler)
func NewSubstituteEs6Syntax(compiler AbstractCompiler) *SubstituteEs6Syntax {
	s := &SubstituteEs6Syntax{}
	s.compiler = compiler
	return s
}

// Java: SubstituteEs6Syntax.process(Node, Node)
func (s *SubstituteEs6Syntax) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(s.compiler, root, s)
}

// Java: SubstituteEs6Syntax.visit(NodeTraversal, Node, Node)
func (s *SubstituteEs6Syntax) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		if n.IsArrowFunction() {
			s.MaybeSimplifyArrowFunctionBody(n, n.GetLastChild())
		}
	case rhino.Token_STRING_KEY:
		if n.GetFirstChild().IsName() && n.GetFirstChild().GetString() == n.GetString() {
			n.SetShorthandProperty(true)
			s.objectLiteralShorthandWasAdded = true
		}
	case rhino.Token_SCRIPT:
		if s.objectLiteralShorthandWasAdded {
			NodeUtil_AddFeatureToScript(n, parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, s.compiler)
			s.objectLiteralShorthandWasAdded = false // false for the next script
		}
	default:
	}
}

// MaybeSimplifyArrowFunctionBody: if possible, replace functions of the
// form ()=>{ return x; } with ()=>x
//
// Java: SubstituteEs6Syntax.maybeSimplifyArrowFunctionBody(Node, Node)
func (s *SubstituteEs6Syntax) MaybeSimplifyArrowFunctionBody(arrowFunction *rhino.Node, body *rhino.Node) {
	jutil.CheckArgument(arrowFunction.IsArrowFunction())
	if !body.IsBlock() || !body.HasOneChild() || !body.GetFirstChild().IsReturn() {
		return
	}
	returnValue := body.GetFirstChild().RemoveFirstChild()
	var replacement *rhino.Node
	if returnValue != nil {
		replacement = returnValue
	} else {
		replacement = rhino.IR_Name("undefined")
	}
	body.ReplaceWith(replacement)
	s.compiler.ReportChangeToEnclosingScope(replacement)
}
