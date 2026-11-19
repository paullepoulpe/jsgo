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
// Port of com/google/javascript/jscomp/RewriteJsonToModule.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.RewriteJsonToModule
type RewriteJsonToModule struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	packageJsonMainEntries j2gort.OrderedMap[string, string] `gson:"packageJsonMainEntries"`
	compiler               AbstractCompiler                  `gson:"compiler"`
}

func (this_ *RewriteJsonToModule) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteJsonToModule) IsJavaNil_() bool { return this_ == nil }
func (*RewriteJsonToModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteJsonToModule", "RewriteJsonToModule"
}
func (this_ *RewriteJsonToModule) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RewriteJsonToModule) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteJsonToModule", "com.google.javascript.jscomp.RewriteJsonToModule", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteJsonToModule](v_)
	})
}
func (*RewriteJsonToModule) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*RewriteJsonToModule) Implements_jscomp_CompilerPass()                            {}
func (*RewriteJsonToModule) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN *DiagnosticType

func StaticGet_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN() *DiagnosticType {
	clinit_RewriteJsonToModule()
	return static_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN
}
func StaticSet_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN(v *DiagnosticType) {
	clinit_RewriteJsonToModule()
	static_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN = v
}

var clinitDone_RewriteJsonToModule uint32
var clinitState_RewriteJsonToModule j2gort.ClinitState

func clinit_RewriteJsonToModule() {
	if atomic.LoadUint32(&clinitDone_RewriteJsonToModule) != 0 {
		return
	}
	clinitSlow_RewriteJsonToModule()
}

func clinitSlow_RewriteJsonToModule() {
	if !j2gort.EnterClinit(&clinitState_RewriteJsonToModule, &clinitDone_RewriteJsonToModule) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RewriteJsonToModule, &clinitDone_RewriteJsonToModule)
	static_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN = DiagnosticType_Error("JSC_JSON_UNEXPECTED_TOKEN", "Unexpected JSON token")
}

func (this_ *RewriteJsonToModule) initFields_RewriteJsonToModule() {
	_ = this_
}

// Java constructor: RewriteJsonToModule(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteJsonToModule(compiler AbstractCompiler) *RewriteJsonToModule {
	clinit_RewriteJsonToModule()
	this_ := &RewriteJsonToModule{}
	InitRewriteJsonToModule(this_, compiler)
	return this_
}

func InitRewriteJsonToModule(this_ *RewriteJsonToModule, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RewriteJsonToModule()
	this_.compiler = compiler
	this_.packageJsonMainEntries = j2gort.NewOrderedMap[string, string]()
}

// Java method: RewriteJsonToModule.getPackageJsonMainEntries()
func (this_ *RewriteJsonToModule) GetPackageJsonMainEntries() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.packageJsonMainEntries.ToImmutableMap()
}

// Java method: RewriteJsonToModule.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteJsonToModule) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(root.IsScript()) {
		j2gort.CheckState(false)
	}
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: RewriteJsonToModule.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteJsonToModule) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		if !(n.HasOneChild()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN(), []string{}))
		} else {
			this_.priv_RewriteJsonToModule_visitScript_5004(t, n)
		}
		return
	case int32(47) /* Token.OBJECTLIT */, int32(46) /* Token.ARRAYLIT */, int32(33) /* Token.NUMBER */, int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */, int32(36) /* Token.NULL */, int32(35) /* Token.STRINGLIT */ :
	case int32(100) /* Token.STRING_KEY */ :
		if !(n.IsQuotedStringKey()) || !(n.HasOneChild()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN(), []string{}))
		}
	case int32(93) /* Token.EXPR_RESULT */ :
		if !(parent.IsScript()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN(), []string{}))
		}
	default:
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN(), []string{}))
	}
	if n.GetLineno() == int32(1) {
		n.SetLinenoCharno(int32(1), (n.GetCharno() - int32(1)))
		t.ReportCodeChange_0()
	}
}

// Java method: RewriteJsonToModule.visitScript(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteJsonToModule) priv_RewriteJsonToModule_visitScript_5004(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !(n.HasOneChild()) || !(n.GetFirstChild().IsExprResult()) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_RewriteJsonToModule_JSON_UNEXPECTED_TOKEN(), []string{}))
		return
	}
	var jsdoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = jsdoc
	jsdoc.RecordFileOverview("Suppresses undefined var goog error")
	jsdoc.RecordSuppression("undefinedVars")
	n.SetJSDocInfo(jsdoc.Build_0())
	var jsonObject pkg_rhino.Node = n.GetFirstFirstChild().Detach()
	_ = jsonObject
	n.RemoveFirstChild()
	var moduleName string = t.GetInput().GetPath().ToModuleName()
	_ = moduleName
	n.AddChildToFront(pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(moduleName).Srcref(jsonObject), jsonObject).Srcref(jsonObject))
	n.AddChildToFront(pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "provide"), []pkg_rhino.Node{pkg_rhino.IR_String(moduleName)})).SrcrefTreeIfMissing(n))
	var inputPath string = t.GetInput().GetSourceFile().GetName()
	_ = inputPath
	if javalang.StringEndsWith(j2gort.CheckNotNullString(inputPath), "/package.json") && jsonObject.IsObjectLit() {
		var possibleMainEntries j2gort.Ext_java_util_List = this_.compiler.GetOptions().GetPackageJsonEntryNames()
		_ = possibleMainEntries
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(possibleMainEntries))
			for cur_1.HasNext() {
				entryName := cur_1.NextStr()
				_ = entryName
				var entry pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(jsonObject, entryName)
				_ = entry
				if (entry != 0) && (entry.IsStringLit() || entry.IsObjectLit()) {
					var dirName string = j2gort.JStrSub(inputPath, int32(0), (j2gort.JStrLen(inputPath) - j2gort.JStrLen("package.json")))
					_ = dirName
					if entry.IsStringLit() {
						this_.packageJsonMainEntries.Put(inputPath, j2gort.StrValString(dirName)+j2gort.StrValString(entry.GetString()))
						break
					} else if entry.IsObjectLit() {
						if !(j2gort.CheckNotNullString(entryName) == "browser") {
							j2gort.CheckStateStr(false, entryName)
						}
						this_.priv_RewriteJsonToModule_processBrowserFieldAdvancedUsage_5005(dirName, entry)
					}
				}
			}
		}
	}
	t.ReportCodeChange_0()
}

// Java method: RewriteJsonToModule.processBrowserFieldAdvancedUsage(java.lang.String dirName, com.google.javascript.rhino.Node entry)
func (this_ *RewriteJsonToModule) priv_RewriteJsonToModule_processBrowserFieldAdvancedUsage_5005(dirName string, entry pkg_rhino.Node) {
	_ = dirName
	_ = entry
	{
		var child pkg_rhino.Node = entry.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			var value pkg_rhino.Node = child.GetFirstChild()
			_ = value
			if !(child.IsStringKey() && (value.IsStringLit() || value.IsFalse())) {
				j2gort.CheckState(false)
			}
			var path string = child.GetString()
			_ = path
			if javalang.StringStartsWith(j2gort.CheckNotNullString(path), "./") {
				path = j2gort.JStrSubFrom(path, j2gort.JStrLen("./"))
			}
			var replacement string
			_ = replacement
			if value.IsStringLit() {
				replacement = (j2gort.StrValString(dirName) + j2gort.StrValString(value.GetString()))
			} else {
				replacement = "$jscomp$browser$skiplisted"
			}
			this_.packageJsonMainEntries.Put(j2gort.StrValString(dirName)+j2gort.StrValString(path), replacement)
		}
	}
}
func (*RewriteJsonToModule) AlwaysShouldTraverse_()      {}
func (*RewriteJsonToModule) CanReleaseSyntacticScopes_() {}
