// Port of com/google/javascript/jscomp/AstManipulations.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AstManipulations holds utilities to abstract away certain common AST
// manipulations.
type AstManipulations struct{}

// AstManipulations_FuseExpressions returns a single node equivalent to
// executing <expr1, expr2>.
//
// Requires that expr1 and expr2 are detached (i.e. have no parent nodes).
// If exp2 is EMPTY this just returns exp1.
//
// Java: AstManipulations.fuseExpressions(Node, Node)
func AstManipulations_FuseExpressions(exp1 *rhino.Node, exp2 *rhino.Node) *rhino.Node {
	jutil.CheckArgumentMsg1(exp1.GetParent() == nil, "Expected detached node, got %s", exp1)
	jutil.CheckArgumentMsg1(exp2.GetParent() == nil, "Expected detached node, got %s", exp2)
	if exp2.IsEmpty() {
		return exp1
	}

	comma := rhino.NewNodeWithChild(rhino.Token_COMMA, exp1)
	comma.SrcrefIfMissing(exp2)

	// We can just join the new comma expression with another comma but
	// lets keep all the comma's in a straight line. That way we can use
	// tree comparison.
	if exp2.IsComma() {
		leftMostChild := exp2
		for leftMostChild.IsComma() {
			leftMostChild = leftMostChild.GetFirstChild()
		}
		parent := leftMostChild.GetParent()
		comma.AddChildToBack(leftMostChild.Detach())
		parent.AddChildToFront(comma)
		return exp2
	} else {
		comma.AddChildToBack(exp2)
		return comma
	}
}
