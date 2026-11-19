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
// Port of com/google/javascript/refactoring/RefactoringUtils.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.refactoring.RefactoringUtils
type RefactoringUtils struct {
	idHash_ atomic.Int32
}

func (this_ *RefactoringUtils) BindSelf_(s any) {
	_ = s
}
func (this_ *RefactoringUtils) IsJavaNil_() bool { return this_ == nil }
func (*RefactoringUtils) JavaClassInfo_() (string, string) { return "com.google.javascript.refactoring.RefactoringUtils", "RefactoringUtils" }
func (this_ *RefactoringUtils) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RefactoringUtils) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.RefactoringUtils", "com.google.javascript.refactoring.RefactoringUtils", func(v_ any) bool {
		return j2gort.InstanceOf[*RefactoringUtils](v_)
	})
}

func (this_ *RefactoringUtils) initFields_RefactoringUtils() {
	_ = this_
}

// Java constructor: RefactoringUtils()
func NewRefactoringUtils() *RefactoringUtils {
	this_ := &RefactoringUtils{}
	InitRefactoringUtils(this_)
	return this_
}

func InitRefactoringUtils(this_ *RefactoringUtils) {
	this_.initFields_RefactoringUtils()
}

// Java static method: RefactoringUtils.isInClosurizedFile(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func RefactoringUtils_IsInClosurizedFile(node pkg_rhino.Node, metadata *NodeMetadata) bool {
	_ = node
	_ = metadata
	var script pkg_rhino.Node = NodeUtil_GetEnclosingScript(node)
	_ = script
	if ((script == 0)) {
		return false
	}
	var child pkg_rhino.Node = script.GetFirstChild()
	_ = child
	for ((child != 0)) {
		if (NodeUtil_IsExprCall(child)) {
			var firstChild pkg_rhino.Node = child.GetFirstChild()
			_ = firstChild
			if ((Matchers_GoogRequire().Matches(firstChild, metadata) || Matchers_GoogModuleOrProvide().Matches(firstChild, metadata))) {
				return true
			}
		} else if ((child.IsVar() && child.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE()))) {
			return true
		}
		child = child.GetNext()
	}
	return false
}

