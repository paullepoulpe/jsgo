/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckGoogModuleTypeScriptName.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.lint.CheckGoogModuleTypeScriptName
type CheckGoogModuleTypeScriptName struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
	finished bool `gson:"finished"`
}

func (this_ *CheckGoogModuleTypeScriptName) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckGoogModuleTypeScriptName) IsJavaNil_() bool { return this_ == nil }
func (*CheckGoogModuleTypeScriptName) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckGoogModuleTypeScriptName", "CheckGoogModuleTypeScriptName"
}
func (this_ *CheckGoogModuleTypeScriptName) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckGoogModuleTypeScriptName) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckGoogModuleTypeScriptName", "com.google.javascript.jscomp.lint.CheckGoogModuleTypeScriptName", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckGoogModuleTypeScriptName](v_)
	})
}
func (*CheckGoogModuleTypeScriptName) Implements_jscomp_NodeTraversal_Callback() {}
func (*CheckGoogModuleTypeScriptName) Implements_jscomp_CompilerPass()           {}

var static_CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE *DiagnosticType

func StaticGet_CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE() *DiagnosticType {
	clinit_CheckGoogModuleTypeScriptName()
	return static_CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE
}
func StaticSet_CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE(v *DiagnosticType) {
	clinit_CheckGoogModuleTypeScriptName()
	static_CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE = v
}

var static_CheckGoogModuleTypeScriptName_allowedDirectories j2gort.Ext_java_util_Set

func StaticGet_CheckGoogModuleTypeScriptName_AllowedDirectories() j2gort.Ext_java_util_Set {
	clinit_CheckGoogModuleTypeScriptName()
	return static_CheckGoogModuleTypeScriptName_allowedDirectories
}
func StaticSet_CheckGoogModuleTypeScriptName_AllowedDirectories(v j2gort.Ext_java_util_Set) {
	clinit_CheckGoogModuleTypeScriptName()
	static_CheckGoogModuleTypeScriptName_allowedDirectories = v
}

var clinitDone_CheckGoogModuleTypeScriptName uint32
var clinitState_CheckGoogModuleTypeScriptName j2gort.ClinitState

func clinit_CheckGoogModuleTypeScriptName() {
	if atomic.LoadUint32(&clinitDone_CheckGoogModuleTypeScriptName) != 0 {
		return
	}
	clinitSlow_CheckGoogModuleTypeScriptName()
}

func clinitSlow_CheckGoogModuleTypeScriptName() {
	if !j2gort.EnterClinit(&clinitState_CheckGoogModuleTypeScriptName, &clinitDone_CheckGoogModuleTypeScriptName) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckGoogModuleTypeScriptName, &clinitDone_CheckGoogModuleTypeScriptName)
	static_CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE = DiagnosticType_Disabled("JSC_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE", "goog.module namespace does not match the future TypeScript namespace, which is generated from the file path. The correct namespace is: \"{0}\"")
	static_CheckGoogModuleTypeScriptName_allowedDirectories = j2gort.NewExt_java_util_LinkedHashSet()
	j2gort.SetAddStr(StaticGet_CheckGoogModuleTypeScriptName_AllowedDirectories(), "google3/gws/")
	j2gort.SetAddStr(StaticGet_CheckGoogModuleTypeScriptName_AllowedDirectories(), "google3/java/com/google/gws/")
	j2gort.SetAddStr(StaticGet_CheckGoogModuleTypeScriptName_AllowedDirectories(), "google3/javascript/search/")
}

func (this_ *CheckGoogModuleTypeScriptName) initFields_CheckGoogModuleTypeScriptName() {
	_ = this_

}

// Java constructor: CheckGoogModuleTypeScriptName(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckGoogModuleTypeScriptName(compiler AbstractCompiler) *CheckGoogModuleTypeScriptName {
	clinit_CheckGoogModuleTypeScriptName()
	this_ := &CheckGoogModuleTypeScriptName{}
	InitCheckGoogModuleTypeScriptName(this_, compiler)
	return this_
}

func InitCheckGoogModuleTypeScriptName(this_ *CheckGoogModuleTypeScriptName, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckGoogModuleTypeScriptName()
	this_.compiler = compiler
}

// Java method: CheckGoogModuleTypeScriptName.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckGoogModuleTypeScriptName) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckGoogModuleTypeScriptName.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckGoogModuleTypeScriptName) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(this_.finished) && ((((parent == 0) || parent.IsRoot()) || parent.IsScript()) || parent.IsModuleBody()))
}

// Java method: CheckGoogModuleTypeScriptName.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckGoogModuleTypeScriptName) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsGoogModuleCall(n) {
		this_.priv_CheckGoogModuleTypeScriptName_checkGoogModuleNamespace_6535(t, n)
		this_.finished = true
	}
}

// Java method: CheckGoogModuleTypeScriptName.checkGoogModuleNamespace(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckGoogModuleTypeScriptName) priv_CheckGoogModuleTypeScriptName_checkGoogModuleNamespace_6535(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var originalNamespace string = n.GetFirstChild().GetSecondChild().GetString()
	_ = originalNamespace
	var sourceName string = t.GetSourceName()
	_ = sourceName
	if (originalNamespace == j2gort.NullString) || (sourceName == j2gort.NullString) {
		return
	}
	var google3Index int32 = j2gort.JStrIndexOf(sourceName, "google3")
	_ = google3Index
	if google3Index == int32(-1) {
		sourceName = ("google3/" + j2gort.StrValString(sourceName))
	} else if google3Index != int32(0) {
		sourceName = j2gort.JStrSubFrom(sourceName, google3Index)
	}
	var replacementNamespace string = j2gort.JStrSub(j2gort.JStrReplaceChar(sourceName, uint16(47), uint16(46)), int32(0), (j2gort.JStrLen(sourceName) - j2gort.JStrLen(".js")))
	_ = replacementNamespace
	if !(j2gort.CheckNotNullString(originalNamespace) == replacementNamespace) {
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_CheckGoogModuleTypeScriptName_AllowedDirectories()))
			for cur_1.HasNext() {
				allowedDirectory := cur_1.NextStr()
				_ = allowedDirectory
				if javalang.StringStartsWith(j2gort.CheckNotNullString(sourceName), j2gort.CheckNotNullString(allowedDirectory)) {
					t.Report_3(n, StaticGet_CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE(), []string{replacementNamespace})
					return
				}
			}
		}
	}
}
func (*CheckGoogModuleTypeScriptName) CanReleaseSyntacticScopes_() {}
