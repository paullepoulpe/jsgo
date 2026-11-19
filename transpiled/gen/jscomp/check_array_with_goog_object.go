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
// Port of com/google/javascript/jscomp/lint/CheckArrayWithGoogObject.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckArrayWithGoogObject
type CheckArrayWithGoogObject struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckArrayWithGoogObject) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckArrayWithGoogObject) IsJavaNil_() bool { return this_ == nil }
func (*CheckArrayWithGoogObject) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckArrayWithGoogObject", "CheckArrayWithGoogObject"
}
func (this_ *CheckArrayWithGoogObject) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckArrayWithGoogObject) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckArrayWithGoogObject", "com.google.javascript.jscomp.lint.CheckArrayWithGoogObject", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckArrayWithGoogObject](v_)
	})
}
func (*CheckArrayWithGoogObject) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckArrayWithGoogObject) Implements_jscomp_CompilerPass()                            {}
func (*CheckArrayWithGoogObject) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckArrayWithGoogObject_GOOG_OBJECT_METHODS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CheckArrayWithGoogObject_GOOG_OBJECT_METHODS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CheckArrayWithGoogObject()
	return static_CheckArrayWithGoogObject_GOOG_OBJECT_METHODS
}
func StaticSet_CheckArrayWithGoogObject_GOOG_OBJECT_METHODS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CheckArrayWithGoogObject()
	static_CheckArrayWithGoogObject_GOOG_OBJECT_METHODS = v
}

var static_CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT *DiagnosticType

func StaticGet_CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT() *DiagnosticType {
	clinit_CheckArrayWithGoogObject()
	return static_CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT
}
func StaticSet_CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT(v *DiagnosticType) {
	clinit_CheckArrayWithGoogObject()
	static_CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT = v
}

var clinitDone_CheckArrayWithGoogObject uint32
var clinitState_CheckArrayWithGoogObject j2gort.ClinitState

func clinit_CheckArrayWithGoogObject() {
	if atomic.LoadUint32(&clinitDone_CheckArrayWithGoogObject) != 0 {
		return
	}
	clinitSlow_CheckArrayWithGoogObject()
}

func clinitSlow_CheckArrayWithGoogObject() {
	if !j2gort.EnterClinit(&clinitState_CheckArrayWithGoogObject, &clinitDone_CheckArrayWithGoogObject) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckArrayWithGoogObject, &clinitDone_CheckArrayWithGoogObject)
	static_CheckArrayWithGoogObject_GOOG_OBJECT_METHODS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("goog.object.forEach", "goog.object.filter", "goog.object.map", "goog.object.some", "goog.object.every", "goog.object.getCount", []any{"goog.object.getAnyKey", "goog.object.getAnyValue", "goog.object.contains", "goog.object.getValues", "goog.object.getKeys", "goog.object.findKey", "goog.object.findValue", "goog.object.isEmpty", "goog.object.clear", "goog.object.remove", "goog.object.equals", "goog.object.clone", "goog.object.transpose"})
	static_CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT = DiagnosticType_Warning("JSC_ARRAY_PASSED_TO_GOOG_OBJECT", "{0} expects an object, not an array. Did you mean to use goog.array?")
}

func (this_ *CheckArrayWithGoogObject) initFields_CheckArrayWithGoogObject() {
	_ = this_
}

// Java constructor: CheckArrayWithGoogObject(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckArrayWithGoogObject(compiler AbstractCompiler) *CheckArrayWithGoogObject {
	clinit_CheckArrayWithGoogObject()
	this_ := &CheckArrayWithGoogObject{}
	InitCheckArrayWithGoogObject(this_, compiler)
	return this_
}

func InitCheckArrayWithGoogObject(this_ *CheckArrayWithGoogObject, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckArrayWithGoogObject()
	this_.compiler = compiler
}

// Java method: CheckArrayWithGoogObject.isGoogObjectIterationOverArray(com.google.javascript.rhino.Node n)
func (this_ *CheckArrayWithGoogObject) IsGoogObjectIterationOverArray(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsCall()) {
		return false
	}
	if !(n.GetFirstChild().IsQualifiedName()) {
		return false
	}
	var name string = n.GetFirstChild().GetQualifiedName()
	_ = name
	if !(j2gort.SetContainsStr(StaticGet_CheckArrayWithGoogObject_GOOG_OBJECT_METHODS(), name)) {
		return false
	}
	var firstArg pkg_rhino.Node = n.GetSecondChild()
	_ = firstArg
	if firstArg == 0 {
		return false
	}
	var type_ pkg_rhino.JSType = firstArg.GetJSType()
	_ = type_
	return ((type_ != nil && !type_.IsJavaNil_()) && this_.priv_CheckArrayWithGoogObject_containsArray_6540(type_))
}

// Java method: CheckArrayWithGoogObject.containsArray(com.google.javascript.rhino.jstype.JSType type)
func (this_ *CheckArrayWithGoogObject) priv_CheckArrayWithGoogObject_containsArray_6540(type_ pkg_rhino.JSType) bool {
	_ = type_
	if type_.IsArrayType() {
		return true
	}
	var templatizedType *pkg_rhino.TemplatizedType = type_.ToMaybeTemplatizedType()
	_ = templatizedType
	if (templatizedType != nil) && templatizedType.GetReferencedType().IsArrayType() {
		return true
	}
	if type_.GetJSTypeBase().IsUnionType() {
		var arrayType pkg_rhino.JSType = this_.compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE())
		_ = arrayType
		var union *pkg_rhino.UnionType = type_.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				if union.GetAlternateAt(i).GetJSTypeBase().IsSubtypeOf_1(arrayType) {
					return true
				}
			}
		}
	}
	return false
}

// Java method: CheckArrayWithGoogObject.visit(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckArrayWithGoogObject) Visit_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = unused
	_ = n
	_ = parent
	if this_.IsGoogObjectIterationOverArray(n) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT(), []string{n.GetFirstChild().GetQualifiedName()}))
	}
}

// Java method: CheckArrayWithGoogObject.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckArrayWithGoogObject) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}
func (*CheckArrayWithGoogObject) AlwaysShouldTraverse_()      {}
func (*CheckArrayWithGoogObject) CanReleaseSyntacticScopes_() {}
