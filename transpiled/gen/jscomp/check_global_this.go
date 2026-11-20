/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckGlobalThis.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CheckGlobalThis
type CheckGlobalThis struct {
	compiler       AbstractCompiler `gson:"compiler"`
	assignLhsChild pkg_rhino.Node   `gson:"assignLhsChild"`
	idHash_        atomic.Int32
}

func (this_ *CheckGlobalThis) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckGlobalThis) IsJavaNil_() bool { return this_ == nil }
func (*CheckGlobalThis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckGlobalThis", "CheckGlobalThis"
}
func (this_ *CheckGlobalThis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckGlobalThis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckGlobalThis", "com.google.javascript.jscomp.CheckGlobalThis", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckGlobalThis](v_)
	})
}
func (*CheckGlobalThis) Implements_jscomp_NodeTraversal_Callback() {}

var static_CheckGlobalThis_GLOBAL_THIS *DiagnosticType

func StaticGet_CheckGlobalThis_GLOBAL_THIS() *DiagnosticType {
	clinit_CheckGlobalThis()
	return static_CheckGlobalThis_GLOBAL_THIS
}
func StaticSet_CheckGlobalThis_GLOBAL_THIS(v *DiagnosticType) {
	clinit_CheckGlobalThis()
	static_CheckGlobalThis_GLOBAL_THIS = v
}

var clinitDone_CheckGlobalThis uint32
var clinitState_CheckGlobalThis j2gort.ClinitState

func clinit_CheckGlobalThis() {
	if atomic.LoadUint32(&clinitDone_CheckGlobalThis) != 0 {
		return
	}
	clinitSlow_CheckGlobalThis()
}

func clinitSlow_CheckGlobalThis() {
	if !j2gort.EnterClinit(&clinitState_CheckGlobalThis, &clinitDone_CheckGlobalThis) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckGlobalThis, &clinitDone_CheckGlobalThis)
	static_CheckGlobalThis_GLOBAL_THIS = DiagnosticType_Warning("JSC_USED_GLOBAL_THIS", "dangerous use of the global 'this' object")
}

func (this_ *CheckGlobalThis) initFields_CheckGlobalThis() {
	_ = this_

}

// Java constructor: CheckGlobalThis(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckGlobalThis(compiler AbstractCompiler) *CheckGlobalThis {
	clinit_CheckGlobalThis()
	this_ := &CheckGlobalThis{}
	InitCheckGlobalThis(this_, compiler)
	return this_
}

func InitCheckGlobalThis(this_ *CheckGlobalThis, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckGlobalThis()
	this_.compiler = compiler
}

// Java method: CheckGlobalThis.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckGlobalThis) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsFunction() {
		if n.IsArrowFunction() {
			return true
		}
		var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = jsDoc
		if (jsDoc != nil) && (((jsDoc.IsConstructor() || jsDoc.IsInterface()) || jsDoc.HasThisType()) || jsDoc.IsOverride()) {
			return false
		}
		if jsDoc != nil {
			var functionType pkg_rhino.JSTypeExpression = jsDoc.GetType()
			_ = functionType
			if !(functionType).IsJavaNil_() {
				var functionNode pkg_rhino.Node = functionType.GetRoot()
				_ = functionNode
				if (functionNode != 0) && functionNode.IsFunction() {
					var thisNode pkg_rhino.Node = functionNode.GetFirstChild()
					_ = thisNode
					if (thisNode != 0) && thisNode.IsThis() {
						return false
					}
				}
			}
		}
		var pType pkg_rhino.Token = parent.GetToken()
		_ = pType
		if !(((((pType == pkg_rhino.StaticGet_Token_BLOCK()) || (pType == pkg_rhino.StaticGet_Token_SCRIPT())) || (pType == pkg_rhino.StaticGet_Token_NAME())) || (pType == pkg_rhino.StaticGet_Token_ASSIGN())) || (pType == pkg_rhino.StaticGet_Token_STRING_KEY())) {
			return false
		}
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		if NodeUtil_MayBeObjectLitKey(parent) {
			var maybeLends *pkg_rhino.JSDocInfo = grandparent.GetJSDocInfo()
			_ = maybeLends
			if ((maybeLends != nil) && maybeLends.HasLendsName()) && javalang.StringEndsWith(j2gort.CheckNotNullString(maybeLends.GetLendsName().GetRoot().GetString()), ".prototype") {
				return false
			}
		}
	}
	if NodeUtil_IsClassStaticBlock(n) {
		return false
	}
	if (n.IsMemberFieldDef() || n.IsComputedFieldDef()) && n.GetParent().IsClassMembers() {
		return false
	}
	if (parent != 0) && parent.IsAssign() {
		var lhs pkg_rhino.Node = parent.GetFirstChild()
		_ = lhs
		if n == lhs {
			if this_.assignLhsChild == 0 {
				this_.assignLhsChild = lhs
			}
		} else {
			if NodeUtil_IsNormalGet(lhs) {
				if lhs.IsGetProp() && (j2gort.CheckNotNullString(lhs.GetString()) == "prototype") {
					return false
				}
				var llhs pkg_rhino.Node = lhs.GetFirstChild()
				_ = llhs
				if llhs.IsGetProp() && (j2gort.CheckNotNullString(llhs.GetString()) == "prototype") {
					return false
				}
			}
		}
	}
	return true
}

// Java method: CheckGlobalThis.visit(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckGlobalThis) Visit_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = unused
	_ = n
	_ = parent
	if n.IsThis() && this_.priv_CheckGlobalThis_shouldReportThis_1056(n) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CheckGlobalThis_GLOBAL_THIS(), []string{}))
	}
	if n == this_.assignLhsChild {
		this_.assignLhsChild = pkg_rhino.Node(0)
	}
}

// Java method: CheckGlobalThis.shouldReportThis(com.google.javascript.rhino.Node n)
func (this_ *CheckGlobalThis) priv_CheckGlobalThis_shouldReportThis_1056(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if this_.assignLhsChild != 0 {
		return true
	}
	return ((parent != 0) && NodeUtil_IsNormalGet(parent))
}
func (*CheckGlobalThis) CanReleaseSyntacticScopes_() {}
