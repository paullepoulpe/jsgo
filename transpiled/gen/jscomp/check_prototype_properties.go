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
// Port of com/google/javascript/jscomp/lint/CheckPrototypeProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckPrototypeProperties
type CheckPrototypeProperties struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *CheckPrototypeProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckPrototypeProperties) IsJavaNil_() bool { return this_ == nil }
func (*CheckPrototypeProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckPrototypeProperties", "CheckPrototypeProperties"
}
func (this_ *CheckPrototypeProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckPrototypeProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckPrototypeProperties", "com.google.javascript.jscomp.lint.CheckPrototypeProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckPrototypeProperties](v_)
	})
}
func (*CheckPrototypeProperties) Implements_jscomp_CompilerPass()           {}
func (*CheckPrototypeProperties) Implements_jscomp_NodeTraversal_Callback() {}

var static_CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER *DiagnosticType

func StaticGet_CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER() *DiagnosticType {
	clinit_CheckPrototypeProperties()
	return static_CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER
}
func StaticSet_CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER(v *DiagnosticType) {
	clinit_CheckPrototypeProperties()
	static_CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER = v
}

var clinitDone_CheckPrototypeProperties uint32
var clinitState_CheckPrototypeProperties j2gort.ClinitState

func clinit_CheckPrototypeProperties() {
	if atomic.LoadUint32(&clinitDone_CheckPrototypeProperties) != 0 {
		return
	}
	clinitSlow_CheckPrototypeProperties()
}

func clinitSlow_CheckPrototypeProperties() {
	if !j2gort.EnterClinit(&clinitState_CheckPrototypeProperties, &clinitDone_CheckPrototypeProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckPrototypeProperties, &clinitDone_CheckPrototypeProperties)
	static_CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER = DiagnosticType_Disabled("JSC_ILLEGAL_PROTOTYPE_MEMBER", "Prototype property {0} should be a primitive, not an Array or Object.")
}

func (this_ *CheckPrototypeProperties) initFields_CheckPrototypeProperties() {
	_ = this_
}

// Java constructor: CheckPrototypeProperties(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckPrototypeProperties(compiler AbstractCompiler) *CheckPrototypeProperties {
	clinit_CheckPrototypeProperties()
	this_ := &CheckPrototypeProperties{}
	InitCheckPrototypeProperties(this_, compiler)
	return this_
}

func InitCheckPrototypeProperties(this_ *CheckPrototypeProperties, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckPrototypeProperties()
	this_.compiler = compiler
}

// Java method: CheckPrototypeProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckPrototypeProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckPrototypeProperties.visit(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckPrototypeProperties) Visit_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = unused
	_ = n
	_ = parent
	if NodeUtil_IsPrototypePropertyDeclaration(n) {
		var assign pkg_rhino.Node = n.GetFirstChild()
		_ = assign
		var rhs pkg_rhino.Node = assign.GetLastChild()
		_ = rhs
		if rhs.IsArrayLit() || rhs.IsObjectLit() {
			var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(rhs)
			_ = jsDoc
			if (jsDoc != nil) && jsDoc.HasEnumParameterType() {
				return
			}
			var propName string = assign.GetFirstChild().GetLastChild().GetString()
			_ = propName
			this_.compiler.Report_1(JSError_Make_3(assign, StaticGet_CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER(), []string{propName}))
		}
	}
}

// Java method: CheckPrototypeProperties.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckPrototypeProperties) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return true
}
func (*CheckPrototypeProperties) AlwaysShouldTraverse_()      {}
func (*CheckPrototypeProperties) CanReleaseSyntacticScopes_() {}
