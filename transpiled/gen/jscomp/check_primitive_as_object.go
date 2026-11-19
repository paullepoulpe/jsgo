/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckPrimitiveAsObject.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckPrimitiveAsObject
type CheckPrimitiveAsObject struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckPrimitiveAsObject) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckPrimitiveAsObject) IsJavaNil_() bool { return this_ == nil }
func (*CheckPrimitiveAsObject) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckPrimitiveAsObject", "CheckPrimitiveAsObject"
}
func (this_ *CheckPrimitiveAsObject) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckPrimitiveAsObject) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckPrimitiveAsObject", "com.google.javascript.jscomp.lint.CheckPrimitiveAsObject", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckPrimitiveAsObject](v_)
	})
}
func (*CheckPrimitiveAsObject) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckPrimitiveAsObject) Implements_jscomp_CompilerPass()                            {}
func (*CheckPrimitiveAsObject) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT *DiagnosticType

func StaticGet_CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT() *DiagnosticType {
	clinit_CheckPrimitiveAsObject()
	return static_CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT
}
func StaticSet_CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT(v *DiagnosticType) {
	clinit_CheckPrimitiveAsObject()
	static_CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT = v
}

var static_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION *DiagnosticType

func StaticGet_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION() *DiagnosticType {
	clinit_CheckPrimitiveAsObject()
	return static_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION
}
func StaticSet_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION(v *DiagnosticType) {
	clinit_CheckPrimitiveAsObject()
	static_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION = v
}

var static_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CheckPrimitiveAsObject()
	return static_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS
}
func StaticSet_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CheckPrimitiveAsObject()
	static_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS = v
}

var clinitDone_CheckPrimitiveAsObject uint32
var clinitState_CheckPrimitiveAsObject j2gort.ClinitState

func clinit_CheckPrimitiveAsObject() {
	if atomic.LoadUint32(&clinitDone_CheckPrimitiveAsObject) != 0 {
		return
	}
	clinitSlow_CheckPrimitiveAsObject()
}

func clinitSlow_CheckPrimitiveAsObject() {
	if !j2gort.EnterClinit(&clinitState_CheckPrimitiveAsObject, &clinitDone_CheckPrimitiveAsObject) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckPrimitiveAsObject, &clinitDone_CheckPrimitiveAsObject)
	static_CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT = DiagnosticType_Warning("JSC_PRIMITIVE_OBJECT", "Explicit creation of a {0} object.")
	static_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION = DiagnosticType_Warning("JSC_PRIMITIVE_OBJECT_DECLARATION", "Declaration of {0} object instead of primitive type.")
	static_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3("Boolean", "Number", "String")
}

func (this_ *CheckPrimitiveAsObject) initFields_CheckPrimitiveAsObject() {
	_ = this_
}

// Java constructor: CheckPrimitiveAsObject(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckPrimitiveAsObject(compiler AbstractCompiler) *CheckPrimitiveAsObject {
	clinit_CheckPrimitiveAsObject()
	this_ := &CheckPrimitiveAsObject{}
	InitCheckPrimitiveAsObject(this_, compiler)
	return this_
}

func InitCheckPrimitiveAsObject(this_ *CheckPrimitiveAsObject, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckPrimitiveAsObject()
	this_.compiler = compiler
}

// Java method: CheckPrimitiveAsObject.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckPrimitiveAsObject) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckPrimitiveAsObject.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckPrimitiveAsObject) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.priv_CheckPrimitiveAsObject_checkForPrimitiveObjectConstructor_6535(t, n)
	this_.priv_CheckPrimitiveAsObject_checkForPrimitiveObjectDeclaration_6533(t, n)
}

// Java method: CheckPrimitiveAsObject.checkForPrimitiveObjectDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckPrimitiveAsObject) priv_CheckPrimitiveAsObject_checkForPrimitiveObjectDeclaration_6533(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var jsDocInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsDocInfo
	if jsDocInfo != nil {
		{
			var i int32 = int32(0)
			_ = i
			var count int32 = jsDocInfo.GetTypeNodeCount()
			_ = count
			for ; i < count; i = i + 1 {
				this_.priv_CheckPrimitiveAsObject_checkTypeNodeForPrimitiveObjectDeclaration_6534(t, jsDocInfo.GetTypeNodeAt(i))
			}
		}
	}
}

// Java method: CheckPrimitiveAsObject.checkTypeNodeForPrimitiveObjectDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node typeRoot)
func (this_ *CheckPrimitiveAsObject) priv_CheckPrimitiveAsObject_checkTypeNodeForPrimitiveObjectDeclaration_6534(t *NodeTraversal, typeRoot pkg_rhino.Node) {
	_ = t
	_ = typeRoot
	NodeUtil_VisitPreOrder_2(typeRoot, NewLambda_NodeUtil_Visitor(func(node pkg_rhino.Node) {
		_ = node
		if node.IsStringLit() {
			var typeName string = node.GetString()
			_ = typeName
			if j2gort.SetContainsStr(StaticGet_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS(), typeName) {
				t.Report_3(node, StaticGet_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION(), []string{typeName})
			}
		}
	}))
}

// Java method: CheckPrimitiveAsObject.checkForPrimitiveObjectConstructor(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckPrimitiveAsObject) priv_CheckPrimitiveAsObject_checkForPrimitiveObjectConstructor_6535(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if n.IsNew() {
		var constructorFunction pkg_rhino.Node = n.GetFirstChild()
		_ = constructorFunction
		if constructorFunction.IsName() {
			var constructorName string = constructorFunction.GetString()
			_ = constructorName
			if j2gort.SetContainsStr(StaticGet_CheckPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS(), constructorName) {
				t.Report_3(n, StaticGet_CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT(), []string{constructorName})
			}
		}
	}
}
func (*CheckPrimitiveAsObject) AlwaysShouldTraverse_()      {}
func (*CheckPrimitiveAsObject) CanReleaseSyntacticScopes_() {}
