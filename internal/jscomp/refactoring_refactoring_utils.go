// Port of com/google/javascript/refactoring/RefactoringUtils.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/rhino"

// RefactoringUtils holds utility methods for refactoring Java code.
//
// Java: public class RefactoringUtils
type RefactoringUtils struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: RefactoringUtils()
func NewRefactoringUtils() *RefactoringUtils {
	return &RefactoringUtils{}
}

// RefactoringUtils_IsInClosurizedFile looks for a goog.require(),
// goog.provide() or goog.module() call in the node's file.
//
// Java: RefactoringUtils.isInClosurizedFile(Node, NodeMetadata)
func RefactoringUtils_IsInClosurizedFile(node rhino.Node, metadata *NodeMetadata) bool {
	script := NodeUtil_GetEnclosingScript(node)

	if script.IsNil() {
		return false
	}

	child := script.GetFirstChild()
	for !child.IsNil() {
		if NodeUtil_IsExprCall(child) {
			firstChild := child.GetFirstChild()
			if Matchers_GoogRequire().Matches(firstChild, metadata) ||
				Matchers_GoogModuleOrProvide().Matches(firstChild, metadata) {
				return true
			}
			// goog.require or goog.module.
		} else if child.IsVar() && child.GetBooleanProp(rhino.Node_IS_NAMESPACE) {
			return true
		}
		child = child.GetNext()
	}
	return false
}
