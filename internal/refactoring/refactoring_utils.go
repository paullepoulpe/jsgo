// Port of com/google/javascript/refactoring/RefactoringUtils.java @ closure-compiler v20261021

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

package refactoring

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RefactoringUtils holds utility methods for refactoring code.
type RefactoringUtils struct{}

// Java: RefactoringUtils()
func NewRefactoringUtils() *RefactoringUtils {
	return &RefactoringUtils{}
}

// RefactoringUtils_IsInClosurizedFile looks for a goog.require(),
// goog.provide() or goog.module() call in the node's file.
//
// Java: RefactoringUtils.isInClosurizedFile(Node, NodeMetadata)
func RefactoringUtils_IsInClosurizedFile(node rhino.Node, metadata *NodeMetadata) bool {
	script := jscomp.NodeUtil_GetEnclosingScript(node)
	if script.IsNil() {
		return false
	}

	reqMatcher := Matchers_GoogRequire()
	modOrProvideMatcher := Matchers_GoogModuleOrProvide()
	for child := script.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if exprStmt, ok := ast.AsExprStmt(child); ok {
			expr := exprStmt.Expr().Node()
			if expr.IsCall() && (reqMatcher.Matches(expr, metadata) || modOrProvideMatcher.Matches(expr, metadata)) {
				return true
			}
		} else if child.IsVar() && child.GetBooleanProp(rhino.Node_IS_NAMESPACE) {
			return true
		}
	}
	return false
}
